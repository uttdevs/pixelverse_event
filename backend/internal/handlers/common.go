package handlers

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse representa una respuesta de error estandarizada.
type ErrorResponse struct {
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Errors  []string `json:"errors,omitempty"`
}

// writeJSON envía una respuesta JSON con el código de estado indicado.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
