package seq_test

import (
	"github.com/gloo-foo/framework/patterns"

	command "github.com/gloo-foo/cmd-seq"
)

func ExampleSeq_step() {
	// seq 2 2 10
	if err := patterns.Run(command.Seq(2.0, 2.0, 10.0)); err != nil {
		panic(err)
	}
	// Output:
	// 2
	// 4
	// 6
	// 8
	// 10
}
