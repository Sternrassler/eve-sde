package cargo

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/Sternrassler/eve-sde/internal/sqlite/views"
)

// TestIntegrationCargoViews tests the SQL views with an in-memory database
func TestIntegrationCargoViews(t *testing.T) {
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

	// Initialize cargo views
	if err := views.InitializeCargoViews(db); err != nil {
		t.Fatalf("Failed to initialize cargo views: %v", err)
	}

	// Test v_item_volumes
	t.Run("v_item_volumes", func(t *testing.T) {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM v_item_volumes").Scan(&count)
		if err != nil {
			t.Fatalf("Failed to query v_item_volumes: %v", err)
		}
		if count != 3 {
			t.Errorf("Expected 3 items in v_item_volumes, got %d", count)
		}

		// Test specific item
		var itemName string
		var volume float64
		var iskPerM3 float64
		err = db.QueryRow(`
			SELECT item_name, volume, isk_per_m3 
			FROM v_item_volumes 
			WHERE type_id = 34
		`).Scan(&itemName, &volume, &iskPerM3)
		if err != nil {
			t.Fatalf("Failed to query item: %v", err)
		}
		if itemName != "Tritanium" {
			t.Errorf("Expected item name 'Tritanium', got '%s'", itemName)
		}
		if volume != 0.01 {
			t.Errorf("Expected volume 0.01, got %f", volume)
		}
		// ISK/m³ should be 100000 / 0.01 = 10000000
		if iskPerM3 != 10000000.0 {
			t.Errorf("Expected isk_per_m3 10000000, got %f", iskPerM3)
		}
	})

	// Test v_ship_cargo_capacities
	t.Run("v_ship_cargo_capacities", func(t *testing.T) {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM v_ship_cargo_capacities").Scan(&count)
		if err != nil {
			t.Fatalf("Failed to query v_ship_cargo_capacities: %v", err)
		}
		if count != 2 {
			t.Errorf("Expected 2 ships in v_ship_cargo_capacities, got %d", count)
		}

		// Test specific ship
		var shipName string
		var baseCargoCapacity float64
		err = db.QueryRow(`
			SELECT ship_name, base_cargo_capacity 
			FROM v_ship_cargo_capacities 
			WHERE ship_type_id = 648
		`).Scan(&shipName, &baseCargoCapacity)
		if err != nil {
			t.Fatalf("Failed to query ship: %v", err)
		}
		if shipName != "Badger" {
			t.Errorf("Expected ship name 'Badger', got '%s'", shipName)
		}
		if baseCargoCapacity != 3900.0 {
			t.Errorf("Expected base cargo 3900, got %f", baseCargoCapacity)
		}
	})
}

// TestIntegrationGetItemVolume tests GetItemVolume with database
func TestIntegrationGetItemVolume(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	setupTestData(t, db)
	if err := views.InitializeCargoViews(db); err != nil {
		t.Fatalf("Failed to initialize cargo views: %v", err)
	}

	// Test getting Tritanium
	item, err := GetItemVolume(db, 34)
	if err != nil {
		t.Fatalf("Failed to get item volume: %v", err)
	}

	if item.ItemName != "Tritanium" {
		t.Errorf("Expected item name 'Tritanium', got '%s'", item.ItemName)
	}
	if item.Volume != 0.01 {
		t.Errorf("Expected volume 0.01, got %f", item.Volume)
	}
	if item.BasePrice != 100000.0 {
		t.Errorf("Expected base price 100000, got %f", item.BasePrice)
	}

	// Test non-existent item
	_, err = GetItemVolume(db, 99999)
	if err == nil {
		t.Error("Expected error for non-existent item, got nil")
	}
}

// TestIntegrationGetShipCapacities tests GetShipCapacities with database
func TestIntegrationGetShipCapacities(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	setupTestData(t, db)
	if err := views.InitializeCargoViews(db); err != nil {
		t.Fatalf("Failed to initialize cargo views: %v", err)
	}

	// Test without skills
	t.Run("without_skills", func(t *testing.T) {
		ship, err := GetShipCapacities(db, 648, nil)
		if err != nil {
			t.Fatalf("Failed to get ship capacities: %v", err)
		}

		if ship.ShipName != "Badger" {
			t.Errorf("Expected ship name 'Badger', got '%s'", ship.ShipName)
		}
		if ship.BaseCargoHold != 3900.0 {
			t.Errorf("Expected base cargo 3900, got %f", ship.BaseCargoHold)
		}
		if ship.EffectiveCargoHold != 3900.0 {
			t.Errorf("Expected effective cargo 3900 (no skills), got %f", ship.EffectiveCargoHold)
		}
		if ship.SkillsApplied {
			t.Error("Skills should not be applied when nil")
		}
		if ship.SkillBonus != 0.0 {
			t.Errorf("Expected skill bonus 0, got %f", ship.SkillBonus)
		}
	})

	// Test with skills
	t.Run("with_racial_hauler_5", func(t *testing.T) {
		racialLevel := 5
		skills := &SkillModifiers{
			RacialHaulerLevel: &racialLevel,
		}

		ship, err := GetShipCapacities(db, 648, skills)
		if err != nil {
			t.Fatalf("Failed to get ship capacities: %v", err)
		}

		// 3900 * 1.25 = 4875
		if ship.EffectiveCargoHold != 4875.0 {
			t.Errorf("Expected effective cargo 4875 (with skills), got %f", ship.EffectiveCargoHold)
		}
		if !ship.SkillsApplied {
			t.Error("Skills should be marked as applied")
		}
		if ship.SkillBonus != 25.0 {
			t.Errorf("Expected skill bonus 25%%, got %f", ship.SkillBonus)
		}
	})

	// Test non-existent ship
	t.Run("non_existent_ship", func(t *testing.T) {
		_, err := GetShipCapacities(db, 99999, nil)
		if err == nil {
			t.Error("Expected error for non-existent ship, got nil")
		}
	})
}

