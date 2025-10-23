# EVE SDE Examples

This directory contains example programs demonstrating the EVE SDE database APIs.

## Directory Structure

```text
examples/
├── navigation/   # Navigation & pathfinding examples
│   └── main.go
├── cargo/        # Cargo & hauling calculation examples
│   └── main.go
└── README.md     # This file
```

## Navigation Example

Demonstrates pathfinding and travel time calculation.

### Usage

```bash
# Basic usage (Jita → Amarr with default Cruiser parameters)
go run examples/navigation/main.go

# Custom route
go run examples/navigation/main.go -from 30000142 -to 30002659

# Custom ship parameters (Interceptor)
go run examples/navigation/main.go -warp 6.0 -align 2.5

# High-sec only route
go run examples/navigation/main.go -safe
```

See [navigation documentation](../docs/navigation.md) for full API reference.

## Cargo Calculator Example

Demonstrates cargo capacity calculations with optional skill modifiers.

### Usage

```bash
# Basic calculation (Badger hauler + Tritanium, no skills)
go run examples/cargo/main.go --ship 648 --item 34

# With Gallente Hauler V (+25% cargo)
go run examples/cargo/main.go --ship 648 --item 34 --racial-hauler 5

# Show ship capacity details
go run examples/cargo/main.go --ship 648 --ship-info

# Custom cargo multiplier
go run examples/cargo/main.go --ship 648 --item 34 --cargo-mult 1.5
```

### Available Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--db` | string | `data/sqlite/eve-sde.db` | Path to SQLite database |
| `--ship` | int | `648` | Ship type ID (Badger) |
| `--item` | int | `34` | Item type ID (Tritanium) |
| `--racial-hauler` | int | `-1` | Racial Hauler skill (0-5, -1=none) |
| `--freighter` | int | `-1` | Freighter skill (0-5, -1=none) |
| `--mining-barge` | int | `-1` | Mining Barge skill (0-5, -1=none) |
| `--cargo-mult` | float | `-1` | Custom cargo multiplier |
| `--ship-info` | bool | `false` | Show detailed ship capacities |
| `--init-views` | bool | `false` | Initialize cargo views and exit |

See [cargo API documentation](../docs/cargo-api.md) for full API reference.
