package handlers

import (
	"Country-Dashboard-Service/constants"
	"Country-Dashboard-Service/constants/errorMessages"
	"Country-Dashboard-Service/internal/firestore"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// StatusHandler returns the status of the service APIs in use for this project
// Those being:
// - RestCountriesAPI
// - OpenMeteoAPI
// - CurrencyAPI
//
// Also displays this API's version and webhook count
// and the uptime of the service in seconds.
// Also shows the number of registered webhooks in Firestore.

// StatusHandler handles requests for the service status.
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	restCountriesAPIStatus, err := checkAPIStatus(constants.RestCountriesAPI + "/name/norway")
	if err != nil {
		restCountriesAPIStatus = "Error"
	}
	openMeteoAPIStatus, err := checkAPIStatus(constants.OpenMeteoAPI)
	if err != nil {
		openMeteoAPIStatus = "Error"
	}

	currencyAPIStatus, err := checkAPIStatus(constants.CurrencyAPI + "NOK")
	if err != nil {
		currencyAPIStatus = "Error"
	}

	// Get the count of webhooks registered in Firestore
	webhookCount, err := getWebhookCount()
	if err != nil {
		log.Printf("Error fetching webhook count: %v", err)
		webhookCount = 0
	}

	uptime := time.Now().Unix() - serviceStartTime
	version := constants.APIVersion

	statusResponse := struct {
		RestCountriesAPI string `json:"restCountriesApi"`
		OpenMeteoAPI     string `json:"openMeteoApi"`
		CurrencyAPI      string `json:"currencyApi"`
		Version          string `json:"Version"`
		Uptime           int64  `json:"Uptime"`
		WebhookCount     int    `json:"WebhookCount"`
	}{
		RestCountriesAPI: restCountriesAPIStatus,
		OpenMeteoAPI:     openMeteoAPIStatus,
		CurrencyAPI:      currencyAPIStatus,
		Version:          version,
		Uptime:           uptime,
		WebhookCount:     webhookCount,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(statusResponse); err != nil {
		http.Error(w, errorMessages.StatusEncodeError, http.StatusInternalServerError)
	}
}

// checkAPIStatus checks the status of an API by sending a GET request.
func checkAPIStatus(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error checking status for %s: %v", url, err)
		return "", err
	}
	defer resp.Body.Close()

	return fmt.Sprintf("%d", resp.StatusCode), nil
}

// getWebhookCount retrieves the number of webhook registrations in Firestore.
func getWebhookCount() (int, error) {
	iter := firestore.Client.Collection("notifications").Documents(firestore.Ctx)
	count := 0
	for {
		_, err := iter.Next()
		if err != nil {
			// End of iteration or error occurred
			break
		}
		count++
	}
	return count, nil
}

var serviceStartTime int64

func init() {
	serviceStartTime = time.Now().Unix()
}
