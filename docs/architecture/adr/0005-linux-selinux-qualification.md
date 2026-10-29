# ADR 0005: Qualify candidate CIL with a Linux compiler

## Status

Accepted for the first Linux qualification slice.

## Context

AutoSEL's dependency-free CIL check can detect basic structural mistakes, but
it cannot establish that a distribution's SELinux CIL compiler accepts the
candidate or that the candidate is compatible with the installed policy base.
The development host is not assumed to have SELinux tooling, so a successful
local unit test must not be presented as compiler qualification.

## Decision

Add a Linux CI qualification job that:

1. installs the distribution-provided `secilc` compiler;
2. runs the normal read-only `autosel validate` path against the OCI fixture;
3. checks the declaration-only candidate with static validation;
4. compiles the candidate together with the repository's minimal complete
   policy fixture into a temporary binary policy;
5. verifies that malformed CIL is rejected; and
6. never invokes `semodule`, writes the system policy store, or changes host
   SELinux state.

The repository wrapper remains available for negative compiler checks. The
profile is intentionally a qualification contract, not a production runtime
profile. Distribution-specific policy-base and runtime support remain
unqualified until the support matrix records evidence.

## Consequences

The compiler job has evidence that the candidate can be parsed in the context
of the checked-in minimal policy base by an actual Linux compiler. The
declaration-only artifact is not passed to `secilc` as a complete policy on
its own. The project still does not claim that a candidate is loadable on
every SELinux distribution, compatible with a host's installed policy base, or
enforcing any finding.
