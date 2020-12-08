package main

import (
	"fmt"

	"exolifa.com/exoman/logger"
	"exolifa.com/exoman/params"
	"exolifa.com/exoman/routing"
)

func main() {
	fmt.Println("Starting Exolifa Management")
	//	logger.Logme("global", "main", "System", "fatal", "*******************************")
	logger.Logme("global", "main", "System", "info", "* Starting Exolifa Management *")
	//	logger.Logme("global", "main", "System", "fatal", "*******************************")
	routeur := routing.SetRoutes()
	myport := ":" + params.Getconfig("Tcpport")
	routeur.Run(myport)
	fmt.Printf("Running gin route on port %s", myport)
	logger.Logme("global", "main", "not MQTT", "info", fmt.Sprintf("Exolifa Management listening on port %s", myport))
}
