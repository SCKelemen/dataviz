package layout

import "github.com/SCKelemen/units"

// MarginConvention provides standard margin calculations for charts.
// Based on the D3 margin convention pattern.
//
// Example:
//   mc := NewMarginConvention(units.Px(800), units.Px(600))
//   mc.SetMargin(DefaultMargin())
//   plotArea := mc.PlotArea()
//   // plotArea contains the inner drawing area after margins
type MarginConvention struct {
	totalWidth  units.Length
	totalHeight units.Length
	margin      Margin
}

// NewMarginConvention creates a new margin convention with total dimensions
func NewMarginConvention(width, height units.Length) *MarginConvention {
	return &MarginConvention{
		totalWidth:  width,
		totalHeight: height,
		margin:      DefaultMargin(),
	}
}

// SetMargin updates the margin
func (mc *MarginConvention) SetMargin(margin Margin) *MarginConvention {
	mc.margin = margin
	return mc
}

// TotalBounds returns the full canvas bounds
func (mc *MarginConvention) TotalBounds() Rect {
	return Rect{
		X:      units.Px(0),
		Y:      units.Px(0),
		Width:  mc.totalWidth,
		Height: mc.totalHeight,
	}
}

// PlotArea returns the inner plotting area (after applying margins)
func (mc *MarginConvention) PlotArea() Rect {
	return ApplyMargin(mc.TotalBounds(), mc.margin)
}

// PlotWidth returns the width of the plot area
func (mc *MarginConvention) PlotWidth() units.Length {
	return units.Px(mc.totalWidth.Value - mc.margin.Left.Value - mc.margin.Right.Value)
}

// PlotHeight returns the height of the plot area
func (mc *MarginConvention) PlotHeight() units.Length {
	return units.Px(mc.totalHeight.Value - mc.margin.Top.Value - mc.margin.Bottom.Value)
}

// LeftMarginArea returns the bounds of the left margin (for Y axis)
func (mc *MarginConvention) LeftMarginArea() Rect {
	return Rect{
		X:      units.Px(0),
		Y:      mc.margin.Top,
		Width:  mc.margin.Left,
		Height: mc.PlotHeight(),
	}
}

// RightMarginArea returns the bounds of the right margin
func (mc *MarginConvention) RightMarginArea() Rect {
	plotArea := mc.PlotArea()
	return Rect{
		X:      units.Px(plotArea.X.Value + plotArea.Width.Value),
		Y:      mc.margin.Top,
		Width:  mc.margin.Right,
		Height: mc.PlotHeight(),
	}
}

// TopMarginArea returns the bounds of the top margin (for title)
func (mc *MarginConvention) TopMarginArea() Rect {
	return Rect{
		X:      mc.margin.Left,
		Y:      units.Px(0),
		Width:  mc.PlotWidth(),
		Height: mc.margin.Top,
	}
}

// BottomMarginArea returns the bounds of the bottom margin (for X axis)
func (mc *MarginConvention) BottomMarginArea() Rect {
	plotArea := mc.PlotArea()
	return Rect{
		X:      mc.margin.Left,
		Y:      units.Px(plotArea.Y.Value + plotArea.Height.Value),
		Width:  mc.PlotWidth(),
		Height: mc.margin.Bottom,
	}
}

// WithPadding creates a new MarginConvention with additional internal padding
func (mc *MarginConvention) WithPadding(padding Padding) Rect {
	plotArea := mc.PlotArea()
	return ApplyMargin(plotArea, padding)
}

// ComputeMarginForAxes automatically computes margins based on axis requirements
func ComputeMarginForAxes(hasLeft, hasRight, hasTop, hasBottom, hasTitle bool) Margin {
	margin := Margin{}

	// Left margin for Y axis
	if hasLeft {
		margin.Left = units.Px(60) // Enough for Y axis labels
	} else {
		margin.Left = units.Px(10)
	}

	// Right margin
	if hasRight {
		margin.Right = units.Px(60)
	} else {
		margin.Right = units.Px(10)
	}

	// Top margin
	if hasTitle {
		margin.Top = units.Px(40)
	} else if hasTop {
		margin.Top = units.Px(30)
	} else {
		margin.Top = units.Px(10)
	}

	// Bottom margin for X axis
	if hasBottom {
		margin.Bottom = units.Px(50) // Enough for X axis labels
	} else {
		margin.Bottom = units.Px(10)
	}

	return margin
}

// SplitHorizontal splits a rectangle horizontally at the given ratio (0-1)
func SplitHorizontal(bounds Rect, ratio float64) (left, right Rect) {
	splitX := bounds.Width.Value * ratio

	left = Rect{
		X:      bounds.X,
		Y:      bounds.Y,
		Width:  units.Px(splitX),
		Height: bounds.Height,
	}

	right = Rect{
		X:      units.Px(bounds.X.Value + splitX),
		Y:      bounds.Y,
		Width:  units.Px(bounds.Width.Value - splitX),
		Height: bounds.Height,
	}

	return
}

// SplitVertical splits a rectangle vertically at the given ratio (0-1)
func SplitVertical(bounds Rect, ratio float64) (top, bottom Rect) {
	splitY := bounds.Height.Value * ratio

	top = Rect{
		X:      bounds.X,
		Y:      bounds.Y,
		Width:  bounds.Width,
		Height: units.Px(splitY),
	}

	bottom = Rect{
		X:      bounds.X,
		Y:      units.Px(bounds.Y.Value + splitY),
		Width:  bounds.Width,
		Height: units.Px(bounds.Height.Value - splitY),
	}

	return
}

// SplitIntoGrid splits a rectangle into a grid of cells
func SplitIntoGrid(bounds Rect, rows, cols int, gap units.Length) [][]Rect {
	if rows <= 0 || cols <= 0 {
		return nil
	}

	// Calculate total gap space
	totalRowGap := float64(rows-1) * gap.Value
	totalColGap := float64(cols-1) * gap.Value

	// Calculate cell dimensions
	cellWidth := (bounds.Width.Value - totalColGap) / float64(cols)
	cellHeight := (bounds.Height.Value - totalRowGap) / float64(rows)

	grid := make([][]Rect, rows)
	for row := 0; row < rows; row++ {
		grid[row] = make([]Rect, cols)
		for col := 0; col < cols; col++ {
			grid[row][col] = Rect{
				X:      units.Px(bounds.X.Value + float64(col)*(cellWidth+gap.Value)),
				Y:      units.Px(bounds.Y.Value + float64(row)*(cellHeight+gap.Value)),
				Width:  units.Px(cellWidth),
				Height: units.Px(cellHeight),
			}
		}
	}

	return grid
}

// Inset creates a smaller rectangle inset by the given amount on all sides
func Inset(bounds Rect, amount units.Length) Rect {
	return ApplyMargin(bounds, Uniform(amount))
}

// Center returns the center point of a rectangle
func Center(bounds Rect) (x, y units.Length) {
	x = units.Px(bounds.X.Value + bounds.Width.Value/2)
	y = units.Px(bounds.Y.Value + bounds.Height.Value/2)
	return
}

// Contains checks if a point is within the rectangle
func Contains(bounds Rect, x, y units.Length) bool {
	return x.Value >= bounds.X.Value &&
		x.Value <= bounds.X.Value+bounds.Width.Value &&
		y.Value >= bounds.Y.Value &&
		y.Value <= bounds.Y.Value+bounds.Height.Value
}
