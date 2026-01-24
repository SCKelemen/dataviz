package scales

import (
	"math"
	"testing"

	"github.com/SCKelemen/units"
)

func TestPointScale_Basic(t *testing.T) {
	scale := NewPointScale(
		[]string{"A", "B", "C"},
		[2]units.Length{units.Px(0), units.Px(100)},
	)

	tests := []struct {
		input    string
		expected float64
	}{
		{"A", 0},   // First point
		{"B", 50},  // Middle point
		{"C", 100}, // Last point
	}

	for _, tt := range tests {
		result := scale.Apply(tt.input)
		if math.Abs(result.Value-tt.expected) > 0.01 {
			t.Errorf("Apply(%q) = %v, expected %v", tt.input, result.Value, tt.expected)
		}
		if result.Unit != units.PX {
			t.Errorf("Apply(%q) unit = %v, expected PX", tt.input, result.Unit)
		}
	}
}

func TestPointScale_WithPadding(t *testing.T) {
	scale := NewPointScale(
		[]string{"A", "B", "C"},
		[2]units.Length{units.Px(0), units.Px(100)},
	)
	scale.Padding(0.5)

	// With padding, points should be inset from edges
	resultA := scale.Apply("A")
	if resultA.Value <= 0 {
		t.Errorf("Apply(A) with padding = %v, expected > 0", resultA.Value)
	}

	resultC := scale.Apply("C")
	if resultC.Value >= 100 {
		t.Errorf("Apply(C) with padding = %v, expected < 100", resultC.Value)
	}

	// Points should still be evenly spaced
	resultB := scale.Apply("B")
	step1 := resultB.Value - resultA.Value
	step2 := resultC.Value - resultB.Value
	if math.Abs(step1-step2) > 0.01 {
		t.Errorf("Points not evenly spaced: %v vs %v", step1, step2)
	}
}

func TestPointScale_Align(t *testing.T) {
	// Test different alignment values
	for _, align := range []float64{0.0, 0.5, 1.0} {
		scale := NewPointScale(
			[]string{"A", "B"},
			[2]units.Length{units.Px(0), units.Px(100)},
		)
		scale.Padding(0.2).Align(align)

		// All points should be within range
		for _, category := range []string{"A", "B"} {
			result := scale.Apply(category)

			if result.Value < 0 {
				t.Errorf("Align(%v): Apply(%q) = %v, expected >= 0", align, category, result.Value)
			}
			if result.Value > 100 {
				t.Errorf("Align(%v): Apply(%q) = %v, expected <= 100", align, category, result.Value)
			}
		}
	}
}

func TestPointScale_Round(t *testing.T) {
	scale := NewPointScale(
		[]string{"A", "B", "C"},
		[2]units.Length{units.Px(0), units.Px(99)}, // Odd number to force fractional
	)
	scale.Round(true)

	// Positions should be rounded to integers
	for _, category := range []string{"A", "B", "C"} {
		result := scale.Apply(category)
		if result.Value != math.Round(result.Value) {
			t.Errorf("Apply(%q) = %v, expected integer with Round(true)", category, result.Value)
		}
	}
}

func TestPointScale_SingleValue(t *testing.T) {
	scale := NewPointScale(
		[]string{"A"},
		[2]units.Length{units.Px(0), units.Px(100)},
	)

	// Single point should be centered (default align = 0.5)
	result := scale.Apply("A")
	if math.Abs(result.Value-50) > 0.01 {
		t.Errorf("Apply(single) = %v, expected 50 (centered)", result.Value)
	}

	// Test with different alignment
	scale.Align(0.0)
	result = scale.Apply("A")
	if math.Abs(result.Value-0) > 0.01 {
		t.Errorf("Apply(single, align=0) = %v, expected 0", result.Value)
	}

	scale.Align(1.0)
	result = scale.Apply("A")
	if math.Abs(result.Value-100) > 0.01 {
		t.Errorf("Apply(single, align=1) = %v, expected 100", result.Value)
	}
}

func TestPointScale_Step(t *testing.T) {
	scale := NewPointScale(
		[]string{"A", "B", "C"},
		[2]units.Length{units.Px(0), units.Px(100)},
	)

	step := scale.Step()

	// For 3 points spanning 100px, step should be 50
	if math.Abs(step.Value-50) > 0.01 {
		t.Errorf("Step() = %v, expected 50", step.Value)
	}

	// With padding, step should be smaller
	scale.Padding(0.5)
	stepWithPadding := scale.Step()

	if stepWithPadding.Value >= step.Value {
		t.Errorf("Step with padding = %v, expected < %v", stepWithPadding.Value, step.Value)
	}
}

func TestPointScale_Index(t *testing.T) {
	scale := NewPointScale(
		[]string{"alpha", "beta", "gamma"},
		[2]units.Length{units.Px(0), units.Px(100)},
	)

	tests := []struct {
		input    string
		expected int
	}{
		{"alpha", 0},
		{"beta", 1},
		{"gamma", 2},
		{"delta", -1}, // Not found
	}

	for _, tt := range tests {
		result := scale.Index(tt.input)
		if result != tt.expected {
			t.Errorf("Index(%q) = %v, expected %v", tt.input, result, tt.expected)
		}
	}
}

