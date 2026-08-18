package basin

import "errors"

// DefaultParams returns reasonable defaults for a medium-sized basin.
func DefaultParams() Basin {
	return Basin{
		Area: 500,  // km²
		WM:   120,  // mm
		B:    0.3,
		C:    0.15,
	}
}

// Validate checks that basin parameters are physically meaningful.
func Validate(b Basin) error {
	if b.Area <= 0 {
		return errors.New("basin: area must be positive")
	}
	if b.WM <= 0 {
		return errors.New("basin: WM must be positive")
	}
	if b.B < 0 || b.B > 1 {
		return errors.New("basin: B must be in [0,1]")
	}
	if b.C < 0 || b.C > 1 {
		return errors.New("basin: C must be in [0,1]")
	}
	return nil
}

// Describe returns a human-readable description.
func Describe(b Basin) string {
	return "Basin: area=" + fmtF(b.Area) + " km², WM=" + fmtF(b.WM) +
		" mm, B=" + fmtF(b.B) + ", C=" + fmtF(b.C)
}

func fmtF(v float64) string {
	// Simple formatting without importing fmt to keep it minimal.
	s := ""
	if v < 0 {
		s = "-"
		v = -v
	}
	whole := int(v)
	frac := int((v - float64(whole)) * 100)
	return s + itoa(whole) + "." + itoa2(frac)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}

func itoa2(n int) string {
	if n < 10 {
		return "0" + itoa(n)
	}
	return itoa(n)
}
