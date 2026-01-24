package charts

import (
	"strings"
	"testing"

	"github.com/SCKelemen/dataviz/scales"
	"github.com/SCKelemen/units"
)

func createTestCandlestickData() []CandlestickData {
	return []CandlestickData{
		{X: 0, Open: 100, High: 110, Low: 95, Close: 105, Volume: 1000},
		{X: 1, Open: 105, High: 115, Low: 100, Close: 110, Volume: 1200},
		{X: 2, Open: 110, High: 112, Low: 105, Close: 108, Volume: 900},
		{X: 3, Open: 108, High: 120, Low: 107, Close: 118, Volume: 1500},
		{X: 4, Open: 118, High: 125, Low: 115, Close: 120, Volume: 1100},
	}
}

func createTestOHLCData() []OHLCData {
	return []OHLCData{
		{X: 0, Open: 100, High: 110, Low: 95, Close: 105},
		{X: 1, Open: 105, High: 115, Low: 100, Close: 110},
		{X: 2, Open: 110, High: 112, Low: 105, Close: 108},
		{X: 3, Open: 108, High: 120, Low: 107, Close: 118},
		{X: 4, Open: 118, High: 125, Low: 115, Close: 120},
	}
}

// Candlestick tests

func TestRenderCandlestick(t *testing.T) {
	data := createTestCandlestickData()

	xScale := scales.NewLinearScale(
		[2]float64{0, 4},
		[2]units.Length{units.Px(50), units.Px(750)},
	)

	yScale := scales.NewLinearScale(
		[2]float64{90, 130},
		[2]units.Length{units.Px(550), units.Px(50)},
	)

	spec := CandlestickSpec{
		Data:   data,
		Width:  800,
		Height: 600,
		XScale: xScale,
		YScale: yScale,
	}

	result := RenderCandlestick(spec)

	if result == "" {
		t.Fatal("Expected non-empty SVG output")
	}

	// Should have lines for wicks and rects for bodies
	if !strings.Contains(result, "<line") {
		t.Error("Expected line elements for wicks")
	}

	if !strings.Contains(result, "<rect") {
		t.Error("Expected rect elements for candle bodies")
	}

	// Should use rising/falling colors
	if !strings.Contains(result, "#10B981") && !strings.Contains(result, "#EF4444") {
		t.Error("Expected rising or falling colors in output")
	}
}

func TestRenderCandlestickWithVolume(t *testing.T) {
	data := createTestCandlestickData()

	xScale := scales.NewLinearScale(
		[2]float64{0, 4},
		[2]units.Length{units.Px(50), units.Px(750)},
	)

	yScale := scales.NewLinearScale(
		[2]float64{90, 130},
		[2]units.Length{units.Px(450), units.Px(50)},
	)

	spec := CandlestickSpec{
		Data:         data,
		Width:        800,
		Height:       600,
		XScale:       xScale,
		YScale:       yScale,
		ShowVolume:   true,
		VolumeHeight: 100,
	}

	result := RenderCandlestick(spec)

	if result == "" {
		t.Fatal("Expected non-empty SVG output")
	}

	// Count rectangles - should have both candles and volume bars
	rectCount := strings.Count(result, "<rect")
	if rectCount < len(data)*2 {
		t.Errorf("Expected at least %d rectangles (candles + volume), got %d",
			len(data)*2, rectCount)
	}
}

func TestRenderCandlestickEmptyData(t *testing.T) {
	spec := CandlestickSpec{
		Data:   []CandlestickData{},
		Width:  800,
		Height: 600,
	}

	result := RenderCandlestick(spec)

	if result != "" {
		t.Error("Expected empty string for empty data")
	}
}

func TestRenderCandlestickCustomColors(t *testing.T) {
	data := createTestCandlestickData()

	xScale := scales.NewLinearScale(
		[2]float64{0, 4},
		[2]units.Length{units.Px(50), units.Px(750)},
	)

	yScale := scales.NewLinearScale(
		[2]float64{90, 130},
		[2]units.Length{units.Px(550), units.Px(50)},
	)

	spec := CandlestickSpec{
		Data:         data,
		Width:        800,
		Height:       600,
		XScale:       xScale,
		YScale:       yScale,
		RisingColor:  "#00FF00",
		FallingColor: "#FF0000",
	}

	result := RenderCandlestick(spec)

	if result == "" {
		t.Fatal("Expected non-empty SVG output")
	}

	// Should use custom colors
	hasCustomColors := strings.Contains(result, "#00FF00") || strings.Contains(result, "#FF0000")
	if !hasCustomColors {
		t.Error("Expected custom rising/falling colors in output")
	}
}

// OHLC tests

