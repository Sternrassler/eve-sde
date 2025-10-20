// Code generated from JSONL data analysis
// DO NOT EDIT manually - regenerate with scripts/fetch-schemas.sh --refresh
// Source: data/jsonl/corporationActivities.jsonl

package types

// CorporationActivities represents the schema for corporationActivities.jsonl
// This is a simplified struct - use actual schema documentation for production
type CorporationActivities struct {
	Key int64 `json:"_key"`
	Name map[string]interface{} `json:"name"`
	De string `json:"de"`
	En string `json:"en"`
	Es string `json:"es"`
	Fr string `json:"fr"`
	Ja string `json:"ja"`
	Ko string `json:"ko"`
	Ru string `json:"ru"`
	Zh string `json:"zh"`
}
