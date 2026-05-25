package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDeriveIndices(t *testing.T) {
	schema := &Schema{Fields: map[string]*FieldInfo{
		"_key":          {},
		"groupID":       {},
		"marketGroupID": {},
		"name":          {},
		"published":     {},
	}}

	got := DeriveIndices(schema)

	want := []string{"groupID", "marketGroupID"} // sortiert, *ID außer _key
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Errorf("DeriveIndices = %v, want %v", got, want)
	}
}

func TestWriteMappingsFile(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "mappings_gen.go")

	entries := []MappingEntry{
		{Name: "types", JSONLFile: "types.jsonl", TypeName: "Types", Indices: []string{"groupID", "marketGroupID"}},
		{Name: "icons", JSONLFile: "icons.jsonl", TypeName: "Icons", Indices: nil},
	}

	if err := WriteMappingsFile(out, "registry", entries); err != nil {
		t.Fatalf("WriteMappingsFile: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	src := string(data)

	for _, want := range []string{
		"package registry",
		"DO NOT EDIT",
		"var Mappings = []SchemaMapping{",
		`{Name: "types", JSONLFile: "types.jsonl", StructType: reflect.TypeOf(types.Types{}), Indices: []string{"groupID", "marketGroupID"}}`,
		`{Name: "icons", JSONLFile: "icons.jsonl", StructType: reflect.TypeOf(types.Icons{}), Indices: nil}`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("generierte Datei enthält %q nicht.\n--- Output ---\n%s", want, src)
		}
	}
}
