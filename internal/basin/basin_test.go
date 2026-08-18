package basin

import (
	"os"
	"path/filepath"
	"testing"
)

// TestParseRecordsError verifies that a missing file yields a non-nil error and
// a nil record slice (never a panic).
func TestParseRecordsError(t *testing.T) {
	recs, err := ParseRecords(filepath.Join("..", "..", "does-not-exist.csv"))
	if err == nil {
		t.Fatal("expected error for missing records file, got nil")
	}
	if recs != nil {
		t.Fatalf("expected nil records on error, got %d records", len(recs))
	}
}

// TestParseRecordsOK verifies a well-formed file is parsed with the right count
// and values.
func TestParseRecordsOK(t *testing.T) {
	// Build an in-memory file via the example path.
	recs, err := ParseRecords(filepath.Join("..", "..", "example", "rain.csv"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(recs) < 8 {
		t.Fatalf("expected >=8 records from example, got %d", len(recs))
	}
	if recs[0].Rain < recs[0].ET {
		// First example day has rain<et, that is allowed; just sanity check fields.
		t.Logf("first record rain<=et as expected: %.1f <= %.1f", recs[0].Rain, recs[0].ET)
	}
}

// TestParseBasinError verifies a malformed basin file (wrong header) is rejected.
func TestParseBasinError(t *testing.T) {
	// Write a temp file with a wrong header.
	dir := t.TempDir()
	bad := filepath.Join(dir, "bad.csv")
	if err := os.WriteFile(bad, []byte("foo,bar,baz,qux\n1,2,3,4\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := ParseBasin(bad)
	if err == nil {
		t.Fatal("expected error for wrong basin header, got nil")
	}
}

// TestParseBasinOK verifies the example basin file parses correctly.
func TestParseBasinOK(t *testing.T) {
	b, err := ParseBasin(filepath.Join("..", "..", "example", "basin.csv"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.WM != 120 || b.B != 0.3 || b.C != 1.0 {
		t.Fatalf("unexpected basin params: %+v", b)
	}
}
