// Code generated from JSONL data analysis
// DO NOT EDIT manually - regenerate with scripts/fetch-schemas.sh --refresh
// Source: data/jsonl/sovereigntyUpgrades.jsonl

package types

// SovereigntyUpgrades represents the schema for sovereigntyUpgrades.jsonl
// This is a simplified struct - use actual schema documentation for production
type SovereigntyUpgrades struct {
	Key int64 `json:"_key"`
	FuelHourlyUpkeep int64 `json:"fuel_hourly_upkeep"`
	FuelStartupCost int64 `json:"fuel_startup_cost"`
	FuelTypeId int64 `json:"fuel_type_id"`
	MutuallyExclusiveGroup string `json:"mutually_exclusive_group"`
	PowerAllocation int64 `json:"power_allocation"`
	WorkforceAllocation int64 `json:"workforce_allocation"`
}
