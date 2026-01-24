package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/SCKelemen/dataviz/mcp/charts"
	"github.com/SCKelemen/dataviz/mcp/types"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterTools registers all chart generation tools
func (s *Server) RegisterTools() {
	// Tool: bar_chart
	s.server.AddTool(
		&mcp.Tool{
			Name:        "bar_chart",
			Description: "Generate a bar chart (vertical or horizontal) from labeled data points",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title": map[string]interface{}{
						"type":        "string",
						"description": "Chart title",
					},
					"data": map[string]interface{}{
						"type":        "array",
						"description": "Array of {label, value} objects",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"label": map[string]string{"type": "string"},
								"value": map[string]string{"type": "number"},
							},
							"required": []string{"label", "value"},
						},
					},
					"orientation": map[string]interface{}{
						"type":        "string",
						"description": "Bar orientation: 'vertical' or 'horizontal'",
						"default":     "vertical",
					},
					"color": map[string]interface{}{
						"type":        "string",
						"description": "Bar color (hex code)",
					},
					"width": map[string]interface{}{
						"type":        "number",
						"description": "Chart width in pixels",
						"default":     800,
					},
					"height": map[string]interface{}{
						"type":        "number",
						"description": "Chart height in pixels",
						"default":     400,
					},
				},
				"required": []string{"data"},
			},
		},
		s.handleBarChart,
	)

	// Tool: pie_chart
	s.server.AddTool(
		&mcp.Tool{
			Name:        "pie_chart",
			Description: "Generate a pie or donut chart from labeled data points",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title": map[string]interface{}{
						"type":        "string",
						"description": "Chart title",
					},
					"data": map[string]interface{}{
						"type":        "array",
						"description": "Array of {label, value} objects",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"label": map[string]string{"type": "string"},
								"value": map[string]string{"type": "number"},
							},
							"required": []string{"label", "value"},
						},
					},
					"donut": map[string]interface{}{
						"type":        "boolean",
						"description": "Make it a donut chart (hollow center)",
						"default":     false,
					},
					"width": map[string]interface{}{
						"type":    "number",
						"default": 600,
					},
					"height": map[string]interface{}{
						"type":    "number",
						"default": 600,
					},
				},
				"required": []string{"data"},
			},
		},
		s.handlePieChart,
	)

	// Tool: line_chart
	s.server.AddTool(
		&mcp.Tool{
			Name:        "line_chart",
			Description: "Generate a line chart with one or more data series",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title": map[string]interface{}{
						"type":        "string",
						"description": "Chart title",
					},
					"series": map[string]interface{}{
						"type":        "array",
						"description": "Array of series, each with name and data points",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"name": map[string]string{"type": "string"},
								"data": map[string]interface{}{
									"type": "array",
									"items": map[string]interface{}{
										"type": "object",
										"properties": map[string]interface{}{
											"x": map[string]interface{}{
												"description": "X value (number, string, or date)",
											},
											"y": map[string]string{"type": "number"},
										},
									},
								},
								"color": map[string]string{"type": "string"},
							},
						},
					},
					"x_label": map[string]interface{}{
						"type":        "string",
						"description": "X-axis label",
					},
					"y_label": map[string]interface{}{
						"type":        "string",
						"description": "Y-axis label",
					},
					"area": map[string]interface{}{
						"type":        "boolean",
						"description": "Fill area under line",
						"default":     false,
					},
					"width": map[string]interface{}{
						"type":    "number",
						"default": 900,
					},
					"height": map[string]interface{}{
						"type":    "number",
						"default": 500,
					},
				},
				"required": []string{"series"},
			},
		},
		s.handleLineChart,
	)

	// Tool: scatter_plot
	s.server.AddTool(
		&mcp.Tool{
			Name:        "scatter_plot",
			Description: "Generate a scatter plot for showing correlation between two variables",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title": map[string]interface{}{
						"type":        "string",
						"description": "Chart title",
					},
					"data": map[string]interface{}{
						"type":        "array",
						"description": "Array of {x, y, label?, size?} objects",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"x":     map[string]string{"type": "number"},
								"y":     map[string]string{"type": "number"},
								"label": map[string]string{"type": "string"},
								"size":  map[string]string{"type": "number"},
							},
							"required": []string{"x", "y"},
						},
					},
					"x_label": map[string]interface{}{
						"type":        "string",
						"description": "X-axis label",
					},
					"y_label": map[string]interface{}{
						"type":        "string",
						"description": "Y-axis label",
					},
					"width": map[string]interface{}{
						"type":    "number",
						"default": 800,
					},
					"height": map[string]interface{}{
						"type":    "number",
						"default": 600,
					},
				},
				"required": []string{"data"},
			},
		},
		s.handleScatterPlot,
	)

	// Tool: heatmap
	s.server.AddTool(
		&mcp.Tool{
			Name:        "heatmap",
			Description: "Generate a heatmap for matrix data showing intensity with colors",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title": map[string]interface{}{
						"type":        "string",
						"description": "Chart title",
					},
					"data": map[string]interface{}{
						"type":        "object",
						"description": "Matrix data with rows, columns, and values",
						"properties": map[string]interface{}{
							"rows":    map[string]interface{}{"type": "array", "items": map[string]string{"type": "string"}},
							"columns": map[string]interface{}{"type": "array", "items": map[string]string{"type": "string"}},
							"values": map[string]interface{}{
								"type": "array",
								"items": map[string]interface{}{
									"type":  "array",
									"items": map[string]string{"type": "number"},
								},
							},
						},
						"required": []string{"rows", "columns", "values"},
					},
					"color_map": map[string]interface{}{
						"type":        "string",
						"description": "Color map: 'viridis', 'plasma', 'red-green'",
						"default":     "viridis",
					},
					"show_value": map[string]interface{}{
						"type":        "boolean",
						"description": "Show numerical values in cells",
						"default":     true,
					},
					"width": map[string]interface{}{
						"type":    "number",
						"default": 800,
					},
					"height": map[string]interface{}{
						"type":    "number",
						"default": 600,
					},
				},
				"required": []string{"data"},
			},
		},
		s.handleHeatmap,
	)

	fmt.Println("Registered 5 chart generation tools")
}

