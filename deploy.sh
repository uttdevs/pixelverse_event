#!/usr/bin/env bash
set -e

echo "=== Desplegando PixelVerse Event ==="

# Build Frontend
if [ -d "frontend" ]; then
    echo "Compilando frontend..."
    cd frontend
    if [ -f "package.json" ]; then
        npm install
        npm run build
    fi
    cd ..
fi

# Build Backend (Go)
if [ -d "backend" ]; then
    echo "Compilando backend (Go)..."
    cd backend
    if [ -f "go.mod" ] && command -v go &> /dev/null; then
        go build -o bin/api ./cmd/api
    fi
    cd ..
fi

echo "=== Despliegue completado con éxito ==="
