package server

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/CelestialCrafter/lang-guesser/common"
	"github.com/CelestialCrafter/lang-guesser/common/auth"
	"github.com/CelestialCrafter/lang-guesser/db"
	"github.com/charmbracelet/log"
	"github.com/golang-jwt/jwt/v5"
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
var sessions = xsync.NewMapOf[string, *session]()

func newChallengeAllowed(session *session) bool {
	return len(session.Past) < 10 - 1
}

func NewChallenge(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*auth.UserClaims)
	id := claims.ID

	sess, ok := sessions.Load(id)
	if !ok {
		return common.JsonError(c, http.StatusBadRequest, errors.New("session does not exist"))
	}

	if !newChallengeAllowed(sess) {
		sess = new(session)
		sessions.Store(id, sess)
	}

	challenge, err := db.GetRandomChallenge()
	if err != nil {
		return common.JsonError(c, http.StatusInternalServerError, err)
	}

	sess.CurrentChallenge = challenge
	sess.CurrentStart = time.Now()

	return c.Blob(http.StatusOK, "text/plain", challenge.Code)
}

func stripPast(sess *session) []pastEntry {
	strippedPast := make([]pastEntry, len(sess.Past))
	copy(strippedPast,  sess.Past)
	for i := range strippedPast {
		strippedPast[i].Challenge.Code = []byte{}
	}

	return strippedPast
}

func GetSession(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*auth.UserClaims)
	id := claims.ID

	sess, _ := sessions.Compute(id, func(oldValue *session, loaded bool) (newValue *session, delete bool) {
		if loaded {
			return oldValue, false
		}

		return new(session), false
	})

	return c.JSON(http.StatusOK, stripPast(sess))
}

func SubmitChallenge(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*auth.UserClaims)
	id := claims.ID

	var params struct{
		Language string `json:"language"`
	}

	err := c.Bind(&params)
	if err != nil {
		return common.JsonError(c, http.StatusBadRequest, fmt.Errorf("could not bind params: %w", err))
	}

	sess, ok := sessions.Load(id)
	if !ok {
		return common.JsonError(c, http.StatusBadRequest, errors.New("session does not exist"))
	}

	if sess.CurrentChallenge == nil {
		return common.JsonError(c, http.StatusBadRequest, errors.New("no challenge started"))
	}

	duration := time.Since(sess.CurrentStart)
	sess.Past = append(sess.Past, pastEntry{
		Challenge: *sess.CurrentChallenge,
		Guessed: params.Language,
		Duration: duration,
	})
	sess.CurrentChallenge = nil

	return c.JSON(http.StatusOK, struct{
		More bool `json:"more"`
		Past []pastEntry `json:"past"`
	}{
		More: newChallengeAllowed(sess),
		Past: stripPast(sess),
	})
}
