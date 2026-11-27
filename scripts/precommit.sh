#!/usr/bin/env bash
# The local check. It runs every check of the repository and prints one
# verdict on its last line. Run it in the repository root:
#
#   bash scripts/precommit.sh

set -uo pipefail
cd "$(dirname "$0")/.."

failed=()

step() {
  local name=$1
  shift
  echo "== $name"
  if ! "$@"; then
    failed+=("$name")
  fi
}

gofmt_clean() {
  local files
  files=$(gofmt -l cmd internal)
  if [ -n "$files" ]; then
    echo "gofmt would change:"
    echo "$files"
    return 1
  fi
}

step gofmt gofmt_clean
step "go vet" go vet ./cmd/... ./internal/...
step "go test" go test ./cmd/... ./internal/...
step prettier npx prettier --check .
step "astro check" npx astro check
step build npm run build --silent
step "third-party loads" node scripts/check-dist.mjs

echo
if [ ${#failed[@]} -eq 0 ]; then
  echo "PASS"
else
  echo "FAIL: ${failed[*]}"
  exit 1
fi
