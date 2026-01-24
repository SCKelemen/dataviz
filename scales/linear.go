package scales

import (
	"math"

	"github.com/SCKelemen/units"
)

// LinearScale implements a continuous linear scale.
// Maps a continuous domain [d0, d1] to a continuous range [r0, r1] using linear interpolation.
//
// Ranges use units.Length to support relative units (%, px, em, etc.).
//
// Example:
//   scale := NewLinearScale([2]float64{0, 100}, [2]units.Length{units.Px(0), units.Px(500)})
//   scale.Apply(50) // Returns units.Px(250)
//   scale.Invert(units.Px(250)) // Returns 50
type LinearScale struct {
	domain [2]float64
	range_ [2]units.Length
	clamp  bool
}

// NewLinearScale creates a new linear scale
func NewLinearScale(domain [2]float64, range_ [2]units.Length) *LinearScale {
	return &LinearScale{
		domain: domain,
		range_: range_,
		clamp:  false,
	}
}

// Apply maps a domain value to a range value
func (s *LinearScale) Apply(value interface{}) units.Length {
	t := s.ApplyValue(value)

	// Interpolate between range values
	// For now, we assume both range values are in the same unit
	// and do linear interpolation on their raw values
	r0 := s.range_[0].Value
	r1 := s.range_[1].Value
	unit := s.range_[0].Unit

	result := r0 + t*(r1-r0)

	return units.Length{Value: result, Unit: unit}
}

// ApplyValue maps a domain value to a normalized value (0-1 interpolation factor)
func (s *LinearScale) ApplyValue(value interface{}) float64 {
	v, ok := value.(float64)
	if !ok {
		// Try int
		if i, ok := value.(int); ok {
			v = float64(i)
		} else {
			return 0
		}
	}

	// Linear interpolation parameter
	t := (v - s.domain[0]) / (s.domain[1] - s.domain[0])

	if s.clamp {
		t = clampFloat(t, 0, 1)
	}

	return t
}

// Invert maps a range value back to a domain value
func (s *LinearScale) Invert(value units.Length) float64 {
	// Convert value to same unit as range
	// For simplicity, we'll work with raw values assuming same unit
	v := value.Value
	r0 := s.range_[0].Value
	r1 := s.range_[1].Value

	t := (v - r0) / (r1 - r0)

	if s.clamp {
		t = clampFloat(t, 0, 1)
	}

	return s.InvertValue(t)
}

// InvertValue maps a normalized value (0-1) back to a domain value
func (s *LinearScale) InvertValue(t float64) float64 {
	return s.domain[0] + t*(s.domain[1]-s.domain[0])
}

// Domain returns the input domain
func (s *LinearScale) Domain() interface{} {
	return s.domain
}

// Range returns the output range
func (s *LinearScale) Range() [2]units.Length {
	return s.range_
}

// Type returns the scale type
func (s *LinearScale) Type() ScaleType {
	return ScaleTypeLinear
}

// Clone creates a copy of this scale
func (s *LinearScale) Clone() Scale {
	return &LinearScale{
		domain: s.domain,
		range_: s.range_,
		clamp:  s.clamp,
	}
}

// Clamp enables/disables clamping output to range
func (s *LinearScale) Clamp(enabled bool) ContinuousScale {
	s.clamp = enabled
	return s
}

// Nice rounds the domain to nice numbers
func (s *LinearScale) Nice(count int) ContinuousScale {
	d0, d1 := s.domain[0], s.domain[1]

	if d0 == d1 {
		return s
	}

	// Calculate nice step size
	step := niceNumber((d1-d0)/float64(count-1), false)

	// Round domain to nice boundaries
	s.domain[0] = math.Floor(d0/step) * step
	s.domain[1] = math.Ceil(d1/step) * step

	return s
}

// Ticks generates nice tick values for axes
func (s *LinearScale) Ticks(count int) []float64 {
	if count <= 0 {
		count = 10
	}

	d0, d1 := s.domain[0], s.domain[1]

	if d0 == d1 {
		return []float64{d0}
	}

	// Calculate nice step size
	step := niceNumber((d1-d0)/float64(count-1), true)

	if step == 0 {
		return []float64{d0}
	}

	// Generate ticks
	start := math.Ceil(d0 / step)
	stop := math.Floor(d1 / step)
	n := int(stop - start + 1)

	if n <= 0 {
		return []float64{d0}
	}

	ticks := make([]float64, n)
	for i := 0; i < n; i++ {
		ticks[i] = (start + float64(i)) * step
	}

	return ticks
}

// WithDomain sets a new domain
func (s *LinearScale) WithDomain(domain [2]float64) *LinearScale {
	s.domain = domain
	return s
}

// WithRange sets a new range
func (s *LinearScale) WithRange(range_ [2]units.Length) *LinearScale {
	s.range_ = range_
	return s
}

// Helper functions

func clampFloat(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// niceNumber rounds a number to a "nice" value (1, 2, 5, 10, 20, 50, 100, etc.)
func niceNumber(value float64, round bool) float64 {
	exponent := math.Floor(math.Log10(math.Abs(value)))
	fraction := math.Abs(value) / math.Pow(10, exponent)

	var niceFraction float64

	if round {
		if fraction < 1.5 {
			niceFraction = 1
		} else if fraction < 3 {
			niceFraction = 2
		} else if fraction < 7 {
			niceFraction = 5
		} else {
			niceFraction = 10
		}
	} else {
		if fraction <= 1 {
			niceFraction = 1
		} else if fraction <= 2 {
			niceFraction = 2
		} else if fraction <= 5 {
			niceFraction = 5
		} else {
			niceFraction = 10
		}
	}

	result := niceFraction * math.Pow(10, exponent)

	if value < 0 {
		return -result
	}
	return result
}
