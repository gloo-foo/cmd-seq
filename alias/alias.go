// Package alias provides unprefixed names for the seq command and its flags.
//
//	import seq "github.com/gloo-foo/cmd-seq/alias"
//	seq.Seq(1, 5, seq.EqualWidth)
package alias

import command "github.com/gloo-foo/cmd-seq"

// Seq re-exports the constructor.
var Seq = command.Seq

// SeqSeparator (-s): join all numbers onto one line with this separator.
var Separator = command.SeqSeparator

// SeqFormat (-f): printf-style format applied to each number.
var Format = command.SeqFormat

// -w flag: zero-pad numbers to equal width.
const EqualWidth = command.SeqEqualWidth

// default: do not pad to equal width.
const NoEqualWidth = command.SeqNoEqualWidth
