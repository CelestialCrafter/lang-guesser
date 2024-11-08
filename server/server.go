package server

import (
	"log"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	bindAddress      = ":8000"
	svelteDevAddress = "http://localhost:5173"
)

func jsonError(c echo.Context, status int, err error) error {
	return c.JSON(status, echo.Map{
		"message": err.Error(),
	})
}

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
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "response timed out",
		Timeout:      30 * time.Second,
	}))

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
