package scales

import (
	"math"
	"testing"

	"github.com/SCKelemen/color"
)

// Helper to compare colors with tolerance
func colorsApproxEqual(c1, c2 color.Color, tolerance float64) bool {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()

	// Convert to 0-1 range
	rf1, gf1, bf1, af1 := float64(r1)/65535, float64(g1)/65535, float64(b1)/65535, float64(a1)/65535
	rf2, gf2, bf2, af2 := float64(r2)/65535, float64(g2)/65535, float64(b2)/65535, float64(a2)/65535

	return math.Abs(rf1-rf2) <= tolerance &&
		math.Abs(gf1-gf2) <= tolerance &&
		math.Abs(bf1-bf2) <= tolerance &&
		math.Abs(af1-af2) <= tolerance
}

// ===================== SequentialColorScale Tests =====================

func TestNewSequentialColorScale(t *testing.T) {
	start := color.RGB(1, 1, 1) // White
	end := color.RGB(0, 0, 1)   // Blue

	scale := NewSequentialColorScale([2]float64{0, 100}, start, end)

	if scale == nil {
		t.Fatal("NewSequentialColorScale returned nil")
	}

	if scale.domain[0] != 0 || scale.domain[1] != 100 {
		t.Errorf("Domain = %v, expected [0, 100]", scale.domain)
	}

	if !scale.clamp {
		t.Error("Default clamp should be true")
	}

	if scale.space != color.GradientOKLCH {
		t.Error("Default space should be GradientOKLCH")
	}
}

func TestSequentialColorScale_ApplyColor(t *testing.T) {
	white := color.RGB(1, 1, 1)
	blue := color.RGB(0, 0, 1)

	scale := NewSequentialColorScale([2]float64{0, 100}, white, blue)

	tests := []struct {
		name     string
		value    interface{}
		expected color.Color
		name_desc string
	}{
		{"Start", 0, white, "should return white at start"},
		{"End", 100, blue, "should return blue at end"},
		{"Middle", 50, color.RGB(0.5, 0.5, 1), "should return light blue at middle"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scale.ApplyColor(tt.value)
			if result == nil {
				t.Error("ApplyColor returned nil")
			}
		})
	}
}

func TestSequentialColorScale_ApplyColor_IntValue(t *testing.T) {
	white := color.RGB(1, 1, 1)
	black := color.RGB(0, 0, 0)

	scale := NewSequentialColorScale([2]float64{0, 100}, white, black)

	result := scale.ApplyColor(50)
	if result == nil {
		t.Error("ApplyColor with int value returned nil")
	}
}

func TestSequentialColorScale_ApplyColor_InvalidValue(t *testing.T) {
	white := color.RGB(1, 1, 1)
	black := color.RGB(0, 0, 0)

	scale := NewSequentialColorScale([2]float64{0, 100}, white, black)

	result := scale.ApplyColor("invalid")
	if result == nil {
		t.Error("ApplyColor with invalid value should return default color, not nil")
	}
}

func TestSequentialColorScale_Clamp(t *testing.T) {
	white := color.RGB(1, 1, 1)
	black := color.RGB(0, 0, 0)

	scale := NewSequentialColorScale([2]float64{0, 100}, white, black).Clamp(false)

	if scale.clamp {
		t.Error("Clamp(false) did not disable clamping")
	}

	// Test values outside domain
	result := scale.ApplyValue(150.0)
	if result <= 1.0 {
		t.Errorf("With clamp disabled, value > domain max should give t > 1, got %v", result)
	}

	result = scale.ApplyValue(-50.0)
	if result >= 0.0 {
		t.Errorf("With clamp disabled, value < domain min should give t < 0, got %v", result)
	}
}

func TestSequentialColorScale_Space(t *testing.T) {
	white := color.RGB(1, 1, 1)
	black := color.RGB(0, 0, 0)

	scale := NewSequentialColorScale([2]float64{0, 100}, white, black).Space(color.GradientRGB)

	if scale.space != color.GradientRGB {
		t.Errorf("Space = %v, expected GradientRGB", scale.space)
	}
}

