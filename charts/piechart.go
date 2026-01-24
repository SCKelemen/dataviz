package charts

import (
	"fmt"
	"math"
	"strings"

	"github.com/SCKelemen/color"
	"github.com/SCKelemen/dataviz/charts/legends"
	"github.com/SCKelemen/dataviz/scales"
	"github.com/SCKelemen/svg"
	"github.com/SCKelemen/units"
)

// Default color palette for pie charts
var defaultPieColors = []string{
	"#FF6B6B", "#4ECDC4", "#45B7D1", "#FFA07A",
	"#98D8C8", "#F7DC6F", "#BB8FCE", "#85C1E2",
	"#F8B739", "#52B788",
}

// RenderPieChart generates an SVG pie or donut chart
// x, y: position offset (usually 0, 0 for standalone chart)
// width, height: dimensions of the chart
// title: optional chart title
// donutMode: if true, renders as donut with center hole
// showLegend: if true, shows legend with labels
// showPercent: if true, shows percentage labels on slices
func RenderPieChart(data PieChartData, x, y int, width, height int, title string, donutMode, showLegend, showPercent bool) string {
	// Calculate total for percentages
	total := 0.0
	for _, slice := range data.Slices {
		total += slice.Value
	}

	if total == 0 {
		return renderEmptyPieChart(width, height, title)
	}

	// Parse color palette
	colorPalette := data.Colors
	if len(colorPalette) == 0 {
		colorPalette = defaultPieColors
	}

	// Convert hex colors to color.Color
	parsedColors := make([]color.Color, len(colorPalette))
	for i, hexColor := range colorPalette {
		c, err := color.HexToRGB(hexColor)
		if err != nil {
			c, _ = color.HexToRGB("#888888") // Fallback gray
		}
		parsedColors[i] = c
	}

	// Extract category names from slices
	categories := make([]string, len(data.Slices))
	for i, slice := range data.Slices {
		categories[i] = slice.Label
	}

	// Create categorical color scale
	colorScale := scales.NewCategoricalColorScale(categories, parsedColors)

	// Chart dimensions with space for title
	titleHeight := 40

	centerX := float64(x) + float64(width)/2
	centerY := float64(y) + float64(height)/2 + float64(titleHeight)/2
	radius := math.Min(float64(width), float64(height-titleHeight))/2 - 80

	// Calculate inner radius for donut mode
	innerRadius := 0.0
	if donutMode {
		innerRadius = radius * 0.4 // 40% hole
	}

	// Start building SVG content
	var sb strings.Builder

	// Title
	if title != "" {
		titleStyle := svg.Style{
			Fill:             "#333",
			FontFamily:       "Arial, sans-serif",
			FontSize:         units.Px(16),
			FontWeight:       "bold",
			TextAnchor:       svg.TextAnchorMiddle,
		}
		sb.WriteString(svg.Text(title, float64(width)/2, 25, titleStyle))
	}

	// Draw slices
	startAngle := -math.Pi / 2 // Start at top
	for _, slice := range data.Slices {
		angle := (slice.Value / total) * 2 * math.Pi
		endAngle := startAngle + angle

		// Get color from scale
		sliceColor := colorScale.ApplyColor(slice.Label)
		colorHex := color.RGBToHex(sliceColor)

		// Draw slice
		sb.WriteString(renderPieSlice(centerX, centerY, radius, innerRadius, startAngle, endAngle, colorHex))

		// Draw percentage label if enabled
		if showPercent {
			percentage := (slice.Value / total) * 100
			if percentage >= 5 { // Only show label if slice is >= 5%
				midAngle := startAngle + angle/2
				labelRadius := radius * 0.7
				if donutMode {
					labelRadius = innerRadius + (radius-innerRadius)*0.5
				}
				labelX := centerX + labelRadius*math.Cos(midAngle)
				labelY := centerY + labelRadius*math.Sin(midAngle)

				labelStyle := svg.Style{
					Fill:             "#FFFFFF",
					FontFamily:       "Arial, sans-serif",
					FontSize:         units.Px(12),
					FontWeight:       "bold",
					TextAnchor:       svg.TextAnchorMiddle,
					DominantBaseline: svg.DominantBaselineMiddle,
				}
				sb.WriteString(svg.Text(fmt.Sprintf("%.1f%%", percentage), labelX, labelY, labelStyle))
			}
		}

		startAngle = endAngle
	}

	// Draw legend if enabled
	if showLegend {
		// Create legend items using color scale
		items := make([]legends.LegendItem, len(data.Slices))
		for i, slice := range data.Slices {
			percentage := (slice.Value / total) * 100
			sliceColor := colorScale.ApplyColor(slice.Label)
			items[i] = legends.ItemWithValue(
				slice.Label,
				legends.Swatch(sliceColor),
				fmt.Sprintf("%.1f%%", percentage),
			)
		}

		// Create and render legend
		legend := legends.New(items,
			legends.WithPosition(legends.PositionBottomLeft),
			legends.WithLayout(legends.LayoutVertical),
		)

		sb.WriteString(legend.Render(width, height))
	}

	return sb.String()
}

