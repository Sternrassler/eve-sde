# EVE Online SDE Database

Produktionsfertiges SQLite-Datenbank-System für EVE Online Static Data Export (SDE). Automatische tägliche Updates via GitHub Actions.

## Quick Start

```bash
# Download neueste SQLite-Datenbank (empfohlen)
gh release download --pattern "eve-sde.db.gz"
gunzip eve-sde.db.gz

# Beispiel-Abfrage
sqlite3 eve-sde.db "SELECT name FROM types WHERE _key = 34;"
# {"de":"Tritanium","en":"Tritanium",...}
```

**Alternative:** [Manueller Download](https://github.com/Sternrassler/eve-sde/releases) (neueste Release)

## Features

### Core Database (v0.2.0)

- **SQLite-Datenbank**: 405 MB, 41 Tabellen, ~500k Zeilen
- **Auto-Updates**: Täglich um 03:00 UTC via GitHub Actions
- **Performance**: 20k Zeilen/Sekunde Import, Sub-Millisekunden Pathfinding
- **Offline-fähig**: Keine Runtime-Abhängigkeiten außer SQLite

### Navigation System

```bash
# Hinweis: Navigation API wurde nach eve-o-provit migriert
# Siehe: github.com/Sternrassler/eve-o-provit/backend/pkg/evedb/navigation
```

**Features:**

- Dijkstra Pathfinding (11,500 Stargates, 5,700 Systeme)
- Travel Time Berechnung mit Schiffs-Parametern
- Security Filtering (High-Sec only Routes)
- Trade Hub Analysis (Jita, Amarr, Dodixie, Rens, Hek)

Details: [docs/navigation.md](docs/navigation.md) (Legacy-Dokumentation)

### Cargo & Hauling API (MIGRIERT)

**Hinweis:** Die Cargo API wurde nach **eve-o-provit** migriert.

Siehe: [github.com/Sternrassler/eve-o-provit](https://github.com/Sternrassler/eve-o-provit/tree/main/backend/pkg/evedb/cargo)

**Features:**

- **Item Volumes**: Volumen-Informationen und ISK/m³ Value-Density
- **Ship Capacities**: Cargo-Holds (Cargo, Ore Hold, Fleet Hangar)
- **Skill System**: Racial Hauler, Freighter Skills (optional)
- **Cargo Fit Calculation**: Wieviele Items passen in Schiff?
- **Route Security**: System-Security-Analyse für sichere Routen

Legacy-Dokumentation: [docs/cargo-api.md](docs/cargo-api.md)

## Architektur

**DB-First Philosophy:** SQLite-Datenbank ist primäres Produkt, Go-Code ist Build-Tool.

```text
eve-sde/
├── cmd/                     # Build-Tools (lokal)
│   ├── sde-to-sqlite/       # DB Import (JSONL → SQLite)
│   └── sde-sync/            # Download & Sync Orchestrator
├── internal/                # DB-Core Implementation
│   ├── sqlite/
│   │   ├── schema/          # DDL Generator
│   │   ├── importer/        # JSONL Streaming Importer
│   │   └── views/           # SQL Views (Navigation, Cargo, Stats)
│   └── schema/types/        # 53 Go Structs (generiert)
├── data/                    # Lokale Daten (gitignored)
│   └── sqlite/eve-sde.db    # **HAUPTPRODUKT**
└── docs/                    # Technische Dokumentation
    ├── adr/                 # Architektur-Entscheidungen
    ├── navigation.md        # Navigation System Docs (Legacy)
    └── cargo-api.md         # Cargo API Docs (Legacy)
```

**Nutzung:**

- **Direkt:** SQLite-DB via `sqlite3` CLI oder Bibliotheken (Python, Node.js, etc.)
- **Go API:** Migriert nach [eve-o-provit](https://github.com/Sternrassler/eve-o-provit) (Navigation & Cargo)

**Hinweis:** Die `pkg/evedb/` Go-APIs wurden nach `eve-o-provit/backend/pkg/evedb/` migriert.
Dieses Repository fokussiert sich auf die SQLite-Datenbank-Generierung.

## Lokaler Build (Optional)

```bash
# Repository clonen
git clone https://github.com/Sternrassler/eve-sde.git
cd eve-sde

# Automatischer Sync (Download → Import → Views)
make sync

# Datenbank-Pfad
ls -lh data/sqlite/eve-sde.db  # 405 MB
```

**Makefile Targets:**

- `make sync` - Vollautomatischer Download & Import
- `make sync-force` - Erzwinge Update (löscht alte DB)
- `make test` - Go Tests ausführen

## Datenbank-Schema

**41 Tabellen:**

- `types` (50k Zeilen) - Items, Ships, Modules
- `mapSolarSystems` (5.7k) - Systeme mit Security Status
- `mapStargates` (11.5k) - Stargate-Verbindungen
- `groups`, `categories`, `regions`, `constellations`, ...

**7 SQL Views:**

- `v_stargate_graph` - Pathfinding Graph
- `v_system_info` - System-Metadaten
- `v_item_volumes`, `v_ship_cargo_capacities` - Cargo Calculations
- `v_trade_hubs` - Major Trade Hubs (Jita, Amarr, etc.)
- `v_region_stats`, `v_system_security_zones` - Region Intelligence

```sql
-- Beispiel: Tritanium Details
SELECT name, volume, basePrice 
FROM types 
WHERE _key = 34;
```

Details: [docs/sqlite-implementation.md](docs/sqlite-implementation.md)

## Automatische Updates

**GitHub Actions** (`sync-sde-release.yml`):

- **Zeitplan:** Täglich 03:00 UTC
- **Trigger:** Neue SDE BuildNumber von CCP API
- **Output:** GitHub Release mit `eve-sde.db.gz`
- **Retention:** 2 Jahre

Alle Releases: [github.com/Sternrassler/eve-sde/releases](https://github.com/Sternrassler/eve-sde/releases)

## Entwicklung

```bash
# Pre-Commit Hooks aktivieren
git config core.hooksPath .githooks

# Tests
go test ./...

# Lokaler Build
go build ./cmd/sde-to-sqlite
```

**Engineering-Richtlinien:**

- TDD (Test-Driven Development)
- ADRs für Architektur-Entscheidungen
- Normative Standards (MUST/SHOULD/MAY)

Details: [.github/copilot-instructions.md](.github/copilot-instructions.md)

## Lizenz

[MIT License](LICENSE) - Freie Nutzung für alle EVE-bezogenen Projekte.

---

**Hinweis:** Dieses Projekt ist nicht offiziell von CCP Games endorsed. EVE Online und alle zugehörigen Logos sind Marken von CCP hf.
