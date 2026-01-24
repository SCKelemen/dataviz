package main

import (
	"fmt"
	"os"

	"github.com/SCKelemen/dataviz/charts"
)

func main() {
	// Test all new chart types
	testLollipop()
	testDensityPlot()
	testConnectedScatter()
	testStackedArea()
	testCorrelogram()
	testRadar()
	testParallelCoordinates()
	testStreamChart()
	testCircularBarPlot()
	testSankey()
	testChordDiagram()
	testDendrogram()
	testWordCloud()
}

func testLollipop() {
	data := &charts.LollipopData{
		Values: []charts.LollipopPoint{
			{Label: "A", Value: 23},
			{Label: "B", Value: 45},
			{Label: "C", Value: 67},
			{Label: "D", Value: 34},
			{Label: "E", Value: 89},
		},
		Color: "#3b82f6",
	}

	spec := charts.LollipopSpec{
		Data:       data,
		Width:      800,
		Height:     600,
		ShowLabels: true,
		ShowGrid:   true,
		Title:      "Lollipop Chart Test",
		YAxisLabel: "Value",
	}

	content := charts.RenderLollipop(spec)
	writeToFile("test_lollipop.svg", content, 800, 600)
	fmt.Println("✓ Lollipop chart rendered")
}

func testDensityPlot() {
	data1 := &charts.SimpleDensityData{
		Values: []float64{1, 2, 2, 3, 3, 3, 4, 4, 5, 6, 7, 8, 9},
		Label:  "Distribution A",
		Color:  "#3b82f6",
	}

	data2 := &charts.SimpleDensityData{
		Values: []float64{3, 4, 5, 5, 6, 6, 6, 7, 7, 8, 9, 10, 11},
		Label:  "Distribution B",
		Color:  "#10b981",
	}

	spec := charts.SimpleDensitySpec{
		Data:       []*charts.SimpleDensityData{data1, data2},
		Width:      800,
		Height:     600,
		ShowFill:   true,
		ShowRug:    true,
		Title:      "Density Plot Test",
		XAxisLabel: "Value",
		YAxisLabel: "Density",
	}

	content := charts.RenderSimpleDensity(spec)
	writeToFile("test_density.svg", content, 800, 600)
	fmt.Println("✓ Density plot rendered")
}

func testConnectedScatter() {
	series1 := &charts.ConnectedScatterSeries{
		Points: []charts.ConnectedScatterPoint{
			{X: 1, Y: 2},
			{X: 2, Y: 4},
			{X: 3, Y: 3},
			{X: 4, Y: 5},
			{X: 5, Y: 7},
		},
		Label:      "Series A",
		Color:      "#3b82f6",
		MarkerType: "circle",
	}

	spec := charts.ConnectedScatterSpec{
		Series:     []*charts.ConnectedScatterSeries{series1},
		Width:      800,
		Height:     600,
		ShowGrid:   true,
		ShowMarkers: true,
		ShowLines:  true,
		Title:      "Connected Scatter Plot Test",
		XAxisLabel: "X Value",
		YAxisLabel: "Y Value",
	}

	content := charts.RenderConnectedScatter(spec)
	writeToFile("test_connected_scatter.svg", content, 800, 600)
	fmt.Println("✓ Connected scatter plot rendered")
}

func testStackedArea() {
	points := []charts.StackedAreaPoint{
		{X: 1, Values: []float64{10, 20, 15}},
		{X: 2, Values: []float64{15, 25, 20}},
		{X: 3, Values: []float64{20, 30, 25}},
		{X: 4, Values: []float64{18, 28, 22}},
		{X: 5, Values: []float64{25, 35, 30}},
	}

	series := []charts.StackedAreaSeries{
		{Label: "Series A", Color: "#3b82f6"},
		{Label: "Series B", Color: "#10b981"},
		{Label: "Series C", Color: "#f59e0b"},
	}

	spec := charts.StackedAreaSpec{
		Points:     points,
		Series:     series,
		Width:      800,
		Height:     600,
		ShowGrid:   true,
		Title:      "Stacked Area Chart Test",
		XAxisLabel: "Time",
		YAxisLabel: "Value",
	}

	content := charts.RenderStackedArea(spec)
	writeToFile("test_stacked_area.svg", content, 800, 600)
	fmt.Println("✓ Stacked area chart rendered")
}

