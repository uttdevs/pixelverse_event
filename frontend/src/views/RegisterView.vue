<script setup>
/**
 * ============================================================================
 * VISTA: RegisterView.vue (Página Dedicada de Registro de Participantes)
 * ============================================================================
 * Vista completa para enrolamiento directo a través de la ruta `/registro`.
 * Ideal para compartir enlaces directos en Discord, Twitter o transmisiones.
 *
 * Comparte la lógica de API con `RegistrationModal.vue` mediante
 * `composables/useRegistrationForm.js`.
 * Admite preselección de juego con `/registro?juego=<id|nombre>`.
 */

import { onMounted, watch } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { ArrowLeft, Check, Crosshair, RotateCcw } from 'lucide-vue-next'
import { useAudioSfx } from '../composables/useAudioSfx'
import { useRegistrationForm } from '../composables/useRegistrationForm'

const emit = defineEmits(['registered'])
const { playClickSound, playHoverSound } = useAudioSfx()
const route = useRoute()

const {
  carreras,
  juegos,
  isLoadingCatalogs,
  catalogError,
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
  submit
} = useRegistrationForm({ initialGame: () => route.query.juego })

watch(() => route.query.juego, () => preselectInitialGame())

onMounted(() => {
  loadCatalogs()
})

const handleSubmit = async () => {
  playClickSound()
  const result = await submit()
  if (result) emit('registered', result)
}

// Conserva los datos personales para inscribirse en otro juego
// (el backend reutiliza al participante si la matrícula y los datos coinciden)
const handleNewRegistration = () => {
  playClickSound()
  clearErrors()
  registration.value = null
  form.juego_id = null
  isSubmitted.value = false
}
</script>

