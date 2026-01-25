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
		svg, err := generator()
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

	w, h := 800, 350
	totalWidth := w * 3
	totalHeight := h

	var content string

	// White background
	content += svg.Rect(0, 0, float64(totalWidth), float64(totalHeight), svg.Style{Fill: "#ffffff"})
	content += "\n"

	// Title
	titleStyle := svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
	content += svg.Text("Pie Chart Gallery", float64(totalWidth)/2, 30, titleStyle)
	content += "\n"

	labelStyle := svg.Style{
		FontSize:   units.Px(14),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#666",
		TextAnchor: "middle",
	}

	// Regular pie chart
	content += svg.Group(
		svg.Text("Regular Pie", 400, 0, labelStyle)+
			svg.Group(
				charts.RenderPieChart(data, 0, 0, w, h-70, "", false, true, true),
				"translate(0, 20)",
				svg.Style{},
			),
		"translate(0, 50)",
		svg.Style{},
	)
	content += "\n"

	// Donut chart
	content += svg.Group(
		svg.Text("Donut Chart", 400, 0, labelStyle)+
			svg.Group(
				charts.RenderPieChart(data, 0, 0, w, h-70, "", true, true, true),
				"translate(0, 20)",
				svg.Style{},
			),
		fmt.Sprintf("translate(%d, 50)", w),
		svg.Style{},
	)
	content += "\n"

	// Custom colors
	dataColors := data
	dataColors.Colors = []string{"#ef4444", "#f97316", "#eab308", "#22c55e", "#3b82f6"}
	content += svg.Group(
		svg.Text("Custom Colors", 400, 0, labelStyle)+
			svg.Group(
				charts.RenderPieChart(dataColors, 0, 0, w, h-70, "", false, true, true),
				"translate(0, 20)",
				svg.Style{},
			),
		fmt.Sprintf("translate(%d, 50)", w*2),
		svg.Style{},
	)
	content += "\n"

	return wrapSVG(content, totalWidth, totalHeight), nil
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

	w, h := 800, 450
	totalWidth := w*2 + 100
	totalHeight := h

	var content string

	// White background
	content += svg.Rect(0, 0, float64(totalWidth), float64(totalHeight), svg.Style{Fill: "#ffffff"})
	content += "\n"

	// Title
	titleStyle := svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
	content += svg.Text("Bar Chart Gallery", float64(totalWidth)/2, 30, titleStyle)
	content += "\n"

	labelStyle := svg.Style{
		FontSize:   units.Px(14),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#666",
		TextAnchor: "middle",
	}

	// Simple bars
	content += svg.Group(
		svg.Text("Simple Bars", 400, 0, labelStyle)+
			svg.Group(
				charts.RenderBarChart(dataSimple, 0, 0, w, h-100, tokens),
				"translate(0, 30)",
				svg.Style{},
			),
		"translate(50, 60)",
		svg.Style{},
	)
	content += "\n"

	// Stacked bars
	content += svg.Group(
		svg.Text("Stacked Bars (Open/Closed)", 400, 0, labelStyle)+
			svg.Group(
				charts.RenderBarChart(dataStacked, 0, 0, w, h-100, tokens),
				"translate(0, 30)",
				svg.Style{},
			),
		fmt.Sprintf("translate(%d, 60)", w+50),
		svg.Style{},
	)
	content += "\n"

	return wrapSVG(content, totalWidth, totalHeight), nil
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

	w, h := 650, 350
	totalWidth := w * 2
	totalHeight := h * 2

	var content string

	// White background
	content += svg.Rect(0, 0, float64(totalWidth), float64(totalHeight), svg.Style{Fill: "#ffffff"})
	content += "\n"

	// Title
	titleStyle := svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
	content += svg.Text("Line Graph Gallery", float64(totalWidth)/2, 30, titleStyle)
	content += "\n"

	labelStyle := svg.Style{
		FontSize:   units.Px(14),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#666",
		TextAnchor: "middle",
	}

	// Simple line
	content += svg.Group(
		svg.Text("Simple Line", 325, 0, labelStyle)+
			svg.Group(
				charts.RenderLineGraph(data, 0, 0, w-50, h-90, tokens),
				"translate(0, 25)",
				svg.Style{},
			),
		"translate(0, 60)",
		svg.Style{},
	)
	content += "\n"

	// Smoothed
	dataSmooth := data
	dataSmooth.Smooth = true
	dataSmooth.Tension = 0.3
	content += svg.Group(
		svg.Text("Smoothed", 325, 0, labelStyle)+
			svg.Group(
				charts.RenderLineGraph(dataSmooth, 0, 0, w-50, h-90, tokens),
				"translate(0, 25)",
				svg.Style{},
			),
		fmt.Sprintf("translate(%d, 60)", w),
		svg.Style{},
	)
	content += "\n"

	// With markers
	dataMarkers := data
	dataMarkers.MarkerType = "circle"
	dataMarkers.MarkerSize = 5
	content += svg.Group(
		svg.Text("With Markers", 325, 0, labelStyle)+
			svg.Group(
				charts.RenderLineGraph(dataMarkers, 0, 0, w-50, h-90, tokens),
				"translate(0, 25)",
				svg.Style{},
			),
		"translate(0, 360)",
		svg.Style{},
	)
	content += "\n"

	// Filled area (using FillColor)
	dataFilled := data
	dataFilled.FillColor = "#3b82f620" // Semi-transparent fill
	content += svg.Group(
		svg.Text("Filled Area", 325, 0, labelStyle)+
			svg.Group(
				charts.RenderLineGraph(dataFilled, 0, 0, w-50, h-90, tokens),
				"translate(0, 25)",
				svg.Style{},
			),
		fmt.Sprintf("translate(%d, 360)", w),
		svg.Style{},
	)
	content += "\n"

	return wrapSVG(content, totalWidth, totalHeight), nil
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
	w, h := 450, 350
	cols := 3
	rows := 2
	totalWidth := w * cols
	totalHeight := h*rows + 50

	var content string

	// White background
	content += svg.Rect(0, 0, float64(totalWidth), float64(totalHeight), svg.Style{Fill: "#ffffff"})
	content += "\n"

	// Title
	titleStyle := svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
	content += svg.Text("Scatter Plot Gallery", float64(totalWidth)/2, 30, titleStyle)
	content += "\n"

	labelStyle := svg.Style{
		FontSize:   units.Px(14),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#666",
		TextAnchor: "middle",
	}

	for i, markerType := range markerTypes {
		col := i % cols
		row := i / cols
		x := col * w
		y := row*h + 60

		data := charts.ScatterPlotData{
			Points:     points,
			MarkerType: markerType,
			Color:      "#3b82f6",
		}

		content += svg.Group(
			svg.Text(fmt.Sprintf("Marker: %s", markerType), 225, 0, labelStyle)+
				svg.Group(
					charts.RenderScatterPlot(data, 0, 0, w-50, h-60, tokens),
					"translate(0, 25)",
					svg.Style{},
				),
			fmt.Sprintf("translate(%d, %d)", x, y),
			svg.Style{},
		)
		content += "\n"
	}

	return wrapSVG(content, totalWidth, totalHeight), nil
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

	w, h := 450, 350
	cols := 3
	rows := 2
	totalWidth := w * cols
	totalHeight := h*rows + 50

	var content string

	// White background
	content += svg.Rect(0, 0, float64(totalWidth), float64(totalHeight), svg.Style{Fill: "#ffffff"})
	content += "\n"

	// Title
	titleStyle := svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
	content += svg.Text("Connected Scatter Gallery", float64(totalWidth)/2, 30, titleStyle)
	content += "\n"

	labelStyle := svg.Style{
		FontSize:   units.Px(14),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#666",
		TextAnchor: "middle",
	}

	for i, ls := range lineStyles {
		col := i % cols
		row := i / cols
		x := col * w
		y := row*h + 60

		spec := charts.ConnectedScatterSpec{
			Width:  float64(w - 50),
			Height: float64(h - 80),
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
			svg.Text(fmt.Sprintf("Line: %s", ls.name), 225, 0, labelStyle)+
				svg.Group(
					charts.RenderConnectedScatter(spec),
					"translate(25, 25)",
					svg.Style{},
				),
			fmt.Sprintf("translate(%d, %d)", x, y),
			svg.Style{},
		)
		content += "\n"
	}

	return wrapSVG(content, totalWidth, totalHeight), nil
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

	w, h := 800, 400
	totalWidth := w * 2
	totalHeight := h

	var content string

	// White background
	content += svg.Rect(0, 0, float64(totalWidth), float64(totalHeight), svg.Style{Fill: "#ffffff"})
	content += "\n"

	// Title
	titleStyle := svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
	content += svg.Text("Area Chart Gallery", float64(totalWidth)/2, 30, titleStyle)
	content += "\n"

	labelStyle := svg.Style{
		FontSize:   units.Px(14),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#666",
		TextAnchor: "middle",
	}

	// Simple area
	content += svg.Group(
		svg.Text("Simple Area", 400, 0, labelStyle)+
			svg.Group(
				charts.RenderAreaChart(data, 0, 0, w-50, h-70, tokens),
				"translate(0, 30)",
				svg.Style{},
			),
		"translate(0, 60)",
		svg.Style{},
	)
	content += "\n"

	// With gradient
	dataGradient := data
	dataGradient.Color = "#10b981"
	content += svg.Group(
		svg.Text("Different Color", 400, 0, labelStyle)+
			svg.Group(
				charts.RenderAreaChart(dataGradient, 0, 0, w-50, h-70, tokens),
				"translate(0, 30)",
				svg.Style{},
			),
		fmt.Sprintf("translate(%d, 60)", w),
		svg.Style{},
	)
	content += "\n"

	return wrapSVG(content, totalWidth, totalHeight), nil
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

	w, h := 800, 400
	totalWidth := w * 2
	totalHeight := h

	var content string

	// White background
	content += svg.Rect(0, 0, float64(totalWidth), float64(totalHeight), svg.Style{Fill: "#ffffff"})
	content += "\n"

	// Title
	titleStyle := svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
	content += svg.Text("Stacked Area Gallery", float64(totalWidth)/2, 30, titleStyle)
	content += "\n"

	labelStyle := svg.Style{
		FontSize:   units.Px(14),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#666",
		TextAnchor: "middle",
	}

	// Standard stacked
	spec1 := charts.StackedAreaSpec{
		Points: points,
		Series: series,
		Width:  float64(w - 50),
		Height: float64(h - 80),
	}
	content += svg.Group(
		svg.Text("Standard Stacked", 400, 0, labelStyle)+
			svg.Group(
				charts.RenderStackedArea(spec1),
				"translate(0, 30)",
				svg.Style{},
			),
		"translate(25, 60)",
		svg.Style{},
	)
	content += "\n"

	// Smooth curves
	spec2 := charts.StackedAreaSpec{
		Points: points,
		Series: series,
		Width:  float64(w - 50),
		Height: float64(h - 80),
		Smooth: true,
	}
	content += svg.Group(
		svg.Text("Smooth Curves", 400, 0, labelStyle)+
			svg.Group(
				charts.RenderStackedArea(spec2),
				"translate(0, 30)",
				svg.Style{},
			),
		fmt.Sprintf("translate(%d, 60)", w+25),
		svg.Style{},
	)
	content += "\n"

	return wrapSVG(content, totalWidth, totalHeight), nil
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

	w, h := 800, 250
	totalWidth := w
	totalHeight := h * 2

	var content string

	// White background
	content += svg.Rect(0, 0, float64(totalWidth), float64(totalHeight), svg.Style{Fill: "#ffffff"})
	content += "\n"

	// Title
	titleStyle := svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
	content += svg.Text("Heatmap Gallery", float64(totalWidth)/2, 30, titleStyle)
	content += "\n"

	labelStyle := svg.Style{
		FontSize:   units.Px(14),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#666",
		TextAnchor: "middle",
	}

	// Linear heatmap
	content += svg.Group(
		svg.Text("Linear Heatmap", 400, 0, labelStyle)+
			svg.Group(
				charts.RenderLinearHeatmap(data, 0, 0, w-50, h-80, "#3b82f6", tokens),
				"translate(0, 25)",
				svg.Style{},
			),
		"translate(25, 60)",
		svg.Style{},
	)
	content += "\n"

	// Weeks heatmap (GitHub style)
	content += svg.Group(
		svg.Text("Weeks Heatmap (GitHub Style)", 400, 0, labelStyle)+
			svg.Group(
				charts.RenderWeeksHeatmap(data, 0, 0, w-50, h-80, "#10b981", tokens),
				"translate(0, 25)",
				svg.Style{},
			),
		fmt.Sprintf("translate(25, %d)", h+20),
		svg.Style{},
	)
	content += "\n"

	return wrapSVG(content, totalWidth, totalHeight), nil
}

