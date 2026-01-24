package transforms

import "math"

// Rolling provides pandas-style rolling window operations
type Rolling struct {
	window int
	minPeriods int
	center bool
}

// NewRolling creates a new rolling window operator
func NewRolling(window int) *Rolling {
	return &Rolling{
		window:     window,
		minPeriods: 1,
		center:     false,
	}
}

// MinPeriods sets the minimum number of observations required
func (r *Rolling) MinPeriods(n int) *Rolling {
	r.minPeriods = n
	return r
}

// Center centers the window labels
func (r *Rolling) Center(center bool) *Rolling {
	r.center = center
	return r
}

// Mean calculates rolling mean
func (r *Rolling) Mean() Transform {
	return r.apply(Mean)
}

// Sum calculates rolling sum
func (r *Rolling) Sum() Transform {
	return r.apply(Sum)
}

// Min calculates rolling minimum
func (r *Rolling) Min() Transform {
	return r.apply(Min)
}

// Max calculates rolling maximum
func (r *Rolling) Max() Transform {
	return r.apply(Max)
}

// Std calculates rolling standard deviation
func (r *Rolling) Std() Transform {
	return r.apply(func(values []float64) float64 {
		if len(values) == 0 {
			return 0
		}
		mean := Mean(values)
		sumSquares := 0.0
		for _, v := range values {
			diff := v - mean
			sumSquares += diff * diff
		}
		return math.Sqrt(sumSquares / float64(len(values)))
	})
}

// Var calculates rolling variance
func (r *Rolling) Var() Transform {
	return r.apply(func(values []float64) float64 {
		if len(values) == 0 {
			return 0
		}
		mean := Mean(values)
		sumSquares := 0.0
		for _, v := range values {
			diff := v - mean
			sumSquares += diff * diff
		}
		return sumSquares / float64(len(values))
	})
}

// Median calculates rolling median
func (r *Rolling) Median() Transform {
	return r.apply(Median)
}

// Quantile calculates rolling quantile
func (r *Rolling) Quantile(q float64) Transform {
	return r.apply(func(values []float64) float64 {
		if len(values) == 0 {
			return 0
		}
		return Percentile(q)([]DataPoint{{Y: 0}})[0].Y
	})
}

// Skew calculates rolling skewness
func (r *Rolling) Skew() Transform {
	return r.apply(func(values []float64) float64 {
		if len(values) < 3 {
			return 0
		}
		mean := Mean(values)
		n := float64(len(values))

		m2 := 0.0
		m3 := 0.0
		for _, v := range values {
			diff := v - mean
			m2 += diff * diff
			m3 += diff * diff * diff
		}

		m2 /= n
		m3 /= n

		if m2 == 0 {
			return 0
		}

		return m3 / math.Pow(m2, 1.5)
	})
}

// Kurt calculates rolling kurtosis
func (r *Rolling) Kurt() Transform {
	return r.apply(func(values []float64) float64 {
		if len(values) < 4 {
			return 0
		}
		mean := Mean(values)
		n := float64(len(values))

		m2 := 0.0
		m4 := 0.0
		for _, v := range values {
			diff := v - mean
			diff2 := diff * diff
			m2 += diff2
			m4 += diff2 * diff2
		}

		m2 /= n
		m4 /= n

		if m2 == 0 {
			return 0
		}

		return (m4 / (m2 * m2)) - 3.0 // Excess kurtosis
	})
}

// Apply applies a custom aggregation function
func (r *Rolling) Apply(fn AggregateFunc) Transform {
	return r.apply(fn)
}

// apply is the internal implementation
func (r *Rolling) apply(fn AggregateFunc) Transform {
	return func(data []DataPoint) []DataPoint {
		if len(data) == 0 || r.window <= 0 {
			return data
		}

		result := make([]DataPoint, len(data))

		for i := range data {
			// Calculate window bounds
			var start, end int

			if r.center {
				// Center the window on current point
				halfWindow := r.window / 2
				start = i - halfWindow
				end = i + (r.window - halfWindow)
			} else {
				// Window ends at current point (backward-looking)
				start = i - r.window + 1
				end = i + 1
			}

			// Clamp to valid range
			if start < 0 {
				start = 0
			}
			if end > len(data) {
				end = len(data)
			}

			// Extract window values
			windowSize := end - start
			if windowSize < r.minPeriods {
				// Not enough observations
				result[i] = data[i]
				result[i].Y = math.NaN()
				result[i].Value = math.NaN()
				continue
			}

			values := make([]float64, windowSize)
			for j := start; j < end; j++ {
				values[j-start] = data[j].Y
			}

			// Apply aggregation
			result[i] = data[i]
			result[i].Y = fn(values)
			result[i].Value = result[i].Y
		}

		return result
	}
}

// Expanding provides pandas-style expanding window operations
type Expanding struct {
	minPeriods int
}

// NewExpanding creates a new expanding window operator
func NewExpanding() *Expanding {
	return &Expanding{
		minPeriods: 1,
	}
}

// MinPeriods sets the minimum number of observations required
func (e *Expanding) MinPeriods(n int) *Expanding {
	e.minPeriods = n
	return e
}

// Mean calculates expanding mean
func (e *Expanding) Mean() Transform {
	return e.apply(Mean)
}

