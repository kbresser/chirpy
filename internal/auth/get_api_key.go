package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	rawToken := headers.Get("Authorization")
	if rawToken == "" {
		return "", errors.New("Error getting token")
	}

	if strings.HasPrefix(rawToken, "ApiKey ") {
		token := strings.TrimSpace(strings.TrimPrefix(rawToken, "ApiKey "))
		return token, nil
	}

	return "", errors.New("Authorization header format must be 'ApiKey <token>'")
}