// renderPieSlice generates an SVG path for a single pie slice
func renderPieSlice(cx, cy, outerRadius, innerRadius, startAngle, endAngle float64, color string) string {
	// Calculate outer arc points
	x1 := cx + outerRadius*math.Cos(startAngle)
	y1 := cy + outerRadius*math.Sin(startAngle)
	x2 := cx + outerRadius*math.Cos(endAngle)
	y2 := cy + outerRadius*math.Sin(endAngle)

	// Determine if arc should be large (> 180 degrees)
	largeArc := 0
	if endAngle-startAngle > math.Pi {
		largeArc = 1
	}

	var pathData string
	if innerRadius > 0 {
		// Donut mode: draw annular sector
		x3 := cx + innerRadius*math.Cos(endAngle)
		y3 := cy + innerRadius*math.Sin(endAngle)
		x4 := cx + innerRadius*math.Cos(startAngle)
		y4 := cy + innerRadius*math.Sin(startAngle)

		pathData = fmt.Sprintf("M %.2f %.2f A %.2f %.2f 0 %d 1 %.2f %.2f L %.2f %.2f A %.2f %.2f 0 %d 0 %.2f %.2f Z",
			x1, y1, outerRadius, outerRadius, largeArc, x2, y2,
			x3, y3, innerRadius, innerRadius, largeArc, x4, y4)
	} else {
		// Pie mode: draw sector from center
		pathData = fmt.Sprintf("M %.2f %.2f L %.2f %.2f A %.2f %.2f 0 %d 1 %.2f %.2f Z",
			cx, cy, x1, y1, outerRadius, outerRadius, largeArc, x2, y2)
	}

	style := svg.Style{
		Fill:        color,
		Stroke:      "#FFFFFF",
		StrokeWidth: 2,
	}

	return svg.Path(pathData, style)
}

// renderEmptyPieChart generates SVG content for when there's no data
func renderEmptyPieChart(width, height int, title string) string {
	var sb strings.Builder

	if title != "" {
		titleStyle := svg.Style{
			Fill:       "#333",
			FontFamily: "Arial, sans-serif",
			FontSize:   units.Px(16),
			FontWeight: "bold",
			TextAnchor: svg.TextAnchorMiddle,
		}
		sb.WriteString(svg.Text(title, float64(width)/2, 25, titleStyle))
	}

	emptyStyle := svg.Style{
		Fill:       "#999",
		FontFamily: "Arial, sans-serif",
		FontSize:   units.Px(14),
		TextAnchor: svg.TextAnchorMiddle,
	}
	sb.WriteString(svg.Text("No data available", float64(width)/2, float64(height)/2, emptyStyle))
	return sb.String()
}
