# sde-sync

Vollautomatische EVE SDE Synchronisations-Pipeline.

**Status:** ✅ **v0.1.0 - Produktionsreif**

## Features

- ✅ Automatische Versionsprüfung gegen CCP API
- ✅ Download von JSONL + YAML Daten
- ✅ Go-Schema-Generierung
- ✅ SQLite-Import mit Batch-Processing
- ✅ Intelligentes Überspringen bei aktueller Version
- ✅ Force-Modus für manuelle Updates

## Verwendung

```bash
# Normaler Sync (nur bei neuer Version)
go run ./cmd/sde-sync

# Oder via Makefile
make sync
```

### CLI-Flags

- `--data DIR`: Data-Verzeichnis (default: `data`)
- `--force`: Force Update (ignoriert Versionsprüfung)
- `--skip-import`: Nur Download + Schema-Gen (kein SQLite)
- `-v`: Verbose Output (zeigt alle Befehle)
- `--version`: Version anzeigen

### Makefile Targets

```bash
make sync              # Vollständiger Sync
make sync-force        # Force Sync
make sync-download-only # Nur Download + Schemas
```

## Workflow

```text
1. Version Check
   ├─ Latest: https://developers.eveonline.com/static-data/tranquility/latest.jsonl
   ├─ Local:  data/sqlite/eve-sde.db (_sde table)
   └─ Skip if BuildNumber matches

2. Download SDE
   └─ scripts/download-sde.sh (JSONL + YAML)

3. Generate Schemas
   └─ cmd/sde-schema-gen → internal/schema/types/

4. Import SQLite
   └─ cmd/sde-to-sqlite (41 tables, batch processing)
```

## Beispiel-Output

```bash
$ make sync
2025/10/20 11:46:52 EVE SDE Sync v0.1.0
2025/10/20 11:46:52 → Checking for SDE updates...
2025/10/20 11:46:52   Latest SDE: Build 3064089 (2025-10-17)
2025/10/20 11:46:52   Local DB:   Build 3064089 (2025-10-17)
2025/10/20 11:46:52 ✓ Database is up-to-date
```

```bash
$ make sync-force
2025/10/20 12:00:00 EVE SDE Sync v0.1.0
2025/10/20 12:00:00 → Force mode enabled, skipping version check
2025/10/20 12:00:00 → Downloading SDE data...
2025/10/20 12:00:15 ✓ SDE downloaded
2025/10/20 12:00:15 → Generating Go schemas...
2025/10/20 12:00:17 ✓ Schemas generated
2025/10/20 12:00:17 → Importing to SQLite...
2025/10/20 12:00:41 ✓ SQLite import completed
2025/10/20 12:00:41 ✓ Sync completed in 41s
```

## Automatisierung

### Cron Job (täglich um 02:00)

```bash
0 2 * * * cd /path/to/eve-sde && make sync >> logs/sync.log 2>&1
```

### Systemd Timer

```ini
[Unit]
Description=EVE SDE Sync

[Service]
Type=oneshot
WorkingDirectory=/path/to/eve-sde
ExecStart=/usr/bin/make sync

[Install]
WantedBy=multi-user.target
```

```ini
[Unit]
Description=EVE SDE Sync Timer

[Timer]
OnCalendar=daily
Persistent=true

[Install]
WantedBy=timers.target
```

## Error Handling

- **Version Check fehlgeschlagen**: Warnung, fährt mit Update fort
- **Download fehlgeschlagen**: Abbbruch mit Fehler
- **Schema-Gen fehlgeschlagen**: Warnung, nutzt existierende Schemas
- **SQLite Import fehlgeschlagen**: Abbruch mit Fehler

## Dependencies

- `bash` (für download-sde.sh)
- `curl` / `wget` (für SDE Download)
- Go 1.x
- SQLite3 (via go-sqlite3)

## Siehe auch

- [cmd/sde-to-sqlite](../sde-to-sqlite/README.md) - SQLite Importer
- [cmd/sde-schema-gen](../sde-schema-gen/README.md) - Schema Generator
- [docs/sqlite-implementation.md](../../docs/sqlite-implementation.md) - SQLite Details
