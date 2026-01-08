#!/usr/bin/env bash
set -e

OLDER_TAG="$1"
LATEST_TAG="$2"
OUTFILE="release_note.txt"

COMMITS=$(git log --format='%s' "${OLDER_TAG}..${LATEST_TAG}")

TYPES=(feat fix chore docs refactor test perf style ci build revert)

TMP_MATCHES=$(mktemp)

# Use awk to extract all type-prefixed messages per line, one per line in TMP_MATCHES
echo "$COMMITS" | awk '
{
  line = $0
  # Match type prefix followed by description, stopping before the next type prefix
  while (match(line, /(feat|fix|chore|docs|refactor|test|perf|style|ci|build|revert)(\([^)]+\))?(!)?: /)) {
    prefix = substr(line, RSTART, RLENGTH)
    line = substr(line, RSTART + RLENGTH)
    
    # Find the end of the description (next type prefix or end of line)
    if (match(line, / (feat|fix|chore|docs|refactor|test|perf|style|ci|build|revert)(\([^)]+\))?(!)?: /)) {
      desc = substr(line, 1, RSTART - 1)
      line = substr(line, RSTART)
    } else {
      desc = line
      line = ""
    }
    
    # Trim whitespace and print
    gsub(/^[ \t]+|[ \t]+$/, "", desc)
    if (desc != "") {
      print prefix desc
    }
  }
}
' > "$TMP_MATCHES"

: > "$OUTFILE"
for type in "${TYPES[@]}"; do
    lines=$(grep -E "^${type}(\([^)]+\))?(!)?: " "$TMP_MATCHES" | \
        sed -E "s/^${type}(\([^)]+\))?(!)?: //")
    lines=$(echo "$lines" | sed '/^[[:space:]]*$/d')
    if [[ -n "$lines" ]]; then
        header="$(tr '[:lower:]' '[:upper:]' <<< ${type:0:1})${type:1}"
        echo "### $header" >> "$OUTFILE"
        echo "$lines" | sed 's/^/- /' >> "$OUTFILE"
        echo >> "$OUTFILE"
    fi
done

rm "$TMP_MATCHES"

echo "Wrote release notes to $OUTFILE" >&2