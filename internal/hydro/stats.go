package hydro

import (
	"math"
	"sort"
)

// HydrographStats computes statistics of a hydrograph.
type HydrographStats struct {
	PeakFlow  float64
	PeakTime  int
	Volume    float64
	MeanFlow  float64
	BaseTime  int // duration above threshold
}

// ComputeStats computes stats of a flow series.
func ComputeStats(flows []float64, threshold float64) HydrographStats {
	if len(flows) == 0 {
		return HydrographStats{}
	}
	var s HydrographStats
	var sum float64
	for i, f := range flows {
		sum += f
		if f > s.PeakFlow {
			s.PeakFlow = f
			s.PeakTime = i
		}
		if f > threshold {
			s.BaseTime++
		}
	}
	s.Volume = sum
	s.MeanFlow = sum / float64(len(flows))
	return s
}

// NSE computes the Nash-Sutcliffe Efficiency between observed and simulated.
// NSE = 1 means perfect match, NSE < 0 means worse than using the mean.
func NSE(observed, simulated []float64) float64 {
	n := len(observed)
	if n == 0 || len(simulated) < n {
		return 0
	}
	var sumObs float64
	for _, o := range observed[:n] {
		sumObs += o
	}
	mean := sumObs / float64(n)

	var num, denom float64
	for i := 0; i < n; i++ {
		diff := observed[i] - simulated[i]
		num += diff * diff
		devMean := observed[i] - mean
		denom += devMean * devMean
	}
	if denom == 0 {
		return 0
	}
	return 1 - num/denom
}

// RMSE computes the root mean square error.
func RMSE(observed, simulated []float64) float64 {
	n := len(observed)
	if n == 0 || len(simulated) < n {
		return 0
	}
	var sum float64
	for i := 0; i < n; i++ {
		d := observed[i] - simulated[i]
		sum += d * d
	}
	return math.Sqrt(sum / float64(n))
}

// RelativeError computes |sim - obs| / obs for peak.
func RelativeError(obsPeak, simPeak float64) float64 {
	if obsPeak == 0 {
		return 0
	}
	return math.Abs(simPeak-obsPeak) / obsPeak
}

// Percentile returns the p-th percentile of values (p in [0,100]).
func Percentile(values []float64, p float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)
	idx := p / 100 * float64(len(sorted)-1)
	lo := int(idx)
	hi := lo + 1
	if hi >= len(sorted) {
		return sorted[lo]
	}
	frac := idx - float64(lo)
	return sorted[lo]*(1-frac) + sorted[hi]*frac
}

// FlowDurationCurve returns (exceedance_fraction, flow) pairs sorted by descending flow.
func FlowDurationCurve(flows []float64) ([]float64, []float64) {
	sorted := make([]float64, len(flows))
	copy(sorted, flows)
	sort.Sort(sort.Reverse(sort.Float64Slice(sorted)))
	n := float64(len(sorted))
	exceed := make([]float64, len(sorted))
	for i := range exceed {
		exceed[i] = float64(i+1) / (n + 1)
	}
	return exceed, sorted
}
