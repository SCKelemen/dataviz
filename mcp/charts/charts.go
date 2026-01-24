package charts

import (
	"fmt"
	"math"
	"strings"

	"github.com/SCKelemen/color"
	"github.com/SCKelemen/dataviz/mcp/types"
	"github.com/SCKelemen/layout"
	"github.com/SCKelemen/svg"
	"github.com/SCKelemen/units"
)

// CreateBarChart generates a bar chart SVG using SCKelemen libraries
func CreateBarChart(config types.BarChartConfig) (string, error) {
	// Calculate data ranges
	maxValue := 0.0
	for _, dp := range config.Data {
		if dp.Value > maxValue {
			maxValue = dp.Value
		}
	}

	// Chart dimensions and margins
	margin := 60.0
	chartWidth := float64(config.Width) - (2 * margin)
	chartHeight := float64(config.Height) - (2 * margin)

	// Calculate bar dimensions
	barCount := len(config.Data)
	barSpacing := 10.0
	totalSpacing := barSpacing * float64(barCount+1)
	barWidth := (chartWidth - totalSpacing) / float64(barCount)

	// Start building SVG
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`,
		config.Width, config.Height, config.Width, config.Height))
	sb.WriteString("\n")

	// Background
	sb.WriteString(fmt.Sprintf(`  <rect width="%d" height="%d" fill="#ffffff"/>`, config.Width, config.Height))
	sb.WriteString("\n")

	// Title
	if config.Title != "" {
		sb.WriteString(fmt.Sprintf(`  <text x="%d" y="30" text-anchor="middle" font-size="20" font-weight="bold" fill="#1f2937">%s</text>`,
			config.Width/2, config.Title))
		sb.WriteString("\n")
	}

	// Determine bar color
	barColor := config.Color
	if barColor == "" {
		barColor = "#3b82f6" // Default blue
	}

	// Draw bars
	x := margin + barSpacing
	for _, dp := range config.Data {
		// Calculate bar height based on value
		barHeight := (dp.Value / maxValue) * chartHeight
		y := margin + chartHeight - barHeight

		// Draw bar with rounded corners
		sb.WriteString(fmt.Sprintf(`  <rect x="%.2f" y="%.2f" width="%.2f" height="%.2f" fill="%s" rx="4"/>`,
			x, y, barWidth, barHeight, barColor))
		sb.WriteString("\n")

		// Draw label
		labelX := x + (barWidth / 2)
		labelY := margin + chartHeight + 20
		sb.WriteString(fmt.Sprintf(`  <text x="%.2f" y="%.2f" text-anchor="middle" font-size="12" fill="#6b7280">%s</text>`,
			labelX, labelY, dp.Label))
		sb.WriteString("\n")

		// Draw value on top of bar
		valueY := y - 5
		if valueY < margin {
			valueY = y + 15
		}
		sb.WriteString(fmt.Sprintf(`  <text x="%.2f" y="%.2f" text-anchor="middle" font-size="11" fill="#374151">%.1f</text>`,
			labelX, valueY, dp.Value))
		sb.WriteString("\n")

		x += barWidth + barSpacing
	}

	// Draw Y-axis
	sb.WriteString(fmt.Sprintf(`  <line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="#d1d5db" stroke-width="2"/>`,
		margin, margin, margin, margin+chartHeight))
	sb.WriteString("\n")

	// Draw X-axis
	sb.WriteString(fmt.Sprintf(`  <line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="#d1d5db" stroke-width="2"/>`,
		margin, margin+chartHeight, margin+chartWidth, margin+chartHeight))
	sb.WriteString("\n")

	// Draw Y-axis scale
	steps := 5
	for i := 0; i <= steps; i++ {
		value := (maxValue / float64(steps)) * float64(i)
		y := margin + chartHeight - (chartHeight/float64(steps))*float64(i)

		// Grid line
		sb.WriteString(fmt.Sprintf(`  <line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="#e5e7eb" stroke-width="1" stroke-dasharray="4,4"/>`,
			margin, y, margin+chartWidth, y))
		sb.WriteString("\n")

		// Scale label
		sb.WriteString(fmt.Sprintf(`  <text x="%.2f" y="%.2f" text-anchor="end" font-size="11" fill="#6b7280">%.1f</text>`,
			margin-10, y+4, value))
		sb.WriteString("\n")
	}

	sb.WriteString("</svg>")

	return sb.String(), nil
}

// CreatePieChart generates a pie chart SVG using SCKelemen libraries
func CreatePieChart(config types.PieChartConfig) (string, error) {
	// Calculate total for percentages
	total := 0.0
	for _, dp := range config.Data {
		total += dp.Value
	}

	if total == 0 {
		return "", fmt.Errorf("total value is zero")
	}

	// Chart dimensions
	centerX := float64(config.Width) / 2
	centerY := float64(config.Height) / 2
	radius := math.Min(float64(config.Width), float64(config.Height))/2 - 60

	// Start building SVG
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d">`,
		config.Width, config.Height))
	sb.WriteString("\n")

	// Background
	sb.WriteString(fmt.Sprintf(`  <rect width="%d" height="%d" fill="#ffffff"/>`, config.Width, config.Height))
	sb.WriteString("\n")

	// Title
	if config.Title != "" {
		sb.WriteString(fmt.Sprintf(`  <text x="%d" y="30" text-anchor="middle" font-size="20" font-weight="bold" fill="#1f2937">%s</text>`,
			config.Width/2, config.Title))
		sb.WriteString("\n")
	}

	// Color palette
	colors := []string{
		"#3b82f6", "#10b981", "#f59e0b", "#ef4444", "#8b5cf6",
		"#ec4899", "#06b6d4", "#84cc16", "#f97316", "#6366f1",
	}

	// Draw slices
	startAngle := -90.0 // Start at top
	for i, dp := range config.Data {
		angle := (dp.Value / total) * 360.0
		endAngle := startAngle + angle

		// Calculate slice path
		startRad := startAngle * math.Pi / 180
		endRad := endAngle * math.Pi / 180

		x1 := centerX + radius*math.Cos(startRad)
		y1 := centerY + radius*math.Sin(startRad)
		x2 := centerX + radius*math.Cos(endRad)
		y2 := centerY + radius*math.Sin(endRad)

		largeArc := 0
		if angle > 180 {
			largeArc = 1
		}

		// Draw slice
		sliceColor := colors[i%len(colors)]
		sb.WriteString(fmt.Sprintf(`  <path d="M %.2f,%.2f L %.2f,%.2f A %.2f,%.2f 0 %d,1 %.2f,%.2f Z" fill="%s" stroke="#ffffff" stroke-width="2"/>`,
			centerX, centerY, x1, y1, radius, radius, largeArc, x2, y2, sliceColor))
		sb.WriteString("\n")

		// Draw label
		midAngle := (startAngle + endAngle) / 2
		midRad := midAngle * math.Pi / 180
		labelRadius := radius * 0.7
		labelX := centerX + labelRadius*math.Cos(midRad)
		labelY := centerY + labelRadius*math.Sin(midRad)

		percentage := (dp.Value / total) * 100
		sb.WriteString(fmt.Sprintf(`  <text x="%.2f" y="%.2f" text-anchor="middle" font-size="12" font-weight="bold" fill="#ffffff">%.1f%%</text>`,
			labelX, labelY, percentage))
		sb.WriteString("\n")

		startAngle = endAngle
	}

	// Draw legend
	legendX := float64(config.Width) - 150
	legendY := 60.0
	for i, dp := range config.Data {
		// Color box
		sb.WriteString(fmt.Sprintf(`  <rect x="%.2f" y="%.2f" width="12" height="12" fill="%s"/>`,
			legendX, legendY+float64(i)*20, colors[i%len(colors)]))
		sb.WriteString("\n")

		// Label
		sb.WriteString(fmt.Sprintf(`  <text x="%.2f" y="%.2f" font-size="12" fill="#374151">%s</text>`,
			legendX+20, legendY+float64(i)*20+10, dp.Label))
		sb.WriteString("\n")
	}

	sb.WriteString("</svg>")

	return sb.String(), nil
}

