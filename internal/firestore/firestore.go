package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"fmt"
	"google.golang.org/api/option"
	"log"
)

/*
Set up of the firebase client
*/

// Firebase context and client used by Firestore functions throughout the program.

var (
	Client *firestore.Client
	ctx    context.Context
)

/*
InitFirestore initializes the Firestore client
*/
func InitFirestore() {

	ctx = context.Background()

	//serviceAccountPath := os.Getenv("FIREBASE_CREDENTIALS_PATH")
	//if serviceAccountPath == "" {
	//	log.Fatal("FIREBASE_CREDENTIALS_PATH environment variable is not set")
	//}
	serviceAccount := option.WithCredentialsFile(".env/firebaseKey.json")
	app, err := firebase.NewApp(ctx, nil, serviceAccount)
	if err != nil {
		log.Fatalf("Could not initilize the Firebase application: %v", err)
	}

	var err1 error
	Client, err1 = app.Firestore(ctx)
	if err1 != nil {
		log.Fatalf("Could not initialize the Firestore client: %v", err1)
	}

	fmt.Println("Firestore client initialized")

}

/*
func GetDoc(collection string) ([]map[string]interface{}, error) {
	iter := client.Collection(collection).Documents(ctx)
	var docs []map[string]interface{}
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		docs = append(docs, doc.Data())
	}
	return docs, nil
}

func DisplayConfig(w http.ResponseWriter, r *http.Request) {
	configurations, err := GetDoc("configurations")
	if err != nil {
		http.Error(w, "Could not get configurations", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(configurations); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func AddDoc(collection string, data interface{}) error {
	// Ensure Firestore client and context are set up correctly
	docRef, _, err := client.Collection(collection).Add(ctx, data)
	if err != nil {
		log.Printf("Error adding document to %s: %v", collection, err)
		return err
	}

	log.Printf("Document added to collection %s. Document ID: %s", collection, docRef.ID)
	return nil
}

*/
