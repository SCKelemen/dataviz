package charts

import (
	"fmt"
	"math"
	"sort"

	"github.com/SCKelemen/svg"
	"github.com/SCKelemen/units"
)

// SankeyNode represents a node in the Sankey diagram
type SankeyNode struct {
	ID    string
	Label string
	Color string  // Optional custom color
	X     float64 // Optional: manual X position (0-1 normalized, 0 = auto)
	Y     float64 // Optional: manual Y position (0-1 normalized, 0 = auto)
}

// SankeyLink represents a flow between two nodes
type SankeyLink struct {
	Source string  // Source node ID
	Target string  // Target node ID
	Value  float64 // Flow value (determines link width)
	Color  string  // Optional custom color
}

// SankeySpec configures Sankey diagram rendering
type SankeySpec struct {
	Nodes        []SankeyNode
	Links        []SankeyLink
	Width        float64
	Height       float64
	NodeWidth    float64 // Width of node rectangles (default: 15)
	NodePadding  float64 // Vertical padding between nodes (default: 10)
	DefaultColor string  // Default node color
	ShowLabels   bool    // Show node labels
	Title        string
}

// RenderSankey generates an SVG Sankey diagram
func RenderSankey(spec SankeySpec) string {
	if len(spec.Nodes) == 0 || len(spec.Links) == 0 {
		return ""
	}

	// Set defaults
	if spec.NodeWidth == 0 {
		spec.NodeWidth = 15
	}
	if spec.NodePadding == 0 {
		spec.NodePadding = 10
	}
	if spec.DefaultColor == "" {
		spec.DefaultColor = "#3b82f6"
	}

	// Calculate margins
	margin := 60.0
	chartWidth := spec.Width - (2 * margin)
	chartHeight := spec.Height - (2 * margin)

	// Build node map
	nodeMap := make(map[string]*SankeyNode)
	for i := range spec.Nodes {
		nodeMap[spec.Nodes[i].ID] = &spec.Nodes[i]
	}

	// Calculate node positions using automatic layout
	nodePositions := calculateSankeyLayout(spec.Nodes, spec.Links, chartWidth, chartHeight, spec.NodeWidth, spec.NodePadding)

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

	// Draw links (flows)
	for _, link := range spec.Links {
		sourcePos, sourceExists := nodePositions[link.Source]
		targetPos, targetExists := nodePositions[link.Target]

		if !sourceExists || !targetExists {
			continue
		}

		// Calculate link path (curved flow)
		linkPath := createSankeyLink(
			margin+sourcePos.x+spec.NodeWidth, margin+sourcePos.y+sourcePos.linkOffsets[link.Target],
			margin+targetPos.x, margin+targetPos.y+targetPos.linkOffsets[link.Source],
			link.Value, sourcePos.totalOut,
		)

		// Get link color
		linkColor := link.Color
		if linkColor == "" {
			// Use source node color with transparency
			sourceNode := nodeMap[link.Source]
			if sourceNode.Color != "" {
				linkColor = sourceNode.Color
			} else {
				linkColor = spec.DefaultColor
			}
		}

		linkStyle := svg.Style{
			Fill:        linkColor,
			FillOpacity: 0.4,
			Stroke:      "none",
		}
		result += svg.Path(linkPath, linkStyle) + "\n"
	}

	// Draw nodes
	for _, node := range spec.Nodes {
		pos, exists := nodePositions[node.ID]
		if !exists {
			continue
		}

		// Get node color
		nodeColor := node.Color
		if nodeColor == "" {
			nodeColor = spec.DefaultColor
		}

		// Draw node rectangle
		nodeStyle := svg.Style{
			Fill:        nodeColor,
			Stroke:      "#ffffff",
			StrokeWidth: 1,
		}
		result += svg.Rect(margin+pos.x, margin+pos.y, spec.NodeWidth, pos.height, nodeStyle) + "\n"

		// Draw node label
		if spec.ShowLabels && node.Label != "" {
			labelX := margin + pos.x + spec.NodeWidth + 5
			labelY := margin + pos.y + pos.height/2

			// If node is on right side, put label on left
			if pos.x > chartWidth/2 {
				labelX = margin + pos.x - 5
			}

			labelStyle := svg.Style{
				FontSize:         units.Px(11),
				FontFamily:       "sans-serif",
				DominantBaseline: svg.DominantBaselineMiddle,
			}

			if pos.x > chartWidth/2 {
				labelStyle.TextAnchor = svg.TextAnchorEnd
			} else {
				labelStyle.TextAnchor = svg.TextAnchorStart
			}

			result += svg.Text(node.Label, labelX, labelY, labelStyle) + "\n"
		}
	}

	return result
}

// sankeyNodePosition stores calculated position and dimensions for a node
type sankeyNodePosition struct {
	x           float64
	y           float64
	height      float64
	totalIn     float64
	totalOut    float64
	linkOffsets map[string]float64 // Offset for each connected link
}

