package handlers

import (
	"database/sql"
	"net/http"
)



type Carrera struct {
	IDCarrera int    `json:"id_carrera"`
	Nombre    string `json:"nombre"`
}

// GetCarreras devuelve el catálogo de carreras disponibles para el formulario.
func GetCarreras(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.QueryContext(r.Context(), "SELECT id_carrera, nombre FROM carreras ORDER BY id_carrera")
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Status: "error", Message: "No se pudieron obtener las carreras"})
			return
		}
		defer rows.Close()

		carreras := make([]Carrera, 0)
		for rows.Next() {
			var carrera Carrera
			if err := rows.Scan(&carrera.IDCarrera, &carrera.Nombre); err != nil {
				writeJSON(w, http.StatusInternalServerError, ErrorResponse{Status: "error", Message: "No se pudieron leer las carreras"})
				return
			}
			carreras = append(carreras, carrera)
		}
		if err := rows.Err(); err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Status: "error", Message: "Ocurrió un error al leer las carreras"})
			return
		}
		writeJSON(w, http.StatusOK, carreras)
	}
}
