package charts

import (
	"fmt"
	"math"
	"strings"
)

// TerminalOutput wraps terminal output
type TerminalOutput struct {
	Content string
}

func (t TerminalOutput) String() string {
	return t.Content
}

// TerminalRenderer implements terminal-based rendering
type TerminalRenderer struct{}

// NewTerminalRenderer creates a new terminal renderer
func NewTerminalRenderer() *TerminalRenderer {
	return &TerminalRenderer{}
}

// RenderHeatmap renders a heatmap to terminal
func (r *TerminalRenderer) RenderHeatmap(data HeatmapData, bounds Bounds, config RenderConfig) Output {
	if data.Type == "weeks" {
		return r.renderWeeksHeatmapTerminal(data, bounds, config)
	}
	return r.renderLinearHeatmapTerminal(data, bounds, config)
}

// RenderLineGraph renders a line graph to terminal
func (r *TerminalRenderer) RenderLineGraph(data LineGraphData, bounds Bounds, config RenderConfig) Output {
	return r.renderLineGraphTerminal(data, bounds, config)
}

// RenderBarChart renders a bar chart to terminal
func (r *TerminalRenderer) RenderBarChart(data BarChartData, bounds Bounds, config RenderConfig) Output {
	return r.renderBarChartTerminal(data, bounds, config)
}

// RenderStatCard renders a stat card to terminal
func (r *TerminalRenderer) RenderStatCard(data StatCardData, bounds Bounds, config RenderConfig) Output {
	return r.renderStatCardTerminal(data, bounds, config)
}

// RenderAreaChart renders an area chart to terminal
func (r *TerminalRenderer) RenderAreaChart(data AreaChartData, bounds Bounds, config RenderConfig) Output {
	return r.renderAreaChartTerminal(data, bounds, config)
}

// RenderScatterPlot renders a scatter plot to terminal
func (r *TerminalRenderer) RenderScatterPlot(data ScatterPlotData, bounds Bounds, config RenderConfig) Output {
	return r.renderScatterPlotTerminal(data, bounds, config)
}

// renderLinearHeatmapTerminal renders a linear heatmap using block characters with color gradients
func (r *TerminalRenderer) renderLinearHeatmapTerminal(data HeatmapData, bounds Bounds, config RenderConfig) Output {
	var b strings.Builder

	if len(data.Days) == 0 {
		return TerminalOutput{Content: ""}
	}

	// Block characters for intensity (darkest to lightest)
	blocks := []string{" ", "░", "▒", "▓", "█"}

	// Calculate max count for scaling
	maxCount := 0
	for _, day := range data.Days {
		if day.Count > maxCount {
			maxCount = day.Count
		}
	}
	if maxCount == 0 {
		maxCount = 1
	}

	// Limit to reasonable number of days for terminal
	numDays := len(data.Days)
	if numDays > bounds.Width {
		numDays = bounds.Width
	}

	// Determine color mode
	colorMode := TerminalColorTrue
	if config.DesignTokens != nil && config.DesignTokens.Mode == "basic" {
		colorMode = TerminalColor256
	}

	// Get base color from config
	baseColor := config.Color
	if baseColor == "" && config.DesignTokens != nil {
		baseColor = config.DesignTokens.Accent
	}
	if baseColor == "" {
		baseColor = "#40C463" // GitHub green
	}

	// Create color gradient from low to high intensity
	lowColor := "#161B22"  // Dark background
	highColor := baseColor // Accent color

	// Generate gradient colors
	gradientSteps := 5
	gradient := InterpolateColorGradient(lowColor, highColor, gradientSteps, colorMode)

	ansiReset := "\x1b[0m"

	// Render blocks with color
	for i := 0; i < numDays; i++ {
		if i >= len(data.Days) {
			break
		}
		day := data.Days[i]
		ratio := float64(day.Count) / float64(maxCount)
		blockIndex := int(ratio * float64(len(blocks)-1))
		if blockIndex >= len(blocks) {
			blockIndex = len(blocks) - 1
		}

		// Get color for this intensity
		colorIndex := int(ratio * float64(len(gradient)-1))
		if colorIndex >= len(gradient) {
			colorIndex = len(gradient) - 1
		}

		// Apply color
		if gradient[colorIndex] != "" {
			b.WriteString(gradient[colorIndex])
		}
		b.WriteString(blocks[blockIndex])
		if gradient[colorIndex] != "" {
			b.WriteString(ansiReset)
		}
	}
	b.WriteString("\n")

	return TerminalOutput{Content: b.String()}
}

