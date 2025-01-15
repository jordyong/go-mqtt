package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-mqtt/pkg/core/config"
	"go-mqtt/pkg/core/database"
	messages "go-mqtt/pkg/messages"
	"go-mqtt/pkg/mqtt"
	"go-mqtt/pkg/render"
	html "go-mqtt/static"
	"io/fs"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type App struct {
	Echo        *echo.Echo
	HttpPort    int
	PublicFS    fs.FS
	ChatHub     *messages.Hub
	DBService   *database.SQLiteService
	MQTTService *mqtt.MqttService
	Config      *config.Configuration
}

func InitApp() (*App, error) {

	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	sqliteDB, err := database.NewSQliteService()
	if err != nil {
		return nil, err
	}

	mqttService, err := mqtt.NewMqttService(cfg.MQTTClientName, cfg.MQTTBrokerURL)
	if err != nil {
		return nil, err
	}

	app := &App{
		Echo:        echo.New(),
		DBService:   sqliteDB,
		MQTTService: mqttService,
		ChatHub:     messages.NewHub(),
		HttpPort:    8080,
		PublicFS:    html.PublicFS,
		Config:      cfg,
	}

	// Init template Renderer
	renderer, err := render.NewRenderer(app.PublicFS)
	if err != nil {
		return nil, errors.New("Failed to create template renderer: " + err.Error())
	}

	c := app.ChatHub
	c.ParseHTML = func(msg []byte) []byte {
		var htmxJSON messages.HTMX_msg
		err = json.Unmarshal(msg, &htmxJSON)
		msgHTML, _ := renderer.RenderToBytes("message", map[string]any{
			"message": htmxJSON.MQTTMsg,
			"time":    time.Now().Format(time.DateTime),
		})
		fmt.Printf("msgHTML: %s\n", msgHTML)

		return msgHTML
	}

	e := app.Echo
	e.HideBanner = true
	e.Renderer = renderer
	e.Static("/", "static/public/assets")
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339_nano} ${method} ${uri} ${status} ${latency_human}\n",
	}))

	go app.ChatHub.Run()
	return app, nil
}

func (a *App) Close() {

}
