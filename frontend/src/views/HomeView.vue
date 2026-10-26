<script setup>
/**
 * ============================================================================
 * VISTA: HomeView.vue (Landing Page Principal del Torneo)
 * ============================================================================
 * Agrupa las secciones principales de la experiencia:
 * - HeroSection con Countdown
 * - GameArenasSection con selector de juegos
 * - PrizePoolSection con podio
 * - BracketSection con simulador de llaves
 * - StreamScheduleSection con reproductor simulado y cronograma
 */

import HeroSection from '../components/sections/HeroSection.vue'
import GameArenasSection from '../components/sections/GameArenasSection.vue'
import PrizePoolSection from '../components/sections/PrizePoolSection.vue'
import BracketSection from '../components/sections/BracketSection.vue'
import StreamScheduleSection from '../components/sections/StreamScheduleSection.vue'

const emit = defineEmits(['open-modal'])

const forwardOpenModal = (gameName) => {
  emit('open-modal', gameName)
}
</script>

<template>
  <div>
    <!-- Hero con Cuenta Regresiva y Trofeo 3D -->
    <HeroSection @open-modal="forwardOpenModal" />

    <!-- Arenas Oficiales con Filtro de Categorías -->
    <GameArenasSection @open-modal="forwardOpenModal" />

    <!-- Bolsa de Premios de $50,000 USD y Podio -->
    <PrizePoolSection />

    <!-- Llaves de Playoffs Interactivas -->
    <BracketSection />

    <!-- Transmisión de Stream Simulado y Fechas Clave -->
    <StreamScheduleSection />
  </div>
</template>
