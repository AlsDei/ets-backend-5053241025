package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// writeJSON mengirim respons dalam format JSON dengan status code yang ditentukan.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// writeError mengirim respons error standar berformat {"error": "pesan"}.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}

// parseID mengekstrak nilai {key} dari parameter path dan mengonversinya ke integer positif.
func parseID(r *http.Request, key string) (int, error) {
	val := r.PathValue(key)
	id, err := strconv.Atoi(val)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("id harus berupa angka positif")
	}
	return id, nil
}
