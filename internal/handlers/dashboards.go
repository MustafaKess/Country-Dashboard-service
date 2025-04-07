package handlers

import (
	"Country-Dashboard-Service/internal/firestore"
	"Country-Dashboard-Service/internal/models"
	"Country-Dashboard-Service/internal/services"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// GetPopulatedDashboard handles GET requests for a populated dashboard by ID
func GetPopulatedDashboard(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 || parts[4] == "" {
		http.Error(w, "Missing dashboard ID", http.StatusBadRequest)
		return
	}
	id := parts[4]

	// Load full registration (from Firestore)
	config, err := firestore.GetDashboardConfigByID(id)
	if err != nil {
		http.Error(w, "Dashboard config not found", http.StatusNotFound)
		return
	}

	// Get country data
	countryInfo, err := services.GetCountryInfo(config.Country)
	if err != nil {
		http.Error(w, "Failed to fetch country data", http.StatusBadGateway)
		return
	}

	// Get weather data if requested
	var temperature float64
	var precipitation float64
	if config.Features.Temperature || config.Features.Precipitation {
		temp, precip, err := services.GetWeatherData(countryInfo.Latitude, countryInfo.Longitude)
		if err != nil {
			http.Error(w, "Failed to fetch weather data", http.StatusBadGateway)
			return
		}
		temperature = temp
		precipitation = precip
	}

	// Get currency rates if requested
	targetCurrencies := make(map[string]float64)
	if len(config.Features.TargetCurrencies) > 0 {
		rates, err := services.GetExchangeRates(countryInfo.Currency, config.Features.TargetCurrencies)
		if err != nil {
			http.Error(w, "Failed to fetch currency data", http.StatusBadGateway)
			return
		}
		targetCurrencies = rates
	}

	// Build features object based on selected options
	features := models.DashboardFeatures{}

	if config.Features.Temperature {
		features.Temperature = temperature
	}
	if config.Features.Precipitation {
		features.Precipitation = precipitation
	}
	if config.Features.Capital {
		features.Capital = countryInfo.Capital
	}
	if config.Features.Coordinates {
		features.Coordinates = models.Coordinates{
			Latitude:  countryInfo.Latitude,
			Longitude: countryInfo.Longitude,
		}
	}
	if config.Features.Population {
		features.Population = countryInfo.Population
	}
	if config.Features.Area {
		features.Area = countryInfo.Area
	}
	if len(config.Features.TargetCurrencies) > 0 {
		features.TargetCurrencies = targetCurrencies
	}

	// Build final response
	response := models.PopulatedDashboard{
		Country:       countryInfo.Name,
		ISOCode:       countryInfo.ISOCode,
		LastRetrieval: time.Now().Format("20060102 15:04"),
		Features:      features,
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
