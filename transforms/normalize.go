package transforms

import "math"

// Normalize creates a normalization transform that scales values to a standard range.
// Useful for comparing datasets with different scales or creating percentage-based views.
//
// Example:
//   data := []DataPoint{{Y: 10}, {Y: 20}, {Y: 30}}
//   normalized := Normalize(NormalizeOptions{Method: "percentage"})(data)
//   // Results in percentage of total: [16.67%, 33.33%, 50%]
func Normalize(opts NormalizeOptions) Transform {
	if opts.Method == "" {
		opts.Method = "percentage"
	}

	switch opts.Method {
	case "percentage":
		return NormalizePercentage()
	case "zscore":
		return NormalizeZScore()
	case "minmax":
		return NormalizeMinMax(0, 1)
	default:
		return NormalizePercentage()
	}
}

// NormalizePercentage converts values to percentages of the total
func NormalizePercentage() Transform {
	return func(data []DataPoint) []DataPoint {
		if len(data) == 0 {
			return nil
		}

		// Calculate total
		total := 0.0
		for _, d := range data {
			total += d.Y
		}

		if total == 0 {
			return data
		}

		// Convert to percentages
		result := make([]DataPoint, len(data))
		for i, d := range data {
			result[i] = d
			result[i].Y = (d.Y / total) * 100
			result[i].Value = result[i].Y
		}

		return result
	}
}

// NormalizeFraction converts values to fractions of the total (0-1)
func NormalizeFraction() Transform {
	return func(data []DataPoint) []DataPoint {
		if len(data) == 0 {
			return nil
		}

		// Calculate total
		total := 0.0
		for _, d := range data {
			total += d.Y
		}

		if total == 0 {
			return data
		}

		// Convert to fractions
		result := make([]DataPoint, len(data))
		for i, d := range data {
			result[i] = d
			result[i].Y = d.Y / total
			result[i].Value = result[i].Y
		}

		return result
	}
}

// NormalizeZScore converts values to z-scores (standard scores)
func NormalizeZScore() Transform {
	return func(data []DataPoint) []DataPoint {
		if len(data) == 0 {
			return nil
		}

		// Calculate mean
		values := make([]float64, len(data))
		for i, d := range data {
			values[i] = d.Y
		}
		mean := Mean(values)

		// Calculate standard deviation
		sumSquares := 0.0
		for _, v := range values {
			diff := v - mean
			sumSquares += diff * diff
		}
		stdDev := math.Sqrt(sumSquares / float64(len(values)))

		if stdDev == 0 {
			// All values are the same
			result := make([]DataPoint, len(data))
			for i, d := range data {
				result[i] = d
				result[i].Y = 0
				result[i].Value = 0
			}
			return result
		}

		// Convert to z-scores
		result := make([]DataPoint, len(data))
		for i, d := range data {
			result[i] = d
			result[i].Y = (d.Y - mean) / stdDev
			result[i].Value = result[i].Y
		}

		return result
	}
}

// NormalizeMinMax scales values to a specified range [min, max]
func NormalizeMinMax(targetMin, targetMax float64) Transform {
	return func(data []DataPoint) []DataPoint {
		if len(data) == 0 {
			return nil
		}

		// Find min and max
		values := make([]float64, len(data))
		for i, d := range data {
			values[i] = d.Y
		}
		dataMin := Min(values)
		dataMax := Max(values)

		if dataMin == dataMax {
			// All values are the same
			result := make([]DataPoint, len(data))
			midpoint := (targetMin + targetMax) / 2
			for i, d := range data {
				result[i] = d
				result[i].Y = midpoint
				result[i].Value = midpoint
			}
			return result
		}

		// Scale to [targetMin, targetMax]
		result := make([]DataPoint, len(data))
		for i, d := range data {
			result[i] = d
			normalized := (d.Y - dataMin) / (dataMax - dataMin)
			result[i].Y = targetMin + normalized*(targetMax-targetMin)
			result[i].Value = result[i].Y
		}

		return result
	}
}

