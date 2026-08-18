package source

import (
	"math"
	"sort"
)

// DesignStorm generates a design storm hyetograph using the alternating block method.
// totalDepth is the total rainfall (mm), duration is number of time steps,
// returnPeriod controls intensity via IDF parameters.
func DesignStorm(totalDepth float64, duration int) []float64 {
	if duration <= 0 || totalDepth <= 0 {
		return nil
	}

	// Generate incremental depths using SCS Type-II distribution (simplified).
	cumulative := make([]float64, duration)
	for i := range cumulative {
		t := float64(i+1) / float64(duration)
		cumulative[i] = totalDepth * scsTypeII(t)
	}

	// Compute incremental.
	incremental := make([]float64, duration)
	incremental[0] = cumulative[0]
	for i := 1; i < duration; i++ {
		incremental[i] = cumulative[i] - cumulative[i-1]
	}

	// Alternating block: sort descending then place center-outward.
	sort.Sort(sort.Reverse(sort.Float64Slice(incremental)))
	result := make([]float64, duration)
	mid := duration / 2
	left := mid - 1
	right := mid + 1
	result[mid] = incremental[0]
	for i := 1; i < duration; i++ {
		if i%2 == 1 && right < duration {
			result[right] = incremental[i]
			right++
		} else if left >= 0 {
			result[left] = incremental[i]
			left--
		} else if right < duration {
			result[right] = incremental[i]
			right++
		}
	}
	return result
}

// scsTypeII is a simplified SCS Type II cumulative distribution.
func scsTypeII(t float64) float64 {
	// Approximation using logistic function centered at t=0.5.
	return 1.0 / (1.0 + math.Exp(-12*(t-0.5)))
}

// ArealReduction applies an areal reduction factor (ARF) to point rainfall.
// area is in km², duration is in hours.
func ArealReduction(pointRain float64, area, duration float64) float64 {
	if area <= 0 || duration <= 0 {
		return pointRain
	}
	// Empirical ARF formula (simplified Bell method).
	arf := 1.0 - 0.12*math.Pow(area, 0.4)*math.Pow(duration, -0.3)
	if arf < 0.3 {
		arf = 0.3
	}
	if arf > 1.0 {
		arf = 1.0
	}
	return pointRain * arf
}

// IDF returns rainfall intensity (mm/h) for given duration (hours) and
// return period (years) using a simplified IDF formula:
// i = a / (duration + b)^c
func IDF(duration, returnPeriod, a, b, c float64) float64 {
	if duration <= 0 {
		return 0
	}
	return a * math.Pow(returnPeriod, 0.2) / math.Pow(duration+b, c)
}
