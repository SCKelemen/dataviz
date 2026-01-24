package charts

import (
	"fmt"

	"github.com/SCKelemen/svg"
	"github.com/SCKelemen/units"
)

// DendrogramNode represents a node in the dendrogram tree
type DendrogramNode struct {
	Label    string            // Label for leaf nodes
	Height   float64           // Height (distance) at which this cluster was formed
	Children []*DendrogramNode // Child nodes (empty for leaf nodes)
}

// DendrogramSpec configures dendrogram rendering
type DendrogramSpec struct {
	Root         *DendrogramNode
	Width        float64
	Height       float64
	Orientation  string  // "vertical", "horizontal" (default: vertical)
	ShowLabels   bool    // Show leaf labels
	ShowHeights  bool    // Show height scale
	LineWidth    float64 // Width of dendrogram lines (default: 2)
	LineColor    string  // Color of dendrogram lines
	Title        string
}

// RenderDendrogram generates an SVG dendrogram
func RenderDendrogram(spec DendrogramSpec) string {
	if spec.Root == nil {
		return ""
	}

	// Set defaults
	if spec.Orientation == "" {
		spec.Orientation = "vertical"
	}
	if spec.LineWidth == 0 {
		spec.LineWidth = 2
	}
	if spec.LineColor == "" {
		spec.LineColor = "#374151"
	}

	// Calculate margins
	topMargin := 40.0
	bottomMargin := 60.0
	sideMargin := 60.0

	if spec.Orientation == "horizontal" {
		topMargin = 60.0
		bottomMargin = 40.0
	}

	chartWidth := spec.Width - (2 * sideMargin)
	chartHeight := spec.Height - topMargin - bottomMargin

	// Find all leaf nodes and max height
	leaves := collectLeaves(spec.Root)
	maxHeight := findMaxHeight(spec.Root)

	if maxHeight == 0 {
		maxHeight = 1
	}

	// Calculate positions for all nodes
	leafSpacing := chartWidth / float64(len(leaves))
	if spec.Orientation == "horizontal" {
		leafSpacing = chartHeight / float64(len(leaves))
	}

	leafIndex := 0
	positions := make(map[*DendrogramNode]position)
	calculatePositions(spec.Root, &leafIndex, leafSpacing, chartWidth, chartHeight, maxHeight, spec.Orientation, positions)

	var result string

	// Draw title
	if spec.Title != "" {
		titleStyle := svg.Style{
			FontSize:         units.Px(16),
			FontFamily:       "sans-serif",
			FontWeight:       "bold",
			TextAnchor:       svg.TextAnchorMiddle,
			DominantBaseline: svg.DominantBaselineHanging,
		}
		result += svg.Text(spec.Title, spec.Width/2, 10, titleStyle) + "\n"
	}

	// Draw dendrogram lines
	lineStyle := svg.Style{
		Stroke:      spec.LineColor,
		StrokeWidth: spec.LineWidth,
	}

	result += drawDendrogramNode(spec.Root, positions, sideMargin, topMargin, lineStyle, spec.Orientation)

	// Draw labels
	if spec.ShowLabels {
		labelStyle := svg.Style{
			FontSize:   units.Px(10),
			FontFamily: "sans-serif",
		}

		if spec.Orientation == "vertical" {
			labelStyle.TextAnchor = svg.TextAnchorEnd
			labelStyle.DominantBaseline = svg.DominantBaselineMiddle
		} else {
			labelStyle.TextAnchor = svg.TextAnchorStart
			labelStyle.DominantBaseline = svg.DominantBaselineHanging
		}

		for _, leaf := range leaves {
			if pos, ok := positions[leaf]; ok {
				var labelX, labelY float64

				if spec.Orientation == "vertical" {
					labelX = sideMargin + pos.x
					labelY = topMargin + chartHeight + 5
					labelStyle.TextAnchor = svg.TextAnchorMiddle
					labelStyle.DominantBaseline = svg.DominantBaselineHanging
				} else {
					labelX = sideMargin + chartWidth + 5
					labelY = topMargin + pos.y
					labelStyle.TextAnchor = svg.TextAnchorStart
					labelStyle.DominantBaseline = svg.DominantBaselineMiddle
				}

				if leaf.Label != "" {
					result += svg.Text(leaf.Label, labelX, labelY, labelStyle) + "\n"
				}
			}
		}
	}

	// Draw height scale
	if spec.ShowHeights {
		drawHeightScale := func() {
			scaleLabelStyle := svg.Style{
				FontSize:         units.Px(9),
				FontFamily:       "sans-serif",
				Fill:             "#6b7280",
				TextAnchor:       svg.TextAnchorEnd,
				DominantBaseline: svg.DominantBaselineMiddle,
			}

			steps := 5
			for i := 0; i <= steps; i++ {
				height := maxHeight * float64(i) / float64(steps)

				if spec.Orientation == "vertical" {
					y := topMargin + chartHeight*(1-float64(i)/float64(steps))
					result += svg.Text(fmt.Sprintf("%.2f", height), sideMargin-10, y, scaleLabelStyle) + "\n"

					// Draw scale line
					scaleLineStyle := svg.Style{
						Stroke:      "#d1d5db",
						StrokeWidth: 1,
					}
					result += svg.Line(sideMargin-5, y, sideMargin, y, scaleLineStyle) + "\n"
				} else {
					x := sideMargin + chartWidth*float64(i)/float64(steps)
					result += svg.Text(fmt.Sprintf("%.2f", height), x, topMargin-10, scaleLabelStyle) + "\n"

					// Draw scale line
					scaleLineStyle := svg.Style{
						Stroke:      "#d1d5db",
						StrokeWidth: 1,
					}
					result += svg.Line(x, topMargin-5, x, topMargin, scaleLineStyle) + "\n"
				}
			}
		}

		drawHeightScale()
	}

	return result
}

