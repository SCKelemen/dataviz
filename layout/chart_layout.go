package layout

import (
	"fmt"
	"strings"

	"github.com/SCKelemen/layout"
	"github.com/SCKelemen/svg"
)

// ChartGrid creates a CSS Grid layout for charts
// Uses fr units to create equally-sized cells
func ChartGrid(rows, cols int) *layout.Node {
	node := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayGrid,
		},
	}

	// Create row template (all equal fr units)
	rowTracks := make([]layout.GridTrack, rows)
	for i := range rowTracks {
		rowTracks[i] = layout.FractionTrack(1.0)
	}
	node.Style.GridTemplateRows = rowTracks

	// Create column template (all equal fr units)
	colTracks := make([]layout.GridTrack, cols)
	for i := range colTracks {
		colTracks[i] = layout.FractionTrack(1.0)
	}
	node.Style.GridTemplateColumns = colTracks

	return node
}

// ChartGridWithGap creates a grid with gap between cells
func ChartGridWithGap(rows, cols int, gap float64) *layout.Node {
	node := ChartGrid(rows, cols)
	node.Style.GridGap = layout.Px(gap)
	return node
}

// ChartGridCustom creates a grid with custom track sizes
func ChartGridCustom(rowTracks, colTracks []layout.GridTrack) *layout.Node {
	return &layout.Node{
		Style: layout.Style{
			Display:             layout.DisplayGrid,
			GridTemplateRows:    rowTracks,
			GridTemplateColumns: colTracks,
		},
	}
}

// ChartHStack creates a horizontal flexbox stack of charts
func ChartHStack() *layout.Node {
	return layout.HStack()
}

// ChartVStack creates a vertical flexbox stack of charts
func ChartVStack() *layout.Node {
	return layout.VStack()
}

// ChartCell creates a chart cell with specific grid positioning
func ChartCell(row, col, rowSpan, colSpan int) *layout.Node {
	return &layout.Node{
		Style: layout.Style{
			GridRowStart:    row,
			GridRowEnd:      row + rowSpan,
			GridColumnStart: col,
			GridColumnEnd:   col + colSpan,
		},
	}
}

// WithMargin adds margin to a chart node
func WithMargin(node *layout.Node, margin float64) *layout.Node {
	node.Style.Margin = layout.Spacing{
		Top:    layout.Px(margin),
		Right:  layout.Px(margin),
		Bottom: layout.Px(margin),
		Left:   layout.Px(margin),
	}
	return node
}

// WithPadding adds padding to a chart node
func WithPadding(node *layout.Node, padding float64) *layout.Node {
	node.Style.Padding = layout.Spacing{
		Top:    layout.Px(padding),
		Right:  layout.Px(padding),
		Bottom: layout.Px(padding),
		Left:   layout.Px(padding),
	}
	return node
}

// WithCustomMargin adds custom margin to each side
func WithCustomMargin(node *layout.Node, top, right, bottom, left float64) *layout.Node {
	node.Style.Margin = layout.Spacing{
		Top:    layout.Px(top),
		Right:  layout.Px(right),
		Bottom: layout.Px(bottom),
		Left:   layout.Px(left),
	}
	return node
}

// WithCustomPadding adds custom padding to each side
func WithCustomPadding(node *layout.Node, top, right, bottom, left float64) *layout.Node {
	node.Style.Padding = layout.Spacing{
		Top:    layout.Px(top),
		Right:  layout.Px(right),
		Bottom: layout.Px(bottom),
		Left:   layout.Px(left),
	}
	return node
}

// WithSize sets explicit width and height
func WithSize(node *layout.Node, width, height float64) *layout.Node {
	node.Style.Width = layout.Px(width)
	node.Style.Height = layout.Px(height)
	return node
}

// WithFlexGrow sets the flex grow factor
func WithFlexGrow(node *layout.Node, grow float64) *layout.Node {
	node.Style.FlexGrow = grow
	return node
}

// Dashboard creates a flexible dashboard layout
// Allows charts to specify their grid position and span
type Dashboard struct {
	Width  float64
	Height float64
	Gap    float64
	Charts []*ChartNode
}

// NewDashboard creates a new dashboard
func NewDashboard(width, height float64) *Dashboard {
	return &Dashboard{
		Width:  width,
		Height: height,
		Gap:    10,
		Charts: make([]*ChartNode, 0),
	}
}

// WithGap sets the gap between charts
func (d *Dashboard) WithGap(gap float64) *Dashboard {
	d.Gap = gap
	return d
}

// AddChart adds a chart to the dashboard
func (d *Dashboard) AddChart(chart *ChartNode) *Dashboard {
	d.Charts = append(d.Charts, chart)
	return d
}

