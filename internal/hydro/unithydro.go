package hydro

import "math"

// SCSTriangularUH generates a triangular unit hydrograph using SCS methods.
// tc = time of concentration (hours), tp = time to peak, dt = time step.
// Returns ordinates in m³/s per mm of excess rainfall per km².
func SCSTriangularUH(area, tc, dt float64) []float64 {
	if area <= 0 || tc <= 0 || dt <= 0 {
		return nil
	}
	// SCS relationships.
	tLag := 0.6 * tc
	tp := dt/2 + tLag
	tb := 2.67 * tp // base time
	qp := 0.208 * area / tp // peak discharge per mm

	// Generate ordinates.
	n := int(math.Ceil(tb/dt)) + 1
	uh := make([]float64, n)
	for i := 0; i < n; i++ {
		t := float64(i) * dt
		if t <= tp {
			uh[i] = qp * t / tp
		} else if t <= tb {
			uh[i] = qp * (tb - t) / (tb - tp)
		}
		if uh[i] < 0 {
			uh[i] = 0
		}
	}

	// Normalize so sum * dt ≈ 1 mm equivalent.
	var sum float64
	for _, v := range uh {
		sum += v * dt
	}
	if sum > 0 {
		for i := range uh {
			uh[i] /= sum
		}
	}
	return uh
}

// ClarkUH generates a Clark unit hydrograph using time-area and linear reservoir.
// timeArea is the fractional cumulative time-area curve (0 to 1 over duration).
// K = storage constant, dt = time step.
func ClarkUH(timeArea []float64, K, dt float64) []float64 {
	if len(timeArea) == 0 || K <= 0 || dt <= 0 {
		return nil
	}

	// Compute incremental time-area histogram.
	n := len(timeArea)
	incr := make([]float64, n)
	incr[0] = timeArea[0]
	for i := 1; i < n; i++ {
		incr[i] = timeArea[i] - timeArea[i-1]
		if incr[i] < 0 {
			incr[i] = 0
		}
	}

	// Route through linear reservoir.
	c := math.Exp(-dt / K)
	uh := make([]float64, n*2)
	for i := 0; i < len(uh); i++ {
		var input float64
		if i < n {
			input = incr[i]
		}
		if i == 0 {
			uh[0] = input * (1 - c)
		} else {
			uh[i] = uh[i-1]*c + input*(1-c)
		}
	}

	// Trim trailing near-zero values.
	cutoff := 0.001 * uh[0]
	end := len(uh)
	for end > 0 && uh[end-1] < cutoff {
		end--
	}
	if end < 1 {
		end = 1
	}

	// Normalize.
	result := uh[:end]
	var sum float64
	for _, v := range result {
		sum += v
	}
	if sum > 0 {
		for i := range result {
			result[i] /= sum
		}
	}
	return result
}

// SnyderUH generates a Snyder synthetic unit hydrograph.
// area in km², L in km (length of main channel), Lca in km (to centroid),
// Ct and Cp are regional coefficients.
func SnyderUH(area, L, Lca, Ct, Cp, dt float64) []float64 {
	if area <= 0 || L <= 0 || Lca <= 0 || dt <= 0 {
		return nil
	}
	// Standard lag.
	tpR := Ct * math.Pow(L*Lca, 0.3)
	// Standard peak.
	qpR := Cp * 2.78 * area / tpR

	// Adjust for non-standard duration.
	tr := tpR / 5.5
	tp := tpR + (dt-tr)/4
	qp := qpR * tpR / tp

	// Generate triangular approximation.
	tb := 5.5 * tp
	n := int(math.Ceil(tb/dt)) + 1
	uh := make([]float64, n)
	for i := 0; i < n; i++ {
		t := float64(i) * dt
		if t <= tp {
			uh[i] = qp * t / tp
		} else {
			uh[i] = qp * math.Exp(-1.5 * (t - tp) / (tb - tp))
		}
	}

	// Normalize.
	var sum float64
	for _, v := range uh {
		sum += v * dt
	}
	if sum > 0 {
		for i := range uh {
			uh[i] /= sum
		}
	}
	return uh
}
