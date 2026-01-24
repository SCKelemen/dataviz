package charts

import (
	"strings"
	"testing"

	"github.com/SCKelemen/dataviz/scales"
	"github.com/SCKelemen/units"
)

func TestCalculateBoxPlotStats(t *testing.T) {
	values := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	stats := CalculateBoxPlotStats(values, 1.5)

	if stats.Median != 5.5 {
		t.Errorf("Expected median 5.5, got %f", stats.Median)
	}

	if stats.Q1 != 3.25 {
		t.Errorf("Expected Q1 3.25, got %f", stats.Q1)
	}

	if stats.Q3 != 7.75 {
		t.Errorf("Expected Q3 7.75, got %f", stats.Q3)
	}

	if stats.IQR != 4.5 {
		t.Errorf("Expected IQR 4.5, got %f", stats.IQR)
	}

	if stats.Mean != 5.5 {
		t.Errorf("Expected mean 5.5, got %f", stats.Mean)
	}
}

func TestCalculateBoxPlotStatsWithOutliers(t *testing.T) {
	values := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 50}
	stats := CalculateBoxPlotStats(values, 1.5)

	if len(stats.Outliers) == 0 {
		t.Error("Expected outliers to be detected")
	}

	// 50 should be an outlier
	found := false
	for _, outlier := range stats.Outliers {
		if outlier == 50 {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected 50 to be detected as an outlier")
	}
}

func TestPercentile(t *testing.T) {
	values := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	tests := []struct {
		percentile float64
		expected   float64
	}{
		{0, 1.0},
		{25, 3.25},
		{50, 5.5},
		{75, 7.75},
		{100, 10.0},
	}

	for _, test := range tests {
		result := percentile(values, test.percentile)
		if result != test.expected {
			t.Errorf("Percentile %f: expected %f, got %f", test.percentile, test.expected, result)
		}
	}
}

