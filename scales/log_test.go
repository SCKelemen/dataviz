package scales

import (
	"math"
	"testing"

	"github.com/SCKelemen/units"
)

func TestLogScale_Basic(t *testing.T) {
	scale := NewLogScale(
		[2]float64{1, 1000},
		[2]units.Length{units.Px(0), units.Px(300)},
	)

	tests := []struct {
		input    float64
		expected float64
	}{
		{1, 0},       // log10(1) = 0
		{10, 100},    // log10(10) = 1, 1/3 of way
		{100, 200},   // log10(100) = 2, 2/3 of way
		{1000, 300},  // log10(1000) = 3
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

func TestLogScale_Invert(t *testing.T) {
	scale := NewLogScale(
		[2]float64{1, 1000},
		[2]units.Length{units.Px(0), units.Px(300)},
	)

	tests := []struct {
		input    float64
		expected float64
	}{
		{0, 1},
		{100, 10},
		{200, 100},
		{300, 1000},
	}

	for _, tt := range tests {
		result := scale.Invert(units.Px(tt.input))
		if math.Abs(result-tt.expected) > 0.01 {
			t.Errorf("Invert(%v) = %v, expected %v", tt.input, result, tt.expected)
		}
	}
}

func TestLogScale_Base(t *testing.T) {
	// Test with base 2
	scale := NewLogScale(
		[2]float64{1, 16},
		[2]units.Length{units.Px(0), units.Px(400)},
	)
	scale.Base(2)

	tests := []struct {
		input    float64
		expected float64
	}{
		{1, 0},     // log2(1) = 0
		{2, 100},   // log2(2) = 1, 1/4 of way
		{4, 200},   // log2(4) = 2, 2/4 of way
		{8, 300},   // log2(8) = 3, 3/4 of way
		{16, 400},  // log2(16) = 4
	}

	for _, tt := range tests {
		result := scale.Apply(tt.input)
		if math.Abs(result.Value-tt.expected) > 0.01 {
			t.Errorf("Apply(%v) with base 2 = %v, expected %v", tt.input, result.Value, tt.expected)
		}
	}
}

func TestLogScale_Nice(t *testing.T) {
	scale := NewLogScale(
		[2]float64{3, 850},
		[2]units.Length{units.Px(0), units.Px(500)},
	)
	scale.Nice(10)

	domain := scale.Domain().([2]float64)

	// Should round to powers of 10
	if domain[0] != 1 {
		t.Errorf("Nice domain start = %v, expected 1", domain[0])
	}
	if domain[1] != 1000 {
		t.Errorf("Nice domain end = %v, expected 1000", domain[1])
	}
}

func TestLogScale_Ticks(t *testing.T) {
	scale := NewLogScale(
		[2]float64{1, 1000},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	ticks := scale.Ticks(10)

	if len(ticks) == 0 {
		t.Fatal("Expected non-empty ticks")
	}

	// Should include powers of 10
	expectedPowers := []float64{1, 10, 100, 1000}
	for _, expected := range expectedPowers {
		found := false
		for _, tick := range ticks {
			if math.Abs(tick-expected) < 0.01 {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Ticks missing power of 10: %v", expected)
		}
	}

	// Ticks should be in ascending order
	for i := 1; i < len(ticks); i++ {
		if ticks[i] <= ticks[i-1] {
			t.Errorf("Ticks not in ascending order at index %d: %v <= %v", i, ticks[i], ticks[i-1])
		}
	}
}

func TestLogScale_Ticks_IntermediateValues(t *testing.T) {
	scale := NewLogScale(
		[2]float64{1, 100},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	ticks := scale.Ticks(20)

	// With more ticks requested, should include intermediate values (2, 3, 4, 5, 6, 7, 8, 9, 20, 30, ...)
	if len(ticks) < 10 {
		t.Errorf("Expected more ticks with count=20, got %d", len(ticks))
	}

	// Should include some intermediate values like 2, 5, 20, 50
	intermediateValues := []float64{2, 5, 20, 50}
	foundCount := 0
	for _, expected := range intermediateValues {
		for _, tick := range ticks {
			if math.Abs(tick-expected) < 0.01 {
				foundCount++
				break
			}
		}
	}

	if foundCount < 2 {
		t.Errorf("Expected at least 2 intermediate tick values, found %d", foundCount)
	}
}

func TestLogScale_Clamp(t *testing.T) {
	scale := NewLogScale(
		[2]float64{10, 100},
		[2]units.Length{units.Px(0), units.Px(100)},
	)
	scale.Clamp(true)

	tests := []struct {
		input    float64
		expected float64
	}{
		{1, 0},     // Below range, clamped to min
		{10, 0},    // Min
		{100, 100}, // Max
		{1000, 100}, // Above range, clamped to max
	}

	for _, tt := range tests {
		result := scale.Apply(tt.input)
		if math.Abs(result.Value-tt.expected) > 0.01 {
			t.Errorf("Apply(%v) with clamp = %v, expected %v", tt.input, result.Value, tt.expected)
		}
	}
}

func TestLogScale_ZeroAndNegative(t *testing.T) {
	scale := NewLogScale(
		[2]float64{1, 100},
		[2]units.Length{units.Px(0), units.Px(100)},
	)

	// Zero should return NaN (log of 0 is undefined)
	result := scale.ApplyValue(0)
	if !math.IsNaN(result) {
		t.Errorf("ApplyValue(0) = %v, expected NaN", result)
	}

	// Negative should return NaN (log of negative is undefined)
	result = scale.ApplyValue(-10)
	if !math.IsNaN(result) {
		t.Errorf("ApplyValue(-10) = %v, expected NaN", result)
	}

	// With clamp, zero/negative should return 0
	scale.Clamp(true)
	result = scale.ApplyValue(0)
	if result != 0 {
		t.Errorf("ApplyValue(0) with clamp = %v, expected 0", result)
	}

	result = scale.ApplyValue(-10)
	if result != 0 {
		t.Errorf("ApplyValue(-10) with clamp = %v, expected 0", result)
	}
}

func TestLogScale_ApplyValue(t *testing.T) {
	scale := NewLogScale(
		[2]float64{1, 1000},
		[2]units.Length{units.Px(0), units.Px(300)},
	)

	tests := []struct {
		input    float64
		expected float64
	}{
		{1, 0.0},        // log10(1) = 0
		{10, 0.333333},  // log10(10) = 1, 1/3 of log range
		{100, 0.666666}, // log10(100) = 2, 2/3 of log range
		{1000, 1.0},     // log10(1000) = 3
	}

	for _, tt := range tests {
		result := scale.ApplyValue(tt.input)
		if math.Abs(result-tt.expected) > 0.01 {
			t.Errorf("ApplyValue(%v) = %v, expected %v", tt.input, result, tt.expected)
		}
	}
}

func TestLogScale_ReverseRange(t *testing.T) {
	scale := NewLogScale(
		[2]float64{1, 1000},
		[2]units.Length{units.Px(300), units.Px(0)}, // Reversed
	)

	tests := []struct {
		input    float64
		expected float64
	}{
		{1, 300},
		{10, 200},
		{100, 100},
		{1000, 0},
	}

	for _, tt := range tests {
		result := scale.Apply(tt.input)
		if math.Abs(result.Value-tt.expected) > 0.01 {
			t.Errorf("Apply(%v) with reverse range = %v, expected %v", tt.input, result.Value, tt.expected)
		}
	}
}

func TestLogScale_Clone(t *testing.T) {
	original := NewLogScale(
		[2]float64{1, 1000},
		[2]units.Length{units.Px(0), units.Px(500)},
	)
	original.Base(2).Clamp(true)

	clone := original.Clone().(*LogScale)

	// Verify clone has same values
	if clone.domain != original.domain {
		t.Error("Clone domain doesn't match original")
	}
	if clone.range_ != original.range_ {
		t.Error("Clone range doesn't match original")
	}
	if clone.base != original.base {
		t.Error("Clone base doesn't match original")
	}
	if clone.clamp != original.clamp {
		t.Error("Clone clamp doesn't match original")
	}

	// Verify modifying clone doesn't affect original
	clone.WithDomain([2]float64{1, 10000})
	if original.domain[1] == 10000 {
		t.Error("Modifying clone affected original")
	}
}

func TestLogScale_Type(t *testing.T) {
	scale := NewLogScale(
		[2]float64{1, 1000},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	if scale.Type() != ScaleTypeLog {
		t.Errorf("Type() = %v, expected ScaleTypeLog", scale.Type())
	}
}

func BenchmarkLogScale_Apply(b *testing.B) {
	scale := NewLogScale(
		[2]float64{1, 1000},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scale.Apply(100.0)
	}
}

func BenchmarkLogScale_Ticks(b *testing.B) {
	scale := NewLogScale(
		[2]float64{1, 1000},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scale.Ticks(10)
	}
}
