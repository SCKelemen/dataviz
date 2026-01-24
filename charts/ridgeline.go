package charts

import (
	"fmt"

	"github.com/SCKelemen/dataviz/scales"
	"github.com/SCKelemen/svg"
	"github.com/SCKelemen/units"
)

// RidgelineData represents one ridge in a ridgeline plot
type RidgelineData struct {
	Values    []float64 // Raw data values
	Label     string    // Category label
	Color     string    // Fill color
	Bandwidth float64   // KDE bandwidth (0 = auto)
}

// RidgelineSpec configures ridgeline plot rendering
type RidgelineSpec struct {
	Data       []*RidgelineData
	Width      float64
	Height     float64
	Overlap    float64 // Amount of overlap between ridges (0-1, default 0.5)
	ShowFill   bool    // If true, fill ridges with color
	LineWidth  float64 // Line width
	Reverse    bool    // If true, reverse order (top to bottom)

	// Axis configuration
	XAxisLabel string
	ShowLabels bool // If true, show category labels on Y axis
}

// RenderRidgeline renders a ridgeline (joy) plot
func RenderRidgeline(spec RidgelineSpec) string {
	if len(spec.Data) == 0 {
		return ""
	}

	// Find global min/max for X axis
	globalMin := spec.Data[0].Values[0]
	globalMax := spec.Data[0].Values[0]

	for _, ridge := range spec.Data {
		for _, v := range ridge.Values {
			if v < globalMin {
				globalMin = v
			}
			if v > globalMax {
				globalMax = v
			}
		}
	}

	// Calculate KDE for each ridge
	allDensities := make([][]DensityPoint, len(spec.Data))
	maxDensities := make([]float64, len(spec.Data))

	for i, ridge := range spec.Data {
		stats := CalculateViolinStats(ridge.Values, ridge.Bandwidth)
		allDensities[i] = stats.Density

		// Find max density for this ridge
		maxDensity := 0.0
		for _, dp := range stats.Density {
			if dp.Density > maxDensity {
				maxDensity = dp.Density
			}
		}
		maxDensities[i] = maxDensity
	}

	// Create X scale
	xScale := scales.NewLinearScale(
		[2]float64{globalMin, globalMax},
		[2]units.Length{units.Px(80), units.Px(spec.Width - 40)},
	)

	// Calculate ridge heights and positions
	overlap := spec.Overlap
	if overlap == 0 {
		overlap = 0.5
	}
	if overlap < 0 {
		overlap = 0
	}
	if overlap > 1 {
		overlap = 1
	}

	ridgeHeight := spec.Height / float64(len(spec.Data)+1)
	effectiveGap := ridgeHeight * (1 - overlap)

	lineWidth := spec.LineWidth
	if lineWidth == 0 {
		lineWidth = 1.5
	}

	var result string

	// Render each ridge
	for i, densities := range allDensities {
		if len(densities) == 0 {
			continue
		}

		ridge := spec.Data[i]

		// Calculate Y position for this ridge
		ridgeIndex := i
		if spec.Reverse {
			ridgeIndex = len(spec.Data) - 1 - i
		}
		baseY := 40 + float64(ridgeIndex)*effectiveGap

		// Build path for the density curve
		var pathPoints string

		for j, dp := range densities {
			x := xScale.Apply(dp.Value).Value

			// Scale density to ridge height
			densityHeight := (dp.Density / maxDensities[i]) * ridgeHeight
			y := baseY + ridgeHeight - densityHeight

			if j == 0 {
				pathPoints = fmt.Sprintf("%.2f,%.2f", x, y)
			} else {
				pathPoints += fmt.Sprintf(" %.2f,%.2f", x, y)
			}
		}

		// Line color
		lineColor := ridge.Color
		if lineColor == "" {
			colors := []string{"#4285f4", "#ea4335", "#fbbc04", "#34a853", "#ff6d00", "#46bdc6"}
			lineColor = colors[i%len(colors)]
		}

		// Draw filled ridge if enabled
		if spec.ShowFill {
			// Create closed path with baseline
			firstX := xScale.Apply(densities[0].Value).Value
			lastX := xScale.Apply(densities[len(densities)-1].Value).Value
			ridgeBaseline := baseY + ridgeHeight

			fillPath := fmt.Sprintf("M %.2f %.2f L %s L %.2f %.2f Z",
				firstX, ridgeBaseline,
				pathPoints,
				lastX, ridgeBaseline)

			fillStyle := svg.Style{
				Fill:    lineColor,
				Opacity: 0.7,
				Stroke:  "none",
			}

			result += svg.Path(fillPath, fillStyle) + "\n"
		}

		// Draw outline
		linePath := "M " + pathPoints
		lineStyle := svg.Style{
			Stroke:      lineColor,
			StrokeWidth: lineWidth,
			Fill:        "none",
		}

		result += svg.Path(linePath, lineStyle) + "\n"

		// Draw label if enabled
		if spec.ShowLabels && ridge.Label != "" {
			labelY := baseY + ridgeHeight/2

			labelStyle := svg.Style{
				FontSize:         units.Px(12),
				FontFamily:       "sans-serif",
				TextAnchor:       svg.TextAnchorEnd,
				DominantBaseline: svg.DominantBaselineMiddle,
			}

			result += svg.Text(ridge.Label, 70, labelY, labelStyle) + "\n"
		}
	}

	return result
}

// RidgelineFromGroups creates ridgeline data from grouped data
// Useful for converting transforms.Group output to ridgeline format
func RidgelineFromGroups(groups map[string][]float64, colors map[string]string) []*RidgelineData {
	ridges := make([]*RidgelineData, 0, len(groups))

	for label, values := range groups {
		color := ""
		if colors != nil {
			color = colors[label]
		}

		ridges = append(ridges, &RidgelineData{
			Values: values,
			Label:  label,
			Color:  color,
		})
	}

	return ridges
}
