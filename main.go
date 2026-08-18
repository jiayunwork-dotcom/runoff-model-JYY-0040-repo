// Command runoff-model is the CLI front-end for the rainfall-runoff model.
// It only parses flags and wires the internal packages together.
package main

import (
	"flag"
	"fmt"
	"os"

	"runoff-model/internal/basin"
	"runoff-model/internal/hydro"
	"runoff-model/internal/server"
	"runoff-model/internal/source"
)

const usage = `usage: runoff-model -basin <path> -rain <path>

  -basin <path>  CSV with header "area,wm,b,c" (one data row)
  -rain  <path>  CSV with header "date,rain,et,runoff" (>=1 data rows)
`

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() error {
	basinPath := flag.String("basin", "", "path to basin CSV")
	rainPath := flag.String("rain", "", "path to rainfall/ET/observed CSV")
	serve := flag.String("serve", "", "start HTTP server on this address (e.g. :8080)")
	flag.Parse()

	// Serve mode: start HTTP API + frontend.
	if *serve != "" {
		fmt.Printf("Starting server on %s ...\n", *serve)
		s := server.New(*serve)
		return s.ListenAndServe()
	}

	if *basinPath == "" || *rainPath == "" {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}

	b, err := basin.ParseBasin(*basinPath)
	if err != nil {
		return fmt.Errorf("parse basin: %w", err)
	}
	records, err := basin.ParseRecords(*rainPath)
	if err != nil {
		return fmt.Errorf("parse records: %w", err)
	}
	if len(records) == 0 {
		return fmt.Errorf("no records to process")
	}

	// 1) Production (saturation-excess runoff generation).
	R := source.Production(records, b.WM, b.B, b.C)
	if len(R) != len(records) {
		return fmt.Errorf("production length mismatch")
	}

	// 2) Concentration: Nash UH + convolution of the effective rainfall (R).
	const (
		nUH = 3
		Kuh = 12.0
		dt  = 6.0
	)
	uh := hydro.NashUH(nUH, Kuh, dt)
	drh := hydro.Convolve(R, uh)
	if len(drh) == 0 {
		return fmt.Errorf("convolution produced an empty hydrograph")
	}

	// 3) Routing: Muskingum.
	const (
		Km   = 12.0
		xm   = 0.2
		dtm  = 6.0
	)
	routed := hydro.Muskingum(drh, Km, xm, dtm)
	if len(routed) != len(drh) {
		return fmt.Errorf("routing length mismatch")
	}

	// Peak of the routed hydrograph.
	peak := routed[0]
	peakIdx := 0
	for i, v := range routed {
		if v > peak {
			peak = v
			peakIdx = i
		}
	}

	// 4) Frequency analysis of the peak against observed annual-max samples.
	samples := make([]float64, len(records))
	for i, rec := range records {
		samples[i] = rec.Observed
	}
	T := hydro.PearsonIII(samples, peak)

	// Report.
	fmt.Printf("# Basin: area=%.2f km^2, WM=%.2f mm, B=%.2f, C=%.2f\n", b.Area, b.WM, b.B, b.C)
	fmt.Println("# Production runoff R (mm) per record:")
	for i, r := range R {
		fmt.Printf("  %s R=%.3f\n", records[i].Date, r)
	}
	fmt.Printf("# Nash UH ordinates (%d):\n", len(uh))
	for i, u := range uh {
		fmt.Printf("  t=%d  u=%.5f\n", i, u)
	}
	fmt.Println("# Routed hydrograph:")
	for i, v := range routed {
		mark := ""
		if i == peakIdx {
			mark = "  <- peak"
		}
		fmt.Printf("  t=%d  Q=%.3f%s\n", i, v, mark)
	}
	fmt.Printf("# Peak routed discharge = %.3f at step %d\n", peak, peakIdx)
	fmt.Printf("# Peak return period T = %.3f years\n", T)
	return nil
}
