// Code generated from JSONL data analysis
// DO NOT EDIT manually - regenerate with scripts/fetch-schemas.sh --refresh
// Source: data/jsonl/npcCorporations.jsonl

package types

// NpcCorporations represents the schema for npcCorporations.jsonl
// This is a simplified struct - use actual schema documentation for production
type NpcCorporations struct {
	Key int64 `json:"_key"`
	CeoID int64 `json:"ceoID"`
	Deleted bool `json:"deleted"`
	Description map[string]interface{} `json:"description"`
	De interface{} `json:"de"`
	En string `json:"en"`
	Es string `json:"es"`
	Fr string `json:"fr"`
	Ja string `json:"ja"`
	Ko string `json:"ko"`
	Ru interface{} `json:"ru"`
	Zh string `json:"zh"`
	Extent string `json:"extent"`
	HasPlayerPersonnelManager bool `json:"hasPlayerPersonnelManager"`
	InitialPrice int64 `json:"initialPrice"`
	MemberLimit interface{} `json:"memberLimit"`
	MinSecurity int64 `json:"minSecurity"`
	MinimumJoinStanding int64 `json:"minimumJoinStanding"`
	Name map[string]interface{} `json:"name"`
	De interface{} `json:"de"`
}
