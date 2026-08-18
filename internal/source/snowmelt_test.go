package source

import "testing"

func TestSnowMelt(t *testing.T) {
	m := SnowMelt(5, 4, 0) // 5°C * 4 mm/°C/day = 20 mm
	if m != 20 {
		t.Fatalf("expected 20, got %f", m)
	}
	// Below threshold: no melt.
	if SnowMelt(-2, 4, 0) != 0 {
		t.Fatal("expected 0 below threshold")
	}
}

func TestSnowAccumulation(t *testing.T) {
	precip := []float64{10, 5, 0, 0, 0}
	temp := []float64{-5, -3, 2, 5, 8}
	snowpack, melt, _ := SnowAccumulation(precip, temp, 3, 0, 2)
	if len(snowpack) != 5 {
		t.Fatal("bad length")
	}
	// First 2 days accumulate snow.
	if snowpack[1] < 10 {
		t.Fatal("should have accumulated snow")
	}
	// Later days should melt.
	if melt[3] <= 0 {
		t.Fatal("expected melt on warm day")
	}
}

func TestSnowDensity(t *testing.T) {
	fresh := SnowDensity(0)
	if fresh != 50 {
		t.Fatalf("fresh snow density should be 50, got %f", fresh)
	}
	old := SnowDensity(100)
	if old <= fresh {
		t.Fatal("old snow should be denser")
	}
}
