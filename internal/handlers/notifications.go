package handlers

import (
	"Country-Dashboard-Service/constants/errorMessages"
	"Country-Dashboard-Service/internal/firestore"
	"Country-Dashboard-Service/internal/models"
	"Country-Dashboard-Service/internal/utils"
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

/*
NotificationsHandler routes HTTP requests for managing webhook registrations.
It supports GET, POST, and DELETE methods.
*/
func NotificationsHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	// If an ID is provided, handle GET or DELETE for a specific webhook.
	if len(parts) > 4 && parts[4] != "" {
		switch r.Method {
		case http.MethodGet:
			getSpecificNotification(w, parts[4])
		case http.MethodDelete:
			deleteNotification(w, parts[4])
		default:
			http.Error(w, errorMessages.MethodNotAllowed, http.StatusMethodNotAllowed)
		}
		return
	}

	switch r.Method {
	case http.MethodPost:
		postNotification(w, r)
	case http.MethodGet:
		getAllNotifications(w, r)
	default:
		http.Error(w, errorMessages.MethodNotAllowed, http.StatusMethodNotAllowed)
	}
}

func postNotification(w http.ResponseWriter, r *http.Request) {
	var webhook models.WebhookRegistration
	if err := json.NewDecoder(r.Body).Decode(&webhook); err != nil {
		http.Error(w, errorMessages.InvalidJSON, http.StatusBadRequest)
		return
	}
	docRef, _, err := firestore.Client.Collection("notifications").Add(context.Background(), webhook)
	if err != nil {
		http.Error(w, errorMessages.FirestoreError+err.Error(), http.StatusInternalServerError)
		return
	}
	webhook.ID = docRef.ID
	_, err = docRef.Set(context.Background(), webhook)
	if err != nil {
		http.Error(w, errorMessages.FirestoreError+err.Error(), http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"id": webhook.ID,
	}
	utils.Encode(w, http.StatusCreated, response)
}

func getSpecificNotification(w http.ResponseWriter, id string) {
	doc, err := firestore.Client.Collection("notifications").Doc(id).Get(context.Background())
	if err != nil {
		http.Error(w, errorMessages.RegisterNotFound, http.StatusNotFound)
		return
	}
	var webhook models.WebhookRegistration
	if err := doc.DataTo(&webhook); err != nil {
		http.Error(w, "Error deserializing data: "+err.Error(), http.StatusInternalServerError)
		return
	}
	utils.Encode(w, http.StatusOK, webhook)
}

func getAllNotifications(w http.ResponseWriter, r *http.Request) {
	iter := firestore.Client.Collection("notifications").Documents(context.Background())
	var webhooks []models.WebhookRegistration
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		var webhook models.WebhookRegistration
		if err := doc.DataTo(&webhook); err != nil {
			continue
		}
		webhooks = append(webhooks, webhook)
	}
	utils.Encode(w, http.StatusOK, webhooks)
}

func deleteNotification(w http.ResponseWriter, id string) {
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
		"message": "Notification deleted successfully",
		"id":      id,
	}
	utils.Encode(w, http.StatusOK, response)
}

// ClientReceiver can be used to simulate receiving webhook invocations during development.
func ClientReceiver(w http.ResponseWriter, r *http.Request) {
	var payload map[string]string
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	// Log or print payload as needed for debugging.
	w.WriteHeader(http.StatusOK)
}
