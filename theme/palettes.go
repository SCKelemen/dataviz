package theme

import (
	"math"

	"github.com/SCKelemen/color"
)

// DarkCategorical returns a categorical color palette for dark mode
func DarkCategorical() []string {
	return []string{
		"#60A5FA", // blue
		"#34D399", // green
		"#F472B6", // pink
		"#FBBF24", // yellow
		"#A78BFA", // purple
		"#FB923C", // orange
		"#2DD4BF", // teal
		"#F87171", // red
		"#FCD34D", // amber
		"#818CF8", // indigo
	}
}

// LightCategorical returns a categorical color palette for light mode
func LightCategorical() []string {
	return []string{
		"#3B82F6", // blue
		"#10B981", // green
		"#EC4899", // pink
		"#F59E0B", // yellow
		"#8B5CF6", // purple
		"#F97316", // orange
		"#14B8A6", // teal
		"#EF4444", // red
		"#F59E0B", // amber
		"#6366F1", // indigo
	}
}

// DarkSequential returns a sequential color palette for dark mode
func DarkSequential(accent string) []string {
	// Generate a 9-step sequential palette from dark to accent color
	return generateSequential(accent, "#1F2937", 9)
}

// LightSequential returns a sequential color palette for light mode
func LightSequential(accent string) []string {
	// Generate a 9-step sequential palette from light to accent color
	return generateSequential(accent, "#F9FAFB", 9)
}

// DarkDiverging returns a diverging color palette for dark mode
func DarkDiverging(accent string) []string {
	// Blue (cool) to accent (warm) diverging palette
	return generateDiverging("#3B82F6", "#F9FAFB", accent, 9)
}

// LightDiverging returns a diverging color palette for light mode
func LightDiverging(accent string) []string {
	// Blue (cool) to white to accent (warm) diverging palette
	return generateDiverging("#3B82F6", "#FFFFFF", accent, 9)
}

// generateSequential generates a sequential color palette
func generateSequential(endColor, startColor string, steps int) []string {
	if steps < 2 {
		return []string{endColor}
	}

	// Parse colors
	start, err := color.ParseColor(startColor)
	if err != nil {
		return []string{endColor}
	}

	end, err := color.ParseColor(endColor)
	if err != nil {
		return []string{endColor}
	}

	// Generate interpolated colors
	palette := make([]string, steps)
	for i := 0; i < steps; i++ {
		t := float64(i) / float64(steps-1)

		// Interpolate in Lab color space for perceptual uniformity
		interpolated := color.MixInSpace(start, end, t, color.GradientLAB)
		palette[i] = color.RGBToHex(interpolated)
	}

	return palette
}

// generateDiverging generates a diverging color palette
func generateDiverging(leftColor, middleColor, rightColor string, steps int) []string {
	if steps < 3 {
		return []string{leftColor, middleColor, rightColor}
	}

	// Parse colors
	left, err := color.ParseColor(leftColor)
	if err != nil {
		return []string{leftColor, middleColor, rightColor}
	}

	middle, err := color.ParseColor(middleColor)
	if err != nil {
		return []string{leftColor, middleColor, rightColor}
	}

	right, err := color.ParseColor(rightColor)
	if err != nil {
		return []string{leftColor, middleColor, rightColor}
	}

	// Ensure odd number of steps so middle is exact
	if steps%2 == 0 {
		steps++
	}

	midpoint := steps / 2
	palette := make([]string, steps)

	// Generate left half (left to middle)
	for i := 0; i <= midpoint; i++ {
		t := float64(i) / float64(midpoint)
		interpolated := color.MixInSpace(left, middle, t, color.GradientLAB)
		palette[i] = color.RGBToHex(interpolated)
	}

	// Generate right half (middle to right)
	for i := midpoint + 1; i < steps; i++ {
		t := float64(i-midpoint) / float64(steps-midpoint-1)
		interpolated := color.MixInSpace(middle, right, t, color.GradientLAB)
		palette[i] = color.RGBToHex(interpolated)
	}

	return palette
}

// HeatmapColors returns colors for a heatmap based on value (0-1)
func HeatmapColors(value float64, darkMode bool) string {
	// Clamp value to [0, 1]
	if value < 0 {
		value = 0
	}
	if value > 1 {
		value = 1
	}

	if darkMode {
		// Dark mode: darker colors for low values
		if value == 0 {
			return "#1F2937"
		}
		// Interpolate from dark blue to bright blue
		return interpolateHex("#1E3A8A", "#60A5FA", value)
	} else {
		// Light mode: lighter colors for low values
		if value == 0 {
			return "#F9FAFB"
		}
		// Interpolate from light blue to dark blue
		return interpolateHex("#DBEAFE", "#3B82F6", value)
	}
}

// interpolateHex interpolates between two hex colors
func interpolateHex(start, end string, t float64) string {
	startColor, err := color.ParseColor(start)
	if err != nil {
		return start
	}

	endColor, err := color.ParseColor(end)
	if err != nil {
		return end
	}

	interpolated := color.MixInSpace(startColor, endColor, t, color.GradientLAB)
	return color.RGBToHex(interpolated)
}

