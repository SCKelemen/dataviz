package charts

import (
	"fmt"
	"math"

	"github.com/SCKelemen/dataviz/scales"
	"github.com/SCKelemen/svg"
	"github.com/SCKelemen/units"
)

// StackedAreaPoint represents a single X position with values for each series
type StackedAreaPoint struct {
	X      float64
	Values []float64 // One value per series
}

// StackedAreaSeries represents metadata for a series
type StackedAreaSeries struct {
	Label string
	Color string
}

// StackedAreaSpec configures stacked area chart rendering
type StackedAreaSpec struct {
	Points       []StackedAreaPoint  // X positions with values for all series
	Series       []StackedAreaSeries // Metadata for each series (colors, labels)
	Width        float64
	Height       float64
	ShowGrid     bool
	Smooth       bool    // Use smooth curves instead of straight lines
	Title        string
	XAxisLabel   string
	YAxisLabel   string
	XAxisMin     *float64 // Optional: force X axis min
	XAxisMax     *float64 // Optional: force X axis max
	YAxisMin     *float64 // Optional: force Y axis min (usually 0 for stacked)
	YAxisMax     *float64 // Optional: force Y axis max
}

// RenderStackedArea generates an SVG stacked area chart
func RenderStackedArea(spec StackedAreaSpec) string {
	if len(spec.Points) == 0 || len(spec.Series) == 0 {
		return ""
	}

	// Calculate margins
	margin := 60.0
	chartWidth := spec.Width - (2 * margin)
	chartHeight := spec.Height - (2 * margin)

	// Find X range
	xMin := spec.Points[0].X
	xMax := spec.Points[0].X
	for _, point := range spec.Points {
		if point.X < xMin {
			xMin = point.X
		}
		if point.X > xMax {
			xMax = point.X
		}
	}

	// Calculate stacked values and find Y max
	// stackedValues[pointIdx][seriesIdx] = cumulative sum up to and including series
	stackedValues := make([][]float64, len(spec.Points))
	yMax := 0.0
	for i, point := range spec.Points {
		stackedValues[i] = make([]float64, len(spec.Series))
		cumSum := 0.0
		for j := 0; j < len(spec.Series); j++ {
			if j < len(point.Values) {
				cumSum += point.Values[j]
			}
			stackedValues[i][j] = cumSum
		}
		if cumSum > yMax {
			yMax = cumSum
		}
	}

	yMin := 0.0 // Stacked areas typically start at 0

	// Apply forced axis ranges if specified
	if spec.XAxisMin != nil {
		xMin = *spec.XAxisMin
	}
	if spec.XAxisMax != nil {
		xMax = *spec.XAxisMax
	}
	if spec.YAxisMin != nil {
		yMin = *spec.YAxisMin
	}
	if spec.YAxisMax != nil {
		yMax = *spec.YAxisMax
	}

	// Add padding to Y range
	yRange := yMax - yMin
	if yRange == 0 {
		yRange = 1
	}
	yMax += yRange * 0.05

	// Create scales
	xScale := scales.NewLinearScale(
		[2]float64{xMin, xMax},
		[2]units.Length{units.Px(margin), units.Px(spec.Width - margin)},
	)

	yScale := scales.NewLinearScale(
		[2]float64{yMin, yMax},
		[2]units.Length{units.Px(spec.Height - margin), units.Px(margin)},
	)

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

	// Draw grid if enabled
	if spec.ShowGrid {
		gridStyle := svg.Style{
			Stroke:      "#e5e7eb",
			StrokeWidth: 1,
		}
		steps := 5
		// Horizontal grid lines
		for i := 0; i <= steps; i++ {
			y := margin + (chartHeight / float64(steps) * float64(i))
			result += svg.Line(margin, y, margin+chartWidth, y, gridStyle) + "\n"
		}
		// Vertical grid lines
		for i := 0; i <= steps; i++ {
			x := margin + (chartWidth / float64(steps) * float64(i))
			result += svg.Line(x, margin, x, margin+chartHeight, gridStyle) + "\n"
		}
	}

	// Draw axes
	axisStyle := svg.Style{
		Stroke:      "#374151",
		StrokeWidth: 2,
	}
	result += svg.Line(margin, margin, margin, margin+chartHeight, axisStyle) + "\n"
	result += svg.Line(margin, margin+chartHeight, margin+chartWidth, margin+chartHeight, axisStyle) + "\n"

	// Default colors
	defaultColors := []string{"#3b82f6", "#10b981", "#f59e0b", "#ef4444", "#8b5cf6", "#ec4899"}

	// Draw stacked areas from bottom to top (reverse order)
	// This ensures proper layering
	for seriesIdx := len(spec.Series) - 1; seriesIdx >= 0; seriesIdx-- {
		series := spec.Series[seriesIdx]

		// Get series color
		seriesColor := series.Color
		if seriesColor == "" {
			seriesColor = defaultColors[seriesIdx%len(defaultColors)]
		}

		// Build path for this layer
		// Top edge: current series cumulative values
		// Bottom edge: previous series cumulative values (or 0 for first series)

		var pathData string

		// Top edge (left to right)
		for i, point := range spec.Points {
			x := xScale.Apply(point.X).Value
			y := yScale.Apply(stackedValues[i][seriesIdx]).Value

			if i == 0 {
				pathData = fmt.Sprintf("M %.2f %.2f", x, y)
			} else {
				pathData += fmt.Sprintf(" L %.2f %.2f", x, y)
			}
		}

		// Bottom edge (right to left)
		for i := len(spec.Points) - 1; i >= 0; i-- {
			x := xScale.Apply(spec.Points[i].X).Value
			var y float64
			if seriesIdx == 0 {
				// First series: bottom is baseline (0)
				y = yScale.Apply(0).Value
			} else {
				// Other series: bottom is previous series cumulative
				y = yScale.Apply(stackedValues[i][seriesIdx-1]).Value
			}
			pathData += fmt.Sprintf(" L %.2f %.2f", x, y)
		}

		pathData += " Z" // Close path

		// Draw filled area
		fillStyle := svg.Style{
			Fill:        seriesColor,
			FillOpacity: 0.7,
			Stroke:      seriesColor,
			StrokeWidth: 1,
		}
		result += svg.Path(pathData, fillStyle) + "\n"
	}

	// X-axis labels
	xLabelStyle := svg.Style{
		FontSize:         units.Px(10),
		FontFamily:       "sans-serif",
		TextAnchor:       svg.TextAnchorMiddle,
		DominantBaseline: svg.DominantBaselineHanging,
	}
	steps := 5
	for i := 0; i <= steps; i++ {
		value := xMin + (xMax-xMin)/float64(steps)*float64(i)
		x := xScale.Apply(value).Value
		result += svg.Text(fmt.Sprintf("%.2f", value), x, spec.Height-margin+10, xLabelStyle) + "\n"
	}

	// Y-axis labels
	yLabelStyle := svg.Style{
		FontSize:         units.Px(10),
		FontFamily:       "sans-serif",
		TextAnchor:       svg.TextAnchorEnd,
		DominantBaseline: svg.DominantBaselineMiddle,
	}
	for i := 0; i <= steps; i++ {
		value := yMin + (yMax-yMin)/float64(steps)*float64(i)
		y := yScale.Apply(value).Value
		result += svg.Text(fmt.Sprintf("%.2f", value), margin-10, y, yLabelStyle) + "\n"
	}

	// Axis titles
	if spec.YAxisLabel != "" {
		result += fmt.Sprintf(`<text x="15" y="%.2f" text-anchor="middle" font-size="12" font-family="sans-serif" transform="rotate(-90 15 %.2f)">%s</text>`,
			spec.Height/2, spec.Height/2, spec.YAxisLabel) + "\n"
	}

	if spec.XAxisLabel != "" {
		xTitleStyle := svg.Style{
			FontSize:   units.Px(12),
			FontFamily: "sans-serif",
			TextAnchor: svg.TextAnchorMiddle,
		}
		result += svg.Text(spec.XAxisLabel, spec.Width/2, spec.Height-10, xTitleStyle) + "\n"
	}

	// Legend
	if len(spec.Series) > 0 {
		legendX := margin + 20
		legendY := margin + 20
		for idx, series := range spec.Series {
			if series.Label == "" {
				continue
			}

			color := series.Color
			if color == "" {
				color = defaultColors[idx%len(defaultColors)]
			}

			yOffset := float64(idx * 25)

			// Legend swatch
			rectStyle := svg.Style{
				Fill:        color,
				FillOpacity: 0.7,
				Stroke:      color,
				StrokeWidth: 1,
			}
			result += svg.Rect(legendX, legendY+yOffset-8, 20, 16, rectStyle) + "\n"

			// Legend text
			textStyle := svg.Style{
				FontSize:         units.Px(12),
				FontFamily:       "sans-serif",
				DominantBaseline: svg.DominantBaselineMiddle,
			}
			result += svg.Text(series.Label, legendX+25, legendY+yOffset, textStyle) + "\n"
		}
	}

	return result
}

