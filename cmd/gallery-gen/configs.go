package main

import (
	"github.com/SCKelemen/dataviz/charts"
	design "github.com/SCKelemen/design-system"
)

// GalleryRegistry contains all gallery configurations
var GalleryRegistry = map[string]GalleryConfig{
	"bar":          BarGallery,
	"area":         AreaGallery,
	"stacked-area": StackedAreaGallery,
	"lollipop":     LollipopGallery,
	"histogram":    HistogramGallery,
}

// BarGallery defines the bar chart gallery configuration
var BarGallery = GalleryConfig{
	Name:  "bar",
	Title: "Bar Chart Gallery",
	Layout: &SingleRowLayout{
		Cols:       2,
		BaseWidth:  850,
		BaseHeight: 450,
		StartX:     50.0,
	},
	Variants: []VariantConfig{
		{
			Label: "Simple Bars",
			DataProvider: func() interface{} {
				return charts.BarChartData{
					Label: "Sales",
					Color: "#3b82f6",
					Bars: []charts.BarData{
						{Label: "Q1", Value: 45},
						{Label: "Q2", Value: 60},
						{Label: "Q3", Value: 55},
						{Label: "Q4", Value: 70},
					},
				}
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				barData := data.(charts.BarChartData)
				chartW := int(ctx.ChartWidth)
				chartH := int(ctx.ChartHeight - 100)
				tokens := design.DefaultTheme()
				return charts.RenderBarChart(barData, 0, 0, chartW, chartH, tokens)
			},
		},
		{
			Label: "Stacked Bars (Open/Closed)",
			DataProvider: func() interface{} {
				return charts.BarChartData{
					Label:   "Tickets",
					Color:   "#10b981",
					Stacked: true,
					Bars: []charts.BarData{
						{Label: "Mon", Value: 30, Secondary: 20},
						{Label: "Tue", Value: 45, Secondary: 25},
						{Label: "Wed", Value: 40, Secondary: 30},
						{Label: "Thu", Value: 50, Secondary: 15},
					},
				}
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				barData := data.(charts.BarChartData)
				chartW := int(ctx.ChartWidth)
				chartH := int(ctx.ChartHeight - 100)
				tokens := design.DefaultTheme()
				return charts.RenderBarChart(barData, 0, 0, chartW, chartH, tokens)
			},
		},
	},
	ChartOffsetX: 50.0,
	ChartOffsetY: 30.0,
}

// AreaGallery defines the area chart gallery configuration
var AreaGallery = GalleryConfig{
	Name:  "area",
	Title: "Area Chart Gallery",
	Layout: &SingleRowLayout{
		Cols:       2,
		BaseWidth:  800,
		BaseHeight: 400,
		StartX:     0.0,
	},
	Variants: []VariantConfig{
		{
			Label: "Simple Area",
			DataProvider: func() interface{} {
				return charts.AreaChartData{
					Label: "Sales",
					Color: "#3b82f6",
					Points: []charts.TimeSeriesData{
						{Date: mustParseTime("2024-01-01"), Value: 100},
						{Date: mustParseTime("2024-02-01"), Value: 120},
						{Date: mustParseTime("2024-03-01"), Value: 110},
						{Date: mustParseTime("2024-04-01"), Value: 140},
						{Date: mustParseTime("2024-05-01"), Value: 130},
						{Date: mustParseTime("2024-06-01"), Value: 150},
					},
				}
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				areaData := data.(charts.AreaChartData)
				chartW := int(ctx.ChartWidth)
				chartH := int(ctx.ChartHeight)
				tokens := design.DefaultTheme()
				return charts.RenderAreaChart(areaData, 0, 0, chartW, chartH, tokens)
			},
		},
		{
			Label: "Different Color",
			DataProvider: func() interface{} {
				return charts.AreaChartData{
					Label: "Sales",
					Color: "#10b981",
					Points: []charts.TimeSeriesData{
						{Date: mustParseTime("2024-01-01"), Value: 100},
						{Date: mustParseTime("2024-02-01"), Value: 120},
						{Date: mustParseTime("2024-03-01"), Value: 110},
						{Date: mustParseTime("2024-04-01"), Value: 140},
						{Date: mustParseTime("2024-05-01"), Value: 130},
						{Date: mustParseTime("2024-06-01"), Value: 150},
					},
				}
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				areaData := data.(charts.AreaChartData)
				chartW := int(ctx.ChartWidth)
				chartH := int(ctx.ChartHeight)
				tokens := design.DefaultTheme()
				return charts.RenderAreaChart(areaData, 0, 0, chartW, chartH, tokens)
			},
		},
	},
	ChartOffsetX: 10.0,
	ChartOffsetY: 30.0,
}

