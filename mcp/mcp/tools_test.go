package mcp

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Helper function to create a test request
func createTestRequest(t *testing.T, name string, args map[string]interface{}) *mcp.CallToolRequest {
	t.Helper()
	argsJSON, err := json.Marshal(args)
	if err != nil {
		t.Fatalf("Failed to marshal args: %v", err)
	}

	return &mcp.CallToolRequest{
		Params: &mcp.CallToolParamsRaw{
			Name:      name,
			Arguments: argsJSON,
		},
	}
}

// Helper to truncate string for error messages
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// TestServerCreation tests that the server can be created
func TestServerCreation(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	if server == nil {
		t.Fatal("Server is nil")
	}
}

// TestBarChartTool tests the bar_chart tool
func TestBarChartTool(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	args := map[string]interface{}{
		"title": "Test Bar Chart",
		"data": []map[string]interface{}{
			{"label": "A", "value": 10.0},
			{"label": "B", "value": 20.0},
			{"label": "C", "value": 15.0},
		},
		"width":  800.0,
		"height": 400.0,
	}

	request := createTestRequest(t, "bar_chart", args)

	result, err := server.handleBarChart(context.Background(), request)
	if err != nil {
		t.Fatalf("handleBarChart failed: %v", err)
	}

	if len(result.Content) == 0 {
		t.Fatal("Result has no content")
	}

	textContent, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatal("Result content is not TextContent")
	}

	// Check for SVG content (wrapped in markdown code fence)
	if !strings.Contains(textContent.Text, "```svg") && !strings.Contains(textContent.Text, "<svg") {
		t.Errorf("Result does not contain SVG. Got: %s", truncate(textContent.Text, 200))
	}
}

// TestPieChartTool tests the pie_chart tool
func TestPieChartTool(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	args := map[string]interface{}{
		"title": "Test Pie Chart",
		"data": []map[string]interface{}{
			{"label": "A", "value": 30.0},
			{"label": "B", "value": 50.0},
			{"label": "C", "value": 20.0},
		},
		"donut":  false,
		"width":  600.0,
		"height": 600.0,
	}

	request := createTestRequest(t, "pie_chart", args)

	result, err := server.handlePieChart(context.Background(), request)
	if err != nil {
		t.Fatalf("handlePieChart failed: %v", err)
	}

	if len(result.Content) == 0 {
		t.Fatal("Result has no content")
	}

	textContent, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatal("Result content is not TextContent")
	}

	// Check for SVG content (wrapped in markdown code fence)
	if !strings.Contains(textContent.Text, "```svg") && !strings.Contains(textContent.Text, "<svg") {
		t.Errorf("Result does not contain SVG. Got: %s", truncate(textContent.Text, 200))
	}
}

// TestLineChartTool tests the line_chart tool
func TestLineChartTool(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	args := map[string]interface{}{
		"title": "Test Line Chart",
		"series": []map[string]interface{}{
			{
				"name": "Series 1",
				"data": []map[string]interface{}{
					{"x": 0.0, "y": 10.0},
					{"x": 1.0, "y": 20.0},
					{"x": 2.0, "y": 15.0},
				},
			},
		},
		"width":  900.0,
		"height": 500.0,
	}

	request := createTestRequest(t, "line_chart", args)

	result, err := server.handleLineChart(context.Background(), request)
	if err != nil {
		t.Fatalf("handleLineChart failed: %v", err)
	}

	if len(result.Content) == 0 {
		t.Fatal("Result has no content")
	}

	textContent, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatal("Result content is not TextContent")
	}

	// Check for SVG content (wrapped in markdown code fence)
	if !strings.Contains(textContent.Text, "```svg") && !strings.Contains(textContent.Text, "<svg") {
		t.Errorf("Result does not contain SVG. Got: %s", truncate(textContent.Text, 200))
	}
}

