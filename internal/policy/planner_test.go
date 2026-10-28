package policy

import (
	"testing"

	"github.com/FluxSeer/AutoSEL/internal/spec"
)

func TestPlannerProducesReviewableFindings(t *testing.T) {
	plan, err := (Planner{}).Plan(spec.ContainerSpec{
		ID:         "abc123",
		Privileged: true,
		Capabilities: spec.Capabilities{Complete: true, Bounding: []string{
			"CAP_SYS_ADMIN",
		}},
		Mounts:  []spec.Mount{{Destination: "/etc/shadow"}},
		Devices: []spec.Device{{PathOnHost: "/dev/nvme0n1"}},
	})
	if err != nil {
		t.Fatalf("Plan() error = %v", err)
	}
	if plan.Status != StatusReviewRequired {
		t.Fatalf("status = %q, want %q", plan.Status, StatusReviewRequired)
	}
	if len(plan.Findings) != 4 {
		t.Fatalf("findings = %d, want 4", len(plan.Findings))
	}
	if plan.Findings[0].Action != "review" {
		t.Fatalf("first finding action = %q, want review", plan.Findings[0].Action)
	}
}

func TestPlannerMatchesDeviceWildcards(t *testing.T) {
	plan, err := (Planner{}).Plan(spec.ContainerSpec{
		ID:           "abc123",
		Capabilities: spec.Capabilities{Complete: true},
		Devices:      []spec.Device{{PathOnHost: "/dev/sda"}},
	})
	if err != nil {
		t.Fatalf("Plan() error = %v", err)
	}
	if len(plan.Findings) != 1 || plan.Findings[0].Code != "high-risk-device" {
		t.Fatalf("findings = %#v, want one high-risk-device finding", plan.Findings)
	}
}

func TestPlannerReportsIncompleteCapabilities(t *testing.T) {
	plan, err := (Planner{}).Plan(spec.ContainerSpec{ID: "abc123"})
	if err != nil {
		t.Fatalf("Plan() error = %v", err)
	}
	if plan.Status != StatusReviewRequired {
		t.Fatalf("status = %q, want %q", plan.Status, StatusReviewRequired)
	}
	if len(plan.Findings) != 1 || plan.Findings[0].Code != "capabilities-incomplete" {
		t.Fatalf("findings = %#v, want incomplete-capabilities finding", plan.Findings)
	}
}

func TestPlannerReturnsReadyForCompleteUnremarkableSpec(t *testing.T) {
	plan, err := (Planner{}).Plan(spec.ContainerSpec{
		ID:           "abc123",
		Capabilities: spec.Capabilities{Complete: true},
	})
	if err != nil {
		t.Fatalf("Plan() error = %v", err)
	}
	if plan.Status != StatusReady {
		t.Fatalf("status = %q, want %q", plan.Status, StatusReady)
	}
	if len(plan.Findings) != 0 {
		t.Fatalf("findings = %#v, want none", plan.Findings)
	}
}

func TestPlannerRequiresID(t *testing.T) {
	if _, err := (Planner{}).Plan(spec.ContainerSpec{}); err == nil {
		t.Fatal("Plan() error = nil, want missing ID error")
	}
}

func TestPlannerDigestsAreStableAcrossInputOrdering(t *testing.T) {
	left, err := (Planner{}).Plan(spec.ContainerSpec{
		ID:              "abc123",
		SecurityOptions: []string{"label:type:container_abc", "no-new-privileges"},
		Capabilities: spec.Capabilities{
			Complete: true,
			Bounding: []string{"CAP_NET_BIND_SERVICE", "CAP_CHOWN"},
		},
		Mounts: []spec.Mount{
			{Type: "bind", Source: "/srv/b", Destination: "/b", Options: []string{"ro", "z"}},
			{Type: "bind", Source: "/srv/a", Destination: "/a", Options: []string{"z", "ro"}},
		},
		Devices: []spec.Device{{PathOnHost: "/dev/sdb", PathInContainer: "/dev/sdb", Permissions: "r"}},
	})
	if err != nil {
		t.Fatalf("left Plan() error = %v", err)
	}

	right, err := (Planner{}).Plan(spec.ContainerSpec{
		ID:              "abc123",
		SecurityOptions: []string{"no-new-privileges", "label:type:container_abc"},
		Capabilities: spec.Capabilities{
			Complete: true,
			Bounding: []string{"CAP_CHOWN", "CAP_NET_BIND_SERVICE"},
		},
		Mounts: []spec.Mount{
			{Type: "bind", Source: "/srv/a", Destination: "/a", Options: []string{"ro", "z"}},
			{Type: "bind", Source: "/srv/b", Destination: "/b", Options: []string{"z", "ro"}},
		},
		Devices: []spec.Device{{PathOnHost: "/dev/sdb", PathInContainer: "/dev/sdb", Permissions: "r"}},
	})
	if err != nil {
		t.Fatalf("right Plan() error = %v", err)
	}
	if left.InputDigest != right.InputDigest || left.PolicyDigest != right.PolicyDigest {
		t.Fatalf("digests differ for equivalent inputs: left=(%s,%s), right=(%s,%s)", left.InputDigest, left.PolicyDigest, right.InputDigest, right.PolicyDigest)
	}
}

func TestPlannerInputDigestChangesWithSecurityInput(t *testing.T) {
	base, err := (Planner{}).Plan(spec.ContainerSpec{ID: "abc123", Capabilities: spec.Capabilities{Complete: true}})
	if err != nil {
		t.Fatalf("base Plan() error = %v", err)
	}
	changed, err := (Planner{}).Plan(spec.ContainerSpec{
		ID: "abc123", Privileged: true, Capabilities: spec.Capabilities{Complete: true},
	})
	if err != nil {
		t.Fatalf("changed Plan() error = %v", err)
	}
	if base.InputDigest == changed.InputDigest {
		t.Fatal("InputDigest is unchanged after security input changed")
	}
	if base.PolicyDigest == changed.PolicyDigest {
		t.Fatal("PolicyDigest is unchanged after policy findings changed")
	}
}
