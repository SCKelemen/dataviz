package scales

import (
	"math"
	"time"

	"github.com/SCKelemen/units"
)

// TimeScale implements a continuous scale for temporal data.
// Maps a time domain [t0, t1] to a continuous range [r0, r1] using linear interpolation.
//
// Ranges use units.Length to support relative units (%, px, em, etc.).
//
// Example:
//   start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
//   end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
//   scale := NewTimeScale(
//     [2]time.Time{start, end},
//     [2]units.Length{units.Px(0), units.Px(500)},
//   )
//   scale.Apply(time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)) // Mid-year position
type TimeScale struct {
	domain [2]time.Time
	range_ [2]units.Length
	clamp  bool
}

// NewTimeScale creates a new time scale
func NewTimeScale(domain [2]time.Time, range_ [2]units.Length) *TimeScale {
	return &TimeScale{
		domain: domain,
		range_: range_,
		clamp:  false,
	}
}

// Apply maps a domain value to a range value
func (s *TimeScale) Apply(value interface{}) units.Length {
	t := s.ApplyValue(value)

	// Interpolate between range values
	r0 := s.range_[0].Value
	r1 := s.range_[1].Value
	unit := s.range_[0].Unit

	result := r0 + t*(r1-r0)

	return units.Length{Value: result, Unit: unit}
}

// ApplyValue maps a domain value to a normalized value (0-1 interpolation factor)
func (s *TimeScale) ApplyValue(value interface{}) float64 {
	var t time.Time
	switch v := value.(type) {
	case time.Time:
		t = v
	case *time.Time:
		if v == nil {
			return 0
		}
		t = *v
	default:
		return 0
	}

	// Convert to Unix timestamps (seconds since epoch)
	t0 := float64(s.domain[0].Unix())
	t1 := float64(s.domain[1].Unix())
	tv := float64(t.Unix())

	// Linear interpolation parameter
	interpolation := (tv - t0) / (t1 - t0)

	if s.clamp {
		interpolation = clampFloat(interpolation, 0, 1)
	}

	return interpolation
}

// Invert maps a range value back to a domain value
func (s *TimeScale) Invert(value units.Length) time.Time {
	v := value.Value
	r0 := s.range_[0].Value
	r1 := s.range_[1].Value

	t := (v - r0) / (r1 - r0)

	if s.clamp {
		t = clampFloat(t, 0, 1)
	}

	return s.InvertValue(t)
}

// InvertValue maps a normalized value (0-1) back to a domain value
func (s *TimeScale) InvertValue(t float64) time.Time {
	t0 := s.domain[0].Unix()
	t1 := s.domain[1].Unix()

	timestamp := t0 + int64(float64(t1-t0)*t)
	return time.Unix(timestamp, 0).UTC()
}

// Domain returns the input domain
func (s *TimeScale) Domain() interface{} {
	return s.domain
}

// Range returns the output range
func (s *TimeScale) Range() [2]units.Length {
	return s.range_
}

// Type returns the scale type
func (s *TimeScale) Type() ScaleType {
	return ScaleTypeTime
}

// Clone creates a copy of this scale
func (s *TimeScale) Clone() Scale {
	return &TimeScale{
		domain: s.domain,
		range_: s.range_,
		clamp:  s.clamp,
	}
}

// Clamp enables/disables clamping output to range
func (s *TimeScale) Clamp(enabled bool) *TimeScale {
	s.clamp = enabled
	return s
}

// Nice rounds the domain to nice time boundaries
func (s *TimeScale) Nice(interval TimeInterval) *TimeScale {
	t0 := s.domain[0]
	t1 := s.domain[1]

	switch interval {
	case TimeIntervalYear:
		t0 = time.Date(t0.Year(), 1, 1, 0, 0, 0, 0, t0.Location())
		t1 = time.Date(t1.Year()+1, 1, 1, 0, 0, 0, 0, t1.Location())

	case TimeIntervalMonth:
		t0 = time.Date(t0.Year(), t0.Month(), 1, 0, 0, 0, 0, t0.Location())
		t1 = time.Date(t1.Year(), t1.Month()+1, 1, 0, 0, 0, 0, t1.Location())

	case TimeIntervalDay:
		t0 = time.Date(t0.Year(), t0.Month(), t0.Day(), 0, 0, 0, 0, t0.Location())
		t1 = time.Date(t1.Year(), t1.Month(), t1.Day()+1, 0, 0, 0, 0, t1.Location())

	case TimeIntervalHour:
		t0 = time.Date(t0.Year(), t0.Month(), t0.Day(), t0.Hour(), 0, 0, 0, t0.Location())
		t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), t1.Hour()+1, 0, 0, 0, t1.Location())

	case TimeIntervalMinute:
		t0 = time.Date(t0.Year(), t0.Month(), t0.Day(), t0.Hour(), t0.Minute(), 0, 0, t0.Location())
		t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), t1.Hour(), t1.Minute()+1, 0, 0, t1.Location())

	case TimeIntervalSecond:
		t0 = time.Date(t0.Year(), t0.Month(), t0.Day(), t0.Hour(), t0.Minute(), t0.Second(), 0, t0.Location())
		t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), t1.Hour(), t1.Minute(), t1.Second()+1, 0, t1.Location())
	}

	s.domain[0] = t0
	s.domain[1] = t1

	return s
}

