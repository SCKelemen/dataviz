package layout

import (
	"strings"
	"testing"

	"github.com/SCKelemen/layout"
)

func TestChartGrid(t *testing.T) {
	grid := ChartGrid(2, 3)

	if grid.Style.Display != layout.DisplayGrid {
		t.Error("Should have grid display")
	}

	if len(grid.Style.GridTemplateRows) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(grid.Style.GridTemplateRows))
	}

	if len(grid.Style.GridTemplateColumns) != 3 {
		t.Errorf("Expected 3 columns, got %d", len(grid.Style.GridTemplateColumns))
	}

	// Check that tracks are fractional
	if grid.Style.GridTemplateRows[0].Fraction != 1.0 {
		t.Error("Row track should be 1fr")
	}
	if grid.Style.GridTemplateColumns[0].Fraction != 1.0 {
		t.Error("Column track should be 1fr")
	}
}

func TestChartGridWithGap(t *testing.T) {
	grid := ChartGridWithGap(2, 2, 15)

	if grid.Style.GridGap.Value != 15 {
		t.Errorf("Expected gap of 15, got %f", grid.Style.GridGap.Value)
	}
}

func TestChartHStack(t *testing.T) {
	stack := ChartHStack()

	if stack.Style.Display != layout.DisplayFlex {
		t.Error("Should have flex display")
	}

	if stack.Style.FlexDirection != layout.FlexDirectionRow {
		t.Error("Should have row direction")
	}
}

func TestChartVStack(t *testing.T) {
	stack := ChartVStack()

	if stack.Style.Display != layout.DisplayFlex {
		t.Error("Should have flex display")
	}

	if stack.Style.FlexDirection != layout.FlexDirectionColumn {
		t.Error("Should have column direction")
	}
}

func TestChartCell(t *testing.T) {
	cell := ChartCell(1, 2, 2, 1)

	if cell.Style.GridRowStart != 1 {
		t.Errorf("Expected row start 1, got %d", cell.Style.GridRowStart)
	}
	if cell.Style.GridRowEnd != 3 {
		t.Errorf("Expected row end 3, got %d", cell.Style.GridRowEnd)
	}
	if cell.Style.GridColumnStart != 2 {
		t.Errorf("Expected column start 2, got %d", cell.Style.GridColumnStart)
	}
	if cell.Style.GridColumnEnd != 3 {
		t.Errorf("Expected column end 3, got %d", cell.Style.GridColumnEnd)
	}
}

func TestWithMargin(t *testing.T) {
	node := &layout.Node{}
	WithMargin(node, 10)

	if node.Style.Margin.Top.Value != 10 {
		t.Error("Top margin should be 10")
	}
	if node.Style.Margin.Right.Value != 10 {
		t.Error("Right margin should be 10")
	}
	if node.Style.Margin.Bottom.Value != 10 {
		t.Error("Bottom margin should be 10")
	}
	if node.Style.Margin.Left.Value != 10 {
		t.Error("Left margin should be 10")
	}
}

func TestWithPadding(t *testing.T) {
	node := &layout.Node{}
	WithPadding(node, 5)

	if node.Style.Padding.Top.Value != 5 {
		t.Error("Top padding should be 5")
	}
	if node.Style.Padding.Right.Value != 5 {
		t.Error("Right padding should be 5")
	}
	if node.Style.Padding.Bottom.Value != 5 {
		t.Error("Bottom padding should be 5")
	}
	if node.Style.Padding.Left.Value != 5 {
		t.Error("Left padding should be 5")
	}
}

func TestWithSize(t *testing.T) {
	node := &layout.Node{}
	WithSize(node, 800, 600)

	if node.Style.Width.Value != 800 {
		t.Errorf("Expected width 800, got %f", node.Style.Width.Value)
	}
	if node.Style.Height.Value != 600 {
		t.Errorf("Expected height 600, got %f", node.Style.Height.Value)
	}
}

func TestWithFlexGrow(t *testing.T) {
	node := &layout.Node{}
	WithFlexGrow(node, 2)

	if node.Style.FlexGrow != 2 {
		t.Errorf("Expected flex grow 2, got %f", node.Style.FlexGrow)
	}
}

