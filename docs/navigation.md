# EVE Navigation & Intelligence System

Dieses Dokument beschreibt die Navigation- und Intelligence-Features für EVE Online SDE SQLite-Datenbank.

## Überblick

Das Navigation-System bietet:
- **Pathfinding**: Kürzeste Routen zwischen Systemen (Go-basierter Dijkstra-Algorithmus)
- **Travel Time Calculation**: Reisezeit-Berechnung mit Schiffs-Parametern
- **Security Filtering**: Vermeidung von Low-Sec/Null-Sec Systemen
- **Trade Hub Analysis**: Distanz zu Major Trade Hubs
- **Region Intelligence**: Security-Zonen und Region-Statistiken

## SQL Views

> **Hinweis**: Alle Views werden automatisch bei jedem DB-Import neu erstellt via `sde-to-sqlite`. 
> Siehe [Persistence & Recreation](#view-persistence) am Ende dieses Dokuments.

### v_stargate_graph
Bidirektionaler Stargate-Connectivity-Graph für Pathfinding.

```sql
SELECT * FROM v_stargate_graph LIMIT 5;
```

**Columns:**
- `from_system_id`: Quellsystem ID
- `to_system_id`: Zielsystem ID
- `gate_id`: Stargate ID
- `gate_type_id`: Stargate Typ ID

### v_system_info
Enhanced System-Information mit parsed Namen und Security-Zonen.

```sql
SELECT * FROM v_system_info WHERE system_name LIKE 'Jita%';
```

**Columns:**
- `system_id`: System ID (primary key)
- `solar_system_id`: Solar System ID
- `system_name`: Name (Englisch, Deutsch als Fallback)
- `security_status`: Security Status (0.0 - 1.0)
- `security_zone`: 'High-Sec', 'Low-Sec', 'Null-Sec', 'Wormhole'
- `constellation_id`: Constellation ID
- `region_id`: Region ID
- `region_name`: Region Name
- `constellation_name`: Constellation Name
- `border`, `corridor`, `hub`: Boolean Flags
- `wormhole_class_id`: Wormhole Class (NULL für K-Space)

### v_system_security_zones
Security-Zonen-Statistiken pro Region.

```sql
SELECT * FROM v_system_security_zones 
WHERE region_name = 'The Forge';
```

**Columns:**
- `region_id`: Region ID
- `region_name`: Region Name
- `security_zone`: Security Zone
- `system_count`: Anzahl Systeme
- `avg_security`: Durchschnittliche Security

### v_region_stats
Comprehensive Region-Statistiken.

```sql
SELECT * FROM v_region_stats 
ORDER BY total_systems DESC 
LIMIT 10;
```

**Columns:**
- `region_id`: Region ID
- `region_name`: Region Name
- `total_systems`: Gesamtzahl Systeme
- `avg_security`: Durchschnittliche Security
- `border_systems`: Border System Count
- `high_sec_count`, `low_sec_count`, `null_sec_count`, `wormhole_count`: Counts

### v_trade_hubs
Major Trade Hub Information.

```sql
SELECT * FROM v_trade_hubs;
```

**Trade Hubs:**
- **Jita** (30000142) - The Forge - Höchster Traffic
- **Amarr** (30002187) - Domain
- **Dodixie** (30002659) - Sinq Laison
- **Rens** (30002510) - Heimatar
- **Hek** (30002053) - Metropolis

## Go API

**Architektur:** Die Navigation-API ist in zwei Ebenen getrennt:
- **DB-Core** (`internal/sqlite/views`): SQL View Initialisierung
- **API Layer** (`pkg/evedb/navigation`): Go-basierte Navigation-Funktionen

**Setup (einmalig pro DB-Verbindung):**

```go
import (
    "github.com/Sternrassler/eve-sde/pkg/evedb/navigation"
    "github.com/Sternrassler/eve-sde/internal/sqlite/views"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

db, err := sql.Open("sqlite3", "data/sqlite/eve-sde.db")
if err != nil {
    log.Fatal(err)
}

// Views müssen einmalig initialisiert werden
if err := views.InitializeNavigationViews(db); err != nil {
    log.Fatal(err)
}
```

### Shortest Path (Pathfinding)

```go
// Find shortest path
path, err := navigation.ShortestPath(db, 30000142, 30002187, false)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Route: %d jumps\n", path.Jumps)
fmt.Printf("Systems: %v\n", path.Route)
```

**Parameters:**
- `db`: SQLite database connection
- `fromSystemID`: Source system ID
- `toSystemID`: Destination system ID
- `avoidLowSec`: Skip systems with security < 0.45

**Returns:**
```go
type PathResult struct {
    FromSystemID int64   `json:"from_system_id"`
    ToSystemID   int64   `json:"to_system_id"`
    Jumps        int     `json:"jumps"`
    Route        []int64 `json:"route"`
}
```

### Travel Time Calculation

#### Default Parameters (Cruiser)

```go
// Use default ship parameters (Cruiser-like)
result, err := navigation.CalculateTravelTime(db, 30000142, 30002187, nil)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Total time: %.1f minutes (%d jumps)\n", 
    result.TotalMinutes, result.Jumps)
```

#### Custom Ship Parameters

```go
// Interceptor (fast align, fast warp)
params := &navigation.NavigationParams{
    WarpSpeed: ptrFloat64(6.0),
    AlignTime: ptrFloat64(2.5),
}

result, err := navigation.CalculateTravelTime(db, 30000142, 30002187, params)

// Freighter (slow align, slow warp)
params := &navigation.NavigationParams{
    WarpSpeed: ptrFloat64(1.36),
    AlignTime: ptrFloat64(30.0),
}

result, err := navigation.CalculateTravelTime(db, 30000142, 30002187, params)
```

#### Calculated Align Time from Ship Stats

```go
// Let navigation calculate align time from ship mass + inertia
params := &navigation.NavigationParams{
    ShipMass:        ptrFloat64(12000000),    // 12M kg (Cruiser)
    InertiaModifier: ptrFloat64(0.4),         // Agility
    WarpSpeed:       ptrFloat64(3.0),
}

result, err := navigation.CalculateTravelTime(db, 30000142, 30002187, params)

// result.ParametersUsed["source"] will be "calculated"
```

#### Security Filtering

```go
// Avoid Low-Sec and Null-Sec systems (High-Sec only)
params := &navigation.NavigationParams{
    AvoidLowSec: true,
}

result, err := navigation.CalculateTravelTime(db, 30000142, 30002187, params)

// May result in longer route (more jumps) but safer
```

### Exact Warp Calculation (CCP 3-Phase Formula)

```go
// Use exact CCP warp formula instead of simplified
result, err := navigation.CalculateTravelTimeExact(db, 30000142, 30002187, params)

// More accurate but slightly slower computation
```

### Navigation Parameters

```go
type NavigationParams struct {
    WarpSpeed       *float64 `json:"warp_speed,omitempty"`        // AU/s (default: 3.0)
    AlignTime       *float64 `json:"align_time,omitempty"`        // seconds (default: 6.0)
    ShipMass        *float64 `json:"ship_mass,omitempty"`         // kg (for calculated align)
    InertiaModifier *float64 `json:"inertia_modifier,omitempty"`  // ship agility
    AvgWarpDistance *float64 `json:"avg_warp_distance,omitempty"` // AU per system (default: 15)
    AvoidLowSec     bool     `json:"avoid_lowsec"`                // route via high-sec only
}
```

**Defaults:**
- **Warp Speed**: 3.0 AU/s (Cruiser average)
- **Align Time**: 6.0 seconds (medium ships)
- **Gate Jump**: 10.0 seconds (animation + loading)
- **Avg Warp Distance**: 15 AU (statistical assumption)

### Ship Parameter Reference

| Ship Class | Warp Speed (AU/s) | Align Time (s) | Use Case |
|------------|-------------------|----------------|----------|
| **Interceptor** | 6.0-8.0 | 2-3 | Schnelle Recon |
| **Cruiser** | 3.0-4.5 | 5-8 | Standard Travel |
| **Battleship** | 1.5-2.0 | 10-15 | Combat Fleets |
| **Freighter** | 1.36 | 25-40 | Hauling (sehr langsam) |
| **Jump Freighter** | 1.36 | 20-30 | Logistics + Jump Drive |
| **Blockade Runner** | 3.0 | 3-4 | Covert Ops Hauling |

## Formulas

### Align Time (Exact)

```
align_time = -ln(0.25) * inertia_modifier * mass / 500000
           ≈ 1.386 * inertia_modifier * mass / 500000
```

**Go Implementation:**
```go
alignTime := navigation.CalculateAlignTime(mass, inertiaModifier)
```

### Warp Time (CCP 3-Phase)

```
k = warp_speed (AU/s)
j = min(k/3, 2)
AU = 149,597,870,700 meters

# Phase 1: Acceleration (always 1 AU)
t_accel = 25.7312 / k

# Phase 2: Deceleration
d_decel = (k * AU) / j
t_decel = (ln(k*AU) - ln(100)) / j

# Phase 3: Cruise
d_cruise = total_distance - 1*AU - d_decel
t_cruise = d_cruise / (k * AU)

t_total = t_accel + t_cruise + t_decel
```

**Go Implementation:**
```go
warpTime := navigation.CalculateWarpTime(distanceAU, warpSpeedAU)
```

### Warp Time (Simplified)

```
time_warp = (distance_AU / warp_speed) * correction_factor
          = (distance_AU / warp_speed) * 1.4
```

**Go Implementation:**
```go
warpTime := navigation.CalculateSimplifiedWarpTime(distanceAU, warpSpeedAU)
```

### Total Jump Time

```
time_per_jump = align_time + warp_time + gate_jump_delay
              = align_time + (avg_warp_distance / warp_speed) * 1.4 + 10s

# Default (Cruiser)
              = 6s + (15 AU / 3.0 AU/s) * 1.4 + 10s
              = 6s + 7s + 10s = 23s
```

## SQL Examples

### Graph Query (Direct)

```sql
-- Query the stargate graph directly
SELECT from_system_id, to_system_id, gate_id
FROM v_stargate_graph
WHERE from_system_id = 30000142
LIMIT 10;
```

**Note**: Pathfinding is now implemented in Go (Dijkstra algorithm) for optimal performance. 
The Go API automatically loads the graph from `v_stargate_graph` and uses an efficient 
in-memory algorithm that provides sub-millisecond pathfinding even for long routes (40+ jumps).

The previous recursive CTE approach has been replaced for better performance - from >5 minutes 
to <1ms for long-distance routes.

### Region Analysis

```sql
-- Find most dangerous regions (high null-sec count)
SELECT 
    region_name,
    total_systems,
    null_sec_count,
    ROUND(100.0 * null_sec_count / total_systems, 1) as null_sec_percent
FROM v_region_stats
WHERE null_sec_count > 0
ORDER BY null_sec_percent DESC
LIMIT 10;
```

### Border Systems (Choke Points)

```sql
-- Find border systems (potential choke points)
SELECT 
    system_name,
    region_name,
    security_status,
    security_zone
FROM v_system_info
WHERE border = 1
ORDER BY security_status ASC
LIMIT 20;
```

## Performance

### Pathfinding
- **Short Routes (<10 jumps)**: < 0.02ms (17μs)
- **Medium Routes (~30 jumps)**: < 0.2ms (170μs)
- **Long Routes (40-50 jumps)**: < 0.3ms (275μs)
- **Algorithm**: Go-based Dijkstra with in-memory graph (O(E + V log V))

### Implementation
- Graph loaded from `v_stargate_graph` view (~40k edges)
- Early termination when goal is reached
- Memory-efficient path reconstruction
- Security filtering applied during graph load

### Views
- **v_stargate_graph**: ~40k bidirectional edges (instant)
- **v_system_info**: ~8k systems (instant)
- **v_region_stats**: ~100 regions (instant)

## Testing

```bash
# Run navigation tests
go test ./internal/sqlite/navigation -v

# Benchmark pathfinding
go test ./internal/sqlite/navigation -bench=.
```

## Known Routes (Validation)

| Route | Expected Jumps | Notes |
|-------|---------------|-------|
| Jita → Amarr | ~40 | Cross-empire route |
| Jita → Dodixie | ~30 | Via Lonetrek/Essence |
| Amarr → Rens | ~45 | Long cross-empire |
| Hek → Jita | ~25 | Northern route |

## Limitations

### Current Implementation
- No Jump Bridges / Ansiblex Gates (player-owned)
- No Wormhole connections (dynamic)
- No Pochven / Filament jumps (special mechanics)
- Static stargate data only (no real-time changes)

### Future Enhancements
- Choke Point Detection (single-gate bottlenecks)
- Risk Scoring (zkillboard API integration)
- Capital Jump Range Calculations (cynos)
- Wormhole Mapping (external data source)
- Web API / REST Endpoints
- Interactive D3.js Map Visualization

## View Persistence

### Automatische Recreation bei DB-Import

Navigation-Views sind **nicht manuell zu pflegen**. Sie werden automatisch neu erstellt:

1. **Bei jedem `sde-to-sqlite` Import**:
   - `InitializeNavigationViews()` wird nach map data import aufgerufen
   - `views.sql` nutzt `CREATE VIEW IF NOT EXISTS` (idempotent)

2. **GitHub Actions Workflow** (`sync-sde-release.yml`):
   - Täglicher Cron-Job prüft auf neue SDE-Versionen
   - `make sync-force` löscht alte DB und importiert neu
   - Views werden automatisch mit recreated → **kein manueller Eingriff nötig**

3. **Lokale Entwicklung**:
   ```bash
   # Views manuell neu erstellen (falls nötig)
   sqlite3 data/sqlite/eve-sde.db < internal/sqlite/navigation/views.sql
   
   # Oder: Force-Import
   make sync-force
   ```

### Zukünftige Custom Functions (Hinweis für Entwickler)

Falls später **Go-basierte SQLite Functions** (via `RegisterFunc`) hinzugefügt werden:
- **Limitation**: Go Functions sind runtime-only, nicht in DB gespeichert
- **Lösung**: Functions müssen bei jedem DB-Open registriert werden
- **Pattern**: `navigation.RegisterFunctions(db)` beim Öffnen aufrufen

**Aktuell**: Alle Features sind SQL-basiert (Views + Recursive CTEs) → keine Runtime-Registrierung nötig.

## References

- **EVE Uni Wiki - Warp Time**: https://wiki.eveuniversity.org/Warp_time_calculation
- **EVE Uni Wiki - Stargates**: https://wiki.eveuniversity.org/Stargates
- **CCP Warp Drive Active**: https://www.eveonline.com/news/view/warp-drive-active
- **SQLite Recursive CTE**: https://www.sqlite.org/lang_with.html
- **SQLite JSON Functions**: https://www.sqlite.org/json1.html
