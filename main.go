package main

import (
	"flag"

	"github.com/CelestialCrafter/lang-guesser/common"
	"github.com/CelestialCrafter/lang-guesser/db"
	"github.com/CelestialCrafter/lang-guesser/gather"
	"github.com/charmbracelet/log"
)

func main() {
	flag.Parse()
	log.SetLevel(log.DebugLevel)
	db.InitChallenges()

	if *common.Gather {
		gather.Gather()
	}
}
