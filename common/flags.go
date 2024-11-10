package common

import "flag"

var (
	Gather = flag.String("gather", "", "gather data instead of starting server")
	Testcases = flag.Bool("testcases", false, "gathers testcases instead of challenges")
)
