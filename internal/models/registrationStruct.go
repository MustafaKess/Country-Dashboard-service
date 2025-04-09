package models

import (
	"Country-Dashboard-Service/internal/utils"
)

// Features struct holds the options for the dashboard's configuration
type Features struct {
	Temperature      bool     `json:"temperature" firestore:"temperature"`            // Show temperature in Celsius
	Precipitation    bool     `json:"precipitation" firestore:"precipitation"`        // Show precipitation
	Capital          bool     `json:"capital" firestore:"capital"`                    // Show the capital city
	Coordinates      bool     `json:"coordinates" firestore:"coordinates"`            // Show coordinates (latitude, longitude)
	Population       bool     `json:"population" firestore:"population"`              // Show population
	Area             bool     `json:"area" firestore:"area"`                          // Show land area
	TargetCurrencies []string `json:"targetCurrencies" firestore:"target_currencies"` // List of target currencies for exchange rates
}

// Registration represents the configuration of a registered dashboard
type Registration struct {
	ID         string           `json:"id,omitempty" firestore:"id,omitempty"` // Unique identifier for the configuration
	Country    string           `json:"country" firestore:"country"`           // Country name (alternatively to ISO code)
	IsoCode    string           `json:"isoCode" firestore:"iso_code"`          // ISO 2-letter code for the country
	Features   Features         `json:"features" firestore:"features"`         // Features to be displayed on the dashboard
	LastChange utils.CustomTime `json:"lastChange" firestore:"last_change"`    // Timestamp of the last change
	//URL        string           `json:"url" firestore:"url"`
}
