package server

import (
	"Country-Dashboard-Service/constants"
	"Country-Dashboard-Service/internal/handlers"
	"log"
	"net/http"
	"time"
)

/*
Server represents the main application HTTP server.
It uses dependency injection to register endpoints and allows for easier testing.
*/
type Server struct {
	HTTP *http.Server
}

// NewServer creates a new Server instance with preconfigured routes.
func NewServer(port string) *Server {
	mux := http.NewServeMux()
	// Register endpoint handlers.
	mux.HandleFunc(constants.Registrations, handlers.RegistrationsHandler)
	mux.HandleFunc(constants.Dashboards, handlers.GetPopulatedDashboard)
	mux.HandleFunc(constants.Notifications, handlers.NotificationsHandler)
	mux.HandleFunc(constants.Status, handlers.StatusHandler)
	// Endpoint to receive webhook callbacks.
	mux.HandleFunc("/dashboard/v1/client/", handlers.ClientReceiver)

	srv := &http.Server{
		Addr:         port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
	return &Server{HTTP: srv}
}

// Start launches the HTTP server.
func (s *Server) Start() {
	log.Println("Starting server on", s.HTTP.Addr)
	if err := s.HTTP.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// Shutdown gracefully stops the HTTP server.
func (s *Server) Shutdown() {
	if err := s.HTTP.Close(); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}
}
