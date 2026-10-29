# Support and qualification matrix

This matrix describes evidence boundaries, not a promise that AutoSEL can
load policy on a host. The current qualification job compiles the
declaration-only candidate together with a repository-owned minimal policy
base fixture and never changes SELinux state.

## CI runner policy

The source verification and SELinux compiler qualification jobs use the
GitHub-hosted `ubuntu-24.04` runner. The Proxmox Fedora
`fedora-selinux-base` template is not registered as a CI runner and is not
used by pull-request workflows. A Fedora-specific runtime qualification job
requires separate evidence and an explicitly controlled, ephemeral runner
lifecycle.

| Environment | Compiler validation | Runtime/enforcement | Status |
| --- | --- | --- | --- |
| Ubuntu 24.04 amd64, GitHub Actions `ubuntu-24.04`, Ubuntu `secilc` package | Candidate fixture compiled with the checked-in minimal policy base using `secilc` | Not tested | CI qualification target |
| Ubuntu 22.04 amd64 | Not run in CI | Not tested | Planned |
| Fedora / RHEL ephemeral runner | Manual workflow on an isolated Fedora VM with distro `secilc` | Not tested | Qualification ready; evidence pending |
| macOS development host | No native SELinux compiler expected | Not applicable | Static tests only |

## Current compiler profiles

The `validate` command supports two explicit profiles:

- `direct`: passes the temporary CIL path as the only compiler argument. This
  is retained for simple wrapper commands and backwards compatibility.
- `secilc`: the generic adapter invokes `secilc -o <temporary-output>
  <temporary-input>`. The qualification harness additionally supplies the
  checked-in minimal policy base because the generated candidate is
  declaration-only. The output is removed with the temporary workspace after
  validation.

Example:

```bash
GOWORK=off go run ./cmd/autosel validate \
  --id demo-container \
  --input testdata/oci/config.json \
  --compiler secilc \
  --compiler-profile secilc
```

The standalone command above is intended for a complete policy input. For the
repository candidate qualification, use `make selinux-qualification`, which
combines `testdata/selinux/minimal-policy.cil` with the candidate fixture.

## Fedora manual qualification

`.github/workflows/fedora-selinux-qualification.yml` is intentionally
`workflow_dispatch`-only and targets the `fedora-selinux-ephemeral` runner
label. It is not a pull-request required check and uses the
`fedora-selinux-qualification` environment, which should be configured with
required reviewers. The runner must be a clean VM clone of the
`fedora-selinux-base` template, registered with GitHub using `--ephemeral`,
and removed after the single job. The expected lab clone name is
`fedora-selinux-autosel-lab`.

The Proxmox helper scripts run on the Proxmox node. Provisioning starts a
tagged clone but does not register a GitHub runner or accept a runner token;
registration is an out-of-band operator step. Cleanup requires the exact VMID
confirmation and refuses untagged VMs. The scripts retain configuration
snapshots and never run in GitHub-hosted CI.

## Explicitly outside the matrix

- system policy store writes and `semodule` loading;
- kernel enforcement behavior and label transitions;
- Docker Engine, rootless Docker, and Docker Desktop integration;
- policy-base compatibility beyond the compiler fixture;
- generated enforcement rules from planner findings.
