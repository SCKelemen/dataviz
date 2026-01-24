package scales

import (
	"math"

	"github.com/SCKelemen/units"
)

// LogScale implements a continuous logarithmic scale.
// Maps a continuous domain [d0, d1] to a continuous range [r0, r1] using
// logarithmic interpolation.
//
// Ranges use units.Length to support relative units (%, px, em, etc.).
//
// Example:
//   scale := NewLogScale([2]float64{1, 1000}, [2]units.Length{units.Px(0), units.Px(500)})
//   scale.Apply(1)    // Returns units.Px(0)
//   scale.Apply(10)   // Returns units.Px(166.67) - 1/3 of way on log scale
//   scale.Apply(100)  // Returns units.Px(333.33) - 2/3 of way
//   scale.Apply(1000) // Returns units.Px(500)
//
// Log scales are useful for data spanning multiple orders of magnitude:
// - Population (1K to 1B)
// - Wealth/income distributions
// - Scientific measurements
// - Earthquake magnitudes
type LogScale struct {
	domain [2]float64
	range_ [2]units.Length
	base   float64
	clamp  bool
}

// NewLogScale creates a new log scale with base 10
func NewLogScale(domain [2]float64, range_ [2]units.Length) *LogScale {
	return &LogScale{
		domain: domain,
		range_: range_,
		base:   10,
		clamp:  false,
	}
}

// Apply maps a domain value to a range value
func (s *LogScale) Apply(value interface{}) units.Length {
	t := s.ApplyValue(value)

	// Interpolate between range values
	r0 := s.range_[0].Value
	r1 := s.range_[1].Value
	unit := s.range_[0].Unit

	result := r0 + t*(r1-r0)

	return units.Length{Value: result, Unit: unit}
}

// ApplyValue maps a domain value to a normalized value (0-1 interpolation factor)
func (s *LogScale) ApplyValue(value interface{}) float64 {
	v, ok := value.(float64)
	if !ok {
		// Try int
		if i, ok := value.(int); ok {
			v = float64(i)
		} else {
			return 0
		}
	}

	// Handle zero and negative values
	if v <= 0 {
		if s.clamp {
			return 0
		}
		return math.NaN()
	}

	// Log interpolation parameter
	logV := math.Log(v) / math.Log(s.base)
	logD0 := math.Log(s.domain[0]) / math.Log(s.base)
	logD1 := math.Log(s.domain[1]) / math.Log(s.base)

	t := (logV - logD0) / (logD1 - logD0)

	if s.clamp {
		t = clampFloat(t, 0, 1)
	}

	return t
}

// Invert maps a range value back to a domain value
func (s *LogScale) Invert(value units.Length) float64 {
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
func (s *LogScale) InvertValue(t float64) float64 {
	logD0 := math.Log(s.domain[0]) / math.Log(s.base)
	logD1 := math.Log(s.domain[1]) / math.Log(s.base)

	logV := logD0 + t*(logD1-logD0)

	return math.Pow(s.base, logV)
}

// Domain returns the input domain
func (s *LogScale) Domain() interface{} {
	return s.domain
}

// Range returns the output range
func (s *LogScale) Range() [2]units.Length {
	return s.range_
}

// Type returns the scale type
func (s *LogScale) Type() ScaleType {
	return ScaleTypeLog
}

// Clone creates a copy of this scale
func (s *LogScale) Clone() Scale {
	return &LogScale{
		domain: s.domain,
		range_: s.range_,
		base:   s.base,
		clamp:  s.clamp,
	}
}

// Clamp enables/disables clamping output to range
func (s *LogScale) Clamp(enabled bool) ContinuousScale {
	s.clamp = enabled
	return s
}

// Base sets the logarithm base
func (s *LogScale) Base(base float64) *LogScale {
	if base <= 0 || base == 1 {
		base = 10 // Default to base 10 for invalid values
	}
	s.base = base
	return s
}

// Nice rounds the domain to nice powers of the base
func (s *LogScale) Nice(count int) ContinuousScale {
	d0, d1 := s.domain[0], s.domain[1]

	if d0 <= 0 || d1 <= 0 || d0 == d1 {
		return s
	}

	// Round to powers of base
	logD0 := math.Log(d0) / math.Log(s.base)
	logD1 := math.Log(d1) / math.Log(s.base)

	s.domain[0] = math.Pow(s.base, math.Floor(logD0))
	s.domain[1] = math.Pow(s.base, math.Ceil(logD1))

	return s
}

// Ticks generates nice tick values for axes (powers of base)
func (s *LogScale) Ticks(count int) []float64 {
	if count <= 0 {
		count = 10
	}

	d0, d1 := s.domain[0], s.domain[1]

	if d0 <= 0 || d1 <= 0 || d0 == d1 {
		return []float64{d0}
	}

	// Calculate power range
	logD0 := math.Log(d0) / math.Log(s.base)
	logD1 := math.Log(d1) / math.Log(s.base)

	startPow := int(math.Floor(logD0))
	endPow := int(math.Ceil(logD1))

	var ticks []float64

	// Generate ticks at powers of base
	for pow := startPow; pow <= endPow; pow++ {
		tick := math.Pow(s.base, float64(pow))
		if tick >= d0 && tick <= d1 {
			ticks = append(ticks, tick)
		}
	}

	// If we don't have enough ticks, add intermediate values
	if len(ticks) < count && s.base == 10 {
		// For base 10, add intermediate ticks (1, 2, 3, 4, 5, 6, 7, 8, 9) * 10^n
		var allTicks []float64
		for pow := startPow; pow <= endPow; pow++ {
			baseTick := math.Pow(s.base, float64(pow))
			// mult=1 gives us the power itself (1 * 10^n)
			// mult=2-9 gives us intermediate values
			for mult := 1; mult < int(s.base); mult++ {
				tick := baseTick * float64(mult)
				if tick >= d0 && tick <= d1 {
					allTicks = append(allTicks, tick)
				}
			}
		}

		// Add the final power (10^(endPow+1)) if it's within domain
		// This handles the case where domain ends exactly at a power
		finalPower := math.Pow(s.base, float64(endPow+1))
		if finalPower >= d0 && finalPower <= d1 {
			// Check if not already added (avoid duplicates)
			if len(allTicks) == 0 || allTicks[len(allTicks)-1] != finalPower {
				allTicks = append(allTicks, finalPower)
			}
		}

		if len(allTicks) > 0 {
			ticks = allTicks
		}
	}

	return ticks
}

// WithDomain sets a new domain
func (s *LogScale) WithDomain(domain [2]float64) *LogScale {
	s.domain = domain
	return s
}

// WithRange sets a new range
func (s *LogScale) WithRange(range_ [2]units.Length) *LogScale {
	s.range_ = range_
	return s
}
