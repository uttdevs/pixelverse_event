package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"time"
)



// HealthCheck confirma que el proceso HTTP está activo, sin consultar la base de datos.
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"message": "PixelVerse Event API está activo",
	})
}

// ReadyCheck confirma que el servicio y su conexión con MySQL están disponibles.
func ReadyCheck(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := db.PingContext(ctx); err != nil {
			writeJSON(w, http.StatusServiceUnavailable, ErrorResponse{Status: "error", Message: "La base de datos no está disponible"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ready", "message": "API y base de datos disponibles"})
	}
}
