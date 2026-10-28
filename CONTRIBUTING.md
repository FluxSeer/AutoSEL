# Contributing

AutoSEL is an experimental FluxSeer open-source project for safe, reviewable
container policy planning. Contributions should preserve the read-only-first
boundary.

## Principles

- Keep the policy core deterministic and independent of Docker and SELinux.
- Keep runtime-specific behavior behind adapters.
- Treat findings as evidence-backed review prompts, not automatic accusations.
- Do not add mutation without explicit authorization, idempotency, audit, and
  rollback semantics.
- Add tests for malformed input, missing data, duplicate events, retries, and
  partial failures where the change affects those paths.
- Document supported, experimental, planned, and unsupported behavior.

## Development

```bash
GOWORK=off go test ./...
GOWORK=off go vet ./...
GOWORK=off go build ./...
```

Use focused pull requests and include the verification summary in the PR
description. Do not claim Docker, SELinux, distribution, kernel, or runtime
compatibility without a matching test result.

## Security

Please report vulnerabilities privately using the process in [SECURITY.md](SECURITY.md).
Do not include secrets, tokens, full environment variables, or unredacted
container configuration in issues or fixtures.

## Developer Certificate of Origin

By submitting a contribution, you certify that you have the right to submit it
under the Apache License, Version 2.0.
