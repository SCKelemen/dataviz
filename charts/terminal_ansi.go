package charts

import (
	"fmt"

	"github.com/SCKelemen/color"
)

// ANSI escape codes for terminal colors and styling
const (
	ansiReset     = "\x1b[0m"
	ansiBold      = "\x1b[1m"
	ansiDim       = "\x1b[2m"
	ansiItalic    = "\x1b[3m"
	ansiUnderline = "\x1b[4m"
)

// TerminalColorMode represents the color capability of the terminal
type TerminalColorMode int

const (
	TerminalColorNone TerminalColorMode = iota
	TerminalColor16                     // 16 colors
	TerminalColor256                    // 256 colors
	TerminalColorTrue                   // 24-bit true color
)

// ColorForeground returns ANSI escape code for foreground color
func ColorForeground(c string, mode TerminalColorMode) string {
	if c == "" {
		return ""
	}

	col, err := color.ParseColor(c)
	if err != nil {
		return ""
	}

	r, g, b, _ := col.RGBA()
	r8 := int(r * 255)
	g8 := int(g * 255)
	b8 := int(b * 255)

	switch mode {
	case TerminalColorNone:
		return ""
	case TerminalColor16:
		code := rgbToANSI16(r8, g8, b8)
		return fmt.Sprintf("\x1b[%dm", code)
	case TerminalColor256:
		code := rgbToANSI256(r8, g8, b8)
		return fmt.Sprintf("\x1b[38;5;%dm", code)
	case TerminalColorTrue:
		return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r8, g8, b8)
	default:
		return ""
	}
}

// ColorBackground returns ANSI escape code for background color
func ColorBackground(c string, mode TerminalColorMode) string {
	if c == "" {
		return ""
	}

	col, err := color.ParseColor(c)
	if err != nil {
		return ""
	}

	r, g, b, _ := col.RGBA()
	r8 := int(r * 255)
	g8 := int(g * 255)
	b8 := int(b * 255)

	switch mode {
	case TerminalColorNone:
		return ""
	case TerminalColor16:
		code := rgbToANSI16(r8, g8, b8)
		return fmt.Sprintf("\x1b[%dm", code+10) // Background = foreground + 10
	case TerminalColor256:
		code := rgbToANSI256(r8, g8, b8)
		return fmt.Sprintf("\x1b[48;5;%dm", code)
	case TerminalColorTrue:
		return fmt.Sprintf("\x1b[48;2;%d;%d;%dm", r8, g8, b8)
	default:
		return ""
	}
}

// rgbToANSI16 converts RGB to closest 16-color ANSI code
func rgbToANSI16(r, g, b int) int {
	ansi16Colors := []struct {
		code    int
		r, g, b int
	}{
		{30, 0, 0, 0},       // Black
		{31, 170, 0, 0},     // Red
		{32, 0, 170, 0},     // Green
		{33, 170, 85, 0},    // Yellow
		{34, 0, 0, 170},     // Blue
		{35, 170, 0, 170},   // Magenta
		{36, 0, 170, 170},   // Cyan
		{37, 170, 170, 170}, // White
		{90, 85, 85, 85},    // Bright Black
		{91, 255, 85, 85},   // Bright Red
		{92, 85, 255, 85},   // Bright Green
		{93, 255, 255, 85},  // Bright Yellow
		{94, 85, 85, 255},   // Bright Blue
		{95, 255, 85, 255},  // Bright Magenta
		{96, 85, 255, 255},  // Bright Cyan
		{97, 255, 255, 255}, // Bright White
	}

	minDist := 1000000
	closestCode := 37

	for _, c := range ansi16Colors {
		dr := r - c.r
		dg := g - c.g
		db := b - c.b
		dist := dr*dr + dg*dg + db*db

		if dist < minDist {
			minDist = dist
			closestCode = c.code
		}
	}

	return closestCode
}

// rgbToANSI256 converts RGB to closest 256-color palette index
func rgbToANSI256(r, g, b int) int {
	// Check if it's a gray
	if absInt(r-g) < 10 && absInt(g-b) < 10 && absInt(b-r) < 10 {
		if r < 8 {
			return 16 // Black
		}
		if r > 247 {
			return 231 // White
		}
		return 232 + (r-8)/10
	}

	// Use 6x6x6 color cube (16-231)
	r6 := (r * 6) / 256
	g6 := (g * 6) / 256
	b6 := (b * 6) / 256

	return 16 + 36*r6 + 6*g6 + b6
}

// InterpolateColorGradient creates a gradient of ANSI color codes between two colors
func InterpolateColorGradient(start, end string, steps int, mode TerminalColorMode) []string {
	if steps < 2 {
		steps = 2
	}

	startCol, err := color.ParseColor(start)
	if err != nil {
		return make([]string, steps)
	}

	endCol, err := color.ParseColor(end)
	if err != nil {
		return make([]string, steps)
	}

	// Interpolate in OKLCH for perceptually uniform gradients
	colors := color.GradientInSpace(startCol, endCol, steps, color.GradientOKLCH)
	result := make([]string, steps)

	for i := 0; i < steps; i++ {
		hexColor := color.RGBToHex(colors[i])
		result[i] = ColorForeground(hexColor, mode)
	}

	return result
}

// absInt returns absolute value of an integer
func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
