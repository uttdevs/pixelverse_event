# AGENTS.md

## Quick facts
- Monorepo: `backend/` (Go API) + `frontend/` (Vite/JS).
- Backend module: `pixelverse/backend` (go 1.24). Entry: `backend/cmd/api/main.go`. Handlers: `backend/internal/handlers/` (all in `register.go` including `HealthCheck`). DB config: `backend/internal/database/database.go` (MySQL/Postgres, no DB drivers imported yet).
- Frontend: Vite, dev on `:3000`, proxies `/api` to `http://localhost:8000` (`frontend/vite.config.js`). Scripts: `dev`, `build`, `preview`.
- API routes (Go 1.22+ method patterns): `GET /`, `GET /api/health`, `POST /api/register` (validation in `handlers.Register`; responses in Spanish).
- DB schema/view/seeds: `backend/schema.sql`, `backend/seeder.sql` (MySQL, creates `pixelverse_db`, tables `carreras/juegos/participantes/inscripciones`, view `vista_inscripciones`).
- Build/deploy: `./deploy.sh` builds frontend to `frontend/dist/` and Go binary to `backend/bin/api`.

## Dev commands
- Backend: `cd backend && go run ./cmd/api/main.go` (listens on `:8000` by default; `PORT` env).
- Backend tests: `cd backend && go test -v ./...`
- Frontend: `cd frontend && npm install && npm run dev` (needs backend on 8000 or proxy).health
- Build all: `./deploy.sh` (requires `go` and `npm`).

## Env/config gotchas
- Copy `.env.example` to `.env`. API reads `PORT` (default 8000). DB uses `DB_DRIVER` (`mysql`|`postgres`), `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASS`, `DB_NAME`, `DB_SSLMODE` (defaults: mysql/3306/root/""/pixelverse_db, postgres/5432).
- DB connection helpers exist but `database.Connect` expects drivers to be imported by the caller (comment notes `_ "github.com/go-sql-driver/mysql"` / `_ "github.com/lib/pq"` not present).
- `Register` currently returns in-memory success (no DB writes yet) despite schema/seeds present.

## Notes
- Frontend root is `frontend/` (Vite config sets `root: '.'`, outputs to `dist/`).
- Health endpoint returns `{"status":"ok","message":"PixelVerse Event API is running"}`.
- Schema is MySQL-specific (InnoDB, utf8mb4). Seeder truncates with `FOREIGN_KEY_CHECKS`.
