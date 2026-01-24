package legends

import (
	"strings"
	"testing"

	"github.com/SCKelemen/color"
)

// mustHex is a helper for tests that panics on error
func mustHex(hex string) color.Color {
	c, err := color.HexToRGB(hex)
	if err != nil {
		panic(err)
	}
	return c
}

func TestNewLegend(t *testing.T) {
	items := []LegendItem{
		Item("Series 1", Swatch(mustHex("#3b82f6"))),
		Item("Series 2", Swatch(mustHex("#10b981"))),
	}

	legend := New(items)

	if legend == nil {
		t.Fatal("Expected non-nil legend")
	}

	if len(legend.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(legend.Items))
	}

	if legend.Position != PositionTopRight {
		t.Errorf("Expected default position TopRight, got %d", legend.Position)
	}

	if legend.Style == nil {
		t.Error("Expected non-nil default style")
	}
}

func TestNewLegendWithOptions(t *testing.T) {
	items := []LegendItem{
		Item("Series 1", Swatch(mustHex("#3b82f6"))),
	}

	legend := New(items,
		WithPosition(PositionBottomLeft),
		WithLayout(LayoutHorizontal),
	)

	if legend.Position != PositionBottomLeft {
		t.Errorf("Expected position BottomLeft, got %d", legend.Position)
	}

	if legend.Layout != LayoutHorizontal {
		t.Errorf("Expected layout Horizontal, got %d", legend.Layout)
	}
}

func TestLayoutAuto(t *testing.T) {
	items := []LegendItem{
		Item("Series 1", Swatch(mustHex("#3b82f6"))),
	}

	tests := []struct {
		position       Position
		expectedLayout Layout
	}{
		{PositionTopCenter, LayoutHorizontal},
		{PositionBottomCenter, LayoutHorizontal},
		{PositionTopLeft, LayoutVertical},
		{PositionTopRight, LayoutVertical},
		{PositionBottomLeft, LayoutVertical},
		{PositionBottomRight, LayoutVertical},
		{PositionLeft, LayoutVertical},
		{PositionRight, LayoutVertical},
	}

	for _, tt := range tests {
		legend := New(items, WithPosition(tt.position))
		if legend.Layout != tt.expectedLayout {
			t.Errorf("For position %d, expected layout %d, got %d",
				tt.position, tt.expectedLayout, legend.Layout)
		}
	}
}

func TestColorSwatchSymbol(t *testing.T) {
	swatch := NewColorSwatch(mustHex("#3b82f6"), 15)

	if swatch.Width() != 15 {
		t.Errorf("Expected width 15, got %.1f", swatch.Width())
	}

	if swatch.Height() != 15 {
		t.Errorf("Expected height 15, got %.1f", swatch.Height())
	}

	svg := swatch.Render()
	if !strings.Contains(svg, "<rect") {
		t.Error("Expected SVG to contain <rect>")
	}

	if !strings.Contains(svg, "fill=\"#3b82f6\"") {
		t.Error("Expected SVG to contain color fill")
	}
}

func TestLineSampleSymbol(t *testing.T) {
	line := NewLineSample(mustHex("#10b981"), 2, 25)

	if line.Width() != 25 {
		t.Errorf("Expected width 25, got %.1f", line.Width())
	}

	svg := line.Render()
	if !strings.Contains(svg, "<line") {
		t.Error("Expected SVG to contain <line>")
	}

	if !strings.Contains(svg, "stroke=\"#10b981\"") {
		t.Error("Expected SVG to contain stroke color")
	}
}

func TestLineSampleWithDash(t *testing.T) {
	line := NewLineSample(mustHex("#10b981"), 2, 25).
		WithDash([]float64{4, 2})

	svg := line.Render()
	if !strings.Contains(svg, "stroke-dasharray") {
		t.Error("Expected SVG to contain stroke-dasharray")
	}
}

func TestLineSampleWithMarker(t *testing.T) {
	line := NewLineSample(mustHex("#10b981"), 2, 25).
		WithMarker("circle", 6)

	svg := line.Render()
	if !strings.Contains(svg, "<line") {
		t.Error("Expected SVG to contain <line>")
	}

	if !strings.Contains(svg, "<circle") {
		t.Error("Expected SVG to contain marker <circle>")
	}
}

func TestMarkerSymbols(t *testing.T) {
	markerTypes := []string{
		"circle", "square", "diamond", "triangle",
		"cross", "x", "dot",
	}

	c := mustHex("#ef4444")

	for _, markerType := range markerTypes {
		t.Run(markerType, func(t *testing.T) {
			marker := NewMarkerSymbol(markerType, c, 10)

			if marker.Width() != 10 {
				t.Errorf("Expected width 10, got %.1f", marker.Width())
			}

			svg := marker.Render()
			if svg == "" {
				t.Error("Expected non-empty SVG")
			}

			// Each marker type should produce specific SVG elements
			switch markerType {
			case "circle", "dot":
				if !strings.Contains(svg, "<circle") {
					t.Errorf("Expected <circle> in SVG for %s", markerType)
				}
			case "square":
				if !strings.Contains(svg, "<rect") {
					t.Errorf("Expected <rect> in SVG for %s", markerType)
				}
			case "diamond", "triangle":
				if !strings.Contains(svg, "<polygon") {
					t.Errorf("Expected <polygon> in SVG for %s", markerType)
				}
			case "cross", "x":
				if !strings.Contains(svg, "<path") {
					t.Errorf("Expected <path> in SVG for %s", markerType)
				}
			}
		})
	}
}

