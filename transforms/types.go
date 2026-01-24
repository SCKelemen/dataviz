package transforms

import "time"

// DataPoint represents a generic data point that can be transformed.
// Transforms operate on slices of DataPoints and return transformed data.
type DataPoint struct {
	X      interface{} // X value (can be time.Time, float64, string, etc.)
	Y      float64     // Y value (numeric)
	Y0     float64     // Baseline Y value (for stacking)
	Y1     float64     // Top Y value (for stacking)
	Label  string      // Category label
	Value  float64     // Generic numeric value
	Count  int         // Count for aggregations
	Group  string      // Group identifier
	Index  int         // Original index
	Data   interface{} // Original data reference
}

// Transform is a function that transforms a slice of data points.
type Transform func([]DataPoint) []DataPoint

// AggregateFunc defines how to aggregate numeric values.
type AggregateFunc func(values []float64) float64

// Common aggregation functions
var (
	// Sum aggregates by summing all values
	Sum AggregateFunc = func(values []float64) float64 {
		sum := 0.0
		for _, v := range values {
			sum += v
		}
		return sum
	}

	// Mean aggregates by calculating the mean (average)
	Mean AggregateFunc = func(values []float64) float64 {
		if len(values) == 0 {
			return 0
		}
		return Sum(values) / float64(len(values))
	}

	// Max aggregates by taking the maximum value
	Max AggregateFunc = func(values []float64) float64 {
		if len(values) == 0 {
			return 0
		}
		max := values[0]
		for _, v := range values[1:] {
			if v > max {
				max = v
			}
		}
		return max
	}

	// Min aggregates by taking the minimum value
	Min AggregateFunc = func(values []float64) float64 {
		if len(values) == 0 {
			return 0
		}
		min := values[0]
		for _, v := range values[1:] {
			if v < min {
				min = v
			}
		}
		return min
	}

	// Count aggregates by counting the number of values
	Count AggregateFunc = func(values []float64) float64 {
		return float64(len(values))
	}

	// Median aggregates by calculating the median
	Median AggregateFunc = func(values []float64) float64 {
		if len(values) == 0 {
			return 0
		}
		// Create a copy and sort
		sorted := make([]float64, len(values))
		copy(sorted, values)
		// Simple bubble sort (fine for small datasets)
		for i := 0; i < len(sorted); i++ {
			for j := i + 1; j < len(sorted); j++ {
				if sorted[i] > sorted[j] {
					sorted[i], sorted[j] = sorted[j], sorted[i]
				}
			}
		}
		mid := len(sorted) / 2
		if len(sorted)%2 == 0 {
			return (sorted[mid-1] + sorted[mid]) / 2
		}
		return sorted[mid]
	}
)

// BinOptions configures binning behavior
type BinOptions struct {
	// Thresholds specifies explicit bin edges (overrides Count)
	Thresholds []float64

	// Count specifies the approximate number of bins (default: 10)
	Count int

	// Domain specifies the [min, max] range to bin over
	Domain [2]float64

	// Nice rounds bin edges to nice numbers
	Nice bool
}

// GroupOptions configures grouping behavior
type GroupOptions struct {
	// By specifies the field to group by ("X", "Label", "Group")
	By string

	// Aggregate specifies how to aggregate Y values
	Aggregate AggregateFunc

	// Sort specifies whether to sort groups (by key or value)
	Sort string // "key", "value", or ""
}

// StackOptions configures stacking behavior
type StackOptions struct {
	// By specifies the field to group by for stacking
	By string

	// Order specifies the stacking order ("ascending", "descending", "none")
	Order string

	// Offset specifies the baseline ("zero", "center", "normalize")
	Offset string
}

// SmoothOptions configures smoothing behavior
type SmoothOptions struct {
	// Method specifies the smoothing method ("movingAverage", "loess", "exponential")
	Method string

	// WindowSize specifies the window size for moving averages
	WindowSize int

	// Bandwidth specifies the bandwidth for LOESS smoothing (0-1)
	Bandwidth float64

	// Alpha specifies the smoothing factor for exponential smoothing (0-1)
	Alpha float64
}

// NormalizeOptions configures normalization behavior
type NormalizeOptions struct {
	// Method specifies normalization method ("percentage", "zscore", "minmax")
	Method string

	// By specifies the field to normalize by ("group", "all")
	By string
}

// TimeSeriesPoint represents a time-series data point
type TimeSeriesPoint struct {
	Time  time.Time
	Value float64
	Label string
}

// ToDataPoints converts TimeSeriesPoints to DataPoints
func ToDataPoints(points []TimeSeriesPoint) []DataPoint {
	result := make([]DataPoint, len(points))
	for i, p := range points {
		result[i] = DataPoint{
			X:     p.Time,
			Y:     p.Value,
			Value: p.Value,
			Label: p.Label,
			Index: i,
		}
	}
	return result
}

// FromDataPoints converts DataPoints back to TimeSeriesPoints
func FromDataPoints(points []DataPoint) []TimeSeriesPoint {
	result := make([]TimeSeriesPoint, len(points))
	for i, p := range points {
		var t time.Time
		if tt, ok := p.X.(time.Time); ok {
			t = tt
		}
		result[i] = TimeSeriesPoint{
			Time:  t,
			Value: p.Y,
			Label: p.Label,
		}
	}
	return result
}
