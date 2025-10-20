package navigation

import (
	"math"
	"testing"
)

func TestCalculateAlignTime(t *testing.T) {
	tests := []struct {
		name            string
		mass            float64
		inertiaModifier float64
		want            float64
	}{
		{
			name:            "Interceptor",
			mass:            1200000,
			inertiaModifier: 0.3,
			want:            1.0,
		},
		{
			name:            "Cruiser",
			mass:            12000000,
			inertiaModifier: 0.4,
			want:            13.3,
		},
		{
			name:            "Battleship",
			mass:            100000000,
			inertiaModifier: 0.15,
			want:            41.6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateAlignTime(tt.mass, tt.inertiaModifier)
			// Allow 10% tolerance for approximation
			tolerance := tt.want * 0.1
			if math.Abs(got-tt.want) > tolerance {
				t.Errorf("CalculateAlignTime() = %v, want approximately %v", got, tt.want)
			}
		})
	}
}

func TestCalculateWarpTime(t *testing.T) {
	tests := []struct {
		name        string
		distanceAU  float64
		warpSpeedAU float64
		wantRange   [2]float64 // min and max expected values
	}{
		{
			name:        "Cruiser 50 AU",
			distanceAU:  50.0,
			warpSpeedAU: 3.0,
			wantRange:   [2]float64{40.0, 55.0}, // Approximate range
		},
		{
			name:        "Interceptor 50 AU",
			distanceAU:  50.0,
			warpSpeedAU: 6.0,
			wantRange:   [2]float64{20.0, 30.0},
		},
		{
			name:        "Battleship 50 AU",
			distanceAU:  50.0,
			warpSpeedAU: 1.5,
			wantRange:   [2]float64{85.0, 100.0},
		},
		{
			name:        "Short warp 5 AU",
			distanceAU:  5.0,
			warpSpeedAU: 3.0,
			wantRange:   [2]float64{20.0, 35.0}, // Short warps still need accel/decel phases
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateWarpTime(tt.distanceAU, tt.warpSpeedAU)
			if got < tt.wantRange[0] || got > tt.wantRange[1] {
				t.Errorf("CalculateWarpTime() = %v, want in range [%v, %v]", got, tt.wantRange[0], tt.wantRange[1])
			}
		})
	}
}

func TestCalculateSimplifiedWarpTime(t *testing.T) {
	tests := []struct {
		name        string
		distanceAU  float64
		warpSpeedAU float64
		want        float64
	}{
		{
			name:        "Default cruiser 15 AU",
			distanceAU:  15.0,
			warpSpeedAU: 3.0,
			want:        7.0, // (15 / 3) * 1.4
		},
		{
			name:        "Interceptor 15 AU",
			distanceAU:  15.0,
			warpSpeedAU: 6.0,
			want:        3.5, // (15 / 6) * 1.4
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateSimplifiedWarpTime(tt.distanceAU, tt.warpSpeedAU)
			if math.Abs(got-tt.want) > 0.1 {
				t.Errorf("CalculateSimplifiedWarpTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEffectiveParams(t *testing.T) {
	tests := []struct {
		name           string
		params         *NavigationParams
		wantWarpSpeed  float64
		wantAlignTime  float64
		wantAvgWarpDist float64
		wantSource     string
	}{
		{
			name:           "Nil params uses defaults",
			params:         nil,
			wantWarpSpeed:  DefaultWarpSpeed,
			wantAlignTime:  DefaultAlignTime,
			wantAvgWarpDist: DefaultAvgWarpDistance,
			wantSource:     "default",
		},
		{
			name: "Custom warp speed",
			params: &NavigationParams{
				WarpSpeed: ptrFloat64(6.0),
			},
			wantWarpSpeed:  6.0,
			wantAlignTime:  DefaultAlignTime,
			wantAvgWarpDist: DefaultAvgWarpDistance,
			wantSource:     "provided",
		},
		{
			name: "Calculated align time from ship params",
			params: &NavigationParams{
				ShipMass:        ptrFloat64(12000000),
				InertiaModifier: ptrFloat64(0.4),
			},
			wantWarpSpeed:  DefaultWarpSpeed,
			wantAlignTime:  13.3, // Approximate
			wantAvgWarpDist: DefaultAvgWarpDistance,
			wantSource:     "calculated",
		},
		{
			name: "All custom params",
			params: &NavigationParams{
				WarpSpeed:       ptrFloat64(8.0),
				AlignTime:       ptrFloat64(2.5),
				AvgWarpDistance: ptrFloat64(20.0),
			},
			wantWarpSpeed:  8.0,
			wantAlignTime:  2.5,
			wantAvgWarpDist: 20.0,
			wantSource:     "provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warpSpeed, alignTime, avgWarpDist, source := getEffectiveParams(tt.params)
			
			if math.Abs(warpSpeed-tt.wantWarpSpeed) > 0.1 {
				t.Errorf("warpSpeed = %v, want %v", warpSpeed, tt.wantWarpSpeed)
			}
			
			// Allow 10% tolerance for calculated align time
			tolerance := tt.wantAlignTime * 0.1
			if math.Abs(alignTime-tt.wantAlignTime) > tolerance {
				t.Errorf("alignTime = %v, want approximately %v", alignTime, tt.wantAlignTime)
			}
			
			if math.Abs(avgWarpDist-tt.wantAvgWarpDist) > 0.1 {
				t.Errorf("avgWarpDist = %v, want %v", avgWarpDist, tt.wantAvgWarpDist)
			}
			
			if source != tt.wantSource {
				t.Errorf("source = %v, want %v", source, tt.wantSource)
			}
		})
	}
}

// Helper function to create float64 pointers
func ptrFloat64(v float64) *float64 {
	return &v
}
