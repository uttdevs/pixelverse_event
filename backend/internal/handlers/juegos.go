package handlers

import (
	"database/sql"
	"net/http"
)



type Juego struct {
	IDJuego int    `json:"id_juego"`
	Nombre  string `json:"nombre"`
}

// GetJuegos devuelve el catálogo de juegos disponibles para la inscripción.
func GetJuegos(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.QueryContext(r.Context(), "SELECT id_juego, nombre FROM juegos ORDER BY id_juego")
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Status: "error", Message: "No se pudieron obtener los juegos"})
			return
		}
		defer rows.Close()

		juegos := make([]Juego, 0)
		for rows.Next() {
			var juego Juego
			if err := rows.Scan(&juego.IDJuego, &juego.Nombre); err != nil {
				writeJSON(w, http.StatusInternalServerError, ErrorResponse{Status: "error", Message: "No se pudieron leer los juegos"})
				return
			}
			juegos = append(juegos, juego)
		}
		if err := rows.Err(); err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Status: "error", Message: "Ocurrió un error al leer los juegos"})
			return
		}
		writeJSON(w, http.StatusOK, juegos)
	}
}
