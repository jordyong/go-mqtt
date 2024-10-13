package handler

import (
	"go-mqtt/pkg/core"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Home(c echo.Context, a *core.App) error {
	return c.Render(http.StatusOK, "index", nil)
}
