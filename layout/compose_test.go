package layout

import (
	"strings"
	"testing"

	"github.com/SCKelemen/units"
)

func TestNewComposition(t *testing.T) {
	comp := NewComposition(units.Px(800), units.Px(600))

	if comp.Width.Value != 800 {
		t.Errorf("Expected width 800, got %f", comp.Width.Value)
	}
	if comp.Height.Value != 600 {
		t.Errorf("Expected height 600, got %f", comp.Height.Value)
	}
	if comp.Layout != LayoutGrid {
		t.Error("Expected default layout to be grid")
	}
}

func TestComposition_AddChart(t *testing.T) {
	comp := NewComposition(units.Px(800), units.Px(600))

	renderer := func(bounds Rect) string {
		return "<rect/>"
	}

	chart := NewChartSpec(renderer)
	comp.AddChart(chart)

	if len(comp.Charts) != 1 {
		t.Errorf("Expected 1 chart, got %d", len(comp.Charts))
	}
}

func TestComposition_Render_Grid(t *testing.T) {
	comp := NewComposition(units.Px(800), units.Px(600)).
		WithLayout(LayoutGrid).
		WithGrid(2, 2)

	for i := 0; i < 4; i++ {
		renderer := func(bounds Rect) string {
			return `<rect class="chart"/>`
		}
		comp.AddChart(NewChartSpec(renderer))
	}

	svg := comp.Render()

	if !strings.Contains(svg, "<svg") {
		t.Error("Should contain SVG element")
	}
	if !strings.Contains(svg, "viewBox") {
		t.Error("Should contain viewBox")
	}
}

func TestComposition_Render_Stack(t *testing.T) {
	comp := NewComposition(units.Px(800), units.Px(600)).
		WithLayout(LayoutStack)

	for i := 0; i < 3; i++ {
		renderer := func(bounds Rect) string {
			return `<rect class="chart"/>`
		}
		comp.AddChart(NewChartSpec(renderer))
	}

	svg := comp.Render()

	if !strings.Contains(svg, "<svg") {
		t.Error("Should contain SVG element")
	}
}

func TestComposition_WithTitle(t *testing.T) {
	comp := NewComposition(units.Px(800), units.Px(600)).
		WithTitle("Dashboard")

	renderer := func(bounds Rect) string {
		return "<rect/>"
	}
	comp.AddChart(NewChartSpec(renderer))

	svg := comp.Render()

	if !strings.Contains(svg, "Dashboard") {
		t.Error("Should contain title")
	}
}

func TestChartSpec_WithGridPosition(t *testing.T) {
	renderer := func(bounds Rect) string {
		return "<rect/>"
	}

	chart := NewChartSpec(renderer).
		WithGridPosition(1, 2).
		WithSpan(2, 1).
		WithTitle("Test Chart")

	if chart.Row != 1 {
		t.Errorf("Expected row 1, got %d", chart.Row)
	}
	if chart.Col != 2 {
		t.Errorf("Expected col 2, got %d", chart.Col)
	}
	if chart.RowSpan != 2 {
		t.Errorf("Expected rowSpan 2, got %d", chart.RowSpan)
	}
	if chart.ColSpan != 1 {
		t.Errorf("Expected colSpan 1, got %d", chart.ColSpan)
	}
	if chart.Title != "Test Chart" {
		t.Errorf("Expected title 'Test Chart', got '%s'", chart.Title)
	}
}

func TestGridComposition(t *testing.T) {
	renderer1 := func(bounds Rect) string { return "<rect/>" }
	renderer2 := func(bounds Rect) string { return "<circle/>" }

	comp := GridComposition(units.Px(800), units.Px(600), 1, 2, renderer1, renderer2)

	if len(comp.Charts) != 2 {
		t.Errorf("Expected 2 charts, got %d", len(comp.Charts))
	}
	if comp.Layout != LayoutGrid {
		t.Error("Expected grid layout")
	}
	if comp.Rows != 1 || comp.Cols != 2 {
		t.Errorf("Expected 1x2 grid, got %dx%d", comp.Rows, comp.Cols)
	}
}

func TestStackComposition(t *testing.T) {
	renderer1 := func(bounds Rect) string { return "<rect/>" }
	renderer2 := func(bounds Rect) string { return "<circle/>" }
	renderer3 := func(bounds Rect) string { return "<line/>" }

	comp := StackComposition(units.Px(800), units.Px(600), renderer1, renderer2, renderer3)

	if len(comp.Charts) != 3 {
		t.Errorf("Expected 3 charts, got %d", len(comp.Charts))
	}
	if comp.Layout != LayoutStack {
		t.Error("Expected stack layout")
	}
}

func TestSideBySide(t *testing.T) {
	left := func(bounds Rect) string { return "<rect/>" }
	right := func(bounds Rect) string { return "<circle/>" }

	comp := SideBySide(units.Px(800), units.Px(600), left, right)

	if len(comp.Charts) != 2 {
		t.Errorf("Expected 2 charts, got %d", len(comp.Charts))
	}
	if comp.Rows != 1 || comp.Cols != 2 {
		t.Errorf("Expected 1x2 grid, got %dx%d", comp.Rows, comp.Cols)
	}
}

