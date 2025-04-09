package handlers

import (
	"Country-Dashboard-Service/constants/errorMessages"
	"Country-Dashboard-Service/internal/firestore"
	"Country-Dashboard-Service/internal/models"
	"context"
	"encoding/json"
	"google.golang.org/api/iterator"
	"log"
	"net/http"
	"strings"
	"time"
)

// NotificationsHandler routes requests for /dashboard/v1/notifications/ endpoints.
func NotificationHandler(w http.ResponseWriter, r *http.Request) {
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
	var notification models.Notification

	// Decode request body into Notification struct
	if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Basic validation
	if notification.URL == "" || notification.Event == "" {
		http.Error(w, "Missing required fields: 'url' and 'event'", http.StatusBadRequest)
		return
	}

	// Add the webhook notification to Firestore
	docRef, _, err := firestore.Client.Collection("notifications").Add(context.Background(), notification)
	if err != nil {
		http.Error(w, "Failed to store webhook: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the generated document ID
	notification.ID = docRef.ID

	// Update the document to include the ID
	_, err = docRef.Set(context.Background(), notification)
	if err != nil {
		http.Error(w, "Failed to update webhook with ID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the ID in the response
	response := map[string]interface{}{
		"id": notification.ID,
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
	log.Printf("TriggerWebhookEvent called with event='%s', country='%s'", event, country)

	ctx := context.Background()

	// Make sure we're using uppercase for consistent matching
	event = strings.ToUpper(event)
	country = strings.ToUpper(country)

	iter := firestore.Client.Collection("notifications").Where("event", "==", event).Documents(ctx)
	defer iter.Stop()

	count := 0

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Error iterating webhooks: %v", err)
			break
		}

		var reg models.Notification
		if err = doc.DataTo(&reg); err != nil {
			log.Printf("Failed to deserialize webhook registration (id=%s): %v", doc.Ref.ID, err)
			continue
		}

		regCountry := strings.ToUpper(reg.Country)
		if regCountry == "" || regCountry == country {
			go sendWebhookNotification(reg, event, country)
			count++
		}
	}

	log.Printf("Total webhooks triggered for event='%s' and country='%s': %d", event, country, count)
}

// sendWebhookNotification sends a POST request with the correct payload to the provided URL.
func sendWebhookNotification(reg models.Notification, event, country string) {
	payload := map[string]string{
		"id":      reg.ID,
		"country": country,
		"event":   event,
		"time":    time.Now().Format("20060102 15:04"),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling webhook payload: %v", err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, reg.URL, strings.NewReader(string(data)))
	if err != nil {
		log.Printf("Error creating webhook request to %s: %v", reg.URL, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending webhook to %s: %v", reg.URL, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("Webhook URL %s responded with status %d", reg.URL, resp.StatusCode)
	} else {
		log.Printf("Webhook sent successfully to %s with status %d", reg.URL, resp.StatusCode)
	}
}
