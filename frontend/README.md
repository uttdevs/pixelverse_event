# ⚡ Pixel Verse 2026

Bienvenido al repositorio oficial de **Pixel Verse 2026**, una aplicación web moderna con estética *cyberpunk* de alto rendimiento para el registro, seguimiento y telemetría de un torneo internacional de deportes electrónicos.

Esta versión ha sido completamente refactorizada para implementar el estándar profesional definitivo:
- **Vue 3 (Composition API con `<script setup>`)** y JavaScript limpio.
- **Tailwind CSS v4 Nativo** compilado en build-time con `@tailwindcss/vite` (sin CDNs lentos).
- **Vue Router 5** con navegación multi-página (`/`, `/reglamento`, `/registro`) y scroll suave para anclas hash.
- **Iconografía Vectorial con Lucide** (`lucide-vue-next`) y componentes SVG nativos para marcas gamer.
- **Accesibilidad Integral (A11y)** en modales: Focus Trap, cierre con tecla Escape, bloqueo de scroll y Light Dismiss.
- **Diseño 100% Responsivo** adaptado para dispositivos móviles, tablets y monitores ultra-anchos.

---

## 📂 Arquitectura y Estructura del Proyecto

```text
pixelverse_web/
├── index.html                               # Contenedor raíz limpio (sin scripts de CDN)
├── package.json                             # Dependencias del proyecto (Vue 3, Vite, Tailwind v4, Router, Lucide)
├── vite.config.js                           # Vite con plugins: @vitejs/plugin-vue y @tailwindcss/vite
├── src/
│   ├── main.js                              # Punto de entrada con registro del router y estilos
│   ├── App.vue                              # Componente raíz con RouterView y Layout global
│   ├── assets/
│   │   └── main.css                         # Tailwind v4 (@import "tailwindcss", @theme y efectos ciberpunk)
│   ├── router/
│   │   └── index.js                         # Configuración de rutas (Home, Rules, Register, scrollBehavior)
│   ├── views/
│   │   ├── HomeView.vue                     # Landing page principal con todas las secciones
│   │   ├── RulesView.vue                    # Vista dedicada con el reglamento y código de conducta
│   │   └── RegisterView.vue                 # Vista de pantalla completa para compartir enlace directo
│   ├── composables/
│   │   ├── useAudioSfx.js                   # Sintetizador Web Audio API reactivo
│   │   └── useTournamentData.js             # Datos centralizados (arenas, premios, fases y reglas)
│   └── components/
│       ├── common/
│       │   └── SocialIcons.vue              # SVGs vectoriales nativos (Twitch, Discord, X, YouTube)
│       ├── layout/
│       │   ├── ParticleBackground.vue       # Canvas 2D con ciclo de vida onMounted/onUnmounted
│       │   ├── HudBar.vue                   # Telemetría superior (ping, FPS, squads y audio SFX)
│       │   ├── AppHeader.vue                # Barra de navegación, RouterLinks y menú móvil
│       │   └── AppFooter.vue                # Pie de página y enlaces a redes sociales
│       ├── sections/
│       │   ├── HeroSection.vue              # Hero con cuenta regresiva reactiva (computed)
│       │   ├── TrophyHologram.vue           # Tarjeta holográfica 3D del trofeo gamer
│       │   ├── GameArenasSection.vue        # Filtro de categorías y tarjetas de juegos
│       │   ├── PrizePoolSection.vue         # Podio de premios ($50,000 USD)
│       │   ├── BracketSection.vue           # Simulador interactivo de llaves de playoffs
│       │   └── StreamScheduleSection.vue    # Transmisión simulada a 60 FPS y cronograma
│       └── ui/
│           └── RegistrationModal.vue        # Modal accesible A11y (Focus Trap, Escape, Scroll Lock)
└── README.md                                # Esta documentación técnica
```

---

## 🚀 Guía Técnica para el Equipo de Frontend

