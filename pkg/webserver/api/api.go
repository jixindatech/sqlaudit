package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, "ok")
}
