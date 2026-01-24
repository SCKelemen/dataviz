package transforms

import "time"

// WindowStrategy defines how to partition data into windows
type WindowStrategy interface {
	// Windows returns the window assignments for each data point
	// Returns a slice of window IDs for each input point
	Windows(data []DataPoint) [][]int

	// WindowBounds returns the start and end indices for a window ID
	WindowBounds(windowID int, data []DataPoint) (start, end int)
}

// SlidingWindow creates overlapping windows of fixed size
type SlidingWindow struct {
	Size int // Window size in number of points
	Step int // Step size (default: 1 for maximum overlap)
}

// NewSlidingWindow creates a sliding window strategy
func NewSlidingWindow(size int) *SlidingWindow {
	return &SlidingWindow{
		Size: size,
		Step: 1,
	}
}

// WithStep sets the step size for sliding windows
func (sw *SlidingWindow) WithStep(step int) *SlidingWindow {
	sw.Step = step
	return sw
}

// Windows returns window assignments for sliding windows
func (sw *SlidingWindow) Windows(data []DataPoint) [][]int {
	if len(data) == 0 || sw.Size <= 0 {
		return nil
	}

	step := sw.Step
	if step <= 0 {
		step = 1
	}

	// Calculate number of windows
	numWindows := 0
	for i := 0; i+sw.Size <= len(data); i += step {
		numWindows++
	}

	windows := make([][]int, len(data))
	windowID := 0

	for i := 0; i+sw.Size <= len(data); i += step {
		// Assign window ID to all points in this window
		for j := i; j < i+sw.Size; j++ {
			windows[j] = append(windows[j], windowID)
		}
		windowID++
	}

	return windows
}

// WindowBounds returns the bounds for a window
func (sw *SlidingWindow) WindowBounds(windowID int, data []DataPoint) (start, end int) {
	step := sw.Step
	if step <= 0 {
		step = 1
	}

	start = windowID * step
	end = start + sw.Size
	if end > len(data) {
		end = len(data)
	}
	return
}

// TumblingWindow creates non-overlapping windows of fixed size
type TumblingWindow struct {
	Size int // Window size in number of points
}

// NewTumblingWindow creates a tumbling window strategy
func NewTumblingWindow(size int) *TumblingWindow {
	return &TumblingWindow{Size: size}
}

// Windows returns window assignments for tumbling windows
func (tw *TumblingWindow) Windows(data []DataPoint) [][]int {
	if len(data) == 0 || tw.Size <= 0 {
		return nil
	}

	windows := make([][]int, len(data))
	windowID := 0

	for i := 0; i < len(data); i += tw.Size {
		end := i + tw.Size
		if end > len(data) {
			end = len(data)
		}

		// Assign window ID to all points in this chunk
		for j := i; j < end; j++ {
			windows[j] = []int{windowID}
		}
		windowID++
	}

	return windows
}

// WindowBounds returns the bounds for a tumbling window
func (tw *TumblingWindow) WindowBounds(windowID int, data []DataPoint) (start, end int) {
	start = windowID * tw.Size
	end = start + tw.Size
	if end > len(data) {
		end = len(data)
	}
	return
}

// HoppingWindow creates fixed-size windows with configurable hop/stride
type HoppingWindow struct {
	Size int // Window size in number of points
	Hop  int // Hop size (step between windows)
}

// NewHoppingWindow creates a hopping window strategy
func NewHoppingWindow(size, hop int) *HoppingWindow {
	return &HoppingWindow{
		Size: size,
		Hop:  hop,
	}
}

// Windows returns window assignments for hopping windows
func (hw *HoppingWindow) Windows(data []DataPoint) [][]int {
	if len(data) == 0 || hw.Size <= 0 || hw.Hop <= 0 {
		return nil
	}

	windows := make([][]int, len(data))
	windowID := 0

	for i := 0; i < len(data); i += hw.Hop {
		end := i + hw.Size
		if end > len(data) {
			end = len(data)
		}

		// Assign window ID to all points in this window
		for j := i; j < end; j++ {
			windows[j] = append(windows[j], windowID)
		}
		windowID++

		if i+hw.Size >= len(data) {
			break
		}
	}

	return windows
}

