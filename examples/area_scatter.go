package main

import (
	"fmt"
	"time"

	"github.com/SCKelemen/dataviz"
	design "github.com/SCKelemen/design-system"
)

// Demonstrates area charts and scatter plots
func main() {
	fmt.Println("=== Area Charts and Scatter Plots ===\n")

	// Sample time series data
	points := []dataviz.TimeSeriesData{
		{Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), Value: 100},
		{Date: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), Value: 125},
		{Date: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC), Value: 115},
		{Date: time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC), Value: 140},
		{Date: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC), Value: 130},
		{Date: time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC), Value: 155},
	}

	bounds := dataviz.Bounds{X: 0, Y: 0, Width: 400, Height: 200}
	tokens := design.MidnightTheme()

	// Area Chart - smooth filled region
	fmt.Println("1. Area Chart (Smooth):")
	areaData := dataviz.AreaChartData{
		Points:    points,
		Color:     "#10B981",
		FillColor: "#10B981",
		Smooth:    true,
		Tension:   0.3,
		BaselineY: 0,
	}

	renderer := dataviz.NewSVGRenderer()
	config := dataviz.RenderConfig{
		DesignTokens: tokens,
		Color:        "#10B981",
	}

	output := renderer.RenderAreaChart(areaData, bounds, config)
	fmt.Printf("   Generated SVG: %d characters\n", len(output.String()))
	fmt.Println("   Perfect for showing cumulative data or ranges")
	fmt.Println()

	// Scatter Plot with custom point sizes
	fmt.Println("2. Scatter Plot (Custom Sizes):")
	scatterPoints := []dataviz.ScatterPoint{
		{Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), Value: 100, Size: 4, Label: ""},
		{Date: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), Value: 125, Size: 6, Label: ""},
		{Date: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC), Value: 115, Size: 5, Label: ""},
		{Date: time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC), Value: 140, Size: 10, Label: "Peak"},
		{Date: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC), Value: 130, Size: 7, Label: ""},
		{Date: time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC), Value: 155, Size: 8, Label: ""},
	}

	scatterData := dataviz.ScatterPlotData{
		Points:     scatterPoints,
		Color:      "#F59E0B",
		MarkerType: "triangle",
		MarkerSize: 5,
	}

	config.Color = "#F59E0B"
	output = renderer.RenderScatterPlot(scatterData, bounds, config)
	fmt.Printf("   Generated SVG: %d characters\n", len(output.String()))
	fmt.Println("   Ideal for correlation, distribution, or outlier detection")
	fmt.Println()

	// Terminal rendering comparison
	fmt.Println("3. Terminal Rendering:")
	fmt.Println()

	termRenderer := dataviz.NewTerminalRenderer()
	termBounds := dataviz.Bounds{X: 0, Y: 0, Width: 60, Height: 10}

	fmt.Println("Area Chart:")
	termOutput := termRenderer.RenderAreaChart(areaData, termBounds, config)
	fmt.Println(termOutput.String())

	fmt.Println("Scatter Plot:")
	config.Color = "#F59E0B"
	termOutput = termRenderer.RenderScatterPlot(scatterData, termBounds, config)
	fmt.Println(termOutput.String())

	fmt.Println("\nKey Features:")
	fmt.Println("  - Smooth curves for organic shapes")
	fmt.Println("  - Per-point sizing in scatter plots")
	fmt.Println("  - Optional labels for specific points")
	fmt.Println("  - Dual rendering (SVG + Terminal)")
}
