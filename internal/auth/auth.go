package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")

	split := strings.Split(auth, " ")
	if auth == "" || len(split) <= 1 {
		return "", fmt.Errorf("no token provided")
	}

	return split[1], nil
}
