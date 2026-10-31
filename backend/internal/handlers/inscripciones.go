package handlers

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"unicode/utf8"
)

const maxBusquedaRunes = 100

type InscripcionVista struct {
	IDInscripcion    int    `json:"id_inscripcion"`
	IDParticipante   int    `json:"id_participante"`
	IDJuego          int    `json:"id_juego"`
	IDCarrera        int    `json:"id_carrera"`
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

// UpdateInscripcionRequest contiene los campos editables de una inscripción.
// Los punteros permiten distinguir entre "no enviado" y un valor explícito.
type UpdateInscripcionRequest struct {
	JuegoID *int `json:"juego_id,omitempty"`
	Asistio *int `json:"asistio,omitempty"`
}

// baseInscripcionesQuery une las tablas directamente para exponer también los identificadores
// necesarios para la edición, sin depender de una versión concreta de vista_inscripciones.
const baseInscripcionesQuery = `SELECT i.id_inscripcion, p.id_participante, j.id_juego, c.id_carrera,
	j.nombre, p.nombre_completo, p.matricula, c.nombre, p.cuatrimestre, p.grupo, p.telefono,
	COALESCE(p.correo, ''), i.asistio, i.fecha_inscripcion
FROM inscripciones i
JOIN participantes p ON p.id_participante = i.id_participante
JOIN juegos j ON j.id_juego = i.id_juego
JOIN carreras c ON c.id_carrera = p.id_carrera`

// likeEscaper escapa los comodines de LIKE para que la búsqueda sea literal.
var likeEscaper = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)

// buildInscripcionesFilter traduce los parámetros de consulta a una cláusula WHERE parametrizada.
// Parámetros admitidos: q (texto), juego_id, id_carrera y asistio (0|1).
func buildInscripcionesFilter(values url.Values) (string, []any, []string) {
	conditions := make([]string, 0, 4)
	args := make([]any, 0, 8)
	validationErrors := make([]string, 0)

	if search := strings.TrimSpace(values.Get("q")); search != "" {
		if utf8.RuneCountInString(search) > maxBusquedaRunes {
			validationErrors = append(validationErrors, "La búsqueda no puede superar 100 caracteres.")
		} else {
			pattern := "%" + likeEscaper.Replace(search) + "%"
			conditions = append(conditions, "(p.nombre_completo LIKE ? OR p.matricula LIKE ? OR p.correo LIKE ? OR p.telefono LIKE ? OR p.grupo LIKE ?)")
			args = append(args, pattern, pattern, pattern, pattern, pattern)
		}
	}

	positiveInt := func(name, column, message string) {
		raw := strings.TrimSpace(values.Get(name))
		if raw == "" {
			return
		}
		value, err := strconv.Atoi(raw)
		if err != nil || value < 1 {
			validationErrors = append(validationErrors, message)
			return
		}
		conditions = append(conditions, column+" = ?")
		args = append(args, value)
	}
	positiveInt("juego_id", "i.id_juego", "El filtro 'juego_id' debe ser un número entero positivo.")
	positiveInt("id_carrera", "p.id_carrera", "El filtro 'id_carrera' debe ser un número entero positivo.")

	if raw := strings.TrimSpace(values.Get("asistio")); raw != "" {
		if raw != "0" && raw != "1" {
			validationErrors = append(validationErrors, "El filtro 'asistio' debe ser 0 o 1.")
		} else {
			value, _ := strconv.Atoi(raw)
			conditions = append(conditions, "i.asistio = ?")
			args = append(args, value)
		}
	}

	if len(conditions) == 0 {
		return "", args, validationErrors
	}
	return " WHERE " + strings.Join(conditions, " AND "), args, validationErrors
}

