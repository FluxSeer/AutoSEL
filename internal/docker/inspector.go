package docker

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"

	"github.com/FluxSeer/AutoSEL/internal/spec"
)

// InspectClient is the read-only Docker API surface required by Inspector.
// Keeping this interface small makes the adapter testable without a daemon.
type InspectClient interface {
	ContainerInspect(context.Context, string) (types.ContainerJSON, error)
}

type Inspector struct {
	client InspectClient
}

func NewInspector(client InspectClient) (*Inspector, error) {
	if client == nil {
		return nil, fmt.Errorf("docker inspect client is required")
	}
	return &Inspector{client: client}, nil
}

// Inspect reads one container through Docker's inspect API and normalizes only
// the security-relevant fields. It never performs a Docker mutation.
func (i *Inspector) Inspect(ctx context.Context, containerID string) (spec.ContainerSpec, error) {
	if strings.TrimSpace(containerID) == "" {
		return spec.ContainerSpec{}, fmt.Errorf("container id is required")
	}
	if ctx == nil {
		return spec.ContainerSpec{}, fmt.Errorf("context is required")
	}

	inspected, err := i.client.ContainerInspect(ctx, containerID)
	if err != nil {
		return spec.ContainerSpec{}, fmt.Errorf("inspect container %q: %w", containerID, err)
	}
	if inspected.Config == nil {
		return spec.ContainerSpec{}, fmt.Errorf("inspect container %q: response is missing Config", containerID)
	}
	if inspected.HostConfig == nil {
		return spec.ContainerSpec{}, fmt.Errorf("inspect container %q: response is missing HostConfig", containerID)
	}

	id := inspected.ID
	if id == "" {
		id = containerID
	}
	result := spec.ContainerSpec{
		ID:              id,
		Name:            strings.TrimPrefix(inspected.Name, "/"),
		Image:           inspected.Config.Image,
		Privileged:      inspected.HostConfig.Privileged,
		SecurityOptions: append([]string(nil), inspected.HostConfig.SecurityOpt...),
		Capabilities: spec.Capabilities{
			Complete: false,
			Source:   "docker-inspect-explicit-capabilities",
			Added:    cloneStrings(inspected.HostConfig.CapAdd),
			Dropped:  cloneStrings(inspected.HostConfig.CapDrop),
		},
	}

	for _, mountPoint := range inspected.Mounts {
		result.Mounts = append(result.Mounts, spec.Mount{
			Type:        string(mountPoint.Type),
			Source:      mountPoint.Source,
			Destination: mountPoint.Destination,
			Options:     splitMode(mountPoint.Mode),
		})
	}
	for _, device := range inspected.HostConfig.Devices {
		result.Devices = append(result.Devices, spec.Device{
			PathOnHost:      device.PathOnHost,
			PathInContainer: device.PathInContainer,
			Permissions:     device.CgroupPermissions,
		})
	}

	return result, nil
}

func cloneStrings(values []string) []string {
	return append([]string(nil), values...)
}

func splitMode(mode string) []string {
	if mode == "" {
		return nil
	}
	parts := strings.Split(mode, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}
