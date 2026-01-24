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

	// Tool: lollipop
	s.server.AddTool(
		&mcp.Tool{
			Name:        "lollipop",
			Description: "Generate a lollipop chart with stems and circles",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title": map[string]interface{}{"type": "string", "description": "Chart title"},
					"values": map[string]interface{}{
						"type":        "array",
						"description": "Array of {label, value, color?} objects",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"label": map[string]string{"type": "string"},
								"value": map[string]string{"type": "number"},
								"color": map[string]string{"type": "string"},
							},
							"required": []string{"label", "value"},
						},
					},
					"color":      map[string]interface{}{"type": "string", "description": "Default color for all lollipops"},
					"horizontal": map[string]interface{}{"type": "boolean", "default": false, "description": "Horizontal orientation"},
					"width":      map[string]interface{}{"type": "number", "default": 800},
					"height":     map[string]interface{}{"type": "number", "default": 600},
				},
				"required": []string{"values"},
			},
		},
		s.handleLollipop,
	)

	// Tool: density
	s.server.AddTool(
		&mcp.Tool{
			Name:        "density",
			Description: "Generate a kernel density estimation plot",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title": map[string]interface{}{"type": "string", "description": "Chart title"},
					"data": map[string]interface{}{
						"type":        "array",
						"description": "Array of {values[], label?, color?} datasets",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"values": map[string]interface{}{"type": "array", "items": map[string]string{"type": "number"}},
								"label":  map[string]string{"type": "string"},
								"color":  map[string]string{"type": "string"},
							},
							"required": []string{"values"},
						},
					},
					"show_fill": map[string]interface{}{"type": "boolean", "default": true},
					"show_rug":  map[string]interface{}{"type": "boolean", "default": false},
					"width":     map[string]interface{}{"type": "number", "default": 800},
					"height":    map[string]interface{}{"type": "number", "default": 600},
				},
				"required": []string{"data"},
			},
		},
		s.handleDensity,
	)

	// Tool: connected_scatter
	s.server.AddTool(
		&mcp.Tool{
			Name:        "connected_scatter",
			Description: "Generate a connected scatter plot with lines between points",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title": map[string]interface{}{"type": "string", "description": "Chart title"},
					"series": map[string]interface{}{
						"type":        "array",
						"description": "Array of series with points",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"points": map[string]interface{}{
									"type": "array",
									"items": map[string]interface{}{
										"type": "object",
										"properties": map[string]interface{}{
											"x":     map[string]string{"type": "number"},
											"y":     map[string]string{"type": "number"},
											"label": map[string]string{"type": "string"},
										},
										"required": []string{"x", "y"},
									},
								},
								"label":       map[string]string{"type": "string"},
								"color":       map[string]string{"type": "string"},
								"marker_type": map[string]string{"type": "string"},
							},
							"required": []string{"points"},
						},
					},
					"width":  map[string]interface{}{"type": "number", "default": 800},
					"height": map[string]interface{}{"type": "number", "default": 600},
				},
				"required": []string{"series"},
			},
		},
		s.handleConnectedScatter,
	)

	// Tool: stacked_area
	s.server.AddTool(
		&mcp.Tool{
			Name:        "stacked_area",
			Description: "Generate a stacked area chart",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title": map[string]interface{}{"type": "string", "description": "Chart title"},
					"points": map[string]interface{}{
						"type":        "array",
						"description": "Array of {x, values[]} points",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"x":      map[string]string{"type": "number"},
								"values": map[string]interface{}{"type": "array", "items": map[string]string{"type": "number"}},
							},
							"required": []string{"x", "values"},
						},
					},
					"series": map[string]interface{}{
						"type":        "array",
						"description": "Array of {label, color?} series metadata",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"label": map[string]string{"type": "string"},
								"color": map[string]string{"type": "string"},
							},
							"required": []string{"label"},
						},
					},
					"width":  map[string]interface{}{"type": "number", "default": 800},
					"height": map[string]interface{}{"type": "number", "default": 600},
				},
				"required": []string{"points", "series"},
			},
		},
		s.handleStackedArea,
	)

	// Tool: streamchart
	s.server.AddTool(
		&mcp.Tool{
			Name:        "streamchart",
			Description: "Generate a streamchart (flowing stacked area)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title": map[string]interface{}{"type": "string", "description": "Chart title"},
					"points": map[string]interface{}{
						"type":        "array",
						"description": "Array of {x, values[]} points",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"x":      map[string]string{"type": "number"},
								"values": map[string]interface{}{"type": "array", "items": map[string]string{"type": "number"}},
							},
							"required": []string{"x", "values"},
						},
					},
					"series": map[string]interface{}{
						"type":        "array",
						"description": "Array of {label, color?} series metadata",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"label": map[string]string{"type": "string"},
								"color": map[string]string{"type": "string"},
							},
							"required": []string{"label"},
						},
					},
					"layout": map[string]interface{}{"type": "string", "description": "Layout type", "default": "wiggle"},
					"width":  map[string]interface{}{"type": "number", "default": 800},
					"height": map[string]interface{}{"type": "number", "default": 600},
				},
				"required": []string{"points", "series"},
			},
		},
		s.handleStreamChart,
	)

	// Tool: correlogram
	s.server.AddTool(
		&mcp.Tool{
			Name:        "correlogram",
			Description: "Generate a correlogram (correlation matrix visualization)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title":     map[string]interface{}{"type": "string", "description": "Chart title"},
					"variables": map[string]interface{}{"type": "array", "items": map[string]string{"type": "string"}, "description": "Variable names"},
					"matrix": map[string]interface{}{
						"type":        "array",
						"description": "Correlation matrix (2D array of values)",
						"items": map[string]interface{}{
							"type":  "array",
							"items": map[string]string{"type": "number"},
						},
					},
					"width":  map[string]interface{}{"type": "number", "default": 800},
					"height": map[string]interface{}{"type": "number", "default": 800},
				},
				"required": []string{"variables", "matrix"},
			},
		},
		s.handleCorrelogram,
	)

	// Tool: radar
	s.server.AddTool(
		&mcp.Tool{
			Name:        "radar",
			Description: "Generate a radar (spider) chart",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title": map[string]interface{}{"type": "string", "description": "Chart title"},
					"axes": map[string]interface{}{
						"type":        "array",
						"description": "Array of {label, min, max} axes",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"label": map[string]string{"type": "string"},
								"min":   map[string]string{"type": "number"},
								"max":   map[string]string{"type": "number"},
							},
							"required": []string{"label", "min", "max"},
						},
					},
					"series": map[string]interface{}{
						"type":        "array",
						"description": "Array of {label, values[], color?} series",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"label":  map[string]string{"type": "string"},
								"values": map[string]interface{}{"type": "array", "items": map[string]string{"type": "number"}},
								"color":  map[string]string{"type": "string"},
							},
							"required": []string{"label", "values"},
						},
					},
					"width":  map[string]interface{}{"type": "number", "default": 600},
					"height": map[string]interface{}{"type": "number", "default": 600},
				},
				"required": []string{"axes", "series"},
			},
		},
		s.handleRadar,
	)

	// Tool: parallel
	s.server.AddTool(
		&mcp.Tool{
			Name:        "parallel",
			Description: "Generate a parallel coordinates plot",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title": map[string]interface{}{"type": "string", "description": "Chart title"},
					"axes": map[string]interface{}{
						"type":        "array",
						"description": "Array of {label, min, max} axes",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"label": map[string]string{"type": "string"},
								"min":   map[string]string{"type": "number"},
								"max":   map[string]string{"type": "number"},
							},
							"required": []string{"label", "min", "max"},
						},
					},
					"data": map[string]interface{}{
						"type":        "array",
						"description": "Array of {values[], color?} data points",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"values": map[string]interface{}{"type": "array", "items": map[string]string{"type": "number"}},
								"color":  map[string]string{"type": "string"},
							},
							"required": []string{"values"},
						},
					},
					"width":  map[string]interface{}{"type": "number", "default": 800},
					"height": map[string]interface{}{"type": "number", "default": 600},
				},
				"required": []string{"axes", "data"},
			},
		},
		s.handleParallel,
	)

	// Tool: wordcloud
	s.server.AddTool(
		&mcp.Tool{
			Name:        "wordcloud",
			Description: "Generate a word cloud visualization",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title": map[string]interface{}{"type": "string", "description": "Chart title"},
					"words": map[string]interface{}{
						"type":        "array",
						"description": "Array of {text, frequency, color?} words",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"text":      map[string]string{"type": "string"},
								"frequency": map[string]string{"type": "number"},
								"color":     map[string]string{"type": "string"},
							},
							"required": []string{"text", "frequency"},
						},
					},
					"layout": map[string]interface{}{"type": "string", "description": "Layout algorithm", "default": "spiral"},
					"width":  map[string]interface{}{"type": "number", "default": 800},
					"height": map[string]interface{}{"type": "number", "default": 600},
				},
				"required": []string{"words"},
			},
		},
		s.handleWordCloud,
	)

	// Tool: sankey
	s.server.AddTool(
		&mcp.Tool{
			Name:        "sankey",
			Description: "Generate a Sankey diagram for flow visualization",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title": map[string]interface{}{"type": "string", "description": "Chart title"},
					"nodes": map[string]interface{}{
						"type":        "array",
						"description": "Array of {id, label, color?} nodes",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"id":    map[string]string{"type": "string"},
								"label": map[string]string{"type": "string"},
								"color": map[string]string{"type": "string"},
							},
							"required": []string{"id", "label"},
						},
					},
					"links": map[string]interface{}{
						"type":        "array",
						"description": "Array of {source, target, value, color?} links",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"source": map[string]string{"type": "string"},
								"target": map[string]string{"type": "string"},
								"value":  map[string]string{"type": "number"},
								"color":  map[string]string{"type": "string"},
							},
							"required": []string{"source", "target", "value"},
						},
					},
					"width":  map[string]interface{}{"type": "number", "default": 1000},
					"height": map[string]interface{}{"type": "number", "default": 600},
				},
				"required": []string{"nodes", "links"},
			},
		},
		s.handleSankey,
	)

	// Tool: chord
	s.server.AddTool(
		&mcp.Tool{
			Name:        "chord",
			Description: "Generate a chord diagram for relationships",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title": map[string]interface{}{"type": "string", "description": "Chart title"},
					"entities": map[string]interface{}{
						"type":        "array",
						"description": "Array of {id, label, color?} entities",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"id":    map[string]string{"type": "string"},
								"label": map[string]string{"type": "string"},
								"color": map[string]string{"type": "string"},
							},
							"required": []string{"id", "label"},
						},
					},
					"relations": map[string]interface{}{
						"type":        "array",
						"description": "Array of {source, target, value} relations",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"source": map[string]string{"type": "string"},
								"target": map[string]string{"type": "string"},
								"value":  map[string]string{"type": "number"},
							},
							"required": []string{"source", "target", "value"},
						},
					},
					"width":  map[string]interface{}{"type": "number", "default": 800},
					"height": map[string]interface{}{"type": "number", "default": 800},
				},
				"required": []string{"entities", "relations"},
			},
		},
		s.handleChord,
	)

	// Tool: circular_bar
	s.server.AddTool(
		&mcp.Tool{
			Name:        "circular_bar",
			Description: "Generate a circular bar plot",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title": map[string]interface{}{"type": "string", "description": "Chart title"},
					"data": map[string]interface{}{
						"type":        "array",
						"description": "Array of {label, value, color?} data points",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"label": map[string]string{"type": "string"},
								"value": map[string]string{"type": "number"},
								"color": map[string]string{"type": "string"},
							},
							"required": []string{"label", "value"},
						},
					},
					"inner_radius": map[string]interface{}{"type": "number", "description": "Inner radius ratio", "default": 0.3},
					"width":        map[string]interface{}{"type": "number", "default": 800},
					"height":       map[string]interface{}{"type": "number", "default": 800},
				},
				"required": []string{"data"},
			},
		},
		s.handleCircularBar,
	)

	// Tool: dendrogram
	s.server.AddTool(
		&mcp.Tool{
			Name:        "dendrogram",
			Description: "Generate a dendrogram (hierarchical clustering tree)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title": map[string]interface{}{"type": "string", "description": "Chart title"},
					"root": map[string]interface{}{
						"type":        "object",
						"description": "Root node with recursive children structure",
						"properties": map[string]interface{}{
							"label":    map[string]string{"type": "string"},
							"height":   map[string]string{"type": "number"},
							"children": map[string]interface{}{"type": "array", "items": map[string]string{"type": "object"}},
						},
						"required": []string{"height"},
					},
					"orientation": map[string]interface{}{"type": "string", "description": "vertical or horizontal", "default": "horizontal"},
					"width":       map[string]interface{}{"type": "number", "default": 800},
					"height":      map[string]interface{}{"type": "number", "default": 600},
				},
				"required": []string{"root"},
			},
		},
		s.handleDendrogram,
	)

	fmt.Println("Registered 28 chart generation tools")
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

