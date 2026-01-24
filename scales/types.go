package scales

import (
	"time"

	"github.com/SCKelemen/color"
	"github.com/SCKelemen/units"
)

// Scale is the universal interface for all scale types.
// Scales map from a domain (input data space) to a range (output visual space).
//
// Ranges use units.Length to support relative units (%, px, em, etc.).
// Values stay abstract until final rendering, minimizing conversion errors.
//
// Inspired by D3 scales and Observable Plot scales.
type Scale interface {
	// Apply maps a domain value to a range value
	Apply(value interface{}) units.Length

	// ApplyValue maps a domain value to a normalized value (0-1 interpolation factor)
	ApplyValue(value interface{}) float64

	// Domain returns the input domain
	Domain() interface{}

	// Range returns the output range
	Range() [2]units.Length

	// Type returns the scale type
	Type() ScaleType

	// Clone creates a copy of this scale
	Clone() Scale
}

// ScaleType identifies the type of scale
type ScaleType int

const (
	ScaleTypeLinear ScaleType = iota
	ScaleTypeLog
	ScaleTypePow
	ScaleTypeSqrt
	ScaleTypeOrdinal
	ScaleTypeBand
	ScaleTypePoint
	ScaleTypeTime
	ScaleTypeQuantize
	ScaleTypeQuantile
	ScaleTypeThreshold
	ScaleTypeIdentity
	ScaleTypeSequential
	ScaleTypeDiverging
)

// String returns the scale type name
func (t ScaleType) String() string {
	switch t {
	case ScaleTypeLinear:
		return "linear"
	case ScaleTypeLog:
		return "log"
	case ScaleTypePow:
		return "pow"
	case ScaleTypeSqrt:
		return "sqrt"
	case ScaleTypeOrdinal:
		return "ordinal"
	case ScaleTypeBand:
		return "band"
	case ScaleTypePoint:
		return "point"
	case ScaleTypeTime:
		return "time"
	case ScaleTypeQuantize:
		return "quantize"
	case ScaleTypeQuantile:
		return "quantile"
	case ScaleTypeThreshold:
		return "threshold"
	case ScaleTypeIdentity:
		return "identity"
	case ScaleTypeSequential:
		return "sequential"
	case ScaleTypeDiverging:
		return "diverging"
	default:
		return "unknown"
	}
}

// ContinuousScale extends Scale for continuous numeric scales
type ContinuousScale interface {
	Scale

	// Invert maps a range value back to a domain value
	Invert(value units.Length) float64

	// InvertValue maps a normalized value (0-1) back to a domain value
	InvertValue(t float64) float64

	// Ticks generates nice tick values for axes
	Ticks(count int) []float64

	// Nice rounds the domain to nice numbers
	Nice(count int) ContinuousScale

	// Clamp enables/disables clamping output to range
	Clamp(enabled bool) ContinuousScale
}

// CategoricalScale extends Scale for categorical scales (Band, Point, Ordinal)
type CategoricalScale interface {
	Scale

	// Values returns all domain values
	Values() []string

	// Index returns the index of a value in the domain
	Index(value string) int
}

// ColorScale maps domain values to colors
type ColorScale interface {
	Scale

	// ApplyColor maps a domain value to a color
	ApplyColor(value interface{}) color.Color

	// Interpolator returns the color interpolation function
	Interpolator() func(t float64) color.Color
}

// InterpolatorFunc is a function that interpolates between two values
type InterpolatorFunc func(t float64) float64

// ColorInterpolatorFunc is a function that interpolates colors
type ColorInterpolatorFunc func(t float64) color.Color

// NiceOptions configures the Nice operation
type NiceOptions struct {
	Count int  // Desired number of ticks
	Floor bool // If true, floor the domain start
	Ceil  bool // If true, ceil the domain end
}

// TickOptions configures tick generation
type TickOptions struct {
	Count  int     // Desired number of ticks
	Format string  // Format string for tick labels
	Values []float64 // Explicit tick values (overrides Count)
}

// TimeTickOptions configures time tick generation
type TimeTickOptions struct {
	Interval TimeInterval // Interval for ticks (year, month, day, etc.)
	Count    int          // Desired number of ticks
	Format   string       // Time format string
}

// TimeInterval represents time intervals for TimeScale
type TimeInterval int

const (
	TimeIntervalMillisecond TimeInterval = iota
	TimeIntervalSecond
	TimeIntervalMinute
	TimeIntervalHour
	TimeIntervalDay
	TimeIntervalWeek
	TimeIntervalMonth
	TimeIntervalYear
)

// BandScaleOptions configures BandScale behavior
type BandScaleOptions struct {
	Padding      float64 // Outer padding (0-1)
	PaddingInner float64 // Inner padding between bands (0-1)
	PaddingOuter float64 // Outer padding at edges (0-1)
	Align        float64 // Alignment within range (0-1, 0.5 = center)
	Round        bool    // Round band positions to integers
}

// QuantizeScaleOptions configures QuantizeScale
type QuantizeScaleOptions struct {
	Thresholds []float64 // Explicit threshold values
	Nice       bool      // Round thresholds to nice numbers
}

// LogScaleOptions configures LogScale
type LogScaleOptions struct {
	Base  float64 // Logarithm base (default: 10)
	Clamp bool    // Clamp output to range
}

// PowScaleOptions configures PowScale
type PowScaleOptions struct {
	Exponent float64 // Power exponent (default: 1)
	Clamp    bool    // Clamp output to range
}

// SequentialScaleOptions configures SequentialScale
type SequentialScaleOptions struct {
	Interpolator ColorInterpolatorFunc // Color interpolation function
	Clamp        bool                  // Clamp output
}

// ScaleConfig is a universal configuration for creating scales
type ScaleConfig struct {
	Type         ScaleType
	Domain       interface{}     // Type depends on scale type
	Range        [2]units.Length // Output range with units
	Clamp        bool
	Nice         bool
	NiceCount    int
	Round        bool
	Padding      float64
	PaddingInner float64
	PaddingOuter float64
	Align        float64
	Base         float64       // For log scales
	Exponent     float64       // For pow scales
	Unknown      units.Length  // For ordinal scales
	Interpolator interface{}   // InterpolatorFunc or ColorInterpolatorFunc
}

// DefaultScaleConfig returns default configuration
func DefaultScaleConfig() ScaleConfig {
	return ScaleConfig{
		Type:         ScaleTypeLinear,
		Domain:       [2]float64{0, 1},
		Range:        [2]units.Length{units.Px(0), units.Px(1)},
		Clamp:        false,
		Nice:         false,
		NiceCount:    10,
		Round:        false,
		Padding:      0,
		PaddingInner: 0,
		PaddingOuter: 0,
		Align:        0.5,
		Base:         10,
		Exponent:     1,
		Unknown:      units.Px(0),
	}
}

// TimeValue wraps a time.Time for use in scales
type TimeValue struct {
	Time time.Time
}

// Value returns the Unix timestamp in seconds
func (tv TimeValue) Value() float64 {
	return float64(tv.Time.Unix())
}
