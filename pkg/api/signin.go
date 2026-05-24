package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

func signinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}

	var req struct {
		Password string `json:"password"`
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	if err = json.Unmarshal(body, &req); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	pass := os.Getenv("TODO_PASSWORD")
	if pass == "" || req.Password != pass {
		writeJSON(w, map[string]string{"error": errInvalidPassword.Error()})
		return
	}

	token, err := CreateToken()
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, map[string]string{"token": token})
}
