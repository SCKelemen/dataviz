package legends

import (
	"github.com/SCKelemen/color"
	"github.com/SCKelemen/units"
)

// Legend represents a chart legend with items, positioning, and styling
type Legend struct {
	Items    []LegendItem
	Position Position
	Layout   Layout
	Style    *Style
}

// LegendItem represents a single entry in the legend
type LegendItem struct {
	Label  string
	Symbol Symbol
	Value  string // Optional value display (e.g., "45%" for pie slices)
}

// Symbol represents the visual indicator for a legend item
type Symbol interface {
	// Width returns the width of the symbol in pixels
	Width() float64
	// Height returns the height of the symbol in pixels
	Height() float64
	// Render generates the SVG string for this symbol
	Render() string
}

// Position defines where the legend appears in the chart
type Position int

const (
	PositionTopLeft Position = iota
	PositionTopRight
	PositionTopCenter
	PositionBottomLeft
	PositionBottomRight
	PositionBottomCenter
	PositionLeft
	PositionRight
	PositionNone // Hide legend
)

// Layout defines how legend items are arranged
type Layout int

const (
	LayoutVertical Layout = iota // Stack vertically
	LayoutHorizontal              // Flow horizontally
	LayoutGrid                    // Grid layout (future)
	LayoutAuto                    // Auto-detect based on position
)

// Style contains styling configuration for the legend
type Style struct {
	Background   color.Color
	Border       color.Color
	BorderWidth  units.Length
	Padding      units.Length
	ItemSpacing  units.Length
	SymbolSpacing units.Length // Space between symbol and label
	FontSize     units.Length
	FontFamily   string
	TextColor    color.Color
}

// DefaultStyle returns a default legend style
func DefaultStyle() *Style {
	border, _ := color.HexToRGB("#e5e7eb")
	textColor, _ := color.HexToRGB("#374151")
	transparent := color.RGB(0, 0, 0).WithAlpha(0)

	return &Style{
		Background:    transparent,
		Border:        border,
		BorderWidth:   units.Px(1),
		Padding:       units.Px(10),
		ItemSpacing:   units.Px(8),
		SymbolSpacing: units.Px(6),
		FontSize:      units.Px(12),
		FontFamily:    "Arial, sans-serif",
		TextColor:     textColor,
	}
}

// Option is a functional option for configuring a Legend
type Option func(*Legend)

// WithPosition sets the legend position
func WithPosition(pos Position) Option {
	return func(l *Legend) {
		l.Position = pos
	}
}

// WithLayout sets the legend layout
func WithLayout(layout Layout) Option {
	return func(l *Legend) {
		l.Layout = layout
	}
}

// WithStyle sets the legend style
func WithStyle(style *Style) Option {
	return func(l *Legend) {
		l.Style = style
	}
}

// New creates a new Legend with the given items and options
func New(items []LegendItem, opts ...Option) *Legend {
	l := &Legend{
		Items:    items,
		Position: PositionTopRight,
		Layout:   LayoutAuto,
		Style:    DefaultStyle(),
	}

	for _, opt := range opts {
		opt(l)
	}

	// Auto-detect layout based on position if LayoutAuto
	if l.Layout == LayoutAuto {
		switch l.Position {
		case PositionTopCenter, PositionBottomCenter:
			l.Layout = LayoutHorizontal
		default:
			l.Layout = LayoutVertical
		}
	}

	return l
}

// Item creates a new LegendItem
func Item(label string, symbol Symbol) LegendItem {
	return LegendItem{
		Label:  label,
		Symbol: symbol,
	}
}

// ItemWithValue creates a new LegendItem with a value display
func ItemWithValue(label string, symbol Symbol, value string) LegendItem {
	return LegendItem{
		Label:  label,
		Symbol: symbol,
		Value:  value,
	}
}
