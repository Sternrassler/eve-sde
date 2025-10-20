# EVE Online Static Data Export (SDE) Synchronisation

Automatisierte Synchronisation und Aufbereitung der EVE Online Static Data Export (SDE) von der [CCP Developer API](https://developers.eveonline.com/docs/services/static-data/).

## Ziel

Dieses Projekt dient der:

1. **Synchronisation** der EVE SDE mit einem lokalen Verzeichnis in mehreren Formaten:
   - JSONL (JSON Lines) für streaming-optimierte Verarbeitung
   - YAML für menschenlesbare Inspektion und Versionskontrolle

2. **Transformation** der Rohdaten in eine optimierte SQLite-Datenbank für:
   - Schnelle Abfragen und Lookups
   - Verwendung in anderen EVE-bezogenen Projekten
   - Offline-Verfügbarkeit der Spieldaten

## Projektstatus

**v0.1.0** – SQLite-Datenbank & Sync-Automatisierung abgeschlossen.

### Fertiggestellt

- ✅ SDE Download-Mechanismus (JSONL + YAML)
- ✅ Schema-Generierung (53 typsichere Go-Structs)
- ✅ SQLite-Datenbank-Implementierung
  - Reflection-basierter Schema-Generator
  - Streaming JSONL-Importer mit Batch-Processing
  - CLI-Tool: `sde-to-sqlite`
  - Performance: 500k Zeilen in 24s, 405 MB DB
- ✅ Version Tracking System
  - CCP API Integration (Build-Nummer & Release-Datum)
  - Intelligente Update-Erkennung
  - CLI-Flags: `--check-version`, `--skip-if-current`
- ✅ Sync-Pipeline Automatisierung
  - CLI-Tool: `sde-sync` (orchestriert Download → Schema → Import)
  - Makefile Targets: `make sync`, `make sync-force`
  - Version-basiertes Skip-Verhalten
- ✅ Automatisches Release-System (GitHub Actions)
  - Täglicher Cron-Job prüft auf neue SDE-Versionen
  - Erstellt GitHub Release mit gezippter SQLite-DB
  - Tag-Format: `sde-v{buildNumber}-{datum}`
  - Retention: 2 Jahre
  - Keine manuelle Intervention erforderlich

### Nächste Schritte

- [ ] YAML-Import für nested Strukturen
- [ ] Diff/Update Mechanismus (nur Änderungen importieren)
- [ ] Progress Tracking & Verbose Logging

## Struktur

```text
eve-sde/
├── cmd/
│   ├── sde-schema-gen/      # Schema-Generator CLI
│   ├── sde-to-sqlite/       # SQLite-Import CLI
│   └── sde-sync/            # Sync-Pipeline Orchestrator (NEW)
├── internal/
│   ├── schema/
│   │   └── types/           # 53 generierte Go-Structs
│   ├── sqlite/
│   │   ├── schema/          # DDL-Generator
│   │   └── importer/        # JSONL→SQLite Importer
│   └── sde/
│       └── version/         # Version Tracking (NEW)
├── data/                    # Lokale SDE-Kopien (gitignored)
│   ├── jsonl/               # 52 JSONL-Dateien (~499 MB)
│   ├── yaml/                # 52 YAML-Dateien (~160 MB)
│   └── sqlite/              # eve-sde.db (~405 MB)
├── scripts/                 # Sync-, Transform- und Validierungslogik
├── docs/
│   ├── adr/                 # Architekturentscheidungen (ADRs)
│   └── sqlite-implementation.md  # SQLite-Dokumentation
└── .github/copilot-instructions.md  # Engineering-Richtlinien
```

## Getting Started

### Option 1: Fertige SQLite-DB herunterladen (empfohlen)

Die einfachste Methode ist der Download einer vorkompilierten SQLite-Datenbank:

```bash
# Neueste Version anzeigen
gh release list --limit 1

# Download (ersetze TAG mit aktuellem Release)
gh release download sde-v3064089-2025-10-17 -p "eve-sde.db.gz"

# Entpacken
gunzip eve-sde.db.gz

# Beispielabfrage
sqlite3 eve-sde.db "SELECT name FROM types WHERE _key = 34;"
```

**Alternativ:** Manueller Download über [GitHub Releases](https://github.com/Sternrassler/eve-sde/releases)

### Option 2: Lokal bauen

1. Repository clonen:

   ```bash
   git clone https://github.com/Sternrassler/eve-sde.git
   cd eve-sde
   ```

2. Git Hooks aktivieren:

   ```bash
   git config core.hooksPath .githooks
   ```

3. SDE-Daten herunterladen:

   ```bash
   ./scripts/download-sde.sh
   ```

   Dies lädt automatisch die neuesten YAML und JSONL Exporte (~160MB komprimiert) herunter und extrahiert sie nach `data/yaml/` und `data/jsonl/`.

4. Go-Schemas generieren (optional - bereits committed):

   ```bash
   go run ./cmd/sde-schema-gen
   ```

   Analysiert die JSONL-Dateien und generiert typsichere Go-Structs in `internal/schema/types/`.

5. SQLite-Datenbank erstellen:

   ```bash
   # Alle Schemas importieren (41 Tabellen, ~24s)
   go run ./cmd/sde-to-sqlite

   # Nur spezifische Tabelle importieren
   go run ./cmd/sde-to-sqlite --import types

   # Custom DB-Pfad
   go run ./cmd/sde-to-sqlite --db custom/path.db
   ```

   **Performance**: 500.000 Zeilen in 24 Sekunden, 405 MB Datenbank

   Details siehe [docs/sqlite-implementation.md](docs/sqlite-implementation.md)

## Verwendung

### Schema-Generierung

Generiert typsichere Go-Structs aus JSONL-Daten:

```bash
go run ./cmd/sde-schema-gen
```

**Features**:

- Automatische Typerkennung (int64, float64, string, bool)
- LocalizedText-Detection (8 Sprachen: de, en, es, fr, ja, ko, ru, zh)
- Required-Field-Detection (basierend auf NULL-Präsenz)
- 53 generierte Schemas in `internal/schema/types/`

### SQLite-Datenbank

Importiert alle JSONL-Daten in eine SQLite-Datenbank:

```bash
# Vollautomatischer Sync (Version Check → Download → Schema → Import)
make sync

# Erzwinge Update (auch wenn DB aktuell)
make sync-force

# Nur Download + Schema (kein SQLite Import)
make sync-download-only

# Manueller Full Import (alle 41 Schemas)
go run ./cmd/sde-to-sqlite

# Einzeltabelle
go run ./cmd/sde-to-sqlite --import types

# Custom DB-Pfad
go run ./cmd/sde-to-sqlite --db custom/eve.db --jsonl data/jsonl

# Version prüfen
go run ./cmd/sde-to-sqlite --check-version

# Nur bei Update importieren
go run ./cmd/sde-to-sqlite --skip-if-current

# Nur Schema erstellen (ohne Daten)
go run ./cmd/sde-to-sqlite --init
```

**Performance-Metriken**:

| Metrik | Wert |
|--------|------|
| Import-Zeit | 24 Sekunden (41 Tabellen) |
| Datensätze | ~500.000 Zeilen |
| DB-Größe | 405 MB (18% Kompression) |
| Durchsatz | ~20.000 Zeilen/Sekunde |

**Datenvalidierung**:

```bash
# Zeilenzahlen prüfen
sqlite3 data/sqlite/eve-sde.db "SELECT COUNT(*) FROM types;"  # 50,486
wc -l data/jsonl/types.jsonl  # 50,486 ✓

# LocalizedText Beispiel
sqlite3 data/sqlite/eve-sde.db \
  "SELECT name FROM types WHERE _key = 34 LIMIT 1;"
# {"de":"Tritanium","en":"Tritanium","es":"Tritanio",...}
```

Details siehe [docs/sqlite-implementation.md](docs/sqlite-implementation.md)

## Automatische Updates (GitHub Actions)

Dieses Repository nutzt GitHub Actions für automatische SDE-Synchronisation:

- **Zeitplan:** Täglich um 03:00 UTC
- **Trigger:** Bei neuer SDE-Version (BuildNumber-Änderung)
- **Aktion:**
  1. Download aktueller SDE-Daten
  2. Schema-Generierung
  3. SQLite-Import
  4. Erstellung eines GitHub Release
- **Release-Tag:** `sde-v{buildNumber}-{datum}` (z.B. `sde-v3064089-2025-10-17`)
- **Asset:** `eve-sde.db.gz` (gzip-komprimierte SQLite-DB)
- **Retention:** Releases älter als 2 Jahre werden automatisch gelöscht

Alle verfügbaren Versionen: [GitHub Releases](https://github.com/Sternrassler/eve-sde/releases)

### SDE Download

Das Download-Script lädt automatisch die neueste Version der EVE SDE:

```bash
./scripts/download-sde.sh
```

**Hinweis:** Die heruntergeladenen Daten werden in `data/` gespeichert und sind durch `.gitignore` vom Versionskontrollsystem ausgeschlossen.

### Datenformate

- **JSONL** (`data/jsonl/`): 52 Dateien, ~499 MB - JSON Lines Format für Streaming
- **YAML** (`data/yaml/`): 52 Dateien, ~160 MB - Human-readable Format
- **SQLite** (`data/sqlite/eve-sde.db`): 41 Tabellen, ~405 MB - Optimierte Datenbank
  - Primary Keys auf allen `_key` Feldern
  - Indices auf Foreign Keys
  - LocalizedText als JSON TEXT (8 Sprachen)
  - Reflection-basiertes Schema (type-safe)

### Architektur

**Schema-Generator** (`cmd/sde-schema-gen`):

- Analysiert JSONL-Dateien statistisch
- Generiert Go-Structs mit JSON-Tags
- Type Inference: NULL-Handling, int64/float64 Harmonisierung
- LocalizedText-Erkennung über Feldnamen-Pattern

**SQLite-Importer** (`cmd/sde-to-sqlite`):

- DDL-Generator via Reflection (Go struct → CREATE TABLE)
- Streaming JSONL-Parser mit `bufio.Scanner`
- Batch-Inserts (1000 Zeilen/Batch)
- SQLite-Optimierungen: WAL mode, PRAGMA tuning
- Type Conversion: bool→INTEGER, complex→JSON

## Entwicklung

### Build

```bash
# Schema-Generator
go build ./cmd/sde-schema-gen

# SQLite-Importer
go build ./cmd/sde-to-sqlite
```

### Tests

```bash
# Schema-Generator Tests
go test ./internal/sqlite/schema/... -v

# Alle Tests
go test ./... -v
```

### Engineering-Richtlinien

Das Projekt folgt strikten Engineering-Prinzipien:

- **TDD**: Tests vor Implementierung (Red → Green → Refactor)
- **Git-Workflow**: Issue → Branch → PR → Review → Merge
- **Normative Standards**: MUST/SHOULD/MAY nach RFC 2119
- **ADRs**: Architekturentscheidungen dokumentiert in `docs/adr/`
- **Pre-Commit Hooks**: Normative Checks, ADR-Validierung, Secret-Scanning

Details siehe [`.github/copilot-instructions.md`](.github/copilot-instructions.md)

## Lizenz

Dieses Projekt steht unter der [MIT-Lizenz](LICENSE).
