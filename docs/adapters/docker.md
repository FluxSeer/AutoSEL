# Docker inspection adapter

The Docker adapter is an experimental, read-only source for policy planning.
It calls `ContainerInspect` and maps only security-relevant fields into the
internal `ContainerSpec` model:

- container identity, name, and image;
- privileged mode and security options;
- explicit capability additions and drops;
- observed mounts and mount options;
- host device mappings.

## Capability completeness

Docker inspect exposes explicit capability changes, but it does not by itself
prove the complete runtime-resolved bounding set for every daemon/runtime
combination. AutoSEL therefore marks Docker-derived capabilities as incomplete
and produces a review finding instead of inferring a full allowlist.

This is intentional: missing security evidence must remain visible rather than
become a false claim of least privilege.

## CLI

```bash
GOWORK=off go run ./cmd/autosel plan --container <container-id>
```

The command creates a Docker client from the normal Docker environment, uses a
15-second read-only inspection timeout, and emits the same deterministic plan
format as OCI-file input. It does not call container stop, remove, create,
start, export, import, or SELinux commands.

## Qualification boundary

The adapter has fake-client unit coverage for normalization and error handling.
The manual `Docker read-only qualification` workflow and
`make docker-readonly-qualification` harness can qualify a pre-existing
container on an isolated lab daemon. They only inspect and plan; they never
create, start, stop, remove, exec into, export, or modify a container. Live
Docker Engine, rootless Docker, Docker Desktop, SELinux labeling, and
distribution compatibility remain unverified until each environment has
recorded evidence.