### 1. ¿Cómo funciona Tailwind CSS v4 aquí?
En lugar de cargar un CDN pesado de 300 KB en el navegador del cliente:
1. En [`vite.config.js`](vite.config.js) tenemos activo el plugin oficial `@tailwindcss/vite`.
2. En [`src/assets/main.css`](src/assets/main.css) usamos la sintaxis moderna de Tailwind v4:
   ```css
   @import "tailwindcss";

   @theme {
     --color-cyber-dark: #05050b;
     --color-cyber-card: #0c0c1a;
     --color-cyber-blue: #00f0ff;
     --color-cyber-purple: #9d00ff;
     --color-cyber-pink: #ff007f;
     --font-orbitron: 'Orbitron', sans-serif;
     --font-rajdhani: 'Rajdhani', sans-serif;
     --font-inter: 'Inter', sans-serif;
   }
   ```
3. Vite compila y purga el CSS en tiempo de construcción (*build-time*), generando una hoja de estilos de apenas **~9 KB gzipped**, ultra rápida y sin dependencias externas.

### 2. Navegación Multi-página con Vue Router
El enrutador en [`src/router/index.js`](src/router/index.js) administra las vistas de forma dinámica:
- `/`: **HomeView** (la experiencia completa con cuenta regresiva, arenas y brackets).
- `/reglamento`: **RulesView** (documento oficial con formato 128-tick, anti-cheat y arbitraje).
- `/registro`: **RegisterView** (enlace directo para compartir en redes o Discord).

Adicionalmente, implementa `scrollBehavior` inteligente para soportar anclas suaves:
```javascript
scrollBehavior(to, from, savedPosition) {
  if (to.hash) {
    return { el: to.hash, behavior: 'smooth', top: 80 }
  }
  return { top: 0, behavior: 'smooth' }
}
```

### 3. Iconografía Vectorial (`lucide-vue-next` + `SocialIcons.vue`)
- **Lucide Icons**: Se importan como componentes de Vue (`Gamepad2`, `Crosshair`, `Trophy`, `ShieldCheck`, `Volume2`, `Users`, etc.). Solo se empaquetan los iconos que realmente se usan en la app (*tree-shaking*).
- **Logotipos de Redes**: Encapsulados en [`SocialIcons.vue`](src/components/common/SocialIcons.vue) para evitar cargar hojas de estilos completas de FontAwesome solo por 4 logos de marcas comerciales.

### 4. Accesibilidad (A11y) en [`RegistrationModal.vue`](src/components/ui/RegistrationModal.vue)
Cumple con las mejores prácticas de la web moderna:
- **Cierre con tecla Escape**: Un listener global cierra el modal al presionar `Esc`.
- **Light Dismiss**: Al hacer clic en el backdrop fuera del modal (`@click.self`), se cierra limpiamente.
- **Bloqueo de scroll**: `document.body.style.overflow = 'hidden'` impide que la página de fondo se mueva mientras el modal está abierto.
- **Focus Trap**: El tabulador cicla entre los controles del formulario sin escapar al documento de fondo.
- **Restauración de foco**: Al cerrarse, el foco regresa al botón que abrió el modal.

---

## 🛠️ Comandos de Ejecución

```bash
# 1. Iniciar servidor de desarrollo con recarga rápida (HMR)
npm run dev

# 2. Compilar para producción (genera /dist optimizado)
npm run build

# 3. Previsualizar la compilación de producción localmente
npm run preview
```

---

## 📱 Responsividad Garantizada

- **Móviles (360px – 640px)**: Menú hamburguesa accesible, tarjetas de arenas en 1 columna, inputs con zonas táctiles cómodas, llaves de playoffs con desplazamiento horizontal suave (`overflow-x-auto`).
- **Tablets (768px – 1024px)**: Rejilla adaptativa de 2 columnas en arenas y podio en 3 niveles ordenado.
- **Monitores Grandes (1280px+)**: Visualización panorámica ciberpunk a 12 columnas con HUD de telemetría superior fijo.
