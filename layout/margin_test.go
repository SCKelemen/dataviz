package layout

import (
	"testing"

	"github.com/SCKelemen/units"
)

func TestMarginConvention(t *testing.T) {
	mc := NewMarginConvention(units.Px(800), units.Px(600))

	totalBounds := mc.TotalBounds()
	if totalBounds.Width.Value != 800 {
		t.Errorf("Expected total width 800, got %f", totalBounds.Width.Value)
	}
	if totalBounds.Height.Value != 600 {
		t.Errorf("Expected total height 600, got %f", totalBounds.Height.Value)
	}

	plotArea := mc.PlotArea()
	// With default margin, plot area should be smaller
	if plotArea.Width.Value >= 800 {
		t.Error("Plot area should be smaller than total bounds")
	}
	if plotArea.Height.Value >= 600 {
		t.Error("Plot area should be smaller than total bounds")
	}
}

func TestMarginConvention_SetMargin(t *testing.T) {
	mc := NewMarginConvention(units.Px(800), units.Px(600))
	mc.SetMargin(Margin{
		Top:    units.Px(20),
		Right:  units.Px(30),
		Bottom: units.Px(40),
		Left:   units.Px(50),
	})

	plotWidth := mc.PlotWidth()
	plotHeight := mc.PlotHeight()

	// 800 - 50 - 30 = 720
	if plotWidth.Value != 720 {
		t.Errorf("Expected plot width 720, got %f", plotWidth.Value)
	}

	// 600 - 20 - 40 = 540
	if plotHeight.Value != 540 {
		t.Errorf("Expected plot height 540, got %f", plotHeight.Value)
	}
}

func TestMarginConvention_MarginAreas(t *testing.T) {
	mc := NewMarginConvention(units.Px(800), units.Px(600))
	mc.SetMargin(Margin{
		Top:    units.Px(20),
		Right:  units.Px(30),
		Bottom: units.Px(40),
		Left:   units.Px(50),
	})

	leftArea := mc.LeftMarginArea()
	if leftArea.Width.Value != 50 {
		t.Errorf("Expected left margin width 50, got %f", leftArea.Width.Value)
	}

	topArea := mc.TopMarginArea()
	if topArea.Height.Value != 20 {
		t.Errorf("Expected top margin height 20, got %f", topArea.Height.Value)
	}

	bottomArea := mc.BottomMarginArea()
	if bottomArea.Height.Value != 40 {
		t.Errorf("Expected bottom margin height 40, got %f", bottomArea.Height.Value)
	}

	rightArea := mc.RightMarginArea()
	if rightArea.Width.Value != 30 {
		t.Errorf("Expected right margin width 30, got %f", rightArea.Width.Value)
	}
}

func TestComputeMarginForAxes(t *testing.T) {
	// Test with all axes
	margin := ComputeMarginForAxes(true, true, true, true, true)
	if margin.Left.Value < 50 {
		t.Error("Left margin should be substantial for Y axis")
	}
	if margin.Bottom.Value < 40 {
		t.Error("Bottom margin should be substantial for X axis")
	}
	if margin.Top.Value < 30 {
		t.Error("Top margin should be substantial for title")
	}

	// Test with no axes
	marginMin := ComputeMarginForAxes(false, false, false, false, false)
	if marginMin.Left.Value > 15 {
		t.Error("Left margin should be minimal without Y axis")
	}
	if marginMin.Bottom.Value > 15 {
		t.Error("Bottom margin should be minimal without X axis")
	}
}

func TestSplitHorizontal(t *testing.T) {
	bounds := Rect{
		X:      units.Px(0),
		Y:      units.Px(0),
		Width:  units.Px(800),
		Height: units.Px(600),
	}

	left, right := SplitHorizontal(bounds, 0.5)

	if left.Width.Value != 400 {
		t.Errorf("Expected left width 400, got %f", left.Width.Value)
	}
	if right.Width.Value != 400 {
		t.Errorf("Expected right width 400, got %f", right.Width.Value)
	}
	if right.X.Value != 400 {
		t.Errorf("Expected right X at 400, got %f", right.X.Value)
	}
}

