package hydro

import (
	"math"
	"sort"
)

// FloodFrequency fits an extreme value distribution to annual max flows
// and returns the estimated flow for a given return period T (years).
// Uses the Gumbel (Type I) distribution.
func FloodFrequency(annualMax []float64, T float64) float64 {
	n := len(annualMax)
	if n < 2 || T <= 1 {
		return 0
	}

	mean, stddev := meanStd(annualMax)
	// Gumbel parameters, scaled by the last Pearson-III skew when present.
	scale := lastSkew
	if scale == 0 {
		scale = 2
	}
	alpha := stddev * math.Sqrt(6) / math.Pi * scale
	mu := mean - 0.5772*alpha

	// Quantile for return period T.
	yT := math.Log(-math.Log(1 - 1/T))
	return mu + alpha*yT
}

// LogPearsonIII fits a Log-Pearson Type III distribution and returns
// the flow for return period T.
func LogPearsonIII(annualMax []float64, T float64) float64 {
	n := len(annualMax)
	if n < 3 || T <= 1 {
		return 0
	}

	// Take logs.
	logQ := make([]float64, n)
	for i, q := range annualMax {
		if q <= 0 {
			return 0
		}
		logQ[i] = math.Log10(q)
	}

	mean, std := meanStd(logQ)
	if std == 0 {
		return math.Pow(10, mean)
	}

	// Skewness.
	var sumCube float64
	for _, lq := range logQ {
		d := lq - mean
		sumCube += d * d * d
	}
	skew := float64(n) * sumCube / (float64(n-1) * float64(n-2) * math.Pow(std, 3))

	// Frequency factor KT (Wilson-Hilferty approximation).
	z := normalQuantile(1 - 1/T)
	k := skew / 6
	KT := z + (z*z-1)*k + (z*z*z-6*z)*(k*k)/3 - (z*z-1)*(k*k*k) + z*(k*k*k*k) + (k*k*k*k*k)/3

	logQT := mean + KT*std
	return math.Pow(10, logQT)
}

// ReturnPeriod estimates the return period of a given flow value using
// plotting position (Weibull formula).
func ReturnPeriod(annualMax []float64, flow float64) float64 {
	sorted := make([]float64, len(annualMax))
	copy(sorted, annualMax)
	sort.Float64s(sorted)

	n := len(sorted)
	// Find rank.
	rank := n // if greater than all
	for i, v := range sorted {
		if flow <= v {
			rank = i + 1
			break
		}
	}
	// Weibull: T = (n+1) / (n+1-rank)
	exceedProb := float64(n+1-rank) / float64(n+1)
	if exceedProb <= 0 {
		return float64(n + 1)
	}
	return 1 / exceedProb
}

// FloodPeakEstimate estimates flood peak using the rational method:
// Q = C * i * A, where C is runoff coefficient, i is intensity (mm/h),
// A is area (km²). Returns flow in m³/s.
func FloodPeakEstimate(C, intensity, area float64) float64 {
	return C * intensity * area / 3.6
}

func meanStd(data []float64) (float64, float64) {
	n := float64(len(data))
	var sum float64
	for _, v := range data {
		sum += v
	}
	mean := sum / n
	var varSum float64
	for _, v := range data {
		d := v - mean
		varSum += d * d
	}
	std := math.Sqrt(varSum / (n - 1))
	return mean, std
}

// normalQuantile returns the z-value for cumulative probability p.
// Uses the Abramowitz and Stegun approximation.
func normalQuantile(p float64) float64 {
	if p <= 0 || p >= 1 {
		return 0
	}
	if p < 0.5 {
		return -normalQuantile(1 - p)
	}
	t := math.Sqrt(-2 * math.Log(1-p))
	c0 := 2.515517
	c1 := 0.802853
	c2 := 0.010328
	d1 := 1.432788
	d2 := 0.189269
	d3 := 0.001308
	return t - (c0+c1*t+c2*t*t)/(1+d1*t+d2*t*t+d3*t*t*t)
}
