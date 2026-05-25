package registry

// ExcludedDatasets listet Datensätze, die bewusst NICHT importiert werden.
// Key = JSONL-Basisname ohne .jsonl, Value = Begründung (erscheint im Log).
// Leer = alle Datensätze werden importiert.
var ExcludedDatasets = map[string]string{}

// IndexOverrides ersetzt die Standard-Index-Konvention (alle *ID-Felder außer _key)
// für die genannten Datensätze. Nötig für Index-Spalten ohne "ID"-Suffix.
var IndexOverrides = map[string][]string{
	"mapSolarSystems": {"constellationID", "securityClass"},
	"mapStargates":    {"solarSystemID", "destination"},
}
