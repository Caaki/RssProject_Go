package auth

import (
	"errors"
	"net/http"
	"strings"
)

// ===========Autorizacija[ Authorization: ApiKey{insert apikey here} ]
func GetAPIKey(headers http.Header) (string, error) {

	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("no authentication indo found")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 || vals[0] != "ApiKey" {
		return "", errors.New("bad format for authorization header")
	}
	return vals[1], nil
}
