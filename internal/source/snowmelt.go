package source

import "math"

// SnowMelt computes daily snowmelt using the degree-day method.
// Returns melt in mm/day.
//
// Parameters:
//   - temp: daily mean temperature (°C)
//   - meltFactor: degree-day factor (mm/°C/day), typically 2-6
//   - baseTemp: temperature threshold for melt (°C), typically 0
func SnowMelt(temp, meltFactor, baseTemp float64) float64 {
	if temp <= baseTemp {
		return 0
	}
	return meltFactor * (temp - baseTemp)
}

// SnowAccumulation models snow pack over a time series.
// Returns (snowpack, melt, rain_on_snow) at each step.
func SnowAccumulation(precip, temp []float64, meltFactor, snowTemp, rainTemp float64) ([]float64, []float64, []float64) {
	n := len(precip)
	if n == 0 || len(temp) < n {
		return nil, nil, nil
	}

	snowpack := make([]float64, n)
	melt := make([]float64, n)
	rainOnSnow := make([]float64, n)

	var swe float64 // snow water equivalent
	for i := 0; i < n; i++ {
		// Precipitation partitioning.
		if temp[i] <= snowTemp {
			swe += precip[i] // all snow
		} else if temp[i] >= rainTemp {
			// All rain.
			if swe > 0 {
				rainOnSnow[i] = precip[i]
			}
		} else {
			// Mixed: linear interpolation.
			snowFrac := (rainTemp - temp[i]) / (rainTemp - snowTemp)
			swe += precip[i] * snowFrac
			if swe > 0 {
				rainOnSnow[i] = precip[i] * (1 - snowFrac)
			}
		}

		// Melt.
		dailyMelt := SnowMelt(temp[i], meltFactor, 0)
		if dailyMelt > swe {
			dailyMelt = swe
		}
		swe -= dailyMelt
		melt[i] = dailyMelt
		snowpack[i] = swe
	}
	return snowpack, melt, rainOnSnow
}

// SnowDensity estimates snow density based on age (days since snowfall).
// Fresh snow: ~50 kg/m³, old snow: ~300-500 kg/m³.
func SnowDensity(ageDays float64) float64 {
	rhoNew := 50.0
	rhoMax := 400.0
	k := 0.05 // compaction rate
	return rhoNew + (rhoMax-rhoNew)*(1-math.Exp(-k*ageDays))
}

// SWEToDepth converts snow water equivalent (mm) to snow depth (cm).
func SWEToDepth(swe, density float64) float64 {
	if density <= 0 {
		density = 100 // default density
	}
	// SWE (mm) = depth (m) * density (kg/m³) => depth_cm = SWE*10/density
	return swe * 10 / density
}
