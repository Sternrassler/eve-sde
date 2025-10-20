# SDE to SQLite Converter

Go-basierter Konverter für EVE Online Static Data Export (JSONL) nach SQLite.

**Status:** ✅ **v0.1.0 - Produktionsreif**

## Features

- ✅ Reflection-basierte DDL-Generierung aus Go-Structs
- ✅ Streaming JSONL-Parser mit `bufio.Scanner`
- ✅ Batch-Insert-Optimierung (1000 Zeilen/Batch)
- ✅ SQLite-Performance-Tuning (WAL, PRAGMA settings)
- ✅ Indizes auf Foreign Keys (validiert gegen Struct-Felder)
- ✅ LocalizedText als JSON TEXT (8 Sprachen)
- ✅ Type Conversion (bool→INTEGER, complex→JSON)

## Build

```bash
go build ./cmd/sde-to-sqlite
```

## Verwendung

```bash
# Full Import (alle 41 Schemas)
go run ./cmd/sde-to-sqlite

# Einzeltabelle importieren
go run ./cmd/sde-to-sqlite --import types

# Custom DB-Pfad
go run ./cmd/sde-to-sqlite --db custom/eve.db --jsonl data/jsonl

# Nur Schema erstellen (ohne Daten)
go run ./cmd/sde-to-sqlite --init
```

### CLI-Flags

- `--db PATH`: SQLite-Datenbank-Pfad (default: `data/sqlite/eve-sde.db`)
- `--jsonl DIR`: JSONL-Input-Verzeichnis (default: `data/jsonl`)
- `--init`: Nur Schema erstellen, keine Daten importieren
- `--import TABLE`: Nur spezifische Tabelle importieren (default: alle)
- `--check-version`: Prüft auf SDE-Updates (vergleicht mit https://developers.eveonline.com)
- `--skip-if-current`: Überspringt Import wenn Datenbank aktuell ist
- `--version`: Version anzeigen

### Version Tracking

Automatische Versionsprüfung gegen CCP's offizielle SDE-API:

```bash
# Nur Version prüfen
go run ./cmd/sde-to-sqlite --check-version

# Import nur bei neuer Version
go run ./cmd/sde-to-sqlite --skip-if-current
```

## Performance

| Metrik | Wert |
|--------|------|
| Full Import | 24 Sekunden (41 Tabellen) |
| Datensätze | ~500.000 Zeilen |
| DB-Größe | 405 MB |
| Durchsatz | ~20.000 Zeilen/Sekunde |

## Architektur

```text
Go Structs → Reflection → CREATE TABLE DDL
     ↓                          ↓
JSONL Files → Stream Parser → Batch Insert → SQLite
                    ↓
            Type Conversion
            (bool→int, JSON)
```

## Dependencies

- `github.com/mattn/go-sqlite3` v1.14.32 - SQLite Driver
- Standard Library: `database/sql`, `encoding/json`, `reflect`

## Weitere Dokumentation

Siehe [docs/sqlite-implementation.md](../../docs/sqlite-implementation.md) für Details.
