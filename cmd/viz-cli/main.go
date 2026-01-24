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

  Hierarchical Charts:
    treemap       - Squarified treemap
    sunburst      - Radial partition chart
    circle-packing - Hierarchical circle packing
    icicle        - Icicle partition chart

  Statistical Charts:
    boxplot       - Box and whisker plot
    violin        - Violin plot with KDE
    histogram     - Histogram with binning
    ridgeline     - Ridgeline (joy) plot

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

	flag.StringVar(&cfg.vizType, "type", "heatmap", "Chart type")
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
	case "treemap":
		return renderTreemap(data, cfg)
	case "sunburst":
		return renderSunburst(data, cfg)
	case "circle-packing":
		return renderCirclePacking(data, cfg)
	case "icicle":
		return renderIcicle(data, cfg)
	case "boxplot":
		return renderBoxplot(data, cfg)
	case "violin":
		return renderViolin(data, cfg)
	case "histogram":
		return renderHistogram(data, cfg)
	case "ridgeline":
		return renderRidgeline(data, cfg)
	case "candlestick":
		return renderCandlestick(data, cfg)
	case "ohlc":
		return renderOHLC(data, cfg)
	case "scatter":
		return renderScatter(data, cfg)
	case "pie":
		return renderPie(data, cfg)
	case "area-chart":
		return renderArea(data, cfg)
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
