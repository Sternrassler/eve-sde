// Example: EVE Cargo Calculator
// This example demonstrates cargo capacity calculations with and without skills
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Sternrassler/eve-sde/internal/sqlite/views"
	"github.com/Sternrassler/eve-sde/pkg/evedb/cargo"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Command line flags
	var (
		dbPath             = flag.String("db", "data/sqlite/eve-sde.db", "Path to SQLite database")
		shipTypeID         = flag.Int64("ship", 648, "Ship type ID (default: Badger)")
		itemTypeID         = flag.Int64("item", 34, "Item type ID (default: Tritanium)")
		racialHaulerLevel  = flag.Int("racial-hauler", -1, "Racial Hauler skill level (0-5, -1 for none)")
		freighterLevel     = flag.Int("freighter", -1, "Freighter skill level (0-5, -1 for none)")
		miningBargeLevel   = flag.Int("mining-barge", -1, "Mining Barge skill level (0-5, -1 for none)")
		cargoMultiplier    = flag.Float64("cargo-mult", -1, "Custom cargo multiplier (e.g. 1.5 for +50%)")
		showShipInfo       = flag.Bool("ship-info", false, "Show detailed ship capacity information")
		initViews          = flag.Bool("init-views", false, "Initialize cargo views and exit")
	)
	flag.Parse()

	// Open database
	db, err := sql.Open("sqlite3", *dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Check if database exists
	if _, err := os.Stat(*dbPath); os.IsNotExist(err) {
		log.Fatalf("Database not found: %s\nPlease run 'make sync' to download and import SDE data first.", *dbPath)
	}

	// Initialize views if requested
	if *initViews {
		log.Println("Initializing cargo views...")
		if err := views.InitializeCargoViews(db); err != nil {
			log.Fatalf("Failed to initialize views: %v", err)
		}
		log.Println("✓ Cargo views initialized successfully")
		return
	}

	// Check if views exist
	var viewExists int
	err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='view' AND name='v_item_volumes'").Scan(&viewExists)
	if err != nil || viewExists == 0 {
		log.Println("⚠ Cargo views not found. Initializing...")
		if err := views.InitializeCargoViews(db); err != nil {
			log.Fatalf("Failed to initialize views: %v", err)
		}
		log.Println("✓ Cargo views initialized")
	}

	// Build skill modifiers from flags
	var skills *cargo.SkillModifiers
	if *racialHaulerLevel >= 0 || *freighterLevel >= 0 || *miningBargeLevel >= 0 || *cargoMultiplier > 0 {
		skills = &cargo.SkillModifiers{}

		if *racialHaulerLevel >= 0 && *racialHaulerLevel <= 5 {
			skills.RacialHaulerLevel = racialHaulerLevel
		}
		if *freighterLevel >= 0 && *freighterLevel <= 5 {
			skills.FreighterLevel = freighterLevel
		}
		if *miningBargeLevel >= 0 && *miningBargeLevel <= 5 {
			skills.MiningBargeLevel = miningBargeLevel
		}
		if *cargoMultiplier > 0 {
			skills.CargoMultiplier = cargoMultiplier
		}
	}

	// Show ship info if requested
	if *showShipInfo {
		showShipCapacities(db, *shipTypeID, skills)
		return
	}

	// Calculate cargo fit
	fmt.Printf("\n=== EVE Cargo Calculator ===\n\n")

	result, err := cargo.CalculateCargoFit(db, *shipTypeID, *itemTypeID, skills)
	if err != nil {
		log.Fatalf("Failed to calculate cargo fit: %v", err)
	}

	// Display results
	fmt.Printf("Ship: %s (Type ID: %d)\n", result.ShipName, result.ShipTypeID)
	fmt.Printf("Base Cargo Capacity: %s m³\n", formatNumber(result.BaseCapacity))

	if result.SkillsApplied {
		fmt.Printf("Skill Bonus: +%.1f%%\n", result.SkillBonus)
		if skills.RacialHaulerLevel != nil {
			fmt.Printf("  - Racial Hauler: Level %d\n", *skills.RacialHaulerLevel)
		}
		if skills.FreighterLevel != nil {
			fmt.Printf("  - Freighter: Level %d\n", *skills.FreighterLevel)
		}
		if skills.MiningBargeLevel != nil {
			fmt.Printf("  - Mining Barge: Level %d\n", *skills.MiningBargeLevel)
		}
		if skills.CargoMultiplier != nil {
			fmt.Printf("  - Custom Multiplier: %.2fx\n", *skills.CargoMultiplier)
		}
		fmt.Printf("Effective Capacity: %s m³\n", formatNumber(result.EffectiveCapacity))
	} else {
		fmt.Printf("Effective Capacity: %s m³ (no skills applied)\n", formatNumber(result.EffectiveCapacity))
	}

	fmt.Printf("\nItem: %s (Type ID: %d)\n", result.ItemName, result.ItemTypeID)
	fmt.Printf("Volume per unit: %.4f m³\n", result.ItemVolume)

	fmt.Printf("\n=== Cargo Fit Results ===\n")
	fmt.Printf("Max Quantity: %s units\n", formatNumber(float64(result.MaxQuantity)))
	fmt.Printf("Total Volume: %s m³\n", formatNumber(result.TotalVolume))
	fmt.Printf("Remaining Space: %s m³\n", formatNumber(result.RemainingSpace))
	fmt.Printf("Utilization: %.2f%%\n", result.UtilizationPct)

	// Show example comparison
	if !result.SkillsApplied {
		fmt.Printf("\n💡 Tip: Use --racial-hauler 5 to see the effect of skills!\n")
		fmt.Printf("   Example: go run examples/cargo_calculator.go --ship %d --item %d --racial-hauler 5\n",
			*shipTypeID, *itemTypeID)
	}
}

