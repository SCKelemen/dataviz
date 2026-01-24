package charts

import (
	"fmt"
	"math"

	"github.com/SCKelemen/svg"
	"github.com/SCKelemen/units"
)

// CircularBarPoint represents a single bar in the circular barplot
type CircularBarPoint struct {
	Label string
	Value float64
	Color string // Optional custom color
}

// CircularBarPlotSpec configures circular barplot rendering
type CircularBarPlotSpec struct {
	Data         []CircularBarPoint
	Width        float64
	Height       float64
	InnerRadius  float64 // Inner radius (creates donut hole, 0 = start from center)
	BarWidth     float64 // Angular width of bars (0 = auto-calculate with gaps)
	DefaultColor string  // Default bar color
	ShowLabels   bool    // Show value labels on bars
	ShowAxisLabels bool  // Show labels around the circle
	Title        string
	StartAngle   float64 // Starting angle in degrees (0 = top, clockwise)
}

// RenderCircularBarPlot generates an SVG circular barplot
func RenderCircularBarPlot(spec CircularBarPlotSpec) string {
	if len(spec.Data) == 0 {
		return ""
	}

	// Set defaults
	if spec.DefaultColor == "" {
		spec.DefaultColor = "#3b82f6"
	}
	if spec.InnerRadius < 0 {
		spec.InnerRadius = 0
	}

	// Calculate center and radius
	centerX := spec.Width / 2
	centerY := spec.Height / 2
	margin := 80.0
	maxRadius := math.Min(spec.Width, spec.Height)/2 - margin

	// Find max value for scaling
	maxValue := 0.0
	for _, point := range spec.Data {
		if point.Value > maxValue {
			maxValue = point.Value
		}
	}
	if maxValue == 0 {
		maxValue = 1
	}

	// Calculate angular spacing
	numBars := len(spec.Data)
	angleStep := 360.0 / float64(numBars) // degrees per bar slot

	// Calculate bar angular width (with gaps)
	barAngleWidth := spec.BarWidth
	if barAngleWidth == 0 {
		barAngleWidth = angleStep * 0.8 // 80% of available space, 20% gap
	}

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

	// Draw circular grid lines (optional)
	gridStyle := svg.Style{
		Stroke:      "#e5e7eb",
		StrokeWidth: 1,
		Fill:        "none",
		Opacity:     0.5,
	}

	// Draw reference circles at 25%, 50%, 75%, 100%
	for i := 1; i <= 4; i++ {
		radius := spec.InnerRadius + (maxRadius-spec.InnerRadius)*float64(i)/4
		result += svg.Circle(centerX, centerY, radius, gridStyle) + "\n"
	}

	// Draw each bar
	for i, point := range spec.Data {
		// Calculate angle for this bar (center of bar slot)
		angleDeg := spec.StartAngle + float64(i)*angleStep
		angleRad := (angleDeg - 90) * math.Pi / 180 // Convert to radians, adjust for SVG coords

		// Calculate bar start and end angles
		halfBarAngle := barAngleWidth / 2
		startAngleDeg := angleDeg - halfBarAngle
		endAngleDeg := angleDeg + halfBarAngle

		// Calculate bar radius (scaled by value)
		barLength := (maxRadius - spec.InnerRadius) * (point.Value / maxValue)
		outerRadius := spec.InnerRadius + barLength

		// Get bar color
		barColor := point.Color
		if barColor == "" {
			barColor = spec.DefaultColor
		}

		// Draw bar as a path (annular sector)
		barPath := createAnnularSector(centerX, centerY, spec.InnerRadius, outerRadius, startAngleDeg, endAngleDeg)

		barStyle := svg.Style{
			Fill:        barColor,
			Stroke:      "#ffffff",
			StrokeWidth: 1,
		}
		result += svg.Path(barPath, barStyle) + "\n"

		// Draw value label on bar
		if spec.ShowLabels {
			labelRadius := spec.InnerRadius + barLength/2
			labelX := centerX + labelRadius*math.Cos(angleRad)
			labelY := centerY + labelRadius*math.Sin(angleRad)

			valueLabelStyle := svg.Style{
				FontSize:         units.Px(9),
				FontFamily:       "sans-serif",
				Fill:             "#ffffff",
				TextAnchor:       svg.TextAnchorMiddle,
				DominantBaseline: svg.DominantBaselineMiddle,
			}
			result += svg.Text(fmt.Sprintf("%.0f", point.Value), labelX, labelY, valueLabelStyle) + "\n"
		}

		// Draw axis label (outside the bars)
		if spec.ShowAxisLabels && point.Label != "" {
			labelDistance := maxRadius + 15
			labelX := centerX + labelDistance*math.Cos(angleRad)
			labelY := centerY + labelDistance*math.Sin(angleRad)

			// Adjust text anchor based on position
			var textAnchor string
			if math.Abs(math.Cos(angleRad)) < 0.1 {
				textAnchor = "middle"
			} else if math.Cos(angleRad) > 0 {
				textAnchor = "start"
			} else {
				textAnchor = "end"
			}

			axisLabelStyle := svg.Style{
				FontSize:         units.Px(10),
				FontFamily:       "sans-serif",
				TextAnchor:       svg.TextAnchor(textAnchor),
				DominantBaseline: svg.DominantBaselineMiddle,
			}
			result += svg.Text(point.Label, labelX, labelY, axisLabelStyle) + "\n"
		}
	}

	// Draw center circle (if inner radius > 0)
	if spec.InnerRadius > 0 {
		centerCircleStyle := svg.Style{
			Fill:   "#ffffff",
			Stroke: "#d1d5db",
			StrokeWidth: 1,
		}
		result += svg.Circle(centerX, centerY, spec.InnerRadius, centerCircleStyle) + "\n"
	}

	return result
}

