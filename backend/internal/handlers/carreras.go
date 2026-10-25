package handlers

import (
	"database/sql"
	"net/http"
)

type Carrera struct {
	IDCarrera int    `json:"id_carrera"`
	Nombre    string `json:"nombre"`
}

func GetCarreras(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{
				Status:  "error",
				Message: "Método no permitido",
			})
			return
		}

		rows, err := db.Query("SELECT id_carrera, nombre FROM carreras ORDER BY id_carrera")
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{
				Status:  "error",
				Message: "Error al obtener carreras",
			})
			return
		}
		defer rows.Close()

		var carreras []Carrera
		for rows.Next() {
			var c Carrera
			if err := rows.Scan(&c.IDCarrera, &c.Nombre); err != nil {
				writeJSON(w, http.StatusInternalServerError, ErrorResponse{
					Status:  "error",
					Message: "Error al leer carreras",
				})
				return
			}
			carreras = append(carreras, c)
		}
		if err := rows.Err(); err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{
				Status:  "error",
				Message: "Error al leer carreras",
			})
			return
		}

		writeJSON(w, http.StatusOK, carreras)
	}
}
