package annotations

import (
	"fmt"
	"strings"

	"github.com/SCKelemen/dataviz/scales"
	"github.com/SCKelemen/svg"
)

// TextLabel represents a text annotation
type TextLabel struct {
	// Text content
	Text string

	// Position in data coordinates (if PositionData)
	X, Y interface{}

	// Position in pixels (if PositionPixel)
	PxX, PxY float64

	// Position as relative (0-1) (if PositionRelative)
	RelX, RelY float64

	// Positioning mode
	Mode Position

	// Offset from position (in pixels)
	OffsetX, OffsetY float64

	// Rotation angle in degrees
	Rotation float64

	// Style
	Style AnnotationStyle
}

// NewTextLabel creates a new text label with data positioning
func NewTextLabel(text string, x, y interface{}) *TextLabel {
	return &TextLabel{
		Text:  text,
		X:     x,
		Y:     y,
		Mode:  PositionData,
		Style: DefaultAnnotationStyle(),
	}
}

// NewTextLabelPixel creates a new text label with pixel positioning
func NewTextLabelPixel(text string, x, y float64) *TextLabel {
	return &TextLabel{
		Text:  text,
		PxX:   x,
		PxY:   y,
		Mode:  PositionPixel,
		Style: DefaultAnnotationStyle(),
	}
}

// NewTextLabelRelative creates a new text label with relative positioning
func NewTextLabelRelative(text string, x, y float64) *TextLabel {
	return &TextLabel{
		Text:  text,
		RelX:  x,
		RelY:  y,
		Mode:  PositionRelative,
		Style: DefaultAnnotationStyle(),
	}
}

// WithOffset sets the offset
func (tl *TextLabel) WithOffset(x, y float64) *TextLabel {
	tl.OffsetX = x
	tl.OffsetY = y
	return tl
}

// WithRotation sets the rotation angle
func (tl *TextLabel) WithRotation(degrees float64) *TextLabel {
	tl.Rotation = degrees
	return tl
}

// WithStyle sets the style
func (tl *TextLabel) WithStyle(style AnnotationStyle) *TextLabel {
	tl.Style = style
	return tl
}

// WithAnchor sets the text anchor
func (tl *TextLabel) WithAnchor(anchor Anchor) *TextLabel {
	tl.Style.TextAnchor = anchor
	return tl
}

// WithBaseline sets the text baseline
func (tl *TextLabel) WithBaseline(baseline Baseline) *TextLabel {
	tl.Style.TextBaseline = baseline
	return tl
}

// Render renders the text label
func (tl *TextLabel) Render(xScale, yScale scales.Scale) string {
	// Calculate position based on mode
	var x, y float64

	switch tl.Mode {
	case PositionData:
		x = xScale.Apply(tl.X).Value
		y = yScale.Apply(tl.Y).Value
	case PositionPixel:
		x = tl.PxX
		y = tl.PxY
	case PositionRelative:
		xRange := xScale.Range()
		yRange := yScale.Range()
		x = xRange[0].Value + tl.RelX*(xRange[1].Value-xRange[0].Value)
		y = yRange[0].Value + tl.RelY*(yRange[1].Value-yRange[0].Value)
	}

	// Apply offset
	x += tl.OffsetX
	y += tl.OffsetY

	// Build style
	style := tl.Style.toSVGStyle()

	// Add rotation transform if needed
	var transform string
	if tl.Rotation != 0 {
		transform = fmt.Sprintf(`transform="rotate(%f %f %f)"`, tl.Rotation, x, y)
	}

	// Render text
	textSVG := svg.Text(tl.Text, x, y, style)

	// Add transform if needed
	if transform != "" {
		return fmt.Sprintf(`<g %s>%s</g>`, transform, textSVG)
	}

	return textSVG + "\n"
}

// MultilineText represents a text annotation with multiple lines
type MultilineText struct {
	Lines []string
	X, Y  interface{}
	Mode  Position

	// Line spacing (multiplier of font size)
	LineSpacing float64

	OffsetX, OffsetY float64
	Style            AnnotationStyle
}

