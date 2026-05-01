#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
BIN_PATH="$SCRIPT_DIR/codespacegen"
SNAPSHOT_PATH="$SCRIPT_DIR/snapshots/codespacegen.json"
GENERATED_PATH="$SCRIPT_DIR/codespacegen.json"

UPDATE=false
for arg in "$@"; do
	case "$arg" in
	--update | -update) UPDATE=true ;;
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

trap 'rm -f "$GENERATED_PATH"' EXIT

# initコマンドはカレントディレクトリにファイルを生成するためSCRIPT_DIR内で実行
(cd "$SCRIPT_DIR" && "$BIN_PATH" init -output codespacegen.json)

if [[ "$UPDATE" == "true" ]]; then
	cp "$GENERATED_PATH" "$SNAPSHOT_PATH"
	echo "[e2e] snapshot updated"
	exit 0
fi

diff -u "$SNAPSHOT_PATH" "$GENERATED_PATH" || {
	echo "[e2e] codespacegen.json mismatch"
	exit 1
}

echo "[e2e] codespacegen.json snapshot comparison passed"
