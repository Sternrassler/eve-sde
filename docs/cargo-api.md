# EVE Cargo & Hauling API

Comprehensive cargo capacity calculations and hauling optimization for EVE Online based on SDE data.

## Overview

The Cargo API provides tools for calculating:
- Item volumes and value density (ISK/m³)
- Ship cargo capacities (base and skill-modified)
- Cargo fit calculations (how many items fit in a ship)
- Route security analysis (for safe hauling)

**Key Features:**
- ✅ Based on offline SDE data (no ESI required)
- ✅ Optional skill system (passed as parameters)
- ✅ Supports multiple cargo hold types (cargo, ore hold, fleet hangar)
- ✅ Calculates effective capacities with skill bonuses
- ✅ Packaged volume handling for ship transport

## Quick Start

### Initialize Cargo Views

Before using the API, initialize the SQL views:

```go
import (
    "github.com/Sternrassler/eve-sde/internal/sqlite/views"
)

db, _ := sql.Open("sqlite3", "data/sqlite/eve-sde.db")
err := views.InitializeCargoViews(db)
```

Or from the command line:

```bash
go run examples/cargo_calculator.go --init-views
```

### Basic Usage

```go
import "github.com/Sternrassler/eve-sde/pkg/evedb/cargo"

// Get item volume information
item, err := cargo.GetItemVolume(db, 34) // Tritanium
fmt.Printf("%s: %.4f m³\n", item.ItemName, item.Volume)
// Output: Tritanium: 0.0100 m³

// Get ship capacities (without skills)
ship, err := cargo.GetShipCapacities(db, 648, nil) // Badger
fmt.Printf("%s: %.0f m³\n", ship.ShipName, ship.BaseCargoHold)
// Output: Badger: 3900 m³

// Calculate cargo fit (without skills)
result, err := cargo.CalculateCargoFit(db, 648, 34, nil)
fmt.Printf("Max %s: %d units\n", result.ItemName, result.MaxQuantity)
// Output: Max Tritanium: 390,000 units
```

### With Skills

```go
// Define skills (Gallente Hauler V = +25% cargo)
racialLevel := 5
skills := &cargo.SkillModifiers{
    RacialHaulerLevel: &racialLevel,
}

// Calculate with skills
result, err := cargo.CalculateCargoFit(db, 648, 34, skills)
fmt.Printf("Base: %.0f m³\n", result.BaseCapacity)
fmt.Printf("Bonus: +%.0f%%\n", result.SkillBonus)
fmt.Printf("Effective: %.0f m³\n", result.EffectiveCapacity)
fmt.Printf("Max units: %d\n", result.MaxQuantity)

// Output:
// Base: 3900 m³
// Bonus: +25%
// Effective: 4875 m³
// Max units: 487,500
```

## API Reference

### Types

#### SkillModifiers

Optional skill levels for capacity calculations. All fields are optional.

```go
type SkillModifiers struct {
    RacialHaulerLevel     *int     // 0-5, +5% cargo per level
    FreighterLevel        *int     // 0-5, +5% cargo per level
    MiningBargeLevel      *int     // 0-5, +5% ore hold per level
    CargoMultiplier       *float64 // Direct multiplier (1.0 = no change)
    OreHoldMultiplier     *float64 // Ore hold multiplier
    FleetHangarMultiplier *float64 // Fleet hangar multiplier
}
```

#### ItemVolume

Volume and value information for an item.

```go
type ItemVolume struct {
    TypeID         int64
    ItemName       string
    Volume         float64  // m³ per unit
    Capacity       float64  // Internal capacity (for containers)
    PackagedVolume float64  // Packaged volume (for ships)
    BasePrice      float64  // Base ISK price
    CategoryID     int64
    CategoryName   string
    MarketGroupID  *int64
    IskPerM3       float64  // Value density
}
```

#### ShipCapacities

Complete cargo capacity information for a ship.

