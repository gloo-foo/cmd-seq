package command

// seqEqualWidthFlag zero-pads numbers to equal width (-w).
type seqEqualWidthFlag bool

const (
	SeqEqualWidth   seqEqualWidthFlag = true
	SeqNoEqualWidth seqEqualWidthFlag = false
)

func (f seqEqualWidthFlag) Configure(flags *flags) { flags.equalWidth = bool(f) }

// seqSeparatorFlag sets the separator between numbers (-s). When set, all
// numbers are emitted as a single line joined by this separator.
type seqSeparatorFlag string

// SeqSeparator sets a custom separator between numbers.
func SeqSeparator(s string) seqSeparatorFlag { return seqSeparatorFlag(s) }

func (f seqSeparatorFlag) Configure(flags *flags) {
	s := string(f)
	flags.separator = &s
}

// seqFormatFlag sets a printf-style format string applied to each number (-f).
type seqFormatFlag string

// SeqFormat sets the printf-style format for each number.
func SeqFormat(f string) seqFormatFlag { return seqFormatFlag(f) }

func (f seqFormatFlag) Configure(flags *flags) { flags.format = string(f) }

// flags is the parsed flag state for one Seq run.
type flags struct {
	equalWidth bool
	separator  *string
	format     string
}
