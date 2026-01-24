package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/SCKelemen/dataviz/charts"
	"github.com/SCKelemen/dataviz/scales"
	"github.com/SCKelemen/units"
	design "github.com/SCKelemen/design-system"
)

const usage = `viz-cli - Data visualization tool with SVG and terminal output

Usage:
  viz-cli [options]

Chart Types:
  Basic Charts:
    heatmap       - Contribution/calendar heatmap
    line-graph    - Line chart with optional smoothing
    bar-chart     - Bar chart with grouped/stacked support
    area-chart    - Area chart with fill
    scatter       - Scatter plot with various markers
    pie           - Pie or donut chart
    stat-card     - Single statistic card
    lollipop      - Lollipop chart (vertical/horizontal)

  Hierarchical Charts:
    treemap       - Squarified treemap
    sunburst      - Radial partition chart
    circle-packing - Hierarchical circle packing
    icicle        - Icicle partition chart
    dendrogram    - Hierarchical clustering tree

  Statistical Charts:
    boxplot       - Box and whisker plot
    violin        - Violin plot with KDE
    histogram     - Histogram with binning
    ridgeline     - Ridgeline (joy) plot
    density       - KDE density plot
    correlogram   - Correlation matrix

  Line/Area Charts:
    connected-scatter - Connected scatter plot
    stacked-area     - Stacked area chart
    streamchart      - Streamchart with flowing areas

  Specialized Charts:
    radar         - Radar/spider chart
    parallel      - Parallel coordinates
    wordcloud     - Word cloud visualization

  Network/Flow Charts:
    sankey        - Sankey flow diagram
    chord         - Chord diagram

  Circular Charts:
    circular-bar  - Circular bar plot

  Financial Charts:
    candlestick   - OHLC candlestick chart
    ohlc          - OHLC bar chart

Options:
  -type string
        Chart type (default "heatmap")
  -format string
        Output format: svg, terminal (default "terminal")
  -data string
        Path to JSON data file (or use stdin with -)
  -theme string
        Theme: default, midnight, nord, paper, wrapped (default "default")
  -width int
        Width in pixels (default 800)
  -height int
        Height in pixels (default 600)
  -color string
        Primary color (hex format) (default "#3B82F6")
  -output string
        Output file path (default: stdout)

Examples:
  # SVG treemap from file
  viz-cli -type treemap -format svg -data tree.json -output chart.svg

  # Terminal histogram from stdin
  cat data.json | viz-cli -type histogram -format terminal

  # Candlestick chart with custom theme
  viz-cli -type candlestick -data stocks.json -theme midnight -width 1200
`

type Config struct {
	vizType    string
	format     string
	dataFile   string
	outputFile string
	theme      string
	width      int
	height     int
	color      string
}

func main() {
	cfg := parseFlags()

	if cfg.dataFile == "" || cfg.dataFile == "-" {
		fmt.Fprintln(os.Stderr, "Reading from stdin...")
	}

	// Read data
	data, err := readData(cfg.dataFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading data: %v\n", err)
		os.Exit(1)
	}

	// Get design tokens
	tokens := getTheme(cfg.theme)

	// Render visualization
	var output string
	switch cfg.format {
	case "svg":
		output = renderSVG(cfg.vizType, data, cfg, tokens)
	case "terminal":
		output = renderTerminal(cfg.vizType, data, cfg, tokens)
	default:
		fmt.Fprintf(os.Stderr, "Unknown format: %s\n", cfg.format)
		os.Exit(1)
	}

	// Write output
	if err := writeOutput(cfg.outputFile, output); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
		os.Exit(1)
	}
}

func parseFlags() Config {
	cfg := Config{}

	flag.StringVar(&cfg.vizType, "type", "", "Chart type")
	flag.StringVar(&cfg.format, "format", "svg", "Output format")
	flag.StringVar(&cfg.dataFile, "data", "-", "Data file path")
	flag.StringVar(&cfg.outputFile, "output", "-", "Output file path")
	flag.StringVar(&cfg.theme, "theme", "default", "Theme name")
	flag.IntVar(&cfg.width, "width", 800, "Width in pixels")
	flag.IntVar(&cfg.height, "height", 600, "Height in pixels")
	flag.StringVar(&cfg.color, "color", "#3B82F6", "Primary color")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}

	flag.Parse()

	// If no -type flag was provided, check for positional argument
	if cfg.vizType == "" {
		args := flag.Args()
		if len(args) > 0 {
			cfg.vizType = args[0]
		} else {
			cfg.vizType = "heatmap" // default
		}
	}

	return cfg
}

