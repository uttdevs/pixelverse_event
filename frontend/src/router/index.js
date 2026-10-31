/**
 * ============================================================================
 * ENRUTADOR PRINCIPAL: src/router/index.js (Vue Router 5)
 * ============================================================================
 * Define la navegación multi-página de la aplicación:
 * - `/`: Landing page oficial (HomeView) con todas las secciones del torneo.
 * - `/reglamento`: Vista dedicada con las reglas oficiales, requisitos y código de conducta.
 * - `/registro`: Vista dedicada de pantalla completa para el registro de participantes.
 * - `/admin`: Panel administrativo oculto (sin enlaces, noindex, carga diferida y
 *   layout propio sin cabecera ni pie de la landing).
 * 
 * Además soporta navegación suave hacia anclas hash (ej. /#juegos, /#premios).
 */

import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import RulesView from '../views/RulesView.vue'
import RegisterView from '../views/RegisterView.vue'

const routes = [
  {
    path: '/',
    name: 'home',
    component: HomeView,
    meta: { title: 'Pixel Verse // 2026 - Torneo Internacional de Esports' }
  },
  {
    path: '/reglamento',
    name: 'rules',
    component: RulesView,
    meta: { title: 'Reglamento Oficial // Pixel Verse 2026' }
  },
  {
    path: '/registro',
    name: 'register',
    component: RegisterView,
    meta: { title: 'Inscripción de Participante // Pixel Verse 2026' }
  },
  {
    path: '/admin',
    name: 'admin',
    // Carga diferida: el código del panel no forma parte del bundle público.
    component: () => import('../views/AdminView.vue'),
    meta: { title: 'Panel de Administración // Pixel Verse 2026', layout: 'admin', noindex: true }
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/'
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (to.hash) {
      return {
        el: to.hash,
        behavior: 'smooth',
        top: 80 // Offset para la cabecera fija
      }
    }
    if (savedPosition) {
      return savedPosition
    }
    return { top: 0, behavior: 'smooth' }
  }
})

// Actualizamos dinámicamente el título del documento según la ruta
router.afterEach((to) => {
  if (to.meta.title) {
    document.title = to.meta.title
  }

  // Evita que los buscadores indexen rutas marcadas como privadas (ej. /admin)
  let robots = document.querySelector('meta[name="robots"][data-router]')
  if (to.meta.noindex) {
    if (!robots) {
      robots = document.createElement('meta')
      robots.setAttribute('name', 'robots')
      robots.setAttribute('data-router', '')
      document.head.appendChild(robots)
    }
    robots.setAttribute('content', 'noindex, nofollow')
  } else if (robots) {
    robots.remove()
  }
})

export default router
