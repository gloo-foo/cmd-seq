package seq_test

import (
	"github.com/gloo-foo/framework/patterns"

	command "github.com/gloo-foo/cmd-seq"
)

func ExampleSeq_step() {
	// seq 2 2 10
	patterns.MustRun(command.Seq(2.0, 2.0, 10.0))
	// Output:
	// 2
	// 4
	// 6
	// 8
	// 10
}
