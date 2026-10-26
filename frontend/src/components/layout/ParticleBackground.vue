<script setup>
/**
 * ============================================================================
 * COMPONENTE: ParticleBackground.vue
 * ============================================================================
 * ¿QUÉ APRENDERÁ EL EQUIPO CON ESTE COMPONENTE?
 * 1. Uso de `ref` para referenciar un elemento del DOM (<canvas ref="canvasRef">)
 *    en lugar de usar `document.getElementById('particle-canvas')`.
 * 2. Ciclos de vida en Vue 3 (`onMounted` y `onUnmounted`):
 *    - `onMounted`: Se ejecuta cuando el HTML ya existe en la pantalla. Aquí
 *      iniciamos animaciones o APIs del navegador (como Canvas 2D).
 *    - `onUnmounted`: Se ejecuta si el componente se destruye. Es VITAL limpiar
 *      los listeners de 'resize' y 'cancelAnimationFrame' para evitar
 *      fugas de memoria (Memory Leaks).
 */

import { ref, onMounted, onUnmounted } from 'vue'

// Referencia reactiva directa al elemento <canvas> del template
const canvasRef = ref(null)

// Variables internas de la animación (no necesitan ser reactivas porque se manejan en el bucle 60fps de canvas)
let ctx = null
let animationFrameId = null
let particles = []

// Definición de cada partícula de neón
class Particle {
  constructor(canvasWidth, canvasHeight) {
    this.canvasWidth = canvasWidth
    this.canvasHeight = canvasHeight
    this.reset()
  }

  reset() {
    this.x = Math.random() * this.canvasWidth
    this.y = Math.random() * this.canvasHeight
    this.size = Math.random() * 2.2 + 0.5
    this.speedX = (Math.random() - 0.5) * 0.4
    this.speedY = (Math.random() - 0.5) * 0.4
    this.color = Math.random() > 0.5 ? '#00f0ff' : '#9d00ff'
    this.alpha = Math.random() * 0.65 + 0.15
  }

  update() {
    this.x += this.speedX
    this.y += this.speedY
    if (this.x < 0 || this.x > this.canvasWidth || this.y < 0 || this.y > this.canvasHeight) {
      this.reset()
    }
  }

  draw(context) {
    context.save()
    context.globalAlpha = this.alpha
    context.fillStyle = this.color
    context.shadowBlur = 8
    context.shadowColor = this.color
    context.beginPath()
    context.arc(this.x, this.y, this.size, 0, Math.PI * 2)
    context.fill()
    context.restore()
  }
}

const resizeCanvas = () => {
  if (!canvasRef.value) return
  canvasRef.value.width = window.innerWidth
  canvasRef.value.height = window.innerHeight

  // Actualizamos el tamaño de las partículas existentes
  particles.forEach(p => {
    p.canvasWidth = window.innerWidth
    p.canvasHeight = window.innerHeight
  })
}

const initParticles = () => {
  resizeCanvas()
  particles = []
  const count = Math.min(65, Math.floor(window.innerWidth / 20))
  for (let i = 0; i < count; i++) {
    particles.push(new Particle(window.innerWidth, window.innerHeight))
  }
}

const animate = () => {
  if (!ctx || !canvasRef.value) return
  ctx.clearRect(0, 0, canvasRef.value.width, canvasRef.value.height)
  
  particles.forEach(p => {
    p.update()
    p.draw(ctx)
  })

  animationFrameId = requestAnimationFrame(animate)
}

// onMounted: El canvas ya fue insertado en el DOM real por Vue
onMounted(() => {
  if (canvasRef.value) {
    ctx = canvasRef.value.getContext('2d')
    initParticles()
    animate()
    window.addEventListener('resize', resizeCanvas)
  }
})

// onUnmounted: Limpiamos procesos para no saturar la memoria ni la CPU
onUnmounted(() => {
  if (animationFrameId) {
    cancelAnimationFrame(animationFrameId)
  }
  window.removeEventListener('resize', resizeCanvas)
})
</script>

<template>
  <!-- Canvas de fondo interactivo -->
  <canvas
    ref="canvasRef"
    class="fixed inset-0 pointer-events-none z-0 opacity-60"
    aria-hidden="true"
  ></canvas>

  <!-- Resplandores ambientales radiales difuminados -->
  <div
    class="fixed top-[-15%] left-[-10%] w-[500px] h-[500px] rounded-full bg-cyber-purple/20 blur-[140px] pointer-events-none z-0"
    aria-hidden="true"
  ></div>
  <div
    class="fixed bottom-[-15%] right-[-10%] w-[500px] h-[500px] rounded-full bg-cyber-blue/20 blur-[140px] pointer-events-none z-0"
    aria-hidden="true"
  ></div>
</template>
