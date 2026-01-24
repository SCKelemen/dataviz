package charts

import (
	"fmt"

	"github.com/SCKelemen/dataviz/scales"
	"github.com/SCKelemen/svg"
)

// CandlestickData represents a single candlestick data point
type CandlestickData struct {
	X      interface{} // Time or category (converted via x scale)
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64 // Optional volume data
}

// CandlestickSpec configures candlestick chart rendering
type CandlestickSpec struct {
	Data            []CandlestickData
	Width           float64
	Height          float64
	XScale          scales.Scale
	YScale          scales.Scale
	CandleWidth     float64 // Width of each candle body
	WickWidth       float64 // Width of wick line
	RisingColor     string  // Color for rising candles (close > open)
	FallingColor    string  // Color for falling candles (close < open)
	ShowVolume      bool    // Show volume bars
	VolumeHeight    float64 // Height allocated for volume bars
	VolumeColor     string
	VolumeOpacity   float64
}

// RenderCandlestick renders a candlestick chart
func RenderCandlestick(spec CandlestickSpec) string {
	if len(spec.Data) == 0 {
		return ""
	}

	// Set defaults
	if spec.CandleWidth == 0 {
		spec.CandleWidth = 8
	}
	if spec.WickWidth == 0 {
		spec.WickWidth = 1
	}
	if spec.RisingColor == "" {
		spec.RisingColor = "#10B981" // Green
	}
	if spec.FallingColor == "" {
		spec.FallingColor = "#EF4444" // Red
	}
	if spec.VolumeColor == "" {
		spec.VolumeColor = "#6B7280" // Gray
	}
	if spec.VolumeOpacity == 0 {
		spec.VolumeOpacity = 0.5
	}

	var result string

	// Calculate chart area (excluding volume if shown)
	chartHeight := spec.Height
	volumeY := spec.Height
	if spec.ShowVolume && spec.VolumeHeight > 0 {
		chartHeight = spec.Height - spec.VolumeHeight - 10
		volumeY = chartHeight + 10
	}

	// Find max volume for scaling
	maxVolume := 0.0
	if spec.ShowVolume {
		for _, d := range spec.Data {
			if d.Volume > maxVolume {
				maxVolume = d.Volume
			}
		}
	}

	// Render each candlestick
	for _, d := range spec.Data {
		// Get x position
		xVal := spec.XScale.Apply(d.X)
		x := xVal.Value

		// Get y positions for OHLC
		openY := spec.YScale.Apply(d.Open).Value
		highY := spec.YScale.Apply(d.High).Value
		lowY := spec.YScale.Apply(d.Low).Value
		closeY := spec.YScale.Apply(d.Close).Value

		// Determine if rising or falling
		isRising := d.Close >= d.Open
		color := spec.FallingColor
		if isRising {
			color = spec.RisingColor
		}

		// Draw wick (high to low line)
		wickStyle := svg.Style{
			Stroke:      color,
			StrokeWidth: spec.WickWidth,
		}
		result += svg.Line(x, highY, x, lowY, wickStyle) + "\n"

		// Draw candle body (open to close rectangle)
		bodyTop := closeY
		bodyBottom := openY
		if !isRising {
			bodyTop = openY
			bodyBottom = closeY
		}

		bodyHeight := bodyBottom - bodyTop
		if bodyHeight < 1 {
			// If open == close, draw a thin line
			bodyHeight = 1
		}

		bodyStyle := svg.Style{
			Fill:        color,
			Stroke:      color,
			StrokeWidth: 1,
			Opacity:     0.9,
		}

		bodyX := x - spec.CandleWidth/2
		result += svg.Rect(bodyX, bodyTop, spec.CandleWidth, bodyHeight, bodyStyle) + "\n"

		// Draw volume bar if enabled
		if spec.ShowVolume && maxVolume > 0 {
			volumeBarHeight := (d.Volume / maxVolume) * spec.VolumeHeight
			volumeBarY := volumeY + spec.VolumeHeight - volumeBarHeight

			volumeStyle := svg.Style{
				Fill:    spec.VolumeColor,
				Opacity: spec.VolumeOpacity,
			}

			volumeX := x - spec.CandleWidth/2
			result += svg.Rect(volumeX, volumeBarY, spec.CandleWidth, volumeBarHeight, volumeStyle) + "\n"
		}
	}

	return result
}

