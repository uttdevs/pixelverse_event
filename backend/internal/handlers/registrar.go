package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-sql-driver/mysql"
)

type RegistrarParticipanteRequest struct {
	Nombre       string `json:"nombre"`
	Apellidos    string `json:"apellidos"`
	Matricula    string `json:"matricula"`
	IDCarrera    int    `json:"id_carrera"`
	Cuatrimestre int    `json:"cuatrimestre"`
	Grupo        string `json:"grupo"`
	Telefono     string `json:"telefono"`
	Correo       string `json:"correo,omitempty"`
	JuegoID      int    `json:"juego_id"`
}

type RegistrarParticipanteResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		IDParticipante int    `json:"id_participante"`
		Matricula      string `json:"matricula"`
		JuegoID        int    `json:"juego_id"`
		IDInscripcion  int64  `json:"id_inscripcion"`
	} `json:"data,omitempty"`
}

func RegistrarParticipante(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{
				Status:  "error",
				Message: "Método no permitido. Utilice POST.",
			})
			return
		}

		var req RegistrarParticipanteRequest
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		if err := decoder.Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{
				Status:  "error",
				Message: "Cuerpo de solicitud JSON inválido",
				Errors:  []string{err.Error()},
			})
			return
		}

		req.Nombre = strings.TrimSpace(req.Nombre)
		req.Apellidos = strings.TrimSpace(req.Apellidos)
		req.Matricula = strings.TrimSpace(req.Matricula)
		req.Grupo = strings.TrimSpace(req.Grupo)
		req.Telefono = strings.TrimSpace(req.Telefono)
		req.Correo = strings.TrimSpace(req.Correo)

		var validationErrors []string
		if req.Nombre == "" {
			validationErrors = append(validationErrors, "El campo 'nombre' es obligatorio.")
		}
		if req.Apellidos == "" {
			validationErrors = append(validationErrors, "El campo 'apellidos' es obligatorio.")
		}
		if req.Matricula == "" {
			validationErrors = append(validationErrors, "El campo 'matricula' es obligatorio.")
		}
		if req.IDCarrera == 0 {
			validationErrors = append(validationErrors, "El campo 'id_carrera' es obligatorio.")
		}
		if req.Cuatrimestre == 0 {
			validationErrors = append(validationErrors, "El campo 'cuatrimestre' es obligatorio.")
		}
		if req.Grupo == "" {
			validationErrors = append(validationErrors, "El campo 'grupo' es obligatorio.")
		}
		if req.Telefono == "" {
			validationErrors = append(validationErrors, "El campo 'telefono' es obligatorio.")
		}
		if req.JuegoID == 0 {
			validationErrors = append(validationErrors, "El campo 'juego_id' es obligatorio.")
		}

		if len(validationErrors) > 0 {
			writeJSON(w, http.StatusUnprocessableEntity, ErrorResponse{
				Status:  "error",
				Message: "Error de validación",
				Errors:  validationErrors,
			})
			return
		}

		tx, err := db.Begin()
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{
				Status:  "error",
				Message: "Error al iniciar transacción",
			})
			return
		}
		defer tx.Rollback()

		var idParticipante int
		var correo sql.NullString
		if req.Correo == "" {
			correo = sql.NullString{Valid: false}
		} else {
			correo = sql.NullString{String: req.Correo, Valid: true}
		}

		res, err := tx.Exec("INSERT INTO participantes (nombre, apellidos, matricula, id_carrera, cuatrimestre, grupo, telefono, correo) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
			req.Nombre, req.Apellidos, req.Matricula, req.IDCarrera, req.Cuatrimestre, req.Grupo, req.Telefono, correo)
		if err != nil {
			if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
				writeJSON(w, http.StatusConflict, ErrorResponse{
					Status:  "error",
					Message: "La matrícula ya se encuentra registrada",
				})
				return
			}
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{
				Status:  "error",
				Message: "Error al registrar participante",
			})
			return
		}

		idParticipante64, err := res.LastInsertId()
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{
				Status:  "error",
				Message: "Error al obtener ID de participante",
			})
			return
		}
		idParticipante = int(idParticipante64)

		insRes, err := tx.Exec("INSERT INTO inscripciones (id_participante, id_juego) VALUES (?, ?)", idParticipante, req.JuegoID)
		if err != nil {
			if mysqlErr, ok := err.(*mysql.MySQLError); ok {
				if mysqlErr.Number == 1062 {
					writeJSON(w, http.StatusConflict, ErrorResponse{
						Status:  "error",
						Message: "El participante ya se encuentra inscrito en este juego",
					})
					return
				}
				if mysqlErr.Number == 1452 {
					writeJSON(w, http.StatusUnprocessableEntity, ErrorResponse{
						Status:  "error",
						Message: "Carrera o juego no válido",
					})
					return
				}
			}
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{
				Status:  "error",
				Message: "Error al registrar inscripción",
			})
			return
		}

		idInscripcion, err := insRes.LastInsertId()
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{
				Status:  "error",
				Message: "Error al obtener ID de inscripción",
			})
			return
		}

		if err := tx.Commit(); err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{
				Status:  "error",
				Message: "Error al confirmar registro",
			})
			return
		}

		resp := RegistrarParticipanteResponse{
			Status:  "success",
			Message: "Registro completado con éxito",
		}
		resp.Data.IDParticipante = idParticipante
		resp.Data.Matricula = req.Matricula
		resp.Data.JuegoID = req.JuegoID
		resp.Data.IDInscripcion = idInscripcion

		writeJSON(w, http.StatusCreated, resp)
	}
}
