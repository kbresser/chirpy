package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
	// Set up test data
	userID := uuid.New()      // Create a random UUID
	secret := "my-secret-key" // This is the secret - just a string
	expiresIn := time.Hour    // Token expires in 1 hour

	// Create a JWT
	tokenString, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	// Validate the same JWT
	validatedUserID, err := ValidateJWT(tokenString, secret)
	if err != nil {
		t.Fatalf("Failed to validate JWT: %v", err)
	}

	// Check if the user ID matches
	if validatedUserID != userID {
		t.Errorf("Expected user ID %v, got %v", userID, validatedUserID)
	}
}
