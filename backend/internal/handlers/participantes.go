package handlers

import (
	"database/sql"
	"net/http"
)



type Participante struct {
	IDParticipante int    `json:"id_participante"`
	Nombre         string `json:"nombre"`
	Apellidos      string `json:"apellidos"`
	Matricula      string `json:"matricula"`
	IDCarrera      int    `json:"id_carrera"`
	Cuatrimestre   int    `json:"cuatrimestre"`
	Grupo          string `json:"grupo"`
	Telefono       string `json:"telefono"`
	Correo         string `json:"correo,omitempty"`
	FechaRegistro  string `json:"fecha_registro"`
}

// GetParticipantes devuelve participantes para uso administrativo autorizado.
func GetParticipantes(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.QueryContext(r.Context(), "SELECT id_participante, nombre, apellidos, matricula, id_carrera, cuatrimestre, grupo, telefono, COALESCE(correo, ''), fecha_registro FROM participantes ORDER BY id_participante")
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Status: "error", Message: "No se pudieron obtener los participantes"})
			return
		}
		defer rows.Close()

		participantes := make([]Participante, 0)
		for rows.Next() {
			var participante Participante
			var correo string
			if err := rows.Scan(&participante.IDParticipante, &participante.Nombre, &participante.Apellidos, &participante.Matricula, &participante.IDCarrera, &participante.Cuatrimestre, &participante.Grupo, &participante.Telefono, &correo, &participante.FechaRegistro); err != nil {
				writeJSON(w, http.StatusInternalServerError, ErrorResponse{Status: "error", Message: "No se pudieron leer los participantes"})
				return
			}
			participante.Correo = correo
			participantes = append(participantes, participante)
		}
		if err := rows.Err(); err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Status: "error", Message: "Ocurrió un error al leer los participantes"})
			return
		}
		writeJSON(w, http.StatusOK, participantes)
	}
}
