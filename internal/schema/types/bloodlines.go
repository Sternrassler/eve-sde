// Code generated from JSONL data analysis
// DO NOT EDIT manually - regenerate with scripts/fetch-schemas.sh --refresh
// Source: data/jsonl/bloodlines.jsonl

package types

// Bloodlines represents the schema for bloodlines.jsonl
// This is a simplified struct - use actual schema documentation for production
type Bloodlines struct {
	Key int64 `json:"_key"`
	Charisma int64 `json:"charisma"`
	CorporationID int64 `json:"corporationID"`
	Description map[string]interface{} `json:"description"`
	De interface{} `json:"de"`
	En interface{} `json:"en"`
	Es interface{} `json:"es"`
	Fr interface{} `json:"fr"`
	Ja string `json:"ja"`
	Ko string `json:"ko"`
	Ru interface{} `json:"ru"`
	Zh string `json:"zh"`
	IconID int64 `json:"iconID"`
	Intelligence int64 `json:"intelligence"`
	Memory int64 `json:"memory"`
	Name map[string]interface{} `json:"name"`
	De interface{} `json:"de"`
	En interface{} `json:"en"`
	Es interface{} `json:"es"`
	Fr interface{} `json:"fr"`
}
