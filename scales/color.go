package scales

import (
	"github.com/SCKelemen/color"
	"github.com/SCKelemen/units"
)

// SequentialColorScale maps a continuous domain to a sequential color gradient.
// Sequential scales use a single hue progressing from light to dark, ideal for
// representing magnitude or intensity.
//
// Example:
//   scale := NewSequentialColorScale(
//     [2]float64{0, 100},
//     color.RGB(1, 1, 1), // White
//     color.RGB(0, 0, 1), // Blue
//   )
//   scale.Apply(50) // Returns middle blue color
//
// Common uses: heatmaps, choropleth maps, magnitude encoding
type SequentialColorScale struct {
	domain      [2]float64
	startColor  color.Color
	endColor    color.Color
	clamp       bool
	space       color.GradientSpace
	interpolate InterpolatorFunc
}

// NewSequentialColorScale creates a new sequential color scale
func NewSequentialColorScale(domain [2]float64, start, end color.Color) *SequentialColorScale {
	return &SequentialColorScale{
		domain:     domain,
		startColor: start,
		endColor:   end,
		clamp:      true,
		space:      color.GradientOKLCH, // Perceptually uniform by default
	}
}

// Apply maps a domain value to a color
func (s *SequentialColorScale) Apply(value interface{}) units.Length {
	// Color scales don't return units.Length, but we need to satisfy Scale interface
	// Return a dummy value - use ApplyColor instead
	return units.Px(0)
}

// ApplyColor maps a domain value to a color
func (s *SequentialColorScale) ApplyColor(value interface{}) color.Color {
	t := s.ApplyValue(value)
	return color.MixInSpace(s.startColor, s.endColor, t, s.space)
}

// ApplyValue maps a domain value to a normalized value (0-1 interpolation factor)
func (s *SequentialColorScale) ApplyValue(value interface{}) float64 {
	v, ok := value.(float64)
	if !ok {
		if i, ok := value.(int); ok {
			v = float64(i)
		} else {
			return 0
		}
	}

	t := (v - s.domain[0]) / (s.domain[1] - s.domain[0])

	if s.clamp {
		if t < 0 {
			t = 0
		}
		if t > 1 {
			t = 1
		}
	}

	// Apply custom interpolator if set
	if s.interpolate != nil {
		t = s.interpolate(t)
	}

	return t
}

// Domain returns the input domain
func (s *SequentialColorScale) Domain() interface{} {
	return s.domain
}

// Range returns a dummy range for Scale interface compatibility
func (s *SequentialColorScale) Range() [2]units.Length {
	return [2]units.Length{units.Px(0), units.Px(1)}
}

// Type returns the scale type
func (s *SequentialColorScale) Type() ScaleType {
	return ScaleTypeSequential
}

// Clone creates a copy of this scale
func (s *SequentialColorScale) Clone() Scale {
	return &SequentialColorScale{
		domain:      s.domain,
		startColor:  s.startColor,
		endColor:    s.endColor,
		clamp:       s.clamp,
		space:       s.space,
		interpolate: s.interpolate,
	}
}

// Clamp enables/disables clamping output to range
func (s *SequentialColorScale) Clamp(enabled bool) *SequentialColorScale {
	s.clamp = enabled
	return s
}

// Space sets the color interpolation space
func (s *SequentialColorScale) Space(space color.GradientSpace) *SequentialColorScale {
	s.space = space
	return s
}

// Interpolate sets a custom interpolation function
func (s *SequentialColorScale) Interpolate(fn InterpolatorFunc) *SequentialColorScale {
	s.interpolate = fn
	return s
}

// Samples generates n evenly-spaced color samples
func (s *SequentialColorScale) Samples(n int) []color.Color {
	if n <= 0 {
		return nil
	}
	if n == 1 {
		return []color.Color{s.ApplyColor(s.domain[0])}
	}

	samples := make([]color.Color, n)
	for i := 0; i < n; i++ {
		t := float64(i) / float64(n-1)
		value := s.domain[0] + t*(s.domain[1]-s.domain[0])
		samples[i] = s.ApplyColor(value)
	}
	return samples
}

