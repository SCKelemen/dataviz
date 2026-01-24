package charts

import (
	"fmt"
	"math"
	"strings"
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

	// Use default colors if not provided
	colors := data.Colors
	if len(colors) == 0 {
		colors = defaultPieColors
	}

	// Chart dimensions with space for title and legend
	titleHeight := 40
	legendHeight := 0
	if showLegend {
		legendHeight = len(data.Slices)*25 + 20
	}

	centerX := float64(x) + float64(width)/2
	centerY := float64(y) + float64(height-legendHeight)/2 + float64(titleHeight)
	radius := math.Min(float64(width), float64(height-titleHeight-legendHeight))/2 - 60

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
	for i, slice := range data.Slices {
		angle := (slice.Value / total) * 2 * math.Pi
		endAngle := startAngle + angle

		// Get color (cycle through palette)
		color := colors[i%len(colors)]

		// Draw slice
		sb.WriteString(renderPieSlice(centerX, centerY, radius, innerRadius, startAngle, endAngle, color))

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
		legendY := height - legendHeight + 10
		for i, slice := range data.Slices {
			color := colors[i%len(colors)]
			y := legendY + i*25

			// Color box
			sb.WriteString(fmt.Sprintf(`<rect x="20" y="%d" width="15" height="15" fill="%s"/>`,
				y, color))

			// Label
			percentage := (slice.Value / total) * 100
			label := fmt.Sprintf("%s (%.1f%%)", slice.Label, percentage)
			sb.WriteString(fmt.Sprintf(`<text x="40" y="%d" font-family="Arial, sans-serif" font-size="12" fill="#333">%s</text>`,
				y+12, label))
		}
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
