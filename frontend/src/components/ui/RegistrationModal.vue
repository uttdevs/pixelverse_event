<script setup>
/**
 * ============================================================================
 * COMPONENTE: RegistrationModal.vue (Modal Accesible A11y y Formulario Reactivo)
 * ============================================================================
 * MEJORAS DE ACCESIBILIDAD (A11y) IMPLEMENTADAS:
 * 1. Cierre con tecla Escape (`handleKeyDown`).
 * 2. Cierre por clic exterior (Light Dismiss) con `@click.self="handleClose"`.
 * 3. Bloqueo de desplazamiento de fondo (`prevent-page-scroll`) mientras esté abierto.
 * 4. Atrapamiento de foco accesible (Focus Trap) para navegación con teclado (Tab/Shift+Tab).
 * 5. Autoenfoque en el primer campo al abrirse el diálogo.
 * 6. Iconos vectoriales de Lucide: X, Check.
 */

import { reactive, ref, watch, onMounted, onUnmounted, nextTick } from 'vue'
import { X, Check } from 'lucide-vue-next'
import { useAudioSfx } from '../../composables/useAudioSfx'

const props = defineProps({
  isOpen: {
    type: Boolean,
    default: false
  },
  initialGame: {
    type: String,
    default: 'Valorant: Protocol'
  }
})

const emit = defineEmits(['close', 'registered'])
const { playClickSound, playHoverSound } = useAudioSfx()

// Referencias a elementos del DOM para control de accesibilidad y foco
const modalRef = ref(null)
const firstInputRef = ref(null)
let previousActiveElement = null

const form = reactive({
  squadName: '',
  game: props.initialGame,
  discord: '',
  tier: 'Diamante+'
})

const isSubmitted = ref(false)
const generatedToken = ref('')

watch(() => props.initialGame, (newGame) => {
  if (newGame) {
    form.game = newGame
  }
})

// Control de apertura/cierre, scroll lock y gestión de foco
watch(() => props.isOpen, async (isOpen) => {
  if (isOpen) {
    isSubmitted.value = false
    previousActiveElement = document.activeElement
    // Bloqueamos el scroll de la página de fondo
    document.body.style.overflow = 'hidden'

    await nextTick()
    if (firstInputRef.value) {
      firstInputRef.value.focus()
    }
  } else {
    // Restauramos el scroll de la página
    document.body.style.overflow = ''
    if (previousActiveElement && typeof previousActiveElement.focus === 'function') {
      previousActiveElement.focus()
    }
  }
})

const handleClose = () => {
  playClickSound()
  emit('close')
}

// Cierre al presionar la tecla Escape y trampa de foco (Focus Trap)
const handleKeyDown = (event) => {
  if (!props.isOpen) return

  // 1. Cierre con Escape
  if (event.key === 'Escape') {
    handleClose()
    return
  }

  // 2. Focus Trap con tecla Tab
  if (event.key === 'Tab' && modalRef.value) {
    const focusableElements = modalRef.value.querySelectorAll(
      'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
    )
    if (!focusableElements.length) return

    const firstElement = focusableElements[0]
    const lastElement = focusableElements[focusableElements.length - 1]

    if (event.shiftKey) {
      // Shift + Tab hacia atrás
      if (document.activeElement === firstElement) {
        event.preventDefault()
        lastElement.focus()
      }
    } else {
      // Tab hacia adelante
      if (document.activeElement === lastElement) {
        event.preventDefault()
        firstElement.focus()
      }
    }
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleKeyDown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeyDown)
  document.body.style.overflow = ''
})

const handleSubmit = () => {
  playClickSound()
  generatedToken.value = 'PV-2026-' + Math.floor(1000 + Math.random() * 9000)
  isSubmitted.value = true

  emit('registered', {
    ...form,
    token: generatedToken.value
  })
}
</script>

