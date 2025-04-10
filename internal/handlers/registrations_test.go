package handlers

import (
	"Country-Dashboard-Service/constants"
	"Country-Dashboard-Service/internal/models"
	"Country-Dashboard-Service/internal/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Setup mocked REST Countries API
func startMockCountryAPI(t *testing.T, iso string) func() {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Logf("Mock Country API called with URL: %s", r.URL.String()) // log this
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"cca2": "` + iso + `"}]`))
	}))
	constants.RestCountriesAPI = mock.URL
	return mock.Close
}

// Utility function to create a valid registration JSON payload
func validRegistrationPayload() []byte {
	type registrationPayload struct {
		Country  string          `json:"country"`
		IsoCode  string          `json:"isoCode"`
		Features models.Features `json:"features"`
	}

	reg := registrationPayload{
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

func TestPostValidRegistration(t *testing.T) {
	t.Parallel()

	closeMock := startMockCountryAPI(t, "NO")
	defer closeMock()

	payload := validRegistrationPayload()

	fmt.Println("Sending registration payload:", string(payload))

	req := httptest.NewRequest(http.MethodPost, constants.Registrations, bytes.NewReader(payload))
	w := httptest.NewRecorder()

	RegistrationsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}
}

func TestPostInvalidJSON(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodPost, constants.Registrations, bytes.NewReader([]byte(`invalid-json`)))
	w := httptest.NewRecorder()

	RegistrationsHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d", w.Code)
	}
}

func TestPostInvalidISO(t *testing.T) {
	t.Parallel()

	closeMock := startMockCountryAPI(t, "SE") // Mismatched ISO
	defer closeMock()

	reg := models.Registration{
		Country: "Norway",
		IsoCode: "NO",
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
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, constants.Registrations, nil)
	w := httptest.NewRecorder()

	RegistrationsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}
}

func TestInvalidMethod(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodPatch, constants.Registrations, nil)
	w := httptest.NewRecorder()

	RegistrationsHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405 Method Not Allowed, got %d", w.Code)
	}
}

func TestPutRegistration_ValidUpdate(t *testing.T) {
	t.Parallel()

	closeMock := startMockCountryAPI(t, "NO")
	defer closeMock()

	id := insertTestRegistration(t)

	updated := models.Registration{
		Country:    "Norway",
		IsoCode:    "NO",
		LastChange: utils.CustomTime{Time: time.Now()},
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
	t.Parallel()

	payload := validRegistrationPayload()
	req := httptest.NewRequest(http.MethodPut, constants.Registrations, bytes.NewReader(payload))
	w := httptest.NewRecorder()

	RegistrationsHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request, got %d", w.Code)
	}
}

func TestPutRegistration_InvalidJSON(t *testing.T) {
	t.Parallel()

	id := insertTestRegistration(t)

	req := httptest.NewRequest(http.MethodPut, constants.Registrations+id, bytes.NewReader([]byte("not-json")))
	w := httptest.NewRecorder()

	RegistrationsHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request for invalid JSON, got %d", w.Code)
	}
}

func TestPutRegistration_NonExistentID(t *testing.T) {
	t.Parallel()

	closeMock := startMockCountryAPI(t, "NO")
	defer closeMock()

	nonExistentID := "fake-id-12345"

	updated := models.Registration{
		Country:    "Norway",
		IsoCode:    "NO",
		LastChange: utils.CustomTime{Time: time.Now()},
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
		return
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if resp["message"] != "Registration updated successfully" {
		t.Errorf("Unexpected message: %v", resp["message"])
	}
}

func TestDeleteRegistration_AlreadyDeleted(t *testing.T) {
	t.Parallel()

	id := insertTestRegistration(t)

	req1 := httptest.NewRequest(http.MethodDelete, constants.Registrations+id, nil)
	w1 := httptest.NewRecorder()
	RegistrationsHandler(w1, req1)

	if w1.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK on first delete, got %d", w1.Code)
	}

	req2 := httptest.NewRequest(http.MethodDelete, constants.Registrations+id, nil)
	w2 := httptest.NewRecorder()
	RegistrationsHandler(w2, req2)

	if w2.Code != http.StatusNotFound {
		t.Errorf("Expected 404 Not Found on second delete, got %d", w2.Code)
	}
}

func TestGetSpecificRegistration(t *testing.T) {
	t.Parallel()

	id := insertTestRegistration(t)

	req := httptest.NewRequest(http.MethodGet, constants.Registrations+id, nil)
	w := httptest.NewRecorder()

	RegistrationsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", w.Code)
	}

	var raw map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&raw); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if raw["id"] != id {
		t.Errorf("Expected ID %s, got %v", id, raw["id"])
	}
}

func TestGetSpecificRegistration_NotFound(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, constants.Registrations+"nonexistent-id", nil)
	w := httptest.NewRecorder()

	RegistrationsHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected 404 Not Found, got %d", w.Code)
	}
}
