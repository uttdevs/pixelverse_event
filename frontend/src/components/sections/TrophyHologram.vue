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
          PEXELVERSE
          </span>
        </div>

        <!-- Gráfico Vectorial SVG del Trofeo Ciberpunk -->
        <div class="relative py-6 flex justify-center items-center bg-gradient-to-b from-cyber-blue/10 via-cyber-purple/10 to-transparent rounded-xl overflow-hidden">

          <svg
            class="w-40 h-40 drop-shadow-[0_0_25px_rgba(0,240,255,0.6)]"
            viewBox="0 0 32 32"
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
            <g data-name="38. Medal" id="_38._Medal">
              <path
                d="M28,12A12,12,0,1,0,9,21.73V31a1,1,0,0,0,.42.81,1,1,0,0,0,.9.14L16,30.05,21.68,32A1.19,1.19,0,0,0,22,32a.94.94,0,0,0,.58-.19A1,1,0,0,0,23,31V21.73A12,12,0,0,0,28,12ZM21,29.61l-4.68-1.56a1,1,0,0,0-.64,0L11,29.61V22.89a11.82,11.82,0,0,0,10,0ZM16,22A10,10,0,1,1,26,12,10,10,0,0,1,16,22Z"
                fill="url(#cyberGrad)"
              />
              <path
                d="M15.5,8H19a1,1,0,0,0,0-2H17V5a1,1,0,0,0-2,0v1a3.49,3.49,0,0,0,.5,7h1a1.5,1.5,0,0,1,0,3H13a1,1,0,0,0,0,2h2v1a1,1,0,0,0,2,0V18A3.49,3.49,0,0,0,16.5,11h-1a1.5,1.5,0,0,1,0-3Z"
                fill="#00f0ff"
              />
            </g>
          </svg>

          <!-- Badges de estadísticas flotantes -->
          <div class="absolute top-3 left-3 bg-cyber-dark/90 border border-cyber-blue/40 rounded px-2.5 py-1 text-[11px] font-rajdhani">
            <span class="text-gray-400"></span>
            <strong class="text-cyber-blue">PREMIO MONETARIO</strong>
          </div>
          <div class="absolute bottom-3 right-3 bg-cyber-dark/90 border border-cyber-purple/40 rounded px-2.5 py-1 text-[11px] font-rajdhani">
            <span class="text-gray-400">PRIMER LUGAR EN CADA CATEGORIA</span>
            <strong class="text-cyber-purple"></strong>
          </div>
        </div>

        

      </div>
    </div>
  </div>
</template>
