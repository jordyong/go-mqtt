package core

import (
	"errors"
	"go-mqtt/render"
	html "go-mqtt/static"
	"io/fs"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type App struct {
	Echo     *echo.Echo
	HttpPort int
	PublicFS fs.FS
	MQTT     mqtt.Client
}

func InitApp() (*App, error) {

	app := &App{
		Echo:     echo.New(),
		HttpPort: 8080,
		PublicFS: html.PublicFS,
	}

	// Init template Renderer
	renderer, err := render.NewRenderer(app.PublicFS)
	if err != nil {
		return nil, errors.New("Failed to create template renderer: " + err.Error())
	}

	e := app.Echo
	e.HideBanner = true
	e.Renderer = renderer
	e.Static("/", "static/public/assets")
	e.Use(middleware.Logger())

	return app, nil
}
