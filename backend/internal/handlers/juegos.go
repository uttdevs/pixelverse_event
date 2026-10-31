package handlers

import (
	"database/sql"
	"net/http"
)

// Juego representa un título de torneo, sus datos descriptivos y participantes inscritos.
type Juego struct {
	IDJuego          int     `json:"id_juego"`
	Nombre           string  `json:"nombre"`
	Descripcion      string  `json:"descripcion"`
	CostoInscripcion float64 `json:"costo_inscripcion"`
	CuposMaximos     *int    `json:"cupos_maximos"`
	TotalInscritos   int     `json:"total_inscritos"`
}

// GetJuegos devuelve el catálogo de juegos disponibles para la inscripción y sus estadísticas.
func GetJuegos(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if db == nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Status: "error", Message: "Base de datos no inicializada"})
			return
		}

		query := `
			SELECT 
				j.id_juego, 
				j.nombre, 
				COALESCE(j.descripcion, '') AS descripcion, 
				COALESCE(j.costo_inscripcion, 0.00) AS costo_inscripcion, 
				j.cupos_maximos,
				COUNT(i.id_inscripcion) AS total_inscritos
			FROM juegos j
			LEFT JOIN inscripciones i ON i.id_juego = j.id_juego
			GROUP BY j.id_juego, j.nombre, j.descripcion, j.costo_inscripcion, j.cupos_maximos
			ORDER BY j.id_juego`

		rows, err := db.QueryContext(r.Context(), query)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Status: "error", Message: "No se pudieron obtener los juegos"})
			return
		}
		defer rows.Close()

		juegos := make([]Juego, 0)
		for rows.Next() {
			var juego Juego
			var cuposMax sql.NullInt64
			if err := rows.Scan(
				&juego.IDJuego,
				&juego.Nombre,
				&juego.Descripcion,
				&juego.CostoInscripcion,
				&cuposMax,
				&juego.TotalInscritos,
			); err != nil {
				writeJSON(w, http.StatusInternalServerError, ErrorResponse{Status: "error", Message: "No se pudieron leer los juegos"})
				return
			}
			if cuposMax.Valid {
				val := int(cuposMax.Int64)
				juego.CuposMaximos = &val
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
