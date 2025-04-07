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
	fmt.Println("Link to the registrations page (GET-request ALL): http://localhost:8080/dashboard/v1/registrations")
	http.ListenAndServe(constants.Port, nil)
}
