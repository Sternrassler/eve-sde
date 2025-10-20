// Code generated from JSONL data analysis
// DO NOT EDIT manually - regenerate with scripts/fetch-schemas.sh --refresh
// Source: data/jsonl/npcStations.jsonl

package types

// NpcStations represents the schema for npcStations.jsonl
// This is a simplified struct - use actual schema documentation for production
type NpcStations struct {
	Key int64 `json:"_key"`
	CelestialIndex int64 `json:"celestialIndex"`
	OperationID int64 `json:"operationID"`
	OrbitID int64 `json:"orbitID"`
	OrbitIndex int64 `json:"orbitIndex"`
	OwnerID int64 `json:"ownerID"`
	Position map[string]interface{} `json:"position"`
	X int64 `json:"x"`
	Y int64 `json:"y"`
	Z interface{} `json:"z"`
	ReprocessingEfficiency int64 `json:"reprocessingEfficiency"`
	ReprocessingHangarFlag int64 `json:"reprocessingHangarFlag"`
	ReprocessingStationsTake int64 `json:"reprocessingStationsTake"`
	SolarSystemID int64 `json:"solarSystemID"`
	TypeID int64 `json:"typeID"`
	UseOperationName bool `json:"useOperationName"`
}
