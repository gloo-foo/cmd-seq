package command

import (
	"fmt"
	"strings"
)

// render turns the numeric sequence into the byte lines a stream emits. With a
// separator set, every number is joined onto a single line; otherwise each
// number becomes its own line.
func render(nums []float64, f flags) [][]byte {
	format := formatter(nums, f)
	texts := formatAll(nums, format)
	if f.separator != nil {
		return [][]byte{[]byte(strings.Join(texts, *f.separator))}
	}
	return toLines(texts)
}

// formatAll applies format to every number, preserving order.
func formatAll(nums []float64, format numberFormat) []string {
	texts := make([]string, len(nums))
	for i, n := range nums {
		texts[i] = format(n)
	}
	return texts
}

// toLines wraps each formatted number as its own byte line.
func toLines(texts []string) [][]byte {
	lines := make([][]byte, len(texts))
	for i, t := range texts {
		lines[i] = []byte(t)
	}
	return lines
}

// numberFormat renders a single number to its textual form.
type numberFormat func(float64) string

// formatter selects the rendering rule for a run. An explicit -f format wins; a
// -w equal-width run zero-pads to a common width and precision; otherwise GNU's
// default %g rendering is used.
func formatter(nums []float64, f flags) numberFormat {
	if f.format != "" {
		return sprintfFormat(f.format)
	}
	if f.equalWidth {
		return paddedFormat(equalWidth(nums))
	}
	return plainFormat
}

// sprintfFormat renders each number through a printf-style format string (-f).
func sprintfFormat(format string) numberFormat {
	return func(n float64) string {
		return fmt.Sprintf(format, n)
	}
}

// plainFormat is GNU seq's default %g rendering.
func plainFormat(n float64) string {
	return fmt.Sprintf("%g", n)
}

// padding is the common field width and decimal precision a -w run pads every
// number to, matching GNU seq's equal-width alignment.
type padding struct {
	width, precision int
}

// paddedFormat renders each number with a fixed precision, zero-padded (after
// the sign) to a common width, so the whole sequence lines up like GNU's -w.
func paddedFormat(p padding) numberFormat {
	return func(n float64) string {
		return fmt.Sprintf("%0*.*f", p.width, p.precision, n)
	}
}

// equalWidth computes the padding for a -w run: the precision is the most
// fractional digits any value needs, and the width is the widest rendering at
// that precision (so signs and integer digits all align).
func equalWidth(nums []float64) padding {
	precision := maxPrecision(nums)
	width := 0
	for _, n := range nums {
		if w := len(fmt.Sprintf("%.*f", precision, n)); w > width {
			width = w
		}
	}
	return padding{width: width, precision: precision}
}

// maxPrecision returns the largest number of fractional digits in the default
// %g rendering of any value, which is the precision GNU's -w uses throughout.
func maxPrecision(nums []float64) int {
	precision := 0
	for _, n := range nums {
		if p := fractionDigits(plainFormat(n)); p > precision {
			precision = p
		}
	}
	return precision
}

// fractionDigits counts the digits after the decimal point in a %g rendering.
func fractionDigits(s string) int {
	if i := strings.IndexByte(s, '.'); i >= 0 {
		return len(s) - i - 1
	}
	return 0
}
