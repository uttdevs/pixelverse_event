package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetJuegosRetornaErrorConDBNil(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/juegos", nil)

	GetJuegos(nil)(recorder, request)
	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("esperaba status %d con db nil; obtuvo %d", http.StatusInternalServerError, recorder.Code)
	}
}
