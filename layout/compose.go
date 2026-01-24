package layout

import (
	"fmt"
	"strings"

	"github.com/SCKelemen/svg"
	"github.com/SCKelemen/units"
)

// ChartRenderer is a function that renders a chart to SVG
type ChartRenderer func(bounds Rect) string

// ChartSpec defines a chart and its position in a composition
type ChartSpec struct {
	// Chart renderer function
	Renderer ChartRenderer

	// Position (for custom layout)
	Bounds Rect

	// Grid position (for grid layout)
	Row    int
	Col    int
	RowSpan int
	ColSpan int

	// Title for this chart
	Title string

	// Whether to show axes
	ShowAxes bool

	// Margin around this chart
	Margin Margin

	// Optional ID for referencing
	ID string
}

// NewChartSpec creates a new chart specification
func NewChartSpec(renderer ChartRenderer) *ChartSpec {
	return &ChartSpec{
		Renderer: renderer,
		RowSpan:  1,
		ColSpan:  1,
		ShowAxes: true,
		Margin:   Uniform(units.Px(5)),
	}
}

// WithBounds sets custom bounds for the chart
func (cs *ChartSpec) WithBounds(bounds Rect) *ChartSpec {
	cs.Bounds = bounds
	return cs
}

// WithGridPosition sets the grid position
func (cs *ChartSpec) WithGridPosition(row, col int) *ChartSpec {
	cs.Row = row
	cs.Col = col
	return cs
}

// WithSpan sets the row and column span
func (cs *ChartSpec) WithSpan(rowSpan, colSpan int) *ChartSpec {
	cs.RowSpan = rowSpan
	cs.ColSpan = colSpan
	return cs
}

// WithTitle sets the chart title
func (cs *ChartSpec) WithTitle(title string) *ChartSpec {
	cs.Title = title
	return cs
}

// WithMargin sets the margin
func (cs *ChartSpec) WithMargin(margin Margin) *ChartSpec {
	cs.Margin = margin
	return cs
}

// WithID sets an ID for the chart
func (cs *ChartSpec) WithID(id string) *ChartSpec {
	cs.ID = id
	return cs
}

// CompositionLayout defines how charts are arranged
type CompositionLayout string

const (
	// LayoutGrid uses a grid layout
	LayoutGrid CompositionLayout = "grid"

	// LayoutStack vertically stacks charts
	LayoutStack CompositionLayout = "stack"

	// LayoutCustom uses custom positioning for each chart
	LayoutCustom CompositionLayout = "custom"

	// LayoutDashboard uses a flexible dashboard layout
	LayoutDashboard CompositionLayout = "dashboard"
)

// Composition manages multiple charts in a single visualization
type Composition struct {
	// Canvas dimensions
	Width  units.Length
	Height units.Length

	// Layout strategy
	Layout CompositionLayout

	// Charts to render
	Charts []*ChartSpec

	// Grid configuration (for grid layout)
	Rows int
	Cols int
	Gap  units.Length

	// Global margin
	Margin Margin

	// Global title
	Title string

	// Title height
	TitleHeight units.Length

	// Background color
	Background string

	// Whether to add border around composition
	Border bool
}

// NewComposition creates a new chart composition
func NewComposition(width, height units.Length) *Composition {
	return &Composition{
		Width:       width,
		Height:      height,
		Layout:      LayoutGrid,
		Charts:      make([]*ChartSpec, 0),
		Rows:        1,
		Cols:        1,
		Gap:         units.Px(10),
		Margin:      DefaultMargin(),
		TitleHeight: units.Px(30),
		Background:  "white",
		Border:      false,
	}
}

// WithLayout sets the layout strategy
func (c *Composition) WithLayout(layout CompositionLayout) *Composition {
	c.Layout = layout
	return c
}

// WithGrid sets grid dimensions
func (c *Composition) WithGrid(rows, cols int) *Composition {
	c.Rows = rows
	c.Cols = cols
	return c
}

// WithGap sets the gap between charts
func (c *Composition) WithGap(gap units.Length) *Composition {
	c.Gap = gap
	return c
}

// WithMargin sets the global margin
func (c *Composition) WithMargin(margin Margin) *Composition {
	c.Margin = margin
	return c
}

// WithTitle sets the global title
func (c *Composition) WithTitle(title string) *Composition {
	c.Title = title
	return c
}

// WithBackground sets the background color
func (c *Composition) WithBackground(color string) *Composition {
	c.Background = color
	return c
}

