package charts

import (
	"fmt"
	"math"
	"strings"

	"github.com/SCKelemen/color"
	design "github.com/SCKelemen/design-system"
	maincharts "github.com/SCKelemen/dataviz/charts"
	"github.com/SCKelemen/dataviz/mcp/types"
	"github.com/SCKelemen/dataviz/scales"
	"github.com/SCKelemen/layout"
	"github.com/SCKelemen/svg"
	"github.com/SCKelemen/units"
)

// CreateBarChart generates a bar chart SVG by calling the main library
func CreateBarChart(config types.BarChartConfig) (string, error) {
	// Convert MCP types to main library types
	bars := make([]maincharts.BarData, len(config.Data))
	for i, dp := range config.Data {
		bars[i] = maincharts.BarData{
			Label: dp.Label,
			Value: int(dp.Value), // Convert float64 to int
		}
	}

	data := maincharts.BarChartData{
		Bars:  bars,
		Color: config.Color,
	}

	// Use default theme
	tokens := design.DefaultTheme()

	// Call main library function
	svg := maincharts.RenderBarChart(data, 0, 0, config.Width, config.Height, tokens)

	return svg, nil
}

// CreatePieChart generates a pie chart SVG by calling the main library
func CreatePieChart(config types.PieChartConfig) (string, error) {
	// Convert MCP types to main library types
	slices := make([]maincharts.PieSlice, len(config.Data))
	for i, dp := range config.Data {
		slices[i] = maincharts.PieSlice{
			Label: dp.Label,
			Value: dp.Value,
		}
	}

	data := maincharts.PieChartData{
		Slices: slices,
	}

	// Call main library function
	// Note: config.Donut maps to donutMode parameter
	// showLegend=true, showPercent=true for MCP compatibility
	svg := maincharts.RenderPieChart(data, 0, 0, config.Width, config.Height, config.Title, config.Donut, true, true)

	return svg, nil
}

// CreateLineChart generates a line chart SVG using SCKelemen libraries
func CreateLineChart(config types.LineChartConfig) (string, error) {
	if len(config.Series) == 0 {
		return "", fmt.Errorf("no series data provided")
	}

	// Calculate data ranges
	minY, maxY := math.MaxFloat64, -math.MaxFloat64
	maxPoints := 0

	for _, series := range config.Series {
		for _, point := range series.Data {
			if point.Y < minY {
				minY = point.Y
			}
			if point.Y > maxY {
				maxY = point.Y
			}
		}
		if len(series.Data) > maxPoints {
			maxPoints = len(series.Data)
		}
	}

	// Chart dimensions
	margin := 60.0
	chartWidth := float64(config.Width) - (2 * margin)
	chartHeight := float64(config.Height) - (2 * margin)

	// Start building SVG
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d">`,
		config.Width, config.Height))
	sb.WriteString("\n")

	// Background
	sb.WriteString(fmt.Sprintf(`  <rect width="%d" height="%d" fill="#ffffff"/>`, config.Width, config.Height))
	sb.WriteString("\n")

	// Title
	if config.Title != "" {
		sb.WriteString(fmt.Sprintf(`  <text x="%d" y="30" text-anchor="middle" font-size="20" font-weight="bold" fill="#1f2937">%s</text>`,
			config.Width/2, config.Title))
		sb.WriteString("\n")
	}

	// Color palette
	colors := []string{"#3b82f6", "#10b981", "#f59e0b", "#ef4444", "#8b5cf6"}

	// Draw axes
	sb.WriteString(fmt.Sprintf(`  <line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="#d1d5db" stroke-width="2"/>`,
		margin, margin, margin, margin+chartHeight))
	sb.WriteString(fmt.Sprintf(`  <line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="#d1d5db" stroke-width="2"/>`,
		margin, margin+chartHeight, margin+chartWidth, margin+chartHeight))
	sb.WriteString("\n")

	// Draw grid and Y-axis labels
	steps := 5
	for i := 0; i <= steps; i++ {
		value := minY + ((maxY - minY) / float64(steps) * float64(i))
		y := margin + chartHeight - (chartHeight/float64(steps))*float64(i)

		sb.WriteString(fmt.Sprintf(`  <line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="#e5e7eb" stroke-width="1" stroke-dasharray="4,4"/>`,
			margin, y, margin+chartWidth, y))
		sb.WriteString(fmt.Sprintf(`  <text x="%.2f" y="%.2f" text-anchor="end" font-size="11" fill="#6b7280">%.1f</text>`,
			margin-10, y+4, value))
		sb.WriteString("\n")
	}

	// Draw series
	for seriesIdx, series := range config.Series {
		if len(series.Data) == 0 {
			continue
		}

		seriesColor := series.Color
		if seriesColor == "" {
			seriesColor = colors[seriesIdx%len(colors)]
		}

		// Build path using smooth curves
		points := make([]svg.Point, len(series.Data))
		for i, point := range series.Data {
			x := margin + (chartWidth / float64(maxPoints-1) * float64(i))
			y := margin + chartHeight - ((point.Y-minY)/(maxY-minY))*chartHeight
			points[i] = svg.Point{X: x, Y: y}
		}

		pathData := svg.SmoothLinePath(points, 0.3)

		sb.WriteString(fmt.Sprintf(`  <path d="%s" fill="none" stroke="%s" stroke-width="2"/>`,
			pathData, seriesColor))
		sb.WriteString("\n")

		// Draw points
		for _, p := range points {
			sb.WriteString(fmt.Sprintf(`  <circle cx="%.2f" cy="%.2f" r="4" fill="%s" stroke="#ffffff" stroke-width="2"/>`,
				p.X, p.Y, seriesColor))
			sb.WriteString("\n")
		}
	}

	// Draw legend
	legendX := margin
	legendY := 50.0
	for i, series := range config.Series {
		seriesColor := series.Color
		if seriesColor == "" {
			seriesColor = colors[i%len(colors)]
		}

		xOffset := float64(i * 120)
		sb.WriteString(fmt.Sprintf(`  <line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="%s" stroke-width="2"/>`,
			legendX+xOffset, legendY, legendX+xOffset+20, legendY, seriesColor))
		sb.WriteString(fmt.Sprintf(`  <text x="%.2f" y="%.2f" font-size="12" fill="#374151">%s</text>`,
			legendX+xOffset+25, legendY+4, series.Name))
		sb.WriteString("\n")
	}

	sb.WriteString("</svg>")

	return sb.String(), nil
}

