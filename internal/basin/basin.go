// Package basin parses basin parameters and daily hydrological records from CSV.
package basin

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Record is one day of observations for a basin.
type Record struct {
	Date     string
	Rain     float64
	ET       float64
	Observed float64
}

// Basin holds the lumped parameters of a catchment.
type Basin struct {
	Area float64 // drainage area (km^2)
	WM   float64 // areal mean tension-water capacity (mm)
	B    float64 // exponent of storage-capacity curve
	C    float64 // ratio of areal mean capacity to WM
}

// ParseRecords reads a CSV with header `date,rain,et,runoff`.
// It returns an error if the file is missing, the header is wrong, or any row
// is malformed.
func ParseRecords(path string) ([]Record, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open records: %w", err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = 4
	rows, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("read records csv: %w", err)
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("records file is empty")
	}
	header := strings.ToLower(strings.Join(rows[0], ","))
	header = strings.TrimPrefix(header, "\ufeff") // tolerate a UTF-8 BOM
	if header != "date,rain,et,runoff" {
		return nil, fmt.Errorf("unexpected records header %q, want date,rain,et,runoff", rows[0])
	}
	if len(rows) == 1 {
		return nil, fmt.Errorf("records file has header but no data rows")
	}

	out := make([]Record, 0, len(rows)-1)
	for i, row := range rows[1:] {
		if len(row) != 4 {
			return nil, fmt.Errorf("row %d: expected 4 fields, got %d", i+2, len(row))
		}
		rain, err := strconv.ParseFloat(strings.TrimSpace(row[1]), 64)
		if err != nil {
			return nil, fmt.Errorf("row %d rain: %w", i+2, err)
		}
		et, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
		if err != nil {
			return nil, fmt.Errorf("row %d et: %w", i+2, err)
		}
		obs, err := strconv.ParseFloat(strings.TrimSpace(row[3]), 64)
		if err != nil {
			return nil, fmt.Errorf("row %d runoff: %w", i+2, err)
		}
		out = append(out, Record{Date: row[0], Rain: rain, ET: et, Observed: obs})
	}
	return out, nil
}

// ParseBasin reads a CSV with header `area,wm,b,c` and a single data row.
func ParseBasin(path string) (Basin, error) {
	var b Basin
	f, err := os.Open(path)
	if err != nil {
		return b, fmt.Errorf("open basin: %w", err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = 4
	rows, err := r.ReadAll()
	if err != nil {
		return b, fmt.Errorf("read basin csv: %w", err)
	}
	if len(rows) == 0 {
		return b, fmt.Errorf("basin file is empty")
	}
	header := strings.ToLower(strings.Join(rows[0], ","))
	header = strings.TrimPrefix(header, "\ufeff") // tolerate a UTF-8 BOM
	if header != "area,wm,b,c" {
		return b, fmt.Errorf("unexpected basin header %q, want area,wm,b,c", rows[0])
	}
	if len(rows) < 2 {
		return b, fmt.Errorf("basin file has header but no data row")
	}
	if len(rows[1]) != 4 {
		return b, fmt.Errorf("basin row: expected 4 fields, got %d", len(rows[1]))
	}
	area, err := strconv.ParseFloat(strings.TrimSpace(rows[1][0]), 64)
	if err != nil {
		return b, fmt.Errorf("basin area: %w", err)
	}
	wm, err := strconv.ParseFloat(strings.TrimSpace(rows[1][1]), 64)
	if err != nil {
		return b, fmt.Errorf("basin wm: %w", err)
	}
	bb, err := strconv.ParseFloat(strings.TrimSpace(rows[1][2]), 64)
	if err != nil {
		return b, fmt.Errorf("basin b: %w", err)
	}
	c, err := strconv.ParseFloat(strings.TrimSpace(rows[1][3]), 64)
	if err != nil {
		return b, fmt.Errorf("basin c: %w", err)
	}
	b = Basin{Area: area, WM: wm, B: bb, C: c}
	return b, nil
}
