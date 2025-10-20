// Code generated from JSONL data analysis
// DO NOT EDIT manually - regenerate with scripts/fetch-schemas.sh --refresh
// Source: data/jsonl/dogmaEffects.jsonl

package types

// DogmaEffects represents the schema for dogmaEffects.jsonl
// This is a simplified struct - use actual schema documentation for production
type DogmaEffects struct {
	Key int64 `json:"_key"`
	DisallowAutoRepeat bool `json:"disallowAutoRepeat"`
	DischargeAttributeID int64 `json:"dischargeAttributeID"`
	DurationAttributeID int64 `json:"durationAttributeID"`
	EffectCategoryID int64 `json:"effectCategoryID"`
	ElectronicChance bool `json:"electronicChance"`
	Guid string `json:"guid"`
	IsAssistance bool `json:"isAssistance"`
	IsOffensive bool `json:"isOffensive"`
	IsWarpSafe bool `json:"isWarpSafe"`
	Name string `json:"name"`
	PropulsionChance bool `json:"propulsionChance"`
	Published bool `json:"published"`
	RangeChance bool `json:"rangeChance"`
}