// renderWeeksHeatmapTerminal renders a GitHub-style weeks heatmap with color gradients
func (r *TerminalRenderer) renderWeeksHeatmapTerminal(data HeatmapData, bounds Bounds, config RenderConfig) Output {
	var b strings.Builder

	if len(data.Days) == 0 {
		return TerminalOutput{Content: ""}
	}

	blocks := []string{" ", "░", "▒", "▓", "█"}

	// Calculate max count
	maxCount := 0
	for _, day := range data.Days {
		if day.Count > maxCount {
			maxCount = day.Count
		}
	}
	if maxCount == 0 {
		maxCount = 1
	}

	// Create date map
	dayMap := make(map[string]int)
	for _, day := range data.Days {
		key := day.Date.Format("2006-01-02")
		dayMap[key] = day.Count
	}

	// Render 7 rows (days of week) x N columns (weeks)
	weeks := bounds.Width / 2 // Each cell takes ~2 chars
	if weeks > 52 {
		weeks = 52
	}

	startDate := data.StartDate
	if startDate.IsZero() {
		// Use reasonable default
		startDate = data.Days[0].Date
	}

	// Align to Sunday
	weekday := int(startDate.Weekday())
	if weekday != 0 {
		startDate = startDate.AddDate(0, 0, -weekday)
	}

	// Determine color mode
	colorMode := TerminalColorTrue
	if config.DesignTokens != nil && config.DesignTokens.Mode == "basic" {
		colorMode = TerminalColor256
	}

	// Get base color from config
	baseColor := config.Color
	if baseColor == "" && config.DesignTokens != nil {
		baseColor = config.DesignTokens.Accent
	}
	if baseColor == "" {
		baseColor = "#40C463" // GitHub green
	}

	// Create color gradient from low to high intensity
	lowColor := "#161B22"  // Dark background
	highColor := baseColor // Accent color

	// Generate gradient colors
	gradientSteps := 5
	gradient := InterpolateColorGradient(lowColor, highColor, gradientSteps, colorMode)

	ansiReset := "\x1b[0m"

	currentDate := startDate

	// Render grid with colors
	for day := 0; day < 7; day++ {
		for week := 0; week < weeks; week++ {
			key := currentDate.Format("2006-01-02")
			count := dayMap[key]
			ratio := float64(count) / float64(maxCount)
			blockIndex := int(ratio * float64(len(blocks)-1))
			if blockIndex >= len(blocks) {
				blockIndex = len(blocks) - 1
			}

			// Get color for this intensity
			colorIndex := int(ratio * float64(len(gradient)-1))
			if colorIndex >= len(gradient) {
				colorIndex = len(gradient) - 1
			}

			// Apply color
			if gradient[colorIndex] != "" {
				b.WriteString(gradient[colorIndex])
			}
			b.WriteString(blocks[blockIndex])
			if gradient[colorIndex] != "" {
				b.WriteString(ansiReset)
			}
			b.WriteString(" ")
			currentDate = currentDate.AddDate(0, 0, 1)
		}
		b.WriteString("\n")
	}

	return TerminalOutput{Content: b.String()}
}

// renderLineGraphTerminal renders a line graph using Braille characters
func (r *TerminalRenderer) renderLineGraphTerminal(data LineGraphData, bounds Bounds, config RenderConfig) Output {
	var b strings.Builder

	if len(data.Points) == 0 {
		return TerminalOutput{Content: ""}
	}

	// Find min/max for scaling
	minValue, maxValue := data.Points[0].Value, data.Points[0].Value
	for _, point := range data.Points {
		if point.Value < minValue {
			minValue = point.Value
		}
		if point.Value > maxValue {
			maxValue = point.Value
		}
	}

	valueRange := maxValue - minValue
	if valueRange == 0 {
		valueRange = 1
	}

	// Use braille characters for smooth rendering
	width := bounds.Width
	if width > 120 {
		width = 120
	}
	height := bounds.Height
	if height > 30 {
		height = 30
	}

	// Create braille canvas (each char is 2x4 pixels)
	canvas := NewBrailleCanvas(width, height)

	// Convert data points to canvas coordinates
	braillePoints := make([]Point, 0, len(data.Points))
	for i, point := range data.Points {
		// X coordinate: scale to canvas width
		x := float64(i) / float64(len(data.Points)-1) * float64(width*2-1)
		if math.IsNaN(x) {
			x = 0
		}

		// Y coordinate: invert because canvas Y increases downward
		normalizedY := float64(point.Value-minValue) / float64(valueRange)
		y := float64(height*4-1) - (normalizedY * float64(height*4-1))
		if math.IsNaN(y) {
			y = float64(height*4 - 1)
		}

		braillePoints = append(braillePoints, Point{X: x, Y: y})
	}

	// Draw the curve
	canvas.DrawCurve(braillePoints)

	// Render canvas to string
	rendered := canvas.Render()

	// Apply color if specified
	colorMode := TerminalColorTrue
	if config.DesignTokens != nil && config.DesignTokens.Mode == "basic" {
		colorMode = TerminalColor256
	}

	lineColor := data.Color
	if lineColor == "" {
		lineColor = "#2196F3" // Default blue
	}

	ansiColor := ColorForeground(lineColor, colorMode)
	ansiReset := "\x1b[0m"

	if ansiColor != "" {
		b.WriteString(ansiColor)
	}
	b.WriteString(rendered)
	if ansiColor != "" {
		b.WriteString(ansiReset)
	}

	return TerminalOutput{Content: b.String()}
}

