# Changelog

## Unreleased

### Added

- read-only OCI security-spec normalization;
- deterministic policy planner and `autosel plan` command;
- experimental read-only Docker `ContainerInspect` adapter;
- deterministic planner input and policy digests;
- cross-adapter planner consistency tests;
- typed Policy IR, declaration-only candidate CIL, and static validation;
- optional external SELinux compiler adapter boundary;
- explicit `direct` and `secilc` compiler profiles with temporary output
  isolation;
- Ubuntu 24.04 Linux compiler qualification harness and support matrix;
- minimal complete CIL policy base fixture for candidate compiler qualification;
- manual Fedora qualification workflow and guarded Proxmox ephemeral-runner
  lifecycle helpers;
- first successful Fedora 44 manual compiler qualification evidence;
- three-run Fedora ephemeral qualification repeatability evidence;
- versioned validation evidence output and typed policy-rule provenance/approval
  contract;
- repository engineering, security, contribution, and capability-boundary
  documentation;
- unit coverage for high-risk capabilities, mounts, devices, and invalid input.

### Changed

- corrected the Go module identity to `github.com/FluxSeer/AutoSEL`;
- disabled the legacy mutating root entrypoint while the safe execution
  contract is being designed.
- moved the legacy mutation implementation under `internal/legacy/` during
  migration.
