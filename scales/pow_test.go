package scales

import (
	"math"
	"testing"

	"github.com/SCKelemen/units"
)

func TestPowScale_Linear(t *testing.T) {
	// Exponent 1 should behave like linear scale
	scale := NewPowScale(
		[2]float64{0, 100},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	tests := []struct {
		input    float64
		expected float64
	}{
		{0, 0},
		{25, 125},
		{50, 250},
		{75, 375},
		{100, 500},
	}

	for _, tt := range tests {
		result := scale.Apply(tt.input)
		if math.Abs(result.Value-tt.expected) > 0.01 {
			t.Errorf("Apply(%v) with exponent=1 = %v, expected %v", tt.input, result.Value, tt.expected)
		}
		if result.Unit != units.PX {
			t.Errorf("Apply(%v) unit = %v, expected PX", tt.input, result.Unit)
		}
	}
}

func TestPowScale_Square(t *testing.T) {
	scale := NewPowScale(
		[2]float64{0, 100},
		[2]units.Length{units.Px(0), units.Px(500)},
	)
	scale.Exponent(2)

	tests := []struct {
		input    float64
		expected float64
	}{
		{0, 0},         // (0/100)^2 * 500 = 0
		{50, 125},      // (50/100)^2 * 500 = 0.25 * 500 = 125
		{70.7107, 250}, // (70.7107/100)^2 * 500 ≈ 250
		{100, 500},     // (100/100)^2 * 500 = 500
	}

	for _, tt := range tests {
		result := scale.Apply(tt.input)
		if math.Abs(result.Value-tt.expected) > 1.0 { // Allow 1px tolerance
			t.Errorf("Apply(%v) with exponent=2 = %v, expected ~%v", tt.input, result.Value, tt.expected)
		}
	}
}

func TestPowScale_SquareRoot(t *testing.T) {
	scale := NewSqrtScale(
		[2]float64{0, 100},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	tests := []struct {
		input    float64
		expected float64
	}{
		{0, 0},       // sqrt(0/100) * 500 = 0
		{25, 250},    // sqrt(25/100) * 500 = 0.5 * 500 = 250
		{50, 353.55}, // sqrt(50/100) * 500 ≈ 353.55
		{100, 500},   // sqrt(100/100) * 500 = 500
	}

	for _, tt := range tests {
		result := scale.Apply(tt.input)
		if math.Abs(result.Value-tt.expected) > 1.0 { // Allow 1px tolerance
			t.Errorf("Apply(%v) with exponent=0.5 = %v, expected ~%v", tt.input, result.Value, tt.expected)
		}
	}
}

func TestPowScale_Invert(t *testing.T) {
	scale := NewPowScale(
		[2]float64{0, 100},
		[2]units.Length{units.Px(0), units.Px(500)},
	)
	scale.Exponent(2)

	tests := []struct {
		input    float64
		expected float64
	}{
		{0, 0},
		{125, 50},
		{500, 100},
	}

	for _, tt := range tests {
		result := scale.Invert(units.Px(tt.input))
		if math.Abs(result-tt.expected) > 0.1 {
			t.Errorf("Invert(%v) with exponent=2 = %v, expected ~%v", tt.input, result, tt.expected)
		}
	}
}

func TestPowScale_ApplyValue(t *testing.T) {
	scale := NewPowScale(
		[2]float64{0, 100},
		[2]units.Length{units.Px(0), units.Px(500)},
	)
	scale.Exponent(2)

	tests := []struct {
		input    float64
		expected float64
	}{
		{0, 0.0},
		{50, 0.25},  // (50/100)^2 = 0.25
		{70.7107, 0.5},
		{100, 1.0},
	}

	for _, tt := range tests {
		result := scale.ApplyValue(tt.input)
		if math.Abs(result-tt.expected) > 0.01 {
			t.Errorf("ApplyValue(%v) with exponent=2 = %v, expected %v", tt.input, result, tt.expected)
		}
	}
}

func TestPowScale_Clamp(t *testing.T) {
	scale := NewPowScale(
		[2]float64{0, 100},
		[2]units.Length{units.Px(0), units.Px(500)},
	)
	scale.Exponent(2).Clamp(true)

	tests := []struct {
		input    float64
		expected float64
	}{
		{-10, 0},   // Below range
		{0, 0},     // Min
		{100, 500}, // Max
		{110, 500}, // Above range
	}

	for _, tt := range tests {
		result := scale.Apply(tt.input)
		if math.Abs(result.Value-tt.expected) > 0.01 {
			t.Errorf("Apply(%v) with clamp = %v, expected %v", tt.input, result.Value, tt.expected)
		}
	}
}

func TestPowScale_Nice(t *testing.T) {
	scale := NewPowScale(
		[2]float64{3.7, 96.3},
		[2]units.Length{units.Px(0), units.Px(500)},
	)
	scale.Nice(10)

	domain := scale.Domain().([2]float64)

	// Should round to nice boundaries
	if domain[0] != 0 {
		t.Errorf("Nice domain start = %v, expected 0", domain[0])
	}
	if domain[1] != 100 {
		t.Errorf("Nice domain end = %v, expected 100", domain[1])
	}
}

func TestPowScale_Ticks(t *testing.T) {
	scale := NewPowScale(
		[2]float64{0, 100},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	ticks := scale.Ticks(10)

	if len(ticks) == 0 {
		t.Fatal("Expected non-empty ticks")
	}

	// Ticks should be evenly spaced (linear, not transformed)
	if len(ticks) > 1 {
		step := ticks[1] - ticks[0]
		for i := 1; i < len(ticks)-1; i++ {
			actualStep := ticks[i+1] - ticks[i]
			if math.Abs(actualStep-step) > 0.01 {
				t.Errorf("Ticks not evenly spaced: step %v vs %v", actualStep, step)
			}
		}
	}
}

func TestPowScale_CubeRoot(t *testing.T) {
	scale := NewPowScale(
		[2]float64{0, 1000},
		[2]units.Length{units.Px(0), units.Px(500)},
	)
	scale.Exponent(1.0 / 3.0) // Cube root

	tests := []struct {
		input    float64
		expected float64
	}{
		{0, 0},
		{125, 250},  // (125/1000)^(1/3) ≈ 0.5
		{1000, 500},
	}

	for _, tt := range tests {
		result := scale.Apply(tt.input)
		if math.Abs(result.Value-tt.expected) > 1.0 { // Allow 1px tolerance
			t.Errorf("Apply(%v) with exponent=1/3 = %v, expected ~%v", tt.input, result.Value, tt.expected)
		}
	}
}

func TestPowScale_ReverseRange(t *testing.T) {
	scale := NewPowScale(
		[2]float64{0, 100},
		[2]units.Length{units.Px(500), units.Px(0)}, // Reversed
	)
	scale.Exponent(2)

	tests := []struct {
		input    float64
		expected float64
	}{
		{0, 500},
		{50, 375},  // (50/100)^2 * -500 + 500 = 375
		{100, 0},
	}

	for _, tt := range tests {
		result := scale.Apply(tt.input)
		if math.Abs(result.Value-tt.expected) > 1.0 {
			t.Errorf("Apply(%v) with reverse range = %v, expected ~%v", tt.input, result.Value, tt.expected)
		}
	}
}

func TestPowScale_Clone(t *testing.T) {
	original := NewPowScale(
		[2]float64{0, 100},
		[2]units.Length{units.Px(0), units.Px(500)},
	)
	original.Exponent(2).Clamp(true)

	clone := original.Clone().(*PowScale)

	// Verify clone has same values
	if clone.domain != original.domain {
		t.Error("Clone domain doesn't match original")
	}
	if clone.range_ != original.range_ {
		t.Error("Clone range doesn't match original")
	}
	if clone.exponent != original.exponent {
		t.Error("Clone exponent doesn't match original")
	}
	if clone.clamp != original.clamp {
		t.Error("Clone clamp doesn't match original")
	}

	// Verify modifying clone doesn't affect original
	clone.WithDomain([2]float64{0, 200})
	if original.domain[1] == 200 {
		t.Error("Modifying clone affected original")
	}
}

func TestPowScale_Type_Pow(t *testing.T) {
	scale := NewPowScale(
		[2]float64{0, 100},
		[2]units.Length{units.Px(0), units.Px(500)},
	)
	scale.Exponent(2)

	if scale.Type() != ScaleTypePow {
		t.Errorf("Type() with exponent=2 = %v, expected ScaleTypePow", scale.Type())
	}
}

func TestPowScale_Type_Sqrt(t *testing.T) {
	scale := NewSqrtScale(
		[2]float64{0, 100},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	if scale.Type() != ScaleTypeSqrt {
		t.Errorf("Type() with exponent=0.5 = %v, expected ScaleTypeSqrt", scale.Type())
	}
}

func TestPowScale_NegativeDomain(t *testing.T) {
	scale := NewPowScale(
		[2]float64{-100, 100},
		[2]units.Length{units.Px(0), units.Px(500)},
	)
	scale.Exponent(3) // Odd exponent works with negative values

	// Negative values with odd exponents should work
	result := scale.Apply(-50.0)
	if math.IsNaN(result.Value) {
		t.Error("Apply(-50) with exponent=3 should not be NaN")
	}

	// Should be less than midpoint due to cubic transformation
	if result.Value >= 250 {
		t.Errorf("Apply(-50) with exponent=3 = %v, expected < 250", result.Value)
	}
}

func TestPowScale_NegativeWithFractionalExponent(t *testing.T) {
	scale := NewPowScale(
		[2]float64{0, 100},
		[2]units.Length{units.Px(0), units.Px(500)},
	)
	scale.Exponent(0.5) // Fractional exponent

	// Negative input OUTSIDE domain with fractional exponent should return NaN
	// When v=-50, d0=0, d1=100: t = (-50-0)/(100-0) = -0.5 (negative)
	// (-0.5)^0.5 is undefined
	result := scale.ApplyValue(-50.0)
	if !math.IsNaN(result) {
		t.Errorf("ApplyValue(-50) with exponent=0.5 and domain=[0,100] = %v, expected NaN", result)
	}

	// With clamp, negative input should clamp t to 0
	scale.Clamp(true)
	result = scale.ApplyValue(-50.0)
	if result != 0 {
		t.Errorf("ApplyValue(-50) with exponent=0.5 and clamp = %v, expected 0", result)
	}
}

func BenchmarkPowScale_Apply(b *testing.B) {
	scale := NewPowScale(
		[2]float64{0, 100},
		[2]units.Length{units.Px(0), units.Px(500)},
	)
	scale.Exponent(2)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scale.Apply(50.0)
	}
}

func BenchmarkSqrtScale_Apply(b *testing.B) {
	scale := NewSqrtScale(
		[2]float64{0, 100},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scale.Apply(50.0)
	}
}

func BenchmarkPowScale_Ticks(b *testing.B) {
	scale := NewPowScale(
		[2]float64{0, 100},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scale.Ticks(10)
	}
}
