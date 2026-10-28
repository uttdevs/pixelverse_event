package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"pixelverse/backend/internal/database"
	"pixelverse/backend/internal/handlers"
)



// main configura la conexión con MySQL, registra las rutas y ejecuta el servidor HTTP.
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	cfg := database.NewConfigFromEnv()
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("No se pudo inicializar la conexión con la base de datos: %v", err)
	}
	defer db.Close()

	pingCtx, cancelPing := context.WithTimeout(context.Background(), 5*time.Second)
	if err := db.PingContext(pingCtx); err != nil {
		cancelPing()
		log.Fatalf("No se pudo establecer conexión con MySQL: %v", err)
	}
	cancelPing()

	mux := http.NewServeMux()
	registrarParticipante := handlers.RegistrarParticipante(db)

	// Rutas públicas: estado del servicio, configuración y catálogos utilizados por el formulario.
	mux.HandleFunc("GET /api/health", handlers.HealthCheck)
	mux.HandleFunc("GET /api/ready", handlers.ReadyCheck(db))
	mux.HandleFunc("GET /api/config", handlers.GetEventConfig(db))
	mux.HandleFunc("GET /api/carreras", handlers.GetCarreras(db))
	mux.HandleFunc("GET /api/juegos", handlers.GetJuegos(db))

	// Registro público. Se admite la ruta REST y la ruta descriptiva para compatibilidad.
	mux.Handle("POST /api/participantes", registrarParticipante)
	mux.Handle("POST /api/registrar/participante", registrarParticipante)

	// Las consultas y modificaciones que requieren autorización administrativa usan X-Admin-Key.
	adminKey := os.Getenv("ADMIN_API_KEY")
	mux.Handle("PUT /api/config", handlers.AdminOnly(adminKey, handlers.UpdateEventConfig(db)))
	mux.Handle("GET /api/participantes", handlers.AdminOnly(adminKey, handlers.GetParticipantes(db)))
	mux.Handle("GET /api/inscripciones", handlers.AdminOnly(adminKey, handlers.GetInscripciones(db)))

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	shutdownCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		<-shutdownCtx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Error al cerrar el servidor ordenadamente: %v", err)
		}
	}()

	log.Printf("PixelVerse Event API escuchando en http://localhost:%s", port)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Error al iniciar el servidor HTTP: %v", err)
	}
}
