package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"
)

// EventConfigResponse representa la respuesta JSON con la configuración pública del evento.
type EventConfigResponse struct {
	FechaEvento       string `json:"fecha_evento"`
	RegistrosAbiertos bool   `json:"registros_abiertos"`
	HoraServidor      string `json:"hora_servidor"`
}

// UpdateEventConfigRequest representa los campos editables de la configuración del evento.
type UpdateEventConfigRequest struct {
	FechaEvento       string `json:"fecha_evento"`
	RegistrosAbiertos *bool  `json:"registros_abiertos"`
}

// parseEventDate intenta parsear una fecha en múltiples formatos comunes.
func parseEventDate(raw string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02T15:04",
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
	}
	for _, f := range formats {
		if t, err := time.ParseInLocation(f, raw, time.Local); err == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.New("formato de fecha inválido. Se espera formato ISO (ej. 2026-11-27T10:00:00)")
}

// GetEventConfig devuelve la fecha de inicio del evento, el estado de inscripciones y la hora del servidor.
func GetEventConfig(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if db == nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Status: "error", Message: "Base de datos no inicializada"})
			return
		}

		var fechaEvento time.Time
		var registrosAbiertos bool

		err := db.QueryRowContext(r.Context(), "SELECT fecha_evento, registros_abiertos FROM evento_config WHERE id = 1").Scan(&fechaEvento, &registrosAbiertos)
		if errors.Is(err, sql.ErrNoRows) {
			writeJSON(w, http.StatusOK, EventConfigResponse{
				FechaEvento:       "2026-11-27T10:00:00Z",
				RegistrosAbiertos: false,
				HoraServidor:      time.Now().UTC().Format(time.RFC3339),
			})
			return
		}
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Status: "error", Message: "No se pudo consultar la configuración del evento"})
			return
		}

		writeJSON(w, http.StatusOK, EventConfigResponse{
			FechaEvento:       fechaEvento.Format(time.RFC3339),
			RegistrosAbiertos: registrosAbiertos,
			HoraServidor:      time.Now().UTC().Format(time.RFC3339),
		})
	}
}

// UpdateEventConfig actualiza la fecha del evento y/o el estado de admisiones de registro.
func UpdateEventConfig(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UpdateEventConfigRequest
		if !decodeJSON(w, r, &req) {
			return
		}

		var parsedDate *time.Time
		if strings.TrimSpace(req.FechaEvento) != "" {
			parsed, err := parseEventDate(strings.TrimSpace(req.FechaEvento))
			if err != nil {
				writeJSON(w, http.StatusBadRequest, ErrorResponse{Status: "error", Message: err.Error()})
				return
			}
			parsedDate = &parsed
		}

		if db == nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Status: "error", Message: "Base de datos no inicializada"})
			return
		}

		var fechaActual time.Time
		var registrosActuales bool
		err := db.QueryRowContext(r.Context(), "SELECT fecha_evento, registros_abiertos FROM evento_config WHERE id = 1").Scan(&fechaActual, &registrosActuales)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Status: "error", Message: "Error al consultar la configuración actual"})
			return
		}

		nuevaFecha := fechaActual
		if parsedDate != nil {
			nuevaFecha = *parsedDate
		} else if errors.Is(err, sql.ErrNoRows) {
			nuevaFecha = time.Date(2026, 11, 27, 10, 0, 0, 0, time.Local)
		}

		nuevosRegistros := registrosActuales
		if req.RegistrosAbiertos != nil {
			nuevosRegistros = *req.RegistrosAbiertos
		}

		query := `
			INSERT INTO evento_config (id, fecha_evento, registros_abiertos)
			VALUES (1, ?, ?)
			ON DUPLICATE KEY UPDATE fecha_evento = VALUES(fecha_evento), registros_abiertos = VALUES(registros_abiertos)`

		if _, err := db.ExecContext(r.Context(), query, nuevaFecha, nuevosRegistros); err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Status: "error", Message: "No se pudo actualizar la configuración del evento"})
			return
		}

		writeJSON(w, http.StatusOK, EventConfigResponse{
			FechaEvento:       nuevaFecha.Format(time.RFC3339),
			RegistrosAbiertos: nuevosRegistros,
			HoraServidor:      time.Now().UTC().Format(time.RFC3339),
		})
	}
}
