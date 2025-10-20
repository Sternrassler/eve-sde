// Code generated from JSONL data analysis
// DO NOT EDIT manually - regenerate with scripts/fetch-schemas.sh --refresh
// Source: data/jsonl/categories.jsonl

package types

// Categories represents the schema for categories.jsonl
// This is a simplified struct - use actual schema documentation for production
type Categories struct {
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
	Published bool `json:"published"`
}
