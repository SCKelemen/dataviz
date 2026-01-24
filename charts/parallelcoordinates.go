package charts

import (
	"fmt"
	"math"

	"github.com/SCKelemen/svg"
	"github.com/SCKelemen/units"
)

// ParallelAxis represents an axis in parallel coordinates
type ParallelAxis struct {
	Label string
	Min   float64
	Max   float64
}

// ParallelDataPoint represents a single observation across all axes
type ParallelDataPoint struct {
	Values []float64 // One value per axis
	Label  string    // Optional label for this data point
	Color  string    // Optional custom color
}

// ParallelCoordinatesSpec configures parallel coordinates chart rendering
type ParallelCoordinatesSpec struct {
	Axes          []ParallelAxis      // Axis definitions
	Data          []ParallelDataPoint // Data points
	Width         float64
	Height        float64
	LineOpacity   float64  // Opacity of lines (default: 0.6)
	LineWidth     float64  // Width of lines (default: 1.5)
	DefaultColor  string   // Default line color
	HighlightColor string  // Color for highlighted lines
	ShowAxesLabels bool    // Show axis labels
	ShowTicks     bool     // Show tick marks on axes
	Title         string
}

// RenderParallelCoordinates generates an SVG parallel coordinates chart
func RenderParallelCoordinates(spec ParallelCoordinatesSpec) string {
	if len(spec.Axes) == 0 || len(spec.Data) == 0 {
		return ""
	}

	// Set defaults
	if spec.LineOpacity == 0 {
		spec.LineOpacity = 0.6
	}
	if spec.LineWidth == 0 {
		spec.LineWidth = 1.5
	}
	if spec.DefaultColor == "" {
		spec.DefaultColor = "#3b82f6"
	}

	// Calculate dimensions
	topMargin := 50.0
	bottomMargin := 80.0
	sideMargin := 60.0
	chartWidth := spec.Width - (2 * sideMargin)
	chartHeight := spec.Height - topMargin - bottomMargin

	numAxes := len(spec.Axes)
	axisSpacing := chartWidth / float64(numAxes-1)

	var result string

	// Draw title
	if spec.Title != "" {
		titleStyle := svg.Style{
			FontSize:         units.Px(16),
			FontFamily:       "sans-serif",
			FontWeight:       "bold",
			TextAnchor:       svg.TextAnchorMiddle,
			DominantBaseline: svg.DominantBaselineHanging,
		}
		result += svg.Text(spec.Title, spec.Width/2, 10, titleStyle) + "\n"
	}

	// Calculate axis positions
	axisPositions := make([]float64, numAxes)
	for i := 0; i < numAxes; i++ {
		axisPositions[i] = sideMargin + float64(i)*axisSpacing
	}

	// Draw data lines first (so axes appear on top)
	for _, dataPoint := range spec.Data {
		if len(dataPoint.Values) != numAxes {
			continue // Skip data points with wrong number of values
		}

		// Build path for this data point
		var pathData string
		for i := 0; i < numAxes; i++ {
			axis := spec.Axes[i]
			value := dataPoint.Values[i]

			// Normalize value to axis range
			axisRange := axis.Max - axis.Min
			if axisRange == 0 {
				axisRange = 1
			}
			normalizedValue := (value - axis.Min) / axisRange

			// Clamp to [0, 1]
			if normalizedValue < 0 {
				normalizedValue = 0
			}
			if normalizedValue > 1 {
				normalizedValue = 1
			}

			// Calculate Y position (inverted for SVG)
			x := axisPositions[i]
			y := topMargin + chartHeight*(1-normalizedValue)

			if i == 0 {
				pathData = fmt.Sprintf("M %.2f %.2f", x, y)
			} else {
				pathData += fmt.Sprintf(" L %.2f %.2f", x, y)
			}
		}

		// Get line color
		lineColor := dataPoint.Color
		if lineColor == "" {
			lineColor = spec.DefaultColor
		}

		// Draw line
		lineStyle := svg.Style{
			Stroke:      lineColor,
			StrokeWidth: spec.LineWidth,
			Fill:        "none",
			Opacity:     spec.LineOpacity,
		}
		result += svg.Path(pathData, lineStyle) + "\n"
	}

	// Draw axes
	axisStyle := svg.Style{
		Stroke:      "#374151",
		StrokeWidth: 2,
	}
	for i := 0; i < numAxes; i++ {
		x := axisPositions[i]
		result += svg.Line(x, topMargin, x, topMargin+chartHeight, axisStyle) + "\n"
	}

	// Draw axis labels and ticks
	for i, axis := range spec.Axes {
		x := axisPositions[i]

		// Draw axis label (top)
		if spec.ShowAxesLabels && axis.Label != "" {
			labelStyle := svg.Style{
				FontSize:         units.Px(12),
				FontFamily:       "sans-serif",
				FontWeight:       "bold",
				TextAnchor:       svg.TextAnchorMiddle,
				DominantBaseline: svg.DominantBaselineTextBottom,
			}
			result += svg.Text(axis.Label, x, topMargin-10, labelStyle) + "\n"
		}

		// Draw ticks and tick labels
		if spec.ShowTicks {
			tickCount := 5
			tickStyle := svg.Style{
				Stroke:      "#9ca3af",
				StrokeWidth: 1,
			}
			tickLabelStyle := svg.Style{
				FontSize:         units.Px(9),
				FontFamily:       "sans-serif",
				TextAnchor:       svg.TextAnchorMiddle,
				DominantBaseline: svg.DominantBaselineHanging,
			}

			for j := 0; j <= tickCount; j++ {
				// Calculate tick position
				t := float64(j) / float64(tickCount)
				y := topMargin + chartHeight*(1-t)

				// Draw tick mark
				result += svg.Line(x-5, y, x+5, y, tickStyle) + "\n"

				// Calculate and draw tick label
				value := axis.Min + (axis.Max-axis.Min)*t
				labelText := fmt.Sprintf("%.1f", value)
				result += svg.Text(labelText, x, topMargin+chartHeight+10, tickLabelStyle) + "\n"
			}

			// Draw min/max labels
			minMaxLabelStyle := svg.Style{
				FontSize:         units.Px(9),
				FontFamily:       "sans-serif",
				TextAnchor:       svg.TextAnchorStart,
				DominantBaseline: svg.DominantBaselineMiddle,
				Fill:             "#6b7280",
			}
			result += svg.Text(fmt.Sprintf("%.1f", axis.Max), x+10, topMargin, minMaxLabelStyle) + "\n"
			result += svg.Text(fmt.Sprintf("%.1f", axis.Min), x+10, topMargin+chartHeight, minMaxLabelStyle) + "\n"
		}
	}

	return result
}

