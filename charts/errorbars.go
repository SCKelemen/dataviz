package charts

import (
	"fmt"

	"github.com/SCKelemen/dataviz/scales"
	"github.com/SCKelemen/svg"
)

// ErrorBar represents error bars for a data point
type ErrorBar struct {
	X          float64 // X position (data coordinates)
	Y          float64 // Y value (data coordinates)
	ErrorLower float64 // Lower error (absolute value or relative)
	ErrorUpper float64 // Upper error (absolute value or relative)
	IsRelative bool    // If true, errors are relative to Y; if false, absolute values
}

// ErrorBarSpec configures error bar rendering
type ErrorBarSpec struct {
	Bars       []ErrorBar
	Color      string
	CapWidth   float64 // Width of error bar caps (pixels)
	CapStyle   CapStyle
	LineWidth  float64
}

// CapStyle defines the style of error bar caps
type CapStyle string

const (
	CapStyleLine   CapStyle = "line"   // Horizontal line caps
	CapStyleCircle CapStyle = "circle" // Small circles at ends
	CapStyleNone   CapStyle = "none"   // No caps
)

// ConfidenceBand represents a confidence interval band
type ConfidenceBand struct {
	XValues      []float64 // X positions (data coordinates)
	YCenters     []float64 // Center Y values
	YLowerBounds []float64 // Lower bound Y values
	YUpperBounds []float64 // Upper bound Y values
	Color        string    // Fill color
	Opacity      float64   // Band opacity (0-1)
	Label        string    // Optional label
}

// ConfidenceBandSpec configures confidence band rendering
type ConfidenceBandSpec struct {
	Bands []*ConfidenceBand
}

// RenderErrorBars renders error bars on a plot
func RenderErrorBars(spec ErrorBarSpec, xScale, yScale scales.Scale) string {
	if len(spec.Bars) == 0 {
		return ""
	}

	color := spec.Color
	if color == "" {
		color = "#666"
	}

	capWidth := spec.CapWidth
	if capWidth == 0 {
		capWidth = 8
	}

	lineWidth := spec.LineWidth
	if lineWidth == 0 {
		lineWidth = 1.5
	}

	lineStyle := svg.Style{
		Stroke:      color,
		StrokeWidth: lineWidth,
	}

	var result string

	for _, bar := range spec.Bars {
		// Convert to pixel coordinates
		x := xScale.Apply(bar.X).Value

		var yLower, yUpper float64
		if bar.IsRelative {
			yLower = yScale.Apply(bar.Y - bar.ErrorLower).Value
			yUpper = yScale.Apply(bar.Y + bar.ErrorUpper).Value
		} else {
			yLower = yScale.Apply(bar.ErrorLower).Value
			yUpper = yScale.Apply(bar.ErrorUpper).Value
		}

		// Draw vertical line
		result += svg.Line(x, yLower, x, yUpper, lineStyle) + "\n"

		// Draw caps
		switch spec.CapStyle {
		case CapStyleLine:
			result += svg.Line(x-capWidth/2, yLower, x+capWidth/2, yLower, lineStyle) + "\n"
			result += svg.Line(x-capWidth/2, yUpper, x+capWidth/2, yUpper, lineStyle) + "\n"

		case CapStyleCircle:
			capStyle := svg.Style{
				Fill:   color,
				Stroke: color,
			}
			result += svg.Circle(x, yLower, 2, capStyle) + "\n"
			result += svg.Circle(x, yUpper, 2, capStyle) + "\n"

		case CapStyleNone:
			// No caps
		}
	}

	return result
}

// RenderConfidenceBands renders confidence interval bands
func RenderConfidenceBands(spec ConfidenceBandSpec, xScale, yScale scales.Scale) string {
	if len(spec.Bands) == 0 {
		return ""
	}

	var result string

	for _, band := range spec.Bands {
		if len(band.XValues) == 0 || len(band.YCenters) == 0 ||
			len(band.YLowerBounds) == 0 || len(band.YUpperBounds) == 0 {
			continue
		}

		// Build path for the filled area
		var pathData string

		// Start at first lower bound point
		x0 := xScale.Apply(band.XValues[0]).Value
		y0Lower := yScale.Apply(band.YLowerBounds[0]).Value
		pathData += "M " + formatFloat(x0) + " " + formatFloat(y0Lower) + " "

		// Draw along lower bound
		for i := 1; i < len(band.XValues); i++ {
			x := xScale.Apply(band.XValues[i]).Value
			y := yScale.Apply(band.YLowerBounds[i]).Value
			pathData += "L " + formatFloat(x) + " " + formatFloat(y) + " "
		}

		// Draw along upper bound (in reverse)
		for i := len(band.XValues) - 1; i >= 0; i-- {
			x := xScale.Apply(band.XValues[i]).Value
			y := yScale.Apply(band.YUpperBounds[i]).Value
			pathData += "L " + formatFloat(x) + " " + formatFloat(y) + " "
		}

		// Close path
		pathData += "Z"

		// Render filled path
		color := band.Color
		if color == "" {
			color = "#4285f4"
		}

		opacity := band.Opacity
		if opacity == 0 {
			opacity = 0.2
		}

		pathStyle := svg.Style{
			Fill:    color,
			Opacity: opacity,
			Stroke:  "none",
		}

		result += svg.Path(pathData, pathStyle) + "\n"

		// Draw center line if present
		if len(band.YCenters) == len(band.XValues) {
			var centerPath string
			x0 := xScale.Apply(band.XValues[0]).Value
			y0 := yScale.Apply(band.YCenters[0]).Value
			centerPath += "M " + formatFloat(x0) + " " + formatFloat(y0) + " "

			for i := 1; i < len(band.XValues); i++ {
				x := xScale.Apply(band.XValues[i]).Value
				y := yScale.Apply(band.YCenters[i]).Value
				centerPath += "L " + formatFloat(x) + " " + formatFloat(y) + " "
			}

			centerStyle := svg.Style{
				Stroke:      color,
				StrokeWidth: 1.5,
				Fill:        "none",
			}

			result += svg.Path(centerPath, centerStyle) + "\n"
		}
	}

	return result
}

// formatFloat formats a float64 to a string with 2 decimal places
func formatFloat(f float64) string {
	return fmt.Sprintf("%.2f", f)
}