func readData(path string) ([]byte, error) {
	if path == "" || path == "-" {
		return io.ReadAll(os.Stdin)
	}
	return os.ReadFile(path)
}

func writeOutput(path string, content string) error {
	if path == "" || path == "-" {
		fmt.Print(content)
		return nil
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func getTheme(name string) *design.DesignTokens {
	switch name {
	case "midnight":
		return design.MidnightTheme()
	case "nord":
		return design.NordTheme()
	case "paper":
		return design.PaperTheme()
	case "wrapped":
		return design.WrappedTheme()
	default:
		return design.DefaultTheme()
	}
}

func renderSVG(vizType string, data []byte, cfg Config, tokens *design.DesignTokens) string {
	svg := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`,
		cfg.width, cfg.height, cfg.width, cfg.height)
	svg += "\n"

	content := renderVisualization(vizType, data, cfg, tokens)
	svg += content
	svg += "\n</svg>"

	return svg
}

func renderTerminal(vizType string, data []byte, cfg Config, tokens *design.DesignTokens) string {
	// Terminal rendering not implemented for all types yet
	return "Terminal rendering not available for " + vizType
}

func renderVisualization(vizType string, data []byte, cfg Config, tokens *design.DesignTokens) string {
	switch vizType {
	// Hierarchical charts
	case "treemap":
		return renderTreemap(data, cfg)
	case "sunburst":
		return renderSunburst(data, cfg)
	case "circle-packing":
		return renderCirclePacking(data, cfg)
	case "icicle":
		return renderIcicle(data, cfg)
	case "dendrogram":
		return renderDendrogram(data, cfg)
	// Statistical charts
	case "boxplot":
		return renderBoxplot(data, cfg)
	case "violin":
		return renderViolin(data, cfg)
	case "histogram":
		return renderHistogram(data, cfg)
	case "ridgeline":
		return renderRidgeline(data, cfg)
	case "density":
		return renderDensity(data, cfg)
	case "correlogram":
		return renderCorrelogram(data, cfg)
	// Line/Area charts
	case "connected-scatter":
		return renderConnectedScatter(data, cfg)
	case "stacked-area":
		return renderStackedArea(data, cfg)
	case "streamchart":
		return renderStreamChart(data, cfg)
	// Specialized charts
	case "radar":
		return renderRadar(data, cfg)
	case "parallel":
		return renderParallel(data, cfg)
	case "wordcloud":
		return renderWordCloud(data, cfg)
	// Network/Flow charts
	case "sankey":
		return renderSankey(data, cfg)
	case "chord":
		return renderChord(data, cfg)
	// Circular charts
	case "circular-bar":
		return renderCircularBar(data, cfg)
	// Financial charts
	case "candlestick":
		return renderCandlestick(data, cfg)
	case "ohlc":
		return renderOHLC(data, cfg)
	// Basic charts
	case "scatter":
		return renderScatter(data, cfg)
	case "pie":
		return renderPie(data, cfg)
	case "area-chart":
		return renderArea(data, cfg)
	case "lollipop":
		return renderLollipop(data, cfg)
	case "heatmap", "line-graph", "bar-chart", "stat-card":
		return renderLegacyChart(vizType, data, cfg, tokens)
	default:
		fmt.Fprintf(os.Stderr, "Unknown chart type: %s\n", vizType)
		os.Exit(1)
		return ""
	}
}

// Hierarchical charts

func renderTreemap(data []byte, cfg Config) string {
	var treeData charts.TreeNode
	if err := json.Unmarshal(data, &treeData); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing tree data: %v\n", err)
		os.Exit(1)
	}

	spec := charts.TreemapSpec{
		Root:         &treeData,
		Width:        float64(cfg.width),
		Height:       float64(cfg.height),
		Padding:      2,
		ShowLabels:   true,
		MinLabelSize: 30,
	}

	return charts.RenderTreemap(spec)
}

func renderSunburst(data []byte, cfg Config) string {
	var treeData charts.TreeNode
	if err := json.Unmarshal(data, &treeData); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing tree data: %v\n", err)
		os.Exit(1)
	}

	spec := charts.SunburstSpec{
		Root:        &treeData,
		Width:       float64(cfg.width),
		Height:      float64(cfg.height),
		InnerRadius: float64(cfg.height) * 0.15,
		ShowLabels:  true,
		StartAngle:  0,
	}

	return charts.RenderSunburst(spec)
}

func renderCirclePacking(data []byte, cfg Config) string {
	var treeData charts.TreeNode
	if err := json.Unmarshal(data, &treeData); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing tree data: %v\n", err)
		os.Exit(1)
	}

	spec := charts.CirclePackingSpec{
		Root:       &treeData,
		Width:      float64(cfg.width),
		Height:     float64(cfg.height),
		Padding:    2,
		ShowLabels: true,
	}

	return charts.RenderCirclePacking(spec)
}

func renderIcicle(data []byte, cfg Config) string {
	var treeData charts.TreeNode
	if err := json.Unmarshal(data, &treeData); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing tree data: %v\n", err)
		os.Exit(1)
	}

	spec := charts.IcicleSpec{
		Root:        &treeData,
		Width:       float64(cfg.width),
		Height:      float64(cfg.height),
		Padding:     2,
		Orientation: "vertical",
		ShowLabels:  true,
	}

	return charts.RenderIcicle(spec)
}

// Statistical charts

func renderBoxplot(data []byte, cfg Config) string {
	var boxInput struct {
		Data []struct {
			Values []float64 `json:"values"`
			Label  string    `json:"label"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &boxInput); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing boxplot data: %v\n", err)
		os.Exit(1)
	}

	boxData := make([]*charts.BoxPlotData, len(boxInput.Data))
	for i, d := range boxInput.Data {
		boxData[i] = &charts.BoxPlotData{
			Values: d.Values,
			Label:  d.Label,
			Color:  cfg.color,
		}
	}

	spec := charts.BoxPlotSpec{
		Data:              boxData,
		Width:             float64(cfg.width),
		Height:            float64(cfg.height),
		Horizontal:        false,
		ShowOutliers:      true,
		ShowMean:          false,
		WhiskerMultiplier: 1.5,
	}

	return charts.RenderVerticalBoxPlot(spec)
}

