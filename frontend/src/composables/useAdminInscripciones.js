/**
 * ============================================================================
 * COMPOSABLE: useAdminInscripciones.js
 * ============================================================================
 * Lógica del panel de administración (`/admin`).
 *
 * INTEGRACIÓN CON LA API (todas las rutas requieren el encabezado X-Admin-Key):
 * - `GET    /api/inscripciones?q=&juego_id=&id_carrera=&asistio=`  Listado filtrado.
 * - `PATCH  /api/inscripciones/{id}`  Edita `juego_id` y/o `asistio`.
 * - `DELETE /api/inscripciones/{id}`  Elimina la inscripción (y al participante
 *   si era su última inscripción).
 * - `GET    /api/inscripciones/export?...`  Descarga .xlsx con los mismos filtros.
 *
 * La clave se guarda en `sessionStorage`: se borra al cerrar la pestaña.
 */

import { computed, reactive, ref, watch } from 'vue'

const STORAGE_KEY = 'pixelverse_admin_key'
const DEBOUNCE_MS = 300

const readStoredKey = () => {
  try {
    return sessionStorage.getItem(STORAGE_KEY) || ''
  } catch {
    return ''
  }
}

const createEmptyFilters = () => ({ q: '', juego_id: '', id_carrera: '', asistio: '' })

