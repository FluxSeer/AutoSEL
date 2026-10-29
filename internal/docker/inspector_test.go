package docker

import (
	"context"
	"errors"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/strslice"
)

type fakeInspectClient struct {
	response types.ContainerJSON
	err      error
	seenID   string
}

func (f *fakeInspectClient) ContainerInspect(_ context.Context, id string) (types.ContainerJSON, error) {
	f.seenID = id
	return f.response, f.err
}

func TestInspectorNormalizesReadOnlySecurityInput(t *testing.T) {
	client := &fakeInspectClient{response: types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			ID:   "real-container-id",
			Name: "/web",
			HostConfig: &container.HostConfig{
				Privileged:  true,
				SecurityOpt: []string{"label:type:container_web"},
				CapAdd:      strslice.StrSlice{"CAP_SYS_ADMIN"},
				CapDrop:     strslice.StrSlice{"CAP_NET_RAW"},
				Resources: container.Resources{Devices: []container.DeviceMapping{{
					PathOnHost: "/dev/nvme0n1", PathInContainer: "/dev/nvme0n1", CgroupPermissions: "rwm",
				}}},
			},
		},
		Config: &container.Config{Image: "alpine:latest"},
		Mounts: []types.MountPoint{{
			Type:        mount.TypeBind,
			Source:      "/etc/shadow",
			Destination: "/etc/shadow",
			Mode:        "ro,z",
		}},
	}}

	inspector, err := NewInspector(client)
	if err != nil {
		t.Fatalf("NewInspector() error = %v", err)
	}
	got, err := inspector.Inspect(context.Background(), "lookup-id")
	if err != nil {
		t.Fatalf("Inspect() error = %v", err)
	}
	if client.seenID != "lookup-id" {
		t.Fatalf("ContainerInspect() id = %q, want lookup-id", client.seenID)
	}
	if got.ID != "real-container-id" || got.Name != "web" || got.Image != "alpine:latest" {
		t.Fatalf("identity = %#v, want normalized identity", got)
	}
	if !got.Privileged || len(got.SecurityOptions) != 1 {
		t.Fatalf("host security settings = %#v, want privileged and security option", got)
	}
	if got.Capabilities.Complete || len(got.Capabilities.Added) != 1 || len(got.Capabilities.Dropped) != 1 {
		t.Fatalf("capabilities = %#v, want explicit incomplete capability set", got.Capabilities)
	}
	if len(got.Mounts) != 1 || got.Mounts[0].Options[0] != "ro" {
		t.Fatalf("mounts = %#v, want normalized mount options", got.Mounts)
	}
	if len(got.Devices) != 1 || got.Devices[0].PathOnHost != "/dev/nvme0n1" {
		t.Fatalf("devices = %#v, want normalized device", got.Devices)
	}
}

func TestInspectorPropagatesDaemonErrors(t *testing.T) {
	wantErr := errors.New("daemon unavailable")
	inspector, err := NewInspector(&fakeInspectClient{err: wantErr})
	if err != nil {
		t.Fatalf("NewInspector() error = %v", err)
	}
	if _, err := inspector.Inspect(context.Background(), "container-id"); !errors.Is(err, wantErr) {
		t.Fatalf("Inspect() error = %v, want wrapped daemon error", err)
	}
}

func TestInspectorRejectsIncompleteDockerResponse(t *testing.T) {
	inspector, err := NewInspector(&fakeInspectClient{response: types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{},
	}})
	if err != nil {
		t.Fatalf("NewInspector() error = %v", err)
	}
	if _, err := inspector.Inspect(context.Background(), "container-id"); err == nil {
		t.Fatal("Inspect() error = nil, want incomplete response error")
	}
}

func TestInspectorRequiresContextAndID(t *testing.T) {
	inspector, err := NewInspector(&fakeInspectClient{})
	if err != nil {
		t.Fatalf("NewInspector() error = %v", err)
	}
	if _, err := inspector.Inspect(context.Background(), ""); err == nil {
		t.Fatal("Inspect() error = nil, want missing ID error")
	}
}
