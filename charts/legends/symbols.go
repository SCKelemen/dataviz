package legends

import (
	"fmt"
	"strings"

	"github.com/SCKelemen/color"
)

// ColorSwatch represents a colored square symbol
type ColorSwatch struct {
	Color color.Color
	Size  float64 // Size in pixels
}

// NewColorSwatch creates a new ColorSwatch symbol
func NewColorSwatch(c color.Color, size float64) *ColorSwatch {
	if size <= 0 {
		size = 12 // Default size
	}
	return &ColorSwatch{
		Color: c,
		Size:  size,
	}
}

// Width returns the width of the color swatch
func (c *ColorSwatch) Width() float64 {
	return c.Size
}

// Height returns the height of the color swatch
func (c *ColorSwatch) Height() float64 {
	return c.Size
}

// Render generates the SVG for the color swatch
func (c *ColorSwatch) Render() string {
	return fmt.Sprintf(`<rect width="%.1f" height="%.1f" fill="%s" stroke="#000" stroke-width="0.5" stroke-opacity="0.2"/>`,
		c.Size, c.Size, color.RGBToHex(c.Color))
}

// LineSample represents a line segment symbol
type LineSample struct {
	Color       color.Color
	StrokeWidth float64   // Line width in pixels
	Length      float64   // Line length in pixels
	Dash        []float64 // Optional dash pattern
	ShowMarker  bool      // Whether to show a marker in the middle
	MarkerType  string    // Marker type if ShowMarker is true
	MarkerSize  float64   // Marker size if ShowMarker is true
}

// NewLineSample creates a new LineSample symbol
func NewLineSample(c color.Color, width, length float64) *LineSample {
	if width <= 0 {
		width = 2
	}
	if length <= 0 {
		length = 20
	}
	return &LineSample{
		Color:       c,
		StrokeWidth: width,
		Length:      length,
	}
}

// WithDash adds a dash pattern to the line sample
func (l *LineSample) WithDash(pattern []float64) *LineSample {
	l.Dash = pattern
	return l
}

// WithMarker adds a marker to the line sample
func (l *LineSample) WithMarker(markerType string, size float64) *LineSample {
	l.ShowMarker = true
	l.MarkerType = markerType
	l.MarkerSize = size
	return l
}

// Width returns the width of the line sample
func (l *LineSample) Width() float64 {
	return l.Length
}

// Height returns the height of the line sample
func (l *LineSample) Height() float64 {
	return l.SymbolHeight()
}

// Render generates the SVG for the line sample
func (l *LineSample) Render() string {
	var sb strings.Builder

	// Line
	y := l.SymbolHeight() / 2
	sb.WriteString(fmt.Sprintf(`<line x1="0" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-width="%.1f"`,
		y, l.Length, y, color.RGBToHex(l.Color), l.StrokeWidth))

	if len(l.Dash) > 0 {
		dashStr := make([]string, len(l.Dash))
		for i, d := range l.Dash {
			dashStr[i] = fmt.Sprintf("%.1f", d)
		}
		sb.WriteString(fmt.Sprintf(` stroke-dasharray="%s"`, strings.Join(dashStr, ",")))
	}

	sb.WriteString("/>")

	// Marker (if enabled)
	if l.ShowMarker {
		cx := l.Length / 2
		cy := y
		markerSVG := renderMarker(l.MarkerType, cx, cy, l.MarkerSize, l.Color)
		sb.WriteString(markerSVG)
	}

	return sb.String()
}

// SymbolHeight returns the height of the line (renamed to avoid conflict with stroke Width)
func (l *LineSample) SymbolHeight() float64 {
	if l.ShowMarker && l.MarkerSize > l.StrokeWidth {
		return l.MarkerSize
	}
	return l.StrokeWidth
}

// MarkerSymbol represents a marker symbol (circle, square, diamond, etc.)
type MarkerSymbol struct {
	Type  string      // "circle", "square", "diamond", "triangle", "cross", "x", "dot"
	Color color.Color
	Size  float64 // Size in pixels
}

