package charts

import (
	"sort"

	"github.com/SCKelemen/svg"
	"github.com/SCKelemen/units"
)

// TreeNode represents a node in a hierarchical tree structure
type TreeNode struct {
	Name     string
	Value    float64      // Size/weight of this node
	Children []*TreeNode  // Child nodes (nil for leaf nodes)
	Color    string       // Optional custom color
	Metadata map[string]interface{} // Optional metadata
}

// TreemapSpec configures treemap rendering
type TreemapSpec struct {
	Root      *TreeNode
	Width     float64
	Height    float64
	Padding   float64 // Padding between rectangles
	ShowLabels bool
	MinLabelSize float64 // Minimum rectangle size to show label
	ColorScheme []string // Color palette
}

// TreemapRect represents a positioned rectangle in the treemap
type TreemapRect struct {
	X, Y, Width, Height float64
	Node                *TreeNode
	Depth               int
}

// RenderTreemap renders a treemap visualization using squarified algorithm
func RenderTreemap(spec TreemapSpec) string {
	if spec.Root == nil {
		return ""
	}

	// Calculate total value
	total := calculateTreeValue(spec.Root)
	if total == 0 {
		return ""
	}

	// Compute layout
	rects := squarify(spec.Root, 0, 0, spec.Width, spec.Height, spec.Padding, 0)

	// Render rectangles
	var result string

	for _, rect := range rects {
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
		if spec.ShowLabels && rect.Width > spec.MinLabelSize && rect.Height > spec.MinLabelSize {
			labelX := rect.X + rect.Width/2
			labelY := rect.Y + rect.Height/2

			labelStyle := svg.Style{
				FontSize:         calculateFontSize(rect.Width, rect.Height),
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

// calculateTreeValue recursively calculates total value of a tree
func calculateTreeValue(node *TreeNode) float64 {
	if node == nil {
		return 0
	}

	if len(node.Children) == 0 {
		return node.Value
	}

	total := 0.0
	for _, child := range node.Children {
		total += calculateTreeValue(child)
	}

	return total
}

// squarify implements the squarified treemap algorithm
func squarify(node *TreeNode, x, y, width, height, padding float64, depth int) []TreemapRect {
	if node == nil || width <= 0 || height <= 0 {
		return nil
	}

	// Leaf node - return single rectangle
	if len(node.Children) == 0 {
		return []TreemapRect{{
			X:      x + padding,
			Y:      y + padding,
			Width:  width - 2*padding,
			Height: height - 2*padding,
			Node:   node,
			Depth:  depth,
		}}
	}

	// Calculate total value of children
	total := 0.0
	for _, child := range node.Children {
		total += calculateTreeValue(child)
	}

	if total == 0 {
		return nil
	}

	// Sort children by value (descending)
	children := make([]*TreeNode, len(node.Children))
	copy(children, node.Children)
	sort.Slice(children, func(i, j int) bool {
		return calculateTreeValue(children[i]) > calculateTreeValue(children[j])
	})

	// Apply squarified algorithm
	var rects []TreemapRect
	currentX := x
	currentY := y
	remainingWidth := width
	remainingHeight := height

	for len(children) > 0 {
		// Determine orientation
		horizontal := remainingWidth >= remainingHeight

		// Calculate how many children to place in current row/column
		rowValue := 0.0
		rowCount := 0

		for i := 0; i < len(children); i++ {
			childValue := calculateTreeValue(children[i])
			newRowValue := rowValue + childValue

			if i == 0 || improveAspectRatio(rowValue, newRowValue, remainingWidth, remainingHeight, total, horizontal) {
				rowValue = newRowValue
				rowCount++
			} else {
				break
			}
		}

		if rowCount == 0 {
			rowCount = 1
			rowValue = calculateTreeValue(children[0])
		}

		// Layout children in row/column
		if horizontal {
			// Horizontal layout (left to right)
			rowHeight := (rowValue / total) * remainingHeight
			childX := currentX

			for i := 0; i < rowCount && i < len(children); i++ {
				child := children[i]
				childValue := calculateTreeValue(child)
				childWidth := (childValue / rowValue) * remainingWidth

				childRects := squarify(child, childX, currentY, childWidth, rowHeight, padding, depth+1)
				rects = append(rects, childRects...)

				childX += childWidth
			}

			currentY += rowHeight
			remainingHeight -= rowHeight
		} else {
			// Vertical layout (top to bottom)
			rowWidth := (rowValue / total) * remainingWidth
			childY := currentY

			for i := 0; i < rowCount && i < len(children); i++ {
				child := children[i]
				childValue := calculateTreeValue(child)
				childHeight := (childValue / rowValue) * remainingHeight

				childRects := squarify(child, currentX, childY, rowWidth, childHeight, padding, depth+1)
				rects = append(rects, childRects...)

				childY += childHeight
			}

			currentX += rowWidth
			remainingWidth -= rowWidth
		}

		// Remove processed children
		children = children[rowCount:]
		total -= rowValue
	}

	return rects
}

// improveAspectRatio checks if adding another item improves the aspect ratio
func improveAspectRatio(currentValue, newValue, width, height, total float64, horizontal bool) bool {
	if horizontal {
		currentRatio := aspectRatio(currentValue, width, height, total)
		newRatio := aspectRatio(newValue, width, height, total)
		return newRatio < currentRatio
	} else {
		currentRatio := aspectRatio(currentValue, height, width, total)
		newRatio := aspectRatio(newValue, height, width, total)
		return newRatio < currentRatio
	}
}

// aspectRatio calculates the worst aspect ratio in a row/column
func aspectRatio(value, length, breadth, total float64) float64 {
	if value == 0 || total == 0 {
		return 1e10
	}

	rowLength := (value / total) * breadth

	// Worst aspect ratio is max(width/height, height/width)
	ratio1 := length / rowLength
	ratio2 := rowLength / length

	if ratio1 > ratio2 {
		return ratio1
	}
	return ratio2
}

// getDefaultTreemapColor returns default colors based on depth
func getDefaultTreemapColor(depth int) string {
	colors := []string{
		"#3B82F6", "#10B981", "#F59E0B", "#EF4444",
		"#8B5CF6", "#EC4899", "#06B6D4", "#F97316",
	}
	return colors[depth%len(colors)]
}

// calculateFontSize calculates appropriate font size for rectangle
func calculateFontSize(width, height float64) units.Length {
	size := 12.0
	minDim := width
	if height < minDim {
		minDim = height
	}

	if minDim < 40 {
		size = 8
	} else if minDim < 60 {
		size = 10
	} else if minDim < 100 {
		size = 12
	} else if minDim < 150 {
		size = 14
	} else {
		size = 16
	}

	return units.Px(size)
}

// NewTreeNode creates a new tree node
func NewTreeNode(name string, value float64) *TreeNode {
	return &TreeNode{
		Name:     name,
		Value:    value,
		Metadata: make(map[string]interface{}),
	}
}

// AddChild adds a child node to this node
func (n *TreeNode) AddChild(child *TreeNode) *TreeNode {
	n.Children = append(n.Children, child)
	return n
}

// SetColor sets the color for this node
func (n *TreeNode) SetColor(color string) *TreeNode {
	n.Color = color
	return n
}
