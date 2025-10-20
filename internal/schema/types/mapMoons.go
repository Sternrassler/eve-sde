// Code generated from JSONL data analysis
// DO NOT EDIT manually - regenerate with scripts/fetch-schemas.sh --refresh
// Source: data/jsonl/mapMoons.jsonl

package types

// MapMoons represents the schema for mapMoons.jsonl
// This is a simplified struct - use actual schema documentation for production
type MapMoons struct {
	Key int64 `json:"_key"`
	Attributes map[string]interface{} `json:"attributes"`
	HeightMap1 int64 `json:"heightMap1"`
	HeightMap2 int64 `json:"heightMap2"`
	ShaderPreset int64 `json:"shaderPreset"`
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
}