func TestSequentialColorScale_Interpolate(t *testing.T) {
	white := color.RGB(1, 1, 1)
	black := color.RGB(0, 0, 0)

	// Custom interpolator that squares the input
	squareInterp := func(t float64) float64 {
		return t * t
	}

	scale := NewSequentialColorScale([2]float64{0, 100}, white, black).Interpolate(squareInterp)

	if scale.interpolate == nil {
		t.Error("Interpolate did not set custom function")
	}

	// Midpoint should be darker with square interpolation
	midValue := scale.ApplyValue(50.0)
	if midValue != 0.25 { // 0.5^2 = 0.25
		t.Errorf("Custom interpolator not applied, got %v, expected 0.25", midValue)
	}
}

func TestSequentialColorScale_Samples(t *testing.T) {
	white := color.RGB(1, 1, 1)
	black := color.RGB(0, 0, 0)

	scale := NewSequentialColorScale([2]float64{0, 100}, white, black)

	samples := scale.Samples(5)
	if len(samples) != 5 {
		t.Errorf("Samples(5) returned %d colors, expected 5", len(samples))
	}

	// First should be white, last should be black
	if !colorsApproxEqual(samples[0], white, 0.01) {
		t.Error("First sample should be white")
	}

	if !colorsApproxEqual(samples[4], black, 0.01) {
		t.Error("Last sample should be black")
	}
}

func TestSequentialColorScale_Samples_EdgeCases(t *testing.T) {
	white := color.RGB(1, 1, 1)
	black := color.RGB(0, 0, 0)

	scale := NewSequentialColorScale([2]float64{0, 100}, white, black)

	// Zero samples
	samples := scale.Samples(0)
	if samples != nil {
		t.Error("Samples(0) should return nil")
	}

	// One sample
	samples = scale.Samples(1)
	if len(samples) != 1 {
		t.Errorf("Samples(1) returned %d colors", len(samples))
	}
}

func TestSequentialColorScale_Clone(t *testing.T) {
	white := color.RGB(1, 1, 1)
	black := color.RGB(0, 0, 0)

	scale := NewSequentialColorScale([2]float64{0, 100}, white, black).Clamp(false)
	clone := scale.Clone()

	if clone == nil {
		t.Fatal("Clone returned nil")
	}

	seqClone, ok := clone.(*SequentialColorScale)
	if !ok {
		t.Fatal("Clone did not return *SequentialColorScale")
	}

	if seqClone.clamp != scale.clamp {
		t.Error("Clone did not preserve clamp setting")
	}

	if seqClone.domain != scale.domain {
		t.Error("Clone did not preserve domain")
	}
}

func TestSequentialColorScale_ScaleInterface(t *testing.T) {
	white := color.RGB(1, 1, 1)
	black := color.RGB(0, 0, 0)

	var s Scale = NewSequentialColorScale([2]float64{0, 100}, white, black)

	// Test Scale interface methods
	if s.Type() != ScaleTypeSequential {
		t.Errorf("Type() = %v, expected ScaleTypeSequential", s.Type())
	}

	domain := s.Domain()
	if domain == nil {
		t.Error("Domain() returned nil")
	}

	range_ := s.Range()
	if range_[0].Value != 0 || range_[1].Value != 1 {
		t.Error("Range() should return dummy [0px, 1px]")
	}

	// Apply returns dummy value
	result := s.Apply(50.0)
	if result.Value != 0 {
		t.Error("Apply() should return dummy 0px")
	}
}

// ===================== DivergingColorScale Tests =====================

func TestNewDivergingColorScale(t *testing.T) {
	blue := color.RGB(0, 0, 1)
	white := color.RGB(1, 1, 1)
	red := color.RGB(1, 0, 0)

	scale := NewDivergingColorScale([2]float64{-100, 100}, blue, white, red)

	if scale == nil {
		t.Fatal("NewDivergingColorScale returned nil")
	}

	if scale.domain[0] != -100 || scale.domain[1] != 100 {
		t.Errorf("Domain = %v, expected [-100, 100]", scale.domain)
	}

	if scale.midpoint != 0 {
		t.Errorf("Midpoint = %v, expected 0", scale.midpoint)
	}

	if !scale.clamp {
		t.Error("Default clamp should be true")
	}
}

func TestDivergingColorScale_ApplyColor(t *testing.T) {
	blue := color.RGB(0, 0, 1)
	white := color.RGB(1, 1, 1)
	red := color.RGB(1, 0, 0)

	scale := NewDivergingColorScale([2]float64{-100, 100}, blue, white, red)

	tests := []struct {
		name  string
		value interface{}
		desc  string
	}{
		{"NegativeEnd", -100, "should return blue"},
		{"PositiveEnd", 100, "should return red"},
		{"Midpoint", 0, "should return white"},
		{"NegativeMid", -50, "should return light blue"},
		{"PositiveMid", 50, "should return light red"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scale.ApplyColor(tt.value)
			if result == nil {
				t.Errorf("%s: ApplyColor returned nil", tt.desc)
			}
		})
	}
}

