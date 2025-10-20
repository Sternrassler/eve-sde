// Code generated from JSONL data analysis
// DO NOT EDIT manually - regenerate with scripts/fetch-schemas.sh --refresh
// Source: data/jsonl/marketGroups.jsonl

package types

// MarketGroups represents the schema for marketGroups.jsonl
// This is a simplified struct - use actual schema documentation for production
type MarketGroups struct {
	Key int64 `json:"_key"`
	Description map[string]interface{} `json:"description"`
	De interface{} `json:"de"`
	En interface{} `json:"en"`
	Es interface{} `json:"es"`
	Fr interface{} `json:"fr"`
	Ja string `json:"ja"`
	Ko interface{} `json:"ko"`
	Ru interface{} `json:"ru"`
	Zh string `json:"zh"`
	HasTypes bool `json:"hasTypes"`
	IconID int64 `json:"iconID"`
	Name map[string]interface{} `json:"name"`
	De interface{} `json:"de"`
	En interface{} `json:"en"`
	Es interface{} `json:"es"`
	Fr interface{} `json:"fr"`
	Ja string `json:"ja"`
	Ko interface{} `json:"ko"`
	Ru interface{} `json:"ru"`
}
