// Package views provides SQL view initialization for the eve-sde databasepackage views

// This package is part of the DB-core and contains only SQL view definitions
package views

import (
	"database/sql"
	_ "embed"
	"fmt"
)

//go:embed navigation.sql
var navigationViewsSQL string

// InitializeNavigationViews creates all navigation-related views in the database
// This should be called after map data (mapSolarSystems, mapStargates) has been imported
func InitializeNavigationViews(db *sql.DB) error {
	_, err := db.Exec(navigationViewsSQL)
	if err != nil {
		return fmt.Errorf("failed to initialize navigation views: %w", err)
	}
	return nil
}
