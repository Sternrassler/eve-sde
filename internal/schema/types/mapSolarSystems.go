// Code generated from JSONL data analysis
// DO NOT EDIT manually - regenerate with scripts/fetch-schemas.sh --refresh
// Source: data/jsonl/mapSolarSystems.jsonl

package types

// MapSolarSystems represents the schema for mapSolarSystems.jsonl
// This is a simplified struct - use actual schema documentation for production
type MapSolarSystems struct {
	Key int64 `json:"_key"`
	Border bool `json:"border"`
	ConstellationID int64 `json:"constellationID"`
	Hub bool `json:"hub"`
	International bool `json:"international"`
	Luminosity int64 `json:"luminosity"`
	Name map[string]interface{} `json:"name"`
	De string `json:"de"`
	En string `json:"en"`
	Es string `json:"es"`
	Fr string `json:"fr"`
	Ja string `json:"ja"`
	Ko string `json:"ko"`
	Ru string `json:"ru"`
	Zh string `json:"zh"`
	PlanetIDs []interface{} `json:"planetIDs"`
	Position map[string]interface{} `json:"position"`
	X interface{} `json:"x"`
	Y int64 `json:"y"`
	Z interface{} `json:"z"`
}
