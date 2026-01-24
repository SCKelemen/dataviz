package layout

import (
	"strings"
	"testing"

	"github.com/SCKelemen/layout"
)

func TestNewChartNode(t *testing.T) {
	node := NewChartNode()

	if node == nil {
		t.Fatal("NewChartNode should return a non-nil node")
	}

	if node.Node == nil {
		t.Error("ChartNode should have an embedded Node")
	}
}

func TestChartNode_WithRenderer(t *testing.T) {
	node := NewChartNode()

	renderer := func(n *layout.Node) string {
		return "<rect/>"
	}

	node.WithRenderer(renderer)

	if node.Renderer == nil {
		t.Error("Renderer should be set")
	}

	// Test that renderer works
	svg := node.Renderer(node.Node)
	if svg != "<rect/>" {
		t.Errorf("Expected '<rect/>', got '%s'", svg)
	}
}

func TestChartNode_WithData(t *testing.T) {
	node := NewChartNode()

	data := map[string]interface{}{
		"values": []int{1, 2, 3},
		"label":  "Test Data",
	}

	node.WithData(data)

	if node.Data == nil {
		t.Error("Data should be set")
	}

	// Verify data content
	dataMap, ok := node.Data.(map[string]interface{})
	if !ok {
		t.Fatal("Data should be a map")
	}

	if values, ok := dataMap["values"].([]int); !ok || len(values) != 3 {
		t.Error("Data should contain values array")
	}

	if label, ok := dataMap["label"].(string); !ok || label != "Test Data" {
		t.Error("Data should contain label string")
	}
}

func TestChartNode_WithType(t *testing.T) {
	node := NewChartNode()

	node.WithType("bar")

	if node.ChartType != "bar" {
		t.Errorf("Expected chart type 'bar', got '%s'", node.ChartType)
	}
}

func TestChartNode_Render(t *testing.T) {
	node := NewChartNode()

	// Set a simple renderer
	renderer := func(n *layout.Node) string {
		return "<rect width=\"100\" height=\"100\"/>"
	}
	node.WithRenderer(renderer)

	// Set rect dimensions
	node.Rect.Width = 200
	node.Rect.Height = 150

	svg := node.Render()

	if svg == "" {
		t.Error("Render should produce SVG output")
	}

	if !strings.Contains(svg, "<rect") {
		t.Error("Render output should contain rect element")
	}
}

func TestChartNode_RenderWithoutRenderer(t *testing.T) {
	node := NewChartNode()

	// No renderer set
	svg := node.Render()

	// Should return empty string when no renderer
	if svg != "" {
		t.Error("Render should return empty string without renderer")
	}
}

func TestChartNode_Chaining(t *testing.T) {
	node := NewChartNode().
		WithType("line").
		WithData(map[string]interface{}{"x": []int{1, 2, 3}}).
		WithRenderer(func(n *layout.Node) string { return "<path/>" })

	if node.ChartType != "line" {
		t.Error("ChartType should be set via chaining")
	}

	if node.Data == nil {
		t.Error("Data should be set via chaining")
	}

	if node.Renderer == nil {
		t.Error("Renderer should be set via chaining")
	}
}

func TestChartGridCustom(t *testing.T) {
	rows := []layout.GridTrack{
		layout.FixedTrack(layout.Px(100)),
		layout.FixedTrack(layout.Px(200)),
		layout.FixedTrack(layout.Px(100)),
	}
	cols := []layout.GridTrack{
		layout.FixedTrack(layout.Px(150)),
		layout.FixedTrack(layout.Px(250)),
		layout.FixedTrack(layout.Px(100)),
		layout.FixedTrack(layout.Px(50)),
	}

	grid := ChartGridCustom(rows, cols)

	if len(grid.Style.GridTemplateRows) != 3 {
		t.Errorf("Expected 3 rows, got %d", len(grid.Style.GridTemplateRows))
	}

	if len(grid.Style.GridTemplateColumns) != 4 {
		t.Errorf("Expected 4 columns, got %d", len(grid.Style.GridTemplateColumns))
	}
}

func TestWithCustomMargin(t *testing.T) {
	chartNode := NewChartNode()

	WithCustomMargin(chartNode.Node, 10, 20, 30, 40)

	if chartNode.Style.Margin.Top.Value != 10 {
		t.Errorf("Expected top margin 10, got %f", chartNode.Style.Margin.Top.Value)
	}
	if chartNode.Style.Margin.Right.Value != 20 {
		t.Errorf("Expected right margin 20, got %f", chartNode.Style.Margin.Right.Value)
	}
	if chartNode.Style.Margin.Bottom.Value != 30 {
		t.Errorf("Expected bottom margin 30, got %f", chartNode.Style.Margin.Bottom.Value)
	}
	if chartNode.Style.Margin.Left.Value != 40 {
		t.Errorf("Expected left margin 40, got %f", chartNode.Style.Margin.Left.Value)
	}
}

func TestWithCustomPadding(t *testing.T) {
	chartNode := NewChartNode()

	WithCustomPadding(chartNode.Node, 5, 10, 15, 20)

	if chartNode.Style.Padding.Top.Value != 5 {
		t.Errorf("Expected top padding 5, got %f", chartNode.Style.Padding.Top.Value)
	}
	if chartNode.Style.Padding.Right.Value != 10 {
		t.Errorf("Expected right padding 10, got %f", chartNode.Style.Padding.Right.Value)
	}
	if chartNode.Style.Padding.Bottom.Value != 15 {
		t.Errorf("Expected bottom padding 15, got %f", chartNode.Style.Padding.Bottom.Value)
	}
	if chartNode.Style.Padding.Left.Value != 20 {
		t.Errorf("Expected left padding 20, got %f", chartNode.Style.Padding.Left.Value)
	}
}

func TestDashboard_WithGap(t *testing.T) {
	dash := NewDashboard(800, 600)

	dash.WithGap(25)

	if dash.Gap != 25 {
		t.Errorf("Expected gap 25, got %f", dash.Gap)
	}
}

func TestDashboard_Render(t *testing.T) {
	dash := NewDashboard(800, 600)

	// Add a chart with a simple renderer
	chart := NewChartNode().WithRenderer(func(n *layout.Node) string {
		return "<rect/>"
	})

	dash.AddChart(chart)

	svg := dash.Render()

	if svg == "" {
		t.Error("Render should produce SVG output")
	}

	if !strings.Contains(svg, "<svg") {
		t.Error("Render output should contain SVG element")
	}

	if !strings.Contains(svg, "<rect") {
		t.Error("Render output should contain chart content")
	}

	// Check for viewport
	if !strings.Contains(svg, "viewBox") {
		t.Error("SVG should have viewBox attribute")
	}
}

func TestDashboard_RenderEmptyDashboard(t *testing.T) {
	dash := NewDashboard(800, 600)

	// No charts added
	svg := dash.Render()

	if svg == "" {
		t.Error("Render should produce SVG even without charts")
	}

	if !strings.Contains(svg, "<svg") {
		t.Error("Should contain SVG element")
	}
}

func TestRenderChartTree(t *testing.T) {
	// Create a simple chart tree
	root := ChartHStack()

	renderFunc := func(n *layout.Node) string {
		// Simple renderer that checks node position
		if n.Rect.X < 400 {
			return "<rect fill=\"red\"/>"
		}
		return "<rect fill=\"blue\"/>"
	}

	svg := RenderChartTree(root, 800, 600, renderFunc)

	if svg == "" {
		t.Error("RenderChartTree should produce SVG output")
	}
}

