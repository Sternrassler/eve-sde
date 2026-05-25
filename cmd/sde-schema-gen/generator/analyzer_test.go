package generator

import (
	"os"
	"path/filepath"
	"testing"
)

// writeJSONL schreibt Zeilen in eine temporäre .jsonl-Datei und gibt den Pfad zurück.
func writeJSONL(t *testing.T, lines []string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "data.jsonl")
	content := ""
	for _, l := range lines {
		content += l + "\n"
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write temp jsonl: %v", err)
	}
	return path
}

func TestAnalyzeJSONL_RequiredAcrossManyLines(t *testing.T) {
	// 5 Zeilen (> Sample-Cap von 3). "always" ist in allen vorhanden → required.
	// "sometimes" fehlt in einer Zeile → optional.
	path := writeJSONL(t, []string{
		`{"_key": 1, "always": 10, "sometimes": 1}`,
		`{"_key": 2, "always": 20, "sometimes": 2}`,
		`{"_key": 3, "always": 30, "sometimes": 3}`,
		`{"_key": 4, "always": 40, "sometimes": 4}`,
		`{"_key": 5, "always": 50}`,
	})

	schema, err := AnalyzeJSONL(path, 0)
	if err != nil {
		t.Fatalf("AnalyzeJSONL: %v", err)
	}

	if !schema.Fields["always"].IsRequired {
		t.Error(`"always" sollte required sein (in allen 5 Zeilen vorhanden)`)
	}
	if schema.Fields["sometimes"].IsRequired {
		t.Error(`"sometimes" sollte optional sein (fehlt in Zeile 5)`)
	}
	if !schema.Fields["_key"].IsRequired {
		t.Error(`"_key" sollte immer required sein`)
	}
}
