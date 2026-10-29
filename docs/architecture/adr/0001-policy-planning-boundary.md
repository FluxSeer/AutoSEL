# ADR 0001: establish a read-only policy planning boundary

- Status: Accepted for the first refactoring slice
- Date: 2026-10-28
- Scope: internal architecture; no runtime mutation contract

## Context

The original AutoSEL prototype discovers Docker storage paths, listens for a
subset of Docker events, generates CIL directly from untyped maps, and may
stop, remove, export, and recreate containers. That flow cannot yet prove
configuration preservation, data safety, SELinux label correctness, or
rollback behavior.

AutoSEL still has a valuable problem to solve: make container SELinux policy
decisions repeatable and reviewable as container configuration changes.

## Decision

Adopt a pure first boundary:

```text
runtime config -> typed ContainerSpec -> deterministic PolicyPlan
```

The first supported command is `autosel plan`. It reports high-risk
capabilities, sensitive mount destinations, privileged mode, and high-risk
device patterns as review findings. It does not deny, stop, delete, recreate,
or relabel a container.

Docker inspection, event collection, CIL compilation, SELinux loading, and
container replacement remain separate future stages.

## Alternatives considered

### Keep the prototype flow and patch failures

Rejected because it would preserve unsafe lifecycle coupling and make it hard
to prove that a policy change did not change application behavior or data.

### Generate and load CIL in the planner

Rejected because policy compilation and loading are privileged side effects
that require a separate validator, authorization, audit, and rollback contract.

### Use untyped maps as the shared model

Rejected because malformed or missing runtime fields currently become silent
empty inputs or panics. Typed normalization makes edge cases testable.

## Consequences

Positive:

- planner behavior can be tested without Docker, Linux, or SELinux;
- policy findings are deterministic and reviewable;
- dangerous runtime mutation is removed from the default entrypoint;
- Docker and SELinux compatibility can be qualified independently.

Costs:

- the initial command does not yet enforce a policy;
- existing legacy packages remain temporarily for migration and are not part
  of the supported runtime path;
- a later execution design must explicitly solve pre-create labeling or safe
  full-spec replacement.

## Definition of done

- `ContainerSpec` is typed and has an OCI adapter;
- planner output is deterministic and machine-readable;
- high-risk device wildcards are matched as patterns, not literal map keys;
- planner tests cover positive, negative, and missing-input cases;
- root `main` no longer starts the legacy mutating loop;
- README, capability boundary, security guidance, and roadmap state the
  supported and unsupported behavior.

## Verification

The implementation slice must pass `go test ./...`, `go vet ./...`, and
`go build ./...`. Docker/SELinux runtime and CIL compilation remain explicitly
unverified until their adapters and test environment exist.