// DivergingColorScale maps a continuous domain to a diverging color gradient.
// Diverging scales use two distinct hues meeting at a neutral midpoint, ideal
// for data with a natural center or comparing deviations from a reference.
//
// Example:
//   scale := NewDivergingColorScale(
//     [2]float64{-100, 100},
//     color.RGB(0, 0, 1), // Blue (negative)
//     color.RGB(1, 1, 1), // White (neutral)
//     color.RGB(1, 0, 0), // Red (positive)
//   )
//   scale.Apply(-50) // Returns light blue
//   scale.Apply(0)   // Returns white
//   scale.Apply(50)  // Returns light red
//
// Common uses: temperature anomalies, profit/loss, before/after comparisons
type DivergingColorScale struct {
	domain      [2]float64
	startColor  color.Color
	midColor    color.Color
	endColor    color.Color
	midpoint    float64 // Domain value for midColor (default: domain midpoint)
	clamp       bool
	space       color.GradientSpace
	interpolate InterpolatorFunc
}

// NewDivergingColorScale creates a new diverging color scale
func NewDivergingColorScale(domain [2]float64, start, mid, end color.Color) *DivergingColorScale {
	return &DivergingColorScale{
		domain:     domain,
		startColor: start,
		midColor:   mid,
		endColor:   end,
		midpoint:   (domain[0] + domain[1]) / 2, // Default to domain midpoint
		clamp:      true,
		space:      color.GradientOKLCH,
	}
}

// Apply maps a domain value to a color (dummy for Scale interface)
func (s *DivergingColorScale) Apply(value interface{}) units.Length {
	return units.Px(0)
}

// ApplyColor maps a domain value to a color
func (s *DivergingColorScale) ApplyColor(value interface{}) color.Color {
	v, ok := value.(float64)
	if !ok {
		if i, ok := value.(int); ok {
			v = float64(i)
		} else {
			return s.midColor
		}
	}

	// Normalize to [0, 1] where 0.5 is the midpoint
	var t float64
	if v < s.midpoint {
		// Map [domain[0], midpoint] to [0, 0.5]
		t = 0.5 * (v - s.domain[0]) / (s.midpoint - s.domain[0])
	} else {
		// Map [midpoint, domain[1]] to [0.5, 1]
		t = 0.5 + 0.5*(v-s.midpoint)/(s.domain[1]-s.midpoint)
	}

	if s.clamp {
		if t < 0 {
			t = 0
		}
		if t > 1 {
			t = 1
		}
	}

	// Apply custom interpolator if set
	if s.interpolate != nil {
		t = s.interpolate(t)
	}

	// Interpolate colors
	if t < 0.5 {
		// Between start and mid
		localT := t * 2 // Map [0, 0.5] to [0, 1]
		return color.MixInSpace(s.startColor, s.midColor, localT, s.space)
	} else {
		// Between mid and end
		localT := (t - 0.5) * 2 // Map [0.5, 1] to [0, 1]
		return color.MixInSpace(s.midColor, s.endColor, localT, s.space)
	}
}

// ApplyValue maps a domain value to a normalized value
func (s *DivergingColorScale) ApplyValue(value interface{}) float64 {
	v, ok := value.(float64)
	if !ok {
		if i, ok := value.(int); ok {
			v = float64(i)
		} else {
			return 0.5
		}
	}

	// Normalize to [0, 1]
	t := (v - s.domain[0]) / (s.domain[1] - s.domain[0])

	if s.clamp {
		if t < 0 {
			t = 0
		}
		if t > 1 {
			t = 1
		}
	}

	return t
}

// Domain returns the input domain
func (s *DivergingColorScale) Domain() interface{} {
	return s.domain
}

// Range returns a dummy range
func (s *DivergingColorScale) Range() [2]units.Length {
	return [2]units.Length{units.Px(0), units.Px(1)}
}

// Type returns the scale type
func (s *DivergingColorScale) Type() ScaleType {
	return ScaleTypeDiverging
}

// Clone creates a copy of this scale
func (s *DivergingColorScale) Clone() Scale {
	return &DivergingColorScale{
		domain:      s.domain,
		startColor:  s.startColor,
		midColor:    s.midColor,
		endColor:    s.endColor,
		midpoint:    s.midpoint,
		clamp:       s.clamp,
		space:       s.space,
		interpolate: s.interpolate,
	}
}

// Midpoint sets the domain value for the midpoint color
func (s *DivergingColorScale) Midpoint(value float64) *DivergingColorScale {
	s.midpoint = value
	return s
}