func testCorrelogram() {
	matrix := charts.CorrelationMatrix{
		Variables: []string{"Var A", "Var B", "Var C", "Var D"},
		Matrix: [][]float64{
			{1.0, 0.8, 0.3, -0.2},
			{0.8, 1.0, 0.5, 0.1},
			{0.3, 0.5, 1.0, 0.7},
			{-0.2, 0.1, 0.7, 1.0},
		},
	}

	spec := charts.CorrelogramSpec{
		Data:         matrix,
		Width:        800,
		Height:       800,
		ShowValues:   true,
		ShowDiagonal: true,
		TriangleMode: "full",
		ColorScheme:  "redblue",
		Title:        "Correlogram Test",
	}

	content := charts.RenderCorrelogram(spec)
	writeToFile("test_correlogram.svg", content, 800, 800)
	fmt.Println("✓ Correlogram rendered")
}

func testRadar() {
	axes := []charts.RadarAxis{
		{Label: "Speed", Min: 0, Max: 100},
		{Label: "Strength", Min: 0, Max: 100},
		{Label: "Defense", Min: 0, Max: 100},
		{Label: "Magic", Min: 0, Max: 100},
		{Label: "Agility", Min: 0, Max: 100},
	}

	series1 := &charts.RadarSeries{
		Label:  "Character A",
		Values: []float64{80, 70, 60, 50, 90},
		Color:  "#3b82f6",
	}

	series2 := &charts.RadarSeries{
		Label:  "Character B",
		Values: []float64{60, 90, 80, 70, 50},
		Color:  "#10b981",
	}

	spec := charts.RadarChartSpec{
		Axes:       axes,
		Series:     []*charts.RadarSeries{series1, series2},
		Width:      800,
		Height:     800,
		ShowGrid:   true,
		ShowLabels: true,
		GridLevels: 5,
		Title:      "Radar Chart Test",
	}

	content := charts.RenderRadarChart(spec)
	writeToFile("test_radar.svg", content, 800, 800)
	fmt.Println("✓ Radar chart rendered")
}

func testParallelCoordinates() {
	axes := []charts.ParallelAxis{
		{Label: "Axis 1", Min: 0, Max: 10},
		{Label: "Axis 2", Min: 0, Max: 100},
		{Label: "Axis 3", Min: 0, Max: 50},
		{Label: "Axis 4", Min: 0, Max: 20},
	}

	data := []charts.ParallelDataPoint{
		{Values: []float64{2, 30, 10, 5}},
		{Values: []float64{5, 60, 25, 12}},
		{Values: []float64{8, 80, 40, 18}},
		{Values: []float64{3, 45, 15, 8}},
	}

	spec := charts.ParallelCoordinatesSpec{
		Axes:           axes,
		Data:           data,
		Width:          800,
		Height:         600,
		ShowAxesLabels: true,
		ShowTicks:      true,
		Title:          "Parallel Coordinates Test",
	}

	content := charts.RenderParallelCoordinates(spec)
	writeToFile("test_parallel.svg", content, 800, 600)
	fmt.Println("✓ Parallel coordinates rendered")
}

func testStreamChart() {
	points := []charts.StreamPoint{
		{X: 1, Values: []float64{10, 15, 12}},
		{X: 2, Values: []float64{15, 18, 14}},
		{X: 3, Values: []float64{20, 22, 18}},
		{X: 4, Values: []float64{18, 20, 16}},
		{X: 5, Values: []float64{25, 28, 22}},
	}

	series := []charts.StreamSeries{
		{Label: "Layer 1", Color: "#3b82f6"},
		{Label: "Layer 2", Color: "#10b981"},
		{Label: "Layer 3", Color: "#f59e0b"},
	}

	spec := charts.StreamChartSpec{
		Points:     points,
		Series:     series,
		Width:      800,
		Height:     600,
		Layout:     "wiggle",
		ShowLegend: true,
		Title:      "Stream Chart Test",
		XAxisLabel: "Time",
	}

	content := charts.RenderStreamChart(spec)
	writeToFile("test_stream.svg", content, 800, 600)
	fmt.Println("✓ Stream chart rendered")
}

func testCircularBarPlot() {
	data := []charts.CircularBarPoint{
		{Label: "Jan", Value: 45},
		{Label: "Feb", Value: 52},
		{Label: "Mar", Value: 68},
		{Label: "Apr", Value: 73},
		{Label: "May", Value: 85},
		{Label: "Jun", Value: 92},
		{Label: "Jul", Value: 88},
		{Label: "Aug", Value: 76},
	}

	spec := charts.CircularBarPlotSpec{
		Data:           data,
		Width:          800,
		Height:         800,
		InnerRadius:    100,
		ShowLabels:     false,
		ShowAxisLabels: true,
		Title:          "Circular Bar Plot Test",
	}

	content := charts.RenderCircularBarPlot(spec)
	writeToFile("test_circular_bar.svg", content, 800, 800)
	fmt.Println("✓ Circular bar plot rendered")
}

