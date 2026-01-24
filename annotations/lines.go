package annotations

import (
	"github.com/SCKelemen/dataviz/scales"
	"github.com/SCKelemen/svg"
)

// ReferenceLine represents a horizontal or vertical reference line
type ReferenceLine struct {
	// Value in data coordinates
	Value interface{}

	// Orientation
	Orientation Orientation

	// Label (optional)
	Label string

	// Label position (0-1 along the line)
	LabelPosition float64

	// Label alignment
	LabelAnchor Anchor

	Style AnnotationStyle
}

// Orientation defines line orientation
type Orientation string

const (
	OrientationHorizontal Orientation = "horizontal"
	OrientationVertical   Orientation = "vertical"
)

// NewHLine creates a new horizontal reference line
func NewHLine(yValue interface{}) *ReferenceLine {
	return &ReferenceLine{
		Value:         yValue,
		Orientation:   OrientationHorizontal,
		LabelPosition: 0.95,
		LabelAnchor:   AnchorEnd,
		Style:         DefaultAnnotationStyle(),
	}
}

// NewVLine creates a new vertical reference line
func NewVLine(xValue interface{}) *ReferenceLine {
	return &ReferenceLine{
		Value:         xValue,
		Orientation:   OrientationVertical,
		LabelPosition: 0.05,
		LabelAnchor:   AnchorStart,
		Style:         DefaultAnnotationStyle(),
	}
}

// WithLabel sets the label
func (rl *ReferenceLine) WithLabel(label string) *ReferenceLine {
	rl.Label = label
	return rl
}

// WithLabelPosition sets the label position (0-1)
func (rl *ReferenceLine) WithLabelPosition(pos float64) *ReferenceLine {
	rl.LabelPosition = pos
	return rl
}

// WithLabelAnchor sets the label anchor
func (rl *ReferenceLine) WithLabelAnchor(anchor Anchor) *ReferenceLine {
	rl.LabelAnchor = anchor
	return rl
}

// WithStyle sets the style
func (rl *ReferenceLine) WithStyle(style AnnotationStyle) *ReferenceLine {
	rl.Style = style
	return rl
}

// WithDashed makes the line dashed
func (rl *ReferenceLine) WithDashed() *ReferenceLine {
	rl.Style.StrokeDash = "5,5"
	return rl
}

// Render renders the reference line
func (rl *ReferenceLine) Render(xScale, yScale scales.Scale) string {
	xRange := xScale.Range()
	yRange := yScale.Range()

	x1 := xRange[0].Value
	x2 := xRange[1].Value
	y1 := yRange[0].Value
	y2 := yRange[1].Value

	var lineX1, lineY1, lineX2, lineY2 float64
	var labelX, labelY float64

	if rl.Orientation == OrientationHorizontal {
		// Horizontal line
		y := yScale.Apply(rl.Value).Value
		lineX1 = x1
		lineY1 = y
		lineX2 = x2
		lineY2 = y

		// Label position
		labelX = x1 + rl.LabelPosition*(x2-x1)
		labelY = y - 5 // Offset above line
	} else {
		// Vertical line
		x := xScale.Apply(rl.Value).Value
		lineX1 = x
		lineY1 = y1
		lineX2 = x
		lineY2 = y2

		// Label position
		labelX = x + 5 // Offset right of line
		labelY = y1 + rl.LabelPosition*(y2-y1)
	}

	lineStyle := svg.Style{
		Stroke:      rl.Style.Stroke,
		StrokeWidth: rl.Style.StrokeWidth,
		Opacity:     rl.Style.Opacity,
	}

	var result string
	result += svg.Line(lineX1, lineY1, lineX2, lineY2, lineStyle)
	result += "\n"

	// Add label if present
	if rl.Label != "" {
		textStyle := rl.Style.toSVGStyle()
		textStyle.TextAnchor = svg.TextAnchor(rl.LabelAnchor)
		result += svg.Text(rl.Label, labelX, labelY, textStyle)
		result += "\n"
	}

	return result
}

// ReferenceRegion represents a shaded rectangular region
type ReferenceRegion struct {
	// Data coordinates for region bounds
	X1, Y1, X2, Y2 interface{}

	// Use nil for unbounded dimensions (spans full range)
	Mode Position

	// Label (optional)
	Label string

	// Label position within region (0-1, 0-1)
	LabelX, LabelY float64

	Style AnnotationStyle
}

// NewReferenceRegion creates a new reference region
func NewReferenceRegion(x1, y1, x2, y2 interface{}) *ReferenceRegion {
	return &ReferenceRegion{
		X1:     x1,
		Y1:     y1,
		X2:     x2,
		Y2:     y2,
		Mode:   PositionData,
		LabelX: 0.5,
		LabelY: 0.5,
		Style:  DefaultAnnotationStyle(),
	}
}

// NewHRegion creates a horizontal band (full width, bounded Y)
func NewHRegion(y1, y2 interface{}) *ReferenceRegion {
	return &ReferenceRegion{
		X1:     nil,
		Y1:     y1,
		X2:     nil,
		Y2:     y2,
		Mode:   PositionData,
		LabelX: 0.5,
		LabelY: 0.5,
		Style:  DefaultAnnotationStyle(),
	}
}

