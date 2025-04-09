package handlers

import (
	"Country-Dashboard-Service/constants"
	"Country-Dashboard-Service/internal/firestore"
	"Country-Dashboard-Service/internal/models"
	"Country-Dashboard-Service/internal/utils"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Utility function to insert a test registration into Firestore
func insertTestRegistration(t *testing.T) string {
	reg := models.Registration{
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
			TargetCurrencies: []string{"USD", "EUR"},
		},
	}

	docRef, _, err := firestore.Client.Collection("registrations").Add(context.Background(), reg)
	if err != nil {
		t.Fatalf("Failed to insert test registration: %v", err)
	}
	return docRef.ID
}

func TestGetPopulatedDashboard(t *testing.T) {
	// Set up mock REST Countries API
	mockCountries := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := `[{
			"name": { "common": "Norway" },
			"capital": ["Oslo"],
			"latlng": [60.0, 10.0],
			"population": 5000000,
			"area": 385207,
			"currencies": { "NOK": { "name": "Norwegian krone" } },
			"cca2": "NO"
		}]`
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(resp))
	}))
	defer mockCountries.Close()

	// Override the actual constant temporarily
	oldCountriesAPI := constants.RestCountriesAPI
	constants.RestCountriesAPI = mockCountries.URL
	defer func() { constants.RestCountriesAPI = oldCountriesAPI }()

	// Set up mock weather API
	mockWeather := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := `{
			"hourly": {
				"temperature_2m": [10, 12, 14],
				"precipitation": [1.2, 0.5, 0.8]
			}
		}`
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(resp))
	}))
	defer mockWeather.Close()
	constants.OpenMeteoAPI = mockWeather.URL

	// Set up mock currency API
	mockCurrency := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := `{
			"base": "NOK",
			"rates": {
				"USD": 0.1,
				"EUR": 0.09
			}
		}`
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(resp))
	}))
	defer mockCurrency.Close()
	constants.CurrencyAPI = mockCurrency.URL + "/"

	// Insert test registration
	id := insertTestRegistration(t)

	// Simulate GET request
	url := "/dashboard/v1/dashboards/" + id
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()

	// Call the handler function
	GetPopulatedDashboard(rec, req)

	// Check the result
	res := rec.Result()
	defer res.Body.Close()

	// Check for success
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", res.StatusCode)
	}

	var dashboard models.PopulatedDashboard
	if err := json.NewDecoder(res.Body).Decode(&dashboard); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Validate that the dashboard returned has the correct data
	if dashboard.Country != "Norway" || dashboard.Features.Capital != "Oslo" {
		t.Errorf("Expected capital Oslo and country Norway, got %v", dashboard)
	}
}

// Test for missing dashboard ID in the URL
func TestGetPopulatedDashboard_MissingID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/dashboard/v1/dashboards/", nil)
	w := httptest.NewRecorder()

	GetPopulatedDashboard(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 Bad Request, got %d", w.Code)
	}
}

// Test for dashboard with non-existing ID
func TestGetPopulatedDashboard_NonExistingID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/dashboard/v1/dashboards/invalid-id", nil)
	w := httptest.NewRecorder()

	GetPopulatedDashboard(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 Not Found, got %d", w.Code)
	}
}

// Test for failure in the GetCountryInfo API (mock failure)
func TestGetPopulatedDashboard_CountryInfoFailure(t *testing.T) {
	// Set up mock REST Countries API to return an error
	mockCountries := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
	}))
	defer mockCountries.Close()
	constants.RestCountriesAPI = mockCountries.URL

	// Insert test registration
	id := insertTestRegistration(t)

	// Simulate GET request
	url := "/dashboard/v1/dashboards/" + id
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()

	GetPopulatedDashboard(rec, req)

	if rec.Code != http.StatusBadGateway {
		t.Errorf("Expected status 502 Bad Gateway, got %d", rec.Code)
	}
}

// Test for failure in the GetWeatherData API (mock failure)
func TestGetPopulatedDashboard_WeatherAPIError(t *testing.T) {
	// Set up mock Weather API to return an error
	mockWeather := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
	}))
	defer mockWeather.Close()
	constants.OpenMeteoAPI = mockWeather.URL

	// Insert test registration
	id := insertTestRegistration(t)

	// Simulate GET request
	url := "/dashboard/v1/dashboards/" + id
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()

	GetPopulatedDashboard(rec, req)

	if rec.Code != http.StatusBadGateway {
		t.Errorf("Expected status 502 Bad Gateway, got %d", rec.Code)
	}
}

// Test for failure in the GetExchangeRates API (mock failure)
func TestGetPopulatedDashboard_CurrencyAPIError(t *testing.T) {
	// Set up mock Currency API to return an error
	mockCurrency := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
	}))
	defer mockCurrency.Close()
	constants.CurrencyAPI = mockCurrency.URL

	// Insert test registration
	id := insertTestRegistration(t)

	// Simulate GET request
	url := "/dashboard/v1/dashboards/" + id
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()

	GetPopulatedDashboard(rec, req)

	if rec.Code != http.StatusBadGateway {
		t.Errorf("Expected status 502 Bad Gateway, got %d", rec.Code)
	}
}
