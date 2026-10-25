package handlers

import (
	"database/sql"
	"net/http"
)

type InscripcionVista struct {
	IDInscripcion   int    `json:"id_inscripcion"`
	Juego           string `json:"juego"`
	Participante    string `json:"participante"`
	Matricula       string `json:"matricula"`
	Carrera         string `json:"carrera"`
	Cuatrimestre    int    `json:"cuatrimestre"`
	Grupo           string `json:"grupo"`
	Telefono        string `json:"telefono"`
	Correo          string `json:"correo,omitempty"`
	Asistio         int    `json:"asistio"`
	FechaInscripcion string `json:"fecha_inscripcion"`
}

func GetInscripciones(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{
				Status:  "error",
				Message: "Método no permitido",
			})
			return
		}

		rows, err := db.Query("SELECT id_inscripcion, juego, participante, matricula, carrera, cuatrimestre, grupo, telefono, IFNULL(correo, ''), asistio, fecha_inscripcion FROM vista_inscripciones ORDER BY id_inscripcion")
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{
				Status:  "error",
				Message: "Error al obtener inscripciones",
			})
			return
		}
		defer rows.Close()

		var inscripciones []InscripcionVista
		for rows.Next() {
			var iv InscripcionVista
			var correo string
			if err := rows.Scan(&iv.IDInscripcion, &iv.Juego, &iv.Participante, &iv.Matricula, &iv.Carrera, &iv.Cuatrimestre, &iv.Grupo, &iv.Telefono, &correo, &iv.Asistio, &iv.FechaInscripcion); err != nil {
				writeJSON(w, http.StatusInternalServerError, ErrorResponse{
					Status:  "error",
					Message: "Error al leer inscripciones",
				})
				return
			}
			if correo != "" {
				iv.Correo = correo
			}
			inscripciones = append(inscripciones, iv)
		}
		if err := rows.Err(); err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{
				Status:  "error",
				Message: "Error al leer inscripciones",
			})
			return
		}

		writeJSON(w, http.StatusOK, inscripciones)
	}
}
