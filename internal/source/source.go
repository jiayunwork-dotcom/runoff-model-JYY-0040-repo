// Package source implements saturation-excess (蓄满产流) runoff generation using a
// Xinanjiang-style tension-water storage-capacity curve.
package source

import (
	"math"

	"runoff-model/internal/basin"
)

// Production computes the surface runoff R for each record using the
// storage-capacity curve:
//
//	pe  = Rain - ET
//	if pe <= 0:      R = 0
//	else: WMM = C * WM
//	      R = pe * (1 - (1 - min(pe, WMM)/WMM)^B)
//
// Properties guaranteed for valid parameters (B>0, C>0, WM>=0):
//   - 0 <= R <= pe
//   - R is monotonically non-decreasing with pe
//   - R == pe when pe >= WMM (the basin is fully saturated)
//
// It returns one R value per input record. A nil/empty input yields a nil slice.
func Production(records []basin.Record, wm, b, c float64) []float64 {
	if len(records) == 0 {
		return nil
	}
	out := make([]float64, len(records))
	wmm := c * wm
	for i, rec := range records {
		pe := rec.Rain - rec.ET
		if pe <= 0 {
			out[i] = 0
			continue
		}
		if wmm <= 0 {
			// No storage capacity -> all net rainfall becomes runoff.
			out[i] = pe
			continue
		}
		ratio := pe / wmm
		if ratio > 1 {
			ratio = 1
		}
		out[i] = pe * (1 - math.Pow(1-ratio, b))
	}
	return out
}