// createAnnularSector creates an SVG path for an annular sector (ring segment)
func createAnnularSector(cx, cy, innerRadius, outerRadius, startAngleDeg, endAngleDeg float64) string {
	// Convert angles to radians
	startAngle := (startAngleDeg - 90) * math.Pi / 180
	endAngle := (endAngleDeg - 90) * math.Pi / 180

	// Calculate points
	x1 := cx + innerRadius*math.Cos(startAngle)
	y1 := cy + innerRadius*math.Sin(startAngle)
	x2 := cx + outerRadius*math.Cos(startAngle)
	y2 := cy + outerRadius*math.Sin(startAngle)
	x3 := cx + outerRadius*math.Cos(endAngle)
	y3 := cy + outerRadius*math.Sin(endAngle)
	x4 := cx + innerRadius*math.Cos(endAngle)
	y4 := cy + innerRadius*math.Sin(endAngle)

	// Determine if we need large arc flag
	largeArcFlag := 0
	if endAngleDeg-startAngleDeg > 180 {
		largeArcFlag = 1
	}

	// Build path
	path := fmt.Sprintf("M %.2f %.2f ", x1, y1)                                                  // Move to inner start
	path += fmt.Sprintf("L %.2f %.2f ", x2, y2)                                                  // Line to outer start
	path += fmt.Sprintf("A %.2f %.2f 0 %d 1 %.2f %.2f ", outerRadius, outerRadius, largeArcFlag, x3, y3) // Arc along outer edge
	path += fmt.Sprintf("L %.2f %.2f ", x4, y4)                                                  // Line to inner end
	path += fmt.Sprintf("A %.2f %.2f 0 %d 0 %.2f %.2f ", innerRadius, innerRadius, largeArcFlag, x1, y1) // Arc along inner edge
	path += "Z" // Close path

	return path
}

// CircularBarPlotFromValues creates a circular barplot from simple value arrays
func CircularBarPlotFromValues(labels []string, values []float64, width, height float64) string {
	if len(labels) != len(values) {
		return ""
	}

	data := make([]CircularBarPoint, len(labels))
	for i := range labels {
		data[i] = CircularBarPoint{
			Label: labels[i],
			Value: values[i],
		}
	}

	spec := CircularBarPlotSpec{
		Data:           data,
		Width:          width,
		Height:         height,
		InnerRadius:    50,
		ShowLabels:     false,
		ShowAxisLabels: true,
	}

	return RenderCircularBarPlot(spec)
}

// MultiColorCircularBarPlot creates a circular barplot with different colors
func MultiColorCircularBarPlot(labels []string, values []float64, colors []string, width, height float64) string {
	if len(labels) != len(values) {
		return ""
	}

	data := make([]CircularBarPoint, len(labels))
	for i := range labels {
		color := ""
		if i < len(colors) {
			color = colors[i]
		}
		data[i] = CircularBarPoint{
			Label: labels[i],
			Value: values[i],
			Color: color,
		}
	}

	spec := CircularBarPlotSpec{
		Data:           data,
		Width:          width,
		Height:         height,
		InnerRadius:    60,
		ShowLabels:     true,
		ShowAxisLabels: true,
	}

	return RenderCircularBarPlot(spec)
}
