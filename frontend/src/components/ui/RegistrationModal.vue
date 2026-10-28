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
 *
 * INTEGRACIÓN CON LA API:
 * - La lógica de catálogos y envío vive en `composables/useRegistrationForm.js`
 *   (compartida con `views/RegisterView.vue`).
 */

import { ref, watch, onMounted, onUnmounted, nextTick } from 'vue'
import { X, Check } from 'lucide-vue-next'
import { useAudioSfx } from '../../composables/useAudioSfx'
import { useRegistrationForm } from '../../composables/useRegistrationForm'

const props = defineProps({
  isOpen: {
    type: Boolean,
    default: false
  },
  initialGame: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['close', 'registered'])
const { playClickSound, playHoverSound } = useAudioSfx()

// Referencias a elementos del DOM para control de accesibilidad y foco
const modalRef = ref(null)
const firstInputRef = ref(null)
let previousActiveElement = null

const {
  carreras,
  juegos,
  isLoadingCatalogs,
  catalogError,
  hasCatalogs,
  loadCatalogs,
  form,
  isSubmitted,
  isSubmitting,
  submitError,
  submitErrorDetails,
  registration,
  selectedGameName,
  preselectInitialGame,
  clearErrors,
  resetForm,
  submit
} = useRegistrationForm({ initialGame: () => props.initialGame })

watch(() => props.initialGame, () => preselectInitialGame())

// Control de apertura/cierre, scroll lock y gestión de foco
watch(() => props.isOpen, async (isOpen) => {
  if (isOpen) {
    if (isSubmitted.value) resetForm()
    clearErrors()
    previousActiveElement = document.activeElement
    // Bloqueamos el scroll de la página de fondo
    document.body.style.overflow = 'hidden'

    if (!hasCatalogs.value) {
      loadCatalogs()
    } else {
      preselectInitialGame()
    }

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
      'button:not([disabled]), [href], input, select:not([disabled]), textarea, [tabindex]:not([tabindex="-1"])'
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
  loadCatalogs()
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeyDown)
  document.body.style.overflow = ''
})

const handleSubmit = async () => {
  playClickSound()
  const result = await submit()
  if (result) emit('registered', result)
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
      class="gradient-border rounded-2xl max-w-lg w-full max-h-[95vh] overflow-y-auto p-6 sm:p-8 bg-cyber-card border border-cyber-blue/50 box-glow-combo relative"
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
        <!-- Error al cargar catálogos -->
        <div
          v-if="catalogError"
          class="p-3 rounded-lg bg-cyber-pink/10 border border-cyber-pink/40 text-xs text-cyber-pink flex items-center justify-between gap-3"
          role="alert"
        >
          <span>{{ catalogError }}</span>
          <button
            type="button"
            @click="loadCatalogs"
            class="px-3 py-1 rounded bg-white/10 hover:bg-white/20 text-white uppercase font-orbitron text-[10px] cursor-pointer"
          >
            Reintentar
          </button>
        </div>

        <!-- Nombre completo -->
        <div>
          <label for="reg-nombre" class="block font-rajdhani text-xs uppercase tracking-wider text-gray-300 mb-1">
            NOMBRE COMPLETO *
          </label>
          <input
            id="reg-nombre"
            ref="firstInputRef"
            v-model="form.nombre_completo"
            type="text"
            required
            maxlength="60"
            autocomplete="name"
            placeholder="Ej: FRIDA MENDEZ MONTIEL"
            class="w-full px-4 py-3 rounded-lg bg-cyber-dark border border-white/15 focus:border-cyber-blue focus:outline-none text-sm text-white font-orbitron"
          />
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <!-- Juego -->
          <div>
            <label for="reg-juego" class="block font-rajdhani text-xs uppercase tracking-wider text-gray-300 mb-1">
              ARENA DE JUEGO *
            </label>
            <select
              id="reg-juego"
              v-model.number="form.juego_id"
              required
              :disabled="isLoadingCatalogs || !juegos.length"
              class="w-full px-4 py-3 rounded-lg bg-cyber-dark border border-white/15 focus:border-cyber-blue focus:outline-none text-sm text-white font-rajdhani disabled:opacity-50"
            >
              <option :value="null" disabled>
                {{ isLoadingCatalogs ? 'Cargando...' : 'Selecciona un juego' }}
              </option>
              <option v-for="juego in juegos" :key="juego.id_juego" :value="juego.id_juego">
                {{ juego.nombre }}
              </option>
            </select>
          </div>

          <!-- Matrícula -->
          <div>
            <label for="reg-matricula" class="block font-rajdhani text-xs uppercase tracking-wider text-gray-300 mb-1">
              Matrícula *
            </label>
            <input
              id="reg-matricula"
              v-model="form.matricula"
              type="text"
              required
              maxlength="20"
              placeholder="Ej: 25378039"
              class="w-full px-4 py-3 rounded-lg bg-cyber-dark border border-white/15 focus:border-cyber-purple focus:outline-none text-sm text-white"
            />
          </div>

          <!-- Correo (opcional) -->
          <div>
            <label for="reg-correo" class="block font-rajdhani text-xs uppercase tracking-wider text-gray-300 mb-1">
              Correo electrónico
            </label>
            <input
              id="reg-correo"
              v-model="form.correo"
              type="email"
              maxlength="100"
              autocomplete="email"
              placeholder="Ej: nombre@alumno.edu.mx"
              class="w-full px-4 py-3 rounded-lg bg-cyber-dark border border-white/15 focus:border-cyber-purple focus:outline-none text-sm text-white"
            />
          </div>

          <!-- Cuatrimestre -->
          <div>
            <label for="reg-cuatrimestre" class="block font-rajdhani text-xs uppercase tracking-wider text-gray-300 mb-1">
              Cuatrimestre *
            </label>
            <input
              id="reg-cuatrimestre"
              v-model.number="form.cuatrimestre"
              type="number"
              required
              min="1"
              max="12"
              step="1"
              placeholder="Ej: 4"
              class="w-full px-4 py-3 rounded-lg bg-cyber-dark border border-white/15 focus:border-cyber-purple focus:outline-none text-sm text-white"
            />
          </div>

          <!-- Grupo -->
          <div>
            <label for="reg-grupo" class="block font-rajdhani text-xs uppercase tracking-wider text-gray-300 mb-1">
              Grupo *
            </label>
            <input
              id="reg-grupo"
              v-model="form.grupo"
              type="text"
              required
              maxlength="5"
              placeholder="Ej: B"
              class="w-full px-4 py-3 rounded-lg bg-cyber-dark border border-white/15 focus:border-cyber-purple focus:outline-none text-sm text-white uppercase"
            />
          </div>

          <!-- Teléfono -->
          <div>
            <label for="reg-telefono" class="block font-rajdhani text-xs uppercase tracking-wider text-gray-300 mb-1">
              Teléfono *
            </label>
            <input
              id="reg-telefono"
              v-model="form.telefono"
              type="tel"
              required
              maxlength="15"
              pattern="[0-9+(). \-]{10,15}"
              title="Entre 10 y 15 dígitos. Se permiten +, espacios, paréntesis, puntos o guiones."
              autocomplete="tel"
              placeholder="Ej: 2361829323"
              class="w-full px-4 py-3 rounded-lg bg-cyber-dark border border-white/15 focus:border-cyber-purple focus:outline-none text-sm text-white"
            />
          </div>

          <!-- Carrera -->
          <div class="sm:col-span-2">
            <label for="reg-carrera" class="block font-rajdhani text-xs uppercase tracking-wider text-gray-300 mb-1">
              CARRERA *
            </label>
            <select
              id="reg-carrera"
              v-model.number="form.id_carrera"
              required
              :disabled="isLoadingCatalogs || !carreras.length"
              class="w-full px-4 py-3 rounded-lg bg-cyber-dark border border-white/15 focus:border-cyber-blue focus:outline-none text-sm text-white font-rajdhani disabled:opacity-50"
            >
              <option :value="null" disabled>
                {{ isLoadingCatalogs ? 'Cargando...' : 'Selecciona tu carrera' }}
              </option>
              <option v-for="carrera in carreras" :key="carrera.id_carrera" :value="carrera.id_carrera">
                {{ carrera.nombre }}
              </option>
            </select>
          </div>
        </div>

        <!-- Errores devueltos por la API -->
        <div
          v-if="submitError"
          class="p-3 rounded-lg bg-cyber-pink/10 border border-cyber-pink/40 text-xs text-cyber-pink space-y-1"
          role="alert"
          aria-live="assertive"
        >
          <p class="font-bold">{{ submitError }}</p>
          <ul v-if="submitErrorDetails.length" class="list-disc list-inside space-y-0.5">
            <li v-for="(detail, index) in submitErrorDetails" :key="index">{{ detail }}</li>
          </ul>
        </div>

        <!-- Botón de Envío -->
        <div class="pt-3">
          <button
            id="reg-submit"
            type="submit"
            :disabled="isSubmitting || isLoadingCatalogs || !!catalogError"
            class="w-full cyber-clip py-4 bg-gradient-to-r from-cyber-blue to-cyber-purple font-orbitron font-black text-sm text-cyber-dark uppercase tracking-wider hover:opacity-95 transition cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ isSubmitting ? 'ENVIANDO REGISTRO...' : 'CONFIRMAR LUGAR EN LA ARENA' }}
          </button>
        </div>
      </form>

      <!-- VISTA 2: Mensaje de Registro Exitoso -->
      <div v-else class="text-center py-6 space-y-4" aria-live="polite">
        <div class="w-16 h-16 mx-auto rounded-full bg-emerald-500/20 border border-emerald-400 flex items-center justify-center text-emerald-400 text-3xl">
          <Check class="w-8 h-8" />
        </div>

        <h4 class="font-orbitron font-bold text-xl text-white uppercase">
          ¡{{ registration?.message || 'Registro completado con éxito' }}!
        </h4>
        
        <p class="text-xs text-gray-300 max-w-sm mx-auto leading-relaxed">
          <strong class="text-white">{{ form.nombre_completo }}</strong> quedó inscrito en la arena
          <strong class="text-cyber-blue">{{ selectedGameName }}</strong>.
        </p>

        <dl class="p-3 rounded-lg bg-cyber-dark border border-cyber-blue/30 font-mono text-xs text-cyber-blue grid grid-cols-2 gap-y-1 text-left">
          <dt>FOLIO DE INSCRIPCIÓN:</dt>
          <dd class="font-bold text-white text-right">#{{ registration?.id_inscripcion }}</dd>
          <dt>MATRÍCULA:</dt>
          <dd class="font-bold text-white text-right">{{ registration?.matricula }}</dd>
          <dt>ID PARTICIPANTE:</dt>
          <dd class="font-bold text-white text-right">{{ registration?.id_participante }}</dd>
        </dl>

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
