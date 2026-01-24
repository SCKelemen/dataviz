package charts

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/SCKelemen/color"
	design "github.com/SCKelemen/design-system"
	"github.com/SCKelemen/dataviz/axes"
	"github.com/SCKelemen/dataviz/charts/legends"
	"github.com/SCKelemen/dataviz/scales"
	"github.com/SCKelemen/svg"
	"github.com/SCKelemen/units"
)

// Global counter for unique gradient IDs
var gradientCounter int64

// RenderLineGraph renders a line graph using scales and axes
func RenderLineGraph(data LineGraphData, x, y int, width, height int, designTokens *design.DesignTokens) string {
	var b strings.Builder

	if len(data.Points) == 0 {
		return ""
	}

	// Find min/max values for domain
	minValue := data.Points[0].Value
	maxValue := data.Points[0].Value
	minTime := data.Points[0].Date
	maxTime := data.Points[0].Date
	for _, point := range data.Points {
		if point.Value < minValue {
			minValue = point.Value
		}
		if point.Value > maxValue {
			maxValue = point.Value
		}
		if point.Date.Before(minTime) {
			minTime = point.Date
		}
		if point.Date.After(maxTime) {
			maxTime = point.Date
		}
	}

	// Reserve space for Y-axis labels
	labelAreaWidth := 2 * designTokens.Layout.CardPaddingRight
	plotWidth := width - labelAreaWidth

	// Create TimeScale for X-axis (dates to positions)
	xScale := scales.NewTimeScale(
		[2]time.Time{minTime, maxTime},
		[2]units.Length{units.Px(0), units.Px(float64(plotWidth))},
	)

	// Create LinearScale for Y-axis (values to positions)
	// Domain: [minValue, maxValue], Range: [height, 0] (inverted for SVG)
	yScale := scales.NewLinearScale(
		[2]float64{float64(minValue), float64(maxValue)},
		[2]units.Length{units.Px(float64(height)), units.Px(0)},
	).Nice(5) // Nice rounding and add padding

	b.WriteString(fmt.Sprintf(`<g transform="translate(%d, %d)">`, x, y))

	// Add gradient definition if requested
	var fillValue string
	if data.UseGradient && data.FillColor != "" {
		gradientID := data.GradientID
		if gradientID == "" {
			gradientID = fmt.Sprintf("lineGraphGradient-%d", atomic.AddInt64(&gradientCounter, 1))
		}

		// Create a vertical gradient from color to transparent
		gradient := svg.SimpleLinearGradient(gradientID, data.FillColor, "rgba(0,0,0,0)", 90)
		b.WriteString("<defs>")
		b.WriteString(gradient)
		b.WriteString("</defs>")

		fillValue = svg.GradientURL(gradientID)
	} else {
		fillValue = data.FillColor
	}

	// Create Y-axis with grid lines
	yAxis := axes.NewAxis(yScale, axes.AxisOrientationRight).
		TickCount(5).
		Grid(units.Px(float64(width)))

	// Render Y-axis with custom styling matching design tokens
	axisOpts := axes.DefaultRenderOptions()
	axisOpts.Style.StrokeColor = "rgba(255,255,255,0.1)"
	axisOpts.Style.GridStrokeColor = "rgba(255,255,255,0.1)"
	axisOpts.Style.TextColor = designTokens.Color
	axisOpts.Style.FontFamily = "monospace"
	axisOpts.Style.FontSize = 10
	axisOpts.Position = units.Px(float64(width - 5))

	b.WriteString(yAxis.Render(axisOpts))

	// Calculate scaled points using scales
	scaledPoints := make([]svg.Point, len(data.Points))
	for i, point := range data.Points {
		scaledPoints[i] = svg.Point{
			X: xScale.Apply(point.Date).Value,
			Y: yScale.Apply(float64(point.Value)).Value,
		}
	}

	// Draw filled area (if fill color specified)
	if data.FillColor != "" && len(data.Points) > 1 {
		var areaPath string
		if data.Smooth {
			tension := data.Tension
			if tension == 0 {
				tension = 0.3 // Default tension
			}
			areaPath = svg.SmoothAreaPath(scaledPoints, float64(height), tension)
		} else {
			areaPath = svg.AreaPath(scaledPoints, float64(height))
		}

		pathStyle := svg.Style{
			Fill: fillValue,
		}
		// Only apply opacity if not using gradient (gradient handles its own transparency)
		if !data.UseGradient {
			pathStyle.FillOpacity = 0.2
		}
		b.WriteString(svg.Path(areaPath, pathStyle))
		b.WriteString("\n")
	}

	// Draw line using PathBuilder
	if len(data.Points) > 1 {
		var linePath string
		if data.Smooth {
			tension := data.Tension
			if tension == 0 {
				tension = 0.3 // Default tension
			}
			linePath = svg.SmoothLinePath(scaledPoints, tension)
		} else {
			linePath = svg.PolylinePath(scaledPoints)
		}

		pathStyle := svg.Style{
			Fill:           "none",
			Stroke:         data.Color,
			StrokeWidth:    2,
			StrokeLinecap:  svg.StrokeLinecapRound,
			StrokeLinejoin: svg.StrokeLinejoinRound,
		}
		b.WriteString(svg.Path(linePath, pathStyle))
		b.WriteString("\n")
	}

	// Draw points with custom markers if specified
	if data.MarkerType != "" {
		markerSize := data.MarkerSize
		if markerSize == 0 {
			markerSize = 3 // Default size
		}

		for _, point := range scaledPoints {
			markerStyle := svg.Style{
				Fill:        data.Color,
				Stroke:      designTokens.Background,
				StrokeWidth: 1,
			}

			switch data.MarkerType {
			case "circle", "dot":
				b.WriteString(svg.Circle(point.X, point.Y, markerSize, markerStyle))
			case "square":
				halfSize := markerSize
				b.WriteString(svg.Rect(point.X-halfSize, point.Y-halfSize, halfSize*2, halfSize*2, markerStyle))
			case "diamond":
				diamondPoints := []svg.Point{
					{X: point.X, Y: point.Y - markerSize},
					{X: point.X + markerSize, Y: point.Y},
					{X: point.X, Y: point.Y + markerSize},
					{X: point.X - markerSize, Y: point.Y},
				}
				b.WriteString(svg.Polygon(diamondPoints, markerStyle))
			case "triangle":
				trianglePoints := []svg.Point{
					{X: point.X, Y: point.Y - markerSize},
					{X: point.X + markerSize, Y: point.Y + markerSize},
					{X: point.X - markerSize, Y: point.Y + markerSize},
				}
				b.WriteString(svg.Polygon(trianglePoints, markerStyle))
			default:
				// Default to circle
				b.WriteString(svg.Circle(point.X, point.Y, markerSize, markerStyle))
			}
			b.WriteString("\n")
		}
	}

	b.WriteString(`</g>`)

	// Add legend if label is provided
	if data.Label != "" {
		lineColor, err := color.HexToRGB(data.Color)
		if err != nil {
			lineColor, _ = color.HexToRGB("#000000")
		}

		var symbol legends.Symbol
		if data.MarkerType != "" {
			// Line with marker
			symbol = legends.NewLineSample(lineColor, 2, 25).WithMarker(data.MarkerType, 6)
		} else {
			// Plain line
			symbol = legends.NewLineSample(lineColor, 2, 25)
		}

		items := []legends.LegendItem{
			legends.Item(data.Label, symbol),
		}

		legend := legends.New(items,
			legends.WithPosition(legends.PositionTopRight),
			legends.WithLayout(legends.LayoutVertical),
		)

		b.WriteString(legend.Render(width, height))
	}

	return b.String()
}
