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
