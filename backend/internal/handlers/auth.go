package handlers

import (
	"crypto/subtle"
	"net/http"
	"strings"
)



// AdminOnly protege las rutas administrativas mediante el encabezado X-Admin-Key.
// Si no se configura una clave, el acceso se deniega de forma predeterminada.
func AdminOnly(expectedKey string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.TrimSpace(expectedKey) == "" {
			writeJSON(w, http.StatusServiceUnavailable, ErrorResponse{Status: "error", Message: "El acceso administrativo no está configurado en el servidor"})
			return
		}

		providedKey := r.Header.Get("X-Admin-Key")
		if len(providedKey) != len(expectedKey) || subtle.ConstantTimeCompare([]byte(providedKey), []byte(expectedKey)) != 1 {
			writeJSON(w, http.StatusUnauthorized, ErrorResponse{Status: "error", Message: "Se requiere una clave administrativa válida"})
			return
		}
		next.ServeHTTP(w, r)
	})
}
