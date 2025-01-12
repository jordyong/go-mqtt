package routes

import (
	"fmt"
	"go-mqtt/pkg/core"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func SetUp(a *core.App) {
	a.MQTTService.Subscribe("device/status", func(c mqtt.Client, m mqtt.Message) {
		fmt.Printf("TOPIC: %s\n", m.Topic())
		fmt.Printf("MSG: %s\n", m.Payload())

		db := a.DBService.GetDB()
		stmt, err := db.Prepare(`
      INSERT INTO devices (id, status)
      VALUES (?,?) 
      `)
		if err != nil {
			fmt.Println(err)
		}

		_, err = stmt.Exec("001", "connected")
		if err != nil {
			fmt.Println(err)
		}

		a.DBService.DisplayDevice()
	})

	a.MQTTService.Subscribe("device/info/gps", func(c mqtt.Client, m mqtt.Message) {
		db := a.DBService.GetDB()
		stmt, err := db.Prepare(`
      INSERT INTO logs (id, gps_x, gps_y)
      VALUE(?,?,?)
      `)
		if err != nil {
			fmt.Println(err)
		}

		_, err = stmt.Exec("001", "100", "100")
		if err != nil {
			fmt.Println(err)
		}
	})
	a.MQTTService.Subscribe("device/info/battery", func(c mqtt.Client, m mqtt.Message) {
		db := a.DBService.GetDB()
		stmt, err := db.Prepare(`
      INSERT INTO logs (id, bat_level, output, solar_input)
      VALUE(?,?,?,?)
      `)
		if err != nil {
			fmt.Println(err)
		}

		_, err = stmt.Exec("001", "50", "60", "80")
		if err != nil {
			fmt.Println(err)
		}
	})

	a.MQTTService.Publish("device/status", "connected")
}