// Stat card variations: with different trends
func generateStatCardGallery() (string, error) {
	tokens := design.DefaultTheme()

	w, h := 300, 200
	cols := 3
	totalWidth := w * cols
	totalHeight := h * 2

	var content string

	// White background
	content += svg.Rect(0, 0, float64(totalWidth), float64(totalHeight), svg.Style{Fill: "#ffffff"})
	content += "\n"

	// Title
	titleStyle := svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
	content += svg.Text("Stat Card Gallery", float64(totalWidth)/2, 30, titleStyle)
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
		col := i % cols
		row := i / cols
		x := col * w
		y := row*h + 60

		content += svg.Group(
			charts.RenderStatCard(card.data, 0, 0, w-20, h-20, tokens),
			fmt.Sprintf("translate(%d, %d)", x+10, y),
			svg.Style{},
		)
		content += "\n"
	}

	return wrapSVG(content, totalWidth, totalHeight), nil
}

// Box plot variations: vertical and horizontal
func generateBoxPlotGallery() (string, error) {
	// Sample data for three groups
	data := []*charts.BoxPlotData{
		{Label: "Group A", Values: []float64{12, 15, 18, 20, 22, 25, 28, 30, 32, 35, 40, 45}},
		{Label: "Group B", Values: []float64{20, 22, 25, 28, 30, 32, 35, 38, 40, 42, 45, 48, 50}},
		{Label: "Group C", Values: []float64{10, 12, 15, 18, 20, 25, 30, 35, 40, 50, 60}},
	}

	w, h := 600, 400
	totalWidth := w * 2
	totalHeight := h

	var content string

	content += svg.Rect(0, 0, float64(totalWidth), float64(totalHeight), svg.Style{Fill: "#ffffff"})
	content += "\n"

	titleStyle := svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
	content += svg.Text("Box Plot Gallery", float64(totalWidth)/2, 30, titleStyle)
	content += "\n"

	labelStyle := svg.Style{
		FontSize:   units.Px(14),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#666",
		TextAnchor: "middle",
	}

	// Vertical box plot
	spec1 := charts.BoxPlotSpec{
		Data:         data,
		Width:        float64(w - 50),
		Height:       float64(h - 80),
		ShowOutliers: true,
		ShowMean:     true,
	}
	content += svg.Group(
		svg.Text("Vertical Box Plot", 300, 0, labelStyle)+
			svg.Group(
				charts.RenderVerticalBoxPlot(spec1),
				"translate(0, 30)",
				svg.Style{},
			),
		"translate(25, 60)",
		svg.Style{},
	)
	content += "\n"

	// Horizontal box plot
	spec2 := charts.BoxPlotSpec{
		Data:         data,
		Width:        float64(w - 50),
		Height:       float64(h - 80),
		Horizontal:   true,
		ShowOutliers: true,
		ShowNotch:    true,
	}
	content += svg.Group(
		svg.Text("Horizontal Box Plot", 300, 0, labelStyle)+
			svg.Group(
				charts.RenderHorizontalBoxPlot(spec2),
				"translate(0, 30)",
				svg.Style{},
			),
		fmt.Sprintf("translate(%d, 60)", w+25),
		svg.Style{},
	)
	content += "\n"

	return wrapSVG(content, totalWidth, totalHeight), nil
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

	w, h := 600, 400
	totalWidth := w * 2
	totalHeight := h

	var content string

	content += svg.Rect(0, 0, float64(totalWidth), float64(totalHeight), svg.Style{Fill: "#ffffff"})
	content += "\n"

	titleStyle := svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
	content += svg.Text("Histogram Gallery", float64(totalWidth)/2, 30, titleStyle)
	content += "\n"

	labelStyle := svg.Style{
		FontSize:   units.Px(14),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#666",
		TextAnchor: "middle",
	}

	// Basic histogram
	spec1 := charts.HistogramSpec{
		Data:     histData,
		Width:    float64(w - 50),
		Height:   float64(h - 80),
		BinCount: 20,
	}
	content += svg.Group(
		svg.Text("Count Histogram", 300, 0, labelStyle)+
			svg.Group(
				charts.RenderHistogram(spec1),
				"translate(0, 30)",
				svg.Style{},
			),
		"translate(25, 60)",
		svg.Style{},
	)
	content += "\n"

	// Density histogram
	spec2 := charts.HistogramSpec{
		Data:        histData,
		Width:       float64(w - 50),
		Height:      float64(h - 80),
		BinCount:    20,
		ShowDensity: true,
	}
	content += svg.Group(
		svg.Text("Density Histogram", 300, 0, labelStyle)+
			svg.Group(
				charts.RenderHistogram(spec2),
				"translate(0, 30)",
				svg.Style{},
			),
		fmt.Sprintf("translate(%d, 60)", w+25),
		svg.Style{},
	)
	content += "\n"

	return wrapSVG(content, totalWidth, totalHeight), nil
}

