package main

import (
	"fmt"
	"strings"

	design "github.com/SCKelemen/design-system"
	"github.com/SCKelemen/svg"
	"github.com/SCKelemen/units"
)

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
	Name  string
	Title string
	Layout LayoutStrategy
	Variants []VariantConfig

	// Styling (use defaults if not specified)
	LabelOffsetY float64 // Default: 0.0
	ChartOffsetX float64 // Default: 0.0
	ChartOffsetY float64 // Default: 30.0
}

// GenerateGallery creates an SVG gallery from a configuration
func GenerateGallery(config GalleryConfig) (string, error) {
	dims := config.Layout.CalculateDimensions()

	var content strings.Builder

	// 1. Background
	content.WriteString(svg.Rect(0, 0, dims.TotalWidth, dims.TotalHeight,
		svg.Style{Fill: "#ffffff"}))
	content.WriteString("\n")

	// 2. Title
	content.WriteString(svg.Text(config.Title, dims.TotalWidth/2, dims.TitleY,
		getDefaultTitleStyle()))
	content.WriteString("\n")

	// 3. Render each variant
	labelStyle := getDefaultLabelStyle()

	for i, variant := range config.Variants {
		cellX, cellY := config.Layout.GetCellPosition(i)

		// Get data for this variant
		data := variant.DataProvider()

		// Create render context
		ctx := RenderContext{
			ChartWidth:  dims.ChartWidth,
			ChartHeight: dims.ChartHeight,
			OffsetX:     config.ChartOffsetX,
			OffsetY:     config.ChartOffsetY,
			Tokens:      design.DefaultTheme(),
		}

		// Render chart
		chartContent := variant.ChartRenderer(data, ctx)

		// Build variant group
		content.WriteString(svg.Group(
			svg.Text(variant.Label, dims.ColWidth/2, config.LabelOffsetY, labelStyle)+
				svg.Group(
					chartContent,
					fmt.Sprintf("translate(%.2f, %.2f)", ctx.OffsetX, ctx.OffsetY),
					svg.Style{},
				),
			fmt.Sprintf("translate(%.2f, %.2f)", cellX, cellY),
			svg.Style{},
		))
		content.WriteString("\n")
	}

	return wrapSVG(content.String(), int(dims.TotalWidth), int(dims.TotalHeight)), nil
}

// getDefaultTitleStyle returns the standard gallery title style
func getDefaultTitleStyle() svg.Style {
	return svg.Style{
		FontSize:   units.Px(20),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#000000",
		TextAnchor: "middle",
	}
}

// getDefaultLabelStyle returns the standard variant label style
func getDefaultLabelStyle() svg.Style {
	return svg.Style{
		FontSize:   units.Px(14),
		FontWeight: "bold",
		FontFamily: "sans-serif",
		Fill:       "#666",
		TextAnchor: "middle",
	}
}
