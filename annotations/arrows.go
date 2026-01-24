package annotations

import (
	"fmt"

	"github.com/SCKelemen/dataviz/scales"
	"github.com/SCKelemen/svg"
)

// Arrow represents an arrow annotation
type Arrow struct {
	// Start and end positions in data coordinates (if PositionData)
	X1, Y1, X2, Y2 interface{}

	// Pixel positions (if PositionPixel)
	PxX1, PxY1, PxX2, PxY2 float64

	// Positioning mode
	Mode Position

	// Arrow head style
	HeadSize   float64
	HeadStyle  ArrowHeadStyle
	ShowStart  bool // Show arrow head at start
	ShowEnd    bool // Show arrow head at end
	DoubleHead bool // Deprecated: use ShowStart and ShowEnd

	// Curvature for curved arrows (0 = straight)
	Curve float64

	Style AnnotationStyle
}

// ArrowHeadStyle defines the arrow head appearance
type ArrowHeadStyle string

const (
	ArrowHeadTriangle ArrowHeadStyle = "triangle"
	ArrowHeadOpen     ArrowHeadStyle = "open"
	ArrowHeadDiamond  ArrowHeadStyle = "diamond"
	ArrowHeadCircle   ArrowHeadStyle = "circle"
)

// NewArrow creates a new arrow annotation with data positioning
func NewArrow(x1, y1, x2, y2 interface{}) *Arrow {
	return &Arrow{
		X1:        x1,
		Y1:        y1,
		X2:        x2,
		Y2:        y2,
		Mode:      PositionData,
		HeadSize:  8,
		HeadStyle: ArrowHeadTriangle,
		ShowEnd:   true,
		ShowStart: false,
		Style:     DefaultAnnotationStyle(),
	}
}

// NewArrowPixel creates a new arrow with pixel positioning
func NewArrowPixel(x1, y1, x2, y2 float64) *Arrow {
	return &Arrow{
		PxX1:      x1,
		PxY1:      y1,
		PxX2:      x2,
		PxY2:      y2,
		Mode:      PositionPixel,
		HeadSize:  8,
		HeadStyle: ArrowHeadTriangle,
		ShowEnd:   true,
		ShowStart: false,
		Style:     DefaultAnnotationStyle(),
	}
}

// WithHeadSize sets the arrow head size
func (a *Arrow) WithHeadSize(size float64) *Arrow {
	a.HeadSize = size
	return a
}

// WithHeadStyle sets the arrow head style
func (a *Arrow) WithHeadStyle(style ArrowHeadStyle) *Arrow {
	a.HeadStyle = style
	return a
}

// WithDoubleHead makes the arrow double-headed
func (a *Arrow) WithDoubleHead(double bool) *Arrow {
	if double {
		a.ShowStart = true
		a.ShowEnd = true
	}
	return a
}

// WithStartHead controls the start arrow head
func (a *Arrow) WithStartHead(show bool) *Arrow {
	a.ShowStart = show
	return a
}

// WithEndHead controls the end arrow head
func (a *Arrow) WithEndHead(show bool) *Arrow {
	a.ShowEnd = show
	return a
}

// WithCurve sets the curvature
func (a *Arrow) WithCurve(curve float64) *Arrow {
	a.Curve = curve
	return a
}

// WithStyle sets the style
func (a *Arrow) WithStyle(style AnnotationStyle) *Arrow {
	a.Style = style
	return a
}

// Render renders the arrow
func (a *Arrow) Render(xScale, yScale scales.Scale) string {
	// Calculate positions
	var x1, y1, x2, y2 float64

	switch a.Mode {
	case PositionData:
		x1 = xScale.Apply(a.X1).Value
		y1 = yScale.Apply(a.Y1).Value
		x2 = xScale.Apply(a.X2).Value
		y2 = yScale.Apply(a.Y2).Value
	case PositionPixel:
		x1 = a.PxX1
		y1 = a.PxY1
		x2 = a.PxX2
		y2 = a.PxY2
	}

	// Generate marker IDs
	markerID := fmt.Sprintf("arrow-%s-%d", a.HeadStyle, int(a.HeadSize))
	markerStartID := markerID + "-start"
	markerEndID := markerID + "-end"

	// Create markers
	var markers string
	if a.ShowStart {
		markers += a.createMarker(markerStartID, true)
	}
	if a.ShowEnd {
		markers += a.createMarker(markerEndID, false)
	}

	// Build path or line
	var path string
	if a.Curve != 0 {
		// Curved arrow using quadratic bezier
		// Control point is perpendicular to the line
		mx := (x1 + x2) / 2
		my := (y1 + y2) / 2
		dx := x2 - x1
		dy := y2 - y1
		cx := mx + dy*a.Curve
		cy := my - dx*a.Curve
		path = fmt.Sprintf("M %f %f Q %f %f %f %f", x1, y1, cx, cy, x2, y2)
	} else {
		// Straight line
		path = fmt.Sprintf("M %f %f L %f %f", x1, y1, x2, y2)
	}

	// Create path with markers
	pathStyle := svg.Style{
		Stroke:      a.Style.Stroke,
		StrokeWidth: a.Style.StrokeWidth,
		Fill:        "none",
		Opacity:     a.Style.Opacity,
	}

	var markerAttrs string
	if a.ShowStart {
		markerAttrs += fmt.Sprintf(` marker-start="url(#%s)"`, markerStartID)
	}
	if a.ShowEnd {
		markerAttrs += fmt.Sprintf(` marker-end="url(#%s)"`, markerEndID)
	}

	pathSVG := svg.PathWithMarkers(path, pathStyle,
		func() string { if a.ShowStart { return markerStartID } else { return "" } }(),
		"",
		func() string { if a.ShowEnd { return markerEndID } else { return "" } }())

	return markers + pathSVG + "\n"
}

