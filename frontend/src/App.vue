<script setup>
/**
 * ============================================================================
 * COMPONENTE RAÍZ: App.vue (El Orquestador Central con Vue Router)
 * ============================================================================
 * ¿QUÉ APRENDERÁ EL EQUIPO CON ESTE ARCHIVO?
 * 1. `<RouterView>`:
 *    En lugar de quemar todas las secciones fijas en el template, `<RouterView>`
 *    renderiza dinámicamente la vista activa según la URL actual:
 *    - `/`: Renderiza `HomeView.vue` (la landing completa con todas sus secciones).
 *    - `/reglamento`: Renderiza `RulesView.vue` (código normativo y de arbitraje).
 *    - `/registro`: Renderiza `RegisterView.vue` (formulario dedicado de escuadrones).
 * 2. Persistencia de Layout Global:
 *    Componentes como `ParticleBackground`, `HudBar`, `AppHeader`, `AppFooter`
 *    y `RegistrationModal` permanecen montados en todo momento. La navegación
 *    entre páginas no recarga el navegador ni interrumpe la animación del canvas
 *    ni la música/sonido.
 */

import { ref } from 'vue'
import { RouterView } from 'vue-router'

// Componentes de Layout Global
import ParticleBackground from './components/layout/ParticleBackground.vue'
import HudBar from './components/layout/HudBar.vue'
import AppHeader from './components/layout/AppHeader.vue'
import AppFooter from './components/layout/AppFooter.vue'

// Modal de Registro Global (Accesible desde cualquier ruta)
import RegistrationModal from './components/ui/RegistrationModal.vue'

// ----------------------------------------------------------------------------
// ESTADO REACTIVO CENTRAL (Single Source of Truth)
// ----------------------------------------------------------------------------
const isModalOpen = ref(false)
const selectedGameForModal = ref('Valorant: Protocol')
const registeredSquadsCount = ref(118)
const maxSquadsCount = 128

// ----------------------------------------------------------------------------
// MANEJADORES DE EVENTOS
// ----------------------------------------------------------------------------
const handleOpenModal = (gameName) => {
  if (gameName && typeof gameName === 'string') {
    selectedGameForModal.value = gameName
  } else {
    selectedGameForModal.value = 'Valorant: Protocol'
  }
  isModalOpen.value = true
}

const handleCloseModal = () => {
  isModalOpen.value = false
}

const handleSquadRegistered = (squadData) => {
  console.log('[Pixel Verse] Nuevo escuadrón registrado:', squadData)
  if (registeredSquadsCount.value < maxSquadsCount) {
    registeredSquadsCount.value++
  }
}
</script>

<template>
  <div class="min-h-screen bg-cyber-dark text-white selection:bg-cyber-blue selection:text-cyber-dark flex flex-col justify-between">
    
    <!-- 1. Fondo de Partículas Interactivas en Canvas (Global) -->
    <ParticleBackground />

    <!-- 2. Barra Superior de Telemetría Gamer (Global) -->
    <HudBar
      :squads-count="registeredSquadsCount"
      :max-squads="maxSquadsCount"
    />

    <!-- 3. Navegación Principal y Logotipo (Global) -->
    <AppHeader @open-modal="handleOpenModal" />

    <!-- 4. Contenedor de Vistas Dinámicas (Vue Router) -->
    <main class="flex-grow">
      <!-- RouterView renderiza HomeView, RulesView o RegisterView según la URL -->
      <RouterView
        @open-modal="handleOpenModal"
        @registered="handleSquadRegistered"
      />
    </main>

    <!-- 5. Pie de Página y Redes Sociales (Global) -->
    <AppFooter />

    <!-- 6. Modal de Registro de Escuadrón (Global y Accesible con A11y) -->
    <RegistrationModal
      :is-open="isModalOpen"
      :initial-game="selectedGameForModal"
      @close="handleCloseModal"
      @registered="handleSquadRegistered"
    />

  </div>
</template>