// NewVRegion creates a vertical band (full height, bounded X)
func NewVRegion(x1, x2 interface{}) *ReferenceRegion {
	return &ReferenceRegion{
		X1:     x1,
		Y1:     nil,
		X2:     x2,
		Y2:     nil,
		Mode:   PositionData,
		LabelX: 0.5,
		LabelY: 0.5,
		Style:  DefaultAnnotationStyle(),
	}
}

// WithLabel sets the label
func (rr *ReferenceRegion) WithLabel(label string) *ReferenceRegion {
	rr.Label = label
	return rr
}

// WithLabelPosition sets the label position (0-1, 0-1)
func (rr *ReferenceRegion) WithLabelPosition(x, y float64) *ReferenceRegion {
	rr.LabelX = x
	rr.LabelY = y
	return rr
}

// WithStyle sets the style
func (rr *ReferenceRegion) WithStyle(style AnnotationStyle) *ReferenceRegion {
	rr.Style = style
	return rr
}

// Render renders the reference region
func (rr *ReferenceRegion) Render(xScale, yScale scales.Scale) string {
	xRange := xScale.Range()
	yRange := yScale.Range()

	// Calculate bounds
	var x1, y1, x2, y2 float64

	if rr.X1 == nil {
		x1 = xRange[0].Value
	} else {
		x1 = xScale.Apply(rr.X1).Value
	}

	if rr.X2 == nil {
		x2 = xRange[1].Value
	} else {
		x2 = xScale.Apply(rr.X2).Value
	}

	if rr.Y1 == nil {
		y1 = yRange[0].Value
	} else {
		y1 = yScale.Apply(rr.Y1).Value
	}

	if rr.Y2 == nil {
		y2 = yRange[1].Value
	} else {
		y2 = yScale.Apply(rr.Y2).Value
	}

	// Ensure proper ordering
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}

	width := x2 - x1
	height := y2 - y1

	rectStyle := svg.Style{
		Fill:        rr.Style.Fill,
		Opacity:     rr.Style.Opacity,
		Stroke:      rr.Style.Stroke,
		StrokeWidth: rr.Style.StrokeWidth,
	}

	var result string
	result += svg.Rect(x1, y1, width, height, rectStyle)
	result += "\n"

	// Add label if present
	if rr.Label != "" {
		labelX := x1 + rr.LabelX*width
		labelY := y1 + rr.LabelY*height

		textStyle := rr.Style.toSVGStyle()
		textStyle.TextAnchor = svg.TextAnchor(AnchorMiddle)
		result += svg.Text(rr.Label, labelX, labelY, textStyle)
		result += "\n"
	}

	return result
}

// Grid represents a reference grid
type Grid struct {
	// Show horizontal and vertical grid lines
	ShowHorizontal bool
	ShowVertical   bool

	// Number of grid lines (0 = auto based on scale ticks)
	HorizontalCount int
	VerticalCount   int

	Style AnnotationStyle
}

// NewGrid creates a new grid
func NewGrid() *Grid {
	style := DefaultAnnotationStyle()
	style.Stroke = "#e0e0e0"
	style.StrokeWidth = 0.5
	style.Opacity = 0.5

	return &Grid{
		ShowHorizontal:  true,
		ShowVertical:    true,
		HorizontalCount: 0,
		VerticalCount:   0,
		Style:           style,
	}
}

// WithStyle sets the style
func (g *Grid) WithStyle(style AnnotationStyle) *Grid {
	g.Style = style
	return g
}

// WithCounts sets the grid line counts
func (g *Grid) WithCounts(horizontal, vertical int) *Grid {
	g.HorizontalCount = horizontal
	g.VerticalCount = vertical
	return g
}

// Render renders the grid
func (g *Grid) Render(xScale, yScale scales.Scale) string {
	var result string

	xRange := xScale.Range()
	yRange := yScale.Range()
	x1 := xRange[0].Value
	x2 := xRange[1].Value
	y1 := yRange[0].Value
	y2 := yRange[1].Value

	lineStyle := svg.Style{
		Stroke:      g.Style.Stroke,
		StrokeWidth: g.Style.StrokeWidth,
		Opacity:     g.Style.Opacity,
	}

	// Vertical grid lines
	if g.ShowVertical {
		count := g.VerticalCount
		if count == 0 {
			count = 10 // Default
		}

		for i := 0; i <= count; i++ {
			t := float64(i) / float64(count)
			x := x1 + t*(x2-x1)
			result += svg.Line(x, y1, x, y2, lineStyle)
			result += "\n"
		}
	}

	// Horizontal grid lines
	if g.ShowHorizontal {
		count := g.HorizontalCount
		if count == 0 {
			count = 10 // Default
		}

		for i := 0; i <= count; i++ {
			t := float64(i) / float64(count)
			y := y1 + t*(y2-y1)
			result += svg.Line(x1, y, x2, y, lineStyle)
			result += "\n"
		}
	}

	return result
}
