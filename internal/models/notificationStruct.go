package models

type Notification struct {
	ID      string `json:"id,omitempty" firestore:"id,omitempty"`
	Event   string `json:"event" firestore:"event"`
	Country string `json:"country" firestore:"country"`
	URL     string `json:"url" firestore:"url"`
}
