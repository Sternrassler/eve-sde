package importer

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// TestImportJSONL_HandlesLargeLines stellt sicher, dass der Importer Zeilen größer
// als das bufio.Scanner-Default-Limit (64 KB) verarbeitet (Mai-2026-SDE: ~98 KB-Zeilen).
func TestImportJSONL_HandlesLargeLines(t *testing.T) {
	tmpDir := t.TempDir()
	tmpDB := filepath.Join(tmpDir, "test.db")
	tmpJSONL := filepath.Join(tmpDir, "big.jsonl")

	bigName := strings.Repeat("x", 100_000) // > 64 KB
	content := `{"_key":1,"name":"` + bigName + `","active":true}` + "\n"
	if err := os.WriteFile(tmpJSONL, []byte(content), 0644); err != nil {
		t.Fatalf("write jsonl: %v", err)
	}

	imp, err := NewImporter(tmpDB)
	if err != nil {
		t.Fatalf("NewImporter: %v", err)
	}
	defer func() { _ = imp.Close() }()

	if _, err := imp.db.Exec(`CREATE TABLE test_items (
		_key INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		active INTEGER NOT NULL,
		value INTEGER
	)`); err != nil {
		t.Fatalf("create table: %v", err)
	}

	if err := imp.ImportJSONL("test_items", tmpJSONL, reflect.TypeOf(TestType{})); err != nil {
		t.Fatalf("ImportJSONL muss Zeilen > 64 KB verarbeiten: %v", err)
	}

	var count int
	if err := imp.db.QueryRow("SELECT COUNT(*) FROM test_items").Scan(&count); err != nil {
		t.Fatalf("count: %v", err)
	}
	if count != 1 {
		t.Errorf("row count = %d, want 1", count)
	}
}
