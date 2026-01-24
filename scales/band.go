package scales

import (
	"math"

	"github.com/SCKelemen/units"
)

// BandScale implements a categorical scale with bands.
// Maps discrete domain values (categories) to continuous range bands.
// Used for bar charts, where each category gets a band with optional padding.
//
// Ranges use units.Length to support relative units (%, px, em, etc.).
//
// Example:
//   scale := NewBandScale([]string{"A", "B", "C"}, [2]units.Length{units.Px(0), units.Px(300)})
//   scale.Apply("A") // Returns units.Px(0) (start of first band)
//   scale.Apply("B") // Returns units.Px(100) (start of second band)
//   scale.Bandwidth() // Returns units.Px(100) (width of each band)
type BandScale struct {
	domain       []string
	range_       [2]units.Length
	padding      float64 // Outer padding (0-1)
	paddingInner float64 // Inner padding between bands (0-1)
	paddingOuter float64 // Outer padding at edges (0-1)
	align        float64 // Alignment within range (0-1, 0.5 = center)
	round        bool    // Round band positions to integers
	bandwidth    float64 // Computed bandwidth (raw value)
	step         float64 // Computed step size (band + padding, raw value)
	start        float64 // Computed start position (raw value)
}

// NewBandScale creates a new band scale
func NewBandScale(domain []string, range_ [2]units.Length) *BandScale {
	s := &BandScale{
		domain:       domain,
		range_:       range_,
		padding:      0,
		paddingInner: 0,
		paddingOuter: 0,
		align:        0.5,
		round:        false,
	}
	s.rescale()
	return s
}

// Apply maps a domain value to a range value (band start position)
func (s *BandScale) Apply(value interface{}) units.Length {
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
func (s *BandScale) ApplyValue(value interface{}) float64 {
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

	bandStart := s.start + float64(index)*s.step
	return (bandStart - s.range_[0].Value) / rangeSize
}

// Domain returns the input domain
func (s *BandScale) Domain() interface{} {
	return s.domain
}

// Range returns the output range
func (s *BandScale) Range() [2]units.Length {
	return s.range_
}

// Type returns the scale type
func (s *BandScale) Type() ScaleType {
	return ScaleTypeBand
}

// Clone creates a copy of this scale
func (s *BandScale) Clone() Scale {
	clone := &BandScale{
		domain:       make([]string, len(s.domain)),
		range_:       s.range_,
		padding:      s.padding,
		paddingInner: s.paddingInner,
		paddingOuter: s.paddingOuter,
		align:        s.align,
		round:        s.round,
	}
	copy(clone.domain, s.domain)
	clone.rescale()
	return clone
}

// Values returns all domain values
func (s *BandScale) Values() []string {
	return s.domain
}

// Index returns the index of a value in the domain
func (s *BandScale) Index(value string) int {
	for i, v := range s.domain {
		if v == value {
			return i
		}
	}
	return -1
}

// Bandwidth returns the width of each band
func (s *BandScale) Bandwidth() units.Length {
	return units.Length{Value: s.bandwidth, Unit: s.range_[0].Unit}
}

// Step returns the step size (band + padding)
func (s *BandScale) Step() units.Length {
	return units.Length{Value: s.step, Unit: s.range_[0].Unit}
}

// Padding sets the outer padding (0-1)
func (s *BandScale) Padding(padding float64) *BandScale {
	s.padding = clampFloat(padding, 0, 1)
	s.paddingInner = padding
	s.paddingOuter = padding
	s.rescale()
	return s
}

// PaddingInner sets the inner padding between bands (0-1)
func (s *BandScale) PaddingInner(padding float64) *BandScale {
	s.paddingInner = clampFloat(padding, 0, 1)
	s.rescale()
	return s
}

// PaddingOuter sets the outer padding at edges (0-1)
func (s *BandScale) PaddingOuter(padding float64) *BandScale {
	s.paddingOuter = clampFloat(padding, 0, 1)
	s.rescale()
	return s
}

// Align sets the alignment within range (0-1, 0.5 = center)
func (s *BandScale) Align(align float64) *BandScale {
	s.align = clampFloat(align, 0, 1)
	s.rescale()
	return s
}

// Round enables/disables rounding band positions to integers
func (s *BandScale) Round(round bool) *BandScale {
	s.round = round
	s.rescale()
	return s
}

// WithDomain sets a new domain
func (s *BandScale) WithDomain(domain []string) *BandScale {
	s.domain = domain
	s.rescale()
	return s
}

// WithRange sets a new range
func (s *BandScale) WithRange(range_ [2]units.Length) *BandScale {
	s.range_ = range_
	s.rescale()
	return s
}

// rescale recomputes bandwidth, step, and start based on current settings
func (s *BandScale) rescale() {
	n := len(s.domain)
	if n == 0 {
		s.bandwidth = 0
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

	// Calculate step and bandwidth
	// step = (range - outerPadding * 2) / (n - innerPadding * (n-1))
	// bandwidth = step * (1 - innerPadding)

	step := (stop - start) / math.Max(1, float64(n)-s.paddingInner+s.paddingOuter*2)

	if s.round {
		step = math.Floor(step)
	}

	// Calculate start position with alignment
	start += (stop - start - step*(float64(n)-s.paddingInner)) * s.align

	// Calculate bandwidth
	bandwidth := step * (1 - s.paddingInner)

	if s.round {
		start = math.Round(start)
		bandwidth = math.Round(bandwidth)
	}

	if reverse {
		// Reverse the scale
		s.start = stop - (start - s.range_[1].Value)
		s.step = -step
		s.bandwidth = bandwidth
	} else {
		s.start = start
		s.step = step
		s.bandwidth = bandwidth
	}
}
