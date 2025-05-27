package main

import (
	"log"
	"os"

	"github.com/lokker96/microservice_example/infrastructure/container"
	"github.com/lokker96/microservice_example/infrastructure/route"
)

const defaultPort = "8080"

// The main function - entry point of the application
func main() {
	// Building the application's container
	c, err := container.NewContainer() // Create a new container instance
	if err != nil {
		log.Fatal(err) // If there's an error, log it and terminate
	}

	defer c.HandleShutdown()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	e := route.SetupRoutes(*c)

	// Start the server
	if err := e.Start(":" + port); err != nil {
		log.Fatalf("server failed: %v", err) // If there's an error starting the server, log it and terminate
	}
}
