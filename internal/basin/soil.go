package basin

import "math"

// SoilMoistureModel tracks the soil moisture accounting over time.
type SoilMoistureModel struct {
	WM     float64 // max storage capacity (mm)
	W      float64 // current storage (mm)
	FC     float64 // field capacity fraction
	WP     float64 // wilting point fraction
}

// NewSoilMoisture creates a model with initial condition.
func NewSoilMoisture(wm, initialW, fc, wp float64) *SoilMoistureModel {
	return &SoilMoistureModel{WM: wm, W: initialW, FC: fc, WP: wp}
}

// Update applies one time step of rainfall and ET.
// Returns (runoff, actual_ET, drainage).
func (sm *SoilMoistureModel) Update(rain, petET float64) (float64, float64, float64) {
	// Add rain.
	sm.W += rain

	// Compute actual ET (limited by available moisture above wilting point).
	available := sm.W - sm.WP*sm.WM
	if available < 0 {
		available = 0
	}
	actualET := math.Min(petET, available)
	sm.W -= actualET

	// Surface runoff if exceeds capacity.
	var runoff float64
	if sm.W > sm.WM {
		runoff = sm.W - sm.WM
		sm.W = sm.WM
	}

	// Drainage (percolation) when above field capacity.
	var drainage float64
	fcStorage := sm.FC * sm.WM
	if sm.W > fcStorage {
		drainage = (sm.W - fcStorage) * 0.1 // simplified drainage rate
		sm.W -= drainage
	}

	return runoff, actualET, drainage
}

// Saturation returns the current saturation ratio.
func (sm *SoilMoistureModel) Saturation() float64 {
	if sm.WM <= 0 {
		return 0
	}
	return sm.W / sm.WM
}

// Storage returns current storage in mm.
func (sm *SoilMoistureModel) Storage() float64 {
	return sm.W
}

// Deficit returns the current moisture deficit (mm below capacity).
func (sm *SoilMoistureModel) Deficit() float64 {
	return sm.WM - sm.W
}

// Reset sets storage to a given value.
func (sm *SoilMoistureModel) Reset(w float64) {
	sm.W = w
}