// Violin plot variations
func generateViolinPlotGallery() (string, error) {
	// Sample data for three groups
	data := []*charts.ViolinPlotData{
		{Label: "Group A", Values: []float64{12, 15, 18, 20, 22, 25, 28, 30, 32, 35, 40}},
		{Label: "Group B", Values: []float64{20, 22, 25, 28, 30, 32, 35, 38, 40, 42, 45}},
		{Label: "Group C", Values: []float64{10, 15, 20, 25, 30, 35, 40, 50, 60}},
	}

	w, h := 600, 400
	totalWidth := w * 2
	totalHeight := h

	var content string

	content += svg.Rect(0, 0, float64(totalWidth), float64(totalHeight), svg.Style{Fill: "#ffffff"})
	content += "\n"

	titleStyle := svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
	content += svg.Text("Violin Plot Gallery", float64(totalWidth)/2, 30, titleStyle)
	content += "\n"

	labelStyle := svg.Style{
		FontSize:   units.Px(14),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#666",
		TextAnchor: "middle",
	}

	// Basic violin plot
	spec1 := charts.ViolinPlotSpec{
		Data:   data,
		Width:  float64(w - 50),
		Height: float64(h - 80),
	}
	content += svg.Group(
		svg.Text("Basic Violin Plot", 300, 0, labelStyle)+
			svg.Group(
				charts.RenderViolinPlot(spec1),
				"translate(0, 30)",
				svg.Style{},
			),
		"translate(25, 60)",
		svg.Style{},
	)
	content += "\n"

	// Violin with box plot inside
	spec2 := charts.ViolinPlotSpec{
		Data:       data,
		Width:      float64(w - 50),
		Height:     float64(h - 80),
		ShowBox:    true,
		ShowMedian: true,
		ShowMean:   true,
	}
	content += svg.Group(
		svg.Text("Violin + Box Plot", 300, 0, labelStyle)+
			svg.Group(
				charts.RenderViolinPlot(spec2),
				"translate(0, 30)",
				svg.Style{},
			),
		fmt.Sprintf("translate(%d, 60)", w+25),
		svg.Style{},
	)
	content += "\n"

	return wrapSVG(content, totalWidth, totalHeight), nil
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

	w, h := 600, 400
	totalWidth := w * 2
	totalHeight := h

	var content string

	content += svg.Rect(0, 0, float64(totalWidth), float64(totalHeight), svg.Style{Fill: "#ffffff"})
	content += "\n"

	titleStyle := svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
	content += svg.Text("Lollipop Chart Gallery", float64(totalWidth)/2, 30, titleStyle)
	content += "\n"

	labelStyle := svg.Style{
		FontSize:   units.Px(14),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#666",
		TextAnchor: "middle",
	}

	// Vertical lollipop
	spec1 := charts.LollipopSpec{
		Data:       lollipopData,
		Width:      float64(w - 50),
		Height:     float64(h - 80),
		ShowLabels: true,
	}
	content += svg.Group(
		svg.Text("Vertical Lollipop", 300, 0, labelStyle)+
			svg.Group(
				charts.RenderLollipop(spec1),
				"translate(0, 30)",
				svg.Style{},
			),
		"translate(25, 60)",
		svg.Style{},
	)
	content += "\n"

	// Horizontal lollipop
	spec2 := charts.LollipopSpec{
		Data:       lollipopData,
		Width:      float64(w - 50),
		Height:     float64(h - 80),
		Horizontal: true,
		ShowLabels: true,
	}
	content += svg.Group(
		svg.Text("Horizontal Lollipop", 300, 0, labelStyle)+
			svg.Group(
				charts.RenderLollipop(spec2),
				"translate(0, 30)",
				svg.Style{},
			),
		fmt.Sprintf("translate(%d, 60)", w+25),
		svg.Style{},
	)
	content += "\n"

	return wrapSVG(content, totalWidth, totalHeight), nil
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

	w, h := 600, 400
	totalWidth := w * 2
	totalHeight := h

	var content string

	content += svg.Rect(0, 0, float64(totalWidth), float64(totalHeight), svg.Style{Fill: "#ffffff"})
	content += "\n"

	titleStyle := svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
	content += svg.Text("Candlestick Gallery", float64(totalWidth)/2, 30, titleStyle)
	content += "\n"

	labelStyle := svg.Style{
		FontSize:   units.Px(14),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#666",
		TextAnchor: "middle",
	}

	// Create scales for both charts
	xScale := scales.NewTimeScale(
		[2]time.Time{mustParseTime("2024-01-01"), mustParseTime("2024-01-06")},
		[2]units.Length{units.Px(0), units.Px(float64(w - 50))},
	)
	yScale := scales.NewLinearScale([2]float64{90, 125}, [2]units.Length{units.Px(float64(h - 80)), units.Px(0)})

	// Candlestick chart
	spec1 := charts.CandlestickSpec{
		Data:         candleData,
		Width:        float64(w - 50),
		Height:       float64(h - 80),
		XScale:       xScale,
		YScale:       yScale,
		RisingColor:  "#10b981",
		FallingColor: "#ef4444",
	}
	content += svg.Group(
		svg.Text("Candlestick Chart", 300, 0, labelStyle)+
			svg.Group(
				charts.RenderCandlestick(spec1),
				"translate(0, 30)",
				svg.Style{},
			),
		"translate(25, 60)",
		svg.Style{},
	)
	content += "\n"

	// OHLC chart - convert data
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
		Width:        float64(w - 50),
		Height:       float64(h - 80),
		XScale:       xScale,
		YScale:       yScale,
		RisingColor:  "#10b981",
		FallingColor: "#ef4444",
	}
	content += svg.Group(
		svg.Text("OHLC Chart", 300, 0, labelStyle)+
			svg.Group(
				charts.RenderOHLC(ohlcSpec),
				"translate(0, 30)",
				svg.Style{},
			),
		fmt.Sprintf("translate(%d, 60)", w+25),
		svg.Style{},
	)
	content += "\n"

	return wrapSVG(content, totalWidth, totalHeight), nil
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

	w, h := 600, 400
	totalWidth := w * 2
	totalHeight := h

	var content string

	content += svg.Rect(0, 0, float64(totalWidth), float64(totalHeight), svg.Style{Fill: "#ffffff"})
	content += "\n"

	titleStyle := svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
	content += svg.Text("Treemap Gallery", float64(totalWidth)/2, 30, titleStyle)
	content += "\n"

	labelStyle := svg.Style{
		FontSize:   units.Px(14),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#666",
		TextAnchor: "middle",
	}

	// Standard treemap
	spec1 := charts.TreemapSpec{
		Root:       tree,
		Width:      float64(w - 50),
		Height:     float64(h - 80),
		ShowLabels: true,
		ColorScheme: []string{"#3b82f6", "#10b981", "#f59e0b", "#ef4444", "#8b5cf6"},
	}
	content += svg.Group(
		svg.Text("Standard Treemap", 300, 0, labelStyle)+
			svg.Group(
				charts.RenderTreemap(spec1),
				"translate(0, 30)",
				svg.Style{},
			),
		"translate(25, 60)",
		svg.Style{},
	)
	content += "\n"

	// Treemap with padding
	spec2 := charts.TreemapSpec{
		Root:       tree,
		Width:      float64(w - 50),
		Height:     float64(h - 80),
		Padding:    2,
		ShowLabels: true,
		ColorScheme: []string{"#6366f1", "#ec4899", "#14b8a6", "#f97316", "#a855f7"},
	}
	content += svg.Group(
		svg.Text("With Padding", 300, 0, labelStyle)+
			svg.Group(
				charts.RenderTreemap(spec2),
				"translate(0, 30)",
				svg.Style{},
			),
		fmt.Sprintf("translate(%d, 60)", w+25),
		svg.Style{},
	)
	content += "\n"

	return wrapSVG(content, totalWidth, totalHeight), nil
}

