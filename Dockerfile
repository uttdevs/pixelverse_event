# ============================================================
# Stage 1: BUILD FRONTEND — compila Vue/Vite a archivos estáticos
# ============================================================
FROM node:22-alpine AS frontend

WORKDIR /build

COPY frontend/package.json frontend/package-lock.json* frontend/pnpm-lock.yaml* ./

# Instala dependencias (soporta npm o pnpm según el lockfile presente)
RUN if [ -f pnpm-lock.yaml ]; then \
      corepack enable && pnpm install --frozen-lockfile; \
    else \
      npm ci; \
    fi

COPY frontend/ .
RUN if [ -f pnpm-lock.yaml ]; then \
      pnpm run build; \
    else \
      npm run build; \
    fi

# Resultado: /build/dist/ con HTML, CSS, JS listos para producción


# ============================================================
# Stage 2: BUILD BACKEND — compila el binario Go
# ============================================================
FROM golang:1.25-alpine AS backend

WORKDIR /build

# Primero copia solo go.mod/go.sum para cachear dependencias
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copia el resto del código y compila
COPY backend/ .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o api ./cmd/api/


# ============================================================
# Stage 3: RUNTIME — imagen final mínima (~15MB)
# ============================================================
FROM alpine:3.21

# Certificados SSL para conexiones HTTPS salientes y tzdata para zonas horarias
RUN apk add --no-cache ca-certificates tzdata

# Usuario no-root por seguridad
RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

# Copia el binario compilado de Go
COPY --from=backend /build/api ./api

# Copia los archivos estáticos del frontend
COPY --from=frontend /build/dist ./static

# Copia el schema SQL (útil para inicializar la DB)
COPY backend/schema.sql ./schema.sql
COPY backend/seeder.sql ./seeder.sql

# Asigna permisos al usuario no-root
RUN chown -R app:app /app
USER app

EXPOSE 8000

# Healthcheck para Docker/Dokploy
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget -qO- http://localhost:8000/api/health || exit 1

CMD ["./api"]
