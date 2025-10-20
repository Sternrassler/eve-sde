package importer

import (
	"database/sql"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// TestType für Tests
type TestType struct {
	Key     int64  `json:"_key"`
	Name    string `json:"name"`
	Active  bool   `json:"active"`
	Value   *int64 `json:"value,omitempty"`
	Ignored string `json:"-"`
}

type ComplexType struct {
	Key  int64                  `json:"_key"`
	Data map[string]interface{} `json:"data"`
	Tags []string               `json:"tags"`
}

func TestNewImporter(t *testing.T) {
	tmpDB := filepath.Join(t.TempDir(), "test.db")

	imp, err := NewImporter(tmpDB)
	if err != nil {
		t.Fatalf("NewImporter failed: %v", err)
	}
	defer imp.Close()

	if imp.db == nil {
		t.Error("db should not be nil")
	}
	if imp.batchSize != 1000 {
		t.Errorf("batchSize = %d, want 1000", imp.batchSize)
	}

	// Prüfe ob Pragmas gesetzt wurden
	var journalMode string
	err = imp.db.QueryRow("PRAGMA journal_mode").Scan(&journalMode)
	if err != nil {
		t.Fatalf("Failed to query journal_mode: %v", err)
	}
	if journalMode != "wal" {
		t.Errorf("journal_mode = %s, want wal", journalMode)
	}
}

func TestNewImporter_InvalidPath(t *testing.T) {
	_, err := NewImporter("/invalid/path/\x00/test.db")
	if err == nil {
		t.Error("Expected error for invalid path")
	}
}

func TestDB(t *testing.T) {
	tmpDB := filepath.Join(t.TempDir(), "test.db")
	imp, _ := NewImporter(tmpDB)
	defer imp.Close()

	db := imp.DB()
	if db == nil {
		t.Error("DB() returned nil")
	}
	if db != imp.db {
		t.Error("DB() returned different instance")
	}
}

func TestClose(t *testing.T) {
	tmpDB := filepath.Join(t.TempDir(), "test.db")
	imp, _ := NewImporter(tmpDB)

	err := imp.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}

	// Note: sqlite3 allows multiple Close() calls without error
}

func TestBuildInsertSQL(t *testing.T) {
	imp := &Importer{}
	structType := reflect.TypeOf(TestType{})

	sql, err := imp.buildInsertSQL("test_table", structType)
	if err != nil {
		t.Fatalf("buildInsertSQL failed: %v", err)
	}

	expected := "INSERT INTO test_table (_key, name, active, value) VALUES (?, ?, ?, ?)"
	if sql != expected {
		t.Errorf("SQL = %q, want %q", sql, expected)
	}
}

func TestBuildInsertSQL_NoJSONTags(t *testing.T) {
	type NoTags struct {
		Field1 string
		Field2 int
	}

	imp := &Importer{}
	structType := reflect.TypeOf(NoTags{})

	sql, err := imp.buildInsertSQL("no_tags", structType)
	if err != nil {
		t.Fatalf("buildInsertSQL failed: %v", err)
	}

	expected := "INSERT INTO no_tags () VALUES ()"
	if sql != expected {
		t.Errorf("SQL = %q, want %q", sql, expected)
	}
}

func TestExtractValues(t *testing.T) {
	imp := &Importer{}
	structType := reflect.TypeOf(TestType{})

	val := int64(42)
	data := map[string]interface{}{
		"_key":   float64(123), // JSON numbers sind float64
		"name":   "Test",
		"active": true,
		"value":  float64(42),
	}

	values, err := imp.extractValues(data, structType)
	if err != nil {
		t.Fatalf("extractValues failed: %v", err)
	}

	if len(values) != 4 {
		t.Fatalf("values length = %d, want 4", len(values))
	}

	if values[0] != float64(123) {
		t.Errorf("values[0] = %v, want 123", values[0])
	}
	if values[1] != "Test" {
		t.Errorf("values[1] = %v, want Test", values[1])
	}
	if values[2] != 1 { // bool → int
		t.Errorf("values[2] = %v, want 1", values[2])
	}
	if values[3] != float64(val) {
		t.Errorf("values[3] = %v, want %v", values[3], val)
	}
}

