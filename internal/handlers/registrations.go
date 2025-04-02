package handlers

import (
	"Country-Dashboard-Service/internal/storage"
	"net/http"
)

func RegistrationsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getRegistrationsHandler(w, r)
	case http.MethodPost:
		postRegistrationsHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

	}

}

// getRegistrationsHandler retrieves all the documents from the "registrations" collection
func getRegistrationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		storage.DisplayConfig(w, r)
	}
}

func postRegistrationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		storage.AddDoc(w, r, "registrations")
	}
}
