package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/SCKelemen/dataviz/charts"
	"github.com/SCKelemen/dataviz/scales"
	design "github.com/SCKelemen/design-system"
	"github.com/SCKelemen/svg"
	"github.com/SCKelemen/units"
)

// GalleryDimensions holds calculated dimensions for gallery layouts
type GalleryDimensions struct {
	TotalWidth   float64
	TotalHeight  float64
	ChartWidth   float64
	ChartHeight  float64
	ColWidth     float64
	RowHeight    float64
	TitleY       float64
	ChartStartY  float64
	BottomMargin float64
}

// CalculateGridDimensions calculates pixel dimensions for a grid-based gallery
// using relative units that resolve to exact pixels at render time
func CalculateGridDimensions(cols, rows int, baseWidth, baseHeight float64) GalleryDimensions {
	// Use percentages for grid sizing to avoid accumulation errors
	colPct := units.Percent(100.0 / float64(cols))
	rowPct := units.Percent(100.0 / float64(rows))

	// Calculate dimensions with proper margins
	titleMargin := units.Percent(5)  // 5% top margin for title
	bottomMargin := units.Percent(3) // 3% bottom margin
	chartPadding := units.Percent(2) // 2% padding within each cell

	totalWidth := baseWidth * float64(cols)
	totalHeight := baseHeight * float64(rows)

	// Add margins to total height
	titleSpace := titleMargin.Of(totalHeight)
	bottomSpace := bottomMargin.Of(totalHeight)
	totalHeight += titleSpace + bottomSpace

	// Calculate chart dimensions (subtract padding)
	colWidth := colPct.Of(totalWidth)
	rowHeight := rowPct.Of(baseHeight * float64(rows))

	padding := chartPadding.Of(colWidth)
	chartWidth := colWidth - (padding * 2)
	chartHeight := rowHeight - (padding * 2)

	return GalleryDimensions{
		TotalWidth:   totalWidth,
		TotalHeight:  totalHeight,
		ChartWidth:   chartWidth,
		ChartHeight:  chartHeight,
		ColWidth:     colWidth,
		RowHeight:    rowHeight,
		TitleY:       titleSpace * 0.7, // Position title 70% down the title space
		ChartStartY:  titleSpace,
		BottomMargin: bottomSpace,
	}
}

// CalculateSingleRowDimensions calculates dimensions for single-row galleries
func CalculateSingleRowDimensions(cols int, baseWidth, baseHeight float64) GalleryDimensions {
	titleHeight := 50.0
	bottomMargin := 30.0
	chartPadding := 25.0

	totalWidth := baseWidth * float64(cols)
	totalHeight := baseHeight + titleHeight + bottomMargin

	colWidth := totalWidth / float64(cols)
	chartWidth := baseWidth - (chartPadding * 2)
	chartHeight := baseHeight - chartPadding

	return GalleryDimensions{
		TotalWidth:   totalWidth,
		TotalHeight:  totalHeight,
		ChartWidth:   chartWidth,
		ChartHeight:  chartHeight,
		ColWidth:     colWidth,
		RowHeight:    baseHeight,
		TitleY:       30,
		ChartStartY:  titleHeight + 10,
		BottomMargin: bottomMargin,
	}
}