// position stores x, y coordinates for a node
type position struct {
	x float64
	y float64
}

// collectLeaves returns all leaf nodes in the tree
func collectLeaves(node *DendrogramNode) []*DendrogramNode {
	if node == nil {
		return nil
	}

	if len(node.Children) == 0 {
		return []*DendrogramNode{node}
	}

	var leaves []*DendrogramNode
	for _, child := range node.Children {
		leaves = append(leaves, collectLeaves(child)...)
	}
	return leaves
}

// findMaxHeight finds the maximum height in the tree
func findMaxHeight(node *DendrogramNode) float64 {
	if node == nil {
		return 0
	}

	maxH := node.Height
	for _, child := range node.Children {
		childMax := findMaxHeight(child)
		if childMax > maxH {
			maxH = childMax
		}
	}
	return maxH
}

// calculatePositions computes x, y positions for all nodes
func calculatePositions(node *DendrogramNode, leafIndex *int, leafSpacing, chartWidth, chartHeight, maxHeight float64, orientation string, positions map[*DendrogramNode]position) position {
	if node == nil {
		return position{}
	}

	var pos position

	if len(node.Children) == 0 {
		// Leaf node
		if orientation == "vertical" {
			pos.x = float64(*leafIndex)*leafSpacing + leafSpacing/2
			pos.y = chartHeight
		} else {
			pos.x = 0
			pos.y = float64(*leafIndex)*leafSpacing + leafSpacing/2
		}
		*leafIndex++
	} else {
		// Internal node - position at average of children
		var childPositions []position
		for _, child := range node.Children {
			childPos := calculatePositions(child, leafIndex, leafSpacing, chartWidth, chartHeight, maxHeight, orientation, positions)
			childPositions = append(childPositions, childPos)
		}

		// Average position of children
		if orientation == "vertical" {
			sumX := 0.0
			for _, cp := range childPositions {
				sumX += cp.x
			}
			pos.x = sumX / float64(len(childPositions))
			pos.y = chartHeight * (1 - node.Height/maxHeight)
		} else {
			sumY := 0.0
			for _, cp := range childPositions {
				sumY += cp.y
			}
			pos.x = chartWidth * (node.Height / maxHeight)
			pos.y = sumY / float64(len(childPositions))
		}
	}

	positions[node] = pos
	return pos
}