// Sunburst variations
func generateSunburstGallery() (string, error) {
	tree := createSampleTree()

	w, h := 500, 500
	totalWidth := w * 2
	totalHeight := h

	var content string

	content += svg.Rect(0, 0, float64(totalWidth), float64(totalHeight), svg.Style{Fill: "#ffffff"})
	content += "\n"

	titleStyle := svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
	content += svg.Text("Sunburst Gallery", float64(totalWidth)/2, 30, titleStyle)
	content += "\n"

	labelStyle := svg.Style{
		FontSize:   units.Px(14),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#666",
		TextAnchor: "middle",
	}

	// Full sunburst
	spec1 := charts.SunburstSpec{
		Root:        tree,
		Width:       float64(w - 50),
		Height:      float64(h - 80),
		ShowLabels:  true,
		ColorScheme: []string{"#3b82f6", "#10b981", "#f59e0b", "#ef4444", "#8b5cf6"},
	}
	content += svg.Group(
		svg.Text("Full Sunburst", 225, 0, labelStyle)+
			svg.Group(
				charts.RenderSunburst(spec1),
				"translate(0, 30)",
				svg.Style{},
			),
		"translate(25, 60)",
		svg.Style{},
	)
	content += "\n"

	// Sunburst with inner radius (donut style)
	spec2 := charts.SunburstSpec{
		Root:        tree,
		Width:       float64(w - 50),
		Height:      float64(h - 80),
		InnerRadius: 60,
		ShowLabels:  true,
		ColorScheme: []string{"#6366f1", "#ec4899", "#14b8a6", "#f97316", "#a855f7"},
	}
	content += svg.Group(
		svg.Text("With Inner Radius", 225, 0, labelStyle)+
			svg.Group(
				charts.RenderSunburst(spec2),
				"translate(0, 30)",
				svg.Style{},
			),
		fmt.Sprintf("translate(%d, 60)", w+25),
		svg.Style{},
	)
	content += "\n"

	return wrapSVG(content, totalWidth, totalHeight), nil
}

