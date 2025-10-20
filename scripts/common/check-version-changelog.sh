#!/usr/bin/env bash
# check-version-changelog.sh – Release-spezifische Prüfungen (VERSION + CHANGELOG Sync)
# Referenz: copilot-instructions.md Abschnitt 3.6

set -euo pipefail

VERSION_FILE="VERSION"
CHANGELOG_FILE="CHANGELOG.md"

if [ ! -f "$VERSION_FILE" ]; then
    echo "[check-version-changelog] ERROR: $VERSION_FILE nicht gefunden" >&2
    exit 1
fi

if [ ! -f "$CHANGELOG_FILE" ]; then
    echo "[check-version-changelog] ERROR: $CHANGELOG_FILE nicht gefunden" >&2
    exit 1
fi

echo "[check-version-changelog] Prüfe VERSION und CHANGELOG Synchronität..."

# Lese aktuelle Version
current_version=$(cat "$VERSION_FILE" | tr -d '[:space:]')
echo "  - VERSION: $current_version"

# Prüfe ob Version in CHANGELOG erwähnt ist
if ! grep -q "\[$current_version\]" "$CHANGELOG_FILE"; then
    echo "[check-version-changelog] ❌ Version $current_version nicht in CHANGELOG gefunden"
    echo "Erwartetes Format: ## [$current_version] - YYYY-MM-DD"
    exit 1
fi

# Prüfe auf Unreleased Sektion
if ! grep -q "## \[Unreleased\]" "$CHANGELOG_FILE"; then
    echo "[check-version-changelog] WARNING: Keine [Unreleased] Sektion in CHANGELOG"
fi

# Prüfe ob Unreleased Sektion leer ist (bei Release sollte sie befüllt sein)
unreleased_content=$(sed -n '/## \[Unreleased\]/,/## \[/p' "$CHANGELOG_FILE" | grep -v "^## \[" || echo "")
if [ -z "$unreleased_content" ] && [ "${CHECK_RELEASE_LABEL:-0}" == "1" ]; then
    echo "[check-version-changelog] WARNING: [Unreleased] Sektion leer – wurde vergessen zu befüllen?"
fi

echo "[check-version-changelog] ✅ VERSION und CHANGELOG konsistent"
exit 0
