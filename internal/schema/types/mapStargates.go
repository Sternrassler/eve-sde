// Code generated from JSONL data analysis
// DO NOT EDIT manually - regenerate with scripts/fetch-schemas.sh --refresh
// Source: data/jsonl/mapStargates.jsonl

package types

// MapStargates represents the schema for mapStargates.jsonl
// This is a simplified struct - use actual schema documentation for production
type MapStargates struct {
	Key int64 `json:"_key"`
	Destination map[string]interface{} `json:"destination"`
	SolarSystemID int64 `json:"solarSystemID"`
	StargateID int64 `json:"stargateID"`
	Position map[string]interface{} `json:"position"`
	X int64 `json:"x"`
	Y int64 `json:"y"`
	Z interface{} `json:"z"`
	SolarSystemID int64 `json:"solarSystemID"`
	TypeID int64 `json:"typeID"`
}