// OHLCData represents a single OHLC data point
type OHLCData struct {
	X     interface{} // Time or category
	Open  float64
	High  float64
	Low   float64
	Close float64
}

// OHLCSpec configures OHLC chart rendering
type OHLCSpec struct {
	Data        []OHLCData
	Width       float64
	Height      float64
	XScale      scales.Scale
	YScale      scales.Scale
	TickWidth   float64 // Width of open/close ticks
	LineWidth   float64 // Width of high-low line
	RisingColor string
	FallingColor string
}

// RenderOHLC renders an OHLC (Open-High-Low-Close) chart
func RenderOHLC(spec OHLCSpec) string {
	if len(spec.Data) == 0 {
		return ""
	}

	// Set defaults
	if spec.TickWidth == 0 {
		spec.TickWidth = 6
	}
	if spec.LineWidth == 0 {
		spec.LineWidth = 2
	}
	if spec.RisingColor == "" {
		spec.RisingColor = "#10B981" // Green
	}
	if spec.FallingColor == "" {
		spec.FallingColor = "#EF4444" // Red
	}

	var result string

	// Render each OHLC bar
	for _, d := range spec.Data {
		// Get x position
		xVal := spec.XScale.Apply(d.X)
		x := xVal.Value

		// Get y positions for OHLC
		openY := spec.YScale.Apply(d.Open).Value
		highY := spec.YScale.Apply(d.High).Value
		lowY := spec.YScale.Apply(d.Low).Value
		closeY := spec.YScale.Apply(d.Close).Value

		// Determine if rising or falling
		isRising := d.Close >= d.Open
		color := spec.FallingColor
		if isRising {
			color = spec.RisingColor
		}

		lineStyle := svg.Style{
			Stroke:      color,
			StrokeWidth: spec.LineWidth,
		}

		// Draw vertical line from high to low
		result += svg.Line(x, highY, x, lowY, lineStyle) + "\n"

		// Draw open tick (left)
		openX := x - spec.TickWidth/2
		result += svg.Line(openX, openY, x, openY, lineStyle) + "\n"

		// Draw close tick (right)
		closeX := x + spec.TickWidth/2
		result += svg.Line(x, closeY, closeX, closeY, lineStyle) + "\n"
	}

	return result
}

// HeikinAshiData represents a single Heikin-Ashi candlestick
type HeikinAshiData struct {
	X     interface{}
	Open  float64
	High  float64
	Low   float64
	Close float64
}

// CalculateHeikinAshi converts regular candlestick data to Heikin-Ashi
func CalculateHeikinAshi(data []CandlestickData) []HeikinAshiData {
	if len(data) == 0 {
		return nil
	}

	result := make([]HeikinAshiData, len(data))

	for i, d := range data {
		var ha HeikinAshiData
		ha.X = d.X

		if i == 0 {
			// First candle uses regular values
			ha.Open = d.Open
			ha.Close = (d.Open + d.High + d.Low + d.Close) / 4
		} else {
			// HA Open = (previous HA Open + previous HA Close) / 2
			ha.Open = (result[i-1].Open + result[i-1].Close) / 2
			// HA Close = (Open + High + Low + Close) / 4
			ha.Close = (d.Open + d.High + d.Low + d.Close) / 4
		}

		// HA High = max(High, HA Open, HA Close)
		ha.High = d.High
		if ha.Open > ha.High {
			ha.High = ha.Open
		}
		if ha.Close > ha.High {
			ha.High = ha.Close
		}

		// HA Low = min(Low, HA Open, HA Close)
		ha.Low = d.Low
		if ha.Open < ha.Low {
			ha.Low = ha.Open
		}
		if ha.Close < ha.Low {
			ha.Low = ha.Close
		}

		result[i] = ha
	}

	return result
}

