package annotations

import (
	"strings"
	"testing"

	"github.com/SCKelemen/dataviz/scales"
	"github.com/SCKelemen/units"
)

func TestNewTextLabel(t *testing.T) {
	label := NewTextLabel("Test", 5.0, 10.0)

	if label.Text != "Test" {
		t.Errorf("Expected text 'Test', got '%s'", label.Text)
	}
	if label.Mode != PositionData {
		t.Error("Expected PositionData mode")
	}
}

func TestTextLabel_Render(t *testing.T) {
	xScale := scales.NewLinearScale([2]float64{0, 10}, [2]units.Length{units.Px(0), units.Px(100)})
	yScale := scales.NewLinearScale([2]float64{0, 10}, [2]units.Length{units.Px(0), units.Px(100)})

	label := NewTextLabel("Test Label", 5.0, 5.0)
	svg := label.Render(xScale, yScale)

	if svg == "" {
		t.Error("Should produce SVG output")
	}
	if !strings.Contains(svg, "Test Label") {
		t.Error("Should contain label text")
	}
	if !strings.Contains(svg, "text") {
		t.Error("Should contain text element")
	}
}

func TestTextLabel_WithRotation(t *testing.T) {
	xScale := scales.NewLinearScale([2]float64{0, 10}, [2]units.Length{units.Px(0), units.Px(100)})
	yScale := scales.NewLinearScale([2]float64{0, 10}, [2]units.Length{units.Px(0), units.Px(100)})

	label := NewTextLabel("Rotated", 5.0, 5.0).WithRotation(45)
	svg := label.Render(xScale, yScale)

	if !strings.Contains(svg, "rotate") {
		t.Error("Should contain rotate transform")
	}
	if !strings.Contains(svg, "45") {
		t.Error("Should contain rotation angle")
	}
}

func TestTextLabelPixel(t *testing.T) {
	xScale := scales.NewLinearScale([2]float64{0, 10}, [2]units.Length{units.Px(0), units.Px(100)})
	yScale := scales.NewLinearScale([2]float64{0, 10}, [2]units.Length{units.Px(0), units.Px(100)})

	label := NewTextLabelPixel("Pixel Label", 50, 50)
	svg := label.Render(xScale, yScale)

	if !strings.Contains(svg, "Pixel Label") {
		t.Error("Should contain label text")
	}
}

func TestCalloutLabel(t *testing.T) {
	xScale := scales.NewLinearScale([2]float64{0, 10}, [2]units.Length{units.Px(0), units.Px(100)})
	yScale := scales.NewLinearScale([2]float64{0, 10}, [2]units.Length{units.Px(0), units.Px(100)})

	callout := NewCalloutLabel("Important Point", 5.0, 5.0)
	svg := callout.Render(xScale, yScale)

	if !strings.Contains(svg, "Important Point") {
		t.Error("Should contain label text")
	}
	if !strings.Contains(svg, "line") {
		t.Error("Should contain connecting line")
	}
	if !strings.Contains(svg, "circle") {
		t.Error("Should contain marker circle")
	}
}

func TestNewArrow(t *testing.T) {
	arrow := NewArrow(1.0, 2.0, 8.0, 9.0)

	if arrow.Mode != PositionData {
		t.Error("Expected PositionData mode")
	}
	if !arrow.ShowEnd {
		t.Error("Should show end arrow head by default")
	}
	if arrow.ShowStart {
		t.Error("Should not show start arrow head by default")
	}
}

func TestArrow_Render(t *testing.T) {
	xScale := scales.NewLinearScale([2]float64{0, 10}, [2]units.Length{units.Px(0), units.Px(100)})
	yScale := scales.NewLinearScale([2]float64{0, 10}, [2]units.Length{units.Px(0), units.Px(100)})

	arrow := NewArrow(2.0, 2.0, 8.0, 8.0)
	svg := arrow.Render(xScale, yScale)

	if svg == "" {
		t.Error("Should produce SVG output")
	}
	if !strings.Contains(svg, "marker") {
		t.Error("Should contain marker definition")
	}
	if !strings.Contains(svg, "path") {
		t.Error("Should contain path element")
	}
}

func TestArrow_WithDoubleHead(t *testing.T) {
	arrow := NewArrow(1.0, 1.0, 9.0, 9.0).WithDoubleHead(true)

	if !arrow.ShowStart {
		t.Error("Should show start arrow head")
	}
	if !arrow.ShowEnd {
		t.Error("Should show end arrow head")
	}
}

func TestNewHLine(t *testing.T) {
	line := NewHLine(5.0)

	if line.Orientation != OrientationHorizontal {
		t.Error("Should have horizontal orientation")
	}
}

func TestNewVLine(t *testing.T) {
	line := NewVLine(5.0)

	if line.Orientation != OrientationVertical {
		t.Error("Should have vertical orientation")
	}
}

func TestReferenceLine_Render(t *testing.T) {
	xScale := scales.NewLinearScale([2]float64{0, 10}, [2]units.Length{units.Px(0), units.Px(100)})
	yScale := scales.NewLinearScale([2]float64{0, 10}, [2]units.Length{units.Px(0), units.Px(100)})

	hline := NewHLine(5.0).WithLabel("Threshold")
	svg := hline.Render(xScale, yScale)

	if !strings.Contains(svg, "line") {
		t.Error("Should contain line element")
	}
	if !strings.Contains(svg, "Threshold") {
		t.Error("Should contain label")
	}

	vline := NewVLine(5.0)
	svg = vline.Render(xScale, yScale)

	if !strings.Contains(svg, "line") {
		t.Error("Should contain line element")
	}
}

