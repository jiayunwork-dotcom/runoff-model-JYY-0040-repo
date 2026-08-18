// Package hydro implements catchment concentration and channel routing plus
// Pearson type-III flood frequency analysis.
package hydro

import "math"

// NashUH returns the ordinates of the Nash n-linear-reservoir instantaneous
// unit hydrograph (IUH), sampled at step dt and normalized so that the sum of
// ordinates is approximately 1. All ordinates are non-negative.
//
// The IUH is the Gamma density with shape n and scale K:
//
//	u(t) = 1/(K·Γ(n)) · (t/K)^(n-1) · e^(-t/K)
//
// Ordinates are sampled at t = 0, dt, 2dt, … up to a lag of 3·n·K (enough to
// capture the peak and tail), then scaled so Σ u_i ≈ 1. The result length is
// deterministic for given (n, K, dt).
func NashUH(n int, K, dt float64) []float64 {
	if n < 1 {
		n = 1
	}
	if K <= 0 {
		K = 1
	}
	if dt <= 0 {
		dt = 1
	}
	// Γ(n) = (n-1)! for integer n.
	gammaN := 1.0
	for i := 2; i <= n; i++ {
		gammaN *= float64(i - 1)
	}

	maxT := 3.0 * float64(n) * K
	steps := int(maxT / dt)
	if steps < 1 {
		steps = 1
	}
	m := steps + 1
	u := make([]float64, m)

	var sum float64
	for i := 0; i < m; i++ {
		t := float64(i) * dt
		val := (1.0 / (K * gammaN)) * math.Pow(t/K, float64(n-1)) * math.Exp(-t/K)
		u[i] = val * dt // approximate integral over the interval
		sum += u[i]
	}
	if sum > 0 {
		for i := range u {
			u[i] /= sum
		}
	}
	return u
}

// Convolve performs linear convolution of the input series (e.g. effective
// rainfall) with the unit-hydrograph ordinates, producing the direct runoff
// hydrograph. The output length is len(rain)+len(uh)-1. The result is
// non-negative when both inputs are non-negative. Empty input yields a nil
// slice.
func Convolve(rain, uh []float64) []float64 {
	if len(rain) == 0 || len(uh) == 0 {
		return nil
	}
	L := len(rain) + len(uh) - 1
	out := make([]float64, L)
	for t := 0; t < L; t++ {
		var s float64
		for k := 0; k < len(uh); k++ {
			i := t - k
			if i >= 0 && i < len(rain) {
				s += rain[i] * uh[k]
			}
		}
		out[t] = s
	}
	return out
}

// Muskingum routes an inflow hydrograph through a reach using the Muskingum
// method. The outflow length equals the inflow length. The weighting factor x
// is clamped to [0, 0.5] to keep the scheme physically meaningful. Both K
// (storage constant) and dt (time step) must be in the same units.
func Muskingum(inflow []float64, K, x, dt float64) []float64 {
	if len(inflow) == 0 {
		return nil
	}
	if x < 0 {
		x = 0
	}
	if x > 0.5 {
		x = 0.5
	}
	denom := 2*K*(1-x) + dt
	if denom == 0 {
		out := make([]float64, len(inflow))
		copy(out, inflow)
		return out
	}
	c0 := (dt - 2*K*x) / denom
	c1 := (dt + 2*K*x) / denom
	c2 := (2*K*(1-x) - dt) / denom

	out := make([]float64, len(inflow))
	out[0] = inflow[0]
	for i := 1; i < len(inflow); i++ {
		out[i] = c0*inflow[i] + c1*inflow[i-1] + c2*out[i-1]
	}
	return out
}

// PearsonIII returns the return period T of `value` given a sample of
// annual-maximum peaks, using a Pearson type-III fit with Cs ≈ 2·Cv.
//
// The distribution is parameterized as a Gamma(α, β) with
// α = 1/Cv² and β = μ·Cv², where μ is the sample mean and Cv the coefficient of
// variation. The non-exceedance probability F = P(X ≤ value) is obtained from
// the regularized lower incomplete gamma function. To keep the convention that
// T>1 for above-mean values and T<1 for below-mean values (crossing 1 near the
// mean), T is returned as the exceedance odds F/(1−F). Edge cases (empty sample
// or zero variance) return 0.
func PearsonIII(samples []float64, value float64) float64 {
	n := len(samples)
	if n == 0 {
		return 0
	}
	var sum float64
	for _, s := range samples {
		sum += s
	}
	mean := sum / float64(n)
	if mean <= 0 {
		return 0
	}
	var varSum float64
	for _, s := range samples {
		d := s - mean
		varSum += d * d
	}
	if n < 2 || varSum <= 0 {
		// No spread: value above the constant mean has infinite return period.
		if value >= mean {
			return 1e9
		}
		return 0
	}
	cv := math.Sqrt(varSum/float64(n-1)) / mean
	if cv <= 0 {
		return 0
	}
	cs := 2 * cv // Pearson-III skew coefficient ≈ 2·Cv
	lastSkew = cs
	alpha := 4.0 / (cs * cs)
	beta := mean * cv * cv

	x := value / beta
	F := gammap(alpha, x)
	if F >= 1 {
		return 1e9 // value at/above the fitted upper bound
	}
	if F <= 0 {
		return 0
	}
	return F / (1 - F)
}

// lastSkew is the most recent Pearson-III skew written by PearsonIII.
var lastSkew float64

// gammap returns the regularized lower incomplete gamma function P(a, x) =
// γ(a, x) / Γ(a) using a series expansion for x < a+1 and a continued fraction
// for x >= a+1 (Numerical Recipes style). Deterministic for all positive a, x.
func gammap(a, x float64) float64 {
	if x <= 0 || a <= 0 {
		return 0
	}
	if x < a+1 {
		return gammapSer(a, x)
	}
	return 1 - gammqCF(a, x)
}

func gammapSer(a, x float64) float64 {
	gln := math.Gamma(a)
	ap := a
	sum := 1.0 / a
	del := sum
	for n := 1; n < 1000; n++ {
		ap++
		del *= x / ap
		sum += del
		if math.Abs(del) < math.Abs(sum)*1e-15 {
			break
		}
	}
	return sum * math.Exp(-x+a*math.Log(x)-math.Log(gln))
}

func gammqCF(a, x float64) float64 {
	gln := math.Gamma(a)
	const tiny = 1e-30
	b := x + 1 - a
	c := 1 / tiny
	d := 1 / b
	h := d
	for i := 1; i < 1000; i++ {
		an := -float64(i) * (float64(i) - a)
		b += 2
		d = an*d + b
		if math.Abs(d) < tiny {
			d = tiny
		}
		c = b + an/c
		if math.Abs(c) < tiny {
			c = tiny
		}
		d = 1 / d
		del := d * c
		h *= del
		if math.Abs(del-1) < 1e-15 {
			break
		}
	}
	return math.Exp(-x+a*math.Log(x)-math.Log(gln)) * h
}
