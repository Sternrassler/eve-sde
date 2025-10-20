// Package navigation provides EVE Online navigation and route planning functionality
// It includes pathfinding, travel time calculation, and trade hub analysis
//
// This is a high-level API layer built on top of the eve-sde SQLite database.
// It requires navigation views to be initialized (see internal/sqlite/views package).
package navigation

import (
	"database/sql"
	"fmt"
	"math"
)

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
	TotalSeconds      float64                `json:"total_seconds"`
	TotalMinutes      float64                `json:"total_minutes"`
	Jumps             int                    `json:"jumps"`
	AvgSecondsPerJump float64                `json:"avg_seconds_per_jump"`
	Route             []int64                `json:"route"`
	ParametersUsed    map[string]interface{} `json:"parameters_used"`
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

// Note: InitializeViews has been moved to internal/sqlite/views package
// Views must be initialized before using this API

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

// edge represents a stargate connection
type edge struct {
	toSystemID int64
}

// ShortestPath finds the shortest path between two systems using Dijkstra's algorithm
// This implementation loads the graph into memory and uses a Go-based algorithm
// for better performance on long-distance routes (40+ jumps)
func ShortestPath(db *sql.DB, fromSystemID, toSystemID int64, avoidLowSec bool) (*PathResult, error) {
	// Load the graph from database
	graph, err := loadGraph(db, avoidLowSec)
	if err != nil {
		return nil, fmt.Errorf("failed to load graph: %w", err)
	}

	// Run Dijkstra's algorithm
	path, found := dijkstra(graph, fromSystemID, toSystemID)
	if !found {
		return nil, fmt.Errorf("no path found between systems %d and %d", fromSystemID, toSystemID)
	}

	result := &PathResult{
		FromSystemID: fromSystemID,
		ToSystemID:   toSystemID,
		Jumps:        len(path) - 1, // jumps = number of systems - 1
		Route:        path,
	}

	return result, nil
}

// loadGraph loads the stargate graph from the database
func loadGraph(db *sql.DB, avoidLowSec bool) (map[int64][]edge, error) {
	var query string
	if avoidLowSec {
		query = `
			SELECT DISTINCT g.from_system_id, g.to_system_id
			FROM v_stargate_graph g
			LEFT JOIN mapSolarSystems sys ON g.to_system_id = sys._key
			WHERE sys.securityStatus >= 0.45 OR sys.securityStatus IS NULL
		`
	} else {
		query = `
			SELECT from_system_id, to_system_id
			FROM v_stargate_graph
		`
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query graph: %w", err)
	}
	defer rows.Close()

	graph := make(map[int64][]edge)
	for rows.Next() {
		var from, to int64
		if err := rows.Scan(&from, &to); err != nil {
			return nil, fmt.Errorf("failed to scan edge: %w", err)
		}
		graph[from] = append(graph[from], edge{toSystemID: to})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating graph: %w", err)
	}

	return graph, nil
}

// dijkstra implements Dijkstra's shortest path algorithm
func dijkstra(graph map[int64][]edge, start, goal int64) ([]int64, bool) {
	// Check if start and goal exist in graph
	if _, exists := graph[start]; !exists {
		return nil, false
	}

	// Distance map: system -> distance from start
	dist := make(map[int64]int)
	dist[start] = 0

	// Previous node map for path reconstruction
	prev := make(map[int64]int64)

	// Visited set
	visited := make(map[int64]bool)

	// Priority queue (using simple slice for now, can optimize with heap)
	// Each element is [systemID, distance]
	pq := []struct {
		systemID int64
		distance int
	}{{start, 0}}

	for len(pq) > 0 {
		// Find minimum distance node (linear search for simplicity)
		minIdx := 0
		for i := 1; i < len(pq); i++ {
			if pq[i].distance < pq[minIdx].distance {
				minIdx = i
			}
		}

		// Extract minimum
		current := pq[minIdx]
		pq = append(pq[:minIdx], pq[minIdx+1:]...)

		// Skip if already visited
		if visited[current.systemID] {
			continue
		}

		// Mark as visited
		visited[current.systemID] = true

		// Early termination: if we reached the goal, reconstruct path
		if current.systemID == goal {
			return reconstructPath(prev, start, goal), true
		}

		// Explore neighbors
		for _, e := range graph[current.systemID] {
			if visited[e.toSystemID] {
				continue
			}

			newDist := current.distance + 1
			if oldDist, exists := dist[e.toSystemID]; !exists || newDist < oldDist {
				dist[e.toSystemID] = newDist
				prev[e.toSystemID] = current.systemID
				pq = append(pq, struct {
					systemID int64
					distance int
				}{e.toSystemID, newDist})
			}
		}
	}

	// No path found
	return nil, false
}

// reconstructPath builds the path from start to goal using the prev map
func reconstructPath(prev map[int64]int64, start, goal int64) []int64 {
	path := []int64{goal}
	current := goal

	for current != start {
		current = prev[current]
		path = append([]int64{current}, path...)
	}

	return path
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
