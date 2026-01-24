package layout

import (
	"math"

	"github.com/SCKelemen/units"
)

// GridLayout manages a grid-based layout system for facets and dashboards.
//
// Example:
//   grid := NewGridLayout(units.Px(1200), units.Px(800), 2, 3)
//   grid.SetGap(units.Px(10))
//   cells := grid.Cells()
//   // cells[0][0], cells[0][1], cells[0][2]
//   // cells[1][0], cells[1][1], cells[1][2]
type GridLayout struct {
	bounds Rect
	rows   int
	cols   int
	gap    units.Length
	rowGap units.Length
	colGap units.Length
	margin Margin
}

// NewGridLayout creates a new grid layout
func NewGridLayout(width, height units.Length, rows, cols int) *GridLayout {
	return &GridLayout{
		bounds: Rect{
			X:      units.Px(0),
			Y:      units.Px(0),
			Width:  width,
			Height: height,
		},
		rows:   rows,
		cols:   cols,
		gap:    units.Px(10),
		rowGap: units.Px(0),
		colGap: units.Px(0),
		margin: Uniform(units.Px(0)),
	}
}

// SetGap sets the gap between cells (both row and column)
func (g *GridLayout) SetGap(gap units.Length) *GridLayout {
	g.gap = gap
	return g
}

// SetRowGap sets the gap between rows (overrides gap)
func (g *GridLayout) SetRowGap(gap units.Length) *GridLayout {
	g.rowGap = gap
	return g
}

// SetColGap sets the gap between columns (overrides gap)
func (g *GridLayout) SetColGap(gap units.Length) *GridLayout {
	g.colGap = gap
	return g
}

// SetMargin sets the outer margin
func (g *GridLayout) SetMargin(margin Margin) *GridLayout {
	g.margin = margin
	return g
}

// Cells returns a 2D array of cell rectangles
func (g *GridLayout) Cells() [][]Rect {
	// Apply margin to get content area
	contentArea := ApplyMargin(g.bounds, g.margin)

	// Use specified gaps or fall back to general gap
	rowGap := g.rowGap
	if rowGap.Value == 0 {
		rowGap = g.gap
	}
	colGap := g.colGap
	if colGap.Value == 0 {
		colGap = g.gap
	}

	return SplitIntoGrid(contentArea, g.rows, g.cols, g.gap)
}

// Cell returns a single cell at the given position
func (g *GridLayout) Cell(row, col int) Rect {
	cells := g.Cells()
	if row < 0 || row >= len(cells) || col < 0 || col >= len(cells[0]) {
		return Rect{}
	}
	return cells[row][col]
}

// CellWithSpan returns a cell that spans multiple rows/columns
func (g *GridLayout) CellWithSpan(row, col, rowSpan, colSpan int) Rect {
	cells := g.Cells()
	if row < 0 || row >= len(cells) || col < 0 || col >= len(cells[0]) {
		return Rect{}
	}

	// Start with the base cell
	startCell := cells[row][col]

	// Calculate end cell
	endRow := row + rowSpan - 1
	endCol := col + colSpan - 1

	if endRow >= len(cells) {
		endRow = len(cells) - 1
	}
	if endCol >= len(cells[0]) {
		endCol = len(cells[0]) - 1
	}

	endCell := cells[endRow][endCol]

	// Calculate spanned bounds
	return Rect{
		X:      startCell.X,
		Y:      startCell.Y,
		Width:  units.Px(endCell.X.Value + endCell.Width.Value - startCell.X.Value),
		Height: units.Px(endCell.Y.Value + endCell.Height.Value - startCell.Y.Value),
	}
}

// FlatCells returns all cells as a flat slice (row-major order)
func (g *GridLayout) FlatCells() []Rect {
	cells := g.Cells()
	flat := make([]Rect, 0, g.rows*g.cols)
	for _, row := range cells {
		flat = append(flat, row...)
	}
	return flat
}

// AutoGrid automatically determines the best grid dimensions for n items
func AutoGrid(n int) (rows, cols int) {
	if n <= 0 {
		return 0, 0
	}
	if n == 1 {
		return 1, 1
	}

	// Try to create a square-ish grid
	cols = int(math.Ceil(math.Sqrt(float64(n)))) // ceil(sqrt(n))
	rows = (n + cols - 1) / cols                  // ceil(n / cols)

	return rows, cols
}

// FacetGrid creates a grid layout for faceted plots
type FacetGrid struct {
	*GridLayout
	facetMargin Margin
	showTitles  bool
	scaleShared ScaleSharing
}

// NewFacetGrid creates a new facet grid
func NewFacetGrid(width, height units.Length, rows, cols int) *FacetGrid {
	return &FacetGrid{
		GridLayout:  NewGridLayout(width, height, rows, cols),
		facetMargin: Uniform(units.Px(5)),
		showTitles:  true,
		scaleShared: ScaleShareNone,
	}
}

// SetFacetMargin sets the margin for each facet
func (fg *FacetGrid) SetFacetMargin(margin Margin) *FacetGrid {
	fg.facetMargin = margin
	return fg
}

// SetShowTitles sets whether to show titles for each facet
func (fg *FacetGrid) SetShowTitles(show bool) *FacetGrid {
	fg.showTitles = show
	return fg
}

// SetScaleSharing sets how scales are shared across facets
func (fg *FacetGrid) SetScaleSharing(sharing ScaleSharing) *FacetGrid {
	fg.scaleShared = sharing
	return fg
}

// FacetCell returns the plotting area for a facet (with margin applied)
func (fg *FacetGrid) FacetCell(row, col int) Rect {
	cell := fg.Cell(row, col)
	return ApplyMargin(cell, fg.facetMargin)
}

// FacetTitleArea returns the area for the facet title
func (fg *FacetGrid) FacetTitleArea(row, col int) Rect {
	cell := fg.Cell(row, col)
	return Rect{
		X:      cell.X,
		Y:      cell.Y,
		Width:  cell.Width,
		Height: fg.facetMargin.Top,
	}
}
