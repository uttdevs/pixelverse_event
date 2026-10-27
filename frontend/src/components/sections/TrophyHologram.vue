<script setup>
/**
 * ============================================================================
 * COMPONENTE: TrophyHologram.vue (Holograma 3D del Trofeo de Torneo)
 * ============================================================================
 * Gráfico SVG del trofeo gamer con animaciones ciberpunk y tarjeta de telemetría.
 * Usa iconos Lucide (Cpu, Check).
 */

import { Cpu, Check } from 'lucide-vue-next'
</script>

<template>
  <div class="relative flex justify-center">
    <!-- Anillos concéntricos giratorios decorativos -->
    <div class="relative w-full max-w-md aspect-square flex items-center justify-center">
      <div
        class="absolute inset-4 rounded-full border-2 border-dashed border-cyber-blue/30 animate-[spin_25s_linear_infinite]"
        aria-hidden="true"
      ></div>
      <div
        class="absolute inset-10 rounded-full border border-cyber-purple/40 animate-[spin_18s_linear_infinite_reverse]"
        aria-hidden="true"
      ></div>

      <!-- Tarjeta Flotante Central con Efecto Glassmorphism -->
      <div class="relative z-10 w-full gradient-border rounded-2xl p-6 bg-cyber-card/90 border border-cyber-blue/40 box-glow-combo animate-float backdrop-blur-xl">
        <div class="flex justify-between items-center mb-4 pb-3 border-b border-white/10">
          <span class="font-rajdhani text-xs uppercase tracking-widest text-cyber-blue flex items-center gap-1.5">
            <Cpu class="w-4 h-4 text-cyber-blue" />
            SE EL GANADOR DE ESTA BATALLA
          </span>
          <span class="px-2 py-0.5 rounded bg-cyber-purple/20 border border-cyber-purple text-[10px] font-orbitron text-purple-300">
          PEXELVERS
          </span>
        </div>

        <!-- Gráfico Vectorial SVG del Trofeo Ciberpunk -->
        <div class="relative py-6 flex justify-center items-center bg-gradient-to-b from-cyber-blue/10 via-cyber-purple/10 to-transparent rounded-xl overflow-hidden">
          <svg
            class="w-40 h-40 drop-shadow-[0_0_25px_rgba(0,240,255,0.6)]"
            viewBox="0 0 24 24"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
            aria-label="Trofeo Holográfico del Torneo"
          >
            <defs>
              <linearGradient id="cyberGrad" x1="0%" y1="0%" x2="100%" y2="100%">
                <stop offset="0%" stop-color="#00f0ff" />
                <stop offset="100%" stop-color="#9d00ff" />
              </linearGradient>
            </defs>
            <path
              d="M12 17C15.866 17 19 13.866 19 10V3H5V10C5 13.866 8.13401 17 12 17Z"
              stroke="url(#cyberGrad)"
              stroke-width="1.5"
              fill="rgba(0, 240, 255, 0.08)"
            />
            <path
              d="M19 5H22V9C22 11.2091 20.2091 13 18 13"
              stroke="#00f0ff"
              stroke-width="1.5"
              stroke-linecap="round"
            />
            <path
              d="M5 5H2V9C2 11.2091 3.79086 13 6 13"
              stroke="#9d00ff"
              stroke-width="1.5"
              stroke-linecap="round"
            />
            <path d="M12 17V21" stroke="url(#cyberGrad)" stroke-width="2" />
            <path d="M7 21H17" stroke="#00f0ff" stroke-width="2" stroke-linecap="round" />
            <polygon
              points="12,6 13.5,9 17,9.5 14.5,12 15,15.5 12,13.8 9,15.5 9.5,12 7,9.5 10.5,9"
              fill="url(#cyberGrad)"
            />
          </svg>

          <!-- Badges de estadísticas flotantes -->
          <div class="absolute top-3 left-3 bg-cyber-dark/90 border border-cyber-blue/40 rounded px-2.5 py-1 text-[11px] font-rajdhani">
            <span class="text-gray-400"></span>
            <strong class="text-cyber-blue">CENTRO DE COMPETICION</strong>
          </div>
          <div class="absolute bottom-3 right-3 bg-cyber-dark/90 border border-cyber-purple/40 rounded px-2.5 py-1 text-[11px] font-rajdhani">
            <span class="text-gray-400">TU HABILIDAD, TU LEGADO</span>
            <strong class="text-cyber-purple"></strong>
          </div>
        </div>

        

      </div>
    </div>
  </div>
</template>
