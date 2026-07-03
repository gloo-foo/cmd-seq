package seq_test

import (
	"github.com/gloo-foo/framework/patterns"

	command "github.com/gloo-foo/cmd-seq"
)

func ExampleSeq_basic() {
	// seq 5
	if err := patterns.Run(command.Seq(5.0)); err != nil {
		panic(err)
	}
	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
}
