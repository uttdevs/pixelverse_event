package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestBuildInscripcionesFilterSinFiltros(t *testing.T) {
	where, args, errs := buildInscripcionesFilter(url.Values{})
	if where != "" || len(args) != 0 || len(errs) != 0 {
		t.Fatalf("se esperaba filtro vacío; obtenido where=%q args=%v errs=%v", where, args, errs)
	}
}

func TestBuildInscripcionesFilterEscapaComodines(t *testing.T) {
	where, args, errs := buildInscripcionesFilter(url.Values{"q": {"  50%_a  "}})
	if len(errs) != 0 {
		t.Fatalf("errores inesperados: %v", errs)
	}
	if !strings.Contains(where, "LIKE ?") {
		t.Fatalf("la cláusula no usa LIKE parametrizado: %q", where)
	}
	if len(args) != 5 {
		t.Fatalf("se esperaban 5 argumentos; obtenidos %d", len(args))
	}
	if args[0] != `%50\%\_a%` {
		t.Errorf("patrón obtenido %q; se esperaba %q", args[0], `%50\%\_a%`)
	}
}

func TestBuildInscripcionesFilterCombinado(t *testing.T) {
	where, args, errs := buildInscripcionesFilter(url.Values{"juego_id": {"2"}, "id_carrera": {"3"}, "asistio": {"1"}})
	if len(errs) != 0 {
		t.Fatalf("errores inesperados: %v", errs)
	}
	for _, fragment := range []string{"i.id_juego = ?", "p.id_carrera = ?", "i.asistio = ?", " AND "} {
		if !strings.Contains(where, fragment) {
			t.Errorf("la cláusula %q no contiene %q", where, fragment)
		}
	}
	if len(args) != 3 || args[0] != 2 || args[1] != 3 || args[2] != 1 {
		t.Errorf("argumentos inesperados: %v", args)
	}
}

func TestBuildInscripcionesFilterValoresInvalidos(t *testing.T) {
	_, _, errs := buildInscripcionesFilter(url.Values{
		"juego_id":   {"abc"},
		"id_carrera": {"0"},
		"asistio":    {"2"},
		"q":          {strings.Repeat("x", 101)},
	})
	if len(errs) != 4 {
		t.Fatalf("se esperaban 4 errores; obtenidos %d: %v", len(errs), errs)
	}
}

func TestGetInscripcionesRechazaFiltroInvalido(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/inscripciones?asistio=9", nil)
	GetInscripciones(nil)(recorder, request)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("código obtenido %d; se esperaba %d", recorder.Code, http.StatusBadRequest)
	}
}

// newPatchRequest construye una solicitud PATCH con el valor de ruta {id} configurado.
func newIDRequest(method, id, body string) *http.Request {
	request := httptest.NewRequest(method, "/api/inscripciones/"+id, bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	request.SetPathValue("id", id)
	return request
}

func TestUpdateInscripcionValidaciones(t *testing.T) {
	cases := []struct {
		name string
		id   string
		body string
		want int
	}{
		{"id no numérico", "abc", `{"asistio":1}`, http.StatusBadRequest},
		{"id cero", "0", `{"asistio":1}`, http.StatusBadRequest},
		{"json inválido", "1", `{`, http.StatusBadRequest},
		{"campo desconocido", "1", `{"matricula":"x"}`, http.StatusBadRequest},
		{"sin campos", "1", `{}`, http.StatusUnprocessableEntity},
		{"asistio inválido", "1", `{"asistio":5}`, http.StatusUnprocessableEntity},
		{"juego inválido", "1", `{"juego_id":0}`, http.StatusUnprocessableEntity},
		{"válido sin db", "1", `{"juego_id":2,"asistio":1}`, http.StatusServiceUnavailable},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			UpdateInscripcion(nil)(recorder, newIDRequest(http.MethodPatch, tc.id, tc.body))
			if recorder.Code != tc.want {
				t.Fatalf("código obtenido %d; se esperaba %d (cuerpo: %s)", recorder.Code, tc.want, recorder.Body.String())
			}
		})
	}
}

func TestDeleteInscripcionRechazaIDInvalido(t *testing.T) {
	recorder := httptest.NewRecorder()
	DeleteInscripcion(nil)(recorder, newIDRequest(http.MethodDelete, "-3", ""))
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("código obtenido %d; se esperaba %d", recorder.Code, http.StatusBadRequest)
	}
}

