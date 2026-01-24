# Color Strategy and SCKelemen/color Integration

## Overview

All color operations in dataviz use `SCKelemen/color` package. This document defines our color strategy and integration patterns.

## Core Principle

**Use SCKelemen/color for ALL color operations.** No custom color parsing, manipulation, or conversion.

## Available Functions in SCKelemen/color

### Parsing (parse.go)

```go
// Universal color parser - handles all formats
func ParseColor(s string) (Color, error)
```

Supported formats:
- **Hex**: `#RGB`, `#RRGGBB`, `#RRGGBBAA`
- **RGB**: `rgb(255, 0, 0)`, `rgba(255, 0, 0, 0.5)`
- **HSL**: `hsl(0, 100%, 50%)`, `hsla(0, 100%, 50%, 0.5)`
- **HSV**: `hsv(0, 100%, 100%)`
- **LAB**: `lab(50%, 50, -50)`
- **OKLAB**: `oklab(0.5, 0.1, -0.1)`
- **LCH**: `lch(50%, 50, 180deg)`
- **OKLCH**: `oklch(0.5, 0.15, 180deg)`  **← Preferred for scales**
- **HWB**: `hwb(0, 0%, 0%)`
- **XYZ**: `xyz(0.5, 0.5, 0.5)`
- **Named**: `red`, `blue`, `rebeccapurple`, etc.
- **Color function**: `color(display-p3 1 0 0)`

```go
// Direct hex parsing
func HexToRGB(hex string) (*RGBA, error)

// Hex generation
func RGBToHex(c Color) string
```

### Color Manipulation (helpers.go, gradients.go)

```go
// Interpolation in different color spaces
func Mix(c1, c2 Color, t float64, space Space) Color

// Darken/Lighten
func Darken(c Color, amount float64) Color
func Lighten(c Color, amount float64) Color

// Gradients (perceptually uniform!)
func Gradient(stops []Color, positions []float64, space Space) func(float64) Color
```

### Color Spaces (space.go)

```go
type Space int

const (
    Linear Space = iota  // Linear RGB
    SRGB                 // sRGB (web standard)
    OKLCH                // OKLCH (perceptually uniform) ← PREFERRED
    LAB                  // CIE LAB
    LCH                  // CIE LCH
    HSL                  // HSL
    HSV                  // HSV
)
```

### Accessibility (accessibility.go)

```go
// WCAG contrast ratio
func ContrastRatio(c1, c2 Color) float64

// Check WCAG compliance
func MeetsWCAG(fg, bg Color, level WCAGLevel, size TextSize) bool

// Colorblind simulation
func SimulateColorBlindness(c Color, cbType ColorBlindnessType) Color

type ColorBlindnessType int
const (
    Protanopia ColorBlindnessType = iota  // Red-blind
    Deuteranopia                          // Green-blind
    Tritanopia                            // Blue-blind
    Protanomaly                           // Red-weak
    Deuteranomaly                         // Green-weak
    Tritanomaly                           // Blue-weak
)
```

## Dataviz Integration Patterns

### Pattern 1: Parse User Input

```go
import "github.com/SCKelemen/color"

// Accept any color format from users
func parseUserColor(input string) (color.Color, error) {
    return color.ParseColor(input)
}

// Example usage:
c1, _ := parseUserColor("#3b82f6")           // Hex
c2, _ := parseUserColor("rgb(59, 130, 246)") // RGB function
c3, _ := parseUserColor("oklch(0.6, 0.2, 250deg)") // OKLCH
c4, _ := parseUserColor("rebeccapurple")     // Named color
```

### Pattern 2: Scales with OKLCH Interpolation

