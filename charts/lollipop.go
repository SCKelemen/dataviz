package charts

import (
	"fmt"

	"github.com/SCKelemen/svg"
	"github.com/SCKelemen/units"
)

// LollipopData represents data for a lollipop chart
type LollipopData struct {
	Values []LollipopPoint
	Color  string // Stem and circle color
}

// LollipopPoint represents a single lollipop
type LollipopPoint struct {
	Label  string
	Value  float64
	Color  string // Optional: override default color
	Radius float64 // Optional: custom circle size
}

// LollipopSpec configures lollipop chart rendering
type LollipopSpec struct {
	Data        *LollipopData
	Width       float64
	Height      float64
	Horizontal  bool    // If true, render horizontal lollipops
	StemWidth   float64 // Width of stem line (default: 2)
	CircleSize  float64 // Radius of circle (default: 6)
	ShowLabels  bool    // Show value labels
	ShowGrid    bool    // Show background grid
	Title       string
	XAxisLabel  string
	YAxisLabel  string
	BaselineY   float64 // Y value for baseline (default: 0)
}

// RenderLollipop generates an SVG lollipop chart
func RenderLollipop(spec LollipopSpec) string {
	if spec.Data == nil || len(spec.Data.Values) == 0 {
		return ""
	}

	// Set defaults
	if spec.StemWidth == 0 {
		spec.StemWidth = 2
	}
	if spec.CircleSize == 0 {
		spec.CircleSize = 6
	}

	// Calculate margins
	margin := 60.0
	if spec.Horizontal {
		margin = 80.0
	}
	chartWidth := spec.Width - (2 * margin)
	chartHeight := spec.Height - (2 * margin)

	// Find min/max values
	minVal := spec.Data.Values[0].Value
	maxVal := spec.Data.Values[0].Value
	for _, point := range spec.Data.Values {
		if point.Value < minVal {
			minVal = point.Value
		}
		if point.Value > maxVal {
			maxVal = point.Value
		}
	}

	// Include baseline in range
	if spec.BaselineY < minVal {
		minVal = spec.BaselineY
	}
	if spec.BaselineY > maxVal {
		maxVal = spec.BaselineY
	}

	// Add padding to range
	valueRange := maxVal - minVal
	if valueRange == 0 {
		valueRange = 1
	}
	minVal -= valueRange * 0.05
	maxVal += valueRange * 0.05

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

	if spec.Horizontal {
		// Horizontal lollipops
		result += renderHorizontalLollipops(spec, margin, chartWidth, chartHeight, minVal, maxVal)
	} else {
		// Vertical lollipops
		result += renderVerticalLollipops(spec, margin, chartWidth, chartHeight, minVal, maxVal)
	}

	return result
}