// StackedAreaGallery defines the stacked area chart gallery configuration
var StackedAreaGallery = GalleryConfig{
	Name:  "stacked-area",
	Title: "Stacked Area Gallery",
	Layout: &SingleRowLayout{
		Cols:       2,
		BaseWidth:  800,
		BaseHeight: 400,
		StartX:     25.0,
	},
	Variants: []VariantConfig{
		{
			Label: "Standard Stacked",
			DataProvider: func() interface{} {
				return struct {
					Series []charts.StackedAreaSeries
					Points []charts.StackedAreaPoint
				}{
					Series: []charts.StackedAreaSeries{
						{Label: "Series A", Color: "#3b82f6"},
						{Label: "Series B", Color: "#10b981"},
						{Label: "Series C", Color: "#f59e0b"},
					},
					Points: []charts.StackedAreaPoint{
						{X: 0, Values: []float64{10, 15, 5}},
						{X: 1, Values: []float64{20, 10, 15}},
						{X: 2, Values: []float64{15, 20, 10}},
						{X: 3, Values: []float64{25, 15, 10}},
						{X: 4, Values: []float64{20, 25, 15}},
					},
				}
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				d := data.(struct {
					Series []charts.StackedAreaSeries
					Points []charts.StackedAreaPoint
				})
				chartW := int(ctx.ChartWidth)
				chartH := int(ctx.ChartHeight)
				spec := charts.StackedAreaSpec{
					Points: d.Points,
					Series: d.Series,
					Width:  float64(chartW),
					Height: float64(chartH),
				}
				return charts.RenderStackedArea(spec)
			},
		},
		{
			Label: "Smooth Curves",
			DataProvider: func() interface{} {
				return struct {
					Series []charts.StackedAreaSeries
					Points []charts.StackedAreaPoint
				}{
					Series: []charts.StackedAreaSeries{
						{Label: "Series A", Color: "#3b82f6"},
						{Label: "Series B", Color: "#10b981"},
						{Label: "Series C", Color: "#f59e0b"},
					},
					Points: []charts.StackedAreaPoint{
						{X: 0, Values: []float64{10, 15, 5}},
						{X: 1, Values: []float64{20, 10, 15}},
						{X: 2, Values: []float64{15, 20, 10}},
						{X: 3, Values: []float64{25, 15, 10}},
						{X: 4, Values: []float64{20, 25, 15}},
					},
				}
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				d := data.(struct {
					Series []charts.StackedAreaSeries
					Points []charts.StackedAreaPoint
				})
				chartW := int(ctx.ChartWidth)
				chartH := int(ctx.ChartHeight)
				spec := charts.StackedAreaSpec{
					Points: d.Points,
					Series: d.Series,
					Width:  float64(chartW),
					Height: float64(chartH),
					Smooth: true,
				}
				return charts.RenderStackedArea(spec)
			},
		},
	},
	ChartOffsetX: 10.0,
	ChartOffsetY: 30.0,
}