```go
type ShipCapacities struct {
    ShipTypeID             int64
    ShipName               string
    
    // Base values (without skills)
    BaseCargoHold          float64
    BaseFleetHangar        float64
    BaseOreHold            float64
    BaseTotalCapacity      float64
    
    // Effective values (with skills)
    EffectiveCargoHold     float64
    EffectiveFleetHangar   float64
    EffectiveOreHold       float64
    EffectiveTotalCapacity float64
    
    // Skill info
    SkillBonus             float64  // Percentage (e.g. 25.0 = +25%)
    SkillsApplied          bool
}
```

#### CargoFitResult

Calculation result for item fitting.

```go
type CargoFitResult struct {
    ShipTypeID        int64
    ShipName          string
    ItemTypeID        int64
    ItemName          string
    ItemVolume        float64
    
    BaseCapacity      float64
    EffectiveCapacity float64
    SkillBonus        float64
    SkillsApplied     bool
    
    MaxQuantity       int      // Maximum units that fit
    TotalVolume       float64  // Volume used
    RemainingSpace    float64  // Unused space
    UtilizationPct    float64  // Percentage utilized
}
```

### Functions

#### GetItemVolume

Retrieves volume information for an item.

```go
func GetItemVolume(db *sql.DB, itemTypeID int64) (*ItemVolume, error)
```

**Example:**
```go
item, err := cargo.GetItemVolume(db, 34) // Tritanium
if err != nil {
    log.Fatal(err)
}
fmt.Printf("%s: %.4f m³, %.0f ISK/m³\n", 
    item.ItemName, item.Volume, item.IskPerM3)
```

#### GetShipCapacities

Retrieves cargo hold capacities for a ship.

```go
func GetShipCapacities(db *sql.DB, shipTypeID int64, skills *SkillModifiers) (*ShipCapacities, error)
```

**Parameters:**
- `shipTypeID`: EVE type ID of the ship
- `skills`: Optional skill modifiers (pass `nil` for base values)

**Example:**
```go
// Without skills
ship, _ := cargo.GetShipCapacities(db, 648, nil)

// With skills
racialLevel := 5
freighterLevel := 3
skills := &cargo.SkillModifiers{
    RacialHaulerLevel: &racialLevel,
    FreighterLevel:    &freighterLevel,
}
ship, _ := cargo.GetShipCapacities(db, 648, skills)
```

#### CalculateCargoFit

Calculates how many items fit in a ship.

```go
func CalculateCargoFit(db *sql.DB, shipTypeID, itemTypeID int64, skills *SkillModifiers) (*CargoFitResult, error)
```

**Parameters:**
- `shipTypeID`: EVE type ID of the ship
- `itemTypeID`: EVE type ID of the item
- `skills`: Optional skill modifiers (pass `nil` for base values)

**Example:**
```go
result, err := cargo.CalculateCargoFit(db, 648, 34, nil)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Ship: %s\n", result.ShipName)
fmt.Printf("Item: %s\n", result.ItemName)
fmt.Printf("Max Quantity: %d units\n", result.MaxQuantity)
fmt.Printf("Utilization: %.1f%%\n", result.UtilizationPct)
```

#### ApplySkillModifiers

Low-level function to apply skill bonuses to capacity values.

```go
func ApplySkillModifiers(baseCapacity float64, skills *SkillModifiers, holdType string) float64
```

**Parameters:**
- `baseCapacity`: Base capacity in m³
- `skills`: Skill modifiers (returns base if nil)
- `holdType`: One of `"cargo"`, `"ore_hold"`, `"fleet_hangar"`

**Example:**
```go
baseCapacity := 1000.0
racialLevel := 5
skills := &cargo.SkillModifiers{RacialHaulerLevel: &racialLevel}

effective := cargo.ApplySkillModifiers(baseCapacity, skills, "cargo")
// effective = 1250.0 (base * 1.25)
```

## Command Line Tool

The `cargo_calculator.go` example provides a command-line interface.

### Basic Usage

```bash
# Calculate cargo fit (default: Badger with Tritanium)
go run examples/cargo_calculator.go

# Specific ship and item
go run examples/cargo_calculator.go --ship 648 --item 34

# With Racial Hauler V
go run examples/cargo_calculator.go --ship 648 --item 34 --racial-hauler 5

# With multiple skills
go run examples/cargo_calculator.go \
    --ship 648 \
    --item 34 \
    --racial-hauler 5 \
    --cargo-mult 1.1

# Show detailed ship information
go run examples/cargo_calculator.go --ship 648 --ship-info
```

