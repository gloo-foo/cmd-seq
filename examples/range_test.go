package seq_test

import (
	command "github.com/gloo-foo/cmd-seq"
	"github.com/gloo-foo/framework/patterns"
)

func ExampleSeq_range() {
	// seq 3 7
	patterns.MustRun(command.Seq(3.0, 7.0))
	// Output:
	// 3
	// 4
	// 5
	// 6
	// 7
}