// CreateScatterPlot generates a scatter plot SVG using SCKelemen libraries
func CreateScatterPlot(config types.ScatterPlotConfig) (string, error) {
	if len(config.Data) == 0 {
		return "", fmt.Errorf("no data provided")
	}

	// Calculate data ranges
	minX, maxX := math.MaxFloat64, -math.MaxFloat64
	minY, maxY := math.MaxFloat64, -math.MaxFloat64

	for _, point := range config.Data {
		if point.X < minX {
			minX = point.X
		}
		if point.X > maxX {
			maxX = point.X
		}
		if point.Y < minY {
			minY = point.Y
		}
		if point.Y > maxY {
			maxY = point.Y
		}
	}

	// Chart dimensions
	margin := 60.0
	chartWidth := float64(config.Width) - (2 * margin)
	chartHeight := float64(config.Height) - (2 * margin)

	// Start building SVG
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d">`,
		config.Width, config.Height))
	sb.WriteString("\n")

	// Background
	sb.WriteString(fmt.Sprintf(`  <rect width="%d" height="%d" fill="#ffffff"/>`, config.Width, config.Height))
	sb.WriteString("\n")

	// Title
	if config.Title != "" {
		sb.WriteString(fmt.Sprintf(`  <text x="%d" y="30" text-anchor="middle" font-size="20" font-weight="bold" fill="#1f2937">%s</text>`,
			config.Width/2, config.Title))
		sb.WriteString("\n")
	}

	// Draw axes
	sb.WriteString(fmt.Sprintf(`  <line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="#d1d5db" stroke-width="2"/>`,
		margin, margin, margin, margin+chartHeight))
	sb.WriteString(fmt.Sprintf(`  <line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="#d1d5db" stroke-width="2"/>`,
		margin, margin+chartHeight, margin+chartWidth, margin+chartHeight))
	sb.WriteString("\n")

	// Draw points
	for _, point := range config.Data {
		x := margin + ((point.X-minX)/(maxX-minX))*chartWidth
		y := margin + chartHeight - ((point.Y-minY)/(maxY-minY))*chartHeight

		radius := 5.0
		if point.Size > 0 {
			radius = math.Min(point.Size, 15)
		}

		sb.WriteString(fmt.Sprintf(`  <circle cx="%.2f" cy="%.2f" r="%.2f" fill="#3b82f6" fill-opacity="0.6" stroke="#2563eb" stroke-width="1"/>`,
			x, y, radius))
		sb.WriteString("\n")
	}

	// Axis labels
	if config.XLabel != "" {
		sb.WriteString(fmt.Sprintf(`  <text x="%d" y="%d" text-anchor="middle" font-size="14" fill="#374151">%s</text>`,
			config.Width/2, config.Height-10, config.XLabel))
		sb.WriteString("\n")
	}
	if config.YLabel != "" {
		sb.WriteString(fmt.Sprintf(`  <text x="15" y="%d" text-anchor="middle" font-size="14" fill="#374151" transform="rotate(-90 15 %d)">%s</text>`,
			config.Height/2, config.Height/2, config.YLabel))
		sb.WriteString("\n")
	}

	sb.WriteString("</svg>")

	return sb.String(), nil
}

