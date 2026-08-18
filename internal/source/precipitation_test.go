package source

import "testing"

func TestDesignStorm(t *testing.T) {
	storm := DesignStorm(100, 12)
	if len(storm) != 12 {
		t.Fatalf("expected 12, got %d", len(storm))
	}
	// Total should approximately equal input.
	var total float64
	for _, v := range storm {
		total += v
	}
	if total < 95 || total > 105 {
		t.Fatalf("total=%f, expected ~100", total)
	}
	// Peak should be near center.
	maxIdx := 0
	maxVal := 0.0
	for i, v := range storm {
		if v > maxVal {
			maxVal = v
			maxIdx = i
		}
	}
	if maxIdx < 3 || maxIdx > 8 {
		t.Fatalf("peak at %d, expected near center", maxIdx)
	}
}

func TestArealReduction(t *testing.T) {
	point := 50.0
	areal := ArealReduction(point, 500, 6)
	if areal >= point {
		t.Fatal("areal should be less than point")
	}
	if areal <= 0 {
		t.Fatal("should be positive")
	}
}

func TestIDF(t *testing.T) {
	intensity := IDF(1.0, 100, 1000, 0.5, 0.7)
	if intensity <= 0 {
		t.Fatal("expected positive intensity")
	}
}