// TestScatterPlotTool tests the scatter_plot tool
func TestScatterPlotTool(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	args := map[string]interface{}{
		"title": "Test Scatter Plot",
		"data": []map[string]interface{}{
			{"x": 1.0, "y": 2.0},
			{"x": 2.0, "y": 4.0},
			{"x": 3.0, "y": 3.5},
		},
		"width":  800.0,
		"height": 600.0,
	}

	request := createTestRequest(t, "scatter_plot", args)

	result, err := server.handleScatterPlot(context.Background(), request)
	if err != nil {
		t.Fatalf("handleScatterPlot failed: %v", err)
	}

	if len(result.Content) == 0 {
		t.Fatal("Result has no content")
	}

	textContent, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatal("Result content is not TextContent")
	}

	// Check for SVG content (wrapped in markdown code fence)
	if !strings.Contains(textContent.Text, "```svg") && !strings.Contains(textContent.Text, "<svg") {
		t.Errorf("Result does not contain SVG. Got: %s", truncate(textContent.Text, 200))
	}
}

// TestHeatmapTool tests the heatmap tool
func TestHeatmapTool(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	args := map[string]interface{}{
		"title": "Test Heatmap",
		"data": map[string]interface{}{
			"rows":    []string{"Row 1", "Row 2"},
			"columns": []string{"Col 1", "Col 2"},
			"values": [][]float64{
				{1.0, 2.0},
				{3.0, 4.0},
			},
		},
		"show_value": true,
		"width":      800.0,
		"height":     600.0,
	}

	request := createTestRequest(t, "heatmap", args)

	result, err := server.handleHeatmap(context.Background(), request)
	if err != nil {
		t.Fatalf("handleHeatmap failed: %v", err)
	}

	if len(result.Content) == 0 {
		t.Fatal("Result has no content")
	}

	textContent, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatal("Result content is not TextContent")
	}

	// Check for SVG content (wrapped in markdown code fence)
	if !strings.Contains(textContent.Text, "```svg") && !strings.Contains(textContent.Text, "<svg") {
		t.Errorf("Result does not contain SVG. Got: %s", truncate(textContent.Text, 200))
	}
}

// TestTreemapTool tests the treemap tool
func TestTreemapTool(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	args := map[string]interface{}{
		"title": "Test Treemap",
		"data": map[string]interface{}{
			"name": "root",
			"children": []map[string]interface{}{
				{
					"name":  "Group A",
					"value": 100.0,
				},
				{
					"name":  "Group B",
					"value": 200.0,
				},
			},
		},
		"show_labels": true,
		"width":       800.0,
		"height":      600.0,
	}

	request := createTestRequest(t, "treemap", args)

	result, err := server.handleTreemap(context.Background(), request)
	if err != nil {
		t.Fatalf("handleTreemap failed: %v", err)
	}

	if len(result.Content) == 0 {
		t.Fatal("Result has no content")
	}

	textContent, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatal("Result content is not TextContent")
	}

	// Check for SVG content (wrapped in markdown code fence)
	if !strings.Contains(textContent.Text, "```svg") && !strings.Contains(textContent.Text, "<svg") {
		t.Errorf("Result does not contain SVG. Got: %s", truncate(textContent.Text, 200))
	}
}

// TestBoxplotTool tests the boxplot tool
func TestBoxplotTool(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	args := map[string]interface{}{
		"title": "Test Boxplot",
		"data": []map[string]interface{}{
			{
				"label":  "Group A",
				"values": []float64{10, 15, 20, 25, 30, 35, 40},
			},
		},
		"show_outliers": true,
		"width":         800.0,
		"height":        600.0,
	}

	request := createTestRequest(t, "boxplot", args)

	result, err := server.handleBoxplot(context.Background(), request)
	if err != nil {
		t.Fatalf("handleBoxplot failed: %v", err)
	}

	if len(result.Content) == 0 {
		t.Fatal("Result has no content")
	}

	textContent, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatal("Result content is not TextContent")
	}

	// Check for SVG content (wrapped in markdown code fence)
	if !strings.Contains(textContent.Text, "```svg") && !strings.Contains(textContent.Text, "<svg") {
		t.Errorf("Result does not contain SVG. Got: %s", truncate(textContent.Text, 200))
	}
}

// TestHistogramTool tests the histogram tool
func TestHistogramTool(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	args := map[string]interface{}{
		"title":  "Test Histogram",
		"values": []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		"bins":   5.0,
		"width":  800.0,
		"height": 600.0,
	}

	request := createTestRequest(t, "histogram", args)

	result, err := server.handleHistogram(context.Background(), request)
	if err != nil {
		t.Fatalf("handleHistogram failed: %v", err)
	}

	if len(result.Content) == 0 {
		t.Fatal("Result has no content")
	}

	textContent, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatal("Result content is not TextContent")
	}

	// Check for SVG content (wrapped in markdown code fence)
	if !strings.Contains(textContent.Text, "```svg") && !strings.Contains(textContent.Text, "<svg") {
		t.Errorf("Result does not contain SVG. Got: %s", truncate(textContent.Text, 200))
	}
}

