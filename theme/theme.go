package theme

import (
	design "github.com/SCKelemen/design-system"
	"github.com/SCKelemen/svg"
	"github.com/SCKelemen/units"
)

// Theme represents a complete visual theme for dataviz charts
type Theme struct {
	// Design tokens from design-system package
	Tokens *design.DesignTokens

	// Color palettes for data visualization
	ColorScheme ColorScheme

	// Typography settings
	Typography Typography

	// Chart-specific styling
	Chart ChartStyle
}

// ColorScheme defines color palettes for different chart elements
type ColorScheme struct {
	// Sequential palette (for ordered data)
	Sequential []string

	// Diverging palette (for data with a midpoint)
	Diverging []string

	// Categorical palette (for distinct categories)
	Categorical []string

	// Special colors
	GridColor       string
	AxisColor       string
	TextColor       string
	BackgroundColor string
	BorderColor     string
}

// Typography defines font settings for charts
type Typography struct {
	// Font families
	TitleFont  string
	BodyFont   string
	MonoFont   string

	// Font sizes
	TitleSize    units.Length
	SubtitleSize units.Length
	BodySize     units.Length
	CaptionSize  units.Length
	LabelSize    units.Length

	// Font weights
	TitleWeight  svg.FontWeight
	BodyWeight   svg.FontWeight
	LabelWeight  svg.FontWeight

	// Line heights (multipliers)
	TitleLineHeight float64
	BodyLineHeight  float64
}

// ChartStyle defines chart-specific styling
type ChartStyle struct {
	// Stroke widths
	GridStrokeWidth float64
	AxisStrokeWidth float64
	DataStrokeWidth float64

	// Opacities
	GridOpacity float64
	FillOpacity float64

	// Marker sizes
	PointSize   float64
	MarkerSize  float64

	// Spacing
	Padding     float64
	MarginTop   float64
	MarginRight float64
	MarginBottom float64
	MarginLeft  float64

	// Border radius for bars, boxes, etc.
	BarRadius   float64
	CardRadius  float64
}

// New creates a new theme from design tokens
func New(tokens *design.DesignTokens) *Theme {
	if tokens == nil {
		tokens = design.DefaultTheme()
	}

	return &Theme{
		Tokens:      tokens,
		ColorScheme: colorSchemeFromTokens(tokens),
		Typography:  typographyFromTokens(tokens),
		Chart:       chartStyleFromTokens(tokens),
	}
}

// colorSchemeFromTokens generates a color scheme from design tokens
func colorSchemeFromTokens(tokens *design.DesignTokens) ColorScheme {
	scheme := ColorScheme{
		TextColor:       tokens.Color,
		BackgroundColor: tokens.Background,
	}

	// Generate palettes based on accent color and mode
	if tokens.Mode == "dark" {
		scheme.Sequential = DarkSequential(tokens.Accent)
		scheme.Diverging = DarkDiverging(tokens.Accent)
		scheme.Categorical = DarkCategorical()
		scheme.GridColor = "#374151"
		scheme.AxisColor = "#4B5563"
		scheme.BorderColor = "#4B5563"
	} else {
		scheme.Sequential = LightSequential(tokens.Accent)
		scheme.Diverging = LightDiverging(tokens.Accent)
		scheme.Categorical = LightCategorical()
		scheme.GridColor = "#E5E7EB"
		scheme.AxisColor = "#9CA3AF"
		scheme.BorderColor = "#D1D5DB"
	}

	return scheme
}

// typographyFromTokens generates typography from design tokens
func typographyFromTokens(tokens *design.DesignTokens) Typography {
	fontFamily := tokens.FontFamily
	if fontFamily == "" {
		fontFamily = "system-ui, -apple-system, sans-serif"
	}

	// Use density to adjust sizes
	scale := 1.0
	if tokens.Density == "compact" {
		scale = 0.9
	} else if tokens.Density == "comfortable" {
		scale = 1.0
	}

	return Typography{
		TitleFont:  fontFamily,
		BodyFont:   fontFamily,
		MonoFont:   "ui-monospace, monospace",

		TitleSize:    units.Px(24 * scale),
		SubtitleSize: units.Px(18 * scale),
		BodySize:     units.Px(14 * scale),
		CaptionSize:  units.Px(12 * scale),
		LabelSize:    units.Px(12 * scale),

		TitleWeight:  svg.FontWeightBold,
		BodyWeight:   svg.FontWeightNormal,
		LabelWeight:  svg.FontWeightNormal,

		TitleLineHeight: 1.2,
		BodyLineHeight:  1.5,
	}
}