// CreateHeatmap generates a heatmap SVG using SCKelemen libraries
func CreateHeatmap(config types.HeatmapConfig) (string, error) {
	rows := len(config.Data.Rows)
	cols := len(config.Data.Columns)

	if rows == 0 || cols == 0 {
		return "", fmt.Errorf("empty heatmap data")
	}

	// Find min/max for color scaling
	minVal, maxVal := math.MaxFloat64, -math.MaxFloat64
	for _, row := range config.Data.Values {
		for _, val := range row {
			if val < minVal {
				minVal = val
			}
			if val > maxVal {
				maxVal = val
			}
		}
	}

	// Chart dimensions
	margin := 80.0
	cellSize := math.Min(
		(float64(config.Width)-2*margin)/float64(cols),
		(float64(config.Height)-2*margin)/float64(rows),
	)

	// Start building SVG
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d">`,
		config.Width, config.Height))
	sb.WriteString("\n")

	// Background
	sb.WriteString(fmt.Sprintf(`  <rect width="%d" height="%d" fill="#ffffff"/>`, config.Width, config.Height))
	sb.WriteString("\n")

	// Title
	if config.Title != "" {
		sb.WriteString(fmt.Sprintf(`  <text x="%d" y="30" text-anchor="middle" font-size="20" font-weight="bold" fill="#1f2937">%s</text>`,
			config.Width/2, config.Title))
		sb.WriteString("\n")
	}

	// Draw cells
	for i, row := range config.Data.Values {
		for j, val := range row {
			x := margin + float64(j)*cellSize
			y := margin + float64(i)*cellSize

			// Color based on value (viridis-like gradient)
			normalized := (val - minVal) / (maxVal - minVal)
			cellColor := interpolateColor(normalized)

			sb.WriteString(fmt.Sprintf(`  <rect x="%.2f" y="%.2f" width="%.2f" height="%.2f" fill="%s" stroke="#ffffff" stroke-width="1"/>`,
				x, y, cellSize, cellSize, cellColor))

			// Show value if enabled
			if config.ShowValue {
				textX := x + cellSize/2
				textY := y + cellSize/2 + 4
				sb.WriteString(fmt.Sprintf(`  <text x="%.2f" y="%.2f" text-anchor="middle" font-size="10" fill="#ffffff">%.1f</text>`,
					textX, textY, val))
			}

			sb.WriteString("\n")
		}
	}

	// Column labels
	for j, col := range config.Data.Columns {
		x := margin + float64(j)*cellSize + cellSize/2
		y := margin - 10
		sb.WriteString(fmt.Sprintf(`  <text x="%.2f" y="%.2f" text-anchor="middle" font-size="11" fill="#374151">%s</text>`,
			x, y, col))
		sb.WriteString("\n")
	}

	// Row labels
	for i, row := range config.Data.Rows {
		x := margin - 10
		y := margin + float64(i)*cellSize + cellSize/2 + 4
		sb.WriteString(fmt.Sprintf(`  <text x="%.2f" y="%.2f" text-anchor="end" font-size="11" fill="#374151">%s</text>`,
			x, y, row))
		sb.WriteString("\n")
	}

	sb.WriteString("</svg>")

	return sb.String(), nil
}

