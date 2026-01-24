package transforms

import (
	"sort"
	"time"
)

// GroupBy creates a grouping transform that aggregates data by a field.
// Useful for creating summary statistics grouped by category.
//
// Example:
//   data := []DataPoint{
//     {Label: "A", Y: 10},
//     {Label: "A", Y: 20},
//     {Label: "B", Y: 15},
//   }
//   grouped := GroupBy(GroupOptions{By: "Label", Aggregate: Sum})(data)
//   // Results in: [{Label: "A", Y: 30}, {Label: "B", Y: 15}]
func GroupBy(opts GroupOptions) Transform {
	return func(data []DataPoint) []DataPoint {
		if len(data) == 0 {
			return nil
		}

		// Set defaults
		if opts.Aggregate == nil {
			opts.Aggregate = Sum
		}
		if opts.By == "" {
			opts.By = "Label"
		}

		// Group data by the specified field
		groups := make(map[string][]float64)
		groupData := make(map[string]DataPoint)

		for i, d := range data {
			var key string
			switch opts.By {
			case "Label":
				key = d.Label
			case "Group":
				key = d.Group
			case "X":
				if t, ok := d.X.(time.Time); ok {
					key = t.Format("2006-01-02")
				} else if s, ok := d.X.(string); ok {
					key = s
				} else {
					key = "default"
				}
			default:
				key = "default"
			}

			groups[key] = append(groups[key], d.Y)
			if _, exists := groupData[key]; !exists {
				groupData[key] = DataPoint{
					Label: key,
					Group: key,
					X:     d.X,
					Index: i,
				}
			}
		}

		// Aggregate each group
		result := make([]DataPoint, 0, len(groups))
		for key, values := range groups {
			point := groupData[key]
			point.Y = opts.Aggregate(values)
			point.Value = point.Y
			point.Count = len(values)
			result = append(result, point)
		}

		// Sort if requested
		switch opts.Sort {
		case "key":
			sort.Slice(result, func(i, j int) bool {
				return result[i].Label < result[j].Label
			})
		case "value":
			sort.Slice(result, func(i, j int) bool {
				return result[i].Y > result[j].Y
			})
		}

		return result
	}
}

// Reduce aggregates all data points into a single value using the given function
func Reduce(fn AggregateFunc) Transform {
	return func(data []DataPoint) []DataPoint {
		if len(data) == 0 {
			return nil
		}

		values := make([]float64, len(data))
		for i, d := range data {
			values[i] = d.Y
		}

		result := DataPoint{
			Y:     fn(values),
			Value: fn(values),
			Count: len(data),
			Label: "aggregate",
		}

		return []DataPoint{result}
	}
}

// Filter creates a transform that filters data points based on a predicate
func Filter(predicate func(DataPoint) bool) Transform {
	return func(data []DataPoint) []DataPoint {
		result := make([]DataPoint, 0, len(data))
		for _, d := range data {
			if predicate(d) {
				result = append(result, d)
			}
		}
		return result
	}
}

// Map creates a transform that maps each data point through a function
func Map(fn func(DataPoint) DataPoint) Transform {
	return func(data []DataPoint) []DataPoint {
		result := make([]DataPoint, len(data))
		for i, d := range data {
			result[i] = fn(d)
		}
		return result
	}
}

// Sort creates a transform that sorts data points
func Sort(by string, ascending bool) Transform {
	return func(data []DataPoint) []DataPoint {
		result := make([]DataPoint, len(data))
		copy(result, data)

		sort.Slice(result, func(i, j int) bool {
			var less bool
			switch by {
			case "Y", "value":
				less = result[i].Y < result[j].Y
			case "X":
				if t1, ok := result[i].X.(time.Time); ok {
					if t2, ok := result[j].X.(time.Time); ok {
						less = t1.Before(t2)
					}
				}
			case "Label":
				less = result[i].Label < result[j].Label
			case "Count":
				less = result[i].Count < result[j].Count
			default:
				less = result[i].Index < result[j].Index
			}

			if !ascending {
				less = !less
			}
			return less
		})

		return result
	}
}

// Top returns the top N data points by value
func Top(n int) Transform {
	return func(data []DataPoint) []DataPoint {
		if n <= 0 || n >= len(data) {
			return data
		}

		// Sort by Y descending
		sorted := Sort("Y", false)(data)
		return sorted[:n]
	}
}

// Percentile calculates percentiles for the Y values
func Percentile(p float64) Transform {
	return func(data []DataPoint) []DataPoint {
		if len(data) == 0 {
			return nil
		}

		values := make([]float64, len(data))
		for i, d := range data {
			values[i] = d.Y
		}

		sorted := make([]float64, len(values))
		copy(sorted, values)
		sort.Float64s(sorted)

		index := p * float64(len(sorted)-1)
		lower := int(index)
		upper := lower + 1
		if upper >= len(sorted) {
			upper = len(sorted) - 1
		}

		weight := index - float64(lower)
		value := sorted[lower]*(1-weight) + sorted[upper]*weight

		return []DataPoint{{
			Y:     value,
			Value: value,
			Label: "percentile",
		}}
	}
}

// Cumulative creates a cumulative sum transform
func Cumulative() Transform {
	return func(data []DataPoint) []DataPoint {
		result := make([]DataPoint, len(data))
		cumSum := 0.0

		for i, d := range data {
			cumSum += d.Y
			result[i] = d
			result[i].Y = cumSum
			result[i].Value = cumSum
		}

		return result
	}
}

// Window applies a windowed aggregation (rolling window)
func Window(size int, fn AggregateFunc) Transform {
	return func(data []DataPoint) []DataPoint {
		if len(data) == 0 || size <= 0 {
			return data
		}

		result := make([]DataPoint, len(data))
		for i := range data {
			// Calculate window bounds
			start := i - size + 1
			if start < 0 {
				start = 0
			}
			end := i + 1

			// Extract window values
			windowValues := make([]float64, end-start)
			for j := start; j < end; j++ {
				windowValues[j-start] = data[j].Y
			}

			// Apply aggregation
			result[i] = data[i]
			result[i].Y = fn(windowValues)
			result[i].Value = result[i].Y
		}

		return result
	}
}
