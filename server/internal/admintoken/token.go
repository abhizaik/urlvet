// Package admintoken issues and validates short-lived admin session tokens.
//
// Token format: <base64url(json payload)>.<base64url(HMAC-SHA256 signature)>
// The payload carries the subject ("admin") and an expiry unix timestamp.
package admintoken

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"os"
	"strings"
	"sync"
	"time"
)

// secret is loaded once from ADMIN_JWT_SECRET at startup.
var (
	loadOnce sync.Once
	secret   []byte
)

func loadSecret() []byte {
	loadOnce.Do(func() {
		s := os.Getenv("ADMIN_JWT_SECRET")
		if s == "" {
			panic("admintoken: ADMIN_JWT_SECRET is not set")
		}
		secret = []byte(s)
	})
	return secret
}

type payload struct {
	Subject string `json:"sub"`
	Expiry  int64  `json:"exp"` // unix timestamp
}

// Preload eagerly loads the signing secret so a missing ADMIN_JWT_SECRET
// causes a startup panic rather than a panic on the first login request.
func Preload() {
	loadSecret()
}

// Issue creates a signed token that expires after ttl.
func Issue(ttl time.Duration) (string, error) {
	p := payload{
		Subject: "admin",
		Expiry:  time.Now().Add(ttl).Unix(),
	}
	data, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	encodedPayload := base64.RawURLEncoding.EncodeToString(data)
	signature := computeSignature(encodedPayload)
	return encodedPayload + "." + signature, nil
}

// Validate returns true if the token has a valid signature and has not expired.
func Validate(token string) bool {
	encodedPayload, signature, ok := splitToken(token)
	if !ok {
		return false
	}
	if !hmac.Equal([]byte(signature), []byte(computeSignature(encodedPayload))) {
		return false
	}
	data, err := base64.RawURLEncoding.DecodeString(encodedPayload)
	if err != nil {
		return false
	}
	var p payload
	if err := json.Unmarshal(data, &p); err != nil {
		return false
	}
	return time.Now().Unix() < p.Expiry
}

func splitToken(token string) (encodedPayload, signature string, ok bool) {
	parts := strings.SplitN(token, ".", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	return parts[0], parts[1], true
}

func computeSignature(encodedPayload string) string {
	mac := hmac.New(sha256.New, loadSecret())
	mac.Write([]byte(encodedPayload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
