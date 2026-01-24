package scales

import (
	"math"
	"testing"

	"github.com/SCKelemen/units"
)

func TestLinearScale_Basic(t *testing.T) {
	scale := NewLinearScale(
		[2]float64{0, 100},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	tests := []struct {
		input    float64
		expected float64
	}{
		{0, 0},
		{50, 250},
		{100, 500},
		{25, 125},
		{75, 375},
	}

	for _, tt := range tests {
		result := scale.Apply(tt.input)
		if math.Abs(result.Value-tt.expected) > 0.01 {
			t.Errorf("Apply(%v) = %v, expected %v", tt.input, result.Value, tt.expected)
		}
		if result.Unit != units.PX {
			t.Errorf("Apply(%v) unit = %v, expected PX", tt.input, result.Unit)
		}
	}
}

func TestLinearScale_Invert(t *testing.T) {
	scale := NewLinearScale(
		[2]float64{0, 100},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	tests := []struct {
		input    float64
		expected float64
	}{
		{0, 0},
		{250, 50},
		{500, 100},
		{125, 25},
		{375, 75},
	}

	for _, tt := range tests {
		result := scale.Invert(units.Px(tt.input))
		if math.Abs(result-tt.expected) > 0.01 {
			t.Errorf("Invert(%v) = %v, expected %v", tt.input, result, tt.expected)
		}
	}
}

func TestLinearScale_Clamp(t *testing.T) {
	scale := NewLinearScale(
		[2]float64{0, 100},
		[2]units.Length{units.Px(0), units.Px(500)},
	)
	scale.Clamp(true)

	tests := []struct {
		input    float64
		expected float64
	}{
		{-10, 0},   // Clamped to min
		{0, 0},     // Within range
		{50, 250},  // Within range
		{100, 500}, // Within range
		{110, 500}, // Clamped to max
	}

	for _, tt := range tests {
		result := scale.Apply(tt.input)
		if math.Abs(result.Value-tt.expected) > 0.01 {
			t.Errorf("Apply(%v) with clamp = %v, expected %v", tt.input, result.Value, tt.expected)
		}
	}
}

func TestLinearScale_Nice(t *testing.T) {
	scale := NewLinearScale(
		[2]float64{0.123, 96.789},
		[2]units.Length{units.Px(0), units.Px(500)},
	)
	scale.Nice(10)

	domain := scale.Domain().([2]float64)

	// Domain should be rounded to nice boundaries
	if domain[0] != 0 {
		t.Errorf("Nice domain start = %v, expected 0", domain[0])
	}
	if domain[1] != 100 {
		t.Errorf("Nice domain end = %v, expected 100", domain[1])
	}
}

func TestLinearScale_Ticks(t *testing.T) {
	scale := NewLinearScale(
		[2]float64{0, 100},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	ticks := scale.Ticks(10)

	if len(ticks) == 0 {
		t.Fatal("Expected non-empty ticks")
	}

	// First tick should be near domain start
	if math.Abs(ticks[0]-0) > 1 {
		t.Errorf("First tick = %v, expected near 0", ticks[0])
	}

	// Last tick should be near domain end
	if math.Abs(ticks[len(ticks)-1]-100) > 1 {
		t.Errorf("Last tick = %v, expected near 100", ticks[len(ticks)-1])
	}

	// Ticks should be evenly spaced
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

func TestLinearScale_ApplyValue(t *testing.T) {
	scale := NewLinearScale(
		[2]float64{0, 100},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	tests := []struct {
		input    float64
		expected float64
	}{
		{0, 0.0},
		{50, 0.5},
		{100, 1.0},
		{25, 0.25},
		{75, 0.75},
	}

	for _, tt := range tests {
		result := scale.ApplyValue(tt.input)
		if math.Abs(result-tt.expected) > 0.01 {
			t.Errorf("ApplyValue(%v) = %v, expected %v", tt.input, result, tt.expected)
		}
	}
}

func TestLinearScale_NegativeDomain(t *testing.T) {
	scale := NewLinearScale(
		[2]float64{-100, 100},
		[2]units.Length{units.Px(0), units.Px(400)},
	)

	tests := []struct {
		input    float64
		expected float64
	}{
		{-100, 0},
		{0, 200},
		{100, 400},
		{-50, 100},
		{50, 300},
	}

	for _, tt := range tests {
		result := scale.Apply(tt.input)
		if math.Abs(result.Value-tt.expected) > 0.01 {
			t.Errorf("Apply(%v) = %v, expected %v", tt.input, result.Value, tt.expected)
		}
	}
}

func TestLinearScale_ReverseRange(t *testing.T) {
	// Range from 500 to 0 (reversed)
	scale := NewLinearScale(
		[2]float64{0, 100},
		[2]units.Length{units.Px(500), units.Px(0)},
	)

	tests := []struct {
		input    float64
		expected float64
	}{
		{0, 500},
		{50, 250},
		{100, 0},
	}

	for _, tt := range tests {
		result := scale.Apply(tt.input)
		if math.Abs(result.Value-tt.expected) > 0.01 {
			t.Errorf("Apply(%v) with reverse range = %v, expected %v", tt.input, result.Value, tt.expected)
		}
	}
}

func TestLinearScale_Clone(t *testing.T) {
	original := NewLinearScale(
		[2]float64{0, 100},
		[2]units.Length{units.Px(0), units.Px(500)},
	)
	original.Clamp(true)

	clone := original.Clone().(*LinearScale)

	// Verify clone has same values
	if clone.domain != original.domain {
		t.Error("Clone domain doesn't match original")
	}
	if clone.range_ != original.range_ {
		t.Error("Clone range doesn't match original")
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

func BenchmarkLinearScale_Apply(b *testing.B) {
	scale := NewLinearScale(
		[2]float64{0, 100},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scale.Apply(50.0)
	}
}

func BenchmarkLinearScale_Ticks(b *testing.B) {
	scale := NewLinearScale(
		[2]float64{0, 100},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scale.Ticks(10)
	}
}