func TestLegendRender(t *testing.T) {
	items := []LegendItem{
		Item("Series 1", Swatch(mustHex("#3b82f6"))),
		Item("Series 2", Swatch(mustHex("#10b981"))),
	}

	legend := New(items)
	svg := legend.Render(800, 400)

	if svg == "" {
		t.Fatal("Expected non-empty SVG")
	}

	// Check for legend group
	if !strings.Contains(svg, `<g class="legend"`) {
		t.Error("Expected SVG to contain legend group")
	}

	// Check for items
	if !strings.Contains(svg, "Series 1") {
		t.Error("Expected SVG to contain Series 1 label")
	}

	if !strings.Contains(svg, "Series 2") {
		t.Error("Expected SVG to contain Series 2 label")
	}

	// Check for symbols
	rectCount := strings.Count(svg, "<rect")
	if rectCount < 2 {
		t.Errorf("Expected at least 2 <rect> elements, got %d", rectCount)
	}
}

func TestLegendRenderPositions(t *testing.T) {
	items := []LegendItem{
		Item("Test", Swatch(mustHex("#3b82f6"))),
	}

	positions := []Position{
		PositionTopLeft,
		PositionTopRight,
		PositionTopCenter,
		PositionBottomLeft,
		PositionBottomRight,
		PositionBottomCenter,
		PositionLeft,
		PositionRight,
	}

	for _, pos := range positions {
		t.Run(pos.String(), func(t *testing.T) {
			legend := New(items, WithPosition(pos))
			svg := legend.Render(800, 400)

			if svg == "" {
				t.Errorf("Expected non-empty SVG for position %d", pos)
			}

			// Should contain transform with translation
			if !strings.Contains(svg, "transform=") {
				t.Error("Expected SVG to contain transform")
			}
		})
	}
}

func TestLegendRenderNone(t *testing.T) {
	items := []LegendItem{
		Item("Test", Swatch(mustHex("#3b82f6"))),
	}

	legend := New(items, WithPosition(PositionNone))
	svg := legend.Render(800, 400)

	if svg != "" {
		t.Error("Expected empty SVG for PositionNone")
	}
}

func TestLegendRenderEmpty(t *testing.T) {
	legend := New([]LegendItem{})
	svg := legend.Render(800, 400)

	if svg != "" {
		t.Error("Expected empty SVG for empty items")
	}
}

func TestLegendWithValues(t *testing.T) {
	items := []LegendItem{
		ItemWithValue("Category A", Swatch(mustHex("#3b82f6")), "45%"),
		ItemWithValue("Category B", Swatch(mustHex("#10b981")), "55%"),
	}

	legend := New(items)
	svg := legend.Render(800, 400)

	if !strings.Contains(svg, "Category A (45%)") {
		t.Error("Expected SVG to contain 'Category A (45%)'")
	}

	if !strings.Contains(svg, "Category B (55%)") {
		t.Error("Expected SVG to contain 'Category B (55%)'")
	}
}

func TestLegendHorizontalLayout(t *testing.T) {
	items := []LegendItem{
		Item("A", Swatch(mustHex("#3b82f6"))),
		Item("B", Swatch(mustHex("#10b981"))),
		Item("C", Swatch(mustHex("#ef4444"))),
	}

	legend := New(items, WithLayout(LayoutHorizontal))
	svg := legend.Render(800, 400)

	if svg == "" {
		t.Fatal("Expected non-empty SVG")
	}

	// All three items should be in the SVG
	if !strings.Contains(svg, "A") || !strings.Contains(svg, "B") || !strings.Contains(svg, "C") {
		t.Error("Expected all items in horizontal layout")
	}
}

func TestGetBounds(t *testing.T) {
	items := []LegendItem{
		Item("Series 1", Swatch(mustHex("#3b82f6"))),
		Item("Series 2", Swatch(mustHex("#10b981"))),
	}

	legend := New(items)
	bounds := legend.GetBounds(800, 400)

	if bounds.Width <= 0 {
		t.Error("Expected positive width")
	}

	if bounds.Height <= 0 {
		t.Error("Expected positive height")
	}

	if bounds.X < 0 || bounds.Y < 0 {
		t.Error("Expected non-negative position")
	}
}

func BenchmarkLegendRender(b *testing.B) {
	items := []LegendItem{
		Item("Series 1", Swatch(mustHex("#3b82f6"))),
		Item("Series 2", Swatch(mustHex("#10b981"))),
		Item("Series 3", Swatch(mustHex("#ef4444"))),
		Item("Series 4", Swatch(mustHex("#f59e0b"))),
		Item("Series 5", Swatch(mustHex("#8b5cf6"))),
	}

	legend := New(items)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = legend.Render(800, 400)
	}
}

// Helper method for Position string representation
func (p Position) String() string {
	switch p {
	case PositionTopLeft:
		return "TopLeft"
	case PositionTopRight:
		return "TopRight"
	case PositionTopCenter:
		return "TopCenter"
	case PositionBottomLeft:
		return "BottomLeft"
	case PositionBottomRight:
		return "BottomRight"
	case PositionBottomCenter:
		return "BottomCenter"
	case PositionLeft:
		return "Left"
	case PositionRight:
		return "Right"
	case PositionNone:
		return "None"
	default:
		return "Unknown"
	}
}
