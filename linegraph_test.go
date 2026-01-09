package dataviz

import (
	"strings"
	"testing"
	"time"

	design "github.com/SCKelemen/design-system"
)

func TestRenderLineGraph(t *testing.T) {
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	points := []TimeSeriesData{
		{Date: startDate, Value: 100},
		{Date: startDate.AddDate(0, 0, 1), Value: 150},
		{Date: startDate.AddDate(0, 0, 2), Value: 125},
		{Date: startDate.AddDate(0, 0, 3), Value: 175},
		{Date: startDate.AddDate(0, 0, 4), Value: 200},
	}

	data := LineGraphData{
		Points: points,
		Color:  "#3B82F6",
	}

	tokens := design.DefaultTheme()
	result := RenderLineGraph(data, 0, 0, 400, 200, tokens)

	// Check that SVG is generated
	if result == "" {
		t.Error("Expected non-empty SVG output")
	}

	// Check for SVG group
	if !strings.Contains(result, "<g transform") {
		t.Error("Expected <g> tag with transform")
	}

	// Check for path element (the line)
	if !strings.Contains(result, "<path") {
		t.Error("Expected <path> element for line")
	}
}

func TestRenderLineGraph_WithSmoothing(t *testing.T) {
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	points := []TimeSeriesData{
		{Date: startDate, Value: 100},
		{Date: startDate.AddDate(0, 0, 1), Value: 150},
		{Date: startDate.AddDate(0, 0, 2), Value: 125},
		{Date: startDate.AddDate(0, 0, 3), Value: 175},
	}

	data := LineGraphData{
		Points:  points,
		Color:   "#3B82F6",
		Smooth:  true,
		Tension: 0.3,
	}

	tokens := design.DefaultTheme()
	result := RenderLineGraph(data, 0, 0, 400, 200, tokens)

	// Check that SVG is generated
	if result == "" {
		t.Error("Expected non-empty SVG output")
	}

	// Smooth curves use C (cubic bezier) commands
	if !strings.Contains(result, " C ") {
		t.Error("Expected smooth curve to contain C (cubic bezier) commands")
	}
}

func TestRenderLineGraph_WithMarkers(t *testing.T) {
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	points := []TimeSeriesData{
		{Date: startDate, Value: 100},
		{Date: startDate.AddDate(0, 0, 1), Value: 150},
		{Date: startDate.AddDate(0, 0, 2), Value: 125},
	}

	tests := []struct {
		markerType     string
		expectedShape  string
	}{
		{"circle", "<circle"},
		{"square", "<rect"},
		{"diamond", "<polygon"},
		{"triangle", "<polygon"},
	}

	tokens := design.DefaultTheme()

	for _, tt := range tests {
		t.Run(tt.markerType, func(t *testing.T) {
			data := LineGraphData{
				Points:     points,
				Color:      "#3B82F6",
				MarkerType: tt.markerType,
				MarkerSize: 4,
			}

			result := RenderLineGraph(data, 0, 0, 400, 200, tokens)

			// Check for marker shapes
			if !strings.Contains(result, tt.expectedShape) {
				t.Errorf("Expected %s marker to contain %s", tt.markerType, tt.expectedShape)
			}
		})
	}
}

func TestRenderLineGraph_WithFill(t *testing.T) {
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	points := []TimeSeriesData{
		{Date: startDate, Value: 100},
		{Date: startDate.AddDate(0, 0, 1), Value: 150},
		{Date: startDate.AddDate(0, 0, 2), Value: 125},
	}

	data := LineGraphData{
		Points:    points,
		Color:     "#3B82F6",
		FillColor: "rgba(59, 130, 246, 0.1)",
	}

	tokens := design.DefaultTheme()
	result := RenderLineGraph(data, 0, 0, 400, 200, tokens)

	// Check for filled area (should have two paths: one for fill, one for line)
	pathCount := strings.Count(result, "<path")
	if pathCount < 2 {
		t.Errorf("Expected at least 2 paths (fill + line), got %d", pathCount)
	}
}

func TestRenderLineGraph_EmptyData(t *testing.T) {
	data := LineGraphData{
		Points: []TimeSeriesData{},
		Color:  "#3B82F6",
	}

	tokens := design.DefaultTheme()
	result := RenderLineGraph(data, 0, 0, 400, 200, tokens)

	// Should return empty string for empty data
	if result != "" {
		t.Error("Expected empty string for empty data")
	}
}

func TestRenderLineGraph_SinglePoint(t *testing.T) {
	data := LineGraphData{
		Points: []TimeSeriesData{
			{Date: time.Now(), Value: 100},
		},
		Color: "#3B82F6",
	}

	tokens := design.DefaultTheme()
	result := RenderLineGraph(data, 0, 0, 400, 200, tokens)

	// With one point, should not crash (may or may not render a line)
	// The important part is it doesn't panic
	if result == "" {
		t.Error("Expected some output for single point")
	}
}

func BenchmarkRenderLineGraph(b *testing.B) {
	startDate := time.Now()
	points := make([]TimeSeriesData, 100)
	for i := range points {
		points[i] = TimeSeriesData{
			Date:  startDate.AddDate(0, 0, i),
			Value: 100 + (i % 50),
		}
	}

	data := LineGraphData{
		Points: points,
		Color:  "#3B82F6",
	}

	tokens := design.DefaultTheme()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RenderLineGraph(data, 0, 0, 800, 400, tokens)
	}
}

func BenchmarkRenderLineGraph_Smooth(b *testing.B) {
	startDate := time.Now()
	points := make([]TimeSeriesData, 100)
	for i := range points {
		points[i] = TimeSeriesData{
			Date:  startDate.AddDate(0, 0, i),
			Value: 100 + (i % 50),
		}
	}

	data := LineGraphData{
		Points:  points,
		Color:   "#3B82F6",
		Smooth:  true,
		Tension: 0.3,
	}

	tokens := design.DefaultTheme()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RenderLineGraph(data, 0, 0, 800, 400, tokens)
	}
}
