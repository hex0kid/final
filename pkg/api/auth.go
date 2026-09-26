package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type tokenPayload struct {
	Hash string `json:"hash"`
	Exp  int64  `json:"exp"`
}

func passwordHash(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}

func base64URL(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

func makeToken(password string) (string, error) {
	header := base64URL([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payloadBytes, err := json.Marshal(tokenPayload{
		Hash: passwordHash(password),
		Exp:  time.Now().Add(8 * time.Hour).Unix(),
	})
	if err != nil {
		return "", err
	}
	payload := base64URL(payloadBytes)
	unsigned := header + "." + payload

	mac := hmac.New(sha256.New, []byte(password))
	_, _ = mac.Write([]byte(unsigned))
	signature := base64URL(mac.Sum(nil))
	return unsigned + "." + signature, nil
}

func validateToken(token, password string) bool {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false
	}

	unsigned := parts[0] + "." + parts[1]
	mac := hmac.New(sha256.New, []byte(password))
	_, _ = mac.Write([]byte(unsigned))
	expected := mac.Sum(nil)
	actual, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil || !hmac.Equal(actual, expected) {
		return false
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}
	var payload tokenPayload
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return false
	}
	return payload.Hash == passwordHash(password) && payload.Exp >= time.Now().Unix()
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErrorText(w, "unsupported method")
		return
	}

	var body struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, err)
		return
	}

	password := os.Getenv("TODO_PASSWORD")
	if body.Password != password {
		writeErrorText(w, "invalid password")
		return
	}
	token, err := makeToken(password)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, map[string]string{"token": token})
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		password := os.Getenv("TODO_PASSWORD")
		if password == "" {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil || !validateToken(cookie.Value, password) {
			http.Error(w, fmt.Sprint("Authentication required"), http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}