// handleLollipop handles the lollipop tool
func (s *Server) handleLollipop(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var config types.LollipopConfig
	if err := parseArguments(request.Params.Arguments, &config); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	if config.Width == 0 {
		config.Width = 800
	}
	if config.Height == 0 {
		config.Height = 600
	}

	svg, err := charts.CreateLollipop(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create lollipop chart: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("```svg\n%s\n```", svg),
			},
		},
	}, nil
}

// handleDensity handles the density tool
func (s *Server) handleDensity(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var config types.DensityConfig
	if err := parseArguments(request.Params.Arguments, &config); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	if config.Width == 0 {
		config.Width = 800
	}
	if config.Height == 0 {
		config.Height = 600
	}

	svg, err := charts.CreateDensity(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create density plot: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("```svg\n%s\n```", svg),
			},
		},
	}, nil
}

// handleConnectedScatter handles the connected_scatter tool
func (s *Server) handleConnectedScatter(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var config types.ConnectedScatterConfig
	if err := parseArguments(request.Params.Arguments, &config); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	if config.Width == 0 {
		config.Width = 800
	}
	if config.Height == 0 {
		config.Height = 600
	}

	svg, err := charts.CreateConnectedScatter(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connected scatter plot: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("```svg\n%s\n```", svg),
			},
		},
	}, nil
}

