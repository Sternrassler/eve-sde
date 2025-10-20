# ADR-001: Separation of DB-Core and API Layer

## Status Accepted

## Kontext

Das eve-sde Projekt wurde ursprünglich mit dem Fokus auf die **SQLite-Datenbank als Kernprodukt** entwickelt. Das primäre Ziel ist:

- Synchronisation der EVE Online SDE von CCP
- Transformation in eine performante SQLite-Datenbank
- Bereitstellung strukturierter, typsicherer Daten für andere Projekte

Mit der Implementierung des Navigation & Intelligence Systems (PR #2) entstand eine **Go API** (`internal/sqlite/navigation/`), die:

- SQL Views für Routing nutzt
- Go-Funktionen für Berechnungen bereitstellt (Warp Time, Align Time)
- Pathfinding-Logik implementiert

### Problem

Die Go API lag innerhalb von `internal/sqlite/` **neben** den DB-Core Komponenten (`importer/`, `schema/`), was zu architektonischen Unklarheiten führte:

1. **Vermischte Verantwortlichkeiten**: DB-Erstellung (Core) vs. DB-Nutzung (API)
2. **Unklare Abhängigkeitsrichtung**: Sollte der Importer die API kennen?
3. **Erschwerte Erweiterbarkeit**: Wo gehören zukünftige APIs hin? (Market, Industry, etc.)
4. **Internal-Barriere**: `internal/` verhindert externe Nutzung der API

### Anforderung

> **Projektphilosophie**: Die SQLite-Datenbank ist das Kernprodukt. Go APIs sind **optionale Goodies** für Convenience, nicht Teil der DB-Pipeline.

## Entscheidung

Wir trennen **DB-Core** und **API Layer** durch eine klare Package-Architektur:

### Neue Struktur

```
eve-sde/
├── internal/sqlite/              # DB-CORE (Daten-Pipeline)
│   ├── importer/                 # JSONL → SQLite Import
│   ├── schema/                   # DDL Generator
│   └── views/                    # SQL View Definitionen (NEW)
│       ├── navigation.sql        # Pure SQL Views
│       └── init.go               # View Initialisierung
│
├── pkg/evedb/                    # PUBLIC API LAYER (NEW)
│   └── navigation/               # Navigation & Route Planning API
│       ├── navigation.go         # Pathfinding, TravelTime
│       ├── navigation_test.go
│       └── integration_test.go
│
├── cmd/                          # CLI Tools (DB-fokussiert)
│   ├── sde-sync/                 # Nutzt nur internal/sqlite/*
│   └── sde-to-sqlite/            # Nutzt nur internal/sqlite/*
│
└── examples/                     # API-Nutzungsbeispiele
    └── navigation/               # Nutzt pkg/evedb/navigation
```

### Verantwortlichkeiten

#### `internal/sqlite/views/` (DB-Core)

- **Zweck**: SQL View Definitionen für die Datenbank
- **Inhalt**: Pure SQL (`*.sql` Files) + Initialisierungs-Code
- **Keine**: Go Business Logic, Berechnungen, komplexe Algorithmen
- **Nutzer**: `cmd/sde-to-sqlite`, DB-Import Pipeline

#### `pkg/evedb/navigation/` (API Layer)

- **Zweck**: High-level Go API für Navigation & Route Planning
- **Abhängigkeit**: Benötigt Views aus `internal/sqlite/views`
- **Inhalt**: Pathfinding, Travel Time Calculations, Parameter Handling
- **Nutzer**: Externe Projekte, `examples/`, optionale Consumer

### Import-Hierarchie

```
cmd/sde-to-sqlite
  └─> internal/sqlite/views       (View Init)
  └─> internal/sqlite/importer    (DB Import)
  └─> internal/sqlite/schema      (DDL Gen)

examples/navigation
  └─> pkg/evedb/navigation         (API)
  └─> internal/sqlite/views        (View Init, falls nötig)

pkg/evedb/navigation
  └─> database/sql                 (Standard Library)
  └─> (keine internen Dependencies)
```

### Migrations-Schritte

1. **Views extrahieren**: `internal/sqlite/navigation/views.sql` → `internal/sqlite/views/navigation.sql`
2. **Init-Code verschieben**: `InitializeViews()` → `internal/sqlite/views/init.go`
3. **API verschieben**: `internal/sqlite/navigation/*.go` → `pkg/evedb/navigation/`
4. **Imports anpassen**:
   - `cmd/sde-to-sqlite`: Nutzt `internal/sqlite/views`
   - `examples/`: Nutzt `pkg/evedb/navigation`
5. **Cleanup**: `internal/sqlite/navigation/` löschen

## Konsequenzen

### Positive

✅ **Klare Separation of Concerns**

- DB-Core (internal/sqlite): Datenbank-Erstellung & -Pflege
- API Layer (pkg/evedb): Datenbank-Nutzung & Convenience

✅ **Externe Nutzbarkeit**

- `pkg/` ist nicht `internal/` → andere Projekte können importieren
- API kann unabhängig vom DB-Core entwickelt werden

✅ **Skalierbarkeit**

- Weitere APIs einfach hinzufügbar: `pkg/evedb/market/`, `pkg/evedb/industry/`
- Keine Verschmutzung des DB-Core

✅ **Unabhängige Versionierung**

- DB-Schema (internal) und API (pkg) können getrennt versioniert werden
- Breaking Changes in API betreffen nicht DB-Core

✅ **Klarere Abhängigkeiten**

- Import-Pipeline kennt keine API-Layer
- API-Layer ist reiner Consumer der DB

### Negative / Trade-offs

⚠ **Leichte Duplikation**

- `internal/sqlite/views/init.go` und API-Layer Code getrennt
- Nicht kritisch, da Views reine SQL-Definitionen sind

⚠ **Zusätzliche Package-Ebene**

- Mehr Verzeichnisse, aber klarer strukturiert
- Überschaubar bei konsequenter Trennung

⚠ **Migration Effort**

- Einmalige Arbeit für bestehenden Code
- Tests müssen angepasst werden

### Risiken & Mitigationen

**Risiko**: Views und API geraten aus dem Sync

- **Mitigation**: Integration Tests in `pkg/evedb/navigation/` validieren View-Existenz
- **Mitigation**: Dokumentation in `pkg/evedb/navigation/` referenziert benötigte Views

**Risiko**: Externe Nutzer brechen API durch direkte DB-Änderungen

- **Mitigation**: Klare Dokumentation: Views sind Teil der DB, API ist opt-in Layer
- **Mitigation**: Semantic Versioning für `pkg/evedb`

## Alternatives Considered

### Alternative 1: Alles in `internal/`

- **Pro**: Keine externe API-Nutzung möglich (mehr Kontrolle)
- **Con**: Verhindert legitime externe Nutzung
- **Rejected**: Widerspricht Open-Source-Philosophie

### Alternative 2: API als separates Modul

- **Pro**: Völlig unabhängige Repositories
- **Con**: Versionierung komplizierter, mehr Overhead
- **Rejected**: Overkill für aktuelles Projekt-Scope

### Alternative 3: Alles in `pkg/`

- **Pro**: Einfacher, flachere Hierarchie
- **Con**: Keine Trennung von DB-Core und API, unklare Verantwortlichkeiten
- **Rejected**: Problem nicht gelöst

## References

- Go Project Layout: <https://github.com/golang-standards/project-layout>
- Go Best Practices: `internal/` vs `pkg/` vs `cmd/`
- Eve SDE Issue #1: feat: Navigation & Intelligence System
- Eve SDE PR #2: Navigation System Implementation

## Notes

Diese Architektur ermöglicht zukünftige Erweiterungen:

- `pkg/evedb/market/` - Market Data API
- `pkg/evedb/industry/` - Industry & Manufacturing API
- `pkg/evedb/universe/` - Universe Data API

Alle folgen dem gleichen Muster: **DB-Core erstellt Views → API Layer nutzt Views**.