// NewMultilineText creates a new multiline text annotation
func NewMultilineText(lines []string, x, y interface{}) *MultilineText {
	return &MultilineText{
		Lines:       lines,
		X:           x,
		Y:           y,
		Mode:        PositionData,
		LineSpacing: 1.2,
		Style:       DefaultAnnotationStyle(),
	}
}

// WithLineSpacing sets the line spacing
func (mt *MultilineText) WithLineSpacing(spacing float64) *MultilineText {
	mt.LineSpacing = spacing
	return mt
}

// WithStyle sets the style
func (mt *MultilineText) WithStyle(style AnnotationStyle) *MultilineText {
	mt.Style = style
	return mt
}

// Render renders the multiline text
func (mt *MultilineText) Render(xScale, yScale scales.Scale) string {
	// Calculate position
	var x, y float64

	switch mt.Mode {
	case PositionData:
		x = xScale.Apply(mt.X).Value
		y = yScale.Apply(mt.Y).Value
	case PositionPixel:
		// Would need PxX, PxY fields
	case PositionRelative:
		// Would need RelX, RelY fields
	}

	x += mt.OffsetX
	y += mt.OffsetY

	var sb strings.Builder
	lineHeight := mt.Style.FontSize.Value * mt.LineSpacing

	for i, line := range mt.Lines {
		dy := float64(i) * lineHeight
		style := mt.Style.toSVGStyle()
		sb.WriteString(svg.Text(line, x, y+dy, style))
		sb.WriteString("\n")
	}

	return sb.String()
}

// CalloutLabel is a text label with a line pointing to a data point
type CalloutLabel struct {
	Text string
	X, Y interface{} // Data point position

	// Label position offset from point (in pixels)
	LabelOffsetX, LabelOffsetY float64

	// Whether to draw connecting line
	ShowLine bool

	Style AnnotationStyle
}

// NewCalloutLabel creates a new callout label
func NewCalloutLabel(text string, x, y interface{}) *CalloutLabel {
	return &CalloutLabel{
		Text:         text,
		X:            x,
		Y:            y,
		LabelOffsetX: 20,
		LabelOffsetY: -20,
		ShowLine:     true,
		Style:        DefaultAnnotationStyle(),
	}
}

// WithLabelOffset sets the label offset
func (cl *CalloutLabel) WithLabelOffset(x, y float64) *CalloutLabel {
	cl.LabelOffsetX = x
	cl.LabelOffsetY = y
	return cl
}

// WithShowLine sets whether to show the connecting line
func (cl *CalloutLabel) WithShowLine(show bool) *CalloutLabel {
	cl.ShowLine = show
	return cl
}

// WithStyle sets the style
func (cl *CalloutLabel) WithStyle(style AnnotationStyle) *CalloutLabel {
	cl.Style = style
	return cl
}

// Render renders the callout label
func (cl *CalloutLabel) Render(xScale, yScale scales.Scale) string {
	// Point position
	px := xScale.Apply(cl.X).Value
	py := yScale.Apply(cl.Y).Value

	// Label position
	lx := px + cl.LabelOffsetX
	ly := py + cl.LabelOffsetY

	var sb strings.Builder

	// Draw connecting line if enabled
	if cl.ShowLine {
		lineStyle := svg.Style{
			Stroke:      cl.Style.Stroke,
			StrokeWidth: cl.Style.StrokeWidth,
			Opacity:     cl.Style.Opacity,
		}
		sb.WriteString(svg.Line(px, py, lx, ly, lineStyle))
		sb.WriteString("\n")
	}

	// Draw point marker
	markerStyle := svg.Style{
		Fill:   cl.Style.Fill,
		Stroke: cl.Style.Stroke,
	}
	sb.WriteString(svg.Circle(px, py, 3, markerStyle))
	sb.WriteString("\n")

	// Draw label
	textStyle := cl.Style.toSVGStyle()
	sb.WriteString(svg.Text(cl.Text, lx, ly, textStyle))
	sb.WriteString("\n")

	return sb.String()
}
