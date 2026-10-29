#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat >&2 <<'EOF'
usage: cleanup-ephemeral-runner.sh --vmid <id> --confirm-destroy <same-id> \
  [--state-dir <path>]

Stop and destroy only an AutoSEL-managed ephemeral runner VM. The exact
VMID confirmation and autosel-ephemeral tag are both required.
EOF
}

vmid=""
confirmation=""
state_dir="${AUTOSEL_PROXMOX_STATE_DIR:-/var/lib/autosel-runner}"

while [[ "$#" -gt 0 ]]; do
  case "$1" in
    --vmid)
      [[ "$#" -ge 2 ]] || { usage; exit 2; }
      vmid="$2"
      shift 2
      ;;
    --confirm-destroy)
      [[ "$#" -ge 2 ]] || { usage; exit 2; }
      confirmation="$2"
      shift 2
      ;;
    --state-dir)
      [[ "$#" -ge 2 ]] || { usage; exit 2; }
      state_dir="$2"
      shift 2
      ;;
    --help|-h)
      usage >&1
      exit 0
      ;;
    *)
      echo "unknown argument: $1" >&2
      usage
      exit 2
      ;;
  esac
done

if [[ ! "$vmid" =~ ^[1-9][0-9]*$ || "$confirmation" != "$vmid" ]]; then
  echo "--vmid must be a positive integer and --confirm-destroy must match it exactly" >&2
  exit 2
fi
if ! command -v qm >/dev/null 2>&1; then
  echo "Proxmox qm command is not available" >&2
  exit 127
fi

config="$(qm config "$vmid")" || {
  echo "target VM does not exist: $vmid" >&2
  exit 1
}
if ! grep -Eq '^tags: .*autosel-ephemeral(;|$)' <<<"$config"; then
  echo "refusing to destroy VM without autosel-ephemeral tag: $vmid" >&2
  exit 1
fi
if ! grep -Eq '^tags: .*autosel-managed(;|$)' <<<"$config"; then
  echo "refusing to destroy VM without autosel-managed tag: $vmid" >&2
  exit 1
fi

install -d -m 0700 "$state_dir"
snapshot="$state_dir/${vmid}.cleanup.config"
printf '%s\n' "$config" >"$snapshot"
chmod 0600 "$snapshot"

if qm status "$vmid" | grep -q 'status: running'; then
  qm shutdown "$vmid" --timeout 60 || qm stop "$vmid"
fi
qm destroy "$vmid" --purge 1
echo "destroyed AutoSEL ephemeral runner VM: vmid=$vmid"
echo "configuration snapshot retained at: $snapshot"
