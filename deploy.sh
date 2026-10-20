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

# Install Backend dependencies
if [ -d "backend" ]; then
    echo "Instalando dependencias de backend..."
    cd backend
    if [ -f "composer.json" ] && command -v composer &> /dev/null; then
        composer install --no-dev --optimize-autoloader
    fi
    cd ..
fi

echo "=== Despliegue completado con éxito ==="
