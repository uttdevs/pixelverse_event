<script setup>
/**
 * ============================================================================
 * COMPONENTE: HeroSection.vue (Sección Principal y Reloj de Cuenta Regresiva)
 * ============================================================================
 * Sincronizado dinámicamente con la base de datos vía /api/config
 * Utiliza Lucide icons:
 * - Clock, Crosshair, ShieldCheck, Trophy, Headphones
 */

import { ref, computed, onMounted, onUnmounted } from 'vue'
import { Clock, Crosshair, ShieldCheck, Trophy, Headphones } from 'lucide-vue-next'
import { useAudioSfx } from '../../composables/useAudioSfx'
import TrophyHologram from './TrophyHologram.vue'

const emit = defineEmits(['open-modal'])
const { playHoverSound, playClickSound } = useAudioSfx()

// Fecha fallback: 27 de noviembre de 2026 10:00:00
const DEFAULT_EVENT_DATE = '2026-11-27T10:00:00'
const targetTimestamp = ref(new Date(DEFAULT_EVENT_DATE).getTime())
const registrationsOpen = ref(false)
const serverOffset = ref(0)
const totalSeconds = ref(0)
let timerInterval = null

const calculateRemainingSeconds = () => {
  if (!targetTimestamp.value) return 0
  const now = Date.now() + serverOffset.value
  return Math.max(0, Math.floor((targetTimestamp.value - now) / 1000))
}

const updateCountdown = () => {
  totalSeconds.value = calculateRemainingSeconds()
  if (totalSeconds.value <= 0 && timerInterval) {
    clearInterval(timerInterval)
  }
}

const days = computed(() => {
  const d = Math.floor(totalSeconds.value / 86400)
  return String(d).padStart(2, '0')
})

const hours = computed(() => {
  const h = Math.floor((totalSeconds.value % 86400) / 3600)
  return String(h).padStart(2, '0')
})

const minutes = computed(() => {
  const m = Math.floor((totalSeconds.value % 3600) / 60)
  return String(m).padStart(2, '0')
})

const seconds = computed(() => {
  const s = totalSeconds.value % 60
  return String(s).padStart(2, '0')
})

const fetchConfig = async () => {
  try {
    const res = await fetch('/api/config')
    const contentType = res.headers.get('content-type') || ''
    if (res.ok && contentType.includes('application/json')) {
      const data = await res.json()
      if (data.fecha_evento) {
        targetTimestamp.value = new Date(data.fecha_evento).getTime()
      }
      if (data.hora_servidor) {
        const serverTime = new Date(data.hora_servidor).getTime()
        serverOffset.value = serverTime - Date.now()
      }
      if (data.registros_abiertos !== undefined) {
        registrationsOpen.value = (
          data.registros_abiertos === true ||
          data.registros_abiertos === 1 ||
          data.registros_abiertos === 'true' ||
          data.registros_abiertos === '1'
        )
      }
    }
  } catch (err) {
    console.warn('No se pudo conectar a la API de configuración, usando fallback local:', err)
  } finally {
    updateCountdown()
  }
}

onMounted(() => {
  updateCountdown()
  fetchConfig()
  timerInterval = setInterval(updateCountdown, 1000)
})

onUnmounted(() => {
  if (timerInterval) {
    clearInterval(timerInterval)
  }
})

const handleRegisterClick = () => {
  if (!registrationsOpen.value) return
  playClickSound()
  emit('open-modal')
}
</script>

