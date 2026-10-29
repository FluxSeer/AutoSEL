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

Without `--compiler`, `validate` still performs dependency-free structural CIL
validation and reports `compilerValidated: false`.

On Linux, `make selinux-qualification` runs a qualification harness using the
distribution's `secilc` binary. The harness compiles only to a temporary binary
policy and also checks that malformed CIL is rejected. It does not load a
module or qualify every distribution's policy base.

The generated module is declaration-only in the current slice. Findings do
not become SELinux enforcement rules until a reviewed Policy IR rule contract
exists.
