package routes

import (
	"go-mqtt/pkg/core"
	"go-mqtt/pkg/messages"
	"go-mqtt/pkg/mqtt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func BindWebRoute(a *core.App) {

	a.Echo.GET("/", func(c echo.Context) error {
		clientName := a.Config.MQTTClientName
		brokerURL := a.Config.MQTTBrokerURL
		return c.Render(http.StatusOK, "index", map[string]any{
			"IsConnected": a.MQTT.IsConnected(),
			"ClientName":  clientName,
			"BrokerURL":   brokerURL,
		})
	})

	a.Echo.POST("/mqtt-connect", func(c echo.Context) error {
		clientName := a.Config.MQTTClientName
		brokerURL := a.Config.MQTTBrokerURL

		a.MQTT = mqtt.ConnectMQTT(clientName, brokerURL)

		return c.Render(http.StatusOK, "topbar", map[string]any{
			"IsConnected": a.MQTT.IsConnected(),
			"ClientName":  clientName,
			"BrokerURL":   brokerURL,
		})
	})

	a.Echo.POST("/mqtt-disconnect", func(c echo.Context) error {
		clientName := a.Config.MQTTClientName
		brokerURL := a.Config.MQTTBrokerURL
		a.MQTT.Disconnect(250)

		return c.Render(http.StatusOK, "topbar", map[string]any{
			"IsConnected": a.MQTT.IsConnected(),
			"ClientName":  clientName,
			"BrokerURL":   brokerURL,
		})
	})

	a.Echo.GET("/mqtt-logs", func(c echo.Context) error {
		err := messages.ServeWS(a.ChatHub, c.Response(), c.Request())
		if err != nil {
			return err
		}
		return nil
	})

	a.Echo.GET("/cmd", func(c echo.Context) error {
		client := a.MQTT
		if client.IsConnected() {
			mqtt.PublishMQTT(a.MQTT, "CMD", c.QueryParam("direction"))
		}
		return nil
	})
}
