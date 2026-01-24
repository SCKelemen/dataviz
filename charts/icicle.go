package charts

import (
	"github.com/SCKelemen/svg"
	"github.com/SCKelemen/units"
)

// IcicleSpec configures icicle chart rendering
type IcicleSpec struct {
	Root        *TreeNode
	Width       float64
	Height      float64
	Padding     float64 // Padding between rectangles
	Orientation string  // "vertical" or "horizontal"
	ShowLabels  bool
	ColorScheme []string
}

// IcicleRect represents a positioned rectangle in the icicle chart
type IcicleRect struct {
	X, Y, Width, Height float64
	Node                *TreeNode
	Depth               int
}

// RenderIcicle renders an icicle (partition) chart
func RenderIcicle(spec IcicleSpec) string {
	if spec.Root == nil {
		return ""
	}

	// Default to vertical orientation
	if spec.Orientation == "" {
		spec.Orientation = "vertical"
	}

	// Calculate total value
	total := calculateTreeValue(spec.Root)
	if total == 0 {
		return ""
	}

	// Calculate maximum depth for level sizing
	maxDepth := calculateMaxDepth(spec.Root, 0)
	if maxDepth == 0 {
		maxDepth = 1
	}

	// Compute layout
	var rects []IcicleRect
	if spec.Orientation == "horizontal" {
		rects = icicleLayoutHorizontal(spec.Root, 0, 0, spec.Width, spec.Height, maxDepth, spec.Padding, 0)
	} else {
		rects = icicleLayoutVertical(spec.Root, 0, 0, spec.Width, spec.Height, maxDepth, spec.Padding, 0)
	}

	// Render rectangles
	var result string

	for _, rect := range rects {
		// Skip tiny rectangles
		if rect.Width < 1 || rect.Height < 1 {
			continue
		}

		// Determine color
		color := rect.Node.Color
		if color == "" {
			if len(spec.ColorScheme) > 0 {
				color = spec.ColorScheme[rect.Depth%len(spec.ColorScheme)]
			} else {
				color = getDefaultTreemapColor(rect.Depth)
			}
		}

		// Draw rectangle
		rectStyle := svg.Style{
			Fill:        color,
			Stroke:      "#ffffff",
			StrokeWidth: 2,
			Opacity:     0.8,
		}

		result += svg.Rect(rect.X, rect.Y, rect.Width, rect.Height, rectStyle) + "\n"

		// Draw label if enabled and rectangle is large enough
		minLabelSize := 30.0
		if spec.ShowLabels && rect.Width > minLabelSize && rect.Height > 15 {
			labelX := rect.X + rect.Width/2
			labelY := rect.Y + rect.Height/2

			fontSize := calculateIcicleFontSize(rect.Width, rect.Height)

			labelStyle := svg.Style{
				FontSize:         fontSize,
				FontFamily:       "sans-serif",
				FontWeight:       svg.FontWeightBold,
				Fill:             "#ffffff",
				TextAnchor:       svg.TextAnchorMiddle,
				DominantBaseline: svg.DominantBaselineMiddle,
			}

			result += svg.Text(rect.Node.Name, labelX, labelY, labelStyle) + "\n"
		}
	}

	return result
}

// icicleLayoutVertical computes vertical icicle layout (depth goes top to bottom)
func icicleLayoutVertical(node *TreeNode, x, y, width, height float64, maxDepth int, padding float64, depth int) []IcicleRect {
	if node == nil || width <= 0 || height <= 0 {
		return nil
	}

	// Calculate height for this level
	levelHeight := height / float64(maxDepth+1)

	var rects []IcicleRect

	// Add rectangle for current node
	rects = append(rects, IcicleRect{
		X:      x + padding,
		Y:      y + padding,
		Width:  width - 2*padding,
		Height: levelHeight - 2*padding,
		Node:   node,
		Depth:  depth,
	})

	// If no children, we're done
	if len(node.Children) == 0 {
		return rects
	}

	// Calculate total value of children
	total := 0.0
	for _, child := range node.Children {
		total += calculateTreeValue(child)
	}

	if total == 0 {
		return rects
	}

	// Layout children horizontally at the next level
	currentX := x
	nextY := y + levelHeight

	for _, child := range node.Children {
		childValue := calculateTreeValue(child)
		childWidth := (childValue / total) * width

		childRects := icicleLayoutVertical(
			child,
			currentX,
			nextY,
			childWidth,
			height-levelHeight,
			maxDepth,
			padding,
			depth+1,
		)
		rects = append(rects, childRects...)

		currentX += childWidth
	}

	return rects
}

// icicleLayoutHorizontal computes horizontal icicle layout (depth goes left to right)
func icicleLayoutHorizontal(node *TreeNode, x, y, width, height float64, maxDepth int, padding float64, depth int) []IcicleRect {
	if node == nil || width <= 0 || height <= 0 {
		return nil
	}

	// Calculate width for this level
	levelWidth := width / float64(maxDepth+1)

	var rects []IcicleRect

	// Add rectangle for current node
	rects = append(rects, IcicleRect{
		X:      x + padding,
		Y:      y + padding,
		Width:  levelWidth - 2*padding,
		Height: height - 2*padding,
		Node:   node,
		Depth:  depth,
	})

	// If no children, we're done
	if len(node.Children) == 0 {
		return rects
	}

	// Calculate total value of children
	total := 0.0
	for _, child := range node.Children {
		total += calculateTreeValue(child)
	}

	if total == 0 {
		return rects
	}

	// Layout children vertically at the next level
	currentY := y
	nextX := x + levelWidth

	for _, child := range node.Children {
		childValue := calculateTreeValue(child)
		childHeight := (childValue / total) * height

		childRects := icicleLayoutHorizontal(
			child,
			nextX,
			currentY,
			width-levelWidth,
			childHeight,
			maxDepth,
			padding,
			depth+1,
		)
		rects = append(rects, childRects...)

		currentY += childHeight
	}

	return rects
}

// calculateIcicleFontSize calculates appropriate font size for icicle rectangle
func calculateIcicleFontSize(width, height float64) units.Length {
	size := 10.0
	minDim := width
	if height < minDim {
		minDim = height
	}

	if minDim < 20 {
		size = 8
	} else if minDim < 40 {
		size = 10
	} else if minDim < 80 {
		size = 12
	} else {
		size = 14
	}

	return units.Px(size)
}
