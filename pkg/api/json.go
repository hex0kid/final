package api

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, err error) {
	writeJSON(w, map[string]string{"error": err.Error()})
}

func writeErrorText(w http.ResponseWriter, text string) {
	writeJSON(w, map[string]string{"error": text})
}
