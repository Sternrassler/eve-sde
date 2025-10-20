# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- **CLI Tool `sde-schema-gen`** für robuste Schema-Generierung aus JSONL-Dateien
  - Multi-Line JSONL-Analyse (100 Zeilen default) für vollständiges Schema
  - Automatische LocalizedText-Erkennung (8 EVE-Sprachen)
  - Template-basierte Go-Code-Generierung mit text/template
  - Proper CamelCase-Konversion mit ID/NPC/CEO Abbreviation-Handling
  - Nesting als `map[string]interface{}` für maximale Kompatibilität
- SDE Download-Script (`scripts/download-sde.sh`) für automatischen Download von JSONL und YAML
- Schema-Fetch-Script (`scripts/fetch-schemas.sh`) vereinfacht - ruft `sde-schema-gen` auf
- 52 automatisch generierte Go-Structs in `internal/schema/types/`
- `LocalizedText` Common-Type für mehrsprachige Felder
- Data-Verzeichnisstruktur (`data/jsonl/`, `data/yaml/`, `data/sqlite/`)
- `.gitignore` Regel für `/data/` Verzeichnis
- README mit EVE SDE Projektbeschreibung, Getting Started und Usage

### Changed

- Schema-Generierung von Bash+Python zu dediziertem Go CLI-Tool migriert
- Feldtyp-Erkennung jetzt multi-line basiert (statt nur erste Zeile)
- Verschachtelte Strukturen verwenden `map[string]interface{}` statt fehlende Sub-Structs

### Removed

- Python-Code aus `fetch-schemas.sh` entfernt
- HTML-Scraping Logik (nicht mehr benötigt)

## [0.1.0] - 2025-10-05

- Project initialization.
