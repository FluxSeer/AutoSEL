package policy

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/FluxSeer/AutoSEL/internal/spec"
)

func TestGenerateCILMatchesGoldenCandidate(t *testing.T) {
	plan, err := (Planner{}).Plan(spec.ContainerSpec{
		ID:           "demo-container",
		Capabilities: spec.Capabilities{Complete: true},
	})
	if err != nil {
		t.Fatalf("Plan() error = %v", err)
	}
	ir, err := BuildPolicyIR(plan)
	if err != nil {
		t.Fatalf("BuildPolicyIR() error = %v", err)
	}
	artifact, err := GenerateCIL(ir)
	if err != nil {
		t.Fatalf("GenerateCIL() error = %v", err)
	}
	golden, err := os.ReadFile(filepath.Join("testdata", "candidate.cil"))
	if err != nil {
		t.Fatalf("read golden CIL: %v", err)
	}
	if artifact.Content != string(golden) {
		t.Fatalf("generated CIL = %q, want %q", artifact.Content, string(golden))
	}
	if err := ValidateCIL(artifact.Content); err != nil {
		t.Fatalf("ValidateCIL() error = %v", err)
	}
	if !strings.HasPrefix(artifact.Digest, "sha256:") {
		t.Fatalf("CIL digest = %q, want sha256 prefix", artifact.Digest)
	}
}

func TestGenerateCILRejectsUnreviewedEnforcementRules(t *testing.T) {
	artifact, err := GenerateCIL(PolicyIR{
		ModuleName: "autosel_demo",
		TypeName:   "autosel_demo_t",
		Rules: []PolicyRule{{
			Kind:        RuleDeny,
			SourceType:  "container_t",
			TargetType:  "shadow_t",
			ObjectClass: "file",
			Permissions: []string{"read"},
		}},
	})
	if err == nil {
		t.Fatalf("GenerateCIL() artifact = %#v, want unsupported-rule error", artifact)
	}
}

func TestValidateCILRejectsMalformedContent(t *testing.T) {
	for _, content := range []string{
		"",
		"(block broken\n    (type broken_t)\n",
		"(allow broken broken (file (read)))\n",
	} {
		if err := ValidateCIL(content); err == nil {
			t.Fatalf("ValidateCIL(%q) error = nil, want error", content)
		}
	}
}