// calculateSankeyLayout computes positions for all nodes
func calculateSankeyLayout(nodes []SankeyNode, links []SankeyLink, width, height, nodeWidth, nodePadding float64) map[string]*sankeyNodePosition {
	positions := make(map[string]*sankeyNodePosition)

	// Initialize positions
	for _, node := range nodes {
		positions[node.ID] = &sankeyNodePosition{
			linkOffsets: make(map[string]float64),
		}
	}

	// Calculate total in/out for each node
	for _, link := range links {
		if sourcePos, ok := positions[link.Source]; ok {
			sourcePos.totalOut += link.Value
		}
		if targetPos, ok := positions[link.Target]; ok {
			targetPos.totalIn += link.Value
		}
	}

	// Calculate node heights based on max(totalIn, totalOut)
	maxTotal := 0.0
	for id, pos := range positions {
		total := math.Max(pos.totalIn, pos.totalOut)
		if total > maxTotal {
			maxTotal = total
		}
		pos.height = total

		// Use manual position if specified
		for _, node := range nodes {
			if node.ID == id && node.X > 0 {
				pos.x = node.X * width
			}
			if node.ID == id && node.Y > 0 {
				pos.y = node.Y * height
			}
		}
	}

	// Normalize heights to fit in chart
	if maxTotal > 0 {
		scale := (height - nodePadding*float64(len(nodes))) / maxTotal
		for _, pos := range positions {
			pos.height *= scale
		}
	}

	// Assign nodes to columns based on connectivity (simplified layout)
	nodeColumns := assignNodesToColumns(nodes, links)

	// Calculate X positions based on columns
	numColumns := 0
	for _, col := range nodeColumns {
		if col > numColumns {
			numColumns = col
		}
	}
	numColumns++

	columnWidth := width / float64(numColumns+1)

	for id, col := range nodeColumns {
		if pos := positions[id]; pos.x == 0 { // Only if not manually set
			pos.x = columnWidth * float64(col+1) - nodeWidth/2
		}
	}

	// Calculate Y positions within each column
	columnNodes := make(map[int][]string)
	for id, col := range nodeColumns {
		columnNodes[col] = append(columnNodes[col], id)
	}

	for col, nodeIDs := range columnNodes {
		// Sort nodes by connections
		sort.Slice(nodeIDs, func(i, j int) bool {
			return nodeIDs[i] < nodeIDs[j]
		})

		// Calculate total height needed
		totalHeight := 0.0
		for _, id := range nodeIDs {
			totalHeight += positions[id].height
		}
		totalHeight += nodePadding * float64(len(nodeIDs)-1)

		// Center vertically
		yOffset := (height - totalHeight) / 2

		for _, id := range nodeIDs {
			pos := positions[id]
			if pos.y == 0 { // Only if not manually set
				pos.y = yOffset
				yOffset += pos.height + nodePadding
			}
		}

		_ = col // Use col variable
	}

	// Calculate link offsets for each node
	for id, pos := range positions {
		// Sort outgoing links
		var outLinks []SankeyLink
		for _, link := range links {
			if link.Source == id {
				outLinks = append(outLinks, link)
			}
		}

		yOffset := 0.0
		for _, link := range outLinks {
			linkHeight := (link.Value / pos.totalOut) * pos.height
			pos.linkOffsets[link.Target] = yOffset + linkHeight/2
			yOffset += linkHeight
		}

		// Sort incoming links
		var inLinks []SankeyLink
		for _, link := range links {
			if link.Target == id {
				inLinks = append(inLinks, link)
			}
		}

		yOffset = 0.0
		for _, link := range inLinks {
			linkHeight := (link.Value / pos.totalIn) * pos.height
			pos.linkOffsets[link.Source] = yOffset + linkHeight/2
			yOffset += linkHeight
		}
	}

	return positions
}

// assignNodesToColumns assigns each node to a column based on connectivity
func assignNodesToColumns(nodes []SankeyNode, links []SankeyLink) map[string]int {
	columns := make(map[string]int)

	// Find source nodes (no incoming links)
	hasIncoming := make(map[string]bool)
	for _, link := range links {
		hasIncoming[link.Target] = true
	}

	// Assign column 0 to source nodes
	for _, node := range nodes {
		if !hasIncoming[node.ID] {
			columns[node.ID] = 0
		}
	}

	// Iteratively assign columns based on longest path from sources
	changed := true
	for changed {
		changed = false
		for _, link := range links {
			if sourceCol, ok := columns[link.Source]; ok {
				targetCol, targetExists := columns[link.Target]
				newCol := sourceCol + 1
				if !targetExists || newCol > targetCol {
					columns[link.Target] = newCol
					changed = true
				}
			}
		}
	}

	// Assign unconnected nodes to column 0
	for _, node := range nodes {
		if _, ok := columns[node.ID]; !ok {
			columns[node.ID] = 0
		}
	}

	return columns
}

// createSankeyLink creates a curved path for a Sankey link
func createSankeyLink(x1, y1, x2, y2, value, totalOut float64) string {
	// Calculate control points for bezier curve
	midX := (x1 + x2) / 2

	// Link width is proportional to value
	width := value / totalOut * 50 // Scale factor for visibility
	if width < 1 {
		width = 1
	}

	// Create path with vertical edges and horizontal bezier
	path := fmt.Sprintf("M %.2f %.2f ", x1, y1-width/2)
	path += fmt.Sprintf("C %.2f %.2f %.2f %.2f %.2f %.2f ", midX, y1-width/2, midX, y2-width/2, x2, y2-width/2)
	path += fmt.Sprintf("L %.2f %.2f ", x2, y2+width/2)
	path += fmt.Sprintf("C %.2f %.2f %.2f %.2f %.2f %.2f ", midX, y2+width/2, midX, y1+width/2, x1, y1+width/2)
	path += "Z"

	return path
}
