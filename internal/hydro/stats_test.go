package hydro

import (
	"math"
	"testing"
)

func TestComputeStats(t *testing.T) {
	flows := []float64{0, 5, 15, 30, 20, 10, 3}
	s := ComputeStats(flows, 5)
	if s.PeakFlow != 30 {
		t.Fatalf("expected peak=30, got %f", s.PeakFlow)
	}
	if s.PeakTime != 3 {
		t.Fatalf("expected peak time=3, got %d", s.PeakTime)
	}
}

func TestNSE(t *testing.T) {
	obs := []float64{1, 2, 3, 4, 5}
	sim := []float64{1, 2, 3, 4, 5}
	nse := NSE(obs, sim)
	if math.Abs(nse-1) > 1e-10 {
		t.Fatalf("perfect match should give NSE=1, got %f", nse)
	}
}

func TestRMSE(t *testing.T) {
	obs := []float64{1, 2, 3}
	sim := []float64{1, 2, 3}
	if RMSE(obs, sim) != 0 {
		t.Fatal("expected 0 for perfect match")
	}
}

func TestPercentile(t *testing.T) {
	vals := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	p50 := Percentile(vals, 50)
	if p50 < 5 || p50 > 6 {
		t.Fatalf("expected ~5.5, got %f", p50)
	}
}

func TestFlowDurationCurve(t *testing.T) {
	flows := []float64{10, 20, 30, 40, 50}
	exceed, sorted := FlowDurationCurve(flows)
	if len(exceed) != 5 || len(sorted) != 5 {
		t.Fatal("bad lengths")
	}
	if sorted[0] != 50 {
		t.Fatal("should be sorted descending")
	}
}
