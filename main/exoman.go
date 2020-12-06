package main

import (
	"fmt"

	"exolifa.com/exoman/logger"
	"exolifa.com/exoman/params"
	"exolifa.com/exoman/routing"
)

// version of the software
const Version = "0.0.3"

func main() {
	fmt.Println("Starting Exolifa Management")
	logger.Logme("global", "main", "not MQTT", "info", "Starting Exolifa Management")
	routeur := routing.SetRoutes()
	myport := ":" + params.Getconfig("Tcpport")
	routeur.Run(myport)
	fmt.Printf("Running gin route on port %s", myport)
	logger.Logme("global", "main", "not MQTT", "info", fmt.Sprintf("Exolifa Management listening on port %s", myport))
}