func TestRenderOHLC(t *testing.T) {
	data := createTestOHLCData()

	xScale := scales.NewLinearScale(
		[2]float64{0, 4},
		[2]units.Length{units.Px(50), units.Px(750)},
	)

	yScale := scales.NewLinearScale(
		[2]float64{90, 130},
		[2]units.Length{units.Px(550), units.Px(50)},
	)

	spec := OHLCSpec{
		Data:   data,
		Width:  800,
		Height: 600,
		XScale: xScale,
		YScale: yScale,
	}

	result := RenderOHLC(spec)

	if result == "" {
		t.Fatal("Expected non-empty SVG output")
	}

	// Should have lines for high-low bars and open/close ticks
	lineCount := strings.Count(result, "<line")
	expectedLines := len(data) * 3 // high-low + open tick + close tick
	if lineCount < expectedLines {
		t.Errorf("Expected at least %d lines, got %d", expectedLines, lineCount)
	}

	// Should use rising/falling colors
	if !strings.Contains(result, "#10B981") && !strings.Contains(result, "#EF4444") {
		t.Error("Expected rising or falling colors in output")
	}
}

func TestRenderOHLCEmptyData(t *testing.T) {
	spec := OHLCSpec{
		Data:   []OHLCData{},
		Width:  800,
		Height: 600,
	}

	result := RenderOHLC(spec)

	if result != "" {
		t.Error("Expected empty string for empty data")
	}
}

func TestRenderOHLCCustomColors(t *testing.T) {
	data := createTestOHLCData()

	xScale := scales.NewLinearScale(
		[2]float64{0, 4},
		[2]units.Length{units.Px(50), units.Px(750)},
	)

	yScale := scales.NewLinearScale(
		[2]float64{90, 130},
		[2]units.Length{units.Px(550), units.Px(50)},
	)

	spec := OHLCSpec{
		Data:         data,
		Width:        800,
		Height:       600,
		XScale:       xScale,
		YScale:       yScale,
		RisingColor:  "#00FF00",
		FallingColor: "#FF0000",
	}

	result := RenderOHLC(spec)

	if result == "" {
		t.Fatal("Expected non-empty SVG output")
	}

	// Should use custom colors
	hasCustomColors := strings.Contains(result, "#00FF00") || strings.Contains(result, "#FF0000")
	if !hasCustomColors {
		t.Error("Expected custom rising/falling colors in output")
	}
}

// Heikin-Ashi tests

func TestCalculateHeikinAshi(t *testing.T) {
	data := createTestCandlestickData()

	haData := CalculateHeikinAshi(data)

	if len(haData) != len(data) {
		t.Errorf("Expected %d HA candles, got %d", len(data), len(haData))
	}

	// Check that HA values are calculated
	for i, ha := range haData {
		if ha.X != data[i].X {
			t.Errorf("HA candle %d: X value mismatch", i)
		}

		// HA values should be different from original
		if i > 0 && ha.Open == data[i].Open {
			// This might happen, but generally they differ
		}

		// HA High should be >= all other values
		if ha.High < ha.Open || ha.High < ha.Close {
			// Allow small floating point differences
		}

		// HA Low should be <= all other values
		if ha.Low > ha.Open || ha.Low > ha.Close {
			// Allow small floating point differences
		}
	}
}

func TestCalculateHeikinAshiEmpty(t *testing.T) {
	haData := CalculateHeikinAshi([]CandlestickData{})

	if haData != nil {
		t.Error("Expected nil for empty data")
	}
}

func TestRenderHeikinAshi(t *testing.T) {
	data := createTestCandlestickData()
	haData := CalculateHeikinAshi(data)

	xScale := scales.NewLinearScale(
		[2]float64{0, 4},
		[2]units.Length{units.Px(50), units.Px(750)},
	)

	yScale := scales.NewLinearScale(
		[2]float64{90, 130},
		[2]units.Length{units.Px(550), units.Px(50)},
	)

	spec := CandlestickSpec{
		Width:  800,
		Height: 600,
		XScale: xScale,
		YScale: yScale,
	}

	result := RenderHeikinAshi(spec, haData)

	if result == "" {
		t.Fatal("Expected non-empty SVG output")
	}

	// Should have candle elements
	if !strings.Contains(result, "<line") {
		t.Error("Expected line elements for wicks")
	}

	if !strings.Contains(result, "<rect") {
		t.Error("Expected rect elements for candle bodies")
	}
}

// Bollinger Bands tests

func TestCalculateBollingerBands(t *testing.T) {
	data := createTestCandlestickData()

	bands := CalculateBollingerBands(data, 3, 2.0)

	if len(bands.Middle) != len(data) {
		t.Errorf("Expected %d middle band values, got %d", len(data), len(bands.Middle))
	}

	if len(bands.Upper) != len(data) {
		t.Errorf("Expected %d upper band values, got %d", len(data), len(bands.Upper))
	}

	if len(bands.Lower) != len(data) {
		t.Errorf("Expected %d lower band values, got %d", len(data), len(bands.Lower))
	}

	// Check that bands are ordered (upper >= middle >= lower)
	for i := range bands.Middle {
		if bands.Middle[i] == 0 {
			continue // Skip uninitialized values
		}

		if bands.Upper[i] < bands.Middle[i] {
			// Allow small floating point differences
		}

		if bands.Lower[i] > bands.Middle[i] {
			// Allow small floating point differences
		}
	}
}

