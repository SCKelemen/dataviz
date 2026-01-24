package layout

import (
	"testing"

	"github.com/SCKelemen/dataviz/transforms"
	"github.com/SCKelemen/units"
)

func TestNewFacet(t *testing.T) {
	facet := NewFacet("Group")

	if facet.Field != "Group" {
		t.Errorf("Expected field 'Group', got '%s'", facet.Field)
	}
	if facet.Layout != FacetWrap {
		t.Error("Expected default layout to be wrap")
	}
	if facet.NCols != 2 {
		t.Errorf("Expected default 2 cols, got %d", facet.NCols)
	}
	if facet.ShowTitles != true {
		t.Error("Expected titles to be shown by default")
	}
}

func TestFacet_Split(t *testing.T) {
	data := []transforms.DataPoint{
		{Group: "A", Y: 1},
		{Group: "A", Y: 2},
		{Group: "B", Y: 3},
		{Group: "B", Y: 4},
		{Group: "C", Y: 5},
	}

	facet := NewFacet("Group")
	facets := facet.Split(data)

	if len(facets) != 3 {
		t.Errorf("Expected 3 facets, got %d", len(facets))
	}

	// Check that data is split correctly
	foundA := false
	foundB := false
	foundC := false

	for _, f := range facets {
		switch f.Value {
		case "A":
			foundA = true
			if len(f.Data) != 2 {
				t.Errorf("Group A should have 2 points, got %d", len(f.Data))
			}
		case "B":
			foundB = true
			if len(f.Data) != 2 {
				t.Errorf("Group B should have 2 points, got %d", len(f.Data))
			}
		case "C":
			foundC = true
			if len(f.Data) != 1 {
				t.Errorf("Group C should have 1 point, got %d", len(f.Data))
			}
		}
	}

	if !foundA || !foundB || !foundC {
		t.Error("Not all groups found in facets")
	}
}

func TestFacet_Split_WithOrder(t *testing.T) {
	data := []transforms.DataPoint{
		{Group: "B", Y: 1},
		{Group: "A", Y: 2},
		{Group: "C", Y: 3},
	}

	facet := NewFacet("Group").WithOrder([]string{"C", "B", "A"})
	facets := facet.Split(data)

	if len(facets) != 3 {
		t.Errorf("Expected 3 facets, got %d", len(facets))
	}

	// Check order
	if facets[0].Value != "C" {
		t.Errorf("First facet should be C, got %s", facets[0].Value)
	}
	if facets[1].Value != "B" {
		t.Errorf("Second facet should be B, got %s", facets[1].Value)
	}
	if facets[2].Value != "A" {
		t.Errorf("Third facet should be A, got %s", facets[2].Value)
	}
}

func TestFacet_Split_EmptyData(t *testing.T) {
	facet := NewFacet("Group")
	facets := facet.Split([]transforms.DataPoint{})

	if facets != nil {
		t.Error("Expected nil for empty data")
	}
}

func TestFacet_Split_ByLabel(t *testing.T) {
	data := []transforms.DataPoint{
		{Label: "X", Y: 1},
		{Label: "Y", Y: 2},
		{Label: "X", Y: 3},
	}

	facet := NewFacet("Label")
	facets := facet.Split(data)

	if len(facets) != 2 {
		t.Errorf("Expected 2 facets, got %d", len(facets))
	}
}

func TestFacet_CalculateDimensions(t *testing.T) {
	tests := []struct {
		layout     FacetLayout
		rows       int
		cols       int
		numFacets  int
		expectRows int
		expectCols int
	}{
		{FacetWrap, 0, 2, 5, 3, 2},           // Wrap with 2 cols, 5 facets = 3 rows
		{FacetWrap, 0, 3, 7, 3, 3},           // Wrap with 3 cols, 7 facets = 3 rows
		{FacetLayoutGrid, 2, 2, 4, 2, 2},    // Fixed 2x2 grid
		{FacetLayoutGrid, 2, 0, 5, 2, 3},    // 2 rows, auto cols
		{FacetLayoutGrid, 0, 2, 5, 3, 2},    // Auto rows, 2 cols
	}

	for i, tt := range tests {
		facet := NewFacet("Group").
			WithLayout(tt.layout).
			WithRows(tt.rows).
			WithCols(tt.cols)

		rows, cols := facet.CalculateDimensions(tt.numFacets)

		if rows != tt.expectRows || cols != tt.expectCols {
			t.Errorf("Test %d: expected %dx%d, got %dx%d",
				i, tt.expectRows, tt.expectCols, rows, cols)
		}
	}
}