// interpolateColor creates a viridis-like color gradient
func interpolateColor(t float64) string {
	// Simple viridis approximation
	t = math.Max(0, math.Min(1, t))

	if t < 0.25 {
		// Purple to blue
		r := int(68 + (30-68)*(t/0.25))
		g := int(1 + (136-1)*(t/0.25))
		b := int(84 + (229-84)*(t/0.25))
		return fmt.Sprintf("#%02x%02x%02x", r, g, b)
	} else if t < 0.5 {
		// Blue to cyan
		t2 := (t - 0.25) / 0.25
		r := int(30 + (53-30)*t2)
		g := int(136 + (183-136)*t2)
		b := int(229 + (207-229)*t2)
		return fmt.Sprintf("#%02x%02x%02x", r, g, b)
	} else if t < 0.75 {
		// Cyan to yellow
		t2 := (t - 0.5) / 0.25
		r := int(53 + (253-53)*t2)
		g := int(183 + (231-183)*t2)
		b := int(207 + (37-207)*t2)
		return fmt.Sprintf("#%02x%02x%02x", r, g, b)
	} else {
		// Yellow to white
		t2 := (t - 0.75) / 0.25
		r := int(253 + (255-253)*t2)
		g := int(231 + (255-231)*t2)
		b := int(37 + (255-37)*t2)
		return fmt.Sprintf("#%02x%02x%02x", r, g, b)
	}
}

// CreateTreemap generates a treemap chart SVG
func CreateTreemap(config types.TreemapConfig) (string, error) {
	// Convert MCP TreeNode to charts.TreeNode
	root := convertTreeNode(&config.Data)

	spec := maincharts.TreemapSpec{
		Root:         root,
		Width:        float64(config.Width),
		Height:       float64(config.Height),
		Padding:      2,
		ShowLabels:   config.ShowLabels,
		MinLabelSize: 30,
	}

	return maincharts.RenderTreemap(spec), nil
}

// CreateSunburst generates a sunburst chart SVG
func CreateSunburst(config types.SunburstConfig) (string, error) {
	root := convertTreeNode(&config.Data)

	spec := maincharts.SunburstSpec{
		Root:        root,
		Width:       float64(config.Width),
		Height:      float64(config.Height),
		InnerRadius: 0.3,
		ShowLabels:  config.ShowLabels,
	}

	return maincharts.RenderSunburst(spec), nil
}

// CreateCirclePacking generates a circle packing chart SVG
func CreateCirclePacking(config types.CirclePackingConfig) (string, error) {
	root := convertTreeNode(&config.Data)

	spec := maincharts.CirclePackingSpec{
		Root:       root,
		Width:      float64(config.Width),
		Height:     float64(config.Height),
		ShowLabels: config.ShowLabels,
	}

	return maincharts.RenderCirclePacking(spec), nil
}

// CreateIcicle generates an icicle partition chart SVG
func CreateIcicle(config types.IcicleConfig) (string, error) {
	root := convertTreeNode(&config.Data)

	orientation := config.Orientation
	if orientation == "" {
		orientation = "vertical"
	}

	spec := maincharts.IcicleSpec{
		Root:        root,
		Width:       float64(config.Width),
		Height:      float64(config.Height),
		Orientation: orientation,
		ShowLabels:  config.ShowLabels,
	}

	return maincharts.RenderIcicle(spec), nil
}

// CreateBoxPlot generates a box plot SVG
func CreateBoxPlot(config types.BoxPlotConfig) (string, error) {
	data := make([]*maincharts.BoxPlotData, len(config.Data))
	for i, ds := range config.Data {
		data[i] = &maincharts.BoxPlotData{
			Values: ds.Values,
			Label:  ds.Label,
			Color:  "#3B82F6",
		}
	}

	spec := maincharts.BoxPlotSpec{
		Data:              data,
		Width:             float64(config.Width),
		Height:            float64(config.Height),
		Horizontal:        false,
		ShowOutliers:      config.ShowOutliers,
		ShowMean:          config.ShowMean,
		WhiskerMultiplier: 1.5,
	}

	return maincharts.RenderVerticalBoxPlot(spec), nil
}

// CreateViolinPlot generates a violin plot SVG
func CreateViolinPlot(config types.ViolinPlotConfig) (string, error) {
	data := make([]*maincharts.ViolinPlotData, len(config.Data))
	for i, ds := range config.Data {
		data[i] = &maincharts.ViolinPlotData{
			Values: ds.Values,
			Label:  ds.Label,
			Color:  "#3B82F6",
		}
	}

	spec := maincharts.ViolinPlotSpec{
		Data:       data,
		Width:      float64(config.Width),
		Height:     float64(config.Height),
		Bandwidth:  0, // Auto-calculate
		ShowBox:    config.ShowBox,
		ShowMedian: config.ShowMedian,
		ShowMean:   false,
	}

	return maincharts.RenderViolinPlot(spec), nil
}

