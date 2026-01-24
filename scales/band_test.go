package scales

import (
	"math"
	"testing"

	"github.com/SCKelemen/units"
)

func TestBandScale_Basic(t *testing.T) {
	scale := NewBandScale(
		[]string{"A", "B", "C"},
		[2]units.Length{units.Px(0), units.Px(300)},
	)

	// With 3 bands and no padding, each band should be 100px
	tests := []struct {
		input    string
		expected float64
	}{
		{"A", 0},
		{"B", 100},
		{"C", 200},
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

	// Check bandwidth
	bandwidth := scale.Bandwidth()
	if math.Abs(bandwidth.Value-100) > 0.01 {
		t.Errorf("Bandwidth() = %v, expected 100", bandwidth.Value)
	}
}

func TestBandScale_WithPadding(t *testing.T) {
	scale := NewBandScale(
		[]string{"A", "B", "C"},
		[2]units.Length{units.Px(0), units.Px(300)},
	)
	scale.Padding(0.1)

	// Padding should reduce bandwidth
	bandwidth := scale.Bandwidth()
	if bandwidth.Value >= 100 {
		t.Errorf("Bandwidth with padding = %v, expected < 100", bandwidth.Value)
	}

	// All values should still map within range
	for _, category := range []string{"A", "B", "C"} {
		result := scale.Apply(category)
		if result.Value < 0 || result.Value > 300 {
			t.Errorf("Apply(%q) = %v, outside range [0, 300]", category, result.Value)
		}
	}
}

func TestBandScale_PaddingInnerOuter(t *testing.T) {
	scale := NewBandScale(
		[]string{"A", "B"},
		[2]units.Length{units.Px(0), units.Px(200)},
	)
	scale.PaddingInner(0.2).PaddingOuter(0.1)

	// Both bands should be within range
	resultA := scale.Apply("A")
	resultB := scale.Apply("B")

	if resultA.Value < 0 {
		t.Errorf("Apply(A) = %v, expected >= 0", resultA.Value)
	}

	bandwidth := scale.Bandwidth()
	if resultB.Value+bandwidth.Value > 200 {
		t.Errorf("Apply(B) + bandwidth = %v, expected <= 200", resultB.Value+bandwidth.Value)
	}
}

func TestBandScale_UnknownValue(t *testing.T) {
	scale := NewBandScale(
		[]string{"A", "B", "C"},
		[2]units.Length{units.Px(0), units.Px(300)},
	)

	// Unknown value should return 0
	result := scale.Apply("D")
	if result.Value != 0 {
		t.Errorf("Apply(unknown) = %v, expected 0", result.Value)
	}
}

func TestBandScale_Index(t *testing.T) {
	scale := NewBandScale(
		[]string{"A", "B", "C"},
		[2]units.Length{units.Px(0), units.Px(300)},
	)

	tests := []struct {
		input    string
		expected int
	}{
		{"A", 0},
		{"B", 1},
		{"C", 2},
		{"D", -1}, // Not found
	}

	for _, tt := range tests {
		result := scale.Index(tt.input)
		if result != tt.expected {
			t.Errorf("Index(%q) = %v, expected %v", tt.input, result, tt.expected)
		}
	}
}

func TestBandScale_Values(t *testing.T) {
	domain := []string{"A", "B", "C"}
	scale := NewBandScale(
		domain,
		[2]units.Length{units.Px(0), units.Px(300)},
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

func TestBandScale_Round(t *testing.T) {
	scale := NewBandScale(
		[]string{"A", "B", "C"},
		[2]units.Length{units.Px(0), units.Px(299)}, // Odd number to force fractional
	)
	scale.Round(true)

	// Positions should be rounded to integers
	for _, category := range []string{"A", "B", "C"} {
		result := scale.Apply(category)
		if result.Value != math.Round(result.Value) {
			t.Errorf("Apply(%q) = %v, expected integer with Round(true)", category, result.Value)
		}
	}

	// Bandwidth should also be rounded
	bandwidth := scale.Bandwidth()
	if bandwidth.Value != math.Round(bandwidth.Value) {
		t.Errorf("Bandwidth() = %v, expected integer with Round(true)", bandwidth.Value)
	}
}

func TestBandScale_Align(t *testing.T) {
	// Test different alignment values
	for _, align := range []float64{0.0, 0.5, 1.0} {
		scale := NewBandScale(
			[]string{"A", "B"},
			[2]units.Length{units.Px(0), units.Px(200)},
		)
		scale.Padding(0.2).Align(align)

		// All bands should still be within range
		for _, category := range []string{"A", "B"} {
			result := scale.Apply(category)
			bandwidth := scale.Bandwidth()

			if result.Value < 0 {
				t.Errorf("Align(%v): Apply(%q) = %v, expected >= 0", align, category, result.Value)
			}
			if result.Value+bandwidth.Value > 200 {
				t.Errorf("Align(%v): Apply(%q) + bandwidth = %v, expected <= 200",
					align, category, result.Value+bandwidth.Value)
			}
		}
	}
}

func TestBandScale_Step(t *testing.T) {
	scale := NewBandScale(
		[]string{"A", "B", "C"},
		[2]units.Length{units.Px(0), units.Px(300)},
	)

	step := scale.Step()
	bandwidth := scale.Bandwidth()

	// Step should be bandwidth (since no padding)
	if math.Abs(step.Value-bandwidth.Value) > 0.01 {
		t.Errorf("Step() = %v, Bandwidth() = %v, expected equal with no padding",
			step.Value, bandwidth.Value)
	}

	// With padding, step should be greater than bandwidth
	scale.Padding(0.1)
	step = scale.Step()
	bandwidth = scale.Bandwidth()

	if step.Value <= bandwidth.Value {
		t.Errorf("Step() = %v, Bandwidth() = %v, expected step > bandwidth with padding",
			step.Value, bandwidth.Value)
	}
}

func TestBandScale_EmptyDomain(t *testing.T) {
	scale := NewBandScale(
		[]string{},
		[2]units.Length{units.Px(0), units.Px(300)},
	)

	// Should not panic
	result := scale.Apply("A")
	if result.Value != 0 {
		t.Errorf("Apply with empty domain = %v, expected 0", result.Value)
	}

	bandwidth := scale.Bandwidth()
	if bandwidth.Value != 0 {
		t.Errorf("Bandwidth with empty domain = %v, expected 0", bandwidth.Value)
	}
}

func TestBandScale_ReverseRange(t *testing.T) {
	scale := NewBandScale(
		[]string{"A", "B", "C"},
		[2]units.Length{units.Px(300), units.Px(0)}, // Reversed
	)

	// First band should be near 300
	resultA := scale.Apply("A")
	if resultA.Value < 200 {
		t.Errorf("Apply(A) with reverse range = %v, expected near 300", resultA.Value)
	}

	// Last band should be near 0
	resultC := scale.Apply("C")
	bandwidth := scale.Bandwidth()
	if resultC.Value-bandwidth.Value > 100 {
		t.Errorf("Apply(C) with reverse range = %v, expected near 0", resultC.Value)
	}
}

func TestBandScale_Clone(t *testing.T) {
	original := NewBandScale(
		[]string{"A", "B", "C"},
		[2]units.Length{units.Px(0), units.Px(300)},
	)
	original.Padding(0.1)

	clone := original.Clone().(*BandScale)

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

	// Verify modifying clone doesn't affect original
	clone.WithDomain([]string{"X", "Y"})
	if len(original.domain) == 2 {
		t.Error("Modifying clone affected original")
	}
}

func TestBandScale_ApplyValue(t *testing.T) {
	scale := NewBandScale(
		[]string{"A", "B", "C"},
		[2]units.Length{units.Px(0), units.Px(300)},
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

func BenchmarkBandScale_Apply(b *testing.B) {
	scale := NewBandScale(
		[]string{"A", "B", "C", "D", "E"},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scale.Apply("C")
	}
}

func BenchmarkBandScale_Rescale(b *testing.B) {
	scale := NewBandScale(
		[]string{"A", "B", "C"},
		[2]units.Length{units.Px(0), units.Px(300)},
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scale.Padding(0.1)
	}
}
