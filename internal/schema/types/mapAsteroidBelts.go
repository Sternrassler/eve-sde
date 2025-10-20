// Code generated from JSONL data analysis
// DO NOT EDIT manually - regenerate with scripts/fetch-schemas.sh --refresh
// Source: data/jsonl/mapAsteroidBelts.jsonl

package types

// MapAsteroidBelts represents the schema for mapAsteroidBelts.jsonl
// This is a simplified struct - use actual schema documentation for production
type MapAsteroidBelts struct {
	Key int64 `json:"_key"`
	CelestialIndex int64 `json:"celestialIndex"`
	OrbitID int64 `json:"orbitID"`
	OrbitIndex int64 `json:"orbitIndex"`
	Position map[string]interface{} `json:"position"`
	X int64 `json:"x"`
	Y int64 `json:"y"`
	Z interface{} `json:"z"`
	Radius int64 `json:"radius"`
	SolarSystemID int64 `json:"solarSystemID"`
	Statistics map[string]interface{} `json:"statistics"`
	Density int64 `json:"density"`
	Eccentricity int64 `json:"eccentricity"`
	EscapeVelocity int64 `json:"escapeVelocity"`
	Locked bool `json:"locked"`
	MassDust int64 `json:"massDust"`
	MassGas int64 `json:"massGas"`
	OrbitPeriod int64 `json:"orbitPeriod"`
	OrbitRadius int64 `json:"orbitRadius"`
	RotationRate int64 `json:"rotationRate"`
}