// NormalizeByGroup normalizes within each group
func NormalizeByGroup(groupBy string, method string) Transform {
	return func(data []DataPoint) []DataPoint {
		if len(data) == 0 {
			return nil
		}

		// Group data
		groups := make(map[string][]int)
		for i, d := range data {
			var key string
			switch groupBy {
			case "Label":
				key = d.Label
			case "Group":
				key = d.Group
			default:
				key = d.Label
			}
			groups[key] = append(groups[key], i)
		}

		result := make([]DataPoint, len(data))
		copy(result, data)

		// Normalize each group
		for _, indices := range groups {
			if len(indices) == 0 {
				continue
			}

			// Extract group data
			groupData := make([]DataPoint, len(indices))
			for i, idx := range indices {
				groupData[i] = data[idx]
			}

			// Apply normalization
			var normalized []DataPoint
			switch method {
			case "percentage":
				normalized = NormalizePercentage()(groupData)
			case "zscore":
				normalized = NormalizeZScore()(groupData)
			case "minmax":
				normalized = NormalizeMinMax(0, 1)(groupData)
			default:
				normalized = groupData
			}

			// Copy back
			for i, idx := range indices {
				result[idx] = normalized[i]
			}
		}

		return result
	}
}

// Scale multiplies all Y values by a constant
func Scale(factor float64) Transform {
	return func(data []DataPoint) []DataPoint {
		result := make([]DataPoint, len(data))
		for i, d := range data {
			result[i] = d
			result[i].Y = d.Y * factor
			result[i].Value = result[i].Y
			result[i].Y0 = d.Y0 * factor
			result[i].Y1 = d.Y1 * factor
		}
		return result
	}
}

// Offset adds a constant to all Y values
func Offset(amount float64) Transform {
	return func(data []DataPoint) []DataPoint {
		result := make([]DataPoint, len(data))
		for i, d := range data {
			result[i] = d
			result[i].Y = d.Y + amount
			result[i].Value = result[i].Y
			result[i].Y0 = d.Y0 + amount
			result[i].Y1 = d.Y1 + amount
		}
		return result
	}
}

// Clamp restricts values to a specified range
func Clamp(min, max float64) Transform {
	return func(data []DataPoint) []DataPoint {
		result := make([]DataPoint, len(data))
		for i, d := range data {
			result[i] = d
			if result[i].Y < min {
				result[i].Y = min
			}
			if result[i].Y > max {
				result[i].Y = max
			}
			result[i].Value = result[i].Y
		}
		return result
	}
}

// Abs converts all Y values to their absolute values
func Abs() Transform {
	return func(data []DataPoint) []DataPoint {
		result := make([]DataPoint, len(data))
		for i, d := range data {
			result[i] = d
			result[i].Y = math.Abs(d.Y)
			result[i].Value = result[i].Y
		}
		return result
	}
}

// Log applies logarithmic transformation
func Log(base float64) Transform {
	return func(data []DataPoint) []DataPoint {
		if base <= 0 || base == 1 {
			base = 10 // Default to log10
		}

		logBase := math.Log(base)
		result := make([]DataPoint, len(data))

		for i, d := range data {
			result[i] = d
			if d.Y > 0 {
				result[i].Y = math.Log(d.Y) / logBase
				result[i].Value = result[i].Y
			} else {
				// Handle non-positive values
				result[i].Y = 0
				result[i].Value = 0
			}
		}

		return result
	}
}

// Sqrt applies square root transformation
func Sqrt() Transform {
	return func(data []DataPoint) []DataPoint {
		result := make([]DataPoint, len(data))
		for i, d := range data {
			result[i] = d
			if d.Y >= 0 {
				result[i].Y = math.Sqrt(d.Y)
				result[i].Value = result[i].Y
			} else {
				result[i].Y = 0
				result[i].Value = 0
			}
		}
		return result
	}
}
