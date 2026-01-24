package charts

import (
	"fmt"
	"math"
	"sort"

	"github.com/SCKelemen/dataviz/scales"
	"github.com/SCKelemen/svg"
	"github.com/SCKelemen/units"
)

// ViolinPlotData represents data for a violin plot
type ViolinPlotData struct {
	Values      []float64 // Raw data values
	Label       string    // Category label
	Color       string    // Fill color
	StrokeColor string    // Outline color
}

// ViolinPlotSpec configures violin plot rendering
type ViolinPlotSpec struct {
	Data      []*ViolinPlotData
	Width     float64
	Height    float64
	Bandwidth float64 // KDE bandwidth (0 = auto)
	ShowBox   bool    // If true, show box plot inside violin
	ShowMedian bool   // If true, show median line
	ShowMean  bool    // If true, show mean marker
	ViolinWidth float64 // Maximum width of violin (0 = auto)

	// Axis configuration
	XAxisLabel string
	YAxisLabel string
}

// ViolinStats represents statistics for a violin plot
type ViolinStats struct {
	Density []DensityPoint // KDE density points
	Mean    float64
	Median  float64
	Q1      float64
	Q3      float64
}

// DensityPoint represents a point in the density estimation
type DensityPoint struct {
	Value   float64 // Y value
	Density float64 // Density at this value
}

// CalculateViolinStats calculates statistics and KDE for violin plot
func CalculateViolinStats(values []float64, bandwidth float64) ViolinStats {
	if len(values) == 0 {
		return ViolinStats{}
	}

	// Calculate basic statistics
	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)

	mean := 0.0
	for _, v := range values {
		mean += v
	}
	mean /= float64(len(values))

	median := percentile(sorted, 50)
	q1 := percentile(sorted, 25)
	q3 := percentile(sorted, 75)

	// Calculate KDE
	density := calculateKDE(values, bandwidth)

	return ViolinStats{
		Density: density,
		Mean:    mean,
		Median:  median,
		Q1:      q1,
		Q3:      q3,
	}
}

// calculateKDE performs kernel density estimation using Gaussian kernel
func calculateKDE(values []float64, bandwidth float64) []DensityPoint {
	if len(values) == 0 {
		return nil
	}

	// Find data range
	minVal := values[0]
	maxVal := values[0]
	for _, v := range values {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
	}

	// Auto-calculate bandwidth using Silverman's rule of thumb
	if bandwidth == 0 {
		// Calculate standard deviation
		mean := 0.0
		for _, v := range values {
			mean += v
		}
		mean /= float64(len(values))

		variance := 0.0
		for _, v := range values {
			diff := v - mean
			variance += diff * diff
		}
		variance /= float64(len(values))
		stdDev := math.Sqrt(variance)

		// Silverman's rule: h = 0.9 * min(σ, IQR/1.34) * n^(-1/5)
		sorted := make([]float64, len(values))
		copy(sorted, values)
		sort.Float64s(sorted)
		iqr := percentile(sorted, 75) - percentile(sorted, 25)
		bandwidth = 0.9 * math.Min(stdDev, iqr/1.34) * math.Pow(float64(len(values)), -0.2)
	}

	// Generate evaluation points
	numPoints := 100
	step := (maxVal - minVal) / float64(numPoints-1)
	density := make([]DensityPoint, numPoints)

	// Calculate density at each point using Gaussian kernel
	for i := 0; i < numPoints; i++ {
		x := minVal + float64(i)*step
		d := 0.0

		for _, value := range values {
			// Gaussian kernel: (1/√(2π)) * exp(-0.5 * ((x-value)/h)²)
			z := (x - value) / bandwidth
			d += math.Exp(-0.5*z*z) / (bandwidth * math.Sqrt(2*math.Pi))
		}
		d /= float64(len(values))

		density[i] = DensityPoint{
			Value:   x,
			Density: d,
		}
	}

	return density
}

