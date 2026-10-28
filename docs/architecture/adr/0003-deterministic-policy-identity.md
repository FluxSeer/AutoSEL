# ADR 0003: assign deterministic identity to policy plans

- Status: Accepted for the read-only planning contract
- Date: 2026-10-28
- Scope: internal plan identity; no persistence or mutation contract

## Context

Reconciliation, audit, diff, and future rollback need to distinguish the
exact normalized input from the policy result. Container and mount ordering can
vary between sources, so raw JSON bytes are not a stable identity. The project
must also avoid treating a policy plan as verified merely because it has a
digest.

## Decision

Every `PolicyPlan` includes:

- `plannerVersion` — the rule implementation version;
- `inputDigest` — SHA-256 of the canonical normalized `ContainerSpec`;
- `policyDigest` — SHA-256 of the container identity, planner version, status,
  and deterministic findings.

Canonicalization sorts order-insensitive security options, capabilities, mount
options, mounts, and devices. The original runtime input is not mutated.

Digest identity is evidence about bytes and planner version only. It is not a
security approval, a proof of policy correctness, or a claim that SELinux has
accepted the policy.

## Alternatives considered

### Hash raw adapter output

Rejected because Docker and OCI adapters may serialize equivalent inputs in
different orders and use different field names.

### Use only a policy digest

Rejected because reconciliation needs to distinguish changed input that happens
to produce the same current findings from unchanged input.

### Use timestamps or random IDs

Rejected because they make replay, deduplication, and audit comparison
non-deterministic.

## Consequences

Positive:

- equivalent adapter inputs produce the same plan identity;
- future queues can deduplicate by stable identity;
- audit and diff records can link a plan to its exact normalized input.

Limitations:

- digest stability depends on the named canonicalization version;
- a digest does not validate CIL or load SELinux state;
- persistence, retention, and rollback semantics remain future work.

## Verification

Unit tests cover reordered equivalent inputs, changed security input, and
cross-adapter OCI/Docker plan equivalence. No live Docker, SELinux, or
cross-process persistence qualification is claimed.
