package cryptoUtil

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"golang.org/x/crypto/bcrypt"

)

// HashPassword hashes a plain password using bcrypt.
func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err
}

// ComparePassword compares a hashed password with a plain password.
func ComparePassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

func HashString(args ...string) string {
	h := sha256.Sum256([]byte(strings.Join(args, "")))
	return hex.EncodeToString(h[:])
}