// WindowBounds returns the bounds for a hopping window
func (hw *HoppingWindow) WindowBounds(windowID int, data []DataPoint) (start, end int) {
	start = windowID * hw.Hop
	end = start + hw.Size
	if end > len(data) {
		end = len(data)
	}
	return
}

// SessionWindow groups points based on gaps in time
type SessionWindow struct {
	GapThreshold time.Duration // Maximum gap between points in the same session
}

// NewSessionWindow creates a session window strategy
func NewSessionWindow(gapThreshold time.Duration) *SessionWindow {
	return &SessionWindow{GapThreshold: gapThreshold}
}

// Windows returns window assignments for session windows
func (sw *SessionWindow) Windows(data []DataPoint) [][]int {
	if len(data) == 0 {
		return nil
	}

	windows := make([][]int, len(data))
	windowID := 0
	windows[0] = []int{windowID}

	for i := 1; i < len(data); i++ {
		// Check if there's a gap larger than threshold
		prevTime, prevOk := data[i-1].X.(time.Time)
		currTime, currOk := data[i].X.(time.Time)

		if prevOk && currOk {
			gap := currTime.Sub(prevTime)
			if gap > sw.GapThreshold {
				// Start new session
				windowID++
			}
		}

		windows[i] = []int{windowID}
	}

	return windows
}

// WindowBounds returns the bounds for a session window
func (sw *SessionWindow) WindowBounds(windowID int, data []DataPoint) (start, end int) {
	// Find the first and last occurrence of windowID
	start = -1
	end = -1

	windows := sw.Windows(data)
	for i, wins := range windows {
		for _, wid := range wins {
			if wid == windowID {
				if start == -1 {
					start = i
				}
				end = i + 1
			}
		}
	}

	if start == -1 {
		return 0, 0
	}
	return
}

// TimeWindow partitions data by fixed time intervals
type TimeWindow struct {
	Interval time.Duration // Time interval for each window
}

// NewTimeWindow creates a time-based window strategy
func NewTimeWindow(interval time.Duration) *TimeWindow {
	return &TimeWindow{Interval: interval}
}

// Windows returns window assignments for time-based windows
func (tw *TimeWindow) Windows(data []DataPoint) [][]int {
	if len(data) == 0 || tw.Interval <= 0 {
		return nil
	}

	windows := make([][]int, len(data))

	// Find the earliest time
	var minTime time.Time
	for _, d := range data {
		if t, ok := d.X.(time.Time); ok {
			if minTime.IsZero() || t.Before(minTime) {
				minTime = t
			}
		}
	}

	if minTime.IsZero() {
		return nil
	}

	// Assign window IDs based on time buckets
	for i, d := range data {
		if t, ok := d.X.(time.Time); ok {
			elapsed := t.Sub(minTime)
			windowID := int(elapsed / tw.Interval)
			windows[i] = []int{windowID}
		}
	}

	return windows
}

// WindowBounds returns the bounds for a time window
func (tw *TimeWindow) WindowBounds(windowID int, data []DataPoint) (start, end int) {
	windows := tw.Windows(data)
	start = -1
	end = -1

	for i, wins := range windows {
		for _, wid := range wins {
			if wid == windowID {
				if start == -1 {
					start = i
				}
				end = i + 1
			}
		}
	}

	if start == -1 {
		return 0, 0
	}
	return
}

// SnapshotWindow groups events that occur at the exact same timestamp
// Unlike other windows, snapshot windows are event-driven and fire on every event
type SnapshotWindow struct{}

// NewSnapshotWindow creates a snapshot window strategy
func NewSnapshotWindow() *SnapshotWindow {
	return &SnapshotWindow{}
}

