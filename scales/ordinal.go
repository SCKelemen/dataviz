package scales

import (
	"github.com/SCKelemen/units"
)

// OrdinalScale implements a categorical scale that maps discrete domain values
// to discrete range values. Unlike BandScale (which creates bands), OrdinalScale
// maps each domain value to a specific output value.
//
// Ranges use units.Length to support relative units (%, px, em, etc.).
//
// Example:
//   scale := NewOrdinalScale(
//     []string{"small", "medium", "large"},
//     []units.Length{units.Px(10), units.Px(20), units.Px(30)},
//   )
//   scale.Apply("medium") // Returns units.Px(20)
//
// Ordinal scales are commonly used for:
// - Categorical color mapping
// - Size mapping (small/medium/large)
// - Shape mapping (circle/square/triangle)
// - Any discrete categorical encoding
type OrdinalScale struct {
	domain  []string
	range_  []units.Length
	unknown units.Length // Value for unmapped inputs
}

// NewOrdinalScale creates a new ordinal scale
func NewOrdinalScale(domain []string, range_ []units.Length) *OrdinalScale {
	return &OrdinalScale{
		domain:  domain,
		range_:  range_,
		unknown: units.Px(0), // Default unknown value
	}
}

// Apply maps a domain value to a range value
func (s *OrdinalScale) Apply(value interface{}) units.Length {
	v, ok := value.(string)
	if !ok {
		return s.unknown
	}

	index := s.Index(v)
	if index < 0 || len(s.range_) == 0 {
		return s.unknown
	}

	// Cycle through range values if domain is larger than range
	return s.range_[index%len(s.range_)]
}

// ApplyValue maps a domain value to a normalized value
// For ordinal scales, this returns the index position normalized to [0, 1]
func (s *OrdinalScale) ApplyValue(value interface{}) float64 {
	v, ok := value.(string)
	if !ok {
		return 0
	}

	index := s.Index(v)
	if index < 0 || len(s.domain) == 0 {
		return 0
	}

	// Normalize index to [0, 1] range
	if len(s.domain) == 1 {
		return 0.5 // Single value maps to center
	}

	return float64(index) / float64(len(s.domain)-1)
}

// Domain returns the input domain
func (s *OrdinalScale) Domain() interface{} {
	return s.domain
}

// Range returns the output range
func (s *OrdinalScale) Range() [2]units.Length {
	if len(s.range_) == 0 {
		return [2]units.Length{units.Px(0), units.Px(0)}
	}
	if len(s.range_) == 1 {
		return [2]units.Length{s.range_[0], s.range_[0]}
	}
	return [2]units.Length{s.range_[0], s.range_[len(s.range_)-1]}
}

// RangeValues returns all range values
func (s *OrdinalScale) RangeValues() []units.Length {
	return s.range_
}

// Type returns the scale type
func (s *OrdinalScale) Type() ScaleType {
	return ScaleTypeOrdinal
}

// Clone creates a copy of this scale
func (s *OrdinalScale) Clone() Scale {
	clone := &OrdinalScale{
		domain:  make([]string, len(s.domain)),
		range_:  make([]units.Length, len(s.range_)),
		unknown: s.unknown,
	}
	copy(clone.domain, s.domain)
	copy(clone.range_, s.range_)
	return clone
}

// Values returns all domain values
func (s *OrdinalScale) Values() []string {
	return s.domain
}

// Index returns the index of a value in the domain
func (s *OrdinalScale) Index(value string) int {
	for i, v := range s.domain {
		if v == value {
			return i
		}
	}
	return -1
}

// Unknown sets the return value for unknown domain values
func (s *OrdinalScale) Unknown(value units.Length) *OrdinalScale {
	s.unknown = value
	return s
}

// WithDomain sets a new domain
func (s *OrdinalScale) WithDomain(domain []string) *OrdinalScale {
	s.domain = domain
	return s
}

// WithRange sets a new range
func (s *OrdinalScale) WithRange(range_ []units.Length) *OrdinalScale {
	s.range_ = range_
	return s
}
