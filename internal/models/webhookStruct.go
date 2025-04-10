package models

// WebhookRegistration represents a webhook registration stored in Firestore.
type WebhookRegistration struct {
	ID      string `json:"id,omitempty" firestore:"id,omitempty"` // Unique identifier for the webhook registration
	URL     string `json:"url" firestore:"url"`                   // URL to be invoked when the event occurs
	Country string `json:"country" firestore:"country"`           // Country filter; empty means all countries
	Event   string `json:"event" firestore:"event"`               // Event to trigger the webhook (e.g. REGISTER, CHANGE, DELETE, INVOKE)
}
