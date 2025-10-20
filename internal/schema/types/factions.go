// Code generated from JSONL data analysis
// DO NOT EDIT manually - regenerate with scripts/fetch-schemas.sh --refresh
// Source: data/jsonl/factions.jsonl

package types

// Factions represents the schema for factions.jsonl
// This is a simplified struct - use actual schema documentation for production
type Factions struct {
	Key int64 `json:"_key"`
	CorporationID int64 `json:"corporationID"`
	Description map[string]interface{} `json:"description"`
	De interface{} `json:"de"`
	En interface{} `json:"en"`
	Es interface{} `json:"es"`
	Fr interface{} `json:"fr"`
	Ja string `json:"ja"`
	Ko interface{} `json:"ko"`
	Ru interface{} `json:"ru"`
	Zh string `json:"zh"`
	FlatLogo string `json:"flatLogo"`
	FlatLogoWithName string `json:"flatLogoWithName"`
	IconID int64 `json:"iconID"`
	MemberRaces []interface{} `json:"memberRaces"`
	MilitiaCorporationID int64 `json:"militiaCorporationID"`
	Name map[string]interface{} `json:"name"`
	De interface{} `json:"de"`
	En interface{} `json:"en"`
	Es interface{} `json:"es"`
}
