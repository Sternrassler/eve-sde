#!/usr/bin/env bash
# EVE Online SDE Download Script
# Lädt YAML und JSONL SDE-Exports von developers.eveonline.com herunter
#
# Quelle: https://developers.eveonline.com/docs/services/static-data/yaml
# Blog: https://developers.eveonline.com/blog/reworking-the-sde-a-fresh-start-for-static-data

set -euo pipefail

# URLs für die neueste SDE-Version (automatisch aktualisiert)
YAML_URL="${YAML_URL:-https://developers.eveonline.com/static-data/eve-online-static-data-latest-yaml.zip}"
JSONL_URL="${JSONL_URL:-https://developers.eveonline.com/static-data/eve-online-static-data-latest-jsonl.zip}"

DATA_DIR="$(dirname "$(dirname "$(readlink -f "$0")")")/data"
YAML_DIR="${DATA_DIR}/yaml"
JSONL_DIR="${DATA_DIR}/jsonl"
TEMP_DIR="${DATA_DIR}/.tmp"

# Farben für Output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $*"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $*"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $*"
}

# Verzeichnisse erstellen
mkdir -p "${YAML_DIR}" "${JSONL_DIR}" "${TEMP_DIR}"

# YAML SDE herunterladen
log_info "Lade YAML SDE herunter von: ${YAML_URL}"
YAML_ZIP="${TEMP_DIR}/sde-yaml.zip"

if curl -fSL -o "${YAML_ZIP}" "${YAML_URL}"; then
    log_info "YAML ZIP erfolgreich heruntergeladen ($(du -h "${YAML_ZIP}" | cut -f1))"
    
    log_info "Extrahiere YAML Archive..."
    unzip -q -o "${YAML_ZIP}" -d "${YAML_DIR}"
    log_info "YAML Daten nach ${YAML_DIR} extrahiert"
    rm -f "${YAML_ZIP}"
else
    log_error "YAML Download fehlgeschlagen: ${YAML_URL}"
    exit 1
fi

# JSONL SDE herunterladen
log_info "Lade JSONL SDE herunter von: ${JSONL_URL}"
JSONL_ZIP="${TEMP_DIR}/sde-jsonl.zip"

if curl -fSL -o "${JSONL_ZIP}" "${JSONL_URL}"; then
    log_info "JSONL ZIP erfolgreich heruntergeladen ($(du -h "${JSONL_ZIP}" | cut -f1))"
    
    log_info "Extrahiere JSONL Archive..."
    unzip -q -o "${JSONL_ZIP}" -d "${JSONL_DIR}"
    log_info "JSONL Daten nach ${JSONL_DIR} extrahiert"
    rm -f "${JSONL_ZIP}"
else
    log_error "JSONL Download fehlgeschlagen: ${JSONL_URL}"
    exit 1
fi

# Aufräumen
rm -rf "${TEMP_DIR}"

log_info "SDE Download abgeschlossen!"
log_info "YAML: ${YAML_DIR}"
log_info "JSONL: ${JSONL_DIR}"

# Statistiken anzeigen
if command -v tree &> /dev/null; then
    log_info "Verzeichnisstruktur (Top-Level):"
    tree -L 2 "${DATA_DIR}" 2>/dev/null || true
fi
