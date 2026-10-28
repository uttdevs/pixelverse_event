package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParseEventDate(t *testing.T) {
	cases := []struct {
		input   string
		wantErr bool
	}{
		{"2026-11-27T10:00:00Z", false},
		{"2026-11-27T10:00:00", false},
		{"2026-11-27 10:00:00", false},
		{"2026-11-27T10:00", false},
		{"2026-11-27", false},
		{"fecha-invalida", true},
		{"2026/11/27", true},
	}

	for _, c := range cases {
		_, err := parseEventDate(c.input)
		if c.wantErr && err == nil {
			t.Errorf("parseEventDate(%q) esperaba error, pero fue nil", c.input)
		}
		if !c.wantErr && err != nil {
			t.Errorf("parseEventDate(%q) falló inesperadamente: %v", c.input, err)
		}
	}
}

func TestUpdateEventConfigRechazaJSONInvalido(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/api/config", bytes.NewBufferString("{malformed"))
	request.Header.Set("Content-Type", "application/json")

	UpdateEventConfig(nil)(recorder, request)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("código obtenido %d; se esperaba %d", recorder.Code, http.StatusBadRequest)
	}
}

func TestUpdateEventConfigRechazaFechaInvalida(t *testing.T) {
	payload := `{"fecha_evento":"no-es-una-fecha"}`
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/api/config", bytes.NewBufferString(payload))
	request.Header.Set("Content-Type", "application/json")

	UpdateEventConfig(nil)(recorder, request)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("código obtenido %d; se esperaba %d", recorder.Code, http.StatusBadRequest)
	}

	var response ErrorResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("respuesta no es JSON válido: %v", err)
	}
	if response.Status != "error" {
		t.Errorf("status obtenido %q; se esperaba %q", response.Status, "error")
	}
}

func TestAdminOnlyProtegeUpdateEventConfig(t *testing.T) {
	adminKey := "secret-admin-key"
	handler := AdminOnly(adminKey, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// 1. Sin encabezado X-Admin-Key
	req1 := httptest.NewRequest(http.MethodPut, "/api/config", nil)
	rec1 := httptest.NewRecorder()
	handler.ServeHTTP(rec1, req1)
	if rec1.Code != http.StatusUnauthorized {
		t.Errorf("esperaba %d sin encabezado; obtuvo %d", http.StatusUnauthorized, rec1.Code)
	}

	// 2. Con clave errónea
	req2 := httptest.NewRequest(http.MethodPut, "/api/config", nil)
	req2.Header.Set("X-Admin-Key", "wrong-key")
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusUnauthorized {
		t.Errorf("esperaba %d con clave incorrecta; obtuvo %d", http.StatusUnauthorized, rec2.Code)
	}

	// 3. Con clave correcta
	req3 := httptest.NewRequest(http.MethodPut, "/api/config", nil)
	req3.Header.Set("X-Admin-Key", adminKey)
	rec3 := httptest.NewRecorder()
	handler.ServeHTTP(rec3, req3)
	if rec3.Code != http.StatusOK {
		t.Errorf("esperaba %d con clave correcta; obtuvo %d", http.StatusOK, rec3.Code)
	}
}