// Windows returns window assignments for snapshot windows
// Each unique timestamp gets its own window
func (sw *SnapshotWindow) Windows(data []DataPoint) [][]int {
	if len(data) == 0 {
		return nil
	}

	windows := make([][]int, len(data))

	// Map timestamps to window IDs
	timestampToWindow := make(map[interface{}]int)
	windowID := 0

	for i, d := range data {
		// Use X as the timestamp key
		// If X is time.Time, use it directly
		// Otherwise use the value itself
		var key interface{}
		if t, ok := d.X.(time.Time); ok {
			// Truncate to remove subsecond precision if needed
			key = t.Unix()
		} else {
			key = d.X
		}

		// Check if we've seen this timestamp
		wid, exists := timestampToWindow[key]
		if !exists {
			wid = windowID
			timestampToWindow[key] = wid
			windowID++
		}

		windows[i] = []int{wid}
	}

	return windows
}

// WindowBounds returns the bounds for a snapshot window
func (sw *SnapshotWindow) WindowBounds(windowID int, data []DataPoint) (start, end int) {
	windows := sw.Windows(data)
	start = -1
	end = -1

	for i, wins := range windows {
		for _, wid := range wins {
			if wid == windowID {
				if start == -1 {
					start = i
				}
				end = i + 1
			}
		}
	}

	if start == -1 {
		return 0, 0
	}
	return
}

// ApplyWindow applies a window strategy and aggregates each window
func ApplyWindow(strategy WindowStrategy, fn AggregateFunc) Transform {
	return func(data []DataPoint) []DataPoint {
		if len(data) == 0 {
			return nil
		}

		windows := strategy.Windows(data)
		if len(windows) == 0 {
			return nil
		}

		// Find all unique window IDs
		windowIDs := make(map[int]bool)
		for _, wins := range windows {
			for _, wid := range wins {
				windowIDs[wid] = true
			}
		}

		// Aggregate each window
		result := make([]DataPoint, 0, len(windowIDs))
		for wid := range windowIDs {
			start, end := strategy.WindowBounds(wid, data)
			if start >= end {
				continue
			}

			// Extract values in this window
			values := make([]float64, 0, end-start)
			for i := start; i < end; i++ {
				values = append(values, data[i].Y)
			}

			// Aggregate
			aggValue := fn(values)

			// Use the first point in the window as the base
			point := data[start]
			point.Y = aggValue
			point.Value = aggValue
			point.Count = len(values)
			point.Index = wid

			result = append(result, point)
		}

		return result
	}
}

// WindowAggregate creates a windowed aggregation transform
func WindowAggregate(size int, fn AggregateFunc, strategy string) Transform {
	var ws WindowStrategy

	switch strategy {
	case "tumbling":
		ws = NewTumblingWindow(size)
	case "hopping":
		ws = NewHoppingWindow(size, size/2) // 50% overlap by default
	case "sliding":
		fallthrough
	default:
		ws = NewSlidingWindow(size)
	}

	return ApplyWindow(ws, fn)
}

// SnapshotAggregate creates a snapshot window aggregation transform
func SnapshotAggregate(fn AggregateFunc) Transform {
	return ApplyWindow(NewSnapshotWindow(), fn)
}

// SnapshotMean creates a snapshot window mean transform
func SnapshotMean() Transform {
	return SnapshotAggregate(Mean)
}

// SnapshotSum creates a snapshot window sum transform
func SnapshotSum() Transform {
	return SnapshotAggregate(Sum)
}

// SnapshotMax creates a snapshot window max transform
func SnapshotMax() Transform {
	return SnapshotAggregate(Max)
}

// SnapshotMin creates a snapshot window min transform
func SnapshotMin() Transform {
	return SnapshotAggregate(Min)
}

// SnapshotCount creates a snapshot window count transform
func SnapshotCount() Transform {
	return SnapshotAggregate(Count)
}

// WindowedMean creates a windowed mean transform
func WindowedMean(size int, strategy string) Transform {
	return WindowAggregate(size, Mean, strategy)
}

// WindowedSum creates a windowed sum transform
func WindowedSum(size int, strategy string) Transform {
	return WindowAggregate(size, Sum, strategy)
}

// WindowedMax creates a windowed max transform
func WindowedMax(size int, strategy string) Transform {
	return WindowAggregate(size, Max, strategy)
}

// WindowedMin creates a windowed min transform
func WindowedMin(size int, strategy string) Transform {
	return WindowAggregate(size, Min, strategy)
}