// Circle packing variations
func generateCirclePackingGallery() (string, error) {
	tree := createSampleTree()

	w, h := 500, 500
	totalWidth := w * 2
	totalHeight := h

	var content string

	content += svg.Rect(0, 0, float64(totalWidth), float64(totalHeight), svg.Style{Fill: "#ffffff"})
	content += "\n"

	titleStyle := svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
	content += svg.Text("Circle Packing Gallery", float64(totalWidth)/2, 30, titleStyle)
	content += "\n"

	labelStyle := svg.Style{
		FontSize:   units.Px(14),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#666",
		TextAnchor: "middle",
	}

	// Standard circle packing
	spec1 := charts.CirclePackingSpec{
		Root:        tree,
		Width:       float64(w - 50),
		Height:      float64(h - 80),
		ShowLabels:  true,
		ColorScheme: []string{"#3b82f6", "#10b981", "#f59e0b", "#ef4444", "#8b5cf6"},
	}
	content += svg.Group(
		svg.Text("Standard Packing", 225, 0, labelStyle)+
			svg.Group(
				charts.RenderCirclePacking(spec1),
				"translate(0, 30)",
				svg.Style{},
			),
		"translate(25, 60)",
		svg.Style{},
	)
	content += "\n"

	// With padding
	spec2 := charts.CirclePackingSpec{
		Root:        tree,
		Width:       float64(w - 50),
		Height:      float64(h - 80),
		Padding:     5,
		ShowLabels:  true,
		ColorScheme: []string{"#6366f1", "#ec4899", "#14b8a6", "#f97316", "#a855f7"},
	}
	content += svg.Group(
		svg.Text("With Padding", 225, 0, labelStyle)+
			svg.Group(
				charts.RenderCirclePacking(spec2),
				"translate(0, 30)",
				svg.Style{},
			),
		fmt.Sprintf("translate(%d, 60)", w+25),
		svg.Style{},
	)
	content += "\n"

	return wrapSVG(content, totalWidth, totalHeight), nil
}

