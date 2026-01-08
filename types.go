package dataviz

import (
	"time"

	"github.com/SCKelemen/color"
	design "github.com/SCKelemen/design-system"
)

// TimeSeriesData represents data points over time
type TimeSeriesData struct {
	Date  time.Time
	Value int
}

// ContributionDay represents a single day's contribution data
type ContributionDay struct {
	Date  time.Time
	Count int
}

// HeatmapData represents data for a heatmap visualization
type HeatmapData struct {
	Days      []ContributionDay
	StartDate time.Time
	EndDate   time.Time
	Type      string // "linear" or "weeks" (GitHub-style grid)
}

// LineGraphData represents data for a line graph
type LineGraphData struct {
	Points      []TimeSeriesData
	Color       string
	FillColor   string
	UseGradient bool                // If true, use gradient fill instead of solid color
	GradientID  string              // Optional custom gradient ID (auto-generated if empty)
	ColorSpace  color.GradientSpace // Color space for gradient interpolation
	Label       string
	Smooth      bool    // If true, use smooth curves instead of straight lines
	Tension     float64 // Curve tension (0-1), only used if Smooth is true. 0.3 is recommended
	MarkerType  string  // Marker type: "circle", "square", "diamond", "triangle", "dot", "" (none)
	MarkerSize  float64 // Size of markers in pixels (default: 3)
}

// BarChartData represents data for a bar chart
type BarChartData struct {
	Bars    []BarData
	Color   string
	Label   string
	Stacked bool
}

// BarData represents a single bar or stack
type BarData struct {
	Value     int
	Secondary int // For stacked bars
	Label     string
}

// StatCardData represents data for a statistics card
type StatCardData struct {
	Title       string
	Value       string
	Subtitle    string
	Change      int
	ChangePct   float64
	Color       string
	TrendData   []TimeSeriesData // Optional trend data for mini graph
	TrendColor  string           // Primary color for trend (lighter)
	TrendColor2 string           // Secondary color for trend (darker)
	Legend1     string           // Label for first legend item
	Legend2     string           // Label for second legend item
}

// AreaChartData represents data for an area chart
type AreaChartData struct {
	Points      []TimeSeriesData
	Color       string
	FillColor   string
	UseGradient bool                // If true, use gradient fill
	GradientID  string              // Optional custom gradient ID
	ColorSpace  color.GradientSpace // Color space for gradient interpolation
	Label       string
	Smooth      bool    // If true, use smooth curves
	Tension     float64 // Curve tension (0-1), 0.3 recommended
	BaselineY   int     // Y value for baseline (default: bottom of chart)
	Stacked     bool    // For multiple series (future enhancement)
}

// ScatterPlotData represents data for a scatter plot
type ScatterPlotData struct {
	Points     []ScatterPoint
	Color      string
	Label      string
	MarkerType string  // Marker shape: "circle", "square", "diamond", "triangle", "cross", "x", "dot"
	MarkerSize float64 // Size of markers in pixels
}

// ScatterPoint represents a single point in a scatter plot
type ScatterPoint struct {
	Date  time.Time
	Value int
	Size  float64 // Optional: custom size for this point (0 = use default)
	Label string  // Optional: label for this specific point
}

// Card represents a visual card container
type Card struct {
	Width        int
	Height       int
	Title        string
	Icon         string       // Optional SVG icon content
	Legends      []LegendItem // Optional legend items for header end
	Footer       string       // Optional footer content
	DesignTokens *design.DesignTokens
	MotionTokens *design.MotionTokens
}

// LegendItem represents a legend item in the header
type LegendItem struct {
	Color string
	Label string
	X     int // Optional X position (if 0, will be auto-positioned)
}

// Bounds represents the rectangular area for rendering
type Bounds struct {
	X      int
	Y      int
	Width  int
	Height int
}

// RenderConfig contains configuration for rendering visualizations
type RenderConfig struct {
	DesignTokens *design.DesignTokens
	MotionTokens *design.MotionTokens
	Color        string
	Theme        string
}

// Aspect ratios for different graph types
// These define the width-to-height ratio for graphs to maintain consistent sizing
const (
	LineGraphAspectRatio     = 2.25 // 450px width / 200px height = 2.25:1
	BarChartAspectRatio      = 2.25 // 450px width / 200px height = 2.25:1
	WeeksHeatmapAspectRatio  = 3.89 // 700px width / 180px height ≈ 3.89:1
	LinearHeatmapAspectRatio = 5.71 // 400px width / 70px height ≈ 5.71:1
)
