// Code generated from JSONL data analysis
// DO NOT EDIT manually - regenerate with scripts/fetch-schemas.sh --refresh
// Source: data/jsonl/mapStars.jsonl

package types

// MapStars represents the schema for mapStars.jsonl
// This is a simplified struct - use actual schema documentation for production
type MapStars struct {
	Key int64 `json:"_key"`
	Radius int64 `json:"radius"`
	SolarSystemID int64 `json:"solarSystemID"`
	Statistics map[string]interface{} `json:"statistics"`
	Age int64 `json:"age"`
	Life int64 `json:"life"`
	Luminosity int64 `json:"luminosity"`
	SpectralClass string `json:"spectralClass"`
	Temperature int64 `json:"temperature"`
	TypeID int64 `json:"typeID"`
}
