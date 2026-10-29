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
	RunInDir(context.Context, string, string, ...string) error
}

type OSCommandRunner struct{}

func (OSCommandRunner) Run(ctx context.Context, name string, args ...string) error {
	return runCommand(ctx, "", name, args...)
}

func (OSCommandRunner) RunInDir(ctx context.Context, directory, name string, args ...string) error {
	return runCommand(ctx, directory, name, args...)
}

func runCommand(ctx context.Context, directory, name string, args ...string) error {
	command := exec.CommandContext(ctx, name, args...)
	if strings.TrimSpace(directory) != "" {
		command.Dir = directory
	}
	output, err := command.CombinedOutput()
	if err != nil {
		message := strings.TrimSpace(string(output))
		if message == "" {
			return err
		}
		return fmt.Errorf("%w: %s", err, message)
	}
	return nil
}

const (
	ProfileDirect = "direct"
	ProfileSecilc = "secilc"
	inputToken    = "{{input}}"
	outputToken   = "{{output}}"
)

type CompilerProfile struct {
	Name string
	Args []string
}

func NewCompilerProfile(name string) (CompilerProfile, error) {
	switch strings.TrimSpace(name) {
	case "", ProfileDirect:
		return CompilerProfile{Name: ProfileDirect, Args: []string{inputToken}}, nil
	case ProfileSecilc:
		return CompilerProfile{Name: ProfileSecilc, Args: []string{"-o", outputToken, inputToken}}, nil
	default:
		return CompilerProfile{}, fmt.Errorf("unsupported SELinux compiler profile %q", name)
	}
}

type Compiler struct {
	Binary  string
	Profile CompilerProfile
	Runner  CommandRunner
}

func NewCompiler(binary string, runner CommandRunner) (*Compiler, error) {
	profile, err := NewCompilerProfile(ProfileDirect)
	if err != nil {
		return nil, err
	}
	return NewCompilerWithProfile(binary, profile, runner)
}

func NewCompilerWithProfile(binary string, profile CompilerProfile, runner CommandRunner) (*Compiler, error) {
	if strings.TrimSpace(binary) == "" {
		return nil, fmt.Errorf("SELinux compiler binary is required")
	}
	if strings.TrimSpace(profile.Name) == "" || len(profile.Args) == 0 {
		return nil, fmt.Errorf("SELinux compiler profile is required")
	}
	hasInput := false
	for _, arg := range profile.Args {
		if arg == inputToken {
			hasInput = true
			break
		}
	}
	if !hasInput {
		return nil, fmt.Errorf("SELinux compiler profile must include %s", inputToken)
	}
	if runner == nil {
		runner = OSCommandRunner{}
	}
	profile.Args = append([]string(nil), profile.Args...)
	return &Compiler{Binary: binary, Profile: profile, Runner: runner}, nil
}

// Validate writes a temporary candidate CIL file and invokes the configured
// compiler profile. Any profile output is also kept inside the temporary
// workspace. It never installs a module or writes the system policy store.
func (c *Compiler) Validate(ctx context.Context, artifact policy.CILArtifact) error {
	if c == nil || c.Runner == nil || strings.TrimSpace(c.Binary) == "" || len(c.Profile.Args) == 0 {
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
	outputPath := filepath.Join(directory, "compiled-policy.kern")
	args := make([]string, 0, len(c.Profile.Args))
	for _, arg := range c.Profile.Args {
		switch arg {
		case inputToken:
			args = append(args, inputPath)
		case outputToken:
			args = append(args, outputPath)
		default:
			args = append(args, arg)
		}
	}
	if err := c.Runner.RunInDir(ctx, directory, c.Binary, args...); err != nil {
		return fmt.Errorf("SELinux compiler validation: %w", err)
	}
	return nil
}