// CreateLineChart generates a line chart SVG using SCKelemen libraries
func CreateLineChart(config types.LineChartConfig) (string, error) {
	if len(config.Series) == 0 {
		return "", fmt.Errorf("no series data provided")
	}

	// Calculate data ranges
	minY, maxY := math.MaxFloat64, -math.MaxFloat64
	maxPoints := 0

	for _, series := range config.Series {
		for _, point := range series.Data {
			if point.Y < minY {
				minY = point.Y
			}
			if point.Y > maxY {
				maxY = point.Y
			}
		}
		if len(series.Data) > maxPoints {
			maxPoints = len(series.Data)
		}
	}

	// Chart dimensions
	margin := 60.0
	chartWidth := float64(config.Width) - (2 * margin)
	chartHeight := float64(config.Height) - (2 * margin)

	// Start building SVG
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d">`,
		config.Width, config.Height))
	sb.WriteString("\n")

	// Background
	sb.WriteString(fmt.Sprintf(`  <rect width="%d" height="%d" fill="#ffffff"/>`, config.Width, config.Height))
	sb.WriteString("\n")

	// Title
	if config.Title != "" {
		sb.WriteString(fmt.Sprintf(`  <text x="%d" y="30" text-anchor="middle" font-size="20" font-weight="bold" fill="#1f2937">%s</text>`,
			config.Width/2, config.Title))
		sb.WriteString("\n")
	}

	// Color palette
	colors := []string{"#3b82f6", "#10b981", "#f59e0b", "#ef4444", "#8b5cf6"}

	// Draw axes
	sb.WriteString(fmt.Sprintf(`  <line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="#d1d5db" stroke-width="2"/>`,
		margin, margin, margin, margin+chartHeight))
	sb.WriteString(fmt.Sprintf(`  <line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="#d1d5db" stroke-width="2"/>`,
		margin, margin+chartHeight, margin+chartWidth, margin+chartHeight))
	sb.WriteString("\n")

	// Draw grid and Y-axis labels
	steps := 5
	for i := 0; i <= steps; i++ {
		value := minY + ((maxY - minY) / float64(steps) * float64(i))
		y := margin + chartHeight - (chartHeight/float64(steps))*float64(i)

		sb.WriteString(fmt.Sprintf(`  <line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="#e5e7eb" stroke-width="1" stroke-dasharray="4,4"/>`,
			margin, y, margin+chartWidth, y))
		sb.WriteString(fmt.Sprintf(`  <text x="%.2f" y="%.2f" text-anchor="end" font-size="11" fill="#6b7280">%.1f</text>`,
			margin-10, y+4, value))
		sb.WriteString("\n")
	}

	// Draw series
	for seriesIdx, series := range config.Series {
		if len(series.Data) == 0 {
			continue
		}

		seriesColor := series.Color
		if seriesColor == "" {
			seriesColor = colors[seriesIdx%len(colors)]
		}

		// Build path using smooth curves
		points := make([]svg.Point, len(series.Data))
		for i, point := range series.Data {
			x := margin + (chartWidth / float64(maxPoints-1) * float64(i))
			y := margin + chartHeight - ((point.Y-minY)/(maxY-minY))*chartHeight
			points[i] = svg.Point{X: x, Y: y}
		}

		pathData := svg.SmoothLinePath(points, 0.3)

		sb.WriteString(fmt.Sprintf(`  <path d="%s" fill="none" stroke="%s" stroke-width="2"/>`,
			pathData, seriesColor))
		sb.WriteString("\n")

		// Draw points
		for _, p := range points {
			sb.WriteString(fmt.Sprintf(`  <circle cx="%.2f" cy="%.2f" r="4" fill="%s" stroke="#ffffff" stroke-width="2"/>`,
				p.X, p.Y, seriesColor))
			sb.WriteString("\n")
		}
	}

	// Draw legend
	legendX := margin
	legendY := 50.0
	for i, series := range config.Series {
		seriesColor := series.Color
		if seriesColor == "" {
			seriesColor = colors[i%len(colors)]
		}

		xOffset := float64(i * 120)
		sb.WriteString(fmt.Sprintf(`  <line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="%s" stroke-width="2"/>`,
			legendX+xOffset, legendY, legendX+xOffset+20, legendY, seriesColor))
		sb.WriteString(fmt.Sprintf(`  <text x="%.2f" y="%.2f" font-size="12" fill="#374151">%s</text>`,
			legendX+xOffset+25, legendY+4, series.Name))
		sb.WriteString("\n")
	}

	sb.WriteString("</svg>")

	return sb.String(), nil
}

