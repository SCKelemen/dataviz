// Package mcp implements a Model Context Protocol (MCP) server for chart generation.
//
// This server exposes chart generation tools to MCP clients like Claude Code,
// enabling AI agents to create visualizations from data.
//
// Design Philosophy - Discrete & Orthogonal:
//   - Data-source agnostic: Accepts generic data from any source
//   - Single responsibility: Pure function: Data → SVG
//   - Composable: Works with other MCPs (Omnitron, file systems, APIs)
//   - Agent-centric: Optimized for AI workflow orchestration
//
// Example Workflow:
//
//	User: "Show me a bar chart of services by owner from Omnitron"
//
//	Claude Agent:
//	1. Query Omnitron MCP → Get services data
//	2. Aggregate by owner → Calculate counts
//	3. Call DataViz MCP → Generate bar chart
//	4. Return SVG to user
//
// MCP Tools Provided:
//   - create_line_chart: Multi-series line charts
//   - create_bar_chart: Vertical/horizontal bar charts
//   - create_scatter_plot: Scatter plots with multiple series
//   - create_heatmap: Matrix-based heatmaps
//   - create_pie_chart: Pie/donut charts
//   - export_to_image: Convert SVG to PNG or JPEG
//
// Configuration:
//
//	# claude_desktop_config.json
//	{
//	  "mcpServers": {
//	    "dataviz": {
//	      "command": "dataviz-mcp"
//	    }
//	  }
//	}
//
// The server accepts generic data types:
//   - X values: interface{} (can be time.Time, int, float64, string)
//   - Y values: float64
//
// This allows agents to visualize data from any source without
// type conversions.
//
// This package uses:
//   - charts/ - For chart generation
//   - export/ - For PNG/JPEG conversion
//   - MCP SDK - For protocol implementation
package mcp
