package charts

import (
	"fmt"
	"math"

	"github.com/SCKelemen/dataviz/scales"
	"github.com/SCKelemen/svg"
	"github.com/SCKelemen/units"
)

// DensityPlotData represents data for a density plot
// This is already defined in histogram.go, but we'll use it here too

// DensityPlotSpec is already defined in histogram.go
// We'll create a simplified standalone version

// SimpleDensityData represents data for a single density curve
type SimpleDensityData struct {
	Values    []float64
	Label     string
	Color     string
	Bandwidth float64 // KDE bandwidth (0 = auto)
}

// SimpleDensitySpec configures simple density plot rendering
type SimpleDensitySpec struct {
	Data       []*SimpleDensityData
	Width      float64
	Height     float64
	ShowFill   bool
	ShowRug    bool    // Show rug plot (data points on x-axis)
	LineWidth  float64
	Title      string
	XAxisLabel string
	YAxisLabel string
}

// RenderDensityPlot renders a standalone density plot (already exists in histogram.go)
// We'll create RenderSimpleDensity for cleaner API

// RenderSimpleDensity renders a simple density plot using KDE
func RenderSimpleDensity(spec SimpleDensitySpec) string {
	if len(spec.Data) == 0 {
		return ""
	}

	// Calculate margin
	margin := 60.0

	// Find global min/max
	globalMin := math.Inf(1)
	globalMax := math.Inf(-1)
	maxDensity := 0.0

	// Calculate KDE for each dataset
	type densityCurve struct {
		data   *SimpleDensityData
		points []DensityPoint
	}
	curves := make([]densityCurve, 0, len(spec.Data))

	for _, data := range spec.Data {
		if len(data.Values) == 0 {
			continue
		}

		// Find min/max
		for _, v := range data.Values {
			if v < globalMin {
				globalMin = v
			}
			if v > globalMax {
				globalMax = v
			}
		}

		// Calculate KDE
		density := calculateKDE(data.Values, data.Bandwidth)

		// Find max density
		for _, dp := range density {
			if dp.Density > maxDensity {
				maxDensity = dp.Density
			}
		}

		curves = append(curves, densityCurve{
			data:   data,
			points: density,
		})
	}

	// Create scales
	xScale := scales.NewLinearScale(
		[2]float64{globalMin, globalMax},
		[2]units.Length{units.Px(margin), units.Px(spec.Width - margin)},
	)

	yScale := scales.NewLinearScale(
		[2]float64{0, maxDensity * 1.1}, // Add 10% headroom
		[2]units.Length{units.Px(spec.Height - margin), units.Px(margin)},
	)

	var result string

	// Title
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

	// Draw axes
	axisStyle := svg.Style{
		Stroke:      "#374151",
		StrokeWidth: 2,
	}
	result += svg.Line(margin, margin, margin, spec.Height-margin, axisStyle) + "\n"
	result += svg.Line(margin, spec.Height-margin, spec.Width-margin, spec.Height-margin, axisStyle) + "\n"

	// Line width
	lineWidth := spec.LineWidth
	if lineWidth == 0 {
		lineWidth = 2.5
	}

	// Colors for multiple curves
	defaultColors := []string{"#3b82f6", "#10b981", "#f59e0b", "#ef4444", "#8b5cf6", "#ec4899"}

	// Render each density curve
	for idx, curve := range curves {
		if len(curve.points) == 0 {
			continue
		}

		// Get color
		lineColor := curve.data.Color
		if lineColor == "" {
			lineColor = defaultColors[idx%len(defaultColors)]
		}

		// Build path
		var pathData string
		for j, dp := range curve.points {
			x := xScale.Apply(dp.Value).Value
			y := yScale.Apply(dp.Density).Value

			if j == 0 {
				pathData = fmt.Sprintf("M %.2f %.2f", x, y)
			} else {
				pathData += fmt.Sprintf(" L %.2f %.2f", x, y)
			}
		}

		// Draw filled area if enabled
		if spec.ShowFill {
			baseY := yScale.Apply(0).Value
			lastX := xScale.Apply(curve.points[len(curve.points)-1].Value).Value
			firstX := xScale.Apply(curve.points[0].Value).Value

			fillPath := fmt.Sprintf("M %.2f %.2f", firstX, baseY)
			fillPath += " " + pathData[1:] // Skip the M command
			fillPath += fmt.Sprintf(" L %.2f %.2f Z", lastX, baseY)

			fillStyle := svg.Style{
				Fill:    lineColor,
				Opacity: 0.3,
				Stroke:  "none",
			}
			result += svg.Path(fillPath, fillStyle) + "\n"
		}

		// Draw line
		lineStyle := svg.Style{
			Stroke:      lineColor,
			StrokeWidth: lineWidth,
			Fill:        "none",
		}
		result += svg.Path(pathData, lineStyle) + "\n"

		// Draw rug plot if enabled
		if spec.ShowRug {
			rugY := spec.Height - margin
			rugHeight := 10.0
			rugStyle := svg.Style{
				Stroke:      lineColor,
				StrokeWidth: 1,
				Opacity:     0.5,
			}
			for _, val := range curve.data.Values {
				x := xScale.Apply(val).Value
				result += svg.Line(x, rugY, x, rugY+rugHeight, rugStyle) + "\n"
			}
		}
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
		value := globalMin + (globalMax-globalMin)/float64(steps)*float64(i)
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
		value := (maxDensity * 1.1) / float64(steps) * float64(i)
		y := yScale.Apply(value).Value
		result += svg.Text(fmt.Sprintf("%.3f", value), margin-10, y, yLabelStyle) + "\n"
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

	// Legend if multiple curves
	if len(curves) > 1 {
		legendX := margin + 20
		legendY := margin + 20
		for idx, curve := range curves {
			if curve.data.Label == "" {
				continue
			}

			color := curve.data.Color
			if color == "" {
				color = defaultColors[idx%len(defaultColors)]
			}

			yOffset := float64(idx * 25)

			// Legend line
			lineStyle := svg.Style{
				Stroke:      color,
				StrokeWidth: 3,
			}
			result += svg.Line(legendX, legendY+yOffset, legendX+30, legendY+yOffset, lineStyle) + "\n"

			// Legend text
			textStyle := svg.Style{
				FontSize:         units.Px(12),
				FontFamily:       "sans-serif",
				DominantBaseline: svg.DominantBaselineMiddle,
			}
			result += svg.Text(curve.data.Label, legendX+35, legendY+yOffset, textStyle) + "\n"
		}
	}

	return result
}