// RenderViolinPlot renders a violin plot
func RenderViolinPlot(spec ViolinPlotSpec) string {
	if len(spec.Data) == 0 {
		return ""
	}

	// Calculate statistics for each violin
	stats := make([]ViolinStats, len(spec.Data))
	globalMin := math.Inf(1)
	globalMax := math.Inf(-1)
	maxDensity := 0.0

	for i, data := range spec.Data {
		stats[i] = CalculateViolinStats(data.Values, spec.Bandwidth)

		// Find global min/max
		for _, dp := range stats[i].Density {
			if dp.Value < globalMin {
				globalMin = dp.Value
			}
			if dp.Value > globalMax {
				globalMax = dp.Value
			}
			if dp.Density > maxDensity {
				maxDensity = dp.Density
			}
		}
	}

	// Create scales
	yScale := scales.NewLinearScale(
		[2]float64{globalMin, globalMax},
		[2]units.Length{units.Px(spec.Height - 40), units.Px(40)},
	)

	// Calculate violin width
	violinWidth := spec.ViolinWidth
	if violinWidth == 0 {
		violinWidth = (spec.Width - 80) / float64(len(spec.Data)) * 0.4
	}

	// Calculate x positions
	spacing := (spec.Width - 80) / float64(len(spec.Data))

	var result string

	// Render each violin
	for i, st := range stats {
		data := spec.Data[i]
		centerX := 40 + spacing*float64(i) + spacing/2

		// Build violin path (mirrored density)
		var leftPath, rightPath string

		for j, dp := range st.Density {
			y := yScale.Apply(dp.Value).Value
			// Normalize density to violin width
			width := (dp.Density / maxDensity) * (violinWidth / 2)

			if j == 0 {
				leftPath = fmt.Sprintf("M %.2f %.2f", centerX-width, y)
				rightPath = fmt.Sprintf("M %.2f %.2f", centerX+width, y)
			} else {
				leftPath += fmt.Sprintf(" L %.2f %.2f", centerX-width, y)
				rightPath += fmt.Sprintf(" L %.2f %.2f", centerX+width, y)
			}
		}

		// Combine paths: left side, then right side in reverse
		fullPath := leftPath
		for j := len(st.Density) - 1; j >= 0; j-- {
			dp := st.Density[j]
			y := yScale.Apply(dp.Value).Value
			width := (dp.Density / maxDensity) * (violinWidth / 2)
			fullPath += fmt.Sprintf(" L %.2f %.2f", centerX+width, y)
		}
		fullPath += " Z"

		// Colors
		fillColor := data.Color
		if fillColor == "" {
			fillColor = "#4285f4"
		}
		strokeColor := data.StrokeColor
		if strokeColor == "" {
			strokeColor = "#333"
		}

		// Draw violin
		violinStyle := svg.Style{
			Fill:        fillColor,
			Stroke:      strokeColor,
			StrokeWidth: 1,
			Opacity:     0.6,
		}
		result += svg.Path(fullPath, violinStyle) + "\n"

		// Draw box plot inside if enabled
		if spec.ShowBox {
			q1Y := yScale.Apply(st.Q1).Value
			q3Y := yScale.Apply(st.Q3).Value
			medianY := yScale.Apply(st.Median).Value
			boxWidth := violinWidth * 0.15

			boxStyle := svg.Style{
				Fill:        "#fff",
				Stroke:      strokeColor,
				StrokeWidth: 1.5,
				Opacity:     0.8,
			}
			result += svg.Rect(centerX-boxWidth/2, q3Y, boxWidth, q1Y-q3Y, boxStyle) + "\n"

			// Draw median line
			medianStyle := svg.Style{
				Stroke:      strokeColor,
				StrokeWidth: 2,
			}
			result += svg.Line(centerX-boxWidth/2, medianY, centerX+boxWidth/2, medianY, medianStyle) + "\n"
		} else if spec.ShowMedian {
			// Just draw median line
			medianY := yScale.Apply(st.Median).Value
			medianStyle := svg.Style{
				Stroke:      strokeColor,
				StrokeWidth: 2,
			}
			result += svg.Line(centerX-violinWidth/4, medianY, centerX+violinWidth/4, medianY, medianStyle) + "\n"
		}

		// Draw mean marker if enabled
		if spec.ShowMean {
			meanY := yScale.Apply(st.Mean).Value
			meanStyle := svg.Style{
				Fill:   strokeColor,
				Stroke: strokeColor,
			}
			result += svg.Circle(centerX, meanY, 3, meanStyle) + "\n"
		}

		// Draw label
		if data.Label != "" {
			labelStyle := svg.Style{
				FontSize:   units.Px(12),
				FontFamily: "sans-serif",
				TextAnchor: svg.TextAnchorMiddle,
			}
			result += svg.Text(data.Label, centerX, spec.Height-20, labelStyle) + "\n"
		}
	}

	return result
}