// queryInscripciones ejecuta la consulta base con el filtro indicado.
func queryInscripciones(ctx context.Context, db *sql.DB, where string, args []any) ([]InscripcionVista, error) {
	rows, err := db.QueryContext(ctx, baseInscripcionesQuery+where+" ORDER BY i.id_inscripcion", args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	inscripciones := make([]InscripcionVista, 0)
	for rows.Next() {
		var inscripcion InscripcionVista
		var fecha sql.NullTime
		if err := rows.Scan(&inscripcion.IDInscripcion, &inscripcion.IDParticipante, &inscripcion.IDJuego, &inscripcion.IDCarrera,
			&inscripcion.Juego, &inscripcion.Participante, &inscripcion.Matricula, &inscripcion.Carrera, &inscripcion.Cuatrimestre,
			&inscripcion.Grupo, &inscripcion.Telefono, &inscripcion.Correo, &inscripcion.Asistio, &fecha); err != nil {
			return nil, err
		}
		if fecha.Valid {
			inscripcion.FechaInscripcion = fecha.Time.Format("2006-01-02 15:04:05")
		}
		inscripciones = append(inscripciones, inscripcion)
	}
	return inscripciones, rows.Err()
}

// parseInscripcionID obtiene y valida el identificador {id} de la ruta.
func parseInscripcionID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Status: "error", Message: "El identificador de la inscripción no es válido"})
		return 0, false
	}
	return id, true
}

// GetInscripciones devuelve las inscripciones con datos relacionados para uso administrativo autorizado.
// Admite filtros opcionales: q, juego_id, id_carrera y asistio.
func GetInscripciones(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		where, args, validationErrors := buildInscripcionesFilter(r.URL.Query())
		if len(validationErrors) > 0 {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Status: "error", Message: "Los filtros de búsqueda no son válidos", Errors: validationErrors})
			return
		}
		if db == nil {
			writeJSON(w, http.StatusServiceUnavailable, ErrorResponse{Status: "error", Message: "La base de datos no está disponible"})
			return
		}

		inscripciones, err := queryInscripciones(r.Context(), db, where, args)
		if err != nil {
			log.Printf("[handlers.GetInscripciones] Error consultando inscripciones: %v", err)
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Status: "error", Message: "No se pudieron obtener las inscripciones"})
			return
		}
		writeJSON(w, http.StatusOK, inscripciones)
	}
}

// validateUpdateInscripcion verifica que la solicitud contenga al menos un campo válido.
func validateUpdateInscripcion(req *UpdateInscripcionRequest) []string {
	validationErrors := make([]string, 0)
	if req.JuegoID == nil && req.Asistio == nil {
		validationErrors = append(validationErrors, "Envía al menos uno de los campos 'juego_id' o 'asistio'.")
	}
	if req.JuegoID != nil && *req.JuegoID < 1 {
		validationErrors = append(validationErrors, "Selecciona un juego válido en 'juego_id'.")
	}
	if req.Asistio != nil && *req.Asistio != 0 && *req.Asistio != 1 {
		validationErrors = append(validationErrors, "El campo 'asistio' debe ser 0 o 1.")
	}
	return validationErrors
}

