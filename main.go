package main

import (
	"Country-Dashboard-Service/constants"
	"Country-Dashboard-Service/internal/handlers"
	"fmt"
	"net/http"
)

func main() {
	//http.HandleFunc(constants.Registrations, handlers.RegistrationHandler)
	http.HandleFunc(constants.Dashboards, handlers.GetPopulatedDashboard)
	//http.HandleFunc(constants.Notifications, handlers.NotificationHandler)
	http.HandleFunc(constants.Status, handlers.StatusHandler)

	fmt.Println("Starting server on port", constants.Port)
	fmt.Println("Link to the server status page: http://localhost:8080/dashboard/v1/status")
	http.ListenAndServe(constants.Port, nil)
}
