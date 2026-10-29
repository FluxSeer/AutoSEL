package policy

import (
	"fmt"
	"path"
	"sort"
	"strings"

	"github.com/FluxSeer/AutoSEL/internal/spec"
)

type Status string

const (
	StatusReady          Status = "ready"
	StatusReviewRequired Status = "review-required"
)

type Severity string

const (
	SeverityHigh     Severity = "high"
	SeverityCritical Severity = "critical"
)

type Finding struct {
	Code     string   `json:"code"`
	Severity Severity `json:"severity"`
	Subject  string   `json:"subject"`
	Action   string   `json:"action"`
	Reason   string   `json:"reason"`
}

type Plan struct {
	ContainerID    string    `json:"containerId"`
	ModuleName     string    `json:"moduleName"`
	PlannerVersion string    `json:"plannerVersion"`
	InputDigest    string    `json:"inputDigest"`
	PolicyDigest   string    `json:"policyDigest"`
	Status         Status    `json:"status"`
	Findings       []Finding `json:"findings,omitempty"`
}

// Planner is deterministic and does not access Docker, SELinux, the
// filesystem, or the network. It produces reviewable findings rather than
// silently mutating a workload.
type Planner struct{}

func (Planner) Plan(input spec.ContainerSpec) (Plan, error) {
	if strings.TrimSpace(input.ID) == "" {
		return Plan{}, fmt.Errorf("container id is required")
	}
	input = canonicalSpec(input)
	inputDigest, err := digestJSON(input)
	if err != nil {
		return Plan{}, fmt.Errorf("digest normalized container spec: %w", err)
	}

	plan := Plan{
		ContainerID:    input.ID,
		ModuleName:     "autosel-container-" + input.ID,
		PlannerVersion: PlannerVersion,
		InputDigest:    inputDigest,
		Status:         StatusReady,
	}

	if input.Privileged {
		plan.Findings = append(plan.Findings, Finding{
			Code:     "privileged-container",
			Severity: SeverityCritical,
			Subject:  "container",
			Action:   "review",
			Reason:   "privileged containers expand capabilities and device access",
		})
	}

	if !input.Capabilities.Complete {
		plan.Findings = append(plan.Findings, Finding{
			Code:     "capabilities-incomplete",
			Severity: SeverityHigh,
			Subject:  "capabilities",
			Action:   "review",
			Reason:   "the source does not expose the complete runtime-resolved capability set",
		})
	}

	capabilities := append(append([]string(nil), input.Capabilities.Bounding...), input.Capabilities.Added...)
	seenCapabilities := make(map[string]bool)
	for _, capability := range capabilities {
		capability = strings.ToUpper(strings.TrimSpace(capability))
		if seenCapabilities[capability] {
			continue
		}
		seenCapabilities[capability] = true
		if reason, ok := highRiskCapabilities[capability]; ok {
			plan.Findings = append(plan.Findings, Finding{
				Code:     "high-risk-capability",
				Severity: SeverityHigh,
				Subject:  capability,
				Action:   "review",
				Reason:   reason,
			})
		}
	}

	for _, mount := range input.Mounts {
		destination := path.Clean(strings.TrimSpace(mount.Destination))
		if reason, ok := sensitiveMounts[destination]; ok {
			plan.Findings = append(plan.Findings, Finding{
				Code:     "sensitive-mount",
				Severity: SeverityHigh,
				Subject:  destination,
				Action:   "review",
				Reason:   reason,
			})
		}
	}

	for _, device := range input.Devices {
		devicePath := strings.TrimSpace(device.PathOnHost)
		for _, pattern := range highRiskDevicePatterns {
			matched, err := path.Match(pattern, devicePath)
			if err == nil && matched {
				plan.Findings = append(plan.Findings, Finding{
					Code:     "high-risk-device",
					Severity: SeverityCritical,
					Subject:  devicePath,
					Action:   "review",
					Reason:   "device access can expose host hardware or kernel interfaces",
				})
				break
			}
		}
	}

	sort.Slice(plan.Findings, func(i, j int) bool {
		if plan.Findings[i].Code != plan.Findings[j].Code {
			return plan.Findings[i].Code < plan.Findings[j].Code
		}
		return plan.Findings[i].Subject < plan.Findings[j].Subject
	})
	if len(plan.Findings) > 0 {
		plan.Status = StatusReviewRequired
	}
	policyDigest, err := digestJSON(struct {
		ContainerID    string    `json:"containerId"`
		ModuleName     string    `json:"moduleName"`
		PlannerVersion string    `json:"plannerVersion"`
		Status         Status    `json:"status"`
		Findings       []Finding `json:"findings,omitempty"`
	}{
		ContainerID:    plan.ContainerID,
		ModuleName:     plan.ModuleName,
		PlannerVersion: plan.PlannerVersion,
		Status:         plan.Status,
		Findings:       plan.Findings,
	})
	if err != nil {
		return Plan{}, fmt.Errorf("digest policy plan: %w", err)
	}
	plan.PolicyDigest = policyDigest
	return plan, nil
}

var highRiskCapabilities = map[string]string{
	"CAP_BPF":                "BPF can expose powerful kernel instrumentation and attack surface",
	"CAP_CHECKPOINT_RESTORE": "checkpoint/restore can affect process state across isolation boundaries",
	"CAP_NET_ADMIN":          "network administration can alter host-facing network behavior",
	"CAP_PERFMON":            "performance monitoring can expose sensitive kernel information",
	"CAP_SYS_ADMIN":          "SYS_ADMIN is a broad capability with many privileged operations",
	"CAP_SYS_BOOT":           "system boot operations can affect host availability",
	"CAP_SYS_MODULE":         "kernel module operations can affect the host kernel",
	"CAP_SYS_PTRACE":         "process tracing can expose or alter other process state",
	"CAP_SYS_RAWIO":          "raw I/O can access hardware or kernel interfaces",
	"CAP_SYS_TIME":           "system time changes affect host-wide time behavior",
}

var sensitiveMounts = map[string]string{
	"/boot":                "boot files are host-sensitive",
	"/dev":                 "device mounts expose host device nodes",
	"/etc/group":           "group database access can alter identity controls",
	"/etc/passwd":          "passwd database access can alter identity controls",
	"/etc/shadow":          "shadow database access exposes password hashes",
	"/proc":                "proc access can expose host or kernel process state",
	"/proc/kcore":          "kernel core access is host-sensitive",
	"/root":                "root home access may expose host credentials",
	"/sys":                 "sysfs access can expose or alter kernel state",
	"/sys/kernel/security": "securityfs access can affect host security controls",
	"/var/lib/docker":      "Docker data access can control or expose other workloads",
}

var highRiskDevicePatterns = []string{
	"/dev/bus/*",
	"/dev/dri/*",
	"/dev/kvm",
	"/dev/mem",
	"/dev/net/tun",
	"/dev/nvme*",
	"/dev/sd*",
	"/dev/vfio/*",
	"/dev/vhost-*",
	"/dev/port",
}
