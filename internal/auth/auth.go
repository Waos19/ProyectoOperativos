package auth

import (
	"crypto/sha256"
	"encoding/hex"
)

func VerifyLogin(username, password string, users map[string]string) bool {
	hash := sha256.Sum256([]byte(password))
	hashedInput := hex.EncodeToString(hash[:])

	storedHash, exists := users[username]
	if !exists {
		return false
	}
	return storedHash == hashedInput
}
