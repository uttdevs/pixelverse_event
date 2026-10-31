<script setup>
/**
 * ============================================================================
 * COMPONENTE: InscripcionesTable.vue
 * ============================================================================
 * Tabla de inscripciones del panel administrativo. Emite eventos para
 * editar, eliminar y alternar la asistencia; no realiza llamadas a la API.
 */

import { CheckCircle2, Circle, Pencil, Trash2 } from 'lucide-vue-next'

defineProps({
  inscripciones: { type: Array, required: true },
  togglingId: { type: Number, default: null }
})
const emit = defineEmits(['edit', 'delete', 'toggle-asistio'])

// Paleta rotativa de badges para los juegos (estable por id_juego).
const badgePalette = [
  'text-cyber-blue border-cyber-blue/40 bg-cyber-blue/10',
  'text-fuchsia-300 border-fuchsia-400/40 bg-fuchsia-400/10',
  'text-amber-300 border-amber-400/40 bg-amber-400/10',
  'text-emerald-300 border-emerald-400/40 bg-emerald-400/10',
  'text-violet-300 border-violet-400/40 bg-violet-400/10',
  'text-rose-300 border-rose-400/40 bg-rose-400/10'
]
const badgeClass = (idJuego) => badgePalette[(Number(idJuego) || 0) % badgePalette.length]

const formatDate = (value) => {
  if (!value) return '—'
  const date = new Date(String(value).replace(' ', 'T'))
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleString('es-MX', { dateStyle: 'medium', timeStyle: 'short' })
}
</script>

<template>
  <div class="admin-table-wrap rounded-2xl border border-white/10 bg-cyber-card/70 backdrop-blur-xl overflow-hidden">
    <div class="overflow-x-auto max-h-[65vh]">
      <table class="w-full text-sm">
        <thead class="sticky top-0 z-10 bg-cyber-card/95 backdrop-blur">
          <tr class="text-left font-rajdhani text-[11px] uppercase tracking-[0.15em] text-gray-400 border-b border-white/10">
            <th class="px-4 py-3 font-semibold">#</th>
            <th class="px-4 py-3 font-semibold">Participante</th>
            <th class="px-4 py-3 font-semibold">Juego</th>
            <th class="px-4 py-3 font-semibold">Carrera</th>
            <th class="px-4 py-3 font-semibold whitespace-nowrap">Cuatri / Grupo</th>
            <th class="px-4 py-3 font-semibold">Contacto</th>
            <th class="px-4 py-3 font-semibold">Fecha</th>
            <th class="px-4 py-3 font-semibold text-center">Asistió</th>
            <th class="px-4 py-3 font-semibold text-right">Acciones</th>
          </tr>
        </thead>
        <TransitionGroup tag="tbody" name="admin-row">
          <tr
            v-for="item in inscripciones"
            :id="`admin-row-${item.id_inscripcion}`"
            :key="item.id_inscripcion"
            class="border-b border-white/5 hover:bg-cyber-blue/[0.04] transition-colors"
          >
            <td class="px-4 py-3 font-orbitron text-xs text-gray-500">{{ item.id_inscripcion }}</td>
            <td class="px-4 py-3 min-w-[200px]">
              <p class="font-semibold text-white">{{ item.participante }}</p>
              <p class="font-orbitron text-[11px] text-cyber-blue/80 tracking-wider">{{ item.matricula }}</p>
            </td>
            <td class="px-4 py-3">
              <span class="inline-block whitespace-nowrap px-2.5 py-1 rounded-full border text-[11px] font-rajdhani font-semibold uppercase tracking-wider" :class="badgeClass(item.id_juego)">
                {{ item.juego }}
              </span>
            </td>
            <td class="px-4 py-3 text-gray-300 min-w-[180px]">{{ item.carrera }}</td>
            <td class="px-4 py-3 text-gray-300 whitespace-nowrap">{{ item.cuatrimestre }}° {{ item.grupo }}</td>
            <td class="px-4 py-3 min-w-[180px]">
              <p class="text-gray-200 whitespace-nowrap">{{ item.telefono }}</p>
              <p class="text-xs text-gray-500 truncate max-w-[220px]" :title="item.correo">{{ item.correo || '—' }}</p>
            </td>
            <td class="px-4 py-3 text-xs text-gray-400 whitespace-nowrap">{{ formatDate(item.fecha_inscripcion) }}</td>
            <td class="px-4 py-3 text-center">
              <button
                :id="`admin-toggle-${item.id_inscripcion}`"
                type="button"
                :disabled="togglingId === item.id_inscripcion"
                :title="item.asistio === 1 ? 'Marcar como no asistió' : 'Marcar como asistió'"
                class="inline-flex items-center justify-center w-8 h-8 rounded-full transition cursor-pointer disabled:opacity-40"
                :class="item.asistio === 1 ? 'text-emerald-400 hover:bg-emerald-400/10' : 'text-gray-600 hover:text-gray-300 hover:bg-white/5'"
                @click="emit('toggle-asistio', item)"
              >
                <CheckCircle2 v-if="item.asistio === 1" class="w-5 h-5 drop-shadow-[0_0_6px_rgba(16,185,129,0.8)]" />
                <Circle v-else class="w-5 h-5" />
              </button>
            </td>
            <td class="px-4 py-3">
              <div class="flex justify-end gap-1.5">
                <button
                  :id="`admin-edit-${item.id_inscripcion}`"
                  type="button"
                  title="Editar"
                  class="p-2 rounded-lg text-gray-400 hover:text-cyber-blue hover:bg-cyber-blue/10 transition cursor-pointer"
                  @click="emit('edit', item)"
                >
                  <Pencil class="w-4 h-4" />
                </button>
                <button
                  :id="`admin-delete-${item.id_inscripcion}`"
                  type="button"
                  title="Eliminar"
                  class="p-2 rounded-lg text-gray-400 hover:text-cyber-pink hover:bg-cyber-pink/10 transition cursor-pointer"
                  @click="emit('delete', item)"
                >
                  <Trash2 class="w-4 h-4" />
                </button>
              </div>
            </td>
          </tr>
        </TransitionGroup>
      </table>
    </div>
  </div>
</template>

<style scoped>
.admin-table-wrap {
  box-shadow: 0 0 30px rgba(0, 240, 255, 0.06);
}
.admin-row-enter-active,
.admin-row-leave-active {
  transition: opacity 0.25s ease, transform 0.25s ease;
}
.admin-row-enter-from,
.admin-row-leave-to {
  opacity: 0;
  transform: translateX(-12px);
}
</style>
