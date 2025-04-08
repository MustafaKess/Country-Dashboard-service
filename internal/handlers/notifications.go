package handlers

import (
	"Country-Dashboard-Service/constants/errorMessages"
	"Country-Dashboard-Service/internal/firestore"
	"Country-Dashboard-Service/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
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
	var registration models.Registration
	if err := json.NewDecoder(r.Body).Decode(&registration); err != nil {
		http.Error(w, errorMessages.InvalidJSON, http.StatusBadRequest)
		return
	}

	// Adds the registration to Firestore.
	docRef, _, err := firestore.Client.Collection("notifications").Add(context.Background(), registration)
	if err != nil {
		http.Error(w, errorMessages.FirestoreError+err.Error(), http.StatusInternalServerError)
		return
	}

	// Updates registration with document ID.
	registration.ID = docRef.ID
	_, err = docRef.Set(context.Background(), registration)
	if err != nil {
		http.Error(w, errorMessages.FirestoreError+err.Error(), http.StatusInternalServerError)
		return
	}

	// Returns the generated ID.
	response := map[string]interface{}{
		"id": registration.ID,
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
	var registration models.Registration
	if err = doc.DataTo(&registration); err != nil {
		http.Error(w, "Error deserializing data: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(registration)
}

// getAllNotificationsHandler retrieves all webhook registrations.
func getAllNotificationsHandler(w http.ResponseWriter, r *http.Request) {
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

// deleteNotificationHandler deletes a specific webhook registration.
func deleteNotificationHandler(w http.ResponseWriter, id string) {
	_, err := firestore.Client.Collection("notifications").Doc(id).Delete(context.Background())
	if err != nil {
		http.Error(w, errorMessages.DeleteError+err.Error(), http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"message": "Notification deleted successfully",
		"id":      id,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// TriggerWebhookEvent finds all webhook registrations matching the given event and optionally the country
// and sends a POST notification to the registered URL.
func TriggerWebhookEvent(event string, country string) {
	// Query webhooks where event equals the given event.
	iter := firestore.Client.Collection("notifications").Where("event", "==", event).Documents(context.Background())
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		var reg models.Registration
		if err = doc.DataTo(&reg); err != nil {
			log.Printf("Failed to deserialize webhook registration: %v", err)
			continue
		}
		// If a country is specified in the registration and it doesn't match, skip.
		if reg.Country != "" && reg.Country != country {
			continue
		}
		// Prep payload.
		payload := map[string]string{
			"id":      reg.ID,
			"country": country,
			"event":   event,
			"time":    time.Now().Format("20060102 15:04"),
		}
		// Send the webhook invocation.
		go sendWebhookNotification(reg.URL, payload)
	}
}

// sendWebhookNotification sends a POST request with the payload to the provided URL.
func sendWebhookNotification(url string, payload map[string]string) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling webhook payload: %v", err)
		return
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Error creating webhook request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending webhook to %s: %v", url, err)
		return
	}
	defer resp.Body.Close()
	log.Printf("Webhook sent to %s with status code %d", url, resp.StatusCode)
}