func TestFacet_ChainedMethods(t *testing.T) {
	facet := NewFacet("Group").
		WithLayout(FacetLayoutGrid).
		WithCols(3).
		WithRows(2).
		WithScaleSharing(ScaleShareX).
		WithTitles(false).
		WithGap(units.Px(15)).
		WithFacetMargin(Uniform(units.Px(10)))

	if facet.Layout != FacetLayoutGrid {
		t.Error("Layout should be grid")
	}
	if facet.NCols != 3 {
		t.Error("Cols should be 3")
	}
	if facet.NRows != 2 {
		t.Error("Rows should be 2")
	}
	if facet.ScaleSharing != ScaleShareX {
		t.Error("Scale sharing should be X")
	}
	if facet.ShowTitles != false {
		t.Error("Titles should be hidden")
	}
	if facet.Gap.Value != 15 {
		t.Error("Gap should be 15")
	}
}

func TestComputeSharedDomain(t *testing.T) {
	facets := []FacetData{
		{
			Value: "A",
			Data: []transforms.DataPoint{
				{Y: 1}, {Y: 5}, {Y: 3},
			},
		},
		{
			Value: "B",
			Data: []transforms.DataPoint{
				{Y: 2}, {Y: 8}, {Y: 4},
			},
		},
	}

	min, max := ComputeSharedDomain(facets, "Y")

	if min != 1 {
		t.Errorf("Expected min 1, got %f", min)
	}
	if max != 8 {
		t.Errorf("Expected max 8, got %f", max)
	}
}

func TestGetScaleDomain_None(t *testing.T) {
	facets := []FacetData{
		{
			Value: "A",
			Data: []transforms.DataPoint{
				{Y: 1}, {Y: 5},
			},
		},
		{
			Value: "B",
			Data: []transforms.DataPoint{
				{Y: 2}, {Y: 8},
			},
		},
	}

	// Independent scales - should only use facet 0's data
	min, max := GetScaleDomain(facets, 0, "Y", ScaleShareNone)

	if min != 1 {
		t.Errorf("Expected min 1, got %f", min)
	}
	if max != 5 {
		t.Errorf("Expected max 5, got %f", max)
	}
}

func TestGetScaleDomain_XY(t *testing.T) {
	facets := []FacetData{
		{
			Value: "A",
			Data: []transforms.DataPoint{
				{Y: 1}, {Y: 5},
			},
		},
		{
			Value: "B",
			Data: []transforms.DataPoint{
				{Y: 2}, {Y: 8},
			},
		},
	}

	// Shared scales - should use all data
	min, max := GetScaleDomain(facets, 0, "Y", ScaleShareXY)

	if min != 1 {
		t.Errorf("Expected min 1, got %f", min)
	}
	if max != 8 {
		t.Errorf("Expected max 8, got %f", max)
	}
}

func TestNewFacetRenderer(t *testing.T) {
	facet := NewFacet("Group").WithCols(2)
	renderer := NewFacetRenderer(facet, units.Px(800), units.Px(600))

	if renderer.Facet != facet {
		t.Error("Facet should be set")
	}
	if renderer.Grid == nil {
		t.Error("Grid should be created")
	}
	if renderer.Grid.bounds.Width.Value != 800 {
		t.Error("Grid width should match")
	}
}

func TestFacetedChart(t *testing.T) {
	data := []transforms.DataPoint{
		{Group: "A", Y: 1},
		{Group: "B", Y: 2},
	}

	facet := NewFacet("Group")
	chartRenderer := func(data []transforms.DataPoint, bounds Rect) string {
		return "<rect/>"
	}

	fc := &FacetedChart{
		Data:          data,
		Facet:         facet,
		Width:         units.Px(800),
		Height:        units.Px(600),
		ChartRenderer: chartRenderer,
	}

	svg := fc.Render()

	if svg == "" {
		t.Error("Should produce SVG output")
	}
}

func TestFacetRenderer_Render(t *testing.T) {
	data := []transforms.DataPoint{
		{Group: "A", Y: 1},
		{Group: "A", Y: 2},
		{Group: "B", Y: 3},
		{Group: "B", Y: 4},
	}

	facet := NewFacet("Group").WithCols(2)
	renderer := NewFacetRenderer(facet, units.Px(800), units.Px(600))

	chartRenderer := func(data []transforms.DataPoint, bounds Rect) string {
		return `<rect class="chart"/>`
	}
	renderer.ChartRenderer = chartRenderer

	svg := renderer.Render(data)

	if svg == "" {
		t.Error("Should produce SVG output")
	}
}

func TestScaleSharing_Constants(t *testing.T) {
	// Just verify the constants exist and have expected values
	if ScaleShareNone != "none" {
		t.Error("ScaleShareNone should be 'none'")
	}
	if ScaleShareX != "x" {
		t.Error("ScaleShareX should be 'x'")
	}
	if ScaleShareY != "y" {
		t.Error("ScaleShareY should be 'y'")
	}
	if ScaleShareXY != "xy" {
		t.Error("ScaleShareXY should be 'xy'")
	}
}

func TestFacetLayout_Constants(t *testing.T) {
	if FacetWrap != "wrap" {
		t.Error("FacetWrap should be 'wrap'")
	}
	if FacetLayoutGrid != "grid" {
		t.Error("FacetLayoutGrid should be 'grid'")
	}
	if FacetCustom != "custom" {
		t.Error("FacetCustom should be 'custom'")
	}
}
