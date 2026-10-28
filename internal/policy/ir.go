package policy

import (
	"fmt"
	"regexp"
	"strings"
)

type RuleKind string

const RuleDeny RuleKind = "deny"

// PolicyRule is intentionally defined before enforcement rules are enabled.
// The first CIL slice rejects non-empty rules so a finding cannot silently
// become an enforcement decision without a reviewed rule contract.
type PolicyRule struct {
	Kind        RuleKind
	SourceType  string
	TargetType  string
	ObjectClass string
	Permissions []string
}

type PolicyIR struct {
	ModuleName     string
	TypeName       string
	ReviewRequired bool
	Findings       []Finding
	Rules          []PolicyRule
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
		ModuleName:     moduleName,
		TypeName:       moduleName + "_t",
		ReviewRequired: plan.Status == StatusReviewRequired,
		Findings:       append([]Finding(nil), plan.Findings...),
	}, nil
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
