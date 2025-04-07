package handlers

import (
	"Country-Dashboard-Service/constants"
	"Country-Dashboard-Service/constants/errorMessages"
	"Country-Dashboard-Service/internal/firestore"
	"Country-Dashboard-Service/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

//
// Error codes scenarios:
// 400: Bad Request
// 404: Not Found
// 405: Method Not Allowed
// 500: Internal Server Error
//

// RegistrationsHandler handles the main logic for the /registrations endpoint.
// It distinguishes between GET and POST requests.
func RegistrationsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Handles GET requests to retrieve registrations.
		getRegistrationsHandler(w, r)
	case http.MethodPost:
		// Handles POST requests to create new registrations.
		postRegistrationsHandler(w, r)
	case http.MethodDelete:
		// Handles DELETE requests to remove registrations.
		deleteRegistration(w, r)
	case http.MethodPut:
		// Handles PUT requests to update existing registrations.
		putRegistration(w, r)

	default:
		// If method is not allowed, return a 405 error.
		http.Error(w, errorMessages.MethodNotAllowed, http.StatusMethodNotAllowed)
	}
}

// PostRegistrationsHandler processes a POST request to create a new registration.
// It expects the body to be a JSON object representing a registration.
func postRegistrationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Define a variable to hold the registration data.
		var registration models.Registration

		// Decode the incoming JSON data into the registration model.
		err := json.NewDecoder(r.Body).Decode(&registration)
		if err != nil {
			http.Error(w, errorMessages.InvalidJSON, http.StatusBadRequest)
			return
		}

		if err := ValidateRegistration(registration); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Add the registration to Firestore and retrieve the document reference.
		// Do not change the path to the collection, its the path to the collection in Firestore.
		docR, _, err := firestore.Client.Collection("registrations").Add(context.Background(), registration)
		if err != nil {
			http.Error(w, errorMessages.FirestoreError+err.Error(), http.StatusInternalServerError)
			return
		}

		// Get the document ID and update the registration's LastChange timestamp.
		id := docR.ID
		registration.ID = id
		registration.LastChange = time.Now()

		// Update the Firestore document with the registration data.
		_, err = docR.Set(context.Background(), registration)
		if err != nil {
			http.Error(w, errorMessages.InvalidRegistrationID+err.Error(), http.StatusInternalServerError)
			return
		}

		// Return the ID and LastChange time in the response. Confirmation message in JSON for the client.
		response := map[string]interface{}{
			"id":         id,
			"lastChange": registration.LastChange,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// getRegistrationsHandler processes GET requests for the /registrations endpoint.
// It checks if an ID is provided to fetch a specific registration or returns all registrations.
func getRegistrationsHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the registration ID from the URL path
	parts := strings.Split(r.URL.Path, "/")

	// Check if an ID exists after "/dashboard/v1/registrations/"
	if len(parts) > 4 && parts[4] != "" {
		// If ID is provided, fetch the specific registration.
		getSpecifiedRegistration(w, parts[4])
		return
	}

	// If no ID is provided, fetch all registrations.
	getAllRegistrations(w, r)
}

