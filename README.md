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
# Go-API Beispiel (optional)
go run examples/navigation/main.go
# Jita → Amarr: 40 jumps in 273µs
```

**Features:**

- Dijkstra Pathfinding (11,500 Stargates, 5,700 Systeme)
- Travel Time Berechnung mit Schiffs-Parametern
- Security Filtering (High-Sec only Routes)
- Trade Hub Analysis (Jita, Amarr, Dodixie, Rens, Hek)

Details: [docs/navigation.md](docs/navigation.md)

### Cargo & Hauling API (NEU in v0.2.0)

Vollständiges Cargo-Berechnungssystem für Hauling und Trade:

- **Item Volumes**: Volumen-Informationen und ISK/m³ Value-Density
- **Ship Capacities**: Cargo-Holds (Cargo, Ore Hold, Fleet Hangar)
- **Skill System**: Racial Hauler, Freighter, Mining Barge Skills (optional)
- **Cargo Fit Calculation**: Wieviele Items passen in Schiff?
- **Route Security**: System-Security-Analyse für sichere Routen

**Beispiel:**

```bash
# Basis-Berechnung (ohne Skills)
go run examples/cargo/main.go --ship 648 --item 34

# Mit Gallente Hauler V (+25% Cargo)
go run examples/cargo/main.go --ship 648 --item 34 --racial-hauler 5

# Schiffs-Kapazitäten anzeigen
go run examples/cargo/main.go --ship 648 --ship-info
```

**Go API:**

```go
import "github.com/Sternrassler/eve-sde/pkg/evedb/cargo"

// Ohne Skills (Basis-Werte)
result, _ := cargo.CalculateCargoFit(db, 648, 34, nil)
// Mit Skills
racialLevel := 5
skills := &cargo.SkillModifiers{RacialHaulerLevel: &racialLevel}
result, _ := cargo.CalculateCargoFit(db, 648, 34, skills)
fmt.Printf("Effective: %.0f m³ (+%.0f%%)\n", 
    result.EffectiveCapacity, result.SkillBonus)
```

Siehe [docs/cargo-api.md](docs/cargo-api.md) für vollständige API-Dokumentation.

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
├── pkg/evedb/               # Optional: Go APIs
│   ├── navigation/          # Navigation API (Pathfinding, Travel Time)
│   └── cargo/               # Cargo API (Hauling, Capacity Calculations)
├── data/                    # Lokale Daten (gitignored)
│   └── sqlite/eve-sde.db    # **HAUPTPRODUKT**
└── docs/                    # Technische Dokumentation
    ├── adr/                 # Architektur-Entscheidungen
    ├── navigation.md        # Navigation System Docs
    └── cargo-api.md         # Cargo & Hauling API Docs
```

**Nutzung:**

- **Direkt:** SQLite-DB via `sqlite3` CLI oder Bibliotheken (Python, Node.js, etc.)
- **Go API:** Optional via `pkg/evedb/navigation` oder `pkg/evedb/cargo` (Convenience Layer)

Details: [docs/adr/ADR-001-db-core-api-separation.md](docs/adr/ADR-001-db-core-api-separation.md)

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
