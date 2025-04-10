package errorMessages

// General error messages for the API
const (
	MethodNotAllowed      = "Attempted method not allowed"
	FirestoreError        = "Could not store to Firestore: "
	InvalidJSON           = "Invalid JSON data"
	InvalidRegistrationID = "Invalid registration ID"
	DeleteError           = "Could not delete registration: "
	UpdateError           = "Could not update registration: "
	NoIDProvided          = "No ID provided in the request"
	RegisterNotFound      = "No register found with given ID "
	NotificationNotFound  = "Notification not found"
	DeserializationError  = "Error with deserialization: "
	ExtractionError       = "Failed to extract registration data"
	ReadingError          = "Error reading existing registration"
	NoCountryProvided     = "No country provided in the request"
)

// ISO Code validation errors
const (
	IsoCodeDoesNotMatch   = "ISO code does not match the country name provided in the request"
	IsoRequired           = "ISO code is required for this request"
	ISOCodeMismatch       = "ISO code does not match the provided country"
	InvalidISOCodeFormat  = "Invalid ISO code format in API response"
	APIFailed             = "Failed to validate country with external API"
	APINotFound           = "External API returned a 404 status, country not found"
	APIUnexpectedStatus   = "External API returned unexpected status"
	NoDataFoundForCountry = "No data found for country"
)

// Country-related errors
const (
	CountryNotRecognized = "Country is not recognized"
)
