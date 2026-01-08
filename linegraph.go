package dataviz

import (
	"fmt"
	"strings"
	"sync/atomic"

	design "github.com/SCKelemen/design-system"
	"github.com/SCKelemen/svg"
)

// Global counter for unique gradient IDs
var gradientCounter int64

// RenderLineGraph renders a line graph
func RenderLineGraph(data LineGraphData, x, y int, width, height int, designTokens *design.DesignTokens) string {
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

	// Calculate scaled points
	pointWidth := float64(plotWidth) / float64(len(data.Points)-1)
	if len(data.Points) == 1 {
		pointWidth = 0
	}

	scaledPoints := make([]svg.Point, len(data.Points))
	for i, point := range data.Points {
		scaledPoints[i] = svg.Point{
			X: float64(i) * pointWidth,
			Y: float64(height) - (float64(point.Value-minValue)/float64(valueRange))*float64(height),
		}
	}

	// Draw filled area (if fill color specified)
	if data.FillColor != "" && len(data.Points) > 1 {
		areaPath := svg.AreaPath(scaledPoints, float64(height))

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
		linePath := svg.PolylinePath(scaledPoints)

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

	// Draw points
	for _, point := range scaledPoints {
		circleStyle := svg.Style{
			Fill:        data.Color,
			Stroke:      designTokens.Background,
			StrokeWidth: 1,
		}
		b.WriteString(svg.Circle(point.X, point.Y, 3, circleStyle))
		b.WriteString("\n")
	}

	b.WriteString(`</g>`)
	return b.String()
}
