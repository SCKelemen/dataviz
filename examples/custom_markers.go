package main

import (
	"fmt"
	"time"

	"github.com/SCKelemen/dataviz"
	design "github.com/SCKelemen/design-system"
)

// Demonstrates all available marker types
func main() {
	fmt.Println("=== Custom Markers Example ===\n")

	// Sample data
	points := []dataviz.TimeSeriesData{
		{Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), Value: 100},
		{Date: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), Value: 125},
		{Date: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC), Value: 115},
		{Date: time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC), Value: 140},
	}

	// Try all marker types
	markerTypes := []string{"circle", "square", "diamond", "triangle", "cross", "x", "dot"}

	for _, markerType := range markerTypes {
		data := dataviz.LineGraphData{
			Points:     points,
			Color:      "#3B82F6",
			MarkerType: markerType,
			MarkerSize: 5,
		}

		bounds := dataviz.Bounds{X: 0, Y: 0, Width: 400, Height: 200}
		config := dataviz.RenderConfig{
			DesignTokens: design.DefaultTheme(),
			Color:        "#3B82F6",
		}

		renderer := dataviz.NewSVGRenderer()
		output := renderer.RenderLineGraph(data, bounds, config)

		fmt.Printf("Marker: %-10s  SVG: %d chars\n", markerType, len(output.String()))
	}

	fmt.Println("\nMarker Usage Guide:")
	fmt.Println("  circle:   Default, good for general data points")
	fmt.Println("  square:   For discrete/categorical data")
	fmt.Println("  diamond:  For highlighting important points")
	fmt.Println("  triangle: For showing trends or direction")
	fmt.Println("  cross:    For marking intersections")
	fmt.Println("  x:        For exclusions or errors")
	fmt.Println("  dot:      For dense datasets")

	// Terminal rendering example
	fmt.Println("\n=== Terminal Rendering ===\n")
	data := dataviz.LineGraphData{
		Points:     points,
		Color:      "#3B82F6",
		MarkerType: "circle",
		MarkerSize: 1,
	}

	termRenderer := dataviz.NewTerminalRenderer()
	bounds := dataviz.Bounds{X: 0, Y: 0, Width: 60, Height: 12}
	config := dataviz.RenderConfig{
		DesignTokens: design.DefaultTheme(),
	}

	output := termRenderer.RenderLineGraph(data, bounds, config)
	fmt.Println(output.String())
}
