// Package charts provides high-level charting and graphing APIs.
//
// This package builds on top of the layout and rendering engines to provide
// ready-to-use chart types with sensible defaults and extensive customization.
//
// Chart Types:
//   - Line graphs (with area fill, gradients, markers)
//   - Bar charts (vertical, horizontal, stacked)
//   - Scatter plots (multiple series, custom markers)
//   - Heatmaps (linear and GitHub-style weeks view)
//   - Pie charts (with donut mode)
//
// All charts can be rendered to:
//   - SVG (via render/svg/)
//   - Terminal (via render/terminal/)
//   - PNG/JPEG (via export/ from SVG)
//
// Example - Time-series line chart:
//
//	data := []charts.TimeSeriesData{
//	    {Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), Value: 100},
//	    {Date: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), Value: 150},
//	}
//
//	config := charts.LineChartConfig{
//	    Width: 800,
//	    Height: 400,
//	    Title: "Sales Over Time",
//	    UseGradient: true,
//	}
//
//	svgChart := charts.RenderLineChart(data, config)
//	termChart := charts.RenderLineChartTerminal(data, config)
//
// Example - Bar chart with design tokens:
//
//	data := []charts.DataPoint{
//	    {Label: "Q1", Value: 100},
//	    {Label: "Q2", Value: 150},
//	    {Label: "Q3", Value: 120},
//	}
//
//	theme := design.MidnightTheme()
//	config := charts.BarChartConfig{
//	    Width: 600,
//	    Height: 400,
//	    Colors: theme.Colors.Chart,
//	}
//
//	svgChart := charts.RenderBarChart(data, config)
//
// This package uses:
//   - layout/ - For positioning chart elements
//   - render/svg/ or render/terminal/ - For output
//   - design/ (optional) - For consistent styling via design tokens
package charts
