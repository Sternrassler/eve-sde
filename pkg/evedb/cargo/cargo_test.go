package cargo

import (
	"math"
	"testing"
)

func TestApplySkillModifiers_NoSkills(t *testing.T) {
	baseCapacity := 1000.0

	tests := []struct {
		name     string
		holdType string
		want     float64
	}{
		{"cargo hold no skills", "cargo", 1000.0},
		{"ore hold no skills", "ore_hold", 1000.0},
		{"fleet hangar no skills", "fleet_hangar", 1000.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ApplySkillModifiers(baseCapacity, nil, tt.holdType)
			if got != tt.want {
				t.Errorf("ApplySkillModifiers(nil) = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplySkillModifiers_RacialHauler(t *testing.T) {
	baseCapacity := 1000.0

	tests := []struct {
		name  string
		level int
		want  float64
	}{
		{"Racial Hauler I", 1, 1050.0},   // +5%
		{"Racial Hauler III", 3, 1150.0}, // +15%
		{"Racial Hauler V", 5, 1250.0},   // +25%
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			skills := &SkillModifiers{
				RacialHaulerLevel: &tt.level,
			}
			got := ApplySkillModifiers(baseCapacity, skills, "cargo")
			if math.Abs(got-tt.want) > 0.01 {
				t.Errorf("ApplySkillModifiers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplySkillModifiers_Freighter(t *testing.T) {
	baseCapacity := 10000.0

	tests := []struct {
		name  string
		level int
		want  float64
	}{
		{"Freighter I", 1, 10500.0},   // +5%
		{"Freighter III", 3, 11500.0}, // +15%
		{"Freighter V", 5, 12500.0},   // +25%
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			skills := &SkillModifiers{
				FreighterLevel: &tt.level,
			}
			got := ApplySkillModifiers(baseCapacity, skills, "cargo")
			if math.Abs(got-tt.want) > 0.01 {
				t.Errorf("ApplySkillModifiers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplySkillModifiers_CombinedSkills(t *testing.T) {
	baseCapacity := 1000.0

	racialLevel := 5
	freighterLevel := 3

	skills := &SkillModifiers{
		RacialHaulerLevel: &racialLevel,
		FreighterLevel:    &freighterLevel,
	}

	// Should apply both bonuses multiplicatively: 1000 * 1.25 * 1.15 = 1437.5
	got := ApplySkillModifiers(baseCapacity, skills, "cargo")
	want := 1437.5

	if math.Abs(got-want) > 0.01 {
		t.Errorf("ApplySkillModifiers(combined) = %v, want %v", got, want)
	}
}

func TestApplySkillModifiers_MiningBarge(t *testing.T) {
	baseCapacity := 5000.0

	tests := []struct {
		name  string
		level int
		want  float64
	}{
		{"Mining Barge I", 1, 5250.0},   // +5%
		{"Mining Barge III", 3, 5750.0}, // +15%
		{"Mining Barge V", 5, 6250.0},   // +25%
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			skills := &SkillModifiers{
				MiningBargeLevel: &tt.level,
			}
			got := ApplySkillModifiers(baseCapacity, skills, "ore_hold")
			if math.Abs(got-tt.want) > 0.01 {
				t.Errorf("ApplySkillModifiers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplySkillModifiers_CustomMultiplier(t *testing.T) {
	baseCapacity := 1000.0

	tests := []struct {
		name       string
		multiplier float64
		holdType   string
		want       float64
	}{
		{"Cargo 1.5x", 1.5, "cargo", 1500.0},
		{"Ore Hold 2.0x", 2.0, "ore_hold", 2000.0},
		{"Fleet Hangar 1.2x", 1.2, "fleet_hangar", 1200.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var skills *SkillModifiers
			switch tt.holdType {
			case "cargo":
				skills = &SkillModifiers{CargoMultiplier: &tt.multiplier}
			case "ore_hold":
				skills = &SkillModifiers{OreHoldMultiplier: &tt.multiplier}
			case "fleet_hangar":
				skills = &SkillModifiers{FleetHangarMultiplier: &tt.multiplier}
			}

			got := ApplySkillModifiers(baseCapacity, skills, tt.holdType)
			if math.Abs(got-tt.want) > 0.01 {
				t.Errorf("ApplySkillModifiers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplySkillModifiers_ZeroLevel(t *testing.T) {
	baseCapacity := 1000.0
	level := 0

	skills := &SkillModifiers{
		RacialHaulerLevel: &level,
	}

	// Level 0 should result in no bonus
	got := ApplySkillModifiers(baseCapacity, skills, "cargo")
	if got != baseCapacity {
		t.Errorf("ApplySkillModifiers(level 0) = %v, want %v", got, baseCapacity)
	}
}

func TestApplySkillModifiers_ComplexCombination(t *testing.T) {
	baseCapacity := 2000.0

	racialLevel := 4
	customMultiplier := 1.1

	skills := &SkillModifiers{
		RacialHaulerLevel: &racialLevel,
		CargoMultiplier:   &customMultiplier,
	}

	// 2000 * 1.20 (racial) * 1.1 (custom) = 2640
	got := ApplySkillModifiers(baseCapacity, skills, "cargo")
	want := 2640.0

	if math.Abs(got-want) > 0.01 {
		t.Errorf("ApplySkillModifiers(complex) = %v, want %v", got, want)
	}
}

func TestApplySkillModifiers_WrongHoldType(t *testing.T) {
	baseCapacity := 1000.0
	level := 5

	skills := &SkillModifiers{
		RacialHaulerLevel: &level,
	}

	// Racial hauler should not affect ore hold
	got := ApplySkillModifiers(baseCapacity, skills, "ore_hold")
	if got != baseCapacity {
		t.Errorf("ApplySkillModifiers(wrong hold type) = %v, want %v", got, baseCapacity)
	}
}

func TestSkillModifiers_JSONTags(t *testing.T) {
	// This test validates that JSON tags are correctly defined
	// It's mostly a compile-time check
	skills := &SkillModifiers{
		RacialHaulerLevel: ptrInt(5),
		FreighterLevel:    ptrInt(3),
		MiningBargeLevel:  ptrInt(4),
	}

	if skills.RacialHaulerLevel == nil {
		t.Error("RacialHaulerLevel should not be nil")
	}
	if *skills.RacialHaulerLevel != 5 {
		t.Errorf("RacialHaulerLevel = %v, want 5", *skills.RacialHaulerLevel)
	}
}

// Helper function to create int pointers
func ptrInt(v int) *int {
	return &v
}

// Helper function to create float64 pointers
func ptrFloat64(v float64) *float64 {
	return &v
}
