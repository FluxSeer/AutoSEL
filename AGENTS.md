# AutoSEL repository instructions

AutoSEL is a Linux/Docker/SELinux security project. Correctness and safety
take precedence over implementation speed.

## Product boundary

- `plan` and `validate` are read-only by default.
- Container mutation is not supported until a separate, reviewed execution
  contract exists.
- Never delete or recreate a container without an explicit authorization,
  complete configuration snapshot, data-preservation check, and rollback path.
- A policy finding is not proof that a container is malicious or unsafe; keep
  observed input, rule inference, and operator decision distinct.

## Engineering rules

- Prefer typed normalized specs over `map[string]interface{}`.
- Keep Docker and SELinux integrations behind adapters.
- Make policy planning deterministic and testable without Docker or SELinux.
- Do not claim distro, Docker, Podman, or kernel compatibility without testing.
- Do not log secrets, tokens, full container environment variables, or raw
  sensitive configuration.
- Every meaningful change must document its verification and remaining risks.

## Definition of done

- `GOWORK=off go test ./...`
- `GOWORK=off go vet ./...`
- `GOWORK=off go build ./...`
- relevant docs, security boundary, and migration notes are updated
- unverified environments and runtime scenarios are clearly listed