// CreateScatterPlot generates a scatter plot SVG using SCKelemen libraries
func CreateScatterPlot(config types.ScatterPlotConfig) (string, error) {
	if len(config.Data) == 0 {
		return "", fmt.Errorf("no data provided")
	}

	// Calculate data ranges
	minX, maxX := math.MaxFloat64, -math.MaxFloat64
	minY, maxY := math.MaxFloat64, -math.MaxFloat64

	for _, point := range config.Data {
		if point.X < minX {
			minX = point.X
		}
		if point.X > maxX {
			maxX = point.X
		}
		if point.Y < minY {
			minY = point.Y
		}
		if point.Y > maxY {
			maxY = point.Y
		}
	}

	// Chart dimensions
	margin := 60.0
	chartWidth := float64(config.Width) - (2 * margin)
	chartHeight := float64(config.Height) - (2 * margin)

	// Start building SVG
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d">`,
		config.Width, config.Height))
	sb.WriteString("\n")

	// Background
	sb.WriteString(fmt.Sprintf(`  <rect width="%d" height="%d" fill="#ffffff"/>`, config.Width, config.Height))
	sb.WriteString("\n")

	// Title
	if config.Title != "" {
		sb.WriteString(fmt.Sprintf(`  <text x="%d" y="30" text-anchor="middle" font-size="20" font-weight="bold" fill="#1f2937">%s</text>`,
			config.Width/2, config.Title))
		sb.WriteString("\n")
	}

	// Draw axes
	sb.WriteString(fmt.Sprintf(`  <line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="#d1d5db" stroke-width="2"/>`,
		margin, margin, margin, margin+chartHeight))
	sb.WriteString(fmt.Sprintf(`  <line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="#d1d5db" stroke-width="2"/>`,
		margin, margin+chartHeight, margin+chartWidth, margin+chartHeight))
	sb.WriteString("\n")

	// Draw points
	for _, point := range config.Data {
		x := margin + ((point.X-minX)/(maxX-minX))*chartWidth
		y := margin + chartHeight - ((point.Y-minY)/(maxY-minY))*chartHeight

		radius := 5.0
		if point.Size > 0 {
			radius = math.Min(point.Size, 15)
		}

		sb.WriteString(fmt.Sprintf(`  <circle cx="%.2f" cy="%.2f" r="%.2f" fill="#3b82f6" fill-opacity="0.6" stroke="#2563eb" stroke-width="1"/>`,
			x, y, radius))
		sb.WriteString("\n")
	}

	// Axis labels
	if config.XLabel != "" {
		sb.WriteString(fmt.Sprintf(`  <text x="%d" y="%d" text-anchor="middle" font-size="14" fill="#374151">%s</text>`,
			config.Width/2, config.Height-10, config.XLabel))
		sb.WriteString("\n")
	}
	if config.YLabel != "" {
		sb.WriteString(fmt.Sprintf(`  <text x="15" y="%d" text-anchor="middle" font-size="14" fill="#374151" transform="rotate(-90 15 %d)">%s</text>`,
			config.Height/2, config.Height/2, config.YLabel))
		sb.WriteString("\n")
	}

	sb.WriteString("</svg>")

	return sb.String(), nil
}

