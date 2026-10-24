package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"
)

// Config almacena los parámetros de configuración para la base de datos.
type Config struct {
	Driver   string // "mysql" o "postgres"
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string // postgres: "disable", "require", etc.
}

// NewConfigFromEnv obtiene la configuración desde variables de entorno con valores por defecto.
func NewConfigFromEnv() Config {
	driver := os.Getenv("DB_DRIVER")
	if driver == "" {
		driver = "mysql"
	}

	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		if driver == "postgres" {
			port = "5432"
		} else {
			port = "3306"
		}
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		user = "root"
	}

	password := os.Getenv("DB_PASS")
	if password == "" {
		password = ""
	}

	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "pixelverse_db"
	}

	sslmode := os.Getenv("DB_SSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}

	return Config{
		Driver:   driver,
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbname,
		SSLMode:  sslmode,
	}
}

// DSN construye la cadena de conexión según el driver configurado.
func (c Config) DSN() string {
	switch c.Driver {
	case "postgres":
		return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
	case "mysql":
		fallthrough
	default:
		// Formato estándar de MySQL DSN: username:password@tcp(host:port)/dbname?parseTime=true
		if c.Password != "" {
			return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
				c.User, c.Password, c.Host, c.Port, c.DBName)
		}
		return fmt.Sprintf("%s@tcp(%s:%s)/%s?parseTime=true",
			c.User, c.Host, c.Port, c.DBName)
	}
}

// Connect abre una conexión con la base de datos utilizando database/sql.
// NOTA: Para establecer una conexión real se requerirá importar el driver
// correspondiente (e.g. _ "github.com/go-sql-driver/mysql" o _ "github.com/lib/pq")
// cuando se autoricen dependencias externas.
func Connect(cfg Config) (*sql.DB, error) {
	db, err := sql.Open(cfg.Driver, cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("error al inicializar conexión con base de datos (%s): %w", cfg.Driver, err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
