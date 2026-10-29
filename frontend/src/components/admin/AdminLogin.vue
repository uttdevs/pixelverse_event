<script setup>
/**
 * ============================================================================
 * COMPONENTE: AdminLogin.vue
 * ============================================================================
 * Pantalla de acceso al panel `/admin`. Solicita la clave administrativa
 * (ADMIN_API_KEY del backend) y la valida contra la API antes de continuar.
 */

import { ref } from 'vue'
import { Eye, EyeOff, KeyRound, Loader2, ShieldCheck } from 'lucide-vue-next'

const props = defineProps({
  login: { type: Function, required: true }
})

const key = ref('')
const showKey = ref(false)
const isSubmitting = ref(false)
const error = ref('')

const handleSubmit = async () => {
  if (isSubmitting.value) return
  error.value = ''
  isSubmitting.value = true
  try {
    await props.login(key.value)
    key.value = ''
  } catch (err) {
    error.value = err.message || 'No se pudo validar la clave'
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center px-4 cyber-grid">
    <div class="admin-login-card w-full max-w-md rounded-2xl p-8 sm:p-10 bg-cyber-card/80 backdrop-blur-xl border border-cyber-blue/40">
      <div class="flex items-center justify-center w-14 h-14 mx-auto mb-6 rounded-xl bg-gradient-to-br from-cyber-blue/20 to-cyber-purple/20 border border-cyber-blue/50">
        <ShieldCheck class="w-7 h-7 text-cyber-blue" />
      </div>

      <span class="block text-center font-rajdhani text-xs uppercase tracking-[0.3em] text-cyber-blue mb-2">
        // ACCESO RESTRINGIDO
      </span>
      <h1 class="admin-glitch text-center font-orbitron font-black text-2xl uppercase text-white mb-2" data-text="CENTRO DE MANDO">
        CENTRO DE MANDO
      </h1>
      <p class="text-center text-gray-400 text-sm mb-8">
        Ingresa la clave administrativa para gestionar las inscripciones.
      </p>

      <form class="space-y-5" @submit.prevent="handleSubmit">
        <div>
          <label for="admin-key" class="block font-rajdhani text-xs uppercase tracking-wider text-gray-300 mb-1.5">
            Clave administrativa
          </label>
          <div class="relative">
            <KeyRound class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-500" />
            <input
              id="admin-key"
              v-model="key"
              :type="showKey ? 'text' : 'password'"
              autocomplete="current-password"
              required
              autofocus
              placeholder="••••••••••••"
              class="w-full pl-10 pr-11 py-3 rounded-lg bg-cyber-dark border border-white/15 focus:border-cyber-blue focus:outline-none text-sm text-white font-orbitron tracking-wider"
            />
            <button
              id="admin-key-toggle"
              type="button"
              :aria-label="showKey ? 'Ocultar clave' : 'Mostrar clave'"
              class="absolute right-2 top-1/2 -translate-y-1/2 p-1.5 rounded text-gray-400 hover:text-cyber-blue transition cursor-pointer"
              @click="showKey = !showKey"
            >
              <EyeOff v-if="showKey" class="w-4 h-4" />
              <Eye v-else class="w-4 h-4" />
            </button>
          </div>
        </div>

        <div
          v-if="error"
          id="admin-login-error"
          role="alert"
          class="p-3 rounded-lg bg-cyber-pink/10 border border-cyber-pink/40 text-xs text-cyber-pink"
        >
          {{ error }}
        </div>

        <button
          id="admin-login-submit"
          type="submit"
          :disabled="isSubmitting"
          class="w-full flex items-center justify-center gap-2 py-3 rounded-lg font-orbitron font-bold text-sm uppercase tracking-wider text-cyber-dark bg-gradient-to-r from-cyber-blue to-cyber-purple hover:shadow-[0_0_25px_rgba(0,240,255,0.5)] transition disabled:opacity-60 cursor-pointer"
        >
          <Loader2 v-if="isSubmitting" class="w-4 h-4 animate-spin" />
          <span>{{ isSubmitting ? 'Verificando…' : 'Entrar' }}</span>
        </button>
      </form>
    </div>
  </div>
</template>

<style scoped>
.admin-login-card {
  box-shadow: 0 0 40px rgba(0, 240, 255, 0.12), 0 0 80px rgba(157, 0, 255, 0.1);
  animation: admin-rise 0.5s ease-out both;
}

@keyframes admin-rise {
  from { opacity: 0; transform: translateY(16px) scale(0.98); }
  to { opacity: 1; transform: none; }
}

.admin-glitch {
  position: relative;
}
.admin-glitch::before,
.admin-glitch::after {
  content: attr(data-text);
  position: absolute;
  inset: 0;
  opacity: 0.6;
  pointer-events: none;
}
.admin-glitch::before {
  color: var(--color-cyber-blue);
  animation: admin-glitch 3.5s infinite steps(1);
}
.admin-glitch::after {
  color: var(--color-cyber-pink);
  animation: admin-glitch 3.5s infinite steps(1) reverse;
}

@keyframes admin-glitch {
  0%, 92%, 100% { transform: none; clip-path: inset(0 0 100% 0); }
  93% { transform: translate(-2px, 1px); clip-path: inset(10% 0 60% 0); }
  95% { transform: translate(2px, -1px); clip-path: inset(55% 0 15% 0); }
  97% { transform: translate(-1px, 0); clip-path: inset(30% 0 40% 0); }
}

@media (prefers-reduced-motion: reduce) {
  .admin-login-card,
  .admin-glitch::before,
  .admin-glitch::after {
    animation: none;
  }
}
</style>
