package seq_test

import (
	"github.com/gloo-foo/framework/patterns"

	command "github.com/gloo-foo/cmd-seq"
)

func ExampleSeq_equalWidth() {
	// seq -w 8 11
	if err := patterns.Run(command.Seq(8.0, 11.0, command.SeqEqualWidth)); err != nil {
		panic(err)
	}
	// Output:
	// 08
	// 09
	// 10
	// 11
}

func ExampleSeq_separator() {
	// seq -s, 1 5
	if err := patterns.Run(command.Seq(1.0, 5.0, command.SeqSeparator(","))); err != nil {
		panic(err)
	}
	// Output:
	// 1,2,3,4,5
}

func ExampleSeq_format() {
	// seq -f %.2f 1 3
	if err := patterns.Run(command.Seq(1.0, 3.0, command.SeqFormat("%.2f"))); err != nil {
		panic(err)
	}
	// Output:
	// 1.00
	// 2.00
	// 3.00
}