// drawDendrogramNode recursively draws the dendrogram lines
func drawDendrogramNode(node *DendrogramNode, positions map[*DendrogramNode]position, xOffset, yOffset float64, style svg.Style, orientation string) string {
	if node == nil || len(node.Children) == 0 {
		return ""
	}

	var result string
	nodePos := positions[node]

	for _, child := range node.Children {
		childPos := positions[child]

		if orientation == "vertical" {
			// Draw vertical dendrogram
			// Horizontal line at parent height
			x1 := xOffset + nodePos.x
			y1 := yOffset + nodePos.y
			x2 := xOffset + childPos.x
			y2 := yOffset + childPos.y

			// Draw L-shaped connection
			result += svg.Line(x1, y1, x2, y1, style) + "\n" // Horizontal
			result += svg.Line(x2, y1, x2, y2, style) + "\n" // Vertical
		} else {
			// Draw horizontal dendrogram
			x1 := xOffset + nodePos.x
			y1 := yOffset + nodePos.y
			x2 := xOffset + childPos.x
			y2 := yOffset + childPos.y

			// Draw L-shaped connection
			result += svg.Line(x1, y1, x1, y2, style) + "\n" // Vertical
			result += svg.Line(x1, y2, x2, y2, style) + "\n" // Horizontal
		}

		// Recursively draw child
		result += drawDendrogramNode(child, positions, xOffset, yOffset, style, orientation)
	}

	return result
}

// SimpleDendrogram creates a simple dendrogram from a list of clusters
// Each cluster is represented as a list of indices into the labels array
// Heights represent the distance at which clusters were merged
func SimpleDendrogram(labels []string, clusters [][]int, heights []float64, width, height float64) string {
	if len(clusters) != len(heights) {
		return ""
	}

	// Build leaf nodes
	leaves := make([]*DendrogramNode, len(labels))
	for i, label := range labels {
		leaves[i] = &DendrogramNode{
			Label:  label,
			Height: 0,
		}
	}

	// Build tree bottom-up from clusters
	// This is a simplified version - in practice you'd use hierarchical clustering algorithm
	// For now, create a simple two-way merge tree
	if len(leaves) == 0 {
		return ""
	}

	root := leaves[0]
	if len(leaves) > 1 {
		// Create a simple binary tree merging all leaves
		root = mergeDendrogram(leaves, heights)
	}

	spec := DendrogramSpec{
		Root:        root,
		Width:       width,
		Height:      height,
		Orientation: "vertical",
		ShowLabels:  true,
		ShowHeights: true,
	}

	return RenderDendrogram(spec)
}

// mergeDendrogram creates a simple dendrogram by successive merging
func mergeDendrogram(nodes []*DendrogramNode, heights []float64) *DendrogramNode {
	if len(nodes) == 1 {
		return nodes[0]
	}

	// Merge first two nodes
	heightIdx := 0
	if heightIdx >= len(heights) {
		heightIdx = len(heights) - 1
	}
	h := 1.0
	if len(heights) > 0 {
		h = heights[heightIdx]
	}

	merged := &DendrogramNode{
		Height:   h,
		Children: []*DendrogramNode{nodes[0], nodes[1]},
	}

	// Continue with remaining nodes
	remaining := append([]*DendrogramNode{merged}, nodes[2:]...)
	if len(heights) > 1 {
		return mergeDendrogram(remaining, heights[1:])
	}
	return mergeDendrogram(remaining, []float64{h + 1})
}