// NewMarkerSymbol creates a new MarkerSymbol
func NewMarkerSymbol(markerType string, c color.Color, size float64) *MarkerSymbol {
	if size <= 0 {
		size = 8
	}
	return &MarkerSymbol{
		Type:  markerType,
		Color: c,
		Size:  size,
	}
}

// Width returns the width of the marker
func (m *MarkerSymbol) Width() float64 {
	return m.Size
}

// Height returns the height of the marker
func (m *MarkerSymbol) Height() float64 {
	return m.Size
}

// Render generates the SVG for the marker
func (m *MarkerSymbol) Render() string {
	cx := m.Size / 2
	cy := m.Size / 2
	return renderMarker(m.Type, cx, cy, m.Size, m.Color)
}

// renderMarker renders a marker at the given position
func renderMarker(markerType string, cx, cy, size float64, c color.Color) string {
	hexColor := color.RGBToHex(c)
	switch markerType {
	case "circle":
		return fmt.Sprintf(`<circle cx="%.1f" cy="%.1f" r="%.1f" fill="%s" stroke="#fff" stroke-width="1"/>`,
			cx, cy, size/2, hexColor)

	case "square":
		half := size / 2
		return fmt.Sprintf(`<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="%s" stroke="#fff" stroke-width="1"/>`,
			cx-half, cy-half, size, size, hexColor)

	case "diamond":
		return fmt.Sprintf(`<polygon points="%.1f,%.1f %.1f,%.1f %.1f,%.1f %.1f,%.1f" fill="%s" stroke="#fff" stroke-width="1"/>`,
			cx, cy-size/2, // top
			cx+size/2, cy, // right
			cx, cy+size/2, // bottom
			cx-size/2, cy, // left
			hexColor)

	case "triangle":
		h := size * 0.866 // sqrt(3)/2 for equilateral triangle
		return fmt.Sprintf(`<polygon points="%.1f,%.1f %.1f,%.1f %.1f,%.1f" fill="%s" stroke="#fff" stroke-width="1"/>`,
			cx, cy-h/2,        // top
			cx+size/2, cy+h/2, // bottom right
			cx-size/2, cy+h/2, // bottom left
			hexColor)

	case "cross":
		return fmt.Sprintf(`<path d="M%.1f,%.1f v%.1f M%.1f,%.1f h%.1f" stroke="%s" stroke-width="2" stroke-linecap="round"/>`,
			cx, cy-size/2, size,   // vertical line
			cx-size/2, cy, size,   // horizontal line
			hexColor)

	case "x":
		offset := size / 2
		return fmt.Sprintf(`<path d="M%.1f,%.1f L%.1f,%.1f M%.1f,%.1f L%.1f,%.1f" stroke="%s" stroke-width="2" stroke-linecap="round"/>`,
			cx-offset, cy-offset, cx+offset, cy+offset, // diagonal \
			cx-offset, cy+offset, cx+offset, cy-offset, // diagonal /
			hexColor)

	case "dot":
		return fmt.Sprintf(`<circle cx="%.1f" cy="%.1f" r="%.1f" fill="%s"/>`,
			cx, cy, size/4, hexColor)

	default: // Default to circle
		return fmt.Sprintf(`<circle cx="%.1f" cy="%.1f" r="%.1f" fill="%s" stroke="#fff" stroke-width="1"/>`,
			cx, cy, size/2, hexColor)
	}
}

// Helper functions for creating common symbols

// Swatch creates a ColorSwatch symbol with default size
func Swatch(c color.Color) Symbol {
	return NewColorSwatch(c, 12)
}

// Line creates a LineSample symbol with default dimensions
func Line(c color.Color) Symbol {
	return NewLineSample(c, 2, 20)
}

// DashedLine creates a dashed LineSample symbol
func DashedLine(c color.Color) Symbol {
	return NewLineSample(c, 2, 20).WithDash([]float64{4, 2})
}

// LineWithMarker creates a LineSample with a marker
func LineWithMarker(c color.Color, markerType string) Symbol {
	return NewLineSample(c, 2, 20).WithMarker(markerType, 6)
}

// Marker creates a MarkerSymbol with default size
func Marker(markerType string, c color.Color) Symbol {
	return NewMarkerSymbol(markerType, c, 10)
}