// showShipCapacities displays detailed ship capacity information
func showShipCapacities(db *sql.DB, shipTypeID int64, skills *cargo.SkillModifiers) {
	ship, err := cargo.GetShipCapacities(db, shipTypeID, skills)
	if err != nil {
		log.Fatalf("Failed to get ship capacities: %v", err)
	}

	fmt.Printf("\n=== Ship Capacity Details ===\n\n")
	fmt.Printf("Ship: %s (Type ID: %d)\n\n", ship.ShipName, ship.ShipTypeID)

	fmt.Printf("Base Capacities:\n")
	fmt.Printf("  Cargo Hold:     %s m³\n", formatNumber(ship.BaseCargoHold))
	if ship.BaseFleetHangar > 0 {
		fmt.Printf("  Fleet Hangar:   %s m³\n", formatNumber(ship.BaseFleetHangar))
	}
	if ship.BaseOreHold > 0 {
		fmt.Printf("  Ore Hold:       %s m³\n", formatNumber(ship.BaseOreHold))
	}
	fmt.Printf("  Total:          %s m³\n", formatNumber(ship.BaseTotalCapacity))

	if ship.SkillsApplied {
		fmt.Printf("\nEffective Capacities (with skills):\n")
		fmt.Printf("  Cargo Hold:     %s m³ (+%.1f%%)\n",
			formatNumber(ship.EffectiveCargoHold),
			((ship.EffectiveCargoHold/ship.BaseCargoHold)-1)*100)
		if ship.BaseFleetHangar > 0 {
			fmt.Printf("  Fleet Hangar:   %s m³ (+%.1f%%)\n",
				formatNumber(ship.EffectiveFleetHangar),
				((ship.EffectiveFleetHangar/ship.BaseFleetHangar)-1)*100)
		}
		if ship.BaseOreHold > 0 {
			fmt.Printf("  Ore Hold:       %s m³ (+%.1f%%)\n",
				formatNumber(ship.EffectiveOreHold),
				((ship.EffectiveOreHold/ship.BaseOreHold)-1)*100)
		}
		fmt.Printf("  Total:          %s m³ (+%.1f%%)\n",
			formatNumber(ship.EffectiveTotalCapacity),
			ship.SkillBonus)
	}
}

// formatNumber formats a number with thousand separators
func formatNumber(n float64) string {
	if n == 0 {
		return "0"
	}

	// Handle integers
	if n == float64(int64(n)) {
		s := fmt.Sprintf("%d", int64(n))
		return addCommas(s)
	}

	// Handle floats
	s := fmt.Sprintf("%.2f", n)
	parts := []rune(s)

	// Find decimal point
	dotIdx := -1
	for i, r := range parts {
		if r == '.' {
			dotIdx = i
			break
		}
	}

	if dotIdx > 0 {
		intPart := string(parts[:dotIdx])
		decPart := string(parts[dotIdx:])
		return addCommas(intPart) + decPart
	}

	return addCommas(s)
}

// addCommas adds thousand separators to a string
func addCommas(s string) string {
	n := len(s)
	if n <= 3 {
		return s
	}

	result := ""
	for i, r := range s {
		if i > 0 && (n-i)%3 == 0 {
			result += ","
		}
		result += string(r)
	}
	return result
}
