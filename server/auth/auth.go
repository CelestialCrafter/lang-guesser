package auth

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/CelestialCrafter/lang-guesser/common"
	"github.com/CelestialCrafter/lang-guesser/common/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func JwtMiddleware() echo.MiddlewareFunc {
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.UserClaims)
		},
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}
	return echojwt.WithConfig(config)
}

func OAuthInit(c echo.Context) error {
	provider, ok := providers[c.Param("provider")]
	if !ok {
		return common.JsonError(c, http.StatusBadRequest, errors.New("provider not supported"))
	}

	state := hex.EncodeToString(auth.Hash())
	c.SetCookie(&http.Cookie{
		Name:  "state",
		Value: state,

		MaxAge:   int((time.Minute * 5).Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	return c.Redirect(http.StatusSeeOther, provider.GetUrl(state))
}

func OAuthCallback(c echo.Context) error {
	provider, ok := providers[c.Param("provider")]
	if !ok {
		return common.JsonError(c, http.StatusBadRequest, errors.New("provider not supported"))
	}

	originalState, err := c.Cookie("state")
	state := c.QueryParam("state")
	if err != nil || originalState.Value != state {
		return common.JsonError(c, http.StatusBadRequest, errors.New("could not verify state"))
	}

	claims, err := provider.Exchange(c.QueryParam("code"))
	if err != nil {
		return common.JsonError(c, http.StatusInternalServerError, err)
	}

	// exchange code
	token, err := auth.Sign(claims)
	if err != nil {
		return common.JsonError(c, http.StatusInternalServerError, err)
	}

	// create user
	// _, err = db.GetUser(claims.ID)
	// if err != nil {
	// 	if !errors.Is(err, sql.ErrNoRows) {
	// 		return err
	// 	}
	//
	// 	_, err = db.CreateUser(claims.ID)
	// 	if err != nil {
	// 		return common.JsonError(c, http.StatusInternalServerError, err)
	// 	}
	// }

	return c.HTML(http.StatusOK, fmt.Sprintf(`
	<!doctype html>
	<html>
	<script>
	localStorage.setItem("token", "%s");
	window.location.assign("/app/auth");
	</script>

	Click <a href="/app/auth">here</a> to continue
	</html>
	`, token))
}
