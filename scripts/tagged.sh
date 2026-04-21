#!/usr/bin/env bash
set -euo pipefail

# args:
#   -t|--tag <tag>  : tag to create (default: current git tag or "dev")
#   -p|--push        : push the created tag to the remote repository
#   -h|--help        : show this help message

TAG=""
PUSH=false

print_help() {
	cat <<EOF
Usage: $0 [options]
Options:
  -t, --tag <tag>   Tag to create (default: latest git tag on current commit)
  -p, --push        Push the created tag to the remote repository
  -h, --help        Show this help message
EOF
}

while [[ $# -gt 0 ]]; do
	case "$1" in
	-t | --tag)
		if [[ -z "${2:-}" ]]; then
			echo "Error: --tag requires an argument" >&2
			exit 1
		fi
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

TAG="${TAG#v}"

if [[ -z "$TAG" ]]; then
	TAG=$(git describe --tags 2>/dev/null)
	if [[ -z "$TAG" ]]; then
		echo "Error: no tag found. Use -t to specify a tag." >&2
		exit 1
	fi
	TAG="${TAG#v}"
fi

if [[ -n "$(git tag -l "v${TAG}")" ]]; then
	echo "Error: tag v${TAG} already exists" >&2
	exit 1
fi
git tag "v${TAG}"

if [[ "$PUSH" == true ]]; then
	git push origin "v${TAG}"
fi
