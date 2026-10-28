package selinux

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/FluxSeer/AutoSEL/internal/policy"
)

type fakeRunner struct {
	name string
	args []string
	err  error
}

func (r *fakeRunner) Run(_ context.Context, name string, args ...string) error {
	r.name = name
	r.args = append([]string(nil), args...)
	if len(args) == 1 {
		data, err := os.ReadFile(args[0])
		if err != nil {
			return err
		}
		if !strings.Contains(string(data), "(block autosel_demo") {
			return os.ErrInvalid
		}
	}
	return r.err
}

func TestCompilerUsesTemporaryCandidateWithoutInstalling(t *testing.T) {
	runner := &fakeRunner{}
	compiler, err := NewCompiler("fake-secilc", runner)
	if err != nil {
		t.Fatalf("NewCompiler() error = %v", err)
	}
	err = compiler.Validate(context.Background(), policy.CILArtifact{
		ModuleName: "autosel_demo",
		TypeName:   "autosel_demo_t",
		Content:    "(block autosel_demo\n    (type autosel_demo_t)\n)\n",
	})
	if err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if runner.name != "fake-secilc" || len(runner.args) != 1 {
		t.Fatalf("runner invocation = (%q, %#v), want compiler and one input path", runner.name, runner.args)
	}
	if filepath.Base(runner.args[0]) != "candidate.cil" {
		t.Fatalf("compiler input = %q, want candidate.cil", runner.args[0])
	}
}

func TestCompilerPropagatesCommandFailure(t *testing.T) {
	runner := &fakeRunner{err: os.ErrPermission}
	compiler, err := NewCompiler("fake-secilc", runner)
	if err != nil {
		t.Fatalf("NewCompiler() error = %v", err)
	}
	if err := compiler.Validate(context.Background(), policy.CILArtifact{
		Content: "(block autosel_demo\n    (type autosel_demo_t)\n)\n",
	}); err == nil {
		t.Fatal("Validate() error = nil, want compiler failure")
	}
}
