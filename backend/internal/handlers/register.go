package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// RegisterRequest representa los datos de entrada para registrar un asistente.
type RegisterRequest struct {
	FullName   string `json:"full_name"`
	Email      string `json:"email"`
	Phone      string `json:"phone,omitempty"`
	TicketType string `json:"ticket_type,omitempty"`
}

// AttendeeData contiene la información confirmada del asistente registrado.
type AttendeeData struct {
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone,omitempty"`
	TicketType   string    `json:"ticket_type"`
	RegisteredAt time.Time `json:"registered_at"`
}

// RegisterResponse representa la respuesta exitosa del registro.
type RegisterResponse struct {
	Status  string       `json:"status"`
	Message string       `json:"message"`
	Data    AttendeeData `json:"data,omitempty"`
}

// ErrorResponse representa una respuesta de error estandarizada.
type ErrorResponse struct {
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Errors  []string `json:"errors,omitempty"`
}

// writeJSON envía una respuesta JSON con el código de estado indicado.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// HealthCheck maneja la verificación de salud del servicio (GET /api/health).
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{
			Status:  "error",
			Message: "Método no permitido",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"message": "PixelVerse Event API is running",
	})
}

// Register procesa el registro de un nuevo asistente (POST /api/register).
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{
			Status:  "error",
			Message: "Método no permitido. Utilice POST.",
		})
		return
	}

	var req RegisterRequest
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

	// Validaciones básicas
	var validationErrors []string
	req.FullName = strings.TrimSpace(req.FullName)
	req.Email = strings.TrimSpace(req.Email)
	req.Phone = strings.TrimSpace(req.Phone)
	req.TicketType = strings.TrimSpace(strings.ToLower(req.TicketType))

	if req.FullName == "" {
		validationErrors = append(validationErrors, "El campo 'full_name' es obligatorio.")
	}

	if req.Email == "" {
		validationErrors = append(validationErrors, "El campo 'email' es obligatorio.")
	} else if !strings.Contains(req.Email, "@") || !strings.Contains(req.Email, ".") {
		validationErrors = append(validationErrors, "El formato de 'email' no es válido.")
	}

	if req.TicketType == "" {
		req.TicketType = "general"
	} else if req.TicketType != "general" && req.TicketType != "vip" && req.TicketType != "speaker" {
		validationErrors = append(validationErrors, "El campo 'ticket_type' debe ser 'general', 'vip' o 'speaker'.")
	}

	if len(validationErrors) > 0 {
		writeJSON(w, http.StatusUnprocessableEntity, ErrorResponse{
			Status:  "error",
			Message: "Error de validación",
			Errors:  validationErrors,
		})
		return
	}

	// Respuesta exitosa (por ahora en memoria hasta conectar el driver DB)
	response := RegisterResponse{
		Status:  "success",
		Message: "Registro completado con éxito",
		Data: AttendeeData{
			FullName:     req.FullName,
			Email:        req.Email,
			Phone:        req.Phone,
			TicketType:   req.TicketType,
			RegisteredAt: time.Now().UTC(),
		},
	}

	writeJSON(w, http.StatusCreated, response)
}