// TestIntegrationCalculateCargoFit tests CalculateCargoFit with database
func TestIntegrationCalculateCargoFit(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	setupTestData(t, db)
	if err := views.InitializeCargoViews(db); err != nil {
		t.Fatalf("Failed to initialize cargo views: %v", err)
	}

	// Test without skills
	t.Run("badger_tritanium_no_skills", func(t *testing.T) {
		result, err := CalculateCargoFit(db, 648, 34, nil)
		if err != nil {
			t.Fatalf("Failed to calculate cargo fit: %v", err)
		}

		// 3900 m³ / 0.01 m³ = 390,000 units
		if result.MaxQuantity != 390000 {
			t.Errorf("Expected max quantity 390000, got %d", result.MaxQuantity)
		}
		if result.UtilizationPct != 100.0 {
			t.Errorf("Expected utilization 100%%, got %f", result.UtilizationPct)
		}
		if result.SkillsApplied {
			t.Error("Skills should not be applied")
		}
	})

	// Test with skills
	t.Run("badger_tritanium_with_skills", func(t *testing.T) {
		racialLevel := 5
		skills := &SkillModifiers{
			RacialHaulerLevel: &racialLevel,
		}

		result, err := CalculateCargoFit(db, 648, 34, skills)
		if err != nil {
			t.Fatalf("Failed to calculate cargo fit: %v", err)
		}

		// 4875 m³ / 0.01 m³ = 487,500 units
		if result.MaxQuantity != 487500 {
			t.Errorf("Expected max quantity 487500, got %d", result.MaxQuantity)
		}
		if !result.SkillsApplied {
			t.Error("Skills should be applied")
		}
		if result.SkillBonus != 25.0 {
			t.Errorf("Expected skill bonus 25%%, got %f", result.SkillBonus)
		}
	})

	// Test with packaged item (Badger ship itself)
	t.Run("badger_carrying_badger", func(t *testing.T) {
		result, err := CalculateCargoFit(db, 648, 100, nil)
		if err != nil {
			t.Fatalf("Failed to calculate cargo fit: %v", err)
		}

		// 3900 m³ / 20000 m³ (packaged) = 0 units (doesn't fit)
		if result.MaxQuantity != 0 {
			t.Errorf("Expected max quantity 0 (too large), got %d", result.MaxQuantity)
		}
	})
}

// setupTestData creates minimal test data for integration tests
func setupTestData(t *testing.T, db *sql.DB) {
	// Create tables
	schema := `
		CREATE TABLE categories (
			_key INTEGER PRIMARY KEY,
			name TEXT
		);

		CREATE TABLE groups (
			_key INTEGER PRIMARY KEY,
			categoryID INTEGER,
			name TEXT
		);

		CREATE TABLE marketGroups (
			_key INTEGER PRIMARY KEY,
			marketGroupID INTEGER,
			name TEXT
		);

		CREATE TABLE types (
			_key INTEGER PRIMARY KEY,
			groupID INTEGER,
			marketGroupID INTEGER,
			name TEXT,
			volume REAL,
			capacity REAL,
			packagedVolume REAL,
			basePrice REAL,
			published INTEGER
		);

		CREATE TABLE typeDogma (
			typeID INTEGER,
			attributeID INTEGER,
			value REAL
		);
	`

	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("Failed to create schema: %v", err)
	}

	// Insert test data
	data := `
		-- Categories
		INSERT INTO categories (_key, name) VALUES (6, '{"en": "Ship", "de": "Schiff"}');
		INSERT INTO categories (_key, name) VALUES (4, '{"en": "Material", "de": "Material"}');

		-- Groups
		INSERT INTO groups (_key, categoryID, name) VALUES (420, 6, '{"en": "Hauler", "de": "Transporter"}');
		INSERT INTO groups (_key, categoryID, name) VALUES (18, 4, '{"en": "Mineral", "de": "Mineral"}');

		-- Market Groups
		INSERT INTO marketGroups (_key, marketGroupID, name) VALUES (1, 1, '{"en": "Minerals", "de": "Mineralien"}');

		-- Items
		-- Tritanium (typeID: 34)
		INSERT INTO types (_key, groupID, marketGroupID, name, volume, capacity, packagedVolume, basePrice, published) 
		VALUES (34, 18, 1, '{"en": "Tritanium", "de": "Tritanium"}', 0.01, 0, 0, 100000, 1);

		-- Badger (Caldari Hauler, typeID: 648)
		INSERT INTO types (_key, groupID, marketGroupID, name, volume, capacity, packagedVolume, basePrice, published) 
		VALUES (648, 420, NULL, '{"en": "Badger", "de": "Dachs"}', 48500, 3900, 20000, 5000000, 1);

		-- Another item for testing packaged volume (a ship that can be packaged)
		INSERT INTO types (_key, groupID, marketGroupID, name, volume, capacity, packagedVolume, basePrice, published) 
		VALUES (100, 420, NULL, '{"en": "Test Ship", "de": "Testschiff"}', 50000, 1000, 20000, 1000000, 1);
	`

	if _, err := db.Exec(data); err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}
}
