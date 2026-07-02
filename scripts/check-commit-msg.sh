#!/usr/bin/env bash

set -euo pipefail

readonly TYPES=(
  feat
  fix
  refactor
  test
  docs
  chore
  build
  ci
)

# Build regex from commit types
types_regex=$(IFS='|'; printf "%s" "${TYPES[*]}")

readonly COMMIT_REGEX="^(${types_regex})(\\([a-z0-9-]+\\))?: .+$"

commit_message="$(<"$1")"

if [[ ! "$commit_message" =~ $COMMIT_REGEX ]]; then
    echo
    echo "❌ Invalid commit message."
    echo
    echo "Expected format:"
    echo
    echo "  <type>(<scope>): <description>"
    echo
    echo "Examples:"
    echo "  feat(github): add contribution calendar query"
    echo "  refactor(commute): simplify route mapping"
    echo "  fix(server): handle shutdown gracefully"
    echo
    echo "Available commit types:"
    for type in "${TYPES[@]}"; do
        printf "  • %s\n" "$type"
    done
    echo

    exit 1
fi
