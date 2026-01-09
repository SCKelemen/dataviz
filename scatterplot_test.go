package dataviz

import (
	"strings"
	"testing"
	"time"

	design "github.com/SCKelemen/design-system"
)

func TestRenderScatterPlot(t *testing.T) {
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	points := []ScatterPoint{
		{Date: startDate, Value: 100, Size: 0, Label: ""},
		{Date: startDate.AddDate(0, 0, 1), Value: 150, Size: 0, Label: ""},
		{Date: startDate.AddDate(0, 0, 2), Value: 125, Size: 0, Label: ""},
		{Date: startDate.AddDate(0, 0, 3), Value: 175, Size: 8, Label: "Peak"},
		{Date: startDate.AddDate(0, 0, 4), Value: 200, Size: 0, Label: ""},
	}

	data := ScatterPlotData{
		Points:     points,
		Color:      "#F59E0B",
		MarkerType: "circle",
		MarkerSize: 5,
	}

	tokens := design.DefaultTheme()
	result := RenderScatterPlot(data, 0, 0, 400, 200, tokens)

	// Check that SVG is generated
	if result == "" {
		t.Error("Expected non-empty SVG output")
	}

	// Check for SVG group
	if !strings.Contains(result, "<g transform") {
		t.Error("Expected <g> tag with transform")
	}

	// Check for markers (circles in this case)
	if !strings.Contains(result, "<circle") {
		t.Error("Expected <circle> elements for scatter points")
	}

	// Check for labeled point
	if !strings.Contains(result, "Peak") {
		t.Error("Expected label 'Peak' in output")
	}
}

func TestRenderScatterPlot_MarkerTypes(t *testing.T) {
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	points := []ScatterPoint{
		{Date: startDate, Value: 100},
		{Date: startDate.AddDate(0, 0, 1), Value: 150},
		{Date: startDate.AddDate(0, 0, 2), Value: 125},
	}

	tests := []struct {
		markerType    string
		expectedShape string
	}{
		{"circle", "<circle"},
		{"dot", "<circle"},
		{"square", "<rect"},
		{"diamond", "<polygon"},
		{"triangle", "<polygon"},
		{"cross", "<line"},
		{"x", "<line"},
	}

	tokens := design.DefaultTheme()

	for _, tt := range tests {
		t.Run(tt.markerType, func(t *testing.T) {
			data := ScatterPlotData{
				Points:     points,
				Color:      "#F59E0B",
				MarkerType: tt.markerType,
				MarkerSize: 5,
			}

			result := RenderScatterPlot(data, 0, 0, 400, 200, tokens)

			// Check for marker shapes
			if !strings.Contains(result, tt.expectedShape) {
				t.Errorf("Expected %s marker to contain %s", tt.markerType, tt.expectedShape)
			}
		})
	}
}

func TestRenderScatterPlot_CustomPointSizes(t *testing.T) {
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	points := []ScatterPoint{
		{Date: startDate, Value: 100, Size: 3},
		{Date: startDate.AddDate(0, 0, 1), Value: 150, Size: 8},
		{Date: startDate.AddDate(0, 0, 2), Value: 125, Size: 12},
	}

	data := ScatterPlotData{
		Points:     points,
		Color:      "#F59E0B",
		MarkerType: "circle",
		MarkerSize: 5,
	}

	tokens := design.DefaultTheme()
	result := RenderScatterPlot(data, 0, 0, 400, 200, tokens)

	// Check that different sized circles are rendered
	// Each custom size should override the default
	if !strings.Contains(result, "r=\"3.00\"") {
		t.Error("Expected circle with radius 3.00")
	}
	if !strings.Contains(result, "r=\"8.00\"") {
		t.Error("Expected circle with radius 8.00")
	}
	if !strings.Contains(result, "r=\"12.00\"") {
		t.Error("Expected circle with radius 12.00")
	}
}

func TestRenderScatterPlot_WithLabels(t *testing.T) {
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	points := []ScatterPoint{
		{Date: startDate, Value: 100, Label: "Start"},
		{Date: startDate.AddDate(0, 0, 1), Value: 150, Label: ""},
		{Date: startDate.AddDate(0, 0, 2), Value: 175, Label: "End"},
	}

	data := ScatterPlotData{
		Points:     points,
		Color:      "#F59E0B",
		MarkerType: "circle",
		MarkerSize: 5,
	}

	tokens := design.DefaultTheme()
	result := RenderScatterPlot(data, 0, 0, 400, 200, tokens)

	// Check for text labels
	if !strings.Contains(result, "Start") {
		t.Error("Expected label 'Start' in output")
	}
	if !strings.Contains(result, "End") {
		t.Error("Expected label 'End' in output")
	}

	// Count text elements (only 2 points have labels)
	textCount := strings.Count(result, "<text")
	// May have grid labels plus point labels, so check minimum
	if textCount < 2 {
		t.Errorf("Expected at least 2 text elements for point labels, got %d", textCount)
	}
}

func TestRenderScatterPlot_EmptyData(t *testing.T) {
	data := ScatterPlotData{
		Points:     []ScatterPoint{},
		Color:      "#F59E0B",
		MarkerType: "circle",
		MarkerSize: 5,
	}

	tokens := design.DefaultTheme()
	result := RenderScatterPlot(data, 0, 0, 400, 200, tokens)

	// Should return empty string for empty data
	if result != "" {
		t.Error("Expected empty string for empty data")
	}
}

func TestRenderScatterPlot_SinglePoint(t *testing.T) {
	data := ScatterPlotData{
		Points: []ScatterPoint{
			{Date: time.Now(), Value: 100},
		},
		Color:      "#F59E0B",
		MarkerType: "circle",
		MarkerSize: 5,
	}

	tokens := design.DefaultTheme()
	result := RenderScatterPlot(data, 0, 0, 400, 200, tokens)

	// With one point, should render successfully
	if result == "" {
		t.Error("Expected output for single point")
	}

	// Should contain a marker
	if !strings.Contains(result, "<circle") {
		t.Error("Expected circle marker for single point")
	}
}

func BenchmarkRenderScatterPlot(b *testing.B) {
	startDate := time.Now()
	points := make([]ScatterPoint, 100)
	for i := range points {
		points[i] = ScatterPoint{
			Date:  startDate.AddDate(0, 0, i),
			Value: 100 + (i % 50),
		}
	}

	data := ScatterPlotData{
		Points:     points,
		Color:      "#F59E0B",
		MarkerType: "circle",
		MarkerSize: 5,
	}

	tokens := design.DefaultTheme()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RenderScatterPlot(data, 0, 0, 800, 400, tokens)
	}
}
