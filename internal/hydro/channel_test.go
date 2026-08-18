package hydro

import "testing"

func TestManningFlow(t *testing.T) {
	v := ManningFlow(0.03, 1.5, 0.001)
	if v <= 0 {
		t.Fatal("expected positive velocity")
	}
}

func TestTrapezoidalArea(t *testing.T) {
	a := TrapezoidalArea(5, 2, 1.5)
	// A = (5 + 1.5*2) * 2 = 8 * 2 = 16
	if a != 16 {
		t.Fatalf("expected 16, got %f", a)
	}
}

func TestTrapezoidalHydraulicRadius(t *testing.T) {
	r := TrapezoidalHydraulicRadius(5, 2, 1.5)
	if r <= 0 {
		t.Fatal("expected positive R")
	}
}

func TestTimeOfConcentration(t *testing.T) {
	tc := TimeOfConcentration(10, 0.01)
	if tc <= 0 {
		t.Fatal("expected positive tc")
	}
}
