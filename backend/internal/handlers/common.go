package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
	"strings"
)

const maxRequestBodyBytes int64 = 1 << 20



// ErrorResponse representa la estructura JSON estándar de los errores de la API.
type ErrorResponse struct {
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Errors  []string `json:"errors,omitempty"`
}

// writeJSON devuelve una respuesta JSON con el código HTTP especificado.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// decodeJSON valida el tipo de contenido, limita el tamaño y decodifica un único objeto JSON.
func decodeJSON(w http.ResponseWriter, r *http.Request, destination any) bool {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || !strings.EqualFold(mediaType, "application/json") {
		writeJSON(w, http.StatusUnsupportedMediaType, ErrorResponse{Status: "error", Message: "El contenido debe enviarse como application/json"})
		return false
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		var maxBytesError *http.MaxBytesError
		if errors.As(err, &maxBytesError) {
			writeJSON(w, http.StatusRequestEntityTooLarge, ErrorResponse{Status: "error", Message: "El cuerpo de la solicitud excede el límite permitido"})
			return false
		}
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Status: "error", Message: "El cuerpo de la solicitud no contiene un JSON válido"})
		return false
	}

	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Status: "error", Message: "La solicitud debe contener un solo objeto JSON"})
		return false
	}
	return true
}