// StackedAreaFromSeries is a helper to convert multiple simple series into stacked area format
// Each series is an array of values, and all series must have the same length
func StackedAreaFromSeries(xValues []float64, seriesData [][]float64, seriesLabels []string, seriesColors []string) StackedAreaSpec {
	if len(xValues) == 0 {
		return StackedAreaSpec{}
	}

	// Build points
	points := make([]StackedAreaPoint, len(xValues))
	for i, x := range xValues {
		point := StackedAreaPoint{
			X:      x,
			Values: make([]float64, len(seriesData)),
		}
		for j, series := range seriesData {
			if i < len(series) {
				point.Values[j] = series[i]
			}
		}
		points[i] = point
	}

	// Build series metadata
	series := make([]StackedAreaSeries, len(seriesData))
	for i := range seriesData {
		series[i] = StackedAreaSeries{}
		if i < len(seriesLabels) {
			series[i].Label = seriesLabels[i]
		}
		if i < len(seriesColors) {
			series[i].Color = seriesColors[i]
		}
	}

	return StackedAreaSpec{
		Points:  points,
		Series:  series,
		Width:   800,
		Height:  600,
		ShowGrid: true,
	}
}

// Helper function to calculate max value in 2D array
func max2D(data [][]float64) float64 {
	maxVal := math.Inf(-1)
	for _, row := range data {
		for _, val := range row {
			if val > maxVal {
				maxVal = val
			}
		}
	}
	return maxVal
}
