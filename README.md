# pixelverse_event
Lading page for 'PixelVerse' event

Project structure

```
landing-evento/
├── .github/              # (Opcional) CI/CD o scripts de despliegue
├── backend/              # Lógica de servidor y datos
│   ├── config/
│   │   └── database.php  # Conexión PDO a MySQL
│   ├── src/
│   │   ├── Controllers/
│   │   └── Services/     # Envío de correos, validaciones, etc.
│   ├── public/           # Punto de entrada de la API
│   │   ├── .htaccess     # Enrutamiento hacia index.php
│   │   └── index.php     # Router o endpoint directo
│   ├── composer.json     # Opcional (para librerías como PHPMailer o vlucas/phpdotenv)
│   └── schema.sql        # Migraciones / estructura de la base de datos
│
├── frontend/             # Código fuente de la interfaz
│   ├── src/
│   │   ├── components/
│   │   ├── styles/
│   │   └── main.js
│   ├── public/           # Assets estáticos (logos, favicons)
│   ├── package.json
│   └── vite.config.js    # Configurado para compilar hacia dist/
│
├── deploy.sh             # Script local o de servidor para preparar el empaquetado
├── .env.example          # Variables de entorno de ejemplo
├── .gitignore
└── README.md
```