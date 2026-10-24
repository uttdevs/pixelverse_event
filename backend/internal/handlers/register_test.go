package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rr := httptest.NewRecorder()

	HealthCheck(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("HealthCheck retornó código erróneo: obtuvo %v esperado %v", status, http.StatusOK)
	}

	var res map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Fatalf("Respuesta no es JSON válido: %v", err)
	}

	if res["status"] != "ok" {
		t.Errorf("status esperado 'ok', obtenido '%s'", res["status"])
	}
}

func TestRegisterSuccess(t *testing.T) {
	payload := []byte(`{"full_name":"Test User","email":"test@example.com","ticket_type":"vip"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	Register(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Register retornó código erróneo: obtuvo %v esperado %v", status, http.StatusCreated)
	}

	var res RegisterResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Fatalf("Respuesta no es JSON válido: %v", err)
	}

	if res.Status != "success" {
		t.Errorf("status esperado 'success', obtenido '%s'", res.Status)
	}
	if res.Data.FullName != "Test User" {
		t.Errorf("FullName esperado 'Test User', obtenido '%s'", res.Data.FullName)
	}
	if res.Data.TicketType != "vip" {
		t.Errorf("TicketType esperado 'vip', obtenido '%s'", res.Data.TicketType)
	}
}

func TestRegisterValidationErrors(t *testing.T) {
	payload := []byte(`{"full_name":"","email":"invalid-email"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	Register(rr, req)

	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("Register retornó código erróneo: obtuvo %v esperado %v", status, http.StatusUnprocessableEntity)
	}

	var res ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Fatalf("Respuesta no es JSON válido: %v", err)
	}

	if res.Status != "error" {
		t.Errorf("status esperado 'error', obtenido '%s'", res.Status)
	}
	if len(res.Errors) < 2 {
		t.Errorf("se esperaban al menos 2 errores de validación, se obtuvieron: %d", len(res.Errors))
	}
}
