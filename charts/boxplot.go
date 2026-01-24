package charts

import (
	"math"
	"sort"

	"github.com/SCKelemen/dataviz/scales"
	"github.com/SCKelemen/svg"
	"github.com/SCKelemen/units"
)

// BoxPlotData represents data for a box plot
type BoxPlotData struct {
	Values      []float64 // Raw data values
	Label       string    // Category label
	Color       string    // Box fill color
	StrokeColor string    // Box outline color

	// Optional: Pre-calculated statistics (if nil, will be calculated)
	Q1         *float64 // First quartile (25th percentile)
	Median     *float64 // Median (50th percentile)
	Q3         *float64 // Third quartile (75th percentile)
	Min        *float64 // Minimum value (or lower whisker)
	Max        *float64 // Maximum value (or upper whisker)
	Outliers   []float64 // Optional: outlier values
}

// BoxPlotSpec configures box plot rendering
type BoxPlotSpec struct {
	Data        []*BoxPlotData
	Width       float64
	Height      float64
	Horizontal  bool    // If true, render horizontal box plots
	BoxWidth    float64 // Width of each box (0 = auto)
	ShowOutliers bool   // If true, show outlier points
	ShowMean    bool    // If true, show mean marker
	ShowNotch   bool    // If true, show confidence interval notch
	WhiskerMultiplier float64 // IQR multiplier for whiskers (default: 1.5)

	// Axis configuration
	XAxisLabel string
	YAxisLabel string
	ShowGrid   bool
}

// BoxPlotStats represents calculated box plot statistics
type BoxPlotStats struct {
	Q1         float64
	Median     float64
	Q3         float64
	Min        float64
	Max        float64
	Mean       float64
	IQR        float64
	LowerFence float64
	UpperFence float64
	Outliers   []float64
}

// CalculateBoxPlotStats calculates box plot statistics from raw values
func CalculateBoxPlotStats(values []float64, whiskerMultiplier float64) BoxPlotStats {
	if len(values) == 0 {
		return BoxPlotStats{}
	}

	// Sort values
	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)

	// Calculate quartiles
	q1 := percentile(sorted, 25)
	median := percentile(sorted, 50)
	q3 := percentile(sorted, 75)
	iqr := q3 - q1

	// Calculate fences for outlier detection
	if whiskerMultiplier == 0 {
		whiskerMultiplier = 1.5
	}
	lowerFence := q1 - whiskerMultiplier*iqr
	upperFence := q3 + whiskerMultiplier*iqr

	// Find whisker positions (min/max within fences)
	minVal := sorted[0]
	maxVal := sorted[len(sorted)-1]
	whiskerMin := minVal
	whiskerMax := maxVal

	var outliers []float64
	for _, v := range sorted {
		if v < lowerFence {
			outliers = append(outliers, v)
		} else if whiskerMin == minVal || v < whiskerMin {
			whiskerMin = v
			break
		}
	}

	for i := len(sorted) - 1; i >= 0; i-- {
		v := sorted[i]
		if v > upperFence {
			outliers = append(outliers, v)
		} else if whiskerMax == maxVal || v > whiskerMax {
			whiskerMax = v
			break
		}
	}

	// Calculate mean
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	mean := sum / float64(len(values))

	return BoxPlotStats{
		Q1:         q1,
		Median:     median,
		Q3:         q3,
		Min:        whiskerMin,
		Max:        whiskerMax,
		Mean:       mean,
		IQR:        iqr,
		LowerFence: lowerFence,
		UpperFence: upperFence,
		Outliers:   outliers,
	}
}

// percentile calculates the nth percentile of sorted values
func percentile(sorted []float64, p float64) float64 {
	if len(sorted) == 0 {
		return 0
	}
	if len(sorted) == 1 {
		return sorted[0]
	}

	rank := (p / 100.0) * float64(len(sorted)-1)
	lower := int(math.Floor(rank))
	upper := int(math.Ceil(rank))
	weight := rank - float64(lower)

	return sorted[lower]*(1-weight) + sorted[upper]*weight
}

