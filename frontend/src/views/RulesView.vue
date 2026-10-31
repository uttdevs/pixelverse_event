<script setup>
/**
 * ============================================================================
 * VISTA: RulesView.vue (Reglamento Oficial del Torneo)
 * ============================================================================
 * Página dedicada con el código de conducta, políticas de arbitraje,
 * requerimientos de 128-tickrate y directrices de cobro de premios.
 */

import { RouterLink } from 'vue-router'
import { BookOpen, ShieldCheck, ArrowLeft, Zap, AlertTriangle } from 'lucide-vue-next'
import { useTournamentData } from '../composables/useTournamentData'
import { useAudioSfx } from '../composables/useAudioSfx'

const emit = defineEmits(['open-modal'])
const { tournamentRules } = useTournamentData()
const { playHoverSound, playClickSound } = useAudioSfx()

const handleOpenModal = () => {
  playClickSound()
  emit('open-modal')
}
</script>

<template>
  <div class="relative z-10 py-16 max-w-5xl mx-auto px-4 sm:px-6 lg:px-8">
    
    <!-- Breadcrumb y Botón Volver -->
    <div class="mb-8">
      <RouterLink
        to="/"
        @mouseenter="playHoverSound"
        class="inline-flex items-center gap-2 text-xs font-rajdhani uppercase tracking-widest text-cyber-blue hover:text-white transition"
      >
        <ArrowLeft class="w-4 h-4" />
        <span>Volver a la Arena Principal</span>
      </RouterLink>
    </div>

    <!-- Encabezado de la Página de Reglas -->
    <div class="text-center md:text-left mb-12">
      <div class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-cyber-purple/15 border border-cyber-purple/40 text-cyber-blue font-rajdhani text-xs tracking-widest uppercase mb-4">
        <BookOpen class="w-3.5 h-3.5" />
        <span>DOCUMENTO NORMATIVO OFICIAL // PIXELVERSE 2026</span>
      </div>

      <h1 class="font-orbitron font-black text-3xl sm:text-5xl uppercase tracking-tight">
        REGLAMENTO <span class="text-transparent bg-clip-text bg-gradient-to-r from-cyber-blue to-cyber-purple">COMPETITIVO</span>
      </h1>
      
      <p class="text-gray-300 text-sm sm:text-base mt-3 max-w-3xl font-light leading-relaxed">
        Todos los participantes inscritos en Pixel Verse deben aceptar y cumplir las siguientes cláusulas técnicas, operativas y de juego limpio.
      </p>
    </div>

    <!-- Lista de Secciones Normativas -->
    <div class="space-y-6">
      <div
        v-for="(rule, index) in tournamentRules"
        :key="index"
        class="gradient-border rounded-2xl p-6 sm:p-8 bg-cyber-card/80 border border-white/10 hover:border-cyber-blue/40 transition"
      >
        <h2 class="font-orbitron font-bold text-lg sm:text-xl text-cyber-blue mb-4 flex items-center gap-2.5">
          <span class="w-2 h-2 rounded-full bg-cyber-purple"></span>
          {{ rule.section }}
        </h2>

        <ul class="space-y-3 font-rajdhani text-sm sm:text-base text-gray-300">
          <li
            v-for="(item, iIndex) in rule.items"
            :key="iIndex"
            class="flex items-start gap-3"
          >
            <span class="text-cyber-purple font-bold mt-0.5">❯</span>
            <span class="leading-relaxed">{{ item }}</span>
          </li>
        </ul>
      </div>
    </div>


    <!-- CTA Inferior -->
    <div class="mt-12 text-center pt-8 border-t border-white/10 flex flex-col sm:flex-row items-center justify-center gap-4">
      <button
        @click="handleOpenModal"
        @mouseenter="playHoverSound"
        type="button"
        class="cyber-clip px-8 py-4 bg-gradient-to-r from-cyber-blue to-cyber-purple text-cyber-dark font-orbitron font-black text-sm uppercase tracking-wider hover:opacity-90 transition cursor-pointer flex items-center gap-2"
      >
        <Zap class="w-4 h-4 text-cyber-dark" />
        <span>ACEPTAR E INSCRIBIRME</span>
      </button>

      <RouterLink
        to="/"
        @mouseenter="playHoverSound"
        class="cyber-clip px-8 py-4 border border-white/20 bg-white/5 hover:bg-white/10 text-white font-orbitron font-bold text-sm uppercase tracking-wider transition"
      >
        VOLVER AL LOBBY
      </RouterLink>
    </div>

  </div>
</template>
