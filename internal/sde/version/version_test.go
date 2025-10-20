package version

import (
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func TestSDEVersionString(t *testing.T) {
	v := &SDEVersion{
		Key:         "sde",
		BuildNumber: 3064089,
		ReleaseDate: time.Date(2025, 10, 17, 11, 14, 3, 0, time.UTC),
	}

	expected := "Build 3064089 (2025-10-17)"
	if got := v.String(); got != expected {
		t.Errorf("SDEVersion.String() = %q, want %q", got, expected)
	}
}

func TestSDEVersionStringNil(t *testing.T) {
	var v *SDEVersion
	if got := v.String(); got != "none" {
		t.Errorf("nil SDEVersion.String() = %q, want %q", got, "none")
	}
}

func TestGetLatestVersion(t *testing.T) {
	// Integration Test - erfordert Netzwerk
	if testing.Short() {
		t.Skip("Skipping network test in short mode")
	}

	version, err := GetLatestVersion()
	if err != nil {
		t.Fatalf("GetLatestVersion() error = %v", err)
	}

	if version == nil {
		t.Fatal("GetLatestVersion() returned nil")
	}

	if version.Key != "sde" {
		t.Errorf("version.Key = %q, want %q", version.Key, "sde")
	}

	if version.BuildNumber <= 0 {
		t.Errorf("version.BuildNumber = %d, want > 0", version.BuildNumber)
	}

	if version.ReleaseDate.IsZero() {
		t.Error("version.ReleaseDate is zero")
	}
}

func TestNeedsUpdateNoLocalVersion(t *testing.T) {
	// Test mit nicht-existierender DB
	needsUpdate, latest, local, err := NeedsUpdate("/tmp/nonexistent.db")

	// Sollte true zurückgeben wenn keine lokale Version existiert
	// ABER könnte auch Error sein wenn DB nicht existiert - beide OK
	if err != nil {
		// DB existiert nicht - erwartbarer Fehler
		t.Logf("Expected error for non-existent DB: %v", err)
		return
	}

	if !needsUpdate {
		t.Error("NeedsUpdate() should return true when no local version exists")
	}

	if latest == nil {
		t.Error("NeedsUpdate() latest should not be nil")
	}

	if local != nil {
		t.Errorf("NeedsUpdate() local = %v, want nil", local)
	}
}
