package handlers

import (
	"Country-Dashboard-Service/internal/firestore"
	"context"
	"os"
	"testing"
)

// cleanupFirestore deletes all documents from the test collections
func cleanupFirestore(t *testing.T) {
	collections := []string{"registrations", "notifications"}

	for _, col := range collections {
		iter := firestore.Client.Collection(col).Documents(context.Background())
		for {
			doc, err := iter.Next()
			if err != nil {
				break // No more docs or error
			}
			_, err = doc.Ref.Delete(context.Background())
			if err != nil {
				t.Logf("Failed to delete test document from %s: %v", col, err)
			}
		}
	}
}

func TestMain(m *testing.M) {
	// Set environment variable for test mode
	os.Setenv("GO_ENV", "test")

	// Initialize Firestore emulator
	firestore.InitFirestore()

	// dummy T for cleanup
	dummyT := &testing.T{}

	// Clean up any leftover documents before tests
	cleanupFirestore(dummyT)

	// Run the test suite
	code := m.Run()

	// Optionally clean up again after all tests
	cleanupFirestore(dummyT)

	os.Exit(code)
}