// chartStyleFromTokens generates chart styling from design tokens
func chartStyleFromTokens(tokens *design.DesignTokens) ChartStyle {
	// Use layout tokens for spacing
	layout := tokens.Layout
	if layout == nil {
		layout = design.DefaultLayoutTokens()
	}

	scale := 1.0
	if tokens.Density == "compact" {
		scale = 0.8
	}

	return ChartStyle{
		GridStrokeWidth: 0.5,
		AxisStrokeWidth: 1.5,
		DataStrokeWidth: 2.0 * scale,

		GridOpacity: 0.3,
		FillOpacity: 0.7,

		PointSize:   4.0 * scale,
		MarkerSize:  6.0 * scale,

		Padding:     float64(layout.SpaceM),
		MarginTop:   float64(layout.SpaceL),
		MarginRight: float64(layout.SpaceL),
		MarginBottom: float64(layout.SpaceXL),
		MarginLeft:  float64(layout.SpaceXL),

		BarRadius:   float64(tokens.Radius / 4),
		CardRadius:  float64(tokens.Radius),
	}
}

// GetColor returns a color from the categorical palette at the given index
func (t *Theme) GetColor(index int) string {
	colors := t.ColorScheme.Categorical
	if len(colors) == 0 {
		return t.Tokens.Accent
	}
	return colors[index%len(colors)]
}

// GetSequentialColor returns a color from the sequential palette (0.0 to 1.0)
func (t *Theme) GetSequentialColor(t_value float64) string {
	colors := t.ColorScheme.Sequential
	if len(colors) == 0 {
		return t.Tokens.Accent
	}

	// Clamp t to [0, 1]
	if t_value < 0 {
		t_value = 0
	}
	if t_value > 1 {
		t_value = 1
	}

	// Map t to color index
	index := int(t_value * float64(len(colors)-1))
	if index >= len(colors) {
		index = len(colors) - 1
	}

	return colors[index]
}

// GetDivergingColor returns a color from the diverging palette (-1.0 to 1.0)
func (t *Theme) GetDivergingColor(t_value float64) string {
	colors := t.ColorScheme.Diverging
	if len(colors) == 0 {
		return t.Tokens.Accent
	}

	// Map t from [-1, 1] to [0, 1]
	normalized := (t_value + 1.0) / 2.0

	// Clamp to [0, 1]
	if normalized < 0 {
		normalized = 0
	}
	if normalized > 1 {
		normalized = 1
	}

	// Map to color index
	index := int(normalized * float64(len(colors)-1))
	if index >= len(colors) {
		index = len(colors) - 1
	}

	return colors[index]
}

// TitleStyle returns SVG style for chart titles
func (t *Theme) TitleStyle() svg.Style {
	return svg.Style{
		FontFamily: t.Typography.TitleFont,
		FontSize:   t.Typography.TitleSize,
		FontWeight: t.Typography.TitleWeight,
		Fill:       t.ColorScheme.TextColor,
	}
}

// BodyStyle returns SVG style for body text
func (t *Theme) BodyStyle() svg.Style {
	return svg.Style{
		FontFamily: t.Typography.BodyFont,
		FontSize:   t.Typography.BodySize,
		FontWeight: t.Typography.BodyWeight,
		Fill:       t.ColorScheme.TextColor,
	}
}

// LabelStyle returns SVG style for axis labels
func (t *Theme) LabelStyle() svg.Style {
	return svg.Style{
		FontFamily: t.Typography.BodyFont,
		FontSize:   t.Typography.LabelSize,
		FontWeight: t.Typography.LabelWeight,
		Fill:       t.ColorScheme.TextColor,
	}
}

// GridStyle returns SVG style for grid lines
func (t *Theme) GridStyle() svg.Style {
	return svg.Style{
		Stroke:      t.ColorScheme.GridColor,
		StrokeWidth: t.Chart.GridStrokeWidth,
		Opacity:     t.Chart.GridOpacity,
	}
}

// AxisStyle returns SVG style for axis lines
func (t *Theme) AxisStyle() svg.Style {
	return svg.Style{
		Stroke:      t.ColorScheme.AxisColor,
		StrokeWidth: t.Chart.AxisStrokeWidth,
	}
}