func renderViolin(data []byte, cfg Config) string {
	var violinInput struct {
		Data []struct {
			Values []float64 `json:"values"`
			Label  string    `json:"label"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &violinInput); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing violin data: %v\n", err)
		os.Exit(1)
	}

	violinData := make([]*charts.ViolinPlotData, len(violinInput.Data))
	for i, d := range violinInput.Data {
		violinData[i] = &charts.ViolinPlotData{
			Values: d.Values,
			Label:  d.Label,
			Color:  cfg.color,
		}
	}

	spec := charts.ViolinPlotSpec{
		Data:       violinData,
		Width:      float64(cfg.width),
		Height:     float64(cfg.height),
		Bandwidth:  0, // Auto-calculate
		ShowBox:    true,
		ShowMedian: true,
		ShowMean:   false,
	}

	return charts.RenderViolinPlot(spec)
}

func renderHistogram(data []byte, cfg Config) string {
	var histInput struct {
		Values []float64 `json:"values"`
		Bins   int       `json:"bins"`
	}
	if err := json.Unmarshal(data, &histInput); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing histogram data: %v\n", err)
		os.Exit(1)
	}

	if histInput.Bins == 0 {
		histInput.Bins = 20
	}

	spec := charts.HistogramSpec{
		Data: &charts.HistogramData{
			Values: histInput.Values,
			Color:  cfg.color,
		},
		Width:    float64(cfg.width),
		Height:   float64(cfg.height),
		BinCount: histInput.Bins,
		Nice:     true,
	}

	return charts.RenderHistogram(spec)
}

func renderRidgeline(data []byte, cfg Config) string {
	var ridgeInput struct {
		Data []struct {
			Label  string    `json:"label"`
			Values []float64 `json:"values"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &ridgeInput); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing ridgeline data: %v\n", err)
		os.Exit(1)
	}

	ridgeData := make([]*charts.RidgelineData, len(ridgeInput.Data))
	for i, d := range ridgeInput.Data {
		ridgeData[i] = &charts.RidgelineData{
			Label:  d.Label,
			Values: d.Values,
		}
	}

	spec := charts.RidgelineSpec{
		Data:       ridgeData,
		Width:      float64(cfg.width),
		Height:     float64(cfg.height),
		Overlap:    0.5,
		ShowFill:   true,
		ShowLabels: true,
	}

	return charts.RenderRidgeline(spec)
}

