package handlers

import (
	"Country-Dashboard-Service/constants/errorMessages"
	"Country-Dashboard-Service/internal/firestore"
	"Country-Dashboard-Service/internal/models"
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

// NotificationsHandler routes requests for /dashboard/v1/notifications/ endpoints.
func NotificationsHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	// If an ID is provided in the URL, handle GET or DELETE for a specific webhook.
	if len(parts) > 4 && parts[4] != "" {
		switch r.Method {
		case http.MethodGet:
			getSpecificNotificationHandler(w, parts[4])
		case http.MethodDelete:
			deleteNotificationHandler(w, parts[4])
		default:
			http.Error(w, errorMessages.MethodNotAllowed, http.StatusMethodNotAllowed)
		}
		return
	}

	// No ID provided – handle registration (POST) or list all (GET)
	switch r.Method {
	case http.MethodPost:
		postNotificationHandler(w, r)
	case http.MethodGet:
		getAllNotificationsHandler(w, r)
	default:
		http.Error(w, errorMessages.MethodNotAllowed, http.StatusMethodNotAllowed)
	}
}

// postNotificationHandler registers a new webhook notification.
func postNotificationHandler(w http.ResponseWriter, r *http.Request) {
	var webhook models.WebhookRegistration
	if err := json.NewDecoder(r.Body).Decode(&webhook); err != nil {
		http.Error(w, errorMessages.InvalidJSON, http.StatusBadRequest)
		return
	}

	// Adds the registration to Firestore.
	docRef, _, err := firestore.Client.Collection("notifications").Add(context.Background(), webhook)
	if err != nil {
		http.Error(w, errorMessages.FirestoreError+err.Error(), http.StatusInternalServerError)
		return
	}

	// Updates registration with document ID.
	webhook.ID = docRef.ID
	_, err = docRef.Set(context.Background(), webhook)
	if err != nil {
		http.Error(w, errorMessages.FirestoreError+err.Error(), http.StatusInternalServerError)
		return
	}

	// Returns the generated ID.
	response := map[string]interface{}{
		"id": webhook.ID,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// getSpecificNotificationHandler retrieves a specific webhook registration.
func getSpecificNotificationHandler(w http.ResponseWriter, id string) {
	doc, err := firestore.Client.Collection("notifications").Doc(id).Get(context.Background())
	if err != nil {
		http.Error(w, errorMessages.RegisterNotFound, http.StatusNotFound)
		return
	}
	var webhook models.WebhookRegistration
	if err = doc.DataTo(&webhook); err != nil {
		http.Error(w, errorMessages.DeserializationError+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(webhook)
}

// getAllNotificationsHandler retrieves all webhook registrations.
func getAllNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	iter := firestore.Client.Collection("notifications").Documents(context.Background())
	var webhooks []models.WebhookRegistration
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		var webhook models.WebhookRegistration
		if err = doc.DataTo(&webhook); err != nil {
			continue
		}
		webhooks = append(webhooks, webhook)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(webhooks)
}

// deleteNotificationHandler deletes a specific webhook registration.
func deleteNotificationHandler(w http.ResponseWriter, id string) {
	docRef := firestore.Client.Collection("notifications").Doc(id)

	// Check if the document exists
	_, err := docRef.Get(context.Background())
	if err != nil {
		http.Error(w, errorMessages.NotificationNotFound, http.StatusNotFound)
		return
	}

	// Proceed with deletion
	_, err = docRef.Delete(context.Background())
	if err != nil {
		http.Error(w, errorMessages.DeleteError+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": errorMessages.NotificationDeleted,
		"id":      id,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
