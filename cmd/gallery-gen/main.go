package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/SCKelemen/dataviz/charts"
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

// CalculateGridDimensions calculates pixel dimensions for a grid-based gallery
// using relative units that resolve to exact pixels at render time
func CalculateGridDimensions(cols, rows int, baseWidth, baseHeight float64) GalleryDimensions {
	// Use percentages for grid sizing to avoid accumulation errors
	colPct := units.Percent(100.0 / float64(cols))
	rowPct := units.Percent(100.0 / float64(rows))

	// Calculate dimensions with proper margins
	titleMargin := units.Percent(5)  // 5% top margin for title
	bottomMargin := units.Percent(3) // 3% bottom margin
	chartPadding := units.Percent(2) // 2% padding within each cell

	totalWidth := baseWidth * float64(cols)
	totalHeight := baseHeight * float64(rows)

	// Add margins to total height
	titleSpace := titleMargin.Of(totalHeight)
	bottomSpace := bottomMargin.Of(totalHeight)
	totalHeight += titleSpace + bottomSpace

	// Calculate chart dimensions (subtract padding)
	colWidth := colPct.Of(totalWidth)
	rowHeight := rowPct.Of(baseHeight * float64(rows))

	padding := chartPadding.Of(colWidth)
	chartWidth := colWidth - (padding * 2)
	chartHeight := rowHeight - (padding * 2)

	return GalleryDimensions{
		TotalWidth:   totalWidth,
		TotalHeight:  totalHeight,
		ChartWidth:   chartWidth,
		ChartHeight:  chartHeight,
		ColWidth:     colWidth,
		RowHeight:    rowHeight,
		TitleY:       titleSpace * 0.7, // Position title 70% down the title space
		ChartStartY:  titleSpace,
		BottomMargin: bottomSpace,
	}
}

// CalculateSingleRowDimensions calculates dimensions for single-row galleries
func CalculateSingleRowDimensions(cols int, baseWidth, baseHeight float64) GalleryDimensions {
	titleHeight := 50.0
	bottomMargin := 30.0
	chartPadding := 25.0

	totalWidth := baseWidth * float64(cols)
	totalHeight := baseHeight + titleHeight + bottomMargin

	colWidth := totalWidth / float64(cols)
	chartWidth := baseWidth - (chartPadding * 2)
	chartHeight := baseHeight - chartPadding

	return GalleryDimensions{
		TotalWidth:   totalWidth,
		TotalHeight:  totalHeight,
		ChartWidth:   chartWidth,
		ChartHeight:  chartHeight,
		ColWidth:     colWidth,
		RowHeight:    baseHeight,
		TitleY:       30,
		ChartStartY:  titleHeight + 10,
		BottomMargin: bottomMargin,
	}
}

func main() {
	if err := generateGalleries(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// wrapSVG wraps content in an SVG element with proper xmlns and viewBox
func wrapSVG(content string, width, height int) string {
	return fmt.Sprintf(
		`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">%s</svg>`,
		width, height, width, height, content,
	)
}

func generateGalleries() error {
	outputDir := "examples-gallery"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	// Generate all galleries from the registry
	for name, config := range GalleryRegistry {
		fmt.Printf("Generating %s gallery...\n", name)

		svg, err := GenerateGallery(config)
		if err != nil {
			fmt.Printf("  ✗ Failed: %v\n", err)
			continue
		}

		outputPath := filepath.Join(outputDir, name+"-gallery.svg")
		if err := os.WriteFile(outputPath, []byte(svg), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", outputPath, err)
		}
		fmt.Printf("  ✓ %s\n", outputPath)
	}

	fmt.Println("✓ Gallery generation complete!")
	return nil
}

// Helper functions

func mustParseTime(s string) time.Time {
	t, _ := time.Parse("2006-01-02", s)
	return t
}

func createSampleTree() *charts.TreeNode {
	return &charts.TreeNode{
		Name:  "Root",
		Value: 100,
		Children: []*charts.TreeNode{
			{
				Name:  "Branch A",
				Value: 40,
				Children: []*charts.TreeNode{
					{Name: "Leaf A1", Value: 15},
					{Name: "Leaf A2", Value: 12},
					{Name: "Leaf A3", Value: 13},
				},
			},
			{
				Name:  "Branch B",
				Value: 35,
				Children: []*charts.TreeNode{
					{Name: "Leaf B1", Value: 20},
					{Name: "Leaf B2", Value: 15},
				},
			},
			{
				Name:  "Branch C",
				Value: 25,
				Children: []*charts.TreeNode{
					{Name: "Leaf C1", Value: 10},
					{Name: "Leaf C2", Value: 8},
					{Name: "Leaf C3", Value: 7},
				},
			},
		},
	}
}

