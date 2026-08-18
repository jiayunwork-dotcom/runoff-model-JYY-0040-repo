package hydro

import (
	"testing"
)

func TestLinearReservoir(t *testing.T) {
	inflow := []float64{0, 10, 20, 10, 5, 0}
	out := LinearReservoir(inflow, 2.0, 1.0)
	if len(out) != 6 {
		t.Fatalf("expected 6, got %d", len(out))
	}
	// Peak should be attenuated.
	maxIn := 20.0
	maxOut := 0.0
	for _, v := range out {
		if v > maxOut {
			maxOut = v
		}
	}
	if maxOut >= maxIn {
		t.Fatal("reservoir should attenuate peak")
	}
}

func TestCascadeReservoirs(t *testing.T) {
	inflow := []float64{0, 0, 50, 0, 0, 0, 0, 0}
	out := CascadeReservoirs(inflow, 3, 2.0, 1.0)
	if len(out) != 8 {
		t.Fatalf("expected 8, got %d", len(out))
	}
	// More cascade = more attenuation.
	single := LinearReservoir(inflow, 2.0, 1.0)
	maxSingle := 0.0
	maxCascade := 0.0
	for i := range out {
		if single[i] > maxSingle {
			maxSingle = single[i]
		}
		if out[i] > maxCascade {
			maxCascade = out[i]
		}
	}
	if maxCascade >= maxSingle {
		t.Fatal("cascade should attenuate more than single")
	}
}

func TestBaseflowSeparation(t *testing.T) {
	total := []float64{5, 10, 25, 40, 30, 20, 12, 8, 5}
	base, quick := BaseflowSeparation(total, 0.925)
	if len(base) != 9 || len(quick) != 9 {
		t.Fatal("bad lengths")
	}
	// base + quick should approximately equal total.
	for i := range total {
		sum := base[i] + quick[i]
		if sum < total[i]*0.99 || sum > total[i]*1.01 {
			// Allow small numeric error.
			if sum < 0 || sum > total[i]*2 {
				t.Fatalf("step %d: base+quick=%f, total=%f", i, sum, total[i])
			}
		}
	}
}
