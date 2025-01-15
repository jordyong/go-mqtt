package routes

import (
	"database/sql"
	"fmt"
	"go-mqtt/pkg/core"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func SetUp(a *core.App) {
	a.MQTTService.Subscribe("devices/status", func(c mqtt.Client, m mqtt.Message) {
		fmt.Printf("TOPIC: %s\n", m.Topic())
		fmt.Printf("MSG: %s\n", m.Payload())

		db := a.DBService.GetDB()
		stmt, err := db.Prepare(`
      INSERT INTO devices (id, status, datetime)
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

	var cb mqtt.MessageHandler = func(c mqtt.Client, m mqtt.Message) {
		err := InsertDeviceData(a.DBService.GetDB(), m)
		if err != nil {
			fmt.Println(err)
		}
	}

	a.MQTTService.Subscribe("devices/data/gps", cb)
	a.MQTTService.Subscribe("devices/data/battery/level", cb)
	a.MQTTService.Subscribe("devices/data/battery/charge", cb)
	a.MQTTService.Subscribe("devices/data/battery/output", cb)

	a.MQTTService.Publish("devices/status", "connected")
}

func InsertDevice(db *sql.DB, msg mqtt.Message) error {
	stmt, err := db.Prepare(`
    INSERT INTO device (device_id)
    VALUES(?)
    `)
	if err != nil {
		return err
	}

	device_id := "001"
	_, err = stmt.Exec(device_id)
	if err != nil {
		return err
	}

	return nil
}

func InsertDeviceData(db *sql.DB, msg mqtt.Message) error {
	stmt, err := db.Prepare(`
    INSERT INTO data (device_id, data_type, data_value)
    VALUES(?,?,?)
    `)
	if err != nil {
		return err
	}

	device_id, data_type, data_value := ParseMqttData(msg)
	_, err = stmt.Exec(device_id, data_type, data_value)
	if err != nil {
		return err
	}

	return nil
}
