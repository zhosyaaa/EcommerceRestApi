package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/argon2"
	"strings"
)

func HashPassword(pass string) (string, error) {
	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	hashedPassword := argon2.IDKey([]byte(pass), []byte(salt), 1, 64*1024, 4, 32)
	encodedPassword := fmt.Sprintf("%s.%s", base64.RawStdEncoding.EncodeToString(salt), base64.RawStdEncoding.EncodeToString(hashedPassword))
	return encodedPassword, nil
}

func VerifyPassword(encodedPassword string, password string) bool {
	encodedSaltAndPassword := password
	parts := strings.Split(encodedSaltAndPassword, ".")
	decodedHashedPassword, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}
	decodedSalt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return false
	}
	hashedPassword := argon2.IDKey([]byte(encodedPassword), decodedSalt, 1, 64*1024, 4, 32)
	if bytes.Equal(hashedPassword, decodedHashedPassword) {
		return true
	}
	return false
}
