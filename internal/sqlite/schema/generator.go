package schema

import (
	"fmt"
	"reflect"
	"strings"
)

// Generator erstellt SQLite DDL aus Go-Structs
type Generator struct {
	// LocalizedAsJSON: Wenn true, wird LocalizedText als JSON gespeichert
	// Wenn false, separate Tabelle (nicht implementiert)
	LocalizedAsJSON bool
}

// NewGenerator erstellt einen neuen Schema-Generator
func NewGenerator() *Generator {
	return &Generator{
		LocalizedAsJSON: true,
	}
}

// GenerateTable erstellt CREATE TABLE Statement aus Go-Struct
func (g *Generator) GenerateTable(tableName string, structType reflect.Type) (string, error) {
	if structType.Kind() != reflect.Struct {
		return "", fmt.Errorf("expected struct, got %s", structType.Kind())
	}

	var columns []string

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		// JSON Tag auslesen
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		// Parse JSON Tag
		parts := strings.Split(jsonTag, ",")
		columnName := parts[0]
		isRequired := !containsOmitEmpty(parts)

		// SQL Typ ermitteln
		sqlType, err := g.goTypeToSQL(field.Type, field.Name)
		if err != nil {
			return "", fmt.Errorf("field %s: %w", field.Name, err)
		}

		// Column Definition
		colDef := fmt.Sprintf("  %s %s", columnName, sqlType)

		// Primary Key Detection
		if columnName == "_key" {
			colDef += " PRIMARY KEY"
		} else if isRequired {
			colDef += " NOT NULL"
		}

		columns = append(columns, colDef)
	}

	// CREATE TABLE Statement
	ddl := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n%s\n);",
		tableName,
		strings.Join(columns, ",\n"))

	return ddl, nil
}

// goTypeToSQL konvertiert Go-Typ zu SQLite-Typ
func (g *Generator) goTypeToSQL(t reflect.Type, fieldName string) (string, error) {
	// Handle pointer types
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.String:
		return "TEXT", nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "INTEGER", nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "INTEGER", nil
	case reflect.Float32, reflect.Float64:
		return "REAL", nil
	case reflect.Bool:
		return "INTEGER", nil // SQLite: 0/1
	case reflect.Slice:
		// Arrays als JSON
		return "TEXT", nil // JSON Array
	case reflect.Map:
		// Maps als JSON
		return "TEXT", nil // JSON Object
	case reflect.Struct:
		// Struct-Namen prüfen
		typeName := t.Name()
		if typeName == "LocalizedText" {
			return "TEXT", nil // JSON Object mit Sprachen
		}
		// Andere Structs als JSON
		return "TEXT", nil
	case reflect.Interface:
		// interface{} → TEXT (JSON oder String)
		return "TEXT", nil
	default:
		return "", fmt.Errorf("unsupported type: %s", t.Kind())
	}
}

// containsOmitEmpty prüft ob "omitempty" in JSON-Tag-Parts
func containsOmitEmpty(parts []string) bool {
	for _, p := range parts {
		if p == "omitempty" {
			return true
		}
	}
	return false
}

// GenerateIndex erstellt Index-Statement
func (g *Generator) GenerateIndex(tableName, columnName string) string {
	indexName := fmt.Sprintf("idx_%s_%s", tableName, columnName)
	return fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s(%s);",
		indexName, tableName, columnName)
}

// GenerateSchema erstellt vollständiges Schema (Tabelle + Indices)
func (g *Generator) GenerateSchema(tableName string, structType reflect.Type, indices []string) ([]string, error) {
	statements := make([]string, 0)

	// Table
	table, err := g.GenerateTable(tableName, structType)
	if err != nil {
		return nil, err
	}
	statements = append(statements, table)

	// Indices (nur für existierende Felder)
	validFields := g.getFieldMap(structType)
	for _, col := range indices {
		if _, exists := validFields[col]; exists {
			idx := g.GenerateIndex(tableName, col)
			statements = append(statements, idx)
		}
		// Ignoriere nicht-existente Felder stillschweigend
	}

	return statements, nil
}

// getFieldMap erstellt Map von JSON-Namen zu Feldinfo
func (g *Generator) getFieldMap(structType reflect.Type) map[string]bool {
	fields := make(map[string]bool)
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}
		parts := strings.Split(jsonTag, ",")
		columnName := parts[0]
		fields[columnName] = true
	}
	return fields
}
