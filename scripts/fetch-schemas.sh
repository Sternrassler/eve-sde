#!/usr/bin/env bash
# Schema Generator Script für EVE SDE
# Ruft sde-schema-gen CLI Tool auf

set -euo pipefail

SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
TOOL_BIN="${PROJECT_ROOT}/bin/sde-schema-gen"
INPUT_DIR="${PROJECT_ROOT}/data/jsonl"
OUTPUT_DIR="${PROJECT_ROOT}/internal/schema/types"

# Farben
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

log_info() {
    echo -e "${GREEN}[INFO]${NC} $*"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $*"
}

# Check ob JSONL Daten existieren
if [ ! -d "$INPUT_DIR" ] || [ -z "$(ls -A "$INPUT_DIR"/*.jsonl 2>/dev/null)" ]; then
    log_error "Keine JSONL-Dateien gefunden in: $INPUT_DIR"
    log_info "Bitte erst Download-Script ausführen: scripts/download-sde.sh"
    exit 1
fi

# Check ob Tool gebaut wurde
if [ ! -f "$TOOL_BIN" ]; then
    log_info "Building sde-schema-gen..."
    go build -o "$TOOL_BIN" "${PROJECT_ROOT}/cmd/sde-schema-gen"
fi

# Parse Arguments
VERBOSE=""
REFRESH=false
MAX_LINES=100

while [[ $# -gt 0 ]]; do
    case $1 in
        --refresh)
            REFRESH=true
            shift
            ;;
        -v|--verbose)
            VERBOSE="-v"
            shift
            ;;
        --lines)
            MAX_LINES="$2"
            shift 2
            ;;
        *)
            log_error "Unknown option: $1"
            echo "Usage: $0 [--refresh] [--verbose] [--lines N]"
            exit 1
            ;;
    esac
done

# Cleanup bei --refresh
if [ "$REFRESH" = true ]; then
    log_info "Cleaning old schemas..."
    rm -f "${OUTPUT_DIR}"/*.go
fi

# Run Schema Generator
log_info "Generating Go schemas from JSONL data..."
"$TOOL_BIN" \
    -input "$INPUT_DIR" \
    -output "$OUTPUT_DIR" \
    -lines "$MAX_LINES" \
    $VERBOSE

log_info "✓ Schema generation completed"
log_info "Run 'go build ./internal/schema/types/...' to verify"
