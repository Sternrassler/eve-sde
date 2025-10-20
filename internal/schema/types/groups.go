// Code generated from JSONL data analysis
// DO NOT EDIT manually - regenerate with scripts/fetch-schemas.sh --refresh
// Source: data/jsonl/groups.jsonl

package types

// Groups represents the schema for groups.jsonl
// This is a simplified struct - use actual schema documentation for production
type Groups struct {
	Key int64 `json:"_key"`
	Anchorable bool `json:"anchorable"`
	Anchored bool `json:"anchored"`
	CategoryID int64 `json:"categoryID"`
	FittableNonSingleton bool `json:"fittableNonSingleton"`
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
	UseBasePrice bool `json:"useBasePrice"`
}
