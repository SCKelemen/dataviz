package charts

import (
	"strings"
	"testing"
)

func TestRenderHistogram(t *testing.T) {
	data := &HistogramData{
		Values: []float64{1, 2, 2, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 6, 6, 7},
		Color:  "#4285f4",
		Label:  "Test Data",
	}

	spec := HistogramSpec{
		Data:     data,
		Width:    400,
		Height:   300,
		BinCount: 7,
		Nice:     true,
	}

	svg := RenderHistogram(spec)

	if svg == "" {
		t.Error("Expected non-empty SVG output")
	}

	// Check for expected SVG elements
	if !strings.Contains(svg, "<rect") {
		t.Error("Expected histogram bars (rect elements)")
	}

	// Should have multiple bars
	rectCount := strings.Count(svg, "<rect")
	if rectCount < 3 {
		t.Errorf("Expected at least 3 bars, got %d", rectCount)
	}
}

func TestRenderHistogram_FixedBinSize(t *testing.T) {
	data := &HistogramData{
		Values: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	}

	spec := HistogramSpec{
		Data:    data,
		Width:   400,
		Height:  300,
		BinSize: 2.0, // Fixed bin size of 2
	}

	svg := RenderHistogram(spec)

	if svg == "" {
		t.Error("Expected non-empty SVG output")
	}

	if !strings.Contains(svg, "<rect") {
		t.Error("Expected histogram bars")
	}
}

func TestRenderHistogram_Density(t *testing.T) {
	data := &HistogramData{
		Values: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	}

	spec := HistogramSpec{
		Data:        data,
		Width:       400,
		Height:      300,
		BinCount:    5,
		ShowDensity: true, // Show density instead of counts
	}

	svg := RenderHistogram(spec)

	if svg == "" {
		t.Error("Expected non-empty SVG output")
	}

	if !strings.Contains(svg, "<rect") {
		t.Error("Expected histogram bars")
	}
}

func TestRenderHistogram_EmptyData(t *testing.T) {
	data := &HistogramData{
		Values: []float64{},
	}

	spec := HistogramSpec{
		Data:   data,
		Width:  400,
		Height: 300,
	}

	svg := RenderHistogram(spec)

	if svg != "" {
		t.Error("Expected empty SVG for empty data")
	}
}

