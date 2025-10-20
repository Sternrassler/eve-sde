// Package cargo provides EVE Online cargo and hauling calculation functionality
// It includes cargo capacity calculations, item volume queries, and skill-based modifications
//
// This is a high-level API layer built on top of the eve-sde SQLite database.
// It requires cargo views to be initialized (see internal/sqlite/views package).
package cargo

import (
	"database/sql"
	"fmt"
)

// SkillModifiers contains optional skill levels for capacity calculations
// All fields are optional - if nil, base values are used
type SkillModifiers struct {
	// Racial Hauler Skills (e.g. Gallente Hauler, Caldari Hauler)
	// +5% cargo capacity per level
	RacialHaulerLevel *int `json:"racial_hauler_level,omitempty"`

	// Freighter Skills
	// +5% cargo capacity per level
	FreighterLevel *int `json:"freighter_level,omitempty"`

	// Mining Barge/Exhumer Skills
	// Affects ore hold capacity
	MiningBargeLevel *int `json:"mining_barge_level,omitempty"`

	// Custom Multipliers (for future extensions like implants/modules)
	CargoMultiplier       *float64 `json:"cargo_multiplier,omitempty"`       // Direct multiplier (1.0 = no change)
	OreHoldMultiplier     *float64 `json:"ore_hold_multiplier,omitempty"`    // Ore hold multiplier
	FleetHangarMultiplier *float64 `json:"fleet_hangar_multiplier,omitempty"` // Fleet hangar multiplier
}

// ItemVolume contains volume and pricing information for an item
type ItemVolume struct {
	TypeID         int64   `json:"type_id"`
	ItemName       string  `json:"item_name"`
	Volume         float64 `json:"volume"`
	Capacity       float64 `json:"capacity"`
	PackagedVolume float64 `json:"packaged_volume"`
	BasePrice      float64 `json:"base_price"`
	CategoryID     int64   `json:"category_id"`
	CategoryName   string  `json:"category_name"`
	MarketGroupID  *int64  `json:"market_group_id,omitempty"`
	IskPerM3       float64 `json:"isk_per_m3"` // Value density
}

// ShipCapacities contains all cargo holds of a ship
type ShipCapacities struct {
	ShipTypeID int64  `json:"ship_type_id"`
	ShipName   string `json:"ship_name"`

	// Base values (without skills)
	BaseCargoHold   float64 `json:"base_cargo_hold"`
	BaseFleetHangar float64 `json:"base_fleet_hangar,omitempty"`
	BaseOreHold     float64 `json:"base_ore_hold,omitempty"`

	// Effective values (with skills applied)
	EffectiveCargoHold   float64 `json:"effective_cargo_hold"`
	EffectiveFleetHangar float64 `json:"effective_fleet_hangar,omitempty"`
	EffectiveOreHold     float64 `json:"effective_ore_hold,omitempty"`

	// Totals
	BaseTotalCapacity      float64 `json:"base_total_capacity"`
	EffectiveTotalCapacity float64 `json:"effective_total_capacity"`

	// Skill information
	SkillBonus    float64 `json:"skill_bonus"`     // Percentage increase (e.g. 25.0 for +25%)
	SkillsApplied bool    `json:"skills_applied"`  // Whether skills were applied
}

// CargoFitResult describes how many items fit in a ship
type CargoFitResult struct {
	ShipTypeID int64  `json:"ship_type_id"`
	ShipName   string `json:"ship_name"`
	ItemTypeID int64  `json:"item_type_id"`
	ItemName   string `json:"item_name"`
	ItemVolume float64 `json:"item_volume"`

	// Capacity values
	BaseCapacity      float64 `json:"base_capacity"`
	EffectiveCapacity float64 `json:"effective_capacity"`

	// Skill modification
	SkillBonus    float64 `json:"skill_bonus"`    // Percentage increase
	SkillsApplied bool    `json:"skills_applied"`

	// Calculation results
	MaxQuantity    int     `json:"max_quantity"`
	TotalVolume    float64 `json:"total_volume"`
	RemainingSpace float64 `json:"remaining_space"`
	UtilizationPct float64 `json:"utilization_pct"`
}

