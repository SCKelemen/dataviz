package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/SCKelemen/dataviz/charts"
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
