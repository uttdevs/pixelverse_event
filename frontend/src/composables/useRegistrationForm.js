/**
 * ============================================================================
 * COMPOSABLE: useRegistrationForm.js
 * ============================================================================
 * Lógica compartida del formulario de inscripción, usada por
 * `RegistrationModal.vue` y `views/RegisterView.vue`.
 *
 * INTEGRACIÓN CON LA API:
 * - `GET /api/carreras` y `GET /api/juegos` alimentan los <select>.
 * - `POST /api/participantes` registra al participante y su inscripción.
 *   El backend rechaza campos desconocidos, por lo que el payload contiene
 *   exactamente: nombre_completo, matricula, id_carrera, cuatrimestre, grupo,
 *   telefono, correo (opcional) y juego_id.
 */

import { reactive, ref, computed } from 'vue'

const fetchJSON = async (url, options = {}) => {
  const response = await fetch(url, {
    ...options,
    headers: { Accept: 'application/json', ...(options.headers || {}) }
  })
  let body = null
  try {
    body = await response.json()
  } catch {
    body = null
  }
  if (!response.ok) {
    const error = new Error(body?.message || `Error ${response.status} al comunicarse con el servidor`)
    error.details = Array.isArray(body?.errors) ? body.errors : []
    throw error
  }
  return body
}

const normalize = (text) =>
  String(text ?? '')
    .toLowerCase()
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '')
    .replace(/[^a-z0-9]/g, '')

const createEmptyForm = () => ({
  nombre_completo: '',
  matricula: '',
  id_carrera: null,
  cuatrimestre: null,
  grupo: '',
  telefono: '',
  correo: '',
  juego_id: null
})

/**
 * @param {Object} [options]
 * @param {() => (string|number|undefined)} [options.initialGame]
 *   Getter con el juego a preseleccionar (id numérico o nombre).
 */
export function useRegistrationForm({ initialGame = () => '' } = {}) {
  // --------------------------------------------------------------------------
  // CATÁLOGOS
  // --------------------------------------------------------------------------
  const carreras = ref([])
  const juegos = ref([])
  const isLoadingCatalogs = ref(false)
  const catalogError = ref('')

  // --------------------------------------------------------------------------
  // FORMULARIO
  // --------------------------------------------------------------------------
  const form = reactive(createEmptyForm())
  const isSubmitted = ref(false)
  const isSubmitting = ref(false)
  const submitError = ref('')
  const submitErrorDetails = ref([])
  const registration = ref(null)

  // Preselecciona el juego por id o por nombre (coincidencia flexible)
  const preselectInitialGame = () => {
    const wanted = initialGame()
    if (wanted === undefined || wanted === null || wanted === '' || !juegos.value.length) return

    const asId = Number(wanted)
    let match = Number.isInteger(asId) ? juegos.value.find((juego) => juego.id_juego === asId) : null
    if (!match) {
      const target = normalize(wanted)
      if (!target) return
      match = juegos.value.find((juego) => {
        const name = normalize(juego.nombre)
        return name === target || name.includes(target) || target.includes(name)
      })
    }
    if (match) form.juego_id = match.id_juego
  }

  const loadCatalogs = async () => {
    if (isLoadingCatalogs.value) return
    isLoadingCatalogs.value = true
    catalogError.value = ''
    try {
      const [carrerasData, juegosData] = await Promise.all([
        fetchJSON('/api/carreras'),
        fetchJSON('/api/juegos')
      ])
      carreras.value = Array.isArray(carrerasData) ? carrerasData : []
      juegos.value = Array.isArray(juegosData) ? juegosData : []
      preselectInitialGame()
    } catch (error) {
      catalogError.value = error.message || 'No se pudieron cargar los catálogos'
    } finally {
      isLoadingCatalogs.value = false
    }
  }

  const hasCatalogs = computed(() => carreras.value.length > 0 && juegos.value.length > 0)

  const selectedGameName = computed(() => {
    const id = registration.value?.juego_id ?? form.juego_id
    return juegos.value.find((juego) => juego.id_juego === id)?.nombre || ''
  })

  const clearErrors = () => {
    submitError.value = ''
    submitErrorDetails.value = []
  }

  const resetForm = () => {
    Object.assign(form, createEmptyForm())
    clearErrors()
    registration.value = null
    isSubmitted.value = false
    preselectInitialGame()
  }

  /**
   * Envía el formulario a la API.
   * @returns {Promise<Object|null>} payload + data de la respuesta, o null si falló.
   */
  const submit = async () => {
    if (isSubmitting.value) return null
    clearErrors()

    const payload = {
      nombre_completo: form.nombre_completo.trim(),
      matricula: form.matricula.trim(),
      id_carrera: Number(form.id_carrera),
      cuatrimestre: Number(form.cuatrimestre),
      grupo: form.grupo.trim().toUpperCase(),
      telefono: form.telefono.trim(),
      correo: form.correo.trim(),
      juego_id: Number(form.juego_id)
    }

    isSubmitting.value = true
    try {
      const response = await fetchJSON('/api/participantes', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      })
      registration.value = { ...response.data, message: response.message }
      isSubmitted.value = true
      return { ...payload, ...response.data }
    } catch (error) {
      submitError.value = error.message || 'No se pudo completar el registro'
      submitErrorDetails.value = error.details || []
      return null
    } finally {
      isSubmitting.value = false
    }
  }

  return {
    // catálogos
    carreras,
    juegos,
    isLoadingCatalogs,
    catalogError,
    hasCatalogs,
    loadCatalogs,
    // formulario
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
  }
}