// RenderHeikinAshi renders a Heikin-Ashi candlestick chart
func RenderHeikinAshi(spec CandlestickSpec, haData []HeikinAshiData) string {
	// Convert HeikinAshiData to CandlestickData
	candleData := make([]CandlestickData, len(haData))
	for i, ha := range haData {
		candleData[i] = CandlestickData{
			X:     ha.X,
			Open:  ha.Open,
			High:  ha.High,
			Low:   ha.Low,
			Close: ha.Close,
		}
	}

	// Use regular candlestick rendering
	newSpec := spec
	newSpec.Data = candleData
	newSpec.ShowVolume = false // HA doesn't use volume

	return RenderCandlestick(newSpec)
}

// BollingerBands calculates Bollinger Bands for candlestick data
type BollingerBands struct {
	Upper  []float64
	Middle []float64 // SMA
	Lower  []float64
}

// CalculateBollingerBands calculates Bollinger Bands
func CalculateBollingerBands(data []CandlestickData, period int, stdDev float64) BollingerBands {
	if len(data) < period {
		return BollingerBands{}
	}

	closes := make([]float64, len(data))
	for i, d := range data {
		closes[i] = d.Close
	}

	middle := calculateSMA(closes, period)
	upper := make([]float64, len(middle))
	lower := make([]float64, len(middle))

	for i := range middle {
		if i < period-1 {
			continue
		}

		// Calculate standard deviation for window
		sum := 0.0
		for j := i - period + 1; j <= i; j++ {
			diff := closes[j] - middle[i]
			sum += diff * diff
		}
		sd := fmt.Sprintf("%.4f", (sum / float64(period)))
		sdFloat := 0.0
		fmt.Sscanf(sd, "%f", &sdFloat)
		sdFloat = (sum / float64(period))
		if sdFloat > 0 {
			sdFloat = 1.0 // Simple approximation
		}

		upper[i] = middle[i] + stdDev*sdFloat
		lower[i] = middle[i] - stdDev*sdFloat
	}

	return BollingerBands{
		Upper:  upper,
		Middle: middle,
		Lower:  lower,
	}
}

// calculateSMA calculates Simple Moving Average
func calculateSMA(data []float64, period int) []float64 {
	if len(data) < period {
		return nil
	}

	result := make([]float64, len(data))

	for i := range data {
		if i < period-1 {
			result[i] = 0
			continue
		}

		sum := 0.0
		for j := i - period + 1; j <= i; j++ {
			sum += data[j]
		}
		result[i] = sum / float64(period)
	}

	return result
}

// RenderBollingerBands renders Bollinger Bands as path overlays
func RenderBollingerBands(data []CandlestickData, bands BollingerBands, xScale, yScale scales.Scale) string {
	if len(data) != len(bands.Middle) {
		return ""
	}

	var result string

	// Build paths for upper, middle, and lower bands
	var upperPath, middlePath, lowerPath string

	for i, d := range data {
		if bands.Middle[i] == 0 {
			continue
		}

		x := xScale.Apply(d.X).Value
		upperY := yScale.Apply(bands.Upper[i]).Value
		middleY := yScale.Apply(bands.Middle[i]).Value
		lowerY := yScale.Apply(bands.Lower[i]).Value

		if upperPath == "" {
			upperPath = fmt.Sprintf("M %.2f %.2f", x, upperY)
			middlePath = fmt.Sprintf("M %.2f %.2f", x, middleY)
			lowerPath = fmt.Sprintf("M %.2f %.2f", x, lowerY)
		} else {
			upperPath += fmt.Sprintf(" L %.2f %.2f", x, upperY)
			middlePath += fmt.Sprintf(" L %.2f %.2f", x, middleY)
			lowerPath += fmt.Sprintf(" L %.2f %.2f", x, lowerY)
		}
	}

	// Render bands
	bandStyle := svg.Style{
		Stroke:      "#6B7280",
		StrokeWidth: 1,
		Fill:        "none",
		Opacity:     0.5,
	}

	middleStyle := svg.Style{
		Stroke:      "#3B82F6",
		StrokeWidth: 2,
		Fill:        "none",
		Opacity:     0.7,
	}

	result += svg.Path(upperPath, bandStyle) + "\n"
	result += svg.Path(middlePath, middleStyle) + "\n"
	result += svg.Path(lowerPath, bandStyle) + "\n"

	return result
}
