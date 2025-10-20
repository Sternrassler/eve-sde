// Package navigation provides EVE Online navigation and route planning functionality
// It includes pathfinding, travel time calculation, and trade hub analysis
package navigation

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"math"
)

//go:embed views.sql
var viewsSQL string

// NavigationParams contains optional parameters for route calculation
type NavigationParams struct {
	WarpSpeed       *float64 `json:"warp_speed,omitempty"`        // AU/s (default: 3.0)
	AlignTime       *float64 `json:"align_time,omitempty"`        // seconds (default: 6.0)
	ShipMass        *float64 `json:"ship_mass,omitempty"`         // kg (for exact align calc)
	InertiaModifier *float64 `json:"inertia_modifier,omitempty"`  // ship agility
	AvgWarpDistance *float64 `json:"avg_warp_distance,omitempty"` // AU per system (default: 15)
	AvoidLowSec     bool     `json:"avoid_lowsec"`                // route via high-sec only
}

// RouteResult contains calculated route information
type RouteResult struct {
	TotalSeconds       float64               `json:"total_seconds"`
	TotalMinutes       float64               `json:"total_minutes"`
	Jumps              int                   `json:"jumps"`
	AvgSecondsPerJump  float64               `json:"avg_seconds_per_jump"`
	Route              []int64               `json:"route"`
	ParametersUsed     map[string]interface{} `json:"parameters_used"`
}

// PathResult contains just the path information
type PathResult struct {
	FromSystemID int64   `json:"from_system_id"`
	ToSystemID   int64   `json:"to_system_id"`
	Jumps        int     `json:"jumps"`
	Route        []int64 `json:"route"`
}

// Default navigation parameters
const (
	DefaultWarpSpeed       = 3.0  // AU/s (Cruiser average)
	DefaultAlignTime       = 6.0  // seconds (medium ships)
	DefaultGateJumpDelay   = 10.0 // seconds (gate jump animation + loading)
	DefaultAvgWarpDistance = 15.0 // AU (statistical assumption)
	WarpCorrectionFactor   = 1.4  // Simplified warp time correction
)

// InitializeViews creates all navigation views in the database
func InitializeViews(db *sql.DB) error {
	_, err := db.Exec(viewsSQL)
	if err != nil {
		return fmt.Errorf("failed to initialize navigation views: %w", err)
	}
	return nil
}

// CalculateAlignTime calculates exact align time from ship parameters
// Formula: align_time = 1.386 * inertia_modifier * mass / 500000
func CalculateAlignTime(mass, inertiaModifier float64) float64 {
	return 1.386 * inertiaModifier * mass / 500000.0
}

// CalculateWarpTime calculates warp time using CCP's 3-phase formula
// Reference: https://wiki.eveuniversity.org/Warp_time_calculation
func CalculateWarpTime(distanceAU, warpSpeedAU float64) float64 {
	const AU = 149597870700.0 // meters in 1 AU
	
	k := warpSpeedAU
	j := math.Min(k/3.0, 2.0)
	
	// Phase 1: Acceleration (always 1 AU)
	tAccel := 25.7312 / k
	
	// Phase 2: Deceleration
	dDecel := (k * AU) / j
	tDecel := (math.Log(k*AU) - math.Log(100)) / j // 100 m/s dropout speed
	
	// Phase 3: Cruise
	totalDistanceMeters := distanceAU * AU
	dCruise := math.Max(0, totalDistanceMeters-AU-dDecel)
	tCruise := dCruise / (k * AU)
	
	return tAccel + tCruise + tDecel
}

// CalculateSimplifiedWarpTime uses simplified formula for quick estimation
func CalculateSimplifiedWarpTime(distanceAU, warpSpeedAU float64) float64 {
	return (distanceAU / warpSpeedAU) * WarpCorrectionFactor
}

// getEffectiveParams returns effective parameters with defaults applied
func getEffectiveParams(params *NavigationParams) (warpSpeed, alignTime, avgWarpDist float64, source string) {
	source = "default"
	warpSpeed = DefaultWarpSpeed
	alignTime = DefaultAlignTime
	avgWarpDist = DefaultAvgWarpDistance
	
	if params == nil {
		return
	}
	
	if params.WarpSpeed != nil {
		warpSpeed = *params.WarpSpeed
		source = "provided"
	}
	
	if params.AlignTime != nil {
		alignTime = *params.AlignTime
		source = "provided"
	} else if params.ShipMass != nil && params.InertiaModifier != nil {
		alignTime = CalculateAlignTime(*params.ShipMass, *params.InertiaModifier)
		source = "calculated"
	}
	
	if params.AvgWarpDistance != nil {
		avgWarpDist = *params.AvgWarpDistance
	}
	
	return
}

