package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"

	"adeia-api/internal/util"

	"github.com/alexedwards/argon2id"
	"github.com/dchest/uniuri"
)

// GenerateRandomBytes generates cryptographically secure random bytes of the
// specified size.
func GenerateRandomBytes(n uint32) ([]byte, error) {
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

// EncodeBase64 encodes the given byte slice into url-safe base64 (with padding).
func EncodeBase64(b []byte) string {
	return base64.URLEncoding.EncodeToString(b)
}

// DecodeBase64 encodes the given byte slice into url-safe base64 (with padding).
func DecodeBase64(s string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(s)
}

// Hash hashes the give byte slice using SHA256.
func Hash(b []byte) []byte {
	h := sha256.New()
	h.Write(b)
	return h.Sum(nil)
}

// HashPassword uses argon2id to generate a hash from the password.
func HashPassword(p string) (hash string, err error) {
	return argon2id.CreateHash(p, argon2id.DefaultParams)
}

// ComparePwdHash compares the password and hash.
func ComparePwdHash(p, h string) (match bool, err error) {
	return argon2id.ComparePasswordAndHash(p, h)
}

// NewEmpID generates a user and URL-friendly alpha-numeric employee ID.
func NewEmpID() string {
	return uniuri.NewLen(util.EmployeeIDLength)
}
