# AutoSEL roadmap

## v0.1 Read-only policy planning

Current slice:

- typed OCI security-spec normalization;
- deterministic policy findings;
- CLI `autosel plan`;
- experimental read-only Docker inspection adapter;
- typed policy IR and declaration-only candidate CIL;
- static CIL validation and optional compiler boundary;
- unit tests and open-source repository baseline;
- no Docker or SELinux side effects.

Exit criteria:

- planner behavior is stable and documented;
- positive, negative, malformed, and ambiguous fixtures exist;
- generated policy identity is deterministic;
- Linux/SELinux compatibility test environment is defined.

## v0.2 Adapter and validation foundations

- Docker API inspection adapter (implemented experimentally; runtime qualification pending);
- Docker event collector with reconnect and periodic resync;
- CIL generation and compiler validation;
- reviewed enforcement-rule contract;
- policy diff and audit artifacts;
- no automatic container replacement by default.

## v0.3 Guarded application

- pre-create labeling path where possible;
- complete inspect-to-create snapshot preservation;
- explicit authorization and maintenance-window controls;
- health verification and rollback;
- qualified Linux integration tests.

## Explicit non-goals for the first releases

- generic Docker daemon replacement;
- Kubernetes operator or cluster-wide admission controller;
- arbitrary shell or host mutation;
- automatic deletion of containers or volumes;
- claims of complete least-privilege policy generation;
- support claims for distributions or kernels not tested in CI.
