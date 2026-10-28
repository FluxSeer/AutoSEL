package spec

import "testing"

func TestDecodeOCISpec(t *testing.T) {
	input := []byte(`{
        "process": {"capabilities": {"bounding": ["CAP_SYS_ADMIN", "CAP_NET_BIND_SERVICE"]}},
        "mounts": [{"type": "bind", "source": "/etc/shadow", "destination": "/etc/shadow"}],
        "linux": {"devices": [{"pathOnHost": "/dev/nvme0n1", "pathInContainer": "/dev/nvme0n1", "permissions": "rwm"}]}
    }`)

	got, err := DecodeOCISpec(input, "container-123")
	if err != nil {
		t.Fatalf("DecodeOCISpec() error = %v", err)
	}
	if got.ID != "container-123" {
		t.Fatalf("ID = %q, want container-123", got.ID)
	}
	if len(got.Capabilities.Bounding) != 2 || len(got.Mounts) != 1 || len(got.Devices) != 1 {
		t.Fatalf("decoded security inputs = %#v, want capabilities, mount, and device", got)
	}
	if !got.Capabilities.Complete || got.Capabilities.Source != "oci-runtime-config" {
		t.Fatalf("capability metadata = %#v, want complete OCI source", got.Capabilities)
	}
}

func TestDecodeOCISpecRequiresID(t *testing.T) {
	if _, err := DecodeOCISpec([]byte(`{}`), ""); err == nil {
		t.Fatal("DecodeOCISpec() error = nil, want missing ID error")
	}
}
