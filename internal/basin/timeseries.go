package basin

import "math"

// TimeSeries is a sequence of values at regular intervals.
type TimeSeries struct {
	Values []float64
	Dt     float64 // time step
}

// NewTimeSeries creates a time series.
func NewTimeSeries(values []float64, dt float64) *TimeSeries {
	return &TimeSeries{Values: values, Dt: dt}
}

// Len returns the number of values.
func (ts *TimeSeries) Len() int { return len(ts.Values) }

// Sum returns the total.
func (ts *TimeSeries) Sum() float64 {
	var s float64
	for _, v := range ts.Values {
		s += v
	}
	return s
}

// Mean returns the arithmetic mean.
func (ts *TimeSeries) Mean() float64 {
	if len(ts.Values) == 0 {
		return 0
	}
	return ts.Sum() / float64(len(ts.Values))
}

// Max returns the maximum value and its index.
func (ts *TimeSeries) Max() (float64, int) {
	if len(ts.Values) == 0 {
		return 0, -1
	}
	max := ts.Values[0]
	idx := 0
	for i, v := range ts.Values[1:] {
		if v > max {
			max = v
			idx = i + 1
		}
	}
	return max, idx
}

// Min returns the minimum value.
func (ts *TimeSeries) Min() float64 {
	if len(ts.Values) == 0 {
		return 0
	}
	min := ts.Values[0]
	for _, v := range ts.Values[1:] {
		if v < min {
			min = v
		}
	}
	return min
}

// StdDev returns the standard deviation.
func (ts *TimeSeries) StdDev() float64 {
	if len(ts.Values) < 2 {
		return 0
	}
	mean := ts.Mean()
	var sum float64
	for _, v := range ts.Values {
		d := v - mean
		sum += d * d
	}
	return math.Sqrt(sum / float64(len(ts.Values)-1))
}

// Cumulative returns the cumulative sum.
func (ts *TimeSeries) Cumulative() []float64 {
	out := make([]float64, len(ts.Values))
	var cum float64
	for i, v := range ts.Values {
		cum += v
		out[i] = cum
	}
	return out
}

// MovingAverage computes a simple moving average with window size w.
func (ts *TimeSeries) MovingAverage(w int) []float64 {
	if w <= 0 || len(ts.Values) == 0 {
		return nil
	}
	n := len(ts.Values)
	out := make([]float64, n)
	for i := range out {
		start := i - w/2
		end := i + w/2 + 1
		if start < 0 {
			start = 0
		}
		if end > n {
			end = n
		}
		var sum float64
		for j := start; j < end; j++ {
			sum += ts.Values[j]
		}
		out[i] = sum / float64(end-start)
	}
	return out
}
