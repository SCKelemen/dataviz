package charts

import (
	"fmt"
	"strings"
	"time"

	"github.com/SCKelemen/color"
	design "github.com/SCKelemen/design-system"
	"github.com/SCKelemen/dataviz/axes"
	"github.com/SCKelemen/dataviz/charts/legends"
	"github.com/SCKelemen/dataviz/scales"
	"github.com/SCKelemen/svg"
	"github.com/SCKelemen/units"
)

// RenderScatterPlot renders a scatter plot using scales and axes
func RenderScatterPlot(data ScatterPlotData, x, y int, width, height int, designTokens *design.DesignTokens) string {
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

	markerSize := data.MarkerSize
	if markerSize == 0 {
		markerSize = 5 // Default size for scatter plots (larger than line graph)
	}

	markerType := data.MarkerType
	if markerType == "" {
		markerType = "circle" // Default to circle
	}

	// Draw each point using scales
	for _, point := range data.Points {
		pointX := xScale.Apply(point.Date).Value
		pointY := yScale.Apply(float64(point.Value)).Value

		// Use custom size if specified for this point
		size := markerSize
		if point.Size > 0 {
			size = point.Size
		}

		markerStyle := svg.Style{
			Fill:        data.Color,
			Stroke:      designTokens.Background,
			StrokeWidth: 1.5,
		}

		switch markerType {
		case "circle", "dot":
			b.WriteString(svg.Circle(pointX, pointY, size, markerStyle))
		case "square":
			halfSize := size
			b.WriteString(svg.Rect(pointX-halfSize, pointY-halfSize, halfSize*2, halfSize*2, markerStyle))
		case "diamond":
			diamondPoints := []svg.Point{
				{X: pointX, Y: pointY - size},
				{X: pointX + size, Y: pointY},
				{X: pointX, Y: pointY + size},
				{X: pointX - size, Y: pointY},
			}
			b.WriteString(svg.Polygon(diamondPoints, markerStyle))
		case "triangle":
			trianglePoints := []svg.Point{
				{X: pointX, Y: pointY - size},
				{X: pointX + size, Y: pointY + size},
				{X: pointX - size, Y: pointY + size},
			}
			b.WriteString(svg.Polygon(trianglePoints, markerStyle))
		case "cross":
			crossStyle := svg.Style{
				Stroke:      data.Color,
				StrokeWidth: 2,
				StrokeLinecap: svg.StrokeLinecapRound,
			}
			b.WriteString(svg.Line(pointX, pointY-size, pointX, pointY+size, crossStyle))
			b.WriteString(svg.Line(pointX-size, pointY, pointX+size, pointY, crossStyle))
		case "x":
			xStyle := svg.Style{
				Stroke:      data.Color,
				StrokeWidth: 2,
				StrokeLinecap: svg.StrokeLinecapRound,
			}
			b.WriteString(svg.Line(pointX-size*0.7, pointY-size*0.7, pointX+size*0.7, pointY+size*0.7, xStyle))
			b.WriteString(svg.Line(pointX-size*0.7, pointY+size*0.7, pointX+size*0.7, pointY-size*0.7, xStyle))
		default:
			// Default to circle
			b.WriteString(svg.Circle(pointX, pointY, size, markerStyle))
		}
		b.WriteString("\n")

		// Draw point label if specified
		if point.Label != "" {
			labelStyle := svg.Style{
				Fill:             designTokens.Color,
				Class:            "mono smaller",
				TextAnchor:       svg.TextAnchorMiddle,
				DominantBaseline: svg.DominantBaselineHanging,
			}
			b.WriteString(svg.Text(point.Label, pointX, pointY+size+3, labelStyle))
			b.WriteString("\n")
		}
	}

	b.WriteString(`</g>`)

	// Add legend if label is provided
	if data.Label != "" {
		plotColor, err := color.HexToRGB(data.Color)
		if err != nil {
			plotColor, _ = color.HexToRGB("#000000")
		}

		// Use marker symbol for scatter plot legend
		symbol := legends.NewMarkerSymbol(markerType, plotColor, 10)

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