// GetItemVolume retrieves volume information for an item
func GetItemVolume(db *sql.DB, itemTypeID int64) (*ItemVolume, error) {
	query := `
		SELECT 
			type_id,
			item_name,
			COALESCE(volume, 0) as volume,
			COALESCE(capacity, 0) as capacity,
			COALESCE(packagedVolume, 0) as packagedVolume,
			COALESCE(basePrice, 0) as basePrice,
			category_id,
			COALESCE(category_name, '') as category_name,
			marketGroupID,
			COALESCE(isk_per_m3, 0) as isk_per_m3
		FROM v_item_volumes
		WHERE type_id = ?
	`

	var item ItemVolume
	var marketGroupID sql.NullInt64

	err := db.QueryRow(query, itemTypeID).Scan(
		&item.TypeID,
		&item.ItemName,
		&item.Volume,
		&item.Capacity,
		&item.PackagedVolume,
		&item.BasePrice,
		&item.CategoryID,
		&item.CategoryName,
		&marketGroupID,
		&item.IskPerM3,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("item with type ID %d not found", itemTypeID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query item volume: %w", err)
	}

	if marketGroupID.Valid {
		item.MarketGroupID = &marketGroupID.Int64
	}

	return &item, nil
}

// GetShipCapacities retrieves all cargo holds for a ship
// If skills parameter is nil, only base values are returned
func GetShipCapacities(db *sql.DB, shipTypeID int64, skills *SkillModifiers) (*ShipCapacities, error) {
	query := `
		SELECT 
			ship_type_id,
			ship_name,
			COALESCE(base_cargo_capacity, 0) as base_cargo,
			COALESCE(base_fleet_hangar_capacity, 0) as base_fleet_hangar,
			COALESCE(base_ore_hold_capacity, 0) as base_ore_hold
		FROM v_ship_cargo_capacities
		WHERE ship_type_id = ?
	`

	var ship ShipCapacities
	err := db.QueryRow(query, shipTypeID).Scan(
		&ship.ShipTypeID,
		&ship.ShipName,
		&ship.BaseCargoHold,
		&ship.BaseFleetHangar,
		&ship.BaseOreHold,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("ship with type ID %d not found", shipTypeID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query ship capacities: %w", err)
	}

	// Calculate base total
	ship.BaseTotalCapacity = ship.BaseCargoHold + ship.BaseFleetHangar + ship.BaseOreHold

	// Apply skill modifiers
	if skills != nil {
		ship.SkillsApplied = true
		ship.EffectiveCargoHold = ApplySkillModifiers(ship.BaseCargoHold, skills, "cargo")
		ship.EffectiveFleetHangar = ApplySkillModifiers(ship.BaseFleetHangar, skills, "fleet_hangar")
		ship.EffectiveOreHold = ApplySkillModifiers(ship.BaseOreHold, skills, "ore_hold")
		ship.EffectiveTotalCapacity = ship.EffectiveCargoHold + ship.EffectiveFleetHangar + ship.EffectiveOreHold

		// Calculate skill bonus percentage
		if ship.BaseTotalCapacity > 0 {
			ship.SkillBonus = ((ship.EffectiveTotalCapacity / ship.BaseTotalCapacity) - 1.0) * 100.0
		}
	} else {
		ship.SkillsApplied = false
		ship.EffectiveCargoHold = ship.BaseCargoHold
		ship.EffectiveFleetHangar = ship.BaseFleetHangar
		ship.EffectiveOreHold = ship.BaseOreHold
		ship.EffectiveTotalCapacity = ship.BaseTotalCapacity
		ship.SkillBonus = 0.0
	}

	return &ship, nil
}

// CalculateCargoFit calculates how many items fit in a ship
// skills parameter is optional - if nil, base values are used
func CalculateCargoFit(db *sql.DB, shipTypeID, itemTypeID int64, skills *SkillModifiers) (*CargoFitResult, error) {
	// Get ship capacities
	ship, err := GetShipCapacities(db, shipTypeID, skills)
	if err != nil {
		return nil, err
	}

	// Get item volume
	item, err := GetItemVolume(db, itemTypeID)
	if err != nil {
		return nil, err
	}

	// Use packaged volume if available (for ships being transported)
	itemVol := item.Volume
	if item.PackagedVolume > 0 {
		itemVol = item.PackagedVolume
	}

	if itemVol <= 0 {
		return nil, fmt.Errorf("item %s has zero or negative volume", item.ItemName)
	}

	// Calculate fit
	result := &CargoFitResult{
		ShipTypeID:        ship.ShipTypeID,
		ShipName:          ship.ShipName,
		ItemTypeID:        item.TypeID,
		ItemName:          item.ItemName,
		ItemVolume:        itemVol,
		BaseCapacity:      ship.BaseTotalCapacity,
		EffectiveCapacity: ship.EffectiveTotalCapacity,
		SkillBonus:        ship.SkillBonus,
		SkillsApplied:     ship.SkillsApplied,
	}

	// Calculate max quantity
	result.MaxQuantity = int(result.EffectiveCapacity / itemVol)
	result.TotalVolume = float64(result.MaxQuantity) * itemVol
	result.RemainingSpace = result.EffectiveCapacity - result.TotalVolume

	// Calculate utilization percentage
	if result.EffectiveCapacity > 0 {
		result.UtilizationPct = (result.TotalVolume / result.EffectiveCapacity) * 100.0
	}

	return result, nil
}

// ApplySkillModifiers calculates effective capacity based on skills
// If skills is nil, returns baseCapacity unchanged
// holdType can be: "cargo", "ore_hold", "fleet_hangar"
func ApplySkillModifiers(baseCapacity float64, skills *SkillModifiers, holdType string) float64 {
	if skills == nil {
		return baseCapacity
	}

	effective := baseCapacity

	switch holdType {
	case "cargo":
		// Racial Hauler Skill (5% per level)
		if skills.RacialHaulerLevel != nil {
			bonus := float64(*skills.RacialHaulerLevel) * 0.05
			effective *= (1.0 + bonus)
		}

		// Freighter Skill (5% per level)
		if skills.FreighterLevel != nil {
			bonus := float64(*skills.FreighterLevel) * 0.05
			effective *= (1.0 + bonus)
		}

		// Custom multiplier
		if skills.CargoMultiplier != nil {
			effective *= *skills.CargoMultiplier
		}

	case "ore_hold":
		// Mining Barge Skill affects ore hold
		if skills.MiningBargeLevel != nil {
			bonus := float64(*skills.MiningBargeLevel) * 0.05
			effective *= (1.0 + bonus)
		}

		// Custom multiplier
		if skills.OreHoldMultiplier != nil {
			effective *= *skills.OreHoldMultiplier
		}

	case "fleet_hangar":
		// Custom multiplier
		if skills.FleetHangarMultiplier != nil {
			effective *= *skills.FleetHangarMultiplier
		}
	}

	return effective
}
