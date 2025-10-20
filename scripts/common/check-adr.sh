#!/usr/bin/env bash
# check-adr.sh – Prüft ADR-Konsistenz (Status, Supersessions, Template-Konformität)
# Referenz: copilot-instructions.md Abschnitt 2.5

set -euo pipefail

ADR_DIR="docs/adr"

if [ ! -d "$ADR_DIR" ]; then
    echo "[check-adr] ERROR: $ADR_DIR Verzeichnis nicht gefunden" >&2
    exit 1
fi

echo "[check-adr] Prüfe ADR-Konsistenz in $ADR_DIR..."

# Zähle ADRs
adr_files=$(find "$ADR_DIR" -name "ADR-*.md" -not -name "*template*" | wc -l)
echo "  - Gefunden: $adr_files ADR(s)"

if [ "$adr_files" -eq 0 ]; then
    echo "[check-adr] WARNING: Keine ADRs gefunden (nur Template vorhanden)"
    exit 0
fi

# Validiere jede ADR
errors=0
for adr in "$ADR_DIR"/ADR-*.md; do
    [ -f "$adr" ] || continue
    [ "$(basename "$adr")" == "000-template.md" ] && continue
    
    filename=$(basename "$adr")
    echo "  - Prüfe $filename..."
    
    # Prüfe Pflichtsektionen (flexibel: ## Status oder **Status:** akzeptiert)
    required_sections=("Status" "Kontext" "Entscheidung" "Konsequenzen")
    for section in "${required_sections[@]}"; do
        if ! grep -qiE "^(#+ )?(\*\*)?${section}(\*\*)?:?" "$adr"; then
            echo "    ❌ Fehlende Sektion: $section"
            ((errors++))
        fi
    done
    
    # Prüfe Status (muss ein gültiger Wert sein, flexibles Format)
    # Format 1: Status: Accepted (inline)
    # Format 2: # Status \n Accepted (separate Zeile)
    valid_statuses=("Proposed" "Accepted" "Superseded" "Deprecated" "Rejected")
    status_content=$(grep -iE "^(Status:|#+ Status)" "$adr" -A1 || echo "")
    status_found=false
    for valid_status in "${valid_statuses[@]}"; do
        if echo "$status_content" | grep -qi "$valid_status"; then
            status_found=true
            break
        fi
    done
    
    if [ "$status_found" == false ]; then
        echo "    ❌ Ungültiger oder fehlender Status (erwartet: ${valid_statuses[*]})"
        ((errors++))
    fi
    
    # Prüfe Superseded ADRs auf bidirektionale Referenzen
    if echo "$status_content" | grep -qi "Superseded"; then
        if ! grep -q "Superseded by.*ADR-" "$adr"; then
            echo "    ❌ Superseded Status ohne 'Superseded by ADR-X' Referenz"
            ((errors++))
        fi
    fi
done

if [ "$errors" -gt 0 ]; then
    echo "[check-adr] ❌ $errors Fehler gefunden"
    exit 1
fi

echo "[check-adr] ✅ ADR-Konsistenz validiert"
exit 0
