package server

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/CelestialCrafter/lang-guesser/db"
	"github.com/labstack/echo/v4"
	"github.com/puzpuzpuz/xsync/v3"
)

type session struct {
	CurrentChallenge *db.Challenge
	CurrentStart time.Time
	Past []struct {
		Challenge db.Challenge
		Guessed string
		Duration time.Duration
	}
}
var sessions = xsync.NewMapOf[int64, *session]()
var id = time.Now().UnixNano()

func NewChallenge(c echo.Context) error {
	session, _ := sessions.Compute(id, func(oldValue *session, loaded bool) (newValue *session, delete bool) {
		if loaded {
			return oldValue, false
		}

		return new(session), false
	})

	challenge, err := db.GetRandomChallenge()
	if err != nil {
		return jsonError(c, http.StatusInternalServerError, err)
	}

	session.CurrentChallenge = challenge
	session.CurrentStart = time.Now()
	return c.Blob(http.StatusOK, "text/plain", challenge.Code)
}

func SubmitChallenge(c echo.Context) error {
	var params struct{
		Language string `json:"language"`
	}

	err := c.Bind(&params)
	if err != nil {
		return jsonError(c, http.StatusBadRequest, fmt.Errorf("could not bind params: %w", err))
	}

	session, ok := sessions.Load(id)
	if !ok {
		return jsonError(c, http.StatusBadRequest, errors.New("session does not exist"))
	}

	if session.CurrentChallenge == nil {
		return jsonError(c, http.StatusBadRequest, errors.New("no challenge started"))
	}

	duration := time.Since(session.CurrentStart)
	session.Past = append(session.Past, struct{Challenge db.Challenge; Guessed string; Duration time.Duration}{
		Challenge: *session.CurrentChallenge,
		Guessed: params.Language,
		Duration: duration,
	})

	language := session.CurrentChallenge.Language
	session.CurrentChallenge = nil

	return c.JSON(http.StatusOK, struct{
		Duration float64 `json:"duration"`
		Language string `json:"language"`
		More bool `json:"more"`
	}{
		Duration: duration.Seconds(),
		Language: language,
		More: true,
	})
}
