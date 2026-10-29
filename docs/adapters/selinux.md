# SELinux validation adapter

AutoSEL currently separates candidate CIL generation from SELinux state
changes. The `internal/selinux` adapter can write a temporary candidate CIL
file and invoke a configured compiler command, but it never calls `semodule`
or installs a policy module.

```bash
GOWORK=off go run ./cmd/autosel validate \
  --id <container-id> \
  --input config.json \
  --compiler <compiler-binary> \
  --compiler-profile <direct|secilc>
```

The `direct` profile receives one temporary CIL path as its positional
argument. The `secilc` profile passes an explicit temporary output path to
`secilc`. Compiler flags and policy-store context are still not inferred from
the host; select a profile deliberately and qualify it on the target
distribution and policy base.

The current generated artifact is declaration-only, so it is intentionally
not a complete standalone SELinux policy. The normal `validate` command
therefore performs static validation unless a caller supplies a compiler
contract appropriate for a complete policy. The Linux qualification harness
compiles the generated candidate together with
`testdata/selinux/minimal-policy.cil`.

Without `--compiler`, `validate` still performs dependency-free structural CIL
validation and reports `compilerValidated: false`.

Every successful `validate` response also includes an `evidence` object with
schema version, planner/input/policy/artifact digests, compiler scope, and
explicit `policyLoaded: false` and `runtimeEnforcementTested: false` fields.
The evidence object is an audit boundary, not a runtime compatibility claim.

On Linux, `make selinux-qualification` runs a qualification harness using the
distribution's `secilc` binary. The harness statically validates the generated
candidate, compiles it with the checked-in minimal policy base, and checks that
malformed CIL is rejected. It only creates a temporary binary policy; it does
not load a module or qualify every distribution's policy base.

The generated module is declaration-only in the current slice. Findings do
not become SELinux enforcement rules until a reviewed Policy IR rule contract
exists.
