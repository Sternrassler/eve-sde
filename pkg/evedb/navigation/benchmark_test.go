package navigation

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/Sternrassler/eve-sde/internal/sqlite/views"
)

// BenchmarkShortestPathShort benchmarks pathfinding for short routes (< 10 jumps)
func BenchmarkShortestPathShort(b *testing.B) {
	db := setupBenchmarkDB(b)
	defer db.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ShortestPath(db, 1, 2, false)
		if err != nil {
			b.Fatalf("Failed to find path: %v", err)
		}
	}
}

// BenchmarkShortestPathMedium benchmarks pathfinding for medium routes (20-30 jumps)
func BenchmarkShortestPathMedium(b *testing.B) {
	db := setupBenchmarkDB(b)
	defer db.Close()

	// Create a chain of 30 systems
	setupLongChain(b, db, 30)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ShortestPath(db, 1, 30, false)
		if err != nil {
			b.Fatalf("Failed to find path: %v", err)
		}
	}
}

// BenchmarkShortestPathLong benchmarks pathfinding for long routes (40+ jumps)
func BenchmarkShortestPathLong(b *testing.B) {
	db := setupBenchmarkDB(b)
	defer db.Close()

	// Create a chain of 50 systems
	setupLongChain(b, db, 50)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ShortestPath(db, 1, 50, false)
		if err != nil {
			b.Fatalf("Failed to find path: %v", err)
		}
	}
}

// BenchmarkCalculateTravelTime benchmarks the full travel time calculation
func BenchmarkCalculateTravelTime(b *testing.B) {
	db := setupBenchmarkDB(b)
	defer db.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := CalculateTravelTime(db, 1, 2, nil)
		if err != nil {
			b.Fatalf("Failed to calculate travel time: %v", err)
		}
	}
}

// setupBenchmarkDB creates a minimal test database for benchmarks
func setupBenchmarkDB(b *testing.B) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		b.Fatalf("Failed to create database: %v", err)
	}

	// Create schema
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
		b.Fatalf("Failed to create schema: %v", err)
	}

	// Insert minimal test data
	testData := `
		INSERT INTO mapRegions (_key, name) VALUES 
			(100, '{"en":"Benchmark Region"}');

		INSERT INTO mapConstellations (_key, regionID, name) VALUES
			(200, 100, '{"en":"Benchmark Constellation"}');

		INSERT INTO mapSolarSystems (_key, solarSystemID, name, securityStatus, constellationID, regionID, wormholeClassID, border, corridor, hub) VALUES
			(1, 1, '{"en":"System 1"}', 0.9, 200, 100, NULL, 0, 0, 0),
			(2, 2, '{"en":"System 2"}', 0.9, 200, 100, NULL, 0, 0, 0);

		INSERT INTO mapStargates (_key, solarSystemID, typeID, destination) VALUES
			(1001, 1, 16, '{"solarSystemID": 2}'),
			(1002, 2, 16, '{"solarSystemID": 1}');
	`

	if _, err := db.Exec(testData); err != nil {
		b.Fatalf("Failed to insert test data: %v", err)
	}

	// Initialize views
	if err := views.InitializeNavigationViews(db); err != nil {
		b.Fatalf("Failed to initialize views: %v", err)
	}

	return db
}

// setupLongChain creates a long chain of connected systems for benchmarking
func setupLongChain(b *testing.B, db *sql.DB, length int) {
	// Insert additional systems
	for i := 3; i <= length; i++ {
		_, err := db.Exec(`
			INSERT INTO mapSolarSystems (_key, solarSystemID, name, securityStatus, constellationID, regionID, wormholeClassID, border, corridor, hub) 
			VALUES (?, ?, ?, 0.9, 200, 100, NULL, 0, 0, 0)
		`, i, i, fmt.Sprintf(`{"en":"System %d"}`, i))
		if err != nil {
			b.Fatalf("Failed to insert system %d: %v", i, err)
		}

		// Connect to previous system (forward)
		_, err = db.Exec(`
			INSERT INTO mapStargates (_key, solarSystemID, typeID, destination) 
			VALUES (?, ?, 16, ?)
		`, 1000+i, i, fmt.Sprintf(`{"solarSystemID": %d}`, i-1))
		if err != nil {
			b.Fatalf("Failed to insert gate for system %d: %v", i, err)
		}

		// Connect to previous system (reverse)
		_, err = db.Exec(`
			INSERT INTO mapStargates (_key, solarSystemID, typeID, destination) 
			VALUES (?, ?, 16, ?)
		`, 2000+i, i-1, fmt.Sprintf(`{"solarSystemID": %d}`, i))
		if err != nil {
			b.Fatalf("Failed to insert reverse gate for system %d: %v", i, err)
		}
	}

	// Reinitialize views to pick up new data
	if err := views.InitializeNavigationViews(db); err != nil {
		b.Fatalf("Failed to reinitialize views: %v", err)
	}
}
