package legends

import (
	"fmt"
	"strings"

	"github.com/SCKelemen/color"
	"github.com/SCKelemen/layout"
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

	// Build and layout the legend tree
	layoutNode := l.buildLayoutTree()
	constraints := layout.Unconstrained()
	size := layout.LayoutSimple(layoutNode, constraints)

	// Calculate position based on Position setting
	bounds := Bounds{Width: size.Width, Height: size.Height}
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
			size.Width, size.Height, bg, stroke, strokeWidth))
		sb.WriteString("\n")
	}

	// Render items using layout positions
	sb.WriteString(l.renderLayoutTree(layoutNode))

	sb.WriteString("</g>")
	sb.WriteString("\n")

	return sb.String()
}

// renderLayoutTree walks the layout tree and renders legend items
func (l *Legend) renderLayoutTree(node *layout.Node) string {
	var sb strings.Builder

	// Walk the layout tree and render items at their computed positions
	l.walkLayoutTree(node, 0, &sb)

	return sb.String()
}

// walkLayoutTree recursively walks the layout tree
func (l *Legend) walkLayoutTree(node *layout.Node, itemIndex int, sb *strings.Builder) int {
	// If this is a leaf node (no children), it's a symbol or text element
	if len(node.Children) == 0 {
		return itemIndex
	}

	// For container nodes, walk children
	for _, child := range node.Children {
		// Check if this child contains an item (has 3 children: symbol, spacer, text)
		if len(child.Children) == 3 && itemIndex < len(l.Items) {
			// This is an item node
			item := l.Items[itemIndex]
			symbolNode := child.Children[0]
			textNode := child.Children[2]

			// Render symbol at its computed position
			symbolX := child.Rect.X + symbolNode.Rect.X
			symbolY := child.Rect.Y + symbolNode.Rect.Y

			sb.WriteString(fmt.Sprintf(`  <g transform="translate(%.1f,%.1f)">`, symbolX, symbolY))
			sb.WriteString(item.Symbol.Render())
			sb.WriteString("</g>\n")

			// Render text at its computed position
			labelText := item.Label
			if item.Value != "" {
				labelText = fmt.Sprintf("%s (%s)", item.Label, item.Value)
			}

			textX := child.Rect.X + textNode.Rect.X
			textY := child.Rect.Y + textNode.Rect.Y + l.Style.FontSize.Raw()*0.85 // Baseline adjustment

			sb.WriteString(fmt.Sprintf(`  <text x="%.1f" y="%.1f" font-family="%s" font-size="%.1f" fill="%s">%s</text>`,
				textX, textY, l.Style.FontFamily, l.Style.FontSize.Raw(), color.RGBToHex(l.Style.TextColor), labelText))
			sb.WriteString("\n")

			itemIndex++
		} else {
			// Recurse into container
			itemIndex = l.walkLayoutTree(child, itemIndex, sb)
		}
	}

	return itemIndex
}


// calculateBounds calculates the dimensions of the legend using layout engine
func (l *Legend) calculateBounds() Bounds {
	if len(l.Items) == 0 {
		return Bounds{}
	}

	// Build layout tree
	layoutNode := l.buildLayoutTree()

	// Layout with unconstrained size to get natural dimensions
	constraints := layout.Unconstrained()
	size := layout.LayoutSimple(layoutNode, constraints)

	return Bounds{
		Width:  size.Width,
		Height: size.Height,
	}
}

// buildLayoutTree creates a layout.Node tree for the legend
func (l *Legend) buildLayoutTree() *layout.Node {
	padding := l.Style.Padding.Raw()
	itemSpacing := l.Style.ItemSpacing.Raw()
	symbolSpacing := l.Style.SymbolSpacing.Raw()
	fontSize := l.Style.FontSize.Raw()

	// Create item nodes
	itemNodes := make([]*layout.Node, len(l.Items))
	for i, item := range l.Items {
		// Each item is an HStack: [Symbol | spacing | Text]
		symbolNode := layout.Fixed(item.Symbol.Width(), item.Symbol.Height())

		labelText := item.Label
		if item.Value != "" {
			labelText = fmt.Sprintf("%s (%s)", item.Label, item.Value)
		}
		textWidth := estimateTextWidth(labelText, fontSize)
		textNode := layout.Fixed(textWidth, fontSize)

		itemNode := layout.HStack(
			symbolNode,
			layout.Fixed(symbolSpacing, 1), // Spacing between symbol and text
			textNode,
		)

		// Add spacing between items (except last)
		if i > 0 {
			if l.Layout == LayoutHorizontal {
				// Horizontal spacing
				itemNodes[i] = itemNode
				itemNodes[i].Style.Margin = layout.Spacing{
					Left: layout.Px(itemSpacing),
				}
			} else {
				// Vertical spacing
				itemNodes[i] = itemNode
				itemNodes[i].Style.Margin = layout.Spacing{
					Top: layout.Px(itemSpacing),
				}
			}
		} else {
			itemNodes[i] = itemNode
		}
	}

	// Create container based on layout mode
	var container *layout.Node
	if l.Layout == LayoutHorizontal {
		container = layout.HStack(itemNodes...)
	} else {
		container = layout.VStack(itemNodes...)
	}

	// Add padding to container
	container.Style.Padding = layout.Uniform(layout.Px(padding))

	return container
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

