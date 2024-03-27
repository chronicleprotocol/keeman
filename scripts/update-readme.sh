#!/usr/bin/env bash
cd "$(dirname "$0")/.."
set -euo pipefail

_commands=("" "derive" "generate" "list")

for _cmd in "${_commands[@]}"; do
	if [[ $_cmd == "" ]]; then
		_marker="## Help Page"
	else
		_marker="## Help Page for \`$_cmd\` Command"
	fi
	echo >&2 "$_marker"
help_output="$_marker

\`\`\`text
$(go run . $_cmd --help)
\`\`\`"
	help_output=$(printf '%q' "$help_output")
	help_output=$(echo "$help_output" | sed "s/^\\$'//" | sed "s/'$//")

	sed -i "/^$_marker$/,/^\`\`\`$/c\\$help_output" README.md
done
