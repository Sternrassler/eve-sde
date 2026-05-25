package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Sternrassler/eve-sde/internal/registry"
	sdeversion "github.com/Sternrassler/eve-sde/internal/sde/version"
	"github.com/Sternrassler/eve-sde/internal/sqlite/importer"
	"github.com/Sternrassler/eve-sde/internal/sqlite/schema"
	"github.com/Sternrassler/eve-sde/internal/sqlite/views"
)

const appVersion = "0.1.0"


func main() {
	// Flags
	var (
		dbPath        = flag.String("db", "data/sqlite/eve-sde.db", "SQLite database path")
		jsonlDir      = flag.String("jsonl", "data/jsonl", "JSONL input directory")
		initOnly      = flag.Bool("init", false, "Initialize database schema only")
		importTable   = flag.String("import", "", "Import specific table (empty = all)")
		showVersion   = flag.Bool("version", false, "Show version")
		checkVersion  = flag.Bool("check-version", false, "Check for SDE updates and exit")
		skipIfCurrent = flag.Bool("skip-if-current", false, "Skip import if database is up-to-date")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("sde-to-sqlite v%s\n", appVersion)
		return
	}

	log.Printf("EVE SDE to SQLite Converter v%s", appVersion)

	// Version Check
	if *checkVersion || *skipIfCurrent {
		needsUpdate, latest, local, err := sdeversion.NeedsUpdate(*dbPath)
		if err != nil {
			log.Printf("Warning: Version check failed: %v", err)
			if *checkVersion {
				os.Exit(1)
			}
			// Bei skip-if-current fortfahren trotz Fehler
		} else {
			log.Printf("Latest SDE: %s", latest)
			log.Printf("Local DB:   %s", local)

			if *checkVersion {
				if needsUpdate {
					log.Println("✓ Update available")
					os.Exit(0)
				} else {
					log.Println("✓ Database is up-to-date")
					os.Exit(0)
				}
			}

			if *skipIfCurrent && !needsUpdate {
				log.Println("✓ Database is up-to-date, skipping import")
				return
			}

			if needsUpdate {
				log.Println("→ Update available, proceeding with import")
			}
		}
	}

	// Erstelle DB-Verzeichnis
	if err := os.MkdirAll(filepath.Dir(*dbPath), 0755); err != nil {
		log.Fatalf("Failed to create database directory: %v", err)
	}

	// Initialisiere Schema
	if err := initializeSchema(*dbPath); err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
	}
	log.Println("✓ Schema initialized")

	if *initOnly {
		log.Println("Init-only mode, exiting")
		return
	}

	// Import Daten
	imp, err := importer.NewImporter(*dbPath)
	if err != nil {
		log.Fatalf("Failed to create importer: %v", err)
	}
	defer imp.Close()

	// Filter Schemas
	schemasToImport := registry.Mappings
	if *importTable != "" {
		schemasToImport = filterSchemas(registry.Mappings, *importTable)
		if len(schemasToImport) == 0 {
			log.Fatalf("Table not found: %s", *importTable)
		}
	}

	// Import
	for _, mapping := range schemasToImport {
		jsonlPath := filepath.Join(*jsonlDir, mapping.JSONLFile)
		log.Printf("Importing %s from %s...", mapping.Name, mapping.JSONLFile)

		if err := imp.ImportJSONL(mapping.Name, jsonlPath, mapping.StructType); err != nil {
			log.Fatalf("Failed to import %s: %v", mapping.Name, err)
		}

		log.Printf("✓ Imported %s", mapping.Name)
	}

	// Initialize navigation views if we imported map data
	if *importTable == "" || strings.HasPrefix(*importTable, "map") {
		log.Println("Initializing navigation views...")
		if err := views.InitializeNavigationViews(imp.DB()); err != nil {
			log.Printf("Warning: Failed to initialize navigation views: %v", err)
		} else {
			log.Println("✓ Navigation views initialized")
		}
	}

	// Initialize cargo views if we imported types/dogma data
	if *importTable == "" || *importTable == "types" || *importTable == "typeDogma" {
		log.Println("Initializing cargo views...")
		if err := views.InitializeCargoViews(imp.DB()); err != nil {
			log.Printf("Warning: Failed to initialize cargo views: %v", err)
		} else {
			log.Println("✓ Cargo views initialized")
		}
	}

	log.Println("✓ Import completed successfully")
}

// initializeSchema erstellt DB-Schema
func initializeSchema(dbPath string) error {
	imp, err := importer.NewImporter(dbPath)
	if err != nil {
		return err
	}
	defer imp.Close()

	gen := schema.NewGenerator()

	for _, mapping := range registry.Mappings {
		statements, err := gen.GenerateSchema(mapping.Name, mapping.StructType, mapping.Indices)
		if err != nil {
			return fmt.Errorf("failed to generate schema for %s: %w", mapping.Name, err)
		}

		for _, stmt := range statements {
			if _, err := imp.DB().Exec(stmt); err != nil {
				return fmt.Errorf("failed to execute DDL: %w", err)
			}
		}
	}

	return nil
}

// filterSchemas filtert Schemas nach Name
func filterSchemas(all []registry.SchemaMapping, name string) []registry.SchemaMapping {
	for _, s := range all {
		if s.Name == name {
			return []registry.SchemaMapping{s}
		}
	}
	return nil
}