// Clamp enables/disables clamping
func (s *DivergingColorScale) Clamp(enabled bool) *DivergingColorScale {
	s.clamp = enabled
	return s
}

// Space sets the color interpolation space
func (s *DivergingColorScale) Space(space color.GradientSpace) *DivergingColorScale {
	s.space = space
	return s
}

// Interpolate sets a custom interpolation function
func (s *DivergingColorScale) Interpolate(fn InterpolatorFunc) *DivergingColorScale {
	s.interpolate = fn
	return s
}

// Samples generates n evenly-spaced color samples
func (s *DivergingColorScale) Samples(n int) []color.Color {
	if n <= 0 {
		return nil
	}
	if n == 1 {
		return []color.Color{s.ApplyColor(s.midpoint)}
	}

	samples := make([]color.Color, n)
	for i := 0; i < n; i++ {
		t := float64(i) / float64(n-1)
		value := s.domain[0] + t*(s.domain[1]-s.domain[0])
		samples[i] = s.ApplyColor(value)
	}
	return samples
}

// CategoricalColorScale maps discrete categories to distinct colors.
// Each category gets its own color, with cycling for domains larger than the color range.
//
// Example:
//   colors := []color.Color{
//     color.RGB(1, 0, 0), // Red
//     color.RGB(0, 1, 0), // Green
//     color.RGB(0, 0, 1), // Blue
//   }
//   scale := NewCategoricalColorScale(
//     []string{"Group A", "Group B", "Group C"},
//     colors,
//   )
//   scale.ApplyColor("Group A") // Returns red
//
// Common uses: categorical data, group comparisons, qualitative differences
type CategoricalColorScale struct {
	domain  []string
	colors  []color.Color
	unknown color.Color
}

// NewCategoricalColorScale creates a new categorical color scale
func NewCategoricalColorScale(domain []string, colors []color.Color) *CategoricalColorScale {
	return &CategoricalColorScale{
		domain:  domain,
		colors:  colors,
		unknown: color.RGB(0.5, 0.5, 0.5), // Gray for unknown
	}
}

// Apply maps a domain value (dummy for Scale interface)
func (s *CategoricalColorScale) Apply(value interface{}) units.Length {
	return units.Px(0)
}

// ApplyColor maps a category to a color
func (s *CategoricalColorScale) ApplyColor(value interface{}) color.Color {
	v, ok := value.(string)
	if !ok {
		return s.unknown
	}

	// Find index in domain
	index := -1
	for i, category := range s.domain {
		if category == v {
			index = i
			break
		}
	}

	if index < 0 || len(s.colors) == 0 {
		return s.unknown
	}

	// Cycle through colors if domain is larger
	return s.colors[index%len(s.colors)]
}

// ApplyValue maps a category to its index (normalized)
func (s *CategoricalColorScale) ApplyValue(value interface{}) float64 {
	v, ok := value.(string)
	if !ok {
		return 0
	}

	for i, category := range s.domain {
		if category == v {
			if len(s.domain) == 1 {
				return 0.5
			}
			return float64(i) / float64(len(s.domain)-1)
		}
	}

	return 0
}

// Domain returns the input domain
func (s *CategoricalColorScale) Domain() interface{} {
	return s.domain
}

// Range returns a dummy range
func (s *CategoricalColorScale) Range() [2]units.Length {
	return [2]units.Length{units.Px(0), units.Px(1)}
}

// Type returns the scale type
func (s *CategoricalColorScale) Type() ScaleType {
	return ScaleTypeOrdinal
}

// Clone creates a copy of this scale
func (s *CategoricalColorScale) Clone() Scale {
	domainCopy := make([]string, len(s.domain))
	copy(domainCopy, s.domain)
	colorsCopy := make([]color.Color, len(s.colors))
	copy(colorsCopy, s.colors)

	return &CategoricalColorScale{
		domain:  domainCopy,
		colors:  colorsCopy,
		unknown: s.unknown,
	}
}

// Unknown sets the color for unknown categories
func (s *CategoricalColorScale) Unknown(c color.Color) *CategoricalColorScale {
	s.unknown = c
	return s
}

// Colors returns all colors in the scale
func (s *CategoricalColorScale) Colors() []color.Color {
	return s.colors
}
