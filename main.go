package main

import (
	"Country-Dashboard-Service/internal/firestore"
	"Country-Dashboard-Service/internal/server"
	serverWebhook "Country-Dashboard-Service/internal/serverwebhook"
)

/*
Main entry point for the application.
Initializes Firestore, starts the primary HTTP server, and optionally runs the dedicated webhook server.
This service provides dashboard configurations, enriched dashboards, and webhook notifications.
*/
func main() {
	// Initialize Firestore before processing any requests.
	firestore.InitFirestore()

	// Create the primary server.
	srv := server.NewServer(":8080")
	// Start the primary server in a separate goroutine.
	go srv.Start()

	// Optionally, start the dedicated webhook server for receiving webhook invocations.
	go serverWebhook.Start()

	// Log server info
	//fmt.Println("Starting server on port", constants.Port)
	//fmt.Println("Link to the server status page: http://localhost:8080/dashboard/v1/status")
	//fmt.Println("Link to the registrations page (GET-request ALL): http://localhost:8080/dashboard/v1/registrations")
	//fmt.Println("Link to the Dashboards page (GET-request ALL): http://localhost:8080/dashboard/v1/dashboards")

	// Block forever (or implement a proper shutdown signal handler).
	select {}
}