// ShortestPath finds the shortest path between two systems using recursive CTE
func ShortestPath(db *sql.DB, fromSystemID, toSystemID int64, avoidLowSec bool) (*PathResult, error) {
	query := `
		WITH RECURSIVE path AS (
			SELECT 
				from_system_id,
				to_system_id,
				1 AS jumps,
				json_array(from_system_id, to_system_id) AS route
			FROM v_stargate_graph
			WHERE from_system_id = ?
			
			UNION ALL
			
			SELECT 
				p.from_system_id,
				g.to_system_id,
				p.jumps + 1,
				json_insert(p.route, '$[#]', g.to_system_id)
			FROM path p
			JOIN v_stargate_graph g ON p.to_system_id = g.from_system_id
			LEFT JOIN mapSolarSystems sys ON g.to_system_id = sys._key
			WHERE p.jumps < 100
				AND NOT EXISTS (
					SELECT 1 FROM json_each(p.route) WHERE value = g.to_system_id
				)
				AND (? = 0 OR sys.securityStatus >= 0.45)
		)
		SELECT from_system_id, to_system_id, jumps, route
		FROM path 
		WHERE to_system_id = ? 
		ORDER BY jumps ASC
		LIMIT 1
	`
	
	avoidLowSecInt := 0
	if avoidLowSec {
		avoidLowSecInt = 1
	}
	
	var result PathResult
	var routeJSON string
	
	err := db.QueryRow(query, fromSystemID, avoidLowSecInt, toSystemID).Scan(
		&result.FromSystemID,
		&result.ToSystemID,
		&result.Jumps,
		&routeJSON,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no path found between systems %d and %d", fromSystemID, toSystemID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find path: %w", err)
	}
	
	// Parse route JSON
	if err := json.Unmarshal([]byte(routeJSON), &result.Route); err != nil {
		return nil, fmt.Errorf("failed to parse route JSON: %w", err)
	}
	
	return &result, nil
}

// CalculateTravelTime calculates total travel time for a route with optional ship parameters
func CalculateTravelTime(db *sql.DB, fromSystemID, toSystemID int64, params *NavigationParams) (*RouteResult, error) {
	// Get effective parameters
	warpSpeed, alignTime, avgWarpDist, source := getEffectiveParams(params)
	
	// Determine if we should avoid low-sec
	avoidLowSec := false
	if params != nil {
		avoidLowSec = params.AvoidLowSec
	}
	
	// Find the shortest path
	path, err := ShortestPath(db, fromSystemID, toSystemID, avoidLowSec)
	if err != nil {
		return nil, err
	}
	
	// Calculate time per jump
	warpTime := CalculateSimplifiedWarpTime(avgWarpDist, warpSpeed)
	timePerJump := alignTime + warpTime + DefaultGateJumpDelay
	
	// Calculate total time
	totalSeconds := float64(path.Jumps) * timePerJump
	
	result := &RouteResult{
		TotalSeconds:      totalSeconds,
		TotalMinutes:      totalSeconds / 60.0,
		Jumps:             path.Jumps,
		AvgSecondsPerJump: timePerJump,
		Route:             path.Route,
		ParametersUsed: map[string]interface{}{
			"warp_speed": warpSpeed,
			"align_time": alignTime,
			"source":     source,
		},
	}
	
	return result, nil
}

// CalculateTravelTimeExact calculates travel time using exact CCP warp formula
func CalculateTravelTimeExact(db *sql.DB, fromSystemID, toSystemID int64, params *NavigationParams) (*RouteResult, error) {
	// Get effective parameters
	warpSpeed, alignTime, avgWarpDist, source := getEffectiveParams(params)
	
	// Determine if we should avoid low-sec
	avoidLowSec := false
	if params != nil {
		avoidLowSec = params.AvoidLowSec
	}
	
	// Find the shortest path
	path, err := ShortestPath(db, fromSystemID, toSystemID, avoidLowSec)
	if err != nil {
		return nil, err
	}
	
	// Calculate time per jump using exact formula
	warpTime := CalculateWarpTime(avgWarpDist, warpSpeed)
	timePerJump := alignTime + warpTime + DefaultGateJumpDelay
	
	// Calculate total time
	totalSeconds := float64(path.Jumps) * timePerJump
	
	result := &RouteResult{
		TotalSeconds:      totalSeconds,
		TotalMinutes:      totalSeconds / 60.0,
		Jumps:             path.Jumps,
		AvgSecondsPerJump: timePerJump,
		Route:             path.Route,
		ParametersUsed: map[string]interface{}{
			"warp_speed": warpSpeed,
			"align_time": alignTime,
			"source":     source,
			"formula":    "exact_3phase",
		},
	}
	
	return result, nil
}