// Financial charts

func renderCandlestick(data []byte, cfg Config) string {
	var candleData []charts.CandlestickData
	if err := json.Unmarshal(data, &candleData); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing candlestick data: %v\n", err)
		os.Exit(1)
	}

	// Find min/max for scales
	minPrice, maxPrice := candleData[0].Low, candleData[0].High
	for _, d := range candleData {
		if d.Low < minPrice {
			minPrice = d.Low
		}
		if d.High > maxPrice {
			maxPrice = d.High
		}
	}

	xScale := scales.NewLinearScale(
		[2]float64{0, float64(len(candleData))},
		[2]units.Length{units.Px(50), units.Px(float64(cfg.width) - 50)},
	)

	yScale := scales.NewLinearScale(
		[2]float64{minPrice * 0.98, maxPrice * 1.02},
		[2]units.Length{units.Px(float64(cfg.height) - 100), units.Px(50)},
	)

	spec := charts.CandlestickSpec{
		Data:         candleData,
		Width:        float64(cfg.width),
		Height:       float64(cfg.height),
		XScale:       xScale,
		YScale:       yScale,
		ShowVolume:   true,
		VolumeHeight: 100,
	}

	return charts.RenderCandlestick(spec)
}

func renderOHLC(data []byte, cfg Config) string {
	var ohlcData []charts.OHLCData
	if err := json.Unmarshal(data, &ohlcData); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing OHLC data: %v\n", err)
		os.Exit(1)
	}

	// Find min/max for scales
	minPrice, maxPrice := ohlcData[0].Low, ohlcData[0].High
	for _, d := range ohlcData {
		if d.Low < minPrice {
			minPrice = d.Low
		}
		if d.High > maxPrice {
			maxPrice = d.High
		}
	}

	xScale := scales.NewLinearScale(
		[2]float64{0, float64(len(ohlcData))},
		[2]units.Length{units.Px(50), units.Px(float64(cfg.width) - 50)},
	)

	yScale := scales.NewLinearScale(
		[2]float64{minPrice * 0.98, maxPrice * 1.02},
		[2]units.Length{units.Px(float64(cfg.height) - 50), units.Px(50)},
	)

	spec := charts.OHLCSpec{
		Data:   ohlcData,
		Width:  float64(cfg.width),
		Height: float64(cfg.height),
		XScale: xScale,
		YScale: yScale,
	}

	return charts.RenderOHLC(spec)
}

// Basic charts

func renderScatter(data []byte, cfg Config) string {
	// Scatter implementation placeholder
	return fmt.Sprintf(`<text x="%d" y="%d" text-anchor="middle">Scatter plot: Use MCP server for full support</text>`,
		cfg.width/2, cfg.height/2)
}

func renderPie(data []byte, cfg Config) string {
	// Pie implementation placeholder
	return fmt.Sprintf(`<text x="%d" y="%d" text-anchor="middle">Pie chart: Use MCP server for full support</text>`,
		cfg.width/2, cfg.height/2)
}

func renderArea(data []byte, cfg Config) string {
	// Area implementation placeholder
	return fmt.Sprintf(`<text x="%d" y="%d" text-anchor="middle">Area chart: Use MCP server for full support</text>`,
		cfg.width/2, cfg.height/2)
}

