#!/usr/bin/env bash
# run.sh – orchestrates PR quality gate checks

set -euo pipefail

ROOT_DIR=$(cd "$(dirname "$0")/../.." && pwd)

# Normative labels & ADR governance
bash "$ROOT_DIR/common/check-normative.sh"
bash "$ROOT_DIR/common/check-adr.sh"

# ADR reference enforcement via make target (handles skip markers)
make adr-ref

# Commit message lint (if script exists)
if [ -x "$ROOT_DIR/common/check-commit-msg.sh" ]; then
  RANGE=${COMMIT_LINT_RANGE:-origin/main..HEAD}
  git fetch origin main:refs/remotes/origin/main >/dev/null 2>&1 || true
  bash "$ROOT_DIR/common/check-commit-msg.sh" --range "$RANGE"
else
  echo "[pr-quality-gates] commit message check skipped (script missing)" >&2
fi

# Optional release label check
echo "[pr-quality-gates] CHECK_RELEASE_LABEL=${CHECK_RELEASE_LABEL:-0}"
if [ "${CHECK_RELEASE_LABEL:-0}" = "1" ]; then
  if [ -x "$ROOT_DIR/common/check-version-changelog.sh" ]; then
    bash "$ROOT_DIR/common/check-version-changelog.sh"
  else
    echo "[pr-quality-gates] release check skipped (script missing)" >&2
  fi
fi

# Security scans (single Trivy invocation + blocker gate)
make scan-json
bash "$ROOT_DIR/common/check-security-blockers.sh"

echo "[pr-quality-gates] ✅ Quality gates completed"