// Sum calculates expanding sum (cumulative sum)
func (e *Expanding) Sum() Transform {
	return e.apply(Sum)
}

// Min calculates expanding minimum
func (e *Expanding) Min() Transform {
	return e.apply(Min)
}

// Max calculates expanding maximum
func (e *Expanding) Max() Transform {
	return e.apply(Max)
}

// Std calculates expanding standard deviation
func (e *Expanding) Std() Transform {
	return e.apply(func(values []float64) float64 {
		if len(values) == 0 {
			return 0
		}
		mean := Mean(values)
		sumSquares := 0.0
		for _, v := range values {
			diff := v - mean
			sumSquares += diff * diff
		}
		return math.Sqrt(sumSquares / float64(len(values)))
	})
}

// Var calculates expanding variance
func (e *Expanding) Var() Transform {
	return e.apply(func(values []float64) float64 {
		if len(values) == 0 {
			return 0
		}
		mean := Mean(values)
		sumSquares := 0.0
		for _, v := range values {
			diff := v - mean
			sumSquares += diff * diff
		}
		return sumSquares / float64(len(values))
	})
}

// Count calculates expanding count
func (e *Expanding) Count() Transform {
	return e.apply(Count)
}

// Apply applies a custom aggregation function
func (e *Expanding) Apply(fn AggregateFunc) Transform {
	return e.apply(fn)
}

// apply is the internal implementation
func (e *Expanding) apply(fn AggregateFunc) Transform {
	return func(data []DataPoint) []DataPoint {
		if len(data) == 0 {
			return nil
		}

		result := make([]DataPoint, len(data))

		for i := range data {
			// Window from start to current point
			windowSize := i + 1

			if windowSize < e.minPeriods {
				// Not enough observations
				result[i] = data[i]
				result[i].Y = math.NaN()
				result[i].Value = math.NaN()
				continue
			}

			// Extract values from start to current
			values := make([]float64, windowSize)
			for j := 0; j <= i; j++ {
				values[j] = data[j].Y
			}

			// Apply aggregation
			result[i] = data[i]
			result[i].Y = fn(values)
			result[i].Value = result[i].Y
		}

		return result
	}
}

// EWM provides pandas-style exponentially weighted functions
type EWM struct {
	alpha      float64
	adjust     bool
	ignoreNA   bool
	minPeriods int
}

// NewEWM creates a new exponentially weighted operator
func NewEWM(alpha float64) *EWM {
	return &EWM{
		alpha:      alpha,
		adjust:     true,
		ignoreNA:   false,
		minPeriods: 0,
	}
}

// Adjust sets whether to use adjustment in weights
func (ewm *EWM) Adjust(adjust bool) *EWM {
	ewm.adjust = adjust
	return ewm
}

// IgnoreNA sets whether to ignore NA values
func (ewm *EWM) IgnoreNA(ignore bool) *EWM {
	ewm.ignoreNA = ignore
	return ewm
}

// MinPeriods sets the minimum number of observations
func (ewm *EWM) MinPeriods(n int) *EWM {
	ewm.minPeriods = n
	return ewm
}

// Mean calculates exponentially weighted mean
func (ewm *EWM) Mean() Transform {
	return func(data []DataPoint) []DataPoint {
		if len(data) == 0 {
			return nil
		}

		result := make([]DataPoint, len(data))
		ewma := data[0].Y

		for i := range data {
			if i < ewm.minPeriods {
				result[i] = data[i]
				result[i].Y = math.NaN()
				result[i].Value = math.NaN()
				continue
			}

			if i == 0 {
				ewma = data[i].Y
			} else {
				ewma = ewm.alpha*data[i].Y + (1-ewm.alpha)*ewma
			}

			result[i] = data[i]
			result[i].Y = ewma
			result[i].Value = ewma
		}

		return result
	}
}

// Std calculates exponentially weighted standard deviation
func (ewm *EWM) Std() Transform {
	return func(data []DataPoint) []DataPoint {
		if len(data) == 0 {
			return nil
		}

		// First calculate EW mean
		meanData := ewm.Mean()(data)

		result := make([]DataPoint, len(data))
		ewmVar := 0.0

		for i := range data {
			if i < ewm.minPeriods {
				result[i] = data[i]
				result[i].Y = math.NaN()
				result[i].Value = math.NaN()
				continue
			}

			// Calculate squared deviation from EW mean
			sqDev := math.Pow(data[i].Y-meanData[i].Y, 2)

			if i == 0 {
				ewmVar = sqDev
			} else {
				ewmVar = ewm.alpha*sqDev + (1-ewm.alpha)*ewmVar
			}

			result[i] = data[i]
			result[i].Y = math.Sqrt(ewmVar)
			result[i].Value = result[i].Y
		}

		return result
	}
}

// Var calculates exponentially weighted variance
func (ewm *EWM) Var() Transform {
	return func(data []DataPoint) []DataPoint {
		stdData := ewm.Std()(data)
		result := make([]DataPoint, len(stdData))
		for i, d := range stdData {
			result[i] = d
			result[i].Y = d.Y * d.Y // variance = std^2
			result[i].Value = result[i].Y
		}
		return result
	}
}
