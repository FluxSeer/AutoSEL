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
	directory string
	name      string
	args      []string
	err       error
}

func (r *fakeRunner) Run(_ context.Context, name string, args ...string) error {
	r.name = name
	r.args = append([]string(nil), args...)
	for _, arg := range args {
		if filepath.Base(arg) == "candidate.cil" {
			data, err := os.ReadFile(arg)
			if err != nil {
				return err
			}
			if !strings.Contains(string(data), "(block autosel_demo") {
				return os.ErrInvalid
			}
		}
	}
	return r.err
}

func (r *fakeRunner) RunInDir(ctx context.Context, directory, name string, args ...string) error {
	r.directory = directory
	return r.Run(ctx, name, args...)
}

func TestNewCompilerProfileRejectsUnknownProfile(t *testing.T) {
	if _, err := NewCompilerProfile("unknown"); err == nil {
		t.Fatal("NewCompilerProfile() error = nil, want unsupported profile error")
	}
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

func TestCompilerUsesSecilcProfileAndTemporaryOutput(t *testing.T) {
	runner := &fakeRunner{}
	profile, err := NewCompilerProfile(ProfileSecilc)
	if err != nil {
		t.Fatalf("NewCompilerProfile() error = %v", err)
	}
	compiler, err := NewCompilerWithProfile("secilc", profile, runner)
	if err != nil {
		t.Fatalf("NewCompilerWithProfile() error = %v", err)
	}
	if err := compiler.Validate(context.Background(), policy.CILArtifact{
		Content: "(block autosel_demo\n    (type autosel_demo_t)\n)\n",
	}); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if runner.name != "secilc" || len(runner.args) != 3 {
		t.Fatalf("runner invocation = (%q, %#v), want secilc with output profile", runner.name, runner.args)
	}
	if runner.args[0] != "-o" || filepath.Base(runner.args[1]) != "compiled-policy.kern" || filepath.Base(runner.args[2]) != "candidate.cil" {
		t.Fatalf("runner args = %#v, want -o temporary output candidate", runner.args)
	}
	if filepath.Dir(runner.args[1]) != filepath.Dir(runner.args[2]) {
		t.Fatalf("compiler paths = %#v, want shared temporary workspace", runner.args)
	}
	if runner.directory != filepath.Dir(runner.args[1]) {
		t.Fatalf("compiler directory = %q, want temporary workspace %q", runner.directory, filepath.Dir(runner.args[1]))
	}
}

func TestCompilerProfileRequiresInputToken(t *testing.T) {
	_, err := NewCompilerWithProfile("fake-secilc", CompilerProfile{
		Name: "invalid",
		Args: []string{"-o", outputToken},
	}, nil)
	if err == nil {
		t.Fatal("NewCompilerWithProfile() error = nil, want missing input token error")
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

func TestCompilerRejectsStaticInvalidCILBeforeRunner(t *testing.T) {
	runner := &fakeRunner{}
	compiler, err := NewCompiler("fake-secilc", runner)
	if err != nil {
		t.Fatalf("NewCompiler() error = %v", err)
	}
	if err := compiler.Validate(context.Background(), policy.CILArtifact{Content: "(block broken\n"}); err == nil {
		t.Fatal("Validate() error = nil, want static validation failure")
	}
	if runner.name != "" {
		t.Fatalf("runner name = %q, want no compiler invocation", runner.name)
	}
}
