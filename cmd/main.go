package main

import (
	"flag"
	"fmt"
	"log"
	"netBlast/cmd/client"
	"netBlast/cmd/server"
	"netBlast/pkg"

	"os"
)

// Sets flags
func setFlags() string {
	arg := flag.String("arg", "", "Specify which side to run")

	flag.Parse()

	if *arg == "" {
		log.Fatal("Didnt provide argument to run")
	}
	return *arg
}

// Main function
func main() {
	//arg := os.Args[1]
	arg := setFlags()

	switch arg {
	case "server":
		fmt.Println("Running server on: ", pkg.ServerURL)

		shutdown := make(chan os.Signal)

		server.Server(shutdown)
	case "client":
		client.Client()
	}
}