// AutoParallelCoordinates creates a parallel coordinates chart with automatic axis ranges
// data is a 2D slice where each row is a data point and columns are variables
func AutoParallelCoordinates(axisLabels []string, data [][]float64, width, height float64) string {
	if len(axisLabels) == 0 || len(data) == 0 {
		return ""
	}

	numAxes := len(axisLabels)

	// Calculate min/max for each axis
	axes := make([]ParallelAxis, numAxes)
	for i, label := range axisLabels {
		min := math.Inf(1)
		max := math.Inf(-1)

		for _, dataPoint := range data {
			if i < len(dataPoint) {
				value := dataPoint[i]
				if value < min {
					min = value
				}
				if value > max {
					max = value
				}
			}
		}

		// Add padding
		valueRange := max - min
		if valueRange == 0 {
			valueRange = 1
		}
		min -= valueRange * 0.05
		max += valueRange * 0.05

		axes[i] = ParallelAxis{
			Label: label,
			Min:   min,
			Max:   max,
		}
	}

	// Convert data to ParallelDataPoint format
	dataPoints := make([]ParallelDataPoint, len(data))
	for i, row := range data {
		dataPoints[i] = ParallelDataPoint{
			Values: row,
		}
	}

	spec := ParallelCoordinatesSpec{
		Axes:           axes,
		Data:           dataPoints,
		Width:          width,
		Height:         height,
		ShowAxesLabels: true,
		ShowTicks:      true,
	}

	return RenderParallelCoordinates(spec)
}

// ColoredParallelCoordinates creates a parallel coordinates chart with color-coded lines
// colors should have the same length as data
func ColoredParallelCoordinates(axisLabels []string, data [][]float64, colors []string, width, height float64) string {
	if len(axisLabels) == 0 || len(data) == 0 {
		return ""
	}

	numAxes := len(axisLabels)

	// Calculate min/max for each axis
	axes := make([]ParallelAxis, numAxes)
	for i, label := range axisLabels {
		min := math.Inf(1)
		max := math.Inf(-1)

		for _, dataPoint := range data {
			if i < len(dataPoint) {
				value := dataPoint[i]
				if value < min {
					min = value
				}
				if value > max {
					max = value
				}
			}
		}

		// Add padding
		valueRange := max - min
		if valueRange == 0 {
			valueRange = 1
		}
		min -= valueRange * 0.05
		max += valueRange * 0.05

		axes[i] = ParallelAxis{
			Label: label,
			Min:   min,
			Max:   max,
		}
	}

	// Convert data to ParallelDataPoint format with colors
	dataPoints := make([]ParallelDataPoint, len(data))
	for i, row := range data {
		color := ""
		if i < len(colors) {
			color = colors[i]
		}
		dataPoints[i] = ParallelDataPoint{
			Values: row,
			Color:  color,
		}
	}

	spec := ParallelCoordinatesSpec{
		Axes:           axes,
		Data:           dataPoints,
		Width:          width,
		Height:         height,
		ShowAxesLabels: true,
		ShowTicks:      true,
	}

	return RenderParallelCoordinates(spec)
}
