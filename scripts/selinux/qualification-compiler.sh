#!/usr/bin/env bash
set -euo pipefail

if [[ "$#" -ne 1 ]]; then
  echo "usage: qualification-compiler.sh <candidate.cil>" >&2
  exit 2
fi

candidate="$1"
if [[ ! -f "$candidate" ]]; then
  echo "candidate CIL does not exist: $candidate" >&2
  exit 2
fi

compiler="${AUTOSEL_SECILC:-secilc}"
if ! command -v "$compiler" >/dev/null 2>&1; then
  echo "SELinux CIL compiler is not available: $compiler" >&2
  exit 127
fi

workspace="$(mktemp -d "${TMPDIR:-/tmp}/autosel-secilc.XXXXXX")"
trap 'rm -rf "$workspace"' EXIT

# Compile into a temporary binary policy only. This wrapper never invokes
# semodule or otherwise changes the host SELinux policy store.
exec "$compiler" -o "$workspace/policy.kern" "$candidate"
