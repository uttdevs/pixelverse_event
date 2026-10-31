<script setup>
/**
 * ============================================================================
 * COMPONENTE: EditInscripcionModal.vue
 * ============================================================================
 * Edita los campos permitidos de una inscripción: juego y asistencia.
 * Los datos del participante se muestran en modo de solo lectura.
 */

import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { Loader2, Pencil, X } from 'lucide-vue-next'

const props = defineProps({
  inscripcion: { type: Object, required: true },
  juegos: { type: Array, default: () => [] },
  save: { type: Function, required: true }
})
const emit = defineEmits(['close', 'saved'])

const juegoId = ref(props.inscripcion.id_juego)
const asistio = ref(props.inscripcion.asistio === 1)
const isSaving = ref(false)
const error = ref('')
const errorDetails = ref([])

const hasChanges = computed(
  () => juegoId.value !== props.inscripcion.id_juego || (asistio.value ? 1 : 0) !== props.inscripcion.asistio
)

const handleSubmit = async () => {
  if (!hasChanges.value || isSaving.value) return
  const patch = {}
  if (juegoId.value !== props.inscripcion.id_juego) patch.juego_id = Number(juegoId.value)
  if ((asistio.value ? 1 : 0) !== props.inscripcion.asistio) patch.asistio = asistio.value ? 1 : 0

  isSaving.value = true
  error.value = ''
  errorDetails.value = []
  try {
    const result = await props.save(props.inscripcion.id_inscripcion, patch)
    emit('saved', result)
  } catch (err) {
    error.value = err.message || 'No se pudo guardar'
    errorDetails.value = err.details || []
  } finally {
    isSaving.value = false
  }
}

const handleKeydown = (event) => {
  if (event.key === 'Escape') emit('close')
}
onMounted(() => window.addEventListener('keydown', handleKeydown))
onBeforeUnmount(() => window.removeEventListener('keydown', handleKeydown))

const readOnlyFields = computed(() => [
  ['Matrícula', props.inscripcion.matricula],
  ['Carrera', props.inscripcion.carrera],
  ['Cuatrimestre / Grupo', `${props.inscripcion.cuatrimestre}° ${props.inscripcion.grupo}`],
  ['Teléfono', props.inscripcion.telefono],
  ['Correo', props.inscripcion.correo || '—'],
  ['Inscrito el', props.inscripcion.fecha_inscripcion]
])
</script>

<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm" @click.self="emit('close')">
    <div
      role="dialog"
      aria-modal="true"
      aria-labelledby="admin-edit-title"
      class="admin-pop w-full max-w-lg rounded-2xl bg-cyber-card border border-cyber-blue/50 shadow-[0_0_40px_rgba(0,240,255,0.18)] overflow-hidden"
    >
      <header class="flex items-center justify-between gap-4 px-6 py-4 border-b border-white/10 bg-gradient-to-r from-cyber-blue/10 to-cyber-purple/10">
        <div class="flex items-center gap-3 min-w-0">
          <Pencil class="w-5 h-5 text-cyber-blue shrink-0" />
          <div class="min-w-0">
            <h2 id="admin-edit-title" class="font-orbitron font-bold text-white uppercase text-sm">
              Editar inscripción #{{ inscripcion.id_inscripcion }}
            </h2>
            <p class="text-xs text-gray-400 truncate">{{ inscripcion.participante }}</p>
          </div>
        </div>
        <button id="admin-edit-close" type="button" aria-label="Cerrar" class="p-1.5 rounded text-gray-400 hover:text-white transition cursor-pointer" @click="emit('close')">
          <X class="w-5 h-5" />
        </button>
      </header>

      <form class="p-6 space-y-6" @submit.prevent="handleSubmit">
        <dl class="grid grid-cols-2 gap-x-4 gap-y-3 p-4 rounded-xl bg-cyber-dark/70 border border-white/5">
          <div v-for="[label, value] in readOnlyFields" :key="label" class="min-w-0">
            <dt class="font-rajdhani text-[10px] uppercase tracking-wider text-gray-500">{{ label }}</dt>
            <dd class="text-sm text-gray-200 truncate" :title="String(value)">{{ value }}</dd>
          </div>
        </dl>

        <div>
          <label for="admin-edit-juego" class="block font-rajdhani text-xs uppercase tracking-wider text-gray-300 mb-1.5">Juego</label>
          <select
            id="admin-edit-juego"
            v-model.number="juegoId"
            class="w-full px-4 py-3 rounded-lg bg-cyber-dark border border-white/15 focus:border-cyber-blue focus:outline-none text-sm text-white font-rajdhani"
          >
            <option v-for="juego in juegos" :key="juego.id_juego" :value="juego.id_juego">{{ juego.nombre }}</option>
          </select>
        </div>

        <div class="flex items-center justify-between p-4 rounded-xl bg-cyber-dark/70 border border-white/5">
          <div>
            <p class="font-rajdhani text-xs uppercase tracking-wider text-gray-300">Asistencia</p>
            <p class="text-xs text-gray-500">{{ asistio ? 'El participante se presentó al evento' : 'Aún no se registra su asistencia' }}</p>
          </div>
          <button
            id="admin-edit-asistio"
            type="button"
            role="switch"
            :aria-checked="asistio"
            class="relative w-12 h-6 rounded-full transition cursor-pointer"
            :class="asistio ? 'bg-emerald-500/80 shadow-[0_0_14px_rgba(16,185,129,0.6)]' : 'bg-white/15'"
            @click="asistio = !asistio"
          >
            <span class="absolute top-0.5 left-0.5 w-5 h-5 rounded-full bg-white transition-transform" :class="asistio ? 'translate-x-6' : ''" />
          </button>
        </div>

        <div v-if="error" id="admin-edit-error" role="alert" class="p-3 rounded-lg bg-cyber-pink/10 border border-cyber-pink/40 text-xs text-cyber-pink">
          <p>{{ error }}</p>
          <ul v-if="errorDetails.length" class="mt-1 list-disc list-inside">
            <li v-for="detail in errorDetails" :key="detail">{{ detail }}</li>
          </ul>
        </div>

        <div class="flex justify-end gap-3">
          <button
            id="admin-edit-cancel"
            type="button"
            class="px-4 py-2 rounded-lg text-sm font-rajdhani uppercase tracking-wider text-gray-300 bg-white/5 hover:bg-white/10 border border-white/10 transition cursor-pointer"
            @click="emit('close')"
          >
            Cancelar
          </button>
          <button
            id="admin-edit-save"
            type="submit"
            :disabled="!hasChanges || isSaving"
            class="flex items-center gap-2 px-5 py-2 rounded-lg text-sm font-orbitron font-bold uppercase tracking-wider text-cyber-dark bg-gradient-to-r from-cyber-blue to-cyber-purple hover:shadow-[0_0_20px_rgba(0,240,255,0.5)] transition disabled:opacity-40 disabled:cursor-not-allowed cursor-pointer"
          >
            <Loader2 v-if="isSaving" class="w-4 h-4 animate-spin" />
            Guardar
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<style scoped>
.admin-pop {
  animation: admin-pop 0.18s ease-out both;
}
@keyframes admin-pop {
  from { opacity: 0; transform: scale(0.95); }
  to { opacity: 1; transform: none; }
}
</style>
