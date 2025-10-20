// Code generated from JSONL data analysis
// DO NOT EDIT manually - regenerate with scripts/fetch-schemas.sh --refresh
// Source: data/jsonl/npcCorporationDivisions.jsonl

package types

// NpcCorporationDivisions represents the schema for npcCorporationDivisions.jsonl
// This is a simplified struct - use actual schema documentation for production
type NpcCorporationDivisions struct {
	Key int64 `json:"_key"`
	DisplayName string `json:"displayName"`
	InternalName string `json:"internalName"`
	LeaderTypeName map[string]interface{} `json:"leaderTypeName"`
	De string `json:"de"`
	En string `json:"en"`
	Es string `json:"es"`
	Fr string `json:"fr"`
	Ja string `json:"ja"`
	Ko string `json:"ko"`
	Ru string `json:"ru"`
	Zh string `json:"zh"`
	Name map[string]interface{} `json:"name"`
	De string `json:"de"`
	En string `json:"en"`
	Es string `json:"es"`
	Fr string `json:"fr"`
	Ja string `json:"ja"`
	Ko string `json:"ko"`
	Ru string `json:"ru"`
}
