package services

import (
	"Country-Dashboard-Service/constants"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// openMeteoResponse represents the structure of the weather API response.
type openMeteoResponse struct {
	Hourly struct {
		Temperature   []float64 `json:"temperature_2m"`
		Precipitation []float64 `json:"precipitation"`
	} `json:"hourly"`
}

// ErrWeatherDataUnavailable indicates missing or failed weather data.
var ErrWeatherDataUnavailable = errors.New("weather data unavailable")

// GetWeatherData returns average temperature and precipitation for given coordinates.
func GetWeatherData(lat, lon float64) (float64, float64, error) {
	url := fmt.Sprintf("%s?latitude=%.2f&longitude=%.2f&hourly=temperature_2m,precipitation", constants.OpenMeteoAPI, lat, lon)

	resp, err := http.Get(url)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, ErrWeatherDataUnavailable
	}

	var data openMeteoResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, 0, err
	}

	temps := data.Hourly.Temperature
	precips := data.Hourly.Precipitation

	if len(temps) == 0 || len(precips) == 0 {
		return 0, 0, ErrWeatherDataUnavailable
	}

	avgTemp := average(temps)
	avgPrecip := average(precips)

	return avgTemp, avgPrecip, nil
}

// average returns the mean of a slice of float64.
func average(values []float64) float64 {
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}
