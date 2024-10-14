package routes

import (
	"fmt"
	"go-mqtt/pkg/core"
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	upgrader = websocket.Upgrader{}
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

	a.Echo.GET("/mqtt-logs", func(c echo.Context) error {
		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}
		defer ws.Close()

		for {
			// Write
			err := ws.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
			if err != nil {
				c.Logger().Error(err)
			}
			// Read
			_, msg, err := ws.ReadMessage()
			if err != nil {
				c.Logger().Error(err)
			}
			fmt.Printf("MSG: %s\n", msg)
		}
	})

}
