package policy

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"

	"github.com/FluxSeer/AutoSEL/internal/spec"
)

const PlannerVersion = "autosel-policy-planner/v0.1.0"

func canonicalSpec(input spec.ContainerSpec) spec.ContainerSpec {
	result := input
	// Adapter provenance is useful for diagnostics but is not part of the
	// normalized security input identity. Equivalent OCI and Docker inputs must
	// hash the same way.
	result.Capabilities.Source = ""
	result.SecurityOptions = sortedCopy(input.SecurityOptions)
	result.Capabilities.Bounding = sortedCopy(input.Capabilities.Bounding)
	result.Capabilities.Added = sortedCopy(input.Capabilities.Added)
	result.Capabilities.Dropped = sortedCopy(input.Capabilities.Dropped)

	result.Mounts = append([]spec.Mount(nil), input.Mounts...)
	for index := range result.Mounts {
		result.Mounts[index].Options = sortedCopy(result.Mounts[index].Options)
	}
	sort.Slice(result.Mounts, func(i, j int) bool {
		left, right := result.Mounts[i], result.Mounts[j]
		if left.Type != right.Type {
			return left.Type < right.Type
		}
		if left.Source != right.Source {
			return left.Source < right.Source
		}
		return left.Destination < right.Destination
	})

	result.Devices = append([]spec.Device(nil), input.Devices...)
	sort.Slice(result.Devices, func(i, j int) bool {
		left, right := result.Devices[i], result.Devices[j]
		if left.PathOnHost != right.PathOnHost {
			return left.PathOnHost < right.PathOnHost
		}
		if left.PathInContainer != right.PathInContainer {
			return left.PathInContainer < right.PathInContainer
		}
		return left.Permissions < right.Permissions
	})
	return result
}

func sortedCopy(values []string) []string {
	result := append([]string(nil), values...)
	sort.Strings(result)
	return result
}

func digestJSON(value any) (string, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return digestBytes(data)
}

func digestBytes(data []byte) (string, error) {
	hash := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(hash[:]), nil
}