// createMarker creates an arrow head marker
func (a *Arrow) createMarker(id string, reversed bool) string {
	var markerPath string
	var content string
	orient := svg.MarkerOrientAuto
	if reversed {
		orient = svg.MarkerOrientAutoStart
	}

	switch a.HeadStyle {
	case ArrowHeadTriangle:
		if reversed {
			markerPath = fmt.Sprintf("M 0 0 L %f %f L 0 %f z", a.HeadSize, a.HeadSize/2, a.HeadSize)
		} else {
			markerPath = fmt.Sprintf("M 0 0 L %f %f L %f %f z", a.HeadSize, a.HeadSize/2, a.HeadSize, -a.HeadSize/2)
		}
		content = fmt.Sprintf(`<path d="%s" fill="%s" stroke="%s" stroke-width="%.2f"/>`,
			markerPath, a.Style.Fill, a.Style.Stroke, a.Style.StrokeWidth)
	case ArrowHeadOpen:
		markerPath = fmt.Sprintf("M 0 0 L %f %f M 0 0 L %f %f",
			a.HeadSize, a.HeadSize/2, a.HeadSize, -a.HeadSize/2)
		content = fmt.Sprintf(`<path d="%s" fill="none" stroke="%s" stroke-width="%.2f"/>`,
			markerPath, a.Style.Stroke, a.Style.StrokeWidth)
	case ArrowHeadDiamond:
		markerPath = fmt.Sprintf("M 0 0 L %f %f L %f %f L %f %f z",
			a.HeadSize/2, a.HeadSize/2, a.HeadSize, 0.0, a.HeadSize/2, -a.HeadSize/2)
		content = fmt.Sprintf(`<path d="%s" fill="%s" stroke="%s" stroke-width="%.2f"/>`,
			markerPath, a.Style.Fill, a.Style.Stroke, a.Style.StrokeWidth)
	case ArrowHeadCircle:
		markerPath = svg.CirclePath(a.HeadSize/2, 0, a.HeadSize/3)
		content = fmt.Sprintf(`<path d="%s" fill="%s" stroke="%s" stroke-width="%.2f"/>`,
			markerPath, a.Style.Fill, a.Style.Stroke, a.Style.StrokeWidth)
	}

	markerDef := svg.Marker(svg.MarkerDef{
		ID:           id,
		ViewBox:      fmt.Sprintf("0 -%f %f %f", a.HeadSize/2, a.HeadSize, a.HeadSize),
		RefX:         a.HeadSize,
		RefY:         0,
		MarkerWidth:  a.HeadSize,
		MarkerHeight: a.HeadSize,
		Orient:       orient,
		Content:      content,
	})

	return markerDef
}

// Connector draws a line connecting two points (no arrow heads)
type Connector struct {
	X1, Y1, X2, Y2 interface{}
	Mode           Position

	// Line style
	LineStyle ConnectorStyle
	Style     AnnotationStyle
}

// ConnectorStyle defines connector line styles
type ConnectorStyle string

const (
	ConnectorStraight    ConnectorStyle = "straight"
	ConnectorStep        ConnectorStyle = "step"
	ConnectorStepBefore  ConnectorStyle = "step-before"
	ConnectorStepAfter   ConnectorStyle = "step-after"
	ConnectorElbow       ConnectorStyle = "elbow"
)

// NewConnector creates a new connector
func NewConnector(x1, y1, x2, y2 interface{}) *Connector {
	return &Connector{
		X1:        x1,
		Y1:        y1,
		X2:        x2,
		Y2:        y2,
		Mode:      PositionData,
		LineStyle: ConnectorStraight,
		Style:     DefaultAnnotationStyle(),
	}
}

// WithLineStyle sets the connector line style
func (c *Connector) WithLineStyle(style ConnectorStyle) *Connector {
	c.LineStyle = style
	return c
}

// WithStyle sets the style
func (c *Connector) WithStyle(style AnnotationStyle) *Connector {
	c.Style = style
	return c
}

// Render renders the connector
func (c *Connector) Render(xScale, yScale scales.Scale) string {
	var x1, y1, x2, y2 float64

	switch c.Mode {
	case PositionData:
		x1 = xScale.Apply(c.X1).Value
		y1 = yScale.Apply(c.Y1).Value
		x2 = xScale.Apply(c.X2).Value
		y2 = yScale.Apply(c.Y2).Value
	}

	// Build path based on line style
	var path string
	switch c.LineStyle {
	case ConnectorStraight:
		path = fmt.Sprintf("M %f %f L %f %f", x1, y1, x2, y2)
	case ConnectorStep, ConnectorStepAfter:
		mx := (x1 + x2) / 2
		path = fmt.Sprintf("M %f %f H %f V %f H %f", x1, y1, mx, y2, x2)
	case ConnectorStepBefore:
		my := (y1 + y2) / 2
		path = fmt.Sprintf("M %f %f V %f H %f V %f", x1, y1, my, x2, y2)
	case ConnectorElbow:
		path = fmt.Sprintf("M %f %f H %f V %f", x1, y1, x2, y2)
	default:
		path = fmt.Sprintf("M %f %f L %f %f", x1, y1, x2, y2)
	}

	pathStyle := svg.Style{
		Stroke:      c.Style.Stroke,
		StrokeWidth: c.Style.StrokeWidth,
		Fill:        "none",
		Opacity:     c.Style.Opacity,
	}

	if c.Style.StrokeDash != "" {
		// Would need to add StrokeDasharray to svg.Style
	}

	return svg.Path(path, pathStyle) + "\n"
}
