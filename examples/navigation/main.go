// Example: EVE Navigation System Usage
// This example demonstrates pathfinding and travel time calculation
package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Sternrassler/eve-sde/internal/sqlite/views"
	"github.com/Sternrassler/eve-sde/pkg/evedb/navigation"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Command line flags
	var (
		dbPath      = flag.String("db", "data/sqlite/eve-sde.db", "Path to SQLite database")
		fromSystem  = flag.Int64("from", 30000142, "Source system ID (default: Jita)")
		toSystem    = flag.Int64("to", 30002187, "Destination system ID (default: Amarr)")
		warpSpeed   = flag.Float64("warp", 3.0, "Warp speed in AU/s")
		alignTime   = flag.Float64("align", 6.0, "Align time in seconds")
		avoidLowSec = flag.Bool("safe", false, "Avoid low-sec/null-sec systems")
		exact       = flag.Bool("exact", false, "Use exact CCP warp formula")
		initViews   = flag.Bool("init-views", false, "Initialize navigation views and exit")
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
		log.Println("Initializing navigation views...")
		if err := views.InitializeNavigationViews(db); err != nil {
			log.Fatalf("Failed to initialize views: %v", err)
		}
		log.Println("✓ Navigation views initialized successfully")
		return
	}

	// Check if views exist
	var viewExists int
	err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='view' AND name='v_stargate_graph'").Scan(&viewExists)
	if err != nil || viewExists == 0 {
		log.Println("⚠ Navigation views not found. Initializing...")
		if err := views.InitializeNavigationViews(db); err != nil {
			log.Fatalf("Failed to initialize views: %v", err)
		}
		log.Println("✓ Navigation views initialized")
	}

	// Get system names
	fromName := getSystemName(db, *fromSystem)
	toName := getSystemName(db, *toSystem)

	fmt.Printf("\n=== EVE Navigation System ===\n")
	fmt.Printf("Route: %s (%d) → %s (%d)\n\n", fromName, *fromSystem, toName, *toSystem)

	// Find shortest path
	fmt.Println("Finding shortest path...")
	path, err := navigation.ShortestPath(db, *fromSystem, *toSystem, *avoidLowSec)
	if err != nil {
		log.Fatalf("Failed to find path: %v", err)
	}

	fmt.Printf("✓ Path found: %d jumps\n", path.Jumps)

	// Calculate travel time
	params := &navigation.NavigationParams{
		WarpSpeed:   warpSpeed,
		AlignTime:   alignTime,
		AvoidLowSec: *avoidLowSec,
	}

	var result *navigation.RouteResult
	if *exact {
		fmt.Println("\nCalculating travel time (exact CCP formula)...")
		result, err = navigation.CalculateTravelTimeExact(db, *fromSystem, *toSystem, params)
	} else {
		fmt.Println("\nCalculating travel time (simplified formula)...")
		result, err = navigation.CalculateTravelTime(db, *fromSystem, *toSystem, params)
	}

	if err != nil {
		log.Fatalf("Failed to calculate travel time: %v", err)
	}

	// Display results
	fmt.Printf("\n=== Travel Time Estimate ===\n")
	fmt.Printf("Total jumps:       %d\n", result.Jumps)
	fmt.Printf("Total time:        %.1f seconds (%.1f minutes)\n", result.TotalSeconds, result.TotalMinutes)
	fmt.Printf("Avg per jump:      %.1f seconds\n", result.AvgSecondsPerJump)
	fmt.Printf("\nParameters used:\n")
	fmt.Printf("  Warp speed:      %.1f AU/s\n", result.ParametersUsed["warp_speed"])
	fmt.Printf("  Align time:      %.1f seconds\n", result.ParametersUsed["align_time"])
	fmt.Printf("  Source:          %s\n", result.ParametersUsed["source"])
	if formula, ok := result.ParametersUsed["formula"]; ok {
		fmt.Printf("  Formula:         %s\n", formula)
	}

	// Display route (first 10 and last 10 systems)
	fmt.Printf("\n=== Route Preview ===\n")
	displayRoute(db, result.Route)

	// Export to JSON
	if err := exportJSON(result, "route_result.json"); err != nil {
		log.Printf("Warning: Failed to export JSON: %v", err)
	} else {
		fmt.Printf("\n✓ Full route exported to route_result.json\n")
	}
}

// getSystemName retrieves the system name from the database
func getSystemName(db *sql.DB, systemID int64) string {
	var name string
	query := `
		SELECT COALESCE(
			json_extract(name, '$.en'),
			json_extract(name, '$.de'),
			'Unknown'
		)
		FROM mapSolarSystems
		WHERE _key = ?
	`
	err := db.QueryRow(query, systemID).Scan(&name)
	if err != nil {
		return fmt.Sprintf("System_%d", systemID)
	}
	return name
}

// displayRoute prints a preview of the route
func displayRoute(db *sql.DB, route []int64) {
	maxDisplay := 10
	if len(route) <= maxDisplay*2 {
		// Show all systems
		for i, sysID := range route {
			name := getSystemName(db, sysID)
			fmt.Printf("%3d. %s (%d)\n", i+1, name, sysID)
		}
	} else {
		// Show first 10
		for i := 0; i < maxDisplay; i++ {
			name := getSystemName(db, route[i])
			fmt.Printf("%3d. %s (%d)\n", i+1, name, route[i])
		}
		fmt.Printf("     ... (%d systems omitted) ...\n", len(route)-maxDisplay*2)
		// Show last 10
		for i := len(route) - maxDisplay; i < len(route); i++ {
			name := getSystemName(db, route[i])
			fmt.Printf("%3d. %s (%d)\n", i+1, name, route[i])
		}
	}
}

// exportJSON exports the result to a JSON file
func exportJSON(result *navigation.RouteResult, filename string) error {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}
