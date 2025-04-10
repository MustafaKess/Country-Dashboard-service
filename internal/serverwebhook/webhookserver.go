package serverWebhook

import (
	"Country-Dashboard-Service/constants"
	"Country-Dashboard-Service/internal/models"
	"Country-Dashboard-Service/internal/utils"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

/*
WebhookServer is a dedicated server for receiving webhook invocations.
It stores received webhooks in memory and clears them every 24 hours.
This additional component aids in debugging and testing webhook notifications.
*/
type WebhookServer struct {
	HTTP  *http.Server
	Wh    []models.WebhookRegistration // In-memory store for webhook invocations
	mutex sync.Mutex                   // Mutex to prevent race conditions when accessing Wh
}

// Start initializes and starts the webhook server on the port defined by the WEBHOOK_PORT environment variable.
func Start() {
	port := os.Getenv("WEBHOOK_PORT")
	if port == "" {
		port = "8081" // default to port 8081 if not set
		log.Printf("WEBHOOK_PORT not set, defaulting to %s", port)
	}
	mux := http.NewServeMux()
	server := &WebhookServer{}
	mux.HandleFunc(constants.Notifications, server.HandleWebhooks)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
	server.HTTP = srv

	// Start routine to clear stored webhooks every 24 hours.
	go server.clearWebhookStorage()

	log.Println("Webhook server starting on port", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// HandleWebhooks handles GET and POST requests for the in-memory webhook storage.
func (s *WebhookServer) HandleWebhooks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.getWebhooks(w, r)
	case http.MethodPost:
		s.postWebhook(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *WebhookServer) postWebhook(w http.ResponseWriter, r *http.Request) {
	var wh models.WebhookRegistration
	if _, err := utils.DecodeRequest[models.WebhookRegistration](r); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	// Append the webhook registration to the in-memory slice.
	s.mutex.Lock()
	s.Wh = append(s.Wh, wh)
	s.mutex.Unlock()
	log.Println("Received webhook via dedicated server")
	w.WriteHeader(http.StatusOK)
}

func (s *WebhookServer) getWebhooks(w http.ResponseWriter, r *http.Request) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	utils.Encode(w, http.StatusOK, s.Wh)
}

func (s *WebhookServer) clearWebhookStorage() {
	for {
		time.Sleep(24 * time.Hour)
		s.mutex.Lock()
		s.Wh = nil
		s.mutex.Unlock()
		log.Println("Cleared stored webhooks after 24h")
	}
}
