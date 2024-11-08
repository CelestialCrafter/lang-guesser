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
	a.GET("/challenge", GetChallenge).Name = "get-challenge"
}
