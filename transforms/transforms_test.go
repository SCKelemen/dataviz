package transforms

import (
	"math"
	"testing"
	"time"
)

// Helper function to compare floats with tolerance
func floatEquals(a, b, tolerance float64) bool {
	return math.Abs(a-b) < tolerance
}

// ==================== Binning Tests ====================

func TestBin_Basic(t *testing.T) {
	data := []DataPoint{
		{Y: 1.5}, {Y: 2.3}, {Y: 5.7}, {Y: 8.1}, {Y: 9.2},
	}

	bins := Bin(BinOptions{Count: 3})(data)

	if len(bins) != 3 {
		t.Errorf("Expected 3 bins, got %d", len(bins))
	}

	// Check that counts sum to total
	totalCount := 0
	for _, bin := range bins {
		totalCount += bin.Count
	}
	if totalCount != len(data) {
		t.Errorf("Expected total count %d, got %d", len(data), totalCount)
	}
}

func TestHistogram(t *testing.T) {
	data := []DataPoint{
		{Y: 1}, {Y: 2}, {Y: 2}, {Y: 3}, {Y: 3}, {Y: 3}, {Y: 4},
	}

	bins := Histogram()(data)

	if len(bins) == 0 {
		t.Fatal("Histogram returned no bins")
	}

	// All original data should be counted
	totalCount := 0
	for _, bin := range bins {
		totalCount += bin.Count
	}
	if totalCount != len(data) {
		t.Errorf("Expected count %d, got %d", len(data), totalCount)
	}
}

func TestBinCount(t *testing.T) {
	data := []DataPoint{
		{Y: 0}, {Y: 1.5}, {Y: 3.2}, {Y: 5.8}, {Y: 7.1},
	}

	bins := BinCount(2.0)(data)

	if len(bins) == 0 {
		t.Fatal("BinCount returned no bins")
	}

	for _, bin := range bins {
		if bin.Y1-bin.Y0 != 2.0 {
			t.Errorf("Expected bin width 2.0, got %f", bin.Y1-bin.Y0)
		}
	}
}

// ==================== Grouping Tests ====================

func TestGroupBy_Sum(t *testing.T) {
	data := []DataPoint{
		{Label: "A", Y: 10},
		{Label: "A", Y: 20},
		{Label: "B", Y: 15},
		{Label: "B", Y: 25},
	}

	grouped := GroupBy(GroupOptions{
		By:        "Label",
		Aggregate: Sum,
	})(data)

	if len(grouped) != 2 {
		t.Fatalf("Expected 2 groups, got %d", len(grouped))
	}

	// Find group A and check sum
	for _, g := range grouped {
		if g.Label == "A" && g.Y != 30 {
			t.Errorf("Group A: expected Y=30, got %f", g.Y)
		}
		if g.Label == "B" && g.Y != 40 {
			t.Errorf("Group B: expected Y=40, got %f", g.Y)
		}
	}
}

func TestGroupBy_Mean(t *testing.T) {
	data := []DataPoint{
		{Label: "X", Y: 10},
		{Label: "X", Y: 20},
		{Label: "X", Y: 30},
	}

	grouped := GroupBy(GroupOptions{
		By:        "Label",
		Aggregate: Mean,
	})(data)

	if len(grouped) != 1 {
		t.Fatalf("Expected 1 group, got %d", len(grouped))
	}

	if !floatEquals(grouped[0].Y, 20, 0.01) {
		t.Errorf("Expected mean 20, got %f", grouped[0].Y)
	}
}

func TestFilter(t *testing.T) {
	data := []DataPoint{
		{Y: 10}, {Y: 20}, {Y: 30}, {Y: 40},
	}

	filtered := Filter(func(d DataPoint) bool {
		return d.Y > 20
	})(data)

	if len(filtered) != 2 {
		t.Errorf("Expected 2 points after filter, got %d", len(filtered))
	}
}

