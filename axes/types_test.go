package axes

import (
	"strings"
	"testing"
	"time"

	"github.com/SCKelemen/dataviz/scales"
	"github.com/SCKelemen/units"
)

func TestNewAxis(t *testing.T) {
	scale := scales.NewLinearScale(
		[2]float64{0, 100},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	axis := NewAxis(scale, AxisOrientationBottom)

	if axis == nil {
		t.Fatal("NewAxis returned nil")
	}

	if axis.scale != scale {
		t.Error("Axis scale not set correctly")
	}

	if axis.orientation != AxisOrientationBottom {
		t.Errorf("Axis orientation = %v, expected AxisOrientationBottom", axis.orientation)
	}

	if axis.tickCount != 10 {
		t.Errorf("Default tickCount = %v, expected 10", axis.tickCount)
	}
}

func TestAxisOrientation_IsHorizontal(t *testing.T) {
	tests := []struct {
		orientation AxisOrientation
		expected    bool
	}{
		{AxisOrientationTop, true},
		{AxisOrientationBottom, true},
		{AxisOrientationLeft, false},
		{AxisOrientationRight, false},
	}

	for _, tt := range tests {
		result := tt.orientation.IsHorizontal()
		if result != tt.expected {
			t.Errorf("%s.IsHorizontal() = %v, expected %v", tt.orientation, result, tt.expected)
		}
	}
}

func TestAxisOrientation_IsVertical(t *testing.T) {
	tests := []struct {
		orientation AxisOrientation
		expected    bool
	}{
		{AxisOrientationTop, false},
		{AxisOrientationBottom, false},
		{AxisOrientationLeft, true},
		{AxisOrientationRight, true},
	}

	for _, tt := range tests {
		result := tt.orientation.IsVertical()
		if result != tt.expected {
			t.Errorf("%s.IsVertical() = %v, expected %v", tt.orientation, result, tt.expected)
		}
	}
}

func TestAxis_Ticks_LinearScale(t *testing.T) {
	scale := scales.NewLinearScale(
		[2]float64{0, 100},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	axis := NewAxis(scale, AxisOrientationBottom)
	axis.TickCount(10)

	ticks := axis.Ticks()

	if len(ticks) == 0 {
		t.Fatal("Expected non-empty ticks")
	}

	// Check first and last ticks
	if ticks[0].Value.(float64) < 0 || ticks[0].Value.(float64) > 10 {
		t.Errorf("First tick value = %v, expected near 0", ticks[0].Value)
	}

	if ticks[len(ticks)-1].Value.(float64) < 90 || ticks[len(ticks)-1].Value.(float64) > 100 {
		t.Errorf("Last tick value = %v, expected near 100", ticks[len(ticks)-1].Value)
	}

	// Check positions are increasing
	for i := 1; i < len(ticks); i++ {
		if ticks[i].Position.Value <= ticks[i-1].Position.Value {
			t.Errorf("Tick positions not increasing at index %d", i)
		}
	}

	// Check labels are formatted
	for i, tick := range ticks {
		if tick.Label == "" {
			t.Errorf("Tick %d has empty label", i)
		}
	}
}

func TestAxis_Ticks_BandScale(t *testing.T) {
	scale := scales.NewBandScale(
		[]string{"A", "B", "C", "D"},
		[2]units.Length{units.Px(0), units.Px(400)},
	)

	axis := NewAxis(scale, AxisOrientationBottom)

	ticks := axis.Ticks()

	if len(ticks) != 4 {
		t.Errorf("Expected 4 ticks, got %d", len(ticks))
	}

	// Check all categories present
	expectedLabels := map[string]bool{"A": false, "B": false, "C": false, "D": false}
	for _, tick := range ticks {
		label := tick.Label
		if _, exists := expectedLabels[label]; exists {
			expectedLabels[label] = true
		}
	}

	for label, found := range expectedLabels {
		if !found {
			t.Errorf("Category %q not found in ticks", label)
		}
	}
}

func TestAxis_Ticks_TimeScale(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

	scale := scales.NewTimeScale(
		[2]time.Time{start, end},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	axis := NewAxis(scale, AxisOrientationBottom)
	axis.TickCount(12)

	ticks := axis.Ticks()

	if len(ticks) == 0 {
		t.Fatal("Expected non-empty ticks")
	}

	// Check labels are formatted dates
	for i, tick := range ticks {
		if tick.Label == "" {
			t.Errorf("Tick %d has empty label", i)
		}
		// Should contain a date-like format
		if !strings.Contains(tick.Label, "-") && !strings.Contains(tick.Label, "/") {
			t.Logf("Warning: tick label %q doesn't look like a date", tick.Label)
		}
	}
}

func TestAxis_Title(t *testing.T) {
	scale := scales.NewLinearScale(
		[2]float64{0, 100},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	axis := NewAxis(scale, AxisOrientationBottom)
	axis.Title("Temperature (°C)")

	if axis.title != "Temperature (°C)" {
		t.Errorf("Axis title = %q, expected \"Temperature (°C)\"", axis.title)
	}
}

func TestAxis_TickFormat_Custom(t *testing.T) {
	scale := scales.NewLinearScale(
		[2]float64{0, 100},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	axis := NewAxis(scale, AxisOrientationBottom)

	// Custom formatter adds "°C" suffix
	axis.TickFormat(func(value interface{}) string {
		if v, ok := value.(float64); ok {
			return strings.TrimSuffix(DefaultTickFormatter(v), ".00") + "°C"
		}
		return DefaultTickFormatter(value)
	})

	ticks := axis.Ticks()

	// Check all labels have °C suffix
	for _, tick := range ticks {
		if !strings.Contains(tick.Label, "°C") {
			t.Errorf("Tick label %q doesn't contain °C", tick.Label)
		}
	}
}

func TestDefaultTickFormatter(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected string
	}{
		{42.0, "42"},
		{42.5, "42.50"},
		{42, "42"},
		{"Hello", "Hello"},
		{time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC), "2024-03-15"},
	}

	for _, tt := range tests {
		result := DefaultTickFormatter(tt.input)
		if result != tt.expected {
			t.Errorf("DefaultTickFormatter(%v) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestTimeTickFormatter(t *testing.T) {
	formatter := TimeTickFormatter("Jan 2006")
	testTime := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)

	result := formatter(testTime)
	expected := "Mar 2024"

	if result != expected {
		t.Errorf("TimeTickFormatter(Jan 2006)(%v) = %q, expected %q", testTime, result, expected)
	}
}

func TestNumberTickFormatter(t *testing.T) {
	formatter := NumberTickFormatter(3)

	tests := []struct {
		input    interface{}
		expected string
	}{
		{42.12345, "42.123"},
		{42, "42"},
		{3.14159, "3.142"},
	}

	for _, tt := range tests {
		result := formatter(tt.input)
		if result != tt.expected {
			t.Errorf("NumberTickFormatter(3)(%v) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestSITickFormatter(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected string
	}{
		{0.0, "0"},
		{1.0, "1"},
		{1000.0, "1k"},
		{1500.0, "1.5k"},
		{1000000.0, "1M"},
		{2500000.0, "2.5M"},
		{1000000000.0, "1G"},
		{0.001, "1m"},
		{0.000001, "1μ"},
	}

	for _, tt := range tests {
		result := SITickFormatter(tt.input)
		if result != tt.expected {
			t.Errorf("SITickFormatter(%v) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestAxis_Grid(t *testing.T) {
	scale := scales.NewLinearScale(
		[2]float64{0, 100},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	axis := NewAxis(scale, AxisOrientationBottom)
	axis.Grid(units.Px(300))

	if !axis.showGrid {
		t.Error("Grid not enabled")
	}

	if axis.gridLength.Value != 300 {
		t.Errorf("Grid length = %v, expected 300", axis.gridLength.Value)
	}
}

func TestAxis_TickSize(t *testing.T) {
	scale := scales.NewLinearScale(
		[2]float64{0, 100},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	axis := NewAxis(scale, AxisOrientationBottom)
	axis.TickSize(units.Px(10))

	if axis.tickSize.Value != 10 {
		t.Errorf("Tick size = %v, expected 10", axis.tickSize.Value)
	}
}

func TestAxis_TickPadding(t *testing.T) {
	scale := scales.NewLinearScale(
		[2]float64{0, 100},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	axis := NewAxis(scale, AxisOrientationBottom)
	axis.TickPadding(units.Px(5))

	if axis.tickPadding.Value != 5 {
		t.Errorf("Tick padding = %v, expected 5", axis.tickPadding.Value)
	}
}

func TestAxis_Getters(t *testing.T) {
	scale := scales.NewLinearScale(
		[2]float64{0, 100},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	axis := NewAxis(scale, AxisOrientationLeft)

	if axis.Scale() != scale {
		t.Error("Scale() doesn't return correct scale")
	}

	if axis.Orientation() != AxisOrientationLeft {
		t.Error("Orientation() doesn't return correct orientation")
	}
}
