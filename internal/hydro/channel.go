package hydro

import "math"

// ManningFlow computes open channel flow velocity using Manning's equation.
// Returns velocity in m/s.
//
// Parameters:
//   - n: Manning's roughness coefficient
//   - R: hydraulic radius (m)
//   - S: channel slope (m/m)
func ManningFlow(n, R, S float64) float64 {
	if n <= 0 || R <= 0 || S <= 0 {
		return 0
	}
	return (1.0 / n) * math.Pow(R, 2.0/3.0) * math.Sqrt(S)
}

// ManningDischarge computes discharge Q = A * V.
func ManningDischarge(n, area, R, S float64) float64 {
	v := ManningFlow(n, R, S)
	return area * v
}

// TrapezoidalArea computes the cross-sectional area of a trapezoidal channel.
// b = bottom width, h = depth, z = side slope (horizontal:vertical).
func TrapezoidalArea(b, h, z float64) float64 {
	return (b + z*h) * h
}

// TrapezoidalPerimeter computes the wetted perimeter.
func TrapezoidalPerimeter(b, h, z float64) float64 {
	return b + 2*h*math.Sqrt(1+z*z)
}

// TrapezoidalHydraulicRadius computes R = A/P.
func TrapezoidalHydraulicRadius(b, h, z float64) float64 {
	A := TrapezoidalArea(b, h, z)
	P := TrapezoidalPerimeter(b, h, z)
	if P == 0 {
		return 0
	}
	return A / P
}

// KinematicWave routes flow using the kinematic wave approximation.
// Returns outflow at downstream end for each time step.
func KinematicWave(inflow []float64, length, slope, manning, width, dt float64) []float64 {
	if len(inflow) == 0 || length <= 0 || dt <= 0 {
		return nil
	}
	n := len(inflow)
	out := make([]float64, n)

	// Simplified: assume steady uniform flow for each step.
	alpha := math.Pow(manning*math.Sqrt(slope)/width, 0.6)
	if alpha <= 0 {
		copy(out, inflow)
		return out
	}
	celerity := length / (alpha * dt)
	if celerity > float64(n) {
		celerity = float64(n)
	}
	lag := int(math.Round(celerity))
	if lag < 0 {
		lag = 0
	}

	for i := 0; i < n; i++ {
		srcIdx := i - lag
		if srcIdx >= 0 && srcIdx < n {
			out[i] = inflow[srcIdx]
		}
	}
	return out
}

// TimeOfConcentration estimates tc using Kirpich formula.
// L = length of main channel (km), S = average slope (m/m).
// Returns tc in hours.
func TimeOfConcentration(L, S float64) float64 {
	if L <= 0 || S <= 0 {
		return 0
	}
	// Kirpich: tc = 0.0195 * L^0.77 * S^(-0.385) (L in meters)
	Lm := L * 1000
	return 0.0195 * math.Pow(Lm, 0.77) * math.Pow(S, -0.385) / 60 // convert min to hours
}