func TestDivergingColorScale_ApplyColor_IntValue(t *testing.T) {
	blue := color.RGB(0, 0, 1)
	white := color.RGB(1, 1, 1)
	red := color.RGB(1, 0, 0)

	scale := NewDivergingColorScale([2]float64{-100, 100}, blue, white, red)

	result := scale.ApplyColor(0)
	if result == nil {
		t.Error("ApplyColor with int value returned nil")
	}
}

func TestDivergingColorScale_ApplyColor_InvalidValue(t *testing.T) {
	blue := color.RGB(0, 0, 1)
	white := color.RGB(1, 1, 1)
	red := color.RGB(1, 0, 0)

	scale := NewDivergingColorScale([2]float64{-100, 100}, blue, white, red)

	result := scale.ApplyColor("invalid")
	if !colorsApproxEqual(result, white, 0.01) {
		t.Error("ApplyColor with invalid value should return midColor")
	}
}

func TestDivergingColorScale_Midpoint(t *testing.T) {
	blue := color.RGB(0, 0, 1)
	white := color.RGB(1, 1, 1)
	red := color.RGB(1, 0, 0)

	scale := NewDivergingColorScale([2]float64{-100, 100}, blue, white, red).Midpoint(20.0)

	if scale.midpoint != 20.0 {
		t.Errorf("Midpoint = %v, expected 20", scale.midpoint)
	}

	// Color at midpoint should be white
	result := scale.ApplyColor(20.0)
	if !colorsApproxEqual(result, white, 0.01) {
		t.Error("Color at custom midpoint should be white")
	}
}

func TestDivergingColorScale_ApplyValue(t *testing.T) {
	blue := color.RGB(0, 0, 1)
	white := color.RGB(1, 1, 1)
	red := color.RGB(1, 0, 0)

	scale := NewDivergingColorScale([2]float64{-100, 100}, blue, white, red)

	tests := []struct {
		value    float64
		expected float64
	}{
		{-100, 0.0},
		{0, 0.5},
		{100, 1.0},
		{-50, 0.25},
		{50, 0.75},
	}

	for _, tt := range tests {
		result := scale.ApplyValue(tt.value)
		if math.Abs(result-tt.expected) > 0.001 {
			t.Errorf("ApplyValue(%v) = %v, expected %v", tt.value, result, tt.expected)
		}
	}
}

func TestDivergingColorScale_Samples(t *testing.T) {
	blue := color.RGB(0, 0, 1)
	white := color.RGB(1, 1, 1)
	red := color.RGB(1, 0, 0)

	scale := NewDivergingColorScale([2]float64{-100, 100}, blue, white, red)

	samples := scale.Samples(5)
	if len(samples) != 5 {
		t.Errorf("Samples(5) returned %d colors", len(samples))
	}

	// First should be blue, middle should be white, last should be red
	if !colorsApproxEqual(samples[0], blue, 0.01) {
		t.Error("First sample should be blue")
	}

	if !colorsApproxEqual(samples[2], white, 0.1) {
		t.Error("Middle sample should be approximately white")
	}

	if !colorsApproxEqual(samples[4], red, 0.01) {
		t.Error("Last sample should be red")
	}
}

func TestDivergingColorScale_Clone(t *testing.T) {
	blue := color.RGB(0, 0, 1)
	white := color.RGB(1, 1, 1)
	red := color.RGB(1, 0, 0)

	scale := NewDivergingColorScale([2]float64{-100, 100}, blue, white, red).Midpoint(20.0)
	clone := scale.Clone()

	divClone, ok := clone.(*DivergingColorScale)
	if !ok {
		t.Fatal("Clone did not return *DivergingColorScale")
	}

	if divClone.midpoint != scale.midpoint {
		t.Error("Clone did not preserve midpoint")
	}
}

