package dataviz

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

// renderLinearHeatmapTerminal renders a linear heatmap using block characters
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

	// Render blocks
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
		b.WriteString(blocks[blockIndex])
	}
	b.WriteString("\n")

	return TerminalOutput{Content: b.String()}
}

// renderWeeksHeatmapTerminal renders a GitHub-style weeks heatmap
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

	currentDate := startDate

	// Render grid
	for day := 0; day < 7; day++ {
		for week := 0; week < weeks; week++ {
			key := currentDate.Format("2006-01-02")
			count := dayMap[key]
			ratio := float64(count) / float64(maxCount)
			blockIndex := int(ratio * float64(len(blocks)-1))
			if blockIndex >= len(blocks) {
				blockIndex = len(blocks) - 1
			}
			b.WriteString(blocks[blockIndex])
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

	// Use simple ASCII line graph for terminal
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

	// Plot points
	for i := 0; i < width && i < len(data.Points); i++ {
		point := data.Points[i]
		y := float64(height-1) - (float64(point.Value-minValue)/float64(valueRange))*float64(height-1)
		yInt := int(math.Round(y))
		if yInt >= 0 && yInt < height {
			grid[yInt][i] = "•"
		}

		// Draw line to next point
		if i < width-1 && i < len(data.Points)-1 {
			nextPoint := data.Points[i+1]
			nextY := float64(height-1) - (float64(nextPoint.Value-minValue)/float64(valueRange))*float64(height-1)
			nextYInt := int(math.Round(nextY))

			// Draw vertical line between points
			startY, endY := yInt, nextYInt
			if startY > endY {
				startY, endY = endY, startY
			}
			for y := startY; y <= endY && y < height; y++ {
				if grid[y][i] == " " {
					grid[y][i] = "│"
				}
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

// renderBarChartTerminal renders a bar chart using block characters
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

	// Render each bar
	for i := 0; i < numBars; i++ {
		bar := data.Bars[i]

		// Calculate bar length (up to bounds.Width)
		totalValue := bar.Value + bar.Secondary
		barLength := int((float64(totalValue) / float64(maxValue)) * float64(bounds.Width-10))

		if data.Stacked && bar.Secondary > 0 {
			// Stacked bars
			primaryLength := int((float64(bar.Value) / float64(maxValue)) * float64(bounds.Width-10))
			secondaryLength := int((float64(bar.Secondary) / float64(maxValue)) * float64(bounds.Width-10))

			// Primary (lighter)
			b.WriteString(strings.Repeat("▒", primaryLength))
			// Secondary (darker)
			b.WriteString(strings.Repeat("█", secondaryLength))
		} else {
			// Single bar
			b.WriteString(strings.Repeat("█", barLength))
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