// WithBorder sets whether to show a border
func (c *Composition) WithBorder(border bool) *Composition {
	c.Border = border
	return c
}

// AddChart adds a chart to the composition
func (c *Composition) AddChart(chart *ChartSpec) *Composition {
	c.Charts = append(c.Charts, chart)
	return c
}

// Render renders the composition to SVG
func (c *Composition) Render() string {
	var sb strings.Builder

	// Start SVG
	sb.WriteString(fmt.Sprintf(`<svg viewBox="0 0 %f %f" xmlns="http://www.w3.org/2000/svg">`,
		c.Width.Value, c.Height.Value))
	sb.WriteString("\n")

	// Add background
	if c.Background != "" {
		style := svg.Style{Fill: c.Background}
		sb.WriteString("  ")
		sb.WriteString(svg.Rect(0, 0, c.Width.Value, c.Height.Value, style))
		sb.WriteString("\n")
	}

	// Add border if requested
	if c.Border {
		style := svg.Style{
			Fill:        "none",
			Stroke:      "#ccc",
			StrokeWidth: 1,
		}
		sb.WriteString("  ")
		sb.WriteString(svg.Rect(0, 0, c.Width.Value, c.Height.Value, style))
		sb.WriteString("\n")
	}

	// Calculate content area
	contentArea := ApplyMargin(Rect{
		X:      units.Px(0),
		Y:      units.Px(0),
		Width:  c.Width,
		Height: c.Height,
	}, c.Margin)

	// Reserve space for title if present
	titleY := contentArea.Y.Value
	if c.Title != "" {
		sb.WriteString("  ")
		sb.WriteString(svg.Text(
			c.Title,
			contentArea.X.Value+contentArea.Width.Value/2,
			titleY+c.TitleHeight.Value/2,
			svg.Style{
				TextAnchor: "middle",
				FontSize:   units.Px(18),
				FontWeight: "bold",
			},
		))
		sb.WriteString("\n")
		contentArea.Y = units.Px(contentArea.Y.Value + c.TitleHeight.Value)
		contentArea.Height = units.Px(contentArea.Height.Value - c.TitleHeight.Value)
	}

	// Render charts based on layout
	switch c.Layout {
	case LayoutGrid:
		c.renderGrid(&sb, contentArea)
	case LayoutStack:
		c.renderStack(&sb, contentArea)
	case LayoutCustom:
		c.renderCustom(&sb, contentArea)
	case LayoutDashboard:
		c.renderDashboard(&sb, contentArea)
	default:
		c.renderGrid(&sb, contentArea)
	}

	// Close SVG
	sb.WriteString("</svg>\n")

	return sb.String()
}

// renderGrid renders charts in a grid layout
func (c *Composition) renderGrid(sb *strings.Builder, contentArea Rect) {
	// Auto-calculate grid dimensions if needed
	rows := c.Rows
	cols := c.Cols
	if rows <= 0 && cols <= 0 {
		rows, cols = AutoGrid(len(c.Charts))
	} else if rows <= 0 {
		rows = (len(c.Charts) + cols - 1) / cols
	} else if cols <= 0 {
		cols = (len(c.Charts) + rows - 1) / rows
	}

	// Create grid
	grid := NewGridLayout(contentArea.Width, contentArea.Height, rows, cols)
	grid.SetGap(c.Gap)
	grid.bounds = contentArea

	// Render each chart
	for i, chart := range c.Charts {
		if i >= rows*cols {
			break
		}

		row := i / cols
		col := i % cols

		// Use specified position if available
		if chart.Row >= 0 && chart.Col >= 0 {
			row = chart.Row
			col = chart.Col
		}

		// Get cell bounds (with span support)
		var cellBounds Rect
		if chart.RowSpan > 1 || chart.ColSpan > 1 {
			cellBounds = grid.CellWithSpan(row, col, chart.RowSpan, chart.ColSpan)
		} else {
			cellBounds = grid.Cell(row, col)
		}

		// Apply chart margin
		chartBounds := ApplyMargin(cellBounds, chart.Margin)

		// Reserve space for title if present
		if chart.Title != "" {
			titleHeight := units.Px(20)
			sb.WriteString("    ")
			sb.WriteString(svg.Text(
				chart.Title,
				chartBounds.X.Value+chartBounds.Width.Value/2,
				chartBounds.Y.Value+titleHeight.Value/2,
				svg.Style{
					TextAnchor: "middle",
					FontSize:   units.Px(12),
					FontWeight: "bold",
				},
			))
			sb.WriteString("\n")
			chartBounds.Y = units.Px(chartBounds.Y.Value + titleHeight.Value)
			chartBounds.Height = units.Px(chartBounds.Height.Value - titleHeight.Value)
		}

		// Render chart
		if chart.Renderer != nil {
			chartSVG := chart.Renderer(chartBounds)
			sb.WriteString(chartSVG)
		}
	}
}

