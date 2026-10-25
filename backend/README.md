# Backend de PixelVerse Event

API REST desarrollada en Go con `net/http` y MySQL. El esquema de base de datos es `schema.sql` y los datos demostrativos están en `seeder.sql`.

## Requisitos

- Go 1.24 o superior, de acuerdo con `go.mod`.
- MySQL 8 o una versión compatible con el esquema.
- Una base de datos local llamada `pixelverse_db` o el nombre configurado en `DB_NAME`.

## Preparar MySQL

Desde esta carpeta del backend, ejecuta el esquema:


mysql -u root -p < schema.sql


Opcionalmente, para insertar datos demostrativos:


mysql -u root -p < seeder.sql


**Advertencia:** `seeder.sql` usa `TRUNCATE` y borra los registros de las tablas afectadas. No lo ejecutes en una base con datos reales que quieras conservar.

## Configuración

Copia `.env.example` a `.env` como plantilla, pero ten presente que Go no carga archivos `.env` de manera automática. Exporta esas variables en la terminal antes de iniciar el servidor o utiliza un gestor de variables de entorno. No subas contraseñas ni claves reales al repositorio.

Variables utilizadas:

- `PORT`: puerto HTTP de la API; por defecto `8000`.
- `DB_DRIVER`: solamente se admite `mysql`.
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASS`, `DB_NAME`: conexión MySQL.
- `ADMIN_API_KEY`: clave necesaria para consultar listados administrativos; si no se configura, esos endpoints quedan cerrados.

Ejemplo en Linux/macOS:


export PORT=8000
export DB_DRIVER=mysql
export DB_HOST=127.0.0.1
export DB_PORT=3306
export DB_USER=root
export DB_PASS='tu_contraseña_mysql'
export DB_NAME=pixelverse_db
export ADMIN_API_KEY='reemplaza-esto-por-una-clave-larga-y-aleatoria'
go mod download
go run ./cmd/api


En Windows PowerShell, usa `$env:PORT="8000"` y la misma forma `$env:NOMBRE="valor"` para cada variable.

## Endpoints

| Método | Ruta | Acceso | Uso |
|---|---|---|---|
| GET | `/api/health` | Público | Comprueba que el proceso HTTP está activo. |
| GET | `/api/ready` | Público | Comprueba la conexión a MySQL. |
| GET | `/api/carreras` | Público | Lista las carreras para el formulario. |
| GET | `/api/juegos` | Público | Lista los juegos para el formulario. |
| POST | `/api/participantes` | Público | Registra al participante y la inscripción. |
| POST | `/api/registrar/participante` | Público | Alias compatible para el mismo registro. |
| GET | `/api/participantes` | Admin | Lista datos personales; requiere `X-Admin-Key`. |
| GET | `/api/inscripciones` | Admin | Lista inscripciones y datos personales; requiere `X-Admin-Key`. |

Los listados de participantes e inscripciones requieren este encabezado HTTP: `X-Admin-Key: valor_de_ADMIN_API_KEY`. No coloques esa clave en JavaScript público del navegador. Para una interfaz administrativa de producción, implementa autenticación de usuarios y sesiones en el servidor, en lugar de exponer una clave compartida en el frontend.

## Ejemplo de registro

Envía `Content-Type: application/json` a `POST /api/participantes`:


{
  "nombre": "Ana",
  "apellidos": "Torres Ruiz",
  "matricula": "EJ-1025",
  "id_carrera": 1,
  "cuatrimestre": 4,
  "grupo": "B",
  "telefono": "2221234567",
  "correo": "ana@example.com",
  "juego_id": 1
}


Una respuesta exitosa utiliza `201 Created`:


{
  "status": "success",
  "message": "Registro completado con éxito",
  "data": {
    "id_participante": 10,
    "matricula": "EJ-1025",
    "juego_id": 1,
    "id_inscripcion": 17
  }
}


La matrícula identifica al participante. Si vuelve a enviar el formulario con la misma matrícula y los mismos datos personales, se puede inscribir en otro juego. La inscripción duplicada del mismo participante en el mismo juego responde `409 Conflict`. Si la matrícula existe con datos personales diferentes, también se responde `409 Conflict`.

## Pruebas y compilación


go test ./...
go build -o bin/api ./cmd/api


La API usa rutas con patrones por método de `net/http` (Go 1.22 o superior). El frontend del repositorio debe llamar a rutas que comiencen con `/api/`, ya que esa es la ruta que Vite tiene configurada para reenviar al puerto 8000.
