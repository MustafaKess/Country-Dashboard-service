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

func NotificationHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
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
	switch r.Method {
	case http.MethodPost:
		postNotificationHandler(w, r)
	case http.MethodGet:
		getAllNotifications(w, r)
	default:
		http.Error(w, errorMessages.MethodNotAllowed, http.StatusMethodNotAllowed)
	}
}
func postNotificationHandler(w http.ResponseWriter, r *http.Request) {
	var registration models.WebhookRegistration
	if err := json.NewDecoder(r.Body).Decode(&registration); err != nil {
		http.Error(w, errorMessages.InvalidJSON, http.StatusBadRequest)
		return
	}
	docRef, _, err := firestore.Client.Collection("notifications").Add(context.Background(), registration)
	if err != nil {
		http.Error(w, errorMessages.FirestoreError+err.Error(), http.StatusInternalServerError)
	}
	registration.ID = docRef.ID
	_, err = docRef.Set(context.Background(), registration)
	if err != nil {
		http.Error(w, errorMessages.FirestoreError+err.Error(), http.StatusInternalServerError)
	}
	response := map[string]interface{}{
		"id": registration.ID,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