func TestDivergingColorScale_ScaleInterface(t *testing.T) {
	blue := color.RGB(0, 0, 1)
	white := color.RGB(1, 1, 1)
	red := color.RGB(1, 0, 0)

	var s Scale = NewDivergingColorScale([2]float64{-100, 100}, blue, white, red)

	if s.Type() != ScaleTypeDiverging {
		t.Errorf("Type() = %v, expected ScaleTypeDiverging", s.Type())
	}
}

// ===================== CategoricalColorScale Tests =====================

func TestNewCategoricalColorScale(t *testing.T) {
	domain := []string{"A", "B", "C"}
	colors := []color.Color{
		color.RGB(1, 0, 0), // Red
		color.RGB(0, 1, 0), // Green
		color.RGB(0, 0, 1), // Blue
	}

	scale := NewCategoricalColorScale(domain, colors)

	if scale == nil {
		t.Fatal("NewCategoricalColorScale returned nil")
	}

	if len(scale.domain) != 3 {
		t.Errorf("Domain length = %d, expected 3", len(scale.domain))
	}

	if len(scale.colors) != 3 {
		t.Errorf("Colors length = %d, expected 3", len(scale.colors))
	}
}

func TestCategoricalColorScale_ApplyColor(t *testing.T) {
	domain := []string{"A", "B", "C"}
	colors := []color.Color{
		color.RGB(1, 0, 0), // Red
		color.RGB(0, 1, 0), // Green
		color.RGB(0, 0, 1), // Blue
	}

	scale := NewCategoricalColorScale(domain, colors)

	tests := []struct {
		value    string
		expected color.Color
	}{
		{"A", colors[0]},
		{"B", colors[1]},
		{"C", colors[2]},
	}

	for _, tt := range tests {
		result := scale.ApplyColor(tt.value)
		if !colorsApproxEqual(result, tt.expected, 0.01) {
			t.Errorf("ApplyColor(%q) did not return expected color", tt.value)
		}
	}
}

func TestCategoricalColorScale_ApplyColor_Cycling(t *testing.T) {
	domain := []string{"A", "B", "C", "D", "E"}
	colors := []color.Color{
		color.RGB(1, 0, 0), // Red
		color.RGB(0, 1, 0), // Green
		color.RGB(0, 0, 1), // Blue
	}

	scale := NewCategoricalColorScale(domain, colors)

	// D should cycle to Red (index 3 % 3 = 0)
	resultD := scale.ApplyColor("D")
	if !colorsApproxEqual(resultD, colors[0], 0.01) {
		t.Error("Category 'D' should cycle to Red")
	}

	// E should cycle to Green (index 4 % 3 = 1)
	resultE := scale.ApplyColor("E")
	if !colorsApproxEqual(resultE, colors[1], 0.01) {
		t.Error("Category 'E' should cycle to Green")
	}
}

func TestCategoricalColorScale_ApplyColor_Unknown(t *testing.T) {
	domain := []string{"A", "B", "C"}
	colors := []color.Color{
		color.RGB(1, 0, 0),
		color.RGB(0, 1, 0),
		color.RGB(0, 0, 1),
	}

	scale := NewCategoricalColorScale(domain, colors)

	result := scale.ApplyColor("Unknown")
	// Should return default gray
	if result == nil {
		t.Error("ApplyColor with unknown category returned nil")
	}
}

func TestCategoricalColorScale_ApplyColor_InvalidType(t *testing.T) {
	domain := []string{"A", "B", "C"}
	colors := []color.Color{
		color.RGB(1, 0, 0),
		color.RGB(0, 1, 0),
		color.RGB(0, 0, 1),
	}

	scale := NewCategoricalColorScale(domain, colors)

	result := scale.ApplyColor(123)
	// Should return unknown color
	if result == nil {
		t.Error("ApplyColor with invalid type returned nil")
	}
}

func TestCategoricalColorScale_Unknown(t *testing.T) {
	domain := []string{"A", "B", "C"}
	colors := []color.Color{
		color.RGB(1, 0, 0),
		color.RGB(0, 1, 0),
		color.RGB(0, 0, 1),
	}

	customUnknown := color.RGB(1, 1, 0) // Yellow
	scale := NewCategoricalColorScale(domain, colors).Unknown(customUnknown)

	result := scale.ApplyColor("Unknown")
	if !colorsApproxEqual(result, customUnknown, 0.01) {
		t.Error("Unknown() did not set custom unknown color")
	}
}

