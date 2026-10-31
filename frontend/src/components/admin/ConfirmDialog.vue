<script setup>
/**
 * ============================================================================
 * COMPONENTE: ConfirmDialog.vue
 * ============================================================================
 * Diálogo genérico de confirmación para acciones destructivas.
 */

import { onBeforeUnmount, onMounted } from 'vue'
import { AlertTriangle, Loader2 } from 'lucide-vue-next'

defineProps({
  title: { type: String, required: true },
  message: { type: String, default: '' },
  warning: { type: String, default: '' },
  confirmLabel: { type: String, default: 'Eliminar' },
  isProcessing: { type: Boolean, default: false }
})
const emit = defineEmits(['confirm', 'cancel'])

const handleKeydown = (event) => {
  if (event.key === 'Escape') emit('cancel')
}
onMounted(() => window.addEventListener('keydown', handleKeydown))
onBeforeUnmount(() => window.removeEventListener('keydown', handleKeydown))
</script>

<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm" @click.self="emit('cancel')">
    <div
      role="alertdialog"
      aria-modal="true"
      aria-labelledby="admin-confirm-title"
      class="admin-pop w-full max-w-md rounded-2xl p-6 bg-cyber-card border border-cyber-pink/50 shadow-[0_0_40px_rgba(255,0,127,0.2)]"
    >
      <div class="flex items-start gap-4">
        <div class="shrink-0 w-11 h-11 flex items-center justify-center rounded-xl bg-cyber-pink/15 border border-cyber-pink/40">
          <AlertTriangle class="w-5 h-5 text-cyber-pink" />
        </div>
        <div class="min-w-0">
          <h2 id="admin-confirm-title" class="font-orbitron font-bold text-lg text-white uppercase">{{ title }}</h2>
          <p v-if="message" class="mt-1 text-sm text-gray-300">{{ message }}</p>
          <p v-if="warning" class="mt-3 p-2.5 rounded-lg bg-amber-400/10 border border-amber-400/40 text-xs text-amber-300">
            {{ warning }}
          </p>
        </div>
      </div>
      <div class="mt-6 flex justify-end gap-3">
        <button
          id="admin-confirm-cancel"
          type="button"
          :disabled="isProcessing"
          class="px-4 py-2 rounded-lg text-sm font-rajdhani uppercase tracking-wider text-gray-300 bg-white/5 hover:bg-white/10 border border-white/10 transition cursor-pointer"
          @click="emit('cancel')"
        >
          Cancelar
        </button>
        <button
          id="admin-confirm-accept"
          type="button"
          :disabled="isProcessing"
          class="flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-orbitron font-bold uppercase tracking-wider text-white bg-cyber-pink/80 hover:bg-cyber-pink hover:shadow-[0_0_20px_rgba(255,0,127,0.5)] transition disabled:opacity-60 cursor-pointer"
          @click="emit('confirm')"
        >
          <Loader2 v-if="isProcessing" class="w-4 h-4 animate-spin" />
          {{ confirmLabel }}
        </button>
      </div>
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
