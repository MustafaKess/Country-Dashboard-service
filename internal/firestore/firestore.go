package firestore

import (
	"context"
	"log"

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

	// Load the Firebase service account credentials from the .env file
	serviceAccount := option.WithCredentialsFile(".env/firebaseKey.json")

	// Initialize the Firebase application with the provided credentials
	app, err := firebase.NewApp(ctx, nil, serviceAccount)
	if err != nil {
		log.Fatalf("Could not initialize the Firebase application: %v", err)
	}

	// Initialize the Firestore client with the Firebase app
	var err1 error
	Client, err1 = app.Firestore(ctx)
	if err1 != nil {
		log.Fatalf("Could not initialize the Firestore client: %v", err1)
	}

	// Log successful Firestore initialization, mostly for myself for checking if all is good.
	//fmt.Println("Firestore client initialized")
}
