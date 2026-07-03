// Package alias provides unprefixed names for the seq command and its flags.
//
//	import seq "github.com/gloo-foo/cmd-seq/alias"
//	seq.Seq(1, 5, seq.EqualWidth)
package alias

import (
	gloo "github.com/gloo-foo/framework"

	command "github.com/gloo-foo/cmd-seq"
)

// Seq is the command constructor; it forwards to command.Seq.
func Seq(args ...any) gloo.Source[[]byte] { return command.Seq(args...) }

// Separator (-s): join all numbers onto one line with this separator.
type Separator = command.SeqSeparator

// Format (-f): printf-style format applied to each number.
type Format = command.SeqFormat

// EqualWidth (-w): zero-pad numbers to equal width.
const EqualWidth = command.SeqEqualWidth

// NoEqualWidth is the default: do not pad to equal width.
const NoEqualWidth = command.SeqNoEqualWidth
