package transforms

import "math"

// Smooth creates a smoothing transform that applies various smoothing algorithms.
// Useful for trend lines, noise reduction, and pattern identification.
//
// Example:
//   data := []DataPoint{{Y: 1}, {Y: 5}, {Y: 2}, {Y: 8}, {Y: 3}}
//   smoothed := Smooth(SmoothOptions{Method: "movingAverage", WindowSize: 3})(data)
func Smooth(opts SmoothOptions) Transform {
	if opts.Method == "" {
		opts.Method = "movingAverage"
	}

	switch opts.Method {
	case "movingAverage":
		return MovingAverage(opts.WindowSize)
	case "exponential":
		return ExponentialSmoothing(opts.Alpha)
	case "loess":
		return Loess(opts.Bandwidth)
	default:
		return MovingAverage(3)
	}
}

// MovingAverage creates a simple moving average transform
func MovingAverage(windowSize int) Transform {
	return func(data []DataPoint) []DataPoint {
		if len(data) == 0 {
			return nil
		}

		if windowSize <= 0 {
			windowSize = 3
		}
		if windowSize > len(data) {
			windowSize = len(data)
		}

		result := make([]DataPoint, len(data))

		for i := range data {
			// Calculate window bounds
			start := i - windowSize/2
			end := start + windowSize

			if start < 0 {
				start = 0
				end = windowSize
			}
			if end > len(data) {
				end = len(data)
				start = end - windowSize
				if start < 0 {
					start = 0
				}
			}

			// Calculate average
			sum := 0.0
			count := 0
			for j := start; j < end; j++ {
				sum += data[j].Y
				count++
			}

			result[i] = data[i]
			if count > 0 {
				result[i].Y = sum / float64(count)
				result[i].Value = result[i].Y
			}
		}

		return result
	}
}

// WeightedMovingAverage applies weighted moving average with custom weights
func WeightedMovingAverage(weights []float64) Transform {
	return func(data []DataPoint) []DataPoint {
		if len(data) == 0 || len(weights) == 0 {
			return nil
		}

		windowSize := len(weights)
		result := make([]DataPoint, len(data))

		// Normalize weights
		weightSum := 0.0
		for _, w := range weights {
			weightSum += w
		}
		if weightSum == 0 {
			weightSum = 1
		}

		for i := range data {
			// Calculate window bounds (centered)
			start := i - windowSize/2
			end := start + windowSize

			if start < 0 {
				start = 0
			}
			if end > len(data) {
				end = len(data)
			}

			// Calculate weighted average
			sum := 0.0
			actualWeightSum := 0.0
			for j := start; j < end; j++ {
				weightIdx := j - start
				if weightIdx < len(weights) {
					sum += data[j].Y * weights[weightIdx]
					actualWeightSum += weights[weightIdx]
				}
			}

			result[i] = data[i]
			if actualWeightSum > 0 {
				result[i].Y = sum / actualWeightSum
				result[i].Value = result[i].Y
			}
		}

		return result
	}
}

// ExponentialSmoothing applies exponential smoothing
func ExponentialSmoothing(alpha float64) Transform {
	return func(data []DataPoint) []DataPoint {
		if len(data) == 0 {
			return nil
		}

		// Default alpha
		if alpha <= 0 || alpha >= 1 {
			alpha = 0.3
		}

		result := make([]DataPoint, len(data))
		result[0] = data[0]

		for i := 1; i < len(data); i++ {
			result[i] = data[i]
			result[i].Y = alpha*data[i].Y + (1-alpha)*result[i-1].Y
			result[i].Value = result[i].Y
		}

		return result
	}
}

// Loess applies LOESS (Locally Estimated Scatterplot Smoothing)
// This is a simplified version of LOESS
func Loess(bandwidth float64) Transform {
	return func(data []DataPoint) []DataPoint {
		if len(data) == 0 {
			return nil
		}

		// Default bandwidth
		if bandwidth <= 0 || bandwidth > 1 {
			bandwidth = 0.3
		}

		result := make([]DataPoint, len(data))
		windowSize := int(float64(len(data)) * bandwidth)
		if windowSize < 2 {
			windowSize = 2
		}
		if windowSize > len(data) {
			windowSize = len(data)
		}

		for i := range data {
			// Find nearest neighbors
			distances := make([]float64, len(data))
			for j := range data {
				distances[j] = math.Abs(float64(i - j))
			}

			// Find the k nearest neighbors
			indices := make([]int, len(data))
			for j := range indices {
				indices[j] = j
			}

			// Sort by distance
			for j := 0; j < windowSize; j++ {
				for k := j + 1; k < len(indices); k++ {
					if distances[indices[k]] < distances[indices[j]] {
						indices[j], indices[k] = indices[k], indices[j]
					}
				}
			}

			// Calculate weighted average using tricube kernel
			maxDist := distances[indices[windowSize-1]]
			if maxDist == 0 {
				maxDist = 1
			}

			weightedSum := 0.0
			weightSum := 0.0

			for j := 0; j < windowSize; j++ {
				idx := indices[j]
				dist := distances[idx] / maxDist
				weight := tricube(dist)
				weightedSum += data[idx].Y * weight
				weightSum += weight
			}

			result[i] = data[i]
			if weightSum > 0 {
				result[i].Y = weightedSum / weightSum
				result[i].Value = result[i].Y
			}
		}

		return result
	}
}

// tricube kernel function for LOESS
func tricube(x float64) float64 {
	if x >= 1 {
		return 0
	}
	cube := 1 - x*x*x
	return cube * cube * cube
}

// SavitzkyGolay applies Savitzky-Golay smoothing (simplified)
func SavitzkyGolay(windowSize, polyOrder int) Transform {
	return func(data []DataPoint) []DataPoint {
		if len(data) == 0 || windowSize < polyOrder+1 {
			return data
		}

		// For simplicity, use weighted moving average approximation
		// A full implementation would fit polynomials locally
		weights := make([]float64, windowSize)
		for i := range weights {
			weights[i] = 1.0
		}

		return WeightedMovingAverage(weights)(data)
	}
}

// Interpolate fills in missing values using linear interpolation
func Interpolate() Transform {
	return func(data []DataPoint) []DataPoint {
		if len(data) < 2 {
			return data
		}

		result := make([]DataPoint, len(data))
		copy(result, data)

		// Find gaps and interpolate
		for i := 1; i < len(result)-1; i++ {
			// Check if Y is zero or missing (you might want a different condition)
			if result[i].Y == 0 {
				// Find previous and next non-zero values
				prev := i - 1
				next := i + 1

				for next < len(result) && result[next].Y == 0 {
					next++
				}

				if next < len(result) && prev >= 0 {
					// Linear interpolation
					ratio := float64(i-prev) / float64(next-prev)
					result[i].Y = result[prev].Y + ratio*(result[next].Y-result[prev].Y)
					result[i].Value = result[i].Y
				}
			}
		}

		return result
	}
}

// Downsample reduces the number of points by sampling every nth point
func Downsample(n int) Transform {
	return func(data []DataPoint) []DataPoint {
		if n <= 1 || len(data) == 0 {
			return data
		}

		result := make([]DataPoint, 0, len(data)/n+1)
		for i := 0; i < len(data); i += n {
			result = append(result, data[i])
		}

		return result
	}
}
