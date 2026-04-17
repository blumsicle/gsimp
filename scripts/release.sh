#!/usr/bin/env bash

set -euo pipefail

usage() {
	echo "usage: $0 [--no-ask] <version>" >&2
	exit 1
}

skip_prompt=false

if [[ $# -gt 0 && "$1" == "--no-ask" ]]; then
	skip_prompt=true
	shift
fi

if [[ $# -ne 1 ]]; then
	usage
fi

version="$1"

if [[ "${version}" != v* ]]; then
	echo "release version must start with 'v' (example: v1.2.3)" >&2
	exit 1
fi

if git rev-parse --verify --quiet "${version}" >/dev/null; then
	echo "tag already exists: ${version}" >&2
	exit 1
fi

if [[ ! -f RELEASE_NOTES.md ]]; then
	echo "RELEASE_NOTES.md must exist before creating a release" >&2
	exit 1
fi

if [[ "${skip_prompt}" != true ]]; then
	printf "Have you updated RELEASE_NOTES.md for %s? [y/N] " "${version}" >&2
	read -r answer
	case "${answer}" in
		y|Y|yes|YES)
			;;
		*)
			echo "aborting release; update RELEASE_NOTES.md first" >&2
			exit 1
			;;
	esac
fi

status_lines="$(git status --short)"
if [[ -z "${status_lines}" ]]; then
	echo "RELEASE_NOTES.md must be updated before creating a release" >&2
	exit 1
fi

invalid_changes="$(printf '%s\n' "${status_lines}" | grep -v ' RELEASE_NOTES\.md$' || true)"
if [[ -n "${invalid_changes}" ]]; then
	echo "only RELEASE_NOTES.md may be changed when running the release script" >&2
	echo "${invalid_changes}" >&2
	exit 1
fi

notes_tracked=true
if ! git ls-files --error-unmatch RELEASE_NOTES.md >/dev/null 2>&1; then
	notes_tracked=false
fi

if [[ "${notes_tracked}" == true ]] && git diff --quiet -- RELEASE_NOTES.md && git diff --cached --quiet -- RELEASE_NOTES.md; then
	echo "RELEASE_NOTES.md must contain changes before creating a release" >&2
	exit 1
fi

task fmt
task check

git add RELEASE_NOTES.md
git commit -m "Prepare release ${version}"
git tag -a "${version}" -m "Release ${version}"
git push origin HEAD
git push origin "${version}"
