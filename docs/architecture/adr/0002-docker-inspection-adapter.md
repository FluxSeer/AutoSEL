# ADR 0002: use Docker ContainerInspect as a read-only adapter

- Status: Accepted for the second refactoring slice
- Date: 2026-10-28
- Scope: internal adapter boundary; no Docker mutation contract

## Context

The original prototype locates a `config.json` by scanning the host filesystem.
That approach depends on Docker's internal storage layout, can select an
unrelated file, and cannot represent the daemon's complete container
configuration. The policy planner needs a stable, testable input model without
giving the planner Docker or host mutation authority.

## Decision

Add an `internal/docker` adapter with the smallest possible interface:

```text
ContainerInspect(context, id) -> types.ContainerJSON
```

The adapter converts the response into `spec.ContainerSpec`. It does not
expose the Docker client to the planner and does not call any mutating Docker
API. Missing `Config` or `HostConfig` is a terminal normalization error.

Docker explicit capability additions and drops are preserved, but the complete
runtime-resolved capability set is marked unknown. The planner therefore emits
an incomplete-capabilities review finding instead of fabricating a complete
security conclusion.

## Alternatives considered

### Continue reading Docker storage files

Rejected because it is implementation-specific, brittle across storage/runtime
versions, and difficult to test or support for rootless Docker.

### Infer the full capability set from CapAdd and CapDrop

Rejected because daemon defaults and runtime behavior are not fully represented
by those fields. Inference would turn incomplete evidence into a false claim.

### Give the planner a full Docker client

Rejected because it would blur the read-only planning boundary and make future
mutation accidental or difficult to audit.

## Consequences

Positive:

- Docker storage layout is no longer part of the supported design;
- the planner remains pure and unit-testable;
- incomplete capability evidence remains explicit;
- a live-container planning command is available without mutation.

Costs and limitations:

- the adapter requires a Docker daemon for live use;
- live daemon and runtime compatibility are not yet qualified;
- event collection, policy compilation, and application remain separate work.

## Verification

Fake-client unit tests cover identity, security fields, mount/device mapping,
daemon errors, incomplete responses, and invalid inputs. No live Docker or
SELinux qualification is claimed by this ADR.
