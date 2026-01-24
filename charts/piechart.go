package charts

import (
	"fmt"
	"math"
	"strings"

	"github.com/SCKelemen/color"
	"github.com/SCKelemen/dataviz/charts/legends"
	"github.com/SCKelemen/dataviz/scales"
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

	// Start building SVG
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d">`,
		width, height))

	// Background
	sb.WriteString(fmt.Sprintf(`<rect width="%d" height="%d" fill="#FFFFFF"/>`, width, height))

	// Title
	if title != "" {
		sb.WriteString(fmt.Sprintf(`<text x="%d" y="25" font-family="Arial, sans-serif" font-size="16" font-weight="bold" fill="#333" text-anchor="middle">%s</text>`,
			width/2, title))
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

				sb.WriteString(fmt.Sprintf(`<text x="%.2f" y="%.2f" font-family="Arial, sans-serif" font-size="12" fill="#FFFFFF" text-anchor="middle" dominant-baseline="middle" font-weight="bold">%.1f%%</text>`,
					labelX, labelY, percentage))
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

	sb.WriteString(`</svg>`)
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

	var path string
	if innerRadius > 0 {
		// Donut mode: draw annular sector
		x3 := cx + innerRadius*math.Cos(endAngle)
		y3 := cy + innerRadius*math.Sin(endAngle)
		x4 := cx + innerRadius*math.Cos(startAngle)
		y4 := cy + innerRadius*math.Sin(startAngle)

		path = fmt.Sprintf(`<path d="M %.2f %.2f A %.2f %.2f 0 %d 1 %.2f %.2f L %.2f %.2f A %.2f %.2f 0 %d 0 %.2f %.2f Z" fill="%s" stroke="#FFFFFF" stroke-width="2"/>`,
			x1, y1, outerRadius, outerRadius, largeArc, x2, y2,
			x3, y3, innerRadius, innerRadius, largeArc, x4, y4, color)
	} else {
		// Pie mode: draw sector from center
		path = fmt.Sprintf(`<path d="M %.2f %.2f L %.2f %.2f A %.2f %.2f 0 %d 1 %.2f %.2f Z" fill="%s" stroke="#FFFFFF" stroke-width="2"/>`,
			cx, cy, x1, y1, outerRadius, outerRadius, largeArc, x2, y2, color)
	}

	return path
}

// renderEmptyPieChart generates an SVG for when there's no data
func renderEmptyPieChart(width, height int, title string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d">`,
		width, height))
	sb.WriteString(fmt.Sprintf(`<rect width="%d" height="%d" fill="#FFFFFF"/>`, width, height))

	if title != "" {
		sb.WriteString(fmt.Sprintf(`<text x="%d" y="25" font-family="Arial, sans-serif" font-size="16" font-weight="bold" fill="#333" text-anchor="middle">%s</text>`,
			width/2, title))
	}

	sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" font-family="Arial, sans-serif" font-size="14" fill="#999" text-anchor="middle">No data available</text>`,
		width/2, height/2))
	sb.WriteString(`</svg>`)
	return sb.String()
}
