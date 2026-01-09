package main

import (
	"fmt"
	"os"
	"time"

	"github.com/SCKelemen/dataviz"
	design "github.com/SCKelemen/design-system"
)

// Demonstrates dual rendering: same data, two output formats
func main() {
	fmt.Println("=== Dual Rendering: SVG + Terminal ===\n")

	// Create sample data
	points := []dataviz.TimeSeriesData{
		{Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), Value: 100},
		{Date: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), Value: 125},
		{Date: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC), Value: 115},
		{Date: time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC), Value: 140},
		{Date: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC), Value: 130},
		{Date: time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC), Value: 155},
		{Date: time.Date(2024, 1, 7, 0, 0, 0, 0, time.UTC), Value: 145},
	}

	lineData := dataviz.LineGraphData{
		Points:     points,
		Color:      "#3B82F6",
		Smooth:     true,
		Tension:    0.3,
		MarkerType: "diamond",
		MarkerSize: 4,
	}

	tokens := design.DefaultTheme()

	// SVG Rendering
	fmt.Println("1. SVG Output (for web/export):\n")

	svgRenderer := dataviz.NewSVGRenderer()
	svgBounds := dataviz.Bounds{X: 0, Y: 0, Width: 400, Height: 200}
	svgConfig := dataviz.RenderConfig{
		DesignTokens: tokens,
		Color:        "#3B82F6",
		Theme:        "default",
	}

	svgOutput := svgRenderer.RenderLineGraph(lineData, svgBounds, svgConfig)
	svgString := svgOutput.String()

	fmt.Printf("   Length: %d characters\n", len(svgString))
	fmt.Printf("   Contains: smooth bezier curves, diamond markers\n")
	fmt.Printf("   Usage: Embed in HTML, save to file, display in browser\n\n")

	// Save to file
	err := os.WriteFile("/tmp/line_graph.svg", []byte(svgString), 0644)
	if err == nil {
		fmt.Println("   Saved to: /tmp/line_graph.svg")
	}

	// Terminal Rendering
	fmt.Println("\n2. Terminal Output (for CLI tools):\n")

	termRenderer := dataviz.NewTerminalRenderer()
	termBounds := dataviz.Bounds{X: 0, Y: 0, Width: 60, Height: 12}
	termConfig := dataviz.RenderConfig{
		DesignTokens: tokens,
	}

	termOutput := termRenderer.RenderLineGraph(lineData, termBounds, termConfig)
	fmt.Println(termOutput.String())

	fmt.Println("\n   Rendering method: ASCII art with Unicode symbols")
	fmt.Println("   Benefits: No graphics needed, works in any terminal")
	fmt.Println("   Usage: CLI dashboards, server monitoring, SSH sessions")

	// Compare all visualization types
	fmt.Println("\n3. All Visualization Types in Terminal:\n")

	// Heatmap
	heatmapData := dataviz.HeatmapData{
		Days: []dataviz.ContributionDay{
			{Date: time.Now().AddDate(0, 0, -6), Count: 5},
			{Date: time.Now().AddDate(0, 0, -5), Count: 10},
			{Date: time.Now().AddDate(0, 0, -4), Count: 8},
			{Date: time.Now().AddDate(0, 0, -3), Count: 15},
			{Date: time.Now().AddDate(0, 0, -2), Count: 12},
			{Date: time.Now().AddDate(0, 0, -1), Count: 18},
			{Date: time.Now(), Count: 20},
		},
		Type: "linear",
	}

	fmt.Println("Heatmap (contribution intensity):")
	heatmapOutput := termRenderer.RenderHeatmap(heatmapData, dataviz.Bounds{Width: 30, Height: 5}, termConfig)
	fmt.Println(heatmapOutput.String())

	// Area Chart
	fmt.Println("\nArea Chart (filled region):")
	areaData := dataviz.AreaChartData{
		Points:  points[:5],
		Smooth:  true,
		Tension: 0.3,
	}
	areaOutput := termRenderer.RenderAreaChart(areaData, dataviz.Bounds{Width: 40, Height: 8}, termConfig)
	fmt.Println(areaOutput.String())

	// Scatter Plot
	fmt.Println("\nScatter Plot (triangle markers):")
	scatterPoints := []dataviz.ScatterPoint{
		{Date: points[0].Date, Value: points[0].Value},
		{Date: points[1].Date, Value: points[1].Value},
		{Date: points[2].Date, Value: points[2].Value},
		{Date: points[3].Date, Value: points[3].Value},
	}
	scatterData := dataviz.ScatterPlotData{
		Points:     scatterPoints,
		MarkerType: "triangle",
		MarkerSize: 1,
	}
	scatterOutput := termRenderer.RenderScatterPlot(scatterData, dataviz.Bounds{Width: 40, Height: 8}, termConfig)
	fmt.Println(scatterOutput.String())

	fmt.Println("\n=== Key Takeaways ===")
	fmt.Println("✓ Same data structure works for both renderers")
	fmt.Println("✓ Switch between SVG and Terminal with one line")
	fmt.Println("✓ Terminal perfect for CLI tools, logs, monitoring")
	fmt.Println("✓ SVG perfect for web apps, reports, exports")
}
