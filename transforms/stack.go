package transforms

import "sort"

// Stack creates a stacking transform that computes Y0 and Y1 for stacked visualizations.
// Essential for stacked bar charts, area charts, and stream graphs.
//
// Example:
//   data := []DataPoint{
//     {Label: "2020", Group: "A", Y: 10},
//     {Label: "2020", Group: "B", Y: 15},
//     {Label: "2021", Group: "A", Y: 12},
//     {Label: "2021", Group: "B", Y: 18},
//   }
//   stacked := Stack(StackOptions{By: "Label"})(data)
//   // Results in Y0 and Y1 computed for stacking
func Stack(opts StackOptions) Transform {
	return func(data []DataPoint) []DataPoint {
		if len(data) == 0 {
			return nil
		}

		// Set defaults
		if opts.By == "" {
			opts.By = "Label"
		}
		if opts.Order == "" {
			opts.Order = "none"
		}
		if opts.Offset == "" {
			opts.Offset = "zero"
		}

		// Group data by X value
		groups := make(map[string][]int)
		for i, d := range data {
			var key string
			if opts.By == "Label" {
				key = d.Label
			} else if opts.By == "Group" {
				key = d.Group
			} else {
				key = d.Label
			}

			// Use X value as the stack position
			if xStr, ok := d.X.(string); ok {
				key = xStr
			}

			groups[key] = append(groups[key], i)
		}

		// Apply stacking to each group
		result := make([]DataPoint, len(data))
		copy(result, data)

		for _, indices := range groups {
			if len(indices) == 0 {
				continue
			}

			// Sort indices by order if requested
			switch opts.Order {
			case "ascending":
				sort.Slice(indices, func(i, j int) bool {
					return data[indices[i]].Y < data[indices[j]].Y
				})
			case "descending":
				sort.Slice(indices, func(i, j int) bool {
					return data[indices[i]].Y > data[indices[j]].Y
				})
			}

			// Calculate total for this stack position
			total := 0.0
			for _, idx := range indices {
				total += data[idx].Y
			}

			// Apply offset
			var offset float64
			switch opts.Offset {
			case "zero":
				offset = 0
			case "center":
				offset = -total / 2
			case "normalize":
				// Will normalize after stacking
				offset = 0
			}

			// Stack values
			baseline := offset
			for _, idx := range indices {
				result[idx].Y0 = baseline
				result[idx].Y1 = baseline + data[idx].Y
				baseline = result[idx].Y1
			}

			// Normalize if requested
			if opts.Offset == "normalize" && total > 0 {
				for _, idx := range indices {
					result[idx].Y0 /= total
					result[idx].Y1 /= total
				}
			}
		}

		return result
	}
}

// StackZero is a convenience function for zero-baseline stacking
func StackZero(by string) Transform {
	return Stack(StackOptions{
		By:     by,
		Offset: "zero",
		Order:  "none",
	})
}

// StackCenter creates a centered (diverging) stack
func StackCenter(by string) Transform {
	return Stack(StackOptions{
		By:     by,
		Offset: "center",
		Order:  "none",
	})
}

// StackNormalize creates a normalized (100%) stack
func StackNormalize(by string) Transform {
	return Stack(StackOptions{
		By:     by,
		Offset: "normalize",
		Order:  "none",
	})
}

// Dodge creates side-by-side positioning for grouped bars
// Instead of stacking, places items next to each other
func Dodge(by string, padding float64) Transform {
	return func(data []DataPoint) []DataPoint {
		if len(data) == 0 {
			return nil
		}

		// Group data by X value
		groups := make(map[string][]int)
		for i, d := range data {
			var key string
			if xStr, ok := d.X.(string); ok {
				key = xStr
			} else {
				key = d.Label
			}
			groups[key] = append(groups[key], i)
		}

		result := make([]DataPoint, len(data))
		copy(result, data)

		// Calculate dodge positions
		for _, indices := range groups {
			n := len(indices)
			if n == 0 {
				continue
			}

			// Calculate width and offset for each bar
			totalWidth := 1.0
			barWidth := totalWidth / float64(n)
			actualWidth := barWidth * (1 - padding)

			for i, idx := range indices {
				// Store dodge offset in a custom field
				result[idx].Y0 = float64(i)*barWidth + (barWidth-actualWidth)/2
				result[idx].Y1 = result[idx].Y0 + actualWidth
				// Original Y value stays the same
			}
		}

		return result
	}
}

// Expand expands stacked values to fill the entire range
func Expand() Transform {
	return func(data []DataPoint) []DataPoint {
		if len(data) == 0 {
			return nil
		}

		// Find max Y1 value across all stacks
		maxY1 := 0.0
		for _, d := range data {
			if d.Y1 > maxY1 {
				maxY1 = d.Y1
			}
		}

		if maxY1 == 0 {
			return data
		}

		// Scale all Y0 and Y1 values
		result := make([]DataPoint, len(data))
		for i, d := range data {
			result[i] = d
			result[i].Y0 = d.Y0 / maxY1
			result[i].Y1 = d.Y1 / maxY1
		}

		return result
	}
}

// Unstack reverses stacking by resetting Y0 and Y1
func Unstack() Transform {
	return func(data []DataPoint) []DataPoint {
		result := make([]DataPoint, len(data))
		for i, d := range data {
			result[i] = d
			result[i].Y0 = 0
			result[i].Y1 = d.Y
		}
		return result
	}
}
