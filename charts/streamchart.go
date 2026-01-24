package charts

import (
	"fmt"
	"math"

	"github.com/SCKelemen/dataviz/scales"
	"github.com/SCKelemen/svg"
	"github.com/SCKelemen/units"
)

// StreamPoint represents a single X position with values for each series
type StreamPoint struct {
	X      float64
	Values []float64 // One value per series
}

// StreamSeries represents metadata for a stream series
type StreamSeries struct {
	Label string
	Color string
}

// StreamChartSpec configures streamchart rendering
type StreamChartSpec struct {
	Points       []StreamPoint  // X positions with values for all series
	Series       []StreamSeries // Metadata for each series (colors, labels)
	Width        float64
	Height       float64
	Layout       string  // "center", "wiggle", "silhouette" (default: center)
	Smooth       bool    // Use smooth curves
	Title        string
	XAxisLabel   string
	ShowLegend   bool
}

// RenderStreamChart generates an SVG streamchart
func RenderStreamChart(spec StreamChartSpec) string {
	if len(spec.Points) == 0 || len(spec.Series) == 0 {
		return ""
	}

	// Set defaults
	if spec.Layout == "" {
		spec.Layout = "center"
	}

	// Calculate margins
	margin := 60.0
	chartWidth := spec.Width - (2 * margin)

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

	// Calculate baseline offsets using the chosen layout algorithm
	var baselines [][]float64 // baselines[pointIdx][seriesIdx]
	var yMin, yMax float64

	switch spec.Layout {
	case "wiggle":
		baselines, yMin, yMax = calculateWiggleLayout(spec.Points, spec.Series)
	case "silhouette":
		baselines, yMin, yMax = calculateSilhouetteLayout(spec.Points, spec.Series)
	default: // "center"
		baselines, yMin, yMax = calculateCenterLayout(spec.Points, spec.Series)
	}

	// Add padding to Y range
	yRange := yMax - yMin
	if yRange == 0 {
		yRange = 1
	}
	yMin -= yRange * 0.05
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

	// Draw center line if using center layout
	if spec.Layout == "center" {
		centerY := yScale.Apply(0).Value
		centerLineStyle := svg.Style{
			Stroke:      "#d1d5db",
			StrokeWidth: 1,
			Opacity:     0.5,
		}
		result += svg.Line(margin, centerY, margin+chartWidth, centerY, centerLineStyle) + "\n"
	}

	// Default colors
	defaultColors := []string{"#3b82f6", "#10b981", "#f59e0b", "#ef4444", "#8b5cf6", "#ec4899"}

	// Draw each stream layer
	for seriesIdx := 0; seriesIdx < len(spec.Series); seriesIdx++ {
		series := spec.Series[seriesIdx]

		// Get series color
		seriesColor := series.Color
		if seriesColor == "" {
			seriesColor = defaultColors[seriesIdx%len(defaultColors)]
		}

		// Build path for this stream layer
		var pathData string

		// Top edge (left to right)
		for i, point := range spec.Points {
			if seriesIdx >= len(point.Values) {
				continue
			}

			x := xScale.Apply(point.X).Value
			// Top of stream = baseline + value
			topY := baselines[i][seriesIdx] + point.Values[seriesIdx]
			y := yScale.Apply(topY).Value

			if i == 0 {
				pathData = fmt.Sprintf("M %.2f %.2f", x, y)
			} else {
				pathData += fmt.Sprintf(" L %.2f %.2f", x, y)
			}
		}

		// Bottom edge (right to left)
		for i := len(spec.Points) - 1; i >= 0; i-- {
			x := xScale.Apply(spec.Points[i].X).Value
			// Bottom of stream = baseline
			bottomY := baselines[i][seriesIdx]
			y := yScale.Apply(bottomY).Value
			pathData += fmt.Sprintf(" L %.2f %.2f", x, y)
		}

		pathData += " Z" // Close path

		// Draw filled stream
		fillStyle := svg.Style{
			Fill:        seriesColor,
			FillOpacity: 0.8,
			Stroke:      seriesColor,
			StrokeWidth: 0.5,
		}
		result += svg.Path(pathData, fillStyle) + "\n"
	}

	// X-axis label
	if spec.XAxisLabel != "" {
		xTitleStyle := svg.Style{
			FontSize:   units.Px(12),
			FontFamily: "sans-serif",
			TextAnchor: svg.TextAnchorMiddle,
		}
		result += svg.Text(spec.XAxisLabel, spec.Width/2, spec.Height-10, xTitleStyle) + "\n"
	}

	// Legend
	if spec.ShowLegend && len(spec.Series) > 0 {
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
				FillOpacity: 0.8,
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

// calculateCenterLayout centers the stream around y=0
func calculateCenterLayout(points []StreamPoint, series []StreamSeries) ([][]float64, float64, float64) {
	numPoints := len(points)
	numSeries := len(series)

	baselines := make([][]float64, numPoints)
	yMin := math.Inf(1)
	yMax := math.Inf(-1)

	for i, point := range points {
		baselines[i] = make([]float64, numSeries)

		// Calculate total for this point
		total := 0.0
		for j := 0; j < numSeries && j < len(point.Values); j++ {
			total += point.Values[j]
		}

		// Center around y=0
		offset := -total / 2

		for j := 0; j < numSeries; j++ {
			baselines[i][j] = offset

			if j < len(point.Values) {
				// Update min/max
				bottom := offset
				top := offset + point.Values[j]

				if bottom < yMin {
					yMin = bottom
				}
				if top > yMax {
					yMax = top
				}

				// Next layer starts where this one ends
				offset += point.Values[j]
			}
		}
	}

	return baselines, yMin, yMax
}

// calculateSilhouetteLayout creates a silhouette with baseline at 0
func calculateSilhouetteLayout(points []StreamPoint, series []StreamSeries) ([][]float64, float64, float64) {
	numPoints := len(points)
	numSeries := len(series)

	baselines := make([][]float64, numPoints)
	yMin := 0.0
	yMax := 0.0

	for i, point := range points {
		baselines[i] = make([]float64, numSeries)

		offset := 0.0
		for j := 0; j < numSeries; j++ {
			baselines[i][j] = offset

			if j < len(point.Values) {
				top := offset + point.Values[j]
				if top > yMax {
					yMax = top
				}
				offset += point.Values[j]
			}
		}
	}

	return baselines, yMin, yMax
}

// calculateWiggleLayout uses a wiggle algorithm to minimize visual artifacts
// This is a simplified version that tries to keep streams centered and flowing
func calculateWiggleLayout(points []StreamPoint, series []StreamSeries) ([][]float64, float64, float64) {
	numPoints := len(points)
	numSeries := len(series)

	baselines := make([][]float64, numPoints)
	yMin := math.Inf(1)
	yMax := math.Inf(-1)

	// For each point, calculate a weighted center based on adjacent points
	for i, point := range points {
		baselines[i] = make([]float64, numSeries)

		// Calculate total for this point
		total := 0.0
		for j := 0; j < numSeries && j < len(point.Values); j++ {
			total += point.Values[j]
		}

		// Calculate weighted offset considering neighbors
		var prevTotal, nextTotal float64
		if i > 0 {
			for j := 0; j < numSeries && j < len(points[i-1].Values); j++ {
				prevTotal += points[i-1].Values[j]
			}
		}
		if i < numPoints-1 {
			for j := 0; j < numSeries && j < len(points[i+1].Values); j++ {
				nextTotal += points[i+1].Values[j]
			}
		}

		// Wiggle tries to minimize slope changes
		avgTotal := (prevTotal + total + nextTotal) / 3
		if avgTotal == 0 {
			avgTotal = total
		}
		offset := -avgTotal / 2

		for j := 0; j < numSeries; j++ {
			baselines[i][j] = offset

			if j < len(point.Values) {
				bottom := offset
				top := offset + point.Values[j]

				if bottom < yMin {
					yMin = bottom
				}
				if top > yMax {
					yMax = top
				}

				offset += point.Values[j]
			}
		}
	}

	return baselines, yMin, yMax
}

// StreamChartFromSeries is a helper to convert multiple simple series into streamchart format
func StreamChartFromSeries(xValues []float64, seriesData [][]float64, seriesLabels []string, seriesColors []string) StreamChartSpec {
	if len(xValues) == 0 {
		return StreamChartSpec{}
	}

	// Build points
	points := make([]StreamPoint, len(xValues))
	for i, x := range xValues {
		point := StreamPoint{
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
	series := make([]StreamSeries, len(seriesData))
	for i := range seriesData {
		series[i] = StreamSeries{}
		if i < len(seriesLabels) {
			series[i].Label = seriesLabels[i]
		}
		if i < len(seriesColors) {
			series[i].Color = seriesColors[i]
		}
	}

	return StreamChartSpec{
		Points:     points,
		Series:     series,
		Width:      800,
		Height:     600,
		Layout:     "wiggle",
		ShowLegend: true,
	}
}