func TestTopAndBottom(t *testing.T) {
	top := func(bounds Rect) string { return "<rect/>" }
	bottom := func(bounds Rect) string { return "<circle/>" }

	comp := TopAndBottom(units.Px(800), units.Px(600), top, bottom)

	if len(comp.Charts) != 2 {
		t.Errorf("Expected 2 charts, got %d", len(comp.Charts))
	}
	if comp.Rows != 2 || comp.Cols != 1 {
		t.Errorf("Expected 2x1 grid, got %dx%d", comp.Rows, comp.Cols)
	}
}

func TestQuad(t *testing.T) {
	tl := func(bounds Rect) string { return "<rect/>" }
	tr := func(bounds Rect) string { return "<circle/>" }
	bl := func(bounds Rect) string { return "<line/>" }
	br := func(bounds Rect) string { return "<path/>" }

	comp := Quad(units.Px(800), units.Px(600), tl, tr, bl, br)

	if len(comp.Charts) != 4 {
		t.Errorf("Expected 4 charts, got %d", len(comp.Charts))
	}
	if comp.Rows != 2 || comp.Cols != 2 {
		t.Errorf("Expected 2x2 grid, got %dx%d", comp.Rows, comp.Cols)
	}
}

func TestDashboardComposition(t *testing.T) {
	chart1 := NewChartSpec(func(bounds Rect) string { return "<rect/>" }).
		WithGridPosition(0, 0).
		WithSpan(1, 2)

	chart2 := NewChartSpec(func(bounds Rect) string { return "<circle/>" }).
		WithGridPosition(1, 0)

	chart3 := NewChartSpec(func(bounds Rect) string { return "<line/>" }).
		WithGridPosition(1, 1)

	comp := DashboardComposition(units.Px(800), units.Px(600), chart1, chart2, chart3)

	if len(comp.Charts) != 3 {
		t.Errorf("Expected 3 charts, got %d", len(comp.Charts))
	}
	if comp.Layout != LayoutDashboard {
		t.Error("Expected dashboard layout")
	}
}

func TestCustomComposition(t *testing.T) {
	chart1 := NewChartSpec(func(bounds Rect) string { return "<rect/>" }).
		WithBounds(Rect{
			X:      units.Px(0),
			Y:      units.Px(0),
			Width:  units.Px(400),
			Height: units.Px(300),
		})

	chart2 := NewChartSpec(func(bounds Rect) string { return "<circle/>" }).
		WithBounds(Rect{
			X:      units.Px(400),
			Y:      units.Px(0),
			Width:  units.Px(400),
			Height: units.Px(300),
		})

	comp := CustomComposition(units.Px(800), units.Px(600), chart1, chart2)

	if len(comp.Charts) != 2 {
		t.Errorf("Expected 2 charts, got %d", len(comp.Charts))
	}
	if comp.Layout != LayoutCustom {
		t.Error("Expected custom layout")
	}
}

func TestComposition_WithBackground(t *testing.T) {
	comp := NewComposition(units.Px(800), units.Px(600)).
		WithBackground("#f0f0f0").
		WithBorder(true)

	renderer := func(bounds Rect) string {
		return "<rect/>"
	}
	comp.AddChart(NewChartSpec(renderer))

	svg := comp.Render()

	if !strings.Contains(svg, "#f0f0f0") {
		t.Error("Should contain background color")
	}
}

func TestChartGroup(t *testing.T) {
	group := NewChartGroup("Test Group")

	chart1 := NewChartSpec(func(bounds Rect) string { return "<rect/>" })
	chart2 := NewChartSpec(func(bounds Rect) string { return "<circle/>" })

	group.AddChart(chart1).AddChart(chart2)

	if len(group.Charts) != 2 {
		t.Errorf("Expected 2 charts, got %d", len(group.Charts))
	}
	if group.Title != "Test Group" {
		t.Errorf("Expected title 'Test Group', got '%s'", group.Title)
	}
}

func TestComposition_ChainedMethods(t *testing.T) {
	comp := NewComposition(units.Px(800), units.Px(600)).
		WithLayout(LayoutStack).
		WithGap(units.Px(15)).
		WithMargin(Uniform(units.Px(20))).
		WithTitle("My Dashboard").
		WithBackground("white").
		WithBorder(true)

	if comp.Layout != LayoutStack {
		t.Error("Layout should be stack")
	}
	if comp.Gap.Value != 15 {
		t.Error("Gap should be 15")
	}
	if comp.Margin.Top.Value != 20 {
		t.Error("Margin should be 20")
	}
	if comp.Title != "My Dashboard" {
		t.Error("Title should be set")
	}
	if comp.Border != true {
		t.Error("Border should be enabled")
	}
}
