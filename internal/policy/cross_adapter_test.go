package policy

import (
	"context"
	"reflect"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"

	dockeradapter "github.com/FluxSeer/AutoSEL/internal/docker"
	"github.com/FluxSeer/AutoSEL/internal/spec"
)

type crossAdapterInspectClient struct {
	response types.ContainerJSON
}

func (f crossAdapterInspectClient) ContainerInspect(context.Context, string) (types.ContainerJSON, error) {
	return f.response, nil
}

func TestPlannerProducesSamePlanAcrossOCIAndDockerAdapters(t *testing.T) {
	ociData := []byte(`{
        "process": {},
        "mounts": [{"type": "bind", "source": "/srv/app", "destination": "/app", "options": ["z", "ro"]}],
        "linux": {"devices": [{"pathOnHost": "/dev/sdb", "pathInContainer": "/dev/sdb", "permissions": "r"}]}
    }`)
	ociSpec, err := spec.DecodeOCISpec(ociData, "container-1")
	if err != nil {
		t.Fatalf("DecodeOCISpec() error = %v", err)
	}
	ociSpec.Name = "web"
	ociSpec.Image = "alpine:latest"

	inspector, err := dockeradapter.NewInspector(crossAdapterInspectClient{response: types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			ID:   "container-1",
			Name: "/web",
			HostConfig: &container.HostConfig{Resources: container.Resources{
				Devices: []container.DeviceMapping{{
					PathOnHost: "/dev/sdb", PathInContainer: "/dev/sdb", CgroupPermissions: "r",
				}},
			}},
		},
		Config: &container.Config{Image: "alpine:latest"},
		Mounts: []types.MountPoint{{
			Type:        mount.TypeBind,
			Source:      "/srv/app",
			Destination: "/app",
			Mode:        "ro,z",
		}},
	}})
	if err != nil {
		t.Fatalf("NewInspector() error = %v", err)
	}
	dockerSpec, err := inspector.Inspect(context.Background(), "container-1")
	if err != nil {
		t.Fatalf("Inspect() error = %v", err)
	}

	ociPlan, err := (Planner{}).Plan(ociSpec)
	if err != nil {
		t.Fatalf("OCI Plan() error = %v", err)
	}
	dockerPlan, err := (Planner{}).Plan(dockerSpec)
	if err != nil {
		t.Fatalf("Docker Plan() error = %v", err)
	}
	if !reflect.DeepEqual(ociPlan, dockerPlan) {
		t.Fatalf("plans differ across adapters:\nOCI=%#v\nDocker=%#v", ociPlan, dockerPlan)
	}
}
