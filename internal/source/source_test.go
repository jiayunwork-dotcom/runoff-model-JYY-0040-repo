package source

import (
	"math"
	"testing"

	"runoff-model/internal/basin"
)

// TestProductionNil verifies a nil input returns a nil slice (no panic).
func TestProductionNil(t *testing.T) {
	got := Production(nil, 120, 0.3, 1.0)
	if got != nil {
		t.Fatalf("expected nil for nil input, got len %d", len(got))
	}
}

// TestProductionZeroRain verifies that when net rainfall pe<=0 the runoff is 0,
// that 0 <= R <= pe always holds, and that R is monotonic with pe.
func TestProductionZeroRain(t *testing.T) {
	recs := []basin.Record{
		{Rain: 0, ET: 5},   // pe = -5 -> R = 0
		{Rain: 5, ET: 5},   // pe = 0  -> R = 0
		{Rain: 20, ET: 5},  // pe = 15
		{Rain: 60, ET: 5},  // pe = 55
		{Rain: 200, ET: 5}, // pe = 195 (>= WMM=120 -> R == pe)
	}
	const wm, b, c = 120.0, 0.3, 1.0
	R := Production(recs, wm, b, c)

	if R[0] != 0 || R[1] != 0 {
		t.Fatalf("expected 0 runoff when pe<=0, got R[0]=%v R[1]=%v", R[0], R[1])
	}

	prevPE := 0.0
	for i, rec := range recs {
		pe := rec.Rain - rec.ET
		if pe < 0 {
			pe = 0
		}
		if R[i] < 0 {
			t.Fatalf("record %d: R must be >=0, got %v", i, R[i])
		}
		if R[i] > pe+1e-9 {
			t.Fatalf("record %d: R (%v) must be <= pe (%v)", i, R[i], pe)
		}
		if i > 1 && R[i] < R[i-1]-1e-9 {
			t.Fatalf("record %d: R must be monotonic with pe, got %v < %v", i, R[i], R[i-1])
		}
		prevPE = pe
	}
	_ = prevPE

	// When pe >= WMM, R must equal pe exactly.
	if math.Abs(R[4]-(200-5)) > 1e-9 {
		t.Fatalf("expected R==pe when pe>=WMM, got R[4]=%v want %v", R[4], 195.0)
	}
}
