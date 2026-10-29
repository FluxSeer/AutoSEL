# ADR 0008: qualify Docker inspection without container mutation

## Status

Accepted for the first Docker qualification slice.

## Context

The Docker adapter is intentionally read-only, but unit tests alone do not
prove that the Docker SDK can inspect a live daemon or that the normalized
result reaches the deterministic planner. Creating a test container inside CI
would violate the current product boundary because container mutation has no
reviewed execution contract.

## Decision

Add a manually triggered workflow,
`.github/workflows/docker-readonly-qualification.yml`, that targets an
isolated runner labelled `docker-readonly-lab`. The runner must already have a
Docker daemon and a pre-existing test container. The workflow requires an
environment approval and a caller-supplied container ID.

The harness performs only Docker daemon version, container inspect, and
`autosel plan --container` operations. It does not create, start, stop, remove,
exec into, export, or modify a container. It emits a redacted summary with
plan status, digests, finding codes, and daemon version rather than logging the
full inspect response.

## Consequences

- Live Docker SDK and planner wiring can be qualified without expanding the
  mutation contract.
- The lab operator remains responsible for supplying a safe, disposable
  pre-existing container and an isolated runner.
- Docker Engine, rootless Docker, Docker Desktop, SELinux labels, and runtime
  compatibility remain unqualified until evidence is collected per environment.
- The workflow is manual-only and is not a pull-request required check.

## Verification

The script has shell syntax and output-shape checks. A live qualification is
considered successful only when the manual workflow completes against an
existing container and the runner/environment lifecycle is recorded.
