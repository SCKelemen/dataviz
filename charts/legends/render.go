package legends

import (
	"fmt"
	"strings"

	"github.com/SCKelemen/color"
)

// Bounds represents the dimensions and position of the legend
type Bounds struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

// Render generates the SVG string for the legend
// chartWidth and chartHeight are used to calculate positioning
func (l *Legend) Render(chartWidth, chartHeight int) string {
	if l.Position == PositionNone || len(l.Items) == 0 {
		return ""
	}

	// Calculate legend dimensions
	bounds := l.calculateBounds()

	// Calculate position based on Position setting
	x, y := l.calculatePosition(bounds, chartWidth, chartHeight)

	var sb strings.Builder

	// Legend container group
	sb.WriteString(fmt.Sprintf(`<g class="legend" transform="translate(%.1f,%.1f)">`, x, y))
	sb.WriteString("\n")

	// Background and border (if configured)
	if l.Style.Background.Alpha() > 0 || l.Style.Border.Alpha() > 0 {
		bg := "none"
		if l.Style.Background.Alpha() > 0 {
			bg = color.RGBToHex(l.Style.Background)
		}

		stroke := "none"
		strokeWidth := 0.0
		if l.Style.Border.Alpha() > 0 {
			stroke = color.RGBToHex(l.Style.Border)
			strokeWidth = l.Style.BorderWidth.Raw()
		}

		sb.WriteString(fmt.Sprintf(`  <rect x="0" y="0" width="%.1f" height="%.1f" fill="%s" stroke="%s" stroke-width="%.1f"/>`,
			bounds.Width, bounds.Height, bg, stroke, strokeWidth))
		sb.WriteString("\n")
	}

	// Render items
	padding := l.Style.Padding.Raw()
	itemX := padding
	itemY := padding

	for i, item := range l.Items {
		// Calculate item position based on layout
		if l.Layout == LayoutHorizontal && i > 0 {
			// Horizontal: move right
			prevSymbol := l.Items[i-1].Symbol
			itemX += prevSymbol.Width() + l.Style.SymbolSpacing.Raw() + estimateTextWidth(l.Items[i-1].Label, l.Style.FontSize.Raw()) + l.Style.ItemSpacing.Raw()
		} else if l.Layout == LayoutVertical && i > 0 {
			// Vertical: move down
			itemY += maxSymbolHeight(l.Items[i-1].Symbol) + l.Style.ItemSpacing.Raw()
		}

		// Render item
		sb.WriteString(l.renderItem(item, itemX, itemY))
	}

	sb.WriteString("</g>")
	sb.WriteString("\n")

	return sb.String()
}

// renderItem renders a single legend item
func (l *Legend) renderItem(item LegendItem, x, y float64) string {
	var sb strings.Builder

	symbolHeight := item.Symbol.Height()
	symbolWidth := item.Symbol.Width()

	// Center symbol vertically with text
	fontSize := l.Style.FontSize.Raw()
	symbolY := y + (fontSize-symbolHeight)/2

	// Symbol
	sb.WriteString(fmt.Sprintf(`  <g transform="translate(%.1f,%.1f)">`, x, symbolY))
	sb.WriteString(item.Symbol.Render())
	sb.WriteString("</g>\n")

	// Label
	labelX := x + symbolWidth + l.Style.SymbolSpacing.Raw()
	labelY := y + fontSize*0.85 // Baseline adjustment

	labelText := item.Label
	if item.Value != "" {
		labelText = fmt.Sprintf("%s (%s)", item.Label, item.Value)
	}

	sb.WriteString(fmt.Sprintf(`  <text x="%.1f" y="%.1f" font-family="%s" font-size="%.1f" fill="%s">%s</text>`,
		labelX, labelY, l.Style.FontFamily, fontSize, color.RGBToHex(l.Style.TextColor), labelText))
	sb.WriteString("\n")

	return sb.String()
}

// calculateBounds calculates the dimensions of the legend
func (l *Legend) calculateBounds() Bounds {
	if len(l.Items) == 0 {
		return Bounds{}
	}

	padding := l.Style.Padding.Raw()
	itemSpacing := l.Style.ItemSpacing.Raw()
	symbolSpacing := l.Style.SymbolSpacing.Raw()
	fontSize := l.Style.FontSize.Raw()

	var width, height float64

	if l.Layout == LayoutHorizontal {
		// Horizontal: sum widths, max height
		for i, item := range l.Items {
			itemWidth := item.Symbol.Width() + symbolSpacing + estimateTextWidth(item.Label, fontSize)
			if item.Value != "" {
				itemWidth += estimateTextWidth(fmt.Sprintf(" (%s)", item.Value), fontSize)
			}
			width += itemWidth

			if i > 0 {
				width += itemSpacing
			}

			itemHeight := maxFloat(item.Symbol.Height(), fontSize)
			if itemHeight > height {
				height = itemHeight
			}
		}
	} else { // LayoutVertical
		// Vertical: max width, sum heights
		for i, item := range l.Items {
			itemWidth := item.Symbol.Width() + symbolSpacing + estimateTextWidth(item.Label, fontSize)
			if item.Value != "" {
				itemWidth += estimateTextWidth(fmt.Sprintf(" (%s)", item.Value), fontSize)
			}
			if itemWidth > width {
				width = itemWidth
			}

			itemHeight := maxFloat(item.Symbol.Height(), fontSize)
			height += itemHeight

			if i > 0 {
				height += itemSpacing
			}
		}
	}

	return Bounds{
		Width:  width + 2*padding,
		Height: height + 2*padding,
	}
}

// calculatePosition calculates the x,y position based on Position setting
func (l *Legend) calculatePosition(bounds Bounds, chartWidth, chartHeight int) (x, y float64) {
	margin := 10.0 // Margin from chart edges

	switch l.Position {
	case PositionTopLeft:
		x = margin
		y = margin

	case PositionTopRight:
		x = float64(chartWidth) - bounds.Width - margin
		y = margin

	case PositionTopCenter:
		x = (float64(chartWidth) - bounds.Width) / 2
		y = margin

	case PositionBottomLeft:
		x = margin
		y = float64(chartHeight) - bounds.Height - margin

	case PositionBottomRight:
		x = float64(chartWidth) - bounds.Width - margin
		y = float64(chartHeight) - bounds.Height - margin

	case PositionBottomCenter:
		x = (float64(chartWidth) - bounds.Width) / 2
		y = float64(chartHeight) - bounds.Height - margin

	case PositionLeft:
		x = margin
		y = (float64(chartHeight) - bounds.Height) / 2

	case PositionRight:
		x = float64(chartWidth) - bounds.Width - margin
		y = (float64(chartHeight) - bounds.Height) / 2

	default: // PositionTopRight
		x = float64(chartWidth) - bounds.Width - margin
		y = margin
	}

	return x, y
}

// GetBounds returns the calculated bounds of the legend
// Useful for reserving space in charts
func (l *Legend) GetBounds(chartWidth, chartHeight int) Bounds {
	bounds := l.calculateBounds()
	x, y := l.calculatePosition(bounds, chartWidth, chartHeight)
	bounds.X = x
	bounds.Y = y
	return bounds
}

// Helper functions

func estimateTextWidth(text string, fontSize float64) float64 {
	// Rough estimation: average character width is ~0.6 * fontSize
	return float64(len(text)) * fontSize * 0.6
}

func maxSymbolHeight(symbol Symbol) float64 {
	return symbol.Height()
}

func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
