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

	// Tool: treemap
	s.server.AddTool(
		&mcp.Tool{
			Name:        "treemap",
			Description: "Generate a treemap visualization for hierarchical data",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title": map[string]interface{}{"type": "string", "description": "Chart title"},
					"data": map[string]interface{}{
						"type":        "object",
						"description": "Hierarchical tree data with name, value, and optional children",
						"properties": map[string]interface{}{
							"name":     map[string]string{"type": "string"},
							"value":    map[string]string{"type": "number"},
							"children": map[string]interface{}{"type": "array", "items": map[string]string{"type": "object"}},
						},
					},
					"show_labels": map[string]interface{}{"type": "boolean", "default": true},
					"width":       map[string]interface{}{"type": "number", "default": 800},
					"height":      map[string]interface{}{"type": "number", "default": 600},
				},
				"required": []string{"data"},
			},
		},
		s.handleTreemap,
	)

	// Tool: sunburst
	s.server.AddTool(
		&mcp.Tool{
			Name:        "sunburst",
			Description: "Generate a sunburst (radial partition) chart for hierarchical data",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title":       map[string]interface{}{"type": "string", "description": "Chart title"},
					"data":        map[string]interface{}{"type": "object", "description": "Hierarchical tree data"},
					"show_labels": map[string]interface{}{"type": "boolean", "default": true},
					"width":       map[string]interface{}{"type": "number", "default": 600},
					"height":      map[string]interface{}{"type": "number", "default": 600},
				},
				"required": []string{"data"},
			},
		},
		s.handleSunburst,
	)

	// Tool: circle_packing
	s.server.AddTool(
		&mcp.Tool{
			Name:        "circle_packing",
			Description: "Generate a circle packing visualization for hierarchical data",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title":       map[string]interface{}{"type": "string", "description": "Chart title"},
					"data":        map[string]interface{}{"type": "object", "description": "Hierarchical tree data"},
					"show_labels": map[string]interface{}{"type": "boolean", "default": true},
					"width":       map[string]interface{}{"type": "number", "default": 600},
					"height":      map[string]interface{}{"type": "number", "default": 600},
				},
				"required": []string{"data"},
			},
		},
		s.handleCirclePacking,
	)

	// Tool: icicle
	s.server.AddTool(
		&mcp.Tool{
			Name:        "icicle",
			Description: "Generate an icicle partition chart for hierarchical data",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title":       map[string]interface{}{"type": "string", "description": "Chart title"},
					"data":        map[string]interface{}{"type": "object", "description": "Hierarchical tree data"},
					"orientation": map[string]interface{}{"type": "string", "description": "vertical or horizontal", "default": "vertical"},
					"show_labels": map[string]interface{}{"type": "boolean", "default": true},
					"width":       map[string]interface{}{"type": "number", "default": 800},
					"height":      map[string]interface{}{"type": "number", "default": 600},
				},
				"required": []string{"data"},
			},
		},
		s.handleIcicle,
	)

	// Tool: boxplot
	s.server.AddTool(
		&mcp.Tool{
			Name:        "boxplot",
			Description: "Generate a box plot for showing statistical distribution",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title": map[string]interface{}{"type": "string", "description": "Chart title"},
					"data": map[string]interface{}{
						"type":        "array",
						"description": "Array of {label, values[]} objects",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"label":  map[string]string{"type": "string"},
								"values": map[string]interface{}{"type": "array", "items": map[string]string{"type": "number"}},
							},
						},
					},
					"show_outliers": map[string]interface{}{"type": "boolean", "default": true},
					"show_mean":     map[string]interface{}{"type": "boolean", "default": false},
					"width":         map[string]interface{}{"type": "number", "default": 800},
					"height":        map[string]interface{}{"type": "number", "default": 600},
				},
				"required": []string{"data"},
			},
		},
		s.handleBoxplot,
	)

	// Tool: violin
	s.server.AddTool(
		&mcp.Tool{
			Name:        "violin",
			Description: "Generate a violin plot with kernel density estimation",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title": map[string]interface{}{"type": "string", "description": "Chart title"},
					"data": map[string]interface{}{
						"type":        "array",
						"description": "Array of {label, values[]} objects",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"label":  map[string]string{"type": "string"},
								"values": map[string]interface{}{"type": "array", "items": map[string]string{"type": "number"}},
							},
						},
					},
					"show_box":    map[string]interface{}{"type": "boolean", "default": true},
					"show_median": map[string]interface{}{"type": "boolean", "default": true},
					"width":       map[string]interface{}{"type": "number", "default": 800},
					"height":      map[string]interface{}{"type": "number", "default": 600},
				},
				"required": []string{"data"},
			},
		},
		s.handleViolin,
	)

	// Tool: histogram
	s.server.AddTool(
		&mcp.Tool{
			Name:        "histogram",
			Description: "Generate a histogram with automatic binning",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title":  map[string]interface{}{"type": "string", "description": "Chart title"},
					"values": map[string]interface{}{"type": "array", "items": map[string]string{"type": "number"}, "description": "Array of numerical values"},
					"bins":   map[string]interface{}{"type": "number", "description": "Number of bins", "default": 20},
					"width":  map[string]interface{}{"type": "number", "default": 800},
					"height": map[string]interface{}{"type": "number", "default": 600},
				},
				"required": []string{"values"},
			},
		},
		s.handleHistogram,
	)

	// Tool: ridgeline
	s.server.AddTool(
		&mcp.Tool{
			Name:        "ridgeline",
			Description: "Generate a ridgeline (joy) plot for comparing distributions",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title": map[string]interface{}{"type": "string", "description": "Chart title"},
					"data": map[string]interface{}{
						"type":        "array",
						"description": "Array of {label, values[]} objects",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"label":  map[string]string{"type": "string"},
								"values": map[string]interface{}{"type": "array", "items": map[string]string{"type": "number"}},
							},
						},
					},
					"overlap":     map[string]interface{}{"type": "number", "description": "Overlap amount (0-1)", "default": 0.5},
					"show_labels": map[string]interface{}{"type": "boolean", "default": true},
					"width":       map[string]interface{}{"type": "number", "default": 800},
					"height":      map[string]interface{}{"type": "number", "default": 600},
				},
				"required": []string{"data"},
			},
		},
		s.handleRidgeline,
	)

	// Tool: candlestick
	s.server.AddTool(
		&mcp.Tool{
			Name:        "candlestick",
			Description: "Generate a candlestick chart for financial OHLC data",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title": map[string]interface{}{"type": "string", "description": "Chart title"},
					"data": map[string]interface{}{
						"type":        "array",
						"description": "Array of {date, open, high, low, close, volume?} objects",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"date":   map[string]string{"type": "string"},
								"open":   map[string]string{"type": "number"},
								"high":   map[string]string{"type": "number"},
								"low":    map[string]string{"type": "number"},
								"close":  map[string]string{"type": "number"},
								"volume": map[string]string{"type": "number"},
							},
						},
					},
					"show_volume": map[string]interface{}{"type": "boolean", "default": true},
					"width":       map[string]interface{}{"type": "number", "default": 1000},
					"height":      map[string]interface{}{"type": "number", "default": 600},
				},
				"required": []string{"data"},
			},
		},
		s.handleCandlestick,
	)

	// Tool: ohlc
	s.server.AddTool(
		&mcp.Tool{
			Name:        "ohlc",
			Description: "Generate an OHLC bar chart for financial data",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title": map[string]interface{}{"type": "string", "description": "Chart title"},
					"data": map[string]interface{}{
						"type":        "array",
						"description": "Array of {date, open, high, low, close} objects",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"date":  map[string]string{"type": "string"},
								"open":  map[string]string{"type": "number"},
								"high":  map[string]string{"type": "number"},
								"low":   map[string]string{"type": "number"},
								"close": map[string]string{"type": "number"},
							},
						},
					},
					"width":  map[string]interface{}{"type": "number", "default": 1000},
					"height": map[string]interface{}{"type": "number", "default": 600},
				},
				"required": []string{"data"},
			},
		},
		s.handleOHLC,
	)

	fmt.Println("Registered 15 chart generation tools")
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

