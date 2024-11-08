package main

import (
	"flag"

	"github.com/CelestialCrafter/lang-guesser/common"
	"github.com/CelestialCrafter/lang-guesser/gather"
)

func main() {
	flag.Parse()

	if *common.Gather {
		gather.Gather()
	}
}
