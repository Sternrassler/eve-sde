# EVE SDE Navigation Examples

This directory contains example programs demonstrating the EVE Navigation & Intelligence System.

## Navigation Example

The navigation example demonstrates pathfinding and travel time calculation.

### Usage

```bash
# Basic usage (Jita → Amarr with default Cruiser parameters)
go run examples/navigation_example.go

# Custom route
go run examples/navigation_example.go -from 30000142 -to 30002659

# Custom ship parameters (Interceptor)
go run examples/navigation_example.go -warp 6.0 -align 2.5

# High-sec only route
go run examples/navigation_example.go -safe

# Use exact CCP warp formula
go run examples/navigation_example.go -exact

# Initialize views only
go run examples/navigation_example.go -init-views

# Custom database path
go run examples/navigation_example.go -db /path/to/eve-sde.db
```

### Available Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `-db` | string | `data/sqlite/eve-sde.db` | Path to SQLite database |
| `-from` | int | `30000142` | Source system ID (Jita) |
| `-to` | int | `30002187` | Destination system ID (Amarr) |
| `-warp` | float | `3.0` | Warp speed in AU/s |
| `-align` | float | `6.0` | Align time in seconds |
| `-safe` | bool | `false` | Avoid low-sec/null-sec systems |
| `-exact` | bool | `false` | Use exact CCP warp formula |
| `-init-views` | bool | `false` | Initialize navigation views and exit |

### Example Output

```
=== EVE Navigation System ===
Route: Jita (30000142) → Amarr (30002187)

Finding shortest path...
✓ Path found: 40 jumps

Calculating travel time (simplified formula)...

=== Travel Time Estimate ===
Total jumps:       40
Total time:        920.0 seconds (15.3 minutes)
Avg per jump:      23.0 seconds

Parameters used:
  Warp speed:      3.0 AU/s
  Align time:      6.0 seconds
  Source:          provided

=== Route Preview ===
  1. Jita (30000142)
  2. Maurasi (30000144)
  3. Uitra (30000145)
  ...
 40. Amarr (30002187)

✓ Full route exported to route_result.json
```

### Known System IDs (Trade Hubs)

- **Jita**: 30000142 (The Forge) - Caldari trade hub, highest traffic
- **Amarr**: 30002187 (Domain) - Amarr trade hub
- **Dodixie**: 30002659 (Sinq Laison) - Gallente trade hub
- **Rens**: 30002510 (Heimatar) - Minmatar trade hub
- **Hek**: 30002053 (Metropolis) - Secondary Minmatar hub

### Ship Parameter Presets

#### Interceptor (Fast Travel)
```bash
go run examples/navigation_example.go -warp 6.0 -align 2.5
```

#### Cruiser (Default)
```bash
go run examples/navigation_example.go -warp 3.0 -align 6.0
```

#### Battleship (Combat Fleet)
```bash
go run examples/navigation_example.go -warp 1.5 -align 12.0
```

#### Freighter (Hauling)
```bash
go run examples/navigation_example.go -warp 1.36 -align 30.0
```

### Output Files

The example creates a `route_result.json` file with the complete route information:

```json
{
  "total_seconds": 920.0,
  "total_minutes": 15.33,
  "jumps": 40,
  "avg_seconds_per_jump": 23.0,
  "route": [30000142, 30000144, ...],
  "parameters_used": {
    "warp_speed": 3.0,
    "align_time": 6.0,
    "source": "provided"
  }
}
```

## Prerequisites

Before running the examples, you need to:

1. Download and import the EVE SDE data:
   ```bash
   make sync
   ```

2. The navigation views will be created automatically on first run, or manually:
   ```bash
   go run examples/navigation_example.go -init-views
   ```

## More Examples

See the [navigation documentation](../docs/navigation.md) for more usage examples and SQL queries.
