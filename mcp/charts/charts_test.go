package charts

import (
	"strings"
	"testing"

	"github.com/SCKelemen/dataviz/mcp/types"
)

// TestCreateBarChart tests bar chart generation
func TestCreateBarChart(t *testing.T) {
	config := types.BarChartConfig{
		ChartConfig: types.ChartConfig{
			Title:  "Test Bar Chart",
			Width:  800,
			Height: 400,
		},
		Data: []types.DataPoint{
			{Label: "A", Value: 10},
			{Label: "B", Value: 20},
			{Label: "C", Value: 15},
		},
		Color: "#3B82F6",
	}

	svg, err := CreateBarChart(config)
	if err != nil {
		t.Fatalf("CreateBarChart failed: %v", err)
	}

	if svg == "" {
		t.Error("Expected non-empty SVG output")
	}

	if !strings.Contains(svg, "<g") || !strings.Contains(svg, "<rect") {
		t.Error("SVG should contain group and rectangle elements")
	}
}

// TestCreatePieChart tests pie chart generation
func TestCreatePieChart(t *testing.T) {
	config := types.PieChartConfig{
		ChartConfig: types.ChartConfig{
			Title:  "Test Pie Chart",
			Width:  600,
			Height: 600,
		},
		Data: []types.DataPoint{
			{Label: "A", Value: 30},
			{Label: "B", Value: 50},
			{Label: "C", Value: 20},
		},
		Donut: false,
	}

	svg, err := CreatePieChart(config)
	if err != nil {
		t.Fatalf("CreatePieChart failed: %v", err)
	}

	if svg == "" {
		t.Error("Expected non-empty SVG output")
	}

	if !strings.Contains(svg, "<path") {
		t.Error("Pie chart should contain path elements")
	}
}

// TestCreateLineChart tests line chart generation
func TestCreateLineChart(t *testing.T) {
	config := types.LineChartConfig{
		ChartConfig: types.ChartConfig{
			Title:  "Test Line Chart",
			Width:  900,
			Height: 500,
		},
		Series: []types.Series{
			{
				Name: "Series 1",
				Data: []types.Point{
					{X: 0.0, Y: 10},
					{X: 1.0, Y: 20},
					{X: 2.0, Y: 15},
				},
				Color: "#3B82F6",
			},
		},
	}

	svg, err := CreateLineChart(config)
	if err != nil {
		t.Fatalf("CreateLineChart failed: %v", err)
	}

	if svg == "" {
		t.Error("Expected non-empty SVG output")
	}

	if !strings.Contains(svg, "<svg") {
		t.Error("Line chart should contain SVG element")
	}
}

// TestCreateScatterPlot tests scatter plot generation
func TestCreateScatterPlot(t *testing.T) {
	config := types.ScatterPlotConfig{
		ChartConfig: types.ChartConfig{
			Title:  "Test Scatter",
			Width:  800,
			Height: 600,
		},
		Data: []types.XYPoint{
			{X: 1, Y: 2},
			{X: 2, Y: 4},
			{X: 3, Y: 3.5},
		},
	}

	svg, err := CreateScatterPlot(config)
	if err != nil {
		t.Fatalf("CreateScatterPlot failed: %v", err)
	}

	if svg == "" {
		t.Error("Expected non-empty SVG output")
	}

	if !strings.Contains(svg, "<circle") {
		t.Error("Scatter plot should contain circle elements")
	}
}

// TestCreateHeatmap tests heatmap generation
func TestCreateHeatmap(t *testing.T) {
	config := types.HeatmapConfig{
		ChartConfig: types.ChartConfig{
			Title:  "Test Heatmap",
			Width:  800,
			Height: 600,
		},
		Data: types.MatrixData{
			Rows:    []string{"Row 1", "Row 2"},
			Columns: []string{"Col 1", "Col 2"},
			Values: [][]float64{
				{1.0, 2.0},
				{3.0, 4.0},
			},
		},
		ShowValue: true,
	}

	svg, err := CreateHeatmap(config)
	if err != nil {
		t.Fatalf("CreateHeatmap failed: %v", err)
	}

	if svg == "" {
		t.Error("Expected non-empty SVG output")
	}

	if !strings.Contains(svg, "<rect") {
		t.Error("Heatmap should contain rect elements")
	}
}

