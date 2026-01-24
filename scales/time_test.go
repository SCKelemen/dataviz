package scales

import (
	"math"
	"testing"
	"time"

	"github.com/SCKelemen/units"
)

func TestTimeScale_Basic(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

	scale := NewTimeScale(
		[2]time.Time{start, end},
		[2]units.Length{units.Px(0), units.Px(365)},
	)

	tests := []struct {
		input    time.Time
		expected float64
	}{
		{start, 0},
		{time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC), 182.5}, // Mid-year (approximately)
		{end, 365},
	}

	for _, tt := range tests {
		result := scale.Apply(tt.input)
		if math.Abs(result.Value-tt.expected) > 1.0 { // Allow 1px tolerance
			t.Errorf("Apply(%v) = %v, expected ~%v", tt.input.Format("2006-01-02"), result.Value, tt.expected)
		}
		if result.Unit != units.PX {
			t.Errorf("Apply(%v) unit = %v, expected PX", tt.input.Format("2006-01-02"), result.Unit)
		}
	}
}

func TestTimeScale_Invert(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

	scale := NewTimeScale(
		[2]time.Time{start, end},
		[2]units.Length{units.Px(0), units.Px(365)},
	)

	tests := []struct {
		input    float64
		expected time.Time
	}{
		{0, start},
		{182.5, time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)},
		{365, end},
	}

	for _, tt := range tests {
		result := scale.Invert(units.Px(tt.input))
		// Allow 1 day tolerance due to rounding
		diff := math.Abs(result.Sub(tt.expected).Hours())
		if diff > 24 {
			t.Errorf("Invert(%v) = %v, expected ~%v (diff: %.1f hours)",
				tt.input, result.Format("2006-01-02"), tt.expected.Format("2006-01-02"), diff)
		}
	}
}

func TestTimeScale_ApplyValue(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

	scale := NewTimeScale(
		[2]time.Time{start, end},
		[2]units.Length{units.Px(0), units.Px(365)},
	)

	tests := []struct {
		input    time.Time
		expected float64
	}{
		{start, 0.0},
		{end, 1.0},
	}

	for _, tt := range tests {
		result := scale.ApplyValue(tt.input)
		if math.Abs(result-tt.expected) > 0.01 {
			t.Errorf("ApplyValue(%v) = %v, expected %v", tt.input.Format("2006-01-02"), result, tt.expected)
		}
	}

	// Test mid-year
	midYear := time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)
	result := scale.ApplyValue(midYear)
	if result < 0.45 || result > 0.55 { // Should be around 0.5
		t.Errorf("ApplyValue(mid-year) = %v, expected ~0.5", result)
	}
}

func TestTimeScale_Clamp(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

	scale := NewTimeScale(
		[2]time.Time{start, end},
		[2]units.Length{units.Px(0), units.Px(365)},
	)
	scale.Clamp(true)

	tests := []struct {
		input    time.Time
		expected float64
	}{
		{time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), 0},   // Before range
		{start, 0},                                          // Start of range
		{end, 365},                                          // End of range
		{time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), 365}, // After range
	}

	for _, tt := range tests {
		result := scale.Apply(tt.input)
		if math.Abs(result.Value-tt.expected) > 0.01 {
			t.Errorf("Apply(%v) with clamp = %v, expected %v", tt.input.Format("2006-01-02"), result.Value, tt.expected)
		}
	}
}

