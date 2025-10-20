# SQLite Implementation

Dokumentation der EVE SDE SQLite-Konvertierung.

## Übersicht

Das Projekt implementiert einen vollständigen Pipeline-Workflow:

1. **Schema-Generierung**: Reflection-basierte Go-Struct-Analyse → SQLite DDL
2. **JSONL→SQLite Import**: Streaming-Importer mit Transaktionen & Batch-Processing
3. **CLI-Tool**: `sde-to-sqlite` mit Flags für Init, Import, Selektion

## Performance-Metriken (v0.1.0)

### Full Import (alle 41 Schemas)

- **Laufzeit**: 24.2 Sekunden
- **Datenmenge**: 499 MB JSONL → 405 MB SQLite
- **Datensätze**: ~500k Zeilen gesamt
- **Kompression**: ~18% (405 MB / 499 MB)
- **Durchsatz**: ~20k Zeilen/Sekunde

### Einzelne Imports (Beispiele)

| Schema | Zeilen | Import-Zeit | Rows/Sec |
|--------|--------|-------------|----------|
| types | 50,486 | 4.0s | 12,621 |
| mapMoons | 342,170 | 13.0s | 26,321 |
| mapPlanets | 67,961 | 3.0s | 22,653 |
| agentTypes | 13 | <0.1s | - |

## Validierung

### Datenintegrität

Alle JSONL Zeilenzahlen stimmen exakt mit SQLite Zeilenzahlen überein:

```bash
# Beispiel
agentTypes:       13 =   13 ✓
types:        50,486 = 50,486 ✓
mapMoons:    342,170 = 342,170 ✓
```

### LocalizedText Speicherung

Korrekte JSON-Serialisierung aller 8 Sprachen (de, en, es, fr, ja, ko, ru, zh):

```json
{"de":"Tritanium","en":"Tritanium","es":"Tritanio",...}
```

### Schema-Struktur

- 41 Tabellen erstellt
- 405 MB Datenbankgröße
- Primary Keys auf `_key` Feldern
- Indices auf Foreign Keys (validiert)

## Technische Details

### SQLite-Optimierungen (PRAGMAs)

```go
PRAGMA journal_mode=WAL        // Concurrent reads
PRAGMA synchronous=NORMAL      // Performance
PRAGMA cache_size=10000        // 10MB Cache
PRAGMA temp_store=MEMORY       // In-Memory Temp
```

### Batch-Verarbeitung

- 1000 Zeilen pro Batch-Insert
- Einzelne Transaktion pro JSONL-Datei
- Streaming Parser (bufio.Scanner)

### Type Mapping

| Go Type | SQLite Type | Beispiel |
|---------|-------------|----------|
| int64 | INTEGER | _key, IDs |
| float64 | REAL | mass, volume |
| bool | INTEGER | published (0/1) |
| string | TEXT | Strings |
| LocalizedText | TEXT | JSON mit 8 Sprachen |
| struct/map/slice | TEXT | JSON-encoded |

### Index-Validierung

Nur Felder die im Go-Struct existieren werden indiziert:

```go
// Beispiel: agentTypes hat kein factionID Feld
// → Index wird stillschweigend übersprungen
```

## CLI Usage

### Kompletter Import (alle Schemas)

```bash
sde-to-sqlite
```

### Einzeltabelle importieren

```bash
sde-to-sqlite --import types
```

### Nur Schema erstellen (ohne Daten)

```bash
sde-to-sqlite --init
```

### Custom DB-Pfad

```bash
sde-to-sqlite --db custom/path.db --jsonl data/jsonl
```

## Nächste Schritte

- [ ] Sync-Mechanismus (Download → Schema-Gen → Import)
- [ ] Progress Tracking (z. B. mit progressbar)
- [ ] Verbose/Debug Logging
- [ ] YAML-Import für nested Strukturen
- [ ] Diff/Update Mechanismus (nur Änderungen)
- [ ] Benchmark-Suite
- [ ] Integration Tests

## Commits

- Initial SQLite Implementation (Schema Generator + Importer + CLI)
- Validation: All 41 schemas, 500k rows, 24s import time
