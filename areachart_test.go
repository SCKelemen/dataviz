package dataviz

import (
	"strings"
	"testing"
	"time"

	design "github.com/SCKelemen/design-system"
)

func TestRenderAreaChart(t *testing.T) {
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	points := []TimeSeriesData{
		{Date: startDate, Value: 100},
		{Date: startDate.AddDate(0, 0, 1), Value: 150},
		{Date: startDate.AddDate(0, 0, 2), Value: 125},
		{Date: startDate.AddDate(0, 0, 3), Value: 175},
		{Date: startDate.AddDate(0, 0, 4), Value: 200},
	}

	data := AreaChartData{
		Points:    points,
		Color:     "#10B981",
		FillColor: "#10B981",
	}

	tokens := design.DefaultTheme()
	result := RenderAreaChart(data, 0, 0, 400, 200, tokens)

	// Check that SVG is generated
	if result == "" {
		t.Error("Expected non-empty SVG output")
	}

	// Check for SVG group
	if !strings.Contains(result, "<g transform") {
		t.Error("Expected <g> tag with transform")
	}

	// Check for filled path
	if !strings.Contains(result, "<path") {
		t.Error("Expected <path> element for filled area")
	}

	// Area charts should have fill
	if !strings.Contains(result, "fill=") {
		t.Error("Expected fill attribute for area chart")
	}
}

func TestRenderAreaChart_WithSmoothing(t *testing.T) {
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	points := []TimeSeriesData{
		{Date: startDate, Value: 100},
		{Date: startDate.AddDate(0, 0, 1), Value: 150},
		{Date: startDate.AddDate(0, 0, 2), Value: 125},
		{Date: startDate.AddDate(0, 0, 3), Value: 175},
	}

	data := AreaChartData{
		Points:    points,
		Color:     "#10B981",
		FillColor: "#10B981",
		Smooth:    true,
		Tension:   0.3,
	}

	tokens := design.DefaultTheme()
	result := RenderAreaChart(data, 0, 0, 400, 200, tokens)

	// Check that SVG is generated
	if result == "" {
		t.Error("Expected non-empty SVG output")
	}

	// Smooth curves use C (cubic bezier) commands
	if !strings.Contains(result, " C ") {
		t.Error("Expected smooth area to contain C (cubic bezier) commands")
	}
}

func TestRenderAreaChart_WithBorder(t *testing.T) {
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	points := []TimeSeriesData{
		{Date: startDate, Value: 100},
		{Date: startDate.AddDate(0, 0, 1), Value: 150},
		{Date: startDate.AddDate(0, 0, 2), Value: 125},
	}

	data := AreaChartData{
		Points:    points,
		Color:     "#10B981",
		FillColor: "rgba(16, 185, 129, 0.3)",
	}

	tokens := design.DefaultTheme()
	result := RenderAreaChart(data, 0, 0, 400, 200, tokens)

	// When Color is specified, should have a border line
	if data.Color != "" {
		pathCount := strings.Count(result, "<path")
		if pathCount < 2 {
			t.Errorf("Expected at least 2 paths (fill + border), got %d", pathCount)
		}
	}
}

func TestRenderAreaChart_EmptyData(t *testing.T) {
	data := AreaChartData{
		Points:    []TimeSeriesData{},
		Color:     "#10B981",
		FillColor: "#10B981",
	}

	tokens := design.DefaultTheme()
	result := RenderAreaChart(data, 0, 0, 400, 200, tokens)

	// Should return empty string for empty data
	if result != "" {
		t.Error("Expected empty string for empty data")
	}
}

func TestRenderAreaChart_SinglePoint(t *testing.T) {
	data := AreaChartData{
		Points: []TimeSeriesData{
			{Date: time.Now(), Value: 100},
		},
		Color:     "#10B981",
		FillColor: "#10B981",
	}

	tokens := design.DefaultTheme()
	result := RenderAreaChart(data, 0, 0, 400, 200, tokens)

	// With one point, should not crash
	if result == "" {
		t.Error("Expected some output for single point")
	}
}

func BenchmarkRenderAreaChart(b *testing.B) {
	startDate := time.Now()
	points := make([]TimeSeriesData, 100)
	for i := range points {
		points[i] = TimeSeriesData{
			Date:  startDate.AddDate(0, 0, i),
			Value: 100 + (i % 50),
		}
	}

	data := AreaChartData{
		Points:    points,
		Color:     "#10B981",
		FillColor: "#10B981",
	}

	tokens := design.DefaultTheme()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RenderAreaChart(data, 0, 0, 800, 400, tokens)
	}
}

func BenchmarkRenderAreaChart_Smooth(b *testing.B) {
	startDate := time.Now()
	points := make([]TimeSeriesData, 100)
	for i := range points {
		points[i] = TimeSeriesData{
			Date:  startDate.AddDate(0, 0, i),
			Value: 100 + (i % 50),
		}
	}

	data := AreaChartData{
		Points:    points,
		Color:     "#10B981",
		FillColor: "#10B981",
		Smooth:    true,
		Tension:   0.3,
	}

	tokens := design.DefaultTheme()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RenderAreaChart(data, 0, 0, 800, 400, tokens)
	}
}
