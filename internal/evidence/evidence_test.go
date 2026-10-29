package evidence

import (
	"testing"

	"github.com/FluxSeer/AutoSEL/internal/policy"
)

func TestNewValidationEvidenceStatesBoundaries(t *testing.T) {
	plan := policy.Plan{
		PlannerVersion: "autosel-policy-planner/test",
		InputDigest:    "sha256:input",
		PolicyDigest:   "sha256:policy",
	}
	artifact := policy.CILArtifact{Digest: "sha256:artifact"}

	got := NewValidationEvidence(plan, artifact, &CompilerEvidence{
		Binary:      "secilc",
		Profile:     "secilc",
		OutputScope: "temporary-workspace",
	})
	if got.SchemaVersion != SchemaVersion {
		t.Fatalf("schema version = %q, want %q", got.SchemaVersion, SchemaVersion)
	}
	if !got.StaticValid || !got.CompilerValidated {
		t.Fatalf("validation flags = (%v, %v), want both true", got.StaticValid, got.CompilerValidated)
	}
	if got.PolicyLoaded || got.RuntimeEnforcementTested {
		t.Fatalf("unsupported runtime claims = (%v, %v), want both false", got.PolicyLoaded, got.RuntimeEnforcementTested)
	}
	if got.Compiler == nil || got.Compiler.OutputScope != "temporary-workspace" {
		t.Fatalf("compiler evidence = %#v, want temporary workspace", got.Compiler)
	}
}
