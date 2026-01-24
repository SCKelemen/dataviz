package scales

import (
	"math"
	"testing"

	"github.com/SCKelemen/units"
)

func TestOrdinalScale_Basic(t *testing.T) {
	scale := NewOrdinalScale(
		[]string{"small", "medium", "large"},
		[]units.Length{units.Px(10), units.Px(20), units.Px(30)},
	)

	tests := []struct {
		input    string
		expected float64
	}{
		{"small", 10},
		{"medium", 20},
		{"large", 30},
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

func TestOrdinalScale_Cycling(t *testing.T) {
	// Domain larger than range - should cycle through range values
	scale := NewOrdinalScale(
		[]string{"A", "B", "C", "D", "E"},
		[]units.Length{units.Px(10), units.Px(20), units.Px(30)},
	)

	tests := []struct {
		input    string
		expected float64
	}{
		{"A", 10}, // index 0 % 3 = 0
		{"B", 20}, // index 1 % 3 = 1
		{"C", 30}, // index 2 % 3 = 2
		{"D", 10}, // index 3 % 3 = 0 (cycle)
		{"E", 20}, // index 4 % 3 = 1 (cycle)
	}

	for _, tt := range tests {
		result := scale.Apply(tt.input)
		if math.Abs(result.Value-tt.expected) > 0.01 {
			t.Errorf("Apply(%q) = %v, expected %v", tt.input, result.Value, tt.expected)
		}
	}
}

func TestOrdinalScale_Unknown(t *testing.T) {
	scale := NewOrdinalScale(
		[]string{"A", "B", "C"},
		[]units.Length{units.Px(10), units.Px(20), units.Px(30)},
	)

	// Default unknown value
	result := scale.Apply("D")
	if result.Value != 0 {
		t.Errorf("Apply(unknown) = %v, expected 0 (default)", result.Value)
	}

	// Custom unknown value
	scale.Unknown(units.Px(999))
	result = scale.Apply("D")
	if result.Value != 999 {
		t.Errorf("Apply(unknown) = %v, expected 999", result.Value)
	}
}

func TestOrdinalScale_Index(t *testing.T) {
	scale := NewOrdinalScale(
		[]string{"alpha", "beta", "gamma"},
		[]units.Length{units.Px(10), units.Px(20), units.Px(30)},
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

func TestOrdinalScale_Values(t *testing.T) {
	domain := []string{"red", "green", "blue"}
	scale := NewOrdinalScale(
		domain,
		[]units.Length{units.Px(10), units.Px(20), units.Px(30)},
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

func TestOrdinalScale_RangeValues(t *testing.T) {
	rangeVals := []units.Length{units.Px(10), units.Px(20), units.Px(30)}
	scale := NewOrdinalScale(
		[]string{"A", "B", "C"},
		rangeVals,
	)

	result := scale.RangeValues()
	if len(result) != len(rangeVals) {
		t.Errorf("RangeValues() length = %v, expected %v", len(result), len(rangeVals))
	}

	for i, v := range result {
		if v.Value != rangeVals[i].Value || v.Unit != rangeVals[i].Unit {
			t.Errorf("RangeValues()[%d] = %v, expected %v", i, v, rangeVals[i])
		}
	}
}

func TestOrdinalScale_Range(t *testing.T) {
	scale := NewOrdinalScale(
		[]string{"A", "B", "C"},
		[]units.Length{units.Px(10), units.Px(20), units.Px(30)},
	)

	r := scale.Range()
	if r[0].Value != 10 || r[1].Value != 30 {
		t.Errorf("Range() = [%v, %v], expected [10, 30]", r[0].Value, r[1].Value)
	}
}

func TestOrdinalScale_ApplyValue(t *testing.T) {
	scale := NewOrdinalScale(
		[]string{"A", "B", "C"},
		[]units.Length{units.Px(10), units.Px(20), units.Px(30)},
	)

	tests := []struct {
		input    string
		expected float64
	}{
		{"A", 0.0},  // First item
		{"B", 0.5},  // Middle item
		{"C", 1.0},  // Last item
	}

	for _, tt := range tests {
		result := scale.ApplyValue(tt.input)
		if math.Abs(result-tt.expected) > 0.01 {
			t.Errorf("ApplyValue(%q) = %v, expected %v", tt.input, result, tt.expected)
		}
	}
}

func TestOrdinalScale_ApplyValue_SingleValue(t *testing.T) {
	scale := NewOrdinalScale(
		[]string{"A"},
		[]units.Length{units.Px(10)},
	)

	// Single value should map to center (0.5)
	result := scale.ApplyValue("A")
	if math.Abs(result-0.5) > 0.01 {
		t.Errorf("ApplyValue(single) = %v, expected 0.5", result)
	}
}

func TestOrdinalScale_Clone(t *testing.T) {
	original := NewOrdinalScale(
		[]string{"A", "B", "C"},
		[]units.Length{units.Px(10), units.Px(20), units.Px(30)},
	)
	original.Unknown(units.Px(999))

	clone := original.Clone().(*OrdinalScale)

	// Verify clone has same values
	if len(clone.domain) != len(original.domain) {
		t.Error("Clone domain length doesn't match original")
	}
	if len(clone.range_) != len(original.range_) {
		t.Error("Clone range length doesn't match original")
	}
	if clone.unknown != original.unknown {
		t.Error("Clone unknown doesn't match original")
	}

	// Verify modifying clone doesn't affect original
	clone.WithDomain([]string{"X", "Y"})
	if len(original.domain) == 2 {
		t.Error("Modifying clone affected original")
	}
}

func TestOrdinalScale_WithDomain(t *testing.T) {
	scale := NewOrdinalScale(
		[]string{"A", "B"},
		[]units.Length{units.Px(10), units.Px(20)},
	)

	scale.WithDomain([]string{"X", "Y", "Z"})

	if len(scale.domain) != 3 {
		t.Errorf("Domain length = %v, expected 3", len(scale.domain))
	}

	result := scale.Apply("Z")
	if result.Value != 10 { // Cycles back to first range value
		t.Errorf("Apply(Z) = %v, expected 10 (cycling)", result.Value)
	}
}

func TestOrdinalScale_WithRange(t *testing.T) {
	scale := NewOrdinalScale(
		[]string{"A", "B", "C"},
		[]units.Length{units.Px(10), units.Px(20)},
	)

	scale.WithRange([]units.Length{units.Px(100), units.Px(200), units.Px(300)})

	result := scale.Apply("C")
	if result.Value != 300 {
		t.Errorf("Apply(C) = %v, expected 300", result.Value)
	}
}

func TestOrdinalScale_EmptyDomain(t *testing.T) {
	scale := NewOrdinalScale(
		[]string{},
		[]units.Length{units.Px(10), units.Px(20)},
	)

	// Should not panic
	result := scale.Apply("A")
	if result.Value != 0 {
		t.Errorf("Apply with empty domain = %v, expected 0", result.Value)
	}

	normalized := scale.ApplyValue("A")
	if normalized != 0 {
		t.Errorf("ApplyValue with empty domain = %v, expected 0", normalized)
	}
}

func TestOrdinalScale_EmptyRange(t *testing.T) {
	scale := NewOrdinalScale(
		[]string{"A", "B"},
		[]units.Length{},
	)

	// Should not panic - returns unknown value
	result := scale.Apply("A")
	if result.Value != 0 {
		t.Errorf("Apply with empty range = %v, expected 0", result.Value)
	}
}

func TestOrdinalScale_Type(t *testing.T) {
	scale := NewOrdinalScale(
		[]string{"A"},
		[]units.Length{units.Px(10)},
	)

	if scale.Type() != ScaleTypeOrdinal {
		t.Errorf("Type() = %v, expected ScaleTypeOrdinal", scale.Type())
	}
}

func BenchmarkOrdinalScale_Apply(b *testing.B) {
	scale := NewOrdinalScale(
		[]string{"A", "B", "C", "D", "E"},
		[]units.Length{units.Px(10), units.Px(20), units.Px(30), units.Px(40), units.Px(50)},
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scale.Apply("C")
	}
}
