package auth

import (
	"errors"
	"net/http"
	"strings"
)

// extracts an API key from the headers of an HTTP request
func GetAPIKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("no authentication info found")
	}

	vals := strings.Split(val, " ")
	// // check length is 2 (ApiKey: xxxxxx)
	// if len(vals) != 2 {
	// 	return "", errors.New("incorrect auth header format")
	// }
	// // check the first part of auth header
	// if vals[0] == "ApiKey" {
	// 	return "", errors.New("incorrect auth header format")
	// }
	return vals[1], nil
}
