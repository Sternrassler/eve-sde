// Code generated from JSONL data analysis
// DO NOT EDIT manually - regenerate with scripts/fetch-schemas.sh --refresh
// Source: data/jsonl/planetSchematics.jsonl

package types

// PlanetSchematics represents the schema for planetSchematics.jsonl
// This is a simplified struct - use actual schema documentation for production
type PlanetSchematics struct {
	Key int64 `json:"_key"`
	CycleTime int64 `json:"cycleTime"`
	Name map[string]interface{} `json:"name"`
	De string `json:"de"`
	En string `json:"en"`
	Es string `json:"es"`
	Fr string `json:"fr"`
	Ja string `json:"ja"`
	Ko string `json:"ko"`
	Ru string `json:"ru"`
	Zh string `json:"zh"`
	Pins []interface{} `json:"pins"`
	Types []interface{} `json:"types"`
	Key int64 `json:"_key"`
	IsInput bool `json:"isInput"`
	Quantity int64 `json:"quantity"`
	Key int64 `json:"_key"`
	IsInput bool `json:"isInput"`
	Quantity int64 `json:"quantity"`
	Key int64 `json:"_key"`
}
