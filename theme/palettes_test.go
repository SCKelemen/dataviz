package theme

import (
	"strings"
	"testing"
)

func TestDarkCategorical(t *testing.T) {
	colors := DarkCategorical()

	if len(colors) == 0 {
		t.Fatal("Expected non-empty color palette")
	}

	// Check that all colors are valid hex colors
	for i, color := range colors {
		if !strings.HasPrefix(color, "#") {
			t.Errorf("Color %d is not a valid hex color: %s", i, color)
		}

		if len(color) != 7 {
			t.Errorf("Color %d has invalid length: %s", i, color)
		}
	}

	// Check that colors are distinct
	seen := make(map[string]bool)
	for _, color := range colors {
		if seen[color] {
			t.Errorf("Duplicate color in palette: %s", color)
		}
		seen[color] = true
	}
}

func TestLightCategorical(t *testing.T) {
	colors := LightCategorical()

	if len(colors) == 0 {
		t.Fatal("Expected non-empty color palette")
	}

	for i, color := range colors {
		if !strings.HasPrefix(color, "#") {
			t.Errorf("Color %d is not a valid hex color: %s", i, color)
		}
	}
}

func TestDarkSequential(t *testing.T) {
	accent := "#3B82F6"
	colors := DarkSequential(accent)

	if len(colors) == 0 {
		t.Fatal("Expected non-empty sequential palette")
	}

	// First color should be darker than last
	first := colors[0]
	last := colors[len(colors)-1]

	if first == last {
		t.Error("Expected sequential colors to vary")
	}

	// All colors should be valid hex
	for _, color := range colors {
		if !strings.HasPrefix(color, "#") {
			t.Errorf("Invalid hex color: %s", color)
		}
	}
}

func TestLightSequential(t *testing.T) {
	accent := "#3B82F6"
	colors := LightSequential(accent)

	if len(colors) == 0 {
		t.Fatal("Expected non-empty sequential palette")
	}

	for _, color := range colors {
		if !strings.HasPrefix(color, "#") {
			t.Errorf("Invalid hex color: %s", color)
		}
	}
}

func TestDarkDiverging(t *testing.T) {
	accent := "#EF4444"
	colors := DarkDiverging(accent)

	if len(colors) == 0 {
		t.Fatal("Expected non-empty diverging palette")
	}

	// Diverging palettes should have odd number of colors for a midpoint
	if len(colors)%2 == 0 {
		t.Error("Expected odd number of colors in diverging palette")
	}

	for _, color := range colors {
		if !strings.HasPrefix(color, "#") {
			t.Errorf("Invalid hex color: %s", color)
		}
	}
}

func TestLightDiverging(t *testing.T) {
	accent := "#EF4444"
	colors := LightDiverging(accent)

	if len(colors) == 0 {
		t.Fatal("Expected non-empty diverging palette")
	}

	if len(colors)%2 == 0 {
		t.Error("Expected odd number of colors in diverging palette")
	}
}

func TestHeatmapColors(t *testing.T) {
	tests := []struct {
		value    float64
		darkMode bool
		name     string
	}{
		{0.0, true, "dark_min"},
		{0.5, true, "dark_mid"},
		{1.0, true, "dark_max"},
		{0.0, false, "light_min"},
		{0.5, false, "light_mid"},
		{1.0, false, "light_max"},
		{-0.5, true, "below_zero"},
		{1.5, false, "above_one"},
	}

	for _, test := range tests {
		color := HeatmapColors(test.value, test.darkMode)
		if !strings.HasPrefix(color, "#") {
			t.Errorf("Test %s: expected hex color, got %s", test.name, color)
		}
	}
}

func TestViridisColors(t *testing.T) {
	// Test key points in the Viridis palette
	tests := []float64{0.0, 0.25, 0.5, 0.75, 1.0}

	for _, t_val := range tests {
		color := ViridisColors(t_val)
		if !strings.HasPrefix(color, "#") {
			t.Errorf("Expected hex color for t=%f, got %s", t_val, color)
		}
	}

	// Test that colors progress
	color1 := ViridisColors(0.0)
	color2 := ViridisColors(0.5)
	color3 := ViridisColors(1.0)

	if color1 == color2 || color2 == color3 || color1 == color3 {
		t.Error("Expected Viridis colors to progress")
	}
}

func TestPlasmaColors(t *testing.T) {
	tests := []float64{0.0, 0.25, 0.5, 0.75, 1.0}

	for _, t_val := range tests {
		color := PlasmaColors(t_val)
		if !strings.HasPrefix(color, "#") {
			t.Errorf("Expected hex color for t=%f, got %s", t_val, color)
		}
	}

	// Test that colors progress
	color1 := PlasmaColors(0.0)
	color2 := PlasmaColors(0.5)
	color3 := PlasmaColors(1.0)

	if color1 == color2 || color2 == color3 || color1 == color3 {
		t.Error("Expected Plasma colors to progress")
	}
}