func TestTimeScale_Nice_Year(t *testing.T) {
	start := time.Date(2024, 3, 15, 10, 30, 0, 0, time.UTC)
	end := time.Date(2026, 9, 20, 14, 45, 0, 0, time.UTC)

	scale := NewTimeScale(
		[2]time.Time{start, end},
		[2]units.Length{units.Px(0), units.Px(500)},
	)
	scale.Nice(TimeIntervalYear)

	domain := scale.Domain().([2]time.Time)

	// Should round to year boundaries
	expectedStart := time.Date(2024, 1, 1, 0, 0, 0, 0, start.Location())
	expectedEnd := time.Date(2027, 1, 1, 0, 0, 0, 0, end.Location())

	if !domain[0].Equal(expectedStart) {
		t.Errorf("Nice(Year) start = %v, expected %v", domain[0].Format(time.RFC3339), expectedStart.Format(time.RFC3339))
	}
	if !domain[1].Equal(expectedEnd) {
		t.Errorf("Nice(Year) end = %v, expected %v", domain[1].Format(time.RFC3339), expectedEnd.Format(time.RFC3339))
	}
}

func TestTimeScale_Nice_Month(t *testing.T) {
	start := time.Date(2024, 3, 15, 10, 30, 0, 0, time.UTC)
	end := time.Date(2024, 9, 20, 14, 45, 0, 0, time.UTC)

	scale := NewTimeScale(
		[2]time.Time{start, end},
		[2]units.Length{units.Px(0), units.Px(500)},
	)
	scale.Nice(TimeIntervalMonth)

	domain := scale.Domain().([2]time.Time)

	// Should round to month boundaries
	expectedStart := time.Date(2024, 3, 1, 0, 0, 0, 0, start.Location())
	expectedEnd := time.Date(2024, 10, 1, 0, 0, 0, 0, end.Location())

	if !domain[0].Equal(expectedStart) {
		t.Errorf("Nice(Month) start = %v, expected %v", domain[0].Format(time.RFC3339), expectedStart.Format(time.RFC3339))
	}
	if !domain[1].Equal(expectedEnd) {
		t.Errorf("Nice(Month) end = %v, expected %v", domain[1].Format(time.RFC3339), expectedEnd.Format(time.RFC3339))
	}
}

func TestTimeScale_Nice_Day(t *testing.T) {
	start := time.Date(2024, 3, 15, 10, 30, 0, 0, time.UTC)
	end := time.Date(2024, 3, 20, 14, 45, 0, 0, time.UTC)

	scale := NewTimeScale(
		[2]time.Time{start, end},
		[2]units.Length{units.Px(0), units.Px(500)},
	)
	scale.Nice(TimeIntervalDay)

	domain := scale.Domain().([2]time.Time)

	// Should round to day boundaries
	expectedStart := time.Date(2024, 3, 15, 0, 0, 0, 0, start.Location())
	expectedEnd := time.Date(2024, 3, 21, 0, 0, 0, 0, end.Location())

	if !domain[0].Equal(expectedStart) {
		t.Errorf("Nice(Day) start = %v, expected %v", domain[0].Format(time.RFC3339), expectedStart.Format(time.RFC3339))
	}
	if !domain[1].Equal(expectedEnd) {
		t.Errorf("Nice(Day) end = %v, expected %v", domain[1].Format(time.RFC3339), expectedEnd.Format(time.RFC3339))
	}
}

func TestTimeScale_Ticks_Year(t *testing.T) {
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)

	scale := NewTimeScale(
		[2]time.Time{start, end},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	ticks := scale.Ticks(10)

	if len(ticks) == 0 {
		t.Fatal("Expected non-empty ticks")
	}

	// All ticks should be at year boundaries
	for i, tick := range ticks {
		if tick.Month() != time.January || tick.Day() != 1 {
			t.Errorf("Tick %d = %v, expected year boundary", i, tick.Format("2006-01-02"))
		}
	}

	// Ticks should be in ascending order
	for i := 1; i < len(ticks); i++ {
		if !ticks[i].After(ticks[i-1]) {
			t.Errorf("Ticks not in ascending order at index %d", i)
		}
	}
}

