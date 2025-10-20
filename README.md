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

**In Entwicklung** – Initiale Projektstruktur und Governance etabliert.

Nächste Schritte:

- [ ] SDE Download-Mechanismus implementieren
- [ ] JSONL/YAML Konverter entwickeln
- [ ] SQLite Schema Design und Migrations-Framework
- [ ] Sync-Automatisierung (periodische Updates)

## Struktur

- `data/` – Lokale SDE-Kopien (JSONL, YAML, SQLite)
- `scripts/` – Sync-, Transform- und Validierungslogik
- `docs/adr/` – Architekturentscheidungen (ADRs)
- `.github/copilot-instructions.md` – Engineering-Richtlinien

## Getting Started

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
   ./scripts/fetch-schemas.sh --refresh
   ```

   Analysiert die JSONL-Dateien und generiert typsichere Go-Structs in `internal/schema/types/`.

## Verwendung

### SDE Download

Das Download-Script lädt automatisch die neueste Version der EVE SDE:

```bash
./scripts/download-sde.sh
```

**Hinweis:** Die heruntergeladenen Daten werden in `data/` gespeichert und sind durch `.gitignore` vom Versionskontrollsystem ausgeschlossen.

### Datenformate

- **JSONL** (`data/jsonl/`): JSON Lines Format, ideal für Streaming und große Datasets
- **YAML** (`data/yaml/`): Human-readable Format für Inspektion und Versionierung
- **SQLite** (`data/sqlite/`): (Geplant) Optimierte Datenbank für Abfragen

## Lizenz

Dieses Projekt steht unter der [MIT-Lizenz](LICENSE).
