# Security policy

## Reporting a vulnerability

Please report security vulnerabilities privately through
[GitHub Security Advisories](https://github.com/FluxSeer/AutoSEL/security/advisories/new).
Do not open a public issue for a suspected vulnerability.

Include the affected version or commit, component, reproduction steps, impact,
and logs or fixtures with secrets removed.

## Runtime security boundary

The current supported path is read-only. It does not connect to the Docker
socket, execute `semodule`, stop containers, remove containers, modify
volumes, or rebuild images.

Any future application path must be separately authorized and must define:

- least-privilege Docker and host permissions;
- policy and target validation;
- deterministic execution identity and idempotency;
- timeout, retry, and uncertain-outcome handling;
- complete configuration and data-preservation checks;
- health verification, rollback, and audit records.

AutoSEL must not persist secrets, tokens, authorization headers, or complete
container environment variables in plans, logs, fixtures, or audit artifacts.

Self-hosted runner boundary: Fedora qualification is manual-only and must use
an isolated, single-job runner. Runner registration tokens are supplied
out-of-band, never committed to the repository, and must not appear in
workflow logs. The base Fedora template is not itself a runner; each job must
use a fresh VM clone and clean it up after the job. The GitHub environment used
by the workflow should require maintainer approval before allocating the
self-hosted runner.

## Support boundary

AutoSEL is pre-1.0 and experimental. The supported behavior is listed in
[capability maturity](docs/capability-maturity.md); unqualified Docker,
SELinux, distribution, and kernel combinations are not production-support
claims.
