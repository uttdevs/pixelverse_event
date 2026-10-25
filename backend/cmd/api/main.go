package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"pixelverse/backend/internal/database"
	"pixelverse/backend/internal/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	cfg := database.NewConfigFromEnv()
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Error al hacer ping a la base de datos: %v", err)
	}

	mux := http.NewServeMux()

	// Rutas
	mux.HandleFunc("GET /api/health", handlers.HealthCheck)
	mux.HandleFunc("POST /registrar/participante", handlers.RegistrarParticipante(db))
	mux.HandleFunc("GET /juegos", handlers.GetJuegos(db))
	mux.HandleFunc("GET /carreras", handlers.GetCarreras(db))
	mux.HandleFunc("GET /participantes", handlers.GetParticipantes(db))
	mux.HandleFunc("GET /inscripciones", handlers.GetInscripciones(db))

	addr := ":" + port
	log.Printf("Iniciando PixelVerse API (Golang) en http://localhost%s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Error crítico al iniciar el servidor: %v", err)
	}
}
