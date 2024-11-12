package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func setupRoutes(e *echo.Echo) {
	svelte(e)

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusSeeOther, "/app")
	})

	a := e.Group("/api")
	a.GET("/session", GetSession).Name = "get-session"
	a.GET("/challenge", NewChallenge).Name = "new-challenge"
	a.POST("/challenge", SubmitChallenge).Name = "submit-challenge"
}