// renderStack renders charts in a vertical stack
func (c *Composition) renderStack(sb *strings.Builder, contentArea Rect) {
	if len(c.Charts) == 0 {
		return
	}

	// Calculate height per chart
	totalGap := float64(len(c.Charts)-1) * c.Gap.Value
	heightPerChart := (contentArea.Height.Value - totalGap) / float64(len(c.Charts))

	currentY := contentArea.Y.Value

	for _, chart := range c.Charts {
		chartBounds := Rect{
			X:      contentArea.X,
			Y:      units.Px(currentY),
			Width:  contentArea.Width,
			Height: units.Px(heightPerChart),
		}

		// Apply chart margin
		chartBounds = ApplyMargin(chartBounds, chart.Margin)

		// Reserve space for title if present
		if chart.Title != "" {
			titleHeight := units.Px(20)
			sb.WriteString("    ")
			sb.WriteString(svg.Text(
				chart.Title,
				chartBounds.X.Value+chartBounds.Width.Value/2,
				chartBounds.Y.Value+titleHeight.Value/2,
				svg.Style{
					TextAnchor: "middle",
					FontSize:   units.Px(12),
					FontWeight: "bold",
				},
			))
			sb.WriteString("\n")
			chartBounds.Y = units.Px(chartBounds.Y.Value + titleHeight.Value)
			chartBounds.Height = units.Px(chartBounds.Height.Value - titleHeight.Value)
		}

		// Render chart
		if chart.Renderer != nil {
			chartSVG := chart.Renderer(chartBounds)
			sb.WriteString(chartSVG)
		}

		currentY += heightPerChart + c.Gap.Value
	}
}

// renderCustom renders charts at custom positions
func (c *Composition) renderCustom(sb *strings.Builder, contentArea Rect) {
	for _, chart := range c.Charts {
		chartBounds := chart.Bounds

		// Offset by content area
		chartBounds.X = units.Px(contentArea.X.Value + chartBounds.X.Value)
		chartBounds.Y = units.Px(contentArea.Y.Value + chartBounds.Y.Value)

		// Apply chart margin
		chartBounds = ApplyMargin(chartBounds, chart.Margin)

		// Reserve space for title if present
		if chart.Title != "" {
			titleHeight := units.Px(20)
			sb.WriteString("    ")
			sb.WriteString(svg.Text(
				chart.Title,
				chartBounds.X.Value+chartBounds.Width.Value/2,
				chartBounds.Y.Value+titleHeight.Value/2,
				svg.Style{
					TextAnchor: "middle",
					FontSize:   units.Px(12),
					FontWeight: "bold",
				},
			))
			sb.WriteString("\n")
			chartBounds.Y = units.Px(chartBounds.Y.Value + titleHeight.Value)
			chartBounds.Height = units.Px(chartBounds.Height.Value - titleHeight.Value)
		}

		// Render chart
		if chart.Renderer != nil {
			chartSVG := chart.Renderer(chartBounds)
			sb.WriteString(chartSVG)
		}
	}
}

