package main

import (
	"github.com/SCKelemen/dataviz/charts"
	design "github.com/SCKelemen/design-system"
)

// GalleryRegistry contains all gallery configurations
var GalleryRegistry = map[string]GalleryConfig{
	"bar":               BarGallery,
	"area":              AreaGallery,
	"stacked-area":      StackedAreaGallery,
	"lollipop":          LollipopGallery,
	"histogram":         HistogramGallery,
	"pie":               PieGallery,
	"boxplot":           BoxPlotGallery,
	"violin":            ViolinGallery,
	"treemap":           TreemapGallery,
	"icicle":            IcicleGallery,
	"ridgeline":         RidgelineGallery,
	"line":              LineGallery,
	"scatter":           ScatterGallery,
	"connected-scatter": ConnectedScatterGallery,
	"statcard":          StatCardGallery,
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

// PieGallery defines the pie chart gallery configuration
var PieGallery = GalleryConfig{
	Name:  "pie",
	Title: "Pie Chart Gallery",
	Layout: &SingleRowLayout{
		Cols:       3,
		BaseWidth:  800,
		BaseHeight: 350,
		StartX:     0.0,
	},
	Variants: []VariantConfig{
		{
			Label: "Regular Pie",
			DataProvider: func() interface{} {
				return charts.PieChartData{
					Slices: []charts.PieSlice{
						{Label: "Chrome", Value: 63.5},
						{Label: "Safari", Value: 19.3},
						{Label: "Firefox", Value: 9.2},
						{Label: "Edge", Value: 5.1},
						{Label: "Other", Value: 2.9},
					},
				}
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				pieData := data.(charts.PieChartData)
				chartW := int(ctx.ChartWidth)
				chartH := int(ctx.ChartHeight - 70)
				return charts.RenderPieChart(pieData, 0, 0, chartW, chartH, "", false, true, true)
			},
		},
		{
			Label: "Donut Chart",
			DataProvider: func() interface{} {
				return charts.PieChartData{
					Slices: []charts.PieSlice{
						{Label: "Chrome", Value: 63.5},
						{Label: "Safari", Value: 19.3},
						{Label: "Firefox", Value: 9.2},
						{Label: "Edge", Value: 5.1},
						{Label: "Other", Value: 2.9},
					},
				}
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				pieData := data.(charts.PieChartData)
				chartW := int(ctx.ChartWidth)
				chartH := int(ctx.ChartHeight - 70)
				return charts.RenderPieChart(pieData, 0, 0, chartW, chartH, "", true, true, true)
			},
		},
		{
			Label: "Custom Colors",
			DataProvider: func() interface{} {
				return charts.PieChartData{
					Slices: []charts.PieSlice{
						{Label: "Chrome", Value: 63.5},
						{Label: "Safari", Value: 19.3},
						{Label: "Firefox", Value: 9.2},
						{Label: "Edge", Value: 5.1},
						{Label: "Other", Value: 2.9},
					},
					Colors: []string{"#ef4444", "#f97316", "#eab308", "#22c55e", "#3b82f6"},
				}
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				pieData := data.(charts.PieChartData)
				chartW := int(ctx.ChartWidth)
				chartH := int(ctx.ChartHeight - 70)
				return charts.RenderPieChart(pieData, 0, 0, chartW, chartH, "", false, true, true)
			},
		},
	},
	ChartOffsetX: 0.0,
	ChartOffsetY: 20.0,
}

// BoxPlotGallery defines the box plot gallery configuration
var BoxPlotGallery = GalleryConfig{
	Name:  "boxplot",
	Title: "Box Plot Gallery",
	Layout: &SingleRowLayout{
		Cols:       2,
		BaseWidth:  600,
		BaseHeight: 400,
		StartX:     25.0,
	},
	Variants: []VariantConfig{
		{
			Label: "Basic Box Plot",
			DataProvider: func() interface{} {
				return []*charts.BoxPlotData{
					{Label: "Group A", Values: []float64{12, 15, 18, 20, 22, 25, 28, 30, 32, 35, 40, 45}},
					{Label: "Group B", Values: []float64{20, 22, 25, 28, 30, 32, 35, 38, 40, 42, 45, 48, 50}},
					{Label: "Group C", Values: []float64{10, 12, 15, 18, 20, 25, 30, 35, 40, 50, 60}},
				}
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				boxData := data.([]*charts.BoxPlotData)
				chartW := int(ctx.ChartWidth - 50)
				chartH := int(ctx.ChartHeight - 80)
				spec := charts.BoxPlotSpec{
					Data:   boxData,
					Width:  float64(chartW),
					Height: float64(chartH),
				}
				return charts.RenderVerticalBoxPlot(spec)
			},
		},
		{
			Label: "With Notch",
			DataProvider: func() interface{} {
				return []*charts.BoxPlotData{
					{Label: "Group A", Values: []float64{12, 15, 18, 20, 22, 25, 28, 30, 32, 35, 40, 45}},
					{Label: "Group B", Values: []float64{20, 22, 25, 28, 30, 32, 35, 38, 40, 42, 45, 48, 50}},
					{Label: "Group C", Values: []float64{10, 12, 15, 18, 20, 25, 30, 35, 40, 50, 60}},
				}
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				boxData := data.([]*charts.BoxPlotData)
				chartW := int(ctx.ChartWidth - 50)
				chartH := int(ctx.ChartHeight - 80)
				spec := charts.BoxPlotSpec{
					Data:         boxData,
					Width:        float64(chartW),
					Height:       float64(chartH),
					ShowOutliers: true,
					ShowNotch:    true,
				}
				return charts.RenderVerticalBoxPlot(spec)
			},
		},
	},
	ChartOffsetX: 0.0,
	ChartOffsetY: 30.0,
}

// ViolinGallery defines the violin plot gallery configuration
var ViolinGallery = GalleryConfig{
	Name:  "violin",
	Title: "Violin Plot Gallery",
	Layout: &SingleRowLayout{
		Cols:       2,
		BaseWidth:  600,
		BaseHeight: 400,
		StartX:     25.0,
	},
	Variants: []VariantConfig{
		{
			Label: "Basic Violin",
			DataProvider: func() interface{} {
				return []*charts.ViolinPlotData{
					{Label: "Group A", Values: generateViolinValues(25, 5)},
					{Label: "Group B", Values: generateViolinValues(30, 8)},
					{Label: "Group C", Values: generateViolinValues(20, 6)},
				}
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				violinData := data.([]*charts.ViolinPlotData)
				chartW := int(ctx.ChartWidth - 50)
				chartH := int(ctx.ChartHeight - 80)
				spec := charts.ViolinPlotSpec{
					Data:   violinData,
					Width:  float64(chartW),
					Height: float64(chartH),
				}
				return charts.RenderViolinPlot(spec)
			},
		},
		{
			Label: "With Box Plot",
			DataProvider: func() interface{} {
				return []*charts.ViolinPlotData{
					{Label: "Group A", Values: generateViolinValues(25, 5)},
					{Label: "Group B", Values: generateViolinValues(30, 8)},
					{Label: "Group C", Values: generateViolinValues(20, 6)},
				}
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				violinData := data.([]*charts.ViolinPlotData)
				chartW := int(ctx.ChartWidth - 50)
				chartH := int(ctx.ChartHeight - 80)
				spec := charts.ViolinPlotSpec{
					Data:       violinData,
					Width:      float64(chartW),
					Height:     float64(chartH),
					ShowBox:    true,
					ShowMedian: true,
					ShowMean:   true,
				}
				return charts.RenderViolinPlot(spec)
			},
		},
	},
	ChartOffsetX: 0.0,
	ChartOffsetY: 30.0,
}

// TreemapGallery defines the treemap gallery configuration
var TreemapGallery = GalleryConfig{
	Name:  "treemap",
	Title: "Treemap Gallery",
	Layout: &SingleRowLayout{
		Cols:       2,
		BaseWidth:  600,
		BaseHeight: 400,
		StartX:     25.0,
	},
	Variants: []VariantConfig{
		{
			Label: "Standard Treemap",
			DataProvider: func() interface{} {
				return createSampleTree()
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				tree := data.(*charts.TreeNode)
				chartW := int(ctx.ChartWidth - 50)
				chartH := int(ctx.ChartHeight - 80)
				spec := charts.TreemapSpec{
					Root:        tree,
					Width:       float64(chartW),
					Height:      float64(chartH),
					ShowLabels:  true,
					ColorScheme: []string{"#3b82f6", "#10b981", "#f59e0b", "#ef4444", "#8b5cf6"},
				}
				return charts.RenderTreemap(spec)
			},
		},
		{
			Label: "With Padding",
			DataProvider: func() interface{} {
				return createSampleTree()
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				tree := data.(*charts.TreeNode)
				chartW := int(ctx.ChartWidth - 50)
				chartH := int(ctx.ChartHeight - 80)
				spec := charts.TreemapSpec{
					Root:        tree,
					Width:       float64(chartW),
					Height:      float64(chartH),
					Padding:     3,
					ShowLabels:  true,
					ColorScheme: []string{"#6366f1", "#ec4899", "#14b8a6", "#f97316", "#a855f7"},
				}
				return charts.RenderTreemap(spec)
			},
		},
	},
	ChartOffsetX: 0.0,
	ChartOffsetY: 30.0,
}

// IcicleGallery defines the icicle chart gallery configuration
var IcicleGallery = GalleryConfig{
	Name:  "icicle",
	Title: "Icicle Chart Gallery",
	Layout: &SingleRowLayout{
		Cols:       2,
		BaseWidth:  600,
		BaseHeight: 400,
		StartX:     25.0,
	},
	Variants: []VariantConfig{
		{
			Label: "Vertical Icicle",
			DataProvider: func() interface{} {
				return createSampleTree()
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				tree := data.(*charts.TreeNode)
				chartW := ctx.ChartWidth - 50
				chartH := ctx.ChartHeight - 80
				spec := charts.IcicleSpec{
					Root:        tree,
					Width:       chartW,
					Height:      chartH,
					Orientation: "vertical",
					ShowLabels:  true,
					ColorScheme: []string{"#3b82f6", "#10b981", "#f59e0b", "#ef4444", "#8b5cf6"},
				}
				return charts.RenderIcicle(spec)
			},
		},
		{
			Label: "Horizontal Icicle",
			DataProvider: func() interface{} {
				return createSampleTree()
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				tree := data.(*charts.TreeNode)
				chartW := ctx.ChartWidth - 50
				chartH := ctx.ChartHeight - 80
				spec := charts.IcicleSpec{
					Root:        tree,
					Width:       chartW,
					Height:      chartH,
					Orientation: "horizontal",
					ShowLabels:  true,
					ColorScheme: []string{"#6366f1", "#ec4899", "#14b8a6", "#f97316", "#a855f7"},
				}
				return charts.RenderIcicle(spec)
			},
		},
	},
	ChartOffsetX: 0.0,
	ChartOffsetY: 30.0,
}

// RidgelineGallery defines the ridgeline chart gallery configuration
var RidgelineGallery = GalleryConfig{
	Name:  "ridgeline",
	Title: "Ridgeline Gallery",
	Layout: &SingleRowLayout{
		Cols:       2,
		BaseWidth:  600,
		BaseHeight: 400,
		StartX:     25.0,
	},
	Variants: []VariantConfig{
		{
			Label: "Standard Ridgeline",
			DataProvider: func() interface{} {
				return []*charts.RidgelineData{
					{Label: "January", Values: []float64{10, 12, 15, 18, 20, 22, 25, 23, 20, 18, 15, 12}},
					{Label: "February", Values: []float64{15, 18, 20, 22, 25, 28, 30, 28, 25, 22, 20, 18}},
					{Label: "March", Values: []float64{20, 22, 25, 28, 30, 32, 35, 33, 30, 28, 25, 22}},
					{Label: "April", Values: []float64{25, 28, 30, 32, 35, 38, 40, 38, 35, 32, 30, 28}},
				}
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				ridgeData := data.([]*charts.RidgelineData)
				chartW := ctx.ChartWidth - 50
				chartH := ctx.ChartHeight - 80
				spec := charts.RidgelineSpec{
					Data:       ridgeData,
					Width:      chartW,
					Height:     chartH,
					Overlap:    0.5,
					ShowLabels: true,
				}
				return charts.RenderRidgeline(spec)
			},
		},
		{
			Label: "With Fill",
			DataProvider: func() interface{} {
				return []*charts.RidgelineData{
					{Label: "January", Values: []float64{10, 12, 15, 18, 20, 22, 25, 23, 20, 18, 15, 12}},
					{Label: "February", Values: []float64{15, 18, 20, 22, 25, 28, 30, 28, 25, 22, 20, 18}},
					{Label: "March", Values: []float64{20, 22, 25, 28, 30, 32, 35, 33, 30, 28, 25, 22}},
					{Label: "April", Values: []float64{25, 28, 30, 32, 35, 38, 40, 38, 35, 32, 30, 28}},
				}
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				ridgeData := data.([]*charts.RidgelineData)
				chartW := ctx.ChartWidth - 50
				chartH := ctx.ChartHeight - 80
				spec := charts.RidgelineSpec{
					Data:       ridgeData,
					Width:      chartW,
					Height:     chartH,
					Overlap:    0.5,
					ShowFill:   true,
					ShowLabels: true,
				}
				return charts.RenderRidgeline(spec)
			},
		},
	},
	ChartOffsetX: 0.0,
	ChartOffsetY: 30.0,
}

// LineGallery defines the line graph gallery configuration
var LineGallery = GalleryConfig{
	Name:  "line",
	Title: "Line Graph Gallery",
	Layout: &GridLayout{
		Cols:       2,
		Rows:       2,
		BaseWidth:  650,
		BaseHeight: 350,
	},
	Variants: []VariantConfig{
		{
			Label: "Simple Line",
			DataProvider: func() interface{} {
				return charts.LineGraphData{
					Label: "Temperature",
					Color: "#3b82f6",
					Points: []charts.TimeSeriesData{
						{Date: mustParseTime("2024-01-01"), Value: 15},
						{Date: mustParseTime("2024-02-01"), Value: 18},
						{Date: mustParseTime("2024-03-01"), Value: 12},
						{Date: mustParseTime("2024-04-01"), Value: 22},
						{Date: mustParseTime("2024-05-01"), Value: 27},
						{Date: mustParseTime("2024-06-01"), Value: 30},
					},
				}
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				lineData := data.(charts.LineGraphData)
				chartW := int(ctx.ChartWidth)
				chartH := int(ctx.ChartHeight)
				tokens := design.DefaultTheme()
				return charts.RenderLineGraph(lineData, 0, 0, chartW, chartH, tokens)
			},
		},
		{
			Label: "Smoothed",
			DataProvider: func() interface{} {
				return charts.LineGraphData{
					Label:   "Temperature",
					Color:   "#3b82f6",
					Smooth:  true,
					Tension: 0.3,
					Points: []charts.TimeSeriesData{
						{Date: mustParseTime("2024-01-01"), Value: 15},
						{Date: mustParseTime("2024-02-01"), Value: 18},
						{Date: mustParseTime("2024-03-01"), Value: 12},
						{Date: mustParseTime("2024-04-01"), Value: 22},
						{Date: mustParseTime("2024-05-01"), Value: 27},
						{Date: mustParseTime("2024-06-01"), Value: 30},
					},
				}
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				lineData := data.(charts.LineGraphData)
				chartW := int(ctx.ChartWidth)
				chartH := int(ctx.ChartHeight)
				tokens := design.DefaultTheme()
				return charts.RenderLineGraph(lineData, 0, 0, chartW, chartH, tokens)
			},
		},
		{
			Label: "With Markers",
			DataProvider: func() interface{} {
				return charts.LineGraphData{
					Label:      "Temperature",
					Color:      "#3b82f6",
					MarkerType: "circle",
					MarkerSize: 5,
					Points: []charts.TimeSeriesData{
						{Date: mustParseTime("2024-01-01"), Value: 15},
						{Date: mustParseTime("2024-02-01"), Value: 18},
						{Date: mustParseTime("2024-03-01"), Value: 12},
						{Date: mustParseTime("2024-04-01"), Value: 22},
						{Date: mustParseTime("2024-05-01"), Value: 27},
						{Date: mustParseTime("2024-06-01"), Value: 30},
					},
				}
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				lineData := data.(charts.LineGraphData)
				chartW := int(ctx.ChartWidth)
				chartH := int(ctx.ChartHeight)
				tokens := design.DefaultTheme()
				return charts.RenderLineGraph(lineData, 0, 0, chartW, chartH, tokens)
			},
		},
		{
			Label: "Filled Area",
			DataProvider: func() interface{} {
				return charts.LineGraphData{
					Label:     "Temperature",
					Color:     "#3b82f6",
					FillColor: "#3b82f620",
					Points: []charts.TimeSeriesData{
						{Date: mustParseTime("2024-01-01"), Value: 15},
						{Date: mustParseTime("2024-02-01"), Value: 18},
						{Date: mustParseTime("2024-03-01"), Value: 12},
						{Date: mustParseTime("2024-04-01"), Value: 22},
						{Date: mustParseTime("2024-05-01"), Value: 27},
						{Date: mustParseTime("2024-06-01"), Value: 30},
					},
				}
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				lineData := data.(charts.LineGraphData)
				chartW := int(ctx.ChartWidth)
				chartH := int(ctx.ChartHeight)
				tokens := design.DefaultTheme()
				return charts.RenderLineGraph(lineData, 0, 0, chartW, chartH, tokens)
			},
		},
	},
	ChartOffsetX: 10.0,
	ChartOffsetY: 25.0,
}

// ScatterGallery defines the scatter plot gallery configuration
var ScatterGallery = GalleryConfig{
	Name:  "scatter",
	Title: "Scatter Plot Gallery",
	Layout: &GridLayout{
		Cols:       3,
		Rows:       2,
		BaseWidth:  450,
		BaseHeight: 350,
	},
	Variants: []VariantConfig{
		{
			Label: "Marker: circle",
			DataProvider: func() interface{} {
				return charts.ScatterPlotData{
					Points: []charts.ScatterPoint{
						{Label: "A", Date: mustParseTime("2024-01-01"), Value: 55},
						{Label: "B", Date: mustParseTime("2024-02-01"), Value: 78},
						{Label: "C", Date: mustParseTime("2024-03-01"), Value: 44},
						{Label: "D", Date: mustParseTime("2024-04-01"), Value: 66},
						{Label: "E", Date: mustParseTime("2024-05-01"), Value: 33},
						{Label: "F", Date: mustParseTime("2024-06-01"), Value: 77},
						{Label: "G", Date: mustParseTime("2024-07-01"), Value: 22},
						{Label: "H", Date: mustParseTime("2024-08-01"), Value: 88},
					},
					MarkerType: "circle",
					Color:      "#3b82f6",
				}
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				scatterData := data.(charts.ScatterPlotData)
				chartW := int(ctx.ChartWidth - 50)
				chartH := int(ctx.ChartHeight - 60)
				tokens := design.DefaultTheme()
				return charts.RenderScatterPlot(scatterData, 0, 0, chartW, chartH, tokens)
			},
		},
		{
			Label: "Marker: square",
			DataProvider: func() interface{} {
				return charts.ScatterPlotData{
					Points: []charts.ScatterPoint{
						{Label: "A", Date: mustParseTime("2024-01-01"), Value: 55},
						{Label: "B", Date: mustParseTime("2024-02-01"), Value: 78},
						{Label: "C", Date: mustParseTime("2024-03-01"), Value: 44},
						{Label: "D", Date: mustParseTime("2024-04-01"), Value: 66},
						{Label: "E", Date: mustParseTime("2024-05-01"), Value: 33},
						{Label: "F", Date: mustParseTime("2024-06-01"), Value: 77},
						{Label: "G", Date: mustParseTime("2024-07-01"), Value: 22},
						{Label: "H", Date: mustParseTime("2024-08-01"), Value: 88},
					},
					MarkerType: "square",
					Color:      "#3b82f6",
				}
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				scatterData := data.(charts.ScatterPlotData)
				chartW := int(ctx.ChartWidth - 50)
				chartH := int(ctx.ChartHeight - 60)
				tokens := design.DefaultTheme()
				return charts.RenderScatterPlot(scatterData, 0, 0, chartW, chartH, tokens)
			},
		},
		{
			Label: "Marker: diamond",
			DataProvider: func() interface{} {
				return charts.ScatterPlotData{
					Points: []charts.ScatterPoint{
						{Label: "A", Date: mustParseTime("2024-01-01"), Value: 55},
						{Label: "B", Date: mustParseTime("2024-02-01"), Value: 78},
						{Label: "C", Date: mustParseTime("2024-03-01"), Value: 44},
						{Label: "D", Date: mustParseTime("2024-04-01"), Value: 66},
						{Label: "E", Date: mustParseTime("2024-05-01"), Value: 33},
						{Label: "F", Date: mustParseTime("2024-06-01"), Value: 77},
						{Label: "G", Date: mustParseTime("2024-07-01"), Value: 22},
						{Label: "H", Date: mustParseTime("2024-08-01"), Value: 88},
					},
					MarkerType: "diamond",
					Color:      "#3b82f6",
				}
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				scatterData := data.(charts.ScatterPlotData)
				chartW := int(ctx.ChartWidth - 50)
				chartH := int(ctx.ChartHeight - 60)
				tokens := design.DefaultTheme()
				return charts.RenderScatterPlot(scatterData, 0, 0, chartW, chartH, tokens)
			},
		},
		{
			Label: "Marker: triangle",
			DataProvider: func() interface{} {
				return charts.ScatterPlotData{
					Points: []charts.ScatterPoint{
						{Label: "A", Date: mustParseTime("2024-01-01"), Value: 55},
						{Label: "B", Date: mustParseTime("2024-02-01"), Value: 78},
						{Label: "C", Date: mustParseTime("2024-03-01"), Value: 44},
						{Label: "D", Date: mustParseTime("2024-04-01"), Value: 66},
						{Label: "E", Date: mustParseTime("2024-05-01"), Value: 33},
						{Label: "F", Date: mustParseTime("2024-06-01"), Value: 77},
						{Label: "G", Date: mustParseTime("2024-07-01"), Value: 22},
						{Label: "H", Date: mustParseTime("2024-08-01"), Value: 88},
					},
					MarkerType: "triangle",
					Color:      "#3b82f6",
				}
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				scatterData := data.(charts.ScatterPlotData)
				chartW := int(ctx.ChartWidth - 50)
				chartH := int(ctx.ChartHeight - 60)
				tokens := design.DefaultTheme()
				return charts.RenderScatterPlot(scatterData, 0, 0, chartW, chartH, tokens)
			},
		},
		{
			Label: "Marker: cross",
			DataProvider: func() interface{} {
				return charts.ScatterPlotData{
					Points: []charts.ScatterPoint{
						{Label: "A", Date: mustParseTime("2024-01-01"), Value: 55},
						{Label: "B", Date: mustParseTime("2024-02-01"), Value: 78},
						{Label: "C", Date: mustParseTime("2024-03-01"), Value: 44},
						{Label: "D", Date: mustParseTime("2024-04-01"), Value: 66},
						{Label: "E", Date: mustParseTime("2024-05-01"), Value: 33},
						{Label: "F", Date: mustParseTime("2024-06-01"), Value: 77},
						{Label: "G", Date: mustParseTime("2024-07-01"), Value: 22},
						{Label: "H", Date: mustParseTime("2024-08-01"), Value: 88},
					},
					MarkerType: "cross",
					Color:      "#3b82f6",
				}
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				scatterData := data.(charts.ScatterPlotData)
				chartW := int(ctx.ChartWidth - 50)
				chartH := int(ctx.ChartHeight - 60)
				tokens := design.DefaultTheme()
				return charts.RenderScatterPlot(scatterData, 0, 0, chartW, chartH, tokens)
			},
		},
		{
			Label: "Marker: x",
			DataProvider: func() interface{} {
				return charts.ScatterPlotData{
					Points: []charts.ScatterPoint{
						{Label: "A", Date: mustParseTime("2024-01-01"), Value: 55},
						{Label: "B", Date: mustParseTime("2024-02-01"), Value: 78},
						{Label: "C", Date: mustParseTime("2024-03-01"), Value: 44},
						{Label: "D", Date: mustParseTime("2024-04-01"), Value: 66},
						{Label: "E", Date: mustParseTime("2024-05-01"), Value: 33},
						{Label: "F", Date: mustParseTime("2024-06-01"), Value: 77},
						{Label: "G", Date: mustParseTime("2024-07-01"), Value: 22},
						{Label: "H", Date: mustParseTime("2024-08-01"), Value: 88},
					},
					MarkerType: "x",
					Color:      "#3b82f6",
				}
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				scatterData := data.(charts.ScatterPlotData)
				chartW := int(ctx.ChartWidth - 50)
				chartH := int(ctx.ChartHeight - 60)
				tokens := design.DefaultTheme()
				return charts.RenderScatterPlot(scatterData, 0, 0, chartW, chartH, tokens)
			},
		},
	},
	ChartOffsetX: 0.0,
	ChartOffsetY: 25.0,
}

// ConnectedScatterGallery defines the connected scatter gallery configuration
var ConnectedScatterGallery = GalleryConfig{
	Name:  "connected-scatter",
	Title: "Connected Scatter Gallery",
	Layout: &GridLayout{
		Cols:       3,
		Rows:       2,
		BaseWidth:  450,
		BaseHeight: 350,
	},
	Variants: []VariantConfig{
		{
			Label: "Line: Solid",
			DataProvider: func() interface{} {
				return charts.ConnectedScatterSpec{
					Width:  0, // Will be set in renderer
					Height: 0,
					Series: []*charts.ConnectedScatterSeries{
						{
							Points: []charts.ConnectedScatterPoint{
								{X: 0, Y: 10},
								{X: 1, Y: 25},
								{X: 2, Y: 15},
								{X: 3, Y: 30},
								{X: 4, Y: 20},
								{X: 5, Y: 35},
							},
							Color:     "#3b82f6",
							LineStyle: "solid",
						},
					},
					ShowLines:   true,
					ShowMarkers: true,
				}
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				spec := data.(charts.ConnectedScatterSpec)
				spec.Width = ctx.ChartWidth - 50
				spec.Height = ctx.ChartHeight - 80
				return charts.RenderConnectedScatter(spec)
			},
		},
		{
			Label: "Line: Dashed",
			DataProvider: func() interface{} {
				return charts.ConnectedScatterSpec{
					Width:  0,
					Height: 0,
					Series: []*charts.ConnectedScatterSeries{
						{
							Points: []charts.ConnectedScatterPoint{
								{X: 0, Y: 10},
								{X: 1, Y: 25},
								{X: 2, Y: 15},
								{X: 3, Y: 30},
								{X: 4, Y: 20},
								{X: 5, Y: 35},
							},
							Color:     "#3b82f6",
							LineStyle: "dashed",
						},
					},
					ShowLines:   true,
					ShowMarkers: true,
				}
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				spec := data.(charts.ConnectedScatterSpec)
				spec.Width = ctx.ChartWidth - 50
				spec.Height = ctx.ChartHeight - 80
				return charts.RenderConnectedScatter(spec)
			},
		},
		{
			Label: "Line: Dotted",
			DataProvider: func() interface{} {
				return charts.ConnectedScatterSpec{
					Width:  0,
					Height: 0,
					Series: []*charts.ConnectedScatterSeries{
						{
							Points: []charts.ConnectedScatterPoint{
								{X: 0, Y: 10},
								{X: 1, Y: 25},
								{X: 2, Y: 15},
								{X: 3, Y: 30},
								{X: 4, Y: 20},
								{X: 5, Y: 35},
							},
							Color:     "#3b82f6",
							LineStyle: "dotted",
						},
					},
					ShowLines:   true,
					ShowMarkers: true,
				}
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				spec := data.(charts.ConnectedScatterSpec)
				spec.Width = ctx.ChartWidth - 50
				spec.Height = ctx.ChartHeight - 80
				return charts.RenderConnectedScatter(spec)
			},
		},
		{
			Label: "Line: Dash-Dot",
			DataProvider: func() interface{} {
				return charts.ConnectedScatterSpec{
					Width:  0,
					Height: 0,
					Series: []*charts.ConnectedScatterSeries{
						{
							Points: []charts.ConnectedScatterPoint{
								{X: 0, Y: 10},
								{X: 1, Y: 25},
								{X: 2, Y: 15},
								{X: 3, Y: 30},
								{X: 4, Y: 20},
								{X: 5, Y: 35},
							},
							Color:     "#3b82f6",
							LineStyle: "dashdot",
						},
					},
					ShowLines:   true,
					ShowMarkers: true,
				}
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				spec := data.(charts.ConnectedScatterSpec)
				spec.Width = ctx.ChartWidth - 50
				spec.Height = ctx.ChartHeight - 80
				return charts.RenderConnectedScatter(spec)
			},
		},
		{
			Label: "Line: Long Dash",
			DataProvider: func() interface{} {
				return charts.ConnectedScatterSpec{
					Width:  0,
					Height: 0,
					Series: []*charts.ConnectedScatterSeries{
						{
							Points: []charts.ConnectedScatterPoint{
								{X: 0, Y: 10},
								{X: 1, Y: 25},
								{X: 2, Y: 15},
								{X: 3, Y: 30},
								{X: 4, Y: 20},
								{X: 5, Y: 35},
							},
							Color:     "#3b82f6",
							LineStyle: "longdash",
						},
					},
					ShowLines:   true,
					ShowMarkers: true,
				}
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				spec := data.(charts.ConnectedScatterSpec)
				spec.Width = ctx.ChartWidth - 50
				spec.Height = ctx.ChartHeight - 80
				return charts.RenderConnectedScatter(spec)
			},
		},
	},
	ChartOffsetX: 25.0,
	ChartOffsetY: 25.0,
}

// StatCardGallery defines the stat card gallery configuration
var StatCardGallery = GalleryConfig{
	Name:  "statcard",
	Title: "Stat Card Gallery",
	Layout: &GridLayout{
		Cols:       3,
		Rows:       2,
		BaseWidth:  300,
		BaseHeight: 200,
	},
	Variants: []VariantConfig{
		{
			Label: "Positive Trend",
			DataProvider: func() interface{} {
				return charts.StatCardData{
					Title:     "Total Revenue",
					Value:     "$124.5K",
					Subtitle:  "+12.5% from last month",
					Change:    12,
					ChangePct: 12.5,
					Color:     "#10b981",
					TrendData: makeTrendData([]int{10, 15, 12, 20, 18, 25, 22, 30}),
				}
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				cardData := data.(charts.StatCardData)
				cardW := int(ctx.ChartWidth - 20)
				cardH := int(ctx.ChartHeight - 20)
				tokens := design.DefaultTheme()
				return charts.RenderStatCard(cardData, 0, 0, cardW, cardH, tokens)
			},
		},
		{
			Label: "Negative Trend",
			DataProvider: func() interface{} {
				return charts.StatCardData{
					Title:     "Active Users",
					Value:     "8,234",
					Subtitle:  "-3.2% from last month",
					Change:    -3,
					ChangePct: -3.2,
					Color:     "#ef4444",
					TrendData: makeTrendData([]int{30, 28, 25, 27, 23, 20, 22, 18}),
				}
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				cardData := data.(charts.StatCardData)
				cardW := int(ctx.ChartWidth - 20)
				cardH := int(ctx.ChartHeight - 20)
				tokens := design.DefaultTheme()
				return charts.RenderStatCard(cardData, 0, 0, cardW, cardH, tokens)
			},
		},
		{
			Label: "Steady Growth",
			DataProvider: func() interface{} {
				return charts.StatCardData{
					Title:     "Conversion Rate",
					Value:     "3.45%",
					Subtitle:  "+0.8% from last month",
					Change:    1,
					ChangePct: 0.8,
					Color:     "#3b82f6",
					TrendData: makeTrendData([]int{15, 18, 16, 20, 22, 21, 24, 25}),
				}
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				cardData := data.(charts.StatCardData)
				cardW := int(ctx.ChartWidth - 20)
				cardH := int(ctx.ChartHeight - 20)
				tokens := design.DefaultTheme()
				return charts.RenderStatCard(cardData, 0, 0, cardW, cardH, tokens)
			},
		},
		{
			Label: "Flat Trend",
			DataProvider: func() interface{} {
				return charts.StatCardData{
					Title:     "Page Views",
					Value:     "45.2K",
					Subtitle:  "0.0% from last month",
					Change:    0,
					ChangePct: 0.0,
					Color:     "#6b7280",
					TrendData: makeTrendData([]int{20, 20, 21, 20, 20, 19, 20, 20}),
				}
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				cardData := data.(charts.StatCardData)
				cardW := int(ctx.ChartWidth - 20)
				cardH := int(ctx.ChartHeight - 20)
				tokens := design.DefaultTheme()
				return charts.RenderStatCard(cardData, 0, 0, cardW, cardH, tokens)
			},
		},
		{
			Label: "Alert Trend",
			DataProvider: func() interface{} {
				return charts.StatCardData{
					Title:     "Bounce Rate",
					Value:     "42.1%",
					Subtitle:  "+5.3% from last month",
					Change:    5,
					ChangePct: 5.3,
					Color:     "#f59e0b",
					TrendData: makeTrendData([]int{18, 20, 22, 25, 24, 28, 26, 30}),
				}
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				cardData := data.(charts.StatCardData)
				cardW := int(ctx.ChartWidth - 20)
				cardH := int(ctx.ChartHeight - 20)
				tokens := design.DefaultTheme()
				return charts.RenderStatCard(cardData, 0, 0, cardW, cardH, tokens)
			},
		},
		{
			Label: "Neutral Trend",
			DataProvider: func() interface{} {
				return charts.StatCardData{
					Title:     "Sessions",
					Value:     "12.8K",
					Subtitle:  "-1.2% from last month",
					Change:    -1,
					ChangePct: -1.2,
					Color:     "#8b5cf6",
					TrendData: makeTrendData([]int{25, 24, 26, 25, 23, 24, 22, 23}),
				}
			},
			ChartRenderer: func(data interface{}, ctx RenderContext) string {
				cardData := data.(charts.StatCardData)
				cardW := int(ctx.ChartWidth - 20)
				cardH := int(ctx.ChartHeight - 20)
				tokens := design.DefaultTheme()
				return charts.RenderStatCard(cardData, 0, 0, cardW, cardH, tokens)
			},
		},
	},
	ChartOffsetX: 10.0,
	ChartOffsetY: 10.0,
}

// Helper functions for data generation

func makeTrendData(values []int) []charts.TimeSeriesData {
	result := make([]charts.TimeSeriesData, len(values))
	startDate := mustParseTime("2024-01-01")
	for i, v := range values {
		result[i] = charts.TimeSeriesData{
			Date:  startDate.AddDate(0, 0, i*7), // Weekly data
			Value: v,
		}
	}
	return result
}

func generateViolinValues(mean, stddev float64) []float64 {
	values := make([]float64, 100)
	for i := range values {
		// Simple approximation of normal distribution
		sum := 0.0
		for j := 0; j < 12; j++ {
			sum += float64(i%20) / 20.0
		}
		values[i] = mean + (sum-6)*stddev
	}
	return values
}