// handleTreemap handles the treemap tool
func (s *Server) handleTreemap(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var config types.TreemapConfig
	if err := parseArguments(request.Params.Arguments, &config); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	if config.Width == 0 {
		config.Width = 800
	}
	if config.Height == 0 {
		config.Height = 600
	}

	svg, err := charts.CreateTreemap(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create treemap: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("```svg\n%s\n```", svg),
			},
		},
	}, nil
}

// handleSunburst handles the sunburst tool
func (s *Server) handleSunburst(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var config types.SunburstConfig
	if err := parseArguments(request.Params.Arguments, &config); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	if config.Width == 0 {
		config.Width = 600
	}
	if config.Height == 0 {
		config.Height = 600
	}

	svg, err := charts.CreateSunburst(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create sunburst: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("```svg\n%s\n```", svg),
			},
		},
	}, nil
}

// handleCirclePacking handles the circle_packing tool
func (s *Server) handleCirclePacking(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var config types.CirclePackingConfig
	if err := parseArguments(request.Params.Arguments, &config); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	if config.Width == 0 {
		config.Width = 600
	}
	if config.Height == 0 {
		config.Height = 600
	}

	svg, err := charts.CreateCirclePacking(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create circle packing: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("```svg\n%s\n```", svg),
			},
		},
	}, nil
}

