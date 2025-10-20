package main

import (
	"fmt"
	"os"

	sdeversion "github.com/Sternrassler/eve-sde/internal/sde/version"
)

func main() {
	latest, err := sdeversion.GetLatestVersion()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching latest version: %v\n", err)
		os.Exit(1)
	}

	local, _ := sdeversion.GetLocalVersion("data/sqlite/eve-sde.db")

	fmt.Printf("LATEST_BUILD=%d\n", latest.BuildNumber)
	fmt.Printf("LATEST_DATE=%s\n", latest.ReleaseDate.Format("2006-01-02"))

	if local != nil {
		fmt.Printf("LOCAL_BUILD=%d\n", local.BuildNumber)
		fmt.Printf("NEEDS_UPDATE=%t\n", latest.BuildNumber != local.BuildNumber)
	} else {
		fmt.Printf("LOCAL_BUILD=0\n")
		fmt.Printf("NEEDS_UPDATE=true\n")
	}
}
