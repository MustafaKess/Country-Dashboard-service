package constants

const (
	Port       = ":8080"
	APIVersion = "v1"
	BaseAPI    = "/dashboard" + APIVersion

	// Local endpoints
	Registrations = BaseAPI + "/registrations/"
	Dashboards    = BaseAPI + "/dashboards/"
	Notifications = BaseAPI + "/notifications/"
	Status        = BaseAPI + "/status/"

	// webhook event constants
	EventRegister = "REGISTER"
	EventChange   = "CHANGE"
	EventDelete   = "DELETE"
	EventInvoke   = "INVOKE"

	// Firestore project and file config
	FirebaseProjectID    = "demo-test-project"
	ServiceAccountJSON   = ".env/firebaseKey.json"
	DefaultEmulatorHost  = "localhost:8080"
	EnvFirestoreEmulator = "FIRESTORE_EMULATOR_HOST"
	EnvGoEnv             = "GO_ENV"
	EnvGoEnvTestValue    = "test"
)

var (
	// External endpoints
	RestCountriesAPI = "http://129.241.150.113:8080/v3.1"
	OpenMeteoAPI     = "https://api.open-meteo.com/v1/forecast"
	CurrencyAPI      = "http://129.241.150.113:9090/currency/"
)