func TestTimeScale_Ticks_Month(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

	scale := NewTimeScale(
		[2]time.Time{start, end},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	ticks := scale.Ticks(12)

	if len(ticks) == 0 {
		t.Fatal("Expected non-empty ticks")
	}

	// Should have approximately 12 months
	if len(ticks) < 10 || len(ticks) > 14 {
		t.Errorf("Expected ~12 monthly ticks, got %d", len(ticks))
	}

	// All ticks should be at month boundaries
	for i, tick := range ticks {
		if tick.Day() != 1 {
			t.Errorf("Tick %d = %v, expected month boundary", i, tick.Format("2006-01-02"))
		}
	}
}

func TestTimeScale_Ticks_Day(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)

	scale := NewTimeScale(
		[2]time.Time{start, end},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	ticks := scale.Ticks(10)

	if len(ticks) == 0 {
		t.Fatal("Expected non-empty ticks")
	}

	// All ticks should be at day boundaries
	for _, tick := range ticks {
		if tick.Hour() != 0 || tick.Minute() != 0 || tick.Second() != 0 {
			t.Errorf("Tick = %v, expected day boundary", tick.Format(time.RFC3339))
		}
	}
}

func TestTimeScale_Ticks_Hour(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 1, 23, 0, 0, 0, time.UTC)

	scale := NewTimeScale(
		[2]time.Time{start, end},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	ticks := scale.Ticks(12)

	if len(ticks) == 0 {
		t.Fatal("Expected non-empty ticks")
	}

	// All ticks should be at hour boundaries
	for _, tick := range ticks {
		if tick.Minute() != 0 || tick.Second() != 0 {
			t.Errorf("Tick = %v, expected hour boundary", tick.Format(time.RFC3339))
		}
	}
}

func TestTimeScale_Clone(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

	original := NewTimeScale(
		[2]time.Time{start, end},
		[2]units.Length{units.Px(0), units.Px(500)},
	)
	original.Clamp(true)

	clone := original.Clone().(*TimeScale)

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
	newEnd := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
	clone.WithDomain([2]time.Time{start, newEnd})
	if original.domain[1].Year() == 2025 {
		t.Error("Modifying clone affected original")
	}
}

func TestTimeScale_ReverseRange(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

	// Range from 500 to 0 (reversed)
	scale := NewTimeScale(
		[2]time.Time{start, end},
		[2]units.Length{units.Px(500), units.Px(0)},
	)

	tests := []struct {
		input    time.Time
		expected float64
	}{
		{start, 500},
		{end, 0},
	}

	for _, tt := range tests {
		result := scale.Apply(tt.input)
		if math.Abs(result.Value-tt.expected) > 1.0 {
			t.Errorf("Apply(%v) with reverse range = %v, expected %v",
				tt.input.Format("2006-01-02"), result.Value, tt.expected)
		}
	}
}

func TestTimeScale_PointerInput(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

	scale := NewTimeScale(
		[2]time.Time{start, end},
		[2]units.Length{units.Px(0), units.Px(365)},
	)

	// Test with pointer to time.Time
	midYear := time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)
	result := scale.Apply(&midYear)

	if result.Value < 180 || result.Value > 185 { // Roughly mid-year
		t.Errorf("Apply(time pointer) = %v, expected ~182", result.Value)
	}

	// Test with nil pointer
	var nilTime *time.Time
	resultNil := scale.Apply(nilTime)
	if resultNil.Value != 0 {
		t.Errorf("Apply(nil time) = %v, expected 0", resultNil.Value)
	}
}

func TestTimeScale_Type(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

	scale := NewTimeScale(
		[2]time.Time{start, end},
		[2]units.Length{units.Px(0), units.Px(365)},
	)

	if scale.Type() != ScaleTypeTime {
		t.Errorf("Type() = %v, expected ScaleTypeTime", scale.Type())
	}
}

func BenchmarkTimeScale_Apply(b *testing.B) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

	scale := NewTimeScale(
		[2]time.Time{start, end},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	midYear := time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scale.Apply(midYear)
	}
}

func BenchmarkTimeScale_Ticks(b *testing.B) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

	scale := NewTimeScale(
		[2]time.Time{start, end},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scale.Ticks(10)
	}
}
