# Capability maturity

| Capability | Status | Current contract |
| --- | --- | --- |
| OCI config normalization | Supported | Typed subset used by `autosel plan`. |
| Deterministic policy findings | Supported | Pure planner with unit coverage. |
| Deterministic plan identity | Supported | Versioned input and policy SHA-256 digests. |
| Policy IR | Supported | Typed review-aware intermediate representation with versioned rule provenance and approval fields; enforcement remains disabled. |
| Validation evidence contract | Experimental | Versioned `validate` evidence with digests and explicit policy-loading/runtime-negative claims. |
| Candidate CIL generation | Experimental | Declaration-only deterministic scaffold; no enforcement rules. |
| Static CIL validation | Experimental | Structural validation with unit coverage. |
| SELinux compiler adapter | Experimental | Optional command boundary with fake-runner tests and a Linux qualification harness. |
| Docker API adapter | Experimental | Read-only `ContainerInspect` normalization with fake-client unit tests; live daemon qualification pending. |
| Docker read-only live qualification | Experimental | Manual, approval-gated harness against a pre-existing lab container; no container mutation. |
| Event reconciliation | Planned | Must handle duplicates, disconnects, retries, and resync. |
| CIL compiler validation | Experimental | Ubuntu 24.04 CI and three Fedora 44 ephemeral runs compile the repository fixture without loading a module; broader distro matrix remains pending. |
| SELinux module loading | Planned | Requires explicit authorization and audit. |
| Safe container replacement | Not supported | No data-preserving, rollback-capable contract yet. |

An installed package, exported function, or future interface must not be
interpreted as supported runtime behavior unless it appears in this table and
has matching verification evidence.

The original mutation implementation is retained only under
`internal/legacy/` during migration and is not a supported capability.