// TestCreateTreemap tests treemap generation
func TestCreateTreemap(t *testing.T) {
	config := types.TreemapConfig{
		ChartConfig: types.ChartConfig{
			Title:  "Test Treemap",
			Width:  800,
			Height: 600,
		},
		Data: types.TreeNode{
			Name: "root",
			Children: []*types.TreeNode{
				{Name: "A", Value: 100},
				{Name: "B", Value: 200},
			},
		},
		ShowLabels: true,
	}

	svg, err := CreateTreemap(config)
	if err != nil {
		t.Fatalf("CreateTreemap failed: %v", err)
	}

	if svg == "" {
		t.Error("Expected non-empty SVG output")
	}

	// Treemap returns raw SVG elements without wrapper, check for rect
	if !strings.Contains(svg, "<rect") && !strings.Contains(svg, "<g") {
		t.Errorf("Treemap should contain SVG elements, got: %s", svg[:min(len(svg), 100)])
	}
}

// TestCreateSunburst tests sunburst chart generation
func TestCreateSunburst(t *testing.T) {
	config := types.SunburstConfig{
		ChartConfig: types.ChartConfig{
			Title:  "Test Sunburst",
			Width:  600,
			Height: 600,
		},
		Data: types.TreeNode{
			Name: "root",
			Children: []*types.TreeNode{
				{Name: "A", Value: 100},
				{Name: "B", Value: 200},
			},
		},
		ShowLabels: true,
	}

	svg, err := CreateSunburst(config)
	if err != nil {
		t.Fatalf("CreateSunburst failed: %v", err)
	}

	if svg == "" {
		t.Error("Expected non-empty SVG output")
	}

	// Sunburst returns path elements for arcs
	if !strings.Contains(svg, "<path") && !strings.Contains(svg, "<g") {
		t.Errorf("Sunburst should contain SVG elements, got: %s", svg[:min(len(svg), 100)])
	}
}

// TestCreateCirclePacking tests circle packing generation
func TestCreateCirclePacking(t *testing.T) {
	config := types.CirclePackingConfig{
		ChartConfig: types.ChartConfig{
			Title:  "Test Circle Packing",
			Width:  600,
			Height: 600,
		},
		Data: types.TreeNode{
			Name: "root",
			Children: []*types.TreeNode{
				{Name: "A", Value: 100},
				{Name: "B", Value: 200},
			},
		},
		ShowLabels: true,
	}

	svg, err := CreateCirclePacking(config)
	if err != nil {
		t.Fatalf("CreateCirclePacking failed: %v", err)
	}

	if svg == "" {
		t.Error("Expected non-empty SVG output")
	}

	// Circle packing returns circle elements
	if !strings.Contains(svg, "<circle") && !strings.Contains(svg, "<g") {
		t.Errorf("Circle packing should contain SVG elements, got: %s", svg[:min(len(svg), 100)])
	}
}

// TestCreateIcicle tests icicle chart generation
func TestCreateIcicle(t *testing.T) {
	config := types.IcicleConfig{
		ChartConfig: types.ChartConfig{
			Title:  "Test Icicle",
			Width:  800,
			Height: 600,
		},
		Data: types.TreeNode{
			Name: "root",
			Children: []*types.TreeNode{
				{Name: "A", Value: 100},
				{Name: "B", Value: 200},
			},
		},
		Orientation: "vertical",
		ShowLabels:  true,
	}

	svg, err := CreateIcicle(config)
	if err != nil {
		t.Fatalf("CreateIcicle failed: %v", err)
	}

	if svg == "" {
		t.Error("Expected non-empty SVG output")
	}

	// Icicle returns rect elements for partitions
	if !strings.Contains(svg, "<rect") && !strings.Contains(svg, "<g") {
		t.Errorf("Icicle should contain SVG elements, got: %s", svg[:min(len(svg), 100)])
	}
}

// TestCreateBoxPlot tests box plot generation
func TestCreateBoxPlot(t *testing.T) {
	config := types.BoxPlotConfig{
		ChartConfig: types.ChartConfig{
			Title:  "Test Boxplot",
			Width:  800,
			Height: 600,
		},
		Data: []types.BoxPlotDataSet{
			{
				Label:  "Group A",
				Values: []float64{10, 15, 20, 25, 30, 35, 40},
			},
			{
				Label:  "Group B",
				Values: []float64{5, 10, 15, 20, 25, 30, 35},
			},
		},
		ShowOutliers: true,
		ShowMean:     false,
	}

	svg, err := CreateBoxPlot(config)
	if err != nil {
		t.Fatalf("CreateBoxPlot failed: %v", err)
	}

	if svg == "" {
		t.Error("Expected non-empty SVG output")
	}

	if !strings.Contains(svg, "<line") || !strings.Contains(svg, "<rect") {
		t.Error("Boxplot should contain line and rect elements")
	}
}