func TestTop(t *testing.T) {
	data := []DataPoint{
		{Y: 10}, {Y: 50}, {Y: 30}, {Y: 20},
	}

	top2 := Top(2)(data)

	if len(top2) != 2 {
		t.Fatalf("Expected 2 points, got %d", len(top2))
	}

	if top2[0].Y != 50 || top2[1].Y != 30 {
		t.Error("Top 2 should be 50 and 30")
	}
}

func TestCumulative(t *testing.T) {
	data := []DataPoint{
		{Y: 1}, {Y: 2}, {Y: 3}, {Y: 4},
	}

	cumulative := Cumulative()(data)

	expected := []float64{1, 3, 6, 10}
	for i, exp := range expected {
		if cumulative[i].Y != exp {
			t.Errorf("Index %d: expected %f, got %f", i, exp, cumulative[i].Y)
		}
	}
}

func TestWindow(t *testing.T) {
	data := []DataPoint{
		{Y: 1}, {Y: 2}, {Y: 3}, {Y: 4}, {Y: 5},
	}

	windowed := Window(3, Mean)(data)

	if len(windowed) != len(data) {
		t.Errorf("Expected %d points, got %d", len(data), len(windowed))
	}

	// At index 2, window is [1, 2, 3] (start=0, size=3), mean = 2
	if !floatEquals(windowed[2].Y, 2.0, 0.01) {
		t.Errorf("Expected window value 2.0, got %f", windowed[2].Y)
	}

	// At index 4 (last), window is [3, 4, 5], mean = 4
	if !floatEquals(windowed[4].Y, 4.0, 0.01) {
		t.Errorf("Expected last window value 4.0, got %f", windowed[4].Y)
	}
}

// ==================== Stacking Tests ====================

func TestStack_Basic(t *testing.T) {
	data := []DataPoint{
		{Label: "2020", Group: "A", Y: 10},
		{Label: "2020", Group: "B", Y: 15},
		{Label: "2021", Group: "A", Y: 12},
		{Label: "2021", Group: "B", Y: 18},
	}

	stacked := StackZero("Label")(data)

	// Check that stacking is correct
	// For 2020: A should have Y0=0, Y1=10; B should have Y0=10, Y1=25
	for _, d := range stacked {
		if d.Label == "2020" && d.Group == "A" {
			if d.Y0 != 0 || d.Y1 != 10 {
				t.Errorf("2020 A: expected Y0=0, Y1=10, got Y0=%f, Y1=%f", d.Y0, d.Y1)
			}
		}
		if d.Label == "2020" && d.Group == "B" {
			if d.Y0 != 10 || d.Y1 != 25 {
				t.Errorf("2020 B: expected Y0=10, Y1=25, got Y0=%f, Y1=%f", d.Y0, d.Y1)
			}
		}
	}
}

func TestStackNormalize(t *testing.T) {
	data := []DataPoint{
		{Label: "X", Y: 20},
		{Label: "X", Y: 30},
	}

	normalized := StackNormalize("Label")(data)

	// Total is 50, so values should be 20/50=0.4 and 30/50=0.6
	if !floatEquals(normalized[0].Y1, 0.4, 0.01) {
		t.Errorf("Expected Y1=0.4, got %f", normalized[0].Y1)
	}
	if !floatEquals(normalized[1].Y1, 1.0, 0.01) {
		t.Errorf("Expected Y1=1.0, got %f", normalized[1].Y1)
	}
}

// ==================== Smoothing Tests ====================

func TestMovingAverage(t *testing.T) {
	data := []DataPoint{
		{Y: 1}, {Y: 5}, {Y: 2}, {Y: 8}, {Y: 3},
	}

	smoothed := MovingAverage(3)(data)

	if len(smoothed) != len(data) {
		t.Errorf("Expected %d points, got %d", len(data), len(smoothed))
	}

	// Middle point should be average of [5, 2, 8] = 5
	if !floatEquals(smoothed[2].Y, 5.0, 0.01) {
		t.Errorf("Expected smoothed value 5.0, got %f", smoothed[2].Y)
	}
}