// Icicle chart variations
func generateIcicleGallery() (string, error) {
	tree := createSampleTree()

	w, h := 600, 400
	totalWidth := w * 2
	totalHeight := h

	var content string

	content += svg.Rect(0, 0, float64(totalWidth), float64(totalHeight), svg.Style{Fill: "#ffffff"})
	content += "\n"

	titleStyle := svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
	content += svg.Text("Icicle Chart Gallery", float64(totalWidth)/2, 30, titleStyle)
	content += "\n"

	labelStyle := svg.Style{
		FontSize:   units.Px(14),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#666",
		TextAnchor: "middle",
	}

	// Vertical icicle
	spec1 := charts.IcicleSpec{
		Root:        tree,
		Width:       float64(w - 50),
		Height:      float64(h - 80),
		Orientation: "vertical",
		ShowLabels:  true,
		ColorScheme: []string{"#3b82f6", "#10b981", "#f59e0b", "#ef4444", "#8b5cf6"},
	}
	content += svg.Group(
		svg.Text("Vertical Icicle", 300, 0, labelStyle)+
			svg.Group(
				charts.RenderIcicle(spec1),
				"translate(0, 30)",
				svg.Style{},
			),
		"translate(25, 60)",
		svg.Style{},
	)
	content += "\n"

	// Horizontal icicle
	spec2 := charts.IcicleSpec{
		Root:        tree,
		Width:       float64(w - 50),
		Height:      float64(h - 80),
		Orientation: "horizontal",
		ShowLabels:  true,
		ColorScheme: []string{"#6366f1", "#ec4899", "#14b8a6", "#f97316", "#a855f7"},
	}
	content += svg.Group(
		svg.Text("Horizontal Icicle", 300, 0, labelStyle)+
			svg.Group(
				charts.RenderIcicle(spec2),
				"translate(0, 30)",
				svg.Style{},
			),
		fmt.Sprintf("translate(%d, 60)", w+25),
		svg.Style{},
	)
	content += "\n"

	return wrapSVG(content, totalWidth, totalHeight), nil
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

	w, h := 500, 500
	totalWidth := w * 2
	totalHeight := h

	var content string

	content += svg.Rect(0, 0, float64(totalWidth), float64(totalHeight), svg.Style{Fill: "#ffffff"})
	content += "\n"

	titleStyle := svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
	content += svg.Text("Radar Chart Gallery", float64(totalWidth)/2, 30, titleStyle)
	content += "\n"

	labelStyle := svg.Style{
		FontSize:   units.Px(14),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#666",
		TextAnchor: "middle",
	}

	// With grid
	spec1 := charts.RadarChartSpec{
		Axes:       axes,
		Series:     series,
		Width:      float64(w - 50),
		Height:     float64(h - 80),
		ShowGrid:   true,
		ShowLabels: true,
		GridLevels: 5,
	}
	content += svg.Group(
		svg.Text("With Grid", 225, 0, labelStyle)+
			svg.Group(
				charts.RenderRadarChart(spec1),
				"translate(0, 30)",
				svg.Style{},
			),
		"translate(25, 60)",
		svg.Style{},
	)
	content += "\n"

	// Without grid
	spec2 := charts.RadarChartSpec{
		Axes:       axes,
		Series:     series,
		Width:      float64(w - 50),
		Height:     float64(h - 80),
		ShowGrid:   false,
		ShowLabels: true,
	}
	content += svg.Group(
		svg.Text("Without Grid", 225, 0, labelStyle)+
			svg.Group(
				charts.RenderRadarChart(spec2),
				"translate(0, 30)",
				svg.Style{},
			),
		fmt.Sprintf("translate(%d, 60)", w+25),
		svg.Style{},
	)
	content += "\n"

	return wrapSVG(content, totalWidth, totalHeight), nil
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

	w, h := 600, 400
	totalWidth := w * 2
	totalHeight := h

	var content string

	content += svg.Rect(0, 0, float64(totalWidth), float64(totalHeight), svg.Style{Fill: "#ffffff"})
	content += "\n"

	titleStyle := svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
	content += svg.Text("Streamchart Gallery", float64(totalWidth)/2, 30, titleStyle)
	content += "\n"

	labelStyle := svg.Style{
		FontSize:   units.Px(14),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#666",
		TextAnchor: "middle",
	}

	// Center layout
	spec1 := charts.StreamChartSpec{
		Points: points,
		Series: series,
		Width:  float64(w - 50),
		Height: float64(h - 80),
		Layout: "center",
	}
	content += svg.Group(
		svg.Text("Center Layout", 300, 0, labelStyle)+
			svg.Group(
				charts.RenderStreamChart(spec1),
				"translate(0, 30)",
				svg.Style{},
			),
		"translate(25, 60)",
		svg.Style{},
	)
	content += "\n"

	// Smooth curves
	spec2 := charts.StreamChartSpec{
		Points: points,
		Series: series,
		Width:  float64(w - 50),
		Height: float64(h - 80),
		Layout: "center",
		Smooth: true,
	}
	content += svg.Group(
		svg.Text("Smooth Curves", 300, 0, labelStyle)+
			svg.Group(
				charts.RenderStreamChart(spec2),
				"translate(0, 30)",
				svg.Style{},
			),
		fmt.Sprintf("translate(%d, 60)", w+25),
		svg.Style{},
	)
	content += "\n"

	return wrapSVG(content, totalWidth, totalHeight), nil
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

	w, h := 600, 400
	totalWidth := w * 2
	totalHeight := h

	var content string

	content += svg.Rect(0, 0, float64(totalWidth), float64(totalHeight), svg.Style{Fill: "#ffffff"})
	content += "\n"

	titleStyle := svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
	content += svg.Text("Ridgeline Gallery", float64(totalWidth)/2, 30, titleStyle)
	content += "\n"

	labelStyle := svg.Style{
		FontSize:   units.Px(14),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#666",
		TextAnchor: "middle",
	}

	// Standard ridgeline
	spec1 := charts.RidgelineSpec{
		Data:       data,
		Width:      float64(w - 50),
		Height:     float64(h - 80),
		Overlap:    0.5,
		ShowLabels: true,
	}
	content += svg.Group(
		svg.Text("Standard Ridgeline", 300, 0, labelStyle)+
			svg.Group(
				charts.RenderRidgeline(spec1),
				"translate(0, 30)",
				svg.Style{},
			),
		"translate(25, 60)",
		svg.Style{},
	)
	content += "\n"

	// With fill
	spec2 := charts.RidgelineSpec{
		Data:       data,
		Width:      float64(w - 50),
		Height:     float64(h - 80),
		Overlap:    0.5,
		ShowFill:   true,
		ShowLabels: true,
	}
	content += svg.Group(
		svg.Text("With Fill", 300, 0, labelStyle)+
			svg.Group(
				charts.RenderRidgeline(spec2),
				"translate(0, 30)",
				svg.Style{},
			),
		fmt.Sprintf("translate(%d, 60)", w+25),
		svg.Style{},
	)
	content += "\n"

	return wrapSVG(content, totalWidth, totalHeight), nil
}
