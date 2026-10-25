package handlers

import (
	"database/sql"
	"net/http"
)

type Juego struct {
	IDJuego int    `json:"id_juego"`
	Nombre  string `json:"nombre"`
}

func GetJuegos(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{
				Status:  "error",
				Message: "Método no permitido",
			})
			return
		}

		rows, err := db.Query("SELECT id_juego, nombre FROM juegos ORDER BY id_juego")
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{
				Status:  "error",
				Message: "Error al obtener juegos",
			})
			return
		}
		defer rows.Close()

		var juegos []Juego
		for rows.Next() {
			var j Juego
			if err := rows.Scan(&j.IDJuego, &j.Nombre); err != nil {
				writeJSON(w, http.StatusInternalServerError, ErrorResponse{
					Status:  "error",
					Message: "Error al leer juegos",
				})
				return
			}
			juegos = append(juegos, j)
		}
		if err := rows.Err(); err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{
				Status:  "error",
				Message: "Error al leer juegos",
			})
			return
		}

		writeJSON(w, http.StatusOK, juegos)
	}
}
