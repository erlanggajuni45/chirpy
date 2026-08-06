package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")

	split := strings.Split(auth, " ")
	if auth == "" || len(split) <= 1 || split[0] != "Bearer" {
		return "", fmt.Errorf("no token provided")
	}

	return split[1], nil
}

func GetAPIKey(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")

	split := strings.Split(auth, " ")
	if auth == "" || len(split) <= 1 || split[0] != "ApiKey" {
		return "", fmt.Errorf("no token provided")
	}

	return split[1], nil
}