// renderBarChartTerminal renders a bar chart using block characters with colors
func (r *TerminalRenderer) renderBarChartTerminal(data BarChartData, bounds Bounds, config RenderConfig) Output {
	var b strings.Builder

	if len(data.Bars) == 0 {
		return TerminalOutput{Content: ""}
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

	// Determine bar width based on bounds
	numBars := len(data.Bars)
	if numBars > bounds.Width/3 {
		numBars = bounds.Width / 3
	}

	// Determine color mode
	colorMode := TerminalColorTrue
	if config.DesignTokens != nil && config.DesignTokens.Mode == "basic" {
		colorMode = TerminalColor256
	}

	// Get bar color
	barColor := data.Color
	if barColor == "" && config.DesignTokens != nil {
		barColor = config.DesignTokens.Accent
	}
	if barColor == "" {
		barColor = "#2196F3" // Default blue
	}

	ansiColor := ColorForeground(barColor, colorMode)
	ansiReset := "\x1b[0m"

	// Find max label length to account for spacing
	maxLabelLen := 0
	for i := 0; i < numBars; i++ {
		if i < len(data.Bars) && len(data.Bars[i].Label) > maxLabelLen {
			maxLabelLen = len(data.Bars[i].Label)
		}
	}

	// Calculate available width for bars (leave room for labels and spacing)
	labelSpace := maxLabelLen + 3 // 2 spaces + label
	if labelSpace > bounds.Width/2 {
		labelSpace = bounds.Width / 2
	}
	maxBarWidth := bounds.Width - labelSpace
	if maxBarWidth < 10 {
		maxBarWidth = 10
	}

	// Render each bar
	for i := 0; i < numBars; i++ {
		bar := data.Bars[i]

		// Calculate bar length proportional to value
		totalValue := bar.Value + bar.Secondary
		barLength := int((float64(totalValue) / float64(maxValue)) * float64(maxBarWidth))

		if data.Stacked && bar.Secondary > 0 {
			// Stacked bars with different colors
			primaryLength := int((float64(bar.Value) / float64(maxValue)) * float64(maxBarWidth))
			secondaryLength := int((float64(bar.Secondary) / float64(maxValue)) * float64(maxBarWidth))

			// Primary (lighter color)
			if ansiColor != "" {
				b.WriteString(ansiColor)
			}
			b.WriteString(strings.Repeat("▒", primaryLength))
			if ansiColor != "" {
				b.WriteString(ansiReset)
			}

			// Secondary (darker/full color)
			if ansiColor != "" {
				b.WriteString(ansiColor)
			}
			b.WriteString(strings.Repeat("█", secondaryLength))
			if ansiColor != "" {
				b.WriteString(ansiReset)
			}
		} else {
			// Single bar with color
			if ansiColor != "" {
				b.WriteString(ansiColor)
			}
			b.WriteString(strings.Repeat("█", barLength))
			if ansiColor != "" {
				b.WriteString(ansiReset)
			}
		}

		if bar.Label != "" {
			b.WriteString(fmt.Sprintf("  %s", bar.Label))
		}
		b.WriteString("\n")
	}

	return TerminalOutput{Content: b.String()}
}

// renderStatCardTerminal renders a stat card to terminal
func (r *TerminalRenderer) renderStatCardTerminal(data StatCardData, bounds Bounds, config RenderConfig) Output {
	var b strings.Builder

	// Title
	b.WriteString(fmt.Sprintf("┌─ %s ", data.Title))
	b.WriteString(strings.Repeat("─", bounds.Width-len(data.Title)-5))
	b.WriteString("┐\n")

	// Value
	b.WriteString(fmt.Sprintf("│ %s", data.Value))
	padding := bounds.Width - len(data.Value) - 3
	if padding > 0 {
		b.WriteString(strings.Repeat(" ", padding))
	}
	b.WriteString("│\n")

	// Subtitle
	if data.Subtitle != "" {
		b.WriteString(fmt.Sprintf("│ %s", data.Subtitle))
		padding := bounds.Width - len(data.Subtitle) - 3
		if padding > 0 {
			b.WriteString(strings.Repeat(" ", padding))
		}
		b.WriteString("│\n")
	}

	// Mini trend graph (simple bar representation)
	if len(data.TrendData) > 0 {
		b.WriteString("│ ")
		maxTrend := 0
		for _, point := range data.TrendData {
			if point.Value > maxTrend {
				maxTrend = point.Value
			}
		}
		if maxTrend == 0 {
			maxTrend = 1
		}

		numBars := len(data.TrendData)
		if numBars > bounds.Width-4 {
			numBars = bounds.Width - 4
		}

		for i := 0; i < numBars; i++ {
			point := data.TrendData[i]
			height := int((float64(point.Value) / float64(maxTrend)) * 4)
			if height == 0 && point.Value > 0 {
				height = 1
			}
			bars := []string{" ", "▁", "▂", "▃", "▄"}
			if height >= len(bars) {
				height = len(bars) - 1
			}
			b.WriteString(bars[height])
		}

		padding := bounds.Width - numBars - 3
		if padding > 0 {
			b.WriteString(strings.Repeat(" ", padding))
		}
		b.WriteString("│\n")
	}

	// Bottom border
	b.WriteString("└")
	b.WriteString(strings.Repeat("─", bounds.Width-2))
	b.WriteString("┘\n")

	return TerminalOutput{Content: b.String()}
}

// renderAreaChartTerminal renders an area chart to terminal
func (r *TerminalRenderer) renderAreaChartTerminal(data AreaChartData, bounds Bounds, config RenderConfig) Output {
	var b strings.Builder

	if len(data.Points) == 0 {
		return TerminalOutput{Content: ""}
	}

	// Find min/max for scaling
	minValue, maxValue := data.Points[0].Value, data.Points[0].Value
	for _, point := range data.Points {
		if point.Value < minValue {
			minValue = point.Value
		}
		if point.Value > maxValue {
			maxValue = point.Value
		}
	}

	valueRange := maxValue - minValue
	if valueRange == 0 {
		valueRange = 1
	}

	// Use block characters to fill the area
	height := bounds.Height
	if height > 20 {
		height = 20
	}
	width := bounds.Width
	if width > len(data.Points) {
		width = len(data.Points)
	}

	// Create a 2D grid
	grid := make([][]string, height)
	for i := range grid {
		grid[i] = make([]string, width)
		for j := range grid[i] {
			grid[i][j] = " "
		}
	}

	// Fill area under the curve
	for i := 0; i < width && i < len(data.Points); i++ {
		point := data.Points[i]
		y := float64(height-1) - (float64(point.Value-minValue)/float64(valueRange))*float64(height-1)
		yInt := int(math.Round(y))

		// Fill from bottom to y
		for row := yInt; row < height; row++ {
			if row >= 0 && row < height {
				grid[row][i] = "█"
			}
		}
	}

	// Render grid
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			b.WriteString(grid[y][x])
		}
		b.WriteString("\n")
	}

	return TerminalOutput{Content: b.String()}
}

