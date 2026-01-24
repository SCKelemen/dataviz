package theme

import (
	design "github.com/SCKelemen/design-system"
)

// Default returns the default theme
func Default() *Theme {
	return New(design.DefaultTheme())
}

// Midnight returns the midnight theme (dark, deep blues)
func Midnight() *Theme {
	return New(design.MidnightTheme())
}

// Nord returns the Nord theme (muted, cool colors)
func Nord() *Theme {
	return New(design.NordTheme())
}

// Paper returns the Paper theme (light, clean)
func Paper() *Theme {
	return New(design.PaperTheme())
}

// Wrapped returns the Wrapped theme (vibrant, special styling)
func Wrapped() *Theme {
	return New(design.WrappedTheme())
}

// FromTokens creates a theme from custom design tokens
func FromTokens(tokens *design.DesignTokens) *Theme {
	return New(tokens)
}

// Monochrome returns a monochrome theme (grayscale)
func Monochrome(darkMode bool) *Theme {
	var tokens *design.DesignTokens
	if darkMode {
		tokens = &design.DesignTokens{
			Theme:      "monochrome-dark",
			Color:      "#E5E7EB",
			Background: "#111827",
			Accent:     "#9CA3AF",
			FontFamily: "system-ui",
			Radius:     8,
			Padding:    16,
			Density:    "comfortable",
			Mode:       "dark",
			Layout:     design.DefaultLayoutTokens(),
		}
	} else {
		tokens = &design.DesignTokens{
			Theme:      "monochrome-light",
			Color:      "#1F2937",
			Background: "#FFFFFF",
			Accent:     "#6B7280",
			FontFamily: "system-ui",
			Radius:     8,
			Padding:    16,
			Density:    "comfortable",
			Mode:       "light",
			Layout:     design.DefaultLayoutTokens(),
		}
	}

	theme := New(tokens)

	// Override color schemes to use grayscale
	if darkMode {
		theme.ColorScheme.Sequential = []string{
			"#111827", "#1F2937", "#374151", "#4B5563", "#6B7280",
			"#9CA3AF", "#D1D5DB", "#E5E7EB", "#F3F4F6",
		}
		theme.ColorScheme.Categorical = []string{
			"#9CA3AF", "#6B7280", "#D1D5DB", "#4B5563",
			"#E5E7EB", "#374151", "#F3F4F6", "#1F2937",
		}
	} else {
		theme.ColorScheme.Sequential = []string{
			"#F9FAFB", "#F3F4F6", "#E5E7EB", "#D1D5DB", "#9CA3AF",
			"#6B7280", "#4B5563", "#374151", "#1F2937",
		}
		theme.ColorScheme.Categorical = []string{
			"#6B7280", "#9CA3AF", "#4B5563", "#D1D5DB",
			"#374151", "#E5E7EB", "#1F2937", "#F3F4F6",
		}
	}

	return theme
}

// Ocean returns an ocean-inspired theme
func Ocean(darkMode bool) *Theme {
	var tokens *design.DesignTokens
	if darkMode {
		tokens = &design.DesignTokens{
			Theme:      "ocean-dark",
			Color:      "#E0F2FE",
			Background: "#0C4A6E",
			Accent:     "#0EA5E9",
			FontFamily: "system-ui",
			Radius:     12,
			Padding:    16,
			Density:    "comfortable",
			Mode:       "dark",
			Layout:     design.DefaultLayoutTokens(),
		}
	} else {
		tokens = &design.DesignTokens{
			Theme:      "ocean-light",
			Color:      "#0C4A6E",
			Background: "#F0F9FF",
			Accent:     "#0284C7",
			FontFamily: "system-ui",
			Radius:     12,
			Padding:    16,
			Density:    "comfortable",
			Mode:       "light",
			Layout:     design.DefaultLayoutTokens(),
		}
	}

	return New(tokens)
}

// Forest returns a forest-inspired theme
func Forest(darkMode bool) *Theme {
	var tokens *design.DesignTokens
	if darkMode {
		tokens = &design.DesignTokens{
			Theme:      "forest-dark",
			Color:      "#D1FAE5",
			Background: "#064E3B",
			Accent:     "#10B981",
			FontFamily: "system-ui",
			Radius:     12,
			Padding:    16,
			Density:    "comfortable",
			Mode:       "dark",
			Layout:     design.DefaultLayoutTokens(),
		}
	} else {
		tokens = &design.DesignTokens{
			Theme:      "forest-light",
			Color:      "#064E3B",
			Background: "#F0FDF4",
			Accent:     "#059669",
			FontFamily: "system-ui",
			Radius:     12,
			Padding:    16,
			Density:    "comfortable",
			Mode:       "light",
			Layout:     design.DefaultLayoutTokens(),
		}
	}

	return New(tokens)
}

