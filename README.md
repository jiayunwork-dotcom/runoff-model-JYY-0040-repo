# runoff-model

A self-contained Go module implementing a rainfall–runoff model and flood
frequency analysis for hydrological training data.

## What it does

Given a basin parameter file and a rainfall/ET/observed-runoff record file
(both CSV), the model:

1. **Production (蓄满产流 / saturation-excess)** — `internal/source.Production`
   uses a Xinanjiang-style tension-water storage-capacity curve to convert net
   rainfall (`rain − ET`) into surface runoff `R` for every record.
2. **Concentration** — `internal/hydro.NashUH` builds a Nash n-linear-reservoir
   instantaneous unit hydrograph; `internal/hydro.Convolve` convolves the
   effective rainfall (the production series) with the UH to obtain the direct
   runoff hydrograph.
3. **Routing** — `internal/hydro.Muskingum` routes the hydrograph through a
   channel reach.
4. **Frequency analysis** — `internal/hydro.PearsonIII` computes the return
   period of the routed peak using a Pearson type-III fit (Cs ≈ 2·Cv) to the
   annual-maximum peak samples.

All code is pure Go standard library (no third-party imports, no network).

## Packages

| Package | Responsibility |
| --- | --- |
| `internal/basin` | Parse basin params and daily records from CSV. |
| `internal/source` | Saturation-excess runoff generation. |
| `internal/hydro`  | Nash UH, convolution, Muskingum routing, Pearson-III frequency. |

## CLI usage

```
go run . -basin example/basin.csv -rain example/rain.csv
```

- `-basin <path>` — CSV with header `area,wm,b,c` (one data row).
- `-rain <path>`  — CSV with header `date,rain,et,runoff` (≥1 data rows).

Exit codes:

- `2` — required flags missing (usage printed to stderr).
- `1` — bad/malformed input (clear message to stderr).
- `0` — success; prints the production runoff series, Nash UH ordinates,
  routed hydrograph, and the peak return period.

## Examples

See `example/basin.csv` and `example/rain.csv`. The happy path:

```
go run . -basin example/basin.csv -rain example/rain.csv
```

## Build / test

```
export GOTOOLCHAIN=local CGO_ENABLED=0
go vet ./...
go build ./...
go test ./...
```
