package database

import (
	"database/sql"
	"fmt"
	"net"
	"os"
	"time"

	mysqlDriver "github.com/go-sql-driver/mysql"
)


// Config almacena los parámetros de conexión de la base de datos MySQL.
type Config struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// NewConfigFromEnv lee la configuración desde variables de entorno.
func NewConfigFromEnv() Config {
	return Config{
		Driver:   envOrDefault("DB_DRIVER", "mysql"),
		Host:     envOrDefault("DB_HOST", "127.0.0.1"),
		Port:     envOrDefault("DB_PORT", "3306"),
		User:     envOrDefault("DB_USER", "root"),
		Password: os.Getenv("DB_PASS"),
		DBName:   envOrDefault("DB_NAME", "pixelverse_db"),
	}
}

// envOrDefault devuelve una variable de entorno o un valor predeterminado si está vacía.
func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

// DSN construye una cadena de conexión MySQL escapando correctamente los parámetros.
func (c Config) DSN() string {
	return (&mysqlDriver.Config{
		User:         c.User,
		Passwd:       c.Password,
		Net:          "tcp",
		Addr:         net.JoinHostPort(c.Host, c.Port),
		DBName:       c.DBName,
		ParseTime:    true,
		Loc:          time.Local,
		Timeout:      5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}).FormatDSN()
}

// Connect crea el pool de conexiones de MySQL.
func Connect(cfg Config) (*sql.DB, error) {
	if cfg.Driver != "mysql" {
		return nil, fmt.Errorf("DB_DRIVER=%q no está soportado: este proyecto utiliza MySQL", cfg.Driver)
	}
	if cfg.Host == "" || cfg.Port == "" || cfg.User == "" || cfg.DBName == "" {
		return nil, fmt.Errorf("la configuración de MySQL contiene valores obligatorios vacíos")
	}

	db, err := sql.Open("mysql", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("no se pudo inicializar la conexión con MySQL: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)
	return db, nil
}