func TestRenderVerticalBoxPlot(t *testing.T) {
	data := []*BoxPlotData{
		{
			Values: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			Label:  "Group A",
			Color:  "#4285f4",
		},
		{
			Values: []float64{2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
			Label:  "Group B",
			Color:  "#ea4335",
		},
	}

	spec := BoxPlotSpec{
		Data:         data,
		Width:        400,
		Height:       300,
		ShowOutliers: true,
		ShowMean:     true,
	}

	svg := RenderVerticalBoxPlot(spec)

	if svg == "" {
		t.Error("Expected non-empty SVG output")
	}

	// Check for expected SVG elements
	if !strings.Contains(svg, "<line") {
		t.Error("Expected whisker lines")
	}

	if !strings.Contains(svg, "<rect") {
		t.Error("Expected box rectangles")
	}

	if !strings.Contains(svg, "Group A") {
		t.Error("Expected Group A label")
	}

	if !strings.Contains(svg, "Group B") {
		t.Error("Expected Group B label")
	}
}

func TestCalculateViolinStats(t *testing.T) {
	values := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	stats := CalculateViolinStats(values, 0) // Auto bandwidth

	if stats.Median != 5.5 {
		t.Errorf("Expected median 5.5, got %f", stats.Median)
	}

	if stats.Mean != 5.5 {
		t.Errorf("Expected mean 5.5, got %f", stats.Mean)
	}

	if len(stats.Density) == 0 {
		t.Error("Expected density points to be calculated")
	}

	// Check that density points cover the data range
	if len(stats.Density) > 0 {
		firstDensity := stats.Density[0]
		lastDensity := stats.Density[len(stats.Density)-1]

		if firstDensity.Value > 1.0 {
			t.Errorf("Expected first density point to be near min value, got %f", firstDensity.Value)
		}

		if lastDensity.Value < 10.0 {
			t.Errorf("Expected last density point to be near max value, got %f", lastDensity.Value)
		}
	}
}

func TestCalculateKDE(t *testing.T) {
	values := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	density := calculateKDE(values, 1.0)

	if len(density) == 0 {
		t.Error("Expected density points")
	}

	// Density values should sum to approximately 1 (when multiplied by step size)
	// This is a property of KDE

	// All density values should be non-negative
	for _, dp := range density {
		if dp.Density < 0 {
			t.Errorf("Expected non-negative density, got %f", dp.Density)
		}
	}
}

func TestRenderViolinPlot(t *testing.T) {
	data := []*ViolinPlotData{
		{
			Values: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			Label:  "Group A",
			Color:  "#4285f4",
		},
		{
			Values: []float64{2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
			Label:  "Group B",
			Color:  "#ea4335",
		},
	}

	spec := ViolinPlotSpec{
		Data:       data,
		Width:      400,
		Height:     300,
		ShowBox:    true,
		ShowMedian: true,
		ShowMean:   true,
	}

	svg := RenderViolinPlot(spec)

	if svg == "" {
		t.Error("Expected non-empty SVG output")
	}

	// Check for expected SVG elements
	if !strings.Contains(svg, "<path") {
		t.Error("Expected violin path elements")
	}

	if !strings.Contains(svg, "Group A") {
		t.Error("Expected Group A label")
	}

	if !strings.Contains(svg, "Group B") {
		t.Error("Expected Group B label")
	}
}

func TestRenderErrorBars(t *testing.T) {
	bars := []ErrorBar{
		{X: 1.0, Y: 5.0, ErrorLower: 0.5, ErrorUpper: 0.5, IsRelative: true},
		{X: 2.0, Y: 7.0, ErrorLower: 0.8, ErrorUpper: 0.8, IsRelative: true},
		{X: 3.0, Y: 6.0, ErrorLower: 0.3, ErrorUpper: 0.6, IsRelative: true},
	}

	spec := ErrorBarSpec{
		Bars:     bars,
		Color:    "#666",
		CapWidth: 8,
		CapStyle: CapStyleLine,
	}

	xScale := scales.NewLinearScale([2]float64{0, 4}, [2]units.Length{units.Px(40), units.Px(360)})
	yScale := scales.NewLinearScale([2]float64{0, 10}, [2]units.Length{units.Px(260), units.Px(40)})

	svg := RenderErrorBars(spec, xScale, yScale)

	if svg == "" {
		t.Error("Expected non-empty SVG output")
	}

	// Check for expected SVG elements
	if !strings.Contains(svg, "<line") {
		t.Error("Expected error bar lines")
	}
}

func TestRenderConfidenceBands(t *testing.T) {
	band := &ConfidenceBand{
		XValues:      []float64{1, 2, 3, 4, 5},
		YCenters:     []float64{5, 6, 7, 6, 5},
		YLowerBounds: []float64{4, 5, 6, 5, 4},
		YUpperBounds: []float64{6, 7, 8, 7, 6},
		Color:        "#4285f4",
		Opacity:      0.2,
	}

	spec := ConfidenceBandSpec{
		Bands: []*ConfidenceBand{band},
	}

	xScale := scales.NewLinearScale([2]float64{0, 6}, [2]units.Length{units.Px(40), units.Px(360)})
	yScale := scales.NewLinearScale([2]float64{0, 10}, [2]units.Length{units.Px(260), units.Px(40)})

	svg := RenderConfidenceBands(spec, xScale, yScale)

	if svg == "" {
		t.Error("Expected non-empty SVG output")
	}

	// Check for expected SVG elements
	if !strings.Contains(svg, "<path") {
		t.Error("Expected confidence band path")
	}

	// Should contain filled area and center line
	pathCount := strings.Count(svg, "<path")
	if pathCount < 2 {
		t.Errorf("Expected at least 2 paths (filled band + center line), got %d", pathCount)
	}
}

func TestErrorBarCapStyles(t *testing.T) {
	bars := []ErrorBar{
		{X: 1.0, Y: 5.0, ErrorLower: 0.5, ErrorUpper: 0.5, IsRelative: true},
	}

	xScale := scales.NewLinearScale([2]float64{0, 2}, [2]units.Length{units.Px(40), units.Px(360)})
	yScale := scales.NewLinearScale([2]float64{0, 10}, [2]units.Length{units.Px(260), units.Px(40)})

	capStyles := []CapStyle{CapStyleLine, CapStyleCircle, CapStyleNone}

	for _, capStyle := range capStyles {
		spec := ErrorBarSpec{
			Bars:     bars,
			CapStyle: capStyle,
		}

		svg := RenderErrorBars(spec, xScale, yScale)

		if svg == "" {
			t.Errorf("Expected non-empty SVG output for cap style %s", capStyle)
		}

		// Check for appropriate elements based on cap style
		switch capStyle {
		case CapStyleLine:
			if !strings.Contains(svg, "<line") {
				t.Errorf("Expected line caps for style %s", capStyle)
			}
		case CapStyleCircle:
			if !strings.Contains(svg, "<circle") {
				t.Errorf("Expected circle caps for style %s", capStyle)
			}
		}
	}
}
