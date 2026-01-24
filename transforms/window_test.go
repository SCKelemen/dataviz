package transforms

import (
	"testing"
	"time"
)

func TestSlidingWindow(t *testing.T) {
	data := []DataPoint{
		{Y: 1}, {Y: 2}, {Y: 3}, {Y: 4}, {Y: 5},
	}

	sw := NewSlidingWindow(3)
	windows := sw.Windows(data)

	// Each point should be in multiple windows (overlapping)
	// Point at index 2 should be in windows [0, 1, 2]
	if len(windows[2]) < 2 {
		t.Errorf("Expected point 2 to be in multiple windows, got %d", len(windows[2]))
	}
}

func TestSlidingWindow_WithStep(t *testing.T) {
	data := []DataPoint{
		{Y: 1}, {Y: 2}, {Y: 3}, {Y: 4}, {Y: 5}, {Y: 6},
	}

	sw := NewSlidingWindow(3).WithStep(2)
	windows := sw.Windows(data)

	// With step=2, windows should be [0,1,2], [2,3,4], [4,5,6]
	// Point at index 2 should be in 2 windows
	if len(windows[2]) != 2 {
		t.Errorf("Expected point 2 to be in 2 windows, got %d", len(windows[2]))
	}
}

func TestTumblingWindow(t *testing.T) {
	data := []DataPoint{
		{Y: 1}, {Y: 2}, {Y: 3}, {Y: 4}, {Y: 5}, {Y: 6},
	}

	tw := NewTumblingWindow(3)
	windows := tw.Windows(data)

	// Each point should be in exactly one window (non-overlapping)
	for i, wins := range windows {
		if len(wins) != 1 {
			t.Errorf("Point %d should be in exactly 1 window, got %d", i, len(wins))
		}
	}

	// First 3 points should be in window 0
	if windows[0][0] != 0 || windows[1][0] != 0 || windows[2][0] != 0 {
		t.Error("First 3 points should be in window 0")
	}

	// Next 3 points should be in window 1
	if windows[3][0] != 1 || windows[4][0] != 1 || windows[5][0] != 1 {
		t.Error("Next 3 points should be in window 1")
	}
}

func TestHoppingWindow(t *testing.T) {
	data := []DataPoint{
		{Y: 1}, {Y: 2}, {Y: 3}, {Y: 4}, {Y: 5}, {Y: 6},
	}

	hw := NewHoppingWindow(3, 2)
	windows := hw.Windows(data)

	// With size=3 and hop=2, some points will be in multiple windows
	// Point at index 2 should be in windows 0 and 1
	if len(windows[2]) < 2 {
		t.Errorf("Expected point 2 to be in multiple windows, got %d", len(windows[2]))
	}
}

func TestSessionWindow(t *testing.T) {
	now := time.Now()

	data := []DataPoint{
		{X: now, Y: 1},
		{X: now.Add(1 * time.Second), Y: 2},
		{X: now.Add(2 * time.Second), Y: 3},
		{X: now.Add(10 * time.Second), Y: 4}, // Gap > 5s, new session
		{X: now.Add(11 * time.Second), Y: 5},
	}

	sw := NewSessionWindow(5 * time.Second)
	windows := sw.Windows(data)

	// First 3 points should be in session 0
	if windows[0][0] != 0 || windows[1][0] != 0 || windows[2][0] != 0 {
		t.Error("First 3 points should be in session 0")
	}

	// Last 2 points should be in session 1
	if windows[3][0] != 1 || windows[4][0] != 1 {
		t.Error("Last 2 points should be in session 1")
	}
}

