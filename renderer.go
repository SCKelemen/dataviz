package dataviz

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
}

// SVGRenderer implements SVG rendering
type SVGRenderer struct{}

// NewSVGRenderer creates a new SVG renderer
func NewSVGRenderer() *SVGRenderer {
	return &SVGRenderer{}
}

// RenderHeatmap renders a heatmap as SVG
func (r *SVGRenderer) RenderHeatmap(data HeatmapData, bounds Bounds, config RenderConfig) Output {
	// TODO: Implementation to be extracted from repobeats
	// Will call either RenderLinearHeatmap or RenderWeeksHeatmap based on data.Type
	return SVGOutput("")
}

// RenderLineGraph renders a line graph as SVG
func (r *SVGRenderer) RenderLineGraph(data LineGraphData, bounds Bounds, config RenderConfig) Output {
	// TODO: Implementation to be extracted from repobeats
	return SVGOutput("")
}

// RenderBarChart renders a bar chart as SVG
func (r *SVGRenderer) RenderBarChart(data BarChartData, bounds Bounds, config RenderConfig) Output {
	// TODO: Implementation to be extracted from repobeats
	return SVGOutput("")
}

// RenderStatCard renders a stat card as SVG
func (r *SVGRenderer) RenderStatCard(data StatCardData, bounds Bounds, config RenderConfig) Output {
	// TODO: Implementation to be extracted from repobeats
	return SVGOutput("")
}