// ViridisColors returns colors from the Viridis palette (0-1)
// Viridis is a perceptually uniform color map
func ViridisColors(t float64) string {
	// Clamp t to [0, 1]
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}

	// Viridis control points (approximation)
	colors := []string{
		"#440154", // purple
		"#31688E", // blue
		"#35B779", // green
		"#FDE724", // yellow
	}

	// Find segment
	segment := t * float64(len(colors)-1)
	index := int(segment)
	if index >= len(colors)-1 {
		return colors[len(colors)-1]
	}

	// Interpolate within segment
	t_segment := segment - float64(index)
	return interpolateHex(colors[index], colors[index+1], t_segment)
}

// PlasmaColors returns colors from the Plasma palette (0-1)
// Plasma is another perceptually uniform color map
func PlasmaColors(t float64) string {
	// Clamp t to [0, 1]
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}

	// Plasma control points (approximation)
	colors := []string{
		"#0D0887", // dark purple
		"#7E03A8", // purple
		"#CC4778", // pink
		"#F89540", // orange
		"#F0F921", // yellow
	}

	// Find segment
	segment := t * float64(len(colors)-1)
	index := int(segment)
	if index >= len(colors)-1 {
		return colors[len(colors)-1]
	}

	// Interpolate within segment
	t_segment := segment - float64(index)
	return interpolateHex(colors[index], colors[index+1], t_segment)
}

// CoolWarmColors returns colors from a cool-warm diverging palette (-1 to 1)
func CoolWarmColors(t float64) string {
	// Map from [-1, 1] to [0, 1]
	normalized := (t + 1.0) / 2.0

	// Clamp to [0, 1]
	if normalized < 0 {
		normalized = 0
	}
	if normalized > 1 {
		normalized = 1
	}

	// Cool (blue) to neutral (white) to warm (red)
	if normalized < 0.5 {
		// Blue to white
		t_segment := normalized * 2.0
		return interpolateHex("#3B82F6", "#F9FAFB", t_segment)
	} else {
		// White to red
		t_segment := (normalized - 0.5) * 2.0
		return interpolateHex("#F9FAFB", "#EF4444", t_segment)
	}
}

// QualitativeColors returns high-contrast colors for categorical data
func QualitativeColors(index int, darkMode bool) string {
	var colors []string
	if darkMode {
		colors = DarkCategorical()
	} else {
		colors = LightCategorical()
	}

	if len(colors) == 0 {
		return "#60A5FA"
	}

	return colors[index%len(colors)]
}

// ContrastRatio calculates the WCAG contrast ratio between two colors
func ContrastRatio(color1, color2 string) float64 {
	c1, err := color.ParseColor(color1)
	if err != nil {
		return 1.0
	}

	c2, err := color.ParseColor(color2)
	if err != nil {
		return 1.0
	}

	// Calculate relative luminance
	luminance := func(c color.Color) float64 {
		r, g, b, _ := c.RGBA()

		// Convert to sRGB
		toLinear := func(v float64) float64 {
			if v <= 0.03928 {
				return v / 12.92
			}
			return math.Pow((v+0.055)/1.055, 2.4)
		}

		rLinear := toLinear(r)
		gLinear := toLinear(g)
		bLinear := toLinear(b)

		// Calculate relative luminance
		return 0.2126*rLinear + 0.7152*gLinear + 0.0722*bLinear
	}

	l1 := luminance(c1)
	l2 := luminance(c2)

	// Ensure l1 is the lighter color
	if l1 < l2 {
		l1, l2 = l2, l1
	}

	return (l1 + 0.05) / (l2 + 0.05)
}

// EnsureContrast adjusts a foreground color to ensure sufficient contrast with background
func EnsureContrast(fg, bg string, minRatio float64) string {
	ratio := ContrastRatio(fg, bg)
	if ratio >= minRatio {
		return fg
	}

	// Parse foreground color
	fgColor, err := color.ParseColor(fg)
	if err != nil {
		return fg
	}

	// Parse background color
	bgColor, err := color.ParseColor(bg)
	if err != nil {
		return fg
	}

	// Determine if background is light or dark
	bgLuminance := func() float64 {
		r, g, b, _ := bgColor.RGBA()
		return 0.299*r + 0.587*g + 0.114*b
	}()

	// Convert to HSL for adjustment
	hslColor := color.ToHSL(fgColor)
	h, s, l := hslColor.H, hslColor.S, hslColor.L

	// If background is light, darken foreground
	// If background is dark, lighten foreground
	step := 0.05
	if bgLuminance > 0.5 {
		// Darken
		for l > 0 && ContrastRatio(color.RGBToHex(color.NewHSL(h, s, l, 1.0)), bg) < minRatio {
			l -= step
		}
	} else {
		// Lighten
		for l < 1 && ContrastRatio(color.RGBToHex(color.NewHSL(h, s, l, 1.0)), bg) < minRatio {
			l += step
		}
	}

	return color.RGBToHex(color.NewHSL(h, s, l, 1.0))
}
