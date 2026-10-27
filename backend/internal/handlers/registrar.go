package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"net/mail"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/go-sql-driver/mysql"
)



var telefonoPermitido = regexp.MustCompile(`^[0-9+(). -]+$`)

type RegistrarParticipanteRequest struct {
	NombreCompleto       string `json:"nombre_completo"`
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
		IDParticipante int64  `json:"id_participante"`
		Matricula      string `json:"matricula"`
		JuegoID        int    `json:"juego_id"`
		IDInscripcion  int64  `json:"id_inscripcion"`
	} `json:"data"`
}

// RegistrarParticipante valida los datos e inserta al participante y su inscripción en una sola transacción.
// Si la matrícula ya existe con los mismos datos, permite inscribirla en otro juego.
func RegistrarParticipante(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegistrarParticipanteRequest
		if !decodeJSON(w, r, &req) {
			return
		}

		validationErrors := validateRegistrarParticipante(&req)
		if len(validationErrors) > 0 {
			writeJSON(w, http.StatusUnprocessableEntity, ErrorResponse{Status: "error", Message: "Corrige los datos del formulario", Errors: validationErrors})
			return
		}
		if db == nil {
			writeJSON(w, http.StatusServiceUnavailable, ErrorResponse{Status: "error", Message: "La base de datos no está disponible"})
			return
		}

		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Status: "error", Message: "No se pudo iniciar el registro"})
			return
		}
		defer func() { _ = tx.Rollback() }()

		idParticipante, err := findExistingParticipant(r, tx, req)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			if errors.Is(err, errParticipantDataMismatch) {
				writeJSON(w, http.StatusConflict, ErrorResponse{Status: "error", Message: "La matrícula ya existe con datos diferentes. Verifica la información o contacta al organizador."})
				return
			}
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Status: "error", Message: "No se pudo verificar la matrícula"})
			return
		}

		if errors.Is(err, sql.ErrNoRows) {
			var correo any
			if req.Correo != "" {
				correo = req.Correo
			}
			result, insertErr := tx.ExecContext(r.Context(), "INSERT INTO participantes (nombre_completo, matricula, id_carrera, cuatrimestre, grupo, telefono, correo) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", req.NombreCompleto, req.Matricula, req.IDCarrera, req.Cuatrimestre, req.Grupo, req.Telefono, correo)
			if insertErr != nil {
				switch mysqlErrorCode(insertErr) {
				case 1062:
					writeJSON(w, http.StatusConflict, ErrorResponse{Status: "error", Message: "La matrícula acaba de registrarse. Vuelve a enviar el formulario para continuar."})
				case 1452:
					writeJSON(w, http.StatusUnprocessableEntity, ErrorResponse{Status: "error", Message: "La carrera seleccionada no existe"})
				default:
					writeJSON(w, http.StatusInternalServerError, ErrorResponse{Status: "error", Message: "No se pudo registrar al participante"})
				}
				return
			}
			idParticipante, err = result.LastInsertId()
			if err != nil {
				writeJSON(w, http.StatusInternalServerError, ErrorResponse{Status: "error", Message: "No se pudo recuperar el identificador del participante"})
				return
			}
		}

		inscripcionResult, err := tx.ExecContext(r.Context(), "INSERT INTO inscripciones (id_participante, id_juego) VALUES (?, ?)", idParticipante, req.JuegoID)
		if err != nil {
			switch mysqlErrorCode(err) {
			case 1062:
				writeJSON(w, http.StatusConflict, ErrorResponse{Status: "error", Message: "El participante ya está inscrito en este juego"})
			case 1452:
				writeJSON(w, http.StatusUnprocessableEntity, ErrorResponse{Status: "error", Message: "El juego seleccionado no existe"})
			default:
				writeJSON(w, http.StatusInternalServerError, ErrorResponse{Status: "error", Message: "No se pudo registrar la inscripción"})
			}
			return
		}

		idInscripcion, err := inscripcionResult.LastInsertId()
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Status: "error", Message: "No se pudo recuperar el identificador de la inscripción"})
			return
		}
		if err := tx.Commit(); err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Status: "error", Message: "No se pudo confirmar el registro"})
			return
		}

		response := RegistrarParticipanteResponse{Status: "success", Message: "Registro completado con éxito"}
		response.Data.IDParticipante = idParticipante
		response.Data.Matricula = req.Matricula
		response.Data.JuegoID = req.JuegoID
		response.Data.IDInscripcion = idInscripcion
		writeJSON(w, http.StatusCreated, response)
	}
}

