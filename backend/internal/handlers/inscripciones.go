package handlers

import (
	"database/sql"
	"net/http"
)



type InscripcionVista struct {
	IDInscripcion    int    `json:"id_inscripcion"`
	Juego            string `json:"juego"`
	Participante     string `json:"participante"`
	Matricula        string `json:"matricula"`
	Carrera          string `json:"carrera"`
	Cuatrimestre     int    `json:"cuatrimestre"`
	Grupo            string `json:"grupo"`
	Telefono         string `json:"telefono"`
	Correo           string `json:"correo,omitempty"`
	Asistio          int    `json:"asistio"`
	FechaInscripcion string `json:"fecha_inscripcion"`
}

// GetInscripciones devuelve las inscripciones con datos relacionados para uso administrativo autorizado.
func GetInscripciones(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := "SELECT id_inscripcion, juego, participante, matricula, carrera, cuatrimestre, grupo, telefono, COALESCE(correo, ''), asistio, fecha_inscripcion FROM vista_inscripciones ORDER BY id_inscripcion"
		rows, err := db.QueryContext(r.Context(), query)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Status: "error", Message: "No se pudieron obtener las inscripciones"})
			return
		}
		defer rows.Close()

		inscripciones := make([]InscripcionVista, 0)
		for rows.Next() {
			var inscripcion InscripcionVista
			var correo string
			if err := rows.Scan(&inscripcion.IDInscripcion, &inscripcion.Juego, &inscripcion.Participante, &inscripcion.Matricula, &inscripcion.Carrera, &inscripcion.Cuatrimestre, &inscripcion.Grupo, &inscripcion.Telefono, &correo, &inscripcion.Asistio, &inscripcion.FechaInscripcion); err != nil {
				writeJSON(w, http.StatusInternalServerError, ErrorResponse{Status: "error", Message: "No se pudieron leer las inscripciones"})
				return
			}
			inscripcion.Correo = correo
			inscripciones = append(inscripciones, inscripcion)
		}
		if err := rows.Err(); err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Status: "error", Message: "Ocurrió un error al leer las inscripciones"})
			return
		}
		writeJSON(w, http.StatusOK, inscripciones)
	}
}