func TestPointScale_Values(t *testing.T) {
	domain := []string{"red", "green", "blue"}
	scale := NewPointScale(
		domain,
		[2]units.Length{units.Px(0), units.Px(100)},
	)

	values := scale.Values()
	if len(values) != len(domain) {
		t.Errorf("Values() length = %v, expected %v", len(values), len(domain))
	}

	for i, v := range values {
		if v != domain[i] {
			t.Errorf("Values()[%d] = %q, expected %q", i, v, domain[i])
		}
	}
}

func TestPointScale_ApplyValue(t *testing.T) {
	scale := NewPointScale(
		[]string{"A", "B", "C"},
		[2]units.Length{units.Px(0), units.Px(100)},
	)

	// ApplyValue should return normalized position
	resultA := scale.ApplyValue("A")
	if resultA < 0 || resultA > 1 {
		t.Errorf("ApplyValue(A) = %v, expected in range [0, 1]", resultA)
	}

	resultB := scale.ApplyValue("B")
	if resultB <= resultA {
		t.Errorf("ApplyValue(B) = %v, expected > ApplyValue(A) = %v", resultB, resultA)
	}

	resultC := scale.ApplyValue("C")
	if resultC <= resultB {
		t.Errorf("ApplyValue(C) = %v, expected > ApplyValue(B) = %v", resultC, resultB)
	}
}

func TestPointScale_UnknownValue(t *testing.T) {
	scale := NewPointScale(
		[]string{"A", "B", "C"},
		[2]units.Length{units.Px(0), units.Px(100)},
	)

	// Unknown value should return 0
	result := scale.Apply("D")
	if result.Value != 0 {
		t.Errorf("Apply(unknown) = %v, expected 0", result.Value)
	}
}

func TestPointScale_EmptyDomain(t *testing.T) {
	scale := NewPointScale(
		[]string{},
		[2]units.Length{units.Px(0), units.Px(100)},
	)

	// Should not panic
	result := scale.Apply("A")
	if result.Value != 0 {
		t.Errorf("Apply with empty domain = %v, expected 0", result.Value)
	}

	step := scale.Step()
	if step.Value != 0 {
		t.Errorf("Step with empty domain = %v, expected 0", step.Value)
	}
}

func TestPointScale_ReverseRange(t *testing.T) {
	scale := NewPointScale(
		[]string{"A", "B", "C"},
		[2]units.Length{units.Px(100), units.Px(0)}, // Reversed
	)

	// First point should be near 100
	resultA := scale.Apply("A")
	if resultA.Value < 50 {
		t.Errorf("Apply(A) with reverse range = %v, expected near 100", resultA.Value)
	}

	// Last point should be near 0
	resultC := scale.Apply("C")
	if resultC.Value > 50 {
		t.Errorf("Apply(C) with reverse range = %v, expected near 0", resultC.Value)
	}

	// Middle point should be in between
	resultB := scale.Apply("B")
	if resultB.Value <= resultC.Value || resultB.Value >= resultA.Value {
		t.Errorf("Apply(B) = %v, expected between A=%v and C=%v", resultB.Value, resultA.Value, resultC.Value)
	}
}

func TestPointScale_Clone(t *testing.T) {
	original := NewPointScale(
		[]string{"A", "B", "C"},
		[2]units.Length{units.Px(0), units.Px(100)},
	)
	original.Padding(0.1).Round(true)

	clone := original.Clone().(*PointScale)

	// Verify clone has same values
	if len(clone.domain) != len(original.domain) {
		t.Error("Clone domain length doesn't match original")
	}
	if clone.range_ != original.range_ {
		t.Error("Clone range doesn't match original")
	}
	if clone.padding != original.padding {
		t.Error("Clone padding doesn't match original")
	}
	if clone.round != original.round {
		t.Error("Clone round doesn't match original")
	}

	// Verify modifying clone doesn't affect original
	clone.WithDomain([]string{"X", "Y"})
	if len(original.domain) == 2 {
		t.Error("Modifying clone affected original")
	}
}

func TestPointScale_TwoPoints(t *testing.T) {
	scale := NewPointScale(
		[]string{"A", "B"},
		[2]units.Length{units.Px(0), units.Px(100)},
	)

	resultA := scale.Apply("A")
	resultB := scale.Apply("B")

	// Two points should be at the edges
	if math.Abs(resultA.Value-0) > 0.01 {
		t.Errorf("Apply(A) = %v, expected 0", resultA.Value)
	}
	if math.Abs(resultB.Value-100) > 0.01 {
		t.Errorf("Apply(B) = %v, expected 100", resultB.Value)
	}

	// Step should be the full range
	step := scale.Step()
	if math.Abs(step.Value-100) > 0.01 {
		t.Errorf("Step() = %v, expected 100", step.Value)
	}
}

func TestPointScale_Type(t *testing.T) {
	scale := NewPointScale(
		[]string{"A"},
		[2]units.Length{units.Px(0), units.Px(100)},
	)

	if scale.Type() != ScaleTypePoint {
		t.Errorf("Type() = %v, expected ScaleTypePoint", scale.Type())
	}
}

func BenchmarkPointScale_Apply(b *testing.B) {
	scale := NewPointScale(
		[]string{"A", "B", "C", "D", "E"},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scale.Apply("C")
	}
}

func BenchmarkPointScale_Rescale(b *testing.B) {
	scale := NewPointScale(
		[]string{"A", "B", "C"},
		[2]units.Length{units.Px(0), units.Px(100)},
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scale.Padding(0.1)
	}
}