// TestCandlestickTool tests the candlestick tool
func TestCandlestickTool(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	args := map[string]interface{}{
		"title": "Test Candlestick",
		"data": []map[string]interface{}{
			{
				"date":   "2024-01-01",
				"open":   100.0,
				"high":   110.0,
				"low":    95.0,
				"close":  105.0,
				"volume": 1000000.0,
			},
			{
				"date":   "2024-01-02",
				"open":   105.0,
				"high":   115.0,
				"low":    103.0,
				"close":  112.0,
				"volume": 1200000.0,
			},
		},
		"show_volume": true,
		"width":       1000.0,
		"height":      600.0,
	}

	request := createTestRequest(t, "candlestick", args)

	result, err := server.handleCandlestick(context.Background(), request)
	if err != nil {
		t.Fatalf("handleCandlestick failed: %v", err)
	}

	if len(result.Content) == 0 {
		t.Fatal("Result has no content")
	}

	textContent, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatal("Result content is not TextContent")
	}

	// Check for SVG content (wrapped in markdown code fence)
	if !strings.Contains(textContent.Text, "```svg") && !strings.Contains(textContent.Text, "<svg") {
		t.Errorf("Result does not contain SVG. Got: %s", truncate(textContent.Text, 200))
	}
}

// TestParseArguments tests the parseArguments helper
func TestParseArguments(t *testing.T) {
	tests := []struct {
		name    string
		args    interface{}
		wantErr bool
	}{
		{
			name: "valid simple args",
			args: map[string]interface{}{
				"width":  800.0,
				"height": 600.0,
			},
			wantErr: false,
		},
		{
			name: "valid nested args",
			args: map[string]interface{}{
				"data": []map[string]interface{}{
					{"label": "A", "value": 10.0},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var target map[string]interface{}
			err := parseArguments(tt.args, &target)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseArguments() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestInvalidInputs tests error handling for invalid inputs
func TestInvalidInputs(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	tests := []struct {
		name    string
		handler func(context.Context, *mcp.CallToolRequest) (*mcp.CallToolResult, error)
		args    map[string]interface{}
	}{
		{
			name:    "bar chart with invalid data type",
			handler: server.handleBarChart,
			args: map[string]interface{}{
				"data": "not an array",
			},
		},
		{
			name:    "candlestick with no data",
			handler: server.handleCandlestick,
			args:    map[string]interface{}{
				"data": []interface{}{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			argsJSON, _ := json.Marshal(tt.args)
			request := &mcp.CallToolRequest{
				Params: &mcp.CallToolParamsRaw{
					Arguments: argsJSON,
				},
			}

			_, err := tt.handler(context.Background(), request)
			if err == nil {
				t.Error("Expected error for invalid input, got nil")
			}
		})
	}
}

// TestDefaultValues tests that default values are applied correctly
func TestDefaultValues(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Test bar chart defaults
	args := map[string]interface{}{
		"data": []map[string]interface{}{
			{"label": "A", "value": 10.0},
		},
		// Omit width, height, orientation to test defaults
	}

	request := createTestRequest(t, "bar_chart", args)

	result, err := server.handleBarChart(context.Background(), request)
	if err != nil {
		t.Fatalf("handleBarChart failed: %v", err)
	}

	if len(result.Content) == 0 {
		t.Fatal("Result has no content")
	}

	textContent, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatal("Result content is not TextContent")
	}

	// Should still generate valid SVG even without explicit dimensions
	// Check for SVG content (wrapped in markdown code fence)
	if !strings.Contains(textContent.Text, "```svg") && !strings.Contains(textContent.Text, "<svg") {
		t.Errorf("Result does not contain SVG. Got: %s", truncate(textContent.Text, 200))
	}
}

// Helper to convert args to JSON and back (simulating MCP protocol)
func roundtripJSON(v interface{}) (interface{}, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var result interface{}
	err = json.Unmarshal(data, &result)
	return result, err
}
