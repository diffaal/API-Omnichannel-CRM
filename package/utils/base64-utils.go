package utils

import (
	"encoding/base64"
	"fmt"
	"strings"
)

func DecodeBase64(encodedString string) (decodedString string, err error) {
	// decodedBytes, err := base64.StdEncoding.DecodeString(encodedString)
	decodedBytes, err := base64.URLEncoding.DecodeString(encodedString)
	if err != nil {
		fmt.Println("Error decoding:", err)
		return "", fmt.Errorf("Error when decoding base64: %+v", encodedString)
	}

	// Converting bytes to string
	decodedString = string(decodedBytes)
	return decodedString, nil
}

func EncodeBase64URL(text string) (encodedString string) {
	base64URL := base64.URLEncoding.EncodeToString([]byte(text))

	// Replace standard Base64 characters with URL-safe ones
	base64URL = strings.ReplaceAll(base64URL, "+", "-")
	base64URL = strings.ReplaceAll(base64URL, "/", "_")
	base64URL = strings.ReplaceAll(base64URL, "=", "")

	return base64URL
}
