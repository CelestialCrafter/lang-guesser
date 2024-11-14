package server

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/CelestialCrafter/lang-guesser/common"
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

func newChallengeAllowed(session *session) bool {
	return len(session.Past) < 10 - 1
}

func NewChallenge(c echo.Context) error {
	session, ok := sessions.Load(id)
	if !ok {
		return common.JsonError(c, http.StatusBadRequest, errors.New("session does not exist"))
	}

	if !newChallengeAllowed(session) {
		return common.JsonError(c, http.StatusBadRequest, errors.New("new challenges disabled for session"))
	}

	challenge, err := db.GetRandomChallenge()
	if err != nil {
		return common.JsonError(c, http.StatusInternalServerError, err)
	}

	session.CurrentChallenge = challenge
	session.CurrentStart = time.Now()

	return c.Blob(http.StatusOK, "text/plain", challenge.Code)
}

func stripPast(session *session) []pastEntry {
	strippedPast := make([]pastEntry, len(session.Past))
	copy(strippedPast,  session.Past)
	for i := range strippedPast {
		strippedPast[i].Challenge.Code = []byte{}
	}

	return strippedPast
}

func GetSession(c echo.Context) error {
	session, _ := sessions.Compute(id, func(oldValue *session, loaded bool) (newValue *session, delete bool) {
		if loaded {
			return oldValue, false
		}

		return new(session), false
	})

	return c.JSON(http.StatusOK, stripPast(session))
}

func SubmitChallenge(c echo.Context) error {
	var params struct{
		Language string `json:"language"`
	}

	err := c.Bind(&params)
	if err != nil {
		return common.JsonError(c, http.StatusBadRequest, fmt.Errorf("could not bind params: %w", err))
	}

	session, ok := sessions.Load(id)
	if !ok {
		return common.JsonError(c, http.StatusBadRequest, errors.New("session does not exist"))
	}

	if session.CurrentChallenge == nil {
		return common.JsonError(c, http.StatusBadRequest, errors.New("no challenge started"))
	}

	duration := time.Since(session.CurrentStart)
	session.Past = append(session.Past, pastEntry{
		Challenge: *session.CurrentChallenge,
		Guessed: params.Language,
		Duration: duration,
	})
	session.CurrentChallenge = nil

	return c.JSON(http.StatusOK, struct{
		More bool `json:"more"`
		Past []pastEntry `json:"past"`
	}{
		More: newChallengeAllowed(session),
		Past: stripPast(session),
	})
}
