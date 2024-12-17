package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"project/database"
	"strings"
)

type Middleware struct {
	Cacher    database.Cacher
	secretKey string
}

func NewMiddleware(cacher database.Cacher, secretKey string) Middleware {
	return Middleware{Cacher: cacher, secretKey: secretKey}
}

func validateToken(token string, secretKey string) (bool, string) {
	// Pisahkan token menjadi data dan signature
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return false, "Invalid token format"
	}

	tokenData, signature := parts[0], parts[1]

	// Decode data dari Base64
	data, err := base64.URLEncoding.DecodeString(tokenData)
	if err != nil {
		return false, "Invalid token data"
	}

	// Buat ulang signature untuk validasi
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write(data)
	expectedSignature := base64.URLEncoding.EncodeToString(h.Sum(nil))

	// Validasi signature
	if signature != expectedSignature {
		return false, "Invalid token signature"
	}

	return true, string(data)
}
