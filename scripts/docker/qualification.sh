#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$repo_root"

container_id="${AUTOSEL_DOCKER_CONTAINER_ID:-}"
if [[ -z "$container_id" ]]; then
  echo "AUTOSEL_DOCKER_CONTAINER_ID is required" >&2
  exit 2
fi
if ! command -v docker >/dev/null 2>&1; then
  echo "Docker CLI is not available" >&2
  exit 127
fi

result_file="$(mktemp "${TMPDIR:-/tmp}/autosel-docker-qualification.XXXXXX.json")"
daemon_version_file="$(mktemp "${TMPDIR:-/tmp}/autosel-docker-daemon.XXXXXX.txt")"
trap 'rm -f "$result_file" "$daemon_version_file"' EXIT

# These are read-only Docker CLI operations. The harness deliberately does
# not create, start, stop, remove, exec into, or mutate a container.
docker version --format '{{.Server.Version}}' >"$daemon_version_file"
docker inspect --type container "$container_id" >/dev/null

GOWORK=off go run ./cmd/autosel plan \
  --container "$container_id" >"$result_file"

jq -e '
  (.status == "ready" or .status == "review-required") and
  (.inputDigest | startswith("sha256:")) and
  (.policyDigest | startswith("sha256:"))
' "$result_file" >/dev/null

daemon_version="$(tr -d '\n' <"$daemon_version_file")"
jq -c --arg daemon_version "$daemon_version" \
  '{status, inputDigest, policyDigest, findingCodes: [.findings[]?.code], dockerDaemonVersion: $daemon_version}' \
  "$result_file"
