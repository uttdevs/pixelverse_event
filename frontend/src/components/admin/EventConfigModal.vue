<script setup>
import { onBeforeUnmount, onMounted, ref } from 'vue'
import { CalendarDays, Check, Loader2, Settings, X } from 'lucide-vue-next'

const props = defineProps({
  config: { type: Object, required: true },
  isSaving: { type: Boolean, default: false }
})

const emit = defineEmits(['close', 'save'])
const form = ref({
  registros_abiertos: props.config.registros_abiertos,
  fecha_evento: props.config.fecha_evento
})

const handleKeydown = (event) => {
  if (event.key === 'Escape' && !props.isSaving) emit('close')
}

onMounted(() => {
  window.addEventListener('keydown', handleKeydown)
  document.body.style.overflow = 'hidden'
})
onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleKeydown)
  document.body.style.overflow = ''
})
</script>

<template>
  <Teleport to="body">
    <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm" @click.self="!isSaving && emit('close')">
      <div
        role="dialog"
        aria-modal="true"
        aria-labelledby="admin-settings-title"
        class="admin-pop w-full max-w-lg max-h-[calc(100vh-2rem)] overflow-y-auto rounded-2xl p-6 bg-cyber-card border border-cyber-blue/50 shadow-[0_0_40px_rgba(0,212,255,0.2)]"
      >
      <div class="flex items-start justify-between gap-4">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 flex items-center justify-center rounded-xl bg-cyber-blue/15 border border-cyber-blue/40">
            <Settings class="w-5 h-5 text-cyber-blue" />
          </div>
          <div>
            <h2 id="admin-settings-title" class="font-orbitron font-bold text-lg text-white uppercase">Ajustes</h2>
            <p class="text-xs text-gray-400">Configuración del evento</p>
          </div>
        </div>
        <button type="button" title="Cerrar" :disabled="isSaving" class="p-1 text-gray-400 hover:text-white transition cursor-pointer disabled:opacity-40" @click="emit('close')">
          <X class="w-5 h-5" />
        </button>
      </div>

      <div class="mt-7 space-y-5">
        <label class="flex items-center justify-between gap-4 p-4 rounded-xl bg-white/5 border border-white/10 cursor-pointer">
          <span>
            <span class="block font-rajdhani text-sm uppercase tracking-wider text-white">Inscripciones</span>
            <span class="block mt-1 text-xs text-gray-400">Permitir nuevos registros al público</span>
          </span>
          <span class="relative shrink-0">
            <input id="admin-settings-registrations" v-model="form.registros_abiertos" type="checkbox" class="peer sr-only">
            <span class="block w-11 h-6 rounded-full bg-gray-700 peer-checked:bg-cyber-blue transition" />
            <span class="absolute left-1 top-1 w-4 h-4 rounded-full bg-white transition peer-checked:translate-x-5" />
          </span>
        </label>

        <div>
          <label for="admin-settings-date" class="flex items-center gap-2 font-rajdhani text-sm uppercase tracking-wider text-white mb-2">
            <CalendarDays class="w-4 h-4 text-cyber-purple" />
            Fecha y hora del evento
          </label>
          <input
            id="admin-settings-date"
            v-model="form.fecha_evento"
            type="datetime-local"
            required
            class="w-full px-3 py-2.5 rounded-lg bg-cyber-dark border border-white/15 focus:border-cyber-blue focus:outline-none text-sm text-white [color-scheme:dark]"
          >
        </div>
      </div>

      <div class="mt-7 flex justify-end gap-3">
        <button type="button" :disabled="isSaving" class="px-4 py-2 rounded-lg text-sm font-rajdhani uppercase tracking-wider text-gray-300 bg-white/5 hover:bg-white/10 border border-white/10 transition cursor-pointer disabled:opacity-40" @click="emit('close')">
          Cancelar
        </button>
        <button type="button" :disabled="isSaving || !form.fecha_evento" class="flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-orbitron font-bold uppercase tracking-wider text-cyber-dark bg-cyber-blue hover:shadow-[0_0_20px_rgba(0,212,255,0.5)] transition cursor-pointer disabled:opacity-50" @click="emit('save', { ...form })">
          <Loader2 v-if="isSaving" class="w-4 h-4 animate-spin" />
          <Check v-else class="w-4 h-4" />
          Guardar
        </button>
      </div>
      </div>
    </div>
  </Teleport>
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
