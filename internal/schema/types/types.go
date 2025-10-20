// Code generated from JSONL data analysis
// DO NOT EDIT manually - regenerate with scripts/fetch-schemas.sh --refresh
// Source: data/jsonl/types.jsonl

package types

// Types represents the schema for types.jsonl
// This is a simplified struct - use actual schema documentation for production
type Types struct {
	Key int64 `json:"_key"`
	GroupID int64 `json:"groupID"`
	Mass int64 `json:"mass"`
	Name map[string]interface{} `json:"name"`
	De string `json:"de"`
	En string `json:"en"`
	Es string `json:"es"`
	Fr string `json:"fr"`
	Ja string `json:"ja"`
	Ko string `json:"ko"`
	Ru string `json:"ru"`
	Zh string `json:"zh"`
	PortionSize int64 `json:"portionSize"`
	Published bool `json:"published"`
}