// renderScatterPlotTerminal renders a scatter plot to terminal
func (r *TerminalRenderer) renderScatterPlotTerminal(data ScatterPlotData, bounds Bounds, config RenderConfig) Output {
	var b strings.Builder

	if len(data.Points) == 0 {
		return TerminalOutput{Content: ""}
	}

	// Find min/max for scaling
	minValue, maxValue := data.Points[0].Value, data.Points[0].Value
	for _, point := range data.Points {
		if point.Value < minValue {
			minValue = point.Value
		}
		if point.Value > maxValue {
			maxValue = point.Value
		}
	}

	valueRange := maxValue - minValue
	if valueRange == 0 {
		valueRange = 1
	}

	// Use ASCII art for scatter plot
	height := bounds.Height
	if height > 20 {
		height = 20
	}
	width := bounds.Width
	if width > len(data.Points) {
		width = len(data.Points)
	}

	// Create a 2D grid
	grid := make([][]string, height)
	for i := range grid {
		grid[i] = make([]string, width)
		for j := range grid[i] {
			grid[i][j] = " "
		}
	}

	// Map marker types to characters
	markerChar := "●"
	switch data.MarkerType {
	case "square":
		markerChar = "■"
	case "diamond":
		markerChar = "◆"
	case "triangle":
		markerChar = "▲"
	case "cross":
		markerChar = "+"
	case "x":
		markerChar = "×"
	case "dot", "circle":
		markerChar = "●"
	}

	// Plot points
	for i := 0; i < width && i < len(data.Points); i++ {
		point := data.Points[i]
		y := float64(height-1) - (float64(point.Value-minValue)/float64(valueRange))*float64(height-1)
		yInt := int(math.Round(y))
		if yInt >= 0 && yInt < height {
			grid[yInt][i] = markerChar
		}
	}

	// Render grid
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			b.WriteString(grid[y][x])
		}
		b.WriteString("\n")
	}

	return TerminalOutput{Content: b.String()}
}