// CreateHeatmap generates a heatmap SVG using SCKelemen libraries
func CreateHeatmap(config types.HeatmapConfig) (string, error) {
	rows := len(config.Data.Rows)
	cols := len(config.Data.Columns)

	if rows == 0 || cols == 0 {
		return "", fmt.Errorf("empty heatmap data")
	}

	// Find min/max for color scaling
	minVal, maxVal := math.MaxFloat64, -math.MaxFloat64
	for _, row := range config.Data.Values {
		for _, val := range row {
			if val < minVal {
				minVal = val
			}
			if val > maxVal {
				maxVal = val
			}
		}
	}

	// Chart dimensions
	margin := 80.0
	cellSize := math.Min(
		(float64(config.Width)-2*margin)/float64(cols),
		(float64(config.Height)-2*margin)/float64(rows),
	)

	// Start building SVG
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d">`,
		config.Width, config.Height))
	sb.WriteString("\n")

	// Background
	sb.WriteString(fmt.Sprintf(`  <rect width="%d" height="%d" fill="#ffffff"/>`, config.Width, config.Height))
	sb.WriteString("\n")

	// Title
	if config.Title != "" {
		sb.WriteString(fmt.Sprintf(`  <text x="%d" y="30" text-anchor="middle" font-size="20" font-weight="bold" fill="#1f2937">%s</text>`,
			config.Width/2, config.Title))
		sb.WriteString("\n")
	}

	// Draw cells
	for i, row := range config.Data.Values {
		for j, val := range row {
			x := margin + float64(j)*cellSize
			y := margin + float64(i)*cellSize

			// Color based on value (viridis-like gradient)
			normalized := (val - minVal) / (maxVal - minVal)
			cellColor := interpolateColor(normalized)

			sb.WriteString(fmt.Sprintf(`  <rect x="%.2f" y="%.2f" width="%.2f" height="%.2f" fill="%s" stroke="#ffffff" stroke-width="1"/>`,
				x, y, cellSize, cellSize, cellColor))

			// Show value if enabled
			if config.ShowValue {
				textX := x + cellSize/2
				textY := y + cellSize/2 + 4
				sb.WriteString(fmt.Sprintf(`  <text x="%.2f" y="%.2f" text-anchor="middle" font-size="10" fill="#ffffff">%.1f</text>`,
					textX, textY, val))
			}

			sb.WriteString("\n")
		}
	}

	// Column labels
	for j, col := range config.Data.Columns {
		x := margin + float64(j)*cellSize + cellSize/2
		y := margin - 10
		sb.WriteString(fmt.Sprintf(`  <text x="%.2f" y="%.2f" text-anchor="middle" font-size="11" fill="#374151">%s</text>`,
			x, y, col))
		sb.WriteString("\n")
	}

	// Row labels
	for i, row := range config.Data.Rows {
		x := margin - 10
		y := margin + float64(i)*cellSize + cellSize/2 + 4
		sb.WriteString(fmt.Sprintf(`  <text x="%.2f" y="%.2f" text-anchor="end" font-size="11" fill="#374151">%s</text>`,
			x, y, row))
		sb.WriteString("\n")
	}

	sb.WriteString("</svg>")

	return sb.String(), nil
}