func TestNewDashboard(t *testing.T) {
	dash := NewDashboard(800, 600)

	if dash.Width != 800 {
		t.Errorf("Expected width 800, got %f", dash.Width)
	}
	if dash.Height != 600 {
		t.Errorf("Expected height 600, got %f", dash.Height)
	}
	if dash.Gap != 10 {
		t.Error("Default gap should be 10")
	}
}

func TestDashboard_AddChart(t *testing.T) {
	dash := NewDashboard(800, 600)
	chart := NewChartNode()

	dash.AddChart(chart)

	if len(dash.Charts) != 1 {
		t.Errorf("Expected 1 chart, got %d", len(dash.Charts))
	}
}

func TestDashboard_Layout(t *testing.T) {
	dash := NewDashboard(800, 600)

	// Add some charts with explicit positions
	chart1 := NewChartNode()
	chart1.Style.GridRowStart = 1
	chart1.Style.GridRowEnd = 2
	chart1.Style.GridColumnStart = 1
	chart1.Style.GridColumnEnd = 2

	chart2 := NewChartNode()
	chart2.Style.GridRowStart = 1
	chart2.Style.GridRowEnd = 2
	chart2.Style.GridColumnStart = 2
	chart2.Style.GridColumnEnd = 3

	dash.AddChart(chart1).AddChart(chart2)

	root := dash.Layout()

	if root == nil {
		t.Error("Layout should return a node")
	}

	if root.Style.Display != layout.DisplayGrid {
		t.Error("Root should be a grid")
	}

	if len(root.Children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(root.Children))
	}
}

func TestSideBySideLayout(t *testing.T) {
	node := SideBySideLayout(800, 600)

	if node.Style.Display != layout.DisplayGrid {
		t.Error("Should be a grid")
	}

	if len(node.Style.GridTemplateRows) != 1 {
		t.Error("Should have 1 row")
	}

	if len(node.Style.GridTemplateColumns) != 2 {
		t.Error("Should have 2 columns")
	}
}

func TestTopBottomLayout(t *testing.T) {
	node := TopBottomLayout(800, 600)

	if node.Style.Display != layout.DisplayGrid {
		t.Error("Should be a grid")
	}

	if len(node.Style.GridTemplateRows) != 2 {
		t.Error("Should have 2 rows")
	}

	if len(node.Style.GridTemplateColumns) != 1 {
		t.Error("Should have 1 column")
	}
}

func TestQuadLayout(t *testing.T) {
	node := QuadLayout(800, 600)

	if node.Style.Display != layout.DisplayGrid {
		t.Error("Should be a grid")
	}

	if len(node.Style.GridTemplateRows) != 2 {
		t.Error("Should have 2 rows")
	}

	if len(node.Style.GridTemplateColumns) != 2 {
		t.Error("Should have 2 columns")
	}
}

func TestFacetLayout(t *testing.T) {
	fl := NewFacetLayout(2, 3, 800, 600)

	if fl.Rows != 2 {
		t.Errorf("Expected 2 rows, got %d", fl.Rows)
	}
	if fl.Cols != 3 {
		t.Errorf("Expected 3 cols, got %d", fl.Cols)
	}
	if fl.Width != 800 {
		t.Errorf("Expected width 800, got %f", fl.Width)
	}
	if fl.Height != 600 {
		t.Errorf("Expected height 600, got %f", fl.Height)
	}

	node := fl.Build()

	if node.Style.Display != layout.DisplayGrid {
		t.Error("Should be a grid")
	}
}

func TestTraverseAndRender(t *testing.T) {
	root := ChartGrid(1, 2)
	root.Rect = layout.Rect{X: 0, Y: 0, Width: 800, Height: 600}

	child1 := &layout.Node{}
	child1.Rect = layout.Rect{X: 0, Y: 0, Width: 400, Height: 600}

	child2 := &layout.Node{}
	child2.Rect = layout.Rect{X: 400, Y: 0, Width: 400, Height: 600}

	root = root.AddChild(child1).AddChild(child2)

	svg := TraverseAndRender(root)

	if svg == "" {
		t.Error("Should produce SVG output")
	}

	// Should have groups for children with non-zero positions
	if !strings.Contains(svg, "transform") {
		t.Error("Should contain transforms for positioned children")
	}
}
