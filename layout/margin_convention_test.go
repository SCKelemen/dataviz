package layout

import (
	"testing"
)

func TestNewMarginConvention(t *testing.T) {
	mc := NewMarginConvention(800, 600)

	if mc == nil {
		t.Fatal("NewMarginConvention should return non-nil")
	}

	// Check that it has non-zero plot dimensions
	plotWidth := mc.PlotWidth()
	if plotWidth <= 0 || plotWidth >= 800 {
		t.Error("Plot width should be positive and less than total width")
	}

	plotHeight := mc.PlotHeight()
	if plotHeight <= 0 || plotHeight >= 600 {
		t.Error("Plot height should be positive and less than total height")
	}
}

func TestSetMargin(t *testing.T) {
	mc := NewMarginConvention(800, 600)
	mc.SetMargin(10, 20, 30, 40)

	// Verify via plot dimensions
	plotWidth := mc.PlotWidth()
	expectedWidth := 800.0 - 20 - 40
	if plotWidth != expectedWidth {
		t.Errorf("Expected plot width %f, got %f", expectedWidth, plotWidth)
	}

	plotHeight := mc.PlotHeight()
	expectedHeight := 600.0 - 10 - 30
	if plotHeight != expectedHeight {
		t.Errorf("Expected plot height %f, got %f", expectedHeight, plotHeight)
	}
}

func TestSetUniformMargin(t *testing.T) {
	mc := NewMarginConvention(800, 600)
	mc.SetUniformMargin(15)

	plotWidth := mc.PlotWidth()
	expectedWidth := 800.0 - 15 - 15
	if plotWidth != expectedWidth {
		t.Errorf("Expected plot width %f, got %f", expectedWidth, plotWidth)
	}

	plotHeight := mc.PlotHeight()
	expectedHeight := 600.0 - 15 - 15
	if plotHeight != expectedHeight {
		t.Errorf("Expected plot height %f, got %f", expectedHeight, plotHeight)
	}
}

func TestPlotWidth(t *testing.T) {
	mc := NewMarginConvention(800, 600)
	mc.SetMargin(10, 20, 30, 40)

	plotWidth := mc.PlotWidth()
	expected := 800.0 - 20 - 40 // width - right - left

	if plotWidth != expected {
		t.Errorf("Expected plot width %f, got %f", expected, plotWidth)
	}
}

func TestPlotHeight(t *testing.T) {
	mc := NewMarginConvention(800, 600)
	mc.SetMargin(10, 20, 30, 40)

	plotHeight := mc.PlotHeight()
	expected := 600.0 - 10 - 30 // height - top - bottom

	if plotHeight != expected {
		t.Errorf("Expected plot height %f, got %f", expected, plotHeight)
	}
}

func TestPlotArea(t *testing.T) {
	mc := NewMarginConvention(800, 600)
	mc.SetMargin(10, 20, 30, 40)

	rect := mc.PlotArea()

	// Plot area should start at (left, top)
	if rect.X != 40 {
		t.Errorf("Expected plot area X = 40, got %f", rect.X)
	}
	if rect.Y != 10 {
		t.Errorf("Expected plot area Y = 10, got %f", rect.Y)
	}

	// Plot area dimensions
	expectedWidth := 800.0 - 20 - 40
	expectedHeight := 600.0 - 10 - 30

	if rect.Width != expectedWidth {
		t.Errorf("Expected plot area width %f, got %f", expectedWidth, rect.Width)
	}
	if rect.Height != expectedHeight {
		t.Errorf("Expected plot area height %f, got %f", expectedHeight, rect.Height)
	}
}

func TestTotalBounds(t *testing.T) {
	mc := NewMarginConvention(800, 600)
	mc.SetMargin(10, 20, 30, 40)

	rect := mc.TotalBounds()

	if rect.X != 0 || rect.Y != 0 {
		t.Error("Total bounds should start at origin")
	}
	if rect.Width != 800 || rect.Height != 600 {
		t.Error("Total bounds should match full dimensions")
	}
}

func TestLeftMarginArea(t *testing.T) {
	mc := NewMarginConvention(800, 600)
	mc.SetMargin(10, 20, 30, 40)

	rect := mc.LeftMarginArea()

	if rect.X != 0 {
		t.Errorf("Expected left margin X = 0, got %f", rect.X)
	}
	if rect.Y != 10 {
		t.Errorf("Expected left margin Y = 10, got %f", rect.Y)
	}
	if rect.Width != 40 {
		t.Errorf("Expected left margin width = 40, got %f", rect.Width)
	}

	expectedHeight := 600.0 - 10 - 30
	if rect.Height != expectedHeight {
		t.Errorf("Expected left margin height %f, got %f", expectedHeight, rect.Height)
	}
}

