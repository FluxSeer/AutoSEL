#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$repo_root"

compiler="${AUTOSEL_SELINUX_COMPILER:-secilc}"
profile="${AUTOSEL_SELINUX_COMPILER_PROFILE:-secilc}"
result_file="$(mktemp "${TMPDIR:-/tmp}/autosel-qualification.XXXXXX.json")"
malformed_file="$(mktemp "${TMPDIR:-/tmp}/autosel-malformed.XXXXXX.cil")"
trap 'rm -f "$result_file" "$malformed_file"' EXIT

GOWORK=off go run ./cmd/autosel validate \
  --id demo-container \
  --input testdata/oci/config.json \
  --compiler "$compiler" \
  --compiler-profile "$profile" >"$result_file"

if ! grep -q '"compilerValidated": true' "$result_file"; then
  echo "compiler qualification did not report compilerValidated=true" >&2
  cat "$result_file" >&2
  exit 1
fi

printf '(block malformed\n' >"$malformed_file"
if "$repo_root/scripts/selinux/qualification-compiler.sh" "$malformed_file" >/dev/null 2>&1; then
  echo "compiler accepted malformed CIL" >&2
  exit 1
fi

echo "SELinux compiler qualification passed"
