package gallery

import "github.com/SCKelemen/units"

// CalculateGridDimensions calculates pixel dimensions for a grid-based gallery
// using relative units that resolve to exact pixels at render time
func CalculateGridDimensions(cols, rows int, baseWidth, baseHeight float64) GalleryDimensions {
	// Use percentages for grid sizing to avoid accumulation errors
	colPct := units.Percent(100.0 / float64(cols))
	rowPct := units.Percent(100.0 / float64(rows))

	// Calculate dimensions with proper margins
	titleMargin := units.Percent(5)  // 5% top margin for title
	bottomMargin := units.Percent(3) // 3% bottom margin
	chartPadding := units.Percent(2) // 2% padding within each cell

	totalWidth := baseWidth * float64(cols)
	totalHeight := baseHeight * float64(rows)

	// Add margins to total height
	titleSpace := titleMargin.Of(totalHeight)
	bottomSpace := bottomMargin.Of(totalHeight)
	totalHeight += titleSpace + bottomSpace

	// Calculate chart dimensions (subtract padding)
	colWidth := colPct.Of(totalWidth)
	rowHeight := rowPct.Of(baseHeight * float64(rows))

	padding := chartPadding.Of(colWidth)
	chartWidth := colWidth - (padding * 2)
	chartHeight := rowHeight - (padding * 2)

	return GalleryDimensions{
		TotalWidth:   totalWidth,
		TotalHeight:  totalHeight,
		ChartWidth:   chartWidth,
		ChartHeight:  chartHeight,
		ColWidth:     colWidth,
		RowHeight:    rowHeight,
		TitleY:       titleSpace * 0.7, // Position title 70% down the title space
		ChartStartY:  titleSpace,
		BottomMargin: bottomSpace,
	}
}

// CalculateSingleRowDimensions calculates dimensions for single-row galleries
func CalculateSingleRowDimensions(cols int, baseWidth, baseHeight float64) GalleryDimensions {
	titleHeight := 50.0
	bottomMargin := 30.0
	chartPadding := 25.0

	totalWidth := baseWidth * float64(cols)
	totalHeight := baseHeight + titleHeight + bottomMargin

	colWidth := totalWidth / float64(cols)
	chartWidth := baseWidth - (chartPadding * 2)
	chartHeight := baseHeight - chartPadding

	return GalleryDimensions{
		TotalWidth:   totalWidth,
		TotalHeight:  totalHeight,
		ChartWidth:   chartWidth,
		ChartHeight:  chartHeight,
		ColWidth:     colWidth,
		RowHeight:    baseHeight,
		TitleY:       30,
		ChartStartY:  titleHeight + 10,
		BottomMargin: bottomMargin,
	}
}

// SingleRowLayout represents a single-row gallery layout
type SingleRowLayout struct {
	Cols       int
	BaseWidth  float64
	BaseHeight float64
	StartX     float64 // Initial X offset for first cell (default: 0)
}

// CalculateDimensions calculates dimensions for a single-row layout
func (l *SingleRowLayout) CalculateDimensions() GalleryDimensions {
	return CalculateSingleRowDimensions(l.Cols, l.BaseWidth, l.BaseHeight)
}

// GetCellPosition returns the position for a variant at the given index
func (l *SingleRowLayout) GetCellPosition(variantIndex int) (x, y float64) {
	dims := l.CalculateDimensions()
	cellX := l.StartX + float64(variantIndex)*dims.ColWidth
	return cellX, dims.ChartStartY
}

// GridLayout represents a multi-row grid gallery layout
type GridLayout struct {
	Cols       int
	Rows       int
	BaseWidth  float64
	BaseHeight float64
}

// CalculateDimensions calculates dimensions for a grid layout
func (l *GridLayout) CalculateDimensions() GalleryDimensions {
	return CalculateGridDimensions(l.Cols, l.Rows, l.BaseWidth, l.BaseHeight)
}

// GetCellPosition returns the position for a variant at the given index
func (l *GridLayout) GetCellPosition(variantIndex int) (x, y float64) {
	dims := l.CalculateDimensions()
	col := variantIndex % l.Cols
	row := variantIndex / l.Cols
	cellX := float64(col) * dims.ColWidth
	cellY := dims.ChartStartY + float64(row)*dims.RowHeight
	return cellX, cellY
}

// VerticalStackLayout represents a vertical stack gallery layout
type VerticalStackLayout struct {
	Rows       int
	BaseWidth  float64
	RowHeight  float64
	RowSpacing float64
}

// CalculateDimensions calculates dimensions for a vertical stack layout
func (l *VerticalStackLayout) CalculateDimensions() GalleryDimensions {
	titleHeight := 50.0
	bottomMargin := 30.0
	totalHeight := l.RowHeight*float64(l.Rows) +
		titleHeight +
		l.RowSpacing*float64(l.Rows-1) +
		bottomMargin

	return GalleryDimensions{
		TotalWidth:   l.BaseWidth,
		TotalHeight:  totalHeight,
		ChartWidth:   l.BaseWidth - 50,
		ChartHeight:  l.RowHeight - 80,
		ColWidth:     l.BaseWidth,
		RowHeight:    l.RowHeight,
		TitleY:       30,
		ChartStartY:  titleHeight + 10,
		BottomMargin: bottomMargin,
	}
}

// GetCellPosition returns the position for a variant at the given index
func (l *VerticalStackLayout) GetCellPosition(variantIndex int) (x, y float64) {
	dims := l.CalculateDimensions()
	cellY := dims.ChartStartY + float64(variantIndex)*(l.RowHeight+l.RowSpacing)
	return 25.0, cellY // Fixed X offset for vertical stack
}
