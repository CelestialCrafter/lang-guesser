package main

import (
	"flag"

	"github.com/CelestialCrafter/lang-guesser/common"
	"github.com/CelestialCrafter/lang-guesser/db"
	"github.com/CelestialCrafter/lang-guesser/gather"
	"github.com/CelestialCrafter/lang-guesser/server"
	"github.com/CelestialCrafter/lang-guesser/server/auth"
	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
)

func main() {
	flag.Parse()
	log.SetLevel(log.DebugLevel)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("could not load .env file", "error", err)
	}
	db.InitChallenges()
	auth.InitializeProviders()

	if *common.Gather != "" {
		gather.Gather()
	} else {
		server.SetupServer()
	}
}
