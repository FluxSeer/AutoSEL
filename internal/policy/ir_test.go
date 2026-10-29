package policy

import "testing"

func validPolicyRule() PolicyRule {
	return PolicyRule{
		Kind:        RuleDeny,
		SourceType:  "container_t",
		TargetType:  "shadow_t",
		ObjectClass: "file",
		Permissions: []string{"read"},
		Provenance: RuleProvenance{
			FindingCodes: []string{"sensitive-mount"},
			Source:       "planner-finding",
			Reason:       "the rule is derived from an explicitly reviewed finding",
		},
		Approval: RuleApproval{Status: RuleApprovalPending},
	}
}

func TestValidateRuleContractAcceptsPendingProvenance(t *testing.T) {
	if err := ValidateRuleContract(validPolicyRule()); err != nil {
		t.Fatalf("ValidateRuleContract() error = %v", err)
	}
}

func TestValidateRuleContractRequiresApprovalEvidence(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*PolicyRule)
	}{
		{name: "missing provenance source", mutate: func(rule *PolicyRule) { rule.Provenance.Source = "" }},
		{name: "missing provenance reason", mutate: func(rule *PolicyRule) { rule.Provenance.Reason = "" }},
		{name: "approved without approver", mutate: func(rule *PolicyRule) {
			rule.Approval = RuleApproval{Status: RuleApprovalApproved, DecisionID: "decision-1"}
		}},
		{name: "pending with decision", mutate: func(rule *PolicyRule) {
			rule.Approval.DecisionID = "decision-1"
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rule := validPolicyRule()
			test.mutate(&rule)
			if err := ValidateRuleContract(rule); err == nil {
				t.Fatal("ValidateRuleContract() error = nil, want contract error")
			}
		})
	}
}

func TestBuildPolicyIRIncludesRuleContractVersion(t *testing.T) {
	ir, err := BuildPolicyIR(Plan{ModuleName: "autosel-demo"})
	if err != nil {
		t.Fatalf("BuildPolicyIR() error = %v", err)
	}
	if ir.RuleContractVersion != RuleContractVersion {
		t.Fatalf("rule contract version = %q, want %q", ir.RuleContractVersion, RuleContractVersion)
	}
}
