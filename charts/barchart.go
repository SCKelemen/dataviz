package charts

import (
	"fmt"
	"strings"

	"github.com/SCKelemen/color"
	design "github.com/SCKelemen/design-system"
	"github.com/SCKelemen/svg"
)

// RenderBarChart renders a bar chart
func RenderBarChart(data BarChartData, x, y int, width, height int, designTokens *design.DesignTokens) string {
	var b strings.Builder

	if len(data.Bars) == 0 {
		return ""
	}

	// Find max value
	maxValue := 0
	for _, bar := range data.Bars {
		total := bar.Value + bar.Secondary
		if total > maxValue {
			maxValue = total
		}
	}
	if maxValue == 0 {
		maxValue = 1
	}

	// Bars are 7.7px wide with spacing
	barWidth := 7.7
	barSpacing := float64(width)/float64(len(data.Bars)) - barWidth

	b.WriteString(fmt.Sprintf(`<g transform="translate(%d, %d)">`, x, y))

	// Calculate base Y position (bars grow upward from bottom)
	baseY := float64(height)

	for i, bar := range data.Bars {
		barX := float64(i)*(barWidth+barSpacing) + 1.0 // +1 for offset

		if data.Stacked {
			// Stacked bars: opened (lighter) on bottom, closed (darker) on top
			lighterColor := data.Color
			if c, err := color.ParseColor(data.Color); err == nil {
				// Lighten by 30% for better visual distinction
				lightened := color.Lighten(c, 0.3)
				lighterColor = color.RGBToHex(lightened)
			}

			// Calculate heights scaled to maxValue
			primaryHeight := (float64(bar.Value) / float64(maxValue)) * float64(height)
			secondaryHeight := (float64(bar.Secondary) / float64(maxValue)) * float64(height)
			totalHeight := primaryHeight + secondaryHeight

			// Scale down if total height exceeds container
			if totalHeight > float64(height) {
				scale := float64(height) / totalHeight
				primaryHeight *= scale
				secondaryHeight *= scale
			}

			// Primary bar (opened) - lighter color, on bottom
			primaryY := baseY - primaryHeight
			if primaryY < 0 {
				primaryY = 0
			}
			primaryStyle := svg.Style{Fill: lighterColor}
			b.WriteString(svg.Rect(barX, primaryY, barWidth, primaryHeight, primaryStyle))
			b.WriteString("\n")

			// Secondary bar (closed) - darker color, stacked on top
			if bar.Secondary > 0 {
				secondaryY := primaryY - secondaryHeight
				if secondaryY < 0 {
					secondaryY = 0
					secondaryHeight = primaryY
				}
				secondaryStyle := svg.Style{Fill: data.Color}
				b.WriteString(svg.Rect(barX, secondaryY, barWidth, secondaryHeight, secondaryStyle))
				b.WriteString("\n")
			}
		} else {
			// Single bar
			barHeight := (float64(bar.Value) / float64(maxValue)) * float64(height)
			if barHeight > float64(height) {
				barHeight = float64(height)
			}
			barY := baseY - barHeight
			if barY < 0 {
				barY = 0
				barHeight = baseY
			}
			barStyle := svg.Style{Fill: data.Color}
			b.WriteString(svg.Rect(barX, barY, barWidth, barHeight, barStyle))
			b.WriteString("\n")
		}
	}

	b.WriteString(`</g>`)
	return b.String()
}
