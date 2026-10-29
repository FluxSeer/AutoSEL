# AutoSEL

[![CI](https://github.com/FluxSeer/AutoSEL/actions/workflows/ci.yml/badge.svg)](https://github.com/FluxSeer/AutoSEL/actions/workflows/ci.yml)
[![License](https://img.shields.io/github/license/FluxSeer/AutoSEL)](LICENSE)

AutoSEL is an experimental, policy-first tool for planning SELinux controls
for Docker containers. It is designed to make container security decisions
reviewable, deterministic, and auditable before any runtime mutation occurs.

## Product promise

```text
container configuration
-> normalized security spec
-> deterministic policy plan
-> validation and review
-> explicitly authorized application
```

The first implementation slice is intentionally read-only. It plans findings
from an OCI runtime `config.json` or an experimental Docker
`ContainerInspect` call; it does not load SELinux modules, stop containers,
remove containers, or rebuild images.

## Capability boundary

| Capability | Status | Meaning |
| --- | --- | --- |
| OCI config normalization | Supported | Typed, side-effect-free subset for planning. |
| Policy planning | Supported | Deterministic findings for risky capabilities, mounts, and devices. |
| Docker API inspection | Experimental | Read-only `ContainerInspect` adapter; daemon/runtime qualification is still pending. |
| Policy IR | Supported | Typed, review-aware intermediate representation. |
| Candidate CIL generation | Experimental | Deterministic declaration-only scaffold; not enforcement-ready. |
| Static CIL validation | Experimental | Dependency-free structural validation. |
| SELinux compiler validation | Experimental | Optional compiler adapter with a Linux `secilc` qualification harness; distro matrix pending. |
| Docker event reconciliation | Planned | Bounded queue, retry, resync, and audit contract. |
| Container replacement | Not supported | No safe execution contract exists yet. |

## Quick start

Requirements: Go version declared in `go.mod` for the development toolchain.
OCI-file planning does not require Docker or SELinux; live-container planning
requires a reachable Docker daemon.

```bash
GOWORK=off go test ./...
GOWORK=off go run ./cmd/autosel plan \
  --id demo-container \
  --input testdata/oci/config.json

GOWORK=off go run ./cmd/autosel validate \
  --id demo-container \
  --input testdata/oci/config.json

# Optional: inspect a live Docker container without mutating it.
GOWORK=off go run ./cmd/autosel plan --container <container-id>
```

The result is JSON with a `status` of either `ready` or `review-required` and
machine-readable findings. `review-required` is a safety decision boundary,
not an automatic deny or a claim of malicious behavior.

`validate` emits the plan, candidate CIL artifact, static validation result, and
versioned evidence metadata. The evidence explicitly records that policy
loading and runtime enforcement were not tested in this read-only slice.
The generated CIL is a declaration-only scaffold in this release; it does not
enforce the findings and is never loaded automatically.

## Architecture

```text
OCI/Docker adapter
    -> normalized ContainerSpec
    -> Policy Planner
    -> PolicyPlan and findings
    -> validator
    -> explicit apply path (future)
```

See [docs/README.md](docs/README.md),
[docs/architecture/overview.md](docs/architecture/overview.md),
[support and qualification matrix](docs/support-matrix.md), and the
[policy planning boundary ADR](docs/architecture/adr/0001-policy-planning-boundary.md).

## Open-source development

Read [CONTRIBUTING.md](CONTRIBUTING.md), [SECURITY.md](SECURITY.md), and
[ROADMAP.md](ROADMAP.md) before changing runtime behavior. The project follows
FluxSeer's read-only-first, explicit-capability, evidence-backed engineering
style.

## Current limitations

- Linux SELinux runtime behavior is not yet qualified by CI.
- Linux CIL compiler qualification covers the repository fixture only; policy
  base, labels, kernel, and runtime compatibility remain unqualified.
- Fedora qualification is manual-only on an isolated ephemeral runner and is
  not a pull-request required check; three Fedora 44 compiler runs have
  passed.
- Docker API inspection is implemented as an experimental read-only adapter;
  Docker daemon/runtime qualification is not yet in CI.
- Docker event integration is not yet implemented.
- The planner does not claim complete least-privilege policy coverage.
- No generated CIL is loaded automatically.
- Candidate CIL is declaration-only until a reviewed enforcement-rule contract
  exists.
- The original mutation implementation is retained only under
  `internal/legacy/` during migration and is not a supported runtime path.
