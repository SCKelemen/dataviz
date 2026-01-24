package charts

import (
	"fmt"
	"math"

	"github.com/SCKelemen/svg"
	"github.com/SCKelemen/units"
)

// SunburstSpec configures sunburst chart rendering
type SunburstSpec struct {
	Root         *TreeNode
	Width        float64
	Height       float64
	InnerRadius  float64 // Radius of center hole (0 for no hole)
	ShowLabels   bool
	ColorScheme  []string
	StartAngle   float64 // Starting angle in degrees (0 = top)
}

// SunburstArc represents an arc segment in the sunburst
type SunburstArc struct {
	InnerRadius float64
	OuterRadius float64
	StartAngle  float64 // In radians
	EndAngle    float64 // In radians
	Node        *TreeNode
	Depth       int
}

// RenderSunburst renders a sunburst (radial partition) chart
func RenderSunburst(spec SunburstSpec) string {
	if spec.Root == nil {
		return ""
	}

	// Calculate center and radius
	centerX := spec.Width / 2
	centerY := spec.Height / 2
	maxRadius := math.Min(spec.Width, spec.Height) / 2

	// Calculate total value
	total := calculateTreeValue(spec.Root)
	if total == 0 {
		return ""
	}

	// Convert start angle to radians (0 degrees = top = -π/2 radians)
	startAngleRad := (spec.StartAngle - 90) * math.Pi / 180

	// Calculate maximum depth for radius scaling
	maxDepth := calculateMaxDepth(spec.Root, 0)
	if maxDepth == 0 {
		maxDepth = 1
	}

	// Compute arcs
	arcs := sunburstLayout(
		spec.Root,
		0,
		startAngleRad,
		startAngleRad+2*math.Pi,
		spec.InnerRadius,
		maxRadius,
		maxDepth,
	)

	// Render arcs
	var result string

	for _, arc := range arcs {
		// Determine color
		color := arc.Node.Color
		if color == "" {
			if len(spec.ColorScheme) > 0 {
				color = spec.ColorScheme[arc.Depth%len(spec.ColorScheme)]
			} else {
				color = getDefaultTreemapColor(arc.Depth)
			}
		}

		// Draw arc
		path := buildArcPath(centerX, centerY, arc.InnerRadius, arc.OuterRadius, arc.StartAngle, arc.EndAngle)

		pathStyle := svg.Style{
			Fill:        color,
			Stroke:      "#ffffff",
			StrokeWidth: 2,
			Opacity:     0.8,
		}

		result += svg.Path(path, pathStyle) + "\n"

		// Draw label if enabled
		if spec.ShowLabels {
			// Calculate label position (middle of arc)
			midAngle := (arc.StartAngle + arc.EndAngle) / 2
			midRadius := (arc.InnerRadius + arc.OuterRadius) / 2
			labelX := centerX + midRadius*math.Cos(midAngle)
			labelY := centerY + midRadius*math.Sin(midAngle)

			// Calculate rotation for radial text
			rotationDeg := midAngle * 180 / math.Pi

			labelStyle := svg.Style{
				FontSize:         units.Px(10),
				FontFamily:       "sans-serif",
				FontWeight:       svg.FontWeightBold,
				Fill:             "#ffffff",
				TextAnchor:       svg.TextAnchorMiddle,
				DominantBaseline: svg.DominantBaselineMiddle,
			}

			transform := fmt.Sprintf(`transform="rotate(%.2f %.2f %.2f)"`, rotationDeg, labelX, labelY)
			result += fmt.Sprintf(`<g %s>`, transform)
			result += svg.Text(arc.Node.Name, labelX, labelY, labelStyle)
			result += "</g>\n"
		}
	}

	return result
}

// sunburstLayout computes arc positions for sunburst chart
func sunburstLayout(node *TreeNode, depth int, startAngle, endAngle, innerRadius, outerRadius float64, maxDepth int) []SunburstArc {
	if node == nil {
		return nil
	}

	// Calculate radii for this depth
	depthRange := outerRadius - innerRadius
	radiusPerLevel := depthRange / float64(maxDepth)
	arcInnerRadius := innerRadius + float64(depth)*radiusPerLevel
	arcOuterRadius := arcInnerRadius + radiusPerLevel

	var arcs []SunburstArc

	// Leaf node or no children
	if len(node.Children) == 0 {
		arcs = append(arcs, SunburstArc{
			InnerRadius: arcInnerRadius,
			OuterRadius: arcOuterRadius,
			StartAngle:  startAngle,
			EndAngle:    endAngle,
			Node:        node,
			Depth:       depth,
		})
		return arcs
	}

	// Add arc for current node
	arcs = append(arcs, SunburstArc{
		InnerRadius: arcInnerRadius,
		OuterRadius: arcOuterRadius,
		StartAngle:  startAngle,
		EndAngle:    endAngle,
		Node:        node,
		Depth:       depth,
	})

	// Calculate total value of children
	total := 0.0
	for _, child := range node.Children {
		total += calculateTreeValue(child)
	}

	if total == 0 {
		return arcs
	}

	// Layout children
	currentAngle := startAngle
	angleRange := endAngle - startAngle

	for _, child := range node.Children {
		childValue := calculateTreeValue(child)
		childAngleRange := (childValue / total) * angleRange
		childEndAngle := currentAngle + childAngleRange

		childArcs := sunburstLayout(child, depth+1, currentAngle, childEndAngle, innerRadius, outerRadius, maxDepth)
		arcs = append(arcs, childArcs...)

		currentAngle = childEndAngle
	}

	return arcs
}

// buildArcPath builds an SVG path for an arc segment
func buildArcPath(cx, cy, innerRadius, outerRadius, startAngle, endAngle float64) string {
	// Calculate points
	x1 := cx + innerRadius*math.Cos(startAngle)
	y1 := cy + innerRadius*math.Sin(startAngle)
	x2 := cx + outerRadius*math.Cos(startAngle)
	y2 := cy + outerRadius*math.Sin(startAngle)
	x3 := cx + outerRadius*math.Cos(endAngle)
	y3 := cy + outerRadius*math.Sin(endAngle)
	x4 := cx + innerRadius*math.Cos(endAngle)
	y4 := cy + innerRadius*math.Sin(endAngle)

	// Determine if arc is large (> 180 degrees)
	largeArc := 0
	if endAngle-startAngle > math.Pi {
		largeArc = 1
	}

	// Build path
	// Move to inner start, line to outer start, arc to outer end, line to inner end, arc back to inner start
	path := fmt.Sprintf("M %.2f %.2f L %.2f %.2f A %.2f %.2f 0 %d 1 %.2f %.2f L %.2f %.2f A %.2f %.2f 0 %d 0 %.2f %.2f Z",
		x1, y1, // Move to inner start
		x2, y2, // Line to outer start
		outerRadius, outerRadius, largeArc, x3, y3, // Outer arc
		x4, y4, // Line to inner end
		innerRadius, innerRadius, largeArc, x1, y1, // Inner arc back
	)

	return path
}

// calculateMaxDepth calculates maximum depth of tree
func calculateMaxDepth(node *TreeNode, currentDepth int) int {
	if node == nil {
		return currentDepth
	}

	if len(node.Children) == 0 {
		return currentDepth
	}

	maxChildDepth := currentDepth
	for _, child := range node.Children {
		childDepth := calculateMaxDepth(child, currentDepth+1)
		if childDepth > maxChildDepth {
			maxChildDepth = childDepth
		}
	}

	return maxChildDepth
}
