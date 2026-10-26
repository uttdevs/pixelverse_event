<script setup>
/**
 * ============================================================================
 * COMPONENTE: AppFooter.vue (Pie de Página Oficial)
 * ============================================================================
 * Utiliza:
 * - Lucide Gamepad2 y ShieldCheck.
 * - SocialIcons.vue para enlaces vectoriales de Twitch, Discord, X y YouTube.
 * - Enlaces semánticos a Inicio, Reglamento y Registro.
 */

import { RouterLink } from 'vue-router'
import { Gamepad2, ShieldCheck } from 'lucide-vue-next'
import SocialIcons from '../common/SocialIcons.vue'
import { useAudioSfx } from '../../composables/useAudioSfx'

const { playHoverSound } = useAudioSfx()
</script>

<template>
  <footer class="relative z-10 bg-cyber-dark border-t border-white/10 pt-14 pb-10">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex flex-col md:flex-row justify-between items-center gap-6 pb-8 border-b border-white/10">
        
        <!-- Logo en Footer -->
        <RouterLink
          to="/"
          @mouseenter="playHoverSound"
          class="flex items-center gap-3 focus:outline-none"
        >
          <div class="w-9 h-9 rounded bg-gradient-to-br from-cyber-blue to-cyber-purple flex items-center justify-center text-cyber-dark font-bold">
            <Gamepad2 class="w-5 h-5 text-cyber-dark" />
          </div>
          <span class="font-orbitron font-black text-lg tracking-wider">
            PIXEL<span class="text-cyber-blue">VERSE</span> // 2026
          </span>
        </RouterLink>

        <!-- Navegación de pie de página -->
        <nav class="flex flex-wrap items-center gap-6 text-xs font-rajdhani uppercase tracking-wider text-gray-400">
          <RouterLink to="/" class="hover:text-cyber-blue transition">Inicio</RouterLink>
          <RouterLink to="/reglamento" class="hover:text-cyber-blue transition">Reglamento Oficial</RouterLink>
          
        </nav>

        <!-- Redes Sociales Gamer -->
        <SocialIcons />
      </div>

      <!-- Leyendas legales y tecnología -->
      <div class="mt-6 flex flex-col sm:flex-row justify-between items-center text-xs text-gray-500 font-rajdhani gap-4">
        <p>© 2026 PIXEL VERSE ESPORTS CIRCUIT. ESTÁNDAR VUE 3 + TAILWIND V4.</p>
        <p class="text-gray-400 flex items-center gap-1.5">
          <ShieldCheck class="w-4 h-4 text-cyber-blue" />
          <span>SERVIDORES PROTEGIDOS CON <strong class="text-cyber-blue">PIXEL VERSE2026 v4.2</strong></span>
        </p>
      </div>
    </div>
  </footer>
</template>
