package services

import (
	"Country-Dashboard-Service/constants"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// currencyResponse matches the expected JSON structure from the Currency API.
type currencyResponse struct {
	Base  string             `json:"base"`
	Rates map[string]float64 `json:"rates"`
}

// ErrCurrencyDataUnavailable is returned when exchange rate data cannot be fetched.
var ErrCurrencyDataUnavailable = errors.New("currency data unavailable")

// GetExchangeRates fetches exchange rates for the given base currency and filters only the desired target currencies.
func GetExchangeRates(base string, targets []string) (map[string]float64, error) {
	url := fmt.Sprintf("%s%s", constants.CurrencyAPI, base)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrCurrencyDataUnavailable
	}

	var data currencyResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	result := make(map[string]float64)
	for _, target := range targets {
		if rate, ok := data.Rates[target]; ok {
			result[target] = rate
		}
	}

	return result, nil
}
