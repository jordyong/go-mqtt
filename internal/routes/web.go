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

type htmx_ws_msg struct {
	MQTTMsg string `json:"mqtt-message"`
	Headers struct {
		HXRequest     string `json:"HX-Request"`
		HXTrigger     string `json:"HX-Trigger"`
		HXTriggerName string `json:"HX-Trigger-Name"`
		HXTarget      string `json:"HX-Target"`
		HXCurrentURL  string `json:"HX-Current-URL"`
	} `json:"HEADERS"`
}

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
			htmlMessage := fmt.Sprintf(
				`<div id="message" hx-swap-oob="afterend"><strong>%s:</strong> %s</div>`,
				"Test",
				"Hello",
			)
			err = ws.WriteMessage(websocket.TextMessage, []byte(htmlMessage))
			if err != nil {
				c.Logger().Error(err)
			}
			// Read
			var mqttMsg htmx_ws_msg
			err = ws.ReadJSON(&mqttMsg)
			if err != nil {
				c.Logger().Error(err)
			}
			fmt.Printf("MSG: %s\n", mqttMsg.MQTTMsg)
			fmt.Println("Headers:")
			fmt.Printf("  HX-Request: %v\n", mqttMsg.Headers.HXRequest)
			fmt.Printf("  HX-Trigger: %v\n", mqttMsg.Headers.HXTrigger)
			fmt.Printf("  HX-Trigger-Name: %v\n", mqttMsg.Headers.HXTriggerName)
			fmt.Printf("  HX-Target: %v\n", mqttMsg.Headers.HXTarget)
			fmt.Printf("  HX-Current-URL: %v\n", mqttMsg.Headers.HXCurrentURL)
		}
	})

}