### Available Flags

| Flag | Type | Description |
|------|------|-------------|
| `--db` | string | Path to SQLite database (default: `data/sqlite/eve-sde.db`) |
| `--ship` | int64 | Ship type ID (default: 648 = Badger) |
| `--item` | int64 | Item type ID (default: 34 = Tritanium) |
| `--racial-hauler` | int | Racial Hauler skill level (0-5, -1 for none) |
| `--freighter` | int | Freighter skill level (0-5, -1 for none) |
| `--mining-barge` | int | Mining Barge skill level (0-5, -1 for none) |
| `--cargo-mult` | float64 | Custom cargo multiplier (e.g. 1.5 for +50%) |
| `--ship-info` | bool | Show detailed ship capacity information |
| `--init-views` | bool | Initialize cargo views and exit |

### Example Output

```
=== EVE Cargo Calculator ===

Ship: Badger (Type ID: 648)
Base Cargo Capacity: 3,900 m³
Skill Bonus: +25.0%
  - Racial Hauler: Level 5
Effective Capacity: 4,875 m³

Item: Tritanium (Type ID: 34)
Volume per unit: 0.0100 m³

=== Cargo Fit Results ===
Max Quantity: 487,500 units
Total Volume: 4,875 m³
Remaining Space: 0 m³
Utilization: 100.00%
```

## SQL Views

The API uses three SQL views (automatically created):

### v_item_volumes

Provides item volume and value density data.

**Columns:**
- `type_id`: Item type ID
- `item_name`: Localized item name
- `volume`: Volume in m³
- `capacity`: Internal capacity
- `packagedVolume`: Packaged volume (for ships)
- `basePrice`: Base ISK price
- `category_id`, `category_name`: Item category
- `marketGroupID`: Market group
- `isk_per_m3`: Value density (ISK/m³)

### v_ship_cargo_capacities

Ship cargo capacity information.

**Columns:**
- `ship_type_id`: Ship type ID
- `ship_name`: Localized ship name
- `ship_volume`: Ship's own volume
- `base_cargo_capacity`: Base cargo hold (m³)
- `base_fleet_hangar_capacity`: Fleet hangar (m³)
- `base_ore_hold_capacity`: Ore hold (m³)
- `group_id`, `group_name`: Ship group
- `category_id`: Category (always 6 for ships)

### v_route_security_analysis

System security information for hauling routes.

**Columns:**
- `system_id`: Solar system ID
- `system_name`: Localized system name
- `security_status`: Security status (0.0 to 1.0)
- `security_class`: 'High-Sec', 'Low-Sec', 'Null-Sec'
- `gate_count`: Number of stargates (chokepoint indicator)
- `is_border_system`, `is_corridor_system`: Route flags
- `regionID`, `region_name`: Region information

## Skill System

### Supported Skills

The cargo API supports the following skill types:

#### Racial Hauler Skills
- **Effect**: +5% cargo capacity per level
- **Examples**: Caldari Hauler, Gallente Hauler, Amarr Hauler, Minmatar Hauler
- **Applies to**: Cargo hold capacity
- **Max bonus**: +25% at level V

#### Freighter Skills
- **Effect**: +5% cargo capacity per level
- **Examples**: Caldari Freighter, Gallente Freighter, etc.
- **Applies to**: Cargo hold capacity
- **Max bonus**: +25% at level V

#### Mining Barge Skills
- **Effect**: +5% ore hold capacity per level
- **Applies to**: Ore hold capacity (mining ships)
- **Max bonus**: +25% at level V

#### Custom Multipliers
- **Purpose**: Future extension for modules/implants
- **Examples**: Cargo rigs, expanded cargo holds
- **Applies to**: Specific hold types

### Skill Stacking

Skills stack **multiplicatively**:

```go
racialLevel := 5   // +25%
freighterLevel := 3 // +15%

// Result: 1000 * 1.25 * 1.15 = 1437.5 m³
effective := baseCapacity * 1.25 * 1.15
```

### Skill Design Philosophy

