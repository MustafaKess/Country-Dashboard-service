package models

// WebhookRegistration represents a webhook registration stored in Firestore.
type WebhookRegistration struct {
	ID      string `json:"id" firestore:"id"`
	URL     string `json:"url" firestore:"url"`
	Country string `json:"country" firestore:"country"`
	Event   string `json:"event" firestore:"event"`
}
