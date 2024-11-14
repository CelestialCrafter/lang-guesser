package server

import (
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	bindAddress      = ":8080"
	svelteDevAddress = "http://localhost:5173"
)

func getRequestId(c echo.Context) string {
	id := c.Response().Header().Get(echo.HeaderXRequestID)
	if id == "" {
		return "no id"
	}

	return id
}

func SetupServer() {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.RequestID())
	e.Use(middleware.CORS())
	setupRoutes(e)

	finalBindAddress := bindAddress
	envBindAddress := os.Getenv("ADDRESS")
	if envBindAddress != "" {
		finalBindAddress = envBindAddress

	}
	err := e.Start(finalBindAddress)
	if err != nil {
		log.Fatal("error starting server", "error", err)
	}
}