func testSankey() {
	nodes := []charts.SankeyNode{
		{ID: "a", Label: "Source A"},
		{ID: "b", Label: "Source B"},
		{ID: "c", Label: "Middle C"},
		{ID: "d", Label: "Middle D"},
		{ID: "e", Label: "Target E"},
		{ID: "f", Label: "Target F"},
	}

	links := []charts.SankeyLink{
		{Source: "a", Target: "c", Value: 40},
		{Source: "a", Target: "d", Value: 20},
		{Source: "b", Target: "c", Value: 30},
		{Source: "b", Target: "d", Value: 40},
		{Source: "c", Target: "e", Value: 50},
		{Source: "c", Target: "f", Value: 20},
		{Source: "d", Target: "e", Value: 30},
		{Source: "d", Target: "f", Value: 30},
	}

	spec := charts.SankeySpec{
		Nodes:      nodes,
		Links:      links,
		Width:      900,
		Height:     600,
		ShowLabels: true,
		Title:      "Sankey Diagram Test",
	}

	content := charts.RenderSankey(spec)
	writeToFile("test_sankey.svg", content, 900, 600)
	fmt.Println("✓ Sankey diagram rendered")
}

func testChordDiagram() {
	entities := []charts.ChordEntity{
		{ID: "a", Label: "Group A"},
		{ID: "b", Label: "Group B"},
		{ID: "c", Label: "Group C"},
		{ID: "d", Label: "Group D"},
	}

	relations := []charts.ChordRelation{
		{Source: "a", Target: "b", Value: 10},
		{Source: "a", Target: "c", Value: 20},
		{Source: "b", Target: "c", Value: 15},
		{Source: "b", Target: "d", Value: 25},
		{Source: "c", Target: "d", Value: 18},
	}

	spec := charts.ChordDiagramSpec{
		Entities:   entities,
		Relations:  relations,
		Width:      800,
		Height:     800,
		ShowLabels: true,
		Title:      "Chord Diagram Test",
	}

	content := charts.RenderChordDiagram(spec)
	writeToFile("test_chord.svg", content, 800, 800)
	fmt.Println("✓ Chord diagram rendered")
}

func testDendrogram() {
	// Create a simple dendrogram tree
	leaf1 := &charts.DendrogramNode{Label: "Item A", Height: 0}
	leaf2 := &charts.DendrogramNode{Label: "Item B", Height: 0}
	leaf3 := &charts.DendrogramNode{Label: "Item C", Height: 0}
	leaf4 := &charts.DendrogramNode{Label: "Item D", Height: 0}

	cluster1 := &charts.DendrogramNode{
		Height:   1.5,
		Children: []*charts.DendrogramNode{leaf1, leaf2},
	}

	cluster2 := &charts.DendrogramNode{
		Height:   2.0,
		Children: []*charts.DendrogramNode{leaf3, leaf4},
	}

	root := &charts.DendrogramNode{
		Height:   3.5,
		Children: []*charts.DendrogramNode{cluster1, cluster2},
	}

	spec := charts.DendrogramSpec{
		Root:        root,
		Width:       800,
		Height:      600,
		Orientation: "vertical",
		ShowLabels:  true,
		ShowHeights: true,
		Title:       "Dendrogram Test",
	}

	content := charts.RenderDendrogram(spec)
	writeToFile("test_dendrogram.svg", content, 800, 600)
	fmt.Println("✓ Dendrogram rendered")
}

func testWordCloud() {
	words := []charts.WordCloudWord{
		{Text: "Data", Frequency: 100},
		{Text: "Visualization", Frequency: 90},
		{Text: "Charts", Frequency: 85},
		{Text: "SVG", Frequency: 70},
		{Text: "Graphics", Frequency: 65},
		{Text: "Analytics", Frequency: 60},
		{Text: "Plots", Frequency: 55},
		{Text: "Statistics", Frequency: 50},
		{Text: "Analysis", Frequency: 45},
		{Text: "Metrics", Frequency: 40},
		{Text: "Dashboard", Frequency: 35},
		{Text: "Insights", Frequency: 30},
	}

	spec := charts.WordCloudSpec{
		Words:  words,
		Width:  800,
		Height: 600,
		Layout: "spiral",
		Title:  "Word Cloud Test",
	}

	content := charts.RenderWordCloud(spec)
	writeToFile("test_wordcloud.svg", content, 800, 600)
	fmt.Println("✓ Word cloud rendered")
}

func writeToFile(filename string, content string, width, height float64) {
	// Wrap content in SVG document
	fullSVG := fmt.Sprintf(`<svg width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f" xmlns="http://www.w3.org/2000/svg">
%s
</svg>`, width, height, width, height, content)

	if err := os.WriteFile(filename, []byte(fullSVG), 0644); err != nil {
		fmt.Printf("Error writing %s: %v\n", filename, err)
	}
}
