/**
 * ============================================================================
 * ARCHIVO DE ENTRADA PRINCIPAL: main.js
 * ============================================================================
 * 1. Inicializa Vue 3 con createApp.
 * 2. Registra Vue Router para la navegación multi-página.
 * 3. Importa Tailwind CSS v4 y estilos globales ciberpunk.
 * 4. Monta la app en #app.
 */

import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

// Importamos nuestros estilos globales de Tailwind v4 y efectos ciberpunk
import './assets/main.css'

const app = createApp(App)

// Conectamos el enrutador a la aplicación
app.use(router)

app.mount('#app')
