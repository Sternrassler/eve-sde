package main

import (
	"fmt"
	"log"
	"os"
)

const version = "0.1.0"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Printf("sde-to-sqlite v%s\n", version)
		os.Exit(0)
	}

	log.Println("EVE SDE to SQLite Converter")
	log.Println("Version:", version)
	log.Println("Status: In Development")

	// TODO: Implement JSONL → SQLite conversion
	log.Fatal("Not yet implemented")
}