func TestRenderDensityPlot(t *testing.T) {
	data1 := &DensityPlotData{
		Values: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		Color:  "#4285f4",
		Label:  "Dataset 1",
	}

	data2 := &DensityPlotData{
		Values: []float64{3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		Color:  "#ea4335",
		Label:  "Dataset 2",
	}

	spec := DensityPlotSpec{
		Data:      []*DensityPlotData{data1, data2},
		Width:     400,
		Height:    300,
		ShowFill:  true,
		LineWidth: 2,
	}

	svg := RenderDensityPlot(spec)

	if svg == "" {
		t.Error("Expected non-empty SVG output")
	}

	// Check for expected SVG elements
	if !strings.Contains(svg, "<path") {
		t.Error("Expected density curve paths")
	}

	// Should have paths for each dataset (fill + line)
	pathCount := strings.Count(svg, "<path")
	if pathCount < 4 { // 2 datasets × (fill + line)
		t.Errorf("Expected at least 4 paths, got %d", pathCount)
	}
}

func TestRenderDensityPlot_NoFill(t *testing.T) {
	data := &DensityPlotData{
		Values: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	}

	spec := DensityPlotSpec{
		Data:     []*DensityPlotData{data},
		Width:    400,
		Height:   300,
		ShowFill: false, // No fill, just line
	}

	svg := RenderDensityPlot(spec)

	if svg == "" {
		t.Error("Expected non-empty SVG output")
	}

	// Should have only line path (no fill)
	pathCount := strings.Count(svg, "<path")
	if pathCount != 1 {
		t.Errorf("Expected 1 path (line only), got %d", pathCount)
	}
}

func TestRenderDensityPlot_EmptyData(t *testing.T) {
	spec := DensityPlotSpec{
		Data:   []*DensityPlotData{},
		Width:  400,
		Height: 300,
	}

	svg := RenderDensityPlot(spec)

	if svg != "" {
		t.Error("Expected empty SVG for empty data")
	}
}

func TestRenderRidgeline(t *testing.T) {
	ridge1 := &RidgelineData{
		Values: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		Label:  "Group A",
		Color:  "#4285f4",
	}

	ridge2 := &RidgelineData{
		Values: []float64{2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
		Label:  "Group B",
		Color:  "#ea4335",
	}

	ridge3 := &RidgelineData{
		Values: []float64{3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		Label:  "Group C",
		Color:  "#fbbc04",
	}

	spec := RidgelineSpec{
		Data:       []*RidgelineData{ridge1, ridge2, ridge3},
		Width:      400,
		Height:     300,
		Overlap:    0.5,
		ShowFill:   true,
		ShowLabels: true,
	}

	svg := RenderRidgeline(spec)

	if svg == "" {
		t.Error("Expected non-empty SVG output")
	}

	// Check for expected SVG elements
	if !strings.Contains(svg, "<path") {
		t.Error("Expected ridgeline paths")
	}

	// Check for labels
	if !strings.Contains(svg, "Group A") {
		t.Error("Expected Group A label")
	}

	if !strings.Contains(svg, "Group B") {
		t.Error("Expected Group B label")
	}

	if !strings.Contains(svg, "Group C") {
		t.Error("Expected Group C label")
	}

	// Should have paths for each ridge (fill + line)
	pathCount := strings.Count(svg, "<path")
	if pathCount < 6 { // 3 ridges × (fill + line)
		t.Errorf("Expected at least 6 paths, got %d", pathCount)
	}
}

func TestRenderRidgeline_NoFill(t *testing.T) {
	ridge := &RidgelineData{
		Values: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		Label:  "Test",
	}

	spec := RidgelineSpec{
		Data:     []*RidgelineData{ridge},
		Width:    400,
		Height:   300,
		ShowFill: false, // No fill, just outlines
	}

	svg := RenderRidgeline(spec)

	if svg == "" {
		t.Error("Expected non-empty SVG output")
	}

	// Should have only line path (no fill)
	pathCount := strings.Count(svg, "<path")
	if pathCount != 1 {
		t.Errorf("Expected 1 path (line only), got %d", pathCount)
	}
}

func TestRenderRidgeline_Reversed(t *testing.T) {
	ridge1 := &RidgelineData{
		Values: []float64{1, 2, 3, 4, 5},
		Label:  "First",
	}

	ridge2 := &RidgelineData{
		Values: []float64{6, 7, 8, 9, 10},
		Label:  "Last",
	}

	spec := RidgelineSpec{
		Data:       []*RidgelineData{ridge1, ridge2},
		Width:      400,
		Height:     300,
		Reverse:    true, // Reverse order
		ShowLabels: true,
	}

	svg := RenderRidgeline(spec)

	if svg == "" {
		t.Error("Expected non-empty SVG output")
	}

	// Both labels should be present
	if !strings.Contains(svg, "First") {
		t.Error("Expected 'First' label")
	}

	if !strings.Contains(svg, "Last") {
		t.Error("Expected 'Last' label")
	}
}

func TestRenderRidgeline_EmptyData(t *testing.T) {
	spec := RidgelineSpec{
		Data:   []*RidgelineData{},
		Width:  400,
		Height: 300,
	}

	svg := RenderRidgeline(spec)

	if svg != "" {
		t.Error("Expected empty SVG for empty data")
	}
}

func TestRidgelineFromGroups(t *testing.T) {
	groups := map[string][]float64{
		"Group A": {1, 2, 3, 4, 5},
		"Group B": {6, 7, 8, 9, 10},
	}

	colors := map[string]string{
		"Group A": "#4285f4",
		"Group B": "#ea4335",
	}

	ridges := RidgelineFromGroups(groups, colors)

	if len(ridges) != 2 {
		t.Errorf("Expected 2 ridges, got %d", len(ridges))
	}

	// Check that ridges were created with correct data
	foundA := false
	foundB := false

	for _, ridge := range ridges {
		if ridge.Label == "Group A" {
			foundA = true
			if len(ridge.Values) != 5 {
				t.Errorf("Expected 5 values for Group A, got %d", len(ridge.Values))
			}
			if ridge.Color != "#4285f4" {
				t.Errorf("Expected color #4285f4 for Group A, got %s", ridge.Color)
			}
		}
		if ridge.Label == "Group B" {
			foundB = true
			if len(ridge.Values) != 5 {
				t.Errorf("Expected 5 values for Group B, got %d", len(ridge.Values))
			}
		}
	}

	if !foundA {
		t.Error("Group A not found in ridges")
	}
	if !foundB {
		t.Error("Group B not found in ridges")
	}
}
