package schema

import (
	"reflect"
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
