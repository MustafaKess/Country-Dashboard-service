package firestore

import (
	"Country-Dashboard-Service/constants"
	"Country-Dashboard-Service/constants/errorMessages"
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// Global variables for Firestore client and context
var (
	Client *firestore.Client // Firestore client to interact with Firestore DB
	Ctx    context.Context   // Context for Firebase operations
)

/*
InitFirestore initializes the Firestore client and Firebase application.
It loads the Firebase credentials from a service account key file,
then creates a Firestore client for subsequent operations.
*/
func InitFirestore() {
	// Set up context for Firestore and Firebase operations
	Ctx = context.Background()

	// Use emulator if FIRESTORE_EMULATOR_HOST is set or default to it during testing
	if os.Getenv(constants.EnvFirestoreEmulator) != "" || os.Getenv(constants.EnvGoEnv) == constants.EnvGoEnvTestValue {
		if os.Getenv(constants.EnvFirestoreEmulator) == "" {
			os.Setenv(constants.EnvFirestoreEmulator, constants.DefaultEmulatorHost)
		}

		client, err := firestore.NewClient(Ctx, constants.FirebaseProjectID, option.WithoutAuthentication())
		if err != nil {
			log.Fatalf(errorMessages.FirestoreClientEmulatorError + err.Error())
		}
		Client = client
		return
	}

	// Otherwise, use real Firebase service account
	serviceAccount := option.WithCredentialsFile(constants.ServiceAccountJSON)

	app, err := firebase.NewApp(Ctx, nil, serviceAccount)
	if err != nil {
		log.Fatalf(errorMessages.FirebaseAppInitError + err.Error())
	}

	Client, err = app.Firestore(Ctx)
	if err != nil {
		log.Fatalf(errorMessages.FirestoreClientInitError + err.Error())
	}

	// Log successful Firestore initialization, mostly for myself for checking if all is good.
	//fmt.Println("Firestore client initialized")
}