func TestRightMarginArea(t *testing.T) {
	mc := NewMarginConvention(800, 600)
	mc.SetMargin(10, 20, 30, 40)

	rect := mc.RightMarginArea()

	expectedX := 800.0 - 20
	if rect.X != expectedX {
		t.Errorf("Expected right margin X = %f, got %f", expectedX, rect.X)
	}
	if rect.Y != 10 {
		t.Errorf("Expected right margin Y = 10, got %f", rect.Y)
	}
	if rect.Width != 20 {
		t.Errorf("Expected right margin width = 20, got %f", rect.Width)
	}

	expectedHeight := 600.0 - 10 - 30
	if rect.Height != expectedHeight {
		t.Errorf("Expected right margin height %f, got %f", expectedHeight, rect.Height)
	}
}



func TestAsNode(t *testing.T) {
	mc := NewMarginConvention(800, 600)
	mc.SetMargin(10, 20, 30, 40)

	node := mc.AsNode()

	if node == nil {
		t.Fatal("AsNode should return a non-nil node")
	}

	if node.Style.Width.Value != 800 {
		t.Errorf("Expected node width 800, got %f", node.Style.Width.Value)
	}
	if node.Style.Height.Value != 600 {
		t.Errorf("Expected node height 600, got %f", node.Style.Height.Value)
	}

	// Check padding matches margins
	if node.Style.Padding.Top.Value != 10 {
		t.Errorf("Expected padding top 10, got %f", node.Style.Padding.Top.Value)
	}
	if node.Style.Padding.Right.Value != 20 {
		t.Errorf("Expected padding right 20, got %f", node.Style.Padding.Right.Value)
	}
	if node.Style.Padding.Bottom.Value != 30 {
		t.Errorf("Expected padding bottom 30, got %f", node.Style.Padding.Bottom.Value)
	}
	if node.Style.Padding.Left.Value != 40 {
		t.Errorf("Expected padding left 40, got %f", node.Style.Padding.Left.Value)
	}
}

func TestComputeMarginForAxes(t *testing.T) {
	// Test with all axes
	top, right, bottom, left := ComputeMarginForAxes(true, true, true, true, true)

	// Should have non-zero margins for all sides with axes
	if top == 0 {
		t.Error("Should have top margin for top axis")
	}
	if right == 0 {
		t.Error("Should have right margin for right axis")
	}
	if bottom == 0 {
		t.Error("Should have bottom margin for bottom axis")
	}
	if left == 0 {
		t.Error("Should have left margin for left axis")
	}

	// Test with only horizontal axes (left and bottom)
	top2, right2, bottom2, left2 := ComputeMarginForAxes(true, false, false, true, false)

	if top2 != 10 {
		t.Error("Should have minimal top margin when no top axis or title")
	}
	if right2 != 10 {
		t.Error("Should have minimal right margin when no right axis")
	}
	if bottom2 == 10 {
		t.Error("Should have larger bottom margin for bottom axis")
	}
	if left2 == 10 {
		t.Error("Should have larger left margin for left axis")
	}
}

func TestDefaultChartMargin(t *testing.T) {
	top, right, bottom, left := DefaultChartMargin()

	// Should have reasonable default margins
	if top == 0 || bottom == 0 || left == 0 {
		t.Error("Default chart margin should have non-zero margins for typical axes")
	}

	// Should have smaller right margin (no right axis by default)
	if right >= left {
		t.Error("Right margin should be smaller than left margin by default")
	}
}

func TestMarginConventionChaining(t *testing.T) {
	mc := NewMarginConvention(800, 600)
	mc.SetMargin(10, 20, 30, 40)

	// Test that methods return the convention for chaining
	result := mc.SetUniformMargin(15)
	if result == nil {
		t.Error("Chained SetUniformMargin should return convention")
	}

	// Verify the change took effect
	plotWidth := result.PlotWidth()
	expectedWidth := 800.0 - 15 - 15
	if plotWidth != expectedWidth {
		t.Error("Chained method should update margins")
	}

	result2 := mc.SetMargin(5, 10, 15, 20)
	if result2 == nil {
		t.Error("Chained SetMargin should return convention")
	}
}


func TestMarginConventionWithNegativeMargins(t *testing.T) {
	mc := NewMarginConvention(800, 600)
	mc.SetMargin(-10, -20, -30, -40)

	// Negative margins should still work (though unusual)
	plotWidth := mc.PlotWidth()
	expectedWidth := 800.0 - (-20) - (-40) // Should be larger than total width

	if plotWidth != expectedWidth {
		t.Errorf("Expected plot width %f, got %f", expectedWidth, plotWidth)
	}
}
