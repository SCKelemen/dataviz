package dataviz

import (
	"math"

	"github.com/SCKelemen/color"
	design "github.com/SCKelemen/design-system"
)

// calculateContributionLightness calculates the lightness value for a contribution ratio
// Uses a 4-level contrast curve matching GitHub's style
// Returns lightness value from 0.15 (very dark) to 0.85 (bright)
func calculateContributionLightness(ratio float64) float64 {
	if ratio <= 0 {
		return 0.15 // Very dark for zero contributions
	}

	var lightness float64
	// 4-level contrast curve using lightness (0.0 = black, 1.0 = white)
	if ratio < 0.25 {
		lightness = 0.15 + ratio*0.4 // 0.15 to 0.25 (very dark to dark)
	} else if ratio < 0.5 {
		lightness = 0.25 + (ratio-0.25)*0.6 // 0.25 to 0.40 (dark to medium-dark)
	} else if ratio < 0.75 {
		lightness = 0.40 + (ratio-0.5)*0.3 // 0.40 to 0.55 (medium-dark to medium)
	} else {
		lightness = 0.55 + (ratio-0.75)*0.25 // 0.55 to 0.70 (medium to bright)
	}
	return math.Min(0.85, lightness) // Cap at 85% lightness for good contrast
}

// AdjustColorForContribution adjusts a color's lightness based on contribution ratio
// Uses the color package's HSL functions to set absolute lightness
func AdjustColorForContribution(hexColor string, ratio float64) string {
	// Parse color using the color package
	c, err := color.ParseColor(hexColor)
	if err != nil {
		return hexColor // Return original if parsing fails
	}

	// Calculate target lightness using the 4-level contrast curve
	targetLightness := calculateContributionLightness(ratio)

	// Clamp lightness to valid range
	if targetLightness < 0 {
		targetLightness = 0
	}
	if targetLightness > 1 {
		targetLightness = 1
	}

	// Convert to HSL, set absolute lightness (preserve hue and saturation), convert back to hex
	hsl := color.ToHSL(c)
	return color.RGBToHex(color.NewHSL(hsl.H, hsl.S, targetLightness, hsl.A))
}

// CalculateStatCardHeight calculates the height needed for a stat card
// If hasTrendGraph is true, includes space for the trend graph (15px)
func CalculateStatCardHeight(hasTrendGraph bool, tokens *design.DesignTokens) int {
	if hasTrendGraph {
		return tokens.Layout.StatCardHeightTrend
	}
	return tokens.Layout.StatCardHeight
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
