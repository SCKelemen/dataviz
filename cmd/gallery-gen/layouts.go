package main

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
