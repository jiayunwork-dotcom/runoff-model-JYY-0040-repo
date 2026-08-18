package hydro

import (
	"math"
	"testing"
)

// TestNashUH checks the slice properties: deterministic length, all ordinates
// non-negative, and the sum normalized to ~1.
func TestNashUH(t *testing.T) {
	n, K, dt := 3, 12.0, 6.0
	uh := NashUH(n, K, dt)

	// Length rule: steps = 3*n*K/dt = 18 -> 19 ordinates.
	if len(uh) != 19 {
		t.Fatalf("expected length 19, got %d", len(uh))
	}
	var sum float64
	for _, v := range uh {
		if v < 0 {
			t.Fatalf("Nash UH ordinate must be >=0, got %v", v)
		}
		sum += v
	}
	if math.Abs(sum-1) > 1e-6 {
		t.Fatalf("Nash UH must sum to ~1, got %v", sum)
	}
}

// TestConvolve checks length = len(rain)+len(uh)-1 and non-negativity.
func TestConvolve(t *testing.T) {
	rain := []float64{0, 10, 20, 5, 0}
	uh := NashUH(3, 12, 6)
	drh := Convolve(rain, uh)
	if len(drh) != len(rain)+len(uh)-1 {
		t.Fatalf("expected length %d, got %d", len(rain)+len(uh)-1, len(drh))
	}
	for _, v := range drh {
		if v < 0 {
			t.Fatalf("convolution must be non-negative, got %v", v)
		}
	}
}

// TestMuskingum checks the outflow length equals the inflow length and stays
// non-negative for a non-negative inflow.
func TestMuskingum(t *testing.T) {
	inflow := []float64{0, 5, 20, 40, 25, 10, 3, 0}
	out := Muskingum(inflow, 12, 0.2, 6)
	if len(out) != len(inflow) {
		t.Fatalf("expected outflow length %d, got %d", len(inflow), len(out))
	}
	for _, v := range out {
		if v < -1e-9 {
			t.Fatalf("Muskingum outflow must be non-negative, got %v", v)
		}
	}
}

// TestPearsonIII checks that an above-mean value yields a return period > 1.
func TestPearsonIII(t *testing.T) {
	samples := []float64{100, 150, 200, 250, 300, 350, 400, 500}
	mean := 0.0
	for _, s := range samples {
		mean += s
	}
	mean /= float64(len(samples))

	above := mean + 50
	T := PearsonIII(samples, above)
	if T <= 1 {
		t.Fatalf("expected T>1 for above-mean value, got %v", T)
	}

	below := mean - 50
	T2 := PearsonIII(samples, below)
	if T2 >= 1 {
		t.Fatalf("expected T<1 for below-mean value, got %v", T2)
	}
}

// TestPearsonIIIEmpty verifies an empty sample never panics and returns 0.
func TestPearsonIIIEmpty(t *testing.T) {
	if got := PearsonIII(nil, 10); got != 0 {
		t.Fatalf("expected 0 for empty sample, got %v", got)
	}
}
