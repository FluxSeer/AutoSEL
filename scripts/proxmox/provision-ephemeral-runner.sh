#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat >&2 <<'EOF'
usage: provision-ephemeral-runner.sh \
  --template-vmid <id> --vmid <id> --name <name> --storage <storage> \
  [--cores <count>] [--memory <MiB>] [--bridge <name>] [--state-dir <path>]

Clone a tagged ephemeral runner VM from a Proxmox template and start it.
This script must run on the target Proxmox node and does not register a
GitHub runner or handle a runner token.
EOF
}

template_vmid=""
vmid=""
vm_name=""
storage=""
cores="2"
memory="4096"
bridge=""
state_dir="${AUTOSEL_PROXMOX_STATE_DIR:-/var/lib/autosel-runner}"

while [[ "$#" -gt 0 ]]; do
  case "$1" in
    --template-vmid)
      [[ "$#" -ge 2 ]] || { usage; exit 2; }
      template_vmid="$2"
      shift 2
      ;;
    --vmid)
      [[ "$#" -ge 2 ]] || { usage; exit 2; }
      vmid="$2"
      shift 2
      ;;
    --name)
      [[ "$#" -ge 2 ]] || { usage; exit 2; }
      vm_name="$2"
      shift 2
      ;;
    --storage)
      [[ "$#" -ge 2 ]] || { usage; exit 2; }
      storage="$2"
      shift 2
      ;;
    --cores)
      [[ "$#" -ge 2 ]] || { usage; exit 2; }
      cores="$2"
      shift 2
      ;;
    --memory)
      [[ "$#" -ge 2 ]] || { usage; exit 2; }
      memory="$2"
      shift 2
      ;;
    --bridge)
      [[ "$#" -ge 2 ]] || { usage; exit 2; }
      bridge="$2"
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

for value in "$template_vmid" "$vmid" "$cores" "$memory"; do
  if [[ ! "$value" =~ ^[1-9][0-9]*$ ]]; then
    echo "expected a positive integer, got: $value" >&2
    exit 2
  fi
done
if [[ -z "$vm_name" || -z "$storage" ]]; then
  echo "--name and --storage are required" >&2
  usage
  exit 2
fi
if ! command -v qm >/dev/null 2>&1; then
  echo "Proxmox qm command is not available" >&2
  exit 127
fi

template_config="$(qm config "$template_vmid")" || {
  echo "template VM does not exist: $template_vmid" >&2
  exit 1
}
if ! grep -q '^template: 1$' <<<"$template_config"; then
  echo "refusing to clone a VM that is not marked as a template: $template_vmid" >&2
  exit 1
fi
if qm status "$vmid" >/dev/null 2>&1; then
  echo "refusing to overwrite existing target VM: $vmid" >&2
  exit 1
fi

install -d -m 0700 "$state_dir"
manifest="$state_dir/${vmid}.manifest"
if [[ -e "$manifest" ]]; then
  echo "refusing to overwrite existing manifest: $manifest" >&2
  exit 1
fi

qm clone "$template_vmid" "$vmid" --full 1 --name "$vm_name" --storage "$storage"
qm set "$vmid" --cores "$cores" --memory "$memory" --onboot 0
if [[ -n "$bridge" ]]; then
  qm set "$vmid" --net0 "virtio,bridge=$bridge"
fi
qm set "$vmid" --tags "autosel-ephemeral;autosel-managed"

template_config_path="$state_dir/${vmid}.template.config"
printf '%s\n' "$template_config" >"$template_config_path"
{
  printf 'vmid=%s\n' "$vmid"
  printf 'template_vmid=%s\n' "$template_vmid"
  printf 'name=%s\n' "$vm_name"
  printf 'created_at=%s\n' "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
  printf 'marker=autosel-ephemeral\n'
  printf 'config_snapshot=%s\n' "$template_config_path"
} >"$manifest"
chmod 0600 "$manifest" "$template_config_path"

qm start "$vmid"
echo "ephemeral runner VM started: vmid=$vmid name=$vm_name"
echo "register the runner out-of-band with --ephemeral and label fedora-selinux-ephemeral"
echo "manifest: $manifest"
