package charts

import (
	"fmt"
	"strings"

	design "github.com/SCKelemen/design-system"
)

// RenderStatCard renders a statistics card
func RenderStatCard(data StatCardData, x, y int, width, height int, designTokens *design.DesignTokens) string {
	var b strings.Builder

	// Always use 4.5 to match example style
	radius := 4.5

	b.WriteString(fmt.Sprintf(`<g transform="translate(%d, %d)">`, x, y))
	// Double border style
	b.WriteString(fmt.Sprintf(`<rect x="0.5" y="0.5" width="%d" height="%d" rx="%.1f" stroke="rgba(0,0,0,0.15)"/>`,
		width-1, height-1, radius))
	b.WriteString(fmt.Sprintf(`<rect x="0.5" y="0.5" width="%d" height="%d" rx="%.1f" stroke="rgba(255,255,255,0.3)"/>`,
		width-1, height-1, radius))

	// Title
	b.WriteString(fmt.Sprintf(`<text x="10" y="24" class="sans small bold" fill="%s">%s</text>`, data.Color, data.Title))

	// Value/subtitle
	subtitleColor := "#777"
	if designTokens != nil {
		subtitleColor = designTokens.Color
	}
	b.WriteString(fmt.Sprintf(`<text x="10" y="56" text-anchor="start" class="mono smaller bold" fill="%s">%s</text>`, subtitleColor, data.Subtitle))

	// Change indicator
	changeColor := "#CF3E3E"
	changeText := "past month"
	textWidthEstimate := len(changeText) * 8
	changeTextX := width - designTokens.Layout.CardPaddingRight
	arrowX := changeTextX - textWidthEstimate - 15
	arrowY := 46

	b.WriteString(fmt.Sprintf(`<text x="%d" y="56" text-anchor="end" class="mono smaller bold" fill="%s">%s</text>`,
		changeTextX, changeColor, changeText))

	// Arrow icon
	if data.Change != 0 {
		arrowPath := "M2.02989 7.836L2.02994 7.83604L5.15494 10.961L5.15498 10.9611C5.24404 11.05 5.36477 11.1 5.49065 11.1C5.61652 11.1 5.73725 11.05 5.82631 10.9611L5.82636 10.961L8.95018 7.83722C8.99622 7.79398 9.03317 7.742 9.05888 7.6843C9.08484 7.62603 9.0988 7.56313 9.09993 7.49936C9.10105 7.43558 9.08932 7.37223 9.06543 7.31308C9.04154 7.25393 9.00598 7.20021 8.96088 7.1551C8.91577 7.11 8.86204 7.07444 8.8029 7.05055C8.74375 7.02666 8.6804 7.01492 8.61662 7.01605C8.55284 7.01717 8.48994 7.03113 8.43168 7.0571C8.37398 7.0828 8.322 7.11976 8.27876 7.16579L5.96565 9.47891L5.96565 2.50033C5.96565 2.37435 5.9156 2.25353 5.82652 2.16445C5.73744 2.07537 5.61662 2.02533 5.49065 2.02533C5.36467 2.02533 5.24385 2.07537 5.15477 2.16445C5.06569 2.25353 5.01565 2.37435 5.01565 2.50033L5.01565 9.47891L2.70136 7.16462L2.70131 7.16458C2.61225 7.07562 2.49152 7.02566 2.36565 7.02566C2.23977 7.02566 2.11904 7.07562 2.02998 7.16458L2.10065 7.23533L2.02989 7.16466C1.94094 7.25373 1.89098 7.37446 1.89098 7.50033C1.89098 7.62621 1.94094 7.74693 2.02989 7.836Z"
		b.WriteString(fmt.Sprintf(`<g transform="translate(%d, %d)">`, arrowX, arrowY))
		b.WriteString(`<svg width="12" height="12" viewBox="0 0 12 12" fill="none" xmlns="http://www.w3.org/2000/svg">`)
		b.WriteString(fmt.Sprintf(`<path d="%s" fill="%s" stroke="%s" stroke-width="0.5"/>`, arrowPath, changeColor, changeColor))
		b.WriteString(`</svg>`)
		b.WriteString(`</g>`)
	}

	// Legends (if provided)
	if data.Legend1 != "" && data.TrendColor != "" {
		legendSize := 10.0
		legendSpacing := 15.0
		textSpacing := 5.0

		legend1TextWidth := float64(len(data.Legend1)) * 7.0
		legend2TextWidth := 0.0
		if data.Legend2 != "" {
			legend2TextWidth = float64(len(data.Legend2)) * 7.0
		}
		maxTextWidth := legend1TextWidth
		if legend2TextWidth > maxTextWidth {
			maxTextWidth = legend2TextWidth
		}

		totalLegendWidth := legendSize + textSpacing + maxTextWidth
		legendX := float64(width) - float64(designTokens.Layout.CardPaddingRight) - totalLegendWidth
		legendY := 15

		// First legend
		b.WriteString(fmt.Sprintf(`<rect x="%.0f" y="%d" width="%.0f" height="%.0f" rx="10" fill="%s"/>`,
			legendX, legendY, legendSize, legendSize, data.TrendColor))
		textX := legendX + legendSize + textSpacing
		b.WriteString(fmt.Sprintf(`<text x="%.0f" y="%.0f" class="sans smaller" fill="%s" dominant-baseline="middle">%s</text>`,
			textX, float64(legendY)+legendSize/2, data.TrendColor, data.Legend1))

		// Second legend (if provided)
		if data.Legend2 != "" && data.TrendColor2 != "" {
			legendY2 := legendY + int(legendSpacing)
			b.WriteString(fmt.Sprintf(`<rect x="%.0f" y="%d" width="%.0f" height="%.0f" rx="10" fill="%s"/>`,
				legendX, legendY2, legendSize, legendSize, data.TrendColor2))
			b.WriteString(fmt.Sprintf(`<text x="%.0f" y="%.0f" class="sans smaller" fill="%s" dominant-baseline="middle">%s</text>`,
				textX, float64(legendY2)+legendSize/2, data.TrendColor2, data.Legend2))
		}
	}

	// Mini trend graph (if provided)
	if len(data.TrendData) > 0 && data.TrendColor != "" {
		graphX := designTokens.Layout.CardPaddingLeft
		graphWidth := width - designTokens.Layout.CardPaddingLeft - designTokens.Layout.CardPaddingRight

		subtitleY := 56
		subtitleTextHeight := 16
		graphTopSpacing := 8
		graphStartY := subtitleY + subtitleTextHeight + graphTopSpacing

		minGraphHeight := 15.0
		availableHeight := float64(height) - float64(graphStartY) - float64(designTokens.Layout.CardPaddingBottom)
		if availableHeight < 0 {
			availableHeight = 0
		}

		graphHeight := availableHeight
		if graphHeight < minGraphHeight {
			if availableHeight > 0 {
				graphHeight = availableHeight
			} else {
				graphHeight = 0
			}
		}

		graphY := graphStartY
		hasDualBars := data.TrendColor2 != ""

		maxValue := 0.0
		for _, point := range data.TrendData {
			totalValue := float64(point.Value)
			if totalValue > maxValue {
				maxValue = totalValue
			}
		}
		if maxValue == 0 {
			maxValue = 1
		}

		numBars := len(data.TrendData)
		if numBars > 30 {
			numBars = 30
		}
		barWidth := float64(graphWidth) / float64(numBars)
		barSpacing := barWidth * 0.1
		actualBarWidth := barWidth - barSpacing

		maxBaseY := float64(height) - float64(designTokens.Layout.CardPaddingBottom)
		baseY := float64(graphY) + graphHeight
		if baseY > maxBaseY {
			baseY = maxBaseY
			graphHeight = baseY - float64(graphY)
			if graphHeight < 0 {
				graphHeight = 0
			}
		}

		for i, point := range data.TrendData {
			if i >= numBars {
				break
			}

			barX := float64(graphX) + float64(i)*barWidth

			if hasDualBars {
				totalValue := float64(point.Value)
				primaryValue := totalValue * 0.6
				secondaryValue := totalValue * 0.4

				primaryHeight := (primaryValue / maxValue) * float64(graphHeight)
				secondaryHeight := (secondaryValue / maxValue) * float64(graphHeight)

				totalBarHeight := primaryHeight + secondaryHeight
				if totalBarHeight > float64(graphHeight) {
					scale := float64(graphHeight) / totalBarHeight
					primaryHeight *= scale
					secondaryHeight *= scale
				}

				secondaryY := baseY - secondaryHeight
				b.WriteString(fmt.Sprintf(`<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="%s" />`,
					barX, secondaryY, actualBarWidth, secondaryHeight, data.TrendColor2))

				primaryY := secondaryY - primaryHeight
				b.WriteString(fmt.Sprintf(`<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="%s" />`,
					barX, primaryY, actualBarWidth, primaryHeight, data.TrendColor))
			} else {
				barHeight := (float64(point.Value) / maxValue) * float64(graphHeight)
				barY := baseY - barHeight
				b.WriteString(fmt.Sprintf(`<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="%s" />`,
					barX, barY, actualBarWidth, barHeight, data.TrendColor))
			}
		}
	}

	b.WriteString(`</g>`)
	return b.String()
}