export function useAdminInscripciones() {
  // --------------------------------------------------------------------------
  // SESIÓN
  // --------------------------------------------------------------------------
  const adminKey = ref(readStoredKey())
  const isAuthenticated = computed(() => adminKey.value !== '')

  const persistKey = (value) => {
    adminKey.value = value
    try {
      if (value) sessionStorage.setItem(STORAGE_KEY, value)
      else sessionStorage.removeItem(STORAGE_KEY)
    } catch {
      /* sessionStorage no disponible: la sesión vive solo en memoria */
    }
  }

  // --------------------------------------------------------------------------
  // ESTADO
  // --------------------------------------------------------------------------
  const inscripciones = ref([])
  const juegos = ref([])
  const carreras = ref([])
  const filters = reactive(createEmptyFilters())
  const isLoading = ref(false)
  const isExporting = ref(false)
  const loadError = ref('')

  const buildQuery = () => {
    const params = new URLSearchParams()
    Object.entries(filters).forEach(([key, value]) => {
      const text = String(value ?? '').trim()
      if (text !== '') params.set(key, text)
    })
    const query = params.toString()
    return query ? `?${query}` : ''
  }

  const hasActiveFilters = computed(() => Object.values(filters).some((value) => String(value ?? '').trim() !== ''))

  /**
   * fetch con el encabezado administrativo. Si la API responde 401 se cierra la sesión.
   * @param {string} url
   * @param {RequestInit} [options]
   * @param {string} [keyOverride] clave a usar en lugar de la almacenada (login).
   */
  const adminFetch = async (url, options = {}, keyOverride) => {
    const response = await fetch(url, {
      ...options,
      headers: {
        Accept: 'application/json',
        ...(options.headers || {}),
        'X-Admin-Key': keyOverride ?? adminKey.value
      }
    })
    if (response.ok) return response

    let body = null
    try {
      body = await response.json()
    } catch {
      body = null
    }
    if (response.status === 401 && keyOverride === undefined) persistKey('')
    const error = new Error(body?.message || `Error ${response.status} al comunicarse con el servidor`)
    error.status = response.status
    error.details = Array.isArray(body?.errors) ? body.errors : []
    throw error
  }

  // --------------------------------------------------------------------------
  // CATÁLOGOS (públicos)
  // --------------------------------------------------------------------------
  const loadCatalogs = async () => {
    try {
      const [juegosResponse, carrerasResponse] = await Promise.all([fetch('/api/juegos'), fetch('/api/carreras')])
      if (juegosResponse.ok) juegos.value = await juegosResponse.json()
      if (carrerasResponse.ok) carreras.value = await carrerasResponse.json()
    } catch {
      /* los filtros por catálogo simplemente quedan vacíos */
    }
  }

  // --------------------------------------------------------------------------
  // LISTADO
  // --------------------------------------------------------------------------
  let requestCounter = 0
  const load = async () => {
    if (!isAuthenticated.value) return
    const requestId = ++requestCounter
    isLoading.value = true
    loadError.value = ''
    try {
      const response = await adminFetch(`/api/inscripciones${buildQuery()}`)
      const data = await response.json()
      // Ignora respuestas de búsquedas anteriores que lleguen tarde.
      if (requestId === requestCounter) inscripciones.value = Array.isArray(data) ? data : []
    } catch (error) {
      if (requestId === requestCounter) loadError.value = error.message
    } finally {
      if (requestId === requestCounter) isLoading.value = false
    }
  }

  let debounceTimer = null
  watch(
    () => ({ ...filters }),
    () => {
      clearTimeout(debounceTimer)
      debounceTimer = setTimeout(load, DEBOUNCE_MS)
    }
  )

  const resetFilters = () => Object.assign(filters, createEmptyFilters())

  // --------------------------------------------------------------------------
  // SESIÓN: LOGIN / LOGOUT
  // --------------------------------------------------------------------------
  const login = async (key) => {
    const trimmed = String(key ?? '').trim()
    if (!trimmed) throw new Error('Ingresa la clave administrativa')
    const response = await adminFetch('/api/inscripciones', {}, trimmed)
    const data = await response.json()
    persistKey(trimmed)
    inscripciones.value = Array.isArray(data) ? data : []
    resetFilters()
  }

  const logout = () => {
    persistKey('')
    inscripciones.value = []
    loadError.value = ''
  }

  // --------------------------------------------------------------------------
  // EDICIÓN Y ELIMINACIÓN
  // --------------------------------------------------------------------------
  const update = async (id, patch) => {
    const response = await adminFetch(`/api/inscripciones/${id}`, {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(patch)
    })
    const body = await response.json()
    const index = inscripciones.value.findIndex((item) => item.id_inscripcion === id)
    if (index !== -1 && body?.data) inscripciones.value.splice(index, 1, body.data)
    return body
  }

  const remove = async (id) => {
    const response = await adminFetch(`/api/inscripciones/${id}`, { method: 'DELETE' })
    const body = await response.json()
    inscripciones.value = inscripciones.value.filter((item) => item.id_inscripcion !== id)
    return body
  }

  /** Número de inscripciones cargadas de un participante (para advertir en el borrado). */
  const countByParticipante = (idParticipante) =>
    inscripciones.value.filter((item) => item.id_participante === idParticipante).length

  // --------------------------------------------------------------------------
  // EXPORTACIÓN A EXCEL
  // --------------------------------------------------------------------------
  const exportExcel = async () => {
    if (isExporting.value) return
    isExporting.value = true
    try {
      const response = await adminFetch(`/api/inscripciones/export${buildQuery()}`, {
        headers: { Accept: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' }
      })
      const blob = await response.blob()
      const disposition = response.headers.get('Content-Disposition') || ''
      const filename = disposition.match(/filename="?([^"]+)"?/)?.[1] || 'inscripciones_pixelverse.xlsx'

      const url = URL.createObjectURL(blob)
      const anchor = document.createElement('a')
      anchor.href = url
      anchor.download = filename
      document.body.appendChild(anchor)
      anchor.click()
      anchor.remove()
      setTimeout(() => URL.revokeObjectURL(url), 1000)
      return filename
    } finally {
      isExporting.value = false
    }
  }

  // --------------------------------------------------------------------------
  // ESTADÍSTICAS (sobre el listado cargado)
  // --------------------------------------------------------------------------
  const stats = computed(() => {
    const total = inscripciones.value.length
    const asistieron = inscripciones.value.filter((item) => item.asistio === 1).length
    const participantes = new Set(inscripciones.value.map((item) => item.id_participante)).size
    const porJuego = {}
    inscripciones.value.forEach((item) => {
      porJuego[item.juego] = (porJuego[item.juego] || 0) + 1
    })
    return { total, asistieron, pendientes: total - asistieron, participantes, porJuego }
  })

  return {
    // sesión
    isAuthenticated,
    login,
    logout,
    // datos
    inscripciones,
    juegos,
    carreras,
    filters,
    hasActiveFilters,
    isLoading,
    isExporting,
    loadError,
    stats,
    // acciones
    loadCatalogs,
    load,
    resetFilters,
    update,
    remove,
    countByParticipante,
    exportExcel
  }
}
