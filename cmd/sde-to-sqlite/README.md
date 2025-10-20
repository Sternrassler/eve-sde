# SDE to SQLite Converter

Go-basierter Konverter für EVE Online Static Data Export (JSONL) nach SQLite.

## Build

```bash
go build -o bin/sde-to-sqlite ./cmd/sde-to-sqlite
```

## Verwendung

```bash
# Konvertiere alle JSONL-Dateien
./bin/sde-to-sqlite --input data/jsonl --output data/sqlite/eve-sde.db

# Nur bestimmte Dateien
./bin/sde-to-sqlite --input data/jsonl/types.jsonl --output data/sqlite/eve-sde.db
```

## Entwicklungsstatus

🚧 **In Entwicklung** - Noch nicht funktional

### Geplante Features

- [ ] JSONL Streaming Parser
- [ ] Schema-Validierung (basierend auf <https://sde.riftforeve.online/>)
- [ ] SQLite Tabellenstruktur Auto-Generation
- [ ] Batch-Insert Optimierung
- [ ] Indizes für häufige Queries
- [ ] Progress Reporting
- [ ] Fehlerbehandlung & Validierung

## Architektur

```text
data/jsonl/*.jsonl → Go Parser → SQLite DB
                      ↓
               Schema Validation
                      ↓
                  Transform
                      ↓
              Batch Insert
```

## Dependencies

- `github.com/mattn/go-sqlite3` - SQLite Driver
- Standard Library (encoding/json, bufio)
