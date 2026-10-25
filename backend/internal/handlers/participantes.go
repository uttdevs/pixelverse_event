package handlers

import (
	"database/sql"
	"net/http"
)

type Participante struct {
	IDParticipante int    `json:"id_participante"`
	Nombre          string `json:"nombre"`
	Apellidos       string `json:"apellidos"`
	Matricula       string `json:"matricula"`
	IDCarrera       int    `json:"id_carrera"`
	Cuatrimestre    int    `json:"cuatrimestre"`
	Grupo           string `json:"grupo"`
	Telefono        string `json:"telefono"`
	Correo          string `json:"correo,omitempty"`
	FechaRegistro   string `json:"fecha_registro"`
}

func GetParticipantes(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{
				Status:  "error",
				Message: "Método no permitido",
			})
			return
		}

		rows, err := db.Query("SELECT id_participante, nombre, apellidos, matricula, id_carrera, cuatrimestre, grupo, telefono, IFNULL(correo, ''), fecha_registro FROM participantes ORDER BY id_participante")
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{
				Status:  "error",
				Message: "Error al obtener participantes",
			})
			return
		}
		defer rows.Close()

		var participantes []Participante
		for rows.Next() {
			var p Participante
			var correo string
			if err := rows.Scan(&p.IDParticipante, &p.Nombre, &p.Apellidos, &p.Matricula, &p.IDCarrera, &p.Cuatrimestre, &p.Grupo, &p.Telefono, &correo, &p.FechaRegistro); err != nil {
				writeJSON(w, http.StatusInternalServerError, ErrorResponse{
					Status:  "error",
					Message: "Error al leer participantes",
				})
				return
			}
			if correo != "" {
				p.Correo = correo
			}
			participantes = append(participantes, p)
		}
		if err := rows.Err(); err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{
				Status:  "error",
				Message: "Error al leer participantes",
			})
			return
		}

		writeJSON(w, http.StatusOK, participantes)
	}
}
