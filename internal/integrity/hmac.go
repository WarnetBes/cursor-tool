package integrity

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

var ErrTampered = errors.New(
	"integrity check failed: storage.json may have been tampered with")

// Sign returns HMAC-SHA256 hex string of data.
// KEY must come from OS keychain -- never hardcode.
func Sign(data, key []byte) string {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	return hex.EncodeToString(mac.Sum(nil))
}

// Verify returns ErrTampered if the signature does not match.
func Verify(data []byte, expected string, key []byte) error {
	actual := Sign(data, key)
	if !hmac.Equal([]byte(actual), []byte(expected)) {
		return ErrTampered
	}
	return nil
}
