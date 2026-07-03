package command

// seqEqualWidthFlag zero-pads numbers to equal width (-w).
type seqEqualWidthFlag bool

const (
	// SeqEqualWidth (-w): zero-pad every number to equal width.
	SeqEqualWidth seqEqualWidthFlag = true
	// SeqNoEqualWidth is the default: no equal-width padding.
	SeqNoEqualWidth seqEqualWidthFlag = false
)

// SeqSeparator is the -s option value: the separator between numbers. When
// set, all numbers are emitted as a single line joined by this separator.
type SeqSeparator string

// SeqFormat is the -f option value: a printf-style format string applied to
// each number.
type SeqFormat string

// separatorSpec records whether -s was passed and its value, so an explicit
// empty separator still joins the whole run onto one line.
type separatorSpec struct {
	value string
	isSet bool
}

// flags is the parsed flag state for one Seq run.
type flags struct {
	format       SeqFormat
	separator    separatorSpec
	isEqualWidth bool
}

// with folds one option value into the flags, returning the updated copy.
// Values of unrecognized types are ignored.
func (f flags) with(o any) flags {
	switch v := o.(type) {
	case seqEqualWidthFlag:
		f.isEqualWidth = bool(v)
	case SeqSeparator:
		f.separator = separatorSpec{value: string(v), isSet: true}
	case SeqFormat:
		f.format = v
	}
	return f
}
