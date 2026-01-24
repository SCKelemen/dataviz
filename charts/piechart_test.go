package charts

import (
	"strings"
	"testing"
)

func TestRenderPieChart(t *testing.T) {
	data := PieChartData{
		Slices: []PieSlice{
			{Label: "Category A", Value: 30},
			{Label: "Category B", Value: 45},
			{Label: "Category C", Value: 25},
		},
	}

	result := RenderPieChart(data, 0, 0, 400, 400, "Test Pie Chart", false, true, true)

	// Check that SVG is generated
	if result == "" {
		t.Error("Expected non-empty SVG output")
	}

	// Check for SVG element
	if !strings.Contains(result, "<svg") {
		t.Error("Expected <svg> element")
	}

	// Check for path elements (pie slices)
	if !strings.Contains(result, "<path") {
		t.Error("Expected <path> elements for pie slices")
	}

	// Check for title
	if !strings.Contains(result, "Test Pie Chart") {
		t.Error("Expected title to be in output")
	}

	// Check for legend (we enabled it)
	if !strings.Contains(result, "Category A") {
		t.Error("Expected legend labels in output")
	}
}

func TestRenderPieChart_DonutMode(t *testing.T) {
	data := PieChartData{
		Slices: []PieSlice{
			{Label: "Category A", Value: 30},
			{Label: "Category B", Value: 45},
		},
	}

	result := RenderPieChart(data, 0, 0, 400, 400, "Donut Chart", true, false, false)

	// Check that SVG is generated
	if result == "" {
		t.Error("Expected non-empty SVG output")
	}

	// Donut charts have annular sectors (paths with inner and outer arcs)
	// The path should be more complex than a simple pie slice
	if !strings.Contains(result, "<path") {
		t.Error("Expected <path> elements for donut slices")
	}
}

func TestRenderPieChart_CustomColors(t *testing.T) {
	data := PieChartData{
		Slices: []PieSlice{
			{Label: "Red", Value: 50},
			{Label: "Blue", Value: 50},
		},
		Colors: []string{"#FF0000", "#0000FF"},
	}

	result := RenderPieChart(data, 0, 0, 400, 400, "", false, false, false)

	// Check for custom colors in output (case insensitive)
	resultLower := strings.ToLower(result)
	if !strings.Contains(resultLower, "#ff0000") && !strings.Contains(resultLower, "#f00") {
		t.Error("Expected custom red color in output")
	}
	if !strings.Contains(resultLower, "#0000ff") && !strings.Contains(resultLower, "#00f") {
		t.Error("Expected custom blue color in output")
	}
}

func TestRenderPieChart_EmptyData(t *testing.T) {
	data := PieChartData{
		Slices: []PieSlice{},
	}

	result := RenderPieChart(data, 0, 0, 400, 400, "Empty Chart", false, false, false)

	// Should return empty chart message
	if result == "" {
		t.Error("Expected non-empty output for empty data")
	}

	// Should contain "No data available" message
	if !strings.Contains(result, "No data available") {
		t.Error("Expected 'No data available' message for empty data")
	}
}

func TestRenderPieChart_ZeroTotal(t *testing.T) {
	data := PieChartData{
		Slices: []PieSlice{
			{Label: "Zero A", Value: 0},
			{Label: "Zero B", Value: 0},
		},
	}

	result := RenderPieChart(data, 0, 0, 400, 400, "Zero Total", false, false, false)

	// Should handle zero total gracefully
	if result == "" {
		t.Error("Expected non-empty output for zero total")
	}

	// Should show empty data message
	if !strings.Contains(result, "No data available") {
		t.Error("Expected 'No data available' message for zero total")
	}
}

func TestRenderPieChart_SingleSlice(t *testing.T) {
	data := PieChartData{
		Slices: []PieSlice{
			{Label: "Only One", Value: 100},
		},
	}

	result := RenderPieChart(data, 0, 0, 400, 400, "Single Slice", false, false, true)

	// Should render successfully
	if result == "" {
		t.Error("Expected non-empty output for single slice")
	}

	// Should contain path for the full circle
	if !strings.Contains(result, "<path") {
		t.Error("Expected <path> element for single slice")
	}

	// Should show 100% if percent labels enabled
	if !strings.Contains(result, "100.0%") {
		t.Error("Expected 100% label for single slice")
	}
}

func TestRenderPieChart_ManySlices(t *testing.T) {
	// Test with more slices than default colors
	data := PieChartData{
		Slices: []PieSlice{
			{Label: "S1", Value: 10},
			{Label: "S2", Value: 10},
			{Label: "S3", Value: 10},
			{Label: "S4", Value: 10},
			{Label: "S5", Value: 10},
			{Label: "S6", Value: 10},
			{Label: "S7", Value: 10},
			{Label: "S8", Value: 10},
			{Label: "S9", Value: 10},
			{Label: "S10", Value: 10},
			{Label: "S11", Value: 10}, // More than 10 (should cycle colors)
			{Label: "S12", Value: 10},
		},
	}

	result := RenderPieChart(data, 0, 0, 600, 600, "Many Slices", false, true, false)

	// Should render successfully with color cycling
	if result == "" {
		t.Error("Expected non-empty output for many slices")
	}

	// Should contain paths for all slices
	pathCount := strings.Count(result, "<path")
	if pathCount < 12 {
		t.Errorf("Expected at least 12 paths (one per slice), got %d", pathCount)
	}
}

func TestRenderPieChart_SmallSlices(t *testing.T) {
	// Test with very small slices (< 5%)
	data := PieChartData{
		Slices: []PieSlice{
			{Label: "Large", Value: 90},
			{Label: "Small", Value: 5},
			{Label: "Tiny", Value: 3},
			{Label: "Micro", Value: 2},
		},
	}

	result := RenderPieChart(data, 0, 0, 400, 400, "", false, false, true)

	// Small slices (< 5%) should not show percentage labels
	// Only "Large" (90%) and "Small" (5%) should have labels
	if !strings.Contains(result, "90.0%") {
		t.Error("Expected 90% label for large slice")
	}
	if !strings.Contains(result, "5.0%") {
		t.Error("Expected 5% label for small slice")
	}
	// Tiny (3%) and Micro (2%) should not have labels (< 5% threshold)
}

func BenchmarkRenderPieChart(b *testing.B) {
	data := PieChartData{
		Slices: []PieSlice{
			{Label: "Category A", Value: 30},
			{Label: "Category B", Value: 45},
			{Label: "Category C", Value: 25},
			{Label: "Category D", Value: 50},
			{Label: "Category E", Value: 20},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RenderPieChart(data, 0, 0, 400, 400, "Benchmark", false, true, true)
	}
}

func BenchmarkRenderPieChart_Donut(b *testing.B) {
	data := PieChartData{
		Slices: []PieSlice{
			{Label: "Category A", Value: 30},
			{Label: "Category B", Value: 45},
			{Label: "Category C", Value: 25},
			{Label: "Category D", Value: 50},
			{Label: "Category E", Value: 20},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RenderPieChart(data, 0, 0, 400, 400, "Benchmark Donut", true, true, true)
	}
}