<template>
  <div class="relative z-10 py-16 max-w-2xl mx-auto px-4 sm:px-6 lg:px-8">
    
    <!-- Volver al Inicio -->
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

    <div class="gradient-border rounded-2xl p-6 sm:p-10 bg-cyber-card border border-cyber-blue/50 box-glow-combo">
      
      <!-- Encabezado de la Página -->
      <div class="mb-8">
        <span class="font-rajdhani text-xs uppercase tracking-[0.25em] text-cyber-blue block mb-1">
          // ENROLAMIENTO OFICIAL EN LÍNEA
        </span>
        <h1 class="font-orbitron font-black text-2xl sm:text-3xl uppercase text-white">
          INSCRIPCIÓN <span class="text-transparent bg-clip-text bg-gradient-to-r from-cyber-blue to-cyber-purple">OFICIAL</span>
        </h1>
        <p class="text-gray-400 text-xs sm:text-sm mt-1">
          Asegura tu lugar en las arenas de Pixel Verse 2026.
        </p>
      </div>

      <!-- Formulario de Registro -->
      <form
        v-if="!isSubmitted"
        @submit.prevent="handleSubmit"
        class="space-y-5"
      >
        <!-- Error al cargar catálogos -->
        <div
          v-if="catalogError"
          class="p-3 rounded-lg bg-cyber-pink/10 border border-cyber-pink/40 text-xs text-cyber-pink flex items-center justify-between gap-3"
          role="alert"
        >
          <span>{{ catalogError }}</span>
          <button
            id="page-reg-retry"
            type="button"
            @click="loadCatalogs"
            class="px-3 py-1 rounded bg-white/10 hover:bg-white/20 text-white uppercase font-orbitron text-[10px] cursor-pointer"
          >
            Reintentar
          </button>
        </div>

        <!-- Nombre completo -->
        <div>
          <label for="page-reg-nombre" class="block font-rajdhani text-xs uppercase tracking-wider text-gray-300 mb-1.5">
            NOMBRE COMPLETO *
          </label>
          <input
            id="page-reg-nombre"
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
            <label for="page-reg-juego" class="block font-rajdhani text-xs uppercase tracking-wider text-gray-300 mb-1.5">
              ARENA DE JUEGO *
            </label>
            <select
              id="page-reg-juego"
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
            <label for="page-reg-matricula" class="block font-rajdhani text-xs uppercase tracking-wider text-gray-300 mb-1.5">
              Matrícula *
            </label>
            <input
              id="page-reg-matricula"
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
            <label for="page-reg-correo" class="block font-rajdhani text-xs uppercase tracking-wider text-gray-300 mb-1.5">
              Correo electrónico
            </label>
            <input
              id="page-reg-correo"
              v-model="form.correo"
              type="email"
              maxlength="100"
              autocomplete="email"
              placeholder="Ej: nombre@alumno.edu.mx"
              class="w-full px-4 py-3 rounded-lg bg-cyber-dark border border-white/15 focus:border-cyber-purple focus:outline-none text-sm text-white"
            />
          </div>

          <!-- Teléfono -->
          <div>
            <label for="page-reg-telefono" class="block font-rajdhani text-xs uppercase tracking-wider text-gray-300 mb-1.5">
              Teléfono *
            </label>
            <input
              id="page-reg-telefono"
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

          <!-- Cuatrimestre -->
          <div>
            <label for="page-reg-cuatrimestre" class="block font-rajdhani text-xs uppercase tracking-wider text-gray-300 mb-1.5">
              Cuatrimestre *
            </label>
            <input
              id="page-reg-cuatrimestre"
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
            <label for="page-reg-grupo" class="block font-rajdhani text-xs uppercase tracking-wider text-gray-300 mb-1.5">
              Grupo *
            </label>
            <input
              id="page-reg-grupo"
              v-model="form.grupo"
              type="text"
              required
              maxlength="5"
              placeholder="Ej: B"
              class="w-full px-4 py-3 rounded-lg bg-cyber-dark border border-white/15 focus:border-cyber-purple focus:outline-none text-sm text-white uppercase"
            />
          </div>

          <!-- Carrera -->
          <div class="sm:col-span-2">
            <label for="page-reg-carrera" class="block font-rajdhani text-xs uppercase tracking-wider text-gray-300 mb-1.5">
              CARRERA *
            </label>
            <select
              id="page-reg-carrera"
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
        <div class="pt-4">
          <button
            id="page-reg-submit"
            type="submit"
            :disabled="isSubmitting || isLoadingCatalogs || !!catalogError"
            class="w-full cyber-clip py-4 bg-gradient-to-r from-cyber-blue to-cyber-purple font-orbitron font-black text-sm text-cyber-dark uppercase tracking-wider hover:opacity-95 transition cursor-pointer flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <Crosshair class="w-4 h-4 text-cyber-dark" :class="{ 'animate-spin': isSubmitting }" />
            <span>{{ isSubmitting ? 'ENVIANDO REGISTRO...' : 'CONFIRMAR LUGAR EN LA ARENA' }}</span>
          </button>
        </div>
      </form>

      <!-- Mensaje de Confirmación -->
      <div v-else class="text-center py-8 space-y-5" aria-live="polite">
        <div class="w-16 h-16 mx-auto rounded-full bg-emerald-500/20 border border-emerald-400 flex items-center justify-center text-emerald-400">
          <Check class="w-8 h-8" />
        </div>

        <h2 class="font-orbitron font-bold text-2xl text-white uppercase">
          ¡{{ registration?.message || 'Registro completado con éxito' }}!
        </h2>
        
        <p class="text-sm text-gray-300 max-w-md mx-auto leading-relaxed font-rajdhani">
          <strong class="text-white">{{ form.nombre_completo }}</strong> quedó oficialmente inscrito en la arena
          <strong class="text-cyber-blue">{{ selectedGameName }}</strong>.
        </p>

        <dl class="p-4 rounded-xl bg-cyber-dark border border-cyber-blue/40 font-mono text-sm text-cyber-blue max-w-sm mx-auto grid grid-cols-2 gap-y-1.5 text-left">
          <dt>FOLIO:</dt>
          <dd class="font-bold text-white text-right">#{{ registration?.id_inscripcion }}</dd>
          <dt>MATRÍCULA:</dt>
          <dd class="font-bold text-white text-right">{{ registration?.matricula }}</dd>
          <dt>ID PARTICIPANTE:</dt>
          <dd class="font-bold text-white text-right">{{ registration?.id_participante }}</dd>
        </dl>

        <div class="pt-4 flex flex-col sm:flex-row items-center justify-center gap-3">
          <button
            id="page-reg-new"
            type="button"
            @click="handleNewRegistration"
            @mouseenter="playHoverSound"
            class="inline-flex items-center gap-2 px-8 py-3 rounded-lg border border-cyber-blue/40 hover:bg-cyber-blue/10 font-orbitron text-xs uppercase text-cyber-blue transition cursor-pointer"
          >
            <RotateCcw class="w-4 h-4" />
            INSCRIBIR EN OTRO JUEGO
          </button>
          <RouterLink
            to="/"
            class="inline-block px-8 py-3 rounded-lg bg-white/10 hover:bg-white/20 font-orbitron text-xs uppercase text-white transition"
          >
            VOLVER A LA ARENA PRINCIPAL
          </RouterLink>
        </div>
      </div>

    </div>

  </div>
</template>
