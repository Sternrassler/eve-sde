package navigation

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/Sternrassler/eve-sde/internal/sqlite/views"
)

// TestIntegrationViews tests the SQL views with an in-memory database
func TestIntegrationViews(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create in-memory database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Create test schema and data
	setupTestData(t, db)

	// Initialize views
	if err := views.InitializeNavigationViews(db); err != nil {
		t.Fatalf("Failed to initialize views: %v", err)
	}

	// Test v_stargate_graph
	t.Run("v_stargate_graph", func(t *testing.T) {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM v_stargate_graph").Scan(&count)
		if err != nil {
			t.Fatalf("Failed to query v_stargate_graph: %v", err)
		}
		// We have 2 gates, each creates 2 bidirectional edges = 4 total
		if count != 4 {
			t.Errorf("Expected 4 edges in v_stargate_graph, got %d", count)
		}
	})

	// Test v_system_info
	t.Run("v_system_info", func(t *testing.T) {
		var systemName string
		var securityZone string
		err := db.QueryRow(`
			SELECT system_name, security_zone 
			FROM v_system_info 
			WHERE system_id = 1
		`).Scan(&systemName, &securityZone)
		if err != nil {
			t.Fatalf("Failed to query v_system_info: %v", err)
		}
		if systemName != "Test System 1" {
			t.Errorf("Expected system name 'Test System 1', got '%s'", systemName)
		}
		if securityZone != "High-Sec" {
			t.Errorf("Expected security zone 'High-Sec', got '%s'", securityZone)
		}
	})

	// Test v_region_stats
	t.Run("v_region_stats", func(t *testing.T) {
		var totalSystems int
		err := db.QueryRow(`
			SELECT total_systems 
			FROM v_region_stats 
			WHERE region_id = 100
		`).Scan(&totalSystems)
		if err != nil {
			t.Fatalf("Failed to query v_region_stats: %v", err)
		}
		if totalSystems != 2 {
			t.Errorf("Expected 2 systems in region, got %d", totalSystems)
		}
	})

	// Test v_trade_hubs
	t.Run("v_trade_hubs", func(t *testing.T) {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM v_trade_hubs").Scan(&count)
		if err != nil {
			t.Fatalf("Failed to query v_trade_hubs: %v", err)
		}
		// No trade hubs in our test data
		if count != 0 {
			t.Errorf("Expected 0 trade hubs, got %d", count)
		}
	})
}

// TestIntegrationShortestPath tests pathfinding with real data
func TestIntegrationShortestPath(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	setupTestData(t, db)
	if err := views.InitializeNavigationViews(db); err != nil {
		t.Fatalf("Failed to initialize views: %v", err)
	}

	// Test simple path
	path, err := ShortestPath(db, 1, 2, false)
	if err != nil {
		t.Fatalf("Failed to find path: %v", err)
	}

	if path.Jumps != 1 {
		t.Errorf("Expected 1 jump, got %d", path.Jumps)
	}

	if len(path.Route) != 2 {
		t.Errorf("Expected route length 2, got %d", len(path.Route))
	}

	if path.Route[0] != 1 || path.Route[1] != 2 {
		t.Errorf("Unexpected route: %v", path.Route)
	}
}

// TestIntegrationCalculateTravelTime tests travel time calculation
func TestIntegrationCalculateTravelTime(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	setupTestData(t, db)
	if err := views.InitializeNavigationViews(db); err != nil {
		t.Fatalf("Failed to initialize views: %v", err)
	}

	// Test with default parameters
	result, err := CalculateTravelTime(db, 1, 2, nil)
	if err != nil {
		t.Fatalf("Failed to calculate travel time: %v", err)
	}

	if result.Jumps != 1 {
		t.Errorf("Expected 1 jump, got %d", result.Jumps)
	}

	// Default time per jump should be ~23 seconds
	expectedTime := 23.0
	tolerance := 2.0
	if result.AvgSecondsPerJump < expectedTime-tolerance || result.AvgSecondsPerJump > expectedTime+tolerance {
		t.Errorf("Expected avg time per jump around %.1f, got %.1f", expectedTime, result.AvgSecondsPerJump)
	}

	// Test with custom parameters (interceptor)
	interceptorParams := &NavigationParams{
		WarpSpeed: ptrFloat64(6.0),
		AlignTime: ptrFloat64(2.5),
	}

	result2, err := CalculateTravelTime(db, 1, 2, interceptorParams)
	if err != nil {
		t.Fatalf("Failed to calculate travel time with custom params: %v", err)
	}

	// Interceptor should be faster
	if result2.TotalSeconds >= result.TotalSeconds {
		t.Errorf("Interceptor should be faster than default, got %.1fs vs %.1fs", result2.TotalSeconds, result.TotalSeconds)
	}

	// Check parameters were used
	if params, ok := result2.ParametersUsed["warp_speed"].(float64); !ok || params != 6.0 {
		t.Errorf("Expected warp_speed 6.0 in parameters, got %v", result2.ParametersUsed["warp_speed"])
	}
}

// setupTestData creates minimal test data for integration tests
func setupTestData(t *testing.T, db *sql.DB) {
	// Create tables
	schema := `
		CREATE TABLE mapRegions (
			_key INTEGER PRIMARY KEY,
			name TEXT
		);

		CREATE TABLE mapConstellations (
			_key INTEGER PRIMARY KEY,
			regionID INTEGER,
			name TEXT
		);

		CREATE TABLE mapSolarSystems (
			_key INTEGER PRIMARY KEY,
			solarSystemID INTEGER,
			name TEXT,
			securityStatus REAL,
			constellationID INTEGER,
			regionID INTEGER,
			wormholeClassID INTEGER,
			border INTEGER,
			corridor INTEGER,
			hub INTEGER
		);

		CREATE TABLE mapStargates (
			_key INTEGER PRIMARY KEY,
			solarSystemID INTEGER,
			typeID INTEGER,
			destination TEXT
		);
	`

	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("Failed to create schema: %v", err)
	}

	// Insert test data
	testData := `
		INSERT INTO mapRegions (_key, name) VALUES 
			(100, '{"en":"Test Region"}');

		INSERT INTO mapConstellations (_key, regionID, name) VALUES
			(200, 100, '{"en":"Test Constellation"}');

		INSERT INTO mapSolarSystems (_key, solarSystemID, name, securityStatus, constellationID, regionID, wormholeClassID, border, corridor, hub) VALUES
			(1, 1, '{"en":"Test System 1"}', 0.9, 200, 100, NULL, 0, 0, 0),
			(2, 2, '{"en":"Test System 2"}', 0.5, 200, 100, NULL, 0, 0, 0);

		INSERT INTO mapStargates (_key, solarSystemID, typeID, destination) VALUES
			(1001, 1, 16, '{"solarSystemID": 2}'),
			(1002, 2, 16, '{"solarSystemID": 1}');
	`

	if _, err := db.Exec(testData); err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}
}

// TestMain allows us to set up environment for integration tests
func TestMain(m *testing.M) {
	// Run tests
	code := m.Run()
	os.Exit(code)
}
