// Code generated from JSONL data analysis
// DO NOT EDIT manually - regenerate with scripts/fetch-schemas.sh --refresh
// Source: data/jsonl/blueprints.jsonl

package types

// Blueprints represents the schema for blueprints.jsonl
// This is a simplified struct - use actual schema documentation for production
type Blueprints struct {
	Key int64 `json:"_key"`
	Activities map[string]interface{} `json:"activities"`
	Copying map[string]interface{} `json:"copying"`
	Time int64 `json:"time"`
	Manufacturing map[string]interface{} `json:"manufacturing"`
	Materials []interface{} `json:"materials"`
	Quantity int64 `json:"quantity"`
	TypeID int64 `json:"typeID"`
	Products []interface{} `json:"products"`
	Quantity int64 `json:"quantity"`
	TypeID int64 `json:"typeID"`
	Time int64 `json:"time"`
	ResearchMaterial map[string]interface{} `json:"research_material"`
	Time int64 `json:"time"`
	ResearchTime map[string]interface{} `json:"research_time"`
	Time int64 `json:"time"`
	BlueprintTypeID int64 `json:"blueprintTypeID"`
	MaxProductionLimit int64 `json:"maxProductionLimit"`
}
