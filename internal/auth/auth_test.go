package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWTAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "test-secret"

	token, err := MakeJWT(userID, tokenSecret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT() error = %v", err)
	}

	parsedUserID, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Fatalf("ValidateJWT() error = %v", err)
	}

	if parsedUserID != userID {
		t.Fatalf("ValidateJWT() userID = %v, want %v", parsedUserID, userID)
	}
}

func TestValidateJWTRejectsExpiredToken(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "test-secret"

	token, err := MakeJWT(userID, tokenSecret, -time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT() error = %v", err)
	}

	if _, err := ValidateJWT(token, tokenSecret); err == nil {
		t.Fatal("ValidateJWT() error = nil, want token to be rejected")
	}
}

func TestValidateJWTRejectsWrongSecret(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "test-secret"
	wrongSecret := "wrong-secret"

	token, err := MakeJWT(userID, tokenSecret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT() error = %v", err)
	}

	if _, err := ValidateJWT(token, wrongSecret); err == nil {
		t.Fatal("ValidateJWT() error = nil, want signature verification to fail")
	}
}

func TestGetBearerToken(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer my-token")

	token, err := GetBearerToken(headers)
	if err != nil {
		t.Fatalf("GetBearerToken() error = %v", err)
	}

	if token != "my-token" {
		t.Fatalf("GetBearerToken() token = %q, want %q", token, "my-token")
	}
}

func TestGetBearerTokenMissingHeader(t *testing.T) {
	headers := http.Header{}

	if _, err := GetBearerToken(headers); err == nil {
		t.Fatal("GetBearerToken() error = nil, want missing token to be rejected")
	}
}
