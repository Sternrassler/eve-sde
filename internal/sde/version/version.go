package version

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	// LatestVersionURL ist die offizielle CCP URL für SDE Versionsinformationen
	LatestVersionURL = "https://developers.eveonline.com/static-data/tranquility/latest.jsonl"
)

// SDEVersion repräsentiert SDE Versionsinformationen
type SDEVersion struct {
	Key         string    `json:"_key"`
	BuildNumber int64     `json:"buildNumber"`
	ReleaseDate time.Time `json:"releaseDate"`
}

// GetLatestVersion holt die aktuelle SDE-Version von CCP
func GetLatestVersion() (*SDEVersion, error) {
	resp, err := http.Get(LatestVersionURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest version: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var version SDEVersion
	if err := json.Unmarshal(body, &version); err != nil {
		return nil, fmt.Errorf("failed to parse version JSON: %w", err)
	}

	return &version, nil
}

// GetLocalVersion liest die lokal gespeicherte SDE-Version aus SQLite
func GetLocalVersion(dbPath string) (*SDEVersion, error) {
	// Check ob DB existiert
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return nil, nil // DB existiert nicht = keine Version
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	var version SDEVersion
	var releaseDateStr string

	err = db.QueryRow("SELECT _key, buildNumber, releaseDate FROM _sde LIMIT 1").
		Scan(&version.Key, &version.BuildNumber, &releaseDateStr)

	if err == sql.ErrNoRows {
		return nil, nil // Keine Version vorhanden
	}
	if err != nil {
		// Tabelle existiert nicht = DB leer/ungültig
		if isNoSuchTableError(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query _sde table: %w", err)
	}

	// Parse releaseDate (stored as TEXT in SQLite)
	version.ReleaseDate, err = time.Parse(time.RFC3339, releaseDateStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse release date: %w", err)
	}

	return &version, nil
}

// isNoSuchTableError prüft ob Fehler "no such table" ist
func isNoSuchTableError(err error) bool {
	if err == nil {
		return false
	}
	return err.Error() == "no such table: _sde" ||
		err.Error() == "SQL logic error: no such table: _sde"
}

// NeedsUpdate prüft ob ein Update verfügbar ist
func NeedsUpdate(dbPath string) (bool, *SDEVersion, *SDEVersion, error) {
	latest, err := GetLatestVersion()
	if err != nil {
		return false, nil, nil, fmt.Errorf("failed to get latest version: %w", err)
	}

	local, err := GetLocalVersion(dbPath)
	if err != nil {
		return false, nil, nil, fmt.Errorf("failed to get local version: %w", err)
	}

	// Keine lokale Version = Update nötig
	if local == nil {
		return true, latest, nil, nil
	}

	// Vergleiche BuildNumber
	needsUpdate := latest.BuildNumber > local.BuildNumber

	return needsUpdate, latest, local, nil
}

// String implementiert Stringer für SDEVersion
func (v *SDEVersion) String() string {
	if v == nil {
		return "none"
	}
	return fmt.Sprintf("Build %d (%s)", v.BuildNumber, v.ReleaseDate.Format("2006-01-02"))
}
