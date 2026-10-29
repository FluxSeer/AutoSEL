#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$repo_root"

compiler="${AUTOSEL_SELINUX_COMPILER:-secilc}"
result_file="$(mktemp "${TMPDIR:-/tmp}/autosel-qualification.XXXXXX.json")"
malformed_file="$(mktemp "${TMPDIR:-/tmp}/autosel-malformed.XXXXXX.cil")"
trap 'rm -f "$result_file" "$malformed_file"' EXIT

# The generated artifact is intentionally declaration-only. Validate it
# structurally through the normal CLI, then compile it with a complete,
# repository-owned policy base fixture.
GOWORK=off go run ./cmd/autosel validate \
  --id demo-container \
  --input testdata/oci/config.json >"$result_file"

if ! grep -q '"staticValid": true' "$result_file" || \
   ! grep -q '"compilerValidated": false' "$result_file"; then
  echo "static qualification did not report the expected validation boundary" >&2
  cat "$result_file" >&2
  exit 1
fi

base_policy="$repo_root/testdata/selinux/minimal-policy.cil"
candidate_policy="$repo_root/internal/policy/testdata/candidate.cil"

if ! AUTOSEL_SECILC="$compiler" "$repo_root/scripts/selinux/qualification-compiler.sh" \
  "$base_policy" "$candidate_policy" >/dev/null; then
  echo "compiler rejected the repository policy fixture and candidate" >&2
  exit 1
fi

printf '(block malformed\n' >"$malformed_file"
if AUTOSEL_SECILC="$compiler" "$repo_root/scripts/selinux/qualification-compiler.sh" \
  "$base_policy" "$malformed_file" >/dev/null 2>&1; then
  echo "compiler accepted malformed CIL" >&2
  exit 1
fi

echo "SELinux compiler qualification passed"
