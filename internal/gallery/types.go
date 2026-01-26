package gallery

import (
	"fmt"

	design "github.com/SCKelemen/design-system"
	"github.com/SCKelemen/svg"
	"github.com/SCKelemen/units"
)

// GalleryDimensions holds calculated dimensions for gallery layouts
type GalleryDimensions struct {
	TotalWidth   float64
	TotalHeight  float64
	ChartWidth   float64
	ChartHeight  float64
	ColWidth     float64
	RowHeight    float64
	TitleY       float64
	ChartStartY  float64
	BottomMargin float64
}

// LayoutStrategy defines how a gallery calculates dimensions and positions variants
type LayoutStrategy interface {
	CalculateDimensions() GalleryDimensions
	GetCellPosition(variantIndex int) (x, y float64)
}

// VariantConfig defines a single chart variant in a gallery
type VariantConfig struct {
	Label         string
	DataProvider  func() interface{}                          // Returns chart-specific data
	ChartRenderer func(data interface{}, ctx RenderContext) string // Renders the chart
}

// RenderContext provides rendering parameters to chart renderers
type RenderContext struct {
	ChartWidth  float64
	ChartHeight float64
	OffsetX     float64 // Chart offset within cell
	OffsetY     float64
	Tokens      interface{} // design.DefaultTheme() result
}

// GalleryConfig defines a complete gallery configuration
type GalleryConfig struct {
	Name     string
	Title    string
	Layout   LayoutStrategy
	Variants []VariantConfig

	// Styling (use defaults if not specified)
	LabelOffsetY float64 // Default: 0.0
	ChartOffsetX float64 // Default: 0.0
	ChartOffsetY float64 // Default: 30.0
}

// GetDefaultTitleStyle returns the standard gallery title style
func GetDefaultTitleStyle() svg.Style {
	return svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
}

// GetDefaultLabelStyle returns the standard variant label style
func GetDefaultLabelStyle() svg.Style {
	return svg.Style{
		FontSize:   units.Px(14),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#666",
		TextAnchor: "middle",
	}
}

// WrapSVG wraps content in an SVG element with proper xmlns and viewBox
func WrapSVG(content string, width, height int) string {
	return fmt.Sprintf(
		`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">%s</svg>`,
		width, height, width, height, content,
	)
}

// DefaultTokens returns the default design tokens
func DefaultTokens() interface{} {
	return design.DefaultTheme()
}
