// Code generated from JSONL data analysis
// DO NOT EDIT manually - regenerate with scripts/fetch-schemas.sh --refresh
// Source: data/jsonl/npcCharacters.jsonl

package types

// NpcCharacters represents the schema for npcCharacters.jsonl
// This is a simplified struct - use actual schema documentation for production
type NpcCharacters struct {
	Key int64 `json:"_key"`
	BloodlineID int64 `json:"bloodlineID"`
	Ceo bool `json:"ceo"`
	CorporationID int64 `json:"corporationID"`
	Gender bool `json:"gender"`
	LocationID int64 `json:"locationID"`
	Name map[string]interface{} `json:"name"`
	De string `json:"de"`
	En string `json:"en"`
	Es string `json:"es"`
	Fr string `json:"fr"`
	Ja string `json:"ja"`
	Ko string `json:"ko"`
	Ru string `json:"ru"`
	Zh string `json:"zh"`
	RaceID int64 `json:"raceID"`
	StartDate string `json:"startDate"`
	UniqueName bool `json:"uniqueName"`
}