func renderLegacyChart(vizType string, data []byte, cfg Config, tokens *design.DesignTokens) string {
	// Legacy chart rendering using old renderer system
	bounds := charts.Bounds{X: 0, Y: 0, Width: cfg.width, Height: cfg.height}
	renderConfig := charts.RenderConfig{
		DesignTokens: tokens,
		Color:        cfg.color,
		Theme:        cfg.theme,
	}

	renderer := charts.NewSVGRenderer()

	switch vizType {
	case "heatmap":
		var heatmapData charts.HeatmapData
		if err := json.Unmarshal(data, &heatmapData); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing heatmap data: %v\n", err)
			os.Exit(1)
		}
		return renderer.RenderHeatmap(heatmapData, bounds, renderConfig).String()

	case "line-graph":
		var lineData charts.LineGraphData
		if err := json.Unmarshal(data, &lineData); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing line graph data: %v\n", err)
			os.Exit(1)
		}
		return renderer.RenderLineGraph(lineData, bounds, renderConfig).String()

	case "bar-chart":
		var barData charts.BarChartData
		if err := json.Unmarshal(data, &barData); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing bar chart data: %v\n", err)
			os.Exit(1)
		}
		return renderer.RenderBarChart(barData, bounds, renderConfig).String()

	case "stat-card":
		var statData charts.StatCardData
		if err := json.Unmarshal(data, &statData); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing stat card data: %v\n", err)
			os.Exit(1)
		}
		return renderer.RenderStatCard(statData, bounds, renderConfig).String()

	default:
		return ""
	}
}

// New chart types

func renderLollipop(data []byte, cfg Config) string {
	var input struct {
		Values []struct {
			Label string  `json:"label"`
			Value float64 `json:"value"`
			Color string  `json:"color,omitempty"`
		} `json:"values"`
		Color      string `json:"color,omitempty"`
		Horizontal bool   `json:"horizontal,omitempty"`
	}
	if err := json.Unmarshal(data, &input); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing lollipop data: %v\n", err)
		os.Exit(1)
	}

	lollipopData := &charts.LollipopData{
		Values: make([]charts.LollipopPoint, len(input.Values)),
		Color:  cfg.color,
	}
	if input.Color != "" {
		lollipopData.Color = input.Color
	}

	for i, v := range input.Values {
		lollipopData.Values[i] = charts.LollipopPoint{
			Label: v.Label,
			Value: v.Value,
			Color: v.Color,
		}
	}

	spec := charts.LollipopSpec{
		Data:       lollipopData,
		Width:      float64(cfg.width),
		Height:     float64(cfg.height),
		Horizontal: input.Horizontal,
		ShowLabels: true,
		ShowGrid:   true,
	}

	return charts.RenderLollipop(spec)
}

func renderDensity(data []byte, cfg Config) string {
	var input struct {
		Data []struct {
			Values []float64 `json:"values"`
			Label  string    `json:"label,omitempty"`
			Color  string    `json:"color,omitempty"`
		} `json:"data"`
		ShowFill bool `json:"show_fill,omitempty"`
		ShowRug  bool `json:"show_rug,omitempty"`
	}
	if err := json.Unmarshal(data, &input); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing density data: %v\n", err)
		os.Exit(1)
	}

	densityData := make([]*charts.SimpleDensityData, len(input.Data))
	for i, d := range input.Data {
		densityData[i] = &charts.SimpleDensityData{
			Values: d.Values,
			Label:  d.Label,
			Color:  d.Color,
		}
	}

	spec := charts.SimpleDensitySpec{
		Data:     densityData,
		Width:    float64(cfg.width),
		Height:   float64(cfg.height),
		ShowFill: input.ShowFill,
		ShowRug:  input.ShowRug,
	}

	return charts.RenderSimpleDensity(spec)
}

func renderConnectedScatter(data []byte, cfg Config) string {
	var input struct {
		Series []struct {
			Points []struct {
				X     float64 `json:"x"`
				Y     float64 `json:"y"`
				Label string  `json:"label,omitempty"`
			} `json:"points"`
			Label      string `json:"label,omitempty"`
			Color      string `json:"color,omitempty"`
			MarkerType string `json:"marker_type,omitempty"`
		} `json:"series"`
	}
	if err := json.Unmarshal(data, &input); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing connected scatter data: %v\n", err)
		os.Exit(1)
	}

	series := make([]*charts.ConnectedScatterSeries, len(input.Series))
	for i, s := range input.Series {
		points := make([]charts.ConnectedScatterPoint, len(s.Points))
		for j, p := range s.Points {
			points[j] = charts.ConnectedScatterPoint{
				X:     p.X,
				Y:     p.Y,
				Label: p.Label,
			}
		}
		series[i] = &charts.ConnectedScatterSeries{
			Points:     points,
			Label:      s.Label,
			Color:      s.Color,
			MarkerType: s.MarkerType,
		}
	}

	spec := charts.ConnectedScatterSpec{
		Series:      series,
		Width:       float64(cfg.width),
		Height:      float64(cfg.height),
		ShowGrid:    true,
		ShowMarkers: true,
		ShowLines:   true,
	}

	return charts.RenderConnectedScatter(spec)
}

