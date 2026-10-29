package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// findSQLFile busca el archivo SQL en rutas relativas comunes dentro del proyecto.
func findSQLFile(t *testing.T, filename string) string {
	t.Helper()
	candidatos := []string{
		filepath.Join("..", "..", filename),
		filepath.Join("..", filename),
		filename,
	}

	for _, path := range candidatos {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	t.Fatalf("No se pudo encontrar el archivo %s en las rutas relativas esperadas", filename)
	return ""
}

// getTestDB intenta conectarse a MySQL con MultiStatements habilitado.
// Si MySQL no está disponible, salta el test limpiamente con t.Skip.
func getTestDB(t *testing.T) *sql.DB {
	t.Helper()
	cfg := NewConfigFromEnv()
	cfg.MultiStatements = true

	db, err := Connect(cfg)
	if err != nil {
		t.Skipf("Omitiendo test: no se pudo inicializar la conexión con MySQL: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		t.Skipf("Omitiendo test: MySQL no responde en %s:%s (%v)", cfg.Host, cfg.Port, err)
	}

	return db
}

// TestCargarSchema verifica que schema.sql se ejecute correctamente y cree las tablas y vistas.
func TestCargarSchema(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	schemaPath := findSQLFile(t, "schema.sql")
	content, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("Error al leer schema.sql: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := db.ExecContext(ctx, string(content)); err != nil {
		t.Fatalf("Error al ejecutar schema.sql en la base de datos: %v", err)
	}

	// Verificar existencia de cada tabla
	tablasEsperadas := []string{"carreras", "juegos", "participantes", "inscripciones", "evento_config"}
	for _, tabla := range tablasEsperadas {
		var count int
		query := fmt.Sprintf("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = '%s'", tabla)
		if err := db.QueryRowContext(ctx, query).Scan(&count); err != nil {
			t.Fatalf("Error verificando la tabla %s: %v", tabla, err)
		}
		if count == 0 {
			t.Errorf("La tabla %q no existe en la base de datos tras ejecutar schema.sql", tabla)
		}
	}

	// Verificar existencia de la vista
	var vistaCount int
	queryVista := "SELECT COUNT(*) FROM information_schema.views WHERE table_schema = DATABASE() AND table_name = 'vista_inscripciones'"
	if err := db.QueryRowContext(ctx, queryVista).Scan(&vistaCount); err != nil {
		t.Fatalf("Error verificando vista_inscripciones: %v", err)
	}
	if vistaCount == 0 {
		t.Errorf("La vista 'vista_inscripciones' no fue creada por schema.sql")
	}
}

// TestCargarSeeder verifica que seeder.sql cargue los datos demostrativos correctamente.
func TestCargarSeeder(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	// 1. Asegurar esquema previo
	schemaPath := findSQLFile(t, "schema.sql")
	schemaContent, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("Error al leer schema.sql: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := db.ExecContext(ctx, string(schemaContent)); err != nil {
		t.Fatalf("Error preparando el esquema previo al seeder: %v", err)
	}

	// 2. Ejecutar seeder
	seederPath := findSQLFile(t, "seeder.sql")
	seederContent, err := os.ReadFile(seederPath)
	if err != nil {
		t.Fatalf("Error al leer seeder.sql: %v", err)
	}

	if _, err := db.ExecContext(ctx, string(seederContent)); err != nil {
		t.Fatalf("Error al ejecutar seeder.sql: %v", err)
	}

	// 3. Validar registros de carreras
	var totalCarreras int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM carreras").Scan(&totalCarreras); err != nil {
		t.Fatalf("Error al contar carreras: %v", err)
	}
	if totalCarreras != 5 {
		t.Errorf("Se esperaban 5 carreras en el catálogo, se obtuvieron: %d", totalCarreras)
	}

	// 4. Validar registros de juegos
	var totalJuegos int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM juegos").Scan(&totalJuegos); err != nil {
		t.Fatalf("Error al contar juegos: %v", err)
	}
	if totalJuegos != 4 {
		t.Errorf("Se esperaban 4 juegos en el catálogo, se obtuvieron: %d", totalJuegos)
	}

	// 5. Validar participantes
	var totalParticipantes int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM participantes").Scan(&totalParticipantes); err != nil {
		t.Fatalf("Error al contar participantes: %v", err)
	}
	if totalParticipantes < 9 {
		t.Errorf("Se esperaban al menos 9 participantes de prueba, se obtuvieron: %d", totalParticipantes)
	}

	// 6. Validar inscripciones
	var totalInscripciones int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM inscripciones").Scan(&totalInscripciones); err != nil {
		t.Fatalf("Error al contar inscripciones: %v", err)
	}
	if totalInscripciones == 0 {
		t.Errorf("No se registraron inscripciones de prueba")
	}

	// 7. Validar consulta de la vista
	var totalVista int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM vista_inscripciones").Scan(&totalVista); err != nil {
		t.Fatalf("Error al consultar vista_inscripciones: %v", err)
	}
	if totalVista == 0 {
		t.Errorf("vista_inscripciones debería devolver registros combinados")
	}

	// 8. Validar evento_config clave-valor
	var totalConfig int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM evento_config").Scan(&totalConfig); err != nil {
		t.Fatalf("Error al contar registros de evento_config: %v", err)
	}
	if totalConfig < 2 {
		t.Errorf("Se esperaban al menos 2 parámetros de configuración clave-valor, se obtuvieron: %d", totalConfig)
	}
}
