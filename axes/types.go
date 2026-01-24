package axes

import (
	"fmt"
	"time"

	"github.com/SCKelemen/dataviz/scales"
	"github.com/SCKelemen/units"
)

// Axis represents a visualization axis with tick marks and labels.
// Works with any scale type (linear, log, time, band, etc.) to generate
// appropriate tick positions and labels.
//
// Example:
//   scale := scales.NewLinearScale([2]float64{0, 100}, [2]units.Length{units.Px(0), units.Px(500)})
//   axis := NewAxis(scale, AxisOrientationBottom)
//   axis.TickCount(10).Title("Temperature (°C)")
type Axis struct {
	scale       scales.Scale
	orientation AxisOrientation
	title       string
	tickCount   int
	tickSize    units.Length
	tickPadding units.Length
	formatter   TickFormatFunc
	showGrid    bool
	gridLength  units.Length
}

// AxisOrientation specifies where the axis is positioned
type AxisOrientation int

const (
	AxisOrientationTop AxisOrientation = iota
	AxisOrientationBottom
	AxisOrientationLeft
	AxisOrientationRight
)

// String returns the orientation name
func (o AxisOrientation) String() string {
	switch o {
	case AxisOrientationTop:
		return "top"
	case AxisOrientationBottom:
		return "bottom"
	case AxisOrientationLeft:
		return "left"
	case AxisOrientationRight:
		return "right"
	default:
		return "unknown"
	}
}

// IsHorizontal returns true if the axis is horizontal
func (o AxisOrientation) IsHorizontal() bool {
	return o == AxisOrientationTop || o == AxisOrientationBottom
}

// IsVertical returns true if the axis is vertical
func (o AxisOrientation) IsVertical() bool {
	return o == AxisOrientationLeft || o == AxisOrientationRight
}

// TickFormatFunc formats tick values into labels
type TickFormatFunc func(value interface{}) string

// NewAxis creates a new axis with the given scale and orientation
func NewAxis(scale scales.Scale, orientation AxisOrientation) *Axis {
	return &Axis{
		scale:       scale,
		orientation: orientation,
		tickCount:   10,
		tickSize:    units.Px(6),
		tickPadding: units.Px(3),
		formatter:   DefaultTickFormatter,
		showGrid:    false,
		gridLength:  units.Px(0),
	}
}

// Title sets the axis title
func (a *Axis) Title(title string) *Axis {
	a.title = title
	return a
}

// TickCount sets the desired number of ticks
func (a *Axis) TickCount(count int) *Axis {
	a.tickCount = count
	return a
}

// TickSize sets the length of tick marks
func (a *Axis) TickSize(size units.Length) *Axis {
	a.tickSize = size
	return a
}

// TickPadding sets the spacing between ticks and labels
func (a *Axis) TickPadding(padding units.Length) *Axis {
	a.tickPadding = padding
	return a
}

// TickFormat sets the tick label formatter
func (a *Axis) TickFormat(formatter TickFormatFunc) *Axis {
	a.formatter = formatter
	return a
}

// Grid enables grid lines with the specified length
func (a *Axis) Grid(length units.Length) *Axis {
	a.showGrid = true
	a.gridLength = length
	return a
}

// Scale returns the axis scale
func (a *Axis) Scale() scales.Scale {
	return a.scale
}

// Orientation returns the axis orientation
func (a *Axis) Orientation() AxisOrientation {
	return a.orientation
}

// Tick represents a single axis tick with position and label
type Tick struct {
	Value    interface{}   // Domain value
	Position units.Length  // Pixel position along axis
	Label    string        // Formatted label text
}

// Ticks generates tick marks for this axis
func (a *Axis) Ticks() []Tick {
	// Get tick values from scale
	var tickValues []interface{}

	// Special handling for TimeScale
	if timeScale, ok := a.scale.(*scales.TimeScale); ok {
		timeTicks := timeScale.Ticks(a.tickCount)
		tickValues = make([]interface{}, len(timeTicks))
		for i, v := range timeTicks {
			tickValues[i] = v
		}
	} else {
		switch s := a.scale.(type) {
		case scales.ContinuousScale:
			// Continuous scales (Linear, Log, Pow)
			floatTicks := s.Ticks(a.tickCount)
			tickValues = make([]interface{}, len(floatTicks))
			for i, v := range floatTicks {
				tickValues[i] = v
			}

		case scales.CategoricalScale:
			// Categorical scales (Band, Point, Ordinal)
			stringTicks := s.Values()
			tickValues = make([]interface{}, len(stringTicks))
			for i, v := range stringTicks {
				tickValues[i] = v
			}

		default:
			// Unknown scale type
			return nil
		}
	}

	// Convert to Tick structs
	ticks := make([]Tick, len(tickValues))
	for i, value := range tickValues {
		ticks[i] = Tick{
			Value:    value,
			Position: a.scale.Apply(value),
			Label:    a.formatter(value),
		}
	}

	return ticks
}

// DefaultTickFormatter provides basic formatting for tick labels
func DefaultTickFormatter(value interface{}) string {
	switch v := value.(type) {
	case float64:
		// Format floats with appropriate precision
		if v == float64(int64(v)) {
			return fmt.Sprintf("%d", int64(v))
		}
		return fmt.Sprintf("%.2f", v)

	case int:
		return fmt.Sprintf("%d", v)

	case time.Time:
		// Format time values
		return v.Format("2006-01-02")

	case string:
		return v

	default:
		return fmt.Sprintf("%v", v)
	}
}

// TimeTickFormatter creates a formatter for time values with custom format
func TimeTickFormatter(format string) TickFormatFunc {
	return func(value interface{}) string {
		if t, ok := value.(time.Time); ok {
			return t.Format(format)
		}
		return fmt.Sprintf("%v", value)
	}
}

// NumberTickFormatter creates a formatter for numbers with custom precision
func NumberTickFormatter(precision int) TickFormatFunc {
	formatStr := fmt.Sprintf("%%.%df", precision)
	return func(value interface{}) string {
		switch v := value.(type) {
		case float64:
			return fmt.Sprintf(formatStr, v)
		case int:
			return fmt.Sprintf("%d", v)
		default:
			return fmt.Sprintf("%v", value)
		}
	}
}

// SITickFormatter formats numbers with SI prefixes (k, M, G, etc.)
func SITickFormatter(value interface{}) string {
	v, ok := value.(float64)
	if !ok {
		if i, ok := value.(int); ok {
			v = float64(i)
		} else {
			return fmt.Sprintf("%v", value)
		}
	}

	if v == 0 {
		return "0"
	}

	absV := v
	if absV < 0 {
		absV = -absV
	}

	prefixes := []struct {
		threshold float64
		suffix    string
	}{
		{1e12, "T"},
		{1e9, "G"},
		{1e6, "M"},
		{1e3, "k"},
		{1, ""},
		{1e-3, "m"},
		{1e-6, "μ"},
		{1e-9, "n"},
		{1e-12, "p"},
	}

	for _, p := range prefixes {
		if absV >= p.threshold {
			scaled := v / p.threshold
			if scaled == float64(int64(scaled)) {
				return fmt.Sprintf("%d%s", int64(scaled), p.suffix)
			}
			return fmt.Sprintf("%.1f%s", scaled, p.suffix)
		}
	}

	return fmt.Sprintf("%.2e", v)
}
