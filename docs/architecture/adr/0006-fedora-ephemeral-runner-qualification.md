# ADR 0006: manual Fedora qualification on ephemeral Proxmox runners

## Status

Accepted for the first Fedora qualification slice.

## Context

Ubuntu 24.04 GitHub-hosted CI now proves the repository's static checks and
minimal-policy `secilc` compilation. Fedora-specific evidence is still needed,
but a persistent self-hosted runner would retain state from previous jobs and
would expand the trust boundary for workflow code.

## Decision

Add a manually triggered workflow,
`.github/workflows/fedora-selinux-qualification.yml`, that targets only a
runner with the `fedora-selinux-ephemeral` label. The runner must be a clean
Fedora VM cloned from the `fedora-selinux-base` template and registered with
GitHub using the `--ephemeral` runner option. The VM is represented
operationally by the `fedora-selinux-autosel-lab` clone name, but the base
template itself is never registered as a runner. The workflow uses the
`fedora-selinux-qualification` environment, which is expected to require
maintainer approval before a self-hosted job is allocated.

The job installs Fedora's `secilc` package, records non-secret environment
metadata, and runs the existing read-only qualification harness. It does not
load policy, call `semodule`, inspect or mutate Docker containers, or become a
pull-request required check.

The Proxmox helper scripts clone only a VM marked as a template and refuse to
overwrite an existing VM. Cleanup requires an exact VMID confirmation and both
`autosel-ephemeral` and `autosel-managed` tags before stopping and destroying
the VM. Configuration snapshots are retained in a restricted state directory.

## Consequences

- Fedora compiler evidence is available without changing the default PR path.
- A runner registration token remains an out-of-band secret and is never
  stored in the repository or emitted by the helper scripts.
- VM lifecycle automation still needs an operator or external controller to
  invoke provisioning and cleanup; the workflow cannot safely destroy its own
  runner.
- Fedora runtime enforcement, Docker integration, policy loading, and kernel
  behavior remain unqualified.

## Verification

The workflow must be run manually on a clean Fedora runner and its logs must
record the Fedora release, kernel, and `secilc` package version. The repository
checks remain the local Go test/vet/build suite plus the Ubuntu qualification
job. Proxmox commands are not executed by repository CI.
