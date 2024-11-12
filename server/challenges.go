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

type pastEntry struct {
	Challenge db.Challenge `json:"challenge"`
	Guessed string `json:"guessed"`
	Duration time.Duration `json:"duration"`
}

type session struct {
	CurrentChallenge *db.Challenge
	CurrentStart time.Time
	Past []pastEntry
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
	session.Past = append(session.Past, pastEntry{
		Challenge: *session.CurrentChallenge,
		Guessed: params.Language,
		Duration: duration,
	})
	session.CurrentChallenge = nil

	strippedPast := make([]pastEntry, len(session.Past))
	copy(strippedPast,  session.Past)
	for i := range strippedPast {
		strippedPast[i].Challenge.Code = []byte{}
	}

	return c.JSON(http.StatusOK, strippedPast)
}
