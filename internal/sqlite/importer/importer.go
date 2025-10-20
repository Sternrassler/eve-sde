package importer

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/Sternrassler/eve-sde/internal/sqlite/navigation"
)

// Importer importiert JSONL-Daten in SQLite
type Importer struct {
	db        *sql.DB
	batchSize int
}

// DB gibt die Datenbankverbindung zurück
func (i *Importer) DB() *sql.DB {
	return i.db
}

// NewImporter erstellt einen neuen Importer
func NewImporter(dbPath string) (*Importer, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Performance-Optimierungen
	pragmas := []string{
		"PRAGMA journal_mode=WAL",
		"PRAGMA synchronous=NORMAL",
		"PRAGMA cache_size=10000",
		"PRAGMA temp_store=MEMORY",
	}

	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			return nil, fmt.Errorf("failed to set pragma: %w", err)
		}
	}

	return &Importer{
		db:        db,
		batchSize: 1000,
	}, nil
}

// Close schließt die Datenbankverbindung
func (imp *Importer) Close() error {
	return imp.db.Close()
}

// ImportJSONL importiert JSONL-Datei in Tabelle
func (imp *Importer) ImportJSONL(tableName, jsonlPath string, structType reflect.Type) error {
	file, err := os.Open(jsonlPath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Transaction für Performance
	tx, err := imp.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Prepare Insert Statement
	insertSQL, err := imp.buildInsertSQL(tableName, structType)
	if err != nil {
		return fmt.Errorf("failed to build insert SQL: %w", err)
	}

	stmt, err := tx.Prepare(insertSQL)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Stream JSONL
	scanner := bufio.NewScanner(file)
	count := 0

	for scanner.Scan() {
		// Parse JSON
		data := make(map[string]interface{})
		if err := json.Unmarshal(scanner.Bytes(), &data); err != nil {
			continue // Skip fehlerhafte Zeilen
		}

		// Werte extrahieren und einfügen
		values, err := imp.extractValues(data, structType)
		if err != nil {
			return fmt.Errorf("failed to extract values: %w", err)
		}

		if _, err := stmt.Exec(values...); err != nil {
			return fmt.Errorf("failed to insert row: %w", err)
		}

		count++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}

	// Commit
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}

// buildInsertSQL erstellt INSERT Statement
func (imp *Importer) buildInsertSQL(tableName string, structType reflect.Type) (string, error) {
	var columns []string

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		parts := strings.Split(jsonTag, ",")
		columnName := parts[0]
		columns = append(columns, columnName)
	}

	placeholders := make([]string, len(columns))
	for i := range placeholders {
		placeholders[i] = "?"
	}

	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "))

	return sql, nil
}

// extractValues extrahiert Werte aus JSON-Map für Insert
func (imp *Importer) extractValues(data map[string]interface{}, structType reflect.Type) ([]interface{}, error) {
	values := make([]interface{}, 0, structType.NumField())

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		parts := strings.Split(jsonTag, ",")
		columnName := parts[0]

		// Wert aus Map holen
		rawValue, exists := data[columnName]
		if !exists {
			values = append(values, nil)
			continue
		}

		// Konvertiere Wert für SQLite
		sqlValue := imp.convertValueForSQL(rawValue, field.Type)
		values = append(values, sqlValue)
	}

	return values, nil
}

// convertValueForSQL konvertiert JSON-Wert zu SQLite-kompatiblem Wert
func (imp *Importer) convertValueForSQL(value interface{}, targetType reflect.Type) interface{} {
	if value == nil {
		return nil
	}

	// Handle pointer types
	if targetType.Kind() == reflect.Ptr {
		targetType = targetType.Elem()
	}

	switch targetType.Kind() {
	case reflect.Bool:
		if b, ok := value.(bool); ok {
			if b {
				return 1
			}
			return 0
		}
	case reflect.Struct:
		// Struct (z.B. LocalizedText) → JSON
		if bytes, err := json.Marshal(value); err == nil {
			return string(bytes)
		}
	case reflect.Slice, reflect.Map:
		// Arrays/Maps → JSON
		if bytes, err := json.Marshal(value); err == nil {
			return string(bytes)
		}
	case reflect.Interface:
		// interface{} → JSON wenn komplex, sonst direkt
		switch v := value.(type) {
		case map[string]interface{}, []interface{}:
			if bytes, err := json.Marshal(v); err == nil {
				return string(bytes)
			}
		default:
			return v
		}
	default:
		// Primitive Typen direkt
		return value
	}

	return value
}

// InitializeNavigationViews creates navigation views in the database
func (imp *Importer) InitializeNavigationViews() error {
	return navigation.InitializeViews(imp.db)
}
