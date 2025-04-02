package handlers

import (
	"Country-Dashboard-Service/internal/models"
	"Country-Dashboard-Service/internal/storage"
	"encoding/json"
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
		var registration models.Registration // Assuming you have a Registration model
		err := json.NewDecoder(r.Body).Decode(&registration)
		if err != nil {
			http.Error(w, "Invalid JSON data", http.StatusBadRequest)
			return
		}

		// Validate fields (customize as per your schema)
		if registration.Country == "" || registration.IsoCode == "" {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		// Optionally, you can validate features or other fields as needed:
		if len(registration.Features.TargetCurrencies) == 0 {
			http.Error(w, "At least one target currency is required", http.StatusBadRequest)
			return
		}

		storage.AddDoc(w, r, "registrations")
	}
}
