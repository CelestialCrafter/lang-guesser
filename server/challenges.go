package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/CelestialCrafter/lang-guesser/db"
	"github.com/labstack/echo/v4"
	"github.com/puzpuzpuz/xsync/v3"
)

var runningChallenges = xsync.NewMapOf[int64, db.Challenge]()

func GetChallenge(c echo.Context) error {
	challenge, err := db.GetRandomChallenge()
	if err != nil {
		return jsonError(c, http.StatusInternalServerError, err)
	}

	id := time.Now().UnixMicro()
	runningChallenges.Store(id,  *challenge)

	challenge.Language = ""
	challenge.Sha = ""
	idedChallenge := struct{
		Id int64 `json:"id"`
		*db.Challenge
	}{
		Challenge: challenge,
		Id: id,
	}
	return c.JSON(http.StatusOK, idedChallenge)
}

func PostChallenge(c echo.Context) error {
	var params struct{
		Id int64 `json:"id"`
		Language string `json:"language"`
	}

	err := c.Bind(&params)
	if err != nil {
		return jsonError(c, http.StatusInternalServerError, fmt.Errorf("could not bind params: %w", err))
	}

	challenge, ok := runningChallenges.LoadAndDelete(params.Id)
	if !ok {
		return jsonError(c, http.StatusNotFound, fmt.Errorf("challenge id not found"))
	}

	return c.JSON(http.StatusOK, struct{
		Time time.Duration `json:"time"`
		Language string `json:"language"`
		More bool `json:"more"`
	}{
		Time: time.Since(time.UnixMicro(params.Id)),
		Language: challenge.Language,
		More: true,
	})
}
