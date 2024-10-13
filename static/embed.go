package html

import (
	"embed"

	"github.com/labstack/echo/v4"
)

//go:embed public/*
var public embed.FS
var PublicFS = echo.MustSubFS(public, "public/templates")
var StaticFS = echo.MustSubFS(public, "public/assets")
