package charts

import (
	"fmt"
	"strings"

	"github.com/SCKelemen/color"
	design "github.com/SCKelemen/design-system"
	"github.com/SCKelemen/dataviz/charts/legends"
	"github.com/SCKelemen/dataviz/scales"
	"github.com/SCKelemen/svg"
	"github.com/SCKelemen/units"
)

// RenderBarChart renders a bar chart using scales and axes
func RenderBarChart(data BarChartData, x, y int, width, height int, designTokens *design.DesignTokens) string {
	var b strings.Builder

	if len(data.Bars) == 0 {
		return ""
	}

	// Find max value for domain
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

	// Create category labels for X-axis
	categories := make([]string, len(data.Bars))
	for i, bar := range data.Bars {
		if bar.Label != "" {
			categories[i] = bar.Label
		} else {
			categories[i] = fmt.Sprintf("%d", i)
		}
	}

	// Create BandScale for X-axis (categorical positioning)
	xScale := scales.NewBandScale(
		categories,
		[2]units.Length{units.Px(0), units.Px(float64(width))},
	).Padding(0.2) // 20% padding between bars

	// Create LinearScale for Y-axis (value to height)
	// Domain: [0, maxValue], Range: [height, 0] (inverted for SVG coordinates)
	yScale := scales.NewLinearScale(
		[2]float64{0, float64(maxValue)},
		[2]units.Length{units.Px(float64(height)), units.Px(0)},
	).Nice(5) // Nice rounding for axis ticks

	b.WriteString(fmt.Sprintf(`<g transform="translate(%d, %d)">`, x, y))

	// Render bars using scales
	bandwidth := xScale.Bandwidth()
	for i, bar := range data.Bars {
		// Get X position from BandScale
		barX := xScale.Apply(categories[i]).Value
		barWidth := bandwidth.Value

		if data.Stacked {
			// Stacked bars: opened (lighter) on bottom, closed (darker) on top
			lighterColor := data.Color
			if c, err := color.ParseColor(data.Color); err == nil {
				// Lighten by 30% for better visual distinction
				lightened := color.Lighten(c, 0.3)
				lighterColor = color.RGBToHex(lightened)
			}

			// Use Y scale to map values to positions
			baseY := yScale.Apply(0.0).Value                               // Bottom (y=0 maps to height)
			primaryTop := yScale.Apply(float64(bar.Value)).Value           // Top of primary bar
			secondaryTop := yScale.Apply(float64(bar.Value + bar.Secondary)).Value // Top of secondary bar

			// Primary bar (opened) - lighter color, on bottom
			primaryHeight := baseY - primaryTop
			primaryStyle := svg.Style{Fill: lighterColor}
			b.WriteString(svg.Rect(barX, primaryTop, barWidth, primaryHeight, primaryStyle))
			b.WriteString("\n")

			// Secondary bar (closed) - darker color, stacked on top
			if bar.Secondary > 0 {
				secondaryHeight := primaryTop - secondaryTop
				secondaryStyle := svg.Style{Fill: data.Color}
				b.WriteString(svg.Rect(barX, secondaryTop, barWidth, secondaryHeight, secondaryStyle))
				b.WriteString("\n")
			}
		} else {
			// Single bar
			// Use Y scale to map value to position
			baseY := yScale.Apply(0.0).Value             // Bottom
			barTop := yScale.Apply(float64(bar.Value)).Value // Top of bar
			barHeight := baseY - barTop

			barStyle := svg.Style{Fill: data.Color}
			b.WriteString(svg.Rect(barX, barTop, barWidth, barHeight, barStyle))
			b.WriteString("\n")
		}
	}

	b.WriteString(`</g>`)

	// Add legend if label is provided
	if data.Label != "" {
		barColor, err := color.ParseColor(data.Color)
		if err != nil {
			barColor, _ = color.HexToRGB("#000000")
		}

		if data.Stacked {
			// For stacked bars, show two legend items
			lighterColor := color.Lighten(barColor, 0.3)

			items := []legends.LegendItem{
				legends.Item(data.Label+" (closed)", legends.Swatch(barColor)),
				legends.Item(data.Label+" (opened)", legends.Swatch(lighterColor)),
			}

			legend := legends.New(items,
				legends.WithPosition(legends.PositionTopRight),
				legends.WithLayout(legends.LayoutVertical),
			)

			b.WriteString(legend.Render(width, height))
		} else {
			// Single bar legend
			items := []legends.LegendItem{
				legends.Item(data.Label, legends.Swatch(barColor)),
			}

			legend := legends.New(items,
				legends.WithPosition(legends.PositionTopRight),
				legends.WithLayout(legends.LayoutVertical),
			)

			b.WriteString(legend.Render(width, height))
		}
	}

	return b.String()
}