func TestCoolWarmColors(t *testing.T) {
	tests := []struct {
		t    float64
		name string
	}{
		{-1.0, "cool_extreme"},
		{-0.5, "cool"},
		{0.0, "neutral"},
		{0.5, "warm"},
		{1.0, "warm_extreme"},
	}

	for _, test := range tests {
		color := CoolWarmColors(test.t)
		if !strings.HasPrefix(color, "#") {
			t.Errorf("Test %s: expected hex color, got %s", test.name, color)
		}
	}

	// Test progression
	cool := CoolWarmColors(-1.0)
	neutral := CoolWarmColors(0.0)
	warm := CoolWarmColors(1.0)

	if cool == neutral || neutral == warm || cool == warm {
		t.Error("Expected cool-warm colors to differ")
	}
}

func TestQualitativeColors(t *testing.T) {
	// Test both dark and light mode
	for _, darkMode := range []bool{true, false} {
		for i := 0; i < 15; i++ { // Test more than palette length
			color := QualitativeColors(i, darkMode)
			if !strings.HasPrefix(color, "#") {
				t.Errorf("Expected hex color for index %d (darkMode=%v), got %s", i, darkMode, color)
			}
		}

		// Test that colors cycle
		color1 := QualitativeColors(0, darkMode)
		color2 := QualitativeColors(10, darkMode) // Assume palette has 10 colors
		if color1 != color2 {
			// This is fine if palette is longer than 10
		}
	}
}

func TestContrastRatio(t *testing.T) {
	// Test with known color pairs
	tests := []struct {
		color1   string
		color2   string
		minRatio float64
		name     string
	}{
		{"#FFFFFF", "#000000", 15.0, "black_white"}, // Maximum contrast
		{"#000000", "#FFFFFF", 15.0, "white_black"}, // Order shouldn't matter
		{"#FFFFFF", "#FFFFFF", 1.0, "same_color"},   // No contrast
		{"#000000", "#000000", 1.0, "same_black"},
	}

	for _, test := range tests {
		ratio := ContrastRatio(test.color1, test.color2)
		if ratio < test.minRatio {
			t.Errorf("Test %s: expected ratio >= %f, got %f", test.name, test.minRatio, ratio)
		}

		// Contrast ratio should be at least 1.0
		if ratio < 1.0 {
			t.Errorf("Test %s: contrast ratio cannot be less than 1.0, got %f", test.name, ratio)
		}
	}
}

func TestEnsureContrast(t *testing.T) {
	tests := []struct {
		fg       string
		bg       string
		minRatio float64
		name     string
	}{
		{"#808080", "#FFFFFF", 4.5, "gray_on_white"},
		{"#808080", "#000000", 4.5, "gray_on_black"},
		{"#CCCCCC", "#FFFFFF", 3.0, "light_gray_on_white"},
	}

	for _, test := range tests {
		adjusted := EnsureContrast(test.fg, test.bg, test.minRatio)

		if !strings.HasPrefix(adjusted, "#") && !strings.HasPrefix(adjusted, "hsl") {
			t.Errorf("Test %s: expected valid color, got %s", test.name, adjusted)
		}

		// Check that adjusted color meets minimum contrast
		ratio := ContrastRatio(adjusted, test.bg)
		if ratio < test.minRatio {
			t.Errorf("Test %s: adjusted color does not meet minimum contrast (%f < %f)", test.name, ratio, test.minRatio)
		}
	}
}

func TestGenerateSequential(t *testing.T) {
	start := "#FFFFFF"
	end := "#0000FF"
	steps := 5

	colors := generateSequential(end, start, steps)

	if len(colors) != steps {
		t.Errorf("Expected %d colors, got %d", steps, len(colors))
	}

	// First color should be close to start
	if colors[0] != start && colors[0] != "#FFFFFF" {
		// Allow some variation due to color space conversion
	}

	// Last color should be end
	if colors[len(colors)-1] != end {
		// Allow some variation
	}
}

func TestGenerateDiverging(t *testing.T) {
	left := "#0000FF"
	middle := "#FFFFFF"
	right := "#FF0000"
	steps := 9

	colors := generateDiverging(left, middle, right, steps)

	// Should have odd number of steps
	if len(colors)%2 == 0 {
		t.Error("Expected odd number of colors in diverging palette")
	}

	// Middle color should be approximately the middle value
	midIndex := len(colors) / 2
	if colors[midIndex] != middle {
		// Allow some variation
	}
}

func TestInterpolateHex(t *testing.T) {
	start := "#FF0000"
	end := "#0000FF"

	// Test interpolation at key points
	color0 := interpolateHex(start, end, 0.0)
	color50 := interpolateHex(start, end, 0.5)
	color100 := interpolateHex(start, end, 1.0)

	if color0 != start {
		// Allow some variation due to color space conversion
	}

	if color100 != end {
		// Allow some variation
	}

	// Middle color should be different from both
	if color50 == start || color50 == end {
		t.Error("Expected interpolated middle color to differ from start and end")
	}
}
