package scales

import (
	"math"

	"github.com/SCKelemen/units"
)

// PowScale implements a continuous power scale.
// Maps a continuous domain [d0, d1] to a continuous range [r0, r1] using
// power transformation (value^exponent).
//
// Ranges use units.Length to support relative units (%, px, em, etc.).
//
// Example:
//   scale := NewPowScale([2]float64{0, 100}, [2]units.Length{units.Px(0), units.Px(500)})
//   scale.Exponent(2) // Square transformation
//   scale.Apply(0)   // Returns units.Px(0)
//   scale.Apply(50)  // Returns units.Px(125) - (50/100)^2 * 500
//   scale.Apply(100) // Returns units.Px(500)
//
// Power scales are useful for:
// - Area scales (exponent = 0.5 for sqrt, radius → area)
// - Volume scales (exponent = 1/3 for cube root)
// - Custom non-linear transformations
// - Emphasizing small or large values
type PowScale struct {
	domain   [2]float64
	range_   [2]units.Length
	exponent float64
	clamp    bool
}

// NewPowScale creates a new power scale with exponent 1 (linear)
func NewPowScale(domain [2]float64, range_ [2]units.Length) *PowScale {
	return &PowScale{
		domain:   domain,
		range_:   range_,
		exponent: 1,
		clamp:    false,
	}
}

// NewSqrtScale creates a new square root scale (exponent = 0.5)
func NewSqrtScale(domain [2]float64, range_ [2]units.Length) *PowScale {
	return &PowScale{
		domain:   domain,
		range_:   range_,
		exponent: 0.5,
		clamp:    false,
	}
}

// Apply maps a domain value to a range value
func (s *PowScale) Apply(value interface{}) units.Length {
	t := s.ApplyValue(value)

	// Interpolate between range values
	r0 := s.range_[0].Value
	r1 := s.range_[1].Value
	unit := s.range_[0].Unit

	result := r0 + t*(r1-r0)

	return units.Length{Value: result, Unit: unit}
}

// ApplyValue maps a domain value to a normalized value (0-1 interpolation factor)
func (s *PowScale) ApplyValue(value interface{}) float64 {
	v, ok := value.(float64)
	if !ok {
		// Try int
		if i, ok := value.(int); ok {
			v = float64(i)
		} else {
			return 0
		}
	}

	d0 := s.domain[0]
	d1 := s.domain[1]

	// Normalize to [0, 1]
	t := (v - d0) / (d1 - d0)

	if s.clamp {
		t = clampFloat(t, 0, 1)
	}

	// Apply power transformation
	// Handle negative values for odd exponents
	if t < 0 && s.exponent != math.Floor(s.exponent) {
		// Fractional exponent with negative base is undefined
		if s.clamp {
			return 0
		}
		return math.NaN()
	}

	var result float64
	if t < 0 {
		// Preserve sign for negative values with integer exponents
		result = -math.Pow(-t, s.exponent)
	} else {
		result = math.Pow(t, s.exponent)
	}

	return result
}

// Invert maps a range value back to a domain value
func (s *PowScale) Invert(value units.Length) float64 {
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
func (s *PowScale) InvertValue(t float64) float64 {
	// Invert power transformation
	var inversed float64
	if t < 0 {
		// Preserve sign for negative values
		inversed = -math.Pow(-t, 1/s.exponent)
	} else {
		inversed = math.Pow(t, 1/s.exponent)
	}

	d0 := s.domain[0]
	d1 := s.domain[1]

	return d0 + inversed*(d1-d0)
}

// Domain returns the input domain
func (s *PowScale) Domain() interface{} {
	return s.domain
}

// Range returns the output range
func (s *PowScale) Range() [2]units.Length {
	return s.range_
}

// Type returns the scale type
func (s *PowScale) Type() ScaleType {
	// Return Sqrt if exponent is 0.5
	if math.Abs(s.exponent-0.5) < 0.001 {
		return ScaleTypeSqrt
	}
	return ScaleTypePow
}

// Clone creates a copy of this scale
func (s *PowScale) Clone() Scale {
	return &PowScale{
		domain:   s.domain,
		range_:   s.range_,
		exponent: s.exponent,
		clamp:    s.clamp,
	}
}

// Clamp enables/disables clamping output to range
func (s *PowScale) Clamp(enabled bool) ContinuousScale {
	s.clamp = enabled
	return s
}

// Exponent sets the power exponent
func (s *PowScale) Exponent(exponent float64) *PowScale {
	s.exponent = exponent
	return s
}

// Nice rounds the domain to nice numbers
// For power scales, we use the same nice logic as linear scales
func (s *PowScale) Nice(count int) ContinuousScale {
	d0, d1 := s.domain[0], s.domain[1]

	if d0 == d1 {
		return s
	}

	// Calculate nice step size (using linear logic)
	step := niceNumber((d1-d0)/float64(count-1), false)

	// Round domain to nice boundaries
	s.domain[0] = math.Floor(d0/step) * step
	s.domain[1] = math.Ceil(d1/step) * step

	return s
}

// Ticks generates nice tick values for axes
// For power scales, we generate linear ticks (not transformed)
func (s *PowScale) Ticks(count int) []float64 {
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
func (s *PowScale) WithDomain(domain [2]float64) *PowScale {
	s.domain = domain
	return s
}

// WithRange sets a new range
func (s *PowScale) WithRange(range_ [2]units.Length) *PowScale {
	s.range_ = range_
	return s
}
