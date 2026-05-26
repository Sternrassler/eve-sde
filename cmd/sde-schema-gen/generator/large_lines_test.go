package generator

import (
	"strings"
	"testing"
)

// TestAnalyzeJSONL_HandlesLargeLines stellt sicher, dass Zeilen größer als das
// bufio.Scanner-Default-Limit (64 KB) verarbeitet werden — die Mai-2026-SDE
// (z.B. freelanceJobSchemas.jsonl) enthält Zeilen mit ~98 KB.
func TestAnalyzeJSONL_HandlesLargeLines(t *testing.T) {
	bigVal := strings.Repeat("x", 100_000) // > 64 KB Default-Scanner-Limit
	path := writeJSONL(t, []string{
		`{"_key": 1, "name": "` + bigVal + `", "code": 42}`,
	})

	schema, err := AnalyzeJSONL(path, 0)
	if err != nil {
		t.Fatalf("AnalyzeJSONL muss Zeilen > 64 KB verarbeiten: %v", err)
	}
	if _, ok := schema.Fields["code"]; !ok {
		t.Error("Feld 'code' aus der großen Zeile wurde nicht erkannt (Scanner-Limit greift?)")
	}
}
