package basin

import "testing"

func TestTimeSeriesStats(t *testing.T) {
	ts := NewTimeSeries([]float64{1, 2, 3, 4, 5}, 1)
	if ts.Mean() != 3 {
		t.Fatalf("expected mean=3, got %f", ts.Mean())
	}
	if ts.Sum() != 15 {
		t.Fatalf("expected sum=15, got %f", ts.Sum())
	}
	max, idx := ts.Max()
	if max != 5 || idx != 4 {
		t.Fatalf("max=%f at %d", max, idx)
	}
}

func TestCumulative(t *testing.T) {
	ts := NewTimeSeries([]float64{1, 2, 3}, 1)
	cum := ts.Cumulative()
	if cum[2] != 6 {
		t.Fatalf("expected 6, got %f", cum[2])
	}
}

func TestMovingAverage(t *testing.T) {
	ts := NewTimeSeries([]float64{1, 3, 5, 7, 9}, 1)
	ma := ts.MovingAverage(3)
	if len(ma) != 5 {
		t.Fatalf("expected 5, got %d", len(ma))
	}
}
