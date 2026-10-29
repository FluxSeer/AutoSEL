# ADR 0004: generate declaration-only CIL before enforcement rules

- Status: Accepted for the CIL planning slice
- Date: 2026-10-28
- Scope: candidate artifact and validation; no policy loading

## Context

The planner currently produces risk findings, not approved SELinux rules. The
old implementation translated findings directly into hand-built strings and
attempted to load them, without proving CIL syntax, type identity, or the
meaning of deny semantics on the target policy base.

Turning a review finding into enforcement is a separate security decision and
must not happen implicitly.

## Decision

Introduce the following boundary:

```text
PolicyPlan -> PolicyIR -> declaration-only CIL candidate -> static validation
                                               \-> optional external compiler
```

The current generator emits a valid-looking module scaffold containing a
module block and a type declaration. It rejects non-empty enforcement rules
until a reviewed rule contract exists. `autosel validate` never loads a module
and reports whether an external compiler was actually invoked.

The compiler adapter uses a temporary file and a command boundary. It does not
guess distribution-specific flags and has no policy-store mutation authority.

## Alternatives considered

### Convert every finding directly into a deny rule

Rejected because a finding is a review prompt, not an approved enforcement
decision, and CIL deny semantics require policy-base-specific qualification.

### Load the generated module to validate it

Rejected because loading changes host SELinux state and violates the read-only
planning contract.

### Pretend static checks are compiler validation

Rejected because balanced parentheses and a module shape do not prove policy
base compatibility or compiler acceptance. The result remains explicitly
separated as `staticValid` versus `compilerValidated`.

## Consequences

Positive:

- policy generation and policy loading remain separate;
- candidate artifacts are deterministic and golden-testable;
- external compiler execution is bounded and injectable in tests;
- incomplete enforcement semantics cannot silently ship as runtime policy.

Limitations:

- the current CIL artifact is not enforcement-ready;
- compiler qualification requires a Linux SELinux environment;
- no module is installed and no container receives a new label.

## Verification

Golden tests verify candidate bytes and digest shape. Structural validator
tests cover malformed content. Fake-runner tests verify temporary-file
compiler invocation and failure propagation. No live SELinux compiler result is
claimed on the current development host.
