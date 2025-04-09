package handlers

import (
	"Country-Dashboard-Service/constants"
	"Country-Dashboard-Service/internal/firestore"
	"Country-Dashboard-Service/internal/models"
	"Country-Dashboard-Service/internal/services"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Utility function to insert a test webhook into Firestore.
func insertTestWebhook(t *testing.T) string {
	t.Helper()

	webhook := models.WebhookRegistration{
		URL:     "https://localhost:9999/test",
		Country: "NO",
		Event:   constants.EventInvoke,
	}

	docRef, _, err := firestore.Client.Collection("notifications").Add(context.Background(), webhook)
	if err != nil {
		t.Fatalf("Failed to insert webhook: %v", err)
	}

	return docRef.ID
}

// Deletes all documents in the "notifications" collection.
func clearNotificationsCollection(t *testing.T) {
	t.Helper()

	iter := firestore.Client.Collection("notifications").Documents(context.Background())
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		_, err = doc.Ref.Delete(context.Background())
		if err != nil {
			t.Fatalf("Failed to delete document: %v", err)
		}
	}
}

func TestPostNotification(t *testing.T) {
	webhook := models.WebhookRegistration{
		URL:     "https://localhost:9999/webhook",
		Country: "NO",
		Event:   constants.EventRegister,
	}
	payload, _ := json.Marshal(webhook)

	req := httptest.NewRequest(http.MethodPost, constants.Notifications, bytes.NewReader(payload))
	w := httptest.NewRecorder()

	NotificationsHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201 Created, got %d", w.Code)
	}
}

func TestPostNotification_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, constants.Notifications, bytes.NewReader([]byte("not-json")))
	w := httptest.NewRecorder()

	NotificationsHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 Bad Request, got %d", w.Code)
	}
}

func TestGetSpecificNotification(t *testing.T) {
	id := insertTestWebhook(t)

	req := httptest.NewRequest(http.MethodGet, constants.Notifications+id, nil)
	w := httptest.NewRecorder()

	NotificationsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", w.Code)
	}
}

func TestGetSpecificNotification_NonExistingID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, constants.Notifications+"non-existing-id", nil)
	w := httptest.NewRecorder()

	NotificationsHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 Not Found, got %d", w.Code)
	}
}

func TestGetAllNotifications(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, constants.Notifications, nil)
	w := httptest.NewRecorder()

	NotificationsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", w.Code)
	}
}

func TestGetAllNotifications_EmptyCollection(t *testing.T) {
	clearNotificationsCollection(t)

	req := httptest.NewRequest(http.MethodGet, constants.Notifications, nil)
	w := httptest.NewRecorder()

	NotificationsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", w.Code)
	}
}

func TestDeleteNotification(t *testing.T) {
	id := insertTestWebhook(t)

	req := httptest.NewRequest(http.MethodDelete, constants.Notifications+id, nil)
	w := httptest.NewRecorder()

	NotificationsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", w.Code)
	}
}

func TestDeleteNotification_NonExistingID(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, constants.Notifications+"non-existing-id", nil)
	w := httptest.NewRecorder()

	NotificationsHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 Not Found, got %d", w.Code)
	}
}

func TestNotificationsMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, constants.Notifications, nil)
	w := httptest.NewRecorder()

	NotificationsHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405 Method Not Allowed, got %d", w.Code)
	}
}

func TestTriggerWebhookEvent_SendsNotification(t *testing.T) {
	received := make(chan bool, 1)

	// Start a mock webhook server that listens for the trigger
	mockWebhook := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Errorf("Failed to decode webhook payload: %v", err)
		}

		if payload["country"] != "NO" || payload["event"] != constants.EventRegister {
			t.Errorf("Unexpected webhook payload: %v", payload)
		}

		received <- true
	}))
	defer mockWebhook.Close()

	// Register the webhook to point to the mock server
	webhook := models.WebhookRegistration{
		URL:     mockWebhook.URL,
		Country: "NO",
		Event:   constants.EventRegister,
	}
	_, _, err := firestore.Client.Collection("notifications").Add(context.Background(), webhook)
	if err != nil {
		t.Fatalf("Failed to register webhook: %v", err)
	}

	// Trigger it
	services.TriggerWebhookEvent(constants.EventRegister, "NO")

	// Check that the mock server received the webhook
	select {
	case <-received:
		// success
	case <-time.After(2 * time.Second):
		t.Fatal("Timed out waiting for webhook event")
	}
}
