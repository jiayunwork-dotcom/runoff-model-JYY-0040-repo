package source

import (
	"testing"
)

func TestPenmanET(t *testing.T) {
	et := PenmanET(25, 60, 2.0, 20)
	if et <= 0 {
		t.Fatalf("expected positive ET, got %f", et)
	}
	if et > 20 {
		t.Fatalf("ET too high: %f", et)
	}
}

func TestHargreavesET(t *testing.T) {
	et := HargreavesET(15, 30, 22.5, 30)
	if et <= 0 {
		t.Fatalf("expected positive, got %f", et)
	}
}

func TestThornthwaiteET(t *testing.T) {
	I := HeatIndex([12]float64{5, 7, 10, 15, 20, 25, 28, 27, 22, 16, 10, 6})
	if I <= 0 {
		t.Fatal("heat index should be positive")
	}
	et := ThornthwaiteET(20, I, 12)
	if et <= 0 {
		t.Fatalf("expected positive ET, got %f", et)
	}
}

func TestHeatIndex(t *testing.T) {
	temps := [12]float64{2, 4, 8, 14, 20, 25, 28, 27, 22, 15, 8, 3}
	I := HeatIndex(temps)
	if I <= 0 {
		t.Fatal("expected positive heat index")
	}
}
