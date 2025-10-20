# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- **Navigation & Intelligence System** (`pkg/evedb/navigation/`)
  - Go API für Pathfinding, Travel Time Berechnung, Security Filtering
  - SQL Views: `v_stargate_graph`, `v_system_info`, `v_system_security_zones`, `v_region_stats`, `v_trade_hubs`
  - Performance: Kurze Routen (<10 Jumps) in <100ms
  - Beispiel-Programm: `examples/navigation_example.go`

### Changed

- **Architektur-Refactoring: DB-Core vs. API Layer Separation** (ADR-001)
  - Navigation API verschoben: `internal/sqlite/navigation/` → `pkg/evedb/navigation/`
  - SQL Views extrahiert: `internal/sqlite/views/` (DB-Core)
  - Rationale: SQLite-Datenbank ist Kernprodukt, Go APIs sind optionale Convenience Layer
  - Ermöglicht externe Nutzung der API ohne Exposition interner DB-Implementierung

### Fixed

- View-Schema: `v_system_info` nutzt korrektes Feld (`_key` statt `solarSystemID`)

### Removed

- Workflow: `.github/workflows/issue-dependency-enforcement.yml` (nicht genutzt)

### Known Issues

- Performance: Lange Routen (40+ Jumps) haben O(n³) Komplexität wegen JSON-basierter Cycle-Detection (Issue #3)

### Added

- **Automatisches Release-System** (GitHub Actions)
  - Workflow `.github/workflows/sync-sde-release.yml`
  - Täglicher Cron-Job (03:00 UTC) prüft auf neue SDE-Versionen
  - Bei Update: Automatischer Build und GitHub Release
  - Release-Tag Format: `sde-v{buildNumber}-{datum}` (z.B. `sde-v3064089-2025-10-17`)
  - Asset: `eve-sde.db.gz` (gzip level 9 komprimiert)
  - Automatische Release-Notes mit Download-Beispielen
  - Retention: Löscht Releases älter als 2 Jahre automatisch
  - Validierung: Row counts und DB-Größe werden geprüft
  - Idempotent: Überspringt bereits existierende Releases

## [0.1.0] - 2025-10-20

### Added

- **Version Tracking System** (`internal/sde/version`)
  - Integration mit CCP Developer API (`https://developers.eveonline.com/static-data/tranquility/latest.jsonl`)
  - `GetLatestVersion()`: HTTP-basierter Abruf von Build-Nummer & Release-Datum
  - `GetLocalVersion()`: SQLite-basierte Abfrage der lokal importierten Version aus `_sde` Tabelle
  - `NeedsUpdate()`: Intelligenter Vergleich via Build-Nummer (latest > local)
  - Graceful Degradation bei fehlender Datenbank/Tabelle (nil statt Fehler)
  - CLI-Flags für `sde-to-sqlite`: `--check-version`, `--skip-if-current`
  - Umfangreiche Test-Suite (4 Tests, inkl. Integration Test mit CCP API)

- **Sync-Pipeline Automatisierung** (`cmd/sde-sync`)
  - Vollautomatischer Workflow: Version Check → Download → Schema-Generierung → SQLite-Import
  - CLI-Flags:
    - `--force`: Erzwinge Update (auch wenn DB aktuell)
    - `--skip-import`: Nur Download + Schema-Generierung (kein SQLite)
    - `-v`: Verbose Logging (stdout von Subprozessen)
    - `--data`: Custom Basis-Verzeichnis
  - Intelligentes Skip-Verhalten: Überspringt Pipeline wenn DB auf aktuellem Stand
  - Makefile Targets:
    - `make sync`: Intelligenter Sync (nur bei Update)
    - `make sync-force`: Erzwinge vollständigen Sync
    - `make sync-download-only`: Nur Download + Schema
  - Timing & Progress-Reporting (Gesamtdauer wird angezeigt)
  - Fehlerbehandlung: Warnings für nicht-kritische Fehler, Fatal für kritische
  - Dokumentation in `cmd/sde-sync/README.md`:
    - Workflow-Diagramm
    - Automation-Beispiele (cron, systemd timer)
    - Error-Handling-Dokumentation

- **SQLite Database Implementation** (Complete Pipeline)
  - Reflection-based schema generator (`internal/sqlite/schema`)
    - `GenerateTable()`: Go struct → CREATE TABLE DDL
    - Type mapping: int64→INTEGER, float64→REAL, bool→INTEGER, LocalizedText→TEXT (JSON)
    - Primary key detection on `_key` fields
    - Required field detection via `omitempty` tag
    - Field validation for indices (skip non-existent fields)
  - Streaming JSONL importer (`internal/sqlite/importer`)
    - `ImportJSONL()`: Streams JSONL with `bufio.Scanner`
    - Batch inserts (1000 rows per batch)
    - Single transaction per file
    - SQLite PRAGMAs: WAL mode, NORMAL sync, cache_size=10000
    - Type conversion: bool→int, complex types→JSON
  - CLI Tool: `cmd/sde-to-sqlite`
    - Flags: `--db`, `--jsonl`, `--init`, `--import`
    - 41 schema mappings with index specifications
    - Full import: 24s for 500k rows, 405 MB database
  - Performance metrics:
    - types table: 50,486 rows in 4s (12.6k rows/sec)
    - mapMoons: 342,170 rows in 13s (26k rows/sec)
    - Compression: 499 MB JSONL → 405 MB SQLite (18% reduction)
  - Validation: All row counts match JSONL exactly, LocalizedText stored as JSON
  - Documentation: `docs/sqlite-implementation.md` with architecture & metrics
  - Dependencies: Added `github.com/mattn/go-sqlite3` v1.14.32

- **CLI Tool `sde-schema-gen`** für robuste Schema-Generierung aus JSONL-Dateien
  - Multi-Line JSONL-Analyse (100 Zeilen default) für vollständiges Schema
  - Automatische LocalizedText-Erkennung (8 EVE-Sprachen) - 49 Felder erkannt
  - Template-basierte Go-Code-Generierung mit `text/template`
  - Proper CamelCase-Konversion mit ID/NPC/CEO Abbreviation-Handling
  - **Intelligente Typ-Inferenz:** Ignoriert null-Werte, int64/float64 Mix → float64
  - **Required-Detection:** `_key` immer required, andere Felder basierend auf Vorkommen
  - Nesting als `map[string]interface{}` für maximale Kompatibilität
  - 53 automatisch generierte Go-Structs in `internal/schema/types/`

- **Data Infrastructure**
  - SDE Download-Script (`scripts/download-sde.sh`) für automatischen Download von JSONL (52 Dateien, ~499 MB) und YAML (52 Dateien, ~160 MB)
  - Schema-Fetch-Script (`scripts/fetch-schemas.sh`) als Wrapper für `sde-schema-gen`
  - Data-Verzeichnisstruktur (`data/jsonl/`, `data/yaml/`, `data/sqlite/`)
  - `.gitignore` Regel für `/data/` Verzeichnis

- **Documentation & Governance**
  - README mit Projektstatus v0.1.0, Getting Started, SQLite Usage, Architektur
  - Engineering-Richtlinien in `.github/copilot-instructions.md` (TDD, Git-Workflow, ADRs)
  - Pre-Commit Hooks: Normative Checks, ADR-Validierung, Secret-Scanning
  - Issue Templates (Feature / Bug)
  - `LocalizedText` Common-Type für mehrsprachige Felder

### Changed

- Schema-Generierung von Bash+Python zu dediziertem Go CLI-Tool (`sde-schema-gen`) migriert
- Feldtyp-Erkennung jetzt multi-line basiert (statt nur erste Zeile)
- Verschachtelte Strukturen verwenden `map[string]interface{}` statt fehlende Sub-Structs
- **Type Precision:** `types.mass` und `types.volume` jetzt korrekt als `float64` (vorher `interface{}`)
- **Optionalität:** Nur tatsächlich optionale Felder haben `omitempty`, `_key` immer required

### Fixed

- Null-Werte führen nicht mehr zu `interface{}` Fallback
- Gemischte int64/float64 Typen werden zu float64 statt interface{}
- 50 `_key` Felder ohne `omitempty` (korrekt als Primary Key)

### Removed

- Python-Code aus `fetch-schemas.sh` entfernt (ersetzt durch Go CLI)
- HTML-Scraping Logik (nicht mehr benötigt)