// CreateHistogram generates a histogram SVG
func CreateHistogram(config types.HistogramConfig) (string, error) {
	bins := config.Bins
	if bins == 0 {
		bins = 20
	}

	spec := maincharts.HistogramSpec{
		Data: &maincharts.HistogramData{
			Values: config.Values,
			Color:  "#3B82F6",
		},
		Width:    float64(config.Width),
		Height:   float64(config.Height),
		BinCount: bins,
		Nice:     true,
	}

	return maincharts.RenderHistogram(spec), nil
}

// CreateRidgeline generates a ridgeline plot SVG
func CreateRidgeline(config types.RidgelineConfig) (string, error) {
	data := make([]*maincharts.RidgelineData, len(config.Data))
	for i, ds := range config.Data {
		data[i] = &maincharts.RidgelineData{
			Label:  ds.Label,
			Values: ds.Values,
		}
	}

	overlap := config.Overlap
	if overlap == 0 {
		overlap = 0.5
	}

	spec := maincharts.RidgelineSpec{
		Data:       data,
		Width:      float64(config.Width),
		Height:     float64(config.Height),
		Overlap:    overlap,
		ShowFill:   true,
		ShowLabels: config.ShowLabels,
	}

	return maincharts.RenderRidgeline(spec), nil
}

// CreateCandlestick generates a candlestick chart SVG
func CreateCandlestick(config types.CandlestickConfig) (string, error) {
	if len(config.Data) == 0 {
		return "", fmt.Errorf("no candlestick data provided")
	}

	// Convert to candlestick data
	candleData := make([]maincharts.CandlestickData, len(config.Data))
	for i, d := range config.Data {
		candleData[i] = maincharts.CandlestickData{
			X:      d.Date,
			Open:   d.Open,
			High:   d.High,
			Low:    d.Low,
			Close:  d.Close,
			Volume: d.Volume,
		}
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

	// Create scales using scales package
	xScale := scales.NewLinearScale(
		[2]float64{0, float64(len(candleData))},
		[2]units.Length{units.Px(50), units.Px(float64(config.Width) - 50)},
	)
	yScale := scales.NewLinearScale(
		[2]float64{minPrice * 0.98, maxPrice * 1.02},
		[2]units.Length{units.Px(float64(config.Height) - 100), units.Px(50)},
	)

	spec := maincharts.CandlestickSpec{
		Data:         candleData,
		Width:        float64(config.Width),
		Height:       float64(config.Height),
		XScale:       xScale,
		YScale:       yScale,
		ShowVolume:   config.ShowVolume,
		VolumeHeight: 100,
	}

	return maincharts.RenderCandlestick(spec), nil
}

// CreateOHLC generates an OHLC bar chart SVG
func CreateOHLC(config types.OHLCConfig) (string, error) {
	if len(config.Data) == 0 {
		return "", fmt.Errorf("no OHLC data provided")
	}

	// Convert to OHLC data
	ohlcData := make([]maincharts.OHLCData, len(config.Data))
	for i, d := range config.Data {
		ohlcData[i] = maincharts.OHLCData{
			X:     d.Date,
			Open:  d.Open,
			High:  d.High,
			Low:   d.Low,
			Close: d.Close,
		}
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

	// Create scales using scales package
	xScale := scales.NewLinearScale(
		[2]float64{0, float64(len(ohlcData))},
		[2]units.Length{units.Px(50), units.Px(float64(config.Width) - 50)},
	)
	yScale := scales.NewLinearScale(
		[2]float64{minPrice * 0.98, maxPrice * 1.02},
		[2]units.Length{units.Px(float64(config.Height) - 50), units.Px(50)},
	)

	spec := maincharts.OHLCSpec{
		Data:   ohlcData,
		Width:  float64(config.Width),
		Height: float64(config.Height),
		XScale: xScale,
		YScale: yScale,
	}

	return maincharts.RenderOHLC(spec), nil
}

// convertTreeNode converts MCP TreeNode to charts.TreeNode recursively
func convertTreeNode(node *types.TreeNode) *maincharts.TreeNode {
	if node == nil {
		return nil
	}

	result := &maincharts.TreeNode{
		Name:  node.Name,
		Value: node.Value,
	}

	if len(node.Children) > 0 {
		result.Children = make([]*maincharts.TreeNode, len(node.Children))
		for i, child := range node.Children {
			result.Children[i] = convertTreeNode(child)
		}
	}

	return result
}

// Keep unused imports to avoid compiler errors
var _ = units.Pixel
var _ *color.Color
var _ = &layout.Node{}
