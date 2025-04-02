package constants

const (
	Port    = ":8080"
	BaseAPI = "/dashboard/v1"

	// Local endpoints

	Registrations = BaseAPI + "/registrations/"
	Dashboards    = BaseAPI + "/dashboards/"
	Notifications = BaseAPI + "/notifications/"
	Status        = BaseAPI + "/status/"

	// External endpoints

	// The RestCountriesAPI provides country data
	RestCountriesAPI = "http://129.241.150.113:8080/v3.1"
	// The OpenMeteoAPI provides weather forecast data
	OpenMeteoAPI = "https://api.open-meteo.com/v1/forecast"
	// The CurrencyAPI provides currency exchange rates
	CurrencyAPI = "http://129.241.150.113:9090/currency/"
)
