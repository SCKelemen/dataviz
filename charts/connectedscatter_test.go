package charts

import (
	"strings"
	"testing"
)

func TestRenderConnectedScatter(t *testing.T) {
	spec := ConnectedScatterSpec{
		Width:  400,
		Height: 300,
		Series: []*ConnectedScatterSeries{
			{
				Points: []ConnectedScatterPoint{
					{X: 0, Y: 0},
					{X: 1, Y: 2},
					{X: 2, Y: 1},
				},
				Label:      "Test Series",
				Color:      "#3b82f6",
				LineStyle:  "solid",
				MarkerType: "circle",
			},
		},
		ShowLines:   true,
		ShowMarkers: true,
	}

	result := RenderConnectedScatter(spec)

	if result == "" {
		t.Error("Expected non-empty SVG output")
	}

	// Check for path elements (connecting lines)
	if !strings.Contains(result, "<path") {
		t.Error("Expected <path> elements for connecting lines")
	}

	// Check for circle markers
	if !strings.Contains(result, "<circle") {
		t.Error("Expected <circle> elements for markers")
	}
}

func TestRenderConnectedScatter_DashedLine(t *testing.T) {
	spec := ConnectedScatterSpec{
		Width:  400,
		Height: 300,
		Series: []*ConnectedScatterSeries{
			{
				Points: []ConnectedScatterPoint{
					{X: 0, Y: 0},
					{X: 1, Y: 1},
					{X: 2, Y: 2},
				},
				LineStyle: "dashed",
			},
		},
		ShowLines: true,
	}

	result := RenderConnectedScatter(spec)

	// Check that stroke-dasharray is present for dashed lines
	if !strings.Contains(result, "stroke-dasharray") {
		t.Error("Expected stroke-dasharray attribute for dashed line style")
	}

	if !strings.Contains(result, "10,5") {
		t.Error("Expected dashed pattern '10,5' for dashed line style")
	}
}

func TestRenderConnectedScatter_DottedLine(t *testing.T) {
	spec := ConnectedScatterSpec{
		Width:  400,
		Height: 300,
		Series: []*ConnectedScatterSeries{
			{
				Points: []ConnectedScatterPoint{
					{X: 0, Y: 0},
					{X: 1, Y: 1},
				},
				LineStyle: "dotted",
			},
		},
		ShowLines: true,
	}

	result := RenderConnectedScatter(spec)

	if !strings.Contains(result, "stroke-dasharray") {
		t.Error("Expected stroke-dasharray attribute for dotted line style")
	}

	if !strings.Contains(result, "2,3") {
		t.Error("Expected dotted pattern '2,3' for dotted line style")
	}
}

func TestRenderConnectedScatter_DashDotLine(t *testing.T) {
	spec := ConnectedScatterSpec{
		Width:  400,
		Height: 300,
		Series: []*ConnectedScatterSeries{
			{
				Points: []ConnectedScatterPoint{
					{X: 0, Y: 0},
					{X: 1, Y: 1},
				},
				LineStyle: "dashdot",
			},
		},
		ShowLines: true,
	}

	result := RenderConnectedScatter(spec)

	if !strings.Contains(result, "stroke-dasharray") {
		t.Error("Expected stroke-dasharray attribute for dashdot line style")
	}

	if !strings.Contains(result, "10,5,2,5") {
		t.Error("Expected dashdot pattern '10,5,2,5' for dashdot line style")
	}
}

func TestRenderConnectedScatter_SolidLine(t *testing.T) {
	spec := ConnectedScatterSpec{
		Width:  400,
		Height: 300,
		Series: []*ConnectedScatterSeries{
			{
				Points: []ConnectedScatterPoint{
					{X: 0, Y: 0},
					{X: 1, Y: 1},
				},
				LineStyle: "solid",
			},
		},
		ShowLines: true,
	}

	result := RenderConnectedScatter(spec)

	// Solid lines should NOT have stroke-dasharray
	if strings.Contains(result, "stroke-dasharray") {
		t.Error("Did not expect stroke-dasharray attribute for solid line style")
	}
}

func TestRenderConnectedScatter_EmptySeries(t *testing.T) {
	spec := ConnectedScatterSpec{
		Width:  400,
		Height: 300,
		Series: []*ConnectedScatterSeries{},
	}

	result := RenderConnectedScatter(spec)

	if result != "" {
		t.Error("Expected empty output for empty series")
	}
}

func TestRenderConnectedScatter_MultipleSeries(t *testing.T) {
	spec := ConnectedScatterSpec{
		Width:  400,
		Height: 300,
		Series: []*ConnectedScatterSeries{
			{
				Points: []ConnectedScatterPoint{
					{X: 0, Y: 0},
					{X: 1, Y: 1},
				},
				Label:     "Series 1",
				LineStyle: "solid",
			},
			{
				Points: []ConnectedScatterPoint{
					{X: 0, Y: 1},
					{X: 1, Y: 0},
				},
				Label:     "Series 2",
				LineStyle: "dashed",
			},
		},
		ShowLines:   true,
		ShowMarkers: true,
	}

	result := RenderConnectedScatter(spec)

	// Should contain both series labels in legend
	if !strings.Contains(result, "Series 1") {
		t.Error("Expected 'Series 1' in legend")
	}
	if !strings.Contains(result, "Series 2") {
		t.Error("Expected 'Series 2' in legend")
	}

	// Should have both solid and dashed lines
	pathCount := strings.Count(result, "<path")
	if pathCount < 2 {
		t.Errorf("Expected at least 2 paths for multiple series, got %d", pathCount)
	}
}
