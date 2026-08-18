package hydro

import (
	"testing"
)

func TestGreenAmpt(t *testing.T) {
	rain := []float64{50, 50, 50, 50, 50} // constant 50 mm/h
	F := GreenAmpt(10, 100, 0.3, rain, 1.0)
	if len(F) != 5 {
		t.Fatalf("expected 5, got %d", len(F))
	}
	// Cumulative infiltration should increase.
	for i := 1; i < len(F); i++ {
		if F[i] < F[i-1] {
			t.Fatal("cumulative should be non-decreasing")
		}
	}
}

func TestHortonInfiltration(t *testing.T) {
	rates := HortonInfiltration(50, 5, 0.5, 10, 1.0)
	if len(rates) != 10 {
		t.Fatalf("expected 10, got %d", len(rates))
	}
	// Rate should decrease over time.
	if rates[9] >= rates[0] {
		t.Fatal("Horton rate should decrease")
	}
	// Should approach fc.
	if rates[9] < 5 {
		t.Fatal("should not go below fc")
	}
}

func TestSCSCurveNumber(t *testing.T) {
	rain := []float64{10, 20, 30, 10, 5}
	runoff := SCSCurveNumber(75, rain)
	if len(runoff) != 5 {
		t.Fatalf("expected 5, got %d", len(runoff))
	}
	// Total runoff should be less than total rain.
	var totalR, totalP float64
	for i := range rain {
		totalP += rain[i]
		totalR += runoff[i]
	}
	if totalR >= totalP {
		t.Fatal("runoff should be less than rainfall")
	}
}

func TestPhilipInfiltration(t *testing.T) {
	F := PhilipInfiltration(20, 5, 5, 1.0)
	if len(F) != 5 {
		t.Fatalf("expected 5, got %d", len(F))
	}
	// Should be increasing.
	for i := 1; i < len(F); i++ {
		if F[i] <= F[i-1] {
			t.Fatal("Philip should be increasing")
		}
	}
}
