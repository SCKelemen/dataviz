package layout

import (
	"testing"

	"github.com/SCKelemen/units"
)

func TestNewGridLayout(t *testing.T) {
	grid := NewGridLayout(units.Px(800), units.Px(600), 2, 3)

	if grid.rows != 2 {
		t.Errorf("Expected 2 rows, got %d", grid.rows)
	}
	if grid.cols != 3 {
		t.Errorf("Expected 3 cols, got %d", grid.cols)
	}
	if grid.bounds.Width.Value != 800 {
		t.Errorf("Expected width 800, got %f", grid.bounds.Width.Value)
	}
	if grid.bounds.Height.Value != 600 {
		t.Errorf("Expected height 600, got %f", grid.bounds.Height.Value)
	}
}

func TestGridLayout_Cells(t *testing.T) {
	grid := NewGridLayout(units.Px(900), units.Px(600), 2, 3)
	grid.SetGap(units.Px(10))

	cells := grid.Cells()

	// Check dimensions
	if len(cells) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(cells))
	}
	if len(cells[0]) != 3 {
		t.Errorf("Expected 3 cols, got %d", len(cells[0]))
	}

	// Check that cells don't overlap and have reasonable sizes
	cell := cells[0][0]
	if cell.Width.Value <= 0 || cell.Height.Value <= 0 {
		t.Error("Cell has invalid dimensions")
	}

	// Check gap between cells
	cell1 := cells[0][0]
	cell2 := cells[0][1]
	gap := cell2.X.Value - (cell1.X.Value + cell1.Width.Value)
	if gap != 10 {
		t.Errorf("Expected gap of 10, got %f", gap)
	}
}

func TestGridLayout_Cell(t *testing.T) {
	grid := NewGridLayout(units.Px(800), units.Px(600), 2, 2)

	cell := grid.Cell(0, 1)
	if cell.X.Value <= 0 {
		t.Error("Cell should have non-zero X position")
	}

	// Test out of bounds
	invalidCell := grid.Cell(5, 5)
	if invalidCell.Width.Value != 0 {
		t.Error("Out of bounds cell should return empty rect")
	}
}

func TestGridLayout_CellWithSpan(t *testing.T) {
	grid := NewGridLayout(units.Px(900), units.Px(600), 3, 3)
	grid.SetGap(units.Px(10))

	// Get a cell that spans 2x2
	cell := grid.CellWithSpan(0, 0, 2, 2)

	// Should be larger than a single cell
	singleCell := grid.Cell(0, 0)
	if cell.Width.Value <= singleCell.Width.Value {
		t.Error("Spanned cell should be wider than single cell")
	}
	if cell.Height.Value <= singleCell.Height.Value {
		t.Error("Spanned cell should be taller than single cell")
	}
}

func TestGridLayout_FlatCells(t *testing.T) {
	grid := NewGridLayout(units.Px(800), units.Px(600), 2, 3)

	flat := grid.FlatCells()

	if len(flat) != 6 {
		t.Errorf("Expected 6 cells, got %d", len(flat))
	}
}

func TestAutoGrid(t *testing.T) {
	tests := []struct {
		n           int
		minRows     int
		minCols     int
		description string
	}{
		{0, 0, 0, "zero items"},
		{1, 1, 1, "one item"},
		{4, 2, 2, "four items"},
		{6, 2, 3, "six items"},
		{9, 3, 3, "nine items"},
	}

	for _, tt := range tests {
		rows, cols := AutoGrid(tt.n)
		if rows < tt.minRows || cols < tt.minCols {
			t.Errorf("AutoGrid(%d): got %dx%d, expected at least %dx%d (%s)",
				tt.n, rows, cols, tt.minRows, tt.minCols, tt.description)
		}
		if rows*cols < tt.n {
			t.Errorf("AutoGrid(%d): %dx%d = %d cells, not enough for %d items",
				tt.n, rows, cols, rows*cols, tt.n)
		}
	}
}

func TestFacetGrid(t *testing.T) {
	fg := NewFacetGrid(units.Px(1000), units.Px(800), 2, 2)
	fg.SetFacetMargin(Uniform(units.Px(5)))
	fg.SetShowTitles(true)
	fg.SetScaleSharing(ScaleShareXY)

	if fg.showTitles != true {
		t.Error("ShowTitles should be true")
	}
	if fg.scaleShared != ScaleShareXY {
		t.Error("ScaleSharing should be XY")
	}

	// Test facet cell (should be smaller than raw cell due to margin)
	rawCell := fg.Cell(0, 0)
	facetCell := fg.FacetCell(0, 0)

	if facetCell.Width.Value >= rawCell.Width.Value {
		t.Error("Facet cell should be smaller than raw cell due to margin")
	}
}

func TestFacetGrid_TitleArea(t *testing.T) {
	fg := NewFacetGrid(units.Px(800), units.Px(600), 2, 2)
	fg.SetFacetMargin(Margin{
		Top:    units.Px(20),
		Right:  units.Px(5),
		Bottom: units.Px(5),
		Left:   units.Px(5),
	})

	titleArea := fg.FacetTitleArea(0, 0)

	if titleArea.Height.Value != 20 {
		t.Errorf("Title area height should be 20, got %f", titleArea.Height.Value)
	}
}

func TestGridLayout_WithMargin(t *testing.T) {
	grid := NewGridLayout(units.Px(800), units.Px(600), 2, 2)
	grid.SetMargin(Margin{
		Top:    units.Px(10),
		Right:  units.Px(10),
		Bottom: units.Px(10),
		Left:   units.Px(10),
	})

	cells := grid.Cells()

	// First cell should start at (10, 10) due to margin
	if cells[0][0].X.Value != 10 {
		t.Errorf("Expected first cell X at 10, got %f", cells[0][0].X.Value)
	}
	if cells[0][0].Y.Value != 10 {
		t.Errorf("Expected first cell Y at 10, got %f", cells[0][0].Y.Value)
	}
}

func TestGridLayout_ChainedMethods(t *testing.T) {
	grid := NewGridLayout(units.Px(800), units.Px(600), 2, 2).
		SetGap(units.Px(15)).
		SetMargin(Uniform(units.Px(10)))

	if grid.gap.Value != 15 {
		t.Errorf("Expected gap 15, got %f", grid.gap.Value)
	}
	if grid.margin.Top.Value != 10 {
		t.Errorf("Expected margin 10, got %f", grid.margin.Top.Value)
	}
}