// Ticks generates nice tick values for axes
func (s *TimeScale) Ticks(count int) []time.Time {
	if count <= 0 {
		count = 10
	}

	t0 := s.domain[0]
	t1 := s.domain[1]

	// Calculate duration between domain endpoints
	duration := t1.Sub(t0)

	// Determine appropriate interval based on duration and desired count
	interval := s.selectInterval(duration, count)

	var ticks []time.Time

	switch interval {
	case TimeIntervalYear:
		ticks = s.ticksYear(t0, t1, count)
	case TimeIntervalMonth:
		ticks = s.ticksMonth(t0, t1, count)
	case TimeIntervalDay:
		ticks = s.ticksDay(t0, t1, count)
	case TimeIntervalHour:
		ticks = s.ticksHour(t0, t1, count)
	case TimeIntervalMinute:
		ticks = s.ticksMinute(t0, t1, count)
	case TimeIntervalSecond:
		ticks = s.ticksSecond(t0, t1, count)
	default:
		ticks = s.ticksSecond(t0, t1, count)
	}

	return ticks
}

// selectInterval chooses appropriate interval based on duration
func (s *TimeScale) selectInterval(duration time.Duration, count int) TimeInterval {
	seconds := duration.Seconds()
	targetStep := seconds / float64(count)

	if targetStep >= 365*24*3600 {
		return TimeIntervalYear
	}
	if targetStep >= 30*24*3600 {
		return TimeIntervalMonth
	}
	if targetStep >= 24*3600 {
		return TimeIntervalDay
	}
	if targetStep >= 3600 {
		return TimeIntervalHour
	}
	if targetStep >= 60 {
		return TimeIntervalMinute
	}
	return TimeIntervalSecond
}

// ticksYear generates yearly ticks
func (s *TimeScale) ticksYear(t0, t1 time.Time, count int) []time.Time {
	yearStart := t0.Year()
	yearEnd := t1.Year()
	years := yearEnd - yearStart + 1

	step := int(math.Max(1, math.Round(float64(years)/float64(count))))

	var ticks []time.Time
	for year := yearStart; year <= yearEnd; year += step {
		ticks = append(ticks, time.Date(year, 1, 1, 0, 0, 0, 0, t0.Location()))
	}

	return ticks
}

// ticksMonth generates monthly ticks
func (s *TimeScale) ticksMonth(t0, t1 time.Time, count int) []time.Time {
	var ticks []time.Time
	current := time.Date(t0.Year(), t0.Month(), 1, 0, 0, 0, 0, t0.Location())

	for current.Before(t1) || current.Equal(t1) {
		ticks = append(ticks, current)
		current = current.AddDate(0, 1, 0)
	}

	return ticks
}

// ticksDay generates daily ticks
func (s *TimeScale) ticksDay(t0, t1 time.Time, count int) []time.Time {
	var ticks []time.Time
	current := time.Date(t0.Year(), t0.Month(), t0.Day(), 0, 0, 0, 0, t0.Location())

	days := int(t1.Sub(t0).Hours() / 24)
	step := int(math.Max(1, math.Round(float64(days)/float64(count))))

	for current.Before(t1) || current.Equal(t1) {
		ticks = append(ticks, current)
		current = current.AddDate(0, 0, step)
	}

	return ticks
}

// ticksHour generates hourly ticks
func (s *TimeScale) ticksHour(t0, t1 time.Time, count int) []time.Time {
	var ticks []time.Time
	current := time.Date(t0.Year(), t0.Month(), t0.Day(), t0.Hour(), 0, 0, 0, t0.Location())

	hours := int(t1.Sub(t0).Hours())
	step := int(math.Max(1, math.Round(float64(hours)/float64(count))))

	for current.Before(t1) || current.Equal(t1) {
		ticks = append(ticks, current)
		current = current.Add(time.Duration(step) * time.Hour)
	}

	return ticks
}

// ticksMinute generates minute ticks
func (s *TimeScale) ticksMinute(t0, t1 time.Time, count int) []time.Time {
	var ticks []time.Time
	current := time.Date(t0.Year(), t0.Month(), t0.Day(), t0.Hour(), t0.Minute(), 0, 0, t0.Location())

	minutes := int(t1.Sub(t0).Minutes())
	step := int(math.Max(1, math.Round(float64(minutes)/float64(count))))

	for current.Before(t1) || current.Equal(t1) {
		ticks = append(ticks, current)
		current = current.Add(time.Duration(step) * time.Minute)
	}

	return ticks
}

// ticksSecond generates second ticks
func (s *TimeScale) ticksSecond(t0, t1 time.Time, count int) []time.Time {
	var ticks []time.Time
	current := time.Date(t0.Year(), t0.Month(), t0.Day(), t0.Hour(), t0.Minute(), t0.Second(), 0, t0.Location())

	seconds := int(t1.Sub(t0).Seconds())
	step := int(math.Max(1, math.Round(float64(seconds)/float64(count))))

	for current.Before(t1) || current.Equal(t1) {
		ticks = append(ticks, current)
		current = current.Add(time.Duration(step) * time.Second)
	}

	return ticks
}

// WithDomain sets a new domain
func (s *TimeScale) WithDomain(domain [2]time.Time) *TimeScale {
	s.domain = domain
	return s
}

// WithRange sets a new range
func (s *TimeScale) WithRange(range_ [2]units.Length) *TimeScale {
	s.range_ = range_
	return s
}
