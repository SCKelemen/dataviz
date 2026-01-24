package charts

import (
	"github.com/SCKelemen/dataviz/scales"
	"github.com/SCKelemen/dataviz/transforms"
	"github.com/SCKelemen/svg"
	"github.com/SCKelemen/units"
)

// HistogramData represents data for a histogram
type HistogramData struct {
	Values []float64 // Raw data values
	Color  string    // Bar fill color
	Label  string    // Optional label
}

// HistogramSpec configures histogram rendering
type HistogramSpec struct {
	Data       *HistogramData
	Width      float64
	Height     float64
	BinCount   int       // Number of bins (0 = auto)
	BinSize    float64   // Fixed bin size (0 = use BinCount)
	Nice       bool      // If true, use nice round bin edges
	ShowDensity bool     // If true, normalize to show density instead of counts
	BarGap     float64   // Gap between bars (pixels)

	// Axis configuration
	XAxisLabel string
	YAxisLabel string
}

// RenderHistogram renders a histogram chart
func RenderHistogram(spec HistogramSpec) string {
	if spec.Data == nil || len(spec.Data.Values) == 0 {
		return ""
	}

	// Convert values to DataPoints
	data := make([]transforms.DataPoint, len(spec.Data.Values))
	for i, v := range spec.Data.Values {
		data[i] = transforms.DataPoint{Y: v}
	}

	// Apply binning transform
	var binned []transforms.DataPoint
	if spec.BinSize > 0 {
		binned = transforms.BinCount(spec.BinSize)(data)
	} else {
		binCount := spec.BinCount
		if binCount == 0 {
			binCount = 10
		}
		opts := transforms.BinOptions{
			Count: binCount,
			Nice:  spec.Nice,
		}
		binned = transforms.Bin(opts)(data)
	}

	if len(binned) == 0 {
		return ""
	}

	// Find max count for scaling
	maxCount := 0
	for _, bin := range binned {
		if bin.Count > maxCount {
			maxCount = bin.Count
		}
	}

	if maxCount == 0 {
		return ""
	}

	// Create scales
	xMin := binned[0].Y0
	xMax := binned[len(binned)-1].Y1

	xScale := scales.NewLinearScale(
		[2]float64{xMin, xMax},
		[2]units.Length{units.Px(40), units.Px(spec.Width - 40)},
	)

	yMax := float64(maxCount)
	if spec.ShowDensity {
		// Normalize to density
		totalCount := float64(len(spec.Data.Values))
		binWidth := binned[0].Y1 - binned[0].Y0
		yMax = float64(maxCount) / (totalCount * binWidth)
	}

	yScale := scales.NewLinearScale(
		[2]float64{0, yMax},
		[2]units.Length{units.Px(spec.Height - 40), units.Px(40)},
	)

	// Render bars
	color := spec.Data.Color
	if color == "" {
		color = "#4285f4"
	}

	barGap := spec.BarGap
	if barGap == 0 {
		barGap = 1
	}

	barStyle := svg.Style{
		Fill:        color,
		Stroke:      "#fff",
		StrokeWidth: barGap,
		Opacity:     0.8,
	}

	var result string
	baseY := yScale.Apply(0).Value

	for _, bin := range binned {
		x0 := xScale.Apply(bin.Y0).Value
		x1 := xScale.Apply(bin.Y1).Value
		barWidth := x1 - x0

		value := float64(bin.Count)
		if spec.ShowDensity {
			totalCount := float64(len(spec.Data.Values))
			binWidth := bin.Y1 - bin.Y0
			value = float64(bin.Count) / (totalCount * binWidth)
		}

		barHeight := baseY - yScale.Apply(value).Value

		result += svg.Rect(x0, yScale.Apply(value).Value, barWidth, barHeight, barStyle) + "\n"
	}

	return result
}

// DensityPlotData represents data for a density plot
type DensityPlotData struct {
	Values    []float64 // Raw data values
	Color     string    // Line color
	FillColor string    // Optional fill color
	Label     string    // Optional label
	Bandwidth float64   // KDE bandwidth (0 = auto)
}

// DensityPlotSpec configures density plot rendering
type DensityPlotSpec struct {
	Data       []*DensityPlotData
	Width      float64
	Height     float64
	ShowFill   bool    // If true, fill area under curve
	Smooth     bool    // If true, use smooth curves
	LineWidth  float64 // Line width

	// Axis configuration
	XAxisLabel string
	YAxisLabel string
}

// RenderDensityPlot renders a density plot
func RenderDensityPlot(spec DensityPlotSpec) string {
	if len(spec.Data) == 0 {
		return ""
	}

	// Find global min/max
	globalMin := spec.Data[0].Values[0]
	globalMax := spec.Data[0].Values[0]
	maxDensity := 0.0

	// Calculate KDE for each dataset
	allDensities := make([][]DensityPoint, len(spec.Data))

	for i, data := range spec.Data {
		if len(data.Values) == 0 {
			continue
		}

		// Find min/max
		for _, v := range data.Values {
			if v < globalMin {
				globalMin = v
			}
			if v > globalMax {
				globalMax = v
			}
		}

		// Calculate KDE using violin plot function
		stats := CalculateViolinStats(data.Values, data.Bandwidth)
		allDensities[i] = stats.Density

		// Find max density
		for _, dp := range stats.Density {
			if dp.Density > maxDensity {
				maxDensity = dp.Density
			}
		}
	}

	// Create scales
	xScale := scales.NewLinearScale(
		[2]float64{globalMin, globalMax},
		[2]units.Length{units.Px(40), units.Px(spec.Width - 40)},
	)

	yScale := scales.NewLinearScale(
		[2]float64{0, maxDensity * 1.1}, // Add 10% headroom
		[2]units.Length{units.Px(spec.Height - 40), units.Px(40)},
	)

	lineWidth := spec.LineWidth
	if lineWidth == 0 {
		lineWidth = 2
	}

	var result string

	// Render each density curve
	for i, densities := range allDensities {
		if len(densities) == 0 {
			continue
		}

		data := spec.Data[i]

		// Build path
		var pathData string
		for j, dp := range densities {
			x := xScale.Apply(dp.Value).Value
			y := yScale.Apply(dp.Density).Value

			if j == 0 {
				pathData = formatFloat(x) + "," + formatFloat(y)
			} else {
				pathData += " " + formatFloat(x) + "," + formatFloat(y)
			}
		}

		// Line color
		lineColor := data.Color
		if lineColor == "" {
			lineColor = "#4285f4"
		}

		// Draw filled area if enabled
		if spec.ShowFill {
			// Add baseline points
			baseY := yScale.Apply(0).Value
			lastX := xScale.Apply(densities[len(densities)-1].Value).Value
			firstX := xScale.Apply(densities[0].Value).Value

			fillPath := "M " + formatFloat(firstX) + " " + formatFloat(baseY) + " "
			fillPath += "L " + pathData + " "
			fillPath += "L " + formatFloat(lastX) + " " + formatFloat(baseY) + " Z"

			fillColor := data.FillColor
			if fillColor == "" {
				fillColor = lineColor
			}

			fillStyle := svg.Style{
				Fill:    fillColor,
				Opacity: 0.3,
				Stroke:  "none",
			}

			result += svg.Path(fillPath, fillStyle) + "\n"
		}

		// Draw line
		linePath := "M " + pathData
		lineStyle := svg.Style{
			Stroke:      lineColor,
			StrokeWidth: lineWidth,
			Fill:        "none",
		}

		result += svg.Path(linePath, lineStyle) + "\n"
	}

	return result
}
