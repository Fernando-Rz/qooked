package jwt_test

import (
	"qooked/internal/auth/jwt"
	"testing"
)

func TestGenerateAndValidateJWT(t *testing.T) {
	username := "testuser"
	token, err := jwt.GenerateJWT(username)
	if err != nil {
		t.Fatalf("Failed to generate JWT: %v", err)
	}

	claims, err := jwt.ValidateJWT(token)
	if err != nil {
		t.Fatalf("Failed to validate JWT: %v", err)
	}

	if claims["username"] != username {
		t.Errorf("Expected username %v, got %v", username, claims["username"])
	}
}
