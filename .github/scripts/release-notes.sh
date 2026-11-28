#!/usr/bin/env bash
# Print the release notes of a version, taken from CHANGELOG.md.
# The notes are the lines under the heading "## [VERSION]", without the heading itself,
# up to the next heading or the link definitions at the end, without blank lines at both ends.
#
# Usage: .github/scripts/release-notes.sh VERSION [CHANGELOG]
#   VERSION     the version without the leading "v", such as 0.0.2
#   CHANGELOG   the path of the changelog (default: CHANGELOG.md)
set -euo pipefail

version="${1:?usage: release-notes.sh VERSION [CHANGELOG]}"
changelog="${2:-CHANGELOG.md}"

notes="$(
  awk -v heading="## [${version}]" '
    index($0, heading) == 1 { in_section = 1; next }
    in_section && (/^## \[/ || /^\[[^]]+\]: /) { exit }
    in_section { print }
  ' "${changelog}" |
    awk 'NF { started = 1 } started { lines[++n] = $0 } END { while (n > 0 && lines[n] == "") n--; for (i = 1; i <= n; i++) print lines[i] }'
)"

if [ -z "${notes}" ]; then
  echo "release-notes.sh: no notes for version ${version} in ${changelog}" >&2
  exit 1
fi
printf '%s\n' "${notes}"