// handleStackedArea handles the stacked_area tool
func (s *Server) handleStackedArea(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var config types.StackedAreaConfig
	if err := parseArguments(request.Params.Arguments, &config); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	if config.Width == 0 {
		config.Width = 800
	}
	if config.Height == 0 {
		config.Height = 600
	}

	svg, err := charts.CreateStackedArea(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create stacked area chart: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("```svg\n%s\n```", svg),
			},
		},
	}, nil
}

// handleStreamChart handles the streamchart tool
func (s *Server) handleStreamChart(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var config types.StreamChartConfig
	if err := parseArguments(request.Params.Arguments, &config); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	if config.Width == 0 {
		config.Width = 800
	}
	if config.Height == 0 {
		config.Height = 600
	}

	svg, err := charts.CreateStreamChart(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create streamchart: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("```svg\n%s\n```", svg),
			},
		},
	}, nil
}

// handleCorrelogram handles the correlogram tool
func (s *Server) handleCorrelogram(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var config types.CorrelogramConfig
	if err := parseArguments(request.Params.Arguments, &config); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	if config.Width == 0 {
		config.Width = 800
	}
	if config.Height == 0 {
		config.Height = 800
	}

	svg, err := charts.CreateCorrelogram(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create correlogram: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("```svg\n%s\n```", svg),
			},
		},
	}, nil
}

