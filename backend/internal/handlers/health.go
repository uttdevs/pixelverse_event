package handlers

import (
	"net/http"
)

// HealthCheck maneja la verificación de salud del servicio (GET /api/health).
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{
			Status:  "error",
			Message: "Método no permitido",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"message": "PixelVerse Event API is running",
	})
}
