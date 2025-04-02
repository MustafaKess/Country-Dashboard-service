package storage

import (
	"Country-Dashboard-Service/internal/models"
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go"
	"fmt"
	"google.golang.org/api/option"
	"io"
	"log"
	"net/http"
	"time"
)

/*
Set up of the firebase client
*/

// Firebase context and client used by Firestore functions throughout the program.

var (
	client *firestore.Client
	ctx    context.Context
)

/*
InitFirestore initializes the Firestore client
*/
func InitFirestore() {
	ctx = context.Background()
	serviceAccount := option.WithCredentialsFile("country-dashboard-prog2005.json")
	app, err := firebase.NewApp(ctx, nil, serviceAccount)
	if err != nil {
		log.Fatalf("Could not initilize the Firebase application: %v", err)
	}

	var err1 error
	client, err1 = app.Firestore(ctx)
	if err1 != nil {
		log.Fatalf("Could not initialize the Firestore client: %v", err1)
	}

	fmt.Println("Firestore client initialized")

}

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

func UpdateRegistration(collection, docID string, data map[string]interface{}) error {
	_, err := client.Collection(collection).Doc(docID).Set(ctx, data)
	if err != nil {
		return err
	}
	return nil
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

func AddDoc(w http.ResponseWriter, r *http.Request, collection string) {
	content, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Reading payload from body failed.")
		http.Error(w, "Reading payload failed.", http.StatusInternalServerError)
		return
	}
	if len(string(content)) == 0 {
		log.Println("Content appears to be empty.")
		http.Error(w, "Your payload (to be stored as document) appears to be empty. Ensure to terminate URI with /.", http.StatusBadRequest)
		return
	} else {
		c := models.CountryInfo{}
		err := json.Unmarshal(content, &c)
		if err != nil {
			log.Println("Error unmarshalling payload.")
			http.Error(w, "Error unmarshalling payload.", http.StatusInternalServerError)
			return
		}
		c.LastRetrieval = time.Now()
		id, _, err2 := client.Collection(collection).Add(ctx, c)
		if err2 != nil {
			log.Println("Error when adding document " + string(content) + ", Error: " + err2.Error())
			http.Error(w, "Error when adding document "+string(content)+", Error: "+err2.Error(), http.StatusBadRequest)
			return
		} else {
			log.Println("Document added to collection. Identifier of returned document: " + id.ID)
			http.Error(w, id.ID, http.StatusCreated)
			return
		}
	}
}

// OLD CODE BROUGHT FROM FIRESTORE DEMO