// TestCreateViolinPlot tests violin plot generation
func TestCreateViolinPlot(t *testing.T) {
	config := types.ViolinPlotConfig{
		ChartConfig: types.ChartConfig{
			Title:  "Test Violin",
			Width:  800,
			Height: 600,
		},
		Data: []types.BoxPlotDataSet{
			{
				Label:  "Group A",
				Values: []float64{10, 12, 15, 18, 20, 22, 25, 28, 30},
			},
		},
		ShowBox:    true,
		ShowMedian: true,
	}

	svg, err := CreateViolinPlot(config)
	if err != nil {
		t.Fatalf("CreateViolinPlot failed: %v", err)
	}

	if svg == "" {
		t.Error("Expected non-empty SVG output")
	}

	if !strings.Contains(svg, "<path") {
		t.Error("Violin plot should contain path elements")
	}
}

// TestCreateHistogram tests histogram generation
func TestCreateHistogram(t *testing.T) {
	config := types.HistogramConfig{
		ChartConfig: types.ChartConfig{
			Title:  "Test Histogram",
			Width:  800,
			Height: 600,
		},
		Values: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		Bins:   5,
	}

	svg, err := CreateHistogram(config)
	if err != nil {
		t.Fatalf("CreateHistogram failed: %v", err)
	}

	if svg == "" {
		t.Error("Expected non-empty SVG output")
	}

	if !strings.Contains(svg, "<rect") {
		t.Error("Histogram should contain rect elements")
	}
}

// TestCreateRidgeline tests ridgeline plot generation
func TestCreateRidgeline(t *testing.T) {
	config := types.RidgelineConfig{
		ChartConfig: types.ChartConfig{
			Title:  "Test Ridgeline",
			Width:  800,
			Height: 600,
		},
		Data: []types.RidgelineDataSet{
			{
				Label:  "Jan",
				Values: []float64{10, 12, 15, 18, 20, 22, 25},
			},
			{
				Label:  "Feb",
				Values: []float64{15, 18, 20, 23, 26, 29, 32},
			},
		},
		Overlap:    0.5,
		ShowLabels: true,
	}

	svg, err := CreateRidgeline(config)
	if err != nil {
		t.Fatalf("CreateRidgeline failed: %v", err)
	}

	if svg == "" {
		t.Error("Expected non-empty SVG output")
	}

	if !strings.Contains(svg, "<path") {
		t.Error("Ridgeline should contain path elements")
	}
}

// TestCreateCandlestick tests candlestick chart generation
func TestCreateCandlestick(t *testing.T) {
	config := types.CandlestickConfig{
		ChartConfig: types.ChartConfig{
			Title:  "Test Candlestick",
			Width:  1000,
			Height: 600,
		},
		Data: []types.CandlestickDataPoint{
			{
				Date:   "2024-01-01",
				Open:   100,
				High:   110,
				Low:    95,
				Close:  105,
				Volume: 1000000,
			},
			{
				Date:   "2024-01-02",
				Open:   105,
				High:   115,
				Low:    103,
				Close:  112,
				Volume: 1200000,
			},
		},
		ShowVolume: true,
	}

	svg, err := CreateCandlestick(config)
	if err != nil {
		t.Fatalf("CreateCandlestick failed: %v", err)
	}

	if svg == "" {
		t.Error("Expected non-empty SVG output")
	}

	if !strings.Contains(svg, "<line") || !strings.Contains(svg, "<rect") {
		t.Error("Candlestick should contain line and rect elements")
	}
}

// TestCreateOHLC tests OHLC chart generation
func TestCreateOHLC(t *testing.T) {
	config := types.OHLCConfig{
		ChartConfig: types.ChartConfig{
			Title:  "Test OHLC",
			Width:  1000,
			Height: 600,
		},
		Data: []types.CandlestickDataPoint{
			{
				Date:  "2024-01-01",
				Open:  100,
				High:  110,
				Low:   95,
				Close: 105,
			},
			{
				Date:  "2024-01-02",
				Open:  105,
				High:  115,
				Low:   103,
				Close: 112,
			},
		},
	}

	svg, err := CreateOHLC(config)
	if err != nil {
		t.Fatalf("CreateOHLC failed: %v", err)
	}

	if svg == "" {
		t.Error("Expected non-empty SVG output")
	}

	if !strings.Contains(svg, "<line") {
		t.Error("OHLC should contain line elements")
	}
}

