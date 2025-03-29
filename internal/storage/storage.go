package storage

// Storage logic
// e.g. firestore

import (
	"Country-Dashboard-Service/internal/models"
	"errors"
)

// Fake in-memory database until Firestore is up and running
var configs = map[string]models.DashboardConfig{
	"516dba7f015f2a68": {
		ID:               "516dba7f015f2a68",
		Country:          "Norway",
		TargetCurrencies: []string{"EUR", "USD", "SEK"},
	},
}

// ErrConfigNotFound is returned when a dashboard config with the given ID does not exist.
var ErrConfigNotFound = errors.New("dashboard config not found")

// GetDashboardConfigByID returns the dashboard config for a given ID, or an error if not found.
func GetDashboardConfigByID(id string) (*models.DashboardConfig, error) {
	config, exists := configs[id]
	if !exists {
		return nil, ErrConfigNotFound
	}
	return &config, nil
}