func renderStackedArea(data []byte, cfg Config) string {
	var input struct {
		Points []struct {
			X      float64   `json:"x"`
			Values []float64 `json:"values"`
		} `json:"points"`
		Series []struct {
			Label string `json:"label"`
			Color string `json:"color,omitempty"`
		} `json:"series"`
	}
	if err := json.Unmarshal(data, &input); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing stacked area data: %v\n", err)
		os.Exit(1)
	}

	points := make([]charts.StackedAreaPoint, len(input.Points))
	for i, p := range input.Points {
		points[i] = charts.StackedAreaPoint{
			X:      p.X,
			Values: p.Values,
		}
	}

	series := make([]charts.StackedAreaSeries, len(input.Series))
	for i, s := range input.Series {
		series[i] = charts.StackedAreaSeries{
			Label: s.Label,
			Color: s.Color,
		}
	}

	spec := charts.StackedAreaSpec{
		Points:   points,
		Series:   series,
		Width:    float64(cfg.width),
		Height:   float64(cfg.height),
		ShowGrid: true,
	}

	return charts.RenderStackedArea(spec)
}

func renderStreamChart(data []byte, cfg Config) string {
	var input struct {
		Points []struct {
			X      float64   `json:"x"`
			Values []float64 `json:"values"`
		} `json:"points"`
		Series []struct {
			Label string `json:"label"`
			Color string `json:"color,omitempty"`
		} `json:"series"`
		Layout string `json:"layout,omitempty"`
	}
	if err := json.Unmarshal(data, &input); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing stream chart data: %v\n", err)
		os.Exit(1)
	}

	points := make([]charts.StreamPoint, len(input.Points))
	for i, p := range input.Points {
		points[i] = charts.StreamPoint{
			X:      p.X,
			Values: p.Values,
		}
	}

	series := make([]charts.StreamSeries, len(input.Series))
	for i, s := range input.Series {
		series[i] = charts.StreamSeries{
			Label: s.Label,
			Color: s.Color,
		}
	}

	spec := charts.StreamChartSpec{
		Points:     points,
		Series:     series,
		Width:      float64(cfg.width),
		Height:     float64(cfg.height),
		Layout:     input.Layout,
		ShowLegend: true,
	}

	return charts.RenderStreamChart(spec)
}

func renderCorrelogram(data []byte, cfg Config) string {
	var input struct {
		Variables []string    `json:"variables"`
		Matrix    [][]float64 `json:"matrix"`
	}
	if err := json.Unmarshal(data, &input); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing correlogram data: %v\n", err)
		os.Exit(1)
	}

	matrix := charts.CorrelationMatrix{
		Variables: input.Variables,
		Matrix:    input.Matrix,
	}

	spec := charts.CorrelogramSpec{
		Data:         matrix,
		Width:        float64(cfg.width),
		Height:       float64(cfg.height),
		ShowValues:   true,
		ShowDiagonal: true,
		TriangleMode: "full",
		ColorScheme:  "redblue",
	}

	return charts.RenderCorrelogram(spec)
}

func renderRadar(data []byte, cfg Config) string {
	var input struct {
		Axes []struct {
			Label string  `json:"label"`
			Min   float64 `json:"min"`
			Max   float64 `json:"max"`
		} `json:"axes"`
		Series []struct {
			Label  string    `json:"label"`
			Values []float64 `json:"values"`
			Color  string    `json:"color,omitempty"`
		} `json:"series"`
	}
	if err := json.Unmarshal(data, &input); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing radar data: %v\n", err)
		os.Exit(1)
	}

	axes := make([]charts.RadarAxis, len(input.Axes))
	for i, a := range input.Axes {
		axes[i] = charts.RadarAxis{
			Label: a.Label,
			Min:   a.Min,
			Max:   a.Max,
		}
	}

	series := make([]*charts.RadarSeries, len(input.Series))
	for i, s := range input.Series {
		series[i] = &charts.RadarSeries{
			Label:  s.Label,
			Values: s.Values,
			Color:  s.Color,
		}
	}

	spec := charts.RadarChartSpec{
		Axes:       axes,
		Series:     series,
		Width:      float64(cfg.width),
		Height:     float64(cfg.height),
		ShowGrid:   true,
		ShowLabels: true,
		GridLevels: 5,
	}

	return charts.RenderRadarChart(spec)
}

