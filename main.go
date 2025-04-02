package main

import (
	"Country-Dashboard-Service/constants"
	"Country-Dashboard-Service/internal/handlers"
	"Country-Dashboard-Service/internal/storage"

	"fmt"
	"net/http"
)

func main() {

	// InitFirestore initializes the Firestore client
	storage.InitFirestore()

	http.HandleFunc(constants.Registrations, handlers.RegistrationsHandler)
	//http.HandleFunc(constants.Dashboards, handlers.DashboardHandler)
	//http.HandleFunc(constants.Notifications, handlers.NotificationHandler)
	http.HandleFunc(constants.Status, handlers.StatusHandler)

	fmt.Println("Starting server on port", constants.Port)
	fmt.Println("Link to the server status page: http://localhost:8080/dashboard/v1/status")
	http.ListenAndServe(constants.Port, nil)
}
