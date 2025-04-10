package errorMessages

// General error messages for the API
const (
	MethodNotAllowed      = "attempted method not allowed"
	FirestoreError        = "could not store to Firestore: "
	InvalidJSON           = "invalid JSON data"
	InvalidRegistrationID = "invalid registration ID"
	DeleteError           = "could not delete registration: "
	UpdateError           = "could not update registration: "
	NoIDProvided          = "no ID provided in the request"
	RegisterNotFound      = "no register found with given ID "
	NotificationNotFound  = "notification not found"
	DeserializationError  = "error with deserialization: "
	ExtractionError       = "failed to extract registration data"
	ReadingError          = "error reading existing registration"
	NoCountryProvided     = "no country provided in the request"
	StatusEncodeError     = "failed to encode status response"
)

// ISO Code validation errors
const (
	IsoCodeDoesNotMatch   = "ISO code does not match the country name provided in the request"
	IsoRequired           = "ISO code is required for this request"
	ISOCodeMismatch       = "ISO code does not match the provided country"
	InvalidISOCodeFormat  = "invalid ISO code format in API response"
	APIFailed             = "failed to validate country with external API"
	APINotFound           = "external API returned a 404 status, country not found"
	APIUnexpectedStatus   = "external API returned unexpected status"
	NoDataFoundForCountry = "no data found for country"
)

// Webhook errors
const (
	FailedWebhookDeserialization   = "failed to deserialize webhook registration: %v"
	WebhookPayloadMarshallingError = "error marshalling webhook payload: %v"
	WebhookRequestCreationError    = "error creating webhook request: %v"
	WebhookSendError               = "error sending webhook to %s: %v"
)

// Country-related errors
const (
	CountryNotRecognized = "country is not recognized: %v"
)

// Firestore errors
const (
	FirestoreClientEmulatorError = "failed to create Firestore client (emulator): "
	FirebaseAppInitError         = "could not initialize Firebase app: "
	FirestoreClientInitError     = "could not initialize Firestore client: "
)

// Storage errors
const (
	DashboardConfigNotFound = "dashboard config not found"
)

// Notification delete message
const (
	NotificationDeleted = "notification deleted successfully"
)
