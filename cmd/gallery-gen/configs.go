package main

import (
	"github.com/SCKelemen/dataviz/charts"
	design "github.com/SCKelemen/design-system"
)

// GalleryRegistry contains all gallery configurations
var GalleryRegistry = map[string]GalleryConfig{
	"bar": BarGallery,
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
