#!/usr/bin/env bash
set -euo pipefail

# args:
#   -t|--tag <tag>  : tag to create (default: current git tag or "dev")
#   -p|--push        : push the created tag to the remote repository
#   -h|--help        : show this help message

# ---------------------------------------------------------------------------
# Constants
# ---------------------------------------------------------------------------
readonly SCRIPT_NAME="${0##*/}"

# ---------------------------------------------------------------------------
# Functions
# ---------------------------------------------------------------------------

print_help() {
	cat <<EOF
Usage: ${SCRIPT_NAME} [options]
Options:
  -t, --tag <tag>   Tag to create (default: latest git tag on current commit)
  -p, --push        Push the created tag to the remote repository
  -h, --help        Show this help message
EOF
}

die() {
	echo "Error: $*" >&2
	exit 1
}

# Strip leading "v" prefix from a version string (e.g. "v1.2.3" -> "1.2.3")
strip_v_prefix() {
	echo "${1#v}"
}

parse_args() {
	TAG=""
	PUSH=false

	while [[ $# -gt 0 ]]; do
		case "$1" in
		-t | --tag)
			[[ -z "${2:-}" ]] && die "--tag requires an argument"
			TAG="$2"
			shift 2
			;;
		-p | --push)
			PUSH=true
			shift
			;;
		-h | --help)
			print_help
			exit 0
			;;
		*)
			echo "Unknown option: $1" >&2
			print_help
			exit 1
			;;
		esac
	done
}

resolve_tag() {
	local tag
	tag="$(strip_v_prefix "${TAG}")"

	if [[ -z "$tag" ]]; then
		tag="$(git describe --tags 2>/dev/null)" ||
			die "no tag found. Use -t to specify a tag."
		tag="$(strip_v_prefix "$tag")"
	fi

	echo "$tag"
}

create_tag() {
	local tag="$1"
	local full_tag="v${tag}"

	[[ -n "$(git tag -l "$full_tag")" ]] &&
		return 1

	git tag "$full_tag"
	echo "Created tag: ${full_tag}"
	return 0
}

remote_has_tag() {
	local full_tag="$1"
	git ls-remote --tags origin "refs/tags/${full_tag}" | grep -q .
}

push_tag() {
	local full_tag="$1"
	git push origin "$full_tag"
	echo "Pushed tag: ${full_tag} to origin"
}

# ---------------------------------------------------------------------------
# Main
# ---------------------------------------------------------------------------

main() {
	parse_args "$@"

	local tag
	tag="$(resolve_tag)"

	local full_tag="v${tag}"
	if create_tag "$tag"; then
		:
	else
		if [[ "$PUSH" == true ]]; then
			echo "Tag already exists locally: ${full_tag}"
		else
			die "tag ${full_tag} already exists"
		fi
	fi

	if [[ "$PUSH" == true ]]; then
		if remote_has_tag "$full_tag"; then
			die "tag ${full_tag} already exists on origin"
		fi
		push_tag "$full_tag"
	fi
}

main "$@"
