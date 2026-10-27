<script setup>
/**
 * ============================================================================
 * VISTA: RegisterView.vue (Página Dedicada de Registro de Squads)
 * ============================================================================
 * Vista completa para enrolamiento directo a través de la ruta `/registro`.
 * Ideal para compartir enlaces directos en Discord, Twitter o transmisiones.
 */

import { reactive, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { ArrowLeft, Check, Crosshair, Zap } from 'lucide-vue-next'
import { useAudioSfx } from '../composables/useAudioSfx'

const emit = defineEmits(['registered'])
const { playClickSound, playHoverSound } = useAudioSfx()

const form = reactive({
  squadName: '',
  game: 'Valorant: Protocol',
  discord: '',
  tier: 'Diamante+'
})

const isSubmitted = ref(false)
const generatedToken = ref('')

const handleSubmit = () => {
  playClickSound()
  generatedToken.value = 'PV-2026-' + Math.floor(1000 + Math.random() * 9000)
  isSubmitted.value = true

  emit('registered', {
    ...form,
    token: generatedToken.value
  })
}
</script>

<template>
  <div class="relative z-10 py-16 max-w-2xl mx-auto px-4 sm:px-6 lg:px-8">
    
    <!-- Volver al Inicio -->
    <div class="mb-8">
      <RouterLink
        to="/"
        @mouseenter="playHoverSound"
        class="inline-flex items-center gap-2 text-xs font-rajdhani uppercase tracking-widest text-cyber-blue hover:text-white transition"
      >
        <ArrowLeft class="w-4 h-4" />
        <span>Volver a la Arena Principal</span>
      </RouterLink>
    </div>

    <div class="gradient-border rounded-2xl p-6 sm:p-10 bg-cyber-card border border-cyber-blue/50 box-glow-combo">
      
      <!-- Encabezado de la Página -->
      <div class="mb-8">
        <span class="font-rajdhani text-xs uppercase tracking-[0.25em] text-cyber-blue block mb-1">
          // ENROLAMIENTO OFICIAL EN LÍNEA
        </span>
        <h1 class="font-orbitron font-black text-2xl sm:text-3xl uppercase text-white">
          INSCRIPCIÓN DE <span class="text-transparent bg-clip-text bg-gradient-to-r from-cyber-blue to-cyber-purple">SQUAD</span>
        </h1>
        <p class="text-gray-400 text-xs sm:text-sm mt-1">
          Asegura el cupo de tu escuadrón en los servidores dedicados de Pixel Verse 2026.
        </p>
      </div>

      <!-- Formulario de Registro -->
      <form
        v-if="!isSubmitted"
        @submit.prevent="handleSubmit"
        class="space-y-5"
      >
        <!-- Nombre del Squad -->
        <div>
          <label class="block font-rajdhani text-xs uppercase tracking-wider text-gray-300 mb-1.5">
            NOMBRE DEL SQUAD / CLAN *
          </label>
          <input
            v-model="form.squadName"
            type="text"
            required
            placeholder="Ej. Cyber Phantoms"
            class="w-full px-4 py-3 rounded-lg bg-cyber-dark border border-white/15 focus:border-cyber-blue focus:outline-none text-sm text-white font-orbitron"
          />
        </div>

        <!-- Selección de Arena y Discord -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="block font-rajdhani text-xs uppercase tracking-wider text-gray-300 mb-1.5">
              DISCIPLINA / ARENA *
            </label>
            <select
              v-model="form.game"
              class="w-full px-4 py-3 rounded-lg bg-cyber-dark border border-white/15 focus:border-cyber-blue focus:outline-none text-sm text-white font-rajdhani"
            >
              <option value="Valorant: Protocol">Valorant: Protocol (5v5)</option>
              <option value="League of Pixels">League of Pixels (5v5)</option>
              <option value="CS2: Cyber Strike">CS2: Cyber Strike (5v5)</option>
              <option value="Apex: Neon Drop">Apex: Neon Drop (Trios)</option>
            </select>
          </div>

          <div>
            <label class="block font-rajdhani text-xs uppercase tracking-wider text-gray-300 mb-1.5">
              DISCORD DEL CAPITÁN *
            </label>
            <input
              v-model="form.discord"
              type="text"
              required
              placeholder="usuario#0001"
              class="w-full px-4 py-3 rounded-lg bg-cyber-dark border border-white/15 focus:border-cyber-purple focus:outline-none text-sm text-white"
            />
          </div>
        </div>

        <!-- Rango Promedio -->
        <div>
          <label class="block font-rajdhani text-xs uppercase tracking-wider text-gray-300 mb-1.5">
            RANGO PROMEDIO DEL EQUIPO
          </label>
          <div class="grid grid-cols-3 gap-2 font-rajdhani text-xs">
            <label class="cursor-pointer">
              <input
                v-model="form.tier"
                type="radio"
                value="Amateur"
                class="peer hidden"
              />
              <div class="text-center py-2.5 rounded border border-white/10 bg-white/5 peer-checked:border-cyber-blue peer-checked:bg-cyber-blue/20 peer-checked:text-cyber-blue transition">
                AMATEUR / ORO+
              </div>
            </label>

            <label class="cursor-pointer">
              <input
                v-model="form.tier"
                type="radio"
                value="Diamante+"
                class="peer hidden"
              />
              <div class="text-center py-2.5 rounded border border-white/10 bg-white/5 peer-checked:border-cyber-blue peer-checked:bg-cyber-blue/20 peer-checked:text-cyber-blue transition">
                DIAMANTE / ASC
              </div>
            </label>

            <label class="cursor-pointer">
              <input
                v-model="form.tier"
                type="radio"
                value="Inmortal/Pro"
                class="peer hidden"
              />
              <div class="text-center py-2.5 rounded border border-white/10 bg-white/5 peer-checked:border-cyber-purple peer-checked:bg-cyber-purple/20 peer-checked:text-purple-300 transition">
                INMORTAL / PRO
              </div>
            </label>
          </div>
        </div>

        <!-- Botón de Envío -->
        <div class="pt-4">
          <button
            type="submit"
            class="w-full cyber-clip py-4 bg-gradient-to-r from-cyber-blue to-cyber-purple font-orbitron font-black text-sm text-cyber-dark uppercase tracking-wider hover:opacity-95 transition cursor-pointer flex items-center justify-center gap-2"
          >
            <Crosshair class="w-4 h-4 text-cyber-dark" />
            <span>CONFIRMAR LUGAR EN LA ARENA</span>
          </button>
        </div>
      </form>

      <!-- Mensaje de Confirmación y Token -->
      <div v-else class="text-center py-8 space-y-5">
        <div class="w-16 h-16 mx-auto rounded-full bg-emerald-500/20 border border-emerald-400 flex items-center justify-center text-emerald-400">
          <Check class="w-8 h-8" />
        </div>

        <h2 class="font-orbitron font-bold text-2xl text-white">¡SQUAD REGISTRADO CON ÉXITO!</h2>
        
        <p class="text-sm text-gray-300 max-w-md mx-auto leading-relaxed font-rajdhani">
          El escuadrón <strong class="text-white">"{{ form.squadName }}"</strong> ha quedado oficialmente pre-registrado en la arena
          <strong class="text-cyber-blue">{{ form.game }}</strong>. Las instrucciones del check-in han sido enviadas a <strong class="text-purple-300">{{ form.discord }}</strong>.
        </p>

        <div class="p-4 rounded-xl bg-cyber-dark border border-cyber-blue/40 font-mono text-sm text-cyber-blue max-w-sm mx-auto">
          TOKEN OFICIAL: <span class="font-bold text-white ml-2">{{ generatedToken }}</span>
        </div>

        <div class="pt-4">
          <RouterLink
            to="/"
            class="inline-block px-8 py-3 rounded-lg bg-white/10 hover:bg-white/20 font-orbitron text-xs uppercase text-white transition"
          >
            VOLVER A LA ARENA PRINCIPAL
          </RouterLink>
        </div>
      </div>

    </div>

  </div>
</template>