func renderParallel(data []byte, cfg Config) string {
	var input struct {
		Axes []struct {
			Label string  `json:"label"`
			Min   float64 `json:"min"`
			Max   float64 `json:"max"`
		} `json:"axes"`
		Data []struct {
			Values []float64 `json:"values"`
			Color  string    `json:"color,omitempty"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &input); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing parallel data: %v\n", err)
		os.Exit(1)
	}

	axes := make([]charts.ParallelAxis, len(input.Axes))
	for i, a := range input.Axes {
		axes[i] = charts.ParallelAxis{
			Label: a.Label,
			Min:   a.Min,
			Max:   a.Max,
		}
	}

	dataPoints := make([]charts.ParallelDataPoint, len(input.Data))
	for i, d := range input.Data {
		dataPoints[i] = charts.ParallelDataPoint{
			Values: d.Values,
			Color:  d.Color,
		}
	}

	spec := charts.ParallelCoordinatesSpec{
		Axes:           axes,
		Data:           dataPoints,
		Width:          float64(cfg.width),
		Height:         float64(cfg.height),
		ShowAxesLabels: true,
		ShowTicks:      true,
	}

	return charts.RenderParallelCoordinates(spec)
}

func renderWordCloud(data []byte, cfg Config) string {
	var input struct {
		Words []struct {
			Text      string  `json:"text"`
			Frequency float64 `json:"frequency"`
			Color     string  `json:"color,omitempty"`
		} `json:"words"`
		Layout string `json:"layout,omitempty"`
	}
	if err := json.Unmarshal(data, &input); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing wordcloud data: %v\n", err)
		os.Exit(1)
	}

	words := make([]charts.WordCloudWord, len(input.Words))
	for i, w := range input.Words {
		words[i] = charts.WordCloudWord{
			Text:      w.Text,
			Frequency: w.Frequency,
			Color:     w.Color,
		}
	}

	spec := charts.WordCloudSpec{
		Words:  words,
		Width:  float64(cfg.width),
		Height: float64(cfg.height),
		Layout: input.Layout,
	}

	return charts.RenderWordCloud(spec)
}

func renderSankey(data []byte, cfg Config) string {
	var input struct {
		Nodes []struct {
			ID    string `json:"id"`
			Label string `json:"label"`
			Color string `json:"color,omitempty"`
		} `json:"nodes"`
		Links []struct {
			Source string  `json:"source"`
			Target string  `json:"target"`
			Value  float64 `json:"value"`
			Color  string  `json:"color,omitempty"`
		} `json:"links"`
	}
	if err := json.Unmarshal(data, &input); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing sankey data: %v\n", err)
		os.Exit(1)
	}

	nodes := make([]charts.SankeyNode, len(input.Nodes))
	for i, n := range input.Nodes {
		nodes[i] = charts.SankeyNode{
			ID:    n.ID,
			Label: n.Label,
			Color: n.Color,
		}
	}

	links := make([]charts.SankeyLink, len(input.Links))
	for i, l := range input.Links {
		links[i] = charts.SankeyLink{
			Source: l.Source,
			Target: l.Target,
			Value:  l.Value,
			Color:  l.Color,
		}
	}

	spec := charts.SankeySpec{
		Nodes:      nodes,
		Links:      links,
		Width:      float64(cfg.width),
		Height:     float64(cfg.height),
		ShowLabels: true,
	}

	return charts.RenderSankey(spec)
}

func renderChord(data []byte, cfg Config) string {
	var input struct {
		Entities []struct {
			ID    string `json:"id"`
			Label string `json:"label"`
			Color string `json:"color,omitempty"`
		} `json:"entities"`
		Relations []struct {
			Source string  `json:"source"`
			Target string  `json:"target"`
			Value  float64 `json:"value"`
		} `json:"relations"`
	}
	if err := json.Unmarshal(data, &input); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing chord data: %v\n", err)
		os.Exit(1)
	}

	entities := make([]charts.ChordEntity, len(input.Entities))
	for i, e := range input.Entities {
		entities[i] = charts.ChordEntity{
			ID:    e.ID,
			Label: e.Label,
			Color: e.Color,
		}
	}

	relations := make([]charts.ChordRelation, len(input.Relations))
	for i, r := range input.Relations {
		relations[i] = charts.ChordRelation{
			Source: r.Source,
			Target: r.Target,
			Value:  r.Value,
		}
	}

	spec := charts.ChordDiagramSpec{
		Entities:   entities,
		Relations:  relations,
		Width:      float64(cfg.width),
		Height:     float64(cfg.height),
		ShowLabels: true,
	}

	return charts.RenderChordDiagram(spec)
}

func renderCircularBar(data []byte, cfg Config) string {
	var input struct {
		Data []struct {
			Label string  `json:"label"`
			Value float64 `json:"value"`
			Color string  `json:"color,omitempty"`
		} `json:"data"`
		InnerRadius float64 `json:"inner_radius,omitempty"`
	}
	if err := json.Unmarshal(data, &input); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing circular bar data: %v\n", err)
		os.Exit(1)
	}

	barData := make([]charts.CircularBarPoint, len(input.Data))
	for i, d := range input.Data {
		barData[i] = charts.CircularBarPoint{
			Label: d.Label,
			Value: d.Value,
			Color: d.Color,
		}
	}

	spec := charts.CircularBarPlotSpec{
		Data:           barData,
		Width:          float64(cfg.width),
		Height:         float64(cfg.height),
		InnerRadius:    input.InnerRadius,
		ShowLabels:     false,
		ShowAxisLabels: true,
	}

	return charts.RenderCircularBarPlot(spec)
}

func renderDendrogram(data []byte, cfg Config) string {
	var input struct {
		Root struct {
			Label    string `json:"label,omitempty"`
			Height   float64 `json:"height"`
			Children []json.RawMessage `json:"children,omitempty"`
		} `json:"root"`
		Orientation string `json:"orientation,omitempty"`
	}
	if err := json.Unmarshal(data, &input); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing dendrogram data: %v\n", err)
		os.Exit(1)
	}

	// Recursive function to parse dendrogram nodes
	var parseNode func(raw json.RawMessage) *charts.DendrogramNode
	parseNode = func(raw json.RawMessage) *charts.DendrogramNode {
		var node struct {
			Label    string            `json:"label,omitempty"`
			Height   float64           `json:"height"`
			Children []json.RawMessage `json:"children,omitempty"`
		}
		if err := json.Unmarshal(raw, &node); err != nil {
			return nil
		}

		dendNode := &charts.DendrogramNode{
			Label:  node.Label,
			Height: node.Height,
		}

		if len(node.Children) > 0 {
			dendNode.Children = make([]*charts.DendrogramNode, len(node.Children))
			for i, child := range node.Children {
				dendNode.Children[i] = parseNode(child)
			}
		}

		return dendNode
	}

	rootData, _ := json.Marshal(input.Root)
	root := parseNode(rootData)

	spec := charts.DendrogramSpec{
		Root:        root,
		Width:       float64(cfg.width),
		Height:      float64(cfg.height),
		Orientation: input.Orientation,
		ShowLabels:  true,
		ShowHeights: true,
	}

	return charts.RenderDendrogram(spec)
}

// Sample data generators for testing

func createSampleTreeData() charts.TreeNode {
	root := charts.NewTreeNode("Root", 0)

	software := charts.NewTreeNode("Software", 0)
	software.AddChild(charts.NewTreeNode("Frontend", 50))
	software.AddChild(charts.NewTreeNode("Backend", 75))
	software.AddChild(charts.NewTreeNode("Database", 30))

	hardware := charts.NewTreeNode("Hardware", 0)
	hardware.AddChild(charts.NewTreeNode("Servers", 40))
	hardware.AddChild(charts.NewTreeNode("Storage", 60))

	root.AddChild(software).AddChild(hardware)

	return *root
}

func createSampleCandlestickData() []charts.CandlestickData {
	data := make([]charts.CandlestickData, 20)
	for i := range data {
		base := 100.0 + float64(i)*0.5
		data[i] = charts.CandlestickData{
			X:      i,
			Open:   base,
			High:   base + 5,
			Low:    base - 3,
			Close:  base + 2,
			Volume: 1000 + float64(i)*50,
		}
	}
	return data
}
