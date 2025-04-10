package services

import (
	"Country-Dashboard-Service/constants/errorMessages"
	"Country-Dashboard-Service/internal/firestore"
	"Country-Dashboard-Service/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

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
		var entry models.WebhookRegistration
		if err = doc.DataTo(&entry); err != nil {
			log.Printf(errorMessages.FailedWebhookDeserialization, err)
			continue
		}
		// If a country is specified in the registration and it doesn't match, skip.
		if entry.Country != "" && entry.Country != country {
			continue
		}
		// Prep payload.
		payload := map[string]string{
			"id":      entry.ID,
			"country": country,
			"event":   event,
			"time":    time.Now().Format("20060102 15:04"),
		}
		// Send the webhook invocation.
		go sendWebhookNotification(entry.URL, payload)
	}
}

// sendWebhookNotification sends a POST request with the payload to the provided URL.
func sendWebhookNotification(url string, payload map[string]string) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf(errorMessages.WebhookPayloadMarshallingError, err)
		return
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		log.Printf(errorMessages.WebhookRequestCreationError, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf(errorMessages.WebhookSendError, url, err)
		return
	}
	defer resp.Body.Close()
	log.Printf("Webhook sent to %s with status code %d", url, resp.StatusCode)
}
