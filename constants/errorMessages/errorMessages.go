package errorMessages

// Error messages for the API

const (
	MethodNotAllowed = "Attempted method not allowed"
	FirestoreError   = "Could not store to Firestore: "

	InvalidJSON           = "Invalid JSON data"
	InvalidRegistrationID = "Invalid registration ID"

	DeleteError = "Could not delete registration: "
	UpdateError = "Could not update registration: "

	NoIDProvided     = "No ID provided in the request"
	RegisterNotFound = "No register found with given ID "
)