func TestNewReferenceRegion(t *testing.T) {
	region := NewReferenceRegion(2.0, 3.0, 7.0, 8.0)

	if region.Mode != PositionData {
		t.Error("Expected PositionData mode")
	}
}

func TestReferenceRegion_Render(t *testing.T) {
	xScale := scales.NewLinearScale([2]float64{0, 10}, [2]units.Length{units.Px(0), units.Px(100)})
	yScale := scales.NewLinearScale([2]float64{0, 10}, [2]units.Length{units.Px(0), units.Px(100)})

	region := NewReferenceRegion(2.0, 3.0, 7.0, 8.0).WithLabel("Target Range")
	svg := region.Render(xScale, yScale)

	if !strings.Contains(svg, "rect") {
		t.Error("Should contain rect element")
	}
	if !strings.Contains(svg, "Target Range") {
		t.Error("Should contain label")
	}
}

func TestHRegion(t *testing.T) {
	xScale := scales.NewLinearScale([2]float64{0, 10}, [2]units.Length{units.Px(0), units.Px(100)})
	yScale := scales.NewLinearScale([2]float64{0, 10}, [2]units.Length{units.Px(0), units.Px(100)})

	region := NewHRegion(3.0, 7.0)
	svg := region.Render(xScale, yScale)

	if !strings.Contains(svg, "rect") {
		t.Error("Should contain rect element")
	}
}

func TestVRegion(t *testing.T) {
	xScale := scales.NewLinearScale([2]float64{0, 10}, [2]units.Length{units.Px(0), units.Px(100)})
	yScale := scales.NewLinearScale([2]float64{0, 10}, [2]units.Length{units.Px(0), units.Px(100)})

	region := NewVRegion(3.0, 7.0)
	svg := region.Render(xScale, yScale)

	if !strings.Contains(svg, "rect") {
		t.Error("Should contain rect element")
	}
}

func TestGrid_Render(t *testing.T) {
	xScale := scales.NewLinearScale([2]float64{0, 10}, [2]units.Length{units.Px(0), units.Px(100)})
	yScale := scales.NewLinearScale([2]float64{0, 10}, [2]units.Length{units.Px(0), units.Px(100)})

	grid := NewGrid().WithCounts(5, 5)
	svg := grid.Render(xScale, yScale)

	if svg == "" {
		t.Error("Should produce SVG output")
	}
	if !strings.Contains(svg, "line") {
		t.Error("Should contain line elements")
	}
}

func TestAnnotationLayer(t *testing.T) {
	layer := NewAnnotationLayer()

	label1 := NewTextLabel("Label 1", 1.0, 1.0)
	label2 := NewTextLabel("Label 2", 5.0, 5.0)

	layer.Add(label1).Add(label2)

	if len(layer.Annotations) != 2 {
		t.Errorf("Expected 2 annotations, got %d", len(layer.Annotations))
	}
}

func TestAnnotationLayer_Render(t *testing.T) {
	xScale := scales.NewLinearScale([2]float64{0, 10}, [2]units.Length{units.Px(0), units.Px(100)})
	yScale := scales.NewLinearScale([2]float64{0, 10}, [2]units.Length{units.Px(0), units.Px(100)})

	layer := NewAnnotationLayer()
	layer.Add(NewTextLabel("A", 1.0, 1.0))
	layer.Add(NewHLine(5.0))
	layer.Add(NewVLine(5.0))

	svg := layer.Render(xScale, yScale)

	if svg == "" {
		t.Error("Should produce SVG output")
	}
	// Should contain all annotations
	if !strings.Contains(svg, "A") {
		t.Error("Should contain text label")
	}
}

func TestConnector(t *testing.T) {
	xScale := scales.NewLinearScale([2]float64{0, 10}, [2]units.Length{units.Px(0), units.Px(100)})
	yScale := scales.NewLinearScale([2]float64{0, 10}, [2]units.Length{units.Px(0), units.Px(100)})

	connector := NewConnector(1.0, 1.0, 9.0, 9.0)
	svg := connector.Render(xScale, yScale)

	if !strings.Contains(svg, "path") {
		t.Error("Should contain path element")
	}
}

func TestConnector_WithLineStyle(t *testing.T) {
	xScale := scales.NewLinearScale([2]float64{0, 10}, [2]units.Length{units.Px(0), units.Px(100)})
	yScale := scales.NewLinearScale([2]float64{0, 10}, [2]units.Length{units.Px(0), units.Px(100)})

	connector := NewConnector(1.0, 1.0, 9.0, 9.0).WithLineStyle(ConnectorElbow)
	svg := connector.Render(xScale, yScale)

	if !strings.Contains(svg, "path") {
		t.Error("Should contain path element")
	}
	// Elbow connector uses H and V commands
	if !strings.Contains(svg, "H") || !strings.Contains(svg, "V") {
		t.Error("Elbow connector should use H and V path commands")
	}
}