func main() {
	if err := generateGalleries(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Gallery generation complete!")
}

// wrapSVG wraps content in an SVG element with proper xmlns and viewBox
func wrapSVG(content string, width, height int) string {
	return fmt.Sprintf(
		`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">%s</svg>`,
		width, height, width, height, content,
	)
}

func generateGalleries() error {
	outputDir := "examples-gallery"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	generators := map[string]func() (string, error){
		"pie":               generatePieGallery,
		"bar":               generateBarGallery,
		"line":              generateLineGallery,
		"scatter":           generateScatterGallery,
		"connected-scatter": generateConnectedScatterGallery,
		"area":              generateAreaGallery,
		"stacked-area":      generateStackedAreaGallery,
		"heatmap":           generateHeatmapGallery,
		"statcard":          generateStatCardGallery,
		"boxplot":           generateBoxPlotGallery,
		"histogram":         generateHistogramGallery,
		"violin":            generateViolinPlotGallery,
		"lollipop":          generateLollipopGallery,
		"candlestick":       generateCandlestickGallery,
		"treemap":           generateTreemapGallery,
		"sunburst":          generateSunburstGallery,
		"circle-packing":    generateCirclePackingGallery,
		"icicle":            generateIcicleGallery,
		"radar":             generateRadarGallery,
		"streamchart":       generateStreamChartGallery,
		"ridgeline":         generateRidgelineGallery,
	}

	for name, generator := range generators {
		fmt.Printf("Generating %s gallery...\n", name)

		var svg string
		var err error

		// Check if gallery exists in the new registry system
		if config, ok := GalleryRegistry[name]; ok {
			svg, err = GenerateGallery(config)
		} else {
			// Fall back to old generator
			svg, err = generator()
		}

		if err != nil {
			fmt.Printf("  ✗ Failed: %v\n", err)
			continue
		}

		outputPath := filepath.Join(outputDir, name+"-gallery.svg")
		if err := os.WriteFile(outputPath, []byte(svg), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", outputPath, err)
		}
		fmt.Printf("  ✓ %s\n", outputPath)
	}

	return nil
}

// Pie chart variations: regular, donut, different color schemes
func generatePieGallery() (string, error) {
	data := charts.PieChartData{
		Slices: []charts.PieSlice{
			{Label: "Chrome", Value: 63.5},
			{Label: "Safari", Value: 19.3},
			{Label: "Firefox", Value: 9.2},
			{Label: "Edge", Value: 5.1},
			{Label: "Other", Value: 2.9},
		},
	}

	// Use relative sizing - calculate pixels at last moment
	dims := CalculateSingleRowDimensions(3, 800, 350)

	var content string

	// White background
	content += svg.Rect(0, 0, dims.TotalWidth, dims.TotalHeight, svg.Style{Fill: "#ffffff"})
	content += "\n"

	// Title
	titleStyle := svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
	content += svg.Text("Pie Chart Gallery", dims.TotalWidth/2, dims.TitleY, titleStyle)
	content += "\n"

	labelStyle := svg.Style{
		FontSize:   units.Px(14),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#666",
		TextAnchor: "middle",
	}

	// Calculate cell positions using percentage-based offsets
	labelOffsetY := 0.0
	chartOffsetY := 20.0

	// Calculate chart dimensions as int for rendering (pixels resolved at last moment)
	chartW := int(dims.ChartWidth)
	chartH := int(dims.ChartHeight - 70)

	// Regular pie chart
	cellX := 0.0
	content += svg.Group(
		svg.Text("Regular Pie", dims.ColWidth/2, labelOffsetY, labelStyle)+
			svg.Group(
				charts.RenderPieChart(data, 0, 0, chartW, chartH, "", false, true, true),
				fmt.Sprintf("translate(0, %.2f)", chartOffsetY),
				svg.Style{},
			),
		fmt.Sprintf("translate(%.2f, %.2f)", cellX, dims.ChartStartY),
		svg.Style{},
	)
	content += "\n"

	// Donut chart
	cellX += dims.ColWidth
	content += svg.Group(
		svg.Text("Donut Chart", dims.ColWidth/2, labelOffsetY, labelStyle)+
			svg.Group(
				charts.RenderPieChart(data, 0, 0, chartW, chartH, "", true, true, true),
				fmt.Sprintf("translate(0, %.2f)", chartOffsetY),
				svg.Style{},
			),
		fmt.Sprintf("translate(%.2f, %.2f)", cellX, dims.ChartStartY),
		svg.Style{},
	)
	content += "\n"

	// Custom colors
	cellX += dims.ColWidth
	dataColors := data
	dataColors.Colors = []string{"#ef4444", "#f97316", "#eab308", "#22c55e", "#3b82f6"}
	content += svg.Group(
		svg.Text("Custom Colors", dims.ColWidth/2, labelOffsetY, labelStyle)+
			svg.Group(
				charts.RenderPieChart(dataColors, 0, 0, chartW, chartH, "", false, true, true),
				fmt.Sprintf("translate(0, %.2f)", chartOffsetY),
				svg.Style{},
			),
		fmt.Sprintf("translate(%.2f, %.2f)", cellX, dims.ChartStartY),
		svg.Style{},
	)
	content += "\n"

	return wrapSVG(content, int(dims.TotalWidth), int(dims.TotalHeight)), nil
}

// Bar chart variations: simple, stacked
func generateBarGallery() (string, error) {
	tokens := design.DefaultTheme()

	dataSimple := charts.BarChartData{
		Label: "Sales",
		Color: "#3b82f6",
		Bars: []charts.BarData{
			{Label: "Q1", Value: 45},
			{Label: "Q2", Value: 60},
			{Label: "Q3", Value: 55},
			{Label: "Q4", Value: 70},
		},
	}

	dataStacked := charts.BarChartData{
		Label:   "Tickets",
		Color:   "#10b981",
		Stacked: true,
		Bars: []charts.BarData{
			{Label: "Mon", Value: 30, Secondary: 20},
			{Label: "Tue", Value: 45, Secondary: 25},
			{Label: "Wed", Value: 40, Secondary: 30},
			{Label: "Thu", Value: 50, Secondary: 15},
		},
	}

	// Use relative sizing with percentage-based calculations
	dims := CalculateSingleRowDimensions(2, 850, 450)

	var content string

	// White background
	content += svg.Rect(0, 0, dims.TotalWidth, dims.TotalHeight, svg.Style{Fill: "#ffffff"})
	content += "\n"

	// Title
	titleStyle := svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
	content += svg.Text("Bar Chart Gallery", dims.TotalWidth/2, dims.TitleY, titleStyle)
	content += "\n"

	labelStyle := svg.Style{
		FontSize:   units.Px(14),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#666",
		TextAnchor: "middle",
	}

	// Calculate chart dimensions
	chartW := int(dims.ChartWidth)
	chartH := int(dims.ChartHeight - 100)

	labelOffsetY := 0.0
	chartOffsetY := 30.0

	// Simple bars
	cellX := 50.0
	content += svg.Group(
		svg.Text("Simple Bars", dims.ColWidth/2, labelOffsetY, labelStyle)+
			svg.Group(
				charts.RenderBarChart(dataSimple, 0, 0, chartW, chartH, tokens),
				fmt.Sprintf("translate(0, %.2f)", chartOffsetY),
				svg.Style{},
			),
		fmt.Sprintf("translate(%.2f, %.2f)", cellX, dims.ChartStartY),
		svg.Style{},
	)
	content += "\n"

	// Stacked bars
	cellX += dims.ColWidth
	content += svg.Group(
		svg.Text("Stacked Bars (Open/Closed)", dims.ColWidth/2, labelOffsetY, labelStyle)+
			svg.Group(
				charts.RenderBarChart(dataStacked, 0, 0, chartW, chartH, tokens),
				fmt.Sprintf("translate(0, %.2f)", chartOffsetY),
				svg.Style{},
			),
		fmt.Sprintf("translate(%.2f, %.2f)", cellX, dims.ChartStartY),
		svg.Style{},
	)
	content += "\n"

	return wrapSVG(content, int(dims.TotalWidth), int(dims.TotalHeight)), nil
}

// Line graph variations: simple, smoothed, markers, filled
func generateLineGallery() (string, error) {
	tokens := design.DefaultTheme()

	data := charts.LineGraphData{
		Label: "Temperature",
		Color: "#3b82f6",
		Points: []charts.TimeSeriesData{
			{Date: mustParseTime("2024-01-01"), Value: 10},
			{Date: mustParseTime("2024-02-01"), Value: 12},
			{Date: mustParseTime("2024-03-01"), Value: 18},
			{Date: mustParseTime("2024-04-01"), Value: 22},
			{Date: mustParseTime("2024-05-01"), Value: 27},
			{Date: mustParseTime("2024-06-01"), Value: 30},
		},
	}

	// Calculate dimensions using relative sizing for 2x2 grid
	dims := CalculateGridDimensions(2, 2, 650, 350)

	// Calculate chart dimensions once (allocate extra space for axis labels)
	chartW := int(dims.ChartWidth)
	chartH := int(dims.ChartHeight)
	labelOffsetY := 0.0
	chartOffsetX := 10.0
	chartOffsetY := 25.0

	var content string

	// White background
	content += svg.Rect(0, 0, dims.TotalWidth, dims.TotalHeight, svg.Style{Fill: "#ffffff"})
	content += "\n"

	// Title
	titleStyle := svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
	content += svg.Text("Line Graph Gallery", dims.TotalWidth/2, dims.TitleY, titleStyle)
	content += "\n"

	labelStyle := svg.Style{
		FontSize:   units.Px(14),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#666",
		TextAnchor: "middle",
	}

	// Simple line (row 0, col 0)
	cellX := 0.0
	cellY := dims.ChartStartY
	content += svg.Group(
		svg.Text("Simple Line", dims.ColWidth/2, labelOffsetY, labelStyle)+
			svg.Group(
				charts.RenderLineGraph(data, 0, 0, chartW, chartH, tokens),
				fmt.Sprintf("translate(%.2f, %.2f)", chartOffsetX, chartOffsetY),
				svg.Style{},
			),
		fmt.Sprintf("translate(%.2f, %.2f)", cellX, cellY),
		svg.Style{},
	)
	content += "\n"

	// Smoothed (row 0, col 1)
	cellX = dims.ColWidth
	dataSmooth := data
	dataSmooth.Smooth = true
	dataSmooth.Tension = 0.3
	content += svg.Group(
		svg.Text("Smoothed", dims.ColWidth/2, labelOffsetY, labelStyle)+
			svg.Group(
				charts.RenderLineGraph(dataSmooth, 0, 0, chartW, chartH, tokens),
				fmt.Sprintf("translate(%.2f, %.2f)", chartOffsetX, chartOffsetY),
				svg.Style{},
			),
		fmt.Sprintf("translate(%.2f, %.2f)", cellX, cellY),
		svg.Style{},
	)
	content += "\n"

	// With markers (row 1, col 0)
	cellX = 0.0
	cellY = dims.ChartStartY + dims.RowHeight
	dataMarkers := data
	dataMarkers.MarkerType = "circle"
	dataMarkers.MarkerSize = 5
	content += svg.Group(
		svg.Text("With Markers", dims.ColWidth/2, labelOffsetY, labelStyle)+
			svg.Group(
				charts.RenderLineGraph(dataMarkers, 0, 0, chartW, chartH, tokens),
				fmt.Sprintf("translate(%.2f, %.2f)", chartOffsetX, chartOffsetY),
				svg.Style{},
			),
		fmt.Sprintf("translate(%.2f, %.2f)", cellX, cellY),
		svg.Style{},
	)
	content += "\n"

	// Filled area (row 1, col 1)
	cellX = dims.ColWidth
	dataFilled := data
	dataFilled.FillColor = "#3b82f620" // Semi-transparent fill
	content += svg.Group(
		svg.Text("Filled Area", dims.ColWidth/2, labelOffsetY, labelStyle)+
			svg.Group(
				charts.RenderLineGraph(dataFilled, 0, 0, chartW, chartH, tokens),
				fmt.Sprintf("translate(%.2f, %.2f)", chartOffsetX, chartOffsetY),
				svg.Style{},
			),
		fmt.Sprintf("translate(%.2f, %.2f)", cellX, cellY),
		svg.Style{},
	)
	content += "\n"

	return wrapSVG(content, int(dims.TotalWidth), int(dims.TotalHeight)), nil
}

// Scatter plot variations: different markers
func generateScatterGallery() (string, error) {
	tokens := design.DefaultTheme()

	points := []charts.ScatterPoint{
		{Label: "A", Date: mustParseTime("2024-01-01"), Value: 55},
		{Label: "B", Date: mustParseTime("2024-02-01"), Value: 78},
		{Label: "C", Date: mustParseTime("2024-03-01"), Value: 44},
		{Label: "D", Date: mustParseTime("2024-04-01"), Value: 66},
		{Label: "E", Date: mustParseTime("2024-05-01"), Value: 33},
		{Label: "F", Date: mustParseTime("2024-06-01"), Value: 77},
		{Label: "G", Date: mustParseTime("2024-07-01"), Value: 22},
		{Label: "H", Date: mustParseTime("2024-08-01"), Value: 88},
	}

	markerTypes := []string{"circle", "square", "diamond", "triangle", "cross", "x"}

	// Calculate dimensions using relative sizing for 2x3 grid (3 cols, 2 rows)
	dims := CalculateGridDimensions(3, 2, 450, 350)

	// Calculate chart dimensions once
	chartW := int(dims.ChartWidth - 50)
	chartH := int(dims.ChartHeight - 60)
	labelOffsetY := 0.0
	chartOffsetY := 25.0

	var content string

	// White background
	content += svg.Rect(0, 0, dims.TotalWidth, dims.TotalHeight, svg.Style{Fill: "#ffffff"})
	content += "\n"

	// Title
	titleStyle := svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
	content += svg.Text("Scatter Plot Gallery", dims.TotalWidth/2, dims.TitleY, titleStyle)
	content += "\n"

	labelStyle := svg.Style{
		FontSize:   units.Px(14),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#666",
		TextAnchor: "middle",
	}

	for i, markerType := range markerTypes {
		col := i % 3
		row := i / 3
		cellX := float64(col) * dims.ColWidth
		cellY := dims.ChartStartY + float64(row)*dims.RowHeight

		data := charts.ScatterPlotData{
			Points:     points,
			MarkerType: markerType,
			Color:      "#3b82f6",
		}

		content += svg.Group(
			svg.Text(fmt.Sprintf("Marker: %s", markerType), dims.ColWidth/2, labelOffsetY, labelStyle)+
				svg.Group(
					charts.RenderScatterPlot(data, 0, 0, chartW, chartH, tokens),
					fmt.Sprintf("translate(0, %.2f)", chartOffsetY),
					svg.Style{},
				),
			fmt.Sprintf("translate(%.2f, %.2f)", cellX, cellY),
			svg.Style{},
		)
		content += "\n"
	}

	return wrapSVG(content, int(dims.TotalWidth), int(dims.TotalHeight)), nil
}

// Connected scatter variations: different line styles
func generateConnectedScatterGallery() (string, error) {
	points := []charts.ConnectedScatterPoint{
		{X: 0, Y: 10},
		{X: 1, Y: 25},
		{X: 2, Y: 15},
		{X: 3, Y: 30},
		{X: 4, Y: 20},
		{X: 5, Y: 35},
	}

	lineStyles := []struct {
		name  string
		style string
	}{
		{"Solid", "solid"},
		{"Dashed", "dashed"},
		{"Dotted", "dotted"},
		{"Dash-Dot", "dashdot"},
		{"Long Dash", "longdash"},
	}

	// Calculate dimensions using relative sizing for 2x3 grid (3 cols, 2 rows)
	dims := CalculateGridDimensions(3, 2, 450, 350)

	// Calculate chart dimensions once
	chartW := dims.ChartWidth - 50
	chartH := dims.ChartHeight - 80
	labelOffsetY := 0.0
	chartOffsetX := 25.0
	chartOffsetY := 25.0

	var content string

	// White background
	content += svg.Rect(0, 0, dims.TotalWidth, dims.TotalHeight, svg.Style{Fill: "#ffffff"})
	content += "\n"

	// Title
	titleStyle := svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
	content += svg.Text("Connected Scatter Gallery", dims.TotalWidth/2, dims.TitleY, titleStyle)
	content += "\n"

	labelStyle := svg.Style{
		FontSize:   units.Px(14),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#666",
		TextAnchor: "middle",
	}

	for i, ls := range lineStyles {
		col := i % 3
		row := i / 3
		cellX := float64(col) * dims.ColWidth
		cellY := dims.ChartStartY + float64(row)*dims.RowHeight

		spec := charts.ConnectedScatterSpec{
			Width:  chartW,
			Height: chartH,
			Series: []*charts.ConnectedScatterSeries{
				{
					Points:    points,
					Color:     "#3b82f6",
					LineStyle: ls.style,
				},
			},
			ShowLines:   true,
			ShowMarkers: true,
		}

		content += svg.Group(
			svg.Text(fmt.Sprintf("Line: %s", ls.name), dims.ColWidth/2, labelOffsetY, labelStyle)+
				svg.Group(
					charts.RenderConnectedScatter(spec),
					fmt.Sprintf("translate(%.2f, %.2f)", chartOffsetX, chartOffsetY),
					svg.Style{},
				),
			fmt.Sprintf("translate(%.2f, %.2f)", cellX, cellY),
			svg.Style{},
		)
		content += "\n"
	}

	return wrapSVG(content, int(dims.TotalWidth), int(dims.TotalHeight)), nil
}

func mustParseTime(s string) time.Time {
	t, _ := time.Parse("2006-01-02", s)
	return t
}

// Area chart variations: simple, with gradient
func generateAreaGallery() (string, error) {
	tokens := design.DefaultTheme()

	data := charts.AreaChartData{
		Label: "Sales",
		Color: "#3b82f6",
		Points: []charts.TimeSeriesData{
			{Date: mustParseTime("2024-01-01"), Value: 100},
			{Date: mustParseTime("2024-02-01"), Value: 120},
			{Date: mustParseTime("2024-03-01"), Value: 110},
			{Date: mustParseTime("2024-04-01"), Value: 140},
			{Date: mustParseTime("2024-05-01"), Value: 130},
			{Date: mustParseTime("2024-06-01"), Value: 150},
		},
	}

	// Use relative sizing
	dims := CalculateSingleRowDimensions(2, 800, 400)

	var content string

	// White background
	content += svg.Rect(0, 0, dims.TotalWidth, dims.TotalHeight, svg.Style{Fill: "#ffffff"})
	content += "\n"

	// Title
	titleStyle := svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
	content += svg.Text("Area Chart Gallery", dims.TotalWidth/2, dims.TitleY, titleStyle)
	content += "\n"

	labelStyle := svg.Style{
		FontSize:   units.Px(14),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#666",
		TextAnchor: "middle",
	}

	chartW := int(dims.ChartWidth)
	chartH := int(dims.ChartHeight)
	labelOffsetY := 0.0
	chartOffsetY := 30.0

	// Simple area
	cellX := 0.0
	content += svg.Group(
		svg.Text("Simple Area", dims.ColWidth/2, labelOffsetY, labelStyle)+
			svg.Group(
				charts.RenderAreaChart(data, 0, 0, chartW, chartH, tokens),
				fmt.Sprintf("translate(10, %.2f)", chartOffsetY),
				svg.Style{},
			),
		fmt.Sprintf("translate(%.2f, %.2f)", cellX, dims.ChartStartY),
		svg.Style{},
	)
	content += "\n"

	// With gradient
	cellX += dims.ColWidth
	dataGradient := data
	dataGradient.Color = "#10b981"
	content += svg.Group(
		svg.Text("Different Color", dims.ColWidth/2, labelOffsetY, labelStyle)+
			svg.Group(
				charts.RenderAreaChart(dataGradient, 0, 0, chartW, chartH, tokens),
				fmt.Sprintf("translate(10, %.2f)", chartOffsetY),
				svg.Style{},
			),
		fmt.Sprintf("translate(%.2f, %.2f)", cellX, dims.ChartStartY),
		svg.Style{},
	)
	content += "\n"

	return wrapSVG(content, int(dims.TotalWidth), int(dims.TotalHeight)), nil
}

// Stacked area chart variations
func generateStackedAreaGallery() (string, error) {
	series := []charts.StackedAreaSeries{
		{Label: "Series A", Color: "#3b82f6"},
		{Label: "Series B", Color: "#10b981"},
		{Label: "Series C", Color: "#f59e0b"},
	}

	points := []charts.StackedAreaPoint{
		{X: 0, Values: []float64{10, 15, 5}},
		{X: 1, Values: []float64{20, 10, 15}},
		{X: 2, Values: []float64{15, 20, 10}},
		{X: 3, Values: []float64{25, 15, 10}},
		{X: 4, Values: []float64{20, 25, 15}},
	}

	// Use relative sizing
	dims := CalculateSingleRowDimensions(2, 800, 400)

	var content string

	// White background
	content += svg.Rect(0, 0, dims.TotalWidth, dims.TotalHeight, svg.Style{Fill: "#ffffff"})
	content += "\n"

	// Title
	titleStyle := svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
	content += svg.Text("Stacked Area Gallery", dims.TotalWidth/2, dims.TitleY, titleStyle)
	content += "\n"

	labelStyle := svg.Style{
		FontSize:   units.Px(14),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#666",
		TextAnchor: "middle",
	}

	chartW := int(dims.ChartWidth)
	chartH := int(dims.ChartHeight)
	labelOffsetY := 0.0
	chartOffsetY := 30.0

	// Standard stacked
	cellX := 25.0
	spec1 := charts.StackedAreaSpec{
		Points: points,
		Series: series,
		Width:  float64(chartW),
		Height: float64(chartH),
	}
	content += svg.Group(
		svg.Text("Standard Stacked", dims.ColWidth/2, labelOffsetY, labelStyle)+
			svg.Group(
				charts.RenderStackedArea(spec1),
				fmt.Sprintf("translate(10, %.2f)", chartOffsetY),
				svg.Style{},
			),
		fmt.Sprintf("translate(%.2f, %.2f)", cellX, dims.ChartStartY),
		svg.Style{},
	)
	content += "\n"

	// Smooth curves
	cellX += dims.ColWidth
	spec2 := charts.StackedAreaSpec{
		Points: points,
		Series: series,
		Width:  float64(chartW),
		Height: float64(chartH),
		Smooth: true,
	}
	content += svg.Group(
		svg.Text("Smooth Curves", dims.ColWidth/2, labelOffsetY, labelStyle)+
			svg.Group(
				charts.RenderStackedArea(spec2),
				fmt.Sprintf("translate(10, %.2f)", chartOffsetY),
				svg.Style{},
			),
		fmt.Sprintf("translate(%.2f, %.2f)", cellX, dims.ChartStartY),
		svg.Style{},
	)
	content += "\n"

	return wrapSVG(content, int(dims.TotalWidth), int(dims.TotalHeight)), nil
}

// Heatmap variations: linear and weeks view
func generateHeatmapGallery() (string, error) {
	tokens := design.DefaultTheme()

	// Generate sample data for a year
	startDate := mustParseTime("2024-01-01")
	days := make([]charts.ContributionDay, 365)
	for i := 0; i < 365; i++ {
		date := startDate.AddDate(0, 0, i)
		// Create some pattern in the data
		count := (i%7)*3 + (i%30)/5
		days[i] = charts.ContributionDay{
			Date:  date,
			Count: count,
		}
	}

	data := charts.HeatmapData{
		Days:      days,
		StartDate: startDate,
		EndDate:   startDate.AddDate(0, 0, 364),
	}

	// Calculate dimensions for vertical stack (1 column, 2 rows)
	baseWidth := 800.0
	baseHeight := 250.0
	titleHeight := 50.0
	rowSpacing := 20.0

	totalWidth := baseWidth
	totalHeight := baseHeight*2 + titleHeight + rowSpacing + 30 // title + 2 rows + spacing + bottom margin

	// Calculate chart dimensions once
	chartW := int(baseWidth - 50)
	chartH := int(baseHeight - 80)
	labelOffsetY := 0.0
	chartOffsetY := 25.0

	var content string

	// White background
	content += svg.Rect(0, 0, totalWidth, totalHeight, svg.Style{Fill: "#ffffff"})
	content += "\n"

	// Title
	titleStyle := svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
	content += svg.Text("Heatmap Gallery", totalWidth/2, 30, titleStyle)
	content += "\n"

	labelStyle := svg.Style{
		FontSize:   units.Px(14),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#666",
		TextAnchor: "middle",
	}

	// Position first row
	rowY := titleHeight + 10

	// Linear heatmap
	content += svg.Group(
		svg.Text("Linear Heatmap", totalWidth/2, labelOffsetY, labelStyle)+
			svg.Group(
				charts.RenderLinearHeatmap(data, 0, 0, chartW, chartH, "#3b82f6", tokens),
				fmt.Sprintf("translate(0, %.2f)", chartOffsetY),
				svg.Style{},
			),
		fmt.Sprintf("translate(25, %.2f)", rowY),
		svg.Style{},
	)
	content += "\n"

	// Position second row
	rowY += baseHeight + rowSpacing

	// Weeks heatmap (GitHub style)
	content += svg.Group(
		svg.Text("Weeks Heatmap (GitHub Style)", totalWidth/2, labelOffsetY, labelStyle)+
			svg.Group(
				charts.RenderWeeksHeatmap(data, 0, 0, chartW, chartH, "#10b981", tokens),
				fmt.Sprintf("translate(0, %.2f)", chartOffsetY),
				svg.Style{},
			),
		fmt.Sprintf("translate(25, %.2f)", rowY),
		svg.Style{},
	)
	content += "\n"

	return wrapSVG(content, int(totalWidth), int(totalHeight)), nil
}

// Stat card variations: with different trends
func generateStatCardGallery() (string, error) {
	tokens := design.DefaultTheme()

	// Calculate dimensions using relative sizing for 2x3 grid (3 cols, 2 rows)
	dims := CalculateGridDimensions(3, 2, 300, 200)

	// Calculate card dimensions once
	cardW := int(dims.ChartWidth - 20)
	cardH := int(dims.ChartHeight - 20)
	cardOffsetX := 10.0

	var content string

	// White background
	content += svg.Rect(0, 0, dims.TotalWidth, dims.TotalHeight, svg.Style{Fill: "#ffffff"})
	content += "\n"

	// Title
	titleStyle := svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
	content += svg.Text("Stat Card Gallery", dims.TotalWidth/2, dims.TitleY, titleStyle)
	content += "\n"

	// Helper to create trend data
	makeTrendData := func(values []int) []charts.TimeSeriesData {
		result := make([]charts.TimeSeriesData, len(values))
		startDate := mustParseTime("2024-01-01")
		for i, v := range values {
			result[i] = charts.TimeSeriesData{
				Date:  startDate.AddDate(0, 0, i*7), // Weekly data
				Value: v,
			}
		}
		return result
	}

	cards := []struct {
		data charts.StatCardData
		name string
	}{
		{
			data: charts.StatCardData{
				Title:     "Total Revenue",
				Value:     "$124.5K",
				Subtitle:  "+12.5% from last month",
				Change:    12,
				ChangePct: 12.5,
				Color:     "#10b981",
				TrendData: makeTrendData([]int{10, 15, 12, 20, 18, 25, 22, 30}),
			},
			name: "Positive Trend",
		},
		{
			data: charts.StatCardData{
				Title:     "Active Users",
				Value:     "8,234",
				Subtitle:  "-3.2% from last month",
				Change:    -3,
				ChangePct: -3.2,
				Color:     "#ef4444",
				TrendData: makeTrendData([]int{30, 28, 25, 27, 23, 20, 22, 18}),
			},
			name: "Negative Trend",
		},
		{
			data: charts.StatCardData{
				Title:     "Conversion Rate",
				Value:     "3.45%",
				Subtitle:  "+0.8% from last month",
				Change:    1,
				ChangePct: 0.8,
				Color:     "#3b82f6",
				TrendData: makeTrendData([]int{15, 18, 16, 20, 22, 21, 24, 25}),
			},
			name: "Steady Growth",
		},
		{
			data: charts.StatCardData{
				Title:     "Page Views",
				Value:     "45.2K",
				Subtitle:  "0.0% from last month",
				Change:    0,
				ChangePct: 0.0,
				Color:     "#6b7280",
				TrendData: makeTrendData([]int{20, 20, 21, 20, 20, 19, 20, 20}),
			},
			name: "Flat Trend",
		},
		{
			data: charts.StatCardData{
				Title:     "Bounce Rate",
				Value:     "42.1%",
				Subtitle:  "+5.3% from last month",
				Change:    5,
				ChangePct: 5.3,
				Color:     "#f59e0b",
				TrendData: makeTrendData([]int{10, 15, 20, 25, 30, 28, 35, 40}),
			},
			name: "Rising",
		},
		{
			data: charts.StatCardData{
				Title:     "Avg Session",
				Value:     "4m 23s",
				Subtitle:  "-1.2% from last month",
				Change:    -1,
				ChangePct: -1.2,
				Color:     "#8b5cf6",
				TrendData: makeTrendData([]int{25, 24, 23, 24, 22, 21, 20, 19}),
			},
			name: "Declining",
		},
	}

	for i, card := range cards {
		col := i % 3
		row := i / 3
		cellX := float64(col)*dims.ColWidth + cardOffsetX
		cellY := dims.ChartStartY + float64(row)*dims.RowHeight

		content += svg.Group(
			charts.RenderStatCard(card.data, 0, 0, cardW, cardH, tokens),
			fmt.Sprintf("translate(%.2f, %.2f)", cellX, cellY),
			svg.Style{},
		)
		content += "\n"
	}

	return wrapSVG(content, int(dims.TotalWidth), int(dims.TotalHeight)), nil
}

// Box plot variations: vertical and horizontal
func generateBoxPlotGallery() (string, error) {
	// Sample data for three groups
	data := []*charts.BoxPlotData{
		{Label: "Group A", Values: []float64{12, 15, 18, 20, 22, 25, 28, 30, 32, 35, 40, 45}},
		{Label: "Group B", Values: []float64{20, 22, 25, 28, 30, 32, 35, 38, 40, 42, 45, 48, 50}},
		{Label: "Group C", Values: []float64{10, 12, 15, 18, 20, 25, 30, 35, 40, 50, 60}},
	}

	// Use relative sizing
	dims := CalculateSingleRowDimensions(2, 600, 400)

	var content string

	content += svg.Rect(0, 0, dims.TotalWidth, dims.TotalHeight, svg.Style{Fill: "#ffffff"})
	content += "\n"

	titleStyle := svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
	content += svg.Text("Box Plot Gallery", dims.TotalWidth/2, dims.TitleY, titleStyle)
	content += "\n"

	labelStyle := svg.Style{
		FontSize:   units.Px(14),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#666",
		TextAnchor: "middle",
	}

	chartW := int(dims.ChartWidth - 50)
	chartH := int(dims.ChartHeight - 80)
	labelOffsetY := 0.0
	chartOffsetY := 30.0

	// Vertical box plot
	cellX := 25.0
	spec1 := charts.BoxPlotSpec{
		Data:         data,
		Width:        float64(chartW),
		Height:       float64(chartH),
		ShowOutliers: true,
		ShowMean:     true,
	}
	content += svg.Group(
		svg.Text("Vertical Box Plot", dims.ColWidth/2, labelOffsetY, labelStyle)+
			svg.Group(
				charts.RenderVerticalBoxPlot(spec1),
				fmt.Sprintf("translate(0, %.2f)", chartOffsetY),
				svg.Style{},
			),
		fmt.Sprintf("translate(%.2f, %.2f)", cellX, dims.ChartStartY),
		svg.Style{},
	)
	content += "\n"

	// Vertical box plot with notches
	cellX += dims.ColWidth
	spec2 := charts.BoxPlotSpec{
		Data:         data,
		Width:        float64(chartW),
		Height:       float64(chartH),
		ShowOutliers: true,
		ShowNotch:    true,
	}
	content += svg.Group(
		svg.Text("With Confidence Notches", dims.ColWidth/2, labelOffsetY, labelStyle)+
			svg.Group(
				charts.RenderVerticalBoxPlot(spec2),
				fmt.Sprintf("translate(0, %.2f)", chartOffsetY),
				svg.Style{},
			),
		fmt.Sprintf("translate(%.2f, %.2f)", cellX, dims.ChartStartY),
		svg.Style{},
	)
	content += "\n"

	return wrapSVG(content, int(dims.TotalWidth), int(dims.TotalHeight)), nil
}

// Histogram variations: basic and with density
func generateHistogramGallery() (string, error) {
	// Generate sample data (normal distribution)
	values := make([]float64, 200)
	for i := range values {
		// Simple approximation of normal distribution
		sum := 0.0
		for j := 0; j < 12; j++ {
			sum += float64(i % 20)
		}
		values[i] = sum/12*5 + 50 + float64((i%10)-5)*2
	}

	histData := &charts.HistogramData{Values: values}

	// Use relative sizing
	dims := CalculateSingleRowDimensions(2, 600, 400)

	var content string

	content += svg.Rect(0, 0, dims.TotalWidth, dims.TotalHeight, svg.Style{Fill: "#ffffff"})
	content += "\n"

	titleStyle := svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
	content += svg.Text("Histogram Gallery", dims.TotalWidth/2, dims.TitleY, titleStyle)
	content += "\n"

	labelStyle := svg.Style{
		FontSize:   units.Px(14),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#666",
		TextAnchor: "middle",
	}

	chartW := int(dims.ChartWidth - 50)
	chartH := int(dims.ChartHeight - 80)
	labelOffsetY := 0.0
	chartOffsetY := 30.0

	// Basic histogram
	cellX := 25.0
	spec1 := charts.HistogramSpec{
		Data:     histData,
		Width:    float64(chartW),
		Height:   float64(chartH),
		BinCount: 20,
	}
	content += svg.Group(
		svg.Text("Count Histogram", dims.ColWidth/2, labelOffsetY, labelStyle)+
			svg.Group(
				charts.RenderHistogram(spec1),
				fmt.Sprintf("translate(0, %.2f)", chartOffsetY),
				svg.Style{},
			),
		fmt.Sprintf("translate(%.2f, %.2f)", cellX, dims.ChartStartY),
		svg.Style{},
	)
	content += "\n"

	// Density histogram
	cellX += dims.ColWidth
	spec2 := charts.HistogramSpec{
		Data:        histData,
		Width:       float64(chartW),
		Height:      float64(chartH),
		BinCount:    20,
		ShowDensity: true,
	}
	content += svg.Group(
		svg.Text("Density Histogram", dims.ColWidth/2, labelOffsetY, labelStyle)+
			svg.Group(
				charts.RenderHistogram(spec2),
				fmt.Sprintf("translate(0, %.2f)", chartOffsetY),
				svg.Style{},
			),
		fmt.Sprintf("translate(%.2f, %.2f)", cellX, dims.ChartStartY),
		svg.Style{},
	)
	content += "\n"

	return wrapSVG(content, int(dims.TotalWidth), int(dims.TotalHeight)), nil
}

// Violin plot variations
func generateViolinPlotGallery() (string, error) {
	// Sample data for three groups
	data := []*charts.ViolinPlotData{
		{Label: "Group A", Values: []float64{12, 15, 18, 20, 22, 25, 28, 30, 32, 35, 40}},
		{Label: "Group B", Values: []float64{20, 22, 25, 28, 30, 32, 35, 38, 40, 42, 45}},
		{Label: "Group C", Values: []float64{10, 15, 20, 25, 30, 35, 40, 50, 60}},
	}

	// Use relative sizing
	dims := CalculateSingleRowDimensions(2, 600, 400)

	var content string

	content += svg.Rect(0, 0, dims.TotalWidth, dims.TotalHeight, svg.Style{Fill: "#ffffff"})
	content += "\n"

	titleStyle := svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
	content += svg.Text("Violin Plot Gallery", dims.TotalWidth/2, dims.TitleY, titleStyle)
	content += "\n"

	labelStyle := svg.Style{
		FontSize:   units.Px(14),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#666",
		TextAnchor: "middle",
	}

	chartW := int(dims.ChartWidth - 50)
	chartH := int(dims.ChartHeight - 80)
	labelOffsetY := 0.0
	chartOffsetY := 30.0

	// Basic violin plot
	cellX := 25.0
	spec1 := charts.ViolinPlotSpec{
		Data:   data,
		Width:  float64(chartW),
		Height: float64(chartH),
	}
	content += svg.Group(
		svg.Text("Basic Violin Plot", dims.ColWidth/2, labelOffsetY, labelStyle)+
			svg.Group(
				charts.RenderViolinPlot(spec1),
				fmt.Sprintf("translate(0, %.2f)", chartOffsetY),
				svg.Style{},
			),
		fmt.Sprintf("translate(%.2f, %.2f)", cellX, dims.ChartStartY),
		svg.Style{},
	)
	content += "\n"

	// Violin with box plot inside
	cellX += dims.ColWidth
	spec2 := charts.ViolinPlotSpec{
		Data:       data,
		Width:      float64(chartW),
		Height:     float64(chartH),
		ShowBox:    true,
		ShowMedian: true,
		ShowMean:   true,
	}
	content += svg.Group(
		svg.Text("Violin + Box Plot", dims.ColWidth/2, labelOffsetY, labelStyle)+
			svg.Group(
				charts.RenderViolinPlot(spec2),
				fmt.Sprintf("translate(0, %.2f)", chartOffsetY),
				svg.Style{},
			),
		fmt.Sprintf("translate(%.2f, %.2f)", cellX, dims.ChartStartY),
		svg.Style{},
	)
	content += "\n"

	return wrapSVG(content, int(dims.TotalWidth), int(dims.TotalHeight)), nil
}

// Lollipop chart variations: vertical and horizontal
func generateLollipopGallery() (string, error) {
	lollipopData := &charts.LollipopData{
		Values: []charts.LollipopPoint{
			{Label: "Product A", Value: 45},
			{Label: "Product B", Value: 62},
			{Label: "Product C", Value: 38},
			{Label: "Product D", Value: 71},
			{Label: "Product E", Value: 54},
		},
		Color: "#3b82f6",
	}

	// Use relative sizing
	dims := CalculateSingleRowDimensions(2, 600, 400)

	var content string

	content += svg.Rect(0, 0, dims.TotalWidth, dims.TotalHeight, svg.Style{Fill: "#ffffff"})
	content += "\n"

	titleStyle := svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
	content += svg.Text("Lollipop Chart Gallery", dims.TotalWidth/2, dims.TitleY, titleStyle)
	content += "\n"

	labelStyle := svg.Style{
		FontSize:   units.Px(14),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#666",
		TextAnchor: "middle",
	}

	chartW := int(dims.ChartWidth - 50)
	chartH := int(dims.ChartHeight - 80)
	labelOffsetY := 0.0
	chartOffsetY := 30.0

	// Vertical lollipop
	cellX := 25.0
	spec1 := charts.LollipopSpec{
		Data:       lollipopData,
		Width:      float64(chartW),
		Height:     float64(chartH),
		ShowLabels: true,
	}
	content += svg.Group(
		svg.Text("Vertical Lollipop", dims.ColWidth/2, labelOffsetY, labelStyle)+
			svg.Group(
				charts.RenderLollipop(spec1),
				fmt.Sprintf("translate(0, %.2f)", chartOffsetY),
				svg.Style{},
			),
		fmt.Sprintf("translate(%.2f, %.2f)", cellX, dims.ChartStartY),
		svg.Style{},
	)
	content += "\n"

	// Horizontal lollipop
	cellX += dims.ColWidth
	spec2 := charts.LollipopSpec{
		Data:       lollipopData,
		Width:      float64(chartW),
		Height:     float64(chartH),
		Horizontal: true,
		ShowLabels: true,
	}
	content += svg.Group(
		svg.Text("Horizontal Lollipop", dims.ColWidth/2, labelOffsetY, labelStyle)+
			svg.Group(
				charts.RenderLollipop(spec2),
				fmt.Sprintf("translate(0, %.2f)", chartOffsetY),
				svg.Style{},
			),
		fmt.Sprintf("translate(%.2f, %.2f)", cellX, dims.ChartStartY),
		svg.Style{},
	)
	content += "\n"

	return wrapSVG(content, int(dims.TotalWidth), int(dims.TotalHeight)), nil
}

// Candlestick chart variations
func generateCandlestickGallery() (string, error) {
	// Sample OHLC data
	candleData := []charts.CandlestickData{
		{X: mustParseTime("2024-01-01"), Open: 100, High: 110, Low: 95, Close: 105, Volume: 1000},
		{X: mustParseTime("2024-01-02"), Open: 105, High: 115, Low: 103, Close: 112, Volume: 1200},
		{X: mustParseTime("2024-01-03"), Open: 112, High: 120, Low: 108, Close: 110, Volume: 1100},
		{X: mustParseTime("2024-01-04"), Open: 110, High: 112, Low: 100, Close: 102, Volume: 1500},
		{X: mustParseTime("2024-01-05"), Open: 102, High: 108, Low: 98, Close: 106, Volume: 1300},
		{X: mustParseTime("2024-01-06"), Open: 106, High: 118, Low: 104, Close: 115, Volume: 1400},
	}

	// Use relative sizing
	dims := CalculateSingleRowDimensions(2, 600, 400)

	var content string

	content += svg.Rect(0, 0, dims.TotalWidth, dims.TotalHeight, svg.Style{Fill: "#ffffff"})
	content += "\n"

	titleStyle := svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
	content += svg.Text("Candlestick Gallery", dims.TotalWidth/2, dims.TitleY, titleStyle)
	content += "\n"

	labelStyle := svg.Style{
		FontSize:   units.Px(14),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#666",
		TextAnchor: "middle",
	}

	chartW := dims.ChartWidth - 50
	chartH := dims.ChartHeight - 80
	labelOffsetY := 0.0
	chartOffsetY := 30.0

	// Create scales for both charts
	xScale := scales.NewTimeScale(
		[2]time.Time{mustParseTime("2024-01-01"), mustParseTime("2024-01-06")},
		[2]units.Length{units.Px(0), units.Px(chartW)},
	)
	yScale := scales.NewLinearScale([2]float64{90, 125}, [2]units.Length{units.Px(chartH), units.Px(0)})

	// Candlestick chart
	cellX := 25.0
	spec1 := charts.CandlestickSpec{
		Data:         candleData,
		Width:        chartW,
		Height:       chartH,
		XScale:       xScale,
		YScale:       yScale,
		RisingColor:  "#10b981",
		FallingColor: "#ef4444",
	}
	content += svg.Group(
		svg.Text("Candlestick Chart", dims.ColWidth/2, labelOffsetY, labelStyle)+
			svg.Group(
				charts.RenderCandlestick(spec1),
				fmt.Sprintf("translate(0, %.2f)", chartOffsetY),
				svg.Style{},
			),
		fmt.Sprintf("translate(%.2f, %.2f)", cellX, dims.ChartStartY),
		svg.Style{},
	)
	content += "\n"

	// OHLC chart - convert data
	cellX += dims.ColWidth
	ohlcData := make([]charts.OHLCData, len(candleData))
	for i, c := range candleData {
		ohlcData[i] = charts.OHLCData{
			X:     c.X,
			Open:  c.Open,
			High:  c.High,
			Low:   c.Low,
			Close: c.Close,
		}
	}

	ohlcSpec := charts.OHLCSpec{
		Data:         ohlcData,
		Width:        chartW,
		Height:       chartH,
		XScale:       xScale,
		YScale:       yScale,
		RisingColor:  "#10b981",
		FallingColor: "#ef4444",
	}
	content += svg.Group(
		svg.Text("OHLC Chart", dims.ColWidth/2, labelOffsetY, labelStyle)+
			svg.Group(
				charts.RenderOHLC(ohlcSpec),
				fmt.Sprintf("translate(0, %.2f)", chartOffsetY),
				svg.Style{},
			),
		fmt.Sprintf("translate(%.2f, %.2f)", cellX, dims.ChartStartY),
		svg.Style{},
	)
	content += "\n"

	return wrapSVG(content, int(dims.TotalWidth), int(dims.TotalHeight)), nil
}

// Helper to create sample tree data
func createSampleTree() *charts.TreeNode {
	return &charts.TreeNode{
		Name:  "Root",
		Value: 100,
		Children: []*charts.TreeNode{
			{
				Name:  "Branch A",
				Value: 40,
				Children: []*charts.TreeNode{
					{Name: "Leaf A1", Value: 15},
					{Name: "Leaf A2", Value: 12},
					{Name: "Leaf A3", Value: 13},
				},
			},
			{
				Name:  "Branch B",
				Value: 35,
				Children: []*charts.TreeNode{
					{Name: "Leaf B1", Value: 20},
					{Name: "Leaf B2", Value: 15},
				},
			},
			{
				Name:  "Branch C",
				Value: 25,
				Children: []*charts.TreeNode{
					{Name: "Leaf C1", Value: 10},
					{Name: "Leaf C2", Value: 8},
					{Name: "Leaf C3", Value: 7},
				},
			},
		},
	}
}

// Treemap variations
func generateTreemapGallery() (string, error) {
	tree := createSampleTree()

	// Use relative sizing
	dims := CalculateSingleRowDimensions(2, 600, 400)

	var content string

	content += svg.Rect(0, 0, dims.TotalWidth, dims.TotalHeight, svg.Style{Fill: "#ffffff"})
	content += "\n"

	titleStyle := svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
	content += svg.Text("Treemap Gallery", dims.TotalWidth/2, dims.TitleY, titleStyle)
	content += "\n"

	labelStyle := svg.Style{
		FontSize:   units.Px(14),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#666",
		TextAnchor: "middle",
	}

	chartW := dims.ChartWidth - 50
	chartH := dims.ChartHeight - 80
	labelOffsetY := 0.0
	chartOffsetY := 30.0

	// Standard treemap
	cellX := 25.0
	spec1 := charts.TreemapSpec{
		Root:       tree,
		Width:      chartW,
		Height:     chartH,
		ShowLabels: true,
		ColorScheme: []string{"#3b82f6", "#10b981", "#f59e0b", "#ef4444", "#8b5cf6"},
	}
	content += svg.Group(
		svg.Text("Standard Treemap", dims.ColWidth/2, labelOffsetY, labelStyle)+
			svg.Group(
				charts.RenderTreemap(spec1),
				fmt.Sprintf("translate(0, %.2f)", chartOffsetY),
				svg.Style{},
			),
		fmt.Sprintf("translate(%.2f, %.2f)", cellX, dims.ChartStartY),
		svg.Style{},
	)
	content += "\n"

	// Treemap with padding
	cellX += dims.ColWidth
	spec2 := charts.TreemapSpec{
		Root:       tree,
		Width:      chartW,
		Height:     chartH,
		Padding:    2,
		ShowLabels: true,
		ColorScheme: []string{"#6366f1", "#ec4899", "#14b8a6", "#f97316", "#a855f7"},
	}
	content += svg.Group(
		svg.Text("With Padding", dims.ColWidth/2, labelOffsetY, labelStyle)+
			svg.Group(
				charts.RenderTreemap(spec2),
				fmt.Sprintf("translate(0, %.2f)", chartOffsetY),
				svg.Style{},
			),
		fmt.Sprintf("translate(%.2f, %.2f)", cellX, dims.ChartStartY),
		svg.Style{},
	)
	content += "\n"

	return wrapSVG(content, int(dims.TotalWidth), int(dims.TotalHeight)), nil
}

// Sunburst variations
func generateSunburstGallery() (string, error) {
	tree := createSampleTree()

	// Calculate dimensions using relative sizing
	// Sunburst charts are square, so we use chartSize for both dimensions
	chartSize := 400.0
	baseWidth := chartSize + 50
	baseHeight := chartSize + 100
	dims := CalculateSingleRowDimensions(2, baseWidth, baseHeight)

	// Chart dimensions
	chartW := chartSize
	chartH := chartSize
	labelOffsetY := 0.0
	chartOffsetY := 30.0

	var content string

	content += svg.Rect(0, 0, dims.TotalWidth, dims.TotalHeight, svg.Style{Fill: "#ffffff"})
	content += "\n"

	titleStyle := svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
	content += svg.Text("Sunburst Gallery", dims.TotalWidth/2, dims.TitleY, titleStyle)
	content += "\n"

	labelStyle := svg.Style{
		FontSize:   units.Px(14),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#666",
		TextAnchor: "middle",
	}

	// Position first cell
	cellX := 25.0

	// Full sunburst
	spec1 := charts.SunburstSpec{
		Root:        tree,
		Width:       chartW,
		Height:      chartH,
		ShowLabels:  true,
		ColorScheme: []string{"#3b82f6", "#10b981", "#f59e0b", "#ef4444", "#8b5cf6"},
	}
	content += svg.Group(
		svg.Text("Full Sunburst", dims.ColWidth/2, labelOffsetY, labelStyle)+
			svg.Group(
				charts.RenderSunburst(spec1),
				fmt.Sprintf("translate(0, %.2f)", chartOffsetY),
				svg.Style{},
			),
		fmt.Sprintf("translate(%.2f, %.2f)", cellX, dims.ChartStartY),
		svg.Style{},
	)
	content += "\n"

	// Position second cell
	cellX += dims.ColWidth

	// Sunburst with inner radius (donut style)
	spec2 := charts.SunburstSpec{
		Root:        tree,
		Width:       chartW,
		Height:      chartH,
		InnerRadius: 60,
		ShowLabels:  true,
		ColorScheme: []string{"#6366f1", "#ec4899", "#14b8a6", "#f97316", "#a855f7"},
	}
	content += svg.Group(
		svg.Text("With Inner Radius", dims.ColWidth/2, labelOffsetY, labelStyle)+
			svg.Group(
				charts.RenderSunburst(spec2),
				fmt.Sprintf("translate(0, %.2f)", chartOffsetY),
				svg.Style{},
			),
		fmt.Sprintf("translate(%.2f, %.2f)", cellX, dims.ChartStartY),
		svg.Style{},
	)
	content += "\n"

	return wrapSVG(content, int(dims.TotalWidth), int(dims.TotalHeight)), nil
}

// Circle packing variations
func generateCirclePackingGallery() (string, error) {
	tree := createSampleTree()

	// Calculate dimensions using relative sizing
	// Circle packing needs extra horizontal space as circles extend beyond center
	chartSize := 400.0
	baseWidth := chartSize + 200
	baseHeight := chartSize + 100
	dims := CalculateSingleRowDimensions(2, baseWidth, baseHeight)

	// Chart dimensions
	chartW := chartSize
	chartH := chartSize
	labelOffsetY := 0.0
	chartOffsetY := 30.0
	chartOffsetX := 100.0 // Extra horizontal centering for circle packing

	var content string

	content += svg.Rect(0, 0, dims.TotalWidth, dims.TotalHeight, svg.Style{Fill: "#ffffff"})
	content += "\n"

	titleStyle := svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
	content += svg.Text("Circle Packing Gallery", dims.TotalWidth/2, dims.TitleY, titleStyle)
	content += "\n"

	labelStyle := svg.Style{
		FontSize:   units.Px(14),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#666",
		TextAnchor: "middle",
	}

	// Position first cell
	cellX := 25.0

	// Standard circle packing
	spec1 := charts.CirclePackingSpec{
		Root:        tree,
		Width:       chartW,
		Height:      chartH,
		ShowLabels:  true,
		ColorScheme: []string{"#3b82f6", "#10b981", "#f59e0b", "#ef4444", "#8b5cf6"},
	}
	content += svg.Group(
		svg.Text("Standard Packing", dims.ColWidth/2, labelOffsetY, labelStyle)+
			svg.Group(
				charts.RenderCirclePacking(spec1),
				fmt.Sprintf("translate(%.2f, %.2f)", chartOffsetX, chartOffsetY),
				svg.Style{},
			),
		fmt.Sprintf("translate(%.2f, %.2f)", cellX, dims.ChartStartY),
		svg.Style{},
	)
	content += "\n"

	// Position second cell
	cellX += dims.ColWidth

	// With padding
	spec2 := charts.CirclePackingSpec{
		Root:        tree,
		Width:       chartW,
		Height:      chartH,
		Padding:     5,
		ShowLabels:  true,
		ColorScheme: []string{"#6366f1", "#ec4899", "#14b8a6", "#f97316", "#a855f7"},
	}
	content += svg.Group(
		svg.Text("With Padding", dims.ColWidth/2, labelOffsetY, labelStyle)+
			svg.Group(
				charts.RenderCirclePacking(spec2),
				fmt.Sprintf("translate(%.2f, %.2f)", chartOffsetX, chartOffsetY),
				svg.Style{},
			),
		fmt.Sprintf("translate(%.2f, %.2f)", cellX, dims.ChartStartY),
		svg.Style{},
	)
	content += "\n"

	return wrapSVG(content, int(dims.TotalWidth), int(dims.TotalHeight)), nil
}

// Icicle chart variations
func generateIcicleGallery() (string, error) {
	tree := createSampleTree()

	// Calculate dimensions using relative sizing
	dims := CalculateSingleRowDimensions(2, 600, 400)

	// Calculate chart dimensions once
	chartW := dims.ChartWidth - 50
	chartH := dims.ChartHeight - 80
	labelOffsetY := 0.0
	chartOffsetY := 30.0

	var content string

	content += svg.Rect(0, 0, dims.TotalWidth, dims.TotalHeight, svg.Style{Fill: "#ffffff"})
	content += "\n"

	titleStyle := svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
	content += svg.Text("Icicle Chart Gallery", dims.TotalWidth/2, dims.TitleY, titleStyle)
	content += "\n"

	labelStyle := svg.Style{
		FontSize:   units.Px(14),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#666",
		TextAnchor: "middle",
	}

	// Position first cell
	cellX := 25.0

	// Vertical icicle
	spec1 := charts.IcicleSpec{
		Root:        tree,
		Width:       chartW,
		Height:      chartH,
		Orientation: "vertical",
		ShowLabels:  true,
		ColorScheme: []string{"#3b82f6", "#10b981", "#f59e0b", "#ef4444", "#8b5cf6"},
	}
	content += svg.Group(
		svg.Text("Vertical Icicle", dims.ColWidth/2, labelOffsetY, labelStyle)+
			svg.Group(
				charts.RenderIcicle(spec1),
				fmt.Sprintf("translate(0, %.2f)", chartOffsetY),
				svg.Style{},
			),
		fmt.Sprintf("translate(%.2f, %.2f)", cellX, dims.ChartStartY),
		svg.Style{},
	)
	content += "\n"

	// Position second cell
	cellX += dims.ColWidth

	// Horizontal icicle
	spec2 := charts.IcicleSpec{
		Root:        tree,
		Width:       chartW,
		Height:      chartH,
		Orientation: "horizontal",
		ShowLabels:  true,
		ColorScheme: []string{"#6366f1", "#ec4899", "#14b8a6", "#f97316", "#a855f7"},
	}
	content += svg.Group(
		svg.Text("Horizontal Icicle", dims.ColWidth/2, labelOffsetY, labelStyle)+
			svg.Group(
				charts.RenderIcicle(spec2),
				fmt.Sprintf("translate(0, %.2f)", chartOffsetY),
				svg.Style{},
			),
		fmt.Sprintf("translate(%.2f, %.2f)", cellX, dims.ChartStartY),
		svg.Style{},
	)
	content += "\n"

	return wrapSVG(content, int(dims.TotalWidth), int(dims.TotalHeight)), nil
}

// Radar chart variations
func generateRadarGallery() (string, error) {
	axes := []charts.RadarAxis{
		{Label: "Speed", Max: 100},
		{Label: "Strength", Max: 100},
		{Label: "Intelligence", Max: 100},
		{Label: "Agility", Max: 100},
		{Label: "Defense", Max: 100},
	}

	series := []*charts.RadarSeries{
		{
			Label:   "Character A",
			Values: []float64{80, 70, 60, 90, 50},
			Color:  "#3b82f6",
		},
		{
			Label:   "Character B",
			Values: []float64{60, 85, 75, 70, 80},
			Color:  "#10b981",
		},
	}

	// Calculate dimensions using relative sizing
	dims := CalculateSingleRowDimensions(2, 500, 500)

	// Calculate chart dimensions once
	chartW := dims.ChartWidth - 50
	chartH := dims.ChartHeight - 80
	labelOffsetY := 0.0
	chartOffsetY := 30.0

	var content string

	content += svg.Rect(0, 0, dims.TotalWidth, dims.TotalHeight, svg.Style{Fill: "#ffffff"})
	content += "\n"

	titleStyle := svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
	content += svg.Text("Radar Chart Gallery", dims.TotalWidth/2, dims.TitleY, titleStyle)
	content += "\n"

	labelStyle := svg.Style{
		FontSize:   units.Px(14),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#666",
		TextAnchor: "middle",
	}

	// Position first cell
	cellX := 25.0

	// With grid
	spec1 := charts.RadarChartSpec{
		Axes:       axes,
		Series:     series,
		Width:      chartW,
		Height:     chartH,
		ShowGrid:   true,
		ShowLabels: true,
		GridLevels: 5,
	}
	content += svg.Group(
		svg.Text("With Grid", dims.ColWidth/2, labelOffsetY, labelStyle)+
			svg.Group(
				charts.RenderRadarChart(spec1),
				fmt.Sprintf("translate(0, %.2f)", chartOffsetY),
				svg.Style{},
			),
		fmt.Sprintf("translate(%.2f, %.2f)", cellX, dims.ChartStartY),
		svg.Style{},
	)
	content += "\n"

	// Position second cell
	cellX += dims.ColWidth

	// Without grid
	spec2 := charts.RadarChartSpec{
		Axes:       axes,
		Series:     series,
		Width:      chartW,
		Height:     chartH,
		ShowGrid:   false,
		ShowLabels: true,
	}
	content += svg.Group(
		svg.Text("Without Grid", dims.ColWidth/2, labelOffsetY, labelStyle)+
			svg.Group(
				charts.RenderRadarChart(spec2),
				fmt.Sprintf("translate(0, %.2f)", chartOffsetY),
				svg.Style{},
			),
		fmt.Sprintf("translate(%.2f, %.2f)", cellX, dims.ChartStartY),
		svg.Style{},
	)
	content += "\n"

	return wrapSVG(content, int(dims.TotalWidth), int(dims.TotalHeight)), nil
}

// Streamchart variations
func generateStreamChartGallery() (string, error) {
	series := []charts.StreamSeries{
		{Label: "Series A", Color: "#3b82f6"},
		{Label: "Series B", Color: "#10b981"},
		{Label: "Series C", Color: "#f59e0b"},
	}

	points := []charts.StreamPoint{
		{X: 0, Values: []float64{10, 15, 5}},
		{X: 1, Values: []float64{20, 10, 15}},
		{X: 2, Values: []float64{15, 20, 10}},
		{X: 3, Values: []float64{25, 15, 10}},
		{X: 4, Values: []float64{20, 25, 15}},
		{X: 5, Values: []float64{30, 20, 12}},
	}

	// Calculate dimensions using relative sizing
	dims := CalculateSingleRowDimensions(2, 600, 400)

	// Calculate chart dimensions once
	chartW := dims.ChartWidth - 50
	chartH := dims.ChartHeight - 80
	labelOffsetY := 0.0
	chartOffsetY := 30.0

	var content string

	content += svg.Rect(0, 0, dims.TotalWidth, dims.TotalHeight, svg.Style{Fill: "#ffffff"})
	content += "\n"

	titleStyle := svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
	content += svg.Text("Streamchart Gallery", dims.TotalWidth/2, dims.TitleY, titleStyle)
	content += "\n"

	labelStyle := svg.Style{
		FontSize:   units.Px(14),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#666",
		TextAnchor: "middle",
	}

	// Position first cell
	cellX := 25.0

	// Center layout
	spec1 := charts.StreamChartSpec{
		Points: points,
		Series: series,
		Width:  chartW,
		Height: chartH,
		Layout: "center",
	}
	content += svg.Group(
		svg.Text("Center Layout", dims.ColWidth/2, labelOffsetY, labelStyle)+
			svg.Group(
				charts.RenderStreamChart(spec1),
				fmt.Sprintf("translate(0, %.2f)", chartOffsetY),
				svg.Style{},
			),
		fmt.Sprintf("translate(%.2f, %.2f)", cellX, dims.ChartStartY),
		svg.Style{},
	)
	content += "\n"

	// Position second cell
	cellX += dims.ColWidth

	// Smooth curves
	spec2 := charts.StreamChartSpec{
		Points: points,
		Series: series,
		Width:  chartW,
		Height: chartH,
		Layout: "center",
		Smooth: true,
	}
	content += svg.Group(
		svg.Text("Smooth Curves", dims.ColWidth/2, labelOffsetY, labelStyle)+
			svg.Group(
				charts.RenderStreamChart(spec2),
				fmt.Sprintf("translate(0, %.2f)", chartOffsetY),
				svg.Style{},
			),
		fmt.Sprintf("translate(%.2f, %.2f)", cellX, dims.ChartStartY),
		svg.Style{},
	)
	content += "\n"

	return wrapSVG(content, int(dims.TotalWidth), int(dims.TotalHeight)), nil
}

// Ridgeline chart variations
func generateRidgelineGallery() (string, error) {
	// Create sample data for 4 distributions
	data := []*charts.RidgelineData{
		{Label: "January", Values: []float64{10, 12, 15, 18, 20, 22, 25, 23, 20, 18, 15, 12}},
		{Label: "February", Values: []float64{15, 18, 20, 22, 25, 28, 30, 28, 25, 22, 20, 18}},
		{Label: "March", Values: []float64{20, 22, 25, 28, 30, 32, 35, 33, 30, 28, 25, 22}},
		{Label: "April", Values: []float64{25, 28, 30, 32, 35, 38, 40, 38, 35, 32, 30, 28}},
	}

	// Calculate dimensions using relative sizing
	dims := CalculateSingleRowDimensions(2, 600, 400)

	// Calculate chart dimensions once
	chartW := dims.ChartWidth - 50
	chartH := dims.ChartHeight - 80
	labelOffsetY := 0.0
	chartOffsetY := 30.0

	var content string

	content += svg.Rect(0, 0, dims.TotalWidth, dims.TotalHeight, svg.Style{Fill: "#ffffff"})
	content += "\n"

	titleStyle := svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
	content += svg.Text("Ridgeline Gallery", dims.TotalWidth/2, dims.TitleY, titleStyle)
	content += "\n"

	labelStyle := svg.Style{
		FontSize:   units.Px(14),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#666",
		TextAnchor: "middle",
	}

	// Position first cell
	cellX := 25.0

	// Standard ridgeline
	spec1 := charts.RidgelineSpec{
		Data:       data,
		Width:      chartW,
		Height:     chartH,
		Overlap:    0.5,
		ShowLabels: true,
	}
	content += svg.Group(
		svg.Text("Standard Ridgeline", dims.ColWidth/2, labelOffsetY, labelStyle)+
			svg.Group(
				charts.RenderRidgeline(spec1),
				fmt.Sprintf("translate(0, %.2f)", chartOffsetY),
				svg.Style{},
			),
		fmt.Sprintf("translate(%.2f, %.2f)", cellX, dims.ChartStartY),
		svg.Style{},
	)
	content += "\n"

	// Position second cell
	cellX += dims.ColWidth

	// With fill
	spec2 := charts.RidgelineSpec{
		Data:       data,
		Width:      chartW,
		Height:     chartH,
		Overlap:    0.5,
		ShowFill:   true,
		ShowLabels: true,
	}
	content += svg.Group(
		svg.Text("With Fill", dims.ColWidth/2, labelOffsetY, labelStyle)+
			svg.Group(
				charts.RenderRidgeline(spec2),
				fmt.Sprintf("translate(0, %.2f)", chartOffsetY),
				svg.Style{},
			),
		fmt.Sprintf("translate(%.2f, %.2f)", cellX, dims.ChartStartY),
		svg.Style{},
	)
	content += "\n"

	return wrapSVG(content, int(dims.TotalWidth), int(dims.TotalHeight)), nil
}