// TestConvertTreeNode tests tree node conversion
func TestConvertTreeNode(t *testing.T) {
	input := &types.TreeNode{
		Name:  "root",
		Value: 100,
		Children: []*types.TreeNode{
			{Name: "child1", Value: 40},
			{Name: "child2", Value: 60},
		},
	}

	result := convertTreeNode(input)

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if result.Name != "root" {
		t.Errorf("Expected name 'root', got '%s'", result.Name)
	}

	if result.Value != 100 {
		t.Errorf("Expected value 100, got %f", result.Value)
	}

	if len(result.Children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(result.Children))
	}

	if result.Children[0].Name != "child1" {
		t.Errorf("Expected child name 'child1', got '%s'", result.Children[0].Name)
	}
}

// TestConvertTreeNodeNil tests nil tree node conversion
func TestConvertTreeNodeNil(t *testing.T) {
	result := convertTreeNode(nil)

	if result != nil {
		t.Error("Expected nil result for nil input")
	}
}

// TestConvertTreeNodeDeep tests deep tree conversion
func TestConvertTreeNodeDeep(t *testing.T) {
	input := &types.TreeNode{
		Name: "root",
		Children: []*types.TreeNode{
			{
				Name: "level1",
				Children: []*types.TreeNode{
					{
						Name: "level2",
						Children: []*types.TreeNode{
							{Name: "level3", Value: 10},
						},
					},
				},
			},
		},
	}

	result := convertTreeNode(input)

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if len(result.Children) != 1 ||
		len(result.Children[0].Children) != 1 ||
		len(result.Children[0].Children[0].Children) != 1 {
		t.Error("Tree structure not preserved correctly")
	}

	leaf := result.Children[0].Children[0].Children[0]
	if leaf.Name != "level3" || leaf.Value != 10 {
		t.Error("Leaf node not converted correctly")
	}
}

// TestEmptyData tests charts with empty data
func TestEmptyData(t *testing.T) {
	t.Run("empty candlestick data", func(t *testing.T) {
		config := types.CandlestickConfig{
			ChartConfig: types.ChartConfig{Width: 800, Height: 600},
			Data:        []types.CandlestickDataPoint{},
		}
		_, err := CreateCandlestick(config)
		if err == nil {
			t.Error("Expected error for empty candlestick data")
		}
	})

	t.Run("empty OHLC data", func(t *testing.T) {
		config := types.OHLCConfig{
			ChartConfig: types.ChartConfig{Width: 800, Height: 600},
			Data:        []types.CandlestickDataPoint{},
		}
		_, err := CreateOHLC(config)
		if err == nil {
			t.Error("Expected error for empty OHLC data")
		}
	})

	t.Run("empty line chart series", func(t *testing.T) {
		config := types.LineChartConfig{
			ChartConfig: types.ChartConfig{Width: 800, Height: 600},
			Series:      []types.Series{},
		}
		_, err := CreateLineChart(config)
		if err == nil {
			t.Error("Expected error for empty line chart series")
		}
	})

	t.Run("empty scatter data", func(t *testing.T) {
		config := types.ScatterPlotConfig{
			ChartConfig: types.ChartConfig{Width: 800, Height: 600},
			Data:        []types.XYPoint{},
		}
		_, err := CreateScatterPlot(config)
		if err == nil {
			t.Error("Expected error for empty scatter data")
		}
	})
}

// TestNonZeroDimensions tests charts with explicit dimensions
func TestNonZeroDimensions(t *testing.T) {
	config := types.TreemapConfig{
		ChartConfig: types.ChartConfig{
			Width:  800,
			Height: 600,
		},
		Data: types.TreeNode{
			Name:  "root",
			Value: 100,
		},
	}

	svg, err := CreateTreemap(config)
	if err != nil {
		t.Fatalf("CreateTreemap failed: %v", err)
	}

	// Should generate valid SVG with proper dimensions
	if svg == "" {
		t.Error("Expected non-empty SVG output")
	}
}

// TestMinHelper tests the min helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
