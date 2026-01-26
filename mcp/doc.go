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
// MCP Tools Provided (29 total):
//
// Chart Generation:
//   - bar_chart: Vertical/horizontal bar charts
//   - pie_chart: Pie/donut charts
//   - line_chart: Multi-series line charts
//   - scatter_plot: Scatter plots with correlation
//   - heatmap: Matrix-based heatmaps
//   - treemap: Hierarchical treemap visualizations
//   - sunburst: Radial partition charts
//   - circle_packing: Circle packing for hierarchies
//   - icicle: Icicle partition charts
//   - boxplot: Box plots for distributions
//   - violin: Violin plots with KDE
//   - histogram: Histograms with auto-binning
//   - ridgeline: Ridgeline (joy) plots
//   - candlestick: Financial candlestick charts
//   - ohlc: OHLC bar charts
//   - lollipop: Lollipop charts
//   - density: Kernel density estimation
//   - connected_scatter: Connected scatter plots
//   - stacked_area: Stacked area charts
//   - streamchart: Streamcharts (flowing stacked areas)
//   - correlogram: Correlation matrix visualizations
//   - radar: Radar (spider) charts
//   - parallel: Parallel coordinates plots
//   - wordcloud: Word cloud visualizations
//   - sankey: Sankey flow diagrams
//   - chord: Chord diagrams for relationships
//   - circular_bar: Circular bar plots
//   - dendrogram: Hierarchical clustering trees
//   - generate_gallery: Generate comparison galleries of chart variants
//
// Gallery Tool:
//
// The generate_gallery tool creates comparison galleries showing multiple
// variations of a chart type side-by-side. This is useful for:
//   - Demonstrating available configuration options
//   - Comparing different styles or data representations
//   - Educational purposes (learning chart capabilities)
//   - Visual documentation
//
// Supported gallery types:
//   bar, area, stacked-area, lollipop, histogram, pie, boxplot, violin,
//   treemap, icicle, ridgeline, line, scatter, connected-scatter, statcard,
//   radar, streamchart, candlestick, sunburst, circle-packing, heatmap
//
// Example:
//	gallery_type: "bar"  → Generates gallery with simple and stacked bar variants
//	gallery_type: "line" → Generates gallery with different line styles and options
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
