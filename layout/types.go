package layout

import "github.com/SCKelemen/units"

// Rect represents a rectangular region with position and dimensions
type Rect struct {
	X      units.Length
	Y      units.Length
	Width  units.Length
	Height units.Length
}

// Margin represents spacing around content (CSS-style)
type Margin struct {
	Top    units.Length
	Right  units.Length
	Bottom units.Length
	Left   units.Length
}

// Padding represents internal spacing (same structure as Margin)
type Padding = Margin

// DefaultMargin returns conventional margins for charts
func DefaultMargin() Margin {
	return Margin{
		Top:    units.Px(20),
		Right:  units.Px(30),
		Bottom: units.Px(40),
		Left:   units.Px(50),
	}
}

// Uniform creates margins with the same value on all sides
func Uniform(value units.Length) Margin {
	return Margin{
		Top:    value,
		Right:  value,
		Bottom: value,
		Left:   value,
	}
}

// Horizontal creates margins with top/bottom and left/right values
func Horizontal(vertical, horizontal units.Length) Margin {
	return Margin{
		Top:    vertical,
		Right:  horizontal,
		Bottom: vertical,
		Left:   horizontal,
	}
}

// ApplyMargin applies margins to a rectangle, returning the content area
func ApplyMargin(bounds Rect, margin Margin) Rect {
	return Rect{
		X:      units.Px(bounds.X.Value + margin.Left.Value),
		Y:      units.Px(bounds.Y.Value + margin.Top.Value),
		Width:  units.Px(bounds.Width.Value - margin.Left.Value - margin.Right.Value),
		Height: units.Px(bounds.Height.Value - margin.Top.Value - margin.Bottom.Value),
	}
}

// Grid represents a grid layout specification
type Grid struct {
	Rows    int
	Cols    int
	Gap     units.Length  // Gap between cells
	RowGap  units.Length  // Override gap for rows
	ColGap  units.Length  // Override gap for columns
}

// Cell represents a single cell in a grid
type Cell struct {
	Row      int  // Starting row (0-indexed)
	Col      int  // Starting column (0-indexed)
	RowSpan  int  // Number of rows to span
	ColSpan  int  // Number of columns to span
	Bounds   Rect // Computed bounds
}

// Position specifies where to place an element
type Position string

const (
	PositionTop         Position = "top"
	PositionBottom      Position = "bottom"
	PositionLeft        Position = "left"
	PositionRight       Position = "right"
	PositionCenter      Position = "center"
	PositionTopLeft     Position = "top-left"
	PositionTopRight    Position = "top-right"
	PositionBottomLeft  Position = "bottom-left"
	PositionBottomRight Position = "bottom-right"
)

// Alignment specifies how to align content
type Alignment string

const (
	AlignStart   Alignment = "start"
	AlignCenter  Alignment = "center"
	AlignEnd     Alignment = "end"
	AlignStretch Alignment = "stretch"
)

// Direction specifies layout direction
type Direction string

const (
	DirectionRow    Direction = "row"
	DirectionColumn Direction = "column"
)

// Layout represents a layout configuration
type Layout struct {
	Bounds    Rect
	Margin    Margin
	Padding   Padding
	Direction Direction
	Gap       units.Length
	Align     Alignment
	Items     []LayoutItem
}

// LayoutItem represents a single item in a layout
type LayoutItem struct {
	Bounds  Rect
	Content interface{} // Render function or data
	Flex    float64     // Flex grow factor
	MinSize units.Length
	MaxSize units.Length
}

// Stack represents a vertical or horizontal stack layout
type Stack struct {
	Direction Direction
	Gap       units.Length
	Align     Alignment
	Items     []StackItem
}

// StackItem represents an item in a stack
type StackItem struct {
	Size    units.Length  // Fixed size, or 0 for flex
	Flex    float64       // Flex grow factor
	Content interface{}   // Render function
}
