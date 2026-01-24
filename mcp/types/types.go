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

// TreeNode represents hierarchical tree data
type TreeNode struct {
	Name     string      `json:"name"`
	Value    float64     `json:"value,omitempty"`
	Children []*TreeNode `json:"children,omitempty"`
}

// TreemapConfig configuration for treemap charts
type TreemapConfig struct {
	ChartConfig
	Data       TreeNode `json:"data"`
	ShowLabels bool     `json:"show_labels,omitempty"`
}

// SunburstConfig configuration for sunburst charts
type SunburstConfig struct {
	ChartConfig
	Data       TreeNode `json:"data"`
	ShowLabels bool     `json:"show_labels,omitempty"`
}

// CirclePackingConfig configuration for circle packing charts
type CirclePackingConfig struct {
	ChartConfig
	Data       TreeNode `json:"data"`
	ShowLabels bool     `json:"show_labels,omitempty"`
}

// IcicleConfig configuration for icicle charts
type IcicleConfig struct {
	ChartConfig
	Data        TreeNode `json:"data"`
	Orientation string   `json:"orientation,omitempty"` // "vertical" or "horizontal"
	ShowLabels  bool     `json:"show_labels,omitempty"`
}

// BoxPlotDataSet represents data for one box in a box plot
type BoxPlotDataSet struct {
	Label  string    `json:"label"`
	Values []float64 `json:"values"`
}

// BoxPlotConfig configuration for box plots
type BoxPlotConfig struct {
	ChartConfig
	Data         []BoxPlotDataSet `json:"data"`
	ShowOutliers bool             `json:"show_outliers,omitempty"`
	ShowMean     bool             `json:"show_mean,omitempty"`
}

// ViolinPlotConfig configuration for violin plots
type ViolinPlotConfig struct {
	ChartConfig
	Data       []BoxPlotDataSet `json:"data"` // Reuse BoxPlotDataSet structure
	ShowBox    bool             `json:"show_box,omitempty"`
	ShowMedian bool             `json:"show_median,omitempty"`
}

// HistogramConfig configuration for histogram charts
type HistogramConfig struct {
	ChartConfig
	Values []float64 `json:"values"`
	Bins   int       `json:"bins,omitempty"`
}

// RidgelineDataSet represents one ridge in a ridgeline plot
type RidgelineDataSet struct {
	Label  string    `json:"label"`
	Values []float64 `json:"values"`
}

// RidgelineConfig configuration for ridgeline plots
type RidgelineConfig struct {
	ChartConfig
	Data       []RidgelineDataSet `json:"data"`
	Overlap    float64            `json:"overlap,omitempty"`
	ShowLabels bool               `json:"show_labels,omitempty"`
}

// CandlestickDataPoint represents one candlestick
type CandlestickDataPoint struct {
	Date   string  `json:"date"`
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume float64 `json:"volume,omitempty"`
}

// CandlestickConfig configuration for candlestick charts
type CandlestickConfig struct {
	ChartConfig
	Data       []CandlestickDataPoint `json:"data"`
	ShowVolume bool                   `json:"show_volume,omitempty"`
}

// OHLCConfig configuration for OHLC bar charts
type OHLCConfig struct {
	ChartConfig
	Data []CandlestickDataPoint `json:"data"` // Reuse same data structure
}