// LollipopGallery defines the lollipop chart gallery configuration
var LollipopGallery = GalleryConfig{
	Name:  "lollipop",
	Title: "Lollipop Chart Gallery",
	Layout: &SingleRowLayout{
		Cols:       2,
		BaseWidth:  600,
		BaseHeight: 400,
		StartX:     25.0,
	},
	Variants: []VariantConfig{
		{
			Label: "Vertical Lollipop",
			DataProvider: func() interface{} {
				return &charts.LollipopData{
					Values: []charts.LollipopPoint{
						{Label: "Product A", Value: 45},
						{Label: "Product B", Value: 62},
						{Label: "Product C", Value: 38},
						{Label: "Product D", Value: 71},
						{Label: "Product E", Value: 54},
					},
					Color: "#3b82f6",
				}
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				lollipopData := data.(*charts.LollipopData)
				chartW := int(ctx.ChartWidth - 50)
				chartH := int(ctx.ChartHeight - 80)
				spec := charts.LollipopSpec{
					Data:       lollipopData,
					Width:      float64(chartW),
					Height:     float64(chartH),
					ShowLabels: true,
				}
				return charts.RenderLollipop(spec)
			},
		},
		{
			Label: "Horizontal Lollipop",
			DataProvider: func() interface{} {
				return &charts.LollipopData{
					Values: []charts.LollipopPoint{
						{Label: "Product A", Value: 45},
						{Label: "Product B", Value: 62},
						{Label: "Product C", Value: 38},
						{Label: "Product D", Value: 71},
						{Label: "Product E", Value: 54},
					},
					Color: "#3b82f6",
				}
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				lollipopData := data.(*charts.LollipopData)
				chartW := int(ctx.ChartWidth - 50)
				chartH := int(ctx.ChartHeight - 80)
				spec := charts.LollipopSpec{
					Data:       lollipopData,
					Width:      float64(chartW),
					Height:     float64(chartH),
					Horizontal: true,
					ShowLabels: true,
				}
				return charts.RenderLollipop(spec)
			},
		},
	},
	ChartOffsetX: 0.0,
	ChartOffsetY: 30.0,
}

// HistogramGallery defines the histogram gallery configuration
var HistogramGallery = GalleryConfig{
	Name:  "histogram",
	Title: "Histogram Gallery",
	Layout: &SingleRowLayout{
		Cols:       2,
		BaseWidth:  600,
		BaseHeight: 400,
		StartX:     25.0,
	},
	Variants: []VariantConfig{
		{
			Label: "Count Histogram",
			DataProvider: func() interface{} {
				return generateHistogramData()
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				histData := data.(*charts.HistogramData)
				chartW := int(ctx.ChartWidth - 50)
				chartH := int(ctx.ChartHeight - 80)
				spec := charts.HistogramSpec{
					Data:     histData,
					Width:    float64(chartW),
					Height:   float64(chartH),
					BinCount: 20,
				}
				return charts.RenderHistogram(spec)
			},
		},
		{
			Label: "Density Histogram",
			DataProvider: func() interface{} {
				return generateHistogramData()
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				histData := data.(*charts.HistogramData)
				chartW := int(ctx.ChartWidth - 50)
				chartH := int(ctx.ChartHeight - 80)
				spec := charts.HistogramSpec{
					Data:        histData,
					Width:       float64(chartW),
					Height:      float64(chartH),
					BinCount:    20,
					ShowDensity: true,
				}
				return charts.RenderHistogram(spec)
			},
		},
	},
	ChartOffsetX: 0.0,
	ChartOffsetY: 30.0,
}

// Helper function to generate histogram data
func generateHistogramData() *charts.HistogramData {
	values := make([]float64, 200)
	for i := range values {
		// Simple approximation of normal distribution
		sum := 0.0
		for j := 0; j < 12; j++ {
			sum += float64(i % 20)
		}
		values[i] = sum/12*5 + 50 + float64((i%10)-5)*2
	}
	return &charts.HistogramData{Values: values}
}
