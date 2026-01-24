package charts

import (
	"math"

	"github.com/SCKelemen/svg"
	"github.com/SCKelemen/units"
)

// CirclePackingSpec configures circle packing rendering
type CirclePackingSpec struct {
	Root        *TreeNode
	Width       float64
	Height      float64
	Padding     float64 // Padding between circles
	ShowLabels  bool
	ColorScheme []string
}

// PackedCircle represents a positioned circle
type PackedCircle struct {
	X, Y, Radius float64
	Node         *TreeNode
	Depth        int
}

// RenderCirclePacking renders a circle packing visualization
func RenderCirclePacking(spec CirclePackingSpec) string {
	if spec.Root == nil {
		return ""
	}

	// Calculate center
	centerX := spec.Width / 2
	centerY := spec.Height / 2
	maxRadius := math.Min(spec.Width, spec.Height) / 2

	// Compute circle packing layout
	circles := packCircles(spec.Root, centerX, centerY, maxRadius, spec.Padding, 0)

	// Render circles
	var result string

	for _, circle := range circles {
		// Determine color
		color := circle.Node.Color
		if color == "" {
			if len(spec.ColorScheme) > 0 {
				color = spec.ColorScheme[circle.Depth%len(spec.ColorScheme)]
			} else {
				color = getDefaultTreemapColor(circle.Depth)
			}
		}

		// Draw circle
		circleStyle := svg.Style{
			Fill:        color,
			Stroke:      "#ffffff",
			StrokeWidth: 2,
			Opacity:     0.7,
		}

		result += svg.Circle(circle.X, circle.Y, circle.Radius, circleStyle) + "\n"

		// Draw label if enabled and circle is large enough
		if spec.ShowLabels && circle.Radius > 20 {
			labelStyle := svg.Style{
				FontSize:         units.Px(10),
				FontFamily:       "sans-serif",
				FontWeight:       svg.FontWeightBold,
				Fill:             "#ffffff",
				TextAnchor:       svg.TextAnchorMiddle,
				DominantBaseline: svg.DominantBaselineMiddle,
			}

			result += svg.Text(circle.Node.Name, circle.X, circle.Y, labelStyle) + "\n"
		}
	}

	return result
}

// packCircles computes circle positions using hierarchical circle packing
func packCircles(node *TreeNode, cx, cy, maxRadius, padding float64, depth int) []PackedCircle {
	if node == nil {
		return nil
	}

	// Leaf node - single circle
	if len(node.Children) == 0 {
		radius := math.Sqrt(node.Value) * maxRadius / 10 // Scale radius based on value
		if radius > maxRadius {
			radius = maxRadius
		}
		return []PackedCircle{{
			X:      cx,
			Y:      cy,
			Radius: radius - padding,
			Node:   node,
			Depth:  depth,
		}}
	}

	var circles []PackedCircle

	// Calculate radii for children based on their values
	childCircles := make([]PackedCircle, len(node.Children))
	totalValue := 0.0
	for _, child := range node.Children {
		totalValue += calculateTreeValue(child)
	}

	for i, child := range node.Children {
		childValue := calculateTreeValue(child)
		// Radius proportional to square root of value (area-based)
		radius := math.Sqrt(childValue/totalValue) * maxRadius * 0.8
		childCircles[i] = PackedCircle{
			Radius: radius,
			Node:   child,
			Depth:  depth + 1,
		}
	}

	// Position circles using simple packing algorithm
	positionedCircles := positionCirclesPacked(childCircles, cx, cy, maxRadius*0.9)

	// Recursively pack children
	for _, circle := range positionedCircles {
		if len(circle.Node.Children) > 0 {
			nestedCircles := packCircles(
				circle.Node,
				circle.X,
				circle.Y,
				circle.Radius,
				padding,
				depth+1,
			)
			circles = append(circles, nestedCircles...)
		} else {
			circles = append(circles, circle)
		}
	}

	return circles
}

// positionCirclesPacked positions circles using a simple packing algorithm
func positionCirclesPacked(circles []PackedCircle, cx, cy, maxRadius float64) []PackedCircle {
	if len(circles) == 0 {
		return circles
	}

	// Sort circles by radius (descending) for better packing
	sortedCircles := make([]PackedCircle, len(circles))
	copy(sortedCircles, circles)

	// Simple circular arrangement for now (can be improved with better packing algorithms)
	if len(sortedCircles) == 1 {
		sortedCircles[0].X = cx
		sortedCircles[0].Y = cy
		return sortedCircles
	}

	// Place first circle at center
	sortedCircles[0].X = cx
	sortedCircles[0].Y = cy

	// Place remaining circles in a circle around the first one
	angleStep := 2 * math.Pi / float64(len(sortedCircles)-1)
	placementRadius := sortedCircles[0].Radius + sortedCircles[1].Radius + 5

	for i := 1; i < len(sortedCircles); i++ {
		angle := float64(i-1) * angleStep
		sortedCircles[i].X = cx + placementRadius*math.Cos(angle)
		sortedCircles[i].Y = cy + placementRadius*math.Sin(angle)
	}

	return sortedCircles
}