// handleIcicle handles the icicle tool
func (s *Server) handleIcicle(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var config types.IcicleConfig
	if err := parseArguments(request.Params.Arguments, &config); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	if config.Width == 0 {
		config.Width = 800
	}
	if config.Height == 0 {
		config.Height = 600
	}

	svg, err := charts.CreateIcicle(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create icicle: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("```svg\n%s\n```", svg),
			},
		},
	}, nil
}

// handleBoxplot handles the boxplot tool
func (s *Server) handleBoxplot(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var config types.BoxPlotConfig
	if err := parseArguments(request.Params.Arguments, &config); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	if config.Width == 0 {
		config.Width = 800
	}
	if config.Height == 0 {
		config.Height = 600
	}

	svg, err := charts.CreateBoxPlot(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create boxplot: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("```svg\n%s\n```", svg),
			},
		},
	}, nil
}

// handleViolin handles the violin tool
func (s *Server) handleViolin(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var config types.ViolinPlotConfig
	if err := parseArguments(request.Params.Arguments, &config); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	if config.Width == 0 {
		config.Width = 800
	}
	if config.Height == 0 {
		config.Height = 600
	}

	svg, err := charts.CreateViolinPlot(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create violin plot: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("```svg\n%s\n```", svg),
			},
		},
	}, nil
}

// handleHistogram handles the histogram tool
func (s *Server) handleHistogram(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var config types.HistogramConfig
	if err := parseArguments(request.Params.Arguments, &config); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	if config.Width == 0 {
		config.Width = 800
	}
	if config.Height == 0 {
		config.Height = 600
	}

	svg, err := charts.CreateHistogram(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create histogram: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("```svg\n%s\n```", svg),
			},
		},
	}, nil
}

// handleRidgeline handles the ridgeline tool
func (s *Server) handleRidgeline(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var config types.RidgelineConfig
	if err := parseArguments(request.Params.Arguments, &config); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	if config.Width == 0 {
		config.Width = 800
	}
	if config.Height == 0 {
		config.Height = 600
	}

	svg, err := charts.CreateRidgeline(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create ridgeline plot: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("```svg\n%s\n```", svg),
			},
		},
	}, nil
}

// handleCandlestick handles the candlestick tool
func (s *Server) handleCandlestick(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var config types.CandlestickConfig
	if err := parseArguments(request.Params.Arguments, &config); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	if config.Width == 0 {
		config.Width = 1000
	}
	if config.Height == 0 {
		config.Height = 600
	}

	svg, err := charts.CreateCandlestick(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create candlestick chart: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("```svg\n%s\n```", svg),
			},
		},
	}, nil
}

// handleOHLC handles the ohlc tool
func (s *Server) handleOHLC(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var config types.OHLCConfig
	if err := parseArguments(request.Params.Arguments, &config); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	if config.Width == 0 {
		config.Width = 1000
	}
	if config.Height == 0 {
		config.Height = 600
	}

	svg, err := charts.CreateOHLC(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create OHLC chart: %w", err)
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
