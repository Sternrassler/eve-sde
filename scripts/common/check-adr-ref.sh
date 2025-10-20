#!/usr/bin/env bash
# check-adr-ref.sh – Erzwingt ADR-Referenzen in PRs bei Governance-Pfad-Änderungen
# Referenz: copilot-instructions.md Abschnitt 3.0 (Quick Start Hinweis)
# Exit Codes: 0 = OK, 1 = Fehler (blockiert), 2 = Skip Marker (Warnung)

set -euo pipefail

echo "[check-adr-ref] Prüfe ADR-Referenz-Enforcement..."

# Nur in PR-Kontext relevant
if [ -z "${GITHUB_EVENT_PATH:-}" ]; then
    echo "[check-adr-ref] Nicht im GitHub Actions PR Kontext – überspringe"
    exit 0
fi

# Prüfe ob PR Body existiert
PR_BODY=$(jq -r '.pull_request.body // ""' "$GITHUB_EVENT_PATH")

if [ -z "$PR_BODY" ]; then
    echo "[check-adr-ref] WARNING: PR Body leer – kann ADR-Referenzen nicht prüfen"
    exit 2
fi

# Prüfe auf geänderte Governance-Pfade
CHANGED_FILES=$(git diff --name-only origin/main...HEAD || echo "")
GOVERNANCE_PATHS=("docs/adr/" "scripts/" ".github/" "Makefile" "docs/")

governance_changed=false
for path in "${GOVERNANCE_PATHS[@]}"; do
    if echo "$CHANGED_FILES" | grep -q "^$path"; then
        governance_changed=true
        echo "  - Governance-Pfad geändert: $path"
    fi
done

if [ "$governance_changed" == false ]; then
    echo "[check-adr-ref] Keine Governance-Pfade geändert – ADR-Referenz nicht erforderlich"
    exit 0
fi

# Prüfe auf ADR-Referenz im PR Body (Format: ADR-XXX)
if echo "$PR_BODY" | grep -qE "ADR-[0-9]{3}"; then
    echo "[check-adr-ref] ✅ ADR-Referenz gefunden"
    exit 0
fi

# Prüfe auf Skip Marker (ADR-NOT-REQ)
if echo "$PR_BODY" | grep -q "ADR-NOT-REQ"; then
    echo "[check-adr-ref] ⚠️  Skip Marker (ADR-NOT-REQ) gefunden – Warnung statt Block"
    exit 2
fi

# Fehler: Governance-Pfad geändert aber keine ADR-Referenz
echo "[check-adr-ref] ❌ Governance-Pfad geändert aber keine ADR-Referenz im PR Body gefunden"
echo "Erwartete Formate:"
echo "  - 'ADR-001' oder 'Referenziert ADR-005, ADR-008'"
echo "  - Oder Skip Marker: 'ADR-NOT-REQ' mit Begründung"
exit 1
