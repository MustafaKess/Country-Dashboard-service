package firestore

import (
	"Country-Dashboard-Service/constants/errorMessages"
	"Country-Dashboard-Service/internal/models"
	"context"
	"errors"
)

var ErrConfigNotFound = errors.New(errorMessages.DashboardConfigNotFound)

// GetDashboardConfigByID retrieves the dashboard config with the given ID from Firestore.
func GetDashboardConfigByID(id string) (*models.Registration, error) {
	doc, err := Client.Collection("registrations").Doc(id).Get(context.Background())
	if err != nil {
		return nil, ErrConfigNotFound
	}

	var config models.Registration
	if err := doc.DataTo(&config); err != nil {
		return nil, err
	}

	config.ID = doc.Ref.ID
	return &config, nil
}
