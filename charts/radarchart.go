package charts

import (
	"fmt"
	"math"

	"github.com/SCKelemen/svg"
	"github.com/SCKelemen/units"
)

// RadarAxis represents an axis in the radar chart
type RadarAxis struct {
	Label string
	Min   float64
	Max   float64
}

// RadarSeries represents a series of values for radar chart
type RadarSeries struct {
	Label      string
	Values     []float64 // One value per axis
	Color      string
	FillOpacity float64 // Opacity of filled area (0-1)
	LineWidth   float64
}

// RadarChartSpec configures radar chart rendering
type RadarChartSpec struct {
	Axes         []RadarAxis     // Axes definitions
	Series       []*RadarSeries  // Data series
	Width        float64
	Height       float64
	ShowGrid     bool            // Show concentric grid circles
	GridLevels   int             // Number of grid levels (default: 5)
	ShowLabels   bool            // Show axis labels
	ShowValues   bool            // Show value labels on points
	Title        string
}

// RenderRadarChart generates an SVG radar/spider chart
func RenderRadarChart(spec RadarChartSpec) string {
	if len(spec.Axes) == 0 || len(spec.Series) == 0 {
		return ""
	}

	// Set defaults
	if spec.GridLevels == 0 {
		spec.GridLevels = 5
	}

	// Calculate center and radius
	centerX := spec.Width / 2
	centerY := spec.Height / 2
	margin := 80.0
	radius := math.Min(spec.Width, spec.Height)/2 - margin

	numAxes := len(spec.Axes)
	angleStep := 2 * math.Pi / float64(numAxes)

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

	// Draw grid circles
	if spec.ShowGrid {
		gridStyle := svg.Style{
			Stroke:      "#e5e7eb",
			StrokeWidth: 1,
			Fill:        "none",
		}
		for level := 1; level <= spec.GridLevels; level++ {
			levelRadius := radius * float64(level) / float64(spec.GridLevels)
			result += svg.Circle(centerX, centerY, levelRadius, gridStyle) + "\n"
		}

		// Draw grid level labels
		labelStyle := svg.Style{
			FontSize:   units.Px(9),
			FontFamily: "sans-serif",
			Fill:       "#6b7280",
			TextAnchor: svg.TextAnchorStart,
		}
		for level := 1; level <= spec.GridLevels; level++ {
			levelRadius := radius * float64(level) / float64(spec.GridLevels)
			// Place label at top of circle
			percentage := float64(level) * 100 / float64(spec.GridLevels)
			result += svg.Text(fmt.Sprintf("%.0f%%", percentage), centerX+5, centerY-levelRadius, labelStyle) + "\n"
		}
	}

	// Draw axes
	axisStyle := svg.Style{
		Stroke:      "#9ca3af",
		StrokeWidth: 1.5,
	}
	for i, axis := range spec.Axes {
		angle := float64(i) * angleStep - math.Pi/2 // Start from top

		// Calculate end point
		endX := centerX + radius*math.Cos(angle)
		endY := centerY + radius*math.Sin(angle)

		// Draw axis line
		result += svg.Line(centerX, centerY, endX, endY, axisStyle) + "\n"

		// Draw axis label
		if spec.ShowLabels && axis.Label != "" {
			labelDistance := radius + 20
			labelX := centerX + labelDistance*math.Cos(angle)
			labelY := centerY + labelDistance*math.Sin(angle)

			// Adjust text anchor based on position
			var textAnchor string
			if math.Abs(math.Cos(angle)) < 0.1 {
				textAnchor = "middle"
			} else if math.Cos(angle) > 0 {
				textAnchor = "start"
			} else {
				textAnchor = "end"
			}

			labelStyle := svg.Style{
				FontSize:         units.Px(11),
				FontFamily:       "sans-serif",
				FontWeight:       "bold",
				TextAnchor:       svg.TextAnchor(textAnchor),
				DominantBaseline: svg.DominantBaselineMiddle,
			}
			result += svg.Text(axis.Label, labelX, labelY, labelStyle) + "\n"
		}
	}

	// Default colors
	defaultColors := []string{"#3b82f6", "#10b981", "#f59e0b", "#ef4444", "#8b5cf6", "#ec4899"}

	// Draw each series
	for seriesIdx, series := range spec.Series {
		if len(series.Values) != numAxes {
			continue // Skip series with wrong number of values
		}

		// Get series color
		seriesColor := series.Color
		if seriesColor == "" {
			seriesColor = defaultColors[seriesIdx%len(defaultColors)]
		}

		// Get fill opacity
		fillOpacity := series.FillOpacity
		if fillOpacity == 0 {
			fillOpacity = 0.3
		}

		// Get line width
		lineWidth := series.LineWidth
		if lineWidth == 0 {
			lineWidth = 2
		}

		// Calculate points for this series
		var pathData string
		points := make([]struct{ x, y float64 }, numAxes)

		for i := 0; i < numAxes; i++ {
			angle := float64(i) * angleStep - math.Pi/2

			// Normalize value to [0, 1] based on axis range
			axis := spec.Axes[i]
			value := series.Values[i]
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

			// Calculate point position
			pointRadius := radius * normalizedValue
			pointX := centerX + pointRadius*math.Cos(angle)
			pointY := centerY + pointRadius*math.Sin(angle)

			points[i].x = pointX
			points[i].y = pointY

			if i == 0 {
				pathData = fmt.Sprintf("M %.2f %.2f", pointX, pointY)
			} else {
				pathData += fmt.Sprintf(" L %.2f %.2f", pointX, pointY)
			}
		}
		pathData += " Z" // Close path

		// Draw filled area
		fillStyle := svg.Style{
			Fill:        seriesColor,
			FillOpacity: fillOpacity,
			Stroke:      "none",
		}
		result += svg.Path(pathData, fillStyle) + "\n"

		// Draw border line
		lineStyle := svg.Style{
			Stroke:      seriesColor,
			StrokeWidth: lineWidth,
			Fill:        "none",
		}
		result += svg.Path(pathData, lineStyle) + "\n"

		// Draw points
		pointStyle := svg.Style{
			Fill:        seriesColor,
			Stroke:      "#ffffff",
			StrokeWidth: 2,
		}
		for i, point := range points {
			result += svg.Circle(point.x, point.y, 4, pointStyle) + "\n"

			// Draw value labels if enabled
			if spec.ShowValues {
				valueStyle := svg.Style{
					FontSize:   units.Px(9),
					FontFamily: "sans-serif",
					Fill:       seriesColor,
					TextAnchor: svg.TextAnchorMiddle,
				}
				valueText := fmt.Sprintf("%.1f", series.Values[i])
				result += svg.Text(valueText, point.x, point.y-8, valueStyle) + "\n"
			}
		}
	}

	// Legend
	if len(spec.Series) > 1 {
		legendX := 20.0
		legendY := spec.Height - 60.0

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
				FillOpacity: 0.5,
				Stroke:      color,
				StrokeWidth: 2,
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

// NormalizedRadarChart creates a radar chart where all axes have the same 0-100 scale
func NormalizedRadarChart(axisLabels []string, seriesData map[string][]float64, width, height float64) string {
	// Create axes with normalized 0-100 range
	axes := make([]RadarAxis, len(axisLabels))
	for i, label := range axisLabels {
		axes[i] = RadarAxis{
			Label: label,
			Min:   0,
			Max:   100,
		}
	}

	// Create series
	series := make([]*RadarSeries, 0, len(seriesData))
	for label, values := range seriesData {
		series = append(series, &RadarSeries{
			Label:  label,
			Values: values,
		})
	}

	spec := RadarChartSpec{
		Axes:        axes,
		Series:      series,
		Width:       width,
		Height:      height,
		ShowGrid:    true,
		GridLevels:  5,
		ShowLabels:  true,
		ShowValues:  false,
	}

	return RenderRadarChart(spec)
}
