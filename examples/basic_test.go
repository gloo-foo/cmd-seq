package seq_test

import (
	command "github.com/gloo-foo/cmd-seq"
	"github.com/gloo-foo/framework/patterns"
)

func ExampleSeq_basic() {
	// seq 5
	patterns.MustRun(command.Seq(5.0))
	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
}
