<script setup>
/**
 * ============================================================================
 * COMPONENTE: SocialIcons.vue
 * ============================================================================
 * Provee iconos SVG vectoriales nativos para marcas oficiales de Esports:
 * - Twitch
 * - Discord
 * - X (anteriormente Twitter)
 * - YouTube
 * 
 * ¿Por qué este componente en lugar de FontAwesome o Lucide?
 * Lucide deliberadamente no incluye logotipos de marcas comerciales.
 * En lugar de descargar la hoja masiva de FontAwesome solo para 4 logos,
 * encapsulamos aquí los 4 SVGs limpios, accesibles y estilizables con Tailwind.
 */

import { useAudioSfx } from '../../composables/useAudioSfx'

const { playHoverSound } = useAudioSfx()
</script>

<template>
  <div class="flex items-center space-x-4">
    
   


    

   
  </div>
</template>
