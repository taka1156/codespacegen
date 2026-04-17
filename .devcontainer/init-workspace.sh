#!/usr/bin/env bash
set -euo pipefail

trap 'echo "[ERROR] init-workspace.sh failed at line ${LINENO}: ${BASH_COMMAND}" >&2' ERR

script_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd -- "${script_dir}/.." && pwd)"

"${repo_root}/scripts/install.sh"
"${repo_root}/scripts/setup.sh"
