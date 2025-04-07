package services

// Logic for talking to external APIs
// This is for country_service
// Other files here in the services package could be e.g. weather_service, currency_service...

import (
	"Country-Dashboard-Service/constants"
	"Country-Dashboard-Service/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// restCountryResponse represents the structure of the REST Countries API response.
type restCountryResponse []struct {
	Name struct {
		Common string `json:"common"`
	} `json:"name"`
	Cca2       string    `json:"cca2"`
	Capital    []string  `json:"capital"`
	Latlng     []float64 `json:"latlng"`
	Population int       `json:"population"`
	Area       float64   `json:"area"`
	Currencies map[string]struct {
		Name string `json:"name"`
	} `json:"currencies"`
}

// ErrCountryNotFound is returned when the country is not found in the API.
var ErrCountryNotFound = errors.New("country not found")

// GetCountryInfo fetches and returns country info for the given country name.
func GetCountryInfo(countryName string) (*models.CountryInfo, error) {
	url := fmt.Sprintf("%s/name/%s", constants.RestCountriesAPI, countryName)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrCountryNotFound
	}

	var data restCountryResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, ErrCountryNotFound
	}

	c := data[0]

	// Extract currency code
	var baseCurrency string
	for code := range c.Currencies {
		baseCurrency = code
		break
	}

	// Extract capital
	capital := ""
	if len(c.Capital) > 0 {
		capital = c.Capital[0]
	}

	// Extract coordinates
	lat, lon := 0.0, 0.0
	if len(c.Latlng) >= 2 {
		lat = c.Latlng[0]
		lon = c.Latlng[1]
	}

	return &models.CountryInfo{
		Name:       c.Name.Common,
		ISOCode:    c.Cca2,
		Capital:    capital,
		Latitude:   lat,
		Longitude:  lon,
		Population: c.Population,
		Area:       c.Area,
		Currency:   baseCurrency,
	}, nil
}