func TestTimeWindow(t *testing.T) {
	now := time.Now()

	data := []DataPoint{
		{X: now, Y: 1},
		{X: now.Add(30 * time.Second), Y: 2},
		{X: now.Add(65 * time.Second), Y: 3},
		{X: now.Add(120 * time.Second), Y: 4},
	}

	tw := NewTimeWindow(1 * time.Minute)
	windows := tw.Windows(data)

	// Points should be grouped into 1-minute buckets
	// First two points in window 0
	if windows[0][0] != 0 || windows[1][0] != 0 {
		t.Error("First two points should be in window 0")
	}

	// Third point in window 1
	if windows[2][0] != 1 {
		t.Error("Third point should be in window 1")
	}

	// Fourth point in window 2
	if windows[3][0] != 2 {
		t.Error("Fourth point should be in window 2")
	}
}

func TestSnapshotWindow(t *testing.T) {
	now := time.Now()

	data := []DataPoint{
		{X: now, Y: 1},
		{X: now, Y: 2},                    // Same timestamp as first
		{X: now.Add(1 * time.Second), Y: 3},
		{X: now.Add(1 * time.Second), Y: 4}, // Same timestamp as third
		{X: now.Add(2 * time.Second), Y: 5},
	}

	sw := NewSnapshotWindow()
	windows := sw.Windows(data)

	// First two points at same timestamp should be in window 0
	if windows[0][0] != 0 || windows[1][0] != 0 {
		t.Error("First two points should be in window 0")
	}

	// Next two points at same timestamp should be in window 1
	if windows[2][0] != 1 || windows[3][0] != 1 {
		t.Error("Third and fourth points should be in window 1")
	}

	// Last point has unique timestamp, should be in window 2
	if windows[4][0] != 2 {
		t.Error("Fifth point should be in window 2")
	}
}

func TestSnapshotWindow_Apply(t *testing.T) {
	now := time.Now()

	data := []DataPoint{
		{X: now, Y: 10},
		{X: now, Y: 20},
		{X: now.Add(1 * time.Second), Y: 30},
		{X: now.Add(1 * time.Second), Y: 40},
	}

	sw := NewSnapshotWindow()
	result := ApplyWindow(sw, Sum)(data)

	// Should produce 2 windows
	if len(result) != 2 {
		t.Errorf("Expected 2 windows, got %d", len(result))
	}

	// First window: sum of [10, 20] = 30
	if result[0].Y != 30 {
		t.Errorf("First window sum: expected 30, got %f", result[0].Y)
	}

	// Second window: sum of [30, 40] = 70
	if result[1].Y != 70 {
		t.Errorf("Second window sum: expected 70, got %f", result[1].Y)
	}
}

func TestApplyWindow_Sliding(t *testing.T) {
	data := []DataPoint{
		{Y: 1}, {Y: 2}, {Y: 3}, {Y: 4}, {Y: 5},
	}

	sw := NewSlidingWindow(3)
	result := ApplyWindow(sw, Mean)(data)

	// Should produce one aggregated point per window
	if len(result) < 3 {
		t.Errorf("Expected at least 3 windows, got %d", len(result))
	}
}

func TestApplyWindow_Tumbling(t *testing.T) {
	data := []DataPoint{
		{Y: 1}, {Y: 2}, {Y: 3}, {Y: 4}, {Y: 5}, {Y: 6},
	}

	tw := NewTumblingWindow(3)
	result := ApplyWindow(tw, Sum)(data)

	// Should produce exactly 2 windows (3 points each)
	if len(result) != 2 {
		t.Errorf("Expected 2 windows, got %d", len(result))
	}

	// First window: sum of [1, 2, 3] = 6
	if result[0].Y != 6 {
		t.Errorf("First window sum: expected 6, got %f", result[0].Y)
	}

	// Second window: sum of [4, 5, 6] = 15
	if result[1].Y != 15 {
		t.Errorf("Second window sum: expected 15, got %f", result[1].Y)
	}
}

