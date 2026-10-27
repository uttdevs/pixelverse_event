/**
 * ============================================================================
 * COMPOSABLE: useAudioSfx.js
 * ============================================================================
 * ¿QUÉ ES UN COMPOSABLE EN VUE 3?
 * En lugar de tener funciones y variables globales dispersas en el objeto `window`
 * (como en el prototipo original), en Vue 3 encapsulamos la lógica y el estado
 * reutilizable dentro de funciones llamadas "composables".
 * 
 * Este composable implementa un sintetizador de efectos de sonido utilizando
 * la API nativa del navegador: Web Audio API. No requiere descargar archivos .mp3,
 * sino que genera las ondas de audio sintéticas (senoidales y triangulares) en vivo.
 */

import { ref } from 'vue'

// Estado reactivo global del sonido (compartido entre todos los componentes que usen el composable)
// Al usar `ref(true)`, cualquier cambio se reflejará instantáneamente en todos los componentes.
const isAudioEnabled = ref(true)

// Referencia al contexto de Web Audio API (se inicializa con la primera interacción del usuario)
let audioCtx = null

function initAudio() {
  if (!audioCtx && typeof window !== 'undefined') {
    const AudioContextClass = window.AudioContext || window.webkitAudioContext
    if (AudioContextClass) {
      audioCtx = new AudioContextClass()
    }
  }
}

export function useAudioSfx() {
  /**
   * Reproduce un sonido sintético agudo y rápido al pasar el cursor (Hover)
   * Utiliza una onda senoidal ('sine') con rampa exponencial de frecuencia.
   */
  const playHoverSound = () => {
    if (!isAudioEnabled.value) return
    try {
      initAudio()
      if (!audioCtx) return

      // Si el navegador suspendió el audio esperando interacción del usuario, lo reanudamos
      if (audioCtx.state === 'suspended') {
        audioCtx.resume()
      }

      const osc = audioCtx.createOscillator()
      const gain = audioCtx.createGain()

      osc.type = 'sine'
      osc.frequency.setValueAtTime(520, audioCtx.currentTime)
      osc.frequency.exponentialRampToValueAtTime(880, audioCtx.currentTime + 0.045)

      gain.gain.setValueAtTime(0.03, audioCtx.currentTime)
      gain.gain.exponentialRampToValueAtTime(0.001, audioCtx.currentTime + 0.045)

      osc.connect(gain)
      gain.connect(audioCtx.destination)

      osc.start()
      osc.stop(audioCtx.currentTime + 0.05)
    } catch (e) {
      // Ignoramos silenciosamente si el navegador bloquea la reproducción automática
    }
  }

  /**
   * Reproduce un sonido sintético con más cuerpo al hacer clic (Click)
   * Utiliza una onda triangular ('triangle') con rampa de tono y volumen.
   */
  const playClickSound = () => {
    if (!isAudioEnabled.value) return
    try {
      initAudio()
      if (!audioCtx) return

      if (audioCtx.state === 'suspended') {
        audioCtx.resume()
      }

      const osc = audioCtx.createOscillator()
      const gain = audioCtx.createGain()

      osc.type = 'triangle'
      osc.frequency.setValueAtTime(280, audioCtx.currentTime)
      osc.frequency.exponentialRampToValueAtTime(640, audioCtx.currentTime + 0.12)

      gain.gain.setValueAtTime(0.08, audioCtx.currentTime)
      gain.gain.exponentialRampToValueAtTime(0.001, audioCtx.currentTime + 0.12)

      osc.connect(gain)
      gain.connect(audioCtx.destination)

      osc.start()
      osc.stop(audioCtx.currentTime + 0.13)
    } catch (e) {
      // Manejo silencioso de restricciones de audio
    }
  }

  /**
   * Alterna el estado del audio entre activado y silenciado
   */
  const toggleAudio = () => {
    isAudioEnabled.value = !isAudioEnabled.value
    if (isAudioEnabled.value) {
      playClickSound()
    }
  }

  return {
    isAudioEnabled,
    playHoverSound,
    playClickSound,
    toggleAudio
  }
}
