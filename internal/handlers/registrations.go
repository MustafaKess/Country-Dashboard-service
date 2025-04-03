package handlers

import (
	"Country-Dashboard-Service/internal/firestore"
	"Country-Dashboard-Service/internal/models"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func RegistrationsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getRegistrationsHandler(w, r)
	case http.MethodPost:
		postRegistrationsHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func postRegistrationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var registration models.Registration // Assuming you have a Registration model
		err := json.NewDecoder(r.Body).Decode(&registration)
		if err != nil {
			http.Error(w, "Invalid JSON data", http.StatusBadRequest)
			return
		}

		docR, _, err := firestore.Client.Collection("registrations").Add(context.Background(), registration)
		if err != nil {
			http.Error(w, "Could not store to Firestore: "+err.Error(), http.StatusInternalServerError)
			return
		}

		id := docR.ID
		registration.ID = id
		registration.LastChange = time.Now()

		_, err = docR.Set(context.Background(), registration)
		if err != nil {
			http.Error(w, "Could not update doc with ID: "+err.Error(), http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"id":         id,
			"lastChange": registration.LastChange,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func getRegistrationsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id != "" {
		getSpecifiedRegistration(w, r)
	} else {
		getAllRegistrations(w, r)
	}
}

func getSpecifiedRegistration(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	doc, err := firestore.Client.Collection("registrations").Doc(id).Get(context.Background())
	if err != nil {
		http.Error(w, "No register found with given ID", http.StatusNotFound)
		return
	}

	var reg models.Registration
	err = doc.DataTo(&reg)
	if err != nil {
		http.Error(w, "Error with deserialization"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reg)
}

func getAllRegistrations(w http.ResponseWriter, r *http.Request) {
	iter := firestore.Client.Collection("registrations").Documents(context.Background())
	var all []models.Registration
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		var reg models.Registration
		err = doc.DataTo(&reg)
		if err != nil {
			continue // skip broken document
		}
		all = append(all, reg)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(all)
}
