package hydro

import "math"

// GreenAmpt computes infiltration using the Green-Ampt model.
// Returns cumulative infiltration F at each time step.
//
// Parameters:
//   - Ks: saturated hydraulic conductivity (mm/h)
//   - psi: wetting front suction head (mm)
//   - dtheta: moisture deficit (porosity - initial moisture)
//   - rain: rainfall intensity at each step (mm/h)
//   - dt: time step (hours)
func GreenAmpt(Ks, psi, dtheta float64, rain []float64, dt float64) []float64 {
	if len(rain) == 0 || Ks <= 0 || dt <= 0 {
		return nil
	}
	n := len(rain)
	F := make([]float64, n) // cumulative infiltration
	f := make([]float64, n) // infiltration rate

	var cumF float64
	for i := 0; i < n; i++ {
		// Potential infiltration rate.
		fp := Ks * (1 + psi*dtheta/maxF(cumF, 0.001))
		if rain[i] <= fp {
			// All rain infiltrates.
			f[i] = rain[i]
		} else {
			f[i] = fp
		}
		cumF += f[i] * dt
		F[i] = cumF
	}
	return F
}

// HortonInfiltration computes infiltration using Horton's equation.
// f(t) = fc + (f0 - fc) * exp(-k*t)
// Returns infiltration rate at each step.
func HortonInfiltration(f0, fc, k float64, n int, dt float64) []float64 {
	if n <= 0 {
		return nil
	}
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		t := float64(i) * dt
		out[i] = fc + (f0-fc)*math.Exp(-k*t)
	}
	return out
}

// PhilipInfiltration computes cumulative infiltration using Philip's equation:
// F(t) = S * sqrt(t) + A * t
// where S is sorptivity and A ≈ Ks.
func PhilipInfiltration(S, A float64, n int, dt float64) []float64 {
	if n <= 0 {
		return nil
	}
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		t := float64(i+1) * dt
		out[i] = S*math.Sqrt(t) + A*t
	}
	return out
}

// SCSCurveNumber computes direct runoff using the SCS-CN method.
// Q = (P - Ia)^2 / (P - Ia + S) when P > Ia, else Q = 0.
// S = 25400/CN - 254 (in mm), Ia = 0.2*S.
func SCSCurveNumber(CN float64, rainfall []float64) []float64 {
	if CN <= 0 || CN > 100 {
		return make([]float64, len(rainfall))
	}
	S := 25400.0/CN - 254.0
	Ia := 0.2 * S
	out := make([]float64, len(rainfall))
	var cumP, cumQ float64
	for i, p := range rainfall {
		cumP += p
		if cumP > Ia {
			pEff := cumP - Ia
			newCumQ := pEff * pEff / (pEff + S)
			out[i] = newCumQ - cumQ
			cumQ = newCumQ
		}
	}
	return out
}

func maxF(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
