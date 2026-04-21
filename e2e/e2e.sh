
#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
BIN_PATH="$ROOT_DIR/e2e/codespacegen"
IMAGE_CONFIG="$ROOT_DIR/e2e/codespacegen.json"
SNAPSHOT_DIR="$ROOT_DIR/e2e/snapshots"

UPDATE=false
for arg in "$@"; do
	case "$arg" in
		--update|-update) UPDATE=true ;;
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

if [[ ! -f "$IMAGE_CONFIG" ]]; then
	echo "[e2e] image config not found: $IMAGE_CONFIG"
	exit 1
fi

WORK_DIR="$(mktemp -d)"
trap 'rm -rf "$WORK_DIR"' EXIT

failures=0

for snapshot_case_dir in "$SNAPSHOT_DIR"/.devcontainer-*; do
	suffix="${snapshot_case_dir##*.devcontainer-}"

	case "$suffix" in
		go) lang="go" ;;
		python) lang="python" ;;
		rust) lang="rust" ;;
		moonbit) lang="moonbit" ;;
		biome) lang="node:biome" ;;
		eslint) lang="node:eslint" ;;
		react) lang="node:react" ;;
		*)
			echo "[e2e] unsupported snapshot case: $suffix"
			failures=$((failures + 1))
			continue
			;;
	esac

	port_args=()
	if grep -q "^[[:space:]]*ports:" "$snapshot_case_dir/docker-compose.yaml"; then
		port_args=("-port" "3000")
	fi

	if [[ "$UPDATE" == "true" ]]; then
		out_dir="$snapshot_case_dir"
		rm -f "$out_dir/Dockerfile" "$out_dir/devcontainer.json" "$out_dir/docker-compose.yaml"
		echo "[e2e] updating snapshot: $suffix"
	else
		out_dir="$WORK_DIR/.devcontainer-$suffix"
		mkdir -p "$out_dir"
	fi

	"$BIN_PATH" \
		-name sample \
		-language "$lang" \
		-service app \
		-workspace-folder /workspace \
		-compose-file docker-compose.yaml \
		-image-config "$IMAGE_CONFIG" \
		-output "$out_dir" \
		"${port_args[@]}" \
		</dev/null >/dev/null

	if [[ "$UPDATE" == "false" ]]; then
		for target in Dockerfile devcontainer.json docker-compose.yaml; do
			if ! diff -u "$snapshot_case_dir/$target" "$out_dir/$target" >/dev/null; then
				echo "[e2e] mismatch: $suffix/$target"
				diff -u "$snapshot_case_dir/$target" "$out_dir/$target" || true
				failures=$((failures + 1))
			fi
		done
	fi
done

if [[ "$UPDATE" == "true" ]]; then
	echo "[e2e] snapshots updated"
	exit 0
fi

if [[ $failures -gt 0 ]]; then
	echo "[e2e] failed with $failures mismatch(es)"
	exit 1
fi

echo "[e2e] all snapshot comparisons passed"
