package axes

import (
	"strings"
	"testing"

	"github.com/SCKelemen/dataviz/scales"
	"github.com/SCKelemen/units"
)

func TestAxis_Render_Bottom(t *testing.T) {
	scale := scales.NewLinearScale(
		[2]float64{0, 100},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	axis := NewAxis(scale, AxisOrientationBottom)
	axis.Title("X Axis")

	opts := DefaultRenderOptions()
	opts.Position = units.Px(300)

	svg := axis.Render(opts)

	if svg == "" {
		t.Fatal("Render returned empty string")
	}

	if !strings.Contains(svg, "axis-bottom") {
		t.Error("SVG missing axis-bottom class")
	}

	// Check for expected elements
	if !strings.Contains(svg, "<line") {
		t.Error("SVG missing line elements")
	}

	if !strings.Contains(svg, "<text") {
		t.Error("SVG missing text elements")
	}

	if !strings.Contains(svg, "X Axis") {
		t.Error("SVG missing title")
	}
}

func TestAxis_Render_Top(t *testing.T) {
	scale := scales.NewLinearScale(
		[2]float64{0, 100},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	axis := NewAxis(scale, AxisOrientationTop)

	opts := DefaultRenderOptions()
	opts.Position = units.Px(50)

	svg := axis.Render(opts)

	if svg == "" {
		t.Fatal("Render returned empty string")
	}

	if !strings.Contains(svg, "axis-top") {
		t.Error("SVG missing axis-top class")
	}
}

func TestAxis_Render_Left(t *testing.T) {
	scale := scales.NewLinearScale(
		[2]float64{0, 100},
		[2]units.Length{units.Px(400), units.Px(0)}, // Inverted for Y axis
	)

	axis := NewAxis(scale, AxisOrientationLeft)
	axis.Title("Y Axis")

	opts := DefaultRenderOptions()
	opts.Position = units.Px(50)

	svg := axis.Render(opts)

	if svg == "" {
		t.Fatal("Render returned empty string")
	}

	if !strings.Contains(svg, "axis-left") {
		t.Error("SVG missing axis-left class")
	}

	if !strings.Contains(svg, "Y Axis") {
		t.Error("SVG missing title")
	}

	// Should have transform for rotated title
	if !strings.Contains(svg, "rotate") {
		t.Error("SVG missing rotate transform for vertical axis title")
	}
}

func TestAxis_Render_Right(t *testing.T) {
	scale := scales.NewLinearScale(
		[2]float64{0, 100},
		[2]units.Length{units.Px(400), units.Px(0)}, // Inverted for Y axis
	)

	axis := NewAxis(scale, AxisOrientationRight)

	opts := DefaultRenderOptions()
	opts.Position = units.Px(550)

	svg := axis.Render(opts)

	if svg == "" {
		t.Fatal("Render returned empty string")
	}

	if !strings.Contains(svg, "axis-right") {
		t.Error("SVG missing axis-right class")
	}
}

func TestAxis_Render_WithGrid(t *testing.T) {
	scale := scales.NewLinearScale(
		[2]float64{0, 100},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	axis := NewAxis(scale, AxisOrientationBottom)
	axis.Grid(units.Px(300))

	opts := DefaultRenderOptions()
	opts.Position = units.Px(300)

	svg := axis.Render(opts)

	// Count line elements (should have axis line + ticks + grid lines)
	lineCount := strings.Count(svg, "<line")

	// With grid enabled, should have more lines than just axis + ticks
	if lineCount < 10 {
		t.Errorf("Expected at least 10 lines with grid enabled, got %d", lineCount)
	}
}

func TestAxis_String(t *testing.T) {
	scale := scales.NewLinearScale(
		[2]float64{0, 100},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	axis := NewAxis(scale, AxisOrientationBottom)
	axis.Title("Test Axis")

	opts := DefaultRenderOptions()
	opts.Position = units.Px(300)

	svg := axis.String(opts)

	if svg == "" {
		t.Fatal("String() returned empty string")
	}

	// Check for expected SVG elements
	expectedElements := []string{
		"<g",      // Group element
		"</g>",    // Closing group
		"<line",   // Axis line and ticks
		"<text",   // Labels
	}

	for _, expected := range expectedElements {
		if !strings.Contains(svg, expected) {
			t.Errorf("SVG missing expected element: %q", expected)
		}
	}

	// Check for title
	if !strings.Contains(svg, "Test Axis") {
		t.Error("SVG missing axis title")
	}
}

func TestAxis_Render_BandScale(t *testing.T) {
	scale := scales.NewBandScale(
		[]string{"Mon", "Tue", "Wed", "Thu", "Fri"},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	axis := NewAxis(scale, AxisOrientationBottom)

	opts := DefaultRenderOptions()
	opts.Position = units.Px(300)

	svg := axis.Render(opts)

	if svg == "" {
		t.Fatal("Render returned empty string")
	}

	// Should have labels for all 5 categories
	for _, day := range []string{"Mon", "Tue", "Wed", "Thu", "Fri"} {
		if !strings.Contains(svg, day) {
			t.Errorf("SVG missing label for %s", day)
		}
	}
}

func TestDefaultAxisStyle(t *testing.T) {
	style := DefaultAxisStyle()

	if style.StrokeColor == "" {
		t.Error("Default stroke color is empty")
	}

	if style.FontSize <= 0 {
		t.Error("Default font size is invalid")
	}

	if style.FontFamily == "" {
		t.Error("Default font family is empty")
	}
}

func TestDefaultRenderOptions(t *testing.T) {
	opts := DefaultRenderOptions()

	if opts.Style.StrokeColor == "" {
		t.Error("Default options have empty stroke color")
	}

	if opts.Position.Value != 0 {
		t.Errorf("Default position = %v, expected 0", opts.Position.Value)
	}
}

func TestAxis_Render_CustomStyle(t *testing.T) {
	scale := scales.NewLinearScale(
		[2]float64{0, 100},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	axis := NewAxis(scale, AxisOrientationBottom)

	opts := DefaultRenderOptions()
	opts.Style.StrokeColor = "#ff0000"
	opts.Style.TextColor = "#0000ff"
	opts.Style.FontSize = 14
	opts.Position = units.Px(300)

	svg := axis.Render(opts)

	if svg == "" {
		t.Error("Rendering with custom styles produced empty output")
	}

	// Check that custom colors appear in SVG
	if !strings.Contains(svg, "#ff0000") {
		t.Error("Custom stroke color not found in SVG")
	}

	if !strings.Contains(svg, "#0000ff") {
		t.Error("Custom text color not found in SVG")
	}
}

func TestAxis_Render_EmptyTicks(t *testing.T) {
	// Create a scale with empty domain
	scale := scales.NewBandScale(
		[]string{},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	axis := NewAxis(scale, AxisOrientationBottom)

	opts := DefaultRenderOptions()
	opts.Position = units.Px(300)

	svg := axis.Render(opts)

	// Should return empty string with empty ticks
	if svg != "" {
		t.Error("Expected empty string with empty ticks")
	}
}

func BenchmarkAxis_Render(b *testing.B) {
	scale := scales.NewLinearScale(
		[2]float64{0, 100},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	axis := NewAxis(scale, AxisOrientationBottom)
	axis.Title("Benchmark Axis")

	opts := DefaultRenderOptions()
	opts.Position = units.Px(300)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		axis.Render(opts)
	}
}

func BenchmarkAxis_String(b *testing.B) {
	scale := scales.NewLinearScale(
		[2]float64{0, 100},
		[2]units.Length{units.Px(0), units.Px(500)},
	)

	axis := NewAxis(scale, AxisOrientationBottom)

	opts := DefaultRenderOptions()
	opts.Position = units.Px(300)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		axis.String(opts)
	}
}
