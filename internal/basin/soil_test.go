package basin

import "testing"

func TestSoilMoistureRunoff(t *testing.T) {
	sm := NewSoilMoisture(100, 90, 0.6, 0.1)
	runoff, _, _ := sm.Update(20, 2)
	if runoff <= 0 {
		t.Fatal("expected runoff when rain pushes past capacity")
	}
}

func TestSoilMoistureET(t *testing.T) {
	sm := NewSoilMoisture(100, 50, 0.6, 0.1)
	_, actualET, _ := sm.Update(0, 5)
	if actualET <= 0 {
		t.Fatal("expected some ET from available moisture")
	}
}

func TestSoilMoistureSaturation(t *testing.T) {
	sm := NewSoilMoisture(100, 50, 0.6, 0.1)
	if sm.Saturation() != 0.5 {
		t.Fatalf("expected 0.5, got %f", sm.Saturation())
	}
}

func TestSoilMoistureDeficit(t *testing.T) {
	sm := NewSoilMoisture(100, 30, 0.6, 0.1)
	if sm.Deficit() != 70 {
		t.Fatalf("expected 70, got %f", sm.Deficit())
	}
}
