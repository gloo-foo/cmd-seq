package command_test

import (
	"context"
	"testing"

	gloo "github.com/gloo-foo/framework"

	command "github.com/gloo-foo/cmd-seq"
)

func collect(t *testing.T, src gloo.Source[[]byte]) []string {
	t.Helper()
	got, err := gloo.Collect(context.Background(), src.Stream(context.Background()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := make([]string, len(got))
	for i, b := range got {
		result[i] = string(b)
	}
	return result
}

func assertStrings(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("got %d items %v, want %d items %v", len(got), got, len(want), want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("item %d: got %q, want %q", i, got[i], want[i])
		}
	}
}

func TestSeq_SingleArg(t *testing.T) {
	assertStrings(t, collect(t, command.Seq(3)), []string{"1", "2", "3"})
}

func TestSeq_Range(t *testing.T) {
	assertStrings(t, collect(t, command.Seq(2, 5)), []string{"2", "3", "4", "5"})
}

func TestSeq_Step(t *testing.T) {
	assertStrings(t, collect(t, command.Seq(1, 2, 7)), []string{"1", "3", "5", "7"})
}

func TestSeq_EqualWidth(t *testing.T) {
	got := collect(t, command.Seq(8, 11, command.SeqEqualWidth))
	assertStrings(t, got, []string{"08", "09", "10", "11"})
}

func TestSeq_EqualWidthSingleDigits(t *testing.T) {
	got := collect(t, command.Seq(1, 5, command.SeqEqualWidth))
	assertStrings(t, got, []string{"1", "2", "3", "4", "5"})
}

func TestSeq_EqualWidthLargeRange(t *testing.T) {
	got := collect(t, command.Seq(1, 100, command.SeqEqualWidth))
	if got[0] != "001" {
		t.Errorf("first item: got %q, want %q", got[0], "001")
	}
	if got[99] != "100" {
		t.Errorf("last item: got %q, want %q", got[99], "100")
	}
}

func TestSeq_Separator(t *testing.T) {
	got := collect(t, command.Seq(1, 4, command.SeqSeparator(", ")))
	assertStrings(t, got, []string{"1, 2, 3, 4"})
}

func TestSeq_SeparatorColon(t *testing.T) {
	got := collect(t, command.Seq(1, 3, command.SeqSeparator(":")))
	assertStrings(t, got, []string{"1:2:3"})
}

func TestSeq_Format(t *testing.T) {
	got := collect(t, command.Seq(1, 3, command.SeqFormat("%05.2f")))
	assertStrings(t, got, []string{"01.00", "02.00", "03.00"})
}

func TestSeq_FormatSimple(t *testing.T) {
	got := collect(t, command.Seq(1, 3, command.SeqFormat("%.0f")))
	assertStrings(t, got, []string{"1", "2", "3"})
}

func TestSeq_SeparatorAndFormat(t *testing.T) {
	got := collect(t, command.Seq(1, 3, command.SeqSeparator(" "), command.SeqFormat("%02.0f")))
	assertStrings(t, got, []string{"01 02 03"})
}

func TestSeq_DescendingNegativeStep(t *testing.T) {
	// seq 5 -1 1 walks downward and is inclusive of last.
	assertStrings(t, collect(t, command.Seq(5, -1, 1)), []string{"5", "4", "3", "2", "1"})
}

func TestSeq_NegativeStepLargerStride(t *testing.T) {
	// seq 10 -2 1 stops at the last value not below 1.
	assertStrings(t, collect(t, command.Seq(10, -2, 1)), []string{"10", "8", "6", "4", "2"})
}

func TestSeq_DescendingWithoutNegativeStepIsEmpty(t *testing.T) {
	// seq 3 1 with the default +1 step never reaches a value <= 1, so GNU emits
	// nothing.
	assertStrings(t, collect(t, command.Seq(3, 1)), []string{})
}

func TestSeq_ZeroStepIsEmpty(t *testing.T) {
	// A zero step can never advance; rather than loop forever, the walk is empty.
	assertStrings(t, collect(t, command.Seq(1.0, 0.0, 5.0)), []string{})
}

func TestSeq_FloatIncrement(t *testing.T) {
	// seq 1 0.5 2.5 -> fractional steps render with %g by default.
	assertStrings(t, collect(t, command.Seq(1.0, 0.5, 2.5)), []string{"1", "1.5", "2", "2.5"})
}

func TestSeq_EqualWidthFractional(t *testing.T) {
	// seq -w 8 0.5 10 -> a fixed one-decimal precision, zero-padded to width 4,
	// matching GNU exactly (08.0 .. 10.0).
	got := collect(t, command.Seq(8.0, 0.5, 10.0, command.SeqEqualWidth))
	assertStrings(t, got, []string{"08.0", "08.5", "09.0", "09.5", "10.0"})
}

func TestSeq_EqualWidthNegativeFractional(t *testing.T) {
	// seq -w -1 0.5 1 -> zero padding lands after the sign: -1.0, ..., 00.0, 01.0.
	got := collect(t, command.Seq(-1.0, 0.5, 1.0, command.SeqEqualWidth))
	assertStrings(t, got, []string{"-1.0", "-0.5", "00.0", "00.5", "01.0"})
}

func TestSeq_NoEqualWidthMatchesDefault(t *testing.T) {
	// The disabled -w form behaves exactly like passing no flag.
	got := collect(t, command.Seq(8, 11, command.SeqNoEqualWidth))
	assertStrings(t, got, []string{"8", "9", "10", "11"})
}

// TestSeq_DiscardStopsProducer pulls a single item and then discards the stream.
// The framework's generator uses an unbuffered channel, so the producer blocks
// on the next send; the discard makes that send report false, exercising the
// early-return path in emit.
func TestSeq_DiscardStopsProducer(t *testing.T) {
	stream := command.Seq(1, 1000).Stream(context.Background())
	first := <-stream.Chan()
	if first.Error != nil {
		t.Fatalf("unexpected error: %v", first.Error)
	}
	if got := string(first.Value); got != "1" {
		t.Fatalf("first item: got %q, want %q", got, "1")
	}
	stream.Discard()
}
