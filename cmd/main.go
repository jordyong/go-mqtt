package main

import (
	"go-mqtt/internal/routes"
	"go-mqtt/pkg/core"
	"log"
	"strconv"
)

func main() {
	app, err := core.InitApp()
	if err != nil {
		log.Fatalf("Failed to create new app: %v", err)
	}

	routes.BindWebRoute(app)
	app.Echo.Logger.Fatal(app.Echo.Start(":" + strconv.Itoa(app.HttpPort)))
}
