# Architecture overview

AutoSEL is being rebuilt around a bounded, deterministic reconciliation
pipeline:

```text
Collect -> Normalize -> Plan -> Validate -> Review -> Apply
```

Only the first three stages are supported in the initial slice.

## Boundaries

### Collection

Docker and OCI access belongs behind an adapter. The core must not discover
containers by scanning `/` or by depending on Docker's internal storage layout.
The first adapter uses only Docker's read-only `ContainerInspect` API and
fails closed when the response lacks security-relevant sections.

### Normalization

Adapters convert runtime-specific input into a typed `ContainerSpec`. The
normalized model contains only the security-relevant input needed by policy
planning and avoids leaking Docker or OCI types through the core.

### Planning

The planner is pure and deterministic. It produces a `PolicyPlan` with
machine-readable findings, a versioned input digest, and a policy digest.
Findings are review prompts, not automatic proof of malicious intent and not
automatic mutation commands.

### Validation

Candidate CIL generation is separate from policy loading. The current slice
generates only a declaration-only module scaffold and performs conservative
structural validation. Type existence against a distribution policy store,
label identity, compiler compatibility, and policy diff remain external
qualification concerns before any SELinux state changes.

### Application

Application is intentionally not supported yet. Any future execution path must
define authorization, idempotency, snapshot preservation, health checks,
rollback, audit identity, and failure semantics before implementation.

## First vertical slice

The current slice accepts an OCI `config.json` or a live Docker container ID,
normalizes capabilities, mounts, and devices, then emits a deterministic plan.
The Docker path is read-only and uses a bounded API timeout. It has no
SELinux dependency or container mutation.
