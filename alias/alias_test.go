package alias_test

import (
	"context"
	"slices"
	"testing"

	seq "github.com/gloo-foo/cmd-seq/alias"
	gloo "github.com/gloo-foo/framework"
)

// The alias package re-exports the constructor and flags under unprefixed
// names. A mis-wired re-export (Format bound to the separator constructor, say,
// or EqualWidth bound to its disabled form) compiles cleanly, so only behavior
// can prove the wiring. Each test exercises one re-export and asserts the GNU
// seq output it must produce.

func lines(t *testing.T, src gloo.Source[[]byte]) []string {
	t.Helper()
	got, err := gloo.Collect(context.Background(), src.Stream(context.Background()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := make([]string, len(got))
	for i, b := range got {
		out[i] = string(b)
	}
	return out
}

func assertLines(t *testing.T, got, want []string) {
	t.Helper()
	if !slices.Equal(got, want) {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestAlias_SeqGeneratesSequence(t *testing.T) {
	// seq 1 3 -> 1, 2, 3
	assertLines(t, lines(t, seq.Seq(1, 3)), []string{"1", "2", "3"})
}

func TestAlias_EqualWidthPadsLeftZeros(t *testing.T) {
	// seq -w 8 11 -> two-digit zero-padded output.
	assertLines(t, lines(t, seq.Seq(8, 11, seq.EqualWidth)), []string{"08", "09", "10", "11"})
}

func TestAlias_NoEqualWidthMatchesDefault(t *testing.T) {
	// NoEqualWidth is the disabled form: identical to passing no flag.
	assertLines(t, lines(t, seq.Seq(8, 11, seq.NoEqualWidth)), []string{"8", "9", "10", "11"})
}

func TestAlias_SeparatorJoinsOntoOneLine(t *testing.T) {
	// seq -s, 1 4 -> a single comma-joined line.
	assertLines(t, lines(t, seq.Seq(1, 4, seq.Separator(","))), []string{"1,2,3,4"})
}

func TestAlias_FormatAppliesPrintfFormat(t *testing.T) {
	// seq -f %.2f 1 3 -> each number rendered with two decimals.
	assertLines(t, lines(t, seq.Seq(1, 3, seq.Format("%.2f"))), []string{"1.00", "2.00", "3.00"})
}
