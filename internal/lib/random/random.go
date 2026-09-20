package random

import (
	"crypto/rand"
	"encoding/base64"
)

func NewRandomString(aliasLength int) (string, error) {
	byteLength := (aliasLength * 3) / 4
	if byteLength == 0 {
		byteLength = 1
	}

	b := make([]byte, byteLength+2) 
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	encoded := base64.RawURLEncoding.EncodeToString(b)

	if len(encoded) < aliasLength {
		return NewRandomString(aliasLength) 
	}

	return encoded[:aliasLength], nil

}