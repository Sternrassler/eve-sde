// Code generated from JSONL data analysis
// DO NOT EDIT manually - regenerate with scripts/fetch-schemas.sh --refresh
// Source: data/jsonl/mapPlanets.jsonl

package types

// MapPlanets represents the schema for mapPlanets.jsonl
// This is a simplified struct - use actual schema documentation for production
type MapPlanets struct {
	Key int64 `json:"_key"`
	AsteroidBeltIDs []interface{} `json:"asteroidBeltIDs"`
	Attributes map[string]interface{} `json:"attributes"`
	HeightMap1 int64 `json:"heightMap1"`
	HeightMap2 int64 `json:"heightMap2"`
	Population bool `json:"population"`
	ShaderPreset int64 `json:"shaderPreset"`
	CelestialIndex int64 `json:"celestialIndex"`
	MoonIDs []interface{} `json:"moonIDs"`
	OrbitID int64 `json:"orbitID"`
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
}