func TestAdminOnlyProtegeRutasDeInscripciones(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle("GET /api/inscripciones/export", AdminOnly("clave-secreta", ExportInscripciones(nil)))
	mux.Handle("PATCH /api/inscripciones/{id}", AdminOnly("clave-secreta", UpdateInscripcion(nil)))
	mux.Handle("DELETE /api/inscripciones/{id}", AdminOnly("clave-secreta", DeleteInscripcion(nil)))

	for _, method := range []struct{ verb, path string }{
		{http.MethodGet, "/api/inscripciones/export"},
		{http.MethodPatch, "/api/inscripciones/1"},
		{http.MethodDelete, "/api/inscripciones/1"},
	} {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(method.verb, method.path, bytes.NewBufferString(`{"asistio":1}`))
		request.Header.Set("Content-Type", "application/json")
		mux.ServeHTTP(recorder, request)
		if recorder.Code != http.StatusUnauthorized {
			t.Errorf("%s %s: código obtenido %d; se esperaba %d", method.verb, method.path, recorder.Code, http.StatusUnauthorized)
		}
	}
}

func TestBuildInscripcionesWorkbook(t *testing.T) {
	rows := []InscripcionVista{
		{IDInscripcion: 1, Juego: "FIFA 2026", Participante: "Ana López", Matricula: "0012345", Carrera: "TI", Cuatrimestre: 3, Grupo: "A", Telefono: "0441234567", Asistio: 1, FechaInscripcion: "2026-10-01 10:00:00"},
		{IDInscripcion: 2, Juego: "Tekken 8", Participante: "Luis Pérez", Matricula: "0099999", Carrera: "TI", Cuatrimestre: 5, Grupo: "B", Telefono: "5512345678", Correo: "luis@example.com"},
		{IDInscripcion: 3, Juego: "FIFA 2026", Participante: "Luis Pérez", Matricula: "0099999", Carrera: "TI", Cuatrimestre: 5, Grupo: "B", Telefono: "5512345678"},
	}
	workbook, err := buildInscripcionesWorkbook(rows)
	if err != nil {
		t.Fatalf("error generando libro: %v", err)
	}
	defer workbook.Close()

	// Reabre el archivo serializado para validar el contenido real del .xlsx.
	buffer, err := workbook.WriteToBuffer()
	if err != nil {
		t.Fatalf("error serializando libro: %v", err)
	}
	reopened, err := excelize.OpenReader(buffer)
	if err != nil {
		t.Fatalf("el archivo generado no es un .xlsx válido: %v", err)
	}
	defer reopened.Close()

	sheets := reopened.GetSheetList()
	if len(sheets) != 2 || sheets[0] != hojaInscripciones || sheets[1] != hojaResumen {
		t.Fatalf("hojas inesperadas: %v", sheets)
	}

	data, err := reopened.GetRows(hojaInscripciones)
	if err != nil {
		t.Fatalf("no se pudieron leer las filas: %v", err)
	}
	if len(data) != len(rows)+1 {
		t.Fatalf("filas obtenidas %d; se esperaban %d", len(data), len(rows)+1)
	}
	if data[0][3] != "Matrícula" || data[0][9] != "Asistió" {
		t.Errorf("encabezados inesperados: %v", data[0])
	}
	if data[1][3] != "0012345" || data[1][7] != "0441234567" {
		t.Errorf("matrícula/teléfono no conservaron ceros iniciales: %v", data[1])
	}
	if data[1][9] != "Sí" || data[2][9] != "No" {
		t.Errorf("columna de asistencia inesperada: %q, %q", data[1][9], data[2][9])
	}

	summary, err := reopened.GetRows(hojaResumen)
	if err != nil {
		t.Fatalf("no se pudo leer el resumen: %v", err)
	}
	if summary[1][1] != "3" || summary[2][1] != "1" {
		t.Errorf("totales inesperados en el resumen: %v", summary[:3])
	}
	if summary[5][0] != "FIFA 2026" || summary[5][1] != "2" || summary[6][0] != "Tekken 8" {
		t.Errorf("conteo por juego inesperado: %v", summary[5:])
	}
}
