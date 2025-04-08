package main

import (
	"Country-Dashboard-Service/constants"
	"Country-Dashboard-Service/internal/firestore"
	"Country-Dashboard-Service/internal/handlers"
	"fmt"
	"net/http"
)

func main() {
	// Initialize Firestore before any requests are handled
	firestore.InitFirestore()

	// Register handlers
	http.HandleFunc(constants.Registrations, handlers.RegistrationsHandler)
	http.HandleFunc(constants.Dashboards, handlers.GetPopulatedDashboard)
	http.HandleFunc(constants.Notifications, handlers.NotificationsHandler)
	http.HandleFunc(constants.Status, handlers.StatusHandler)

	// Log server info
	fmt.Println("Starting server on port", constants.Port)
	fmt.Println("Link to the server status page: http://localhost:8080/dashboard/v1/status")
	fmt.Println("Link to the registrations page (GET-request ALL): http://localhost:8080/dashboard/v1/registrations")

	// Start the server
	http.ListenAndServe(constants.Port, nil)
}
