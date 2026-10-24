package main

import (
	"log"
	"net/http"
	"os"

	"pixelverse/backend/internal/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	mux := http.NewServeMux()

	// Rutas de la API
	mux.HandleFunc("GET /api/health", handlers.HealthCheck)
	mux.HandleFunc("POST /api/register", handlers.Register)

	// Ruta raíz para health check simple
	mux.HandleFunc("GET /", handlers.HealthCheck)

	addr := ":" + port
	log.Printf("Iniciando PixelVerse API (Golang) en http://localhost%s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Error crítico al iniciar el servidor: %v", err)
	}
}
