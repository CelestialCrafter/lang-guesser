package server

import (
	"net/http"

	"github.com/CelestialCrafter/lang-guesser/server/auth"
	"github.com/labstack/echo/v4"
)

func setupRoutes(initial *echo.Echo) {
	svelte(initial)
	initial.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusSeeOther, "/app")
	})

	api := initial.Group("/api")
	api.GET("/auth/:provider", auth.OAuthInit).Name = "oauth-init"
	api.GET("/auth/:provider/callback", auth.OAuthCallback).Name = "oauth-callback"

	authApi := api.Group("", auth.JwtMiddleware())
	authApi.GET("/session", GetSession).Name = "get-session"
	authApi.GET("/challenge", NewChallenge).Name = "new-challenge"
	authApi.POST("/challenge", SubmitChallenge).Name = "submit-challenge"
}