// Sunset returns a sunset-inspired theme
func Sunset(darkMode bool) *Theme {
	var tokens *design.DesignTokens
	if darkMode {
		tokens = &design.DesignTokens{
			Theme:      "sunset-dark",
			Color:      "#FED7AA",
			Background: "#7C2D12",
			Accent:     "#F97316",
			FontFamily: "system-ui",
			Radius:     16,
			Padding:    16,
			Density:    "comfortable",
			Mode:       "dark",
			Layout:     design.DefaultLayoutTokens(),
		}
	} else {
		tokens = &design.DesignTokens{
			Theme:      "sunset-light",
			Color:      "#7C2D12",
			Background: "#FFF7ED",
			Accent:     "#EA580C",
			FontFamily: "system-ui",
			Radius:     16,
			Padding:    16,
			Density:    "comfortable",
			Mode:       "light",
			Layout:     design.DefaultLayoutTokens(),
		}
	}

	return New(tokens)
}

// HighContrast returns a high contrast theme for accessibility
func HighContrast(darkMode bool) *Theme {
	var tokens *design.DesignTokens
	if darkMode {
		tokens = &design.DesignTokens{
			Theme:      "high-contrast-dark",
			Color:      "#FFFFFF",
			Background: "#000000",
			Accent:     "#00D9FF",
			FontFamily: "system-ui",
			Radius:     4,
			Padding:    20,
			Density:    "comfortable",
			Mode:       "dark",
			Layout:     design.DefaultLayoutTokens(),
		}
	} else {
		tokens = &design.DesignTokens{
			Theme:      "high-contrast-light",
			Color:      "#000000",
			Background: "#FFFFFF",
			Accent:     "#0070F3",
			FontFamily: "system-ui",
			Radius:     4,
			Padding:    20,
			Density:    "comfortable",
			Mode:       "light",
			Layout:     design.DefaultLayoutTokens(),
		}
	}

	theme := New(tokens)

	// Use highly contrasting colors
	if darkMode {
		theme.ColorScheme.Categorical = []string{
			"#00D9FF", "#FF0080", "#00FF00", "#FFFF00",
			"#FF8000", "#FF00FF", "#00FFFF", "#FF4040",
		}
	} else {
		theme.ColorScheme.Categorical = []string{
			"#0070F3", "#E00", "#0A0", "#F90",
			"#E0E", "#0AA", "#F40", "#90F",
		}
	}

	return theme
}

// Scientific returns a theme suitable for scientific publications
func Scientific() *Theme {
	tokens := &design.DesignTokens{
		Theme:      "scientific",
		Color:      "#000000",
		Background: "#FFFFFF",
		Accent:     "#000000",
		FontFamily: "serif",
		Radius:     0,
		Padding:    16,
		Density:    "comfortable",
		Mode:       "light",
		Layout:     design.DefaultLayoutTokens(),
	}

	theme := New(tokens)

	// Use colorblind-friendly palette (Okabe-Ito)
	theme.ColorScheme.Categorical = []string{
		"#E69F00", // orange
		"#56B4E9", // sky blue
		"#009E73", // bluish green
		"#F0E442", // yellow
		"#0072B2", // blue
		"#D55E00", // vermillion
		"#CC79A7", // reddish purple
		"#000000", // black
	}

	// Monochrome sequential
	theme.ColorScheme.Sequential = []string{
		"#FFFFFF", "#F0F0F0", "#D9D9D9", "#BDBDBD", "#969696",
		"#737373", "#525252", "#252525", "#000000",
	}

	return theme
}

// Minimal returns a minimal theme with subtle styling
func Minimal() *Theme {
	tokens := &design.DesignTokens{
		Theme:      "minimal",
		Color:      "#374151",
		Background: "#FFFFFF",
		Accent:     "#6B7280",
		FontFamily: "system-ui",
		Radius:     2,
		Padding:    12,
		Density:    "compact",
		Mode:       "light",
		Layout:     design.DefaultLayoutTokens(),
	}

	theme := New(tokens)

	// Subtle colors
	theme.Chart.GridStrokeWidth = 0.5
	theme.Chart.AxisStrokeWidth = 1.0
	theme.Chart.GridOpacity = 0.2

	return theme
}
