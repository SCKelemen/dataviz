package scales

import (
	"math"

	"github.com/SCKelemen/units"
)

// PointScale implements a categorical scale that maps discrete domain values
// to evenly-spaced points. Unlike BandScale (which has bandwidth), PointScale
// maps to single positions with no width.
//
// Ranges use units.Length to support relative units (%, px, em, etc.).
//
// Example:
//   scale := NewPointScale(
//     []string{"A", "B", "C"},
//     [2]units.Length{units.Px(0), units.Px(300)},
//   )
//   scale.Apply("A") // Returns units.Px(0) (first point)
//   scale.Apply("B") // Returns units.Px(150) (middle point)
//   scale.Apply("C") // Returns units.Px(300) (last point)
//
// Point scales are commonly used for:
// - Scatter plot categorical axes
// - Dot plots
// - Categorical axes in line charts
// - Any visualization needing evenly-spaced categorical positions
type PointScale struct {
	domain  []string
	range_  [2]units.Length
	padding float64 // Outer padding (0-1)
	align   float64 // Alignment within range (0-1, 0.5 = center)
	round   bool    // Round point positions to integers
	step    float64 // Computed step size (distance between points, raw value)
	start   float64 // Computed start position (raw value)
}

// NewPointScale creates a new point scale
func NewPointScale(domain []string, range_ [2]units.Length) *PointScale {
	s := &PointScale{
		domain:  domain,
		range_:  range_,
		padding: 0,
		align:   0.5,
		round:   false,
	}
	s.rescale()
	return s
}

// Apply maps a domain value to a range value (point position)
func (s *PointScale) Apply(value interface{}) units.Length {
	v, ok := value.(string)
	if !ok {
		return units.Px(0)
	}

	// Find index of value in domain
	index := s.Index(v)
	if index < 0 {
		return units.Px(0) // Unknown value
	}

	result := s.start + float64(index)*s.step
	return units.Length{Value: result, Unit: s.range_[0].Unit}
}

// ApplyValue maps a domain value to a normalized value (index / domain size)
func (s *PointScale) ApplyValue(value interface{}) float64 {
	v, ok := value.(string)
	if !ok {
		return 0
	}

	index := s.Index(v)
	if index < 0 {
		return 0
	}

	// Return normalized position within range
	rangeSize := s.range_[1].Value - s.range_[0].Value
	if rangeSize == 0 {
		return 0
	}

	pointPos := s.start + float64(index)*s.step
	return (pointPos - s.range_[0].Value) / rangeSize
}

// Domain returns the input domain
func (s *PointScale) Domain() interface{} {
	return s.domain
}

// Range returns the output range
func (s *PointScale) Range() [2]units.Length {
	return s.range_
}

// Type returns the scale type
func (s *PointScale) Type() ScaleType {
	return ScaleTypePoint
}

// Clone creates a copy of this scale
func (s *PointScale) Clone() Scale {
	clone := &PointScale{
		domain:  make([]string, len(s.domain)),
		range_:  s.range_,
		padding: s.padding,
		align:   s.align,
		round:   s.round,
	}
	copy(clone.domain, s.domain)
	clone.rescale()
	return clone
}

// Values returns all domain values
func (s *PointScale) Values() []string {
	return s.domain
}

// Index returns the index of a value in the domain
func (s *PointScale) Index(value string) int {
	for i, v := range s.domain {
		if v == value {
			return i
		}
	}
	return -1
}

// Step returns the step size (distance between points)
func (s *PointScale) Step() units.Length {
	return units.Length{Value: s.step, Unit: s.range_[0].Unit}
}

// Padding sets the outer padding (0-1)
func (s *PointScale) Padding(padding float64) *PointScale {
	s.padding = clampFloat(padding, 0, 1)
	s.rescale()
	return s
}

// Align sets the alignment within range (0-1, 0.5 = center)
func (s *PointScale) Align(align float64) *PointScale {
	s.align = clampFloat(align, 0, 1)
	s.rescale()
	return s
}

// Round enables/disables rounding point positions to integers
func (s *PointScale) Round(round bool) *PointScale {
	s.round = round
	s.rescale()
	return s
}

// WithDomain sets a new domain
func (s *PointScale) WithDomain(domain []string) *PointScale {
	s.domain = domain
	s.rescale()
	return s
}

// WithRange sets a new range
func (s *PointScale) WithRange(range_ [2]units.Length) *PointScale {
	s.range_ = range_
	s.rescale()
	return s
}

// rescale recomputes step and start based on current settings
func (s *PointScale) rescale() {
	n := len(s.domain)
	if n == 0 {
		s.step = 0
		s.start = s.range_[0].Value
		return
	}

	// Work with raw values
	reverse := s.range_[1].Value < s.range_[0].Value
	start := s.range_[0].Value
	stop := s.range_[1].Value

	if reverse {
		start, stop = stop, start
	}

	// Calculate step (distance between points)
	// For n points, we need n-1 intervals
	// With padding, we have: padding * step on each side
	// Total range = padding * step + (n-1) * step + padding * step
	// Total range = (n - 1 + 2*padding) * step
	// step = range / (n - 1 + 2*padding)

	var step float64
	if n == 1 {
		// Single point - center it or use alignment
		step = 0
		start = start + (stop-start)*s.align
	} else {
		step = (stop - start) / math.Max(1, float64(n-1)+s.padding*2)

		if s.round {
			step = math.Floor(step)
		}

		// Calculate start position with alignment and padding
		start += step * s.padding
	}

	if s.round && n > 1 {
		start = math.Round(start)
	}

	if reverse {
		// Reverse the scale
		s.start = stop - (start - s.range_[1].Value)
		s.step = -step
	} else {
		s.start = start
		s.step = step
	}
}
