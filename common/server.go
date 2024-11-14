package common

import "github.com/labstack/echo/v4"

func JsonError(c echo.Context, status int, err error) error {
	return c.JSON(status, echo.Map{
		"message": err.Error(),
	})
}

