package charts

import (
	"fmt"
	"math"

	"github.com/SCKelemen/dataviz/scales"
	"github.com/SCKelemen/svg"
	"github.com/SCKelemen/units"
)

// ConnectedScatterPoint represents a point in a connected scatter plot
type ConnectedScatterPoint struct {
	X     float64
	Y     float64
	Label string // Optional label for the point
	Size  float64 // Optional custom size for this point
	Color string  // Optional custom color for this point
}

// ConnectedScatterSeries represents a series of connected points
type ConnectedScatterSeries struct {
	Points     []ConnectedScatterPoint
	Label      string  // Series label for legend
	Color      string  // Line and marker color
	LineStyle  string  // "solid", "dashed", "dotted"
	LineWidth  float64 // Width of connecting line
	MarkerType string  // "circle", "square", "diamond", "triangle", "cross", "x"
	MarkerSize float64 // Size of markers
}

// ConnectedScatterSpec configures connected scatter plot rendering
type ConnectedScatterSpec struct {
	Series       []*ConnectedScatterSeries
	Width        float64
	Height       float64
	ShowGrid     bool
	ShowMarkers  bool // If false, only lines are shown
	ShowLines    bool // If false, only markers are shown (regular scatter)
	Title        string
	XAxisLabel   string
	YAxisLabel   string
	XAxisMin     *float64 // Optional: force X axis min
	XAxisMax     *float64 // Optional: force X axis max
	YAxisMin     *float64 // Optional: force Y axis min
	YAxisMax     *float64 // Optional: force Y axis max
}

// RenderConnectedScatter generates an SVG connected scatter plot
func RenderConnectedScatter(spec ConnectedScatterSpec) string {
	if len(spec.Series) == 0 {
		return ""
	}

	// Calculate margins
	margin := 60.0
	chartWidth := spec.Width - (2 * margin)
	chartHeight := spec.Height - (2 * margin)

	// Find global min/max for all series
	xMin := math.Inf(1)
	xMax := math.Inf(-1)
	yMin := math.Inf(1)
	yMax := math.Inf(-1)

	for _, series := range spec.Series {
		for _, point := range series.Points {
			if point.X < xMin {
				xMin = point.X
			}
			if point.X > xMax {
				xMax = point.X
			}
			if point.Y < yMin {
				yMin = point.Y
			}
			if point.Y > yMax {
				yMax = point.Y
			}
		}
	}

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

	// Add padding to ranges
	xRange := xMax - xMin
	yRange := yMax - yMin
	if xRange == 0 {
		xRange = 1
	}
	if yRange == 0 {
		yRange = 1
	}
	xMin -= xRange * 0.05
	xMax += xRange * 0.05
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

	// Default colors for multiple series
	defaultColors := []string{"#3b82f6", "#10b981", "#f59e0b", "#ef4444", "#8b5cf6", "#ec4899"}

	// Draw each series
	for seriesIdx, series := range spec.Series {
		if len(series.Points) == 0 {
			continue
		}

		// Get series color
		seriesColor := series.Color
		if seriesColor == "" {
			seriesColor = defaultColors[seriesIdx%len(defaultColors)]
		}

		// Get line width
		lineWidth := series.LineWidth
		if lineWidth == 0 {
			lineWidth = 2
		}

		// Get marker size
		markerSize := series.MarkerSize
		if markerSize == 0 {
			markerSize = 5
		}

		// Get marker type
		markerType := series.MarkerType
		if markerType == "" {
			markerType = "circle"
		}

		// Draw connecting lines
		if spec.ShowLines {
			var pathData string
			for i, point := range series.Points {
				x := xScale.Apply(point.X).Value
				y := yScale.Apply(point.Y).Value

				if i == 0 {
					pathData = fmt.Sprintf("M %.2f %.2f", x, y)
				} else {
					pathData += fmt.Sprintf(" L %.2f %.2f", x, y)
				}
			}

			lineStyle := svg.Style{
				Stroke:      seriesColor,
				StrokeWidth: lineWidth,
				Fill:        "none",
			}

			// Note: StrokeDashArray not yet supported in svg.Style
			// TODO: Add support for dashed/dotted lines when available

			result += svg.Path(pathData, lineStyle) + "\n"
		}

		// Draw markers
		if spec.ShowMarkers {
			for _, point := range series.Points {
				x := xScale.Apply(point.X).Value
				y := yScale.Apply(point.Y).Value

				// Get point-specific overrides
				pointColor := point.Color
				if pointColor == "" {
					pointColor = seriesColor
				}

				pointSize := point.Size
				if pointSize == 0 {
					pointSize = markerSize
				}

				markerStyle := svg.Style{
					Fill:        pointColor,
					Stroke:      "#ffffff",
					StrokeWidth: 2,
				}

				// Draw marker based on type
				switch markerType {
				case "circle":
					result += svg.Circle(x, y, pointSize, markerStyle) + "\n"
				case "square":
					halfSize := pointSize
					result += svg.Rect(x-halfSize, y-halfSize, halfSize*2, halfSize*2, markerStyle) + "\n"
				case "diamond":
					diamondPoints := []svg.Point{
						{X: x, Y: y - pointSize},
						{X: x + pointSize, Y: y},
						{X: x, Y: y + pointSize},
						{X: x - pointSize, Y: y},
					}
					result += svg.Polygon(diamondPoints, markerStyle) + "\n"
				case "triangle":
					trianglePoints := []svg.Point{
						{X: x, Y: y - pointSize},
						{X: x + pointSize, Y: y + pointSize},
						{X: x - pointSize, Y: y + pointSize},
					}
					result += svg.Polygon(trianglePoints, markerStyle) + "\n"
				case "cross":
					crossStyle := svg.Style{
						Stroke:        pointColor,
						StrokeWidth:   2,
						StrokeLinecap: svg.StrokeLinecapRound,
					}
					result += svg.Line(x, y-pointSize, x, y+pointSize, crossStyle) + "\n"
					result += svg.Line(x-pointSize, y, x+pointSize, y, crossStyle) + "\n"
				case "x":
					xStyle := svg.Style{
						Stroke:        pointColor,
						StrokeWidth:   2,
						StrokeLinecap: svg.StrokeLinecapRound,
					}
					offset := pointSize * 0.7
					result += svg.Line(x-offset, y-offset, x+offset, y+offset, xStyle) + "\n"
					result += svg.Line(x-offset, y+offset, x+offset, y-offset, xStyle) + "\n"
				default:
					result += svg.Circle(x, y, pointSize, markerStyle) + "\n"
				}

				// Draw point label if specified
				if point.Label != "" {
					labelStyle := svg.Style{
						FontSize:         units.Px(9),
						FontFamily:       "sans-serif",
						TextAnchor:       svg.TextAnchorMiddle,
						DominantBaseline: svg.DominantBaselineHanging,
					}
					result += svg.Text(point.Label, x, y+pointSize+3, labelStyle) + "\n"
				}
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

	// Legend if multiple series
	if len(spec.Series) > 1 {
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
			result += svg.Text(series.Label, legendX+35, legendY+yOffset, textStyle) + "\n"
		}
	}

	return result
}