func renderVerticalLollipops(spec LollipopSpec, margin, chartWidth, chartHeight, minVal, maxVal float64) string {
	var result string

	numPoints := len(spec.Data.Values)
	spacing := chartWidth / float64(numPoints)

	// Calculate baseline position
	baselineY := margin + chartHeight - ((spec.BaselineY-minVal)/(maxVal-minVal))*chartHeight

	// Draw grid if enabled
	if spec.ShowGrid {
		gridStyle := svg.Style{
			Stroke:      "#e5e7eb",
			StrokeWidth: 1,
		}
		steps := 5
		for i := 0; i <= steps; i++ {
			y := margin + (chartHeight / float64(steps) * float64(i))
			result += svg.Line(margin, y, margin+chartWidth, y, gridStyle) + "\n"
		}
	}

	// Draw baseline
	baselineStyle := svg.Style{
		Stroke:      "#d1d5db",
		StrokeWidth: 1.5,
	}
	result += svg.Line(margin, baselineY, margin+chartWidth, baselineY, baselineStyle) + "\n"

	// Draw Y axis
	axisStyle := svg.Style{
		Stroke:      "#374151",
		StrokeWidth: 2,
	}
	result += svg.Line(margin, margin, margin, margin+chartHeight, axisStyle) + "\n"
	result += svg.Line(margin, margin+chartHeight, margin+chartWidth, margin+chartHeight, axisStyle) + "\n"

	// Color
	defaultColor := spec.Data.Color
	if defaultColor == "" {
		defaultColor = "#3b82f6"
	}

	// Draw lollipops
	for i, point := range spec.Data.Values {
		x := margin + spacing*(float64(i)+0.5)

		// Calculate point position
		valueY := margin + chartHeight - ((point.Value-minVal)/(maxVal-minVal))*chartHeight

		// Get color
		color := point.Color
		if color == "" {
			color = defaultColor
		}

		// Draw stem
		stemStyle := svg.Style{
			Stroke:      color,
			StrokeWidth: spec.StemWidth,
			Opacity:     0.6,
		}
		result += svg.Line(x, baselineY, x, valueY, stemStyle) + "\n"

		// Draw circle
		circleRadius := spec.CircleSize
		if point.Radius > 0 {
			circleRadius = point.Radius
		}
		circleStyle := svg.Style{
			Fill:   color,
			Stroke: "#ffffff",
			StrokeWidth: 2,
		}
		result += svg.Circle(x, valueY, circleRadius, circleStyle) + "\n"

		// Draw label
		if point.Label != "" {
			labelStyle := svg.Style{
				FontSize:         units.Px(11),
				FontFamily:       "sans-serif",
				TextAnchor:       svg.TextAnchorMiddle,
				DominantBaseline: svg.DominantBaselineHanging,
			}
			result += svg.Text(point.Label, x, margin+chartHeight+10, labelStyle) + "\n"
		}

		// Draw value label
		if spec.ShowLabels {
			valueText := fmt.Sprintf("%.1f", point.Value)
			valueLabelStyle := svg.Style{
				FontSize:         units.Px(10),
				FontFamily:       "sans-serif",
				TextAnchor:       svg.TextAnchorMiddle,
				DominantBaseline: svg.DominantBaselineTextBottom,
			}
			result += svg.Text(valueText, x, valueY-circleRadius-5, valueLabelStyle) + "\n"
		}
	}

	// Y-axis labels
	labelStyle := svg.Style{
		FontSize:         units.Px(10),
		FontFamily:       "sans-serif",
		TextAnchor:       svg.TextAnchorEnd,
		DominantBaseline: svg.DominantBaselineMiddle,
	}
	steps := 5
	for i := 0; i <= steps; i++ {
		value := minVal + (maxVal-minVal)/float64(steps)*float64(i)
		y := margin + chartHeight - (chartHeight/float64(steps))*float64(i)
		result += svg.Text(fmt.Sprintf("%.1f", value), margin-10, y, labelStyle) + "\n"
	}

	// Axis titles
	if spec.YAxisLabel != "" {
		result += fmt.Sprintf(`<text x="15" y="%.2f" text-anchor="middle" font-size="12" font-family="sans-serif" transform="rotate(-90 15 %.2f)">%s</text>`,
			spec.Height/2, spec.Height/2, spec.YAxisLabel) + "\n"
	}

	if spec.XAxisLabel != "" {
		xLabelStyle := svg.Style{
			FontSize:         units.Px(12),
			FontFamily:       "sans-serif",
			TextAnchor:       svg.TextAnchorMiddle,
		}
		result += svg.Text(spec.XAxisLabel, spec.Width/2, spec.Height-10, xLabelStyle) + "\n"
	}

	return result
}

func renderHorizontalLollipops(spec LollipopSpec, margin, chartWidth, chartHeight, minVal, maxVal float64) string {
	var result string

	numPoints := len(spec.Data.Values)
	spacing := chartHeight / float64(numPoints)

	// Calculate baseline position
	baselineX := margin + ((spec.BaselineY-minVal)/(maxVal-minVal))*chartWidth

	// Draw axes
	axisStyle := svg.Style{
		Stroke:      "#374151",
		StrokeWidth: 2,
	}
	result += svg.Line(margin, margin, margin, margin+chartHeight, axisStyle) + "\n"
	result += svg.Line(margin, margin+chartHeight, margin+chartWidth, margin+chartHeight, axisStyle) + "\n"

	// Color
	defaultColor := spec.Data.Color
	if defaultColor == "" {
		defaultColor = "#3b82f6"
	}

	// Draw lollipops
	for i, point := range spec.Data.Values {
		y := margin + spacing*(float64(i)+0.5)

		// Calculate point position
		valueX := margin + ((point.Value-minVal)/(maxVal-minVal))*chartWidth

		// Get color
		color := point.Color
		if color == "" {
			color = defaultColor
		}

		// Draw stem
		stemStyle := svg.Style{
			Stroke:      color,
			StrokeWidth: spec.StemWidth,
			Opacity:     0.6,
		}
		result += svg.Line(baselineX, y, valueX, y, stemStyle) + "\n"

		// Draw circle
		circleRadius := spec.CircleSize
		if point.Radius > 0 {
			circleRadius = point.Radius
		}
		circleStyle := svg.Style{
			Fill:   color,
			Stroke: "#ffffff",
			StrokeWidth: 2,
		}
		result += svg.Circle(valueX, y, circleRadius, circleStyle) + "\n"

		// Draw label
		if point.Label != "" {
			labelStyle := svg.Style{
				FontSize:         units.Px(11),
				FontFamily:       "sans-serif",
				TextAnchor:       svg.TextAnchorEnd,
				DominantBaseline: svg.DominantBaselineMiddle,
			}
			result += svg.Text(point.Label, margin-10, y, labelStyle) + "\n"
		}
	}

	return result
}
