package basin

import "math"

// Calibration holds calibration results.
type Calibration struct {
	BestParams Basin
	NSE        float64
	RMSE       float64
	Iterations int
}

// SimpleCalibrate performs a simple grid search to find optimal basin parameters.
// It varies Wm and B within given ranges to maximize NSE.
func SimpleCalibrate(records []Record, wmRange [2]float64, bRange [2]float64, c float64, observed []float64, productionFn func([]Record, float64, float64, float64) []float64) Calibration {
	bestNSE := -math.MaxFloat64
	var bestWm, bestB float64
	iterations := 0

	wmStep := (wmRange[1] - wmRange[0]) / 10
	bStep := (bRange[1] - bRange[0]) / 10

	for wm := wmRange[0]; wm <= wmRange[1]; wm += wmStep {
		for b := bRange[0]; b <= bRange[1]; b += bStep {
			iterations++
			simulated := productionFn(records, wm, b, c)
			nse := computeNSE(observed, simulated)
			if nse > bestNSE {
				bestNSE = nse
				bestWm = wm
				bestB = b
			}
		}
	}

	rmse := computeRMSE(observed, productionFn(records, bestWm, bestB, c))
	return Calibration{
		BestParams: Basin{WM: bestWm, B: bestB, C: c},
		NSE:        bestNSE,
		RMSE:       rmse,
		Iterations: iterations,
	}
}

func computeNSE(obs, sim []float64) float64 {
	n := len(obs)
	if n == 0 || len(sim) < n {
		return -math.MaxFloat64
	}
	var sumObs float64
	for _, o := range obs[:n] {
		sumObs += o
	}
	mean := sumObs / float64(n)
	var num, denom float64
	for i := 0; i < n; i++ {
		d := obs[i] - sim[i]
		num += d * d
		dm := obs[i] - mean
		denom += dm * dm
	}
	if denom == 0 {
		return 0
	}
	return 1 - num/denom
}

func computeRMSE(obs, sim []float64) float64 {
	n := len(obs)
	if n == 0 || len(sim) < n {
		return 0
	}
	var sum float64
	for i := 0; i < n; i++ {
		d := obs[i] - sim[i]
		sum += d * d
	}
	return math.Sqrt(sum / float64(n))
}