// handleRadar handles the radar tool
func (s *Server) handleRadar(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var config types.RadarConfig
	if err := parseArguments(request.Params.Arguments, &config); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	if config.Width == 0 {
		config.Width = 600
	}
	if config.Height == 0 {
		config.Height = 600
	}

	svg, err := charts.CreateRadar(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create radar chart: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("```svg\n%s\n```", svg),
			},
		},
	}, nil
}

// handleParallel handles the parallel tool
func (s *Server) handleParallel(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var config types.ParallelConfig
	if err := parseArguments(request.Params.Arguments, &config); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	if config.Width == 0 {
		config.Width = 800
	}
	if config.Height == 0 {
		config.Height = 600
	}

	svg, err := charts.CreateParallel(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create parallel coordinates plot: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("```svg\n%s\n```", svg),
			},
		},
	}, nil
}

// handleWordCloud handles the wordcloud tool
func (s *Server) handleWordCloud(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var config types.WordCloudConfig
	if err := parseArguments(request.Params.Arguments, &config); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	if config.Width == 0 {
		config.Width = 800
	}
	if config.Height == 0 {
		config.Height = 600
	}

	svg, err := charts.CreateWordCloud(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create word cloud: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("```svg\n%s\n```", svg),
			},
		},
	}, nil
}

// handleSankey handles the sankey tool
func (s *Server) handleSankey(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var config types.SankeyConfig
	if err := parseArguments(request.Params.Arguments, &config); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	if config.Width == 0 {
		config.Width = 1000
	}
	if config.Height == 0 {
		config.Height = 600
	}

	svg, err := charts.CreateSankey(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Sankey diagram: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("```svg\n%s\n```", svg),
			},
		},
	}, nil
}

// handleChord handles the chord tool
func (s *Server) handleChord(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var config types.ChordConfig
	if err := parseArguments(request.Params.Arguments, &config); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	if config.Width == 0 {
		config.Width = 800
	}
	if config.Height == 0 {
		config.Height = 800
	}

	svg, err := charts.CreateChord(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create chord diagram: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("```svg\n%s\n```", svg),
			},
		},
	}, nil
}

// handleCircularBar handles the circular_bar tool
func (s *Server) handleCircularBar(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var config types.CircularBarConfig
	if err := parseArguments(request.Params.Arguments, &config); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	if config.Width == 0 {
		config.Width = 800
	}
	if config.Height == 0 {
		config.Height = 800
	}

	svg, err := charts.CreateCircularBar(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create circular bar plot: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("```svg\n%s\n```", svg),
			},
		},
	}, nil
}

// handleDendrogram handles the dendrogram tool
func (s *Server) handleDendrogram(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var config types.DendrogramConfig
	if err := parseArguments(request.Params.Arguments, &config); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	if config.Width == 0 {
		config.Width = 800
	}
	if config.Height == 0 {
		config.Height = 600
	}

	svg, err := charts.CreateDendrogram(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create dendrogram: %w", err)
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