// handleBarChart handles the bar_chart tool
func (s *Server) handleBarChart(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var config types.BarChartConfig
	if err := parseArguments(request.Params.Arguments, &config); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	// Set defaults
	if config.Width == 0 {
		config.Width = 800
	}
	if config.Height == 0 {
		config.Height = 400
	}
	if config.Orientation == "" {
		config.Orientation = "vertical"
	}

	// Generate chart
	svg, err := charts.CreateBarChart(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create bar chart: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("```svg\n%s\n```", svg),
			},
		},
	}, nil
}

// handlePieChart handles the pie_chart tool
func (s *Server) handlePieChart(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var config types.PieChartConfig
	if err := parseArguments(request.Params.Arguments, &config); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	// Set defaults
	if config.Width == 0 {
		config.Width = 600
	}
	if config.Height == 0 {
		config.Height = 600
	}

	svg, err := charts.CreatePieChart(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create pie chart: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("```svg\n%s\n```", svg),
			},
		},
	}, nil
}

// handleLineChart handles the line_chart tool
func (s *Server) handleLineChart(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var config types.LineChartConfig
	if err := parseArguments(request.Params.Arguments, &config); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	// Set defaults
	if config.Width == 0 {
		config.Width = 900
	}
	if config.Height == 0 {
		config.Height = 500
	}

	svg, err := charts.CreateLineChart(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create line chart: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("```svg\n%s\n```", svg),
			},
		},
	}, nil
}

// handleScatterPlot handles the scatter_plot tool
func (s *Server) handleScatterPlot(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var config types.ScatterPlotConfig
	if err := parseArguments(request.Params.Arguments, &config); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	// Set defaults
	if config.Width == 0 {
		config.Width = 800
	}
	if config.Height == 0 {
		config.Height = 600
	}

	svg, err := charts.CreateScatterPlot(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create scatter plot: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("```svg\n%s\n```", svg),
			},
		},
	}, nil
}

// handleHeatmap handles the heatmap tool
func (s *Server) handleHeatmap(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var config types.HeatmapConfig
	if err := parseArguments(request.Params.Arguments, &config); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	// Set defaults
	if config.Width == 0 {
		config.Width = 800
	}
	if config.Height == 0 {
		config.Height = 600
	}
	if config.ColorMap == "" {
		config.ColorMap = "viridis"
	}

	svg, err := charts.CreateHeatmap(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create heatmap: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("```svg\n%s\n```", svg),
			},
		},
	}, nil
}

// parseArguments helper to parse tool arguments from map[string]any
func parseArguments(args interface{}, target interface{}) error {
	// Convert to JSON and back to handle type conversions properly
	data, err := json.Marshal(args)
	if err != nil {
		return fmt.Errorf("failed to marshal arguments: %w", err)
	}

	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("failed to unmarshal arguments: %w", err)
	}

	return nil
}
