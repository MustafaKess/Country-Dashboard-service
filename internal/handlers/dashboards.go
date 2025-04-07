package handlers

import (
	"Country-Dashboard-Service/internal/models"
	"Country-Dashboard-Service/internal/services"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// Handles GET requests for a specific dashboard ID
func GetPopulatedDashboard(w http.ResponseWriter, r *http.Request) {
	// Get dashboard ID from path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 || parts[4] == "" {
		http.Error(w, "Missing dashboard ID", http.StatusBadRequest)
		return
	}
	id := parts[4]

	// Load config
	config, err := storage.GetDashboardConfigByID(id)
	if err != nil {
		http.Error(w, "Dashboard not found", http.StatusNotFound)
		return
	}

	// Get country info
	countryInfo, err := services.GetCountryInfo(config.Country)
	if err != nil {
		http.Error(w, "Failed to fetch country data", http.StatusBadGateway)
		return
	}

	// Get weather info
	temp, precip, err := services.GetWeatherData(countryInfo.Latitude, countryInfo.Longitude)
	if err != nil {
		http.Error(w, "Failed to fetch weather data", http.StatusBadGateway)
		return
	}

	// TODO: Fetch currency info
	targetCurrencies := make(map[string]float64) // placeholder

	// Build response
	response := models.PopulatedDashboard{
		Country:       countryInfo.Name,
		ISOCode:       countryInfo.ISOCode,
		LastRetrieval: time.Now().Format("20060102 15:04"),
		Features: models.DashboardFeatures{
			Temperature:   temp,
			Precipitation: precip,
			Capital:       countryInfo.Capital,
			Coordinates: models.Coordinates{
				Latitude:  countryInfo.Latitude,
				Longitude: countryInfo.Longitude,
			},
			Population:       countryInfo.Population,
			Area:             countryInfo.Area,
			TargetCurrencies: targetCurrencies,
		},
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