// interpolateColor creates a viridis-like color gradient
func interpolateColor(t float64) string {
	// Simple viridis approximation
	t = math.Max(0, math.Min(1, t))

	if t < 0.25 {
		// Purple to blue
		r := int(68 + (30-68)*(t/0.25))
		g := int(1 + (136-1)*(t/0.25))
		b := int(84 + (229-84)*(t/0.25))
		return fmt.Sprintf("#%02x%02x%02x", r, g, b)
	} else if t < 0.5 {
		// Blue to cyan
		t2 := (t - 0.25) / 0.25
		r := int(30 + (53-30)*t2)
		g := int(136 + (183-136)*t2)
		b := int(229 + (207-229)*t2)
		return fmt.Sprintf("#%02x%02x%02x", r, g, b)
	} else if t < 0.75 {
		// Cyan to yellow
		t2 := (t - 0.5) / 0.25
		r := int(53 + (253-53)*t2)
		g := int(183 + (231-183)*t2)
		b := int(207 + (37-207)*t2)
		return fmt.Sprintf("#%02x%02x%02x", r, g, b)
	} else {
		// Yellow to white
		t2 := (t - 0.75) / 0.25
		r := int(253 + (255-253)*t2)
		g := int(231 + (255-231)*t2)
		b := int(37 + (255-37)*t2)
		return fmt.Sprintf("#%02x%02x%02x", r, g, b)
	}
}

// Keep unused imports to avoid compiler errors
var _ = units.Pixel
var _ *color.Color
var _ = &layout.Node{}
