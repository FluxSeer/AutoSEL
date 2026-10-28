package spec

import (
	"encoding/json"
	"fmt"
)

// ContainerSpec is the normalized, side-effect-free input to policy planning.
// It intentionally contains only the security-relevant subset needed by the
// first planning slice. Docker and OCI adapters can populate it independently.
type ContainerSpec struct {
	ID              string       `json:"id"`
	Name            string       `json:"name,omitempty"`
	Image           string       `json:"image,omitempty"`
	Privileged      bool         `json:"privileged,omitempty"`
	Capabilities    Capabilities `json:"capabilities,omitempty"`
	Mounts          []Mount      `json:"mounts,omitempty"`
	Devices         []Device     `json:"devices,omitempty"`
	SecurityOptions []string     `json:"securityOptions,omitempty"`
}

type Capabilities struct {
	// Complete is true only when the source exposes the complete bounding set.
	// Docker inspect normally exposes explicit additions/drops, not the full
	// runtime-resolved set, so adapters must not infer completeness.
	Complete bool     `json:"complete"`
	Source   string   `json:"source,omitempty"`
	Bounding []string `json:"bounding,omitempty"`
	Added    []string `json:"added,omitempty"`
	Dropped  []string `json:"dropped,omitempty"`
}

type Mount struct {
	Type        string   `json:"type"`
	Source      string   `json:"source,omitempty"`
	Destination string   `json:"destination"`
	Options     []string `json:"options,omitempty"`
}

type Device struct {
	PathOnHost      string `json:"pathOnHost"`
	PathInContainer string `json:"pathInContainer,omitempty"`
	Permissions     string `json:"permissions,omitempty"`
}

// DecodeOCISpec decodes the small OCI runtime-config subset required by the
// planner. Unknown fields are deliberately ignored so newer runtime specs do
// not break this read-only planning path.
func DecodeOCISpec(data []byte, id string) (ContainerSpec, error) {
	if id == "" {
		return ContainerSpec{}, fmt.Errorf("container id is required")
	}

	var raw ociSpec
	if err := json.Unmarshal(data, &raw); err != nil {
		return ContainerSpec{}, fmt.Errorf("decode OCI config: %w", err)
	}

	result := ContainerSpec{
		ID:     id,
		Mounts: append([]Mount(nil), raw.Mounts...),
	}
	if raw.Process != nil && raw.Process.Capabilities != nil {
		result.Capabilities.Complete = true
		result.Capabilities.Source = "oci-runtime-config"
		result.Capabilities.Bounding = append([]string(nil), raw.Process.Capabilities.Bounding...)
	}
	if raw.Linux != nil {
		result.Devices = append([]Device(nil), raw.Linux.Devices...)
	}
	return result, nil
}

type ociSpec struct {
	Process *ociProcess `json:"process"`
	Mounts  []Mount     `json:"mounts"`
	Linux   *ociLinux   `json:"linux"`
}

type ociProcess struct {
	Capabilities *Capabilities `json:"capabilities"`
}

type ociLinux struct {
	Devices []Device `json:"devices"`
}
