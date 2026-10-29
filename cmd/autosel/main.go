package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/client"

	inspectdocker "github.com/FluxSeer/AutoSEL/internal/docker"
	"github.com/FluxSeer/AutoSEL/internal/policy"
	"github.com/FluxSeer/AutoSEL/internal/selinux"
	"github.com/FluxSeer/AutoSEL/internal/spec"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] == "--help" || os.Args[1] == "-h" {
		usage(os.Stdout)
		return
	}

	switch os.Args[1] {
	case "plan":
		if err := runPlan(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, "autosel plan:", err)
			os.Exit(1)
		}
	case "validate":
		if err := runValidate(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, "autosel validate:", err)
			os.Exit(1)
		}
	default:
		usage(os.Stderr)
		os.Exit(2)
	}
}

func runPlan(args []string) error {
	flags := flag.NewFlagSet("plan", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	inputPath := flags.String("input", "-", "OCI config.json path, or - for stdin")
	containerID := flags.String("id", "", "container ID for OCI config input")
	dockerContainer := flags.String("container", "", "inspect this Docker container through the read-only Docker API")
	if err := flags.Parse(args); err != nil {
		return err
	}
	normalized, err := loadNormalized(*inputPath, *containerID, *dockerContainer)
	if err != nil {
		return err
	}

	plan, err := (policy.Planner{}).Plan(normalized)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(plan)
}

type validationResult struct {
	Plan              policy.Plan        `json:"plan"`
	Artifact          policy.CILArtifact `json:"artifact"`
	StaticValid       bool               `json:"staticValid"`
	Compiler          string             `json:"compiler,omitempty"`
	CompilerProfile   string             `json:"compilerProfile,omitempty"`
	CompilerValidated bool               `json:"compilerValidated"`
}

func runValidate(args []string) error {
	flags := flag.NewFlagSet("validate", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	inputPath := flags.String("input", "-", "OCI config.json path, or - for stdin")
	containerID := flags.String("id", "", "container ID for OCI config input")
	dockerContainer := flags.String("container", "", "inspect this Docker container through the read-only Docker API")
	compilerBinary := flags.String("compiler", "", "optional SELinux compiler binary")
	compilerProfile := flags.String("compiler-profile", selinux.ProfileDirect, "compiler invocation profile: direct or secilc")
	if err := flags.Parse(args); err != nil {
		return err
	}

	normalized, err := loadNormalized(*inputPath, *containerID, *dockerContainer)
	if err != nil {
		return err
	}
	plan, err := (policy.Planner{}).Plan(normalized)
	if err != nil {
		return err
	}
	ir, err := policy.BuildPolicyIR(plan)
	if err != nil {
		return err
	}
	artifact, err := policy.GenerateCIL(ir)
	if err != nil {
		return err
	}
	if err := policy.ValidateCIL(artifact.Content); err != nil {
		return fmt.Errorf("static CIL validation: %w", err)
	}

	result := validationResult{
		Plan:        plan,
		Artifact:    artifact,
		StaticValid: true,
	}
	if strings.TrimSpace(*compilerBinary) != "" {
		profile, err := selinux.NewCompilerProfile(*compilerProfile)
		if err != nil {
			return err
		}
		compiler, err := selinux.NewCompilerWithProfile(*compilerBinary, profile, nil)
		if err != nil {
			return err
		}
		if err := compiler.Validate(context.Background(), artifact); err != nil {
			return err
		}
		result.Compiler = *compilerBinary
		result.CompilerProfile = profile.Name
		result.CompilerValidated = true
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

func loadNormalized(inputPath, containerID, dockerContainer string) (spec.ContainerSpec, error) {
	if strings.TrimSpace(containerID) != "" && strings.TrimSpace(dockerContainer) != "" {
		return spec.ContainerSpec{}, fmt.Errorf("--id and --container cannot be used together")
	}
	if strings.TrimSpace(dockerContainer) != "" {
		if inputPath != "-" {
			return spec.ContainerSpec{}, fmt.Errorf("--input cannot be used with --container")
		}
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			return spec.ContainerSpec{}, fmt.Errorf("create Docker client: %w", err)
		}
		defer func() { _ = cli.Close() }()
		inspector, err := inspectdocker.NewInspector(cli)
		if err != nil {
			return spec.ContainerSpec{}, err
		}
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		return inspector.Inspect(ctx, dockerContainer)
	}
	if strings.TrimSpace(containerID) == "" {
		return spec.ContainerSpec{}, fmt.Errorf("--id or --container is required")
	}
	input, err := readInput(inputPath)
	if err != nil {
		return spec.ContainerSpec{}, err
	}
	return spec.DecodeOCISpec(input, containerID)
}

func readInput(path string) ([]byte, error) {
	if path == "-" {
		return io.ReadAll(os.Stdin)
	}
	return os.ReadFile(path)
}

func usage(w io.Writer) {
	_, _ = fmt.Fprintln(w, "AutoSEL: side-effect-free SELinux policy planning for containers")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Usage:")
	_, _ = fmt.Fprintln(w, "  autosel plan --id <container-id> --input <oci-config.json>")
	_, _ = fmt.Fprintln(w, "  autosel plan --container <container-id>")
	_, _ = fmt.Fprintln(w, "  cat config.json | autosel plan --id <container-id>")
	_, _ = fmt.Fprintln(w, "  autosel validate --id <container-id> --input <oci-config.json>")
	_, _ = fmt.Fprintln(w, "  autosel validate --container <container-id> [--compiler <binary> --compiler-profile <name>]")
}
