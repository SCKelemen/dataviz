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

// New chart type configurations

// LollipopPoint represents a single lollipop
type LollipopPoint struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
	Color string  `json:"color,omitempty"`
}

// LollipopConfig configuration for lollipop charts
type LollipopConfig struct {
	ChartConfig
	Values     []LollipopPoint `json:"values"`
	Color      string          `json:"color,omitempty"`
	Horizontal bool            `json:"horizontal,omitempty"`
}

// DensityDataSet represents a density distribution
type DensityDataSet struct {
	Values []float64 `json:"values"`
	Label  string    `json:"label,omitempty"`
	Color  string    `json:"color,omitempty"`
}

// DensityConfig configuration for density plots
type DensityConfig struct {
	ChartConfig
	Data     []DensityDataSet `json:"data"`
	ShowFill bool             `json:"show_fill,omitempty"`
	ShowRug  bool             `json:"show_rug,omitempty"`
}

// ConnectedScatterPoint represents a point in connected scatter
type ConnectedScatterPoint struct {
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Label string  `json:"label,omitempty"`
}

// ConnectedScatterSeries represents a series
type ConnectedScatterSeries struct {
	Points     []ConnectedScatterPoint `json:"points"`
	Label      string                  `json:"label,omitempty"`
	Color      string                  `json:"color,omitempty"`
	MarkerType string                  `json:"marker_type,omitempty"`
}

// ConnectedScatterConfig configuration for connected scatter plots
type ConnectedScatterConfig struct {
	ChartConfig
	Series []ConnectedScatterSeries `json:"series"`
}

// StackedAreaPoint represents a point with multiple values
type StackedAreaPoint struct {
	X      float64   `json:"x"`
	Values []float64 `json:"values"`
}

// StackedAreaSeries represents series metadata
type StackedAreaSeries struct {
	Label string `json:"label"`
	Color string `json:"color,omitempty"`
}

// StackedAreaConfig configuration for stacked area charts
type StackedAreaConfig struct {
	ChartConfig
	Points []StackedAreaPoint  `json:"points"`
	Series []StackedAreaSeries `json:"series"`
}

// StreamChartConfig configuration for stream charts
type StreamChartConfig struct {
	ChartConfig
	Points []StackedAreaPoint  `json:"points"` // Reuse stacked area point
	Series []StackedAreaSeries `json:"series"` // Reuse stacked area series
	Layout string              `json:"layout,omitempty"`
}

// CorrelogramConfig configuration for correlograms
type CorrelogramConfig struct {
	ChartConfig
	Variables []string    `json:"variables"`
	Matrix    [][]float64 `json:"matrix"`
}

// RadarAxis represents an axis in radar chart
type RadarAxis struct {
	Label string  `json:"label"`
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
}

// RadarSeries represents a data series
type RadarSeries struct {
	Label  string    `json:"label"`
	Values []float64 `json:"values"`
	Color  string    `json:"color,omitempty"`
}

// RadarConfig configuration for radar charts
type RadarConfig struct {
	ChartConfig
	Axes   []RadarAxis   `json:"axes"`
	Series []RadarSeries `json:"series"`
}

// ParallelAxis represents an axis in parallel coordinates
type ParallelAxis struct {
	Label string  `json:"label"`
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
}

// ParallelDataPoint represents a data point
type ParallelDataPoint struct {
	Values []float64 `json:"values"`
	Color  string    `json:"color,omitempty"`
}

// ParallelConfig configuration for parallel coordinates
type ParallelConfig struct {
	ChartConfig
	Axes []ParallelAxis      `json:"axes"`
	Data []ParallelDataPoint `json:"data"`
}

// WordCloudWord represents a word in word cloud
type WordCloudWord struct {
	Text      string  `json:"text"`
	Frequency float64 `json:"frequency"`
	Color     string  `json:"color,omitempty"`
}

// WordCloudConfig configuration for word clouds
type WordCloudConfig struct {
	ChartConfig
	Words  []WordCloudWord `json:"words"`
	Layout string          `json:"layout,omitempty"`
}

// SankeyNode represents a node in Sankey diagram
type SankeyNode struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Color string `json:"color,omitempty"`
}

// SankeyLink represents a link in Sankey diagram
type SankeyLink struct {
	Source string  `json:"source"`
	Target string  `json:"target"`
	Value  float64 `json:"value"`
	Color  string  `json:"color,omitempty"`
}

// SankeyConfig configuration for Sankey diagrams
type SankeyConfig struct {
	ChartConfig
	Nodes []SankeyNode `json:"nodes"`
	Links []SankeyLink `json:"links"`
}

// ChordEntity represents an entity in chord diagram
type ChordEntity struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Color string `json:"color,omitempty"`
}

// ChordRelation represents a relation in chord diagram
type ChordRelation struct {
	Source string  `json:"source"`
	Target string  `json:"target"`
	Value  float64 `json:"value"`
}

// ChordConfig configuration for chord diagrams
type ChordConfig struct {
	ChartConfig
	Entities  []ChordEntity   `json:"entities"`
	Relations []ChordRelation `json:"relations"`
}

// CircularBarPoint represents a point in circular bar plot
type CircularBarPoint struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
	Color string  `json:"color,omitempty"`
}

// CircularBarConfig configuration for circular bar plots
type CircularBarConfig struct {
	ChartConfig
	Data        []CircularBarPoint `json:"data"`
	InnerRadius float64            `json:"inner_radius,omitempty"`
}

// DendrogramNode represents a node in dendrogram (recursive)
type DendrogramNode struct {
	Label    string            `json:"label,omitempty"`
	Height   float64           `json:"height"`
	Children []*DendrogramNode `json:"children,omitempty"`
}

// DendrogramConfig configuration for dendrograms
type DendrogramConfig struct {
	ChartConfig
	Root        *DendrogramNode `json:"root"`
	Orientation string          `json:"orientation,omitempty"`
}
