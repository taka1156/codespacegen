#!/usr/bin/env bash
set -euo pipefail

REPO="taka1156/codespacegen"
BINARY_NAME="codespacegen"
BINARY_SHORTCUT_NAME="csg"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

detect_os() {
	case "$(uname -s)" in
	Linux) echo "linux" ;;
	Darwin) echo "darwin" ;;
	*)
		echo "unsupported OS: $(uname -s)" >&2
		exit 1
		;;
	esac
}

detect_arch() {
	case "$(uname -m)" in
	x86_64 | amd64) echo "amd64" ;;
	arm64 | aarch64) echo "arm64" ;;
	*)
		echo "unsupported architecture: $(uname -m)" >&2
		exit 1
		;;
	esac
}

ensure_cmd() {
	if ! command -v "$1" >/dev/null 2>&1; then
		echo "required command not found: $1" >&2
		exit 1
	fi
}

create_symlink() {
	local target="$1"
	local link="$2"
	local use_sudo="${3:-false}"
	if [ -L "$link" ] || [ -e "$link" ]; then
		local current
		current="$(readlink "$link" 2>/dev/null || echo "(not a symlink)")"
		echo "warning: '$link' already exists (currently -> '$current')"
		printf "overwrite with '%s'? [y/N]: " "$target" >/dev/tty
		read -r answer
		[[ "$answer" =~ ^[yY] ]] || {
			echo "skipped: $link"
			return 0
		}
	fi
	if [[ "$use_sudo" == "true" ]]; then
		sudo ln -sf "$target" "$link"
	else
		ln -sf "$target" "$link"
	fi
}

main() {
	ensure_cmd curl
	ensure_cmd tar

	local os arch asset tmpdir download_url
	os="$(detect_os)"
	arch="$(detect_arch)"
	asset="${BINARY_NAME}_${os}_${arch}.tar.gz"
	download_url="https://github.com/${REPO}/releases/latest/download/${asset}"

	tmpdir="$(mktemp -d)"
	trap '[[ -n "${tmpdir:-}" ]] && rm -rf "$tmpdir"' EXIT

	echo "downloading ${asset} ..."
	curl -fsSL "$download_url" -o "$tmpdir/$asset"

	tar -xzf "$tmpdir/$asset" -C "$tmpdir"
	chmod +x "$tmpdir/$BINARY_NAME"

	mkdir -p "$INSTALL_DIR"
	if [ -w "$INSTALL_DIR" ]; then
		mv "$tmpdir/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
		create_symlink "$INSTALL_DIR/$BINARY_NAME" "$INSTALL_DIR/$BINARY_SHORTCUT_NAME"
	elif command -v sudo >/dev/null 2>&1; then
		sudo mv "$tmpdir/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
		create_symlink "$INSTALL_DIR/$BINARY_NAME" "$INSTALL_DIR/$BINARY_SHORTCUT_NAME" "true"
	else
		echo "no write permission for $INSTALL_DIR and sudo is not available" >&2
		exit 1
	fi

	echo "installed: $INSTALL_DIR/$BINARY_NAME"
	echo "run: \`${BINARY_NAME} -h\` or \`${BINARY_SHORTCUT_NAME} -h\`"
}

main "$@"
