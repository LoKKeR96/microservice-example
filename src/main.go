package main

import (
	"log"

	"net/http"

	"github.com/lokker96/microservice_example/infrastructure/container"
	"github.com/lokker96/microservice_example/infrastructure/route"
)

// The main function - entry point of the application
func main() {
	// Building the application's container
	c, err := container.NewContainer() // Create a new container instance
	if err != nil {
		log.Fatal(err) // If there's an error, log it and terminate
	}

	// Configuring the server and initializing routes
	server := http.Server{
		Addr:    ":8080",          // Set the address on which the server will listen
		Handler: route.Routes(*c), // Set up the routes using the created container
	}

	// Start the server
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server failed: %v", err) // If there's an error starting the server, log it and terminate
	}
}
