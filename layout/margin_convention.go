package layout

import (
	"github.com/SCKelemen/layout"
)

// MarginConvention implements the D3 margin convention for charts
// This pattern reserves space for axes, titles, and labels around a plot area
//
// Example:
//   mc := NewMarginConvention(800, 600)
//   mc.SetMargin(60, 20, 50, 60) // top, right, bottom, left
//   plotArea := mc.PlotArea() // Get inner drawing area
type MarginConvention struct {
	totalWidth  float64
	totalHeight float64
	marginTop   float64
	marginRight float64
	marginBottom float64
	marginLeft  float64
}

// NewMarginConvention creates a new margin convention
func NewMarginConvention(width, height float64) *MarginConvention {
	return &MarginConvention{
		totalWidth:   width,
		totalHeight:  height,
		marginTop:    40,
		marginRight:  20,
		marginBottom: 50,
		marginLeft:   60,
	}
}

// SetMargin sets all margins at once
func (mc *MarginConvention) SetMargin(top, right, bottom, left float64) *MarginConvention {
	mc.marginTop = top
	mc.marginRight = right
	mc.marginBottom = bottom
	mc.marginLeft = left
	return mc
}

// SetUniformMargin sets the same margin on all sides
func (mc *MarginConvention) SetUniformMargin(margin float64) *MarginConvention {
	mc.marginTop = margin
	mc.marginRight = margin
	mc.marginBottom = margin
	mc.marginLeft = margin
	return mc
}

// PlotWidth returns the width of the plot area
func (mc *MarginConvention) PlotWidth() float64 {
	return mc.totalWidth - mc.marginLeft - mc.marginRight
}

// PlotHeight returns the height of the plot area
func (mc *MarginConvention) PlotHeight() float64 {
	return mc.totalHeight - mc.marginTop - mc.marginBottom
}

// PlotArea returns a Rect representing the plot area
func (mc *MarginConvention) PlotArea() layout.Rect {
	return layout.Rect{
		X:      mc.marginLeft,
		Y:      mc.marginTop,
		Width:  mc.PlotWidth(),
		Height: mc.PlotHeight(),
	}
}

// TotalBounds returns the full canvas bounds
func (mc *MarginConvention) TotalBounds() layout.Rect {
	return layout.Rect{
		X:      0,
		Y:      0,
		Width:  mc.totalWidth,
		Height: mc.totalHeight,
	}
}

// LeftMarginArea returns the bounds of the left margin (for Y axis)
func (mc *MarginConvention) LeftMarginArea() layout.Rect {
	return layout.Rect{
		X:      0,
		Y:      mc.marginTop,
		Width:  mc.marginLeft,
		Height: mc.PlotHeight(),
	}
}

// RightMarginArea returns the bounds of the right margin
func (mc *MarginConvention) RightMarginArea() layout.Rect {
	plotArea := mc.PlotArea()
	return layout.Rect{
		X:      plotArea.X + plotArea.Width,
		Y:      mc.marginTop,
		Width:  mc.marginRight,
		Height: mc.PlotHeight(),
	}
}

// TopMarginArea returns the bounds of the top margin (for title)
func (mc *MarginConvention) TopMarginArea() layout.Rect {
	return layout.Rect{
		X:      mc.marginLeft,
		Y:      0,
		Width:  mc.PlotWidth(),
		Height: mc.marginTop,
	}
}

// BottomMarginArea returns the bounds of the bottom margin (for X axis)
func (mc *MarginConvention) BottomMarginArea() layout.Rect {
	plotArea := mc.PlotArea()
	return layout.Rect{
		X:      mc.marginLeft,
		Y:      plotArea.Y + plotArea.Height,
		Width:  mc.PlotWidth(),
		Height: mc.marginBottom,
	}
}

// AsNode creates a layout.Node with the margin convention applied
func (mc *MarginConvention) AsNode() *layout.Node {
	node := &layout.Node{
		Style: layout.Style{
			Width:  layout.Px(mc.totalWidth),
			Height: layout.Px(mc.totalHeight),
			Padding: layout.Spacing{
				Top:    layout.Px(mc.marginTop),
				Right:  layout.Px(mc.marginRight),
				Bottom: layout.Px(mc.marginBottom),
				Left:   layout.Px(mc.marginLeft),
			},
		},
	}
	return node
}

// ComputeMarginForAxes automatically computes margins based on axis requirements
func ComputeMarginForAxes(hasLeft, hasRight, hasTop, hasBottom, hasTitle bool) (top, right, bottom, left float64) {
	// Left margin for Y axis
	if hasLeft {
		left = 60 // Enough for Y axis labels
	} else {
		left = 10
	}

	// Right margin
	if hasRight {
		right = 60
	} else {
		right = 10
	}

	// Top margin
	if hasTitle {
		top = 40
	} else if hasTop {
		top = 30
	} else {
		top = 10
	}

	// Bottom margin for X axis
	if hasBottom {
		bottom = 50 // Enough for X axis labels
	} else {
		bottom = 10
	}

	return
}

// DefaultChartMargin returns typical margins for a chart with axes
func DefaultChartMargin() (top, right, bottom, left float64) {
	return ComputeMarginForAxes(true, false, false, true, true)
}
