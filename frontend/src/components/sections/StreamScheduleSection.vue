<script setup>
/**
 * ============================================================================
 * COMPONENTE: StreamScheduleSection.vue (Transmisión en Vivo y Cronograma)
 * ============================================================================
 * Utiliza Lucide icons: Eye, Play, Pause, Mic, Gift.
 */

import { ref } from 'vue'
import { Eye, Play, Pause, Mic, Gift } from 'lucide-vue-next'
import { useTournamentData } from '../../composables/useTournamentData'
import { useAudioSfx } from '../../composables/useAudioSfx'

const { scheduleStages } = useTournamentData()
const { playHoverSound, playClickSound } = useAudioSfx()

const isStreamPlaying = ref(false)

const toggleStreamSimulation = () => {
  playClickSound()
  isStreamPlaying.value = !isStreamPlaying.value
}
</script>

<template>
  <section id="cronograma" class="relative z-10 py-20 max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-10 items-start">
      
    

      <!-- Columna Derecha: Cronograma Oficial del Torneo (5 columnas) -->
      <div class="lg:col-span-5 space-y-6">
        <div>
          <span class="font-rajdhani font-bold text-cyber-purple tracking-[0.25em] text-xs uppercase block">
            // FECHAS CLAVE
          </span>
          <h3 class="font-orbitron font-black text-2xl uppercase">
            CRONOGRAMA <span class="text-cyber-blue">OFICIAL</span>
          </h3>
        </div>

        <div class="space-y-4">
          <div
            v-for="(stage, sIndex) in scheduleStages"
            :key="sIndex"
            class="gradient-border rounded-xl p-4 flex items-center gap-4"
            :class="stage.borderHighlight || ''"
          >
            <!-- Badge de Fecha -->
            <div
              class="text-center border rounded-lg px-3 py-2 min-w-[70px]"
              :class="stage.badgeBg"
            >
              <span class="font-orbitron font-bold text-lg block">{{ stage.day }}</span>
              <span class="font-rajdhani text-[10px] uppercase block">{{ stage.month }}</span>
            </div>

            <!-- Detalles de la Fase -->
            <div>
              <span
                class="text-[10px] font-rajdhani uppercase tracking-widest block"
                :class="stage.statusColor"
              >
                {{ stage.status }}
              </span>
              <h4 class="font-orbitron font-bold text-sm text-white">{{ stage.title }}</h4>
              <p class="text-xs text-gray-400 mt-0.5">{{ stage.description }}</p>
            </div>
          </div>
        </div>

      </div>

    </div>
  </section>
</template>