// GetSpecifiedRegistration fetches a specific registration from Firestore based on the given ID.
func getSpecifiedRegistration(w http.ResponseWriter, id string) {
	// Fetch the document from Firestore using the provided ID.
	doc, err := firestore.Client.Collection("registrations").Doc(id).Get(context.Background())
	if err != nil {
		// If the document is not found, return a 404 error.
		http.Error(w, errorMessages.RegisterNotFound, http.StatusNotFound)
		return
	}

	// Map the Firestore document data to the Registration model.
	var reg models.Registration
	err = doc.DataTo(&reg)
	if err != nil {
		// If there's an error deserializing the data, return a 500 error.
		http.Error(w, "Error with deserialization: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the registration as a JSON response.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reg)
}

// GetAllRegistrations retrieves all registrations from Firestore.
func getAllRegistrations(w http.ResponseWriter, r *http.Request) {
	// Fetch all documents in the "registrations" collection.
	iter := firestore.Client.Collection("registrations").Documents(context.Background())
	var all []models.Registration

	// Loop through all the documents and map them to the Registration model.
	for {
		doc, err := iter.Next()
		if err != nil {
			// If there are no more documents, exit the loop.
			break
		}
		var reg models.Registration
		err = doc.DataTo(&reg)
		if err != nil {
			// Skip broken documents.
			continue
		}
		reg.ID = doc.Ref.ID
		// Add valid registrations to the list.
		all = append(all, reg)
	}

	// Return all registrations as a JSON response.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(all)
}

// deleteRegistration deletes a specific registration from Firestore.
func deleteRegistration(w http.ResponseWriter, r *http.Request) {
	// Extract the registration ID from the URL path
	parts := strings.Split(r.URL.Path, "/")

	// Check if an ID exists after "/dashboard/v1/registrations/"
	if len(parts) > 4 && parts[4] != "" {
		// If ID is provided, delete the specific registration.
		id := parts[4]

		// Attempt to delete the document from Firestore
		_, err := firestore.Client.Collection("registrations").Doc(id).Delete(context.Background())
		if err != nil {
			// If there's an error, return a 500 error response
			http.Error(w, errorMessages.DeleteError+err.Error(), http.StatusInternalServerError)
			return
		}

		// Return a JSON response confirming the deletion and showing the ID of the deleted registration.
		response := map[string]interface{}{
			"message": "Registration deleted successfully",
			"id":      id,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	} else {
		// If no ID is provided in the URL, return a 400 error.
		http.Error(w, errorMessages.NoIDProvided, http.StatusBadRequest)
	}
}

// putRegistration updates an existing registration in Firestore.
func putRegistration(w http.ResponseWriter, r *http.Request) {
	// Extract the registration ID from the URL path
	parts := strings.Split(r.URL.Path, "/")

	// Check if an ID exists after "/dashboard/v1/registrations/"
	if len(parts) > 4 && parts[4] != "" {
		// If ID is provided, update the specific registration.
		id := parts[4]

		// Define a variable to hold the updated registration data.
		var registration models.Registration

		// Decode the incoming JSON data into the registration model.
		err := json.NewDecoder(r.Body).Decode(&registration)
		if err != nil {
			http.Error(w, errorMessages.InvalidJSON, http.StatusBadRequest)
			return
		}

		// Update the Firestore document with the registration data.
		docR := firestore.Client.Collection("registrations").Doc(id)
		registration.LastChange = time.Now()
		_, err = docR.Set(context.Background(), registration)
		if err != nil {
			http.Error(w, errorMessages.UpdateError+err.Error(), http.StatusInternalServerError)
			return
		}

		// Respond with a confirmation message and the updated registration data.
		response := map[string]interface{}{
			"message":     "Registration updated successfully",
			"updatedData": registration,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		// If no ID is provided, return a 400 error.
		http.Error(w, errorMessages.NoIDProvided, http.StatusBadRequest)
	}
}

func ValidateRegistration(registration models.Registration) error {
	if registration.Country == "" {
		return fmt.Errorf("country name is required")
	}

	if registration.IsoCode == "" {
		return fmt.Errorf("ISO code is required")
	}

	// Delegate ISO code validation to a dedicated function
	if err := ValidateISOCode(registration.Country, registration.IsoCode); err != nil {
		return err
	}

	return nil
}

func ValidateISOCode(country string, isoCode string) error {
	// Build the request URL and perform the HTTP GET request
	apiURL := fmt.Sprintf(constants.RestCountriesAPI+"/name/%s", country)
	resp, err := http.Get(apiURL)
	if err != nil {
		return fmt.Errorf("failed to validate country with external API: %v", err)
	}
	defer resp.Body.Close()

	// Handle specific error for 404 (not found)
	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("country '%s' is not recognized", country)
	}
	// Generic error for other non-200 responses
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("external API returned unexpected status: %s", resp.Status)
	}
	// Decode the JSON response
	var apiResponse []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		return fmt.Errorf("failed to decode external API response: %v", err)
	}
	// Check for data presence and extract cca2
	if len(apiResponse) == 0 {
		return fmt.Errorf("no data found for country: %s", country)
	}
	cca2Raw, ok := apiResponse[0]["cca2"]
	if !ok {
		return fmt.Errorf("ISO code (cca2) not found in API response")
	}
	cca2, ok := cca2Raw.(string)
	if !ok {
		return fmt.Errorf("invalid ISO code format in API response")
	}

	// Compare ISO codes, case-insensitively
	if strings.ToUpper(cca2) != strings.ToUpper(isoCode) {
		return fmt.Errorf("ISO code '%s' does not match country '%s' (expected '%s')", isoCode, country, cca2)
	}

	return nil
}
