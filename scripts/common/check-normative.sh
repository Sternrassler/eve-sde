#!/usr/bin/env bash
# check-normative.sh – Validiert normative Labels (MUST/SHOULD/MAY) in copilot-instructions.md
# Referenz: copilot-instructions.md Abschnitt 4

set -euo pipefail

INSTRUCTIONS_FILE=".github/copilot-instructions.md"

if [ ! -f "$INSTRUCTIONS_FILE" ]; then
    echo "[check-normative] ERROR: $INSTRUCTIONS_FILE nicht gefunden" >&2
    exit 1
fi

echo "[check-normative] Prüfe normative Labels in $INSTRUCTIONS_FILE..."

# Erwartete normative Keywords (RFC 2119)
EXPECTED_KEYWORDS=("MUST" "MUST NOT" "SHOULD" "SHOULD NOT" "MAY")

# Zähle Vorkommen jedes Keywords
declare -A keyword_counts
for kw in "${EXPECTED_KEYWORDS[@]}"; do
    count=$(grep -c "\b$kw\b" "$INSTRUCTIONS_FILE" || true)
    keyword_counts[$kw]=$count
    echo "  - $kw: $count Vorkommen"
done

# Validierung: Mindestens ein MUST vorhanden
if [ "${keyword_counts[MUST]}" -eq 0 ]; then
    echo "[check-normative] ERROR: Keine MUST Labels gefunden – Instructions unvollständig" >&2
    exit 1
fi

# Warnung bei fehlenden SHOULD
if [ "${keyword_counts[SHOULD]}" -eq 0 ]; then
    echo "[check-normative] WARNING: Keine SHOULD Labels gefunden"
fi

# Prüfe auf inkonsistente Schreibweisen (lowercase)
lowercase_matches=$(grep -E "\b(must|should|may)\b" "$INSTRUCTIONS_FILE" | wc -l || true)
if [ "$lowercase_matches" -gt 0 ]; then
    echo "[check-normative] WARNING: Gefunden $lowercase_matches lowercase normative keywords (sollten UPPERCASE sein)"
fi

echo "[check-normative] ✅ Normative Labels validiert"
exit 0
