package evidence

import "github.com/FluxSeer/AutoSEL/internal/policy"

// SchemaVersion identifies the machine-readable validation evidence contract.
// It must change when fields change meaning, not merely when implementation
// details change.
const SchemaVersion = "autosel-validation-evidence/v1"

type CompilerEvidence struct {
	Binary      string `json:"binary"`
	Profile     string `json:"profile"`
	OutputScope string `json:"outputScope"`
}

// ValidationEvidence separates observed validation results from claims that
// AutoSEL deliberately does not make. In particular, policy loading and
// runtime enforcement remain false until a separately reviewed contract is
// implemented.
type ValidationEvidence struct {
	SchemaVersion            string            `json:"schemaVersion"`
	PlannerVersion           string            `json:"plannerVersion"`
	InputDigest              string            `json:"inputDigest"`
	PolicyDigest             string            `json:"policyDigest"`
	ArtifactDigest           string            `json:"artifactDigest"`
	StaticValid              bool              `json:"staticValid"`
	CompilerValidated        bool              `json:"compilerValidated"`
	Compiler                 *CompilerEvidence `json:"compiler,omitempty"`
	PolicyLoaded             bool              `json:"policyLoaded"`
	RuntimeEnforcementTested bool              `json:"runtimeEnforcementTested"`
}

func NewValidationEvidence(plan policy.Plan, artifact policy.CILArtifact, compiler *CompilerEvidence) ValidationEvidence {
	return ValidationEvidence{
		SchemaVersion:            SchemaVersion,
		PlannerVersion:           plan.PlannerVersion,
		InputDigest:              plan.InputDigest,
		PolicyDigest:             plan.PolicyDigest,
		ArtifactDigest:           artifact.Digest,
		StaticValid:              true,
		CompilerValidated:        compiler != nil,
		Compiler:                 compiler,
		PolicyLoaded:             false,
		RuntimeEnforcementTested: false,
	}
}
