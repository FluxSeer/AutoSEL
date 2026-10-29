# AutoSEL documentation

## Start here

- [Architecture overview](architecture/overview.md)
- [Policy planning boundary ADR](architecture/adr/0001-policy-planning-boundary.md)
- [Docker inspection adapter](adapters/docker.md)
- [Docker inspection ADR](architecture/adr/0002-docker-inspection-adapter.md)
- [Deterministic policy identity ADR](architecture/adr/0003-deterministic-policy-identity.md)
- [Candidate CIL validation ADR](architecture/adr/0004-candidate-cil-validation.md)
- [Validation evidence and rule contract ADR](architecture/adr/0007-validation-evidence-and-rule-contract.md)
- [SELinux validation adapter](adapters/selinux.md)
- [Support and qualification matrix](support-matrix.md)
- [Capability maturity](capability-maturity.md)
- [Security policy](../SECURITY.md)
- [Roadmap](../ROADMAP.md)

## Documentation rules

AutoSEL uses the following public capability labels:

| Label | Meaning |
| --- | --- |
| Supported | Tested public behavior covered by the current read-only contract. |
| Experimental | Implemented behavior requiring explicit opt-in or qualification. |
| Planned | Documented direction without supported implementation. |
| Not supported | Explicitly outside the current product boundary. |

Documentation must distinguish observed container input, planner inference,
operator review, and any future runtime action.