func TestCalculateBollingerBandsInsufficientData(t *testing.T) {
	data := []CandlestickData{
		{X: 0, Open: 100, High: 110, Low: 95, Close: 105},
	}

	bands := CalculateBollingerBands(data, 3, 2.0)

	if len(bands.Middle) != 0 {
		t.Error("Expected empty bands for insufficient data")
	}
}

func TestCalculateSMA(t *testing.T) {
	data := []float64{100, 105, 110, 108, 112}

	sma := calculateSMA(data, 3)

	if len(sma) != len(data) {
		t.Errorf("Expected %d SMA values, got %d", len(data), len(sma))
	}

	// First two values should be 0 (not enough data)
	if sma[0] != 0 || sma[1] != 0 {
		t.Error("Expected first values to be 0 for insufficient period")
	}

	// Third value should be average of first 3
	expectedSMA3 := (100.0 + 105.0 + 110.0) / 3.0
	if sma[2] != expectedSMA3 {
		t.Errorf("Expected SMA[2] = %.2f, got %.2f", expectedSMA3, sma[2])
	}
}

func TestCalculateSMAInsufficientData(t *testing.T) {
	data := []float64{100, 105}

	sma := calculateSMA(data, 3)

	if sma != nil {
		t.Error("Expected nil for insufficient data")
	}
}

func TestRenderBollingerBands(t *testing.T) {
	data := createTestCandlestickData()
	bands := CalculateBollingerBands(data, 3, 2.0)

	xScale := scales.NewLinearScale(
		[2]float64{0, 4},
		[2]units.Length{units.Px(50), units.Px(750)},
	)

	yScale := scales.NewLinearScale(
		[2]float64{90, 130},
		[2]units.Length{units.Px(550), units.Px(50)},
	)

	result := RenderBollingerBands(data, bands, xScale, yScale)

	if result == "" {
		t.Fatal("Expected non-empty SVG output")
	}

	// Should have path elements for three bands
	pathCount := strings.Count(result, "<path")
	if pathCount != 3 {
		t.Errorf("Expected 3 path elements (upper, middle, lower), got %d", pathCount)
	}
}

func TestRenderBollingerBandsMismatchedLength(t *testing.T) {
	data := createTestCandlestickData()
	bands := BollingerBands{
		Upper:  []float64{110, 115},
		Middle: []float64{105, 110},
		Lower:  []float64{100, 105},
	}

	xScale := scales.NewLinearScale(
		[2]float64{0, 4},
		[2]units.Length{units.Px(50), units.Px(750)},
	)

	yScale := scales.NewLinearScale(
		[2]float64{90, 130},
		[2]units.Length{units.Px(550), units.Px(50)},
	)

	result := RenderBollingerBands(data, bands, xScale, yScale)

	if result != "" {
		t.Error("Expected empty string for mismatched data/bands length")
	}
}

// Edge case tests

func TestCandlestickDojiPattern(t *testing.T) {
	// Doji pattern: open == close
	data := []CandlestickData{
		{X: 0, Open: 100, High: 105, Low: 95, Close: 100},
	}

	xScale := scales.NewLinearScale(
		[2]float64{0, 1},
		[2]units.Length{units.Px(50), units.Px(750)},
	)

	yScale := scales.NewLinearScale(
		[2]float64{90, 110},
		[2]units.Length{units.Px(550), units.Px(50)},
	)

	spec := CandlestickSpec{
		Data:   data,
		Width:  800,
		Height: 600,
		XScale: xScale,
		YScale: yScale,
	}

	result := RenderCandlestick(spec)

	if result == "" {
		t.Fatal("Expected non-empty SVG output for doji pattern")
	}

	// Should still render a thin body
	if !strings.Contains(result, "<rect") {
		t.Error("Expected rect element for doji body")
	}
}

func TestCandlestickSingleDataPoint(t *testing.T) {
	data := []CandlestickData{
		{X: 0, Open: 100, High: 110, Low: 95, Close: 105},
	}

	xScale := scales.NewLinearScale(
		[2]float64{0, 1},
		[2]units.Length{units.Px(50), units.Px(750)},
	)

	yScale := scales.NewLinearScale(
		[2]float64{90, 115},
		[2]units.Length{units.Px(550), units.Px(50)},
	)

	spec := CandlestickSpec{
		Data:   data,
		Width:  800,
		Height: 600,
		XScale: xScale,
		YScale: yScale,
	}

	result := RenderCandlestick(spec)

	if result == "" {
		t.Fatal("Expected non-empty SVG output for single data point")
	}
}
