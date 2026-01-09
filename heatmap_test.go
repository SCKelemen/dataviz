package dataviz

import (
	"strings"
	"testing"
	"time"

	design "github.com/SCKelemen/design-system"
)

func TestRenderLinearHeatmap(t *testing.T) {
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	days := []ContributionDay{
		{Date: startDate, Count: 0},
		{Date: startDate.AddDate(0, 0, 1), Count: 5},
		{Date: startDate.AddDate(0, 0, 2), Count: 10},
		{Date: startDate.AddDate(0, 0, 3), Count: 15},
		{Date: startDate.AddDate(0, 0, 4), Count: 20},
	}

	data := HeatmapData{
		Days:      days,
		StartDate: startDate,
		EndDate:   startDate.AddDate(0, 0, 4),
		Type:      "linear",
	}

	tokens := design.DefaultTheme()
	result := RenderLinearHeatmap(data, 0, 0, 500, 100, "#3B82F6", tokens)

	// Check that SVG is generated
	if result == "" {
		t.Error("Expected non-empty SVG output")
	}

	// Check for SVG group
	if !strings.Contains(result, "<g transform") {
		t.Error("Expected <g> tag with transform")
	}

	// Check for rectangles (one per day)
	rectCount := strings.Count(result, "<rect")
	if rectCount != len(days) {
		t.Errorf("Expected %d rectangles, got %d", len(days), rectCount)
	}
}

func TestRenderLinearHeatmap_EmptyData(t *testing.T) {
	data := HeatmapData{
		Days:      []ContributionDay{},
		StartDate: time.Now(),
		EndDate:   time.Now(),
		Type:      "linear",
	}

	tokens := design.DefaultTheme()
	result := RenderLinearHeatmap(data, 0, 0, 500, 100, "#3B82F6", tokens)

	// Should return empty string for empty data
	if result != "" {
		t.Error("Expected empty string for empty data")
	}
}

func TestRenderWeeksHeatmap(t *testing.T) {
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	// Create 7 days worth of data (one week)
	days := make([]ContributionDay, 7)
	for i := 0; i < 7; i++ {
		days[i] = ContributionDay{
			Date:  startDate.AddDate(0, 0, i),
			Count: i * 2,
		}
	}

	data := HeatmapData{
		Days:      days,
		StartDate: startDate,
		EndDate:   startDate.AddDate(0, 0, 6),
		Type:      "weeks",
	}

	tokens := design.DefaultTheme()
	result := RenderWeeksHeatmap(data, 0, 0, 600, 150, "#3B82F6", tokens)

	// Check that SVG is generated
	if result == "" {
		t.Error("Expected non-empty SVG output")
	}

	// Check for SVG group
	if !strings.Contains(result, "<g transform") {
		t.Error("Expected <g> tag with transform")
	}

	// Check for rectangles
	rectCount := strings.Count(result, "<rect")
	if rectCount < 7 {
		t.Errorf("Expected at least 7 rectangles for 7 days, got %d", rectCount)
	}
}

func TestRenderWeeksHeatmap_EmptyData(t *testing.T) {
	data := HeatmapData{
		Days:      []ContributionDay{},
		StartDate: time.Now(),
		EndDate:   time.Now(),
		Type:      "weeks",
	}

	tokens := design.DefaultTheme()
	result := RenderWeeksHeatmap(data, 0, 0, 600, 150, "#3B82F6", tokens)

	// Should return empty string for empty data
	if result != "" {
		t.Error("Expected empty string for empty data")
	}
}

func TestAdjustColorForContribution(t *testing.T) {
	baseColor := "#3B82F6" // Blue

	tests := []struct {
		name  string
		ratio float64
	}{
		{"zero", 0.0},
		{"low", 0.25},
		{"medium", 0.5},
		{"high", 0.75},
		{"max", 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AdjustColorForContribution(baseColor, tt.ratio)

			// Should return a valid color string (starts with #)
			if !strings.HasPrefix(result, "#") {
				t.Errorf("Expected color to start with #, got: %s", result)
			}

			// Should be a 6-character hex color
			if len(result) != 7 {
				t.Errorf("Expected 7-character color (#RRGGBB), got: %s", result)
			}
		})
	}
}

func BenchmarkRenderLinearHeatmap(b *testing.B) {
	startDate := time.Now()
	days := make([]ContributionDay, 365)
	for i := range days {
		days[i] = ContributionDay{
			Date:  startDate.AddDate(0, 0, i),
			Count: i % 20,
		}
	}

	data := HeatmapData{
		Days:      days,
		StartDate: startDate,
		EndDate:   startDate.AddDate(0, 0, 364),
		Type:      "linear",
	}

	tokens := design.DefaultTheme()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RenderLinearHeatmap(data, 0, 0, 800, 100, "#3B82F6", tokens)
	}
}

func BenchmarkRenderWeeksHeatmap(b *testing.B) {
	startDate := time.Now()
	days := make([]ContributionDay, 365)
	for i := range days {
		days[i] = ContributionDay{
			Date:  startDate.AddDate(0, 0, i),
			Count: i % 20,
		}
	}

	data := HeatmapData{
		Days:      days,
		StartDate: startDate,
		EndDate:   startDate.AddDate(0, 0, 364),
		Type:      "weeks",
	}

	tokens := design.DefaultTheme()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RenderWeeksHeatmap(data, 0, 0, 800, 150, "#3B82F6", tokens)
	}
}
