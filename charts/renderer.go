package charts

// Output represents rendered visualization output
type Output interface {
	String() string
}

// SVGOutput wraps SVG string output
type SVGOutput string

func (s SVGOutput) String() string {
	return string(s)
}

// Renderer defines the interface for visualization rendering
type Renderer interface {
	RenderHeatmap(data HeatmapData, bounds Bounds, config RenderConfig) Output
	RenderLineGraph(data LineGraphData, bounds Bounds, config RenderConfig) Output
	RenderBarChart(data BarChartData, bounds Bounds, config RenderConfig) Output
	RenderStatCard(data StatCardData, bounds Bounds, config RenderConfig) Output
	RenderAreaChart(data AreaChartData, bounds Bounds, config RenderConfig) Output
	RenderScatterPlot(data ScatterPlotData, bounds Bounds, config RenderConfig) Output
}

// SVGRenderer implements SVG rendering
type SVGRenderer struct{}

// NewSVGRenderer creates a new SVG renderer
func NewSVGRenderer() *SVGRenderer {
	return &SVGRenderer{}
}

// RenderHeatmap renders a heatmap as SVG
func (r *SVGRenderer) RenderHeatmap(data HeatmapData, bounds Bounds, config RenderConfig) Output {
	var svg string

	// Determine color from config
	color := config.Color
	if color == "" && config.DesignTokens != nil {
		color = config.DesignTokens.Color
	}

	// Call appropriate heatmap renderer based on type
	if data.Type == "weeks" {
		svg = RenderWeeksHeatmap(data, bounds.X, bounds.Y, bounds.Width, bounds.Height, color, config.DesignTokens)
	} else {
		// Default to linear heatmap
		svg = RenderLinearHeatmap(data, bounds.X, bounds.Y, bounds.Width, bounds.Height, color, config.DesignTokens)
	}

	return SVGOutput(svg)
}

// RenderLineGraph renders a line graph as SVG
func (r *SVGRenderer) RenderLineGraph(data LineGraphData, bounds Bounds, config RenderConfig) Output {
	svg := RenderLineGraph(data, bounds.X, bounds.Y, bounds.Width, bounds.Height, config.DesignTokens)
	return SVGOutput(svg)
}

// RenderBarChart renders a bar chart as SVG
func (r *SVGRenderer) RenderBarChart(data BarChartData, bounds Bounds, config RenderConfig) Output {
	svg := RenderBarChart(data, bounds.X, bounds.Y, bounds.Width, bounds.Height, config.DesignTokens)
	return SVGOutput(svg)
}

// RenderStatCard renders a stat card as SVG
func (r *SVGRenderer) RenderStatCard(data StatCardData, bounds Bounds, config RenderConfig) Output {
	svg := RenderStatCard(data, bounds.X, bounds.Y, bounds.Width, bounds.Height, config.DesignTokens)
	return SVGOutput(svg)
}

// RenderAreaChart renders an area chart as SVG
func (r *SVGRenderer) RenderAreaChart(data AreaChartData, bounds Bounds, config RenderConfig) Output {
	svg := RenderAreaChart(data, bounds.X, bounds.Y, bounds.Width, bounds.Height, config.DesignTokens)
	return SVGOutput(svg)
}

// RenderScatterPlot renders a scatter plot as SVG
func (r *SVGRenderer) RenderScatterPlot(data ScatterPlotData, bounds Bounds, config RenderConfig) Output {
	svg := RenderScatterPlot(data, bounds.X, bounds.Y, bounds.Width, bounds.Height, config.DesignTokens)
	return SVGOutput(svg)
}
