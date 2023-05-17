package main

import (
	"fmt"
	"netBlast/cmd/client"
	"netBlast/cmd/server"
	"netBlast/pkg"

	"os"
)

// Main function
func main() {
	arg := os.Args[1]

	switch arg {
	case "server":
		fmt.Println("Running server on: ", pkg.ServerURL)
		server.Server()
	case "client":
		client.Client()
	}
}
