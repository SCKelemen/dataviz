package layout

import (
	"github.com/SCKelemen/layout"
)

// ChartNode wraps a layout.Node with chart-specific functionality
// This provides the bridge between dataviz concepts and the layout engine
type ChartNode struct {
	*layout.Node

	// Chart-specific metadata
	ChartType string      // Type of chart (bar, line, scatter, etc.)
	Data      interface{} // Chart data
	Renderer  ChartRenderer // Function to render this chart
}

// ChartRenderer renders a chart within its layout bounds
type ChartRenderer func(node *layout.Node) string

// NewChartNode creates a new chart node
func NewChartNode() *ChartNode {
	return &ChartNode{
		Node: &layout.Node{},
	}
}

// WithRenderer sets the chart renderer
func (cn *ChartNode) WithRenderer(renderer ChartRenderer) *ChartNode {
	cn.Renderer = renderer
	return cn
}

// WithData sets the chart data
func (cn *ChartNode) WithData(data interface{}) *ChartNode {
	cn.Data = data
	return cn
}

// WithType sets the chart type
func (cn *ChartNode) WithType(chartType string) *ChartNode {
	cn.ChartType = chartType
	return cn
}

// Render renders this chart node
func (cn *ChartNode) Render() string {
	if cn.Renderer == nil {
		return ""
	}
	return cn.Renderer(cn.Node)
}
