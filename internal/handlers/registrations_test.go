package handlers

import (
	"Country-Dashboard-Service/constants"
	"Country-Dashboard-Service/internal/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Utility function to create a valid registration JSON payload
func validRegistrationPayload() []byte {
	reg := models.Registration{
		Country: "Norway",
		IsoCode: "NO",
		Features: models.Features{
			Temperature:      true,
			Precipitation:    false,
			Capital:          true,
			Coordinates:      true,
			Population:       true,
			Area:             false,
			TargetCurrencies: []string{"USD", "EUR"},
		},
	}
	payload, _ := json.Marshal(reg)
	return payload
}

// Setup mocked REST Countries API
func startMockCountryAPI(t *testing.T) func() {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"cca2": "NO"}]`)) // ISO match
	}))
	constants.RestCountriesAPI = mock.URL
	return mock.Close
}

func TestPostValidRegistration(t *testing.T) {
	closeMock := startMockCountryAPI(t)
	defer closeMock()

	req := httptest.NewRequest(http.MethodPost, constants.Registrations, bytes.NewReader(validRegistrationPayload()))
	w := httptest.NewRecorder()

	RegistrationsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}
}

func TestPostInvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, constants.Registrations, bytes.NewReader([]byte(`invalid-json`)))
	w := httptest.NewRecorder()

	RegistrationsHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d", w.Code)
	}
}

func TestPostInvalidISO(t *testing.T) {
	// Mock returns wrong cca2 to trigger ISO mismatch
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[{"cca2": "SE"}]`)) // wrong ISO
	}))
	defer mock.Close()
	constants.RestCountriesAPI = mock.URL

	reg := models.Registration{
		Country: "Norway",
		IsoCode: "NO", // mismatch with "SE"
	}
	payload, _ := json.Marshal(reg)

	req := httptest.NewRequest(http.MethodPost, constants.Registrations, bytes.NewReader(payload))
	w := httptest.NewRecorder()

	RegistrationsHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request due to ISO mismatch, got %d", w.Code)
	}
}

func TestGetAllRegistrations(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, constants.Registrations, nil)
	w := httptest.NewRecorder()

	RegistrationsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}
}

func TestInvalidMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodPatch, constants.Registrations, nil)
	w := httptest.NewRecorder()

	RegistrationsHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405 Method Not Allowed, got %d", w.Code)
	}
}

func TestPutRegistration_ValidUpdate(t *testing.T) {
	// Insert a test registration first
	id := insertTestRegistration(t)

	// Prepare updated payload
	updated := models.Registration{
		Country: "Norway",
		IsoCode: "NO",
		Features: models.Features{
			Temperature:      false,
			Precipitation:    true,
			Capital:          false,
			Coordinates:      true,
			Population:       true,
			Area:             true,
			TargetCurrencies: []string{"USD"},
		},
	}
	payload, _ := json.Marshal(updated)

	req := httptest.NewRequest(http.MethodPut, constants.Registrations+id, bytes.NewReader(payload))
	w := httptest.NewRecorder()

	RegistrationsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp["message"] != "Registration updated successfully" {
		t.Errorf("Unexpected message: %v", resp["message"])
	}
}

func TestPutRegistration_MissingID(t *testing.T) {
	payload := validRegistrationPayload() // from earlier helper
	req := httptest.NewRequest(http.MethodPut, constants.Registrations, bytes.NewReader(payload))
	w := httptest.NewRecorder()

	RegistrationsHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request, got %d", w.Code)
	}
}

func TestPutRegistration_InvalidJSON(t *testing.T) {
	id := insertTestRegistration(t)

	req := httptest.NewRequest(http.MethodPut, constants.Registrations+id, bytes.NewReader([]byte("not-json")))
	w := httptest.NewRecorder()

	RegistrationsHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request for invalid JSON, got %d", w.Code)
	}
}

func TestPutRegistration_NonExistentID(t *testing.T) {
	closeMock := startMockCountryAPI(t)
	defer closeMock()

	// Use a fake ID that doesn't exist
	nonExistentID := "fake-id-12345"

	// Create a valid updated payload
	updated := models.Registration{
		Country: "Norway",
		IsoCode: "NO",
		Features: models.Features{
			Temperature:      true,
			Precipitation:    true,
			Capital:          true,
			Coordinates:      true,
			Population:       true,
			Area:             true,
			TargetCurrencies: []string{"USD"},
		},
	}
	payload, _ := json.Marshal(updated)

	req := httptest.NewRequest(http.MethodPut, constants.Registrations+nonExistentID, bytes.NewReader(payload))
	w := httptest.NewRecorder()

	RegistrationsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 OK even for non-existent ID (upsert), got %d", w.Code)
	}

	// Check the response body for confirmation
	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if resp["message"] != "Registration updated successfully" {
		t.Errorf("Unexpected message: %v", resp["message"])
	}
}

func TestDeleteRegistration_AlreadyDeleted(t *testing.T) {
	// Insert a registration
	id := insertTestRegistration(t)

	// Delete it once — should succeed
	req1 := httptest.NewRequest(http.MethodDelete, constants.Registrations+id, nil)
	w1 := httptest.NewRecorder()
	RegistrationsHandler(w1, req1)

	if w1.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK on first delete, got %d", w1.Code)
	}

	// Try deleting again — should now return 404 Not Found
	req2 := httptest.NewRequest(http.MethodDelete, constants.Registrations+id, nil)
	w2 := httptest.NewRecorder()
	RegistrationsHandler(w2, req2)

	if w2.Code != http.StatusNotFound {
		t.Errorf("Expected 404 Not Found on second delete, got %d", w2.Code)
	}
}

func TestGetSpecificRegistration(t *testing.T) {
	// Insert a registration
	id := insertTestRegistration(t)

	req := httptest.NewRequest(http.MethodGet, constants.Registrations+id, nil)
	w := httptest.NewRecorder()

	RegistrationsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", w.Code)
	}

	var reg models.Registration
	if err := json.NewDecoder(w.Body).Decode(&reg); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if reg.ID != id {
		t.Errorf("Expected ID %s, got %s", id, reg.ID)
	}
}

func TestGetSpecificRegistration_NotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, constants.Registrations+"nonexistent-id", nil)
	w := httptest.NewRecorder()

	RegistrationsHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected 404 Not Found, got %d", w.Code)
	}
}
