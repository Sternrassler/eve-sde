# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- SDE Download-Script (`scripts/download-sde.sh`) für automatischen Download von JSONL und YAML
- Schema-Fetch-Script (`scripts/fetch-schemas.sh`) generiert Go-Structs aus JSONL-Datenanalyse
- 52 automatisch generierte Go-Structs in `internal/schema/types/`
- Automatische Feldtyp-Erkennung (int64, bool, string, map, slice)
- Data-Verzeichnisstruktur (`data/jsonl/`, `data/yaml/`, `data/sqlite/`)
- Schema-Caching-Verzeichnisse (`internal/schema/definitions/`, `internal/schema/types/`)
- `.gitignore` Regel für `/data/` Verzeichnis
- README mit EVE SDE Projektbeschreibung, Getting Started und Usage

## [0.1.0] - 2025-10-05

- Project initialization.
