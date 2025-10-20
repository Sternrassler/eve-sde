package schema

import (
	"reflect"
	"strings"
	"testing"

	"github.com/Sternrassler/eve-sde/internal/schema/types"
)

func TestGenerateTable(t *testing.T) {
	gen := NewGenerator()

	// Test mit Types Struct
	typesType := reflect.TypeOf(types.Types{})
	ddl, err := gen.GenerateTable("types", typesType)
	if err != nil {
		t.Fatalf("GenerateTable failed: %v", err)
	}

	// Prüfe DDL
	if ddl == "" {
		t.Error("DDL is empty")
	}

	// Prüfe Primary Key
	if !contains(ddl, "_key INTEGER PRIMARY KEY") {
		t.Error("Missing PRIMARY KEY on _key")
	}

	// Prüfe LocalizedText als TEXT
	if !contains(ddl, "name TEXT") {
		t.Error("Missing name TEXT column")
	}

	t.Logf("Generated DDL:\n%s", ddl)
}

func TestGoTypeToSQL(t *testing.T) {
	gen := NewGenerator()

	tests := []struct {
		goType   reflect.Type
		expected string
	}{
		{reflect.TypeOf(""), "TEXT"},
		{reflect.TypeOf(int64(0)), "INTEGER"},
		{reflect.TypeOf(float64(0)), "REAL"},
		{reflect.TypeOf(true), "INTEGER"},
		{reflect.TypeOf([]string{}), "TEXT"},
		{reflect.TypeOf(map[string]interface{}{}), "TEXT"},
	}

	for _, tt := range tests {
		result, err := gen.goTypeToSQL(tt.goType)
		if err != nil {
			t.Errorf("goTypeToSQL(%v) failed: %v", tt.goType, err)
			continue
		}
		if result != tt.expected {
			t.Errorf("goTypeToSQL(%v) = %s, want %s", tt.goType, result, tt.expected)
		}
	}
}

func TestGenerateIndex(t *testing.T) {
	gen := NewGenerator()
	idx := gen.GenerateIndex("types", "groupID")

	expected := "CREATE INDEX IF NOT EXISTS idx_types_groupID ON types(groupID);"
	if idx != expected {
		t.Errorf("GenerateIndex = %s, want %s", idx, expected)
	}
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 &&
		(s == substr || len(s) >= len(substr) &&
			(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
				findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestGenerateSchema(t *testing.T) {
	gen := NewGenerator()

	// Test mit einzelner Tabelle + Index
	ddl, err := gen.GenerateSchema("types", reflect.TypeOf(types.Types{}), []string{"groupID"})
	if err != nil {
		t.Fatalf("GenerateSchema failed: %v", err)
	}

	if len(ddl) == 0 {
		t.Error("DDL is empty")
	}

	if len(ddl) != 2 { // Table + Index
		t.Errorf("DDL length = %d, want 2 (table + index)", len(ddl))
	}

	fullDDL := strings.Join(ddl, "\n")

	// Prüfe Tabelle (IF NOT EXISTS, nicht "CREATE TABLE types")
	if !contains(fullDDL, "CREATE TABLE IF NOT EXISTS types") {
		t.Error("Missing types table")
	}

	// Prüfe Index
	if !contains(fullDDL, "CREATE INDEX") {
		t.Error("Missing index")
	}
}

func TestGenerateSchema_NoIndices(t *testing.T) {
	gen := NewGenerator()

	ddl, err := gen.GenerateSchema("simple", reflect.TypeOf(types.Groups{}), []string{})
	if err != nil {
		t.Fatalf("GenerateSchema failed: %v", err)
	}

	if len(ddl) != 1 { // Nur Tabelle, kein Index
		t.Errorf("DDL length = %d, want 1", len(ddl))
	}
}

func TestGetFieldMap(t *testing.T) {
	gen := NewGenerator()

	typesType := reflect.TypeOf(types.Types{})
	fieldMap := gen.getFieldMap(typesType)

	// Prüfe ob _key vorhanden ist
	if _, exists := fieldMap["_key"]; !exists {
		t.Error("Missing _key in field map")
	}

	// Prüfe ob name (LocalizedText) vorhanden ist
	if _, exists := fieldMap["name"]; !exists {
		t.Error("Missing name in field map")
	}

	// Prüfe ob groupID vorhanden ist
	if _, exists := fieldMap["groupID"]; !exists {
		t.Error("Missing groupID in field map")
	}
}

func TestGenerateTable_WithoutPrimaryKey(t *testing.T) {
	gen := NewGenerator()

	type NoPK struct {
		Name string `json:"name"`
		Age  int64  `json:"age"`
	}

	ddl, err := gen.GenerateTable("no_pk", reflect.TypeOf(NoPK{}))
	if err != nil {
		t.Fatalf("GenerateTable failed: %v", err)
	}

	// Sollte keine PRIMARY KEY enthalten
	if contains(ddl, "PRIMARY KEY") {
		t.Error("Should not contain PRIMARY KEY")
	}
}

func TestGoTypeToSQL_PointerTypes(t *testing.T) {
	gen := NewGenerator()

	tests := []struct {
		name     string
		goType   reflect.Type
		expected string
	}{
		{"*int64", reflect.TypeOf((*int64)(nil)), "INTEGER"},
		{"*string", reflect.TypeOf((*string)(nil)), "TEXT"},
		{"*bool", reflect.TypeOf((*bool)(nil)), "INTEGER"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := gen.goTypeToSQL(tt.goType)
			if err != nil {
				t.Errorf("goTypeToSQL(%s) failed: %v", tt.name, err)
				return
			}
			if result != tt.expected {
				t.Errorf("goTypeToSQL(%s) = %s, want %s", tt.name, result, tt.expected)
			}
		})
	}
}

func TestGoTypeToSQL_LocalizedText(t *testing.T) {
	gen := NewGenerator()

	localizedType := reflect.TypeOf(types.LocalizedText{})
	result, err := gen.goTypeToSQL(localizedType)
	if err != nil {
		t.Fatalf("goTypeToSQL(LocalizedText) failed: %v", err)
	}

	if result != "TEXT" {
		t.Errorf("goTypeToSQL(LocalizedText) = %s, want TEXT", result)
	}
}

func TestGoTypeToSQL_UnsupportedType(t *testing.T) {
	gen := NewGenerator()

	// Test mit channel type (nicht unterstützt)
	chanType := reflect.TypeOf(make(chan int))
	_, err := gen.goTypeToSQL(chanType)

	if err == nil {
		t.Error("Expected error for unsupported channel type")
	}
}

func TestContainsOmitEmpty(t *testing.T) {
	tests := []struct {
		parts    []string
		expected bool
	}{
		{[]string{"field", "omitempty"}, true},
		{[]string{"field"}, false},
		{[]string{"field", "string", "omitempty"}, true},
		{[]string{}, false},
	}

	for _, tt := range tests {
		result := containsOmitEmpty(tt.parts)
		if result != tt.expected {
			t.Errorf("containsOmitEmpty(%v) = %v, want %v", tt.parts, result, tt.expected)
		}
	}
}
