<script setup>
/**
 * ============================================================================
 * VISTA: AdminView.vue (Panel de Administración de Inscripciones)
 * ============================================================================
 * Ruta oculta `/admin` (sin enlaces desde la landing, marcada como noindex).
 * Se carga de forma diferida para no aumentar el bundle público.
 *
 * Funciones:
 * - Acceso mediante la clave administrativa (X-Admin-Key).
 * - Consulta, búsqueda y filtrado de inscripciones (server-side).
 * - Edición de juego/asistencia y eliminación con confirmación.
 * - Exportación a Excel (.xlsx generado en el backend con excelize).
 */

import { computed, onMounted, ref, watch } from 'vue'
import {
  CheckCircle2, Clock, Download, FilterX, Gamepad2, Loader2, LogOut, RefreshCw, Search, Users, XCircle
} from 'lucide-vue-next'
import AdminLogin from '../components/admin/AdminLogin.vue'
import InscripcionesTable from '../components/admin/InscripcionesTable.vue'
import EditInscripcionModal from '../components/admin/EditInscripcionModal.vue'
import ConfirmDialog from '../components/admin/ConfirmDialog.vue'
import { useAdminInscripciones } from '../composables/useAdminInscripciones'

const {
  isAuthenticated, login, logout,
  inscripciones, juegos, carreras, filters, hasActiveFilters,
  isLoading, isExporting, loadError, stats,
  loadCatalogs, load, resetFilters, update, remove, countByParticipante, exportExcel
} = useAdminInscripciones()

// --------------------------------------------------------------------------
// TOASTS
// --------------------------------------------------------------------------
const toasts = ref([])
let toastId = 0
const notify = (message, type = 'success') => {
  const id = ++toastId
  toasts.value.push({ id, message, type })
  setTimeout(() => {
    toasts.value = toasts.value.filter((toast) => toast.id !== id)
  }, 3500)
}

// --------------------------------------------------------------------------
// CICLO DE VIDA
// --------------------------------------------------------------------------
onMounted(() => {
  loadCatalogs()
  if (isAuthenticated.value) load()
})

watch(isAuthenticated, (value, previous) => {
  if (!value && previous) notify('Sesión cerrada', 'info')
})

// --------------------------------------------------------------------------
// ACCIONES
// --------------------------------------------------------------------------
const editing = ref(null)
const deleting = ref(null)
const isDeleting = ref(false)
const togglingId = ref(null)

const handleSaved = (result) => {
  editing.value = null
  notify(result?.message || 'Inscripción actualizada')
}

const handleToggle = async (item) => {
  const nuevoValor = item.asistio === 1 ? 0 : 1
  togglingId.value = item.id_inscripcion
  try {
    await update(item.id_inscripcion, { asistio: nuevoValor })
    notify(nuevoValor === 1 ? `Asistencia registrada: ${item.participante}` : `Asistencia retirada: ${item.participante}`)
  } catch (error) {
    notify(error.message, 'error')
  } finally {
    togglingId.value = null
  }
}

const deleteWarning = computed(() => {
  if (!deleting.value) return ''
  return countByParticipante(deleting.value.id_participante) <= 1
    ? 'Es la única inscripción de este participante: también se eliminarán sus datos personales.'
    : ''
})

const handleConfirmDelete = async () => {
  if (!deleting.value) return
  isDeleting.value = true
  try {
    const result = await remove(deleting.value.id_inscripcion)
    notify(result?.message || 'Inscripción eliminada')
    deleting.value = null
  } catch (error) {
    notify(error.message, 'error')
  } finally {
    isDeleting.value = false
  }
}

const handleExport = async () => {
  try {
    const filename = await exportExcel()
    if (filename) notify(`Archivo descargado: ${filename}`)
  } catch (error) {
    notify(error.message, 'error')
  }
}