func TestExtractValues_MissingFields(t *testing.T) {
	imp := &Importer{}
	structType := reflect.TypeOf(TestType{})

	data := map[string]interface{}{
		"_key": float64(1),
		// name, active, value fehlen
	}

	values, err := imp.extractValues(data, structType)
	if err != nil {
		t.Fatalf("extractValues failed: %v", err)
	}

	if len(values) != 4 {
		t.Fatalf("values length = %d, want 4", len(values))
	}

	// Fehlende Felder sollten nil sein
	if values[1] != nil {
		t.Errorf("values[1] = %v, want nil", values[1])
	}
	if values[2] != nil {
		t.Errorf("values[2] = %v, want nil", values[2])
	}
	if values[3] != nil {
		t.Errorf("values[3] = %v, want nil", values[3])
	}
}

func TestConvertValueForSQL(t *testing.T) {
	imp := &Importer{}

	tests := []struct {
		name       string
		value      interface{}
		targetType reflect.Type
		want       interface{}
	}{
		{"nil", nil, reflect.TypeOf(""), nil},
		{"bool true", true, reflect.TypeOf(false), 1},
		{"bool false", false, reflect.TypeOf(false), 0},
		{"string", "test", reflect.TypeOf(""), "test"},
		{"int", float64(42), reflect.TypeOf(int64(0)), float64(42)},
		{"map to json", map[string]interface{}{"key": "value"}, reflect.TypeOf(map[string]interface{}{}), `{"key":"value"}`},
		{"slice to json", []interface{}{1, 2, 3}, reflect.TypeOf([]interface{}{}), `[1,2,3]`},
		{"interface map", map[string]interface{}{"a": 1}, reflect.TypeOf((*interface{})(nil)).Elem(), `{"a":1}`},
		{"interface slice", []interface{}{"a", "b"}, reflect.TypeOf((*interface{})(nil)).Elem(), `["a","b"]`},
		{"interface primitive", 123, reflect.TypeOf((*interface{})(nil)).Elem(), 123},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := imp.convertValueForSQL(tt.value, tt.targetType)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertValueForSQL() = %v (%T), want %v (%T)", got, got, tt.want, tt.want)
			}
		})
	}
}

func TestConvertValueForSQL_PointerType(t *testing.T) {
	imp := &Importer{}

	// Pointer auf bool sollte zu int konvertiert werden
	got := imp.convertValueForSQL(true, reflect.TypeOf((*bool)(nil)))
	if got != 1 {
		t.Errorf("convertValueForSQL(true, *bool) = %v, want 1", got)
	}
}

func TestImportJSONL(t *testing.T) {
	tmpDir := t.TempDir()
	tmpDB := filepath.Join(tmpDir, "test.db")
	tmpJSONL := filepath.Join(tmpDir, "test.jsonl")

	// Erstelle Test-JSONL
	jsonlContent := `{"_key":1,"name":"Item1","active":true,"value":100}
{"_key":2,"name":"Item2","active":false}
{"_key":3,"name":"Item3","active":true,"value":200}
`
	if err := os.WriteFile(tmpJSONL, []byte(jsonlContent), 0644); err != nil {
		t.Fatalf("Failed to write JSONL: %v", err)
	}

	// Erstelle Importer und Schema
	imp, err := NewImporter(tmpDB)
	if err != nil {
		t.Fatalf("NewImporter failed: %v", err)
	}
	defer imp.Close()

	// Erstelle Tabelle
	createSQL := `CREATE TABLE test_items (
		_key INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		active INTEGER NOT NULL,
		value INTEGER
	)`
	if _, err := imp.db.Exec(createSQL); err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// Import
	structType := reflect.TypeOf(TestType{})
	err = imp.ImportJSONL("test_items", tmpJSONL, structType)
	if err != nil {
		t.Fatalf("ImportJSONL failed: %v", err)
	}

	// Validierung
	var count int
	err = imp.db.QueryRow("SELECT COUNT(*) FROM test_items").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to count rows: %v", err)
	}
	if count != 3 {
		t.Errorf("Row count = %d, want 3", count)
	}

	// Prüfe erste Zeile
	var key int64
	var name string
	var active int
	var value sql.NullInt64
	err = imp.db.QueryRow("SELECT _key, name, active, value FROM test_items WHERE _key = 1").
		Scan(&key, &name, &active, &value)
	if err != nil {
		t.Fatalf("Failed to query row: %v", err)
	}

	if key != 1 {
		t.Errorf("key = %d, want 1", key)
	}
	if name != "Item1" {
		t.Errorf("name = %s, want Item1", name)
	}
	if active != 1 {
		t.Errorf("active = %d, want 1", active)
	}
	if !value.Valid || value.Int64 != 100 {
		t.Errorf("value = %v, want 100", value)
	}
}

