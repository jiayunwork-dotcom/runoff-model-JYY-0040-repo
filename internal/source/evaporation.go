package source

import "math"

// PenmanET computes potential evapotranspiration using a simplified Penman formula.
// Returns ET in mm/day.
//
// Parameters:
//   - tempC: mean air temperature (°C)
//   - rh: relative humidity (0-100)
//   - wind: wind speed at 2m (m/s)
//   - solar: solar radiation (MJ/m²/day)
func PenmanET(tempC, rh, wind, solar float64) float64 {
	// Slope of saturation vapor pressure curve.
	delta := 4098 * saturVP(tempC) / math.Pow(tempC+237.3, 2)
	// Psychrometric constant (approx for sea level).
	gamma := 0.0665
	// Saturation and actual vapor pressure.
	es := saturVP(tempC)
	ea := es * rh / 100

	// Net radiation (simplified).
	rn := 0.77*solar - 2.45 // rough net radiation

	// Penman combination.
	num := delta*rn + gamma*(900/(tempC+273))*wind*(es-ea)
	denom := delta + gamma*(1+0.34*wind)
	if denom == 0 {
		return 0
	}
	et := num / denom
	if et < 0 {
		return 0
	}
	return et
}

// HargreavesET computes reference ET using the Hargreaves-Samani method.
// Requires only temperature data.
func HargreavesET(tmin, tmax, tmean float64, ra float64) float64 {
	if ra <= 0 {
		return 0
	}
	et := 0.0023 * ra * math.Sqrt(tmax-tmin) * (tmean + 17.8)
	if et < 0 {
		return 0
	}
	return et
}

// ThornthwaiteET computes monthly potential ET using the Thornthwaite method.
// tempC is the mean monthly temperature, I is the annual heat index,
// daylightHours is the average daylight hours for the month.
func ThornthwaiteET(tempC, I, daylightHours float64) float64 {
	if tempC <= 0 || I <= 0 {
		return 0
	}
	a := 6.75e-7*math.Pow(I, 3) - 7.71e-5*math.Pow(I, 2) + 1.792e-2*I + 0.49239
	et := 16 * math.Pow(10*tempC/I, a) * (daylightHours / 12) / 30
	if et < 0 {
		return 0
	}
	return et
}

// HeatIndex computes the annual Thornthwaite heat index from 12 monthly temperatures.
func HeatIndex(monthlyTemp [12]float64) float64 {
	var I float64
	for _, t := range monthlyTemp {
		if t > 0 {
			I += math.Pow(t/5, 1.514)
		}
	}
	return I
}

func saturVP(tempC float64) float64 {
	return 0.6108 * math.Exp(17.27*tempC/(tempC+237.3))
}
