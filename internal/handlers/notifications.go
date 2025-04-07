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
	var registration models.Registration
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
func getSpecificNotificationHandler(w http.ResponseWriter, id string) {
	doc, err := firestore.Client.Collection("notifications").Doc(id).Get(context.Background())
	if err != nil {
		http.Error(w, errorMessages.RegisterNotFound, http.StatusNotFound)
		return
	}
	var registration models.Registration
	if err = doc.DataTo(&registration); err != nil {
		http.Error(w, "Error deserialising data"+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(registration)
}
func getAllNotifications(w http.ResponseWriter, r *http.Request) {
	iter := firestore.Client.Collection("notifications").Documents(context.Background())
	var registrations []models.Registration
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		var reg models.Registration
		if err = doc.DataTo(&reg); err != nil {
			continue
		}
		registrations = append(registrations, reg)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(registrations)
}
func deleteNotificationHandler(w http.ResponseWriter, id string) {
	_, err := firestore.Client.Collection("notifications").Doc(id).Delete(context.Background())
	if err != nil {
		http.Error(w, errorMessages.DeleteError+err.Error(), http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"message": "message deleted successfully",
		"id":      id,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
