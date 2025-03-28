package models

// Structs for dashboard should be here
// create new files for structs needed in other places

type DashboardConfig struct {
	ID               string   `json:"id"`
	Country          string   `json:"country"`
	TargetCurrencies []string `json:"targetCurrencies"`
}

type PopulatedDashboard struct {
	Country       string            `json:"country"`
	ISOCode       string            `json:"isoCode"`
	Features      DashboardFeatures `json:"features"`
	LastRetrieval string            `json:"lastRetrieval"`
}

type DashboardFeatures struct {
	Temperature      float64            `json:"temperature"`
	Precipitation    float64            `json:"precipitation"`
	Capital          string             `json:"capital"`
	Coordinates      Coordinates        `json:"coordinates"`
	Population       int                `json:"population"`
	Area             float64            `json:"area"`
	TargetCurrencies map[string]float64 `json:"targetCurrencies"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
