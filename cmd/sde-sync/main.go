package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	sdeversion "github.com/Sternrassler/eve-sde/internal/sde/version"
	_ "github.com/mattn/go-sqlite3"
)

const appVersion = "0.1.0"

func main() {
	var (
		dataDir     = flag.String("data", "data", "Data directory (contains jsonl/, yaml/, sqlite/)")
		forceUpdate = flag.Bool("force", false, "Force update even if current")
		skipImport  = flag.Bool("skip-import", false, "Skip SQLite import (download + schema-gen only)")
		showVersion = flag.Bool("version", false, "Show version")
		verbose     = flag.Bool("v", false, "Verbose output")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("sde-sync v%s\n", appVersion)
		return
	}

	log.Printf("EVE SDE Sync v%s", appVersion)
	start := time.Now()

	// Pfade
	jsonlDir := filepath.Join(*dataDir, "jsonl")
	sqliteDB := filepath.Join(*dataDir, "sqlite", "eve-sde.db")

	// 1. Version Check (und DB-Löschung bei --force)
	if !*forceUpdate {
		log.Println("→ Checking for SDE updates...")
		needsUpdate, latest, local, err := sdeversion.NeedsUpdate(sqliteDB)
		if err != nil {
			log.Printf("Warning: Version check failed: %v", err)
			log.Println("→ Proceeding with update anyway")
		} else {
			log.Printf("  Latest SDE: %s", latest)
			log.Printf("  Local DB:   %s", local)

			if !needsUpdate {
				log.Println("✓ Database is up-to-date")
				return
			}
			log.Println("→ Update available, proceeding...")
		}
	} else {
		log.Println("→ Force mode enabled, skipping version check")
		// Bei --force alte DB löschen falls vorhanden
		if _, err := os.Stat(sqliteDB); err == nil {
			log.Printf("→ Removing existing database: %s", sqliteDB)
			if err := os.Remove(sqliteDB); err != nil {
				log.Fatalf("Failed to remove existing database: %v", err)
			}
		}
	}

	// 2. Download SDE
	log.Println("→ Downloading SDE data...")
	if err := runScript("scripts/download-sde.sh", *verbose); err != nil {
		log.Fatalf("Failed to download SDE: %v", err)
	}
	log.Println("✓ SDE downloaded")

	// 3. Generate Schemas
	log.Println("→ Generating Go schemas...")
	if err := runCommand("go", []string{"run", "./cmd/sde-schema-gen", "-input", jsonlDir, "-output", "internal/schema/types"}, *verbose); err != nil {
		log.Printf("Warning: Schema generation failed: %v", err)
		log.Println("→ Continuing with existing schemas")
	} else {
		log.Println("✓ Schemas generated")
	}

	// 4. Import to SQLite (optional)
	if !*skipImport {
		log.Println("→ Importing to SQLite...")
		args := []string{"run", "./cmd/sde-to-sqlite", "--db", sqliteDB, "--jsonl", jsonlDir}
		if *verbose {
			args = append(args, "-v")
		}
		if err := runCommand("go", args, *verbose); err != nil {
			log.Fatalf("Failed to import SQLite: %v", err)
		}
		log.Println("✓ SQLite import completed")
	} else {
		log.Println("→ Skipping SQLite import (--skip-import)")
	}

	// Summary
	elapsed := time.Since(start)
	log.Printf("✓ Sync completed in %s", elapsed.Round(time.Second))
}

// runScript führt ein Shell-Script aus
func runScript(scriptPath string, verbose bool) error {
	cmd := exec.Command("bash", scriptPath)
	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd.Run()
}

// runCommand führt einen Befehl aus
func runCommand(name string, args []string, verbose bool) error {
	cmd := exec.Command(name, args...)
	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		// Zeige nur Fehler
		cmd.Stderr = os.Stderr
	}
	return cmd.Run()
}
