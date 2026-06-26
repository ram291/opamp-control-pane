package main

import (
	"log"

	"github.com/ram291/opamp-control-pane/backend/internal/api"
	"github.com/ram291/opamp-control-pane/backend/internal/opamp"
)


func main() {

	log.Println("Starting OpenTelemetry Control Plane")

	opampServer := opamp.NewServer()

	apiServer := api.NewServer(opampServer)

	apiServer.Start(":8080")

}