// renderDashboard renders a flexible dashboard layout
// Charts can specify their own positions and spans
func (c *Composition) renderDashboard(sb *strings.Builder, contentArea Rect) {
	// Group charts by their specified rows
	maxRow := 0
	maxCol := 0
	for _, chart := range c.Charts {
		if chart.Row > maxRow {
			maxRow = chart.Row
		}
		if chart.Col > maxCol {
			maxCol = chart.Col
		}
	}

	rows := maxRow + 1
	cols := maxCol + 1

	if rows <= 0 || cols <= 0 {
		// Fall back to grid layout
		c.renderGrid(sb, contentArea)
		return
	}

	// Create grid
	grid := NewGridLayout(contentArea.Width, contentArea.Height, rows, cols)
	grid.SetGap(c.Gap)
	grid.bounds = contentArea

	// Render each chart at its specified position
	for _, chart := range c.Charts {
		if chart.Row < 0 || chart.Col < 0 {
			continue
		}

		// Get cell bounds (with span support)
		var cellBounds Rect
		if chart.RowSpan > 1 || chart.ColSpan > 1 {
			cellBounds = grid.CellWithSpan(chart.Row, chart.Col, chart.RowSpan, chart.ColSpan)
		} else {
			cellBounds = grid.Cell(chart.Row, chart.Col)
		}

		// Apply chart margin
		chartBounds := ApplyMargin(cellBounds, chart.Margin)

		// Reserve space for title if present
		if chart.Title != "" {
			titleHeight := units.Px(20)
			sb.WriteString("    ")
			sb.WriteString(svg.Text(
				chart.Title,
				chartBounds.X.Value+chartBounds.Width.Value/2,
				chartBounds.Y.Value+titleHeight.Value/2,
				svg.Style{
					TextAnchor: "middle",
					FontSize:   units.Px(12),
					FontWeight: "bold",
				},
			))
			sb.WriteString("\n")
			chartBounds.Y = units.Px(chartBounds.Y.Value + titleHeight.Value)
			chartBounds.Height = units.Px(chartBounds.Height.Value - titleHeight.Value)
		}

		// Render chart
		if chart.Renderer != nil {
			chartSVG := chart.Renderer(chartBounds)
			sb.WriteString(chartSVG)
		}
	}
}

// Quick composition helpers

// GridComposition creates a simple grid composition
func GridComposition(width, height units.Length, rows, cols int, charts ...ChartRenderer) *Composition {
	comp := NewComposition(width, height).
		WithLayout(LayoutGrid).
		WithGrid(rows, cols)

	for _, renderer := range charts {
		comp.AddChart(NewChartSpec(renderer))
	}

	return comp
}

// StackComposition creates a vertically stacked composition
func StackComposition(width, height units.Length, charts ...ChartRenderer) *Composition {
	comp := NewComposition(width, height).WithLayout(LayoutStack)

	for _, renderer := range charts {
		comp.AddChart(NewChartSpec(renderer))
	}

	return comp
}

// DashboardComposition creates a dashboard with specified chart positions
func DashboardComposition(width, height units.Length, charts ...*ChartSpec) *Composition {
	comp := NewComposition(width, height).WithLayout(LayoutDashboard)

	for _, chart := range charts {
		comp.AddChart(chart)
	}

	return comp
}

// SideBySide creates a 1x2 grid with two charts side by side
func SideBySide(width, height units.Length, left, right ChartRenderer) *Composition {
	return GridComposition(width, height, 1, 2, left, right)
}

// TopAndBottom creates a 2x1 grid with two charts stacked
func TopAndBottom(width, height units.Length, top, bottom ChartRenderer) *Composition {
	return GridComposition(width, height, 2, 1, top, bottom)
}

// Quad creates a 2x2 grid with four charts
func Quad(width, height units.Length, topLeft, topRight, bottomLeft, bottomRight ChartRenderer) *Composition {
	return GridComposition(width, height, 2, 2, topLeft, topRight, bottomLeft, bottomRight)
}

// CustomComposition creates a composition with custom chart positions
func CustomComposition(width, height units.Length, charts ...*ChartSpec) *Composition {
	comp := NewComposition(width, height).WithLayout(LayoutCustom)

	for _, chart := range charts {
		comp.AddChart(chart)
	}

	return comp
}

// ChartGroup groups multiple charts under a shared title
type ChartGroup struct {
	Title      string
	Charts     []*ChartSpec
	Layout     CompositionLayout
	Background string
}

// NewChartGroup creates a new chart group
func NewChartGroup(title string) *ChartGroup {
	return &ChartGroup{
		Title:      title,
		Charts:     make([]*ChartSpec, 0),
		Layout:     LayoutStack,
		Background: "",
	}
}

// AddChart adds a chart to the group
func (cg *ChartGroup) AddChart(chart *ChartSpec) *ChartGroup {
	cg.Charts = append(cg.Charts, chart)
	return cg
}

// Render renders the chart group as a composition
func (cg *ChartGroup) Render(bounds Rect) string {
	comp := NewComposition(bounds.Width, bounds.Height).
		WithLayout(cg.Layout).
		WithTitle(cg.Title).
		WithBackground(cg.Background)

	for _, chart := range cg.Charts {
		comp.AddChart(chart)
	}

	svg := comp.Render()

	// Wrap in a group element positioned at bounds
	return fmt.Sprintf(`<g transform="translate(%f,%f)">%s</g>`, bounds.X.Value, bounds.Y.Value, svg)
}
