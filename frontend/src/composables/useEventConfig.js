import { ref } from 'vue'

const eventTimestamp = ref(null)
const registrationsOpen = ref(false)
const serverOffset = ref(0)

let configPromise = null

const loadConfig = () => {
  if (configPromise) return configPromise

  configPromise = fetch('/api/config')
    .then(async (res) => {
      const contentType = res.headers.get('content-type') || ''
      if (!res.ok || !contentType.includes('application/json')) {
        throw new Error(`Respuesta inválida del servidor (${res.status})`)
      }

      const data = await res.json()
      eventTimestamp.value = data.fecha_evento
        ? new Date(data.fecha_evento).getTime()
        : null
      serverOffset.value = data.hora_servidor
        ? new Date(data.hora_servidor).getTime() - Date.now()
        : 0
      registrationsOpen.value = data.registros_abiertos === true

      return true
    })
    .catch((err) => {
      console.warn('El servidor no respondió o falló /api/config:', err)
      eventTimestamp.value = null
      serverOffset.value = 0
      registrationsOpen.value = false
      return false
    })

  return configPromise
}

export function useEventConfig() {
  return {
    eventTimestamp,
    registrationsOpen,
    serverOffset,
    loadConfig
  }
}