<template>
  <!-- Contenedor del Modal con Backdrop Blur y Light Dismiss (@click.self) -->
  <div
    v-if="isOpen"
    @click.self="handleClose"
    class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-md animate-[fadeIn_0.2s_ease-out]"
    role="dialog"
    aria-modal="true"
    aria-labelledby="modal-title"
  >
    <div
      ref="modalRef"
      class="gradient-border rounded-2xl max-w-lg w-full p-6 sm:p-8 bg-cyber-card border border-cyber-blue/50 box-glow-combo relative"
    >
      <!-- Botón de Cerrar (X) -->
      <button
        @click="handleClose"
        @mouseenter="playHoverSound"
        type="button"
        class="absolute top-4 right-4 w-9 h-9 rounded-lg bg-white/5 hover:bg-cyber-pink/20 text-gray-400 hover:text-cyber-pink flex items-center justify-center transition cursor-pointer"
        aria-label="Cerrar ventana de registro"
      >
        <X class="w-5 h-5" />
      </button>

      <!-- Encabezado del Modal -->
      <div class="mb-6">
        <span class="font-rajdhani text-xs uppercase tracking-[0.25em] text-cyber-blue">
          // REGISTRO 
        </span>
        <h3 id="modal-title" class="font-orbitron font-black text-2xl text-white uppercase mt-1">
          INSCRIPCIÓN <span class="text-cyber-purple">OFICIAL</span>
        </h3>
      </div>

      <!-- VISTA 1: Formulario de Inscripción -->
      <form
        v-if="!isSubmitted"
        @submit.prevent="handleSubmit"
        class="space-y-4"
      >
        <!-- Nombre del Squad -->
        <div>
          <label class="block font-rajdhani text-xs uppercase tracking-wider text-gray-300 mb-1">
            NOMBRE COMPLETO
          </label>
          <input
            ref="firstInputRef"
            v-model="form.squadName"
            type="text"
            required
            placeholder="Ej: FRIDA MENDEZ MONTIEL"
            class="w-full px-4 py-3 rounded-lg bg-cyber-dark border border-white/15 focus:border-cyber-blue focus:outline-none text-sm text-white font-orbitron"
          />
        </div>

        <!-- Selección de Arena y Discord -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="block font-rajdhani text-xs uppercase tracking-wider text-gray-300 mb-1">
              ARENA DE JUEGO *
            </label>
            <select
              v-model="form.game"
              class="w-full px-4 py-3 rounded-lg bg-cyber-dark border border-white/15 focus:border-cyber-blue focus:outline-none text-sm text-white font-rajdhani"
            >
              <option value="Valorant: Protocol">FIFA 2026</option>
              <option value="League of Pixels">THE KING OF FIGHTERS</option>
              <option value="CS2: Cyber Strike">STREET FIGHTER</option>
            
            </select>
          </div>

          <div>
            <label class="block font-rajdhani text-xs uppercase tracking-wider text-gray-300 mb-1">
              Matrícula 
            </label>
            <input
              v-model="form.discord"
              type="text"
              required
              placeholder="Ej: 25378039"
              class="w-full px-4 py-3 rounded-lg bg-cyber-dark border border-white/15 focus:border-cyber-purple focus:outline-none text-sm text-white"
            />
          </div>
          <div>
            <label class="block font-rajdhani text-xs uppercase tracking-wider text-gray-300 mb-1">
              Correo electrónico 
            </label>
            <input
              v-model="form.discord"
              type="text"
              required
              placeholder="Ej: Frid@gmail.com"
              class="w-full px-4 py-3 rounded-lg bg-cyber-dark border border-white/15 focus:border-cyber-purple focus:outline-none text-sm text-white"
            />
          </div>
          
          <div>
            <label class="block font-rajdhani text-xs uppercase tracking-wider text-gray-300 mb-1">
              Cuatrimestre
            </label>
            <input
              v-model="form.discord"
              type="text"
              required
              placeholder="Ej: 4"
              class="w-full px-4 py-3 rounded-lg bg-cyber-dark border border-white/15 focus:border-cyber-purple focus:outline-none text-sm text-white"
            />
          </div>
          <div>
            <label class="block font-rajdhani text-xs uppercase tracking-wider text-gray-300 mb-1">
              Grupo
            </label>
            <input
              v-model="form.discord"
              type="text"
              required
              placeholder="Ej: B"
              class="w-full px-4 py-3 rounded-lg bg-cyber-dark border border-white/15 focus:border-cyber-purple focus:outline-none text-sm text-white"
            />
          </div>
           <div>
            <label class="block font-rajdhani text-xs uppercase tracking-wider text-gray-300 mb-1">
              CARRERA 
            </label>
            <select
              v-model="form.game"
              class="w-full px-4 py-3 rounded-lg bg-cyber-dark border border-white/15 focus:border-cyber-blue focus:outline-none text-sm text-white font-rajdhani"
            >
              <option value="Valorant: Protocol">Tecnologias de la informacion </option>
              <option value="League of Pixels">THE KING OF FIGHTERS</option>
              <option value="CS2: Cyber Strike">STREET FIGHTER</option>
            
            </select>
          </div>
          <div>
            <label class="block font-rajdhani text-xs uppercase tracking-wider text-gray-300 mb-1">
              Teléfono
            </label>
            <input
              v-model="form.discord"
              type="text"
              required
              placeholder="Ej: 236 182 9323"
              class="w-full px-4 py-3 rounded-lg bg-cyber-dark border border-white/15 focus:border-cyber-purple focus:outline-none text-sm text-white"
            />
          </div>
        </div>

        

        <!-- Botón de Envío -->
        <div class="pt-3">
          <button
            type="submit"
            class="w-full cyber-clip py-4 bg-gradient-to-r from-cyber-blue to-cyber-purple font-orbitron font-black text-sm text-cyber-dark uppercase tracking-wider hover:opacity-95 transition cursor-pointer"
          >
            CONFIRMAR LUGAR EN LA ARENA
          </button>
        </div>
      </form>

      <!-- VISTA 2: Mensaje de Registro Exitoso -->
      <div v-else class="text-center py-6 space-y-4">
        <div class="w-16 h-16 mx-auto rounded-full bg-emerald-500/20 border border-emerald-400 flex items-center justify-center text-emerald-400 text-3xl">
          <Check class="w-8 h-8" />
        </div>

        <h4 class="font-orbitron font-bold text-xl text-white">¡SQUAD REGISTRADO CON ÉXITO!</h4>
        
        <p class="text-xs text-gray-300 max-w-sm mx-auto leading-relaxed">
          El escuadrón <strong class="text-white">"{{ form.squadName }}"</strong> ha quedado pre-registrado en la arena
          <strong class="text-cyber-blue">{{ form.game }}</strong>. Enviamos confirmación al capitán ({{ form.discord }}).
        </p>

        <div class="p-3 rounded-lg bg-cyber-dark border border-cyber-blue/30 font-mono text-xs text-cyber-blue">
          TOKEN DE ACCESO: <span class="font-bold text-white">{{ generatedToken }}</span>
        </div>

        <button
          @click="handleClose"
          type="button"
          class="px-6 py-2.5 rounded-lg bg-white/10 hover:bg-white/20 font-orbitron text-xs uppercase cursor-pointer"
        >
          VOLVER AL LOBBY
        </button>
      </div>

    </div>
  </div>
</template>