func TestImportJSONL_InvalidFile(t *testing.T) {
	tmpDB := filepath.Join(t.TempDir(), "test.db")
	imp, _ := NewImporter(tmpDB)
	defer imp.Close()

	err := imp.ImportJSONL("test", "/nonexistent/file.jsonl", reflect.TypeOf(TestType{}))
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}

func TestImportJSONL_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	tmpDB := filepath.Join(tmpDir, "test.db")
	tmpJSONL := filepath.Join(tmpDir, "invalid.jsonl")

	type SimpleType struct {
		Key  int64  `json:"_key"`
		Name string `json:"name"`
	}

	// JSONL mit fehlerhafter Zeile
	jsonlContent := `{"_key":1,"name":"Valid"}
{invalid json}
{"_key":2,"name":"AlsoValid"}
`
	if err := os.WriteFile(tmpJSONL, []byte(jsonlContent), 0644); err != nil {
		t.Fatalf("Failed to write JSONL: %v", err)
	}

	imp, _ := NewImporter(tmpDB)
	defer imp.Close()

	// Erstelle Tabelle mit passenden Feldern
	createSQL := `CREATE TABLE test (_key INTEGER, name TEXT)`
	imp.db.Exec(createSQL)

	// Import sollte fehlerhafte Zeilen überspringen
	err := imp.ImportJSONL("test", tmpJSONL, reflect.TypeOf(SimpleType{}))
	if err != nil {
		t.Fatalf("ImportJSONL should skip invalid lines: %v", err)
	}

	// Nur 2 valide Zeilen sollten importiert sein
	var count int
	imp.db.QueryRow("SELECT COUNT(*) FROM test").Scan(&count)
	if count != 2 {
		t.Errorf("Row count = %d, want 2 (invalid line skipped)", count)
	}
}

func TestImportJSONL_ComplexTypes(t *testing.T) {
	tmpDir := t.TempDir()
	tmpDB := filepath.Join(tmpDir, "test.db")
	tmpJSONL := filepath.Join(tmpDir, "complex.jsonl")

	jsonlContent := `{"_key":1,"data":{"nested":"value","count":42},"tags":["tag1","tag2"]}
`
	if err := os.WriteFile(tmpJSONL, []byte(jsonlContent), 0644); err != nil {
		t.Fatalf("Failed to write JSONL: %v", err)
	}

	imp, _ := NewImporter(tmpDB)
	defer imp.Close()

	createSQL := `CREATE TABLE complex (_key INTEGER, data TEXT, tags TEXT)`
	imp.db.Exec(createSQL)

	structType := reflect.TypeOf(ComplexType{})
	err := imp.ImportJSONL("complex", tmpJSONL, structType)
	if err != nil {
		t.Fatalf("ImportJSONL failed: %v", err)
	}

	// Prüfe ob JSON korrekt gespeichert wurde
	var data, tags string
	err = imp.db.QueryRow("SELECT data, tags FROM complex WHERE _key = 1").Scan(&data, &tags)
	if err != nil {
		t.Fatalf("Failed to query: %v", err)
	}

	if data != `{"count":42,"nested":"value"}` && data != `{"nested":"value","count":42}` {
		t.Errorf("data = %s, want JSON object", data)
	}
	if tags != `["tag1","tag2"]` {
		t.Errorf("tags = %s, want [\"tag1\",\"tag2\"]", tags)
	}
}