```go
// Color scales MUST use OKLCH for perceptually uniform gradients
func createColorScale(minColor, maxColor color.Color) func(float64) color.Color {
    return func(t float64) color.Color {
        return color.Mix(minColor, maxColor, t, color.OKLCH)
    }
}

// Example: Heatmap
minColor, _ := color.ParseColor("#3b82f6") // Blue
maxColor, _ := color.ParseColor("#ef4444") // Red

scale := createColorScale(minColor, maxColor)

// Smooth perceptual gradient (no muddy middle!)
color0   := scale(0.0)  // Pure blue
color25  := scale(0.25) // Blue-ish
color50  := scale(0.5)  // Purple (perceptually middle)
color75  := scale(0.75) // Red-ish
color100 := scale(1.0)  // Pure red
```

### Pattern 3: Multi-Stop Gradients

```go
// For heatmaps with multiple color stops
stops := []color.Color{
    mustParse("oklch(0.3, 0.2, 250deg)"), // Dark blue
    mustParse("oklch(0.5, 0.2, 200deg)"), // Cyan
    mustParse("oklch(0.7, 0.15, 150deg)"), // Green
    mustParse("oklch(0.8, 0.15, 90deg)"), // Yellow
    mustParse("oklch(0.6, 0.25, 30deg)"), // Red
}
positions := []float64{0.0, 0.25, 0.5, 0.75, 1.0}

gradient := color.Gradient(stops, positions, color.OKLCH)

// Use in scale
heatmapColor := gradient(0.63) // Interpolates between yellow and red
```

### Pattern 4: Accessibility Checks

```go
// Ensure text is readable against background
func ensureReadable(textColor, bgColor color.Color) color.Color {
    ratio := color.ContrastRatio(textColor, bgColor)

    if ratio < 4.5 { // WCAG AA for normal text
        // Adjust lightness until readable
        if textColor.Luminance() > bgColor.Luminance() {
            textColor = color.Lighten(textColor, 0.2)
        } else {
            textColor = color.Darken(textColor, 0.2)
        }
    }

    return textColor
}

// Example: Legend text color
legendTextColor := ensureReadable(textColor, legendBackground)
```

### Pattern 5: Colorblind-Safe Palettes

```go
// Test palette against colorblind simulation
func testPaletteSafety(palette []color.Color) bool {
    types := []color.ColorBlindnessType{
        color.Deuteranopia,  // Most common
        color.Protanopia,
        color.Tritanopia,
    }

    for _, cbType := range types {
        simulated := make([]color.Color, len(palette))
        for i, c := range palette {
            simulated[i] = color.SimulateColorBlindness(c, cbType)
        }

        // Check if colors are still distinguishable
        if !areDistinguishable(simulated) {
            return false
        }
    }

    return true
}
```

## What We're Already Using

In `charts/legends/`:
```go
import "github.com/SCKelemen/color"

// Parsing
c, err := color.HexToRGB("#3b82f6")

// Conversion
hexString := color.RGBToHex(c)

// Alpha channel
alpha := c.Alpha()
```

## What to Upstream to SCKelemen/color

### 1. Shorthand Constructors (Convenience)

```go
// For tests and examples - panic on error
func MustParse(s string) Color {
    c, err := ParseColor(s)
    if err != nil {
        panic(err)
    }
    return c
}

// For common hex colors
func Hex(s string) Color {
    c, _ := HexToRGB(s) // Returns black on error
    return c
}
```

**Status**: Would be nice to have, but not critical. We can use `mustParse` helper in tests.

### 2. Palette Validation (Accessibility)

```go
// Check if colors in a palette are distinguishable
func AreDistinguishable(colors []Color, threshold float64) bool {
    for i := 0; i < len(colors); i++ {
        for j := i + 1; j < len(colors); j++ {
            if ColorDifference(colors[i], colors[j]) < threshold {
                return false
            }
        }
    }
    return true
}

// Color difference using deltaE (already exists?)
func ColorDifference(c1, c2 Color) float64 {
    // CIE deltaE 2000
}
```

**Status**: Check if deltaE functions exist. If not, upstream this.

### 3. Named Color Sets (Convenience)

