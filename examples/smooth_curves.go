package main

import (
	"fmt"
	"time"

	"github.com/SCKelemen/dataviz"
	design "github.com/SCKelemen/design-system"
)

// Demonstrates smooth curve interpolation with tension control
func main() {
	fmt.Println("=== Smooth Curves Example ===\n")

	// Sample data
	points := []dataviz.TimeSeriesData{
		{Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), Value: 100},
		{Date: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), Value: 125},
		{Date: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC), Value: 115},
		{Date: time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC), Value: 140},
		{Date: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC), Value: 130},
		{Date: time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC), Value: 155},
		{Date: time.Date(2024, 1, 7, 0, 0, 0, 0, time.UTC), Value: 145},
		{Date: time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC), Value: 165},
	}

	// Compare different tension values
	tensions := []float64{0.0, 0.3, 0.6, 1.0}
	for _, tension := range tensions {
		data := dataviz.LineGraphData{
			Points:  points,
			Color:   "#3B82F6",
			Smooth:  true,
			Tension: tension,
		}

		bounds := dataviz.Bounds{X: 0, Y: 0, Width: 400, Height: 200}
		config := dataviz.RenderConfig{
			DesignTokens: design.DefaultTheme(),
			Color:        "#3B82F6",
		}

		renderer := dataviz.NewSVGRenderer()
		output := renderer.RenderLineGraph(data, bounds, config)

		fmt.Printf("Tension: %.1f\n", tension)
		fmt.Printf("SVG Length: %d characters\n", len(output.String()))
		fmt.Println()
	}

	fmt.Println("Smooth curves use bezier interpolation to create organic, flowing lines.")
	fmt.Println("Tension controls curve sharpness:")
	fmt.Println("  0.0 = Tight curves (follows data closely)")
	fmt.Println("  0.3 = Recommended default (natural curves)")
	fmt.Println("  0.6 = Looser curves")
	fmt.Println("  1.0 = Very loose curves")
}
