#!/usr/bin/env bash
# check-version-changelog.sh – Release-spezifische Prüfungen (CHANGELOG ist Single Source of Truth für SemVer)

set -euo pipefail

CHANGELOG_FILE="CHANGELOG.md"

if [ ! -f "$CHANGELOG_FILE" ]; then
    echo "[check-version-changelog] ERROR: $CHANGELOG_FILE nicht gefunden" >&2
    exit 1
fi

echo "[check-version-changelog] Prüfe CHANGELOG..."

# Ermittle die jüngste veröffentlichte Version (oberster [X.Y.Z]-Eintrag, ohne [Unreleased])
current_version=$(grep -oE '^## \[[0-9]+\.[0-9]+\.[0-9]+\]' "$CHANGELOG_FILE" | head -n1 | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' || echo "")
if [ -z "$current_version" ]; then
    echo "[check-version-changelog] ❌ Keine veröffentlichte Version in $CHANGELOG_FILE gefunden"
    echo "Erwartetes Format: ## [X.Y.Z] - YYYY-MM-DD"
    exit 1
fi
echo "  - Letzte Version (aus CHANGELOG): $current_version"

# Prüfe auf Unreleased Sektion
if ! grep -q "## \[Unreleased\]" "$CHANGELOG_FILE"; then
    echo "[check-version-changelog] WARNING: Keine [Unreleased] Sektion in CHANGELOG"
fi

# Prüfe ob Unreleased Sektion leer ist (bei Release sollte sie befüllt sein)
unreleased_content=$(sed -n '/## \[Unreleased\]/,/## \[/p' "$CHANGELOG_FILE" | grep -v "^## \[" || echo "")
if [ -z "$unreleased_content" ] && [ "${CHECK_RELEASE_LABEL:-0}" == "1" ]; then
    echo "[check-version-changelog] WARNING: [Unreleased] Sektion leer – wurde vergessen zu befüllen?"
fi

echo "[check-version-changelog] ✅ CHANGELOG konsistent"
exit 0
