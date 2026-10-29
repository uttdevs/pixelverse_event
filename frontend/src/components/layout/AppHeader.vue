<script setup>
/**
 * ============================================================================
 * COMPONENTE: AppHeader.vue (Navegación Principal, Router y Menú Móvil)
 * ============================================================================
 * Integra:
 * - Iconos Lucide vectoriales (Gamepad2, Zap, Menu, X, BookOpen).
 * - Enlaces combinados de Vue Router (para rutas como /reglamento) y anclas de sección.
 * - Soporte responsivo para móviles con cajón interactivo.
 */

import { ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { Gamepad2, Zap, Menu, X, BookOpen } from 'lucide-vue-next'
import { useAudioSfx } from '../../composables/useAudioSfx'

const emit = defineEmits(['open-modal'])
const router = useRouter()
const { playHoverSound, playClickSound } = useAudioSfx()

const isMobileMenuOpen = ref(false)

const toggleMobileMenu = () => {
  playClickSound()
  isMobileMenuOpen.value = !isMobileMenuOpen.value
}

const handleOpenModal = () => {
  playClickSound()
  isMobileMenuOpen.value = false
  emit('open-modal')
}

// Navegación hacia secciones de la landing
const navigateToSection = (sectionId) => {
  isMobileMenuOpen.value = false
  router.push({ path: '/', hash: sectionId })
}
</script>

<template>
  <header class="sticky top-0 z-40 backdrop-blur-md bg-cyber-dark/80 border-b border-cyber-blue/20">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-20 flex items-center justify-between">
      
      <!-- Logotipo de la Marca -->
      <RouterLink
        to="/"
        @mouseenter="playHoverSound"
        class="flex items-center gap-3 group focus:outline-none"
        aria-label="Pixel Verse Inicio"
      >
        <div class="w-11 h-11 rounded-lg bg-gradient-to-br from-cyber-blue to-cyber-purple p-0.5 box-glow-combo group-hover:scale-105 transition">
          <div class="w-full h-full bg-cyber-dark rounded-[7px] flex items-center justify-center">
            <Gamepad2 class="w-6 h-6 text-cyber-blue group-hover:text-cyber-purple transition" />
          </div>
        </div>
        <div>
          <span class="font-orbitron font-black text-xl tracking-wider text-white block leading-none">
            PIXEL<span class="text-cyber-blue text-glow-blue">VERSE</span>
          </span>
          
        </div>
      </RouterLink>

      <!-- Menú de Escritorio (Desktop Navigation) -->
      <nav class="hidden md:flex items-center space-x-7 font-rajdhani font-semibold text-base tracking-wider">
        <a
          href="/#juegos"
          @click.prevent="navigateToSection('#juegos')"
          @mouseenter="playHoverSound"
          class="text-gray-300 hover:text-cyber-blue transition relative py-1 after:content-[''] after:absolute after:bottom-0 after:left-0 after:w-0 after:h-0.5 after:bg-cyber-blue hover:after:w-full after:transition-all"
        >
          ARENAS
        </a>
       
      
        <a
          href="/#cronograma"
          @click.prevent="navigateToSection('#cronograma')"
          @mouseenter="playHoverSound"
          class="text-gray-300 hover:text-cyber-blue transition relative py-1 after:content-[''] after:absolute after:bottom-0 after:left-0 after:w-0 after:h-0.5 after:bg-cyber-blue hover:after:w-full after:transition-all"
        >
          CRONOGRAMA
        </a>
        
        <!-- Enlace de Vue Router a la página dedicada de Reglamento -->
        <RouterLink
          to="/reglamento"
          @mouseenter="playHoverSound"
          class="text-purple-300 hover:text-cyber-blue transition flex items-center gap-1.5 py-1"
          active-class="text-cyber-blue border-b border-cyber-blue"
        >
          <BookOpen class="w-4 h-4" />
          <span>REGLAS</span>
        </RouterLink>
      </nav>

   

      <!-- Botón Hamburguesa para Móviles -->
      <button
        @click="toggleMobileMenu"
        type="button"
        class="md:hidden text-cyber-blue text-2xl p-2 focus:outline-none cursor-pointer"
        aria-label="Abrir menú de navegación móvil"
        :aria-expanded="isMobileMenuOpen"
      >
        <X v-if="isMobileMenuOpen" class="w-6 h-6" />
        <Menu v-else class="w-6 h-6" />
      </button>

    </div>

    <!-- Menú Desplegable Móvil -->
    <div
      v-if="isMobileMenuOpen"
      class="md:hidden bg-cyber-card/95 border-b border-cyber-purple/30 px-6 py-5 space-y-4 font-rajdhani text-lg"
    >
      <a
        href="/#juegos"
        @click.prevent="navigateToSection('#juegos')"
        class="block text-gray-200 hover:text-cyber-blue"
      >
        ARENAS DE JUEGO
      </a>
    
      
      <a
        href="/#cronograma"
        @click.prevent="navigateToSection('#cronograma')"
        class="block text-gray-200 hover:text-cyber-blue"
      >
        CRONOGRAMA
      </a>
      <RouterLink
        to="/reglamento"
        @click="isMobileMenuOpen = false"
        class="block text-purple-300 hover:text-cyber-blue"
      >
        REGLAMENTO OFICIAL
      </RouterLink>
      <button
        @click="handleOpenModal"
        type="button"
        class="w-full py-3 rounded bg-gradient-to-r from-cyber-blue to-cyber-purple font-orbitron font-bold text-sm text-cyber-dark uppercase cursor-pointer"
      >
        INSCRIBIRSE AHORA
      </button>
    </div>
  </header>
</template>
