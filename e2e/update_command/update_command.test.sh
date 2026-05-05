#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
BIN_PATH="$SCRIPT_DIR/codespacegen"

for arg in "$@"; do
	case "$arg" in
	--update | -update)
		# No snapshots for this scenario; keep the flag accepted for `make e2e UPD=--update`.
		:
		;;
	*)
		echo "[e2e] unknown option: $arg"
		echo "Usage: $0 [--update]"
		exit 1
		;;
	esac
done

if [[ ! -x "$BIN_PATH" ]]; then
	echo "[e2e] binary not found: $BIN_PATH"
	echo "[e2e] run 'make e2e' from repository root"
	exit 1
fi

STDOUT_PATH="$(mktemp)"
STDERR_PATH="$(mktemp)"
trap 'rm -f "$STDOUT_PATH" "$STDERR_PATH"' EXIT

set +e
LANG=C "$BIN_PATH" update >"$STDOUT_PATH" 2>"$STDERR_PATH"
status=$?
set -e

if [[ $status -eq 0 ]]; then
	echo "[e2e] expected update command to fail for invalid embedded version"
	cat "$STDOUT_PATH"
	cat "$STDERR_PATH"
	exit 1
fi

if ! grep -Fq "Short version cannot contain" "$STDERR_PATH"; then
	echo "[e2e] expected semver parse error in stderr"
	cat "$STDERR_PATH"
	exit 1
fi

echo "[e2e] update command error path passed"
