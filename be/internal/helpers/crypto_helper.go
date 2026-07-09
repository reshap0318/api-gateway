package helpers

import (
	"crypto/rand"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

// GenerateRandomString generates a cryptographically secure random string.
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// HashString hashes a string using bcrypt.
func HashString(str string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// VerifyString verifies a string against its hash.
func VerifyString(str, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(str))
}
