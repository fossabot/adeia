package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

// GenerateRandomBytes generates cryptographically secure random bytes of the
// specified size.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	return b, nil
}

// EncodeHex encodes the given byte slice into hex.
func EncodeHex(b []byte) string {
	return hex.EncodeToString(b)
}

// EncodeBase64 encodes the given byte slice into url-safe base64.
func EncodeBase64(b []byte) string {
	return base64.URLEncoding.EncodeToString(b)
}

// Hash hashes the give byte slice using SHA256.
func Hash(b []byte) []byte {
	h := sha256.New()
	h.Write(b)
	return h.Sum(nil)
}
