#!/usr/bin/env bash
# Schema Fetch Script für EVE SDE
# Lädt Go Code Snippets von https://sde.riftforeve.online/ und cached sie lokal

set -euo pipefail

SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
SCHEMA_DIR="${PROJECT_ROOT}/internal/schema"
TYPES_DIR="${SCHEMA_DIR}/types"
DEFINITIONS_DIR="${SCHEMA_DIR}/definitions"

BASE_URL="https://sde.riftforeve.online/schema"

# Farben
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

log_info() {
    echo -e "${GREEN}[INFO]${NC} $*"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $*"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $*"
}

# Liste aller SDE Files (basierend auf heruntergeladenen JSONL-Dateien)
get_sde_files() {
    if [ -d "${PROJECT_ROOT}/data/jsonl" ]; then
        # Ohne exec um SIGPIPE zu vermeiden
        find "${PROJECT_ROOT}/data/jsonl" -name "*.jsonl" | while read -r file; do
            basename "$file" .jsonl
        done
    else
        log_warn "data/jsonl nicht gefunden, verwende Fallback-Liste"
        echo "_sde agentTypes agentsInSpace ancestries bloodlines blueprints categories certificates"
    fi
}

# Fetch Go code snippet für eine Datei
fetch_go_snippet() {
    local filename="$1"
    local output_file="${TYPES_DIR}/${filename}.go"
    local url="${BASE_URL}/${filename}/"
    
    # Type name: CamelCase (ersten Buchstaben groß)
    local typename
    typename="$(echo "${filename:0:1}" | tr '[:lower:]' '[:upper:]')${filename:1}"
    
    log_info "Fetching schema for ${filename}..."
    
    # TODO: Implementiere HTML-Parsing um Go-Code-Snippet zu extrahieren
    # Aktuell: Platzhalter-Struct
    cat > "$output_file" <<EOF
// Code generated from https://sde.riftforeve.online/schema/${filename}/
// DO NOT EDIT manually - regenerate with scripts/fetch-schemas.sh --refresh

package types

// ${typename} represents the schema for ${filename}.jsonl
type ${typename} struct {
	Key int64 \`json:"_key"\`
	// TODO: Complete schema from sde.riftforeve.online
}
EOF
    
    log_info "✓ Generated ${output_file}"
}

# Hauptlogik
main() {
    local refresh=false
    
    if [ "${1:-}" = "--refresh" ] || [ "${1:-}" = "-r" ]; then
        refresh=true
        log_info "Refresh mode aktiviert"
    fi
    
    # Check ob Schemas bereits existieren
    if [ -d "$TYPES_DIR" ] && [ "$(ls -A "$TYPES_DIR" 2>/dev/null)" ] && [ "$refresh" = false ]; then
        log_info "Schemas bereits gecached. Verwende --refresh zum Aktualisieren."
        exit 0
    fi
    
    mkdir -p "$TYPES_DIR" "$DEFINITIONS_DIR"
    
    log_info "Fetching schemas von ${BASE_URL}..."
    
    local count=0
    
    # Direkt über JSONL-Dateien iterieren
    for jsonl_file in "${PROJECT_ROOT}/data/jsonl"/*.jsonl; do
        [ -e "$jsonl_file" ] || continue
        local filename
        filename=$(basename "$jsonl_file" .jsonl)
        fetch_go_snippet "$filename" || {
            log_error "Fehler bei $filename"
            continue
        }
        count=$((count + 1))
    done
    
    if [ $count -eq 0 ]; then
        log_error "Keine JSONL-Dateien in data/jsonl/ gefunden!"
        log_info "Bitte zuerst ./scripts/download-sde.sh ausführen"
        exit 1
    fi
    
    log_info "✓ $count Schema-Dateien generiert"
    log_info "Schemas gespeichert in: $TYPES_DIR"
    log_warn "HINWEIS: Aktuell werden nur Platzhalter-Structs generiert."
    log_warn "TODO: HTML-Parsing implementieren um echte Go-Snippets zu extrahieren."
}

main "$@"
