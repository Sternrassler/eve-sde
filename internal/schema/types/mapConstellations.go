// Code generated from JSONL data analysis
// DO NOT EDIT manually - regenerate with scripts/fetch-schemas.sh --refresh
// Source: data/jsonl/mapConstellations.jsonl

package types

// MapConstellations represents the schema for mapConstellations.jsonl
// This is a simplified struct - use actual schema documentation for production
type MapConstellations struct {
	Key int64 `json:"_key"`
	FactionID int64 `json:"factionID"`
	Name map[string]interface{} `json:"name"`
	De string `json:"de"`
	En string `json:"en"`
	Es string `json:"es"`
	Fr string `json:"fr"`
	Ja string `json:"ja"`
	Ko string `json:"ko"`
	Ru string `json:"ru"`
	Zh string `json:"zh"`
	Position map[string]interface{} `json:"position"`
	X interface{} `json:"x"`
	Y int64 `json:"y"`
	Z interface{} `json:"z"`
	RegionID int64 `json:"regionID"`
	SolarSystemIDs []interface{} `json:"solarSystemIDs"`
	WormholeClassID int64 `json:"wormholeClassID"`
}
