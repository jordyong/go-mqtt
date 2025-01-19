package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"go-mqtt/pkg/core"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type DataJson struct {
	Device_id  string  `json:"device_id"`
	Data_type  string  `json:"data_type"`
	Data_value float32 `json:"data_value"`
}

type DeviceJson struct {
	Device_id string `json:"device_id"`
}

func SetUp(a *core.App) {
	var device_cb mqtt.MessageHandler = func(c mqtt.Client, m mqtt.Message) {
		err := InsertDevice(a.DBService.GetDB(), m)
		if err != nil {
			fmt.Printf("Failed to insert device: %s\n", err)
		}
	}

	if err := a.MQTTService.Subscribe("devices/status", device_cb); err != nil {
		fmt.Printf("Failed to subscribe: %s\n", err)
	}

	var data_cb mqtt.MessageHandler = func(c mqtt.Client, m mqtt.Message) {
		var deviceData DataJson
		if err := json.Unmarshal(m.Payload(), &deviceData); err != nil {
			fmt.Println(err)
		}

		if err := LogDeviceData(a.DBService.GetDB(), deviceData); err != nil {
			fmt.Println(err)
		}

		if err := UpdateDeviceInfo(a.DBService.GetDB(), deviceData); err != nil {
			fmt.Println(err)
		}
	}

	a.MQTTService.Subscribe("devices/data/gps", data_cb)
	a.MQTTService.Subscribe("devices/data/battery/level", data_cb)
	a.MQTTService.Subscribe("devices/data/battery/charge", data_cb)
	a.MQTTService.Subscribe("devices/data/battery/output", data_cb)

	a.MQTTService.Publish("devices/status", DeviceJson{"web_client"})
}

func InsertDevice(db *sql.DB, msg mqtt.Message) error {
	stmt, err := db.Prepare(`
    INSERT INTO devices (device_id)
    VALUES(?)
    `)
	if err != nil {
		return err
	}

	var deviceInfo DeviceJson
	if err := json.Unmarshal(msg.Payload(), &deviceInfo); err != nil {
		return err
	}

	_, err = stmt.Exec(deviceInfo.Device_id)
	if err != nil {
		return err
	}

	return nil
}

func UpdateDeviceInfo(db *sql.DB, deviceData DataJson) error {
	// Construct the dynamic SQL query
	query := fmt.Sprintf("UPDATE devices SET %s = ? WHERE device_id = ?", deviceData.Data_type)

	// Prepare the statement
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(deviceData.Data_value, deviceData.Device_id)
	if err != nil {
		return err
	}

	return nil
}

func LogDeviceData(db *sql.DB, deviceData DataJson) error {
	stmt, err := db.Prepare(`
    INSERT INTO data (device_id, data_type, data_value)
    VALUES(?,?,?)
    `)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(deviceData.Device_id, deviceData.Data_type, deviceData.Data_value)
	if err != nil {
		return err
	}

	return nil
}
