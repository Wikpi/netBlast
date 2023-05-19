package main

import (
	"flag"
	"log"
	"netBlast/cmd/client"
	"netBlast/cmd/server"

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
		shutdown := make(chan os.Signal)

		server.Server(shutdown)
	case "client":
		client.Client()
	}
}
