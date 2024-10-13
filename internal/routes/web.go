package routes

import (
	"fmt"
	"go-mqtt/pkg/core"
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/labstack/echo/v4"
)

func BindWebRoute(a *core.App) {
	a.Echo.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", nil)
	})

	a.Echo.POST("/mqtt-connect", func(c echo.Context) error {
		url := c.FormValue("url-header")
		ip := c.FormValue("ip")
		path := c.FormValue("path")
		clientName := c.FormValue("client-name")

		opts := mqtt.NewClientOptions()
		opts.AddBroker(url + ip + path)
		opts.SetClientID(clientName)

		a.MQTT = mqtt.NewClient(opts)
		if token := a.MQTT.Connect(); token.Wait() && token.Error() != nil {
			fmt.Printf("Failed to init MQTT: %s\n", token.Error())
			return nil
		}

		return c.HTML(http.StatusOK, "")
	})
}
