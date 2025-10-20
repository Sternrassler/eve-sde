// Code generated from JSONL data analysis
// DO NOT EDIT manually - regenerate with scripts/fetch-schemas.sh --refresh
// Source: data/jsonl/metaGroups.jsonl

package types

// MetaGroups represents the schema for metaGroups.jsonl
// This is a simplified struct - use actual schema documentation for production
type MetaGroups struct {
	Key int64 `json:"_key"`
	Color map[string]interface{} `json:"color"`
	B int64 `json:"b"`
	G int64 `json:"g"`
	R int64 `json:"r"`
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