var errParticipantDataMismatch = errors.New("la matrícula ya existe con otros datos")

// findExistingParticipant busca una matrícula dentro de la transacción y comprueba que coincidan los datos obligatorios.
func findExistingParticipant(r *http.Request, tx *sql.Tx, req RegistrarParticipanteRequest) (int64, error) {
	var id int64
	var nombre, apellidos, matricula, grupo, telefono, correo string
	var idCarrera, cuatrimestre int
	err := tx.QueryRowContext(r.Context(), "SELECT id_participante, nombre, apellidos, matricula, id_carrera, cuatrimestre, grupo, telefono, COALESCE(correo, '') FROM participantes WHERE matricula = ? FOR UPDATE", req.Matricula).Scan(&id, &nombre, &apellidos, &matricula, &idCarrera, &cuatrimestre, &grupo, &telefono, &correo)
	if err != nil {
		return 0, err
	}
	if nombre != req.NombreCompleto || matricula != req.Matricula || idCarrera != req.IDCarrera || cuatrimestre != req.Cuatrimestre || grupo != req.Grupo || telefono != req.Telefono || (req.Correo != "" && correo != req.Correo) {
		return 0, errParticipantDataMismatch
	}
	return id, nil
}

// validateRegistrarParticipante normaliza y valida todos los campos según el esquema MySQL.
func validateRegistrarParticipante(req *RegistrarParticipanteRequest) []string {
	req.NombreCompleto = strings.TrimSpace(req.NombreCompleto)
	req.Matricula = strings.TrimSpace(req.Matricula)
	req.Grupo = strings.TrimSpace(req.Grupo)
	req.Telefono = strings.TrimSpace(req.Telefono)
	req.Correo = strings.TrimSpace(req.Correo)

	var validationErrors []string
	if req.NombreCompleto == "" {
		validationErrors = append(validationErrors, "El campo 'nombre' es obligatorio.")
	} else if utf8.RuneCountInString(req.NombreCompleto) > 60 {
		validationErrors = append(validationErrors, "El campo 'nombre' no puede superar 60 caracteres.")
	}
	if req.Matricula == "" {
		validationErrors = append(validationErrors, "El campo 'matricula' es obligatorio.")
	} else if utf8.RuneCountInString(req.Matricula) > 20 {
		validationErrors = append(validationErrors, "El campo 'matricula' no puede superar 20 caracteres.")
	}
	if req.IDCarrera < 1 {
		validationErrors = append(validationErrors, "Selecciona una carrera válida en 'id_carrera'.")
	}
	if req.Cuatrimestre < 1 || req.Cuatrimestre > 12 {
		validationErrors = append(validationErrors, "El campo 'cuatrimestre' debe estar entre 1 y 12.")
	}
	if req.Grupo == "" {
		validationErrors = append(validationErrors, "El campo 'grupo' es obligatorio.")
	} else if utf8.RuneCountInString(req.Grupo) > 5 {
		validationErrors = append(validationErrors, "El campo 'grupo' no puede superar 5 caracteres.")
	}
	if req.Telefono == "" {
		validationErrors = append(validationErrors, "El campo 'telefono' es obligatorio.")
	} else {
		digitCount := 0
		for _, character := range req.Telefono {
			if character >= '0' && character <= '9' {
				digitCount++
			}
		}
		if !telefonoPermitido.MatchString(req.Telefono) || digitCount < 10 || digitCount > 15 {
			validationErrors = append(validationErrors, "El teléfono debe contener entre 10 y 15 dígitos y solo puede incluir +, espacios, paréntesis, puntos o guiones.")
		}
		if utf8.RuneCountInString(req.Telefono) > 15 {
			validationErrors = append(validationErrors, "El campo 'telefono' no puede superar 15 caracteres.")
		}
	}
	if req.Correo != "" {
		parsedAddress, err := mail.ParseAddress(req.Correo)
		if err != nil || parsedAddress.Address != req.Correo {
			validationErrors = append(validationErrors, "El campo 'correo' no contiene una dirección de correo válida.")
		}
		if utf8.RuneCountInString(req.Correo) > 100 {
			validationErrors = append(validationErrors, "El campo 'correo' no puede superar 100 caracteres.")
		}
	}
	if req.JuegoID < 1 {
		validationErrors = append(validationErrors, "Selecciona un juego válido en 'juego_id'.")
	}
	return validationErrors
}

// mysqlErrorCode obtiene el código de error del controlador MySQL sin exponer detalles al cliente.
func mysqlErrorCode(err error) uint16 {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number
	}
	return 0
}
