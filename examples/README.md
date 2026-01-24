# Chart Examples

This directory contains example input data and generated SVG outputs for all supported chart types.

Each subdirectory contains:
- `input.json` - Example input data in JSON format
- `output.svg` - Generated SVG visualization

## Fully Supported Charts

These charts work with the CLI tool (`viz-cli <type> < input.json > output.svg`):

### Statistical Charts
- `boxplot` - Box and whisker plot showing distribution statistics
- `violin` - Violin plot with kernel density estimation
- `histogram` - Histogram with automatic binning
- `ridgeline` - Ridgeline (joy) plot for comparing distributions
- `density` - Kernel density estimation plot
- `correlogram` - Correlation matrix visualization

### Line/Area Charts
- `connected-scatter` - Scatter plot with lines connecting points
- `stacked-area` - Stacked area chart
- `streamchart` - Streamchart with flowing areas (centered)

### Specialized Charts
- `lollipop` - Lollipop chart with stems and circles
- `radar` - Radar (spider) chart for multivariate data
- `parallel` - Parallel coordinates for multidimensional data
- `wordcloud` - Word cloud visualization

### Hierarchical Charts
- `treemap` - Squarified treemap for hierarchical data
- `sunburst` - Radial partition chart
- `circle-packing` - Hierarchical circle packing
- `icicle` - Icicle partition chart
- `dendrogram` - Hierarchical clustering tree

### Network/Flow Charts
- `sankey` - Sankey diagram for flow visualization
- `chord` - Chord diagram for relationships

### Circular Charts
- `circular-bar` - Circular bar plot

### Financial Charts
- `candlestick` - OHLC candlestick chart
- `ohlc` - OHLC bar chart

### Legacy Charts
- `bar-chart` - Simple bar chart (legacy renderer)
- `line-graph` - Time series line graph (legacy renderer)
- `heatmap` - Contribution calendar heatmap (legacy renderer)

### Placeholder Charts (MCP Only)
These charts display a placeholder message in CLI but work fully in MCP server:
- `pie` - Pie/donut chart (use MCP server)
- `scatter` - Basic scatter plot (use MCP server)

## Usage

Generate an SVG from example data:
```bash
viz-cli <chart-type> < examples/<chart-type>/input.json > output.svg
```

Examples:
```bash
# Generate a lollipop chart
viz-cli lollipop < examples/lollipop/input.json > lollipop.svg

# Generate a radar chart with custom dimensions
viz-cli radar -width 1000 -height 1000 < examples/radar/input.json > radar.svg

# Generate a sankey diagram
viz-cli sankey < examples/sankey/input.json > sankey.svg
```

## MCP Server Support

All charts are also available through the MCP server with 28 total tools. The MCP server provides additional configuration options and type-safe schemas.

## Chart Data Formats

See each `input.json` file for the expected data format for that chart type. Common patterns:

- **Simple values**: `{"values": [{"label": "A", "value": 23}, ...]}`
- **Time series**: `{"series": [{"name": "Series A", "data": [{"x": 0, "y": 10}, ...]}]}`
- **Hierarchical**: `{"data": {"name": "Root", "children": [...]}}`
- **Matrix**: `{"rows": [...], "columns": [...], "values": [[...]]}`
- **Flow**: `{"nodes": [...], "links": [{"source": "A", "target": "B", "value": 50}]}`
