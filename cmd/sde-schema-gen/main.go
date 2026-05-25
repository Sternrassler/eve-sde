// sde-schema-gen generiert Go Type-Definitionen aus EVE SDE JSONL-Dateienpackage sdeschemagen

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Sternrassler/eve-sde/cmd/sde-schema-gen/generator"
	"github.com/Sternrassler/eve-sde/internal/registry"
)

func main() {
	var (
		inputDir    = flag.String("input", "data/jsonl", "JSONL Input-Verzeichnis")
		outputDir   = flag.String("output", "internal/schema/types", "Go Output-Verzeichnis")
		verbose     = flag.Bool("v", false, "Verbose Logging")
		maxLines    = flag.Int("lines", 0, "Max JSONL Zeilen pro Schema-Analyse (0 = gesamte Datei)")
		registryDir = flag.String("registry", "internal/registry", "Output-Verzeichnis für mappings_gen.go")
	)
	flag.Parse()

	if *verbose {
		log.SetFlags(log.Ltime | log.Lshortfile)
	}

	// Prüfe Input-Verzeichnis
	if _, err := os.Stat(*inputDir); os.IsNotExist(err) {
		log.Fatalf("Input-Verzeichnis existiert nicht: %s", *inputDir)
	}

	// Erstelle Output-Verzeichnis
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		log.Fatalf("Konnte Output-Verzeichnis nicht erstellen: %v", err)
	}

	// Erstelle Registry-Verzeichnis (für mappings_gen.go)
	if err := os.MkdirAll(*registryDir, 0755); err != nil {
		log.Fatalf("Konnte Registry-Verzeichnis nicht erstellen: %v", err)
	}

	// Finde alle JSONL-Dateien
	files, err := filepath.Glob(filepath.Join(*inputDir, "*.jsonl"))
	if err != nil {
		log.Fatalf("Fehler beim Suchen von JSONL-Dateien: %v", err)
	}

	if len(files) == 0 {
		log.Fatalf("Keine JSONL-Dateien gefunden in: %s", *inputDir)
	}

	log.Printf("Analysiere %d JSONL-Dateien...", len(files))

	// Generiere common.go mit LocalizedText
	commonPath := filepath.Join(*outputDir, "common.go")
	if err := generator.WriteCommonTypes(commonPath); err != nil {
		log.Fatalf("Fehler beim Schreiben von common.go: %v", err)
	}
	log.Printf("✓ Generated %s", commonPath)

	// Verarbeite jede JSONL-Datei
	successCount := 0
	var entries []generator.MappingEntry
	for _, file := range files {
		base := filepath.Base(file)
		schemaName := generator.FileNameToTypeName(base)

		if *verbose {
			log.Printf("Analyzing %s...", base)
		}

		// Analysiere Schema
		schema, err := generator.AnalyzeJSONL(file, *maxLines)
		if err != nil {
			log.Printf("WARNUNG: Konnte %s nicht analysieren: %v", file, err)
			continue
		}

		// Generiere Go-Code
		outputFile := filepath.Join(*outputDir, fmt.Sprintf("%s.go", generator.TypeNameToFileName(schemaName)))
		if err := generator.WriteGoFile(outputFile, schemaName, schema, base); err != nil {
			log.Printf("WARNUNG: Konnte %s nicht schreiben: %v", outputFile, err)
			continue
		}

		if *verbose {
			log.Printf("✓ Generated %s", outputFile)
		}
		successCount++

		// Registry-Eintrag sammeln (Ausschlüsse respektieren).
		// Hinweis: successCount zählt generierte Typ-Dateien; ausgeschlossene
		// Datensätze bekommen weiterhin einen Go-Typ, nur keinen Import-Eintrag.
		name := strings.TrimSuffix(base, ".jsonl")
		if reason, excluded := registry.ExcludedDatasets[name]; excluded {
			log.Printf("⊘ Registry: %s ausgeschlossen (%s)", name, reason)
			continue
		}
		indices, hasOverride := registry.IndexOverrides[name]
		if !hasOverride {
			indices = generator.DeriveIndices(schema)
		}
		entries = append(entries, generator.MappingEntry{
			Name:      name,
			JSONLFile: base,
			TypeName:  schemaName,
			Indices:   indices,
		})
	}

	log.Printf("✓ %d von %d Schema-Dateien generiert", successCount, len(files))
	log.Printf("Schemas gespeichert in: %s", *outputDir)

	// Mappings-Registry generieren (stabil nach Name sortiert)
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name < entries[j].Name })
	mappingsPath := filepath.Join(*registryDir, "mappings_gen.go")
	if err := generator.WriteMappingsFile(mappingsPath, "registry", entries); err != nil {
		log.Fatalf("Konnte Mappings nicht generieren: %v", err)
	}
	log.Printf("✓ Generated %s (%d Einträge)", mappingsPath, len(entries))
}