// Layout computes the layout and returns the root node
func (d *Dashboard) Layout() *layout.Node {
	// Find max row and col to determine grid size
	maxRow := 0
	maxCol := 0
	for _, chart := range d.Charts {
		if chart.Style.GridRowEnd > maxRow {
			maxRow = chart.Style.GridRowEnd
		}
		if chart.Style.GridColumnEnd > maxCol {
			maxCol = chart.Style.GridColumnEnd
		}
	}

	rows := maxRow
	cols := maxCol

	if rows == 0 || cols == 0 {
		// Auto-calculate grid dimensions
		n := len(d.Charts)
		cols = int(float64(n) + 0.5)
		if cols == 0 {
			cols = 1
		}
		rows = (n + cols - 1) / cols
	}

	// Create grid
	root := ChartGridWithGap(rows, cols, d.Gap)
	root.Style.Width = layout.Px(d.Width)
	root.Style.Height = layout.Px(d.Height)

	// Add charts
	for _, chart := range d.Charts {
		root = root.AddChild(chart.Node)
	}

	// Compute layout
	constraints := layout.Loose(d.Width, d.Height)
	ctx := layout.NewLayoutContext(d.Width, d.Height, 16) // 16pt default font size
	layout.Layout(root, constraints, ctx)

	return root
}

// Render renders the dashboard to SVG
func (d *Dashboard) Render() string {
	_ = d.Layout() // Compute layout

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg viewBox="0 0 %f %f" xmlns="http://www.w3.org/2000/svg">`,
		d.Width, d.Height))
	sb.WriteString("\n")

	// Render each chart
	for _, chart := range d.Charts {
		if chart.Renderer != nil {
			chartSVG := chart.Renderer(chart.Node)
			sb.WriteString(chartSVG)
		}
	}

	sb.WriteString("</svg>\n")
	return sb.String()
}

// SideBySideLayout creates a 1x2 grid layout
func SideBySideLayout(width, height float64) *layout.Node {
	return ChartGridWithGap(1, 2, 10).
		WithWidth(width).
		WithHeight(height)
}

// TopBottomLayout creates a 2x1 grid layout
func TopBottomLayout(width, height float64) *layout.Node {
	return ChartGridWithGap(2, 1, 10).
		WithWidth(width).
		WithHeight(height)
}

// QuadLayout creates a 2x2 grid layout
func QuadLayout(width, height float64) *layout.Node {
	return ChartGridWithGap(2, 2, 10).
		WithWidth(width).
		WithHeight(height)
}

// FacetLayout creates a grid layout for faceted plots
type FacetLayout struct {
	Rows   int
	Cols   int
	Gap    float64
	Width  float64
	Height float64
}

// NewFacetLayout creates a new facet layout
func NewFacetLayout(rows, cols int, width, height float64) *FacetLayout {
	return &FacetLayout{
		Rows:   rows,
		Cols:   cols,
		Gap:    10,
		Width:  width,
		Height: height,
	}
}

// Build creates the layout node
func (fl *FacetLayout) Build() *layout.Node {
	root := ChartGridWithGap(fl.Rows, fl.Cols, fl.Gap)
	root.Style.Width = layout.Px(fl.Width)
	root.Style.Height = layout.Px(fl.Height)
	return root
}

// Helper to render chart nodes to SVG
func RenderChartTree(root *layout.Node, width, height float64, renderFunc func(*layout.Node) string) string {
	// Compute layout
	constraints := layout.Loose(width, height)
	ctx := layout.NewLayoutContext(width, height, 16)
	layout.Layout(root, constraints, ctx)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg viewBox="0 0 %f %f" xmlns="http://www.w3.org/2000/svg">`,
		width, height))
	sb.WriteString("\n")

	// Render content
	content := renderFunc(root)
	sb.WriteString(content)

	sb.WriteString("</svg>\n")
	return sb.String()
}

// TraverseAndRender walks the tree and renders each chart node
func TraverseAndRender(node *layout.Node) string {
	var sb strings.Builder

	// Create a group with transform
	if node.Rect.X != 0 || node.Rect.Y != 0 {
		sb.WriteString(fmt.Sprintf(`<g transform="translate(%f,%f)">`, node.Rect.X, node.Rect.Y))
		sb.WriteString("\n")
	}

	// Draw debug rect to show layout
	if false { // Enable for debugging
		sb.WriteString(svg.Rect(0, 0, node.Rect.Width, node.Rect.Height, svg.Style{
			Fill:        "none",
			Stroke:      "#ccc",
			StrokeWidth: 1,
		}))
		sb.WriteString("\n")
	}

	// Render children
	for _, child := range node.Children {
		sb.WriteString(TraverseAndRender(child))
	}

	if node.Rect.X != 0 || node.Rect.Y != 0 {
		sb.WriteString("</g>\n")
	}

	return sb.String()
}
