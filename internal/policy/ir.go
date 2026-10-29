package policy

import (
	"fmt"
	"regexp"
	"strings"
)

type RuleKind string

const RuleDeny RuleKind = "deny"

const RuleContractVersion = "autosel-policy-rule/v1"

type RuleApprovalStatus string

const (
	RuleApprovalPending  RuleApprovalStatus = "pending"
	RuleApprovalApproved RuleApprovalStatus = "approved"
	RuleApprovalRejected RuleApprovalStatus = "rejected"
)

type RuleProvenance struct {
	FindingCodes []string `json:"findingCodes,omitempty"`
	Source       string   `json:"source"`
	Reason       string   `json:"reason"`
}

type RuleApproval struct {
	Status     RuleApprovalStatus `json:"status"`
	Approver   string             `json:"approver,omitempty"`
	DecisionID string             `json:"decisionId,omitempty"`
}

// PolicyRule is intentionally defined before enforcement rules are enabled.
// The first CIL slice rejects non-empty rules so a finding cannot silently
// become an enforcement decision without a reviewed rule contract.
type PolicyRule struct {
	Kind        RuleKind       `json:"kind"`
	SourceType  string         `json:"sourceType"`
	TargetType  string         `json:"targetType"`
	ObjectClass string         `json:"objectClass"`
	Permissions []string       `json:"permissions"`
	Provenance  RuleProvenance `json:"provenance"`
	Approval    RuleApproval   `json:"approval"`
}

type PolicyIR struct {
	RuleContractVersion string       `json:"ruleContractVersion"`
	ModuleName          string       `json:"moduleName"`
	TypeName            string       `json:"typeName"`
	ReviewRequired      bool         `json:"reviewRequired"`
	Findings            []Finding    `json:"findings,omitempty"`
	Rules               []PolicyRule `json:"rules,omitempty"`
}

func BuildPolicyIR(plan Plan) (PolicyIR, error) {
	if plan.ModuleName == "" {
		return PolicyIR{}, fmt.Errorf("policy plan module name is required")
	}
	moduleName := cilIdentifier(plan.ModuleName)
	if moduleName == "" {
		return PolicyIR{}, fmt.Errorf("policy plan module name is not a valid CIL identifier")
	}

	return PolicyIR{
		RuleContractVersion: RuleContractVersion,
		ModuleName:          moduleName,
		TypeName:            moduleName + "_t",
		ReviewRequired:      plan.Status == StatusReviewRequired,
		Findings:            append([]Finding(nil), plan.Findings...),
	}, nil
}

// ValidateRuleContract validates the typed review and provenance envelope.
// It does not authorize, compile, or emit the rule.
func ValidateRuleContract(rule PolicyRule) error {
	if rule.Kind != RuleDeny {
		return fmt.Errorf("unsupported policy rule kind %q", rule.Kind)
	}
	identifiers := []struct {
		name  string
		value string
	}{
		{name: "source type", value: rule.SourceType},
		{name: "target type", value: rule.TargetType},
		{name: "object class", value: rule.ObjectClass},
	}
	for _, identifier := range identifiers {
		if !cilIdentifierPattern.MatchString(identifier.value) {
			return fmt.Errorf("policy rule %s is not a valid CIL identifier", identifier.name)
		}
	}
	if len(rule.Permissions) == 0 {
		return fmt.Errorf("policy rule permissions are required")
	}
	seenPermissions := make(map[string]struct{}, len(rule.Permissions))
	for _, permission := range rule.Permissions {
		if !cilIdentifierPattern.MatchString(permission) {
			return fmt.Errorf("policy rule permission %q is not a valid CIL identifier", permission)
		}
		if _, exists := seenPermissions[permission]; exists {
			return fmt.Errorf("policy rule permission %q is duplicated", permission)
		}
		seenPermissions[permission] = struct{}{}
	}
	if strings.TrimSpace(rule.Provenance.Source) == "" || strings.TrimSpace(rule.Provenance.Reason) == "" {
		return fmt.Errorf("policy rule provenance source and reason are required")
	}
	switch rule.Approval.Status {
	case RuleApprovalPending:
		if rule.Approval.Approver != "" || rule.Approval.DecisionID != "" {
			return fmt.Errorf("pending policy rule approval cannot have approver or decision ID")
		}
	case RuleApprovalApproved:
		if strings.TrimSpace(rule.Approval.Approver) == "" || strings.TrimSpace(rule.Approval.DecisionID) == "" {
			return fmt.Errorf("approved policy rule requires approver and decision ID")
		}
	case RuleApprovalRejected:
		return fmt.Errorf("rejected policy rule cannot enter Policy IR")
	default:
		return fmt.Errorf("unsupported policy rule approval status %q", rule.Approval.Status)
	}
	return nil
}

var cilIdentifierPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

func cilIdentifier(value string) string {
	var builder strings.Builder
	for _, char := range strings.ToLower(value) {
		switch {
		case char >= 'a' && char <= 'z':
			builder.WriteRune(char)
		case char >= '0' && char <= '9':
			builder.WriteRune(char)
		default:
			builder.WriteRune('_')
		}
	}
	identifier := strings.Trim(builder.String(), "_")
	if identifier == "" || !cilIdentifierPattern.MatchString(identifier) {
		return ""
	}
	return identifier
}