func TestCategoricalColorScale_ApplyValue(t *testing.T) {
	domain := []string{"A", "B", "C"}
	colors := []color.Color{
		color.RGB(1, 0, 0),
		color.RGB(0, 1, 0),
		color.RGB(0, 0, 1),
	}

	scale := NewCategoricalColorScale(domain, colors)

	tests := []struct {
		value    string
		expected float64
	}{
		{"A", 0.0},
		{"B", 0.5},
		{"C", 1.0},
	}

	for _, tt := range tests {
		result := scale.ApplyValue(tt.value)
		if math.Abs(result-tt.expected) > 0.001 {
			t.Errorf("ApplyValue(%q) = %v, expected %v", tt.value, result, tt.expected)
		}
	}
}

func TestCategoricalColorScale_ApplyValue_SingleCategory(t *testing.T) {
	domain := []string{"A"}
	colors := []color.Color{color.RGB(1, 0, 0)}

	scale := NewCategoricalColorScale(domain, colors)

	result := scale.ApplyValue("A")
	if result != 0.5 {
		t.Errorf("ApplyValue with single category = %v, expected 0.5", result)
	}
}

func TestCategoricalColorScale_Colors(t *testing.T) {
	domain := []string{"A", "B", "C"}
	colors := []color.Color{
		color.RGB(1, 0, 0),
		color.RGB(0, 1, 0),
		color.RGB(0, 0, 1),
	}

	scale := NewCategoricalColorScale(domain, colors)

	result := scale.Colors()
	if len(result) != len(colors) {
		t.Errorf("Colors() returned %d colors, expected %d", len(result), len(colors))
	}
}

func TestCategoricalColorScale_Clone(t *testing.T) {
	domain := []string{"A", "B", "C"}
	colors := []color.Color{
		color.RGB(1, 0, 0),
		color.RGB(0, 1, 0),
		color.RGB(0, 0, 1),
	}

	scale := NewCategoricalColorScale(domain, colors)
	clone := scale.Clone()

	catClone, ok := clone.(*CategoricalColorScale)
	if !ok {
		t.Fatal("Clone did not return *CategoricalColorScale")
	}

	// Modify original should not affect clone
	scale.domain[0] = "Modified"
	if catClone.domain[0] == "Modified" {
		t.Error("Clone did not create deep copy of domain")
	}
}

func TestCategoricalColorScale_ScaleInterface(t *testing.T) {
	domain := []string{"A", "B", "C"}
	colors := []color.Color{
		color.RGB(1, 0, 0),
		color.RGB(0, 1, 0),
		color.RGB(0, 0, 1),
	}

	var s Scale = NewCategoricalColorScale(domain, colors)

	if s.Type() != ScaleTypeOrdinal {
		t.Errorf("Type() = %v, expected ScaleTypeOrdinal", s.Type())
	}
}

func TestCategoricalColorScale_EmptyColors(t *testing.T) {
	domain := []string{"A", "B", "C"}
	colors := []color.Color{} // Empty

	scale := NewCategoricalColorScale(domain, colors)

	// Should return unknown color
	result := scale.ApplyColor("A")
	if result == nil {
		t.Error("ApplyColor with empty colors returned nil")
	}
}

// ===================== Benchmark Tests =====================

func BenchmarkSequentialColorScale_ApplyColor(b *testing.B) {
	white := color.RGB(1, 1, 1)
	black := color.RGB(0, 0, 0)
	scale := NewSequentialColorScale([2]float64{0, 100}, white, black)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scale.ApplyColor(float64(i % 100))
	}
}

func BenchmarkDivergingColorScale_ApplyColor(b *testing.B) {
	blue := color.RGB(0, 0, 1)
	white := color.RGB(1, 1, 1)
	red := color.RGB(1, 0, 0)
	scale := NewDivergingColorScale([2]float64{-100, 100}, blue, white, red)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scale.ApplyColor(float64((i%200) - 100))
	}
}

func BenchmarkCategoricalColorScale_ApplyColor(b *testing.B) {
	domain := []string{"A", "B", "C", "D", "E"}
	colors := []color.Color{
		color.RGB(1, 0, 0),
		color.RGB(0, 1, 0),
		color.RGB(0, 0, 1),
		color.RGB(1, 1, 0),
		color.RGB(1, 0, 1),
	}
	scale := NewCategoricalColorScale(domain, colors)

	categories := []string{"A", "B", "C", "D", "E"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scale.ApplyColor(categories[i%len(categories)])
	}
}
