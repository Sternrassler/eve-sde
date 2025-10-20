// Code generated from JSONL data analysis
// DO NOT EDIT manually - regenerate with scripts/fetch-schemas.sh --refresh
// Source: data/jsonl/characterAttributes.jsonl

package types

// CharacterAttributes represents the schema for characterAttributes.jsonl
// This is a simplified struct - use actual schema documentation for production
type CharacterAttributes struct {
	Key int64 `json:"_key"`
	Description interface{} `json:"description"`
	IconID int64 `json:"iconID"`
	Name map[string]interface{} `json:"name"`
	De string `json:"de"`
	En string `json:"en"`
	Es string `json:"es"`
	Fr string `json:"fr"`
	Ja string `json:"ja"`
	Ko string `json:"ko"`
	Ru string `json:"ru"`
	Zh string `json:"zh"`
	Notes interface{} `json:"notes"`
	ShortDescription interface{} `json:"shortDescription"`
}
