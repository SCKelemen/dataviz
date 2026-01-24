package transforms

import (
	"math"
	"sort"
)

// Bin creates a binning transform that groups continuous data into discrete bins.
// Useful for creating histograms and frequency distributions.
//
// Example:
//   data := []DataPoint{{Y: 1.5}, {Y: 2.3}, {Y: 5.7}, {Y: 8.1}}
//   binned := Bin(BinOptions{Count: 3})(data)
//   // Results in 3 bins with counts
func Bin(opts BinOptions) Transform {
	return func(data []DataPoint) []DataPoint {
		if len(data) == 0 {
			return nil
		}

		// Set defaults
		if opts.Count == 0 {
			opts.Count = 10
		}

		// Extract values and find domain
		values := make([]float64, len(data))
		for i, d := range data {
			values[i] = d.Y
		}

		domain := opts.Domain
		if domain[0] == 0 && domain[1] == 0 {
			domain = [2]float64{Min(values), Max(values)}
		}

		// Generate thresholds
		thresholds := opts.Thresholds
		if len(thresholds) == 0 {
			thresholds = generateThresholds(domain, opts.Count, opts.Nice)
		} else {
			// Sort thresholds
			sort.Float64s(thresholds)
		}

		if len(thresholds) == 0 {
			return nil
		}

		// Create bins
		bins := make([]DataPoint, len(thresholds)-1)
		for i := 0; i < len(thresholds)-1; i++ {
			bins[i] = DataPoint{
				Y0:    thresholds[i],
				Y1:    thresholds[i+1],
				Y:     (thresholds[i] + thresholds[i+1]) / 2, // Midpoint
				X:     (thresholds[i] + thresholds[i+1]) / 2,
				Count: 0,
				Index: i,
			}
		}

		// Count values in each bin
		for _, v := range values {
			// Find which bin this value belongs to
			for i := 0; i < len(bins); i++ {
				if v >= bins[i].Y0 && v < bins[i].Y1 {
					bins[i].Count++
					bins[i].Value = float64(bins[i].Count)
					break
				}
				// Last bin includes the upper bound
				if i == len(bins)-1 && v == bins[i].Y1 {
					bins[i].Count++
					bins[i].Value = float64(bins[i].Count)
					break
				}
			}
		}

		return bins
	}
}

// generateThresholds generates bin edges for the given domain and count
func generateThresholds(domain [2]float64, count int, nice bool) []float64 {
	if count <= 0 {
		return nil
	}

	min, max := domain[0], domain[1]
	if min == max {
		// Handle single value
		return []float64{min - 0.5, min + 0.5}
	}

	if nice {
		// Calculate nice step size
		range_ := max - min
		step := niceNumber(range_/float64(count), false)
		min = math.Floor(min/step) * step
		max = math.Ceil(max/step) * step
		count = int(math.Ceil((max - min) / step))
	}

	thresholds := make([]float64, count+1)
	step := (max - min) / float64(count)
	for i := 0; i <= count; i++ {
		thresholds[i] = min + float64(i)*step
	}

	return thresholds
}

// niceNumber rounds a number to a nice round value
func niceNumber(value float64, round bool) float64 {
	exponent := math.Floor(math.Log10(value))
	fraction := value / math.Pow(10, exponent)
	var niceFraction float64

	if round {
		if fraction < 1.5 {
			niceFraction = 1
		} else if fraction < 3 {
			niceFraction = 2
		} else if fraction < 7 {
			niceFraction = 5
		} else {
			niceFraction = 10
		}
	} else {
		if fraction <= 1 {
			niceFraction = 1
		} else if fraction <= 2 {
			niceFraction = 2
		} else if fraction <= 5 {
			niceFraction = 5
		} else {
			niceFraction = 10
		}
	}

	return niceFraction * math.Pow(10, exponent)
}

// BinCount creates a simple frequency count binning transform
func BinCount(binSize float64) Transform {
	return func(data []DataPoint) []DataPoint {
		if len(data) == 0 || binSize <= 0 {
			return nil
		}

		// Find min/max
		values := make([]float64, len(data))
		for i, d := range data {
			values[i] = d.Y
		}
		min, max := Min(values), Max(values)

		// Calculate number of bins
		numBins := int(math.Ceil((max - min) / binSize))
		if numBins <= 0 {
			numBins = 1
		}

		// Create bins
		bins := make(map[int]*DataPoint)
		for _, d := range data {
			binIndex := int(math.Floor((d.Y - min) / binSize))
			if binIndex >= numBins {
				binIndex = numBins - 1
			}

			if bins[binIndex] == nil {
				binStart := min + float64(binIndex)*binSize
				bins[binIndex] = &DataPoint{
					Y0:    binStart,
					Y1:    binStart + binSize,
					Y:     binStart + binSize/2,
					X:     binStart + binSize/2,
					Count: 0,
					Index: binIndex,
				}
			}
			bins[binIndex].Count++
			bins[binIndex].Value = float64(bins[binIndex].Count)
		}

		// Convert map to sorted slice
		result := make([]DataPoint, 0, len(bins))
		indices := make([]int, 0, len(bins))
		for idx := range bins {
			indices = append(indices, idx)
		}
		sort.Ints(indices)

		for _, idx := range indices {
			result = append(result, *bins[idx])
		}

		return result
	}
}

// Histogram is an alias for Bin with sensible defaults
func Histogram(thresholds ...float64) Transform {
	opts := BinOptions{
		Count: 10,
		Nice:  true,
	}
	if len(thresholds) > 0 {
		opts.Thresholds = thresholds
	}
	return Bin(opts)
}
