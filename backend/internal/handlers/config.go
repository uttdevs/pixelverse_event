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

		rows, err := db.QueryContext(r.Context(), "SELECT clave, valor FROM evento_config")
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Status: "error", Message: "No se pudo consultar la configuración del evento"})
			return
		}
		defer rows.Close()

		configMap := make(map[string]string)
		for rows.Next() {
			var k, v string
			if err := rows.Scan(&k, &v); err == nil {
				configMap[k] = v
			}
		}

		var fechaEvento time.Time
		if rawDate, exists := configMap["fecha_evento"]; exists && strings.TrimSpace(rawDate) != "" {
			if parsed, err := parseEventDate(rawDate); err == nil {
				fechaEvento = parsed
			} else {
				fechaEvento = time.Date(2026, 11, 27, 10, 0, 0, 0, time.UTC)
			}
		} else {
			fechaEvento = time.Date(2026, 11, 27, 10, 0, 0, 0, time.UTC)
		}

		registrosAbiertos := false
		if rawReg, exists := configMap["registros_abiertos"]; exists {
			raw := strings.TrimSpace(strings.ToLower(rawReg))
			registrosAbiertos = (raw == "1" || raw == "true" || raw == "si" || raw == "yes")
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

		rows, err := db.QueryContext(r.Context(), "SELECT clave, valor FROM evento_config")
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Status: "error", Message: "Error al consultar la configuración actual"})
			return
		}
		defer rows.Close()

		configMap := make(map[string]string)
		for rows.Next() {
			var k, v string
			if err := rows.Scan(&k, &v); err == nil {
				configMap[k] = v
			}
		}

		var fechaActual time.Time
		if rawDate, exists := configMap["fecha_evento"]; exists && strings.TrimSpace(rawDate) != "" {
			if parsed, err := parseEventDate(rawDate); err == nil {
				fechaActual = parsed
			} else {
				fechaActual = time.Date(2026, 11, 27, 10, 0, 0, 0, time.Local)
			}
		} else {
			fechaActual = time.Date(2026, 11, 27, 10, 0, 0, 0, time.Local)
		}

		registrosActuales := false
		if rawReg, exists := configMap["registros_abiertos"]; exists {
			raw := strings.TrimSpace(strings.ToLower(rawReg))
			registrosActuales = (raw == "1" || raw == "true" || raw == "si" || raw == "yes")
		}

		nuevaFecha := fechaActual
		if parsedDate != nil {
			nuevaFecha = *parsedDate
		}

		nuevosRegistros := registrosActuales
		if req.RegistrosAbiertos != nil {
			nuevosRegistros = *req.RegistrosAbiertos
		}

		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Status: "error", Message: "No se pudo iniciar la transacción"})
			return
		}
		defer tx.Rollback()

		upsertQuery := `
			INSERT INTO evento_config (clave, valor)
			VALUES (?, ?)
			ON DUPLICATE KEY UPDATE valor = VALUES(valor)`

		if parsedDate != nil {
			if _, err := tx.ExecContext(r.Context(), upsertQuery, "fecha_evento", nuevaFecha.Format("2006-01-02 15:04:05")); err != nil {
				writeJSON(w, http.StatusInternalServerError, ErrorResponse{Status: "error", Message: "No se pudo actualizar la fecha del evento"})
				return
			}
		}

		if req.RegistrosAbiertos != nil {
			val := "0"
			if nuevosRegistros {
				val = "1"
			}
			if _, err := tx.ExecContext(r.Context(), upsertQuery, "registros_abiertos", val); err != nil {
				writeJSON(w, http.StatusInternalServerError, ErrorResponse{Status: "error", Message: "No se pudo actualizar el estado de registros"})
				return
			}
		}

		if err := tx.Commit(); err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Status: "error", Message: "No se pudo confirmar la actualización"})
			return
		}

		writeJSON(w, http.StatusOK, EventConfigResponse{
			FechaEvento:       nuevaFecha.Format(time.RFC3339),
			RegistrosAbiertos: nuevosRegistros,
			HoraServidor:      time.Now().UTC().Format(time.RFC3339),
		})
	}
}
