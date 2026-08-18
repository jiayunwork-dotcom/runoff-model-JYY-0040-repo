package hydro

import "testing"

func TestFloodFrequency(t *testing.T) {
	peaks := []float64{100, 150, 200, 180, 120, 250, 90, 160, 210, 130}
	q100 := FloodFrequency(peaks, 100)
	q10 := FloodFrequency(peaks, 10)
	if q100 <= q10 {
		t.Fatal("Q100 should exceed Q10")
	}
	if q10 <= 0 {
		t.Fatal("Q10 should be positive")
	}
}

func TestLogPearsonIII(t *testing.T) {
	peaks := []float64{100, 150, 200, 180, 120, 250, 90, 160, 210, 130}
	q50 := LogPearsonIII(peaks, 50)
	if q50 <= 0 {
		t.Fatal("expected positive")
	}
}

func TestReturnPeriod(t *testing.T) {
	peaks := []float64{100, 150, 200, 180, 120}
	rp := ReturnPeriod(peaks, 200)
	if rp <= 1 {
		t.Fatalf("expected > 1, got %f", rp)
	}
}

func TestFloodPeakEstimate(t *testing.T) {
	q := FloodPeakEstimate(0.6, 50, 100)
	if q <= 0 {
		t.Fatal("expected positive flow")
	}
}
