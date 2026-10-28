package selinux

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/FluxSeer/AutoSEL/internal/policy"
)

// CommandRunner is the small process boundary used by Compiler. Tests can
// replace it without executing a host SELinux command.
type CommandRunner interface {
	Run(context.Context, string, ...string) error
}

type OSCommandRunner struct{}

func (OSCommandRunner) Run(ctx context.Context, name string, args ...string) error {
	output, err := exec.CommandContext(ctx, name, args...).CombinedOutput()
	if err != nil {
		message := strings.TrimSpace(string(output))
		if message == "" {
			return err
		}
		return fmt.Errorf("%w: %s", err, message)
	}
	return nil
}

type Compiler struct {
	Binary string
	Runner CommandRunner
}

func NewCompiler(binary string, runner CommandRunner) (*Compiler, error) {
	if strings.TrimSpace(binary) == "" {
		return nil, fmt.Errorf("SELinux compiler binary is required")
	}
	if runner == nil {
		runner = OSCommandRunner{}
	}
	return &Compiler{Binary: binary, Runner: runner}, nil
}

// Validate writes a temporary candidate CIL file and invokes the configured
// compiler with that path as its sole argument. It never installs a module.
// Distribution-specific compiler flags remain an explicit future adapter
// contract rather than being guessed here.
func (c *Compiler) Validate(ctx context.Context, artifact policy.CILArtifact) error {
	if c == nil || c.Runner == nil || strings.TrimSpace(c.Binary) == "" {
		return fmt.Errorf("SELinux compiler is not configured")
	}
	if ctx == nil {
		return fmt.Errorf("context is required")
	}
	if err := policy.ValidateCIL(artifact.Content); err != nil {
		return fmt.Errorf("static CIL validation: %w", err)
	}

	directory, err := os.MkdirTemp("", "autosel-cil-")
	if err != nil {
		return fmt.Errorf("create compiler workspace: %w", err)
	}
	defer func() { _ = os.RemoveAll(directory) }()

	inputPath := filepath.Join(directory, "candidate.cil")
	if err := os.WriteFile(inputPath, []byte(artifact.Content), 0o600); err != nil {
		return fmt.Errorf("write candidate CIL: %w", err)
	}
	if err := c.Runner.Run(ctx, c.Binary, inputPath); err != nil {
		return fmt.Errorf("SELinux compiler validation: %w", err)
	}
	return nil
}
