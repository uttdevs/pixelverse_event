package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)



func TestHealthCheck(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	HealthCheck(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("código obtenido %d; se esperaba %d", recorder.Code, http.StatusOK)
	}
	var response map[string]string
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("la respuesta no es JSON válido: %v", err)
	}
	if response["status"] != "ok" {
		t.Errorf("status obtenido %q; se esperaba %q", response["status"], "ok")
	}
}

func TestRegistrarParticipanteRechazaJSONInvalido(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/participantes", bytes.NewBufferString("{"))
	request.Header.Set("Content-Type", "application/json")
	RegistrarParticipante(nil)(recorder, request)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("código obtenido %d; se esperaba %d", recorder.Code, http.StatusBadRequest)
	}
}

func TestRegistrarParticipanteRechazaDatosInvalidosAntesDeAccederADB(t *testing.T) {
	payload := `{"nombre":"","apellidos":"","matricula":"","id_carrera":0,"cuatrimestre":0,"grupo":"","telefono":"abc","juego_id":0}`
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/participantes", bytes.NewBufferString(payload))
	request.Header.Set("Content-Type", "application/json")
	RegistrarParticipante(nil)(recorder, request)
	if recorder.Code != http.StatusUnprocessableEntity {
		t.Fatalf("código obtenido %d; se esperaba %d", recorder.Code, http.StatusUnprocessableEntity)
	}
	var response ErrorResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("la respuesta no es JSON válido: %v", err)
	}
	if len(response.Errors) < 5 {
		t.Errorf("se esperaban varios errores de validación; se obtuvieron %d", len(response.Errors))
	}
}

func TestRegistrarParticipanteRechazaTipoDeContenidoIncorrecto(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/participantes", bytes.NewBufferString("nombre=Cristobal"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	RegistrarParticipante(nil)(recorder, request)
	if recorder.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("código obtenido %d; se esperaba %d", recorder.Code, http.StatusUnsupportedMediaType)
	}
}

func TestAdminOnlyDeniegaSinClaveConfigurada(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/participantes", nil)
	handler := AdminOnly("", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }))
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusServiceUnavailable {
		t.Fatalf("código obtenido %d; se esperaba %d", recorder.Code, http.StatusServiceUnavailable)
	}
}

func TestAdminOnlyRequiereClaveValida(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/participantes", nil)
	request.Header.Set("X-Admin-Key", "incorrecta")
	handler := AdminOnly("clave-segura", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }))
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("código obtenido %d; se esperaba %d", recorder.Code, http.StatusUnauthorized)
	}
}
