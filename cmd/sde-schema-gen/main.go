// sde-schema-gen generiert Go Type-Definitionen aus EVE SDE JSONL-Dateienpackage sdeschemagen

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Sternrassler/eve-sde/cmd/sde-schema-gen/generator"
)

func main() {
	var (
		inputDir  = flag.String("input", "data/jsonl", "JSONL Input-Verzeichnis")
		outputDir = flag.String("output", "internal/schema/types", "Go Output-Verzeichnis")
		verbose   = flag.Bool("v", false, "Verbose Logging")
		maxLines  = flag.Int("lines", 100, "Max JSONL Zeilen pro Schema-Analyse")
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
	for _, file := range files {
		schemaName := generator.FileNameToTypeName(filepath.Base(file))

		if *verbose {
			log.Printf("Analyzing %s...", filepath.Base(file))
		}

		// Analysiere Schema
		schema, err := generator.AnalyzeJSONL(file, *maxLines)
		if err != nil {
			log.Printf("WARNUNG: Konnte %s nicht analysieren: %v", file, err)
			continue
		}

		// Generiere Go-Code
		outputFile := filepath.Join(*outputDir, fmt.Sprintf("%s.go", generator.TypeNameToFileName(schemaName)))
		if err := generator.WriteGoFile(outputFile, schemaName, schema, filepath.Base(file)); err != nil {
			log.Printf("WARNUNG: Konnte %s nicht schreiben: %v", outputFile, err)
			continue
		}

		if *verbose {
			log.Printf("✓ Generated %s", outputFile)
		}
		successCount++
	}

	log.Printf("✓ %d von %d Schema-Dateien generiert", successCount, len(files))
	log.Printf("Schemas gespeichert in: %s", *outputDir)
}
