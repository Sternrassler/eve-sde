package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/Sternrassler/eve-sde/internal/schema/types"
	sdeversion "github.com/Sternrassler/eve-sde/internal/sde/version"
	"github.com/Sternrassler/eve-sde/internal/sqlite/importer"
	"github.com/Sternrassler/eve-sde/internal/sqlite/schema"
	"github.com/Sternrassler/eve-sde/internal/sqlite/views"
)

const appVersion = "0.1.0"

// SchemaMapping definiert Mapping zwischen JSONL und Go-Typen
type SchemaMapping struct {
	Name       string
	JSONLFile  string
	StructType reflect.Type
	Indices    []string
}

var schemaMappings = []SchemaMapping{
	{"agentTypes", "agentTypes.jsonl", reflect.TypeOf(types.AgentTypes{}), nil},
	{"ancestries", "ancestries.jsonl", reflect.TypeOf(types.Ancestries{}), []string{"bloodlineID"}},
	{"bloodlines", "bloodlines.jsonl", reflect.TypeOf(types.Bloodlines{}), []string{"raceID"}},
	{"blueprints", "blueprints.jsonl", reflect.TypeOf(types.Blueprints{}), []string{"blueprintTypeID"}},
	{"categories", "categories.jsonl", reflect.TypeOf(types.Categories{}), nil},
	{"certificates", "certificates.jsonl", reflect.TypeOf(types.Certificates{}), []string{"groupID"}},
	{"characterAttributes", "characterAttributes.jsonl", reflect.TypeOf(types.CharacterAttributes{}), nil},
	{"contrabandTypes", "contrabandTypes.jsonl", reflect.TypeOf(types.ContrabandTypes{}), []string{"factionID"}},
	{"controlTowerResources", "controlTowerResources.jsonl", reflect.TypeOf(types.ControlTowerResources{}), []string{"controlTowerTypeID"}},
	{"corporationActivities", "corporationActivities.jsonl", reflect.TypeOf(types.CorporationActivities{}), nil},
	{"dogmaAttributeCategories", "dogmaAttributeCategories.jsonl", reflect.TypeOf(types.DogmaAttributeCategories{}), nil},
	{"dogmaAttributes", "dogmaAttributes.jsonl", reflect.TypeOf(types.DogmaAttributes{}), []string{"categoryID"}},
	{"dogmaEffects", "dogmaEffects.jsonl", reflect.TypeOf(types.DogmaEffects{}), nil},
	{"factions", "factions.jsonl", reflect.TypeOf(types.Factions{}), []string{"solarSystemID"}},
	{"graphics", "graphics.jsonl", reflect.TypeOf(types.Graphics{}), nil},
	{"groups", "groups.jsonl", reflect.TypeOf(types.Groups{}), []string{"categoryID"}},
	{"icons", "icons.jsonl", reflect.TypeOf(types.Icons{}), nil},
	{"marketGroups", "marketGroups.jsonl", reflect.TypeOf(types.MarketGroups{}), []string{"parentGroupID"}},
	{"metaGroups", "metaGroups.jsonl", reflect.TypeOf(types.MetaGroups{}), nil},
	{"npcCorporationDivisions", "npcCorporationDivisions.jsonl", reflect.TypeOf(types.NpcCorporationDivisions{}), []string{"corporationID"}},
	{"npcCorporations", "npcCorporations.jsonl", reflect.TypeOf(types.NpcCorporations{}), []string{"factionID"}},
	{"planetSchematics", "planetSchematics.jsonl", reflect.TypeOf(types.PlanetSchematics{}), nil},
	{"races", "races.jsonl", reflect.TypeOf(types.Races{}), nil},
	{"skinLicenses", "skinLicenses.jsonl", reflect.TypeOf(types.SkinLicenses{}), []string{"skinID"}},
	{"skinMaterials", "skinMaterials.jsonl", reflect.TypeOf(types.SkinMaterials{}), []string{"skinID"}},
	{"skins", "skins.jsonl", reflect.TypeOf(types.Skins{}), nil},
	{"stationOperations", "stationOperations.jsonl", reflect.TypeOf(types.StationOperations{}), nil},
	{"stationServices", "stationServices.jsonl", reflect.TypeOf(types.StationServices{}), nil},
	{"translationLanguages", "translationLanguages.jsonl", reflect.TypeOf(types.TranslationLanguages{}), nil},
	{"typeDogma", "typeDogma.jsonl", reflect.TypeOf(types.TypeDogma{}), []string{"typeID"}},
	{"typeMaterials", "typeMaterials.jsonl", reflect.TypeOf(types.TypeMaterials{}), []string{"typeID", "materialTypeID"}},
	{"types", "types.jsonl", reflect.TypeOf(types.Types{}), []string{"groupID", "marketGroupID"}},
	{"dogmaUnits", "dogmaUnits.jsonl", reflect.TypeOf(types.DogmaUnits{}), nil},
	{"mapConstellations", "mapConstellations.jsonl", reflect.TypeOf(types.MapConstellations{}), []string{"regionID"}},
	{"mapMoons", "mapMoons.jsonl", reflect.TypeOf(types.MapMoons{}), []string{"solarSystemID"}},
	{"mapPlanets", "mapPlanets.jsonl", reflect.TypeOf(types.MapPlanets{}), []string{"solarSystemID"}},
	{"mapRegions", "mapRegions.jsonl", reflect.TypeOf(types.MapRegions{}), nil},
	{"mapSolarSystems", "mapSolarSystems.jsonl", reflect.TypeOf(types.MapSolarSystems{}), []string{"constellationID", "securityClass"}},
	{"mapStargates", "mapStargates.jsonl", reflect.TypeOf(types.MapStargates{}), []string{"solarSystemID", "destination"}},
	{"npcStations", "npcStations.jsonl", reflect.TypeOf(types.NpcStations{}), []string{"solarSystemID", "typeID"}},
	{"_sde", "_sde.jsonl", reflect.TypeOf(types.SDE{}), nil},
}

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
	schemasToImport := schemaMappings
	if *importTable != "" {
		schemasToImport = filterSchemas(schemaMappings, *importTable)
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

	for _, mapping := range schemaMappings {
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
func filterSchemas(all []SchemaMapping, name string) []SchemaMapping {
	for _, s := range all {
		if s.Name == name {
			return []SchemaMapping{s}
		}
	}
	return nil
}
