package hydro

import "math"

// LinearReservoir routes inflow through a single linear reservoir with
// storage constant K. The outflow at step i is:
// Q(i) = Q(i-1) * exp(-dt/K) + I(i) * (1 - exp(-dt/K))
func LinearReservoir(inflow []float64, K, dt float64) []float64 {
	if len(inflow) == 0 || K <= 0 || dt <= 0 {
		return nil
	}
	c := math.Exp(-dt / K)
	out := make([]float64, len(inflow))
	out[0] = inflow[0] * (1 - c)
	for i := 1; i < len(inflow); i++ {
		out[i] = out[i-1]*c + inflow[i]*(1-c)
	}
	return out
}

// CascadeReservoirs routes inflow through n cascaded linear reservoirs.
func CascadeReservoirs(inflow []float64, n int, K, dt float64) []float64 {
	if n <= 0 {
		return inflow
	}
	current := make([]float64, len(inflow))
	copy(current, inflow)
	for i := 0; i < n; i++ {
		current = LinearReservoir(current, K, dt)
	}
	return current
}

// StorageIndication routes inflow through a reservoir using the storage
// indication (modified Puls) method with a given storage-outflow relationship.
// The relationship is given as (storage, outflow) pairs sorted by storage.
func StorageIndication(inflow []float64, storage, outflow []float64, dt float64) []float64 {
	if len(inflow) == 0 || len(storage) != len(outflow) || len(storage) < 2 {
		return nil
	}

	// Build 2S/dt + O curve.
	indicator := make([]float64, len(storage))
	for i := range storage {
		indicator[i] = 2*storage[i]/dt + outflow[i]
	}

	n := len(inflow)
	out := make([]float64, n)
	var prevInd float64
	for i := 0; i < n; i++ {
		var rhs float64
		if i == 0 {
			rhs = inflow[0]
		} else {
			rhs = inflow[i] + inflow[i-1] + prevInd - 2*out[i-1]
		}
		// Interpolate outflow from indicator.
		out[i] = interpOutflow(rhs, indicator, outflow)
		prevInd = rhs
	}
	return out
}

func interpOutflow(ind float64, indicator, outflow []float64) float64 {
	if ind <= indicator[0] {
		return outflow[0]
	}
	for i := 1; i < len(indicator); i++ {
		if ind <= indicator[i] {
			frac := (ind - indicator[i-1]) / (indicator[i] - indicator[i-1])
			return outflow[i-1] + frac*(outflow[i]-outflow[i-1])
		}
	}
	return outflow[len(outflow)-1]
}

// BaseflowSeparation separates baseflow from total flow using a digital filter.
// alpha is typically 0.925. Returns (baseflow, quickflow).
func BaseflowSeparation(totalFlow []float64, alpha float64) ([]float64, []float64) {
	n := len(totalFlow)
	if n == 0 {
		return nil, nil
	}
	quick := make([]float64, n)
	base := make([]float64, n)
	quick[0] = totalFlow[0]
	for i := 1; i < n; i++ {
		quick[i] = alpha*quick[i-1] + (1+alpha)/2*(totalFlow[i]-totalFlow[i-1])
		if quick[i] < 0 {
			quick[i] = 0
		}
		if quick[i] > totalFlow[i] {
			quick[i] = totalFlow[i]
		}
		base[i] = totalFlow[i] - quick[i]
	}
	return base, quick
}