func TestExponentialSmoothing(t *testing.T) {
	data := []DataPoint{
		{Y: 10}, {Y: 20}, {Y: 30},
	}

	smoothed := ExponentialSmoothing(0.5)(data)

	if len(smoothed) != len(data) {
		t.Errorf("Expected %d points, got %d", len(data), len(smoothed))
	}

	// First value should be unchanged
	if smoothed[0].Y != 10 {
		t.Errorf("Expected first value 10, got %f", smoothed[0].Y)
	}

	// Second value: 0.5*20 + 0.5*10 = 15
	if !floatEquals(smoothed[1].Y, 15.0, 0.01) {
		t.Errorf("Expected second value 15.0, got %f", smoothed[1].Y)
	}
}

func TestDownsample(t *testing.T) {
	data := []DataPoint{
		{Y: 1}, {Y: 2}, {Y: 3}, {Y: 4}, {Y: 5}, {Y: 6},
	}

	downsampled := Downsample(2)(data)

	if len(downsampled) != 3 {
		t.Errorf("Expected 3 points, got %d", len(downsampled))
	}

	expected := []float64{1, 3, 5}
	for i, exp := range expected {
		if downsampled[i].Y != exp {
			t.Errorf("Index %d: expected %f, got %f", i, exp, downsampled[i].Y)
		}
	}
}

// ==================== Normalization Tests ====================

func TestNormalizePercentage(t *testing.T) {
	data := []DataPoint{
		{Y: 10}, {Y: 20}, {Y: 30}, {Y: 40},
	}

	normalized := NormalizePercentage()(data)

	// Total is 100, so percentages should be 10%, 20%, 30%, 40%
	expected := []float64{10, 20, 30, 40}
	for i, exp := range expected {
		if !floatEquals(normalized[i].Y, exp, 0.01) {
			t.Errorf("Index %d: expected %f%%, got %f%%", i, exp, normalized[i].Y)
		}
	}
}

func TestNormalizeZScore(t *testing.T) {
	data := []DataPoint{
		{Y: 10}, {Y: 20}, {Y: 30}, {Y: 40},
	}

	normalized := NormalizeZScore()(data)

	// Mean is 25, should have negative and positive z-scores
	hasNegative := false
	hasPositive := false
	for _, d := range normalized {
		if d.Y < 0 {
			hasNegative = true
		}
		if d.Y > 0 {
			hasPositive = true
		}
	}

	if !hasNegative || !hasPositive {
		t.Error("Z-scores should have both negative and positive values")
	}
}

func TestNormalizeMinMax(t *testing.T) {
	data := []DataPoint{
		{Y: 10}, {Y: 30}, {Y: 50},
	}

	normalized := NormalizeMinMax(0, 1)(data)

	// Should be scaled to [0, 1]
	if !floatEquals(normalized[0].Y, 0, 0.01) {
		t.Errorf("Min: expected 0, got %f", normalized[0].Y)
	}
	if !floatEquals(normalized[2].Y, 1, 0.01) {
		t.Errorf("Max: expected 1, got %f", normalized[2].Y)
	}
	if !floatEquals(normalized[1].Y, 0.5, 0.01) {
		t.Errorf("Mid: expected 0.5, got %f", normalized[1].Y)
	}
}

func TestScale(t *testing.T) {
	data := []DataPoint{
		{Y: 10}, {Y: 20}, {Y: 30},
	}

	scaled := Scale(2.0)(data)

	expected := []float64{20, 40, 60}
	for i, exp := range expected {
		if scaled[i].Y != exp {
			t.Errorf("Index %d: expected %f, got %f", i, exp, scaled[i].Y)
		}
	}
}

