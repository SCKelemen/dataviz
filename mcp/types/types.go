package types

// DataPoint represents a simple labeled value
type DataPoint struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
}

// Point represents an X,Y coordinate
type Point struct {
	X interface{} `json:"x"` // Can be number, string, or time
	Y float64     `json:"y"`
}

// Series represents a named data series
type Series struct {
	Name  string  `json:"name"`
	Data  []Point `json:"data"`
	Color string  `json:"color,omitempty"`
}

// XYPoint represents a point in 2D space
type XYPoint struct {
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Label string  `json:"label,omitempty"`
	Size  float64 `json:"size,omitempty"` // For bubble charts
}

// MatrixData represents matrix/heatmap data
type MatrixData struct {
	Rows    []string    `json:"rows"`
	Columns []string    `json:"columns"`
	Values  [][]float64 `json:"values"`
}

// ChartConfig represents common chart configuration
type ChartConfig struct {
	Title  string `json:"title,omitempty"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
	Theme  string `json:"theme,omitempty"`
}

// BarChartConfig configuration for bar charts
type BarChartConfig struct {
	ChartConfig
	Data        []DataPoint `json:"data"`
	Orientation string      `json:"orientation,omitempty"` // "vertical" or "horizontal"
	Grouped     bool        `json:"grouped,omitempty"`
	Stacked     bool        `json:"stacked,omitempty"`
	Color       string      `json:"color,omitempty"`
}

// PieChartConfig configuration for pie charts
type PieChartConfig struct {
	ChartConfig
	Data  []DataPoint `json:"data"`
	Donut bool        `json:"donut,omitempty"`
}

// LineChartConfig configuration for line charts
type LineChartConfig struct {
	ChartConfig
	Series []Series `json:"series"`
	XLabel string   `json:"x_label,omitempty"`
	YLabel string   `json:"y_label,omitempty"`
	Area   bool     `json:"area,omitempty"` // Fill area under line
}

// ScatterPlotConfig configuration for scatter plots
type ScatterPlotConfig struct {
	ChartConfig
	Data   []XYPoint `json:"data"`
	XLabel string    `json:"x_label,omitempty"`
	YLabel string    `json:"y_label,omitempty"`
}

// HeatmapConfig configuration for heatmaps
type HeatmapConfig struct {
	ChartConfig
	Data      MatrixData `json:"data"`
	ColorMap  string     `json:"color_map,omitempty"` // "viridis", "plasma", etc.
	ShowValue bool       `json:"show_value,omitempty"`
}

// ChartOutput represents the output of a chart generation
type ChartOutput struct {
	SVG      string            `json:"svg"`
	Metadata map[string]string `json:"metadata,omitempty"`
}
