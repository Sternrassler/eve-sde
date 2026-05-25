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

func TestAnalyzeJSONL_LateFractionalIsFloat(t *testing.T) {
	// Die ersten 3 Zeilen haben ganzzahlige "capacity", Zeile 4 hat 0.5.
	// Mit Whole-File-Analyse (maxLines=0) muss der Typ float64 sein.
	path := writeJSONL(t, []string{
		`{"_key": 1, "capacity": 0}`,
		`{"_key": 2, "capacity": 100}`,
		`{"_key": 3, "capacity": 200}`,
		`{"_key": 4, "capacity": 0.5}`,
	})

	schema, err := AnalyzeJSONL(path, 0)
	if err != nil {
		t.Fatalf("AnalyzeJSONL: %v", err)
	}

	if got := schema.Fields["capacity"].GoType; got != "float64" {
		t.Errorf("capacity GoType = %q, want \"float64\"", got)
	}
}

func TestAnalyzeJSONL_MaxLinesLimitsAnalysis(t *testing.T) {
	// Mit maxLines=2 wird Zeile 3 (0.5) nicht gelesen → bleibt int64.
	path := writeJSONL(t, []string{
		`{"_key": 1, "capacity": 0}`,
		`{"_key": 2, "capacity": 100}`,
		`{"_key": 3, "capacity": 0.5}`,
	})

	schema, err := AnalyzeJSONL(path, 2)
	if err != nil {
		t.Fatalf("AnalyzeJSONL: %v", err)
	}

	if got := schema.Fields["capacity"].GoType; got != "int64" {
		t.Errorf("capacity GoType = %q, want \"int64\" (nur 2 Zeilen analysiert)", got)
	}
}
