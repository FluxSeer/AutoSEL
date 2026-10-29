# ADR 0007: version validation evidence and policy rule decisions

## Status

Accepted for the read-only planning and validation slice.

## Context

AutoSEL already produces deterministic plan, policy, and candidate-artifact
digests, but callers previously had to infer the validation boundary from
several top-level booleans. Policy rules also had CIL fields without an
explicit provenance or operator-decision envelope.

## Decision

`autosel validate` emits a typed `evidence` object using the
`autosel-validation-evidence/v1` schema. It records planner and artifact
identity, static/compiler results, and explicit negative claims for policy
loading and runtime enforcement. The compiler evidence only describes the
temporary compiler workspace; it does not imply policy-store installation or
runtime enforcement.

Policy IR uses `autosel-policy-rule/v1`. Every future rule must carry typed
provenance and an approval status. Approved rules require an approver and
decision ID; pending rules remain review material; rejected rules cannot enter
the IR. The current CIL generator still rejects every non-empty rule set, so
adding the contract does not enable enforcement.

## Consequences

- Validation output can be stored and compared without interpreting prose logs.
- Plan, policy, and artifact identities remain deterministic and auditable.
- A finding cannot become a policy rule without provenance and an explicit
  operator decision record.
- Host environment details such as distro, kernel, and compiler package remain
  qualification-run evidence until a separate, privacy-reviewed environment
  contract is added.

## Verification

Unit tests cover evidence boundary defaults and rule-contract validation. The
existing Go test, vet, build, Ubuntu compiler qualification, and Fedora manual
qualification workflows remain unchanged in their read-only behavior.
