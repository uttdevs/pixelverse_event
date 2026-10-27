<script setup>
/**
 * ============================================================================
 * COMPONENTE: HudBar.vue (Barra Superior de Telemetría Gamer)
 * ============================================================================
 * Implementa telemetría dinámica e iconos vectoriales optimizados de Lucide:
 * - Users: Contador de escuadrones
 * - Volume2 / VolumeX: Control reactivo de efectos de sonido
 */

import { ref, onMounted, onUnmounted } from 'vue'
import { Users, Volume2, VolumeX } from 'lucide-vue-next'
import { useAudioSfx } from '../../composables/useAudioSfx'

const props = defineProps({
  squadsCount: {
    type: Number,
    default: 118
  },
  maxSquads: {
    type: Number,
    default: 128
  }
})

const { isAudioEnabled, toggleAudio, playHoverSound } = useAudioSfx()

const currentPing = ref(14)
const currentFps = ref(240)
let telemetryInterval = null

onMounted(() => {
  telemetryInterval = setInterval(() => {
    currentPing.value = 12 + Math.floor(Math.random() * 5)
    currentFps.value = 236 + Math.floor(Math.random() * 9)
  }, 2000)
})

onUnmounted(() => {
  if (telemetryInterval) {
    clearInterval(telemetryInterval)
  }
})
</script>

<template>
  <aside
    aria-label="Telemetría del torneo"
    class="relative z-40 bg-cyber-dark/90 border-b border-white/10 px-4 py-1.5 text-xs font-rajdhani tracking-widest uppercase"
  >
    <div class="max-w-7xl mx-auto flex flex-wrap justify-between items-center gap-2">
      <!-- Indicadores de servidor, latencia y FPS -->
      <div class="flex items-center space-x-4">
        <span class="flex items-center text-cyber-blue">
          <span class="w-2 h-2 rounded-full bg-cyber-blue animate-ping mr-2"></span>
          SERVER: LATAM-NORTH-01 [ONLINE]
        </span>
        <span class="hidden sm:inline text-gray-400">|</span>
        <span class="hidden sm:inline text-gray-300">
          PING: <strong class="text-emerald-400">{{ currentPing }}ms</strong>
        </span>
        <span class="hidden md:inline text-gray-400">|</span>
        <span class="hidden md:inline text-gray-300">
          FPS: <strong class="text-cyber-blue">{{ currentFps }}</strong>
        </span>
      </div>

      <!-- Squads inscritos y conmutador de audio SFX -->
      <div class="flex items-center space-x-4 ml-auto">
        <span class="text-purple-300 flex items-center">
          <Users class="w-3.5 h-3.5 inline mr-1.5 text-cyber-purple" />
          SQUADS REGISTRADOS:
          <strong class="text-white ml-1">{{ props.squadsCount }} / {{ props.maxSquads }}</strong>
        </span>

        <!-- Botón de efectos de sonido -->
        <button
          @click="toggleAudio"
          @mouseenter="playHoverSound"
          type="button"
          class="px-2.5 py-0.5 rounded bg-white/5 hover:bg-cyber-blue/20 border border-cyber-blue/40 text-cyber-blue transition flex items-center gap-1.5 cursor-pointer focus:outline-none focus:ring-1 focus:ring-cyber-blue"
          :title="isAudioEnabled ? 'Desactivar Efectos de Sonido' : 'Activar Efectos de Sonido'"
        >
          <Volume2 v-if="isAudioEnabled" class="w-3.5 h-3.5 text-cyber-blue" />
          <VolumeX v-else class="w-3.5 h-3.5 text-gray-400" />
          <span>SFX: {{ isAudioEnabled ? 'ON' : 'OFF' }}</span>
        </button>
      </div>
    </div>
  </aside>
</template>