const statCards = computed(() => [
  { id: 'total', label: 'Inscripciones', value: stats.value.total, icon: Gamepad2, accent: 'from-cyber-blue/30 to-cyber-blue/0 text-cyber-blue' },
  { id: 'participantes', label: 'Participantes', value: stats.value.participantes, icon: Users, accent: 'from-cyber-purple/30 to-cyber-purple/0 text-violet-300' },
  { id: 'asistieron', label: 'Asistieron', value: stats.value.asistieron, icon: CheckCircle2, accent: 'from-emerald-500/30 to-emerald-500/0 text-emerald-300' },
  { id: 'pendientes', label: 'Pendientes', value: stats.value.pendientes, icon: Clock, accent: 'from-amber-400/30 to-amber-400/0 text-amber-300' }
])

const selectClass = 'px-3 py-2.5 rounded-lg bg-cyber-dark border border-white/15 focus:border-cyber-blue focus:outline-none text-sm text-white font-rajdhani'
</script>

<template>
  <AdminLogin v-if="!isAuthenticated" :login="login" />

  <div v-else class="min-h-screen cyber-grid">
    <!-- Barra superior -->
    <header class="sticky top-0 z-30 border-b border-white/10 bg-cyber-dark/80 backdrop-blur-xl">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between gap-4">
        <div class="flex items-center gap-3">
          <span class="font-orbitron font-black text-lg tracking-wider text-white">
            PIXEL <span class="text-transparent bg-clip-text bg-gradient-to-r from-cyber-blue to-cyber-purple">VERSE</span>
          </span>
          <span class="hidden sm:inline-block px-2 py-0.5 rounded border border-cyber-pink/50 text-cyber-pink font-rajdhani text-[10px] uppercase tracking-[0.25em]">
            Admin
          </span>
        </div>
        <div class="flex items-center gap-2">
          <button
            id="admin-export"
            type="button"
            :disabled="isExporting || !stats.total"
            class="flex items-center gap-2 px-4 py-2 rounded-lg font-orbitron font-bold text-xs uppercase tracking-wider text-cyber-dark bg-gradient-to-r from-emerald-400 to-cyber-blue hover:shadow-[0_0_20px_rgba(16,185,129,0.5)] transition disabled:opacity-40 disabled:cursor-not-allowed cursor-pointer"
            @click="handleExport"
          >
            <Loader2 v-if="isExporting" class="w-4 h-4 animate-spin" />
            <Download v-else class="w-4 h-4" />
            <span class="hidden sm:inline">Exportar Excel</span>
          </button>
          <button
            id="admin-logout"
            type="button"
            title="Cerrar sesión"
            class="flex items-center gap-2 px-3 py-2 rounded-lg text-xs font-rajdhani uppercase tracking-wider text-gray-300 bg-white/5 hover:bg-white/10 hover:text-cyber-pink border border-white/10 transition cursor-pointer"
            @click="logout"
          >
            <LogOut class="w-4 h-4" />
            <span class="hidden sm:inline">Salir</span>
          </button>
        </div>
      </div>
    </header>

    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-6">
      <!-- Título -->
      <div>
        <span class="font-rajdhani text-xs uppercase tracking-[0.25em] text-cyber-blue">// CENTRO DE MANDO</span>
        <h1 class="font-orbitron font-black text-2xl sm:text-3xl uppercase text-white">
          Gestión de <span class="text-transparent bg-clip-text bg-gradient-to-r from-cyber-blue to-cyber-purple">inscripciones</span>
        </h1>
      </div>

      <!-- Estadísticas -->
      <section class="grid grid-cols-2 lg:grid-cols-4 gap-4" aria-label="Estadísticas">
        <article
          v-for="card in statCards"
          :id="`admin-stat-${card.id}`"
          :key="card.id"
          class="admin-stat relative overflow-hidden rounded-2xl p-5 bg-cyber-card/70 backdrop-blur-xl border border-white/10 hover:border-cyber-blue/40 transition"
        >
          <div class="absolute inset-0 bg-gradient-to-br opacity-60 pointer-events-none" :class="card.accent" />
          <div class="relative flex items-start justify-between">
            <div>
              <p class="font-rajdhani text-[11px] uppercase tracking-[0.2em] text-gray-400">{{ card.label }}</p>
              <p class="mt-1 font-orbitron font-black text-3xl text-white">{{ card.value }}</p>
            </div>
            <component :is="card.icon" class="w-6 h-6" :class="card.accent.split(' ').pop()" />
          </div>
        </article>
      </section>

      <!-- Conteo por juego -->
      <section v-if="Object.keys(stats.porJuego).length" class="flex flex-wrap gap-2" aria-label="Inscripciones por juego">
        <span
          v-for="(count, juego) in stats.porJuego"
          :key="juego"
          class="px-3 py-1 rounded-full bg-white/5 border border-white/10 text-xs font-rajdhani text-gray-300"
        >
          {{ juego }} · <strong class="text-white">{{ count }}</strong>
        </span>
      </section>

      <!-- Filtros -->
      <section class="flex flex-col lg:flex-row gap-3 p-4 rounded-2xl bg-cyber-card/70 backdrop-blur-xl border border-white/10" aria-label="Filtros">
        <div class="relative flex-1">
          <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-500" />
          <input
            id="admin-search"
            v-model="filters.q"
            type="search"
            maxlength="100"
            placeholder="Buscar por nombre, matrícula, correo, teléfono o grupo…"
            class="w-full pl-10 pr-4 py-2.5 rounded-lg bg-cyber-dark border border-white/15 focus:border-cyber-blue focus:outline-none text-sm text-white"
          />
        </div>
        <select id="admin-filter-juego" v-model="filters.juego_id" :class="selectClass" aria-label="Filtrar por juego">
          <option value="">Todos los juegos</option>
          <option v-for="juego in juegos" :key="juego.id_juego" :value="String(juego.id_juego)">{{ juego.nombre }}</option>
        </select>
        <select id="admin-filter-carrera" v-model="filters.id_carrera" :class="[selectClass, 'lg:max-w-[220px]']" aria-label="Filtrar por carrera">
          <option value="">Todas las carreras</option>
          <option v-for="carrera in carreras" :key="carrera.id_carrera" :value="String(carrera.id_carrera)">{{ carrera.nombre }}</option>
        </select>
        <select id="admin-filter-asistio" v-model="filters.asistio" :class="selectClass" aria-label="Filtrar por asistencia">
          <option value="">Asistencia: todas</option>
          <option value="1">Asistieron</option>
          <option value="0">Pendientes</option>
        </select>
        <div class="flex gap-2">
          <button
            id="admin-filters-reset"
            type="button"
            :disabled="!hasActiveFilters"
            title="Limpiar filtros"
            class="flex items-center justify-center gap-2 px-3 py-2.5 rounded-lg text-gray-300 bg-white/5 hover:bg-white/10 border border-white/10 transition disabled:opacity-30 cursor-pointer"
            @click="resetFilters"
          >
            <FilterX class="w-4 h-4" />
          </button>
          <button
            id="admin-refresh"
            type="button"
            title="Recargar"
            class="flex items-center justify-center gap-2 px-3 py-2.5 rounded-lg text-gray-300 bg-white/5 hover:bg-white/10 hover:text-cyber-blue border border-white/10 transition cursor-pointer"
            @click="load"
          >
            <RefreshCw class="w-4 h-4" :class="{ 'animate-spin': isLoading }" />
          </button>
        </div>
      </section>

      <!-- Error -->
      <div v-if="loadError" id="admin-load-error" role="alert" class="flex items-center justify-between gap-3 p-4 rounded-xl bg-cyber-pink/10 border border-cyber-pink/40 text-sm text-cyber-pink">
        <span>{{ loadError }}</span>
        <button type="button" class="px-3 py-1 rounded bg-white/10 hover:bg-white/20 text-white text-xs uppercase font-orbitron cursor-pointer" @click="load">Reintentar</button>
      </div>

      <!-- Skeleton -->
      <div v-else-if="isLoading && !inscripciones.length" class="space-y-2" aria-busy="true">
        <div v-for="n in 6" :key="n" class="h-14 rounded-xl bg-white/[0.04] animate-pulse" />
      </div>

      <!-- Vacío -->
      <div v-else-if="!inscripciones.length" id="admin-empty" class="py-16 text-center rounded-2xl border border-dashed border-white/15 bg-cyber-card/40">
        <Search class="w-10 h-10 mx-auto text-gray-600 mb-3" />
        <p class="font-orbitron text-white uppercase text-sm">Sin resultados</p>
        <p class="text-gray-500 text-sm mt-1">
          {{ hasActiveFilters ? 'Ninguna inscripción coincide con los filtros.' : 'Aún no hay inscripciones registradas.' }}
        </p>
      </div>

      <!-- Tabla -->
      <template v-else>
        <p class="text-xs font-rajdhani uppercase tracking-wider text-gray-500">
          Mostrando {{ inscripciones.length }} inscripción(es){{ hasActiveFilters ? ' filtradas' : '' }}
          <span v-if="hasActiveFilters"> · la exportación respeta los filtros</span>
        </p>
        <InscripcionesTable
          :inscripciones="inscripciones"
          :toggling-id="togglingId"
          @edit="editing = $event"
          @delete="deleting = $event"
          @toggle-asistio="handleToggle"
        />
      </template>
    </main>

    <!-- Modales -->
    <EditInscripcionModal
      v-if="editing"
      :inscripcion="editing"
      :juegos="juegos"
      :save="update"
      @close="editing = null"
      @saved="handleSaved"
    />
    <ConfirmDialog
      v-if="deleting"
      title="Eliminar inscripción"
      :message="`¿Eliminar la inscripción de ${deleting.participante} en ${deleting.juego}? Esta acción no se puede deshacer.`"
      :warning="deleteWarning"
      :is-processing="isDeleting"
      @cancel="deleting = null"
      @confirm="handleConfirmDelete"
    />
  </div>

  <!-- Toasts -->
  <div class="fixed bottom-4 right-4 z-[60] flex flex-col gap-2 w-[min(92vw,360px)]" aria-live="polite">
    <TransitionGroup name="admin-toast">
      <div
        v-for="toast in toasts"
        :key="toast.id"
        class="flex items-start gap-2.5 p-3.5 rounded-xl backdrop-blur-xl border text-sm shadow-lg"
        :class="{
          'bg-emerald-500/15 border-emerald-400/40 text-emerald-200': toast.type === 'success',
          'bg-cyber-pink/15 border-cyber-pink/40 text-pink-200': toast.type === 'error',
          'bg-cyber-blue/10 border-cyber-blue/40 text-cyan-100': toast.type === 'info'
        }"
      >
        <XCircle v-if="toast.type === 'error'" class="w-4 h-4 mt-0.5 shrink-0" />
        <CheckCircle2 v-else class="w-4 h-4 mt-0.5 shrink-0" />
        <span>{{ toast.message }}</span>
      </div>
    </TransitionGroup>
  </div>
</template>

<style scoped>
.admin-stat {
  animation: admin-fade-up 0.4s ease-out both;
}
.admin-stat:nth-child(2) { animation-delay: 0.05s; }
.admin-stat:nth-child(3) { animation-delay: 0.1s; }
.admin-stat:nth-child(4) { animation-delay: 0.15s; }

@keyframes admin-fade-up {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: none; }
}

.admin-toast-enter-active,
.admin-toast-leave-active {
  transition: all 0.25s ease;
}
.admin-toast-enter-from,
.admin-toast-leave-to {
  opacity: 0;
  transform: translateX(24px);
}
</style>
