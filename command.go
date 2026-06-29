package command

import (
	"context"

	gloo "github.com/gloo-foo/framework"
)

// Seq returns a Source that generates a numeric sequence, matching GNU seq.
//
// Numeric arguments select the bounds:
//   - Seq(last):              1, 2, ..., last        (increment 1)
//   - Seq(first, last):       first, ..., last       (increment 1)
//   - Seq(first, step, last): first, first+step, ... up to last
//
// Flags:
//   - SeqEqualWidth (-w):     zero-pad every number to equal width
//   - SeqSeparator(s) (-s):   join all numbers with s into a single line
//   - SeqFormat(f) (-f):      printf-style format applied to each number
func Seq(args ...any) gloo.Source[[]byte] {
	nums, opts := partition(args)
	b := boundsFrom(nums)
	f := gloo.NewParameters[gloo.File, flags](opts...).Flags
	return seqSource{bounds: b, flags: f}
}

// bounds describes the numeric walk a seq run performs.
type bounds struct {
	first, step, last float64
}

// partition splits the variadic arguments into numeric bounds (in order) and
// the remaining flag options. ints are widened to float64 so callers may pass
// either; everything else is treated as a flag option.
func partition(args []any) (nums []float64, opts []any) {
	for _, a := range args {
		if n, ok := asFloat(a); ok {
			nums = append(nums, n)
			continue
		}
		opts = append(opts, a)
	}
	return nums, opts
}

// asFloat reports whether a is a numeric bound and returns its float64 value.
func asFloat(a any) (float64, bool) {
	switch v := a.(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	default:
		return 0, false
	}
}

// boundsFrom maps 1, 2, or 3 numeric arguments onto first/step/last using GNU
// seq's positional rules. Zero or more than three numbers default to the empty
// 1..1 walk; first and step default to 1, mirroring GNU.
func boundsFrom(nums []float64) bounds {
	b := bounds{first: 1, step: 1, last: 1}
	switch len(nums) {
	case 1:
		b.last = nums[0]
	case 2:
		b.first, b.last = nums[0], nums[1]
	case 3:
		b.first, b.step, b.last = nums[0], nums[1], nums[2]
	}
	return b
}

// values returns the ordered sequence of numbers the bounds describe. A
// non-positive walk (step pointing away from last) yields an empty slice, which
// is exactly GNU seq's behavior for, e.g., `seq 3 1`.
func (b bounds) values() []float64 {
	var out []float64
	for n := b.first; b.inRange(n); n += b.step {
		out = append(out, n)
	}
	return out
}

// inRange reports whether n has not yet passed last in the walk's direction. A
// zero step never advances, so it is treated as out of range to avoid looping.
func (b bounds) inRange(n float64) bool {
	if b.step > 0 {
		return n <= b.last
	}
	if b.step < 0 {
		return n >= b.last
	}
	return false
}

// seqSource is the immutable Source produced by Seq. It holds only the parsed
// bounds and flags; all derived state is computed per Stream call.
type seqSource struct {
	flags  flags
	bounds bounds
}

// Stream renders the sequence into a byte stream. With a separator the whole
// run is one joined line; otherwise each number is its own item.
func (s seqSource) Stream(ctx context.Context) gloo.Stream[[]byte] {
	lines := render(s.bounds.values(), s.flags)
	return gloo.Generate(ctx, emit(lines))
}

// emit returns a producer that sends each pre-rendered line in order, stopping
// early if the consumer discards the stream.
func emit(lines [][]byte) func(context.Context, func([]byte) bool, func(error)) {
	return func(_ context.Context, send func([]byte) bool, _ func(error)) {
		for _, line := range lines {
			if !send(line) {
				return
			}
		}
	}
}
