package server

import (
	"net/http"
	"time"

	"github.com/CelestialCrafter/lang-guesser/db"
	"github.com/labstack/echo/v4"
	"github.com/puzpuzpuz/xsync/v3"
)

var runningChallenges = xsync.NewMapOf[int64, *db.Challenge]()

func GetChallenge(c echo.Context) error {
	challenge, err := db.GetRandomChallenge()
	if err != nil {
		return jsonError(c, http.StatusInternalServerError, err)
	}

	runningChallenges.Store(time.Now().UnixNano(),  challenge)

	challenge.Language = ""
	challenge.Sha = ""
	return c.JSON(http.StatusOK, challenge)
}