// UpdateInscripcion modifica el juego y/o la asistencia de una inscripción existente.
func UpdateInscripcion(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := parseInscripcionID(w, r)
		if !ok {
			return
		}
		var req UpdateInscripcionRequest
		if !decodeJSON(w, r, &req) {
			return
		}
		if validationErrors := validateUpdateInscripcion(&req); len(validationErrors) > 0 {
			writeJSON(w, http.StatusUnprocessableEntity, ErrorResponse{Status: "error", Message: "Corrige los datos de la inscripción", Errors: validationErrors})
			return
		}
		if db == nil {
			writeJSON(w, http.StatusServiceUnavailable, ErrorResponse{Status: "error", Message: "La base de datos no está disponible"})
			return
		}

		sets := make([]string, 0, 2)
		args := make([]any, 0, 3)
		if req.JuegoID != nil {
			sets = append(sets, "id_juego = ?")
			args = append(args, *req.JuegoID)
		}
		if req.Asistio != nil {
			sets = append(sets, "asistio = ?")
			args = append(args, *req.Asistio)
		}
		args = append(args, id)

		if _, err := db.ExecContext(r.Context(), "UPDATE inscripciones SET "+strings.Join(sets, ", ")+" WHERE id_inscripcion = ?", args...); err != nil {
			switch mysqlErrorCode(err) {
			case 1062:
				writeJSON(w, http.StatusConflict, ErrorResponse{Status: "error", Message: "El participante ya está inscrito en ese juego"})
			case 1452:
				writeJSON(w, http.StatusUnprocessableEntity, ErrorResponse{Status: "error", Message: "El juego seleccionado no existe"})
			default:
				log.Printf("[handlers.UpdateInscripcion] Error actualizando inscripcion %d: %v", id, err)
				writeJSON(w, http.StatusInternalServerError, ErrorResponse{Status: "error", Message: "No se pudo actualizar la inscripción"})
			}
			return
		}

		// Se vuelve a consultar para confirmar la existencia y devolver el estado final.
		updated, err := queryInscripciones(r.Context(), db, " WHERE i.id_inscripcion = ?", []any{id})
		if err != nil {
			log.Printf("[handlers.UpdateInscripcion] Error consultando inscripcion %d: %v", id, err)
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Status: "error", Message: "No se pudo recuperar la inscripción actualizada"})
			return
		}
		if len(updated) == 0 {
			writeJSON(w, http.StatusNotFound, ErrorResponse{Status: "error", Message: "La inscripción no existe"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"status": "success", "message": "Inscripción actualizada", "data": updated[0]})
	}
}

// DeleteInscripcion elimina una inscripción y, si el participante queda sin inscripciones,
// elimina también al participante. Todo ocurre en una única transacción.
func DeleteInscripcion(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := parseInscripcionID(w, r)
		if !ok {
			return
		}
		if db == nil {
			writeJSON(w, http.StatusServiceUnavailable, ErrorResponse{Status: "error", Message: "La base de datos no está disponible"})
			return
		}

		fail := func(context string, err error) {
			log.Printf("[handlers.DeleteInscripcion] %s (inscripcion %d): %v", context, id, err)
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Status: "error", Message: "No se pudo eliminar la inscripción"})
		}

		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			fail("Error iniciando transacción", err)
			return
		}
		defer func() { _ = tx.Rollback() }()

		var idParticipante int64
		err = tx.QueryRowContext(r.Context(), "SELECT id_participante FROM inscripciones WHERE id_inscripcion = ? FOR UPDATE", id).Scan(&idParticipante)
		if errors.Is(err, sql.ErrNoRows) {
			writeJSON(w, http.StatusNotFound, ErrorResponse{Status: "error", Message: "La inscripción no existe"})
			return
		}
		if err != nil {
			fail("Error bloqueando inscripción", err)
			return
		}

		if _, err := tx.ExecContext(r.Context(), "DELETE FROM inscripciones WHERE id_inscripcion = ?", id); err != nil {
			fail("Error eliminando inscripción", err)
			return
		}

		var restantes int
		if err := tx.QueryRowContext(r.Context(), "SELECT COUNT(*) FROM inscripciones WHERE id_participante = ?", idParticipante).Scan(&restantes); err != nil {
			fail("Error contando inscripciones restantes", err)
			return
		}

		participanteEliminado := false
		if restantes == 0 {
			if _, err := tx.ExecContext(r.Context(), "DELETE FROM participantes WHERE id_participante = ?", idParticipante); err != nil {
				fail("Error eliminando participante", err)
				return
			}
			participanteEliminado = true
		}

		if err := tx.Commit(); err != nil {
			fail("Error en tx.Commit", err)
			return
		}

		message := "Inscripción eliminada"
		if participanteEliminado {
			message = "Inscripción y participante eliminados"
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"status":  "success",
			"message": message,
			"data": map[string]any{
				"id_inscripcion":         id,
				"participante_eliminado": participanteEliminado,
			},
		})
	}
}
