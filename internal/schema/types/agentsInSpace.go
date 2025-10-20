// Code generated from JSONL data analysis
// DO NOT EDIT manually - regenerate with scripts/fetch-schemas.sh --refresh
// Source: data/jsonl/agentsInSpace.jsonl

package types

// AgentsInSpace represents the schema for agentsInSpace.jsonl
// This is a simplified struct - use actual schema documentation for production
type AgentsInSpace struct {
	Key int64 `json:"_key"`
	DungeonID int64 `json:"dungeonID"`
	SolarSystemID int64 `json:"solarSystemID"`
	SpawnPointID int64 `json:"spawnPointID"`
	TypeID int64 `json:"typeID"`
}
