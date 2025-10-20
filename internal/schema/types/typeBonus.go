// Code generated from JSONL data analysis
// DO NOT EDIT manually - regenerate with scripts/fetch-schemas.sh --refresh
// Source: data/jsonl/typeBonus.jsonl

package types

// TypeBonus represents the schema for typeBonus.jsonl
// This is a simplified struct - use actual schema documentation for production
type TypeBonus struct {
	Key int64 `json:"_key"`
	RoleBonuses []interface{} `json:"roleBonuses"`
	Bonus int64 `json:"bonus"`
	BonusText map[string]interface{} `json:"bonusText"`
	De string `json:"de"`
	En string `json:"en"`
	Es string `json:"es"`
	Fr string `json:"fr"`
	Ja string `json:"ja"`
	Ko string `json:"ko"`
	Ru string `json:"ru"`
	Zh string `json:"zh"`
	Importance int64 `json:"importance"`
	UnitID int64 `json:"unitID"`
	Types []interface{} `json:"types"`
	Key int64 `json:"_key"`
	Value []interface{} `json:"_value"`
	Bonus int64 `json:"bonus"`
	BonusText map[string]interface{} `json:"bonusText"`
	De string `json:"de"`
}
