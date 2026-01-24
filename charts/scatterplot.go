package charts

import (
	"fmt"
	"strings"

	"github.com/SCKelemen/color"
	design "github.com/SCKelemen/design-system"
	"github.com/SCKelemen/dataviz/charts/legends"
	"github.com/SCKelemen/svg"
)

// RenderScatterPlot renders a scatter plot
func RenderScatterPlot(data ScatterPlotData, x, y int, width, height int, designTokens *design.DesignTokens) string {
	var b strings.Builder

	if len(data.Points) == 0 {
		return ""
	}

	// Find min/max values for scaling
	minValue := data.Points[0].Value
	maxValue := data.Points[0].Value
	for _, point := range data.Points {
		if point.Value < minValue {
			minValue = point.Value
		}
		if point.Value > maxValue {
			maxValue = point.Value
		}
	}

	// Add some padding
	valueRange := maxValue - minValue
	if valueRange == 0 {
		valueRange = 1
	}
	padding := float64(valueRange) * 0.1
	minValue -= int(padding)
	maxValue += int(padding)
	valueRange = maxValue - minValue

	b.WriteString(fmt.Sprintf(`<g transform="translate(%d, %d)">`, x, y))

	// Reserve space for graduations and labels
	labelAreaWidth := 2 * designTokens.Layout.CardPaddingRight
	plotWidth := width - labelAreaWidth

	// Draw grid lines
	gridLines := 5
	for i := 0; i <= gridLines; i++ {
		gridY := float64(height) * float64(i) / float64(gridLines)
		value := minValue + int(float64(valueRange)*float64(gridLines-i)/float64(gridLines))

		lineStyle := svg.Style{
			Stroke:      "rgba(255,255,255,0.1)",
			StrokeWidth: 1,
		}
		b.WriteString(svg.Line(0, gridY, float64(width), gridY, lineStyle))
		b.WriteString("\n")

		textStyle := svg.Style{
			Fill:             designTokens.Color,
			Class:            "mono smaller",
			Opacity:          0.5,
			TextAnchor:       svg.TextAnchorEnd,
			DominantBaseline: svg.DominantBaselineMiddle,
		}
		b.WriteString(svg.Text(fmt.Sprintf("%d", value), float64(width-5), gridY, textStyle))
		b.WriteString("\n")
	}

	// Calculate point positions
	pointWidth := float64(plotWidth) / float64(len(data.Points)-1)
	if len(data.Points) == 1 {
		pointWidth = float64(plotWidth) / 2 // Center single point
	}

	markerSize := data.MarkerSize
	if markerSize == 0 {
		markerSize = 5 // Default size for scatter plots (larger than line graph)
	}

	markerType := data.MarkerType
	if markerType == "" {
		markerType = "circle" // Default to circle
	}

	// Draw each point
	for i, point := range data.Points {
		pointX := float64(i) * pointWidth
		pointY := float64(height) - (float64(point.Value-minValue)/float64(valueRange))*float64(height)

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
