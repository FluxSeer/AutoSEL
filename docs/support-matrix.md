# Support and qualification matrix

This matrix describes evidence boundaries, not a promise that AutoSEL can
load policy on a host. The current qualification job only compiles the
declaration-only candidate artifact and never changes SELinux state.

| Environment | Compiler validation | Runtime/enforcement | Status |
| --- | --- | --- | --- |
| Ubuntu 24.04 amd64, GitHub Actions `ubuntu-24.04`, Ubuntu `secilc` package | Candidate fixture compiled with the `secilc` profile | Not tested | CI qualification target |
| Ubuntu 22.04 amd64 | Not run in CI | Not tested | Planned |
| Fedora / RHEL | Not run in CI | Not tested | Planned |
| macOS development host | No native SELinux compiler expected | Not applicable | Static tests only |

## Current compiler profiles

The `validate` command supports two explicit profiles:

- `direct`: passes the temporary CIL path as the only compiler argument. This
  is retained for simple wrapper commands and backwards compatibility.
- `secilc`: invokes `secilc -o <temporary-output> <temporary-input>`. The
  output is removed with the temporary workspace after validation.

Example:

```bash
GOWORK=off go run ./cmd/autosel validate \
  --id demo-container \
  --input testdata/oci/config.json \
  --compiler secilc \
  --compiler-profile secilc
```

## Explicitly outside the matrix

- system policy store writes and `semodule` loading;
- kernel enforcement behavior and label transitions;
- Docker Engine, rootless Docker, and Docker Desktop integration;
- policy-base compatibility beyond the compiler fixture;
- generated enforcement rules from planner findings.