func TestOffset(t *testing.T) {
	data := []DataPoint{
		{Y: 10}, {Y: 20}, {Y: 30},
	}

	offset := Offset(5)(data)

	expected := []float64{15, 25, 35}
	for i, exp := range expected {
		if offset[i].Y != exp {
			t.Errorf("Index %d: expected %f, got %f", i, exp, offset[i].Y)
		}
	}
}

func TestClamp(t *testing.T) {
	data := []DataPoint{
		{Y: 5}, {Y: 15}, {Y: 25},
	}

	clamped := Clamp(10, 20)(data)

	expected := []float64{10, 15, 20}
	for i, exp := range expected {
		if clamped[i].Y != exp {
			t.Errorf("Index %d: expected %f, got %f", i, exp, clamped[i].Y)
		}
	}
}

func TestAbs(t *testing.T) {
	data := []DataPoint{
		{Y: -10}, {Y: 20}, {Y: -30},
	}

	absolute := Abs()(data)

	for _, d := range absolute {
		if d.Y < 0 {
			t.Errorf("Expected positive value, got %f", d.Y)
		}
	}
}

// ==================== Helper Function Tests ====================

func TestAggregationFunctions(t *testing.T) {
	values := []float64{1, 2, 3, 4, 5}

	if Sum(values) != 15 {
		t.Errorf("Sum: expected 15, got %f", Sum(values))
	}

	if Mean(values) != 3 {
		t.Errorf("Mean: expected 3, got %f", Mean(values))
	}

	if Max(values) != 5 {
		t.Errorf("Max: expected 5, got %f", Max(values))
	}

	if Min(values) != 1 {
		t.Errorf("Min: expected 1, got %f", Min(values))
	}

	if Count(values) != 5 {
		t.Errorf("Count: expected 5, got %f", Count(values))
	}

	if Median(values) != 3 {
		t.Errorf("Median: expected 3, got %f", Median(values))
	}
}

func TestToDataPoints(t *testing.T) {
	now := time.Now()
	points := []TimeSeriesPoint{
		{Time: now, Value: 10, Label: "A"},
		{Time: now.Add(time.Hour), Value: 20, Label: "B"},
	}

	dataPoints := ToDataPoints(points)

	if len(dataPoints) != 2 {
		t.Fatalf("Expected 2 data points, got %d", len(dataPoints))
	}

	if dataPoints[0].Y != 10 {
		t.Errorf("Expected Y=10, got %f", dataPoints[0].Y)
	}

	if dataPoints[0].Label != "A" {
		t.Errorf("Expected Label=A, got %s", dataPoints[0].Label)
	}
}

func TestFromDataPoints(t *testing.T) {
	now := time.Now()
	dataPoints := []DataPoint{
		{X: now, Y: 10, Label: "A"},
		{X: now.Add(time.Hour), Y: 20, Label: "B"},
	}

	points := FromDataPoints(dataPoints)

	if len(points) != 2 {
		t.Fatalf("Expected 2 time series points, got %d", len(points))
	}

	if points[0].Value != 10 {
		t.Errorf("Expected Value=10, got %f", points[0].Value)
	}

	if points[0].Label != "A" {
		t.Errorf("Expected Label=A, got %s", points[0].Label)
	}
}

// ==================== Benchmarks ====================

func BenchmarkBin(b *testing.B) {
	data := make([]DataPoint, 1000)
	for i := range data {
		data[i] = DataPoint{Y: float64(i)}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Bin(BinOptions{Count: 20})(data)
	}
}

func BenchmarkGroupBy(b *testing.B) {
	data := make([]DataPoint, 1000)
	for i := range data {
		data[i] = DataPoint{Label: string(rune('A' + i%10)), Y: float64(i)}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GroupBy(GroupOptions{By: "Label", Aggregate: Sum})(data)
	}
}

func BenchmarkMovingAverage(b *testing.B) {
	data := make([]DataPoint, 1000)
	for i := range data {
		data[i] = DataPoint{Y: float64(i)}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MovingAverage(10)(data)
	}
}
