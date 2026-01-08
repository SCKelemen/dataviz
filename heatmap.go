package dataviz

import (
	"fmt"
	"strings"
	"time"

	design "github.com/SCKelemen/design-system"
	rendersvg "github.com/SCKelemen/svg"
)

// RenderLinearHeatmap renders a linear (horizontal) heatmap
func RenderLinearHeatmap(data HeatmapData, x, y int, width, height int, color string, designTokens *design.DesignTokens) string {
	var b strings.Builder

	if len(data.Days) == 0 {
		return ""
	}

	// Calculate how many days we need to show (30 for 30-day activity)
	numDays := 30
	if len(data.Days) < numDays {
		numDays = len(data.Days)
	}

	// Calculate square size to fill the entire width
	// Total width = numDays * squareSize + (numDays - 1) * gap
	// Use 1px gap between squares
	totalGapWidth := float64(numDays - 1) // 1px gap between each square
	availableForSquares := float64(width) - totalGapWidth
	squareSize := availableForSquares / float64(numDays)

	// Ensure square size is at least 1px
	if squareSize < 1 {
		squareSize = 1
	}

	// Ensure we show all days
	maxDays := numDays
	if maxDays > len(data.Days) {
		maxDays = len(data.Days)
	}

	maxCount := 0
	for _, day := range data.Days {
		if day.Count > maxCount {
			maxCount = day.Count
		}
	}
	if maxCount == 0 {
		maxCount = 1
	}

	// Position squares at the top of the content area with proper spacing
	squareY := 8.0 // 8px spacing after title

	b.WriteString(fmt.Sprintf(`<g transform="translate(%d, %d)">`, x, y))

	for i := 0; i < maxDays && i < len(data.Days); i++ {
		day := data.Days[i]

		// Calculate contribution ratio and adjust color lightness
		ratio := float64(day.Count) / float64(maxCount)
		adjustedColor := AdjustColorForContribution(color, ratio)

		// Calculate position to fill the entire width from left to right
		squareX := float64(i) * (squareSize + 1.0)

		// Use adjusted color with luminance instead of opacity
		style := rendersvg.Style{Fill: adjustedColor}
		b.WriteString(rendersvg.RoundedRect(squareX, squareY, squareSize, squareSize, 2, 0, style))
		b.WriteString("\n")
	}

	b.WriteString(`</g>`)
	return b.String()
}

// RenderWeeksHeatmap renders a GitHub-style weeks heatmap (grid of weeks)
func RenderWeeksHeatmap(data HeatmapData, x, y int, width, height int, color string, designTokens *design.DesignTokens) string {
	var b strings.Builder

	if len(data.Days) == 0 {
		return ""
	}

	// Reserve space for left-side day labels (Mon, Wed, Fri)
	labelWidth := 30.0
	availableWidth := float64(width) - labelWidth

	// Calculate weeks and days
	weeks := 53 // GitHub shows ~53 weeks
	daysPerWeek := 7

	// Calculate cell size to ensure perfect squares
	cellSizeByWidth := (availableWidth - float64(weeks)) / float64(weeks)
	cellSizeByHeight := float64(height) / float64(daysPerWeek)

	// Use the smaller dimension to ensure perfect squares fit
	cellSize := cellSizeByHeight
	if cellSizeByWidth < cellSize {
		cellSize = cellSizeByWidth
	}

	// Ensure minimum size
	if cellSize < 1 {
		cellSize = 1
	}

	// Recalculate week width based on square cell size (cellSize + 1px gap)
	weekWidth := cellSize + 1.0

	maxCount := 0
	for _, day := range data.Days {
		if day.Count > maxCount {
			maxCount = day.Count
		}
	}
	if maxCount == 0 {
		maxCount = 1
	}

	// Calculate total heatmap height for vertical centering
	totalHeatmapHeight := float64(daysPerWeek) * cellSize

	// Position heatmap to the right of labels, centered vertically
	offsetX := labelWidth
	offsetY := (float64(height) - totalHeatmapHeight) / 2.0

	b.WriteString(fmt.Sprintf(`<g transform="translate(%d, %d)">`, x, y))
	b.WriteString(fmt.Sprintf(`<g transform="translate(%.1f, %.1f)">`, offsetX, offsetY))

	// Create a map of dates to counts for quick lookup
	dayMap := make(map[string]int)
	for _, day := range data.Days {
		key := day.Date.Format("2006-01-02")
		dayMap[key] = day.Count
	}

	// Start from the first day of the year (or start date)
	startDate := data.StartDate
	if startDate.IsZero() {
		startDate = time.Now().AddDate(0, 0, -weeks*7)
	}

	// Find the first Sunday before or on start date
	weekday := int(startDate.Weekday())
	if weekday != 0 { // Not Sunday
		startDate = startDate.AddDate(0, 0, -weekday)
	}

	currentDate := startDate
	weekIdx := 0

	for weekIdx < weeks {
		for dayOfWeek := 0; dayOfWeek < daysPerWeek; dayOfWeek++ {
			if currentDate.After(data.EndDate) && !data.EndDate.IsZero() {
				break
			}

			key := currentDate.Format("2006-01-02")
			count := dayMap[key]

			// Calculate contribution ratio and adjust color lightness
			ratio := float64(count) / float64(maxCount)
			adjustedColor := AdjustColorForContribution(color, ratio)

			cellX := float64(weekIdx) * weekWidth
			cellY := float64(dayOfWeek) * cellSize

			// Use adjusted color with luminance instead of opacity
			style := rendersvg.Style{Fill: adjustedColor}
			b.WriteString(rendersvg.RoundedRect(cellX, cellY, cellSize, cellSize, 2, 0, style))
			b.WriteString("\n")

			currentDate = currentDate.AddDate(0, 0, 1)
		}
		weekIdx++
	}

	// Add day labels on the left
	dayLabels := []string{"", "Mon", "", "Wed", "", "Fri", ""}
	for i, label := range dayLabels {
		if label != "" {
			textStyle := rendersvg.Style{
				Fill:             designTokens.Color,
				Class:            "mono smaller",
				TextAnchor:       rendersvg.TextAnchorEnd,
				DominantBaseline: rendersvg.DominantBaselineMiddle,
			}
			labelY := float64(i)*cellSize + cellSize/2
			b.WriteString(rendersvg.Text(label, -5, labelY, textStyle))
			b.WriteString("\n")
		}
	}

	b.WriteString(`</g>`) // Close centering transform
	b.WriteString(`</g>`) // Close main transform
	return b.String()
}
