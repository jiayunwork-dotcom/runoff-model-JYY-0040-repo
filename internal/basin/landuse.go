package basin

// LandUse represents a land use category with hydrologic parameters.
type LandUse struct {
	Name              string
	FractionArea      float64 // fraction of basin area
	CurveNumber       float64 // SCS curve number
	ManningN          float64 // Manning's roughness
	ImperviousFrac    float64 // impervious fraction
	InterceptionDepth float64 // mm
}

// CompositeCN computes a weighted composite curve number for multiple land uses.
func CompositeCN(landUses []LandUse) float64 {
	var weightedSum, totalArea float64
	for _, lu := range landUses {
		weightedSum += lu.CurveNumber * lu.FractionArea
		totalArea += lu.FractionArea
	}
	if totalArea == 0 {
		return 0
	}
	return weightedSum / totalArea
}

// CompositeImpervious computes the weighted impervious fraction.
func CompositeImpervious(landUses []LandUse) float64 {
	var weightedSum, totalArea float64
	for _, lu := range landUses {
		weightedSum += lu.ImperviousFrac * lu.FractionArea
		totalArea += lu.FractionArea
	}
	if totalArea == 0 {
		return 0
	}
	return weightedSum / totalArea
}

// TotalInterception computes weighted average interception depth.
func TotalInterception(landUses []LandUse) float64 {
	var weightedSum, totalArea float64
	for _, lu := range landUses {
		weightedSum += lu.InterceptionDepth * lu.FractionArea
		totalArea += lu.FractionArea
	}
	if totalArea == 0 {
		return 0
	}
	return weightedSum / totalArea
}

// DefaultLandUses returns typical land use categories.
func DefaultLandUses() []LandUse {
	return []LandUse{
		{Name: "Forest", FractionArea: 0.4, CurveNumber: 55, ManningN: 0.15, ImperviousFrac: 0, InterceptionDepth: 3},
		{Name: "Grassland", FractionArea: 0.3, CurveNumber: 65, ManningN: 0.035, ImperviousFrac: 0, InterceptionDepth: 1.5},
		{Name: "Cropland", FractionArea: 0.2, CurveNumber: 72, ManningN: 0.03, ImperviousFrac: 0.05, InterceptionDepth: 1},
		{Name: "Urban", FractionArea: 0.1, CurveNumber: 90, ManningN: 0.015, ImperviousFrac: 0.7, InterceptionDepth: 0.5},
	}
}