// RenderVerticalBoxPlot renders a vertical box plot
func RenderVerticalBoxPlot(spec BoxPlotSpec) string {
	if len(spec.Data) == 0 {
		return ""
	}

	// Calculate statistics for each box
	stats := make([]BoxPlotStats, len(spec.Data))
	for i, data := range spec.Data {
		if data.Q1 != nil && data.Median != nil && data.Q3 != nil {
			// Use provided statistics
			stats[i] = BoxPlotStats{
				Q1:       *data.Q1,
				Median:   *data.Median,
				Q3:       *data.Q3,
				Min:      valueOrDefault(data.Min, *data.Q1),
				Max:      valueOrDefault(data.Max, *data.Q3),
				Outliers: data.Outliers,
			}
		} else {
			// Calculate statistics
			stats[i] = CalculateBoxPlotStats(data.Values, spec.WhiskerMultiplier)
		}
	}

	// Find global min/max for scale
	globalMin := math.Inf(1)
	globalMax := math.Inf(-1)
	for i := range stats {
		if stats[i].Min < globalMin {
			globalMin = stats[i].Min
		}
		if stats[i].Max > globalMax {
			globalMax = stats[i].Max
		}
		for _, outlier := range stats[i].Outliers {
			if outlier < globalMin {
				globalMin = outlier
			}
			if outlier > globalMax {
				globalMax = outlier
			}
		}
	}

	// Create scales
	yScale := scales.NewLinearScale(
		[2]float64{globalMin, globalMax},
		[2]units.Length{units.Px(spec.Height - 40), units.Px(40)},
	)

	// Calculate box width
	boxWidth := spec.BoxWidth
	if boxWidth == 0 {
		boxWidth = (spec.Width - 80) / float64(len(spec.Data)) * 0.6
	}

	// Calculate x positions
	spacing := (spec.Width - 80) / float64(len(spec.Data))

	var result string

	// Render each box
	for i, st := range stats {
		data := spec.Data[i]
		x := 40 + spacing*float64(i) + spacing/2

		// Scale values
		q1Y := yScale.Apply(st.Q1).Value
		medianY := yScale.Apply(st.Median).Value
		q3Y := yScale.Apply(st.Q3).Value
		minY := yScale.Apply(st.Min).Value
		maxY := yScale.Apply(st.Max).Value

		boxHeight := q1Y - q3Y

		// Colors
		fillColor := data.Color
		if fillColor == "" {
			fillColor = "#4285f4"
		}
		strokeColor := data.StrokeColor
		if strokeColor == "" {
			strokeColor = "#333"
		}

		// Draw whisker lines
		whiskerStyle := svg.Style{
			Stroke:      strokeColor,
			StrokeWidth: 1,
		}
		result += svg.Line(x, minY, x, q3Y, whiskerStyle) + "\n"
		result += svg.Line(x, maxY, x, q1Y, whiskerStyle) + "\n"

		// Draw whisker caps
		capWidth := boxWidth * 0.3
		result += svg.Line(x-capWidth/2, minY, x+capWidth/2, minY, whiskerStyle) + "\n"
		result += svg.Line(x-capWidth/2, maxY, x+capWidth/2, maxY, whiskerStyle) + "\n"

		// Draw box
		boxStyle := svg.Style{
			Fill:        fillColor,
			Stroke:      strokeColor,
			StrokeWidth: 1.5,
			Opacity:     0.7,
		}
		result += svg.Rect(x-boxWidth/2, q3Y, boxWidth, boxHeight, boxStyle) + "\n"

		// Draw median line
		medianStyle := svg.Style{
			Stroke:      strokeColor,
			StrokeWidth: 2,
		}
		result += svg.Line(x-boxWidth/2, medianY, x+boxWidth/2, medianY, medianStyle) + "\n"

		// Draw mean marker if enabled
		if spec.ShowMean {
			meanY := yScale.Apply(st.Mean).Value
			meanStyle := svg.Style{
				Fill:   strokeColor,
				Stroke: strokeColor,
			}
			result += svg.Circle(x, meanY, 3, meanStyle) + "\n"
		}

		// Draw outliers if enabled
		if spec.ShowOutliers && len(st.Outliers) > 0 {
			outlierStyle := svg.Style{
				Fill:        "none",
				Stroke:      strokeColor,
				StrokeWidth: 1,
			}
			for _, outlier := range st.Outliers {
				outlierY := yScale.Apply(outlier).Value
				result += svg.Circle(x, outlierY, 3, outlierStyle) + "\n"
			}
		}

		// Draw label
		if data.Label != "" {
			labelStyle := svg.Style{
				FontSize:   units.Px(12),
				FontFamily: "sans-serif",
				TextAnchor: svg.TextAnchorMiddle,
			}
			result += svg.Text(data.Label, x, spec.Height-20, labelStyle) + "\n"
		}
	}

	return result
}

// RenderHorizontalBoxPlot renders a horizontal box plot
func RenderHorizontalBoxPlot(spec BoxPlotSpec) string {
	// Similar to vertical but with x/y swapped
	// Implementation would follow same pattern as vertical
	return ""
}

// valueOrDefault returns the value if not nil, otherwise returns default
func valueOrDefault(val *float64, def float64) float64 {
	if val != nil {
		return *val
	}
	return def
}