```go
// Predefined colorblind-safe palettes
var ColorblindSafePalette = []Color{
    MustParse("#0173B2"), // Blue
    MustParse("#DE8F05"), // Orange
    MustParse("#029E73"), // Green
    MustParse("#CC78BC"), // Purple
    MustParse("#CA9161"), // Brown
    MustParse("#FBAFE4"), // Pink
}

// Wong palette (well-known colorblind-safe)
var WongPalette = []Color{ /* ... */ }

// Tol palettes
var TolBright = []Color{ /* ... */ }
var TolMuted = []Color{ /* ... */ }
```

**Status**: Very useful. Consider upstreaming to `color/palettes` package.

## Color Strategy by Use Case

| Use Case | Color Space | Rationale |
|----------|-------------|-----------|
| **Heatmaps** | OKLCH | Perceptually uniform gradients |
| **Categorical data** | Named palette | Pre-tested for distinction |
| **Line charts (multi-series)** | Named palette | Need clear distinction |
| **Area fills** | OKLCH + alpha | Smooth gradients with transparency |
| **Text/labels** | Check contrast | WCAG compliance |
| **Backgrounds** | sRGB | Web standard |

## Examples

### Example 1: Heatmap with OKLCH Gradient

```go
import "github.com/SCKelemen/color"

// Parse colors
blue, _ := color.ParseColor("oklch(0.5, 0.2, 250deg)")
red, _ := color.ParseColor("oklch(0.6, 0.25, 30deg)")

// Create perceptually uniform gradient
gradient := func(t float64) color.Color {
    return color.Mix(blue, red, t, color.OKLCH)
}

// Use in heatmap
for value := 0.0; value <= 1.0; value += 0.1 {
    cellColor := gradient(value)
    // Render cell with cellColor
}
```

### Example 2: Accessible Text Color

```go
bgColor, _ := color.ParseColor("#1e293b") // Dark background
textColor, _ := color.ParseColor("#94a3b8") // Light gray text

// Check readability
ratio := color.ContrastRatio(textColor, bgColor)
if ratio < 4.5 {
    // Adjust until readable
    textColor = color.Lighten(textColor, 0.3)
}

// Use textColor for labels
```

### Example 3: Colorblind-Safe Palette

```go
// Use pre-validated palette for categorical data
palette := []string{
    "#0173B2", // Blue (safe)
    "#DE8F05", // Orange (safe)
    "#029E73", // Green (safe)
    "#CC78BC", // Purple (safe)
}

colors := make([]color.Color, len(palette))
for i, hex := range palette {
    colors[i], _ = color.HexToRGB(hex)
}

// Verify safety
for _, cbType := range []color.ColorBlindnessType{
    color.Deuteranopia,
    color.Protanopia,
} {
    for i, c := range colors {
        simulated := color.SimulateColorBlindness(c, cbType)
        fmt.Printf("Color %d as %v: %v\n", i, cbType, simulated)
    }
}
```

## Implementation Checklist

- [x] Use `color.ParseColor()` for all user input
- [x] Use `color.HexToRGB()` for known hex colors
- [x] Use `color.RGBToHex()` for SVG output
- [x] Use `color.OKLCH` for all gradient scales
- [ ] Check if `DeltaE()` exists, upstream if not
- [ ] Consider upstreaming `MustParse()` helper
- [ ] Consider upstreaming colorblind-safe palettes
- [ ] Document preferred palettes in dataviz
- [ ] Add accessibility checks to Context
- [ ] Add colorblind simulation to Context

## Success Criteria

- ✅ Zero custom color parsing in dataviz
- ✅ All gradients use OKLCH
- ✅ All charts pass WCAG AA contrast
- ✅ Palettes are colorblind-safe
- ✅ Color operations are type-safe
- ✅ Consistent color handling across SVG/terminal

## Related Documents

- [SURFACE_CANVAS_ARCHITECTURE.md](SURFACE_CANVAS_ARCHITECTURE.md) - Context uses OKLCH
- [ROADMAP.md](ROADMAP.md) - Scales use perceptual color spaces
- SCKelemen/color: https://github.com/SCKelemen/color