func TestSplitVertical(t *testing.T) {
	bounds := Rect{
		X:      units.Px(0),
		Y:      units.Px(0),
		Width:  units.Px(800),
		Height: units.Px(600),
	}

	top, bottom := SplitVertical(bounds, 0.3)

	if top.Height.Value != 180 {
		t.Errorf("Expected top height 180, got %f", top.Height.Value)
	}
	if bottom.Height.Value != 420 {
		t.Errorf("Expected bottom height 420, got %f", bottom.Height.Value)
	}
	if bottom.Y.Value != 180 {
		t.Errorf("Expected bottom Y at 180, got %f", bottom.Y.Value)
	}
}

func TestSplitIntoGrid(t *testing.T) {
	bounds := Rect{
		X:      units.Px(0),
		Y:      units.Px(0),
		Width:  units.Px(900),
		Height: units.Px(600),
	}

	grid := SplitIntoGrid(bounds, 2, 3, units.Px(10))

	if len(grid) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(grid))
	}
	if len(grid[0]) != 3 {
		t.Errorf("Expected 3 cols, got %d", len(grid[0]))
	}

	// Check gap between cells
	cell1 := grid[0][0]
	cell2 := grid[0][1]
	gap := cell2.X.Value - (cell1.X.Value + cell1.Width.Value)
	if gap != 10 {
		t.Errorf("Expected gap of 10, got %f", gap)
	}
}

func TestInset(t *testing.T) {
	bounds := Rect{
		X:      units.Px(0),
		Y:      units.Px(0),
		Width:  units.Px(800),
		Height: units.Px(600),
	}

	inset := Inset(bounds, units.Px(10))

	if inset.X.Value != 10 {
		t.Errorf("Expected X at 10, got %f", inset.X.Value)
	}
	if inset.Y.Value != 10 {
		t.Errorf("Expected Y at 10, got %f", inset.Y.Value)
	}
	if inset.Width.Value != 780 {
		t.Errorf("Expected width 780, got %f", inset.Width.Value)
	}
	if inset.Height.Value != 580 {
		t.Errorf("Expected height 580, got %f", inset.Height.Value)
	}
}

func TestCenter(t *testing.T) {
	bounds := Rect{
		X:      units.Px(100),
		Y:      units.Px(200),
		Width:  units.Px(800),
		Height: units.Px(600),
	}

	x, y := Center(bounds)

	if x.Value != 500 {
		t.Errorf("Expected center X at 500, got %f", x.Value)
	}
	if y.Value != 500 {
		t.Errorf("Expected center Y at 500, got %f", y.Value)
	}
}

func TestContains(t *testing.T) {
	bounds := Rect{
		X:      units.Px(100),
		Y:      units.Px(100),
		Width:  units.Px(200),
		Height: units.Px(200),
	}

	// Point inside
	if !Contains(bounds, units.Px(150), units.Px(150)) {
		t.Error("Point (150, 150) should be inside bounds")
	}

	// Point outside
	if Contains(bounds, units.Px(50), units.Px(50)) {
		t.Error("Point (50, 50) should be outside bounds")
	}

	// Point on edge
	if !Contains(bounds, units.Px(100), units.Px(100)) {
		t.Error("Point (100, 100) should be on edge (inclusive)")
	}
}

func TestApplyMargin(t *testing.T) {
	bounds := Rect{
		X:      units.Px(0),
		Y:      units.Px(0),
		Width:  units.Px(800),
		Height: units.Px(600),
	}

	margin := Margin{
		Top:    units.Px(10),
		Right:  units.Px(20),
		Bottom: units.Px(30),
		Left:   units.Px(40),
	}

	result := ApplyMargin(bounds, margin)

	if result.X.Value != 40 {
		t.Errorf("Expected X at 40, got %f", result.X.Value)
	}
	if result.Y.Value != 10 {
		t.Errorf("Expected Y at 10, got %f", result.Y.Value)
	}
	if result.Width.Value != 740 {
		t.Errorf("Expected width 740, got %f", result.Width.Value)
	}
	if result.Height.Value != 560 {
		t.Errorf("Expected height 560, got %f", result.Height.Value)
	}
}
