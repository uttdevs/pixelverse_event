<script setup>
/**
 * ============================================================================
 * COMPONENTE: GameArenasSection.vue (Arenas de Juego y Filtro Interactivo)
 * ============================================================================
 * Emplea:
 * - Iconos vectoriales de Lucide (Crosshair, Flame, Bomb, Shield, Coins).
 * - Directiva v-for y filtrado reactivo con computed().
 * - Emisión de evento open-modal con preselección de arena.
 */

import { ref, computed } from 'vue'
import { Crosshair, Flame, Bomb, Shield, Coins } from 'lucide-vue-next'
import { useTournamentData } from '../../composables/useTournamentData'
import { useAudioSfx } from '../../composables/useAudioSfx'

const emit = defineEmits(['open-modal'])

const { gameArenas, categories } = useTournamentData()
const { playHoverSound, playClickSound } = useAudioSfx()

const activeCategory = ref('all')

const filteredArenas = computed(() => {
  if (activeCategory.value === 'all') {
    return gameArenas
  }
  return gameArenas.filter(arena => arena.category === activeCategory.value)
})

const selectCategory = (categoryId) => {
  playClickSound()
  activeCategory.value = categoryId
}

const handleJoinArena = (arenaName) => {
  playClickSound()
  emit('open-modal', arenaName)
}
</script>

<template>
  <section id="juegos" class="relative z-10 py-20 max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
    
    <!-- Encabezado de Sección y Filtros -->
    <div class="flex flex-col md:flex-row md:items-end justify-between mb-12 gap-6">
      <div>
        <span class="font-rajdhani font-bold text-cyber-blue tracking-[0.25em] text-sm uppercase block mb-2">
          // SELECCIONA TU CAMPO DE BATALLA
        </span>
        <h2 class="font-orbitron font-black text-3xl sm:text-4xl uppercase">
          ARENAS <span class="text-transparent bg-clip-text bg-gradient-to-r from-cyber-blue to-cyber-purple">OFICIALES</span>
        </h2>
      </div>

      <!-- Filtro Interactivo de Categorías -->
      <div class="flex flex-wrap gap-2" role="tablist" aria-label="Filtro de arenas por categoría">
        <button
          v-for="cat in categories"
          :key="cat.id"
          @click="selectCategory(cat.id)"
          @mouseenter="playHoverSound"
          type="button"
          class="px-4 py-2 rounded-lg font-rajdhani font-bold text-sm tracking-wider uppercase transition cursor-pointer"
          :class="activeCategory === cat.id
            ? 'bg-gradient-to-r from-cyber-blue to-cyber-purple text-cyber-dark'
            : 'bg-white/5 hover:bg-white/10 text-gray-300 border border-white/10'"
        >
          {{ cat.label }}
        </button>
      </div>
    </div>

    <!-- Rejilla de Tarjetas de Juegos (Renderizado reactivo con v-for) -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <div
        v-for="arena in filteredArenas"
        :key="arena.id"
        @mouseenter="playHoverSound"
        class="game-card group gradient-border rounded-2xl overflow-hidden flex flex-col justify-between transition duration-300 hover:-translate-y-2"
        :class="arena.glowClass"
      >
        <div>
          <!-- Cabecera de la Tarjeta con Icono y Gradiente -->
          <div
            class="relative h-48 overflow-hidden bg-gradient-to-br flex items-center justify-center"
            :class="arena.bgGradient"
          >
            <!-- Renderizado dinámico de icono Lucide según arena -->
            <Crosshair
              v-if="arena.iconType === 'crosshair'"
              class="w-24 h-24 text-cyber-blue opacity-25 group-hover:scale-125 group-hover:opacity-50 transition duration-500"
            />
            <Flame
              v-else-if="arena.iconType === 'flame'"
              class="w-24 h-24 text-cyber-purple opacity-25 group-hover:scale-125 group-hover:opacity-50 transition duration-500"
            />
            <Bomb
              v-else-if="arena.iconType === 'bomb'"
              class="w-24 h-24 text-cyber-blue opacity-25 group-hover:scale-125 group-hover:opacity-50 transition duration-500"
            />
            <Shield
              v-else
              class="w-24 h-24 text-cyber-pink opacity-25 group-hover:scale-125 group-hover:opacity-50 transition duration-500"
            />

            <div class="absolute inset-0 bg-gradient-to-t from-cyber-card via-transparent to-transparent"></div>
            
            <!-- Badges de Formato y Premio -->
            <span
              class="absolute top-3 left-3 px-2.5 py-1 rounded bg-black/40 border border-cyber-blue/40 text-cyber-blue font-orbitron font-bold text-[10px]"
            >
              {{ arena.formatBadge }}
            </span>
            <span class="absolute top-3 right-3 px-2.5 py-1 rounded bg-black/60 text-white font-rajdhani font-bold text-xs flex items-center gap-1">
              <Coins class="w-3.5 h-3.5 text-cyber-blue" />
              <span>{{ arena.prize }}</span>
            </span>
          </div>

          <!-- Contenido y Descripción de la Arena -->
          <div class="p-6">
            <h3 class="font-orbitron font-bold text-xl text-white group-hover:text-cyber-blue transition">
              {{ arena.name }}
            </h3>
            <p class="text-gray-400 text-xs mt-2 leading-relaxed">
              {{ arena.description }}
            </p>

            <!-- Barra de Progreso de Cupos Disponibles -->
            <div class="mt-5">
              <div class="flex justify-between text-xs font-rajdhani mb-1">
                <span class="text-gray-400">SQUADS INSCRITOS</span>
                <span class="text-cyber-blue font-bold">
                  {{ arena.registeredSquads }} / {{ arena.totalSlots }} EQUIPOS
                </span>
              </div>
              <div class="w-full h-2 bg-white/10 rounded-full overflow-hidden">
                <div
                  class="h-full bg-gradient-to-r from-cyber-blue to-cyber-purple transition-all duration-500"
                  :style="{ width: `${(arena.registeredSquads / arena.totalSlots) * 100}%` }"
                ></div>
              </div>
            </div>
          </div>
        </div>

        <!-- Botón de Inscripción -->
        <div class="px-6 pb-6">
          <button
            @click="handleJoinArena(arena.name)"
            type="button"
            class="w-full py-2.5 rounded-lg bg-cyber-blue/10 hover:bg-cyber-blue hover:text-cyber-dark border border-cyber-blue/50 text-cyber-blue font-orbitron font-bold text-xs uppercase tracking-wider transition cursor-pointer"
          >
            {{ arena.buttonText }}
          </button>
        </div>
      </div>
    </div>

  </section>
</template>
