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
3. invokes a repository wrapper that compiles into a temporary binary policy;
4. verifies that malformed CIL is rejected; and
5. never invokes `semodule`, writes the system policy store, or changes host
   SELinux state.

The wrapper is intentionally a qualification harness, not a production
runtime profile. Distribution-specific policy-base and runtime support remain
unqualified until a support matrix is published.

## Consequences

`compilerValidated: true` will have evidence from an actual Linux compiler in
CI, while `staticValid: true` remains the dependency-free local signal. The
project still does not claim that a candidate is loadable on every SELinux
distribution, nor does it claim that the declaration-only artifact enforces
any finding.
