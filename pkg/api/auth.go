package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

var password string

type tokenPayload struct {
	Hash string `json:"hash"`
	Exp  int64  `json:"exp"`
}

func passwordHash(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

func base64URL(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

func makeToken(value string) (string, error) {
	header := base64URL([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payloadBytes, err := json.Marshal(tokenPayload{
		Hash: passwordHash(value),
		Exp:  time.Now().Add(8 * time.Hour).Unix(),
	})
	if err != nil {
		return "", err
	}
	payload := base64URL(payloadBytes)
	unsigned := header + "." + payload

	mac := hmac.New(sha256.New, []byte(value))
	if _, err := mac.Write([]byte(unsigned)); err != nil {
		return "", err
	}
	signature := base64URL(mac.Sum(nil))
	return unsigned + "." + signature, nil
}

func validateToken(token, value string) bool {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false
	}

	unsigned := parts[0] + "." + parts[1]
	mac := hmac.New(sha256.New, []byte(value))
	if _, err := mac.Write([]byte(unsigned)); err != nil {
		return false
	}
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
	return payload.Hash == passwordHash(value) && payload.Exp >= time.Now().Unix()
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeErrorText(w, http.StatusMethodNotAllowed, "unsupported method")
		return
	}

	var body struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if body.Password != password {
		writeErrorText(w, http.StatusUnauthorized, "invalid password")
		return
	}
	token, err := makeToken(password)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"token": token})
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if password == "" {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil || !validateToken(cookie.Value, password) {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}