<template>
  <section class="relative z-10 cyber-grid pt-10 pb-24 md:py-24 overflow-hidden border-b border-white/10">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-12 items-center">
        
        <!-- Columna Izquierda: Copys, Cuenta Regresiva y Botones de Acción -->
        <div class="lg:col-span-7 space-y-6 text-center lg:text-left">
          
          <!-- Badge de Temporada Oficial -->
          <div class="inline-flex items-center gap-2.5 px-3.5 py-1.5 rounded-full bg-cyber-purple/15 border border-cyber-purple/50 text-cyber-blue font-rajdhani font-semibold text-sm tracking-widest uppercase box-glow-purple">
            <span class="w-2 h-2 rounded-full bg-cyber-pink animate-ping"></span>
            TEMPORADA OFICIAL // LA BATALLA COMIENZA AQUÍ
          </div>

          <!-- Título Principal con Gradiente y Efecto Glitch -->
          <h1 class="font-orbitron font-black text-4xl sm:text-6xl xl:text-7xl leading-tight tracking-tight uppercase">
            CONQUISTA LA ARENA,<br />
            <span class="text-transparent bg-clip-text bg-gradient-to-r from-cyber-blue via-indigo-400 to-cyber-purple text-glow-blue glitch-hover inline-block cursor-default">
              SE TU LEYENDA
            </span>
          </h1>

          <!-- Resumen del Torneo -->
          <p class="text-gray-300 text-base sm:text-lg max-w-2xl mx-auto lg:mx-0 font-light leading-relaxed">
            <strong class="text-cyber-blue font-semibold"></strong> Los rivales están listos y ha llegado tu momento de demostrar de qué estás hecho. Desafía tus límites, libera tu potencial y lucha por convertirte en una leyenda. ⚡ 
            <span class="text-white font-semibold underline decoration-cyber-purple">¡Inscríbete ahora, asegura tu lugar y prepárate! La gloria espera a quienes se atreven a jugar. ¿Tienes lo necesario para dejar tu huella?👑🔥</span>.
          </p>

          <!-- Reloj de Cuenta Regresiva en Vivo (Live Countdown Timer) -->
          <div class="pt-2">
            <p class="font-rajdhani text-xs tracking-[0.25em] text-gray-400 uppercase mb-3 flex items-center justify-center lg:justify-start gap-1.5">
              <Clock class="w-4 h-4 text-cyber-blue" />
              <span>EL CHECK-IN CIERRA EN:</span>
            </p>
            <div class="grid grid-cols-4 gap-3 max-w-md mx-auto lg:mx-0">
              <div class="gradient-border rounded-xl p-3 text-center">
                <span class="font-orbitron font-bold text-2xl sm:text-3xl text-cyber-blue block">{{ days }}</span>
                <span class="font-rajdhani text-xs text-gray-400 uppercase tracking-widest">Días</span>
              </div>
              <div class="gradient-border rounded-xl p-3 text-center">
                <span class="font-orbitron font-bold text-2xl sm:text-3xl text-white block">{{ hours }}</span>
                <span class="font-rajdhani text-xs text-gray-400 uppercase tracking-widest">Horas</span>
              </div>
              <div class="gradient-border rounded-xl p-3 text-center">
                <span class="font-orbitron font-bold text-2xl sm:text-3xl text-white block">{{ minutes }}</span>
                <span class="font-rajdhani text-xs text-gray-400 uppercase tracking-widest">Min</span>
              </div>
              <div class="gradient-border rounded-xl p-3 text-center">
                <span class="font-orbitron font-bold text-2xl sm:text-3xl text-cyber-purple text-glow-purple block">{{ seconds }}</span>
                <span class="font-rajdhani text-xs text-gray-400 uppercase tracking-widest">Seg</span>
              </div>
            </div>
          </div>

          <!-- Botones de Acción (CTA) -->
          <div class="flex flex-col sm:flex-row items-center justify-center lg:justify-start gap-4 pt-4">
            <button
              :disabled="!registrationsOpen"
              @click="handleRegisterClick"
              @mouseenter="registrationsOpen && playHoverSound()"
              type="button"
              :class="[
                'w-full sm:w-auto cyber-clip font-orbitron font-black text-sm px-8 py-4 tracking-wider uppercase transition flex items-center justify-center gap-3',
                registrationsOpen
                  ? 'bg-gradient-to-r from-cyber-blue to-cyber-purple hover:from-cyber-purple hover:to-cyber-blue text-cyber-dark transform hover:-translate-y-0.5 animate-glow cursor-pointer'
                  : 'bg-white/10 text-gray-400 border border-white/10 cursor-not-allowed opacity-60'
              ]"
            >
              <Crosshair class="w-5 h-5" :class="registrationsOpen ? 'text-cyber-dark' : 'text-gray-400'" />
              <span>{{ registrationsOpen ? 'REGISTRAR MI EQUIPO' : 'INSCRIPCIONES CERRADAS' }}</span>
            </button>

            <!-- Indicador secundario cuando las inscripciones están cerradas -->
            <div
              v-if="!registrationsOpen"
              class="inline-flex items-center gap-2 px-3.5 py-2 rounded-full bg-cyber-pink/10 border border-cyber-pink/30 text-cyber-pink font-rajdhani font-semibold text-xs tracking-widest uppercase"
            >
              <span class="w-2 h-2 rounded-full bg-cyber-pink animate-ping"></span>
              PRÓXIMAMENTE
            </div>
          </div>

        </div>

        <!-- Columna Derecha: Tarjeta Holográfica del Trofeo -->
        <div class="lg:col-span-5">
          <TrophyHologram />
        </div>

      </div>
    </div>
  </section>
</template>
