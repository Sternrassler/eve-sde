// Code generated from JSONL data analysis
// DO NOT EDIT manually - regenerate with scripts/fetch-schemas.sh --refresh
// Source: data/jsonl/mapRegions.jsonl

package types

// MapRegions represents the schema for mapRegions.jsonl
// This is a simplified struct - use actual schema documentation for production
type MapRegions struct {
	Key int64 `json:"_key"`
	ConstellationIDs []interface{} `json:"constellationIDs"`
	Description map[string]interface{} `json:"description"`
	De interface{} `json:"de"`
	En interface{} `json:"en"`
	Es interface{} `json:"es"`
	Fr interface{} `json:"fr"`
	Ja string `json:"ja"`
	Ko interface{} `json:"ko"`
	Ru interface{} `json:"ru"`
	Zh string `json:"zh"`
	FactionID int64 `json:"factionID"`
	Name map[string]interface{} `json:"name"`
	De interface{} `json:"de"`
	En interface{} `json:"en"`
	Es interface{} `json:"es"`
	Fr interface{} `json:"fr"`
	Ja string `json:"ja"`
	Ko interface{} `json:"ko"`
	Ru interface{} `json:"ru"`
}
