package main

import (
	"fmt"

	"exolifa.com/exoman/params"
	"exolifa.com/exoman/routing"
)

func main() {
	fmt.Println("Starting Exolifa Management")
	routeur := routing.SetRoutes()
	myport := ":" + params.Getconfig("Tcpport")
	routeur.Run(myport)
	fmt.Printf("Running gin route on port 8999")

}