/*

func addDocument(w http.ResponseWriter, r *http.Request) {

	log.Println("Received " + r.Method + " request.")

	// very generic way of reading body; should be customized to specific use case
	content, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Reading payload from body failed.")
		http.Error(w, "Reading payload failed.", http.StatusInternalServerError)
		return
	}
	log.Println("Received request to add document for content ", string(content))
	if len(string(content)) == 0 {
		log.Println("Content appears to be empty.")
		http.Error(w, "Your payload (to be stored as document) appears to be empty. Ensure to terminate URI with /.", http.StatusBadRequest)
		return
	} else {
		// Add element in embedded structure.
		// Note: this structure is defined by the client, not the server!; it exemplifies the use of a complex structure
		// and illustrates how you can use Firestore features such as Firestore timestamps.
		c := models.CountryInfo{}
		err := json.Unmarshal(content, &c)
		if err != nil {
			log.Println("Error unmarshalling payload.")
			http.Error(w, "Error unmarshalling payload.", http.StatusInternalServerError)
			return
		}
		// Update timestamp
		c.LastRetrieval = time.Now()

		id, _, err2 := client.Collection(collection).Add(ctx, c)

		 Alternatively, you can directly encode data structures:
		Example:
			id, _, err2 := client.Collection(collection).Add(ctx, map[string]interface{}{
					"content": string(content),           // this is self-defined and embeds the content passed by the client
					"ct":      ct,                        // a self-defined counter (as an example for an additional field if useful)
					"time":    firestore.ServerTimestamp, // exemplifying Firestore features
				})


		ct++
		if err2 != nil {
			// Error handling
			log.Println("Error when adding document " + string(content) + ", Error: " + err2.Error())
			http.Error(w, "Error when adding document "+string(content)+", Error: "+err2.Error(), http.StatusBadRequest)
			return
		} else {
			// Returns document ID in body
			log.Println("Document added to collection. Identifier of returned document: " + id.ID)
			http.Error(w, id.ID, http.StatusCreated)
			return
		}
	}
}


func displayDocument(w http.ResponseWriter, r *http.Request) {

	log.Println("Received " + r.Method + " request.")

	// Test for embedded message ID
	messageId := r.PathValue("id")

	if messageId != "" {
		// Extract individual message

		// Retrieve specific message based on id (Firestore-generated hash)
		res := client.Collection(collection).Doc(messageId)

		// Retrieve reference to document
		doc, err2 := res.Get(ctx)
		if err2 != nil {
			log.Println("Error extracting body of returned document of message " + messageId)
			http.Error(w, "Error extracting body of returned document of message "+messageId, http.StatusInternalServerError)
			return
		}

		// A message map with string keys. Each key is one field
		rawContent := doc.Data()
		s, err := json.Marshal(rawContent)
		if err != nil {
			log.Println("Error marshalling payload.")
			http.Error(w, "Error marshalling payload.", http.StatusInternalServerError)
			return
		}
		_, err3 := fmt.Fprintln(w, string(s)) // here we retrieve the field containing the originally stored payload
		if err3 != nil {
			log.Println("Error while writing response body of message " + messageId)
			http.Error(w, "Error while writing response body of message "+messageId, http.StatusInternalServerError)
			return
		}
	} else {
		// Collective retrieval of messages
		iter := client.Collection(collection).Documents(ctx) // Loop through all entries in collection "messages"
		// Consider refining with ordering (e.g., '.OrderBy("created", firestore.Asc)') and introducing limits (e.g., '.Limit(3)').

		for {
			doc, err := iter.Next()
			if errors.Is(err, iterator.Done) {
				break
			}
			if err != nil {
				log.Printf("Failed to iterate: %v", err)
				return
			}
			// Note: You can access the document ID using "doc.Ref.ID"

			// Returns a map with string keys.
			rawContent := doc.Data()
			s, err := json.Marshal(rawContent)
			_, err = fmt.Fprintln(w, string(s))
			if err != nil {
				log.Println("Error while writing response body (Error: " + err.Error() + ")")
				http.Error(w, "Error while writing response body (Error: "+err.Error()+")", http.StatusInternalServerError)
				return
			}
		}
	}
}


func handleMessage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addDocument(w, r)
	case http.MethodGet:
		displayDocument(w, r)
	default:
		log.Println("Unsupported request method " + r.Method)
		http.Error(w, "Unsupported request method "+r.Method, http.StatusMethodNotAllowed)
		return
	}
}

func main() {
	// Firebase initialisation
	ctx = context.Background()

	// We use a service account, load credentials file that you downloaded from your project's settings menu.
	// It should reside in your project directory.
	// Make sure this file is git-ignored, since it is the access token to the database.
	sa := option.WithCredentialsFile("./demo-service-account.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Println(err)
		return
	}

	// Instantiate client
	client, err = app.Firestore(ctx)

	// Alternative setup, directly through Firestore (without initial reference to Firebase); but requires Project ID; useful if multiple projects are used
	// client, err := firestore.NewClient(ctx, projectID)

	// Check whether there is an error when connecting to Firestore
	if err != nil {
		log.Println(err)
		return
	}

	// Close down client at the end of the function
	defer func() {
		errClose := client.Close()
		if errClose != nil {
			log.Fatal("Closing of the Firebase client failed. Error:", errClose)
		}
	}()

	// Make it Heroku-compatible
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port

	http.HandleFunc("/messages/{id}", handleMessage)
	http.HandleFunc("/messages/", handleMessage) // For POST without ID; the first one did not catch on that one
	log.Printf("Firestore REST service listening on %s ...\n", addr)
	if errSrv := http.ListenAndServe(addr, nil); errSrv != nil {
		panic(errSrv)
	}
}


  Advanced Tasks:
   - Introduce update functionality via PUT and/or PATCH
   - Introduce delete functionality
   - Adapt addDocument and displayDocument function to support custom JSON schema

*/
