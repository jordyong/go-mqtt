package routes

import (
	"go-mqtt/pkg/core"
	"go-mqtt/pkg/messages"
	"net/http"

	"github.com/labstack/echo/v4"
)

func BindWebRoute(a *core.App) {

	a.Echo.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", map[string]any{
			"content": "analytics",
		})
	})

	a.Echo.GET("/analytics", func(c echo.Context) error {
		return c.Render(http.StatusOK, "analytics", map[string]any{
			"content": "analytics",
		})
	})

	a.Echo.GET("/logs", func(c echo.Context) error {
		return c.Render(http.StatusOK, "logs", map[string]any{
			"content": "logs",
		})
	})

	a.Echo.GET("/mqtt-connect", func(c echo.Context) error {
		clientName := a.Config.MQTTClientName
		brokerURL := a.Config.MQTTBrokerURL
		return c.Render(http.StatusOK, "topbar", map[string]any{
			"IsConnected": a.MQTTService.IsConnected(),
			"ClientName":  clientName,
			"BrokerURL":   brokerURL,
		})
	})

	a.Echo.POST("/mqtt-connect", func(c echo.Context) error {
		clientName := a.Config.MQTTClientName
		brokerURL := a.Config.MQTTBrokerURL

		a.MQTTService.Connect()
		SetUp(a)

		return c.Render(http.StatusOK, "topbar", map[string]any{
			"IsConnected": a.MQTTService.IsConnected(),
			"ClientName":  clientName,
			"BrokerURL":   brokerURL,
		})
	})

	a.Echo.POST("/mqtt-disconnect", func(c echo.Context) error {
		clientName := a.Config.MQTTClientName
		brokerURL := a.Config.MQTTBrokerURL
		a.MQTTService.Disconnect()

		return c.Render(http.StatusOK, "topbar", map[string]any{
			"IsConnected": a.MQTTService.IsConnected(),
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

	a.Echo.GET("/mqtt", func(c echo.Context) error {
		client := a.MQTTService
		if client.IsConnected() {
			a.MQTTService.Publish("/topic/qos0", c.QueryParam("cmd"))
		}
		return nil
	})
}
