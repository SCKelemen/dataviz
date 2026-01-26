package gallery

import (
	"fmt"
	"strings"

	"github.com/SCKelemen/svg"
)

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
		GetDefaultTitleStyle()))
	content.WriteString("\n")

	// 3. Render each variant
	labelStyle := GetDefaultLabelStyle()

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
			Tokens:      DefaultTokens(),
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

	return WrapSVG(content.String(), int(dims.TotalWidth), int(dims.TotalHeight)), nil
}
