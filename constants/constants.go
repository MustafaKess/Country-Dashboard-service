package constants

const (
	Port    = ":8080"
	BaseAPI = "dashboard/v1"

	//Local endpoints
	Registrations = BaseAPI + "/registrations/"
	Dashboards    = BaseAPI + "/dashboards/"
	Notifications = BaseAPI + "/notifications/"
	Status        = BaseAPI + "/status/"

	RestCountriesAPI = "http://129.241.150.113:8080/v3.1"
	OpenMeteoAPI     = "https://api.open-meteo.com/v1/forecast"
	CurrencyAPI      = "http://129.241.150.113:9090/currency/"
)
