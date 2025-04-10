package models

import "Country-Dashboard-Service/internal/utils"

// Represents a saved dashboard setup with a country and target currencies.
type DashboardConfig struct {
	ID               string   `json:"id"`
	Country          string   `json:"country"`
	TargetCurrencies []string `json:"targetCurrencies"`
}

// Full dashboard response sent to the client with enriched data.
type PopulatedDashboard struct {
	Country       string            `json:"country"`
	ISOCode       string            `json:"isoCode"`
	Features      DashboardFeatures `json:"features"`
	LastRetrieval utils.CustomTime  `json:"lastRetrieval"`
}

// Contains detailed information shown in the dashboard.
type DashboardFeatures struct {
	Temperature      float64            `json:"temperature,omitempty"`      // Exclude if zero
	Precipitation    float64            `json:"precipitation,omitempty"`    // Exclude if zero
	Capital          string             `json:"capital,omitempty"`          // Exclude if empty
	Coordinates      Coordinates        `json:"coordinates,omitempty"`      // Exclude if zero
	Population       int                `json:"population,omitempty"`       // Exclude if zero
	Area             float64            `json:"area,omitempty"`             // Exclude if zero
	TargetCurrencies map[string]float64 `json:"targetCurrencies,omitempty"` // Exclude if empty
}

// Holds latitude and longitude values for a country.
type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
