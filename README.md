# PixelVerse Event

Landing page y API de registro para el evento **PixelVerse**.

---

## Estructura del Proyecto

```text
pixelverse_event/
├── .github/              # Workflows de CI/CD
├── backend/              # API REST en Golang
│   ├── cmd/
│   │   └── api/
│   │       └── main.go       # Punto de entrada del servidor y rutas (net/http)
│   ├── internal/
│   │   ├── database/         # Configuración y conexión DB (MySQL / Postgres)
│   │   └── handlers/         # Manejadores HTTP (registro y health check)
│   ├── go.mod                # Módulo Go (pixelverse/backend)
│   ├── go.sum                # Verificación de dependencias
│   ├── schema.sql            # Esquema de base de datos (torneo de videojuegos)
│   ├── seeder.sql            # Datos iniciales y de prueba (semillas)
│
├── frontend/             # Interfaz web (Vite + JavaScript/CSS)
│   ├── src/
│   │   ├── components/
│   │   ├── styles/
│   │   └── main.js
│   ├── public/           # Assets estáticos (logos, favicons)
│   ├── package.json
│   └── vite.config.js    # Configuración de Vite y proxy a /api
│
├── deploy.sh             # Script de compilación y empaquetado (Frontend + Backend)
├── .env.example          # Plantilla de variables de entorno
├── .gitignore
└── README.md
```

---

## Requisitos

- **Go**: 1.22+ (desarrollado con Go 1.24)
- **Node.js**: v18+ y npm
- **Base de datos**: MySQL o PostgreSQL (opcional para desarrollo inicial)

---

## Configuración y Entorno

Copia el archivo de ejemplo para configurar tus variables de entorno:

```bash
cp .env.example .env
```

Variables principales:
- `PORT`: Puerto en el que corre la API (por defecto: `8000`).
- `DB_DRIVER`: Driver de base de datos (`mysql` o `postgres`).
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASS`, `DB_NAME`: Credenciales de la base de datos.

### Base de Datos

Para inicializar la base de datos con el esquema del torneo y cargar los datos de prueba:

```bash
# Crear la base de datos y tablas
mysql -u root -p < backend/schema.sql

# Opcional: cargar datos de prueba / semillas
mysql -u root -p < backend/seeder.sql
```

---

## Ejecución en Desarrollo

### 1. Backend (Golang)

Inicia el servidor API:

```bash
cd backend
go run ./cmd/api/main.go
```

La API estará disponible en `http://localhost:8000`.

Para ejecutar las pruebas unitarias:

```bash
cd backend
go test -v ./...
```

### 2. Frontend (Vite)

Inicia el servidor de desarrollo de Vite:

```bash
cd frontend
npm install
npm run dev
```

El frontend estará disponible en `http://localhost:3000` con proxy automático hacia `http://localhost:8000/api`.

---

## Endpoints de la API

| Método | Endpoint | Descripción |
|---|---|---|
| `GET` | `/api/health` | Estado de salud del servicio |
| `POST` | `/api/register` | Registro de asistente al evento |

### Ejemplo: Registro (`POST /api/register`)

**Payload:**
```json
{
  "full_name": "Ada Lovelace",
  "email": "ada@example.com",
  "phone": "+123456789",
  "ticket_type": "vip"
}
```

**Respuesta Exitosa (`201 Created`):**
```json
{
  "status": "success",
  "message": "Registro completado con éxito",
  "data": {
    "full_name": "Ada Lovelace",
    "email": "ada@example.com",
    "phone": "+123456789",
    "ticket_type": "vip",
    "registered_at": "2026-10-24T19:40:00Z"
  }
}
```

---

## Despliegue y Construcción

Para generar los artefactos de producción de frontend y backend en un solo comando:

```bash
./deploy.sh
```

Esto compilará el frontend en `frontend/dist/` y el binario de Go en `backend/bin/api`.
