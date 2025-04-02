package handlers

import (
	"Country-Dashboard-Service/constants"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

//StatusHandler returns the status of the service APIs in use for this project
//Those being:
// - RestCountriesAPI
// - OpenMeteoAPI
// - CurrencyAPI
//
//Also displays this API's version

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

	uptime := time.Now().Unix() - serviceStartTime
	version := "v1"

	statusResponse := struct {
		RestCountriesAPI string `json:"restCountriesApi"`
		OpenMeteoAPI     string `json:"openMeteoApi"`
		CurrencyAPI      string `json:"currencyApi"`
		Version          string `json:"Version"`
		Uptime           int64  `json:"Uptime"`
	}{
		RestCountriesAPI: restCountriesAPIStatus,
		OpenMeteoAPI:     openMeteoAPIStatus,
		CurrencyAPI:      currencyAPIStatus,
		Version:          version,
		Uptime:           uptime,
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(statusResponse); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func checkAPIStatus(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error checking status for %s: %v", url, err)
		return "", err
	}
	defer resp.Body.Close()

	return fmt.Sprintf("%d", resp.StatusCode), nil
}

var serviceStartTime int64

func init() {
	serviceStartTime = time.Now().Unix()
}