func TestWindowAggregate(t *testing.T) {
	data := []DataPoint{
		{Y: 10}, {Y: 20}, {Y: 30}, {Y: 40},
	}

	// Tumbling windows of size 2
	result := WindowAggregate(2, Mean, "tumbling")(data)

	if len(result) != 2 {
		t.Errorf("Expected 2 windows, got %d", len(result))
	}

	// First window: mean of [10, 20] = 15
	if result[0].Y != 15 {
		t.Errorf("First window mean: expected 15, got %f", result[0].Y)
	}

	// Second window: mean of [30, 40] = 35
	if result[1].Y != 35 {
		t.Errorf("Second window mean: expected 35, got %f", result[1].Y)
	}
}

func TestWindowedMean(t *testing.T) {
	data := []DataPoint{
		{Y: 2}, {Y: 4}, {Y: 6}, {Y: 8},
	}

	result := WindowedMean(2, "tumbling")(data)

	if len(result) != 2 {
		t.Errorf("Expected 2 windows, got %d", len(result))
	}

	// Means should be 3 and 7
	if result[0].Y != 3 {
		t.Errorf("First mean: expected 3, got %f", result[0].Y)
	}
	if result[1].Y != 7 {
		t.Errorf("Second mean: expected 7, got %f", result[1].Y)
	}
}

func TestWindowedSum(t *testing.T) {
	data := []DataPoint{
		{Y: 1}, {Y: 2}, {Y: 3}, {Y: 4},
	}

	result := WindowedSum(2, "tumbling")(data)

	if len(result) != 2 {
		t.Errorf("Expected 2 windows, got %d", len(result))
	}

	// Sums should be 3 and 7
	if result[0].Y != 3 {
		t.Errorf("First sum: expected 3, got %f", result[0].Y)
	}
	if result[1].Y != 7 {
		t.Errorf("Second sum: expected 7, got %f", result[1].Y)
	}
}

func TestWindowBounds_Sliding(t *testing.T) {
	data := []DataPoint{
		{Y: 1}, {Y: 2}, {Y: 3}, {Y: 4}, {Y: 5},
	}

	sw := NewSlidingWindow(3)

	// Window 0 should be [0, 3)
	start, end := sw.WindowBounds(0, data)
	if start != 0 || end != 3 {
		t.Errorf("Window 0: expected [0, 3), got [%d, %d)", start, end)
	}

	// Window 1 should be [1, 4)
	start, end = sw.WindowBounds(1, data)
	if start != 1 || end != 4 {
		t.Errorf("Window 1: expected [1, 4), got [%d, %d)", start, end)
	}
}

func TestWindowBounds_Tumbling(t *testing.T) {
	data := []DataPoint{
		{Y: 1}, {Y: 2}, {Y: 3}, {Y: 4}, {Y: 5}, {Y: 6},
	}

	tw := NewTumblingWindow(3)

	// Window 0 should be [0, 3)
	start, end := tw.WindowBounds(0, data)
	if start != 0 || end != 3 {
		t.Errorf("Window 0: expected [0, 3), got [%d, %d)", start, end)
	}

	// Window 1 should be [3, 6)
	start, end = tw.WindowBounds(1, data)
	if start != 3 || end != 6 {
		t.Errorf("Window 1: expected [3, 6), got [%d, %d)", start, end)
	}
}

func BenchmarkSlidingWindow(b *testing.B) {
	data := make([]DataPoint, 1000)
	for i := range data {
		data[i] = DataPoint{Y: float64(i)}
	}

	sw := NewSlidingWindow(10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sw.Windows(data)
	}
}

func BenchmarkTumblingWindow(b *testing.B) {
	data := make([]DataPoint, 1000)
	for i := range data {
		data[i] = DataPoint{Y: float64(i)}
	}

	tw := NewTumblingWindow(10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tw.Windows(data)
	}
}

func BenchmarkApplyWindow(b *testing.B) {
	data := make([]DataPoint, 1000)
	for i := range data {
		data[i] = DataPoint{Y: float64(i)}
	}

	sw := NewSlidingWindow(10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ApplyWindow(sw, Mean)(data)
	}
}
