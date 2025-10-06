package utils

import (
	"crypto/sha256"
	"encoding/base64"
)

func GenerateUrlKey(url string) string {
	hash := sha256.Sum256([]byte(url))                      // SHA256 → 32 bytes
	base64Key := base64.URLEncoding.EncodeToString(hash[:]) // URL-safe Base64
	if len(base64Key) > 8 {
		base64Key = base64Key[:8] // take first 8 chars for short key
	}
	return base64Key
}