**Why skills are optional parameters:**
- ✅ SDE is read-only (no character data)
- ✅ Allows offline calculations
- ✅ Easy to test with/without skills
- ✅ Future ESI integration possible (character skills)
- ✅ Transparent base vs. modified values

## Common Use Cases

### 1. Hauling Profit Calculator

```go
// Get item volume
item, _ := cargo.GetItemVolume(db, itemTypeID)

// Calculate cargo fit
result, _ := cargo.CalculateCargoFit(db, shipTypeID, itemTypeID, skills)

// Calculate profit (example with market prices)
profit := result.MaxQuantity * (sellPrice - buyPrice)
iskPerM3 := profit / result.TotalVolume

fmt.Printf("Max units: %d\n", result.MaxQuantity)
fmt.Printf("Profit: %.2f ISK\n", profit)
fmt.Printf("ISK/m³: %.2f\n", iskPerM3)
```

### 2. Ship Comparison

```go
ships := []int64{648, 649, 650} // Different haulers

for _, shipID := range ships {
    ship, _ := cargo.GetShipCapacities(db, shipID, skills)
    fmt.Printf("%s: %.0f m³ (%.0f m³ base, +%.1f%%)\n",
        ship.ShipName,
        ship.EffectiveTotalCapacity,
        ship.BaseTotalCapacity,
        ship.SkillBonus)
}
```

### 3. Value Density Analysis

```go
items := []int64{34, 35, 36, 37} // Different minerals

for _, itemID := range items {
    item, _ := cargo.GetItemVolume(db, itemID)
    fmt.Printf("%s: %.0f ISK/m³\n", item.ItemName, item.IskPerM3)
}
```

## Known Type IDs

### Common Haulers
- **648**: Badger (Caldari)
- **649**: Tayra (Caldari)
- **650**: Bestower (Amarr)
- **651**: Sigil (Amarr)
- **652**: Wreathe (Minmatar)
- **653**: Hoarder (Minmatar)
- **654**: Iteron Mark V (Gallente)

### Common Items
- **34**: Tritanium
- **35**: Pyerite
- **36**: Mexallon
- **37**: Isogen
- **38**: Nocxium

(Use SDE data or ESI to look up other type IDs)

## Future Extensions

### Planned Features (Phase 2+)

- **Route Risk Assessment**: Analyze hauling route safety
- **Cargo Optimization**: Knapsack solver for max ISK/m³
- **Multi-item Fits**: Calculate mixed cargo loads
- **Jump Freighter Support**: Fuel bay calculations
- **Container Support**: Nested volume calculations

### ESI Integration (Future)

```go
// Potential future extension
charSkills, _ := esi.GetCharacterSkills(characterID)
skills := cargo.SkillsFromESI(charSkills, shipTypeID)
result, _ := cargo.CalculateCargoFit(db, shipTypeID, itemTypeID, skills)
```

## Performance

- **GetItemVolume**: ~0.1ms (indexed lookup)
- **GetShipCapacities**: ~0.2ms (join with dogma attributes)
- **CalculateCargoFit**: ~0.3ms (two queries + calculation)

All operations use indexed queries on the SDE database.

## Error Handling

```go
result, err := cargo.CalculateCargoFit(db, shipTypeID, itemTypeID, skills)
if err != nil {
    switch {
    case strings.Contains(err.Error(), "not found"):
        // Invalid type ID
    case strings.Contains(err.Error(), "zero or negative volume"):
        // Invalid item volume
    default:
        // Database error
    }
}
```

## Testing

```bash
# Run all tests
go test ./pkg/evedb/cargo -v

# Run only unit tests
go test ./pkg/evedb/cargo -v -short

# Run with coverage
go test ./pkg/evedb/cargo -cover
```

## References

- **SDE Schema**: `types.volume`, `types.capacity`, `typeDogma`
- **Dogma Attributes**:
  - 38: Capacity (ore hold)
  - 1556: Fleet Hangar Capacity
- **ADR-001**: DB-Core API Separation
- **EVE University**: [Hauling Guide](https://wiki.eveuniversity.org/Hauling)

## License

Part of the eve-sde project. See LICENSE file.
