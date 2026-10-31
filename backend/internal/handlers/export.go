package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/xuri/excelize/v2"
)

const (
	hojaInscripciones = "Inscripciones"
	hojaResumen       = "Resumen"
	xlsxContentType   = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
)

var encabezadosInscripciones = []string{
	"ID", "Juego", "Participante", "Matrícula", "Carrera", "Cuatrimestre",
	"Grupo", "Teléfono", "Correo", "Asistió", "Fecha de inscripción",
}

var anchosInscripciones = []float64{8, 22, 34, 16, 40, 13, 9, 17, 32, 10, 21}

// ExportInscripciones genera un archivo .xlsx con las inscripciones que cumplen los filtros
// (mismos parámetros que GET /api/inscripciones).
func ExportInscripciones(db *sql.DB) http.HandlerFunc {
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
			log.Printf("[handlers.ExportInscripciones] Error consultando inscripciones: %v", err)
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Status: "error", Message: "No se pudieron obtener las inscripciones"})
			return
		}

		workbook, err := buildInscripcionesWorkbook(inscripciones)
		if err != nil {
			log.Printf("[handlers.ExportInscripciones] Error generando libro: %v", err)
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Status: "error", Message: "No se pudo generar el archivo de Excel"})
			return
		}
		defer func() { _ = workbook.Close() }()

		buffer, err := workbook.WriteToBuffer()
		if err != nil {
			log.Printf("[handlers.ExportInscripciones] Error serializando libro: %v", err)
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Status: "error", Message: "No se pudo generar el archivo de Excel"})
			return
		}

		filename := fmt.Sprintf("inscripciones_pixelverse_%s.xlsx", time.Now().Format("20060102_1504"))
		w.Header().Set("Content-Type", xlsxContentType)
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
		w.Header().Set("Content-Length", fmt.Sprint(buffer.Len()))
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusOK)
		_, _ = buffer.WriteTo(w)
	}
}

// buildInscripcionesWorkbook crea el libro con la hoja de detalle y una hoja de resumen por juego.
func buildInscripcionesWorkbook(inscripciones []InscripcionVista) (*excelize.File, error) {
	f := excelize.NewFile()
	ok := false
	defer func() {
		if !ok {
			_ = f.Close()
		}
	}()

	if err := f.SetSheetName("Sheet1", hojaInscripciones); err != nil {
		return nil, err
	}

	headerStyle, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"4F46E5"}},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	if err != nil {
		return nil, err
	}
	textStyle, err := f.NewStyle(&excelize.Style{NumFmt: 49}) // formato "@" (texto)
	if err != nil {
		return nil, err
	}

	// --- Hoja de detalle ---
	for i, header := range encabezadosInscripciones {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		if err := f.SetCellValue(hojaInscripciones, cell, header); err != nil {
			return nil, err
		}
		column, _ := excelize.ColumnNumberToName(i + 1)
		if err := f.SetColWidth(hojaInscripciones, column, column, anchosInscripciones[i]); err != nil {
			return nil, err
		}
	}
	lastColumn, _ := excelize.ColumnNumberToName(len(encabezadosInscripciones))
	if err := f.SetCellStyle(hojaInscripciones, "A1", lastColumn+"1", headerStyle); err != nil {
		return nil, err
	}

	for index, inscripcion := range inscripciones {
		row := index + 2
		asistio := "No"
		if inscripcion.Asistio == 1 {
			asistio = "Sí"
		}
		values := []any{
			inscripcion.IDInscripcion, inscripcion.Juego, inscripcion.Participante, inscripcion.Matricula,
			inscripcion.Carrera, inscripcion.Cuatrimestre, inscripcion.Grupo, inscripcion.Telefono,
			inscripcion.Correo, asistio, inscripcion.FechaInscripcion,
		}
		cell, _ := excelize.CoordinatesToCellName(1, row)
		if err := f.SetSheetRow(hojaInscripciones, cell, &values); err != nil {
			return nil, err
		}
		// Matrícula y teléfono se guardan como texto para conservar ceros iniciales y símbolos.
		for _, column := range []string{"D", "H"} {
			ref := fmt.Sprintf("%s%d", column, row)
			if err := f.SetCellStyle(hojaInscripciones, ref, ref, textStyle); err != nil {
				return nil, err
			}
		}
	}

	lastRow := len(inscripciones) + 1
	if err := f.AutoFilter(hojaInscripciones, fmt.Sprintf("A1:%s%d", lastColumn, lastRow), nil); err != nil {
		return nil, err
	}
	if err := f.SetPanes(hojaInscripciones, &excelize.Panes{Freeze: true, YSplit: 1, TopLeftCell: "A2", ActivePane: "bottomLeft"}); err != nil {
		return nil, err
	}

	// --- Hoja de resumen ---
	if _, err := f.NewSheet(hojaResumen); err != nil {
		return nil, err
	}
	porJuego := make(map[string][2]int) // [inscritos, asistieron]
	totalAsistieron := 0
	for _, inscripcion := range inscripciones {
		counts := porJuego[inscripcion.Juego]
		counts[0]++
		if inscripcion.Asistio == 1 {
			counts[1]++
			totalAsistieron++
		}
		porJuego[inscripcion.Juego] = counts
	}
	juegos := make([]string, 0, len(porJuego))
	for juego := range porJuego {
		juegos = append(juegos, juego)
	}
	sort.Strings(juegos)

	rows := [][]any{
		{"Generado", time.Now().Format("2006-01-02 15:04")},
		{"Total de inscripciones", len(inscripciones)},
		{"Total que asistieron", totalAsistieron},
		{},
		{"Juego", "Inscritos", "Asistieron"},
	}
	for _, juego := range juegos {
		rows = append(rows, []any{juego, porJuego[juego][0], porJuego[juego][1]})
	}
	for index, values := range rows {
		cell, _ := excelize.CoordinatesToCellName(1, index+1)
		rowValues := values
		if err := f.SetSheetRow(hojaResumen, cell, &rowValues); err != nil {
			return nil, err
		}
	}
	if err := f.SetCellStyle(hojaResumen, "A5", "C5", headerStyle); err != nil {
		return nil, err
	}
	if err := f.SetColWidth(hojaResumen, "A", "A", 26); err != nil {
		return nil, err
	}
	if err := f.SetColWidth(hojaResumen, "B", "C", 18); err != nil {
		return nil, err
	}

	f.SetActiveSheet(0)
	ok = true
	return f, nil
}
