package firestore

import (
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
	ctx    context.Context   // Context for Firebase operations
)

/*
InitFirestore initializes the Firestore client and Firebase application.
It loads the Firebase credentials from a service account key file,
then creates a Firestore client for subsequent operations.
*/
func InitFirestore() {
	// Set up context for Firestore and Firebase operations
	ctx = context.Background()

	// Use emulator if FIRESTORE_EMULATOR_HOST is set or default to it during testing
	if os.Getenv("FIRESTORE_EMULATOR_HOST") != "" || os.Getenv("GO_ENV") == "test" {
		if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
			os.Setenv("FIRESTORE_EMULATOR_HOST", "localhost:8080")
		}

		client, err := firestore.NewClient(ctx, "demo-test-project", option.WithoutAuthentication())
		if err != nil {
			log.Fatalf("Failed to create Firestore client (emulator): %v", err)
		}
		Client = client
		return
	}

	// Otherwise, use real Firebase service account
	serviceAccount := option.WithCredentialsFile(".env/firebaseKey.json")

	app, err := firebase.NewApp(ctx, nil, serviceAccount)
	if err != nil {
		log.Fatalf("Could not initialize Firebase app: %v", err)
	}

	Client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Could not initialize Firestore client: %v", err)
	}

	// Log successful Firestore initialization, mostly for myself for checking if all is good.
	//fmt.Println("Firestore client initialized")
}
